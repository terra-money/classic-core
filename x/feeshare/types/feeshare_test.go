package types

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"
)

type FeeShareTestSuite struct {
	suite.Suite
	address1 sdk.AccAddress
	address2 sdk.AccAddress
	contract sdk.AccAddress
}

func TestFeeShareSuite(t *testing.T) {
	suite.Run(t, new(FeeShareTestSuite))
}

func (suite *FeeShareTestSuite) SetupTest() {
	suite.address1 = sdk.AccAddress([]byte("cosmos1"))
	suite.address2 = sdk.AccAddress([]byte("cosmos2"))

	suite.contract = sdk.AccAddress([]byte("cosmos1contract"))
}

func (suite *FeeShareTestSuite) TestFeeNew() {
	testCases := []struct {
		name       string
		contract   sdk.Address
		deployer   sdk.AccAddress
		withdraw   sdk.AccAddress
		expectPass bool
	}{
		{
			"Create feeshare- pass",
			suite.contract,
			suite.address1,
			suite.address2,
			true,
		},
		{
			"Create feeshare- invalid contract address",
			sdk.AccAddress{},
			suite.address1,
			suite.address2,
			false,
		},
		{
			"Create feeshare- invalid deployer address",
			suite.contract,
			sdk.AccAddress{},
			suite.address2,
			false,
		},
	}

	for _, tc := range testCases {
		i := NewFeeShare(tc.contract, tc.deployer, tc.withdraw)
		err := i.Validate()

		if tc.expectPass {
			suite.Require().NoError(err, tc.name)
		} else {
			suite.Require().Error(err, tc.name)
		}
	}
}

func (suite *FeeShareTestSuite) TestFee() {
	testCases := []struct {
		msg        string
		feeshare   FeeShare
		expectPass bool
	}{
		{
			"Create feeshare- pass",
			FeeShare{
				suite.contract.String(),
				suite.address1.String(),
				suite.address2.String(),
			},
			true,
		},
		{
			"Create feeshare- invalid contract address (invalid length 2)",
			FeeShare{
				"terra15u3dt79t6sxxa3x3kpkhzsy56edaa5a66kxmukqjz2sx0hes5sn38g",
				suite.address1.String(),
				suite.address2.String(),
			},
			false,
		},
		{
			"Create feeshare- invalid deployer address",
			FeeShare{
				suite.contract.String(),
				"terra1hj5fveer5cjtn4wd6wstzugjfdxzl0xps73ftl",
				suite.address2.String(),
			},
			false,
		},
		{
			"Create feeshare- invalid withdraw address",
			FeeShare{
				suite.contract.String(),
				suite.address1.String(),
				"terra1hj5fveer5cjtn4wd6wstzugjfdxzl0xps73ftl",
			},
			false,
		},
	}

	for _, tc := range testCases {
		err := tc.feeshare.Validate()

		if tc.expectPass {
			suite.Require().NoError(err, tc.msg)
		} else {
			suite.Require().Error(err, tc.msg)
		}
	}
}

func (suite *FeeShareTestSuite) TestFeeShareGetters() {
	contract := sdk.AccAddress([]byte("cosmos1contract"))
	fs := FeeShare{
		contract.String(),
		suite.address1.String(),
		suite.address2.String(),
	}
	suite.Equal(fs.GetContractAddr(), contract)
	suite.Equal(fs.GetDeployerAddr(), suite.address1)
	suite.Equal(fs.GetWithdrawerAddr(), suite.address2)

	fs = FeeShare{
		contract.String(),
		suite.address1.String(),
		"",
	}
	suite.Equal(fs.GetContractAddr(), contract)
	suite.Equal(fs.GetDeployerAddr(), suite.address1)
	suite.Equal(len(fs.GetWithdrawerAddr()), 0)
}
