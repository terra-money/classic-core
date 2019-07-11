package main

import (
	"encoding/json"
	"io"

	"github.com/terra-project/core/types/util"
	"github.com/terra-project/core/version"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/store"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/cli"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/terra-project/core/app"
	terraInit "github.com/terra-project/core/cmd/init"

	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	terraserver "github.com/terra-project/core/server"
)

const flagAssertInvariantsBlockly = "assert-invariants-blockly"

var assertInvariantsBlockly bool

func main() {
	cdc := app.MakeCodec()

	config := sdk.GetConfig()
	config.SetCoinType(util.CoinType)
	config.SetFullFundraiserPath(util.FullFundraiserPath)
	config.SetBech32PrefixForAccount(util.Bech32PrefixAccAddr, util.Bech32PrefixAccPub)
	config.SetBech32PrefixForValidator(util.Bech32PrefixValAddr, util.Bech32PrefixValPub)
	config.SetBech32PrefixForConsensusNode(util.Bech32PrefixConsAddr, util.Bech32PrefixConsPub)
	config.Seal()

	ctx := server.NewDefaultContext()
	cobra.EnableCommandSorting = false
	rootCmd := &cobra.Command{
		Use:               "terrad",
		Short:             "Terra Daemon (server)",
		PersistentPreRunE: terraserver.PersistentPreRunEFn(ctx),
	}
	rootCmd.AddCommand(terraInit.InitCmd(ctx, cdc))
	rootCmd.AddCommand(terraInit.CollectGenTxsCmd(ctx, cdc))
	rootCmd.AddCommand(terraInit.TestnetFilesCmd(ctx, cdc))
	rootCmd.AddCommand(terraInit.GenTxCmd(ctx, cdc))
	rootCmd.AddCommand(terraInit.AddGenesisAccountCmd(ctx, cdc))
	rootCmd.AddCommand(terraInit.ValidateGenesisCmd(ctx, cdc))
	rootCmd.AddCommand(client.NewCompletionCmd(rootCmd, true))

	// preempting version command
	rootCmd.AddCommand(version.VersionCmd)

	server.AddCommands(ctx, cdc, rootCmd, newApp, exportAppStateAndTMValidators)

	// prepare and add flags
	executor := cli.PrepareBaseCmd(rootCmd, "TE", app.DefaultNodeHome)
	rootCmd.PersistentFlags().BoolVar(&assertInvariantsBlockly, flagAssertInvariantsBlockly,
		false, "Assert registered invariants on a blockly basis")
	err := executor.Execute()
	if err != nil {
		// handle with #870
		panic(err)
	}
}

func newApp(logger log.Logger, db dbm.DB, traceStore io.Writer) abci.Application {
	return app.NewTerraApp(
		logger, db, traceStore, true, assertInvariantsBlockly,
		baseapp.SetPruning(store.NewPruningOptionsFromString(viper.GetString("pruning"))),
		baseapp.SetMinGasPrices(viper.GetString(server.FlagMinGasPrices)),
	)
}

func exportAppStateAndTMValidators(
	logger log.Logger, db dbm.DB, traceStore io.Writer, height int64, forZeroHeight bool, jailWhiteList []string,
) (json.RawMessage, []tmtypes.GenesisValidator, error) {
	if height != -1 {
		tApp := app.NewTerraApp(logger, db, traceStore, false, false)
		err := tApp.LoadHeight(height)
		if err != nil {
			return nil, nil, err
		}
		return tApp.ExportAppStateAndValidators(forZeroHeight, jailWhiteList)
	}
	tApp := app.NewTerraApp(logger, db, traceStore, true, false)
	return tApp.ExportAppStateAndValidators(forZeroHeight, jailWhiteList)
}
