package rest

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	authclient "github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/terra-project/core/x/bank"
	"github.com/terra-project/core/x/msgauth/internal/types"
)

func registerTxRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc(fmt.Sprintf("/msgauth/granters/{%s}/grantees/{%s}/grants/{%s}", RestGranter, RestGrantee, RestMsgType), grantHandler(cliCtx)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/msgauth/granters/{%s}/grantees/{%s}/grants/{%s}/revoke", RestGranter, RestGrantee, RestMsgType), revokeHandler(cliCtx)).Methods("POST")
}

type GrantRequest struct {
	BaseReq rest.BaseReq  `json:"base_req" yaml:"base_req"`
	Limit   string        `json:"limit,omitempty"`
	Period  time.Duration `json:"period"`
}

type RevokeRequest struct {
	BaseReq rest.BaseReq `json:"base_req" yaml:"base_req"`
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

		var authorization types.Authorization
		if msgType == (bank.MsgSend{}.Type()) {
			limit, err := sdk.ParseCoins(args[2])
			if err != nil {
				return err
			}

			authorization = types.NewSendAuthorization(limit)
		} else {
			authorization = types.NewGenericAuthorization(msgType)
		}

		msg := types.NewMsgGrantAuthorization(granter, grantee, authorization, req.Period)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
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

		msg := types.NewMsgRevokeAuthorization(granter, grantee, msgType)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		authclient.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}
