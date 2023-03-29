package feeshare_test

import (
	"encoding/json"

	wasmconfig "github.com/classic-terra/core/x/wasm/config"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	dbm "github.com/tendermint/tm-db"

	app "github.com/classic-terra/core/app"
	"github.com/cosmos/cosmos-sdk/simapp"
	"github.com/cosmos/cosmos-sdk/x/mint/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// returns context and an app with updated mint keeper
func CreateTestApp(isCheckTx bool) (*app.TerraApp, sdk.Context) {
	app := Setup(isCheckTx)

	ctx := app.BaseApp.NewContext(isCheckTx, tmproto.Header{})
	app.MintKeeper.SetParams(ctx, types.DefaultParams())
	app.MintKeeper.SetMinter(ctx, types.DefaultInitialMinter())

	return app, ctx
}

func Setup(isCheckTx bool) *app.TerraApp {
	app, genesisState := GenApp(!isCheckTx, 5)
	if !isCheckTx {
		// init chain must be called to stop deliverState from being nil
		stateBytes, err := json.MarshalIndent(genesisState, "", " ")
		if err != nil {
			panic(err)
		}

		// Initialize the chain
		app.InitChain(
			abci.RequestInitChain{
				Validators:      []abci.ValidatorUpdate{},
				ConsensusParams: simapp.DefaultConsensusParams,
				AppStateBytes:   stateBytes,
			},
		)
	}

	return app
}

func GenApp(withGenesis bool, invCheckPeriod uint) (*app.TerraApp, app.GenesisState) {
	db := dbm.NewMemDB()
	encCdc := app.MakeEncodingConfig()
	terraapp := app.NewTerraApp(
		log.NewNopLogger(),
		db,
		nil,
		true,
		map[int64]bool{},
		app.DefaultNodeHome,
		invCheckPeriod,
		encCdc,
		simapp.EmptyAppOptions{},
		wasmconfig.DefaultConfig(),
	)

	if withGenesis {
		return terraapp, app.NewDefaultGenesisState()
	}

	return terraapp, app.GenesisState{}
}
