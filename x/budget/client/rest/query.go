package rest

import (
	"fmt"
	"net/http"
	"terra/x/budget"
	"terra/x/oracle"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

func registerQueryRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *codec.Codec) {
	r.HandleFunc("/budget/programs/actives", queryActivesHandlerFn(cdc, cliCtx)).Methods("GET")
	r.HandleFunc("/budget/programs/candidates", queryCandidatesHandlerFn(cdc, cliCtx)).Methods("GET")

	r.HandleFunc(fmt.Sprintf("/budget/programs/{%s}", RestProgramID), queryProgramHandlerFn(cdc, cliCtx)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/budget/programs/{%s}/votes", RestProgramID), queryVotesHandlerFn(cdc, cliCtx)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/budget/programs/{%s}/votes/{%s}", RestProgramID, RestVoter), queryVotesHandlerFn(cdc, cliCtx)).Methods("GET")

	r.HandleFunc("/budget/params", queryParamsHandlerFn(cdc, cliCtx)).Methods("GET")
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

		res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s", budget.QuerierRoute, budget.QueryProgram, strProgramID), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cdc, res, cliCtx.Indent)
	}
}

func queryActivesHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", budget.QuerierRoute, budget.QueryActiveList), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cdc, res, cliCtx.Indent)
	}
}

func queryCandidatesHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", budget.QuerierRoute, budget.QueryCandidateList), nil)
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

		if len(strProgramID) == 0 {
			err := errors.New("programID should be specified")
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		programID, ok := rest.ParseUint64OrReturnBadRequest(w, strProgramID)
		if !ok {
			return
		}

		params.ProgramID = programID

		if len(strVoterAddr) != 0 {
			voterAcc, err := sdk.AccAddressFromBech32(strVoterAddr)
			if err != nil {
				err := errors.New("voter address malformed")
				rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}

			params.Voter = voterAcc
		}

		bz, err := cdc.MarshalJSON(params)
		if err != nil {
			return
		}

		res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", budget.QuerierRoute, budget.QueryVotes), bz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		rest.PostProcessResponse(w, cdc, res, cliCtx.Indent)
	}
}

func queryParamsHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", budget.QuerierRoute, oracle.QueryParams), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cdc, res, cliCtx.Indent)
	}
}
