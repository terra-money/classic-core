package treasury

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	core "github.com/terra-project/core/types"
	"github.com/terra-project/core/x/treasury/internal/keeper"
)

func TestEndBlockerIssuanceUpdate(t *testing.T) {
	input := keeper.CreateTestInput(t)

	// Set total staked luna to prevent divide by zero error when computing TRL
	bondedModuleAcc := input.SupplyKeeper.GetModuleAccount(input.Ctx, stakingtypes.BondedPoolName)
	err := bondedModuleAcc.SetCoins(sdk.NewCoins(sdk.NewInt64Coin(core.MicroLunaDenom, 100000000000)))
	require.NoError(t, err)
	input.SupplyKeeper.SetModuleAccount(input.Ctx, bondedModuleAcc)

	targetIssuance := sdk.NewInt(1000)
	input.Ctx = input.Ctx.WithBlockHeight(core.BlocksPerWeek - 1)
	supply := input.SupplyKeeper.GetSupply(input.Ctx)
	supply = supply.SetTotal(sdk.NewCoins(sdk.NewCoin(core.MicroLunaDenom, targetIssuance)))
	input.SupplyKeeper.SetSupply(input.Ctx, supply)
	EndBlocker(input.Ctx, input.TreasuryKeeper)

	issuance := input.TreasuryKeeper.GetEpochInitialIssuance(input.Ctx).AmountOf(core.MicroLunaDenom)
	require.Equal(t, targetIssuance, issuance)
}

func TestUpdate(t *testing.T) {
	input := keeper.CreateTestInput(t)

	windowProbation := input.TreasuryKeeper.WindowProbation(input.Ctx)

	// Set total staked luna to prevent divide by zero error when computing TRL
	bondedModuleAcc := input.SupplyKeeper.GetModuleAccount(input.Ctx, stakingtypes.BondedPoolName)
	err := bondedModuleAcc.SetCoins(sdk.NewCoins(sdk.NewInt64Coin(core.MicroLunaDenom, 100000000000)))
	require.NoError(t, err)
	input.SupplyKeeper.SetModuleAccount(input.Ctx, bondedModuleAcc)

	targetEpoch := windowProbation + 1
	for epoch := int64(0); epoch < targetEpoch; epoch++ {
		input.Ctx = input.Ctx.WithBlockHeight(core.BlocksPerWeek*epoch - 1)
		EndBlocker(input.Ctx, input.TreasuryKeeper)
	}

	// load old tax rate & reward weight
	taxRate := input.TreasuryKeeper.GetTaxRate(input.Ctx)
	rewardWeight := input.TreasuryKeeper.GetRewardWeight(input.Ctx)

	input.Ctx = input.Ctx.WithBlockHeight(core.BlocksPerWeek*targetEpoch - 1)
	EndBlocker(input.Ctx, input.TreasuryKeeper)

	// zero tax proceeds will increase tax rate with change max amount
	newTaxRate := input.TreasuryKeeper.GetTaxRate(input.Ctx)
	require.Equal(t, taxRate.Add(input.TreasuryKeeper.TaxPolicy(input.Ctx).ChangeRateMax), newTaxRate)

	// zero mining rewards will increase reward weight with change max amount
	newRewardWeight := input.TreasuryKeeper.GetRewardWeight(input.Ctx)
	require.Equal(t, rewardWeight.Add(input.TreasuryKeeper.RewardPolicy(input.Ctx).ChangeRateMax), newRewardWeight)
}
