package keeper

import (
	"fmt"
	"time"

	"auction/x/auction/types"

	"cosmossdk.io/store/prefix"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// // BeginBlocker will persist the current header and validator set as a historical entry
// // and prune the oldest entry based on the HistoricalEntries parameter
// func (k *Keeper) BeginBlocker(ctx sdk.Context) {
// 	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyBeginBlocker)
// 	// Add your implementation here
// }

func (k *Keeper) EndBlocker(ctx sdk.Context) {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyEndBlocker)

	if ctx.BlockHeight()%100 == 0 {
		k.Logger().Info("Checking maximum bids for auctions")

		storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
		store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.AuctionKey))
		iterator := store.Iterator(nil, nil)
		defer iterator.Close()

		for ; iterator.Valid(); iterator.Next() {
			if string(iterator.Key()) == "count" {
				continue
			}
			var auction types.Auction
			err := k.cdc.Unmarshal(iterator.Value(), &auction)
			if err != nil {
				k.Logger().Error(fmt.Sprintf("Failed to unmarshal auction: %v", err))
				continue
			}

			if len(auction.Bids) > 0 {
				highestBid := auction.Bids[0]
				for _, bid := range auction.Bids {
					if bid.BidAmount.IsGTE(*highestBid.BidAmount) {
						highestBid = bid
					}
				}

				// Output the maximum bid
				k.Logger().Info(fmt.Sprintf("Auction %s: Highest bid is %s from %s", auction.Id, highestBid.BidAmount.String(), highestBid.Bidder))
			}
		}
	}
}
