package rest

import (
	"fmt"
	"net/http"
	"github.com/terra-project/core/types/assets"
	"github.com/terra-project/core/x/oracle"

	"github.com/cosmos/cosmos-sdk/client/context"
	clientrest "github.com/cosmos/cosmos-sdk/client/rest"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"

	"github.com/gorilla/mux"
)

func resgisterTxRoute(cliCtx context.CLIContext, r *mux.Router, cdc *codec.Codec) {
	r.HandleFunc(fmt.Sprintf("/oracle/denoms/{%s}/votes", RestDenom), submitVoteHandlerFunction(cdc, cliCtx)).Methods("POST")
}

//VoteReq ...
type VoteReq struct {
	BaseReq rest.BaseReq `json:"base_req"`
	Price   sdk.Dec      `json:"price"`
}

func submitVoteHandlerFunction(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		denom := vars[RestDenom]

		if !assets.IsValidDenom(denom) {
			err := fmt.Errorf("The denom is not known: %s", denom)
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		var req VoteReq
		if !rest.ReadRESTReq(w, r, cdc, &req) {
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

		// create the message
		msg := oracle.NewMsgPriceFeed(denom, req.Price, fromAddress)
		err = msg.ValidateBasic()
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		if req.BaseReq.GenerateOnly {
			clientrest.WriteGenerateStdTxResponse(w, cdc, cliCtx, req.BaseReq, []sdk.Msg{msg})
			return
		}

		clientrest.CompleteAndBroadcastTxREST(w, cliCtx, req.BaseReq, []sdk.Msg{msg}, cdc)
	}
}
