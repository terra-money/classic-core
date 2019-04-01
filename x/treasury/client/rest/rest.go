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
	RestDenom = "denom"
	RestEpoch = "epoch"
)

// RegisterRoutes - Central function to define routes that get registered by the main application
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *codec.Codec) {
	r.HandleFunc(fmt.Sprintf("/treasury/%s/{%s}", treasury.QueryTaxRate, RestEpoch), queryTaxRateHandlerFunction(cdc, cliCtx)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/treasury/%s/{%s}", treasury.QueryTaxCap, RestDenom), queryTaxCapHandlerFunction(cdc, cliCtx)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/treasury/%s/{%s}", treasury.QueryMiningRewardWeight, RestEpoch), queryMiningWeightHandlerFunction(cdc, cliCtx)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/treasury/%s/{%s}", treasury.QueryIssuance, RestDenom), queryIssuanceHandlerFunction(cdc, cliCtx)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/treasury/%s/{%s}", treasury.QueryTaxProceeds, RestEpoch), queryTaxProceedsHandlerFunction(cdc, cliCtx)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/treasury/%s/{%s}", treasury.QuerySeigniorageProceeds, RestEpoch), querySgProceedsHandlerFunction(cdc, cliCtx)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/treasury/%s", treasury.QueryActiveClaims), queryActiveClaimsHandlerFunction(cdc, cliCtx)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/treasury/%s", treasury.QueryCurrentEpoch), queryCurrentEpochHandlerFunction(cdc, cliCtx)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/treasury/%s", treasury.QueryParams), queryParamsHandlerFn(cdc, cliCtx)).Methods("GET")
}

func queryTaxRateHandlerFunction(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		epoch := vars[RestEpoch]

		res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s", treasury.QuerierRoute, treasury.QueryTaxRate, epoch), nil)
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

		res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s", treasury.QuerierRoute, treasury.QueryTaxCap, denom), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cdc, res, cliCtx.Indent)
	}
}

func queryMiningWeightHandlerFunction(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		epoch := vars[RestEpoch]

		res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s", treasury.QuerierRoute, treasury.QueryMiningRewardWeight, epoch), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		rest.PostProcessResponse(w, cdc, res, cliCtx.Indent)
	}
}

func queryIssuanceHandlerFunction(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		denom := vars[RestDenom]

		res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s", treasury.QuerierRoute, treasury.QueryIssuance, denom), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		rest.PostProcessResponse(w, cdc, res, cliCtx.Indent)
	}
}

func queryTaxProceedsHandlerFunction(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		epoch := vars[RestEpoch]

		res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s", treasury.QuerierRoute, treasury.QueryTaxProceeds, epoch), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		rest.PostProcessResponse(w, cdc, res, cliCtx.Indent)
	}
}

func querySgProceedsHandlerFunction(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		epoch := vars[RestEpoch]

		res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s", treasury.QuerierRoute, treasury.QuerySeigniorageProceeds, epoch), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		rest.PostProcessResponse(w, cdc, res, cliCtx.Indent)
	}
}

func queryActiveClaimsHandlerFunction(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", treasury.QuerierRoute, treasury.QueryActiveClaims), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		rest.PostProcessResponse(w, cdc, res, cliCtx.Indent)
	}
}

func queryCurrentEpochHandlerFunction(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", treasury.QuerierRoute, treasury.QueryCurrentEpoch), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		rest.PostProcessResponse(w, cdc, res, cliCtx.Indent)
	}
}

func queryParamsHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", treasury.QuerierRoute, treasury.QueryParams), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		rest.PostProcessResponse(w, cdc, res, cliCtx.Indent)
	}
}
