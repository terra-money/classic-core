package testutil

import (
	"bytes"

	"github.com/spf13/cobra"

	"github.com/tendermint/tendermint/libs/cli"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/terra-project/core/app"
	"github.com/terra-project/core/types/util"
)

// PrepareTest prepares codec, rootCmd, txCmd, queryCmd instances
func PrepareCmdTest() (cdc *codec.Codec, rootCmd *cobra.Command, txCmd *cobra.Command, queryCmd *cobra.Command) {
	cdc = app.MakeCodec()

	config := sdk.GetConfig()

	if config.GetBech32AccountAddrPrefix() != util.Bech32PrefixAccAddr {
		config.SetCoinType(330)
		config.SetFullFundraiserPath("44'/330'/0'/0/0")
		config.SetBech32PrefixForAccount(util.Bech32PrefixAccAddr, util.Bech32PrefixAccPub)
		config.SetBech32PrefixForValidator(util.Bech32PrefixValAddr, util.Bech32PrefixValPub)
		config.SetBech32PrefixForConsensusNode(util.Bech32PrefixConsAddr, util.Bech32PrefixConsPub)
		config.Seal()
	}

	rootCmd = &cobra.Command{
		Use:   "terracli",
		Short: "Command line interface for interacting with terrad",
	}

	txCmd = &cobra.Command{
		Use:   "tx",
		Short: "Transactions subcommands",
	}

	queryCmd = &cobra.Command{
		Use:     "query",
		Aliases: []string{"q"},
		Short:   "Querying subcommands",
	}

	rootCmd.PersistentFlags().String(client.FlagChainID, "", "Chain ID of tendermint node")
	rootCmd.AddCommand(txCmd)
	rootCmd.AddCommand(queryCmd)

	return
}

// ExecuteCommand executes command
func ExecuteCommand(rootCmd *cobra.Command, args ...string) (output string, err error) {
	output, err = executeCommandC(rootCmd, args...)
	return output, err
}

func executeCommandC(rootCmd *cobra.Command, args ...string) (output string, err error) {
	buf := new(bytes.Buffer)
	rootCmd.SetOutput(buf)
	rootCmd.SetArgs(args)

	executor := cli.PrepareMainCmd(rootCmd, "TE", app.DefaultCLIHome)
	err = executor.Execute()

	return buf.String(), err
}
