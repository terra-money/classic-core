package market

// import (
// 	"terra/types/assets"
// 	"terra/x/oracle"
// 	"terra/x/treasury"
// 	"testing"
// 	"time"

// 	sdk "github.com/cosmos/cosmos-sdk/types"
// 	"github.com/cosmos/cosmos-sdk/x/auth"
// 	"github.com/cosmos/cosmos-sdk/x/mock"
// 	"github.com/cosmos/cosmos-sdk/x/params"
// 	"github.com/stretchr/testify/require"
// 	"github.com/tendermint/tendermint/crypto"

// 	abci "github.com/tendermint/tendermint/abci/types"
// )

// // initialize the mock application for this module
// func getMockApp(t *testing.T, numGenAccs int) (*mock.App, Keeper, oracle.Keeper, []sdk.AccAddress, []crypto.PubKey, []crypto.PrivKey) {
// 	mapp := mock.NewApp()

// 	stake.RegisterCodec(mapp.Cdc)

// 	keyGlobalParams := sdk.NewKVStoreKey("params")
// 	tkeyGlobalParams := sdk.NewTransientStoreKey("transient_params")
// 	keyStake := sdk.NewKVStoreKey("stake")
// 	tkeyStake := sdk.NewTransientStoreKey("transient_stake")
// 	keyTreasury := sdk.NewKVStoreKey("treasury")
// 	keyOracle := sdk.NewKVStoreKey("oracle")
// 	keyBank := sdk.NewKVStoreKey("bank")
// 	keyMarket := sdk.NewKVStoreKey("market")
// 	keyFeeCollection := sdk.NewKVStoreKey("fee")

// 	pk := params.NewKeeper(mapp.Cdc, keyGlobalParams, tkeyGlobalParams)
// 	fck := auth.NewFeeCollectionKeeper(mapp.Cdc, keyFeeCollection)
// 	ck := tax.NewBaseKeeper(keyBank, mapp.Cdc, mapp.AccountKeeper, fck)
// 	sk := stake.NewKeeper(mapp.Cdc, keyStake, tkeyStake, ck, pk.Subspace(stake.DefaultParamspace), stake.DefaultCodespace)
// 	tk := treasury.NewKeeper(keyTreasury, mapp.Cdc, ck)
// 	ok := oracle.NewKeeper(keyOracle, mapp.Cdc, tk, sk.GetValidatorSet(), pk.Subspace(oracle.DefaultParamspace))
// 	keeper := NewKeeper(ok, tk, ck)

// 	mapp.Router().AddRoute("market", NewHandler(keeper))

// 	genAccs, addrs, pubKeys, privKeys := mock.CreateGenAccounts(numGenAccs, sdk.Coins{sdk.NewInt64Coin(assets.LunaDenom, 42)})

// 	mapp.SetEndBlocker(getEndBlocker(keeper))
// 	mapp.SetInitChainer(getInitChainer(mapp, keeper, sk, ok, addrs, pubKeys))

// 	require.NoError(t, mapp.CompleteSetup(keyStake, tkeyStake,
// 		keyTreasury,
// 		keyOracle, keyBank, keyMarket, keyFeeCollection,
// 		keyGlobalParams, tkeyGlobalParams))

// 	mock.SetGenesis(mapp, genAccs)

// 	return mapp, keeper, ok, addrs, pubKeys, privKeys
// }

// // oracle and stake endblocker
// func getEndBlocker(keeper Keeper) sdk.EndBlocker {
// 	return func(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
// 		return abci.ResponseEndBlock{}
// 	}
// }

// // oracle and stake initchainer
// func getInitChainer(mapp *mock.App, keeper Keeper, sk stake.Keeper, ok oracle.Keeper, addrs []sdk.AccAddress, pubKeys []crypto.PubKey) sdk.InitChainer {
// 	return func(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
// 		mapp.InitChainer(ctx, req)

// 		data := oracle.DefaultGenesisState()
// 		oracle.InitGenesis(ctx, ok, data)

// 		stakeGenesis := stake.GenesisState{
// 			Pool: stake.InitialPool(),
// 			Params: stake.Params{
// 				UnbondingTime: 60 * 60 * 24 * 3 * time.Second,
// 				MaxValidators: 100,
// 				BondDenom:     assets.LunaDenom,
// 			},
// 		}
// 		stakeGenesis.Pool.LooseTokens = sdk.NewDec(100000)

// 		validators, err := stake.InitGenesis(ctx, sk, stakeGenesis)
// 		if err != nil {
// 			panic(err)
// 		}

// 		trashDescription := stake.NewDescription("", "", "", "")
// 		trashCommission := stake.NewCommissionMsg(sdk.ZeroDec(), sdk.ZeroDec(), sdk.ZeroDec())

// 		stakeHandler := stake.NewHandler(sk)
// 		valCreateMsg := stake.NewMsgCreateValidator(
// 			sdk.ValAddress(addrs[0].Bytes()), pubKeys[0],
// 			sdk.NewInt64Coin(assets.LunaDenom, 10), trashDescription, trashCommission,
// 		)

// 		stakeHandler(ctx, valCreateMsg)

// 		_ = sk.ApplyAndReturnValidatorSetUpdates(ctx)

// 		return abci.ResponseInitChain{
// 			Validators: validators,
// 		}
// 	}
// }
