package rest

import (
	"fmt"
	"net/http"
	"terra/x/treasury"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/rest"

	"github.com/gorilla/mux"
)

// REST Variable names
// nolint
const (
	RestDenom  = "denom"
	queryRoute = "treasury"
)

// RegisterRoutes - Central function to define routes that get registered by the main application
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *codec.Codec) {
	r.HandleFunc("/treasury/taxrate", queryTaxRateHandlerFunction(cdc, cliCtx)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/treasury/taxcap/{%s}", RestDenom), queryTaxCapHandlerFunction(cdc, cliCtx)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/treasury/issuance/{%s}", RestDenom), queryIssuanceHandlerFunction(cdc, cliCtx)).Methods("GET")
	r.HandleFunc("/treasury/miningweight", queryMiningWeightHandlerFunction(cdc, cliCtx)).Methods("GET")
	r.HandleFunc("/treasury/balance", queryBalanceHandlerFunction(cdc, cliCtx)).Methods("GET")
	r.HandleFunc("/treasury/activeclaims", queryActiveClaimsHandlerFunction(cdc, cliCtx)).Methods("GET")
	r.HandleFunc("/treasury/params", queryParamsHandlerFn(cdc, cliCtx)).Methods("GET")
}

func queryTaxRateHandlerFunction(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", queryRoute, treasury.QueryTaxRate), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		rest.PostProcessResponse(w, cdc, res, cliCtx.Indent)
	}
}

func queryTaxCapHandlerFunction(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		denom := vars[RestDenom]

		res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s", queryRoute, treasury.QueryTaxCap, denom), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cdc, res, cliCtx.Indent)
	}
}

func queryIssuanceHandlerFunction(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		denom := vars[RestDenom]

		res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s", queryRoute, treasury.QueryIssuance, denom), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cdc, res, cliCtx.Indent)
	}
}

func queryMiningWeightHandlerFunction(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", queryRoute, treasury.QueryMiningRewardWeight), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		rest.PostProcessResponse(w, cdc, res, cliCtx.Indent)
	}
}

func queryBalanceHandlerFunction(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", queryRoute, treasury.QueryBalance), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		rest.PostProcessResponse(w, cdc, res, cliCtx.Indent)
	}
}

func queryActiveClaimsHandlerFunction(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", queryRoute, treasury.QueryActiveClaims), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cdc, res, cliCtx.Indent)
	}
}

func queryParamsHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", queryRoute, treasury.QueryParams), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cdc, res, cliCtx.Indent)
	}
}
