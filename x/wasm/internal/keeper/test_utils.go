// nolint:deadcode unused noalias
package keeper

import (
	"fmt"
	"testing"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	distr "github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/cosmos-sdk/x/supply"

	"github.com/stretchr/testify/require"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"

	core "github.com/terra-money/core/types"
	bankwasm "github.com/terra-money/core/x/bank/wasm"
	"github.com/terra-money/core/x/market"
	marketwasm "github.com/terra-money/core/x/market/wasm"
	"github.com/terra-money/core/x/oracle"
	oraclewasm "github.com/terra-money/core/x/oracle/wasm"
	stakingwasm "github.com/terra-money/core/x/staking/wasm"
	"github.com/terra-money/core/x/treasury"
	treasurywasm "github.com/terra-money/core/x/treasury/wasm"
	"github.com/terra-money/core/x/wasm/config"
	"github.com/terra-money/core/x/wasm/internal/types"
)

func makeTestCodec() *codec.Codec {
	var cdc = codec.New()

	// Register AppAccount
	auth.RegisterCodec(cdc)
	bank.RegisterCodec(cdc)
	sdk.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	staking.RegisterCodec(cdc)
	supply.RegisterCodec(cdc)
	distr.RegisterCodec(cdc)
	oracle.RegisterCodec(cdc)
	market.RegisterCodec(cdc)

	return cdc
}

// TestInput nolint
type TestInput struct {
	Ctx            sdk.Context
	Cdc            *codec.Codec
	AccKeeper      auth.AccountKeeper
	BankKeeper     bank.Keeper
	SupplyKeeper   supply.Keeper
	StakingKeeper  staking.Keeper
	DistrKeeper    distr.Keeper
	OracleKeeper   oracle.Keeper
	MarketKeeper   market.Keeper
	TreasuryKeeper treasury.Keeper
	WasmKeeper     Keeper
}

// CreateTestInput nolint
func CreateTestInput(t *testing.T) TestInput {
	keyContract := sdk.NewKVStoreKey(types.StoreKey)
	keyAcc := sdk.NewKVStoreKey(auth.StoreKey)
	keyParams := sdk.NewKVStoreKey(params.StoreKey)
	tkeyParams := sdk.NewTransientStoreKey(params.TStoreKey)
	keyStaking := sdk.NewKVStoreKey(staking.StoreKey)
	keyDistr := sdk.NewKVStoreKey(distr.StoreKey)
	keySupply := sdk.NewKVStoreKey(supply.StoreKey)
	keyOracle := sdk.NewKVStoreKey(oracle.StoreKey)
	keyMarket := sdk.NewKVStoreKey(market.StoreKey)
	keyTreasury := sdk.NewKVStoreKey(treasury.StoreKey)

	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(keyContract, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyAcc, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyParams, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(tkeyParams, sdk.StoreTypeTransient, db)
	ms.MountStoreWithDB(keyStaking, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyDistr, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keySupply, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyOracle, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyMarket, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyTreasury, sdk.StoreTypeIAVL, db)

	require.NoError(t, ms.LoadLatestVersion())

	blackListAddrs := map[string]bool{
		auth.FeeCollectorName:     true,
		staking.NotBondedPoolName: true,
		staking.BondedPoolName:    true,
		distr.ModuleName:          true,
		oracle.ModuleName:         true,
		market.ModuleName:         true,
		treasury.ModuleName:       true,
	}

	ctx := sdk.NewContext(ms, abci.Header{}, false, log.NewNopLogger())
	cdc := makeTestCodec()

	paramsKeeper := params.NewKeeper(cdc, keyParams, tkeyParams)

	accountKeeper := auth.NewAccountKeeper(
		cdc,    // amino codec
		keyAcc, // target store
		paramsKeeper.Subspace(auth.DefaultParamspace),
		auth.ProtoBaseAccount, // prototype
	)

	bankKeeper := bank.NewBaseKeeper(
		accountKeeper,
		paramsKeeper.Subspace(bank.DefaultParamspace),
		blackListAddrs,
	)
	bankKeeper.SetSendEnabled(ctx, true)

	maccPerms := map[string][]string{
		auth.FeeCollectorName:     nil,
		staking.NotBondedPoolName: {supply.Burner, supply.Staking},
		staking.BondedPoolName:    {supply.Burner, supply.Staking},
		distr.ModuleName:          nil,
		oracle.ModuleName:         nil,
		market.ModuleName:         {supply.Burner, supply.Minter},
		treasury.ModuleName:       {supply.Minter},
	}

	supplyKeeper := supply.NewKeeper(cdc, keySupply, accountKeeper, bankKeeper, maccPerms)
	totalSupply := sdk.NewCoins(sdk.NewCoin(core.MicroLunaDenom, sdk.NewInt(100000000)), sdk.NewCoin(core.MicroSDRDenom, sdk.NewInt(100000000)))
	supplyKeeper.SetSupply(ctx, supply.NewSupply(totalSupply))

	stakingKeeper := staking.NewKeeper(
		cdc,
		keyStaking,
		supplyKeeper, paramsKeeper.Subspace(staking.DefaultParamspace),
	)

	distrKeeper := distr.NewKeeper(
		cdc,
		keyDistr, paramsKeeper.Subspace(distr.DefaultParamspace),
		stakingKeeper, supplyKeeper, auth.FeeCollectorName, blackListAddrs)

	oracleKeeper := oracle.NewKeeper(
		cdc,
		keyOracle, paramsKeeper.Subspace(oracle.DefaultParamspace),
		distrKeeper, stakingKeeper, supplyKeeper, distr.ModuleName,
	)

	marketKeeper := market.NewKeeper(
		cdc,
		keyMarket, paramsKeeper.Subspace(market.DefaultParamspace),
		oracleKeeper, supplyKeeper,
	)

	treasuryKeeper := treasury.NewKeeper(
		cdc,
		keyTreasury, paramsKeeper.Subspace(treasury.DefaultParamspace),
		supplyKeeper, marketKeeper, stakingKeeper, distrKeeper,
		oracle.ModuleName, distr.ModuleName,
	)

	treasuryKeeper.SetParams(ctx, treasury.DefaultParams())

	distrKeeper.SetFeePool(ctx, distr.InitialFeePool())
	distrParams := distr.DefaultParams()
	distrParams.CommunityTax = sdk.NewDecWithPrec(2, 2)
	distrParams.BaseProposerReward = sdk.NewDecWithPrec(1, 2)
	distrParams.BonusProposerReward = sdk.NewDecWithPrec(4, 2)
	distrKeeper.SetParams(ctx, distrParams)

	feeCollectorAcc := supply.NewEmptyModuleAccount(auth.FeeCollectorName)
	notBondedPool := supply.NewEmptyModuleAccount(staking.NotBondedPoolName, supply.Burner, supply.Staking)
	bondPool := supply.NewEmptyModuleAccount(staking.BondedPoolName, supply.Burner, supply.Staking)
	distrAcc := supply.NewEmptyModuleAccount(distr.ModuleName)
	marketAcc := supply.NewEmptyModuleAccount(types.ModuleName, supply.Burner, supply.Minter)

	// funds for huge withdraw
	distrAcc.SetCoins(sdk.NewCoins(sdk.NewInt64Coin(core.MicroLunaDenom, 500000)))
	notBondedPool.SetCoins(totalSupply)

	supplyKeeper.SetModuleAccount(ctx, feeCollectorAcc)
	supplyKeeper.SetModuleAccount(ctx, bondPool)
	supplyKeeper.SetModuleAccount(ctx, notBondedPool)
	supplyKeeper.SetModuleAccount(ctx, distrAcc)
	supplyKeeper.SetModuleAccount(ctx, marketAcc)

	stakingKeeper.SetHooks(staking.NewMultiStakingHooks(distrKeeper.Hooks()))

	genesis := staking.DefaultGenesisState()
	genesis.Params.BondDenom = core.MicroLunaDenom
	_ = staking.InitGenesis(ctx, stakingKeeper, accountKeeper, supplyKeeper, genesis)

	router := baseapp.NewRouter()

	keeper := NewKeeper(
		cdc,
		keyContract,
		paramsKeeper.Subspace(types.DefaultParamspace),
		accountKeeper,
		bankKeeper,
		supplyKeeper,
		treasuryKeeper,
		router,
		types.DefaultFeatures,
		config.DefaultConfig(),
	)

	bankHandler := bank.NewHandler(bankKeeper)
	stakingHandler := staking.NewHandler(stakingKeeper)
	distrHandler := distr.NewHandler(distrKeeper)
	marketHandler := market.NewHandler(marketKeeper)
	router.AddRoute(bank.RouterKey, bankHandler)
	router.AddRoute(staking.RouterKey, stakingHandler)
	router.AddRoute(distr.RouterKey, distrHandler)
	router.AddRoute(market.RouterKey, marketHandler)
	router.AddRoute(types.RouterKey, TestHandler(keeper))

	marketKeeper.SetParams(ctx, market.DefaultParams())
	oracleKeeper.SetParams(ctx, oracle.DefaultParams())

	keeper.SetParams(ctx, types.DefaultParams())
	keeper.RegisterQueriers(map[string]types.WasmQuerierInterface{
		types.WasmQueryRouteBank:     bankwasm.NewWasmQuerier(bankKeeper),
		types.WasmQueryRouteStaking:  stakingwasm.NewWasmQuerier(stakingKeeper),
		types.WasmQueryRouteMarket:   marketwasm.NewWasmQuerier(marketKeeper),
		types.WasmQueryRouteTreasury: treasurywasm.NewWasmQuerier(treasuryKeeper),
		types.WasmQueryRouteWasm:     NewWasmQuerier(keeper),
		types.WasmQueryRouteOracle:   oraclewasm.NewWasmQuerier(oracleKeeper),
	})
	keeper.RegisterMsgParsers(map[string]types.WasmMsgParserInterface{
		types.WasmMsgParserRouteBank:    bankwasm.NewWasmMsgParser(),
		types.WasmMsgParserRouteStaking: stakingwasm.NewWasmMsgParser(),
		types.WasmMsgParserRouteMarket:  marketwasm.NewWasmMsgParser(),
		types.WasmMsgParserRouteWasm:    NewWasmMsgParser(),
	})

	keeper.SetLastCodeID(ctx, 0)
	keeper.SetLastInstanceID(ctx, 0)

	return TestInput{
		ctx, cdc,
		accountKeeper,
		bankKeeper,
		supplyKeeper,
		stakingKeeper,
		distrKeeper,
		oracleKeeper,
		marketKeeper,
		treasuryKeeper,
		keeper}
}

// InitMsg nolint
type InitMsg struct {
	Verifier    sdk.AccAddress `json:"verifier"`
	Beneficiary sdk.AccAddress `json:"beneficiary"`
}

func createFakeFundedAccount(ctx sdk.Context, am auth.AccountKeeper, coins sdk.Coins) sdk.AccAddress {
	_, _, addr := keyPubAddr()
	baseAcct := auth.NewBaseAccountWithAddress(addr)
	_ = baseAcct.SetCoins(coins)
	am.SetAccount(ctx, &baseAcct)

	return addr
}

func keyPubAddr() (crypto.PrivKey, crypto.PubKey, sdk.AccAddress) {
	key := ed25519.GenPrivKey()
	pub := key.PubKey()
	addr := sdk.AccAddress(pub.Address())
	return key, pub, addr
}

// TestHandler returns a wasm handler for tests (to avoid circular imports)
func TestHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case types.MsgInstantiateContract:

			return handleInstantiate(ctx, k, msg)
		case *types.MsgInstantiateContract:
			return handleInstantiate(ctx, k, *msg)

		case types.MsgExecuteContract:
			return handleExecute(ctx, k, msg)
		case *types.MsgExecuteContract:
			return handleExecute(ctx, k, *msg)

		default:
			errMsg := fmt.Sprintf("unrecognized wasm message type: %T", msg)
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
		}
	}
}

func handleInstantiate(ctx sdk.Context, k Keeper, msg types.MsgInstantiateContract) (*sdk.Result, error) {
	contractAddr, err := k.InstantiateContract(ctx, msg.CodeID, msg.Owner, msg.InitMsg, msg.InitCoins, msg.Migratable)
	if err != nil {
		return nil, err
	}

	return &sdk.Result{
		Data:   contractAddr,
		Events: ctx.EventManager().Events(),
	}, nil
}

func handleExecute(ctx sdk.Context, k Keeper, msg types.MsgExecuteContract) (*sdk.Result, error) {
	res, err := k.ExecuteContract(ctx, msg.Contract, msg.Sender, msg.ExecuteMsg, msg.Coins)
	if err != nil {
		return nil, err
	}

	return &sdk.Result{
		Data:   res,
		Events: ctx.EventManager().Events(),
	}, nil
}
