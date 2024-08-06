package keeper

import (
	"encoding/binary"
	"fmt"

	"cosmossdk.io/core/store"
	"cosmossdk.io/store/prefix"
	"github.com/cosmos/cosmos-sdk/runtime"

	// "cosmossdk.io/core/store"
	"cosmossdk.io/log"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"auction/x/auction/types"
)

type (
	Keeper struct {
		cdc          codec.BinaryCodec
		storeService store.KVStoreService
		logger       log.Logger

		// the address capable of executing a MsgUpdateParams message. Typically, this
		// should be the x/gov module account.
		authority string
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeService store.KVStoreService,
	logger log.Logger,
	authority string,

) Keeper {
	if _, err := sdk.AccAddressFromBech32(authority); err != nil {
		panic(fmt.Sprintf("invalid authority address: %s", authority))
	}

	return Keeper{
		cdc:          cdc,
		storeService: storeService,
		authority:    authority,
		logger:       logger,
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
	auctionID := fmt.Sprintf("%s%d", types.AuctionKey, auctionCount)

	auction := types.Auction{
		Creator:     msg.Creator,
		Item:        msg.Item,
		StartingBid: msg.StartingBid,
		Id:          auctionID,
	}

	auctionBytes := k.cdc.MustMarshal(&auction)
	store.Set([]byte(auctionID), auctionBytes)

	// Update the auction count
	k.SetAuctionCount(ctx, auctionCount+1)

	return &types.MsgCreateAuctionResponse{
		AuctionId: auctionID,
	}, nil
}

// PlaceBid places a bid on an auction.
func (k Keeper) PlaceBid(ctx sdk.Context, msg *types.MsgPlaceBid) (*types.MsgPlaceBidResponse, error) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.AuctionKey))

	// Check if the auction exists
	if !k.IsAuctionExists(ctx, msg.AuctionId) {
		return nil, fmt.Errorf("auction does not exist")
	}

	// Check if the bid is valid
	if !k.IsValidBid(ctx, msg.AuctionId, *msg.BidAmount) {
		return nil, fmt.Errorf("invalid bid amount")
	}

	// Get the auction from the store
	auctionBytes := store.Get([]byte(msg.AuctionId))
	var auction types.Auction
	k.cdc.MustUnmarshal(auctionBytes, &auction)

	// Append the bid to the auction
	bid := &types.Bid{
		Bidder:    msg.Bidder,
		BidAmount: msg.BidAmount,
	}
	auction.Bids = append(auction.Bids, bid)

	// Save the updated auction to the store
	auctionBytes = k.cdc.MustMarshal(&auction)
	store.Set([]byte(msg.AuctionId), auctionBytes)

	return &types.MsgPlaceBidResponse{Success: true}, nil
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
func (k Keeper) AppendBid(ctx sdk.Context, auctionID string, bidder string, bidAmount sdk.Coin) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.AuctionKey))

	// Get the auction from the store
	auctionBytes := store.Get([]byte(auctionID))
	var auction types.Auction
	k.cdc.MustUnmarshal(auctionBytes, &auction)

	// Append the bid to the auction
	bid := &types.Bid{
		Bidder:    bidder,
		BidAmount: &bidAmount,
	}
	auction.Bids = append(auction.Bids, bid)

	// Save the updated auction to the store
	auctionBytes = k.cdc.MustMarshal(&auction)
	store.Set([]byte(auctionID), auctionBytes)
}

// GetAuctionCount gets the number of auctions from the store.
func (k Keeper) GetAuctionCount(ctx sdk.Context) int {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.AuctionKey))

	bz := store.Get([]byte(types.AuctionKey + "count"))
	if bz == nil {
		return 0
	}
	return int(binary.BigEndian.Uint64(bz))
}

// SetAuctionCount sets the number of auctions in the store.
func (k Keeper) SetAuctionCount(ctx sdk.Context, count int) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.AuctionKey))

	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, uint64(count))
	store.Set([]byte(types.AuctionKey+"count"), bz)
}
