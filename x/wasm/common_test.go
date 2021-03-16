// nolint:deadcode unused noalias
package wasm_test

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"

	wasmTypes "github.com/CosmWasm/go-cosmwasm/types"

	"github.com/stretchr/testify/require"

	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/ed25519"

	"github.com/terra-project/core/x/wasm/types"
)

var (
	key1, pub1, addr1 = keyPubAddr()
	testContract      []byte
	maskContract      []byte
	oldEscrowContract []byte
)

func loadContracts() {
	testContract = mustLoad("./keeper/testdata/contract.wasm")
	maskContract = mustLoad("./keeper/testdata/mask.wasm")
	oldEscrowContract = mustLoad("./testdata/escrow_0.7.wasm")
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

func createFakeFundedAccount(
	ctx sdk.Context,
	am authkeeper.AccountKeeper,
	bk bankkeeper.Keeper,
	coins sdk.Coins) sdk.AccAddress {
	_, _, addr := keyPubAddr()
	baseAcct := authtypes.NewBaseAccountWithAddress(addr)
	am.SetAccount(ctx, baseAcct)
	bk.SetBalances(ctx, addr, coins)

	return addr
}

func assertContractStore(t *testing.T, models []types.Model, expected state) {
	require.Equal(t, 1, len(models), "#v", models)
	require.Equal(t, []byte("config"), models[0].Key)

	expectedBz, err := json.Marshal(expected)
	require.NoError(t, err)
	require.Equal(t, expectedBz, models[0].Value)
}
