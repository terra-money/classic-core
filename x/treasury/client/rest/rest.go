package rest

import (
	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client"
	govrest "github.com/cosmos/cosmos-sdk/x/gov/client/rest"
)

// Defines wildcard part of the request paths
const (
	RestDenom = "denom"
	RestEpoch = "epoch"
)

// RegisterRoutes registers oracle-related REST handlers to a router
func RegisterRoutes(clientCtx client.Context, r *mux.Router) {
	registerQueryRoute(clientCtx, r)
}

// TaxRateUpdateProposalRESTHandler returns a ProposalRESTHandler that exposes the community pool spend REST handler with a given sub-route.
func TaxRateUpdateProposalRESTHandler(clientCtx client.Context) govrest.ProposalRESTHandler {
	return govrest.ProposalRESTHandler{
		SubRoute: "tax_rate_update",
		Handler:  postTaxRateUpdateProposalHandlerFn(clientCtx),
	}
}

// RewardWeightUpdateProposalRESTHandler returns a ProposalRESTHandler that exposes the community pool spend REST handler with a given sub-route.
func RewardWeightUpdateProposalRESTHandler(clientCtx client.Context) govrest.ProposalRESTHandler {
	return govrest.ProposalRESTHandler{
		SubRoute: "reward_weight_update",
		Handler:  postRewardWeightUpdateProposalHandlerFn(clientCtx),
	}
}
