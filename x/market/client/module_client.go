package client

import (
	"github.com/terra-project/core/x/market/client/cli"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/spf13/cobra"
	amino "github.com/tendermint/go-amino"
)

// ModuleClient exports all client functionality from this module
type ModuleClient struct {
	storeKey string
	cdc      *amino.Codec
}

func NewModuleClient(storeKey string, cdc *amino.Codec) ModuleClient {
	return ModuleClient{storeKey, cdc}
}

// GetQueryCmd returns the cli query commands for this module
func (mc ModuleClient) GetQueryCmd() *cobra.Command {
	marketQueryCmd := &cobra.Command{
		Use:   "market",
		Short: "Market query subcommands",
	}

	marketQueryCmd.AddCommand(client.GetCommands(
		cli.GetCmdQuerySwap(mc.cdc),
		cli.GetCmdQueryParams(mc.cdc),
	)...)

	return marketQueryCmd
}

// GetTxCmd returns the transaction commands for this module
func (mc ModuleClient) GetTxCmd() *cobra.Command {
	marketTxCmd := &cobra.Command{
		Use:   "market",
		Short: "Market transaction subcommands",
	}

	marketTxCmd.AddCommand(client.PostCommands(
		cli.GetSwapCmd(mc.cdc),
	)...)

	return marketTxCmd
}
