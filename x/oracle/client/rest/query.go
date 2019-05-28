package rest

import (
	"fmt"
	"net/http"

	"github.com/terra-project/core/x/oracle"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"
)

func registerQueryRoute(cliCtx context.CLIContext, r *mux.Router, cdc *codec.Codec) {
	r.HandleFunc(fmt.Sprintf("/oracle/denoms/{%s}/votes", RestDenom), queryVotesHandlerFunction(cdc, cliCtx)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/oracle/denoms/{%s}/votes/{%s}", RestDenom, RestVoter), queryVotesHandlerFunction(cdc, cliCtx)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/oracle/denoms/{%s}/price", RestDenom), queryPriceHandlerFunction(cdc, cliCtx)).Methods("GET")
	r.HandleFunc("/oracle/denoms/actives", queryActivesHandlerFunction(cdc, cliCtx)).Methods("GET")
	r.HandleFunc("/oracle/params", queryParamsHandlerFn(cdc, cliCtx)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/oracle/voters/{%s}/delegation", RestVoter), queryFeederDelegationHandlerFn(cdc, cliCtx)).Methods("GET")
}

func queryVotesHandlerFunction(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		denom := vars[RestDenom]

		voter := vars[RestVoter]

		var voterAddress sdk.ValAddress
		params := oracle.NewQueryVoteParams(voterAddress, denom)

		if len(voter) != 0 {

			voterAddress, err := sdk.ValAddressFromBech32(voter)
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}
			params.Voter = voterAddress
		}

		bz, err := cdc.MarshalJSON(params)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", oracle.QuerierRoute, oracle.QueryVotes), bz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		rest.PostProcessResponse(w, cdc, res, cliCtx.Indent)
	}
}

func queryPriceHandlerFunction(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		denom := vars[RestDenom]

		res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s", oracle.QuerierRoute, oracle.QueryPrice, denom), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		rest.PostProcessResponse(w, cdc, res, cliCtx.Indent)
	}
}

func queryActivesHandlerFunction(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", oracle.QuerierRoute, oracle.QueryActive), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		rest.PostProcessResponse(w, cdc, res, cliCtx.Indent)
	}
}

func queryParamsHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", oracle.QuerierRoute, oracle.QueryParams), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		rest.PostProcessResponse(w, cdc, res, cliCtx.Indent)
	}
}

func queryFeederDelegationHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		voter := vars[RestVoter]

		validator, err := sdk.ValAddressFromBech32(voter)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		params := oracle.NewQueryFeederDelegationParams(validator)
		bz, err := cdc.MarshalJSON(params)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", oracle.QuerierRoute, oracle.QueryFeederDelegation), bz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		var delegatee sdk.AccAddress
		cdc.MustUnmarshalJSON(res, &delegatee)

		rest.PostProcessResponse(w, cdc, res, cliCtx.Indent)
	}
}
