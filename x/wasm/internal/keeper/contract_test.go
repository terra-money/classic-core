package keeper

import (
	"bytes"
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

	core "github.com/terra-money/core/types"
	"github.com/terra-money/core/x/wasm/internal/types"
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
	contractID, err := keeper.StoreCode(ctx, creator, wasmCode)
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
	_, err = keeper.StoreCode(ctx, creator, wasmCode)

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

	contractID, err := keeper.StoreCode(ctx, creator, wasmCode)
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

	codeID, err := keeper.StoreCode(ctx, creator, wasmCode)
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
	addr, err := keeper.InstantiateContract(ctx, codeID, creator, initMsgBz, nil, true)
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
	_, err = keeper.InstantiateContract(ctx, nonExistingCodeID, creator, initMsgBz, nil, true)
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

	codeID, err := keeper.StoreCode(ctx, creator, wasmCode)
	require.NoError(t, err)

	// test max init msg size
	initMsgBz := make([]byte, keeper.MaxContractMsgSize(ctx)+1)
	_, err = keeper.InstantiateContract(ctx, codeID, creator, initMsgBz, deposit, true)
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

	codeID, err := keeper.StoreCode(ctx, creator, wasmCode)
	require.NoError(t, err)

	_, _, bob := keyPubAddr()
	initMsg := InitMsg{
		Verifier:    fred,
		Beneficiary: bob,
	}
	initMsgBz, err := json.Marshal(initMsg)
	require.NoError(t, err)

	addr, err := keeper.InstantiateContract(ctx, codeID, creator, initMsgBz, deposit, true)
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
	require.True(t, gasAfter-gasBefore > types.InstanceCost)

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

	codeID, err := keeper.StoreCode(ctx, creator, wasmCode)
	require.NoError(t, err)

	_, _, bob := keyPubAddr()
	initMsg := InitMsg{
		Verifier:    fred,
		Beneficiary: bob,
	}
	initMsgBz, err := json.Marshal(initMsg)
	require.NoError(t, err)

	addr, err := keeper.InstantiateContract(ctx, codeID, creator, initMsgBz, deposit, true)
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

	contractID, err := keeper.StoreCode(ctx, creator, wasmCode)
	require.NoError(t, err)

	_, _, bob := keyPubAddr()
	initMsg := InitMsg{
		Verifier:    fred,
		Beneficiary: bob,
	}
	initMsgBz, err := json.Marshal(initMsg)
	require.NoError(t, err)

	addr, err := keeper.InstantiateContract(ctx, contractID, creator, initMsgBz, deposit, true)
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

	contractID, err := keeper.StoreCode(ctx, creator, wasmCode)
	require.NoError(t, err)

	_, _, bob := keyPubAddr()
	initMsg := InitMsg{
		Verifier:    fred,
		Beneficiary: bob,
	}
	initMsgBz, err := json.Marshal(initMsg)
	require.NoError(t, err)

	addr, err := keeper.InstantiateContract(ctx, contractID, creator, initMsgBz, deposit, true)
	require.NoError(t, err)

	// make sure we set a limit before calling
	var gasLimit uint64 = 400_000
	ctx = ctx.WithGasMeter(sdk.NewGasMeter(gasLimit))
	require.Equal(t, uint64(0), ctx.GasMeter().GasConsumed())

	require.PanicsWithValue(t, sdk.ErrorOutOfGas{Descriptor: "Contract Execution"}, func() {
		_, _ = keeper.ExecuteContract(ctx, addr, fred, []byte(`{"cpu_loop":{}}`), nil)
	})
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

	contractID, err := keeper.StoreCode(ctx, creator, wasmCode)
	require.NoError(t, err)

	_, _, bob := keyPubAddr()
	initMsg := InitMsg{
		Verifier:    fred,
		Beneficiary: bob,
	}
	initMsgBz, err := json.Marshal(initMsg)
	require.NoError(t, err)

	addr, err := keeper.InstantiateContract(ctx, contractID, creator, initMsgBz, deposit, true)
	require.NoError(t, err)

	// make sure we set a limit before calling
	var gasLimit uint64 = 400_000
	ctx = ctx.WithGasMeter(sdk.NewGasMeter(gasLimit))
	require.Equal(t, uint64(0), ctx.GasMeter().GasConsumed())

	// ensure we get an out of gas panic
	require.PanicsWithValue(t, sdk.ErrorOutOfGas{Descriptor: "Contract Execution"}, func() {
		_, err = keeper.ExecuteContract(ctx, addr, fred, []byte(`{"storage_loop":{}}`), nil)
	})
}

func TestMigrate(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "wasm")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)
	viper.Set(flags.FlagHome, tempDir)
	input := CreateTestInput(t)
	ctx, accKeeper, keeper := input.Ctx, input.AccKeeper, input.WasmKeeper

	deposit := sdk.NewCoins(sdk.NewInt64Coin("denom", 100000))
	creator := createFakeFundedAccount(ctx, accKeeper, deposit.Add(deposit...))
	fred := createFakeFundedAccount(ctx, accKeeper, sdk.NewCoins(sdk.NewInt64Coin("denom", 5000)))

	wasmCode, err := ioutil.ReadFile("./testdata/contract.wasm")
	require.NoError(t, err)

	originalContractID, err := keeper.StoreCode(ctx, creator, wasmCode)
	require.NoError(t, err)
	newContractID, err := keeper.StoreCode(ctx, creator, wasmCode)
	require.NoError(t, err)
	require.NotEqual(t, originalContractID, newContractID)

	_, _, anyAddr := keyPubAddr()
	_, _, newVerifierAddr := keyPubAddr()
	initMsg := InitMsg{
		Verifier:    fred,
		Beneficiary: anyAddr,
	}
	initMsgBz, err := json.Marshal(initMsg)
	require.NoError(t, err)

	migMsg := struct {
		Verifier sdk.AccAddress `json:"verifier"`
	}{Verifier: newVerifierAddr}
	migMsgBz, err := json.Marshal(migMsg)
	require.NoError(t, err)

	specs := map[string]struct {
		migratable           bool
		overrideContractAddr sdk.AccAddress
		caller               sdk.AccAddress
		codeID               uint64
		migrateMsg           []byte
		expErr               *sdkerrors.Error
		expVerifier          sdk.AccAddress
	}{
		"all good with same code id": {
			migratable:  true,
			caller:      creator,
			codeID:      originalContractID,
			migrateMsg:  migMsgBz,
			expVerifier: newVerifierAddr,
		},
		"all good with different code id": {
			migratable:  true,
			caller:      creator,
			codeID:      newContractID,
			migrateMsg:  migMsgBz,
			expVerifier: newVerifierAddr,
		},
		"prevent migration when migratable is false on instantiate": {
			migratable: false,
			caller:     creator,
			codeID:     originalContractID,
			expErr:     types.ErrNotMigratable,
		},
		"prevent migration when not sent by owner": {
			migratable: true,
			caller:     fred,
			codeID:     originalContractID,
			expErr:     sdkerrors.ErrUnauthorized,
		},
		"fail with non existing code id": {
			migratable: true,
			caller:     creator,
			codeID:     99999,
			expErr:     types.ErrNotFound,
		},
		"fail with non existing contract addr": {
			migratable:           true,
			caller:               creator,
			overrideContractAddr: anyAddr,
			codeID:               originalContractID,
			expErr:               types.ErrNotFound,
		},
		"fail in contract with invalid migrate msg": {
			migratable: true,
			caller:     creator,
			codeID:     originalContractID,
			migrateMsg: bytes.Repeat([]byte{0x1}, 7),
			expErr:     types.ErrMigrationFailed,
		},
		"fail in contract without migrate msg": {
			migratable: true,
			caller:     creator,
			codeID:     originalContractID,
			expErr:     types.ErrMigrationFailed,
		},
	}

	for msg, spec := range specs {
		t.Run(msg, func(t *testing.T) {
			ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)
			addr, err := keeper.InstantiateContract(ctx, originalContractID, creator, initMsgBz, nil, spec.migratable)
			require.NoError(t, err)
			if spec.overrideContractAddr != nil {
				addr = spec.overrideContractAddr
			}
			_, err = keeper.MigrateContract(ctx, addr, spec.caller, spec.codeID, spec.migrateMsg)
			require.True(t, spec.expErr.Is(err), "expected %v but got %+v", spec.expErr, err)
			if spec.expErr != nil {
				return
			}
			cInfo, err := keeper.GetContractInfo(ctx, addr)
			require.NoError(t, err)
			assert.Equal(t, spec.codeID, cInfo.CodeID)

			m := keeper.queryToStore(ctx, addr, []byte("config"))
			var stored map[string][]byte
			require.NoError(t, json.Unmarshal(m, &stored))
			require.Contains(t, stored, "verifier")
			require.NoError(t, err)
			assert.Equal(t, spec.expVerifier, sdk.AccAddress(stored["verifier"]))
		})
	}
}

func TestMigrateWithDispatchedMessage(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "wasm")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)
	viper.Set(flags.FlagHome, tempDir)
	input := CreateTestInput(t)
	ctx, accKeeper, keeper := input.Ctx, input.AccKeeper, input.WasmKeeper

	deposit := sdk.NewCoins(sdk.NewInt64Coin(core.MicroLunaDenom, 100000))
	creator := createFakeFundedAccount(ctx, accKeeper, deposit.Add(deposit...))
	fred := createFakeFundedAccount(ctx, accKeeper, sdk.NewCoins(sdk.NewInt64Coin(core.MicroLunaDenom, 5000)))

	wasmCode, err := ioutil.ReadFile("./testdata/contract.wasm")
	require.NoError(t, err)
	burnerCode, err := ioutil.ReadFile("./testdata/burner.wasm")
	require.NoError(t, err)

	originalContractID, err := keeper.StoreCode(ctx, creator, wasmCode)
	require.NoError(t, err)
	burnerContractID, err := keeper.StoreCode(ctx, creator, burnerCode)
	require.NoError(t, err)
	require.NotEqual(t, originalContractID, burnerContractID)

	_, _, myPayoutAddr := keyPubAddr()
	initMsg := InitMsg{
		Verifier:    fred,
		Beneficiary: fred,
	}
	initMsgBz, err := json.Marshal(initMsg)
	require.NoError(t, err)

	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)
	contractAddr, err := keeper.InstantiateContract(ctx, originalContractID, creator, initMsgBz, deposit, true)
	require.NoError(t, err)

	migMsg := struct {
		Payout sdk.AccAddress `json:"payout"`
	}{Payout: myPayoutAddr}
	migMsgBz, err := json.Marshal(migMsg)
	require.NoError(t, err)
	keeper.loggingWhitelist[contractAddr.String()] = true
	ctx = ctx.WithEventManager(sdk.NewEventManager()).WithBlockHeight(ctx.BlockHeight() + 1)
	res, err := keeper.MigrateContract(ctx, contractAddr, creator, burnerContractID, migMsgBz)
	require.NoError(t, err)
	assert.Equal(t, "burnt 1 keys", string(res))
	type dict map[string]interface{}
	expEvents := []dict{
		{
			"Type": "from_contract",
			"Attr": []dict{
				{"contract_address": contractAddr},
				{"action": "burn"},
				{"payout": myPayoutAddr},
			},
		},
		{
			"Type": "transfer",
			"Attr": []dict{
				{"recipient": myPayoutAddr},
				{"sender": contractAddr},
				{"amount": "100000" + core.MicroLunaDenom},
			},
		},
		{
			"Type": "message",
			"Attr": []dict{
				{"sender": contractAddr},
			},
		},
		{
			"Type": "message",
			"Attr": []dict{
				{"module": "bank"},
			},
		},
	}
	expJSONEvts := string(mustMarshal(t, expEvents))
	assert.JSONEq(t, expJSONEvts, prettyEvents(t, ctx.EventManager().Events()))

	// all persistent data cleared
	m := keeper.queryToStore(ctx, contractAddr, []byte("config"))
	require.Len(t, m, 0)

	// and all deposit tokens sent to myPayoutAddr
	balance := accKeeper.GetAccount(ctx, myPayoutAddr).GetCoins()
	assert.Equal(t, deposit, balance)
}

func prettyEvents(t *testing.T, events sdk.Events) string {
	t.Helper()
	type prettyEvent struct {
		Type string
		Attr []map[string]string
	}

	r := make([]prettyEvent, len(events))
	for i, e := range events {
		attr := make([]map[string]string, len(e.Attributes))
		for j, a := range e.Attributes {
			attr[j] = map[string]string{string(a.Key): string(a.Value)}
		}
		r[i] = prettyEvent{Type: e.Type, Attr: attr}
	}
	return string(mustMarshal(t, r))
}

func mustMarshal(t *testing.T, r interface{}) []byte {
	t.Helper()
	bz, err := json.Marshal(r)
	require.NoError(t, err)
	return bz
}
