package rest

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"

	"github.com/cosmos/cosmos-sdk/x/bank"
	bankrest "github.com/cosmos/cosmos-sdk/x/bank/client/rest"

	feeutils "github.com/terra-money/core/x/auth/client/utils"
)

// RegisterRoutes - Central function to define routes that get registered by the main application
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc("/bank/accounts/{address}/transfers", SendRequestHandlerFn(cliCtx)).Methods("POST")
	r.HandleFunc("/bank/balances/{address}", bankrest.QueryBalancesRequestHandlerFn(cliCtx)).Methods("GET")
}

// SendReq defines the properties of a send request's body.
type SendReq struct {
	BaseReq rest.BaseReq `json:"base_req" yaml:"base_req"`
	Coins   sdk.Coins    `json:"coins" yaml:"coins"`
}

// SendRequestHandlerFn - http request handler to send coins to a address.
func SendRequestHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		bech32Addr := vars["address"]

		toAddr, err := sdk.AccAddressFromBech32(bech32Addr)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		var req SendReq
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

		msg := bank.MsgSend{FromAddress: fromAddr, ToAddress: toAddr, Amount: req.Coins}

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

		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}
