package types

const (
	// ModuleName defines the module name
	ModuleName = "auction"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey is the message route for slashing
	RouterKey = ModuleName

	// QuerierRoute defines the module's query routing key
	QuerierRoute = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_auction"

	// AuctionKey defines the key to store auctions
	AuctionKey = "auction-"
)

var (
	ParamsKey = []byte("p_auction")
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}
