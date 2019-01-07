package budget

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Params oracle parameters
type Params struct {
	ActiveThreshold sdk.Dec `json:"active_threshold"` // threshold of vote that will transition a program open -> active budget queue
	LegacyThreshold sdk.Dec `json:"legacy_threshold"` // threshold of vote that will transition a program active -> legacy budget queue
	VotePeriod      int64   `json:"vote_period"`      // vote period
	MinDeposit      int64   `json:"min_deposit"`
}

// NewParams creates a new param instance
func NewParams(activeThreshold sdk.Dec, legacyThreshold sdk.Dec, votePeriod int64, minDeposit int64) Params {
	return Params{
		ActiveThreshold: activeThreshold,
		LegacyThreshold: legacyThreshold,
		VotePeriod:      votePeriod,
		MinDeposit:      minDeposit,
	}
}

// DefaultParams creates default oracle module parameters
func DefaultParams() Params {
	return NewParams(
		sdk.NewDecWithPrec(10, 2),
		sdk.NewDecWithPrec(0, 2),
		1209600,
		100,
	)
}

func validateParams(params Params) error {
	if params.ActiveThreshold.LT(sdk.ZeroDec()) {
		return fmt.Errorf("budget active threshold should be greater than 0, is %s", params.ActiveThreshold.String())
	}
	if params.LegacyThreshold.LT(sdk.ZeroDec()) {
		return fmt.Errorf("budget legacy threshold should be greater than 0, is %s", params.LegacyThreshold.String())
	}
	if params.VotePeriod.LT(sdk.ZeroInt()) {
		return fmt.Errorf("oracle parameter VotePeriod must be > 0, is %s", params.VotePeriod.String())
	}
	return nil
}
