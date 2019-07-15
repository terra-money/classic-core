package auth

import (
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/params"

	"github.com/terra-project/core/types/assets"
)

type testInput struct {
	cdc *codec.Codec
	ctx sdk.Context
	ak  auth.AccountKeeper
	fck auth.FeeCollectionKeeper
}

func setupTestInput() testInput {
	db := dbm.NewMemDB()

	cdc := codec.New()
	auth.RegisterBaseAccount(cdc)

	authCapKey := sdk.NewKVStoreKey("authCapKey")
	fckCapKey := sdk.NewKVStoreKey("fckCapKey")
	keyParams := sdk.NewKVStoreKey("params")
	tkeyParams := sdk.NewTransientStoreKey("transient_params")

	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(authCapKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(fckCapKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyParams, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(tkeyParams, sdk.StoreTypeTransient, db)
	err := ms.LoadLatestVersion()
	if err != nil {
		panic(err)
	}

	pk := params.NewKeeper(cdc, keyParams, tkeyParams)
	ak := auth.NewAccountKeeper(cdc, authCapKey, pk.Subspace(auth.DefaultParamspace), auth.ProtoBaseAccount)
	fck := auth.NewFeeCollectionKeeper(cdc, fckCapKey)
	ctx := sdk.NewContext(ms, abci.Header{ChainID: "test-chain-id"}, false, log.NewNopLogger())

	ak.SetParams(ctx, auth.DefaultParams())

	return testInput{cdc: cdc, ctx: ctx, ak: ak, fck: fck}
}

func newTestMsg(addrs ...sdk.AccAddress) *sdk.TestMsg {
	return sdk.NewTestMsg(addrs...)
}

func newStdFee() auth.StdFee {
	return auth.NewStdFee(50000,
		sdk.NewCoins(sdk.NewInt64Coin(assets.MicroSDRDenom, 150)),
	)
}

// coins to more than cover the fee
func newCoins() sdk.Coins {
	return sdk.Coins{
		sdk.NewInt64Coin(assets.MicroSDRDenom, 10000000),
	}
}

func keyPubAddr() (crypto.PrivKey, crypto.PubKey, sdk.AccAddress) {
	key := secp256k1.GenPrivKey()
	pub := key.PubKey()
	addr := sdk.AccAddress(pub.Address())
	return key, pub, addr
}

func newTestTx(ctx sdk.Context, msgs []sdk.Msg, privs []crypto.PrivKey, accNums []uint64, seqs []uint64, fee auth.StdFee) sdk.Tx {
	sigs := make([]auth.StdSignature, len(privs))
	for i, priv := range privs {
		signBytes := auth.StdSignBytes(ctx.ChainID(), accNums[i], seqs[i], fee, msgs, "")

		sig, err := priv.Sign(signBytes)
		if err != nil {
			panic(err)
		}

		sigs[i] = auth.StdSignature{PubKey: priv.PubKey(), Signature: sig}
	}

	tx := auth.NewStdTx(msgs, fee, sigs, "")
	return tx
}

func newTestTxWithMemo(ctx sdk.Context, msgs []sdk.Msg, privs []crypto.PrivKey, accNums []uint64, seqs []uint64, fee auth.StdFee, memo string) sdk.Tx {
	sigs := make([]auth.StdSignature, len(privs))
	for i, priv := range privs {
		signBytes := auth.StdSignBytes(ctx.ChainID(), accNums[i], seqs[i], fee, msgs, memo)

		sig, err := priv.Sign(signBytes)
		if err != nil {
			panic(err)
		}

		sigs[i] = auth.StdSignature{PubKey: priv.PubKey(), Signature: sig}
	}

	tx := auth.NewStdTx(msgs, fee, sigs, memo)
	return tx
}

func newTestTxWithSignBytes(msgs []sdk.Msg, privs []crypto.PrivKey, accNums []uint64, seqs []uint64, fee auth.StdFee, signBytes []byte, memo string) sdk.Tx {
	sigs := make([]auth.StdSignature, len(privs))
	for i, priv := range privs {
		sig, err := priv.Sign(signBytes)
		if err != nil {
			panic(err)
		}

		sigs[i] = auth.StdSignature{PubKey: priv.PubKey(), Signature: sig}
	}

	tx := auth.NewStdTx(msgs, fee, sigs, memo)
	return tx
}

// MockTreasuryKeeper no-lint
type MockTreasuryKeeper struct{}

// NewMockTreasuryKeeper no-lint
func NewMockTreasuryKeeper() MockTreasuryKeeper {
	return MockTreasuryKeeper{}
}

// GetTaxRate no-lint
func (tk MockTreasuryKeeper) GetTaxRate(ctx sdk.Context, epoch sdk.Int) (rate sdk.Dec) {
	return sdk.NewDecWithPrec(1, 3) // 0.1%
}

// GetTaxCap no-lint
func (tk MockTreasuryKeeper) GetTaxCap(ctx sdk.Context, denom string) (taxCap sdk.Int) {
	return sdk.NewInt(1000000) // 1,000,000 {denom}
}

// RecordTaxProceeds no-lint
func (tk MockTreasuryKeeper) RecordTaxProceeds(_ sdk.Context, _ sdk.Coins) {
	return
}
