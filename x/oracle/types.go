package oracle

import (
	"encoding/json"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Whitelist []string

//-------------------------------------------------
//-------------------------------------------------

// PriceFeedMsg - struct for voting on payloads. Note that the Price
// is denominated in Luna. All validators must vote on Terra prices.
type PriceFeedMsg struct {
	Denom        string
	TargetPrice  sdk.Dec // in Luna
	CurrentPrice sdk.Dec // in Luna
	Feeder       sdk.AccAddress
}

// NewPriceFeedMsg creates a PriceFeedMsg instance
func NewPriceFeedMsg(denom string, targetPrice, currentPrice sdk.Dec, feederAddress sdk.AccAddress) PriceFeedMsg {
	return PriceFeedMsg{
		Denom:        denom,
		TargetPrice:  targetPrice,
		CurrentPrice: currentPrice,
		Feeder:       feederAddress,
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

	return nil
}

// String Implements sdk.Msg
func (msg PriceFeedMsg) String() string {
	return fmt.Sprintf("PriceFeedMsg{feeder: %v, denom: %v, target: %v, current: %v}",
		msg.Feeder, msg.Denom, msg.TargetPrice, msg.CurrentPrice)
}

//-------------------------------------------------
//-------------------------------------------------

// PriceVote - struct to store a validator's vote on the price
type PriceVote struct {
	FeedMsg PriceFeedMsg
	Power   sdk.Dec
}

// NewPriceVote creates a PriceVote instance
func NewPriceVote(feedMsg PriceFeedMsg, power sdk.Dec) PriceVote {
	return PriceVote{
		FeedMsg: feedMsg,
		Power:   power,
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
	return pv[i].FeedMsg.CurrentPrice.LT(pv[j].FeedMsg.CurrentPrice)
}
