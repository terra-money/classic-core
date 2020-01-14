package auth

import (
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/terra-project/core/x/params"
	"github.com/terra-project/core/x/supply"
)

type testInput struct {
	cdc *codec.Codec
	ctx sdk.Context
	ak  AccountKeeper
	sk  SupplyKeeper
	tk  TreasuryKeeper
}

// moduleAccount defines an account for modules that holds coins on a pool
type moduleAccount struct {
	*BaseAccount
	name        string   // name of the module
	permissions []string // permissions of module account
}

// HasPermission returns whether or not the module account has permission.
func (ma moduleAccount) HasPermission(permission string) bool {
	for _, perm := range ma.permissions {
		if perm == permission {
			return true
		}
	}
	return false
}

// GetName returns the the name of the holder's module
func (ma moduleAccount) GetName() string {
	return ma.name
}

// GetPermissions returns permissions granted to the module account
func (ma moduleAccount) GetPermissions() []string {
	return ma.permissions
}

func setupTestInput() testInput {
	db := dbm.NewMemDB()

	cdc := codec.New()
	RegisterCodec(cdc)
	cdc.RegisterInterface((*supply.ModuleAccountI)(nil), nil)
	cdc.RegisterConcrete(&moduleAccount{}, "core/ModuleAccount", nil)
	codec.RegisterCrypto(cdc)

	authCapKey := sdk.NewKVStoreKey("authCapKey")
	keyParams := sdk.NewKVStoreKey("subspace")
	tkeyParams := sdk.NewTransientStoreKey("transient_subspace")

	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(authCapKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyParams, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(tkeyParams, sdk.StoreTypeTransient, db)
	err := ms.LoadLatestVersion()
	if err != nil {
		panic(err)
	}

	ps := params.NewSubspace(cdc, keyParams, tkeyParams, DefaultParamspace)
	ak := NewAccountKeeper(cdc, authCapKey, ps, ProtoBaseAccount)
	sk := NewDummySupplyKeeper(ak)

	ctx := sdk.NewContext(ms, abci.Header{ChainID: "test-chain-id"}, false, log.NewNopLogger())

	ak.SetParams(ctx, DefaultParams())

	tk := NewDummyTreasuryKeeper()

	return testInput{cdc: cdc, ctx: ctx, ak: ak, sk: sk, tk: tk}
}

// DummyTreasuryKeeper no-lint
type DummyTreasuryKeeper struct{}

// NewDummyTreasuryKeeper no-lint
func NewDummyTreasuryKeeper() DummyTreasuryKeeper { return DummyTreasuryKeeper{} }

// GetTaxRate for the dummy treasury keeper
func (tk DummyTreasuryKeeper) GetTaxRate(_ sdk.Context) (rate sdk.Dec) {
	return sdk.NewDecWithPrec(1, 3) // 0.1%
}

// GetTaxCap for the dummy treasury keeper
func (tk DummyTreasuryKeeper) GetTaxCap(_ sdk.Context, _ string) (taxCap sdk.Int) {
	return sdk.OneInt()
}

// RecordEpochTaxProceeds for the dummy treasury keeper
func (tk DummyTreasuryKeeper) RecordEpochTaxProceeds(_ sdk.Context, _ sdk.Coins) {
	return
}

// DummySupplyKeeper defines a supply keeper used only for testing to avoid
// circle dependencies
type DummySupplyKeeper struct {
	ak AccountKeeper
}

// NewDummySupplyKeeper creates a DummySupplyKeeper instance
func NewDummySupplyKeeper(ak AccountKeeper) DummySupplyKeeper {
	return DummySupplyKeeper{ak}
}

// SendCoinsFromAccountToModule for the dummy supply keeper
func (sk DummySupplyKeeper) SendCoinsFromAccountToModule(ctx sdk.Context, fromAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) sdk.Error {

	fromAcc := sk.ak.GetAccount(ctx, fromAddr)
	moduleAcc := sk.GetModuleAccount(ctx, recipientModule)

	newFromCoins, hasNeg := fromAcc.GetCoins().SafeSub(amt)
	if hasNeg {
		return sdk.ErrInsufficientCoins(fromAcc.GetCoins().String())
	}

	newToCoins := moduleAcc.GetCoins().Add(amt)

	if err := fromAcc.SetCoins(newFromCoins); err != nil {
		return sdk.ErrInternal(err.Error())
	}

	if err := moduleAcc.SetCoins(newToCoins); err != nil {
		return sdk.ErrInternal(err.Error())
	}

	sk.ak.SetAccount(ctx, fromAcc)
	sk.ak.SetAccount(ctx, moduleAcc)

	return nil
}

// GetModuleAccount for dummy supply keeper
func (sk DummySupplyKeeper) GetModuleAccount(ctx sdk.Context, moduleName string) supply.ModuleAccountI {
	addr := sk.GetModuleAddress(moduleName)

	acc := sk.ak.GetAccount(ctx, addr)
	if acc != nil {
		macc, ok := acc.(supply.ModuleAccountI)
		if ok {
			return macc
		}
	}

	moduleAddress := sk.GetModuleAddress(moduleName)
	baseAcc := NewBaseAccountWithAddress(moduleAddress)

	// create a new module account
	macc := &moduleAccount{
		BaseAccount: &baseAcc,
		name:        moduleName,
		permissions: []string{"basic"},
	}

	maccI := (sk.ak.NewAccount(ctx, macc)).(supply.ModuleAccountI)
	sk.ak.SetAccount(ctx, maccI)
	return maccI
}

// GetModuleAddress for dummy supply keeper
func (sk DummySupplyKeeper) GetModuleAddress(moduleName string) sdk.AccAddress {
	return sdk.AccAddress(crypto.AddressHash([]byte(moduleName)))
}
