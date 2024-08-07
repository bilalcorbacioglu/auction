# Auction
**auction** is a blockchain built using Cosmos SDK and Tendermint and created with [Ignite CLI](https://ignite.com/cli).

## Get started

```
ignite chain serve
```

`serve` command installs dependencies, builds, initializes, and starts your blockchain in development.

## Overview

The Auction module allows users to create auctions and place bids. When bids are placed, the highest bidder's funds are stored in a storage account, and the previous highest bid is refunded. The highest bid in each auction is logged every 100 blocks.

## Example Scenario

Let's consider a scenario with three users: Bob, Alice, and Joe.

### Creating an Auction

1. Bob creates an auction for an item.

```sh
auctiond create-auction "Vintage Car" "10token" --from bob --chain-id auction --fees 10token -y
```

### Placing Bids

2. Alice places a bid on Bob's auction.

```sh
auctiond place-bid "auction-0" "15token" --from alice --chain-id auction --fees 10token -y
```

3. Joe places a higher bid on the same auction.

```sh
auctiond place-bid "auction-0" "20token" --from joe --chain-id auction --fees 10token -y
```

### Fund Handling

- When Alice places her bid, her `15token` is sent to the storage account.
- When Joe places a higher bid, Alice's `15token` is refunded to her, and Joe's `20token` is sent to the storage account.

### Checking Logs

The highest bid in each auction is logged every 100 blocks. This can be checked in the logs.

### Checking Transaction Events

To verify the details of a specific transaction, including events such as bids placed or auctions created, you can use the following command:

```sh
auctiond query tx TX_ID
```

Replace TX_ID with the transaction hash you want to inspect. This command will provide detailed information about the transaction, including the events that were triggered, the involved addresses, and the status of the transaction.

## Logging the Maximum Bid

The maximum bid in each auction is logged every 100 blocks in `EndBlocker`.

```go
func (k Keeper) EndBlocker(ctx sdk.Context) {
    if ctx.BlockHeight()%100 == 0 {
        k.Logger().Info("Checking maximum bids for auctions")
        // Implementation as described
        // ...
        // ...
    }
}
```

## Viewing Logs

To view the logs, start the node and check the logs for the auction module:

```sh
auctiond start
```

Look for log entries related to auction and bid events.

## Further Steps

- **In-memory Storage**: The current implementation uses in-memory storage, which means that all auction information is lost when Ignite is restarted. To avoid this, a database connection can be used to store auction data persistently.
- **Module Account**: Currently, a base account is used as the storage account, which may not be ideal. Instead, a module account can be used, allowing the bank module to interact more securely and appropriately.
