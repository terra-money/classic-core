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
	ParamStoreKeyMinStabilitySpread = []byte("minstabilityspread")
)

// Default parameter values
var (
	DefaultBasePool           = sdk.NewDec(250000 * core.MicroUnit) // 250,000sdr = 250,000,000,000usdr
	DefaultPoolRecoveryPeriod = core.BlocksPerDay                   // 14,400
	DefaultMinStabilitySpread = sdk.NewDecWithPrec(2, 2)            // 2%
)

var _ subspace.ParamSet = &Params{}

// Params market parameters
type Params struct {
	PoolRecoveryPeriod int64   `json:"pool_recovery_period" yaml:"pool_recovery_period"`
	BasePool           sdk.Dec `json:"base_pool" yaml:"base_pool"`
	MinStabilitySpread sdk.Dec `json:"min_spread" yaml:"min_spread"`
}

// DefaultParams creates default market module parameters
func DefaultParams() Params {
	return Params{
		BasePool:           DefaultBasePool,
		PoolRecoveryPeriod: DefaultPoolRecoveryPeriod,
		MinStabilitySpread: DefaultMinStabilitySpread,
	}
}

// Validate a set of params
func (params Params) Validate() error {
	if params.BasePool.IsNegative() {
		return fmt.Errorf("base pool should be positive or zero, is %s", params.BasePool)
	}
	if params.PoolRecoveryPeriod <= 0 {
		return fmt.Errorf("pool recovery period should be positive, is %d", params.PoolRecoveryPeriod)
	}
	if params.MinStabilitySpread.IsNegative() || params.MinStabilitySpread.GT(sdk.OneDec()) {
		return fmt.Errorf("market minimum stability spead should be a value between [0,1], is %s", params.MinStabilitySpread)
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
		{Key: ParamStoreKeyMinStabilitySpread, Value: &params.MinStabilitySpread},
	}
}

// String implements fmt.Stringer interface
func (params Params) String() string {
	return fmt.Sprintf(`Treasury Params:
	BasePool:                   %s
	PoolRecoveryPeriod:         %d
	MinStabilitySpread:         %s
	`, params.BasePool, params.PoolRecoveryPeriod, params.MinStabilitySpread)
}
