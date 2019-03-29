package rest

import (
	"fmt"
	"net/http"

	"terra/x/budget"
	"terra/x/oracle"

	"github.com/cosmos/cosmos-sdk/client/context"
	clientrest "github.com/cosmos/cosmos-sdk/client/rest"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

// REST Variable names
// nolint
const (
	RestParamsType    = "type"
	RestProgramID     = "program-id"
	RestVoter         = "voter"
	RestProgramStatus = "status"
	RestNumLimit      = "limit"

	queryRoute = "budget"
)

// RegisterRoutes - Central function to define routes that get registered by the main application
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *codec.Codec) {
	r.HandleFunc("/budget/program/submit", submitProgramHandlerFn(cdc, cliCtx)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/budget/program/{%s}/withdraw", RestProgramID), voteHandlerFn(cdc, cliCtx)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/budget/program/{%s}/vote", RestProgramID), voteHandlerFn(cdc, cliCtx)).Methods("POST")

	r.HandleFunc(fmt.Sprintf("/budget/program"), queryProgramHandlerFn(cdc, cliCtx)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/budget/program/{%s}", RestProgramID), queryProgramHandlerFn(cdc, cliCtx)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/budget/program/{%s}/actives", RestProgramID), queryActivesHandlerFn(cdc, cliCtx)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/budget/program/{%s}/candidates", RestProgramID), queryCandidatesHandlerFn(cdc, cliCtx)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/budget/program/{%s}/votes/{%s}", RestProgramID, RestVoter), queryVotesHandlerFn(cdc, cliCtx)).Methods("GET")
	r.HandleFunc("/budget/params", queryParamsHandlerFn(cdc, cliCtx)).Methods("GET")
}

type submitProgramReq struct {
	BaseReq     rest.BaseReq   `json:"base_req"`
	Title       string         `json:"title"`       //  Title of the Program
	Description string         `json:"description"` //  Description of the Program
	Submitter   sdk.AccAddress `json:"submitter"`   //  Address of the submitter
	Executor    sdk.AccAddress `json:"executor"`    //  Address of the executor
	Deposit     sdk.Coin       `json:"deposit"`     // Coins to add to the Program's deposit
}

type voteReq struct {
	BaseReq rest.BaseReq   `json:"base_req"`
	Voter   sdk.AccAddress `json:"voter"`  //  address of the voter
	Option  bool           `json:"option"` //  option from OptionSet chosen by the voter
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

		fromAddress := cliCtx.GetFromAddress()
		if !req.Submitter.Equals(fromAddress) {
			rest.WriteErrorResponse(w, http.StatusUnauthorized, "Must use own address")
			return
		}

		// create the message
		msg := budget.NewMsgSubmitProgram(req.Title, req.Description, req.Submitter, req.Executor)
		err := msg.ValidateBasic()
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

		fromAddress := cliCtx.GetFromAddress()

		// create the message
		msg := budget.NewMsgWithdrawProgram(programID, fromAddress)
		err := msg.ValidateBasic()
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

		var req voteReq
		if !rest.ReadRESTReq(w, r, cdc, &req) {
			return
		}

		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		// create the message
		msg := budget.NewMsgVoteProgram(programID, req.Option, req.Voter)
		err := msg.ValidateBasic()
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		clientrest.CompleteAndBroadcastTxREST(w, cliCtx, req.BaseReq, []sdk.Msg{msg}, cdc)
	}
}

func queryProgramHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		strProgramID := vars[RestProgramID]

		if len(strProgramID) == 0 {
			err := errors.New("programID required but not specified")
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s", queryRoute, budget.QueryProgram, strProgramID), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cdc, res, cliCtx.Indent)
	}
}

func queryActivesHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", queryRoute, budget.QueryActiveList), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cdc, res, cliCtx.Indent)
	}
}

func queryCandidatesHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", queryRoute, budget.QueryCandidateList), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cdc, res, cliCtx.Indent)
	}
}

// todo: Split this functionality into helper functions to remove the above
func queryVotesHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		strProgramID := vars[RestProgramID]
		strVoterAddr := vars[RestVoter]

		params := budget.QueryVotesParams{}

		if len(strProgramID) != 0 {
			programID, ok := rest.ParseUint64OrReturnBadRequest(w, strProgramID)
			if !ok {
				err := errors.New("programID must be an unsigned int")
				rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}

			params.ProgramID = programID
		}

		if len(strVoterAddr) != 0 {
			voterAcc, err := sdk.AccAddressFromBech32(strVoterAddr)
			if err != nil {
				err := errors.New("voteraccount malformed")
				rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}

			params.Voter = voterAcc
		}

		bz, err := cdc.MarshalJSON(params)
		if err != nil {
			return
		}

		res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", queryRoute, budget.QueryVotes), bz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		rest.PostProcessResponse(w, cdc, res, cliCtx.Indent)
	}
}

func queryParamsHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", queryRoute, oracle.QueryParams), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cdc, res, cliCtx.Indent)
	}
}
