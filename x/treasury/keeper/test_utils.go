//nolint
package keeper

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	customauth "github.com/terra-project/core/custom/auth"
	custombank "github.com/terra-project/core/custom/bank"
	customdistr "github.com/terra-project/core/custom/distribution"
	customparams "github.com/terra-project/core/custom/params"
	customstaking "github.com/terra-project/core/custom/staking"
	core "github.com/terra-project/core/types"
	"github.com/terra-project/core/x/market"
	marketkeeper "github.com/terra-project/core/x/market/keeper"
	markettypes "github.com/terra-project/core/x/market/types"
	"github.com/terra-project/core/x/oracle"
	oraclekeeper "github.com/terra-project/core/x/oracle/keeper"
	oracletypes "github.com/terra-project/core/x/oracle/types"
	"github.com/terra-project/core/x/treasury/types"

	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"github.com/tendermint/tendermint/libs/log"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	dbm "github.com/tendermint/tm-db"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/simapp"
	simparams "github.com/cosmos/cosmos-sdk/simapp/params"
	"github.com/cosmos/cosmos-sdk/std"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	"github.com/cosmos/cosmos-sdk/x/auth/tx"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

var ModuleBasics = module.NewBasicManager(
	customauth.AppModuleBasic{},
	custombank.AppModuleBasic{},
	customdistr.AppModuleBasic{},
	customstaking.AppModuleBasic{},
	customparams.AppModuleBasic{},
	oracle.AppModuleBasic{},
	market.AppModuleBasic{},
)

// MakeTestCodec nolint
func MakeTestCodec(t *testing.T) codec.Marshaler {
	return MakeEncodingConfig(t).Marshaler
}

// MakeEncodingConfig nolint
func MakeEncodingConfig(_ *testing.T) simparams.EncodingConfig {
	amino := codec.NewLegacyAmino()
	interfaceRegistry := codectypes.NewInterfaceRegistry()
	marshaler := codec.NewProtoCodec(interfaceRegistry)
	txCfg := tx.NewTxConfig(marshaler, tx.DefaultSignModes)

	std.RegisterInterfaces(interfaceRegistry)
	std.RegisterLegacyAminoCodec(amino)

	ModuleBasics.RegisterLegacyAminoCodec(amino)
	ModuleBasics.RegisterInterfaces(interfaceRegistry)
	types.RegisterLegacyAminoCodec(amino)
	types.RegisterInterfaces(interfaceRegistry)

	return simparams.EncodingConfig{
		InterfaceRegistry: interfaceRegistry,
		Marshaler:         marshaler,
		TxConfig:          txCfg,
		Amino:             amino,
	}
}

var (
	ValPubKeys = simapp.CreateTestPubKeys(5)

	PubKeys = []crypto.PubKey{
		secp256k1.GenPrivKey().PubKey(),
		secp256k1.GenPrivKey().PubKey(),
		secp256k1.GenPrivKey().PubKey(),
	}

	Addrs = []sdk.AccAddress{
		sdk.AccAddress(PubKeys[0].Address()),
		sdk.AccAddress(PubKeys[1].Address()),
		sdk.AccAddress(PubKeys[2].Address()),
	}

	ValAddrs = []sdk.ValAddress{
		sdk.ValAddress(PubKeys[0].Address()),
		sdk.ValAddress(PubKeys[1].Address()),
		sdk.ValAddress(PubKeys[2].Address()),
	}

	InitTokens = sdk.TokensFromConsensusPower(200)
	InitCoins  = sdk.NewCoins(sdk.NewCoin(core.MicroLunaDenom, InitTokens))
)

// TestInput nolint
type TestInput struct {
	Ctx            sdk.Context
	Cdc            *codec.LegacyAmino
	TreasuryKeeper Keeper
	StakingKeeper  stakingkeeper.Keeper
	OracleKeeper   types.OracleKeeper
	AccountKeeper  types.AccountKeeper
	BankKeeper     types.BankKeeper
	MarketKeeper   types.MarketKeeper
	DistrKeeper    types.DistributionKeeper
}

// CreateTestInput nolint
func CreateTestInput(t *testing.T) TestInput {
	keyAcc := sdk.NewKVStoreKey(authtypes.StoreKey)
	keyBank := sdk.NewKVStoreKey(banktypes.StoreKey)
	keyParams := sdk.NewKVStoreKey(paramstypes.StoreKey)
	tKeyParams := sdk.NewTransientStoreKey(paramstypes.TStoreKey)
	keyOracle := sdk.NewKVStoreKey(oracletypes.StoreKey)
	keyStaking := sdk.NewKVStoreKey(stakingtypes.StoreKey)
	keyDistr := sdk.NewKVStoreKey(distrtypes.StoreKey)
	keyMarket := sdk.NewKVStoreKey(markettypes.StoreKey)
	keyTreasury := sdk.NewKVStoreKey(types.StoreKey)

	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)
	ctx := sdk.NewContext(ms, tmproto.Header{Time: time.Now().UTC()}, false, log.NewNopLogger())
	encodingConfig := MakeEncodingConfig(t)
	appCodec, legacyAmino := encodingConfig.Marshaler, encodingConfig.Amino

	ms.MountStoreWithDB(keyAcc, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyBank, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(tKeyParams, sdk.StoreTypeTransient, db)
	ms.MountStoreWithDB(keyParams, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyOracle, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyStaking, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyDistr, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyMarket, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyTreasury, sdk.StoreTypeIAVL, db)

	require.NoError(t, ms.LoadLatestVersion())

	blackListAddrs := map[string]bool{
		authtypes.FeeCollectorName:     true,
		stakingtypes.NotBondedPoolName: true,
		stakingtypes.BondedPoolName:    true,
		distrtypes.ModuleName:          true,
		oracletypes.ModuleName:         true,
	}

	maccPerms := map[string][]string{
		authtypes.FeeCollectorName:     nil,
		stakingtypes.NotBondedPoolName: {authtypes.Burner, authtypes.Staking},
		stakingtypes.BondedPoolName:    {authtypes.Burner, authtypes.Staking},
		markettypes.ModuleName:         {authtypes.Burner, authtypes.Minter},
		distrtypes.ModuleName:          nil,
		oracletypes.ModuleName:         nil,
		types.ModuleName:               {authtypes.Burner, authtypes.Minter},
		types.BurnModuleName:           {authtypes.Burner},
	}

	paramsKeeper := paramskeeper.NewKeeper(appCodec, legacyAmino, keyParams, tKeyParams)
	accountKeeper := authkeeper.NewAccountKeeper(appCodec, keyAcc, paramsKeeper.Subspace(authtypes.ModuleName), authtypes.ProtoBaseAccount, maccPerms)
	bankKeeper := bankkeeper.NewBaseKeeper(appCodec, keyBank, accountKeeper, paramsKeeper.Subspace(banktypes.ModuleName), blackListAddrs)

	totalSupply := sdk.NewCoins(sdk.NewCoin(core.MicroLunaDenom, InitTokens.MulRaw(int64(len(Addrs)))))
	bankKeeper.SetSupply(ctx, banktypes.NewSupply(totalSupply))

	stakingKeeper := stakingkeeper.NewKeeper(
		appCodec,
		keyStaking,
		accountKeeper,
		bankKeeper,
		paramsKeeper.Subspace(stakingtypes.ModuleName),
	)

	stakingParams := stakingtypes.DefaultParams()
	stakingParams.BondDenom = core.MicroLunaDenom
	stakingKeeper.SetParams(ctx, stakingParams)

	distrKeeper := distrkeeper.NewKeeper(
		appCodec,
		keyDistr, paramsKeeper.Subspace(distrtypes.ModuleName),
		accountKeeper, bankKeeper, stakingKeeper,
		authtypes.FeeCollectorName, blackListAddrs)

	distrKeeper.SetFeePool(ctx, distrtypes.InitialFeePool())
	distrParams := distrtypes.DefaultParams()
	distrParams.CommunityTax = sdk.NewDecWithPrec(2, 2)
	distrParams.BaseProposerReward = sdk.NewDecWithPrec(1, 2)
	distrParams.BonusProposerReward = sdk.NewDecWithPrec(4, 2)
	distrKeeper.SetParams(ctx, distrParams)
	stakingKeeper.SetHooks(stakingtypes.NewMultiStakingHooks(distrKeeper.Hooks()))

	feeCollectorAcc := authtypes.NewEmptyModuleAccount(authtypes.FeeCollectorName)
	notBondedPool := authtypes.NewEmptyModuleAccount(stakingtypes.NotBondedPoolName, authtypes.Burner, authtypes.Staking)
	bondPool := authtypes.NewEmptyModuleAccount(stakingtypes.BondedPoolName, authtypes.Burner, authtypes.Staking)
	distrAcc := authtypes.NewEmptyModuleAccount(distrtypes.ModuleName)
	oracleAcc := authtypes.NewEmptyModuleAccount(oracletypes.ModuleName)
	marketAcc := authtypes.NewEmptyModuleAccount(markettypes.ModuleName, authtypes.Burner, authtypes.Minter)
	treasuryAcc := authtypes.NewEmptyModuleAccount(types.ModuleName, authtypes.Burner, authtypes.Minter)
	burnAcc := authtypes.NewEmptyModuleAccount(types.BurnModuleName, authtypes.Burner)

	// + 1 for burn account
	bankKeeper.SetBalances(ctx, notBondedPool.GetAddress(), sdk.NewCoins(sdk.NewCoin(core.MicroLunaDenom, InitTokens.MulRaw(int64(len(Addrs)+1)))))

	accountKeeper.SetModuleAccount(ctx, feeCollectorAcc)
	accountKeeper.SetModuleAccount(ctx, bondPool)
	accountKeeper.SetModuleAccount(ctx, notBondedPool)
	accountKeeper.SetModuleAccount(ctx, distrAcc)
	accountKeeper.SetModuleAccount(ctx, oracleAcc)
	accountKeeper.SetModuleAccount(ctx, marketAcc)
	accountKeeper.SetModuleAccount(ctx, treasuryAcc)
	accountKeeper.SetModuleAccount(ctx, burnAcc)

	for _, addr := range Addrs {
		accountKeeper.SetAccount(ctx, authtypes.NewBaseAccountWithAddress(addr))
		err := bankKeeper.AddCoins(ctx, addr, InitCoins)
		require.NoError(t, err)
	}

	// to test burn module account
	err := bankKeeper.AddCoins(ctx, burnAcc.GetAddress(), InitCoins)
	require.NoError(t, err)

	oracleKeeper := oraclekeeper.NewKeeper(
		appCodec,
		keyOracle,
		paramsKeeper.Subspace(oracletypes.ModuleName),
		accountKeeper,
		bankKeeper,
		distrKeeper,
		stakingKeeper,
		distrtypes.ModuleName,
	)
	oracleDefaultParams := oracletypes.DefaultParams()
	oracleKeeper.SetParams(ctx, oracleDefaultParams)

	for _, denom := range oracleDefaultParams.Whitelist {
		oracleKeeper.SetTobinTax(ctx, denom.Name, denom.TobinTax)
	}

	marketKeeper := marketkeeper.NewKeeper(
		appCodec,
		keyMarket, paramsKeeper.Subspace(markettypes.ModuleName),
		accountKeeper,
		bankKeeper,
		oracleKeeper,
	)
	marketKeeper.SetParams(ctx, markettypes.DefaultParams())

	treasuryKeeper := NewKeeper(
		appCodec,
		keyTreasury, paramsKeeper.Subspace(types.ModuleName),
		accountKeeper,
		bankKeeper,
		marketKeeper,
		stakingKeeper,
		distrKeeper,
		distrtypes.ModuleName,
	)

	treasuryKeeper.SetParams(ctx, types.DefaultParams())

	return TestInput{ctx, legacyAmino, treasuryKeeper, stakingKeeper, oracleKeeper, accountKeeper, bankKeeper, marketKeeper, distrKeeper}
}

// NewTestMsgCreateValidator test msg creator
func NewTestMsgCreateValidator(address sdk.ValAddress, pubKey cryptotypes.PubKey, amt sdk.Int) *stakingtypes.MsgCreateValidator {
	commission := stakingtypes.NewCommissionRates(sdk.ZeroDec(), sdk.ZeroDec(), sdk.ZeroDec())
	msg, _ := stakingtypes.NewMsgCreateValidator(
		address, pubKey, sdk.NewCoin(core.MicroLunaDenom, amt),
		stakingtypes.Description{}, commission, sdk.OneInt(),
	)

	return msg
}

func setupValidators(t *testing.T) (TestInput, sdk.Handler) {
	input := CreateTestInput(t)
	sh := staking.NewHandler(input.StakingKeeper)

	// Create Validators
	amt := sdk.TokensFromConsensusPower(100)
	addr, val := ValAddrs[0], ValPubKeys[0]
	addr1, val1 := ValAddrs[1], ValPubKeys[1]
	_, err := sh(input.Ctx, NewTestMsgCreateValidator(addr, val, amt))

	require.NoError(t, err)
	_, err = sh(input.Ctx, NewTestMsgCreateValidator(addr1, val1, amt))
	require.NoError(t, err)

	staking.EndBlocker(input.Ctx, input.StakingKeeper)

	return input, sh
}
