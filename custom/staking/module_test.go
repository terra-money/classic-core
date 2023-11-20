package staking_test

import (
	"testing"

	apptesting "github.com/classic-terra/core/v2/app/testing"
	simapp "github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	disttypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	teststaking "github.com/cosmos/cosmos-sdk/x/staking/teststaking"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/suite"
)

type StakingTestSuite struct {
	apptesting.KeeperTestHelper
}

func TestStakingTestSuite(t *testing.T) {
	suite.Run(t, new(StakingTestSuite))
}

// go test -v -run=TestStakingTestSuite/TestValidatorVPLimit github.com/classic-terra/core/v2/custom/staking
func (s *StakingTestSuite) TestValidatorVPLimit() {
	s.KeeperTestHelper.Setup(s.T())

	// construct new validators, to a total of 10 validators, each with 10% of the total voting power
	num := 9
	addrDels := s.RandomAccountAddresses(num)
	for i, addrDel := range addrDels {
		s.FundAcc(addrDel, sdk.NewCoins(sdk.NewInt64Coin("uluna", 1000000)))
		err := s.App.BankKeeper.DelegateCoinsFromAccountToModule(s.Ctx, addrDels[i], stakingtypes.NotBondedPoolName, sdk.NewCoins(sdk.NewInt64Coin("uluna", 1000000)))
		s.Require().NoError(err)
	}
	valAddrs := simapp.ConvertAddrsToValAddrs(addrDels)
	PKs := simapp.CreateTestPubKeys(num)

	var amts [9]sdk.Int
	for i := range amts {
		amts[i] = sdk.NewInt(1000000)
	}

	var validators [9]stakingtypes.Validator
	for i, amt := range amts {
		validators[i] = teststaking.NewValidator(s.T(), valAddrs[i], PKs[i])
		validators[i], _ = validators[i].AddTokensFromDel(amt)
	}

	for i := range validators {
		validators[i] = stakingkeeper.TestingUpdateValidator(s.App.StakingKeeper, s.Ctx, validators[i], true)
	}

	// delegate to a validator over 20% VP
	s.FundAcc(s.TestAccs[0], sdk.NewCoins(sdk.NewInt64Coin("uluna", 2000000)))
	s.App.DistrKeeper.SetValidatorHistoricalRewards(s.Ctx, valAddrs[0], 1, disttypes.NewValidatorHistoricalRewards(sdk.NewDecCoins(sdk.NewDecCoin("uluna", sdk.NewInt(1))), 2))
	s.App.DistrKeeper.SetValidatorCurrentRewards(s.Ctx, valAddrs[0], disttypes.NewValidatorCurrentRewards(sdk.NewDecCoins(sdk.NewDecCoin("uluna", sdk.NewInt(1))), 2))
	s.App.DistrKeeper.SetDelegatorStartingInfo(s.Ctx, valAddrs[0], s.TestAccs[0], disttypes.NewDelegatorStartingInfo(1, sdk.OneDec(), 1))

	// first delegation should be normal
	// raise voting power of validator 0 by 1 (1+1)/(10+1) = 0.181818 < 0.2
	s.App.StakingKeeper.SetDelegation(s.Ctx, stakingtypes.NewDelegation(s.TestAccs[0], valAddrs[0], sdk.NewDec(1000000)))
	_, err := s.App.StakingKeeper.Delegate(s.Ctx, s.TestAccs[0], sdk.NewInt(1000000), stakingtypes.Unbonded, validators[0], true)
	s.Require().NoError(err)

	// update validator set and validator 0 state
	_, err = s.App.StakingKeeper.ApplyAndReturnValidatorSetUpdates(s.Ctx)
	s.Require().NoError(err)
	validator, found := s.App.StakingKeeper.GetValidator(s.Ctx, valAddrs[0])
	s.Require().True(found)
	validators[0] = validator

	// test that the code panic on a delegation that surpasses the 20% VP allowed
	// raise voting power of validator 0 by 1 (2+1)/(11+1) = 0.250000 > 0.2
	defer func() {
		if r := recover(); r == nil {
			s.T().Errorf("The code did not panic")
		}
	}()

	s.App.StakingKeeper.SetDelegation(s.Ctx, stakingtypes.NewDelegation(s.TestAccs[0], valAddrs[0], sdk.NewDec(1000000)))
	s.App.StakingKeeper.Delegate(s.Ctx, s.TestAccs[0], sdk.NewInt(1000000), stakingtypes.Unbonded, validators[0], true)
}
