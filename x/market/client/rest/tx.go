package rest

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"

	"github.com/terra-money/core/x/market/types"
)

func registerTxHandlers(clientCtx client.Context, rtr *mux.Router) {
	rtr.HandleFunc("/market/swap", submitSwapHandlerFn(clientCtx)).Methods("POST")
}

type (
	swapReq struct {
		BaseReq   rest.BaseReq `json:"base_req"`
		OfferCoin sdk.Coin     `json:"offer_coin"`
		AskDenom  string       `json:"ask_denom"`
		Receiver  string       `json:"receiver,omitempty"`
	}
)

// submitSwapHandlerFn handles a POST vote request
func submitSwapHandlerFn(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req swapReq
		if !rest.ReadRESTReq(w, r, clientCtx.LegacyAmino, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		fromAddress, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if rest.CheckBadRequestError(w, err) {
			return
		}

		// create the message depends on the toAddress existence
		var msg sdk.Msg
		if req.Receiver == "" {
			msg = types.NewMsgSwap(fromAddress, req.OfferCoin, req.AskDenom)
			if rest.CheckBadRequestError(w, msg.ValidateBasic()) {
				return
			}
		} else {
			toAddress, err := sdk.AccAddressFromBech32(req.Receiver)
			if rest.CheckBadRequestError(w, err) {
				return
			}

			msg := types.NewMsgSwapSend(fromAddress, toAddress, req.OfferCoin, req.AskDenom)
			if rest.CheckBadRequestError(w, msg.ValidateBasic()) {
				return
			}
		}

		tx.WriteGeneratedTxResponse(clientCtx, w, req.BaseReq, msg)
	}
}
