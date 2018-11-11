package oracle

import (
	"encoding/json"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	denomWhiteList = []string{
		"terra", "luna",
	}
)

// PriceVote - struct to store a validator's vote on the price
type PriceVote struct {
	Denom  string
	Price  sdk.Dec
	Feeder sdk.AccAddress
	Power  sdk.Dec
}

// NewPriceVote creates a PriceVote instance
func NewPriceVote(feedMsg PriceFeedMsg, power sdk.Dec) PriceVote {
	return PriceVote{
		Denom:  feedMsg.Denom,
		Price:  feedMsg.Price,
		Feeder: feedMsg.Feeder,
		Power:  power,
	}
}

// PriceVotes are a collection of Price Votes
type PriceVotes []PriceVote

func (pv PriceVotes) Len() int {
	return len(pv)
}
func (pv PriceVotes) Swap(i, j int) {
	pv[i], pv[j] = pv[j], pv[i]
}
func (pv PriceVotes) Less(i, j int) bool {
	return pv[i].Price.LT(pv[j].Price)
}

//-------------------------------------------------
//-------------------------------------------------

// PriceFeedMsg - struct for voting on payloads
type PriceFeedMsg struct {
	Denom  string
	Price  sdk.Dec
	Feeder sdk.AccAddress
}

// NewPriceFeedMsg creates a PriceFeedMsg instance
func NewPriceFeedMsg(denom string, price sdk.Dec, feederAddress sdk.AccAddress) PriceFeedMsg {
	return PriceFeedMsg{
		Denom:  denom,
		Price:  price,
		Feeder: feederAddress,
	}
}

// Route Implements Msg
func (msg PriceFeedMsg) Route() string { return "oracle" }

// Type implements sdk.Msg
func (msg PriceFeedMsg) Type() string { return "pricefeed" }

// GetSignBytes implements sdk.Msg
func (msg PriceFeedMsg) GetSignBytes() []byte {
	bz, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(bz)
}

// GetSigners implements sdk.Msg
func (msg PriceFeedMsg) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Feeder}
}

// ValidateBasic Implements sdk.Msg
func (msg PriceFeedMsg) ValidateBasic() sdk.Error {
	if len(msg.Feeder) == 0 {
		return sdk.ErrInvalidAddress("Invalid address: " + msg.Feeder.String())
	}

	whiteListed := false
	for i := 0; i < len(denomWhiteList); i++ {
		if denomWhiteList[i] == msg.Denom {
			whiteListed = true
		}
	}

	if !whiteListed {
		sdk.ErrInvalidCoins("Invalid denom for oracle price vote " + msg.Denom)
	}

	return nil
}

// String Implements sdk.Msg
func (msg PriceFeedMsg) String() string {
	return fmt.Sprintf("PriceFeedMsg{%v, %v, %v}", msg.Feeder, msg.Denom, msg.Price)
}
