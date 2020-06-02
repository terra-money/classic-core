package keeper

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	core "github.com/terra-project/core/types"
	"github.com/terra-project/core/x/wasm/internal/types"
)

func TestStoreCode(t *testing.T) {
	// Create & set temp as home
	tempDir, err := ioutil.TempDir("", "wasmtest")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)
	viper.Set(flags.FlagHome, tempDir)

	input := CreateTestInput(t)
	ctx, accKeeper, keeper := input.Ctx, input.AccKeeper, input.WasmKeeper

	deposit := sdk.NewCoins(sdk.NewInt64Coin(core.MicroLunaDenom, 100000))
	creator := createFakeFundedAccount(ctx, accKeeper, deposit)

	wasmCode, err := ioutil.ReadFile("./testdata/contract.wasm")
	require.NoError(t, err)

	// Create contract
	contractID, err := keeper.StoreCode(ctx, creator, wasmCode, true)
	require.NoError(t, err)
	require.Equal(t, uint64(1), contractID)

	// Verify content
	storedCode, err := keeper.GetByteCode(ctx, contractID)
	require.NoError(t, err)
	require.Equal(t, wasmCode, storedCode)
}

func TestStoreCodeWithHugeCode(t *testing.T) {
	// Create & set temp as home
	tempDir, err := ioutil.TempDir("", "wasmtest")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)
	viper.Set(flags.FlagHome, tempDir)

	input := CreateTestInput(t)
	ctx, keeper := input.Ctx, input.WasmKeeper

	_, _, creator := keyPubAddr()
	wasmCode := make([]byte, keeper.MaxContractSize(ctx)+1)
	_, err = keeper.StoreCode(ctx, creator, wasmCode, true)

	require.Error(t, err)
	require.Contains(t, err.Error(), "contract size is too huge")
}

func TestCreateWithGzippedPayload(t *testing.T) {
	// Create & set temp as home
	tempDir, err := ioutil.TempDir("", "wasmtest")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)
	viper.Set(flags.FlagHome, tempDir)

	input := CreateTestInput(t)
	ctx, accKeeper, keeper := input.Ctx, input.AccKeeper, input.WasmKeeper

	deposit := sdk.NewCoins(sdk.NewInt64Coin(core.MicroLunaDenom, 100000))
	creator := createFakeFundedAccount(ctx, accKeeper, deposit)

	wasmCode, err := ioutil.ReadFile("./testdata/contract.wasm.gzip")
	require.NoError(t, err)

	contractID, err := keeper.StoreCode(ctx, creator, wasmCode, true)
	require.NoError(t, err)
	require.Equal(t, uint64(1), contractID)
	// and verify content
	storedCode, err := keeper.GetByteCode(ctx, contractID)
	require.NoError(t, err)
	rawCode, err := ioutil.ReadFile("./testdata/contract.wasm")
	require.NoError(t, err)
	require.Equal(t, rawCode, storedCode)
}

func TestInstantiate(t *testing.T) {
	// Create & set temp as home
	tempDir, err := ioutil.TempDir("", "wasmtest")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)
	viper.Set(flags.FlagHome, tempDir)
	input := CreateTestInput(t)
	ctx, accKeeper, keeper := input.Ctx, input.AccKeeper, input.WasmKeeper

	deposit := sdk.NewCoins(sdk.NewInt64Coin(core.MicroLunaDenom, 100000))
	creator := createFakeFundedAccount(ctx, accKeeper, deposit)

	wasmCode, err := ioutil.ReadFile("./testdata/contract.wasm")
	require.NoError(t, err)

	codeID, err := keeper.StoreCode(ctx, creator, wasmCode, true)
	require.NoError(t, err)

	_, _, bob := keyPubAddr()
	_, _, fred := keyPubAddr()

	initMsg := InitMsg{
		Verifier:    fred,
		Beneficiary: bob,
	}

	initMsgBz, err := json.Marshal(initMsg)
	require.NoError(t, err)

	// create with no balance is also legal
	addr, err := keeper.InstantiateContract(ctx, codeID, creator, initMsgBz, nil)
	require.NoError(t, err)
	require.Equal(t, "cosmos18vd8fpwxzck93qlwghaj6arh4p7c5n89uzcee5", addr.String())
}

func TestInstantiateWithNonExistingCodeID(t *testing.T) {
	// Create & set temp as home
	tempDir, err := ioutil.TempDir("", "wasmtest")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)
	viper.Set(flags.FlagHome, tempDir)
	input := CreateTestInput(t)
	ctx, accKeeper, keeper := input.Ctx, input.AccKeeper, input.WasmKeeper

	deposit := sdk.NewCoins(sdk.NewInt64Coin(core.MicroLunaDenom, 100000))
	creator := createFakeFundedAccount(ctx, accKeeper, deposit)

	initMsg := InitMsg{}
	initMsgBz, err := json.Marshal(initMsg)
	require.NoError(t, err)

	const nonExistingCodeID = 9999
	_, err = keeper.InstantiateContract(ctx, nonExistingCodeID, creator, initMsgBz, nil)
	require.Error(t, err, sdkerrors.Wrapf(types.ErrNotFound, "codeID %d", nonExistingCodeID))
}

func TestInstantiateWithBigInitMsg(t *testing.T) {
	// Create & set temp as home
	tempDir, err := ioutil.TempDir("", "wasmtest")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)
	viper.Set(flags.FlagHome, tempDir)
	input := CreateTestInput(t)
	ctx, accKeeper, keeper := input.Ctx, input.AccKeeper, input.WasmKeeper

	deposit := sdk.NewCoins(sdk.NewInt64Coin(core.MicroLunaDenom, 100000))
	creator := createFakeFundedAccount(ctx, accKeeper, deposit)

	wasmCode, err := ioutil.ReadFile("./testdata/contract.wasm")
	require.NoError(t, err)

	codeID, err := keeper.StoreCode(ctx, creator, wasmCode, true)
	require.NoError(t, err)

	// test max init msg size
	initMsgBz := make([]byte, keeper.MaxContractMsgSize(ctx)+1)
	_, err = keeper.InstantiateContract(ctx, codeID, creator, initMsgBz, deposit)
	require.Error(t, err)
	require.Contains(t, err.Error(), "init msg size is too huge")
}

func TestExecute(t *testing.T) {
	// Create & set temp as home
	tempDir, err := ioutil.TempDir("", "wasmtest")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)
	viper.Set(flags.FlagHome, tempDir)
	input := CreateTestInput(t)
	ctx, accKeeper, keeper := input.Ctx, input.AccKeeper, input.WasmKeeper

	deposit := sdk.NewCoins(sdk.NewInt64Coin(core.MicroLunaDenom, 100000))
	topUp := sdk.NewCoins(sdk.NewInt64Coin(core.MicroLunaDenom, 5000))
	creator := createFakeFundedAccount(ctx, accKeeper, deposit.Add(deposit...))
	fred := createFakeFundedAccount(ctx, accKeeper, topUp)

	wasmCode, err := ioutil.ReadFile("./testdata/contract.wasm")
	require.NoError(t, err)

	codeID, err := keeper.StoreCode(ctx, creator, wasmCode, true)
	require.NoError(t, err)

	_, _, bob := keyPubAddr()
	initMsg := InitMsg{
		Verifier:    fred,
		Beneficiary: bob,
	}
	initMsgBz, err := json.Marshal(initMsg)
	require.NoError(t, err)

	addr, err := keeper.InstantiateContract(ctx, codeID, creator, initMsgBz, deposit)
	require.NoError(t, err)
	require.Equal(t, "cosmos18vd8fpwxzck93qlwghaj6arh4p7c5n89uzcee5", addr.String())

	// ensure bob doesn't exist
	bobAcct := accKeeper.GetAccount(ctx, bob)
	require.Nil(t, bobAcct)

	// ensure funder has reduced balance
	creatorAcct := accKeeper.GetAccount(ctx, creator)
	require.NotNil(t, creatorAcct)
	// we started at 2*deposit, should have spent one above
	assert.Equal(t, deposit, creatorAcct.GetCoins())

	// ensure contract has updated balance
	contractAcct := accKeeper.GetAccount(ctx, addr)
	require.NotNil(t, contractAcct)
	assert.Equal(t, deposit, contractAcct.GetCoins())

	// unauthorized - trialCtx so we don't change state
	trialCtx := ctx.WithMultiStore(ctx.MultiStore().CacheWrap().(sdk.MultiStore))
	res, err := keeper.ExecuteContract(trialCtx, addr, creator, []byte(`{"release":{}}`), nil)
	require.Error(t, err)
	require.Contains(t, err.Error(), "unauthorized")

	// verifier can execute, and get proper gas amount
	start := time.Now()
	gasBefore := ctx.GasMeter().GasConsumed()

	res, err = keeper.ExecuteContract(ctx, addr, fred, []byte(`{"release":{}}`), topUp)
	diff := time.Now().Sub(start)
	require.NoError(t, err)
	require.NotNil(t, res)

	// make sure gas is properly deducted from ctx
	gasAfter := ctx.GasMeter().GasConsumed()
	require.Equal(t, uint64(0xc739), gasAfter-gasBefore)

	// ensure bob now exists and got both payments released
	bobAcct = accKeeper.GetAccount(ctx, bob)
	require.NotNil(t, bobAcct)
	balance := bobAcct.GetCoins()
	assert.Equal(t, deposit.Add(topUp...), balance)

	// ensure contract has updated balance
	contractAcct = accKeeper.GetAccount(ctx, addr)
	require.NotNil(t, contractAcct)
	assert.Equal(t, sdk.Coins(nil), contractAcct.GetCoins())

	t.Logf("Duration: %v (35619 gas)\n", diff)
}

func TestExecuteWithNonExistingContractAddress(t *testing.T) {
	// Create & set temp as home
	tempDir, err := ioutil.TempDir("", "wasmtest")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)
	viper.Set(flags.FlagHome, tempDir)
	input := CreateTestInput(t)
	ctx, accKeeper, keeper := input.Ctx, input.AccKeeper, input.WasmKeeper

	deposit := sdk.NewCoins(sdk.NewInt64Coin(core.MicroLunaDenom, 100000))
	creator := createFakeFundedAccount(ctx, accKeeper, deposit.Add(deposit...))

	// unauthorized - trialCtx so we don't change state
	nonExistingContractAddress := addrFromUint64(9999)
	_, err = keeper.ExecuteContract(ctx, nonExistingContractAddress, creator, []byte(`{}`), nil)
	require.Error(t, err, sdkerrors.Wrapf(types.ErrNotFound, "contract %s", nonExistingContractAddress))
}

func TestExecuteWithHugeMsg(t *testing.T) {
	// Create & set temp as home
	tempDir, err := ioutil.TempDir("", "wasmtest")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)
	viper.Set(flags.FlagHome, tempDir)
	input := CreateTestInput(t)
	ctx, accKeeper, keeper := input.Ctx, input.AccKeeper, input.WasmKeeper

	deposit := sdk.NewCoins(sdk.NewInt64Coin(core.MicroLunaDenom, 100000))
	topUp := sdk.NewCoins(sdk.NewInt64Coin(core.MicroLunaDenom, 5000))
	creator := createFakeFundedAccount(ctx, accKeeper, deposit.Add(deposit...))
	fred := createFakeFundedAccount(ctx, accKeeper, topUp)

	wasmCode, err := ioutil.ReadFile("./testdata/contract.wasm")
	require.NoError(t, err)

	codeID, err := keeper.StoreCode(ctx, creator, wasmCode, true)
	require.NoError(t, err)

	_, _, bob := keyPubAddr()
	initMsg := InitMsg{
		Verifier:    fred,
		Beneficiary: bob,
	}
	initMsgBz, err := json.Marshal(initMsg)
	require.NoError(t, err)

	addr, err := keeper.InstantiateContract(ctx, codeID, creator, initMsgBz, deposit)
	require.NoError(t, err)
	require.Equal(t, "cosmos18vd8fpwxzck93qlwghaj6arh4p7c5n89uzcee5", addr.String())

	msgBz := make([]byte, keeper.MaxContractMsgSize(ctx)+1)
	_, err = keeper.ExecuteContract(ctx, addr, fred, msgBz, topUp)
	require.Error(t, err)
	require.Contains(t, err.Error(), "execute msg size is too huge")
}

func TestExecuteWithPanic(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "wasm")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)
	viper.Set(flags.FlagHome, tempDir)

	input := CreateTestInput(t)
	ctx, accKeeper, keeper := input.Ctx, input.AccKeeper, input.WasmKeeper

	deposit := sdk.NewCoins(sdk.NewInt64Coin(core.MicroLunaDenom, 100000))
	topUp := sdk.NewCoins(sdk.NewInt64Coin(core.MicroLunaDenom, 5000))
	creator := createFakeFundedAccount(ctx, accKeeper, deposit.Add(deposit...))
	fred := createFakeFundedAccount(ctx, accKeeper, topUp)

	wasmCode, err := ioutil.ReadFile("./testdata/contract.wasm")
	require.NoError(t, err)

	contractID, err := keeper.StoreCode(ctx, creator, wasmCode, true)
	require.NoError(t, err)

	_, _, bob := keyPubAddr()
	initMsg := InitMsg{
		Verifier:    fred,
		Beneficiary: bob,
	}
	initMsgBz, err := json.Marshal(initMsg)
	require.NoError(t, err)

	addr, err := keeper.InstantiateContract(ctx, contractID, creator, initMsgBz, deposit)
	require.NoError(t, err)

	// let's make sure we get a reasonable error, no panic/crash
	_, err = keeper.ExecuteContract(ctx, addr, fred, []byte(`{"panic":{}}`), topUp)
	require.Error(t, err)
}

func TestExecuteWithCpuLoop(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "wasm")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)
	viper.Set(flags.FlagHome, tempDir)
	input := CreateTestInput(t)
	ctx, accKeeper, keeper := input.Ctx, input.AccKeeper, input.WasmKeeper

	deposit := sdk.NewCoins(sdk.NewInt64Coin(core.MicroLunaDenom, 100000))
	topUp := sdk.NewCoins(sdk.NewInt64Coin(core.MicroLunaDenom, 5000))
	creator := createFakeFundedAccount(ctx, accKeeper, deposit.Add(deposit...))
	fred := createFakeFundedAccount(ctx, accKeeper, topUp)

	wasmCode, err := ioutil.ReadFile("./testdata/contract.wasm")
	require.NoError(t, err)

	contractID, err := keeper.StoreCode(ctx, creator, wasmCode, true)
	require.NoError(t, err)

	_, _, bob := keyPubAddr()
	initMsg := InitMsg{
		Verifier:    fred,
		Beneficiary: bob,
	}
	initMsgBz, err := json.Marshal(initMsg)
	require.NoError(t, err)

	addr, err := keeper.InstantiateContract(ctx, contractID, creator, initMsgBz, deposit)
	require.NoError(t, err)

	// make sure we set a limit before calling
	var gasLimit uint64 = 400_000
	ctx = ctx.WithGasMeter(sdk.NewGasMeter(gasLimit))
	require.Equal(t, uint64(0), ctx.GasMeter().GasConsumed())

	// this must fail
	_, err = keeper.ExecuteContract(ctx, addr, fred, []byte(`{"cpu_loop":{}}`), nil)
	require.Error(t, err)
}

func TestExecuteWithStorageLoop(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "wasm")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)
	viper.Set(flags.FlagHome, tempDir)
	input := CreateTestInput(t)
	ctx, accKeeper, keeper := input.Ctx, input.AccKeeper, input.WasmKeeper

	deposit := sdk.NewCoins(sdk.NewInt64Coin(core.MicroLunaDenom, 100000))
	topUp := sdk.NewCoins(sdk.NewInt64Coin(core.MicroLunaDenom, 5000))
	creator := createFakeFundedAccount(ctx, accKeeper, deposit.Add(deposit...))
	fred := createFakeFundedAccount(ctx, accKeeper, topUp)

	wasmCode, err := ioutil.ReadFile("./testdata/contract.wasm")
	require.NoError(t, err)

	contractID, err := keeper.StoreCode(ctx, creator, wasmCode, true)
	require.NoError(t, err)

	_, _, bob := keyPubAddr()
	initMsg := InitMsg{
		Verifier:    fred,
		Beneficiary: bob,
	}
	initMsgBz, err := json.Marshal(initMsg)
	require.NoError(t, err)

	addr, err := keeper.InstantiateContract(ctx, contractID, creator, initMsgBz, deposit)
	require.NoError(t, err)

	// make sure we set a limit before calling
	var gasLimit uint64 = 400_000
	ctx = ctx.WithGasMeter(sdk.NewGasMeter(gasLimit))
	require.Equal(t, uint64(0), ctx.GasMeter().GasConsumed())

	// ensure we get an out of gas panic
	require.PanicsWithValue(t, sdk.ErrorOutOfGas{Descriptor: "ReadFlat"}, func() {
		_, err = keeper.ExecuteContract(ctx, addr, fred, []byte(`{"storage_loop":{}}`), nil)
	})
}