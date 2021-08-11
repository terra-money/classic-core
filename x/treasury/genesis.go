package treasury

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	core "github.com/terra-money/core/types"
	"github.com/terra-money/core/x/treasury/keeper"
	"github.com/terra-money/core/x/treasury/types"
)

// InitGenesis initializes default parameters
// and the keeper's address to pubkey map
func InitGenesis(ctx sdk.Context, keeper keeper.Keeper, data *types.GenesisState) {
	keeper.SetParams(ctx, data.Params)

	keeper.SetTaxRate(ctx, data.TaxRate)
	keeper.SetRewardWeight(ctx, data.RewardWeight)
	keeper.SetEpochTaxProceeds(ctx, data.TaxProceeds)

	// If EpochInitialIssuance is empty, we use current supply as epoch initial issuance
	if data.EpochInitialIssuance.IsZero() {
		keeper.RecordEpochInitialIssuance(ctx)
	} else {
		keeper.SetEpochInitialIssuance(ctx, data.EpochInitialIssuance)
	}

	// store tax caps
	for _, cap := range data.TaxCaps {
		keeper.SetTaxCap(ctx, cap.Denom, cap.TaxCap)
	}

	for _, epochState := range data.EpochStates {
		keeper.SetTR(ctx, int64(epochState.Epoch), epochState.TaxReward)
		keeper.SetSR(ctx, int64(epochState.Epoch), epochState.SeigniorageReward)
		keeper.SetTSL(ctx, int64(epochState.Epoch), epochState.TotalStakedLuna)
	}

	// check if the module account exists
	moduleAcc := keeper.GetTreasuryModuleAccount(ctx)
	if moduleAcc == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.ModuleName))
	}

	// check if the burn module account exists
	burnModuleAcc := keeper.GetBurnModuleAccount(ctx)
	if burnModuleAcc == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.BurnModuleName))
	}
}

// ExportGenesis writes the current store values
// to a genesis file, which can be imported again
// with InitGenesis
func ExportGenesis(ctx sdk.Context, keeper keeper.Keeper) (data *types.GenesisState) {
	params := keeper.GetParams(ctx)

	taxRate := keeper.GetTaxRate(ctx)
	rewardWeight := keeper.GetRewardWeight(ctx)
	taxProceeds := keeper.PeekEpochTaxProceeds(ctx)
	epochInitialIssuance := keeper.GetEpochInitialIssuance(ctx)

	var taxCaps []types.TaxCap
	keeper.IterateTaxCap(ctx, func(denom string, taxCap sdk.Int) bool {
		taxCaps = append(taxCaps, types.TaxCap{
			Denom:  denom,
			TaxCap: taxCap,
		})
		return false
	})

	var epochStates []types.EpochState

	curEpoch := keeper.GetEpoch(ctx)
	for e := int64(0); e < curEpoch ||
		(e == curEpoch && core.IsPeriodLastBlock(ctx, core.BlocksPerWeek)); e++ {
		epochStates = append(epochStates, types.EpochState{
			Epoch:             uint64(e),
			TaxReward:         keeper.GetTR(ctx, e),
			SeigniorageReward: keeper.GetSR(ctx, e),
			TotalStakedLuna:   keeper.GetTSL(ctx, e),
		})
	}

	return types.NewGenesisState(params, taxRate, rewardWeight,
		taxCaps, taxProceeds, epochInitialIssuance, epochStates)
}
