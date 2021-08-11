package keeper

import (
	"testing"

	core "github.com/terra-money/core/types"
	oracletypes "github.com/terra-money/core/x/oracle/types"
	"github.com/terra-money/core/x/treasury/types"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
)

func TestUpdateTaxRate(t *testing.T) {
	input := CreateTestInput(t)
	sh := staking.NewHandler(input.StakingKeeper)

	// Create Validators
	amt := sdk.TokensFromConsensusPower(1, sdk.DefaultPowerReduction)
	addr, val := ValAddrs[0], ValPubKeys[0]
	addr1, val1 := ValAddrs[1], ValPubKeys[1]
	_, err := sh(input.Ctx, NewTestMsgCreateValidator(addr, val, amt))
	require.NoError(t, err)
	_, err = sh(input.Ctx, NewTestMsgCreateValidator(addr1, val1, amt))
	require.NoError(t, err)
	staking.EndBlocker(input.Ctx, input.StakingKeeper)

	windowLong := input.TreasuryKeeper.WindowLong(input.Ctx)
	taxPolicy := input.TreasuryKeeper.TaxPolicy(input.Ctx)

	// zero reward tax proceeds
	for i := uint64(0); i < windowLong; i++ {
		input.Ctx = input.Ctx.WithBlockHeight(int64(i * core.BlocksPerWeek))

		taxProceeds := sdk.NewCoins(sdk.NewCoin(core.MicroSDRDenom, sdk.ZeroInt()))
		input.TreasuryKeeper.RecordEpochTaxProceeds(input.Ctx, taxProceeds)
		input.TreasuryKeeper.UpdateIndicators(input.Ctx)
	}

	input.TreasuryKeeper.UpdateTaxPolicy(input.Ctx)
	taxRate := input.TreasuryKeeper.GetTaxRate(input.Ctx)
	require.Equal(t, types.DefaultTaxRate.Add(taxPolicy.ChangeRateMax), taxRate)
}

func TestUpdateRewardWeight(t *testing.T) {
	input := CreateTestInput(t)
	input.OracleKeeper.SetLunaExchangeRate(input.Ctx, core.MicroSDRDenom, sdk.OneDec())
	sh := staking.NewHandler(input.StakingKeeper)

	// Create Validators
	amt := sdk.TokensFromConsensusPower(1, sdk.DefaultPowerReduction)
	addr, val := ValAddrs[0], ValPubKeys[0]
	addr1, val1 := ValAddrs[1], ValPubKeys[1]
	_, err := sh(input.Ctx, NewTestMsgCreateValidator(addr, val, amt))
	require.NoError(t, err)
	_, err = sh(input.Ctx, NewTestMsgCreateValidator(addr1, val1, amt))
	require.NoError(t, err)
	staking.EndBlocker(input.Ctx, input.StakingKeeper)

	input.TreasuryKeeper.UpdateIndicators(input.Ctx)

	// Case 1: zero seigniorage will increase reward weight as much as possible
	rewardPolicy := input.TreasuryKeeper.RewardPolicy(input.Ctx)
	input.TreasuryKeeper.UpdateRewardPolicy(input.Ctx)
	rewardWeight := input.TreasuryKeeper.GetRewardWeight(input.Ctx)
	require.Equal(t, types.DefaultRewardWeight.Add(rewardPolicy.ChangeRateMax), rewardWeight)

	// Case 2: huge seigniorage rewards will decrease reward weight by %types.DefaultSeigniorageBurdenTarget
	input.TreasuryKeeper.SetEpochInitialIssuance(input.Ctx, sdk.NewCoins(sdk.NewCoin(core.MicroLunaDenom, sdk.NewInt(1000000000000))))
	input.TreasuryKeeper.UpdateIndicators(input.Ctx)
	input.TreasuryKeeper.UpdateRewardPolicy(input.Ctx)
	rewardWeight = input.TreasuryKeeper.GetRewardWeight(input.Ctx)
	require.Equal(t, types.DefaultRewardWeight.Add(rewardPolicy.ChangeRateMax).Mul(types.DefaultSeigniorageBurdenTarget), rewardWeight)
}

func TestUpdateTaxCap(t *testing.T) {
	input := CreateTestInput(t)
	input.OracleKeeper.SetWhitelist(
		input.Ctx,
		oracletypes.DenomList{
			{
				Name: core.MicroLunaDenom,
			},
			{
				Name: core.MicroSDRDenom,
			},
			{
				Name: core.MicroKRWDenom,
			},
		},
	)

	// Create Validators
	sdrPrice := sdk.NewDecWithPrec(13, 1)
	input.OracleKeeper.SetLunaExchangeRate(input.Ctx, core.MicroSDRDenom, sdrPrice)
	krwPrice := sdk.NewDecWithPrec(153412, 2)
	input.OracleKeeper.SetLunaExchangeRate(input.Ctx, core.MicroKRWDenom, krwPrice)
	input.TreasuryKeeper.UpdateTaxCap(input.Ctx)

	krwCap := input.TreasuryKeeper.GetTaxCap(input.Ctx, core.MicroKRWDenom)
	sdrCapAmt := input.TreasuryKeeper.GetParams(input.Ctx).TaxPolicy.Cap.Amount
	require.Equal(t, krwCap, krwPrice.Quo(sdrPrice).MulInt(sdrCapAmt).TruncateInt())
}
