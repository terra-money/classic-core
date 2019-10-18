package keeper

import (
	"testing"

	core "github.com/terra-project/core/types"
	"github.com/terra-project/core/x/treasury/internal/types"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
)

func TestUpdateTaxRate(t *testing.T) {
	input := CreateTestInput(t)
	sh := staking.NewHandler(input.StakingKeeper)

	// Create Validators
	amt := sdk.TokensFromConsensusPower(1)
	addr, val := ValAddrs[0], PubKeys[0]
	addr1, val1 := ValAddrs[1], PubKeys[1]
	res := sh(input.Ctx, NewTestMsgCreateValidator(addr, val, amt))
	require.True(t, res.IsOK())
	res = sh(input.Ctx, NewTestMsgCreateValidator(addr1, val1, amt))
	require.True(t, res.IsOK())
	staking.EndBlocker(input.Ctx, input.StakingKeeper)

	windowLong := input.TreasuryKeeper.WindowLong(input.Ctx)
	taxPolicy := input.TreasuryKeeper.TaxPolicy(input.Ctx)

	// zero reward tax proceeds
	for i := int64(0); i < windowLong; i++ {
		input.Ctx = input.Ctx.WithBlockHeight(i * core.BlocksPerEpoch)

		taxProceeds := sdk.NewCoins(sdk.NewCoin(core.MicroSDRDenom, sdk.ZeroInt()))
		input.TreasuryKeeper.RecordTaxProceeds(input.Ctx, taxProceeds)
	}

	input.TreasuryKeeper.UpdateTaxPolicy(input.Ctx)
	taxRate := input.TreasuryKeeper.GetTaxRate(input.Ctx, core.GetEpoch(input.Ctx)+1)
	require.Equal(t, types.DefaultTaxRate.Add(taxPolicy.ChangeRateMax), taxRate)
}

func TestUpdateRewardWeight(t *testing.T) {
	input := CreateTestInput(t)
	sh := staking.NewHandler(input.StakingKeeper)

	// Create Validators
	amt := sdk.TokensFromConsensusPower(1)
	addr, val := ValAddrs[0], PubKeys[0]
	addr1, val1 := ValAddrs[1], PubKeys[1]
	res := sh(input.Ctx, NewTestMsgCreateValidator(addr, val, amt))
	require.True(t, res.IsOK())
	res = sh(input.Ctx, NewTestMsgCreateValidator(addr1, val1, amt))
	require.True(t, res.IsOK())
	staking.EndBlocker(input.Ctx, input.StakingKeeper)

	rewardPolicy := input.TreasuryKeeper.RewardPolicy(input.Ctx)
	input.TreasuryKeeper.UpdateRewardPolicy(input.Ctx)
	rewardWeight := input.TreasuryKeeper.GetRewardWeight(input.Ctx, core.GetEpoch(input.Ctx)+1)
	require.Equal(t, types.DefaultRewardWeight.Add(rewardPolicy.ChangeRateMax), rewardWeight)
}

func TestUpdateTaxCap(t *testing.T) {
	input := CreateTestInput(t)
	input.SupplyKeeper.SetSupply(input.Ctx,
		input.SupplyKeeper.GetSupply(input.Ctx).SetTotal(
			sdk.NewCoins(
				sdk.NewInt64Coin(core.MicroLunaDenom, 1000000),
				sdk.NewInt64Coin(core.MicroSDRDenom, 1000000),
				sdk.NewInt64Coin(core.MicroKRWDenom, 1000000),
			),
		),
	)

	// Create Validators
	sdrPrice := sdk.NewDecWithPrec(13, 1)
	input.OracleKeeper.SetLunaPrice(input.Ctx, core.MicroSDRDenom, sdrPrice)
	krwPrice := sdk.NewDecWithPrec(153412, 2)
	input.OracleKeeper.SetLunaPrice(input.Ctx, core.MicroKRWDenom, krwPrice)
	input.TreasuryKeeper.UpdateTaxCap(input.Ctx)

	krwCap := input.TreasuryKeeper.GetTaxCap(input.Ctx, core.MicroKRWDenom)
	sdrCapAmt := input.TreasuryKeeper.GetParams(input.Ctx).TaxPolicy.Cap.Amount
	require.Equal(t, krwCap, krwPrice.Quo(sdrPrice).MulInt(sdrCapAmt).TruncateInt())
}
