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

// ValidateRoutes validate each routes and following invariants
// - sum of weights must be smaller than 1
// - all weights should have positive value
func (s SeigniorageRoutes) ValidateRoutes() error {
	routes := s.Routes

	weightsSum := sdk.ZeroDec()
	addrMap := map[string]bool{}
	for _, pc := range routes {
		_, err := sdk.AccAddressFromBech32(pc.Address)
		if err != nil {
			return ErrInvalidAddress
		}

		// each weight must be bigger than zero
		if !pc.Weight.IsPositive() {
			return ErrInvalidWeight
		}

		// check duplicated address
		if _, exists := addrMap[pc.Address]; exists {
			return ErrDuplicateRoute
		}

		weightsSum = weightsSum.Add(pc.Weight)
		addrMap[pc.Address] = true
	}

	// the sum of weights must be smaller than one
	if weightsSum.GTE(sdk.OneDec()) {
		return ErrInvalidWeightsSum
	}

	return nil
}
