package rest

import (
	"bytes"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	feeutils "github.com/terra-project/core/custom/auth/client/utils"
	"github.com/terra-project/core/x/msgauth/types"
)

func registerTxRoutes(clientCtx client.Context, rtr *mux.Router) {
	rtr.HandleFunc(fmt.Sprintf("/msgauth/granters/{%s}/grantees/{%s}/grants/{%s}", RestGranter, RestGrantee, RestMsgType), grantHandler(clientCtx)).Methods("POST")
	rtr.HandleFunc(fmt.Sprintf("/msgauth/granters/{%s}/grantees/{%s}/grants/{%s}/revoke", RestGranter, RestGrantee, RestMsgType), revokeHandler(clientCtx)).Methods("POST")
	rtr.HandleFunc("/msgauth/execute", executeHandler(clientCtx)).Methods("POST")
}

func grantHandler(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req GrantRequest

		if !rest.ReadRESTReq(w, r, clientCtx.LegacyAmino, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		vars := mux.Vars(r)
		granter := vars[RestGranter]
		grantee := vars[RestGrantee]
		msgType := vars[RestMsgType]

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

		var authorization types.AuthorizationI
		if msgType == (banktypes.TypeMsgSend) {
			authorization = types.NewSendAuthorization(req.Limit)
		} else {
			authorization = types.NewGenericAuthorization(msgType)
		}

		msg, err := types.NewMsgGrantAuthorization(granterAddr, granteeAddr, authorization, req.Period)
		if rest.CheckBadRequestError(w, err) {
			return
		}

		if rest.CheckBadRequestError(w, msg.ValidateBasic()) {
			return
		}

		fromAddr, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		if !bytes.Equal(fromAddr, granterAddr) {
			rest.WriteErrorResponse(w, http.StatusUnauthorized, "must use own granter address")
			return
		}

		tx.WriteGeneratedTxResponse(clientCtx, w, req.BaseReq, msg)
	}
}

func revokeHandler(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req RevokeRequest

		if !rest.ReadRESTReq(w, r, clientCtx.LegacyAmino, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		vars := mux.Vars(r)
		granter := vars[RestGranter]
		grantee := vars[RestGrantee]
		msgType := vars[RestMsgType]

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

		msg := types.NewMsgRevokeAuthorization(granterAddr, granteeAddr, msgType)
		if rest.CheckBadRequestError(w, msg.ValidateBasic()) {
			return
		}

		fromAddr, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		if !bytes.Equal(fromAddr, granteeAddr) {
			rest.WriteErrorResponse(w, http.StatusUnauthorized, "must use own grantee address")
			return
		}

		tx.WriteGeneratedTxResponse(clientCtx, w, req.BaseReq, msg)
	}
}

func executeHandler(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req ExecuteRequest

		if !rest.ReadRESTReq(w, r, clientCtx.LegacyAmino, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		fromAddr, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		msg, err := types.NewMsgExecAuthorized(fromAddr, req.Msgs)
		if rest.CheckBadRequestError(w, err) {
			return
		}

		if rest.CheckBadRequestError(w, msg.ValidateBasic()) {
			return
		}

		if req.BaseReq.Fees.IsZero() {
			stdFee, err := feeutils.ComputeFeesWithBaseReq(clientCtx, req.BaseReq, msg)
			if rest.CheckBadRequestError(w, err) {
				return
			}

			// override gas and fees
			req.BaseReq.Gas = strconv.FormatUint(stdFee.Gas, 10)
			req.BaseReq.Fees = stdFee.Amount
			req.BaseReq.GasPrices = sdk.DecCoins{}
		}

		tx.WriteGeneratedTxResponse(clientCtx, w, req.BaseReq, msg)
	}
}
