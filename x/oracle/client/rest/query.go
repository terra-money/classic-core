package rest

import (
	"fmt"
	"github.com/terra-project/core/types/assets"
	"github.com/terra-project/core/x/oracle"
	"net/http"

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
}

func queryVotesHandlerFunction(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		denom := vars[RestDenom]

		if !assets.IsValidDenom(denom) {
			err := fmt.Errorf("The denom is not known: %s", denom)
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		voter := vars[RestVoter]

		var voterAddress sdk.AccAddress
		params := oracle.NewQueryVoteParams(voterAddress, denom)

		if len(voter) != 0 {

			voterAddress, err := sdk.AccAddressFromBech32(voter)
			if err != nil {
				return
			}
			params.Voter = voterAddress
		}

		bz, err := cdc.MarshalJSON(params)
		if err != nil {
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

		if !assets.IsValidDenom(denom) {
			err := fmt.Errorf("The denom is not known: %s", denom)
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s", oracle.QuerierRoute, oracle.QueryPrice, denom), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cdc, res, cliCtx.Indent)
	}
}

func queryActivesHandlerFunction(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", oracle.QuerierRoute, oracle.QueryActive), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cdc, res, cliCtx.Indent)
	}
}

func queryParamsHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", oracle.QuerierRoute, oracle.QueryParams), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cdc, res, cliCtx.Indent)
	}
}
