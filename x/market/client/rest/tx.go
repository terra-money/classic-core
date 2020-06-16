package rest

import (
	"net/http"

	"github.com/terra-project/core/x/market/internal/types"

	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
)

func registerTxRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc("/market/swap", submitSwapHandlerFn(cliCtx)).Methods("POST")
}

// SwapReq defines request body for swap operation
type SwapReq struct {
	BaseReq   rest.BaseReq   `json:"base_req"`
	OfferCoin sdk.Coin       `json:"offer_coin"`
	AskDenom  string         `json:"ask_denom"`
	Receiver  sdk.AccAddress `json:"receiver,omitempty"`
}

// submitSwapHandlerFn handles a POST vote request
func submitSwapHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req SwapReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		fromAddress, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		toAddress := req.Receiver
		if toAddress.Empty() {
			toAddress = fromAddress
		}

		// create the message
		msg := types.NewMsgSwapSend(fromAddress, toAddress, req.OfferCoin, req.AskDenom)
		err = msg.ValidateBasic()
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}
