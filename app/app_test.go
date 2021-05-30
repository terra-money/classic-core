package app

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"

	wasmconfig "github.com/terra-money/core/x/wasm/config"
)

func TestTerraExport(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "wasmtest")
	require.NoError(t, err)
	viper.Set(flags.FlagHome, tempDir)
	defer os.RemoveAll(tempDir)

	db := dbm.NewMemDB()
	tapp := NewTerraApp(log.NewTMLogger(log.NewSyncWriter(os.Stdout)), db, nil, true, 0, map[int64]bool{}, wasmconfig.DefaultConfig())
	err = setGenesis(tapp)
	require.NoError(t, err)

	// Making a new app object with the db, so that initchain hasn't been called
	newTapp := NewTerraApp(log.NewTMLogger(log.NewSyncWriter(os.Stdout)), db, nil, true, 0, map[int64]bool{}, wasmconfig.DefaultConfig())
	_, _, err = newTapp.ExportAppStateAndValidators(false, []string{})
	require.NoError(t, err, "ExportAppStateAndValidators should not have an error")

	_, _, err = newTapp.ExportAppStateAndValidators(true, []string{})
	require.NoError(t, err, "ExportAppStateAndValidators for zero height should not have an error")
}

func setGenesis(tapp *TerraApp) error {
	genesisState := ModuleBasics.DefaultGenesis()
	stateBytes, err := codec.MarshalJSONIndent(tapp.Codec(), genesisState)
	if err != nil {
		return err
	}

	// Initialize the chain
	tapp.InitChain(
		abci.RequestInitChain{
			Validators:    []abci.ValidatorUpdate{},
			AppStateBytes: stateBytes,
		},
	)
	tapp.Commit()
	return nil
}

// ensure that black listed addresses are properly set in bank keeper
func TestBlackListedAddrs(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "wasmtest")
	require.NoError(t, err)
	viper.Set(flags.FlagHome, tempDir)
	defer os.RemoveAll(tempDir)

	db := dbm.NewMemDB()
	app := NewTerraApp(log.NewTMLogger(log.NewSyncWriter(os.Stdout)), db, nil, true, 0, map[int64]bool{}, wasmconfig.DefaultConfig())

	for acc := range maccPerms {
		require.Equal(t, !allowedReceivingModAcc[acc], app.bankKeeper.BlacklistedAddr(app.supplyKeeper.GetModuleAddress(acc)))
	}
}
