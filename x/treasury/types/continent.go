package treasury

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Continent struct {
	Address sdk.AccAddress
	Name    string
	Website string
	Weight  sdk.Dec
}
