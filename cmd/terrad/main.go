package main

import (
	"encoding/json"
	"io"
	"os"
	"terra/version"

	"github.com/cosmos/cosmos-sdk/baseapp"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/cli"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
	tmtypes "github.com/tendermint/tendermint/types"

	"terra/app"
	terraInit "terra/cmd/init"

	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func main() {
	cdc := app.MakeCodec()

	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(sdk.Bech32PrefixAccAddr, sdk.Bech32PrefixAccPub)
	config.SetBech32PrefixForValidator(sdk.Bech32PrefixValAddr, sdk.Bech32PrefixValPub)
	config.SetBech32PrefixForConsensusNode(sdk.Bech32PrefixConsAddr, sdk.Bech32PrefixConsPub)
	config.Seal()

	ctx := server.NewDefaultContext()
	cobra.EnableCommandSorting = false
	rootCmd := &cobra.Command{
		Use:               "terrad",
		Short:             "Terra Daemon (server)",
		PersistentPreRunE: server.PersistentPreRunEFn(ctx),
	}
	rootCmd.AddCommand(terraInit.InitCmd(ctx, cdc))
	rootCmd.AddCommand(terraInit.CollectGenTxsCmd(ctx, cdc))
	rootCmd.AddCommand(terraInit.TestnetFilesCmd(ctx, cdc))
	rootCmd.AddCommand(terraInit.GenTxCmd(ctx, cdc))
	rootCmd.AddCommand(terraInit.AddGenesisAccountCmd(ctx, cdc))

	// preoccupy the version command (that will be added in server.AddCommands) @matthew
	rootCmd.AddCommand(version.VersionCmd)
	server.AddCommands(ctx, cdc, rootCmd, appInit,
		newApp, exportAppStateAndTMValidators)

	// prepare and add flags
	executor := cli.PrepareBaseCmd(rootCmd, "TE", app.DefaultNodeHome)
	err := executor.Execute()
	if err != nil {
		// handle with #870
		panic(err)
	}
}

func newApp(logger log.Logger, db dbm.DB, traceStore io.Writer) abci.Application {
	return app.NewTerraApp(logger, db, traceStore,
		baseapp.SetPruning(viper.GetString("pruning")),
		baseapp.SetMinimumFees(viper.GetString("minimum_fees")),
	)
}

func exportAppStateAndTMValidators(
	logger log.Logger, db dbm.DB, traceStore io.Writer, height int64, forZeroHeight bool,
) (json.RawMessage, []tmtypes.GenesisValidator, error) {
	tApp := app.NewTerraApp(logger, db, traceStore)
	if height != -1 {
		err := tApp.LoadHeight(height)
		if err != nil {
			return nil, nil, err
		}
	}
	return tApp.ExportAppStateAndValidators(forZeroHeight)
}
