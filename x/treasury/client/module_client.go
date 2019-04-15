package client

import (
	"github.com/spf13/cobra"
	"github.com/tendermint/go-amino"

	treasuryCli "github.com/terra-project/core/x/treasury/client/cli"

	"github.com/cosmos/cosmos-sdk/client"
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
	// Group treasury queries under a subcommand
	treasuryQueryCmd := &cobra.Command{
		Use:   "treasury",
		Short: "Querying commands for the treasury module",
	}

	treasuryQueryCmd.AddCommand(client.GetCommands(
		treasuryCli.GetCmdQueryTaxRate(mc.cdc),
		treasuryCli.GetCmdQueryTaxCap(mc.cdc),
		treasuryCli.GetCmdQueryMiningRewardWeight(mc.cdc),
		treasuryCli.GetCmdQueryIssuance(mc.cdc),
		treasuryCli.GetCmdQueryTaxProceeds(mc.cdc),
		treasuryCli.GetCmdQuerySeigniorageProceeds(mc.cdc),
		treasuryCli.GetCmdQueryActiveClaims(mc.cdc),
		treasuryCli.GetCmdQueryCurrentEpoch(mc.cdc),
		treasuryCli.GetCmdQueryParams(mc.cdc),
	)...)

	return treasuryQueryCmd
}

// GetTxCmd The treasury module returns no TX commands.
func (mc ModuleClient) GetTxCmd() *cobra.Command {
	return &cobra.Command{Hidden: true}
}
