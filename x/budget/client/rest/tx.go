package rest

import (
	"fmt"
	"github.com/terra-project/core/x/budget"
	"net/http"
	"strings"

	"github.com/cosmos/cosmos-sdk/client/context"
	clientrest "github.com/cosmos/cosmos-sdk/client/rest"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

func registerTxRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *codec.Codec) {
	r.HandleFunc("/budget/programs/submit", submitProgramHandlerFn(cdc, cliCtx)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/budget/programs/{%s}/withdraw", RestProgramID), withdrawProgramHandlerFn(cdc, cliCtx)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/budget/programs/{%s}/votes", RestProgramID), voteHandlerFn(cdc, cliCtx)).Methods("POST")
}

type submitProgramReq struct {
	BaseReq     rest.BaseReq   `json:"base_req"`
	Title       string         `json:"title"`       //  Title of the Program
	Description string         `json:"description"` //  Description of the Program
	Executor    sdk.AccAddress `json:"executor"`    //  Address of the executor
}

type voteReq struct {
	BaseReq rest.BaseReq `json:"base_req"`
	Option  bool         `json:"option"` //  option from OptionSet chosen by the voter
}

type withdrawReq struct {
	BaseReq rest.BaseReq `json:"base_req"`
}

func submitProgramHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req submitProgramReq
		if !rest.ReadRESTReq(w, r, cdc, &req) {
			return
		}

		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
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

		// Query params to get deposit amount
		res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", budget.QuerierRoute, budget.QueryParams), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		var params budget.Params
		cdc.MustUnmarshalJSON(res, &params)

		if fromAccount.GetCoins().AmountOf(params.Deposit.Denom).LT(params.Deposit.Amount) {
			err := fmt.Errorf(strings.TrimSpace(`
                              account %s has insufficient amount of coins to pay the offered coins.\n
                              Required: %s\n
                              Given:    %s\n`), fromAddress, params.Deposit, fromAccount.GetCoins())

			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		}

		// create the message
		msg := budget.NewMsgSubmitProgram(req.Title, req.Description, fromAddress, req.Executor)
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

func withdrawProgramHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		strProgramID := vars[RestProgramID]

		if len(strProgramID) == 0 {
			err := errors.New("programID required but not specified")
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		programID, ok := rest.ParseUint64OrReturnBadRequest(w, strProgramID)
		if !ok {
			return
		}

		var req withdrawReq
		if !rest.ReadRESTReq(w, r, cdc, &req) {
			return
		}

		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		fromAddress, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// create the message
		msg := budget.NewMsgWithdrawProgram(programID, fromAddress)
		err = msg.ValidateBasic()
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		clientrest.CompleteAndBroadcastTxREST(w, cliCtx, req.BaseReq, []sdk.Msg{msg}, cdc)
	}
}

func voteHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		programIDStr := vars[RestProgramID]

		if len(programIDStr) == 0 {
			err := errors.New("programID required but not specified")
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		programID, ok := rest.ParseUint64OrReturnBadRequest(w, programIDStr)
		if !ok {
			return
		}

		var req voteReq
		if !rest.ReadRESTReq(w, r, cdc, &req) {
			return
		}

		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		fromAddress, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// create the message
		msg := budget.NewMsgVoteProgram(programID, req.Option, fromAddress)
		err = msg.ValidateBasic()
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		clientrest.CompleteAndBroadcastTxREST(w, cliCtx, req.BaseReq, []sdk.Msg{msg}, cdc)
	}
}
