package types

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"
)

type GenesisTestSuite struct {
	suite.Suite
	address1  string
	address2  string
	contractA string
	contractB string
}

func TestGenesisTestSuite(t *testing.T) {
	suite.Run(t, new(GenesisTestSuite))
}

func (suite *GenesisTestSuite) SetupTest() {
	suite.address1 = sdk.AccAddress([]byte("cosmos1")).String()
	suite.address2 = sdk.AccAddress([]byte("cosmos2")).String()

	suite.contractA = "cosmos15u3dt79t6sxxa3x3kpkhzsy56edaa5a66wvt3kxmukqjz2sx0hesh45zsv"
	suite.contractB = "cosmos168ctmpyppk90d34p3jjy658zf5a5l3w8wk35wht6ccqj4mr0yv8skhnwe8"
}

func (suite *GenesisTestSuite) TestValidateGenesis() {
	newGen := NewGenesisState(DefaultParams(), []FeeShare{})
	testCases := []struct {
		name     string
		genState *GenesisState
		expPass  bool
	}{
		{
			name:     "valid genesis constructor",
			genState: newGen,
			expPass:  true,
		},
		{
			name:     "default",
			genState: DefaultGenesisState(),
			expPass:  true,
		},
		{
			name: "valid genesis",
			genState: &GenesisState{
				Params:   DefaultParams(),
				FeeShare: []FeeShare{},
			},
			expPass: true,
		},
		{
			name: "valid genesis - with fee",
			genState: &GenesisState{
				Params: DefaultParams(),
				FeeShare: []FeeShare{
					{
						ContractAddress:   suite.contractA,
						DeployerAddress:   suite.address1,
						WithdrawerAddress: suite.address1,
					},
					{
						ContractAddress:   suite.contractB,
						DeployerAddress:   suite.address2,
						WithdrawerAddress: suite.address2,
					},
				},
			},
			expPass: true,
		},
		{
			name:     "empty genesis",
			genState: &GenesisState{},
			expPass:  false,
		},
		{
			name: "invalid genesis - duplicated fee",
			genState: &GenesisState{
				Params: DefaultParams(),
				FeeShare: []FeeShare{
					{
						ContractAddress: suite.contractA,
						DeployerAddress: suite.address1,
					},
					{
						ContractAddress: suite.contractA,
						DeployerAddress: suite.address1,
					},
				},
			},
			expPass: false,
		},
		{
			name: "invalid genesis - duplicated fee with different deployer address",
			genState: &GenesisState{
				Params: DefaultParams(),
				FeeShare: []FeeShare{
					{
						ContractAddress: suite.contractA,
						DeployerAddress: suite.address1,
					},
					{
						ContractAddress: suite.contractA,
						DeployerAddress: suite.address2,
					},
				},
			},
			expPass: false,
		},
		{
			name: "invalid genesis - invalid withdrawer address",
			genState: &GenesisState{
				Params: DefaultParams(),
				FeeShare: []FeeShare{
					{
						ContractAddress:   suite.contractA,
						DeployerAddress:   suite.address1,
						WithdrawerAddress: "withdraw",
					},
				},
			},
			expPass: false,
		},
	}

	for _, tc := range testCases {
		err := tc.genState.Validate()
		if tc.expPass {
			suite.Require().NoError(err, tc.name)
		} else {
			suite.Require().Error(err, tc.name)
		}
	}
}
