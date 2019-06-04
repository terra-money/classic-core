package market

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Params market parameters
type Params struct {
	DailyLunaDeltaCap sdk.Dec `json:"daily_luna_delta_limit"` // daily % inflation or deflation cap on Luna
	MinSwapSpread     sdk.Dec `json:"min_swap_spread"`        // minimum spread for swaps involving Luna
	MaxSwapSpread     sdk.Dec `json:"max_swap_spread"`        // maximum spread for swaps involving Luna
}

// NewParams creates a new param instance
func NewParams(dailyLunaDeltaCap, minSwapSpread, maxSwapSpread sdk.Dec) Params {
	return Params{
		DailyLunaDeltaCap: dailyLunaDeltaCap,
		MinSwapSpread:     minSwapSpread,
		MaxSwapSpread:     maxSwapSpread,
	}
}

// DefaultParams creates default market module parameters
func DefaultParams() Params {
	return NewParams(
		sdk.NewDecWithPrec(5, 3),  // 0.5%
		sdk.NewDecWithPrec(2, 2),  // 2%
		sdk.NewDecWithPrec(10, 2), // 10%
	)
}

func validateParams(params Params) error {
	if params.DailyLunaDeltaCap.IsNegative() {
		return fmt.Errorf("market daily luna issuance change should be non-negative, is %s", params.DailyLunaDeltaCap.String())
	}
	if params.MinSwapSpread.IsNegative() {
		return fmt.Errorf("market minimum swap spead should be non-negative, is %s", params.MinSwapSpread.String())
	}
	if params.MaxSwapSpread.LT(params.MinSwapSpread) {
		return fmt.Errorf("market maximum swap spead should be larger or equal to the minimum, is %s", params.MaxSwapSpread.String())
	}

	return nil
}

func (params Params) String() string {
	return fmt.Sprintf(`market Params:
	DailyLunaDeltaCap: %v,
	MinSwapSpread:  %v,
	MaxSwapSpread:  %v
  `, params.DailyLunaDeltaCap, params.MinSwapSpread, params.MaxSwapSpread)
}
