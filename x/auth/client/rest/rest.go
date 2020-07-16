package rest

import (
	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/context"
)

// RegisterRoutes registers the auth module REST routes
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc("/txs", QueryTxsRequestHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc("/auth/accounts/{address}/multisign", MultiSignRequestHandlerFn(cliCtx)).Methods("POST")
	r.HandleFunc("/txs/estimate_fee", EstimateTxFeeRequestHandlerFn(cliCtx)).Methods("POST")
}
