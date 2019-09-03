package rest

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/terra-project/core/x/treasury/internal/types"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"
)

func registerQueryRoute(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc("/treasury/tax_rate", queryTaxRateHandlerFunction(cliCtx)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/treasury/tax_rate/{%s}", RestEpoch), queryTaxRateHandlerFunction(cliCtx)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/treasury/tax_cap/{%s}", RestDenom), queryTaxCapHandlerFunction(cliCtx)).Methods("GET")
	r.HandleFunc("/treasury/reward_weight", queryRewardWeightHandlerFunction(cliCtx)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/treasury/reward_weight/{%s}", RestEpoch), queryRewardWeightHandlerFunction(cliCtx)).Methods("GET")
	r.HandleFunc("/treasury/historical_issuance", queryHistoricalIssuanceHandlerFunction(cliCtx)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/treasury/historical_issuance/{%s}", RestEpoch), queryHistoricalIssuanceHandlerFunction(cliCtx)).Methods("GET")
	r.HandleFunc("/treasury/tax_proceeds", queryTaxProceedsHandlerFunction(cliCtx)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/treasury/tax_proceeds/{%s}", RestEpoch), queryTaxProceedsHandlerFunction(cliCtx)).Methods("GET")
	r.HandleFunc("/treasury/seigniorage_proceeds", querySeigniorageProceedsHandlerFunction(cliCtx)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/treasury/seigniorage_proceeds/{%s}", RestEpoch), querySeigniorageProceedsHandlerFunction(cliCtx)).Methods("GET")
	r.HandleFunc("/treasury/current_epoch", queryCurrentEpochHandlerFunction(cliCtx)).Methods("GET")
	r.HandleFunc("/treasury/parameters", queryParametersHandlerFn(cliCtx)).Methods("GET")
}

func queryTaxRateHandlerFunction(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		vars := mux.Vars(r)
		epochStr := vars[RestEpoch]

		var epoch int64
		if len(epochStr) == 0 {
			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryCurrentEpoch), nil)
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
				return
			}

			cliCtx.Codec.MustUnmarshalJSON(res, &epoch)
		} else {
			var err error
			epoch, err = strconv.ParseInt(epochStr, 10, 64)
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusBadRequest, sdk.AppendMsgToErr("Falied to parse epoch", err.Error()))
				return
			}
		}

		params := types.NewQueryTaxRateParams(epoch)
		bz := cliCtx.Codec.MustMarshalJSON(params)

		res, height, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryTaxRate), bz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func queryTaxCapHandlerFunction(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		vars := mux.Vars(r)
		denom := vars[RestDenom]

		params := types.NewQueryTaxCapParams(denom)
		bz := cliCtx.Codec.MustMarshalJSON(params)

		res, height, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryTaxCap), bz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func queryRewardWeightHandlerFunction(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		vars := mux.Vars(r)
		epochStr := vars[RestEpoch]

		var epoch int64
		if len(epochStr) == 0 {
			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryCurrentEpoch), nil)
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
				return
			}

			cliCtx.Codec.MustUnmarshalJSON(res, &epoch)
		} else {
			var err error
			epoch, err = strconv.ParseInt(epochStr, 10, 64)
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusBadRequest, sdk.AppendMsgToErr("Falied to parse epoch", err.Error()))
				return
			}
		}

		params := types.NewQueryRewardWeightParams(epoch)
		bz := cliCtx.Codec.MustMarshalJSON(params)

		res, height, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryRewardWeight), bz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func queryHistoricalIssuanceHandlerFunction(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		vars := mux.Vars(r)
		epochStr := vars[RestEpoch]

		var epoch int64
		if len(epochStr) == 0 {
			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryCurrentEpoch), nil)
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
				return
			}

			cliCtx.Codec.MustUnmarshalJSON(res, &epoch)
		} else {
			var err error
			epoch, err = strconv.ParseInt(epochStr, 10, 64)
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusBadRequest, sdk.AppendMsgToErr("Falied to parse epoch", err.Error()))
				return
			}
		}

		params := types.NewQueryHistoricalIssuanceParams(epoch)
		bz := cliCtx.Codec.MustMarshalJSON(params)

		res, height, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryHistoricalIssuance), bz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func queryTaxProceedsHandlerFunction(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		vars := mux.Vars(r)
		epochStr := vars[RestEpoch]

		var epoch int64
		if len(epochStr) == 0 {
			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryCurrentEpoch), nil)
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
				return
			}

			cliCtx.Codec.MustUnmarshalJSON(res, &epoch)

		} else {
			var err error
			epoch, err = strconv.ParseInt(epochStr, 10, 64)
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusBadRequest, sdk.AppendMsgToErr("Falied to parse epoch", err.Error()))
				return
			}
		}

		params := types.NewQueryTaxProceedsParams(epoch)
		bz := cliCtx.Codec.MustMarshalJSON(params)

		res, height, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryTaxProceeds), bz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func querySeigniorageProceedsHandlerFunction(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		epochStr := vars[RestEpoch]

		var epoch int64
		if len(epochStr) == 0 {
			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryCurrentEpoch), nil)
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
				return
			}

			cliCtx.Codec.MustUnmarshalJSON(res, &epoch)

		} else {
			var err error
			epoch, err = strconv.ParseInt(epochStr, 10, 64)
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusBadRequest, sdk.AppendMsgToErr("Falied to parse epoch", err.Error()))
				return
			}
		}

		params := types.NewQuerySeigniorageParams(epoch)
		bz := cliCtx.Codec.MustMarshalJSON(params)

		res, height, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QuerySeigniorageProceeds), bz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func queryCurrentEpochHandlerFunction(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		res, height, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryCurrentEpoch), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func queryParametersHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		res, height, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryParameters), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}
