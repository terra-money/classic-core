package rest

import (
	"net/http"
	"terra/x/market"

	clientrest "github.com/cosmos/cosmos-sdk/client/rest"
	"github.com/cosmos/cosmos-sdk/types/rest"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gorilla/mux"
)

// RegisterRoutes - Central function to define routes that get registered by the main application
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *codec.Codec) {
	r.HandleFunc("market/swap", submitSwapHandlerFn(cdc, cliCtx)).Methods("POST")
}

//nolint
type SwapReq struct {
	BaseReq       rest.BaseReq   `json:"base_req"`
	OfferCoin     sdk.Coin       `json:"offer_coin"`
	AskDenom      string         `json:"ask_denom"`
	TraderAddress sdk.AccAddress `json:"trader_address"`
}

// submitSwapHandlerFn handles a POST vote request
func submitSwapHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var swapReq SwapReq
		if !rest.ReadRESTReq(w, r, cdc, &swapReq) {
			err := sdk.ErrUnknownRequest("malformed request")
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		baseReq := swapReq.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			err := sdk.ErrUnknownRequest("malformed request")
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// create the message
		msg := market.NewMsgSwap(swapReq.TraderAddress, swapReq.OfferCoin, swapReq.AskDenom)
		err := msg.ValidateBasic()
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		clientrest.CompleteAndBroadcastTxREST(w, cliCtx, swapReq.BaseReq, []sdk.Msg{msg}, cdc)
	}
}
