package rest

import (
	"bytes"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	authclient "github.com/cosmos/cosmos-sdk/x/auth/client/utils"

	feeutils "github.com/terra-money/core/x/auth/client/utils"
	"github.com/terra-money/core/x/msgauth/internal/types"
)

func registerTxRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc(fmt.Sprintf("/msgauth/granters/{%s}/grantees/{%s}/grants/{%s}", RestGranter, RestGrantee, RestMsgType), grantHandler(cliCtx)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/msgauth/granters/{%s}/grantees/{%s}/grants/{%s}/revoke", RestGranter, RestGrantee, RestMsgType), revokeHandler(cliCtx)).Methods("POST")
	r.HandleFunc("/msgauth/execute", executeHandler(cliCtx)).Methods("POST")
}

func grantHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req GrantRequest

		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
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

		var authorization types.Authorization
		if msgType == (types.SendAuthorization{}.MsgType()) {
			authorization = types.NewSendAuthorization(req.Limit)
		} else {
			authorization = types.NewGenericAuthorization(msgType)
		}

		msg := types.NewMsgGrantAuthorization(granterAddr, granteeAddr, authorization, req.Period)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
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

		authclient.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

func revokeHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req RevokeRequest

		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
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
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
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

		authclient.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

func executeHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req ExecuteRequest

		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
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

		msg := types.NewMsgExecAuthorized(fromAddr, req.Msgs)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		if req.BaseReq.Fees.IsZero() {
			fees, gas, err := feeutils.ComputeFees(cliCtx, feeutils.ComputeReqParams{
				Memo:          req.BaseReq.Memo,
				ChainID:       req.BaseReq.ChainID,
				AccountNumber: req.BaseReq.AccountNumber,
				Sequence:      req.BaseReq.Sequence,
				GasPrices:     req.BaseReq.GasPrices,
				Gas:           req.BaseReq.Gas,
				GasAdjustment: req.BaseReq.GasAdjustment,
				Msgs:          []sdk.Msg{msg},
			})

			if err != nil {
				rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}

			// override gas and fees
			req.BaseReq.Gas = strconv.FormatUint(gas, 10)
			req.BaseReq.Fees = fees
			req.BaseReq.GasPrices = sdk.DecCoins{}
		}

		authclient.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}
