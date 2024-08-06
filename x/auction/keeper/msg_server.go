package keeper

import (
	"auction/x/auction/types"
	"context"
	"fmt"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

// var _ types.MsgServer = msgServer{}

// CreateAuction handles the creation of an auction.
func (m msgServer) CreateAuction(goCtx context.Context, msg *types.MsgCreateAuction) (*types.MsgCreateAuctionResponse, error) {
	fmt.Println("CreateAuction fonksiyonuna girildi")
	ctx := sdk.UnwrapSDKContext(goCtx)
	auctionID, err := m.Keeper.AppendAuction(ctx, msg)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	return &types.MsgCreateAuctionResponse{AuctionId: auctionID.AuctionId}, nil
}

// PlaceBid handles the creation of a bid.
func (m msgServer) PlaceBid(goCtx context.Context, msg *types.MsgPlaceBid) (*types.MsgPlaceBidResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if !m.Keeper.IsAuctionExists(ctx, msg.AuctionId) {
		return nil, errorsmod.Wrapf(types.ErrInvalidAuctionId, "auction does not exist")
	}

	if !m.Keeper.IsValidBid(ctx, msg.AuctionId, *msg.BidAmount) {
		return nil, errorsmod.Wrapf(types.ErrInvalidBidAmount, "invalid bid amount")
	}

	m.Keeper.AppendBid(ctx, msg.AuctionId, msg.Bidder, *msg.BidAmount)

	return &types.MsgPlaceBidResponse{Success: true}, nil
}
