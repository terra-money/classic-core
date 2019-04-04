package rest

import (
	"fmt"
	"net/http"
	"github.com/terra-project/core/types/assets"
	"github.com/terra-project/core/x/treasury"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"
)

func registerQueryRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *codec.Codec) {

	r.HandleFunc(fmt.Sprintf("/treasury/%s", treasury.QueryTaxRate), queryTaxRateHandlerFunction(cdc, cliCtx)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/treasury/%s/{%s}", treasury.QueryTaxRate, RestEpoch), queryTaxRateHandlerFunction(cdc, cliCtx)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/treasury/%s/{%s}", treasury.QueryTaxCap, RestDenom), queryTaxCapHandlerFunction(cdc, cliCtx)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/treasury/%s", treasury.QueryMiningRewardWeight), queryMiningWeightHandlerFunction(cdc, cliCtx)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/treasury/%s/{%s}", treasury.QueryMiningRewardWeight, RestDenom), queryMiningWeightHandlerFunction(cdc, cliCtx)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/treasury/%s/{%s}", treasury.QueryIssuance, RestDenom), queryIssuanceHandlerFunction(cdc, cliCtx)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/treasury/%s", treasury.QueryTaxProceeds), queryTaxProceedsHandlerFunction(cdc, cliCtx)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/treasury/%s/{%s}", treasury.QueryTaxProceeds, RestEpoch), queryTaxProceedsHandlerFunction(cdc, cliCtx)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/treasury/%s", treasury.QuerySeigniorageProceeds), querySgProceedsHandlerFunction(cdc, cliCtx)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/treasury/%s/{%s}", treasury.QuerySeigniorageProceeds, RestEpoch), querySgProceedsHandlerFunction(cdc, cliCtx)).Methods("GET")

	r.HandleFunc(fmt.Sprintf("/treasury/%s", treasury.QueryActiveClaims), queryActiveClaimsHandlerFunction(cdc, cliCtx)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/treasury/%s", treasury.QueryCurrentEpoch), queryCurrentEpochHandlerFunction(cdc, cliCtx)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/treasury/%s", treasury.QueryParams), queryParamsHandlerFn(cdc, cliCtx)).Methods("GET")
}

func queryTaxRateHandlerFunction(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		epochStr := vars[RestEpoch]

		var epoch sdk.Int
		if len(epochStr) == 0 {
			res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", treasury.QuerierRoute, treasury.QueryCurrentEpoch), nil)
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
				return
			}

			cdc.MustUnmarshalJSON(res, &epoch)
		} else {
			var ok bool
			epoch, ok = sdk.NewIntFromString(epochStr)
			if !ok {
				err := fmt.Errorf("the given epoch {%s} is not a valid format; epoch should be formatted as an integer", epochStr)
				rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}
		}

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

		if !assets.IsValidDenom(denom) {
			err := fmt.Errorf("given denom {%s} is not a valid one", denom)
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

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
		epochStr := vars[RestEpoch]

		var epoch sdk.Int
		if len(epochStr) == 0 {
			res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", treasury.QuerierRoute, treasury.QueryCurrentEpoch), nil)
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
				return
			}

			cdc.MustUnmarshalJSON(res, &epoch)
		} else {
			var ok bool
			epoch, ok = sdk.NewIntFromString(epochStr)
			if !ok {
				err := fmt.Errorf("the given epoch {%s} is not a valid format; epoch should be formatted as an integer", epochStr)
				rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}
		}

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

		if !assets.IsValidDenom(denom) {
			err := fmt.Errorf("given denom {%s} is not a valid one", denom)
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

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
		epochStr := vars[RestEpoch]

		var epoch sdk.Int
		if len(epochStr) == 0 {
			res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", treasury.QuerierRoute, treasury.QueryCurrentEpoch), nil)
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
				return
			}

			cdc.MustUnmarshalJSON(res, &epoch)
		} else {
			var ok bool
			epoch, ok = sdk.NewIntFromString(epochStr)
			if !ok {
				err := fmt.Errorf("the given epoch {%s} is not a valid format; epoch should be formatted as an integer", epochStr)
				rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}
		}

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
		epochStr := vars[RestEpoch]

		var epoch sdk.Int
		if len(epochStr) == 0 {
			res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", treasury.QuerierRoute, treasury.QueryCurrentEpoch), nil)
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
				return
			}

			cdc.MustUnmarshalJSON(res, &epoch)
		} else {
			var ok bool
			epoch, ok = sdk.NewIntFromString(epochStr)
			if !ok {
				err := fmt.Errorf("the given epoch {%s} is not a valid format; epoch should be formatted as an integer", epochStr)
				rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}
		}

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
