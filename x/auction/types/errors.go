package types

// DONTCOVER

import (
	sdkerrors "cosmossdk.io/errors"
)

// x/auction module sentinel errors
var (
	ErrInvalidSigner    = sdkerrors.Register(ModuleName, 1100, "expected gov account as only signer for proposal message")
	ErrSample           = sdkerrors.Register(ModuleName, 1101, "sample error")
	ErrInvalidAuctionId = sdkerrors.Register(ModuleName, 1102, "invalid auction ID")
	ErrInvalidBidAmount = sdkerrors.Register(ModuleName, 1103, "invalid bid amount")
)
