package keeper

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	core "github.com/terra-money/core/types"
	"github.com/terra-money/core/x/wasm/types"
)

func TestStoreCode(t *testing.T) {

	input := CreateTestInput(t)
	ctx, accKeeper, bankKeeper, keeper := input.Ctx, input.AccKeeper, input.BankKeeper, input.WasmKeeper

	deposit := sdk.NewCoins(sdk.NewInt64Coin(core.MicroLunaDenom, 100000))
	creator := createFakeFundedAccount(ctx, accKeeper, bankKeeper, deposit)

	wasmCode, err := ioutil.ReadFile("./testdata/hackatom.wasm")
	require.NoError(t, err)

	// Create contract
	codeID, err := keeper.StoreCode(ctx, creator, wasmCode)
	require.NoError(t, err)
	require.Equal(t, uint64(1), codeID)

	// Verify content
	storedCode, err := keeper.GetByteCode(ctx, codeID)
	require.NoError(t, err)
	require.Equal(t, wasmCode, storedCode)
}

func TestMigrateCode(t *testing.T) {
	input := CreateTestInput(t)
	ctx, accKeeper, bankKeeper, keeper := input.Ctx, input.AccKeeper, input.BankKeeper, input.WasmKeeper

	deposit := sdk.NewCoins(sdk.NewInt64Coin(core.MicroLunaDenom, 100000))
	creator := createFakeFundedAccount(ctx, accKeeper, bankKeeper, deposit)
	fakeAccount := createFakeFundedAccount(ctx, accKeeper, bankKeeper, deposit)

	codeID := uint64(1)
	keeper.SetCodeInfo(ctx, codeID, types.CodeInfo{
		CodeID:   1,
		CodeHash: []byte{},
		Creator:  creator.String(),
	})

	wasmCode, err := ioutil.ReadFile("./testdata/hackatom.wasm")
	require.NoError(t, err)

	err = keeper.MigrateCode(ctx, codeID, fakeAccount, wasmCode)
	require.Error(t, err)

	err = keeper.MigrateCode(ctx, codeID, creator, wasmCode)
	require.NoError(t, err)
	require.Equal(t, uint64(1), codeID)

	// Verify content
	storedCode, err := keeper.GetByteCode(ctx, codeID)
	require.NoError(t, err)
	require.Equal(t, wasmCode, storedCode)

	// Migration failed for the code which contains valid CodeHash
	err = keeper.MigrateCode(ctx, codeID, creator, wasmCode)
	require.Error(t, err)
}

func TestStoreCodeWithHugeCode(t *testing.T) {
	input := CreateTestInput(t)
	ctx, keeper := input.Ctx, input.WasmKeeper

	_, _, creator := keyPubAddr()
	wasmCode := make([]byte, keeper.MaxContractSize(ctx)+1)
	_, err := keeper.StoreCode(ctx, creator, wasmCode)

	require.Error(t, err)
	require.Contains(t, err.Error(), "contract size is too huge")
}

func TestCreateWithGzippedPayload(t *testing.T) {
	input := CreateTestInput(t)
	ctx, accKeeper, bankKeeper, keeper := input.Ctx, input.AccKeeper, input.BankKeeper, input.WasmKeeper

	deposit := sdk.NewCoins(sdk.NewInt64Coin(core.MicroLunaDenom, 100000))
	creator := createFakeFundedAccount(ctx, accKeeper, bankKeeper, deposit)

	wasmCode, err := ioutil.ReadFile("./testdata/hackatom.wasm.gzip")
	require.NoError(t, err)

	contractID, err := keeper.StoreCode(ctx, creator, wasmCode)
	require.NoError(t, err)
	require.Equal(t, uint64(1), contractID)
	// and verify content
	storedCode, err := keeper.GetByteCode(ctx, contractID)
	require.NoError(t, err)
	rawCode, err := ioutil.ReadFile("./testdata/hackatom.wasm")
	require.NoError(t, err)
	require.Equal(t, rawCode, storedCode)
}

func TestInstantiate(t *testing.T) {
	input := CreateTestInput(t)
	ctx, accKeeper, bankKeeper, keeper := input.Ctx, input.AccKeeper, input.BankKeeper, input.WasmKeeper

	deposit := sdk.NewCoins(sdk.NewInt64Coin(core.MicroLunaDenom, 100000))
	creator := createFakeFundedAccount(ctx, accKeeper, bankKeeper, deposit)

	wasmCode, err := ioutil.ReadFile("./testdata/hackatom.wasm")
	require.NoError(t, err)

	codeID, err := keeper.StoreCode(ctx, creator, wasmCode)
	require.NoError(t, err)

	_, _, bob := keyPubAddr()
	_, _, fred := keyPubAddr()

	initMsg := HackatomExampleInitMsg{
		Verifier:    fred,
		Beneficiary: bob,
	}

	initMsgBz, err := json.Marshal(initMsg)
	require.NoError(t, err)

	// create with no balance is also legal
	addr, _, err := keeper.InstantiateContract(ctx, codeID, creator, sdk.AccAddress{}, initMsgBz, nil)
	require.NoError(t, err)
	require.Equal(t, "cosmos18vd8fpwxzck93qlwghaj6arh4p7c5n89uzcee5", addr.String())
}

func TestInstantiateWithNonExistingCodeID(t *testing.T) {
	input := CreateTestInput(t)
	ctx, accKeeper, bankKeeper, keeper := input.Ctx, input.AccKeeper, input.BankKeeper, input.WasmKeeper

	deposit := sdk.NewCoins(sdk.NewInt64Coin(core.MicroLunaDenom, 100000))
	creator := createFakeFundedAccount(ctx, accKeeper, bankKeeper, deposit)

	initMsg := HackatomExampleInitMsg{}
	initMsgBz, err := json.Marshal(initMsg)
	require.NoError(t, err)

	const nonExistingCodeID = 9999
	_, _, err = keeper.InstantiateContract(ctx, nonExistingCodeID, creator, sdk.AccAddress{}, initMsgBz, nil)
	require.Error(t, err, sdkerrors.Wrapf(types.ErrNotFound, "codeID %d", nonExistingCodeID))
}

func TestInstantiateWithBigInitMsg(t *testing.T) {
	input := CreateTestInput(t)
	ctx, accKeeper, bankKeeper, keeper := input.Ctx, input.AccKeeper, input.BankKeeper, input.WasmKeeper

	deposit := sdk.NewCoins(sdk.NewInt64Coin(core.MicroLunaDenom, 100000))
	creator := createFakeFundedAccount(ctx, accKeeper, bankKeeper, deposit)

	wasmCode, err := ioutil.ReadFile("./testdata/hackatom.wasm")
	require.NoError(t, err)

	codeID, err := keeper.StoreCode(ctx, creator, wasmCode)
	require.NoError(t, err)

	// test max init msg size
	initMsgBz := make([]byte, keeper.MaxContractMsgSize(ctx)+1)
	_, _, err = keeper.InstantiateContract(ctx, codeID, creator, sdk.AccAddress{}, initMsgBz, deposit)
	require.Error(t, err)
	require.Contains(t, err.Error(), "init msg size is too huge")
}

func TestInstantiateIBCEnabled(t *testing.T) {
	input := CreateTestInput(t)
	ctx, accKeeper, bankKeeper, keeper := input.Ctx, input.AccKeeper, input.BankKeeper, input.WasmKeeper

	deposit := sdk.NewCoins(sdk.NewInt64Coin(core.MicroLunaDenom, 100000))
	creator := createFakeFundedAccount(ctx, accKeeper, bankKeeper, deposit)

	ibcWasmCode, err := ioutil.ReadFile("./testdata/ibc_reflect.wasm")
	require.NoError(t, err)

	reflectWasmCode, err := ioutil.ReadFile("./testdata/reflect.wasm")
	require.NoError(t, err)

	ibcCodeID, err := keeper.StoreCode(ctx, creator, ibcWasmCode)
	require.NoError(t, err)

	reflectCodeID, err := keeper.StoreCode(ctx, creator, reflectWasmCode)
	require.NoError(t, err)

	ibcInitMsg := IBCReflectInitMsg{ReflectCodeID: reflectCodeID}
	ibcInitMsgBz, err := json.Marshal(ibcInitMsg)
	require.NoError(t, err)

	// create with no balance is also legal
	addr, _, err := keeper.InstantiateContract(ctx, ibcCodeID, creator, sdk.AccAddress{}, ibcInitMsgBz, nil)
	require.NoError(t, err)
	require.Equal(t, "cosmos18vd8fpwxzck93qlwghaj6arh4p7c5n89uzcee5", addr.String())

	cInfo, err := keeper.GetContractInfo(ctx, addr)
	require.NoError(t, err)
	assert.Equal(t, ibcCodeID, cInfo.CodeID)
	assert.Equal(t, cInfo.IBCPortID, types.PortIDForContract(addr))
}

func TestExecute(t *testing.T) {
	input := CreateTestInput(t)
	ctx, accKeeper, bankKeeper, keeper := input.Ctx, input.AccKeeper, input.BankKeeper, input.WasmKeeper

	deposit := sdk.NewCoins(sdk.NewInt64Coin(core.MicroLunaDenom, 100000))
	topUp := sdk.NewCoins(sdk.NewInt64Coin(core.MicroLunaDenom, 5000))
	creator := createFakeFundedAccount(ctx, accKeeper, bankKeeper, deposit.Add(deposit...))
	fred := createFakeFundedAccount(ctx, accKeeper, bankKeeper, topUp)

	wasmCode, err := ioutil.ReadFile("./testdata/hackatom.wasm")
	require.NoError(t, err)

	codeID, err := keeper.StoreCode(ctx, creator, wasmCode)
	require.NoError(t, err)

	_, _, bob := keyPubAddr()
	initMsg := HackatomExampleInitMsg{
		Verifier:    fred,
		Beneficiary: bob,
	}
	initMsgBz, err := json.Marshal(initMsg)
	require.NoError(t, err)

	addr, _, err := keeper.InstantiateContract(ctx, codeID, creator, sdk.AccAddress{}, initMsgBz, deposit)
	require.NoError(t, err)
	require.Equal(t, "cosmos18vd8fpwxzck93qlwghaj6arh4p7c5n89uzcee5", addr.String())

	// ensure bob doesn't exist
	bobAcct := accKeeper.GetAccount(ctx, bob)
	require.Nil(t, bobAcct)

	// ensure funder has reduced balance
	creatorAcct := accKeeper.GetAccount(ctx, creator)
	require.NotNil(t, creatorAcct)
	// we started at 2*deposit, should have spent one above
	assert.Equal(t, deposit, bankKeeper.GetAllBalances(input.Ctx, creator))

	// ensure contract has updated balance
	contractAcct := accKeeper.GetAccount(ctx, addr)
	require.NotNil(t, contractAcct)
	assert.Equal(t, deposit, bankKeeper.GetAllBalances(input.Ctx, addr))

	// unauthorized - trialCtx so we don't change state
	trialCtx := ctx.WithMultiStore(ctx.MultiStore().CacheWrap().(sdk.MultiStore))
	res, err := keeper.ExecuteContract(trialCtx, addr, creator, []byte(`{"release":{}}`), nil)
	require.Error(t, err)
	require.Equal(t, "Unauthorized: execute wasm contract failed", err.Error())

	// verifier can execute, and get proper gas amount
	start := time.Now()
	gasBefore := ctx.GasMeter().GasConsumed()

	res, err = keeper.ExecuteContract(ctx, addr, fred, []byte(`{"release":{}}`), topUp)
	diff := time.Now().Sub(start)
	require.NoError(t, err)
	require.NotNil(t, res)

	// make sure gas is properly deducted from ctx
	gasAfter := ctx.GasMeter().GasConsumed()
	require.True(t, gasAfter-gasBefore > types.InstantiateContractCosts(0))

	// ensure bob now exists and got both payments released
	bobAcct = accKeeper.GetAccount(ctx, bob)
	require.NotNil(t, bobAcct)
	assert.Equal(t, deposit.Add(topUp...), bankKeeper.GetAllBalances(input.Ctx, bob))

	// ensure contract has updated balance
	contractAcct = accKeeper.GetAccount(ctx, addr)
	require.NotNil(t, contractAcct)
	assert.Equal(t, sdk.Coins{}, bankKeeper.GetAllBalances(input.Ctx, addr))

	t.Logf("Duration: %v (35619 gas)\n", diff)
}

func TestExecuteWithNonExistingContractAddress(t *testing.T) {
	input := CreateTestInput(t)
	ctx, accKeeper, bankKeeper, keeper := input.Ctx, input.AccKeeper, input.BankKeeper, input.WasmKeeper

	deposit := sdk.NewCoins(sdk.NewInt64Coin(core.MicroLunaDenom, 100000))
	creator := createFakeFundedAccount(ctx, accKeeper, bankKeeper, deposit.Add(deposit...))

	// unauthorized - trialCtx so we don't change state
	nonExistingContractAddress := types.GenerateContractAddress(9999, 9999)
	_, err := keeper.ExecuteContract(ctx, nonExistingContractAddress, creator, []byte(`{}`), nil)
	require.Error(t, err, sdkerrors.Wrapf(types.ErrNotFound, "contract %s", nonExistingContractAddress))
}

func TestExecuteWithHugeMsg(t *testing.T) {
	input := CreateTestInput(t)
	ctx, accKeeper, bankKeeper, keeper := input.Ctx, input.AccKeeper, input.BankKeeper, input.WasmKeeper

	deposit := sdk.NewCoins(sdk.NewInt64Coin(core.MicroLunaDenom, 100000))
	topUp := sdk.NewCoins(sdk.NewInt64Coin(core.MicroLunaDenom, 5000))
	creator := createFakeFundedAccount(ctx, accKeeper, bankKeeper, deposit.Add(deposit...))
	fred := createFakeFundedAccount(ctx, accKeeper, bankKeeper, topUp)

	wasmCode, err := ioutil.ReadFile("./testdata/hackatom.wasm")
	require.NoError(t, err)

	codeID, err := keeper.StoreCode(ctx, creator, wasmCode)
	require.NoError(t, err)

	_, _, bob := keyPubAddr()
	initMsg := HackatomExampleInitMsg{
		Verifier:    fred,
		Beneficiary: bob,
	}
	initMsgBz, err := json.Marshal(initMsg)
	require.NoError(t, err)

	addr, _, err := keeper.InstantiateContract(ctx, codeID, creator, sdk.AccAddress{}, initMsgBz, deposit)
	require.NoError(t, err)
	require.Equal(t, "cosmos18vd8fpwxzck93qlwghaj6arh4p7c5n89uzcee5", addr.String())

	msgBz := make([]byte, keeper.MaxContractMsgSize(ctx)+1)
	_, err = keeper.ExecuteContract(ctx, addr, fred, msgBz, topUp)
	require.Error(t, err)
	require.Contains(t, err.Error(), "execute msg size is too huge")
}

func TestExecuteWithPanic(t *testing.T) {
	input := CreateTestInput(t)
	ctx, accKeeper, bankKeeper, keeper := input.Ctx, input.AccKeeper, input.BankKeeper, input.WasmKeeper

	deposit := sdk.NewCoins(sdk.NewInt64Coin(core.MicroLunaDenom, 100000))
	topUp := sdk.NewCoins(sdk.NewInt64Coin(core.MicroLunaDenom, 5000))
	creator := createFakeFundedAccount(ctx, accKeeper, bankKeeper, deposit.Add(deposit...))
	fred := createFakeFundedAccount(ctx, accKeeper, bankKeeper, topUp)

	wasmCode, err := ioutil.ReadFile("./testdata/hackatom.wasm")
	require.NoError(t, err)

	contractID, err := keeper.StoreCode(ctx, creator, wasmCode)
	require.NoError(t, err)

	_, _, bob := keyPubAddr()
	initMsg := HackatomExampleInitMsg{
		Verifier:    fred,
		Beneficiary: bob,
	}
	initMsgBz, err := json.Marshal(initMsg)
	require.NoError(t, err)

	addr, _, err := keeper.InstantiateContract(ctx, contractID, creator, sdk.AccAddress{}, initMsgBz, deposit)
	require.NoError(t, err)

	// let's make sure we get a reasonable error, no panic/crash
	_, err = keeper.ExecuteContract(ctx, addr, fred, []byte(`{"panic":{}}`), topUp)
	require.Error(t, err)
}

func TestExecuteWithCpuLoop(t *testing.T) {
	input := CreateTestInput(t)
	ctx, accKeeper, bankKeeper, keeper := input.Ctx, input.AccKeeper, input.BankKeeper, input.WasmKeeper

	deposit := sdk.NewCoins(sdk.NewInt64Coin(core.MicroLunaDenom, 100000))
	topUp := sdk.NewCoins(sdk.NewInt64Coin(core.MicroLunaDenom, 5000))
	creator := createFakeFundedAccount(ctx, accKeeper, bankKeeper, deposit.Add(deposit...))
	fred := createFakeFundedAccount(ctx, accKeeper, bankKeeper, topUp)

	wasmCode, err := ioutil.ReadFile("./testdata/hackatom.wasm")
	require.NoError(t, err)

	contractID, err := keeper.StoreCode(ctx, creator, wasmCode)
	require.NoError(t, err)

	_, _, bob := keyPubAddr()
	initMsg := HackatomExampleInitMsg{
		Verifier:    fred,
		Beneficiary: bob,
	}
	initMsgBz, err := json.Marshal(initMsg)
	require.NoError(t, err)

	addr, _, err := keeper.InstantiateContract(ctx, contractID, creator, sdk.AccAddress{}, initMsgBz, deposit)
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
	input := CreateTestInput(t)
	ctx, accKeeper, bankKeeper, keeper := input.Ctx, input.AccKeeper, input.BankKeeper, input.WasmKeeper

	deposit := sdk.NewCoins(sdk.NewInt64Coin(core.MicroLunaDenom, 100000))
	topUp := sdk.NewCoins(sdk.NewInt64Coin(core.MicroLunaDenom, 5000))
	creator := createFakeFundedAccount(ctx, accKeeper, bankKeeper, deposit.Add(deposit...))
	fred := createFakeFundedAccount(ctx, accKeeper, bankKeeper, topUp)

	wasmCode, err := ioutil.ReadFile("./testdata/hackatom.wasm")
	require.NoError(t, err)

	contractID, err := keeper.StoreCode(ctx, creator, wasmCode)
	require.NoError(t, err)

	_, _, bob := keyPubAddr()
	initMsg := HackatomExampleInitMsg{
		Verifier:    fred,
		Beneficiary: bob,
	}
	initMsgBz, err := json.Marshal(initMsg)
	require.NoError(t, err)

	addr, _, err := keeper.InstantiateContract(ctx, contractID, creator, sdk.AccAddress{}, initMsgBz, deposit)
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
	input := CreateTestInput(t)
	ctx, accKeeper, bankKeeper, keeper := input.Ctx, input.AccKeeper, input.BankKeeper, input.WasmKeeper

	deposit := sdk.NewCoins(sdk.NewInt64Coin("denom", 100000))
	creator := createFakeFundedAccount(ctx, accKeeper, bankKeeper, deposit.Add(deposit...))
	fred := createFakeFundedAccount(ctx, accKeeper, bankKeeper, sdk.NewCoins(sdk.NewInt64Coin("denom", 5000)))

	wasmCode, err := ioutil.ReadFile("./testdata/hackatom.wasm")
	require.NoError(t, err)

	ibcWasmCode, err := ioutil.ReadFile("./testdata/ibc_reflect.wasm")
	require.NoError(t, err)

	reflectWasmCode, err := ioutil.ReadFile("./testdata/reflect.wasm")
	require.NoError(t, err)

	originalCodeID, err := keeper.StoreCode(ctx, creator, wasmCode)
	require.NoError(t, err)
	newCodeID, err := keeper.StoreCode(ctx, creator, wasmCode)
	require.NoError(t, err)
	ibcCodeID, err := keeper.StoreCode(ctx, creator, ibcWasmCode)
	require.NoError(t, err)
	reflectCodeID, err := keeper.StoreCode(ctx, creator, reflectWasmCode)
	require.NoError(t, err)

	require.NotEqual(t, originalCodeID, newCodeID)
	require.NotEqual(t, originalCodeID, ibcCodeID)
	require.NotEqual(t, newCodeID, ibcCodeID)

	_, _, anyAddr := keyPubAddr()
	_, _, newVerifierAddr := keyPubAddr()
	initMsg := HackatomExampleInitMsg{
		Verifier:    fred,
		Beneficiary: anyAddr,
	}
	initMsgBz, err := json.Marshal(initMsg)
	require.NoError(t, err)

	ibcInitMsg := IBCReflectInitMsg{ReflectCodeID: reflectCodeID}
	ibcInitMsgBz, err := json.Marshal(ibcInitMsg)
	require.NoError(t, err)

	migMsg := struct {
		Verifier sdk.AccAddress `json:"verifier"`
	}{Verifier: newVerifierAddr}
	migMsgBz, err := json.Marshal(migMsg)
	require.NoError(t, err)

	specs := map[string]struct {
		admin                sdk.AccAddress
		overrideContractAddr sdk.AccAddress
		caller               sdk.AccAddress
		fromCodeID           uint64
		toCodeID             uint64
		migrateMsg           []byte
		expErr               *sdkerrors.Error
		expVerifier          sdk.AccAddress
		expIBCPort           bool
		initMsg              []byte
	}{
		"all good with same code id": {
			admin:       creator,
			caller:      creator,
			initMsg:     initMsgBz,
			fromCodeID:  originalCodeID,
			toCodeID:    originalCodeID,
			migrateMsg:  migMsgBz,
			expVerifier: newVerifierAddr,
		},
		"all good with different code id": {
			admin:       creator,
			caller:      creator,
			initMsg:     initMsgBz,
			fromCodeID:  originalCodeID,
			toCodeID:    newCodeID,
			migrateMsg:  migMsgBz,
			expVerifier: newVerifierAddr,
		},
		"all good with admin set": {
			admin:       fred,
			caller:      fred,
			initMsg:     initMsgBz,
			fromCodeID:  originalCodeID,
			toCodeID:    newCodeID,
			migrateMsg:  migMsgBz,
			expVerifier: newVerifierAddr,
		},
		"adds IBC port for IBC enabled contracts": {
			admin:       fred,
			caller:      fred,
			initMsg:     initMsgBz,
			fromCodeID:  originalCodeID,
			toCodeID:    ibcCodeID,
			migrateMsg:  []byte(`{}`),
			expIBCPort:  true,
			expVerifier: fred, // not updated
		},
		"prevent migration when admin was not set on instantiate": {
			caller:     creator,
			initMsg:    initMsgBz,
			fromCodeID: originalCodeID,
			toCodeID:   originalCodeID,
			expErr:     sdkerrors.ErrUnauthorized,
		},
		"prevent migration when not sent by admin": {
			caller:     creator,
			admin:      fred,
			initMsg:    initMsgBz,
			fromCodeID: originalCodeID,
			toCodeID:   originalCodeID,
			expErr:     sdkerrors.ErrUnauthorized,
		},
		"fail with non existing code id": {
			admin:      creator,
			caller:     creator,
			initMsg:    initMsgBz,
			fromCodeID: originalCodeID,
			toCodeID:   99999,
			expErr:     sdkerrors.ErrInvalidRequest,
		},
		"fail with non existing contract addr": {
			admin:                creator,
			caller:               creator,
			initMsg:              initMsgBz,
			overrideContractAddr: anyAddr,
			fromCodeID:           originalCodeID,
			toCodeID:             originalCodeID,
			expErr:               sdkerrors.ErrInvalidRequest,
		},
		"fail in contract with invalid migrate msg": {
			admin:      creator,
			caller:     creator,
			initMsg:    initMsgBz,
			fromCodeID: originalCodeID,
			toCodeID:   originalCodeID,
			migrateMsg: bytes.Repeat([]byte{0x1}, 7),
			expErr:     types.ErrMigrationFailed,
		},
		"fail in contract without migrate msg": {
			admin:      creator,
			caller:     creator,
			initMsg:    initMsgBz,
			fromCodeID: originalCodeID,
			toCodeID:   originalCodeID,
			expErr:     types.ErrMigrationFailed,
		},
		"fail when no IBC callbacks": {
			admin:      fred,
			caller:     fred,
			initMsg:    ibcInitMsgBz,
			fromCodeID: ibcCodeID,
			toCodeID:   newCodeID,
			migrateMsg: migMsgBz,
			expErr:     types.ErrMigrationFailed,
		},
	}

	for msg, spec := range specs {
		t.Run(msg, func(t *testing.T) {
			ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)
			addr, _, err := keeper.InstantiateContract(ctx, spec.fromCodeID, creator, spec.admin, spec.initMsg, nil)
			require.NoError(t, err)
			if spec.overrideContractAddr != nil {
				addr = spec.overrideContractAddr
			}
			_, err = keeper.MigrateContract(ctx, addr, spec.caller, spec.toCodeID, spec.migrateMsg)
			require.True(t, spec.expErr.Is(err), "expected %v but got %+v", spec.expErr, err)
			if spec.expErr != nil {
				return
			}
			cInfo, err := keeper.GetContractInfo(ctx, addr)
			require.NoError(t, err)
			assert.Equal(t, spec.toCodeID, cInfo.CodeID)
			assert.Equal(t, spec.expIBCPort, cInfo.IBCPortID != "", cInfo.IBCPortID)

			m := keeper.queryToStore(ctx, addr, []byte("config"))
			var stored map[string]string
			require.NoError(t, json.Unmarshal(m, &stored))
			require.Contains(t, stored, "verifier")
			require.NoError(t, err)
			assert.Equal(t, spec.expVerifier.String(), stored["verifier"])
		})
	}
}

func TestMigrateWithDispatchedMessage(t *testing.T) {
	input := CreateTestInput(t)
	ctx, accKeeper, bankKeeper, keeper := input.Ctx, input.AccKeeper, input.BankKeeper, input.WasmKeeper

	deposit := sdk.NewCoins(sdk.NewInt64Coin(core.MicroLunaDenom, 100000))
	creator := createFakeFundedAccount(ctx, accKeeper, bankKeeper, deposit.Add(deposit...))
	fred := createFakeFundedAccount(ctx, accKeeper, bankKeeper, sdk.NewCoins(sdk.NewInt64Coin(core.MicroLunaDenom, 5000)))

	wasmCode, err := ioutil.ReadFile("./testdata/hackatom.wasm")
	require.NoError(t, err)
	burnerCode, err := ioutil.ReadFile("./testdata/burner.wasm")
	require.NoError(t, err)

	originalContractID, err := keeper.StoreCode(ctx, creator, wasmCode)
	require.NoError(t, err)
	burnerContractID, err := keeper.StoreCode(ctx, creator, burnerCode)
	require.NoError(t, err)
	require.NotEqual(t, originalContractID, burnerContractID)

	_, _, myPayoutAddr := keyPubAddr()
	initMsg := HackatomExampleInitMsg{
		Verifier:    fred,
		Beneficiary: fred,
	}
	initMsgBz, err := json.Marshal(initMsg)
	require.NoError(t, err)

	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)
	contractAddr, _, err := keeper.InstantiateContract(ctx, originalContractID, creator, creator, initMsgBz, deposit)
	require.NoError(t, err)

	migMsg := struct {
		Payout sdk.AccAddress `json:"payout"`
	}{Payout: myPayoutAddr}
	migMsgBz, err := json.Marshal(migMsg)
	require.NoError(t, err)
	ctx = ctx.WithEventManager(sdk.NewEventManager()).WithBlockHeight(ctx.BlockHeight() + 1)
	res, err := keeper.MigrateContract(ctx, contractAddr, creator, burnerContractID, migMsgBz)
	require.NoError(t, err)
	assert.Equal(t, "burnt 1 keys", string(res))
	type dict map[string]interface{}
	expEvents := []dict{
		{
			"Type": "wasm",
			"Attr": []dict{
				{"contract_address": contractAddr},
				{"action": "burn"},
				{"payout": myPayoutAddr},
			},
		},
		{
			"Type": "from_contract",
			"Attr": []dict{
				{"contract_address": contractAddr},
				{"action": "burn"},
				{"payout": myPayoutAddr},
			},
		},
		{
			"Type": "coin_spent",
			"Attr": []dict{
				{"spender": contractAddr},
				{"amount": "100000" + core.MicroLunaDenom},
			},
		},
		{
			"Type": "coin_received",
			"Attr": []dict{
				{"receiver": myPayoutAddr},
				{"amount": "100000" + core.MicroLunaDenom},
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
	assert.Equal(t, deposit, bankKeeper.GetAllBalances(ctx, myPayoutAddr))
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
