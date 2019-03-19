package oracle

import (
	"terra/types/assets"

	"github.com/cosmos/cosmos-sdk/x/staking"

	"testing"
	"time"

	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/params"
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

	initAmt = sdk.NewInt(1005)
	lunaAmt = sdk.NewInt(10)
)

type testInput struct {
	ctx          sdk.Context
	accKeeper    auth.AccountKeeper
	bankKeeper   bank.Keeper
	oracleKeeper Keeper
	valset       MockValset
}

func newTestCodec() *codec.Codec {
	cdc := codec.New()

	bank.RegisterCodec(cdc)
	RegisterCodec(cdc)
	auth.RegisterCodec(cdc)
	staking.RegisterCodec(cdc)
	sdk.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)

	return cdc
}

func createTestInput(t *testing.T) testInput {
	keyAcc := sdk.NewKVStoreKey(auth.StoreKey)
	keyParams := sdk.NewKVStoreKey(params.StoreKey)
	tKeyParams := sdk.NewTransientStoreKey(params.TStoreKey)
	keyOracle := sdk.NewKVStoreKey(StoreKey)
	keyStaking := sdk.NewKVStoreKey(staking.StoreKey)
	tkeyStaking := sdk.NewKVStoreKey(staking.TStoreKey)

	cdc := newTestCodec()
	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)
	ctx := sdk.NewContext(ms, abci.Header{Time: time.Now().UTC()}, false, log.NewNopLogger())

	ms.MountStoreWithDB(keyAcc, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(tKeyParams, sdk.StoreTypeTransient, db)
	ms.MountStoreWithDB(keyParams, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyOracle, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyStaking, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(tkeyStaking, sdk.StoreTypeIAVL, db)

	require.NoError(t, ms.LoadLatestVersion())

	paramsKeeper := params.NewKeeper(cdc, keyParams, tKeyParams)
	accKeeper := auth.NewAccountKeeper(
		cdc,
		keyAcc,
		paramsKeeper.Subspace(auth.DefaultParamspace),
		auth.ProtoBaseAccount,
	)

	bankKeeper := bank.NewBaseKeeper(
		accKeeper,
		paramsKeeper.Subspace(bank.DefaultParamspace),
		bank.DefaultCodespace,
	)

	valset := NewMockValSet()
	for _, addr := range addrs {
		_, _, err := bankKeeper.AddCoins(ctx, addr, sdk.Coins{
			sdk.NewCoin(assets.SDRDenom, initAmt),
			sdk.NewCoin(assets.LunaDenom, lunaAmt),
		})
		require.NoError(t, err)

		// Add validators
		validator := NewMockValidator(sdk.ValAddress(addr.Bytes()), lunaAmt)
		valset.validators = append(valset.validators, validator)
	}

	oracleKeeper := NewKeeper(
		cdc, keyOracle, valset,
		paramsKeeper.Subspace(DefaultParamspace),
	)

	return testInput{ctx, accKeeper, bankKeeper, oracleKeeper, valset}
}

type MockValset struct {
	sdk.ValidatorSet

	validators []MockValidator
}

type MockValidator struct {
	sdk.Validator

	address sdk.ValAddress
	power   sdk.Int
}

func NewMockValSet() MockValset {
	return MockValset{
		validators: []MockValidator{},
	}
}

func NewMockValidator(address sdk.ValAddress, power sdk.Int) MockValidator {
	return MockValidator{
		address: address,
		power:   power,
	}
}

func (mv MockValidator) GetBondedTokens() sdk.Int {
	return mv.power
}

func (mv MockValset) Validator(ctx sdk.Context, valAddress sdk.ValAddress) sdk.Validator {
	for _, val := range mv.validators {
		if val.address.Equals(valAddress) {
			return val
		}
	}
	return nil
}

func (mv MockValset) TotalBondedTokens(ctx sdk.Context) sdk.Int {
	rval := sdk.ZeroInt()
	for _, val := range mv.validators {
		rval = rval.Add(val.power)
	}
	return rval
}
