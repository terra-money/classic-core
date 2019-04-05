package rest

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/terra-project/core/types/assets"
	"github.com/terra-project/core/x/market"

	clientrest "github.com/cosmos/cosmos-sdk/client/rest"

	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func registerTxRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *codec.Codec) {
	r.HandleFunc("/market/swap", submitSwapHandlerFn(cdc, cliCtx)).Methods("POST")
}

//nolint
type SwapReq struct {
	BaseReq   rest.BaseReq `json:"base_req"`
	OfferCoin sdk.Coin     `json:"offer_coin"`
	AskDenom  string       `json:"ask_denom"`
}

// submitSwapHandlerFn handles a POST vote request
func submitSwapHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req SwapReq
		if !rest.ReadRESTReq(w, r, cdc, &req) {
			err := sdk.ErrUnknownRequest("malformed request")
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		if !assets.IsValidDenom(req.AskDenom) {
			err := fmt.Errorf("The denom is not known: %s", req.AskDenom)
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			err := sdk.ErrUnknownRequest("malformed request")
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		fromAddress, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		fromAccount, err := cliCtx.GetAccount(fromAddress)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		if fromAccount.GetCoins().AmountOf(req.OfferCoin.Denom).LT(req.OfferCoin.Amount) {
			err := fmt.Errorf(strings.TrimSpace(`
                              account %s has insufficient amount of coins to pay the offered coins.\n
                              Required: %s\n
                              Given:    %s\n`), fromAddress, req.OfferCoin, fromAccount.GetCoins())

			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// create the message
		msg := market.NewMsgSwap(fromAddress, req.OfferCoin, req.AskDenom)
		err = msg.ValidateBasic()
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		clientrest.WriteGenerateStdTxResponse(w, cdc, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}
