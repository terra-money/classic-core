package client

import (
	"net/http"

	"github.com/classic-terra/core/v2/x/treasury/client/cli"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/types/rest"
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
	govrest "github.com/cosmos/cosmos-sdk/x/gov/client/rest"
)

// should we support legacy rest?
// general direction of the hub seems to be moving away from legacy rest
var (
	ProposalAddBurnTaxExemptionAddressHandler    = govclient.NewProposalHandler(cli.ProposalAddBurnTaxExemptionAddressCmd, emptyRestHandler)
	ProposalRemoveBurnTaxExemptionAddressHandler = govclient.NewProposalHandler(cli.ProposalRemoveBurnTaxExemptionAddressCmd, emptyRestHandler)
)

func emptyRestHandler(client.Context) govrest.ProposalRESTHandler {
	return govrest.ProposalRESTHandler{
		SubRoute: "unsupported-service",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "Legacy REST Routes are not supported for tax exemption address proposals")
		},
	}
}
