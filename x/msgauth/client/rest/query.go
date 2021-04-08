package rest

import (
	"fmt"
	"net/http"

	"github.com/terra-project/core/x/msgauth/types"

	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/types/rest"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func registerQueryRoutes(clientCtx client.Context, rtr *mux.Router) {
	rtr.HandleFunc(fmt.Sprintf("/msgauth/granters/{%s}/grants", RestGranter), queryAllGrantsHandlerFunction(clientCtx)).Methods("GET")
	rtr.HandleFunc(fmt.Sprintf("/msgauth/granters/{%s}/grantees/{%s}/grants", RestGranter, RestGrantee), queryGrantsHandlerFunction(clientCtx)).Methods("GET")
}

func queryGrantsHandlerFunction(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		clientCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, clientCtx, r)
		if !ok {
			return
		}

		vars := mux.Vars(r)
		granter := vars[RestGranter]
		grantee := vars[RestGrantee]

		granterAddr, err := sdk.AccAddressFromBech32(granter)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		granteeAddr, err := sdk.AccAddressFromBech32(grantee)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		params := types.NewQueryGrantsParams(granterAddr, granteeAddr)

		bz, err := clientCtx.LegacyAmino.MarshalJSON(params)
		if rest.CheckBadRequestError(w, err) {
			return
		}

		res, height, err := clientCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryGrants), bz)
		if rest.CheckInternalServerError(w, err) {
			return
		}

		clientCtx = clientCtx.WithHeight(height)
		rest.PostProcessResponse(w, clientCtx, res)
	}
}

func queryAllGrantsHandlerFunction(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		clientCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, clientCtx, r)
		if !ok {
			return
		}

		vars := mux.Vars(r)
		granter := vars[RestGranter]

		granterAddr, err := sdk.AccAddressFromBech32(granter)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		params := types.NewQueryAllGrantsParams(granterAddr)

		bz, err := clientCtx.LegacyAmino.MarshalJSON(params)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		res, height, err := clientCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryGrants), bz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		clientCtx = clientCtx.WithHeight(height)
		rest.PostProcessResponse(w, clientCtx, res)
	}
}
