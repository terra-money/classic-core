package rest

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/gorilla/mux"

	feeutils "github.com/terra-project/core/x/auth/client/utils"
	wasmUtils "github.com/terra-project/core/x/wasm/client/utils"
	"github.com/terra-project/core/x/wasm/internal/types"
)

func registerTxRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc("/wasm/code/", storeCodeHandlerFn(cliCtx)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/wasm/code/{%s}", RestCodeID), instantiateContractHandlerFn(cliCtx)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/wasm/contract/{%s}", RestContractAddress), executeContractHandlerFn(cliCtx)).Methods("POST")
}

// limit max bytes read to prevent gzip bombs
const maxSize = 400 * 1024

type storeCodeReq struct {
	BaseReq   rest.BaseReq `json:"base_req" yaml:"base_req"`
	WasmBytes []byte       `json:"wasm_bytes"`
}

type instantiateContractReq struct {
	BaseReq   rest.BaseReq `json:"base_req" yaml:"base_req"`
	InitCoins sdk.Coins    `json:"init_coins" yaml:"init_coins"`
	InitMsg   []byte       `json:"init_msg" yaml:"init_msg"`
}

type executeContractReq struct {
	BaseReq rest.BaseReq `json:"base_req" yaml:"base_req"`
	ExecMsg []byte       `json:"exec_msg" yaml:"exec_msg"`
	Amount  sdk.Coins    `json:"coins" yaml:"coins"`
}

func storeCodeHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req storeCodeReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		var err error
		wasm := req.WasmBytes
		if len(wasm) > maxSize {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "Binary size exceeds maximum limit")
			return
		}

		// gzip the wasm file
		if wasmUtils.IsWasm(wasm) {
			wasm, err = wasmUtils.GzipIt(wasm)
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}
		} else if !wasmUtils.IsGzip(wasm) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "Invalid input file, use wasm binary or zip")
			return
		}

		fromAddr, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// build and sign the transaction, then broadcast to Tendermint
		msg := types.MsgStoreCode{
			Sender:       fromAddr,
			WASMByteCode: wasm,
		}

		err = msg.ValidateBasic()
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

func instantiateContractHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req instantiateContractReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}
		vars := mux.Vars(r)
		strCodeID := vars[RestCodeID]

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		// get the id of the code to instantiate
		codeID, err := strconv.ParseUint(strCodeID, 10, 64)
		if err != nil {
			return
		}

		fromAddr, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		msg := types.MsgInstantiateContract{
			Sender:    fromAddr,
			CodeID:    codeID,
			InitCoins: req.InitCoins,
			InitMsg:   req.InitMsg,
		}

		err = msg.ValidateBasic()
		if err != nil {
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

		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

func executeContractHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req executeContractReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}
		vars := mux.Vars(r)
		contractAddr := vars[RestContractAddress]

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		contractAddress, err := sdk.AccAddressFromBech32(contractAddr)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		fromAddr, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		msg := types.MsgExecuteContract{
			Sender:   fromAddr,
			Contract: contractAddress,
			Msg:      req.ExecMsg,
			Coins:    req.Amount,
		}

		err = msg.ValidateBasic()
		if err != nil {
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

		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}
