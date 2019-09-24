package rest

import (
	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/context"
	govrest "github.com/cosmos/cosmos-sdk/x/gov/client/rest"
)

const (
	RestDenom = "denom"
	RestEpoch = "epoch"
)

// RegisterRoutes registers oracle-related REST handlers to a router
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router) {
	// resgisterTxRoute(cliCtx, r)
	registerQueryRoute(cliCtx, r)
}

// TaxRateUpdateProposalRESTHandler returns a ProposalRESTHandler that exposes the community pool spend REST handler with a given sub-route.
func TaxRateUpdateProposalRESTHandler(cliCtx context.CLIContext) govrest.ProposalRESTHandler {
	return govrest.ProposalRESTHandler{
		SubRoute: "tax_rate_update",
		Handler:  postTaxRateUpdateProposalHandlerFn(cliCtx),
	}
}

// RewardWeightUpdateProposalRESTHandler returns a ProposalRESTHandler that exposes the community pool spend REST handler with a given sub-route.
func RewardWeightUpdateProposalRESTHandler(cliCtx context.CLIContext) govrest.ProposalRESTHandler {
	return govrest.ProposalRESTHandler{
		SubRoute: "reward_weight_update",
		Handler:  postRewardWeightUpdateProposalHandlerFn(cliCtx),
	}
}
