package rest

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/terra-project/core/x/market"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func registerQueryRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *codec.Codec) {
	r.HandleFunc("/market/swap", querySwapHandlerFn(cdc, cliCtx)).Methods("GET")
	r.HandleFunc("/market/params", queryParamsHandlerFn(cdc, cliCtx)).Methods("GET")
}

func querySwapHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		err := r.ParseForm()
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest,
				sdk.AppendMsgToErr("could not parse query parameters", err.Error()))
			return
		}

		if len(r.Form) == 0 {
			err := errors.New("ask_denom & offer_coin should be specified")
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		askDenom := r.Form.Get("ask_denom")
		offerCoinStr := r.Form.Get("offer_coin")

		// parse offerCoin
		offerCoin, err := sdk.ParseCoin(offerCoinStr)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		params := market.NewQuerySwapParams(offerCoin)
		bz := cdc.MustMarshalJSON(params)
		res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s", market.QuerierRoute, market.QuerySwap, askDenom), bz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		rest.PostProcessResponse(w, cdc, res, cliCtx.Indent)
	}
}
func queryParamsHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", market.QuerierRoute, market.QueryParams), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cdc, res, cliCtx.Indent)
	}
}
