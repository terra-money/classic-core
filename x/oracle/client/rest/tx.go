package rest

import (
	"bytes"
	"math"
	"net/http"
	"terra/x/oracle"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/utils"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keys"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gorilla/mux"
)

func registerTxRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *codec.Codec, kb keys.Keybase) {
	// POST /vote/{denom}
	r.HandleFunc(
		"/oracle/vote/{denom}",
		submitVoteHandlerFunction(cdc, kb, cliCtx),
	).Methods("POST")
}

func submitVoteHandlerFunction(cdc *codec.Codec, kb keys.Keybase, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req VoteReq
		err := utils.ReadRESTReq(w, r, cdc, &req)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w, cliCtx) {
			return
		}

		info, err := kb.Get(baseReq.Name)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusUnauthorized, err.Error())
			return
		}

		if !bytes.Equal(info.GetPubKey().Address(), []byte(req.VoterAddress)) {
			utils.WriteErrorResponse(w, http.StatusUnauthorized, "Must use own address")
			return
		}

		target := sdk.NewDecWithPrec(int64(math.Round(req.TargetPrice*1000)), 3)
		current := sdk.NewDecWithPrec(int64(math.Round(req.CurrentPrice*1000)), 3)

		// create the message
		msg := oracle.NewPriceFeedMsg(req.Denom, target, current, info.GetAddress())
		err = msg.ValidateBasic()
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.CompleteAndBroadcastTxREST(w, r, cliCtx, baseReq, []sdk.Msg{msg}, cdc)
	}
}
