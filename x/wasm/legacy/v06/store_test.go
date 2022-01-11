package v06_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/terra-money/core/app"
	v05 "github.com/terra-money/core/x/wasm/legacy/v05"
	v06 "github.com/terra-money/core/x/wasm/legacy/v06"
	"github.com/terra-money/core/x/wasm/types"

	"github.com/cosmos/cosmos-sdk/testutil"
	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestContractInfoMigration(t *testing.T) {
	encodingConfig := app.MakeEncodingConfig()
	wasmKey := sdk.NewKVStoreKey(types.StoreKey)
	ctx := testutil.DefaultContext(wasmKey, sdk.NewTransientStoreKey("transient_test"))
	store := ctx.KVStore(wasmKey)

	oldContractInfo := v05.ContractInfo{
		Address: "addr1",
		Creator: "addr2",
		Admin:   "admin",
		CodeID:  10,
		InitMsg: []byte("{\"empty\":{}}"),
	}

	bz, err := encodingConfig.Marshaler.Marshal(&oldContractInfo)
	require.NoError(t, err)

	_, _, addr := testdata.KeyTestPubAddr()

	store.Set(v05.GetContractInfoKey(addr), bz)

	// Run migration
	err = v06.MigrateStore(ctx, wasmKey, encodingConfig.Marshaler)
	require.NoError(t, err)

	bz = store.Get(v05.GetContractInfoKey(addr))

	var contractInfo types.ContractInfo
	err = encodingConfig.Marshaler.Unmarshal(bz, &contractInfo)
	require.NoError(t, err)

	require.Equal(t, types.ContractInfo{
		Address:   "addr1",
		Creator:   "addr2",
		Admin:     "admin",
		CodeID:    10,
		InitMsg:   []byte("{\"empty\":{}}"),
		IBCPortID: "",
	}, contractInfo)

}
