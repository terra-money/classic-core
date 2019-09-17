package rest

import (
	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/context"
)

// RegisterRoutes registers the auth module REST routes
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc("/auth/accounts/{address}/multisign", MultiSignRequestHandlerFn(cliCtx)).Methods("POST")
}

// RegisterTxRoutes registers custom transaction routes on the provided router.
func RegisterTxRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc("/txs/estimate_fee", EstimateTxFeeRequestHandlerFn(cliCtx)).Methods("POST")
}
