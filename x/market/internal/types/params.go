package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params/subspace"

	core "github.com/terra-project/core/types"
)

// DefaultParamspace
const DefaultParamspace = ModuleName

// Parameter keys
var (
	ParamStoreKeyPoolUpdateInterval = []byte("poolupdateinterval")
	// The period required to recover Terra&Luna Pool to BasePool
	ParamStoreKeyPoolRecoveryPeriod = []byte("poolrecoveryperiod")
	// Terra's market cap made available per ${poolrecoveryperiod}
	ParamStoreKeyTerraLiquidityRatio = []byte("dailyterraliquidityratio")
	// Min spread
	ParamStoreKeyMinSpread = []byte("minspread")
	// Tobin tax
	ParmamStoreKeyTobinTax = []byte("tobintax")
)

// Default parameter values
var (
	DefaultPoolUpdateInterval  = core.BlocksPerEpoch       // 100,400
	DefaultPoolRecoveryPeriod  = core.BlocksPerDay         // 14,400
	DefaultTerraLiquidityRatio = sdk.NewDecWithPrec(1, 2)  // 1%
	DefaultMinSpread           = sdk.NewDecWithPrec(2, 2)  // 2%
	DefaultTobinTax            = sdk.NewDecWithPrec(30, 4) // 0.3%
)

var _ subspace.ParamSet = &Params{}

// Params market parameters
type Params struct {
	PoolUpdateInterval  int64   `json:"pool_update_interval" yaml:"pool_update_interval"`
	PoolRecoveryPeriod  int64   `json:"pool_recovery_period" yaml:"pool_recovery_period"`
	TerraLiquidityRatio sdk.Dec `json:"terra_liquidity_ratio" yaml:"terra_liquidity_ratio"`
	MinSpread           sdk.Dec `json:"min_spread" yaml:"min_spread"`
	TobinTax            sdk.Dec `json:"tobin_tax" yaml:"tobin_tax"`
}

// DefaultParams creates default market module parameters
func DefaultParams() Params {
	return Params{
		PoolUpdateInterval:  DefaultPoolUpdateInterval,
		PoolRecoveryPeriod:  DefaultPoolRecoveryPeriod,
		TerraLiquidityRatio: DefaultTerraLiquidityRatio,
		MinSpread:           DefaultMinSpread,
		TobinTax:            DefaultTobinTax,
	}
}

// Validate a set of params
func (params Params) Validate() error {
	if params.PoolUpdateInterval <= 0 {
		return fmt.Errorf("pool update interval should be positive, is %d", params.PoolUpdateInterval)
	}
	if params.PoolRecoveryPeriod <= 0 {
		return fmt.Errorf("pool recovery period should be positive, is %d", params.PoolRecoveryPeriod)
	}
	if params.TerraLiquidityRatio.LT(sdk.ZeroDec()) || params.TerraLiquidityRatio.GT(sdk.OneDec()) {
		return fmt.Errorf("daily terra liquidity ratio should be a value between [0,1], is %s", params.TerraLiquidityRatio.String())
	}
	if params.MinSpread.IsNegative() || params.MinSpread.GT(sdk.OneDec()) {
		return fmt.Errorf("market minimum spead should be a value between [0,1], is %s", params.MinSpread.String())
	}
	if params.TobinTax.IsNegative() || params.TobinTax.GT(sdk.OneDec()) {
		return fmt.Errorf("tobin tax should be a value between [0,1], is %s", params.TobinTax.String())
	}

	return nil
}

// ParamSetPairs implements the ParamSet interface and returns all the key/value pairs
// pairs of market module's parameters.
// nolint
func (params *Params) ParamSetPairs() subspace.ParamSetPairs {
	return subspace.ParamSetPairs{
		{Key: ParamStoreKeyPoolUpdateInterval, Value: &params.PoolUpdateInterval},
		{Key: ParamStoreKeyPoolRecoveryPeriod, Value: &params.PoolRecoveryPeriod},
		{Key: ParamStoreKeyTerraLiquidityRatio, Value: &params.TerraLiquidityRatio},
		{Key: ParamStoreKeyMinSpread, Value: &params.MinSpread},
		{Key: ParmamStoreKeyTobinTax, Value: &params.TobinTax},
	}
}

// String implements fmt.Stringer interface
func (params Params) String() string {
	return fmt.Sprintf(`Treasury Params:
	PoolUpdateInterval:					%d
	PoolRecoveryPeriod:         %d
  TerraLiquidityRatio:		    %s
	MinSpread:            			%s
	TobinTax:                   %s
	`, params.PoolUpdateInterval, params.PoolRecoveryPeriod, params.TerraLiquidityRatio, params.MinSpread, params.TobinTax)
}
