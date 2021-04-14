package types

import (
	"fmt"

	"gopkg.in/yaml.v2"

	sdk "github.com/cosmos/cosmos-sdk/types"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"

	core "github.com/terra-project/core/types"
)

// Parameter keys
var (
	// Terra mint liquidity pool(usdr unit) made available per ${PoolRecoveryPeriod} (usdr unit)
	KeyMintBasePool = []byte("MintBasePool")
	// Terra burn liquidity pool(usdr unit) made available per ${PoolRecoveryPeriod} (usdr unit)
	KeyBurnBasePool = []byte("BurnBasePool")
	// The period required to recover BasePool
	KeyPoolRecoveryPeriod = []byte("PoolRecoveryPeriod")
	// Min spread
	KeyMinStabilitySpread = []byte("MinStabilitySpread")
)

// Default parameter values
var (
	DefaultMintBasePool       = sdk.NewDec(1000000 * core.MicroUnit) // 1000,000sdr = 1000,000,000,000usdr
	DefaultBurnBasePool       = sdk.NewDec(100000 * core.MicroUnit)  // 100,000sdr = 100,000,000,000usdr
	DefaultPoolRecoveryPeriod = core.BlocksPerDay                    // 14,400
	DefaultMinStabilitySpread = sdk.NewDecWithPrec(2, 2)             // 2%
)

var _ paramstypes.ParamSet = &Params{}

// DefaultParams creates default market module parameters
func DefaultParams() Params {
	return Params{
		MintBasePool:       DefaultMintBasePool,
		BurnBasePool:       DefaultBurnBasePool,
		PoolRecoveryPeriod: DefaultPoolRecoveryPeriod,
		MinStabilitySpread: DefaultMinStabilitySpread,
	}
}

// ParamKeyTable returns the parameter key table.
func ParamKeyTable() paramstypes.KeyTable {
	return paramstypes.NewKeyTable().RegisterParamSet(&Params{})
}

// String implements fmt.Stringer interface
func (p Params) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}

// ParamSetPairs implements the ParamSet interface and returns all the key/value pairs
// pairs of market module's parameters.
func (p *Params) ParamSetPairs() paramstypes.ParamSetPairs {
	return paramstypes.ParamSetPairs{
		paramstypes.NewParamSetPair(KeyMintBasePool, &p.MintBasePool, validateMintBasePool),
		paramstypes.NewParamSetPair(KeyBurnBasePool, &p.BurnBasePool, validateBurnBasePool),
		paramstypes.NewParamSetPair(KeyPoolRecoveryPeriod, &p.PoolRecoveryPeriod, validatePoolRecoveryPeriod),
		paramstypes.NewParamSetPair(KeyMinStabilitySpread, &p.MinStabilitySpread, validateMinStabilitySpread),
	}
}

// Validate a set of params
func (p Params) Validate() error {
	if p.MintBasePool.IsNegative() {
		return fmt.Errorf("mint base pool should be positive or zero, is %s", p.MintBasePool)
	}
	if p.BurnBasePool.IsNegative() {
		return fmt.Errorf("burn base pool should be positive or zero, is %s", p.BurnBasePool)
	}
	if p.PoolRecoveryPeriod == 0 {
		return fmt.Errorf("pool recovery period should be positive, is %d", p.PoolRecoveryPeriod)
	}
	if p.MinStabilitySpread.IsNegative() || p.MinStabilitySpread.GT(sdk.OneDec()) {
		return fmt.Errorf("market minimum stability spead should be a value between [0,1], is %s", p.MinStabilitySpread)
	}

	return nil
}

func validateMintBasePool(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNegative() {
		return fmt.Errorf("mint base pool must be positive or zero: %s", v)
	}

	return nil
}

func validateBurnBasePool(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNegative() {
		return fmt.Errorf("burn base pool must be positive or zero: %s", v)
	}

	return nil
}

func validatePoolRecoveryPeriod(i interface{}) error {
	v, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v <= 0 {
		return fmt.Errorf("pool recovery period must be positive: %d", v)
	}

	return nil
}

func validateMinStabilitySpread(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNegative() {
		return fmt.Errorf("min spread must be positive or zero: %s", v)
	}

	if v.GT(sdk.OneDec()) {
		return fmt.Errorf("min spread is too large: %s", v)
	}

	return nil
}
