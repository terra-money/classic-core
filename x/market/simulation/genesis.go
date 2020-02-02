package simulation

// DONTCOVER

import (
	"math/rand"

	sdk "github.com/cosmos/cosmos-sdk/types"

	core "github.com/terra-project/core/types"
	"github.com/terra-project/core/x/market/internal/types"
)

// Simulation parameter constants
const (
	BasePool             = "base_pool"
	PoolRecoveryPeriod   = "pool_recovery_period"
	MinSpread            = "min_spread"
	TobinTax             = "tobin_tax"
	IlliquidTobinTaxList = "illiquid_tobin_tax_list"
)

// GenBasePool randomized BasePool
func GenBasePool(r *rand.Rand) sdk.Dec {
	return sdk.NewDec(100000000).Add(sdk.NewDec(int64(r.Intn(10000000000))))
}

// GenPoolRecoveryPeriod randomized PoolRecoveryPeriod
func GenPoolRecoveryPeriod(r *rand.Rand) int64 {
	return int64(100 + r.Intn(10000000000))
}

// GenMinSpread randomized MinSpread
func GenMinSpread(r *rand.Rand) sdk.Dec {
	return sdk.NewDecWithPrec(1, 2).Add(sdk.NewDecWithPrec(int64(r.Intn(100)), 3))
}

// GenTobinTax randomized TobinTax
func GenTobinTax(r *rand.Rand) sdk.Dec {
	return sdk.NewDecWithPrec(1, 3).Add(sdk.NewDecWithPrec(int64(r.Intn(100)), 3))
}

// GenIlliquidTobinTaxList randomized IlliquidTobinTaxList
func GenIlliquidTobinTaxList(r *rand.Rand) types.TobinTaxList {
	return types.TobinTaxList{
		types.TobinTax{
			Denom:   core.MicroMNTDenom,
			TaxRate: sdk.NewDecWithPrec(1, 2).Add(sdk.NewDecWithPrec(int64(r.Intn(100)), 3)),
		},
	}
}
