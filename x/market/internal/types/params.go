package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params/subspace"
)

// DefaultParamspace
const DefaultParamspace = ModuleName

// Parameter keys
var (
	ParamStoreKeyDailyLunaDeltaCap = []byte("dailylunadeltalimit")
	ParamStoreKeyMaxSwapSpread     = []byte("maxswapspread")
	ParamStoreKeyMinSwapSpread     = []byte("minswapspread")
)

// Default parameter values
var (
	DefaultDailyLunaDeltaCap = sdk.NewDecWithPrec(5, 3) // 0.5%
	DefaultMaxSwapSpread     = sdk.NewDec(1)            // 100%
	DefaultMinSwapSpread     = sdk.NewDecWithPrec(2, 2) // 2%
)

var _ subspace.ParamSet = &Params{}

// Params market parameters
type Params struct {
	DailyLunaDeltaCap sdk.Dec `json:"daily_luna_delta_cap" yaml:"daily_luna_delta_cap"`
	MaxSwapSpread     sdk.Dec `json:"max_swap_spread" yaml:"max_swap_spread"`
	MinSwapSpread     sdk.Dec `json:"min_swap_spread" yaml:"min_swap_spread"`
}

// DefaultParams creates default market module parameters
func DefaultParams() Params {
	return Params{
		DailyLunaDeltaCap: DefaultDailyLunaDeltaCap,
		MaxSwapSpread:     DefaultMaxSwapSpread,
		MinSwapSpread:     DefaultMinSwapSpread,
	}
}

// Validate a set of params
func (params Params) Validate() error {
	if params.DailyLunaDeltaCap.IsNegative() {
		return fmt.Errorf("market daily luna issuance change should be non-negative, is %s", params.DailyLunaDeltaCap.String())
	}
	if params.MinSwapSpread.IsNegative() || params.MinSwapSpread.GT(sdk.OneDec()) {
		return fmt.Errorf("market minimum swap spead should be non-negative, is %s", params.MinSwapSpread.String())
	}
	if params.MaxSwapSpread.LT(params.MinSwapSpread) || params.MaxSwapSpread.GT(sdk.OneDec()) {
		return fmt.Errorf("market maximum swap spead should be larger or equal to the minimum, is %s", params.MaxSwapSpread.String())
	}

	return nil
}

// ParamSetPairs implements the ParamSet interface and returns all the key/value pairs
// pairs of oracle module's parameters.
// nolint
func (params *Params) ParamSetPairs() subspace.ParamSetPairs {
	return subspace.ParamSetPairs{
		{Key: ParamStoreKeyDailyLunaDeltaCap, Value: &params.DailyLunaDeltaCap},
		{Key: ParamStoreKeyMaxSwapSpread, Value: &params.MaxSwapSpread},
		{Key: ParamStoreKeyMinSwapSpread, Value: &params.MinSwapSpread},
	}
}

// implements fmt.Stringer
func (params Params) String() string {
	return fmt.Sprintf(`Market Params:
  DailyLunaDeltaCap:        %s
  MaxSwapSpread:            %s
	MinSwapSpread:            %s
	`, params.DailyLunaDeltaCap, params.MaxSwapSpread, params.MinSwapSpread)
}
