package types

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Auction defines a struct used by the nameservice module to allow for bidding
// on TNS name.
type Auction struct {
	Name         Name           `json:"name" yaml:"name"`
	Status       AuctionStatus  `json:"status" yaml:"status"`
	EndTime      time.Time      `json:"end_time" yaml:"end_time"`
	TopBidder    sdk.AccAddress `json:"top_bidder" yaml:"top_bidder"`
	TopBidAmount sdk.Coins      `json:"top_bid_amount" yaml:"top_bid_amount"`
}

// NewAuction returns Auction instance
func NewAuction(name Name, status AuctionStatus, endTime time.Time) Auction {
	return Auction{
		Name:    name,
		Status:  status,
		EndTime: endTime,
	}
}

// String implements fmt.Stringer interface
func (a Auction) String() string {
	return fmt.Sprintf(`Auction
Name:         %s
Status:       %s
EndTime:      %s
TopBidder:    %s
TopBidAmount: %s
`, a.Name, a.Status, a.EndTime, a.TopBidder, a.TopBidAmount)
}

// Auctions - array of Auction
type Auctions []Auction

// String implements fmt.Stringer interface
func (auctions Auctions) String() (out string) {
	for _, auction := range auctions {
		out += auction.String() + "\n"
	}
	return
}

// Bid - struct to store a user's bidding on the name auction
// The bidder have to pay bigger deposit than actual bidding price
type Bid struct {
	Hash    BidHash        `json:"hash" yaml:"hash"`
	Deposit sdk.Coin       `json:"deposit" yaml:"deposit"`
	Bidder  sdk.AccAddress `json:"bidder" yaml:"bidder"`
}

// NewBid returns Bid instance
func NewBid(hash BidHash, deposit sdk.Coin, bidder sdk.AccAddress) Bid {
	return Bid{
		Hash:    hash,
		Deposit: deposit,
		Bidder:  bidder,
	}
}

// String implements fmt.Stringer interface
func (b Bid) String() string {
	return fmt.Sprintf(`Auction
Hash:     %s
Deposit:  %s
Bidder:   %s 
`, b.Hash, b.Deposit, b.Bidder)
}

// Bids - array of Bid
type Bids []Bid

// String implements fmt.Stringer interface
func (bids Bids) String() (out string) {
	for _, auction := range bids {
		out += auction.String() + "\n"
	}
	return
}
