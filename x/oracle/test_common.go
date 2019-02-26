package oracle

import (
	"terra/types/assets"
	"terra/x/pay"
	"terra/x/treasury"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/mock"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto"

	abci "github.com/tendermint/tendermint/abci/types"
)

// initialize the mock application for this module
func getMockApp(t *testing.T, numGenAccs int) (*mock.App, Keeper, staking.Keeper, []sdk.AccAddress, []crypto.PubKey, []crypto.PrivKey) {
	mapp := mock.NewApp()

	staking.RegisterCodec(mapp.Cdc)

	keyGlobalParams := sdk.NewKVStoreKey("params")
	tkeyGlobalParams := sdk.NewTransientStoreKey("transient_params")
	keystaking := sdk.NewKVStoreKey("staking")
	tkeystaking := sdk.NewTransientStoreKey("transient_staking")
	keyTreasury := sdk.NewKVStoreKey("treasury")
	keyOracle := sdk.NewKVStoreKey("oracle")
	keyBank := sdk.NewKVStoreKey("bank")
	keyFeeCollection := sdk.NewKVStoreKey("fee")

	pk := params.NewKeeper(mapp.Cdc, keyGlobalParams, tkeyGlobalParams)
	fck := auth.NewFeeCollectionKeeper(mapp.Cdc, keyFeeCollection)
	ck := pay.NewKeeper(keyBank, mapp.Cdc, mapp.AccountKeeper, fck)
	sk := staking.NewKeeper(mapp.Cdc, keystaking, tkeystaking, ck, pk.Subspace(staking.DefaultParamspace), staking.DefaultCodespace)
	tk := treasury.NewKeeper(keyTreasury, mapp.Cdc, ck, pk.Subspace(treasury.DefaultParamspace))
	keeper := NewKeeper(keyOracle, mapp.Cdc, tk, sk.GetValidatorSet(), pk.Subspace(DefaultParamspace))

	mapp.Router().AddRoute("oracle", NewHandler(keeper))

	mapp.SetEndBlocker(getEndBlocker(keeper))
	mapp.SetInitChainer(getInitChainer(mapp, keeper, sk))

	require.NoError(t, mapp.CompleteSetup(keystaking, tkeystaking,
		keyTreasury,
		keyOracle, keyBank, keyFeeCollection,
		keyGlobalParams, tkeyGlobalParams))

	genAccs, addrs, pubKeys, privKeys := mock.CreateGenAccounts(numGenAccs, sdk.Coins{sdk.NewInt64Coin(assets.LunaDenom, 42)})

	mock.SetGenesis(mapp, genAccs)

	return mapp, keeper, sk, addrs, pubKeys, privKeys
}

// oracle and staking endblocker
func getEndBlocker(keeper Keeper) sdk.EndBlocker {
	return func(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
		_, _, tags := EndBlocker(ctx, keeper)
		return abci.ResponseEndBlock{
			Tags: tags,
		}
	}
}

// oracle and staking initchainer
func getInitChainer(mapp *mock.App, keeper Keeper, stakingKeeper staking.Keeper) sdk.InitChainer {
	return func(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
		mapp.InitChainer(ctx, req)

		stakingGenesis := staking.GenesisState{
			Pool: staking.InitialPool(),
			Params: staking.Params{
				UnbondingTime: 60 * 60 * 24 * 3 * time.Second,
				MaxValidators: 100,
				BondDenom:     assets.LunaDenom,
			},
		}
		stakingGenesis.Pool.NotBondedTokens = sdk.NewInt(100000)

		validators, err := staking.InitGenesis(ctx, stakingKeeper, stakingGenesis)
		if err != nil {
			panic(err)
		}

		InitGenesis(ctx, keeper, GenesisState{
			Params: NewParams(
				sdk.NewInt(1),             // one block
				sdk.NewDecWithPrec(66, 2), // 66%
				sdk.NewInt(10),
			),
		})

		return abci.ResponseInitChain{
			Validators: validators,
		}
	}
}
