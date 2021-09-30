package rest

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"

	"github.com/terra-money/core/x/wasm/types"

	"github.com/gorilla/mux"
)

func registerQueryRoutes(clientCtx client.Context, rtr *mux.Router) {
	rtr.HandleFunc(fmt.Sprintf("/wasm/codes/{%s}", RestCodeID), queryCodeInfoHandlerFn(clientCtx)).Methods("GET")
	rtr.HandleFunc(fmt.Sprintf("/wasm/contracts/{%s}", RestContractAddress), queryContractInfoHandlerFn(clientCtx)).Methods("GET")
	rtr.HandleFunc(fmt.Sprintf("/wasm/contracts/{%s}/store", RestContractAddress), queryContractStoreHandlerFn(clientCtx)).Methods("GET")
	rtr.HandleFunc(fmt.Sprintf("/wasm/contracts/{%s}/store/raw", RestContractAddress), queryRawStoreHandlerFn(clientCtx)).Methods("GET")
	rtr.HandleFunc("/wasm/parameters", queryParamsHandlerFn(clientCtx)).Methods("GET")
}

func queryCodeInfoHandlerFn(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		clientCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, clientCtx, r)
		if !ok {
			return
		}

		vars := mux.Vars(r)
		codeIDStr := vars[RestCodeID]

		codeID, err := strconv.ParseUint(codeIDStr, 10, 64)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		params := types.NewQueryCodeIDParams(codeID)
		bz, err := clientCtx.LegacyAmino.MarshalJSON(params)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		}

		route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryGetCodeInfo)
		res, height, err := clientCtx.QueryWithData(route, bz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		}

		clientCtx = clientCtx.WithHeight(height)
		rest.PostProcessResponse(w, clientCtx, res)
	}
}

func queryContractInfoHandlerFn(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		clientCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, clientCtx, r)
		if !ok {
			return
		}

		vars := mux.Vars(r)
		contractAddrStr := vars[RestContractAddress]

		addr, err := sdk.AccAddressFromBech32(contractAddrStr)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		params := types.NewQueryContractAddressParams(addr)
		bz, err := clientCtx.LegacyAmino.MarshalJSON(params)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryGetContractInfo)
		res, height, err := clientCtx.QueryWithData(route, bz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		clientCtx = clientCtx.WithHeight(height)
		rest.PostProcessResponse(w, clientCtx, res)
	}
}

func queryContractStoreHandlerFn(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		clientCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, clientCtx, r)
		if !ok {
			return
		}

		vars := mux.Vars(r)
		contractAddrStr := vars[RestContractAddress]
		queryMsg := r.URL.Query().Get("query_msg")
		queryMsgBz := []byte(queryMsg)
		if !json.Valid(queryMsgBz) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "msg must be a json string format")
			return
		}

		addr, err := sdk.AccAddressFromBech32(contractAddrStr)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		params := types.NewQueryContractParams(addr, queryMsgBz)
		bz, err := clientCtx.LegacyAmino.MarshalJSON(params)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryContractStore)
		res, height, err := clientCtx.QueryWithData(route, bz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		clientCtx = clientCtx.WithHeight(height)
		rest.PostProcessResponse(w, clientCtx, res)
	}
}

func queryRawStoreHandlerFn(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		clientCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, clientCtx, r)
		if !ok {
			return
		}

		vars := mux.Vars(r)
		contractAddrStr := vars[RestContractAddress]
		keyBz, err := base64.StdEncoding.DecodeString(r.URL.Query().Get("key"))
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		addr, err := sdk.AccAddressFromBech32(contractAddrStr)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		params := types.NewQueryRawStoreParams(addr, keyBz)
		bz, err := clientCtx.LegacyAmino.MarshalJSON(params)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryRawStore)
		res, height, err := clientCtx.QueryWithData(route, bz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		model := types.Model{
			Key:   keyBz,
			Value: res,
		}

		clientCtx = clientCtx.WithHeight(height)
		rest.PostProcessResponse(w, clientCtx, model)
	}
}

func queryParamsHandlerFn(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		clientCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, clientCtx, r)
		if !ok {
			return
		}

		res, height, err := clientCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryParameters), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		clientCtx = clientCtx.WithHeight(height)
		rest.PostProcessResponse(w, clientCtx, res)
	}
}
