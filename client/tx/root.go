package tx

import (
	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdktx "github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/codec"
)

// register REST routes
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *codec.Codec) {
	r.HandleFunc("/txs/{hash}", sdktx.QueryTxRequestHandlerFn(cdc, cliCtx)).Methods("GET")
	r.HandleFunc("/txs", QueryTxsByTagsRequestHandlerFn(cliCtx, cdc)).Methods("GET")
	r.HandleFunc("/txs", sdktx.BroadcastTxRequest(cliCtx, cdc)).Methods("POST")
	r.HandleFunc("/txs/encode", sdktx.EncodeTxRequestHandlerFn(cdc, cliCtx)).Methods("POST")
	r.HandleFunc("/txs/estimate_fee", EstimateTxFeeRequestHandlerFn(cdc, cliCtx)).Methods("POST")
}
