package rest

import (
	"net/http"
	"terra/x/treasury"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/utils"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/gorilla/mux"
)

// REST Variable names
// nolint
const (
	storeName = "treasury"
)

// RegisterRoutes - Central function to define routes that get registered by the main application
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *codec.Codec) {
	r.HandleFunc("/tresury/assets", queryAssetHandlerFn(cdc, cliCtx)).Methods("GET")
}

func queryAssetHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		res, err := cliCtx.QueryStore(treasury.KeyIncomePool, storeName)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		var coins sdk.Coins
		err = cdc.UnmarshalJSON(res, &coins)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		assetValue := coins[0].Amount

		utils.PostProcessResponse(w, cdc, assetValue, cliCtx.Indent)
	}
}
