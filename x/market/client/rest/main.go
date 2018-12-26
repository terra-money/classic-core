package rest

import (
	"fmt"
	"net/http"
	"terra/x/market"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/utils"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gorilla/mux"
)

// REST Variable names
// nolint
const (
	RestOfferCoin = "offerCoin"
	RestAskDenom  = "askDenom"
	storeName     = "market"
)

// RegisterRoutes - Central function to define routes that get registered by the main application
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *codec.Codec) {
	// POST /vote/{denom}
	r.HandleFunc(fmt.Sprintf("market/swap/{%s}/ask/{%s}", RestOfferCoin, RestAskDenom),
		SubmitSwapHandlerFunction(cdc, cliCtx)).Methods("POST")
}

//nolint
type SwapReq struct {
	BaseReq        utils.BaseReq `json:"base_req"`
	OfferDenom     string        `json:"offer_denom"`
	OfferAmount    sdk.Int       `json:"offer_amount"`
	AskDenom       string        `json:"ask_denom"`
	SwapperAddress string        `json:"swapper_address"`
}

// SubmitSwapHandlerFunction handles a POST vote request
func SubmitSwapHandlerFunction(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		swapReq := SwapReq{}
		err := utils.ReadRESTReq(w, r, cdc, &swapReq)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		baseReq := swapReq.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w, cliCtx) {
			return
		}

		swapAddr, err := cliCtx.GetFromAddress()
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		offerCoin := sdk.NewCoin(swapReq.OfferDenom, swapReq.OfferAmount)
		askDenom := swapReq.AskDenom

		// create the message
		msg := market.NewSwapMsg(swapAddr, offerCoin, askDenom)
		err = msg.ValidateBasic()
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.CompleteAndBroadcastTxREST(w, r, cliCtx, baseReq, []sdk.Msg{msg}, cdc)
	}
}
