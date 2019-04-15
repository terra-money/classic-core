package client

import (
	"github.com/terra-project/core/x/budget/client/cli"

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
	budgetQueryCmd := &cobra.Command{
		Use:   "budget",
		Short: "Querying commands for the budget module",
	}
	budgetQueryCmd.AddCommand(client.GetCommands(
		cli.GetCmdQueryProgram(mc.cdc),
		cli.GetCmdQueryActives(mc.cdc),
		cli.GetCmdQueryCandidates(mc.cdc),
		cli.GetCmdQueryVotes(mc.cdc),
		cli.GetCmdQueryParams(mc.cdc),
	)...)

	return budgetQueryCmd
}

// GetTxCmd returns the transaction commands for this module
func (mc ModuleClient) GetTxCmd() *cobra.Command {
	budgetTxCmd := &cobra.Command{
		Use:   "budget",
		Short: "budget transaction subcommands",
	}

	budgetTxCmd.AddCommand(client.PostCommands(
		cli.GetCmdSubmitProgram(mc.cdc),
		cli.GetCmdWithdrawProgram(mc.cdc),
		cli.GetCmdVote(mc.cdc),
	)...)

	return budgetTxCmd
}
