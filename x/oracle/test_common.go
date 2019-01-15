package oracle

import (
	"terra/types/assets"
	"terra/types/tax"
	"terra/x/treasury"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/mock"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/stake"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto"

	abci "github.com/tendermint/tendermint/abci/types"
)

// initialize the mock application for this module
func getMockApp(t *testing.T, numGenAccs int) (*mock.App, Keeper, stake.Keeper, []sdk.AccAddress, []crypto.PubKey, []crypto.PrivKey) {
	mapp := mock.NewApp()

	stake.RegisterCodec(mapp.Cdc)

	keyGlobalParams := sdk.NewKVStoreKey("params")
	tkeyGlobalParams := sdk.NewTransientStoreKey("transient_params")
	keyStake := sdk.NewKVStoreKey("stake")
	tkeyStake := sdk.NewTransientStoreKey("transient_stake")
	keyTreasury := sdk.NewKVStoreKey("treasury")
	keyOracle := sdk.NewKVStoreKey("oracle")
	keyBank := sdk.NewKVStoreKey("bank")
	keyFeeCollection := sdk.NewKVStoreKey("fee")

	pk := params.NewKeeper(mapp.Cdc, keyGlobalParams, tkeyGlobalParams)
	fck := auth.NewFeeCollectionKeeper(mapp.Cdc, keyFeeCollection)
	ck := tax.NewBaseKeeper(keyBank, mapp.Cdc, mapp.AccountKeeper, fck)
	sk := stake.NewKeeper(mapp.Cdc, keyStake, tkeyStake, ck, pk.Subspace(stake.DefaultParamspace), stake.DefaultCodespace)
	tk := treasury.NewKeeper(keyTreasury, mapp.Cdc, ck)
	keeper := NewKeeper(keyOracle, mapp.Cdc, tk, sk.GetValidatorSet(), pk.Subspace(DefaultParamspace))

	mapp.Router().AddRoute("oracle", NewHandler(keeper))

	mapp.SetEndBlocker(getEndBlocker(keeper))
	mapp.SetInitChainer(getInitChainer(mapp, keeper, sk))

	require.NoError(t, mapp.CompleteSetup(keyStake, tkeyStake,
		keyTreasury,
		keyOracle, keyBank, keyFeeCollection,
		keyGlobalParams, tkeyGlobalParams))

	genAccs, addrs, pubKeys, privKeys := mock.CreateGenAccounts(numGenAccs, sdk.Coins{sdk.NewInt64Coin(assets.LunaDenom, 42)})

	mock.SetGenesis(mapp, genAccs)

	return mapp, keeper, sk, addrs, pubKeys, privKeys
}

// oracle and stake endblocker
func getEndBlocker(keeper Keeper) sdk.EndBlocker {
	return func(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
		_, _, _, tags := EndBlocker(ctx, keeper)
		return abci.ResponseEndBlock{
			Tags: tags,
		}
	}
}

// oracle and stake initchainer
func getInitChainer(mapp *mock.App, keeper Keeper, stakeKeeper stake.Keeper) sdk.InitChainer {
	return func(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
		mapp.InitChainer(ctx, req)

		stakeGenesis := stake.GenesisState{
			Pool: stake.InitialPool(),
			Params: stake.Params{
				UnbondingTime: 60 * 60 * 24 * 3 * time.Second,
				MaxValidators: 100,
				BondDenom:     assets.LunaDenom,
			},
		}
		stakeGenesis.Pool.LooseTokens = sdk.NewDec(100000)

		validators, err := stake.InitGenesis(ctx, stakeKeeper, stakeGenesis)
		if err != nil {
			panic(err)
		}

		InitGenesis(ctx, keeper, GenesisState{
			Params: NewParams(
				assets.GetAllDenoms(),
				sdk.NewInt(1),             // one block
				sdk.NewDecWithPrec(66, 2), // 66%
			),
		})

		return abci.ResponseInitChain{
			Validators: validators,
		}
	}
}
