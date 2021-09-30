package rest

import (
	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client"
	clientrest "github.com/cosmos/cosmos-sdk/client/rest"
)

// RegisterTxRoutes registers registers terra custom transaction routes on the provided router.
func RegisterTxRoutes(clientCtx client.Context, rtr *mux.Router) {
	r := clientrest.WithHTTPDeprecationHeaders(rtr)
	r.HandleFunc("/txs/estimate_fee", EstimateTxFeeRequestHandlerFn(clientCtx)).Methods("POST")
	r.HandleFunc("/txs/encode", EncodeTxRequest(clientCtx)).Methods("POST")
	r.HandleFunc("/txs", BroadcastTxRequest(clientCtx)).Methods("POST")
}
