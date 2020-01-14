package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// PolicyConstraints wraps constraints around updating a key Treasury variable
type PolicyConstraints struct {
	RateMin       sdk.Dec  `json:"rate_min"`
	RateMax       sdk.Dec  `json:"rate_max"`
	Cap           sdk.Coin `json:"cap"`
	ChangeRateMax sdk.Dec  `json:"change_max"`
}

// String implements fmt.Stringer interface
func (pc PolicyConstraints) String() string {
	return fmt.Sprintf(`PolicyConstraints :
 RateMin:       %s
 RateMax:       %s
 Cap:           %s
 ChangeRateMax: %s
	`, pc.RateMin, pc.RateMax, pc.Cap, pc.ChangeRateMax)
}

// Clamp constrains a policy variable update within the policy constraints
func (pc PolicyConstraints) Clamp(prevRate sdk.Dec, newRate sdk.Dec) (clampedRate sdk.Dec) {
	if newRate.LT(pc.RateMin) {
		newRate = pc.RateMin
	} else if newRate.GT(pc.RateMax) {
		newRate = pc.RateMax
	}

	delta := newRate.Sub(prevRate)
	if newRate.GT(prevRate) {
		if delta.GT(pc.ChangeRateMax) {
			newRate = prevRate.Add(pc.ChangeRateMax)
		}
	} else {
		if delta.Abs().GT(pc.ChangeRateMax) {
			newRate = prevRate.Sub(pc.ChangeRateMax)
		}
	}
	return newRate
}
