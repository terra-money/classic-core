package main

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/cli"
	"github.com/tendermint/tendermint/libs/log"
	tmtypes "github.com/tendermint/tendermint/types"
	dbm "github.com/tendermint/tm-db"

	"github.com/terra-money/core/app"
	core "github.com/terra-money/core/types"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client/debug"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	genutilcli "github.com/cosmos/cosmos-sdk/x/genutil/client/cli"

	"github.com/terra-money/core/x/auth"
	coreante "github.com/terra-money/core/x/auth/ante"
	"github.com/terra-money/core/x/staking"
	wasmconfig "github.com/terra-money/core/x/wasm/config"
)

// terrad custom flags
const flagInvCheckPeriod = "inv-check-period"

var invCheckPeriod uint

func main() {
	cdc := app.MakeCodec()

	config := sdk.GetConfig()
	config.SetCoinType(core.CoinType)
	config.SetFullFundraiserPath(core.FullFundraiserPath)
	config.SetBech32PrefixForAccount(core.Bech32PrefixAccAddr, core.Bech32PrefixAccPub)
	config.SetBech32PrefixForValidator(core.Bech32PrefixValAddr, core.Bech32PrefixValPub)
	config.SetBech32PrefixForConsensusNode(core.Bech32PrefixConsAddr, core.Bech32PrefixConsPub)
	config.Seal()

	ctx := server.NewDefaultContext()
	cobra.EnableCommandSorting = false
	rootCmd := &cobra.Command{
		Use:               "terrad",
		Short:             "Terra Daemon (server)",
		PersistentPreRunE: persistentPreRunEFn(ctx),
	}

	rootCmd.AddCommand(genutilcli.InitCmd(ctx, cdc, app.ModuleBasics, app.DefaultNodeHome))
	rootCmd.AddCommand(genutilcli.CollectGenTxsCmd(ctx, cdc, auth.GenesisAccountIterator{}, app.DefaultNodeHome))
	rootCmd.AddCommand(genutilcli.GenTxCmd(ctx, cdc, app.ModuleBasics, staking.AppModuleBasic{},
		auth.GenesisAccountIterator{}, app.DefaultNodeHome, app.DefaultCLIHome))
	rootCmd.AddCommand(genutilcli.ValidateGenesisCmd(ctx, cdc, app.ModuleBasics))
	rootCmd.AddCommand(AddGenesisAccountCmd(ctx, cdc, app.DefaultNodeHome, app.DefaultCLIHome))
	rootCmd.AddCommand(flags.NewCompletionCmd(rootCmd, true))
	rootCmd.AddCommand(testnetCmd(ctx, cdc, app.ModuleBasics, auth.GenesisAccountIterator{}))
	rootCmd.AddCommand(replayCmd())
	rootCmd.AddCommand(debug.Cmd(cdc))

	server.AddCommands(ctx, cdc, rootCmd, newApp, exportAppStateAndTMValidators)

	// prepare and add flags
	executor := cli.PrepareBaseCmd(rootCmd, "TE", app.DefaultNodeHome)
	rootCmd.PersistentFlags().UintVar(&invCheckPeriod, flagInvCheckPeriod,
		0, "Assert registered invariants every N blocks")

	// register tx gas hard cap flag
	rootCmd.PersistentFlags().Uint64(coreante.FlagTxGasHardLimit, uint64(30000000),
		"Transaction hard cap to prevent spamming attack")
	viper.BindPFlag(coreante.FlagTxGasHardLimit, rootCmd.Flags().Lookup(coreante.FlagTxGasHardLimit))

	err := executor.Execute()
	if err != nil {
		panic(err)
	}
}

func newApp(logger log.Logger, db dbm.DB, traceStore io.Writer) abci.Application {
	var cache sdk.MultiStorePersistentCache

	if viper.GetBool(server.FlagInterBlockCache) {
		cache = store.NewCommitKVStoreCacheManager()
	}

	skipUpgradeHeights := make(map[int64]bool)
	for _, h := range viper.GetIntSlice(server.FlagUnsafeSkipUpgrades) {
		skipUpgradeHeights[int64(h)] = true
	}

	pruningOpts, err := server.GetPruningOptionsFromFlags()
	if err != nil {
		panic(err)
	}

	return app.NewTerraApp(
		logger, db, traceStore, true, invCheckPeriod, skipUpgradeHeights,
		&wasmconfig.Config{BaseConfig: wasmconfig.BaseConfig{
			ContractQueryGasLimit:    viper.GetUint64(wasmconfig.FlagContractQueryGasLimit),
			ContractLoggingWhitelist: viper.GetString(wasmconfig.FlagContractLoggingWhitelist),
		}},
		baseapp.SetPruning(pruningOpts),
		baseapp.SetMinGasPrices(viper.GetString(server.FlagMinGasPrices)),
		baseapp.SetHaltHeight(viper.GetUint64(server.FlagHaltHeight)),
		baseapp.SetHaltTime(viper.GetUint64(server.FlagHaltTime)),
		baseapp.SetInterBlockCache(cache),
	)
}

func exportAppStateAndTMValidators(
	logger log.Logger, db dbm.DB, traceStore io.Writer, height int64, forZeroHeight bool, jailWhiteList []string,
) (json.RawMessage, []tmtypes.GenesisValidator, error) {

	if height != -1 {
		tApp := app.NewTerraApp(logger, db, traceStore, false, uint(1), map[int64]bool{}, wasmconfig.DefaultConfig())
		err := tApp.LoadHeight(height)
		if err != nil {
			return nil, nil, err
		}
		return tApp.ExportAppStateAndValidators(forZeroHeight, jailWhiteList)
	}
	tApp := app.NewTerraApp(logger, db, traceStore, true, uint(1), map[int64]bool{}, wasmconfig.DefaultConfig())
	return tApp.ExportAppStateAndValidators(forZeroHeight, jailWhiteList)
}

// custom pre-run to generate wasm config
func persistentPreRunEFn(context *server.Context) func(*cobra.Command, []string) error {
	originPreRun := server.PersistentPreRunEFn(context)
	return func(cmd *cobra.Command, args []string) error {
		if cmd.Name() == version.Cmd.Name() {
			return nil
		}

		err := originPreRun(cmd, args)
		if err != nil {
			return err
		}

		rootDir := viper.GetString(flags.FlagHome)

		// load application DBDir and set wasm DBDir
		wasmconfig.DBDir = filepath.Base(context.Config.DBDir()) + "/wasm"
		wasmConfigFilePath := filepath.Join(rootDir, "config/wasm.toml")
		if _, err := os.Stat(wasmConfigFilePath); os.IsNotExist(err) {
			wasmConf, _ := wasmconfig.ParseConfig()
			wasmconfig.WriteConfigFile(wasmConfigFilePath, wasmConf)
		}

		viper.SetConfigName("wasm")
		err = viper.MergeInConfig()

		return err
	}
}
