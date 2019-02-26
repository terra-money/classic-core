package oracle

import (
	"encoding/json"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// Precision of Oracle vote
	OracleDecPrec = 2

	// RouterKey is they name of the oracle module
	RouterKey = "oracle"
)

//-------------------------------------------------
//-------------------------------------------------

// PriceFeedMsg - struct for voting on payloads. Note that the Price
// is denominated in Luna. All validators must vote on Terra prices.
type PriceFeedMsg struct {
	Denom  string
	Price  sdk.Dec // in Luna
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
	if len(msg.Denom) == 0 {
		return ErrUnknownDenomination(DefaultCodespace, "")
	}

	if len(msg.Feeder) == 0 {
		return sdk.ErrInvalidAddress("Invalid address: " + msg.Feeder.String())
	}

	if msg.Price.LTE(sdk.ZeroDec()) {
		return ErrInvalidPrice(DefaultCodespace, msg.Price)
	}

	return nil
}

// String Implements sdk.Msg
func (msg PriceFeedMsg) String() string {
	return fmt.Sprintf("PriceFeedMsg{feeder: %v, denom: %v, price: %v}",
		msg.Feeder, msg.Denom, msg.Price)
}
