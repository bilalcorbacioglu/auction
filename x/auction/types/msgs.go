package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Ensure MsgCreateAuction implements sdk.Msg interface
var _ sdk.Msg = &MsgCreateAuction{}

// NewMsgCreateAuction creates a new MsgCreateAuction instance
func NewMsgCreateAuction(creator string, item string, startingBid sdk.Coin) *MsgCreateAuction {
	return &MsgCreateAuction{
		Creator:     creator,
		Item:        item,
		StartingBid: &startingBid,
	}
}

// ValidateBasic performs basic validation
func (msg *MsgCreateAuction) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Creator); err != nil {
		return fmt.Errorf("invalid creator address: %w", err)
	}
	if msg.Item == "" {
		return fmt.Errorf("item cannot be empty")
	}
	if !msg.StartingBid.IsValid() {
		return fmt.Errorf("invalid starting bid")
	}
	return nil
}

// Ensure MsgPlaceBid implements sdk.Msg interface
var _ sdk.Msg = &MsgPlaceBid{}

// NewMsgPlaceBid creates a new MsgPlaceBid instance
func NewMsgPlaceBid(bidder string, auctionID string, bidAmount sdk.Coin) *MsgPlaceBid {
	return &MsgPlaceBid{
		Bidder:    bidder,
		AuctionId: auctionID,
		BidAmount: &bidAmount,
	}
}

// ValidateBasic performs basic validation
func (msg *MsgPlaceBid) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Bidder); err != nil {
		return fmt.Errorf("invalid bidder address: %w", err)
	}
	if msg.AuctionId == "" {
		return fmt.Errorf("auction ID cannot be empty")
	}
	if !msg.BidAmount.IsValid() {
		return fmt.Errorf("invalid bid amount")
	}
	return nil
}
