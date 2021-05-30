package client

import (
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
	"github.com/terra-money/core/x/treasury/client/cli"
	"github.com/terra-money/core/x/treasury/client/rest"
)

// param change proposal handler
var (
	TaxRateUpdateProposalHandler      = govclient.NewProposalHandler(cli.GetCmdSubmitTaxRateUpdateProposal, rest.TaxRateUpdateProposalRESTHandler)
	RewardWeightUpdateProposalHandler = govclient.NewProposalHandler(cli.GetCmdSubmitRewardWeightUpdateProposal, rest.RewardWeightUpdateProposalRESTHandler)
)
