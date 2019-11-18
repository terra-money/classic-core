package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params/subspace"

	core "github.com/terra-project/core/types"
)

// DefaultParamspace nolint
const DefaultParamspace = ModuleName

// Parameter keys
var (
	//Terra liquidity pool(usdr unit) made available per ${poolrecoveryperiod} (usdr unit)
	ParamStoreKeyBasePool = []byte("basepool")
	// The period required to recover BasePool
	ParamStoreKeyPoolRecoveryPeriod = []byte("poolrecoveryperiod")
	// Min spread
	ParamStoreKeyMinSpread = []byte("minspread")
	// Tobin tax
	ParmamStoreKeyTobinTax = []byte("tobintax")
)

// Default parameter values
var (
	DefaultBasePool           = sdk.NewDec(250000 * core.MicroUnit) // 250,000sdr = 250,000,000,000usdr
	DefaultPoolRecoveryPeriod = core.BlocksPerDay                   // 14,400
	DefaultMinSpread          = sdk.NewDecWithPrec(2, 2)            // 2%
	DefaultTobinTax           = sdk.NewDecWithPrec(25, 4)           // 0.25%
)

var _ subspace.ParamSet = &Params{}

// Params market parameters
type Params struct {
	PoolRecoveryPeriod int64   `json:"pool_recovery_period" yaml:"pool_recovery_period"`
	BasePool           sdk.Dec `json:"base_pool" yaml:"base_pool"`
	MinSpread          sdk.Dec `json:"min_spread" yaml:"min_spread"`
	TobinTax           sdk.Dec `json:"tobin_tax" yaml:"tobin_tax"`
}

// DefaultParams creates default market module parameters
func DefaultParams() Params {
	return Params{
		BasePool:           DefaultBasePool,
		PoolRecoveryPeriod: DefaultPoolRecoveryPeriod,
		MinSpread:          DefaultMinSpread,
		TobinTax:           DefaultTobinTax,
	}
}

// Validate a set of params
func (params Params) Validate() error {
	if params.BasePool.IsNegative() {
		return fmt.Errorf("base pool should be positive or zero, is %d", params.BasePool)
	}
	if params.PoolRecoveryPeriod <= 0 {
		return fmt.Errorf("pool recovery period should be positive, is %d", params.PoolRecoveryPeriod)
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
		{Key: ParamStoreKeyBasePool, Value: &params.BasePool},
		{Key: ParamStoreKeyPoolRecoveryPeriod, Value: &params.PoolRecoveryPeriod},
		{Key: ParamStoreKeyMinSpread, Value: &params.MinSpread},
		{Key: ParmamStoreKeyTobinTax, Value: &params.TobinTax},
	}
}

// String implements fmt.Stringer interface
func (params Params) String() string {
	return fmt.Sprintf(`Treasury Params:
	BasePool:		            %d
	PoolRecoveryPeriod:         %d
	MinSpread:            			%s
	TobinTax:                   %s
	`, params.BasePool, params.PoolRecoveryPeriod, params.MinSpread, params.TobinTax)
}
