package pay

import (
	"terra/types/assets"

	//"terra/x/treasury"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/mock"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto"

	abci "github.com/tendermint/tendermint/abci/types"
)

// initialize the mock application for this module
func getMockApp(t *testing.T, numGenAccs int) (*mock.App, Keeper, []sdk.AccAddress, []crypto.PubKey, []crypto.PrivKey) {
	mapp := mock.NewApp()

	keyBank := sdk.NewKVStoreKey("bank")
	keyFeeCollection := sdk.NewKVStoreKey("fee")

	fck := auth.NewFeeCollectionKeeper(mapp.Cdc, keyFeeCollection)
	keeper := NewKeeper(keyBank, mapp.Cdc, mapp.AccountKeeper, fck)

	mapp.SetEndBlocker(getEndBlocker(keeper))
	mapp.SetInitChainer(getInitChainer(mapp, keeper))
	require.NoError(t, mapp.CompleteSetup(keyBank, keyFeeCollection))

	genAccs, addrs, pubKeys, privKeys := mock.CreateGenAccounts(numGenAccs, sdk.Coins{sdk.NewInt64Coin(assets.LunaDenom, 42)})
	mock.SetGenesis(mapp, genAccs)

	return mapp, keeper, addrs, pubKeys, privKeys
}

// oracle and staking endblocker
func getEndBlocker(keeper Keeper) sdk.EndBlocker {
	return func(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {

		return abci.ResponseEndBlock{
			//Tags: tags,
		}
	}
}

// oracle and staking initchainer
func getInitChainer(mapp *mock.App, keeper Keeper) sdk.InitChainer {
	return func(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
		mapp.InitChainer(ctx, req)

		// stakingGenesis := staking.GenesisState{
		// 	Pool: staking.InitialPool(),
		// 	Params: staking.Params{
		// 		UnbondingTime: 60 * 60 * 24 * 3 * time.Second,
		// 		MaxValidators: 100,
		// 		BondDenom:     assets.LunaDenom,
		// 	},
		// }
		// stakingGenesis.Pool.NotBondedTokens = sdk.NewInt(100000)

		// validators, err := staking.InitGenesis(ctx, stakingKeeper, stakingGenesis)
		// if err != nil {
		// 	panic(err)
		// }

		// InitGenesis(ctx, keeper, GenesisState{
		// 	Params: NewParams(
		// 		sdk.NewInt(1),             // one block
		// 		sdk.NewDecWithPrec(66, 2), // 66%
		// 		sdk.NewInt(10),
		// 	),
		// })

		return abci.ResponseInitChain{
			//Validators: validators,
		}
	}
}
