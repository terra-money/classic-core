package client

import (
	"terra/x/market/client/cli"

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
		Short: "Querying commands for the market module",
	}
	marketQueryCmd.AddCommand(client.GetCommands(
	// cli.GetCmdQueryDelegation(mc.storeKey, mc.cdc),
	// cli.GetCmdQueryDelegations(mc.storeKey, mc.cdc),
	// cli.GetCmdQueryUnbondingDelegation(mc.storeKey, mc.cdc),
	// cli.GetCmdQueryUnbondingDelegations(mc.storeKey, mc.cdc),
	// cli.GetCmdQueryRedelegation(mc.storeKey, mc.cdc),
	// cli.GetCmdQueryRedelegations(mc.storeKey, mc.cdc),
	// cli.GetCmdQueryValidator(mc.storeKey, mc.cdc),
	// cli.GetCmdQueryValidators(mc.storeKey, mc.cdc),
	// cli.GetCmdQueryValidatorDelegations(mc.storeKey, mc.cdc),
	// cli.GetCmdQueryValidatorUnbondingDelegations(mc.storeKey, mc.cdc),
	// cli.GetCmdQueryValidatorRedelegations(mc.storeKey, mc.cdc),
	// cli.GetCmdQueryParams(mc.storeKey, mc.cdc),
	// cli.GetCmdQueryPool(mc.storeKey, mc.cdc)
	)...)

	return marketQueryCmd

}

// GetTxCmd returns the transaction commands for this module
func (mc ModuleClient) GetTxCmd() *cobra.Command {
	marketTxCmd := &cobra.Command{
		Use:   "market",
		Short: "Staking transaction subcommands",
	}

	marketTxCmd.AddCommand(client.PostCommands(
		cli.GetSwapCmd(mc.cdc),
	)...)

	return marketTxCmd
}
