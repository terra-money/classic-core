package rest

import (
	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client"
)

// RegisterTxRoutes registers registers terra custom transaction routes on the provided router.
func RegisterTxRoutes(clientCtx client.Context, rtr *mux.Router) {
	rtr.HandleFunc("/txs/estimate_fee", EstimateTxFeeRequestHandlerFn(clientCtx)).Methods("POST")
}
