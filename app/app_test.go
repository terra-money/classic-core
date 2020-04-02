package app

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"
	"github.com/terra-project/core/x/oracle"

	"github.com/cosmos/cosmos-sdk/codec"

	abci "github.com/tendermint/tendermint/abci/types"
)

func TestTerraExport(t *testing.T) {
	db := dbm.NewMemDB()
	tapp := NewTerraApp(log.NewTMLogger(log.NewSyncWriter(os.Stdout)), db, nil, true, 0)
	err := setGenesis(tapp)
	require.NoError(t, err)

	// Making a new app object with the db, so that initchain hasn't been called
	newTapp := NewTerraApp(log.NewTMLogger(log.NewSyncWriter(os.Stdout)), db, nil, true, 0)
	_, _, err2 := newTapp.ExportAppStateAndValidators(false, []string{})
	require.NoError(t, err2, "ExportAppStateAndValidators should not have an error")

	_, _, err2 = newTapp.ExportAppStateAndValidators(true, []string{})
	require.NoError(t, err2, "ExportAppStateAndValidators for zero height should not have an error")
}

func setGenesis(tapp *TerraApp) error {

	genesisState := ModuleBasics.DefaultGenesis()
	stateBytes, err := codec.MarshalJSONIndent(tapp.cdc, genesisState)
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
	db := dbm.NewMemDB()
	app := NewTerraApp(log.NewTMLogger(log.NewSyncWriter(os.Stdout)), db, nil, true, 0)

	for acc := range maccPerms {
		if acc == oracle.ModuleName {
			require.False(t, app.bankKeeper.BlacklistedAddr(app.supplyKeeper.GetModuleAddress(acc)))
		} else {
			require.True(t, app.bankKeeper.BlacklistedAddr(app.supplyKeeper.GetModuleAddress(acc)))
		}

	}
}
