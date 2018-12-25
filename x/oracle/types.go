package oracle

import (
	"encoding/json"
	"fmt"

	"gonum.org/v1/gonum/stat"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	OracleDecPrec = 2
)

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

func getTotalVotePower(votes []PriceVote) sdk.Dec {
	votePower := sdk.ZeroDec()
	for _, vote := range votes {
		votePower.Add(vote.Power)
	}

	return votePower
}

func decToFloat64(a sdk.Dec) float64 {
	// roundup
	b := a.MulInt(sdk.NewInt(10 ^ OracleDecPrec))
	c := b.TruncateInt64()

	return float64(c) / (10 ^ OracleDecPrec)
}

func float64ToDec(a float64) sdk.Dec {
	b := int64(a * (10 ^ OracleDecPrec))
	return sdk.NewDecWithPrec(b, 2)
}

func tallyVotes(votes []PriceVote) (targetMode sdk.Dec, observedMode sdk.Dec, rewardees []PriceVote) {
	vTarget := make([]float64, len(votes))
	vPower := make([]float64, len(votes))
	vObserved := make([]float64, len(votes))

	for _, vote := range votes {
		vPower = append(vPower, decToFloat64(vote.Power))
		vTarget = append(vTarget, decToFloat64(vote.FeedMsg.TargetPrice))
		vObserved = append(vObserved, decToFloat64(vote.FeedMsg.ObservedPrice))
	}

	tmode, _ := stat.Mode(vTarget, vPower)
	omode, _ := stat.Mode(vObserved, vPower)

	tsd := stat.StdDev(vTarget, vPower)
	osd := stat.StdDev(vTarget, vPower)

	for i, vote := range votes {
		if vTarget[i] >= tmode-tsd && vTarget[i] <= tmode+tsd &&
			vObserved[i] >= omode-osd && vObserved[i] <= omode+osd {
			rewardees = append(rewardees, vote)
		}
	}

	targetMode = float64ToDec(tmode)
	observedMode = float64ToDec(omode)
	return
}

//-------------------------------------------------
//-------------------------------------------------

// PriceFeedMsg - struct for voting on payloads. Note that the Price
// is denominated in Luna. All validators must vote on Terra prices.
type PriceFeedMsg struct {
	Denom         string
	TargetPrice   sdk.Dec // in Luna
	ObservedPrice sdk.Dec // in Luna
	Feeder        sdk.AccAddress
}

// NewPriceFeedMsg creates a PriceFeedMsg instance
func NewPriceFeedMsg(denom string, targetPrice, observedPrice sdk.Dec, feederAddress sdk.AccAddress) PriceFeedMsg {
	return PriceFeedMsg{
		Denom:         denom,
		TargetPrice:   targetPrice,
		ObservedPrice: observedPrice,
		Feeder:        feederAddress,
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
	return fmt.Sprintf("PriceFeedMsg{feeder: %v, denom: %v, target: %v, observed: %v}",
		msg.Feeder, msg.Denom, msg.TargetPrice, msg.ObservedPrice)
}
