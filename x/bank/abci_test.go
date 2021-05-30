// nolint:deadcode unused noalias
package bank

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	core "github.com/terra-money/core/types"
	"github.com/terra-money/core/x/auth"
	"github.com/terra-money/core/x/bank/internal/types"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/supply"
)

var (
	pubKeys = []crypto.PubKey{
		secp256k1.GenPrivKey().PubKey(),
		secp256k1.GenPrivKey().PubKey(),
		secp256k1.GenPrivKey().PubKey(),
	}

	addrs = []sdk.AccAddress{
		sdk.AccAddress(pubKeys[0].Address()),
		sdk.AccAddress(pubKeys[1].Address()),
		sdk.AccAddress(pubKeys[2].Address()),
	}

	initTokens = sdk.TokensFromConsensusPower(200)
	initCoins  = sdk.NewCoins(sdk.NewCoin(core.MicroLunaDenom, initTokens))
)

// TestInput nolint
type TestInput struct {
	Ctx           sdk.Context
	Cdc           *codec.Codec
	AccountKeeper auth.AccountKeeper
	BankKeeper    bank.Keeper
	SupplyKeeper  supply.Keeper
}

func newTestCodec() *codec.Codec {
	cdc := codec.New()

	auth.RegisterCodec(cdc)
	sdk.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	supply.RegisterCodec(cdc)
	params.RegisterCodec(cdc)
	bank.RegisterCodec(cdc)

	return cdc
}

// CreateTestInput nolint
func createTestInput(t *testing.T) TestInput {
	keyAcc := sdk.NewKVStoreKey(auth.StoreKey)
	keyParams := sdk.NewKVStoreKey(params.StoreKey)
	tKeyParams := sdk.NewTransientStoreKey(params.TStoreKey)
	keySupply := sdk.NewKVStoreKey(supply.StoreKey)

	cdc := newTestCodec()
	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)
	ctx := sdk.NewContext(ms, abci.Header{Time: time.Now().UTC()}, false, log.NewNopLogger())

	ms.MountStoreWithDB(keyAcc, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(tKeyParams, sdk.StoreTypeTransient, db)
	ms.MountStoreWithDB(keyParams, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keySupply, sdk.StoreTypeIAVL, db)

	require.NoError(t, ms.LoadLatestVersion())

	blackListAddrs := map[string]bool{
		auth.FeeCollectorName: true,
	}

	paramsKeeper := params.NewKeeper(cdc, keyParams, tKeyParams)
	accountKeeper := auth.NewAccountKeeper(cdc, keyAcc, paramsKeeper.Subspace(auth.DefaultParamspace), auth.ProtoBaseAccount)
	bankKeeper := bank.NewBaseKeeper(accountKeeper, paramsKeeper.Subspace(bank.DefaultParamspace), blackListAddrs)
	bankKeeper.SetSendEnabled(ctx, true)

	maccPerms := map[string][]string{
		auth.FeeCollectorName: nil,
		types.BurnModuleName:  {supply.Burner},
	}

	supplyKeeper := supply.NewKeeper(cdc, keySupply, accountKeeper, bankKeeper, maccPerms)
	totalSupply := sdk.NewCoins(sdk.NewCoin(core.MicroLunaDenom, initTokens.MulRaw(int64(len(addrs)))))
	supplyKeeper.SetSupply(ctx, supply.NewSupply(totalSupply))

	feeCollectorAcc := supply.NewEmptyModuleAccount(auth.FeeCollectorName)
	burnAcc := supply.NewEmptyModuleAccount(types.BurnModuleName, supply.Burner)
	supplyKeeper.SetModuleAccount(ctx, feeCollectorAcc)
	supplyKeeper.SetModuleAccount(ctx, burnAcc)

	for _, addr := range addrs {
		_, err := bankKeeper.AddCoins(ctx, sdk.AccAddress(addr), initCoins)
		require.NoError(t, err)
	}

	supply := supplyKeeper.GetSupply(ctx)
	supply = supply.SetTotal(sdk.NewCoins(sdk.NewCoin(core.MicroLunaDenom, initTokens.MulRaw(int64(len(addrs))))))
	supplyKeeper.SetSupply(ctx, supply)

	return TestInput{ctx, cdc, accountKeeper, bankKeeper, supplyKeeper}
}

func TestBurnAddress(t *testing.T) {
	input := createTestInput(t)

	h := bank.NewHandler(input.BankKeeper)

	burnAddress := input.SupplyKeeper.GetModuleAddress(types.BurnModuleName)
	msg := bank.NewMsgSend(addrs[0], burnAddress, initCoins)

	_, err := h(input.Ctx, msg)
	require.NoError(t, err)
	require.Equal(t, initCoins, input.AccountKeeper.GetAccount(input.Ctx, burnAddress).GetCoins())

	EndBlocker(input.Ctx, input.BankKeeper, input.SupplyKeeper)
	require.True(t, input.AccountKeeper.GetAccount(input.Ctx, burnAddress).GetCoins().IsZero())
}
