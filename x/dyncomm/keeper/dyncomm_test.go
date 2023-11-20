package keeper

import (
	"testing"
	"time"

	core "github.com/classic-terra/core/v2/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking/teststaking"
	"github.com/stretchr/testify/require"
)

func TestCalculateVotingPower(t *testing.T) {
	input := CreateTestInput(t)
	helper := teststaking.NewHelper(
		t, input.Ctx, input.StakingKeeper,
	)
	helper.Denom = core.MicroLunaDenom
	helper.CreateValidatorWithValPower(ValAddrFrom(0), PubKeys[0], 9, true)
	helper.CreateValidatorWithValPower(ValAddrFrom(1), PubKeys[1], 1, true)
	helper.TurnBlock(time.Now())
	vals := input.StakingKeeper.GetBondedValidatorsByPower(input.Ctx)

	require.Equal(
		t,
		sdk.NewDecWithPrec(90, 0),
		input.DyncommKeeper.CalculateVotingPower(input.Ctx, vals[0]),
	)
}

func TestCalculateDynCommission(t *testing.T) {
	input := CreateTestInput(t)
	helper := teststaking.NewHelper(
		t, input.Ctx, input.StakingKeeper,
	)
	helper.Denom = core.MicroLunaDenom
	helper.CreateValidatorWithValPower(ValAddrFrom(0), PubKeys[0], 950, true)
	helper.CreateValidatorWithValPower(ValAddrFrom(1), PubKeys[1], 46, true)
	helper.CreateValidatorWithValPower(ValAddrFrom(2), PubKeys[2], 4, true)
	helper.TurnBlock(time.Now())
	vals := input.StakingKeeper.GetBondedValidatorsByPower(input.Ctx)

	// capped commission
	require.Equal(
		t,
		sdk.NewDecWithPrec(20, 2),
		input.DyncommKeeper.CalculateDynCommission(input.Ctx, vals[0]),
	)

	// curve
	require.Equal(
		t,
		sdk.NewDecWithPrec(10086, 5),
		input.DyncommKeeper.CalculateDynCommission(input.Ctx, vals[1]),
	)

	// min. commission
	require.Equal(
		t,
		sdk.ZeroDec(),
		input.DyncommKeeper.CalculateDynCommission(input.Ctx, vals[2]),
	)
}
