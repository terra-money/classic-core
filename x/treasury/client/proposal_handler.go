package client

import (
	"github.com/classic-terra/core/v2/x/treasury/client/cli"
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
)

// should we support legacy rest?
// general direction of the hub seems to be moving away from legacy rest
var (
	ProposalAddBurnTaxExemptionAddressHandler    = govclient.NewProposalHandler(cli.ProposalAddBurnTaxExemptionAddressCmd)
	ProposalRemoveBurnTaxExemptionAddressHandler = govclient.NewProposalHandler(cli.ProposalRemoveBurnTaxExemptionAddressCmd)
)
