package budget

import (
	"fmt"
	"terra/types/assets"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Params budget parameters
type Params struct {
	ActiveThreshold sdk.Dec  `json:"active_threshold"` // threshold of vote that will transition a program open -> active budget queue
	LegacyThreshold sdk.Dec  `json:"legacy_threshold"` // threshold of vote that will transition a program active -> legacy budget queue
	VotePeriod      int64    `json:"vote_period"`      // vote period
	Deposit         sdk.Coin `json:"deposit"`          // Minimum deposit in TerraSDR
}

// NewParams creates a new param instance
func NewParams(activeThreshold sdk.Dec, legacyThreshold sdk.Dec, votePeriod int64, deposit sdk.Coin) Params {
	return Params{
		ActiveThreshold: activeThreshold,
		LegacyThreshold: legacyThreshold,
		VotePeriod:      votePeriod,
		Deposit:         deposit,
	}
}

// DefaultParams creates default budget module parameters
func DefaultParams() Params {
	return NewParams(
		sdk.NewDecWithPrec(1, 1), // 10%
		sdk.NewDecWithPrec(0, 2), // 0%
		util.BlocksPerMonth,
		sdk.NewInt64Coin(assets.MicroSDRDenom, sdk.NewInt(100).MulRaw(assets.MicroUnit).Int64()),
	)
}

func validateParams(params Params) error {
	if params.ActiveThreshold.LT(sdk.ZeroDec()) {
		return fmt.Errorf("budget active threshold should be greater than 0, is %s", params.ActiveThreshold.String())
	}
	if params.LegacyThreshold.LT(sdk.ZeroDec()) {
		return fmt.Errorf("budget legacy threshold should be greater than or equal to 0, is %s", params.LegacyThreshold.String())
	}
	if params.VotePeriod < 0 {
		return fmt.Errorf("budget parameter VotePeriod must be > 0, is %d", params.VotePeriod)
	}

	if params.Deposit.Amount.LTE(sdk.ZeroInt()) {
		return fmt.Errorf("budget parameter Deposit must be > 0, is %v", params.Deposit.String())
	}
	return nil
}

func (params Params) String() string {
	return fmt.Sprintf(`Budget Params:
	ActiveThreshold: %s
	LegacyThreshold: %s
	VotePeriod: %d
	Deposit: %s
  `, params.ActiveThreshold, params.LegacyThreshold, params.VotePeriod, params.Deposit)
}
