package types

import (
	"gopkg.in/yaml.v2"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewSeigniorageRoute returns SeigniorageRoute object
func NewSeigniorageRoute(address sdk.AccAddress, weight sdk.Dec) SeigniorageRoute {
	return SeigniorageRoute{
		Address: address.String(),
		Weight:  weight,
	}
}

// String implement stringify
func (v SeigniorageRoute) String() string {
	out, _ := yaml.Marshal(v)
	return string(out)
}
