package keeper

import (
	"fmt"
	"math/rand"
	"testing"

	core "github.com/terra-project/core/types"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestSettle(t *testing.T) {
	input := CreateTestInput(t)

	input.OracleKeeper.SetLunaPrice(input.Ctx, core.MicroSDRDenom, sdk.OneDec())

	issuance := sdk.NewInt(rand.Int63() + 1)
	supply := input.SupplyKeeper.GetSupply(input.Ctx)
	supply = supply.SetTotal(sdk.NewCoins(sdk.NewCoin(core.MicroLunaDenom, issuance)))
	input.SupplyKeeper.SetSupply(input.Ctx, supply)
	input.TreasuryKeeper.RecordHistoricalIssuance(input.Ctx)

	input.Ctx = input.Ctx.WithBlockHeight(core.BlocksPerEpoch)
	supply = supply.SetTotal(sdk.NewCoins(sdk.NewCoin(core.MicroLunaDenom, sdk.ZeroInt())))
	input.SupplyKeeper.SetSupply(input.Ctx, supply)

	// check seigniorage update
	require.Equal(t, issuance, input.TreasuryKeeper.PeekEpochSeigniorage(input.Ctx, 1))

	input.TreasuryKeeper.SettleSeigniorage(input.Ctx)
	oracleAcc := input.SupplyKeeper.GetModuleAccount(input.Ctx, input.TreasuryKeeper.oracleModuleName)
	feePool := input.DistrKeeper.GetFeePool(input.Ctx)
	fmt.Println(oracleAcc)
	fmt.Println(feePool)

	rewardWeight := input.TreasuryKeeper.GetRewardWeight(input.Ctx, 1)
	oracleRewardAmt := rewardWeight.MulInt(issuance).TruncateInt()
	leftAmt := issuance.Sub(oracleRewardAmt)

	require.Equal(t, oracleRewardAmt, oracleAcc.GetCoins().AmountOf(core.MicroSDRDenom))
	require.Equal(t, leftAmt, feePool.CommunityPool.AmountOf(core.MicroSDRDenom).TruncateInt())
}
