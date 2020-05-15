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
	"github.com/cosmos/cosmos-sdk/x/params"

	"github.com/stretchr/testify/require"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"

	bankwasm "github.com/terra-project/core/x/bank/wasm"
	"github.com/terra-project/core/x/wasm/internal/types"
)

func makeTestCodec() *codec.Codec {
	var cdc = codec.New()

	// Register AppAccount
	auth.RegisterCodec(cdc)
	bank.RegisterCodec(cdc)
	sdk.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)

	return cdc
}

// CreateTestInput nolint
func CreateTestInput(t *testing.T) (sdk.Context, auth.AccountKeeper, Keeper) {
	keyContract := sdk.NewKVStoreKey(types.StoreKey)
	keyAcc := sdk.NewKVStoreKey(auth.StoreKey)
	keyParams := sdk.NewKVStoreKey(params.StoreKey)
	tkeyParams := sdk.NewTransientStoreKey(params.TStoreKey)

	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(keyContract, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyAcc, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyParams, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(tkeyParams, sdk.StoreTypeTransient, db)
	err := ms.LoadLatestVersion()
	require.Nil(t, err)

	ctx := sdk.NewContext(ms, abci.Header{}, false, log.NewNopLogger())
	cdc := makeTestCodec()

	pk := params.NewKeeper(cdc, keyParams, tkeyParams)

	accountKeeper := auth.NewAccountKeeper(
		cdc,    // amino codec
		keyAcc, // target store
		pk.Subspace(auth.DefaultParamspace),
		auth.ProtoBaseAccount, // prototype
	)

	bk := bank.NewBaseKeeper(
		accountKeeper,
		pk.Subspace(bank.DefaultParamspace),
		nil,
	)
	bk.SetSendEnabled(ctx, true)

	router := baseapp.NewRouter()

	keeper := NewKeeper(
		cdc,
		keyContract,
		pk.Subspace(types.DefaultParamspace),
		accountKeeper,
		bk,
		router,
		types.FeatureStaking,
		types.DefaultWasmConfig(),
	)

	h := bank.NewHandler(bk)
	router.AddRoute(bank.RouterKey, h)
	router.AddRoute(types.RouterKey, TestHandler(keeper))

	keeper.SetParams(ctx, types.DefaultParams())
	keeper.RegisterQueriers(map[string]types.WasmQuerierInterface{
		types.WasmQueryRouteBank: bankwasm.NewWasmQuerier(bk),
		types.WasmQueryRouteWasm: NewWasmQuerier(keeper),
	})
	keeper.RegisterMsgParsers(map[string]types.WasmMsgParserInterface{
		types.WasmMsgParserRouteBank: bankwasm.NewWasmMsgParser(),
		types.WasmMsgParserRouteWasm: NewWasmMsgParser(),
	})

	return ctx, accountKeeper, keeper
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
	contractAddr, err := k.InstantiateContract(ctx, msg.CodeID, msg.Sender, msg.InitMsg, msg.InitCoins)
	if err != nil {
		return nil, err
	}

	return &sdk.Result{
		Data:   contractAddr,
		Events: ctx.EventManager().Events(),
	}, nil
}

func handleExecute(ctx sdk.Context, k Keeper, msg types.MsgExecuteContract) (*sdk.Result, error) {
	res, err := k.ExecuteContract(ctx, msg.Contract, msg.Sender, msg.Msg, msg.Coins)
	if err != nil {
		return nil, err
	}

	res.Events = ctx.EventManager().Events()
	return &res, nil
}
