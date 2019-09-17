// nolint:deadcode unused DONTCOVER
package oracle

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto/ed25519"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	distr "github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/mock"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/cosmos-sdk/x/supply"

	core "github.com/terra-project/core/types"
	"github.com/terra-project/core/x/oracle/internal/keeper"
	"github.com/terra-project/core/x/oracle/internal/types"
)

var (
	valTokens  = sdk.TokensFromConsensusPower(42)
	initTokens = sdk.TokensFromConsensusPower(100000)
	valCoins   = sdk.NewCoins(sdk.NewCoin(core.MicroLunaDenom, valTokens))
	initCoins  = sdk.NewCoins(sdk.NewCoin(core.MicroLunaDenom, initTokens))
)

type testInput struct {
	mApp          *mock.App
	oracleKeeper  Keeper
	stakingKeeper staking.Keeper
	supplyKeeper  supply.Keeper
	distrKeeper   distr.Keeper
	addrs         []sdk.AccAddress
	pubKeys       []crypto.PubKey
	privKeys      []crypto.PrivKey
}

func getMockApp(t *testing.T, numGenAccs int, genState GenesisState, genAccs []auth.Account) testInput {
	mApp := mock.NewApp()

	staking.RegisterCodec(mApp.Cdc)
	RegisterCodec(mApp.Cdc)
	supply.RegisterCodec(mApp.Cdc)

	keyStaking := sdk.NewKVStoreKey(staking.StoreKey)
	tKeyStaking := sdk.NewTransientStoreKey(staking.TStoreKey)
	keyDistr := sdk.NewKVStoreKey(distr.StoreKey)
	keyOracle := sdk.NewKVStoreKey(StoreKey)
	keySupply := sdk.NewKVStoreKey(supply.StoreKey)

	blackListAddrs := map[string]bool{
		auth.FeeCollectorName:     true,
		staking.NotBondedPoolName: true,
		staking.BondedPoolName:    true,
		distr.ModuleName:          true,
		types.ModuleName:          true,
	}

	pk := mApp.ParamsKeeper

	bk := bank.NewBaseKeeper(mApp.AccountKeeper, mApp.ParamsKeeper.Subspace(bank.DefaultParamspace), bank.DefaultCodespace, blackListAddrs)

	maccPerms := map[string][]string{
		distr.ModuleName:          nil,
		staking.NotBondedPoolName: {supply.Burner, supply.Staking},
		staking.BondedPoolName:    {supply.Burner, supply.Staking},
		types.ModuleName:          nil,
	}

	supplyKeeper := supply.NewKeeper(mApp.Cdc, keySupply, mApp.AccountKeeper, bk, maccPerms)
	stakingKeeper := staking.NewKeeper(mApp.Cdc, keyStaking, tKeyStaking, supplyKeeper, pk.Subspace(staking.DefaultParamspace), staking.DefaultCodespace)
	distrKeeper := distr.NewKeeper(mApp.Cdc, keyDistr, pk.Subspace(distr.DefaultParamspace), stakingKeeper, supplyKeeper, distr.DefaultCodespace, auth.FeeCollectorName, blackListAddrs)

	keeper := NewKeeper(mApp.Cdc, keyOracle, pk.Subspace(DefaultParamspace), distrKeeper, stakingKeeper, supplyKeeper, distr.ModuleName, DefaultCodespace)

	mApp.Router().AddRoute(RouterKey, NewHandler(keeper))
	mApp.QueryRouter().AddRoute(QuerierRoute, NewQuerier(keeper))

	mApp.SetBeginBlocker(getBeginBloker(distrKeeper))
	mApp.SetEndBlocker(getEndBlocker(keeper, stakingKeeper))
	mApp.SetInitChainer(getInitChainer(mApp, keeper, stakingKeeper, supplyKeeper, genAccs, genState))

	require.NoError(t, mApp.CompleteSetup(keyStaking, tKeyStaking, keyOracle, keySupply, keyDistr))

	var (
		addrs    []sdk.AccAddress
		pubKeys  []crypto.PubKey
		privKeys []crypto.PrivKey
	)

	if genAccs == nil || len(genAccs) == 0 {
		genAccs, addrs, pubKeys, privKeys = mock.CreateGenAccounts(numGenAccs, valCoins)
	}

	mock.SetGenesis(mApp, genAccs)

	return testInput{mApp, keeper, stakingKeeper, supplyKeeper, distrKeeper, addrs, pubKeys, privKeys}
}

func getBeginBloker(distrKeeper distr.Keeper) sdk.BeginBlocker {
	return func(ctx sdk.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock {
		distr.BeginBlocker(ctx, req, distrKeeper)

		return abci.ResponseBeginBlock{
			Events: ctx.EventManager().ABCIEvents(),
		}
	}
}

// oracle and staking endblocker
func getEndBlocker(keeper Keeper, sk staking.Keeper) sdk.EndBlocker {
	return func(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
		EndBlocker(ctx, keeper)
		staking.EndBlocker(ctx, sk)
		return abci.ResponseEndBlock{}
	}
}

// gov and staking initchainer
func getInitChainer(mapp *mock.App, keeper Keeper, stakingKeeper staking.Keeper, supplyKeeper supply.Keeper, accs []auth.Account, genState GenesisState) sdk.InitChainer {
	return func(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
		mapp.InitChainer(ctx, req)

		stakingGenesis := staking.DefaultGenesisState()

		totalSupply := sdk.NewCoins(sdk.NewCoin(core.MicroLunaDenom, initTokens.MulRaw(int64(len(mapp.GenesisAccounts)))))
		supplyKeeper.SetSupply(ctx, supply.NewSupply(totalSupply))

		// set module accounts
		govAcc := supply.NewEmptyModuleAccount(types.ModuleName, supply.Burner)
		notBondedPool := supply.NewEmptyModuleAccount(staking.NotBondedPoolName, supply.Burner, supply.Staking)
		bondPool := supply.NewEmptyModuleAccount(staking.BondedPoolName, supply.Burner, supply.Staking)

		supplyKeeper.SetModuleAccount(ctx, govAcc)
		supplyKeeper.SetModuleAccount(ctx, notBondedPool)
		supplyKeeper.SetModuleAccount(ctx, bondPool)

		validators := staking.InitGenesis(ctx, stakingKeeper, mapp.AccountKeeper, supplyKeeper, stakingGenesis)
		if genState.IsEmpty() {
			InitGenesis(ctx, keeper, DefaultGenesisState())
		} else {
			InitGenesis(ctx, keeper, genState)
		}

		return abci.ResponseInitChain{
			Validators: validators,
		}
	}
}

var (
	pubkeys = []crypto.PubKey{
		ed25519.GenPrivKey().PubKey(),
		ed25519.GenPrivKey().PubKey(),
		ed25519.GenPrivKey().PubKey(),
	}

	testDescription     = staking.NewDescription("T", "E", "S", "T")
	testCommissionRates = staking.NewCommissionRates(sdk.ZeroDec(), sdk.ZeroDec(), sdk.ZeroDec())
)

func createValidators(t *testing.T, stakingHandler sdk.Handler, ctx sdk.Context, addrs []sdk.ValAddress, powerAmt []int64) {
	require.True(t, len(addrs) <= len(pubkeys), "Not enough pubkeys specified at top of file.")

	for i := 0; i < len(addrs); i++ {

		valTokens := sdk.TokensFromConsensusPower(powerAmt[i])
		valCreateMsg := staking.NewMsgCreateValidator(
			addrs[i], pubkeys[i], sdk.NewCoin(core.MicroLunaDenom, valTokens),
			testDescription, testCommissionRates, sdk.OneInt(),
		)

		res := stakingHandler(ctx, valCreateMsg)
		require.True(t, res.IsOK())
	}
}

var (
	uSDRAmt    = sdk.NewInt(1005 * core.MicroUnit)
	stakingAmt = sdk.TokensFromConsensusPower(10)

	randomPrice        = sdk.NewDec(1700)
	anotherRandomPrice = sdk.NewDecWithPrec(4882, 2) // swap rate
)

func setup(t *testing.T) (keeper.TestInput, sdk.Handler) {
	input := keeper.CreateTestInput(t)
	params := input.OracleKeeper.GetParams(input.Ctx)
	params.VotePeriod = 1
	input.OracleKeeper.SetParams(input.Ctx, params)
	h := NewHandler(input.OracleKeeper)

	sh := staking.NewHandler(input.StakingKeeper)

	// Validator created
	got := sh(input.Ctx, keeper.NewTestMsgCreateValidator(keeper.ValAddrs[0], keeper.PubKeys[0], stakingAmt))
	require.True(t, got.IsOK())
	got = sh(input.Ctx, keeper.NewTestMsgCreateValidator(keeper.ValAddrs[1], keeper.PubKeys[1], stakingAmt))
	require.True(t, got.IsOK())
	got = sh(input.Ctx, keeper.NewTestMsgCreateValidator(keeper.ValAddrs[2], keeper.PubKeys[2], stakingAmt))
	require.True(t, got.IsOK())
	staking.EndBlocker(input.Ctx, input.StakingKeeper)

	return input, h
}
