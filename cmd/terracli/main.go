package main

import (
	"fmt"
	"net/http"
	"os"
	"path"

	"github.com/rakyll/statik/fs"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/libs/cli"

	"github.com/terra-project/core/app"
	"github.com/terra-project/core/types/util"
	"github.com/terra-project/core/version"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/client/lcd"
	"github.com/cosmos/cosmos-sdk/client/rpc"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	at "github.com/cosmos/cosmos-sdk/x/auth"

	txcustom "github.com/terra-project/core/client/tx"

	authcustom "github.com/terra-project/core/x/auth/client/rest"
	dist "github.com/terra-project/core/x/distribution/client/rest"
	slashing "github.com/terra-project/core/x/slashing/client/rest"
	staking "github.com/terra-project/core/x/staking/client/rest"

	budget "github.com/terra-project/core/x/budget/client/rest"
	market "github.com/terra-project/core/x/market/client/rest"
	oracle "github.com/terra-project/core/x/oracle/client/rest"
	pay "github.com/terra-project/core/x/pay/client/rest"
	treasury "github.com/terra-project/core/x/treasury/client/rest"

	bud "github.com/terra-project/core/x/budget"
	mkt "github.com/terra-project/core/x/market"
	ora "github.com/terra-project/core/x/oracle"
	tre "github.com/terra-project/core/x/treasury"

	dt "github.com/cosmos/cosmos-sdk/x/distribution"
	sl "github.com/cosmos/cosmos-sdk/x/slashing"
	st "github.com/cosmos/cosmos-sdk/x/staking"

	auth "github.com/cosmos/cosmos-sdk/x/auth/client/rest"

	authcmd "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
	paycmd "github.com/terra-project/core/x/pay/client/cli"

	budgetClient "github.com/terra-project/core/x/budget/client"
	distClient "github.com/terra-project/core/x/distribution/client"
	marketClient "github.com/terra-project/core/x/market/client"
	oracleClient "github.com/terra-project/core/x/oracle/client"
	slashingClient "github.com/terra-project/core/x/slashing/client"
	stakingClient "github.com/terra-project/core/x/staking/client"
	treasuryClient "github.com/terra-project/core/x/treasury/client"

	crisisClient "github.com/cosmos/cosmos-sdk/x/crisis/client"

	_ "github.com/terra-project/core/client/lcd/statik"
)

func main() {
	// Configure cobra to sort commands
	cobra.EnableCommandSorting = false

	// Instantiate the codec for the command line application
	cdc := app.MakeCodec()

	// Read in the configuration file for the sdk
	config := sdk.GetConfig()
	config.SetCoinType(330)
	config.SetFullFundraiserPath("44'/330'/0'/0/0")
	config.SetBech32PrefixForAccount(util.Bech32PrefixAccAddr, util.Bech32PrefixAccPub)
	config.SetBech32PrefixForValidator(util.Bech32PrefixValAddr, util.Bech32PrefixValPub)
	config.SetBech32PrefixForConsensusNode(util.Bech32PrefixConsAddr, util.Bech32PrefixConsPub)
	config.Seal()

	// TODO: setup keybase, viper object, etc. to be passed into
	// the below functions and eliminate global vars, like we do
	// with the cdc

	// Module clients hold cli commnads (tx,query) and lcd routes
	// TODO: Make the lcd command take a list of ModuleClient
	mc := []sdk.ModuleClients{
		distClient.NewModuleClient(dt.StoreKey, cdc),
		stakingClient.NewModuleClient(st.StoreKey, cdc),
		slashingClient.NewModuleClient(sl.StoreKey, cdc),
		oracleClient.NewModuleClient(ora.StoreKey, cdc),
		treasuryClient.NewModuleClient(tre.StoreKey, cdc),
		budgetClient.NewModuleClient(bud.StoreKey, cdc),
		marketClient.NewModuleClient(mkt.StoreKey, cdc),
		crisisClient.NewModuleClient(sl.StoreKey, cdc),
	}

	rootCmd := &cobra.Command{
		Use:   "terracli",
		Short: "Command line interface for interacting with terrad",
	}

	// Add --chain-id to persistent flags and mark it required
	rootCmd.PersistentFlags().String(client.FlagChainID, "", "Chain ID of tendermint node")
	rootCmd.PersistentPreRunE = func(_ *cobra.Command, _ []string) error {
		return initConfig(rootCmd)
	}

	// Construct Root Command
	rootCmd.AddCommand(
		rpc.StatusCommand(),
		client.ConfigCmd(app.DefaultCLIHome),
		queryCmd(cdc, mc),
		txCmd(cdc, mc),
		client.LineBreak,
		lcd.ServeCommand(cdc, registerRoutes),
		client.LineBreak,
		keys.Commands(),
		client.LineBreak,
		version.VersionCmd,
		client.NewCompletionCmd(rootCmd, true),
	)

	// Add flags and prefix all env exposed with GA
	executor := cli.PrepareMainCmd(rootCmd, "TE", app.DefaultCLIHome)

	err := executor.Execute()
	if err != nil {
		fmt.Printf("Failed executing CLI command: %s, exiting...\n", err)
		os.Exit(1)
	}
}

func queryCmd(cdc *amino.Codec, mc []sdk.ModuleClients) *cobra.Command {
	queryCmd := &cobra.Command{
		Use:     "query",
		Aliases: []string{"q"},
		Short:   "Querying subcommands",
	}

	queryCmd.AddCommand(
		rpc.ValidatorCommand(cdc),
		rpc.BlockCommand(),
		tx.SearchTxCmd(cdc),
		tx.QueryTxCmd(cdc),
		client.LineBreak,
		authcmd.GetAccountCmd(at.StoreKey, cdc),
	)

	for _, m := range mc {
		mQueryCmd := m.GetQueryCmd()
		if mQueryCmd != nil {
			queryCmd.AddCommand(mQueryCmd)
		}
	}

	return queryCmd
}

func txCmd(cdc *amino.Codec, mc []sdk.ModuleClients) *cobra.Command {
	txCmd := &cobra.Command{
		Use:   "tx",
		Short: "Transactions subcommands",
	}

	txCmd.AddCommand(
		paycmd.SendTxCmd(cdc),
		client.LineBreak,
		authcmd.GetSignCommand(cdc),
		authcmd.GetMultiSignCommand(cdc),
		tx.GetBroadcastCommand(cdc),
		tx.GetEncodeCommand(cdc),
		client.LineBreak,
	)

	for _, m := range mc {
		txCmd.AddCommand(m.GetTxCmd())
	}

	return txCmd
}

// CLIVersionRequestHandler cli version REST handler endpoint
func CLIVersionRequestHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(fmt.Sprintf("{\"version\": \"%s\"}", version.Version)))
}

// registerRoutes registers the routes from the different modules for the LCD.
// NOTE: details on the routes added for each module are in the module documentation
// NOTE: If making updates here you also need to update the test helper in client/lcd/test_helper.go
func registerRoutes(rs *lcd.RestServer) {

	rs.Mux.HandleFunc("/version", CLIVersionRequestHandler).Methods("GET")

	registerSwaggerUI(rs)
	rpc.RegisterRoutes(rs.CliCtx, rs.Mux)
	txcustom.RegisterRoutes(rs.CliCtx, rs.Mux, rs.Cdc)
	auth.RegisterRoutes(rs.CliCtx, rs.Mux, rs.Cdc, at.StoreKey)
	dist.RegisterRoutes(rs.CliCtx, rs.Mux, rs.Cdc, dt.StoreKey)
	staking.RegisterRoutes(rs.CliCtx, rs.Mux, rs.Cdc, rs.KeyBase)
	slashing.RegisterRoutes(rs.CliCtx, rs.Mux, rs.Cdc, rs.KeyBase)

	authcustom.RegisterRoutes(rs.CliCtx, rs.Mux, rs.Cdc, rs.KeyBase)
	pay.RegisterRoutes(rs.CliCtx, rs.Mux, rs.Cdc, rs.KeyBase)
	oracle.RegisterRoutes(rs.CliCtx, rs.Mux, rs.Cdc)
	treasury.RegisterRoutes(rs.CliCtx, rs.Mux, rs.Cdc)
	market.RegisterRoutes(rs.CliCtx, rs.Mux, rs.Cdc)
	budget.RegisterRoutes(rs.CliCtx, rs.Mux, rs.Cdc)
}

func registerSwaggerUI(rs *lcd.RestServer) {
	statikFS, err := fs.New()
	if err != nil {
		panic(err)
	}
	staticServer := http.FileServer(statikFS)
	rs.Mux.PathPrefix("/swagger-ui/").Handler(http.StripPrefix("/swagger-ui/", staticServer))
}

func initConfig(cmd *cobra.Command) error {
	home, err := cmd.PersistentFlags().GetString(cli.HomeFlag)
	if err != nil {
		return err
	}

	cfgFile := path.Join(home, "config", "config.toml")
	if _, err := os.Stat(cfgFile); err == nil {
		viper.SetConfigFile(cfgFile)

		if err := viper.ReadInConfig(); err != nil {
			return err
		}
	}
	if err := viper.BindPFlag(client.FlagChainID, cmd.PersistentFlags().Lookup(client.FlagChainID)); err != nil {
		return err
	}
	if err := viper.BindPFlag(cli.EncodingFlag, cmd.PersistentFlags().Lookup(cli.EncodingFlag)); err != nil {
		return err
	}
	return viper.BindPFlag(cli.OutputFlag, cmd.PersistentFlags().Lookup(cli.OutputFlag))
}
