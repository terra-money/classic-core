package app

import (
	"io"

	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
)

var (
	genesisFile        string
	paramsFile         string
	exportParamsPath   string
	exportParamsHeight int
	exportStatePath    string
	exportStatsPath    string
	seed               int64
	initialBlockHeight int
	numBlocks          int
	blockSize          int
	enabled            bool
	verbose            bool
	lean               bool
	commit             bool
	period             int
	onOperation        bool // TODO Remove in favor of binary search for invariant violation
	allInvariants      bool
	genesisTime        int64
)

// DONTCOVER

// NewTerraAppUNSAFE is used for debugging purposes only.
//
// NOTE: to not use this function with non-test code
func NewTerraAppUNSAFE(logger log.Logger, db dbm.DB, traceStore io.Writer, loadLatest bool,
	invCheckPeriod uint, baseAppOptions ...func(*baseapp.BaseApp),
) (tapp *TerraApp, keyMain, keyStaking *sdk.KVStoreKey, stakingKeeper staking.Keeper) {

	tapp = NewTerraApp(logger, db, traceStore, loadLatest, invCheckPeriod, baseAppOptions...)
	return tapp, tapp.keys[baseapp.MainStoreKey], tapp.keys[staking.StoreKey], tapp.stakingKeeper
}
