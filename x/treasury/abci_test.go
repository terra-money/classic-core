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

	targetIssuance := sdk.NewInt(1000)
	input.Ctx = input.Ctx.WithBlockHeight(core.BlocksPerEpoch - 1)
	supply := input.SupplyKeeper.GetSupply(input.Ctx)
	supply = supply.SetTotal(sdk.NewCoins(sdk.NewCoin(core.MicroLunaDenom, targetIssuance)))
	input.SupplyKeeper.SetSupply(input.Ctx, supply)
	EndBlocker(input.Ctx, input.TreasuryKeeper)

	issuance := input.TreasuryKeeper.GetHistoricalIssuance(input.Ctx, 0).AmountOf(core.MicroLunaDenom)
	require.Equal(t, targetIssuance, issuance)
}

func TestUpdate(t *testing.T) {
	input := keeper.CreateTestInput(t)
	windowProbation := input.TreasuryKeeper.WindowProbation(input.Ctx)
	bondedModuleAcc := input.SupplyKeeper.GetModuleAccount(input.Ctx, stakingtypes.BondedPoolName)
	err := bondedModuleAcc.SetCoins(sdk.NewCoins(sdk.NewInt64Coin(core.MicroLunaDenom, 100000000000)))
	require.NoError(t, err)
	input.SupplyKeeper.SetModuleAccount(input.Ctx, bondedModuleAcc)

	targetEpoch := windowProbation + 1
	input.Ctx = input.Ctx.WithBlockHeight(core.BlocksPerEpoch*targetEpoch - 1)

	// zero tax proceeds will increase tax rate with change max amount
	EndBlocker(input.Ctx, input.TreasuryKeeper)
	taxRate := input.TreasuryKeeper.GetTaxRate(input.Ctx, targetEpoch-1)
	newTaxRate := input.TreasuryKeeper.GetTaxRate(input.Ctx, targetEpoch)
	require.Equal(t, taxRate.Add(input.TreasuryKeeper.TaxPolicy(input.Ctx).ChangeRateMax), newTaxRate)

	// zero mining rewards will increase reward weight with change max amount
	rewardWeight := input.TreasuryKeeper.GetRewardWeight(input.Ctx, targetEpoch-1)
	newRewardWeight := input.TreasuryKeeper.GetRewardWeight(input.Ctx, targetEpoch)
	require.Equal(t, rewardWeight.Add(input.TreasuryKeeper.RewardPolicy(input.Ctx).ChangeRateMax), newRewardWeight)
}
