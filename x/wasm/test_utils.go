// nolint:deadcode unused noalias
package wasm

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"

	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/spf13/viper"

	wasmTypes "github.com/CosmWasm/go-cosmwasm/types"

	"github.com/stretchr/testify/require"

	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/ed25519"

	"github.com/terra-money/core/x/auth"
	"github.com/terra-money/core/x/wasm/internal/keeper"
)

var (
	key1, pub1, addr1 = keyPubAddr()
	testContract      []byte
	maskContract      []byte
	oldEscrowContract []byte
)

type testData struct {
	module     module.AppModule
	ctx        sdk.Context
	acctKeeper auth.AccountKeeper
	keeper     Keeper
}

func loadContracts() {
	testContract = mustLoad("./internal/keeper/testdata/contract.wasm")
	maskContract = mustLoad("./internal/keeper/testdata/mask.wasm")
	oldEscrowContract = mustLoad("./testdata/escrow_0.7.wasm")
}

// Returns a cleanup function, which must be defered on
func setupTest(t *testing.T) (testData, func()) {
	// Create & set temp as home
	tempDir, err := ioutil.TempDir("", "wasmtest")
	require.NoError(t, err)
	viper.Set(flags.FlagHome, tempDir)

	input := keeper.CreateTestInput(t)
	data := testData{
		module:     NewAppModule(input.WasmKeeper, input.AccKeeper, input.BankKeeper),
		ctx:        input.Ctx,
		acctKeeper: input.AccKeeper,
		keeper:     input.WasmKeeper,
	}
	cleanup := func() { os.RemoveAll(tempDir) }
	return data, cleanup
}

func keyPubAddr() (crypto.PrivKey, crypto.PubKey, sdk.AccAddress) {
	key := ed25519.GenPrivKey()
	pub := key.PubKey()
	addr := sdk.AccAddress(pub.Address())
	return key, pub, addr
}

func mustLoad(path string) []byte {
	bz, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return bz
}

type initMsg struct {
	Verifier    string `json:"verifier"`
	Beneficiary string `json:"beneficiary"`
}

type state struct {
	Verifier    wasmTypes.CanonicalAddress `json:"verifier"`
	Beneficiary wasmTypes.CanonicalAddress `json:"beneficiary"`
	Funder      wasmTypes.CanonicalAddress `json:"funder"`
}

func createFakeFundedAccount(ctx sdk.Context, am auth.AccountKeeper, coins sdk.Coins) sdk.AccAddress {
	_, _, addr := keyPubAddr()
	baseAcct := auth.NewBaseAccountWithAddress(addr)
	_ = baseAcct.SetCoins(coins)
	am.SetAccount(ctx, &baseAcct)

	return addr
}

func assertContractStore(t *testing.T, models []Model, expected state) {
	require.Equal(t, 1, len(models), "#v", models)
	require.Equal(t, []byte("config"), models[0].Key.Bytes())

	expectedBz, err := json.Marshal(expected)
	require.NoError(t, err)
	require.Equal(t, expectedBz, models[0].Value.Bytes())
}
