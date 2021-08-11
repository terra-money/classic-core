package rest

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client"
	clientrest "github.com/cosmos/cosmos-sdk/client/rest"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	bankrest "github.com/cosmos/cosmos-sdk/x/bank/client/rest"
	"github.com/cosmos/cosmos-sdk/x/bank/types"

	feeutils "github.com/terra-money/core/custom/auth/client/utils"
)

// RegisterRoutes registers bank-related REST handlers to a router
func RegisterRoutes(clientCtx client.Context, rtr *mux.Router) {
	r := clientrest.WithHTTPDeprecationHeaders(rtr)

	r.HandleFunc("/bank/accounts/{address}/transfers", sendRequestHandlerFn(clientCtx)).Methods("POST")
	bankrest.RegisterHandlers(clientCtx, rtr)
}

// SendReq defines the properties of a send request's body.
type SendReq struct {
	BaseReq rest.BaseReq `json:"base_req" yaml:"base_req"`
	Coins   sdk.Coins    `json:"coins" yaml:"coins"`
}

// sendRequestHandlerFn - http request handler to send coins to a address.
func sendRequestHandlerFn(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		bech32Addr := vars["address"]

		toAddr, err := sdk.AccAddressFromBech32(bech32Addr)
		if rest.CheckBadRequestError(w, err) {
			return
		}

		var req SendReq
		if !rest.ReadRESTReq(w, r, clientCtx.LegacyAmino, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		fromAddr, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if rest.CheckBadRequestError(w, err) {
			return
		}

		msg := types.NewMsgSend(fromAddr, toAddr, req.Coins)
		if rest.CheckBadRequestError(w, msg.ValidateBasic()) {
			return
		}

		if req.BaseReq.Fees.IsZero() {
			stdFee, err := feeutils.ComputeFeesWithBaseReq(clientCtx, req.BaseReq, msg)

			if err != nil {
				rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
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
