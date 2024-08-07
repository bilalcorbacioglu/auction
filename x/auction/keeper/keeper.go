package keeper

import (
	"encoding/binary"
	"fmt"

	"cosmossdk.io/core/store"
	"cosmossdk.io/store/prefix"
	"github.com/cosmos/cosmos-sdk/runtime"

	// "cosmossdk.io/core/store"
	"auction/x/auction/types"

	"cosmossdk.io/log"
	// storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type (
	Keeper struct {
		cdc           codec.BinaryCodec
		storeService  store.KVStoreService
		logger        log.Logger
		bankKeeper    types.BankKeeper
		accountKeeper types.AccountKeeper
		// the address capable of executing a MsgUpdateParams message. Typically, this
		// should be the x/gov module account.
		authority      string
		storageAddress sdk.AccAddress
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeService store.KVStoreService,
	logger log.Logger,
	authority string,
	bankKeeper types.BankKeeper,
	accountKeeper types.AccountKeeper,
	storageAddress sdk.AccAddress,
) Keeper {
	if _, err := sdk.AccAddressFromBech32(authority); err != nil {
		panic(fmt.Sprintf("invalid authority address: %s", authority))
	}

	return Keeper{
		cdc:            cdc,
		storeService:   storeService,
		authority:      authority,
		logger:         logger,
		bankKeeper:     bankKeeper,
		accountKeeper:  accountKeeper,
		storageAddress: storageAddress,
	}
}

// GetAuthority returns the module's authority.
func (k Keeper) GetAuthority() string {
	return k.authority
}

// Logger returns a module-specific logger.
func (k Keeper) Logger() log.Logger {
	return k.logger.With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// AppendAuction creates a new auction.
func (k Keeper) AppendAuction(ctx sdk.Context, msg *types.MsgCreateAuction) (*types.MsgCreateAuctionResponse, error) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.AuctionKey))

	auctionCount := k.GetAuctionCount(ctx)
	auctionID := fmt.Sprintf("auction-%d", auctionCount)

	auction := types.Auction{
		Creator:     msg.Creator,
		Item:        msg.Item,
		StartingBid: msg.StartingBid,
		Id:          auctionID,
		Bids:        []*types.Bid{},
	}

	auctionBytes := k.cdc.MustMarshal(&auction)
	store.Set([]byte(auctionID), auctionBytes)

	// Update the auction count
	k.SetAuctionCount(ctx, auctionCount+1)

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"create_auction",
			sdk.NewAttribute("auction_id", auctionID),
		),
	)

	return &types.MsgCreateAuctionResponse{
		AuctionId: auctionID,
	}, nil
}

// IsAuctionExists checks if the auction exists.
func (k Keeper) IsAuctionExists(ctx sdk.Context, auctionID string) bool {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.AuctionKey))

	return store.Has([]byte(auctionID))
}

// IsValidBid checks if the bid amount is valid for the auction.
func (k Keeper) IsValidBid(ctx sdk.Context, auctionID string, bidAmount sdk.Coin) bool {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.AuctionKey))

	// Get the auction from the store
	auctionBytes := store.Get([]byte(auctionID))
	var auction types.Auction
	k.cdc.MustUnmarshal(auctionBytes, &auction)

	// Check if the bid amount is greater than or equal to the starting bid and all previous bids
	if bidAmount.IsLT(*auction.StartingBid) {
		return false
	}
	for _, bid := range auction.Bids {
		if bidAmount.IsLT(*bid.BidAmount) {
			return false
		}
	}

	return true
}

// AppendBid appends a bid to the auction.
func (k Keeper) AppendBid(ctx sdk.Context, auctionID string, bidder string, bidAmount sdk.Coin) (*types.MsgPlaceBidResponse, error) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.AuctionKey))

	// Get the auction from the store
	auctionBytes := store.Get([]byte(auctionID))
	var auction types.Auction
	k.cdc.MustUnmarshal(auctionBytes, &auction)

	// Check if there's a previous highest bid
	var previousHighestBid *types.Bid
	if len(auction.Bids) > 0 {
		previousHighestBid = auction.Bids[len(auction.Bids)-1]
	}

	// Append the bid to the auction
	bid := &types.Bid{
		Bidder:    bidder,
		BidAmount: &bidAmount,
	}
	auction.Bids = append(auction.Bids, bid)

	// Save the updated auction to the store
	auctionBytes = k.cdc.MustMarshal(&auction)
	store.Set([]byte(auctionID), auctionBytes)

	// Send coins from bidder to storage account
	bidderAddress, _ := sdk.AccAddressFromBech32(bidder)
	err := k.bankKeeper.SendCoins(ctx, bidderAddress, k.storageAddress, sdk.NewCoins(bidAmount))
	if err != nil {
		return nil, err
	}

	// Refund the previous highest bidder if there was one
	if previousHighestBid != nil {
		previousHighestBidderAddress, _ := sdk.AccAddressFromBech32(previousHighestBid.Bidder)
		err = k.bankKeeper.SendCoins(ctx, k.storageAddress, previousHighestBidderAddress, sdk.NewCoins(*previousHighestBid.BidAmount))
		if err != nil {
			return nil, err
		}
	}

	// Broadcast an event
	bidsString := ""
	for _, bid := range auction.Bids {
		bidsString += fmt.Sprintf("%v, ", bid)
	}
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"place_bid",
			sdk.NewAttribute("bids", bidsString),
		),
	)

	return &types.MsgPlaceBidResponse{Success: true}, nil
}

// GetAuctionCount gets the number of auctions from the store.
func (k Keeper) GetAuctionCount(ctx sdk.Context) int {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.AuctionKey))

	bz := store.Get([]byte("count"))
	if bz == nil {
		return 0
	}
	count := int(binary.BigEndian.Uint64(bz))
	return count
}

// SetAuctionCount sets the number of auctions in the store.
func (k Keeper) SetAuctionCount(ctx sdk.Context, count int) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.AuctionKey))

	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, uint64(count))
	store.Set([]byte("count"), bz)
}
