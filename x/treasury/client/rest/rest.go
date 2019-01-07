package rest

import (
	"fmt"
	"net/http"
	"terra/x/treasury"

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
	RestShareID = "share-id"
	storeName   = "treasury"
)

// RegisterRoutes - Central function to define routes that get registered by the main application
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *codec.Codec) {
	r.HandleFunc("/tresury/assets", queryAssetHandlerFn(cdc, cliCtx)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/treasury/share/{%s}", RestShareID), queryShareHandlerFn(cdc, cliCtx)).Methods("GET")
}

func queryShareHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		strShareID := vars[RestShareID]

		if len(strShareID) == 0 {
			err := errors.New("shareID required but not specified")
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		switch strShareID {
		case treasury.OracleShareID:
			break
		case treasury.DebtShareID:
			break
		case treasury.SubsidyShareID:
			break
		default:
			err := errors.New("shareID not one of 'oracle' 'debt' 'subsidy'")
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		res, err := cliCtx.QueryStore(treasury.GetShareKey(strShareID), storeName)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		utils.PostProcessResponse(w, cdc, res, cliCtx.Indent)
	}
}

func queryAssetHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		res, err := cliCtx.QueryStore(treasury.GetIncomePoolKey(), storeName)
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
