package rest

import (
	"fmt"
	"net/http"

	"terra/x/budget"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/utils"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

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
	storeName         = "budget"
)

// RegisterRoutes - Central function to define routes that get registered by the main application
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *codec.Codec) {
	r.HandleFunc("/budget/program", postProgramHandlerFn(cdc, cliCtx)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/budget/program/{%s}/vote", RestProgramID), voteHandlerFn(cdc, cliCtx)).Methods("POST")

	r.HandleFunc(fmt.Sprintf("/budget/program/{%s}", RestProgramID), queryProgramHandlerFn(cdc, cliCtx)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/budget/program/{%s}/tally", RestProgramID), queryTallyOnProgramHandlerFn(cdc, cliCtx)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/budget/program/{%s}/votes", RestProgramID), queryVotesOnProgramHandlerFn(cdc, cliCtx)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/budget/program/{%s}/votes/{%s}", RestProgramID, RestVoter), queryVoteHandlerFn(cdc, cliCtx)).Methods("GET")
}

type postProgramReq struct {
	BaseReq     utils.BaseReq  `json:"base_req"`
	Title       string         `json:"title"`        //  Title of the Program
	Description string         `json:"description"`  //  Description of the Program
	ProgramType string         `json:"Program_type"` //  Type of Program. Initial set {PlainTextProgram, SoftwareUpgradeProgram}
	Submitter   sdk.AccAddress `json:"submitter"`    //  Address of the submitter
	Executor    sdk.AccAddress `json:"executor"`     //  Address of the executor
	Deposit     sdk.Coins      `json:"deposit"`      // Coins to add to the Program's deposit
}

type voteReq struct {
	BaseReq utils.BaseReq  `json:"base_req"`
	Voter   sdk.AccAddress `json:"voter"`  //  address of the voter
	Option  string         `json:"option"` //  option from OptionSet chosen by the voter
}

func postProgramHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req postProgramReq
		err := utils.ReadRESTReq(w, r, cdc, &req)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		cliCtx = cliCtx.WithGenerateOnly(req.BaseReq.GenerateOnly)
		cliCtx = cliCtx.WithSimulation(req.BaseReq.Simulate)

		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w, cliCtx) {
			return
		}

		// create the message
		msg := budget.NewSubmitProgramMsg(req.Title, req.Description, req.Deposit, req.Submitter, req.Executor)
		err = msg.ValidateBasic()
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.CompleteAndBroadcastTxREST(w, r, cliCtx, baseReq, []sdk.Msg{msg}, cdc)
	}
}

func voteHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		strProgramID := vars[RestProgramID]

		if len(strProgramID) == 0 {
			err := errors.New("programID required but not specified")
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		programID, ok := utils.ParseUint64OrReturnBadRequest(w, strProgramID)
		if !ok {
			return
		}

		var req voteReq
		err := utils.ReadRESTReq(w, r, cdc, &req)
		if err != nil {
			return
		}

		cliCtx = cliCtx.WithGenerateOnly(req.BaseReq.GenerateOnly)
		cliCtx = cliCtx.WithSimulation(req.BaseReq.Simulate)

		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w, cliCtx) {
			return
		}

		// TODO: check option and validate
		var voteOption budget.ProgramVote
		switch req.Option {
		case "yes":
			voteOption = budget.YesVote
			break
		case "no":
			voteOption = budget.NoVote
			break
		default:
			voteOption = budget.AbstainVote
			break
		}

		// create the message
		msg := budget.NewVoteMsg(programID, voteOption, req.Voter)
		err = msg.ValidateBasic()
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.CompleteAndBroadcastTxREST(w, r, cliCtx, baseReq, []sdk.Msg{msg}, cdc)
	}
}

func queryProgramHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		strProgramID := vars[RestProgramID]

		if len(strProgramID) == 0 {
			err := errors.New("programID required but not specified")
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		programID, ok := utils.ParseUint64OrReturnBadRequest(w, strProgramID)
		if !ok {
			return
		}

		params := budget.NewQueryProgramParams(programID)

		bz, err := cdc.MarshalJSON(params)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		res, err := cliCtx.QueryWithData("custom/budget/Program", bz)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		utils.PostProcessResponse(w, cdc, res, cliCtx.Indent)
	}
}

func queryVoteHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		strProgramID := vars[RestProgramID]
		bechVoterAddr := vars[RestVoter]

		if len(strProgramID) == 0 {
			err := errors.New("ProgramId required but not specified")
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		ProgramID, ok := utils.ParseUint64OrReturnBadRequest(w, strProgramID)
		if !ok {
			return
		}

		if len(bechVoterAddr) == 0 {
			err := errors.New("voter address required but not specified")
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		voterAddr, err := sdk.AccAddressFromBech32(bechVoterAddr)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		params := budget.NewQueryVoteParams(ProgramID, voterAddr)

		bz, err := cdc.MarshalJSON(params)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		res, err := cliCtx.QueryWithData("custom/budget/vote", bz)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		utils.PostProcessResponse(w, cdc, res, cliCtx.Indent)
	}
}

// todo: Split this functionality into helper functions to remove the above
func queryVotesOnProgramHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		strProgramID := vars[RestProgramID]

		if len(strProgramID) == 0 {
			err := errors.New("ProgramId required but not specified")
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		ProgramID, ok := utils.ParseUint64OrReturnBadRequest(w, strProgramID)
		if !ok {
			return
		}

		params := budget.NewQueryProgramParams(ProgramID)

		bz, err := cdc.MarshalJSON(params)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		res, err := cliCtx.QueryWithData("custom/budget/Program", bz)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		var program budget.Program
		if err := cdc.UnmarshalJSON(res, &program); err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		// For inactive Programs we must query the txs directly to get the votes
		// as they're no longer in state.
		if !(program.State == budget.InactiveProgramState || program.State == budget.ActiveProgramState) {
			//res, err = gcutils.QueryVotesByTxQuery(cdc, cliCtx, params)
		} else {
			res, err = cliCtx.QueryWithData("custom/budget/votes", bz)
		}

		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		utils.PostProcessResponse(w, cdc, res, cliCtx.Indent)
	}
}

// todo: Split this functionality into helper functions to remove the above
func queryTallyOnProgramHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		strProgramID := vars[RestProgramID]

		if len(strProgramID) == 0 {
			err := errors.New("ProgramId required but not specified")
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		ProgramID, ok := utils.ParseUint64OrReturnBadRequest(w, strProgramID)
		if !ok {
			return
		}

		params := budget.NewQueryProgramParams(ProgramID)

		bz, err := cdc.MarshalJSON(params)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		res, err := cliCtx.QueryWithData("custom/budget/tally", bz)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		utils.PostProcessResponse(w, cdc, res, cliCtx.Indent)
	}
}
