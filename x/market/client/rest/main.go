package rest

import (
	"net/http"
	"terra/x/market"
	"terra/x/market/client"

	"github.com/cosmos/cosmos-sdk/types/rest"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/utils"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gorilla/mux"
)

// RegisterRoutes - Central function to define routes that get registered by the main application
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *codec.Codec) {
	r.HandleFunc("market/swap", submitSwapHandlerFunction(cdc, cliCtx)).Methods("POST")

	r.HandleFunc("market/history", queryHistoryHandlerFunction(cdc, cliCtx)).Methods("GET")
}

//nolint
type SwapReq struct {
	BaseReq       utils.BaseReq  `json:"base_req"`
	OfferCoin     sdk.Coin       `json:"offer_coin"`
	AskDenom      string         `json:"ask_denom"`
	TraderAddress sdk.AccAddress `json:"trader_address"`
}

//nolint
type HistoryReq struct {
	BaseReq       utils.BaseReq  `json:"base_req"`
	OfferDenom    string         `json:"offer_denom"`
	AskDenom      string         `json:"ask_denom"`
	TraderAddress sdk.AccAddress `json:"trader_address"`
}

// submitSwapHandlerFunction handles a POST vote request
func submitSwapHandlerFunction(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		swapReq := SwapReq{}
		err := rest.ReadRESTReq(w, r, cdc, &swapReq)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		baseReq := swapReq.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w, cliCtx) {
			return
		}

		// create the message
		msg := market.NewSwapMsg(swapReq.TraderAddress, swapReq.OfferCoin, swapReq.AskDenom)
		err = msg.ValidateBasic()
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		rest.CompleteAndBroadcastTxREST(w, r, cliCtx, baseReq, []sdk.Msg{msg}, cdc)
	}
}

func queryHistoryHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		histReq := HistoryReq{}
		err := rest.ReadRESTReq(w, r, cdc, &histReq)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		params := client.QueryHistoryParams{
			TraderAddress: histReq.TraderAddress,
			AskDenom:      histReq.AskDenom,
			OfferDenom:    histReq.OfferDenom,
		}

		res, err := client.QueryHistoryByTxQuery(cdc, cliCtx, params)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		rest.PostProcessResponse(w, cdc, res, cliCtx.Indent)
	}
}
