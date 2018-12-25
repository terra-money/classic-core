package treasury

import sdk "github.com/cosmos/cosmos-sdk/types"

type Claim struct {
	Account sdk.AccAddress
	Weight  sdk.Dec
}

type Continent struct {
	Address sdk.AccAddress
	Name    string
	Website string
	Weight  sdk.Dec
}
