package market

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Params market parameters
type Params struct {
	DailySwapLimit sdk.Dec `json:"daily_swap_limit"` // daily % inflation cap on a currency from swaps
}

// NewParams creates a new param instance
func NewParams(dailySwapLimit sdk.Dec) Params {
	return Params{
		DailySwapLimit: dailySwapLimit,
	}
}

// DefaultParams creates default market module parameters
func DefaultParams() Params {
	return NewParams(
		sdk.NewDecWithPrec(1, 2), // 1%
	)
}

func validateParams(params Params) error {
	if params.DailySwapLimit.IsNegative() {
		return fmt.Errorf("market daily swap limit should be non-negative, is %s", params.DailySwapLimit.String())
	}
	return nil
}

func (params Params) String() string {
	return fmt.Sprintf(`market Params:
	DailySwapLimit: %v
  `, params.DailySwapLimit)
}
