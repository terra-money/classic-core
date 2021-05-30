package rest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/gorilla/mux"

	feeutils "github.com/terra-money/core/x/auth/client/utils"
	wasmUtils "github.com/terra-money/core/x/wasm/client/utils"
	"github.com/terra-money/core/x/wasm/internal/types"
)

func registerTxRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc("/wasm/codes", storeCodeHandlerFn(cliCtx)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/wasm/codes/{%s}", RestCodeID), instantiateContractHandlerFn(cliCtx)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/wasm/contracts/{%s}", RestContractAddress), executeContractHandlerFn(cliCtx)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/wasm/contract/{%s}/migrate", RestContractAddress), migrateContractHandlerFn(cliCtx)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/wasm/contract/{%s}/owner", RestContractAddress), updateOwnerContractHandlerFn(cliCtx)).Methods("POST")
}

type storeCodeReq struct {
	BaseReq   rest.BaseReq `json:"base_req" yaml:"base_req"`
	WasmBytes []byte       `json:"wasm_bytes"`
}

type instantiateContractReq struct {
	BaseReq    rest.BaseReq `json:"base_req" yaml:"base_req"`
	InitCoins  sdk.Coins    `json:"init_coins" yaml:"init_coins"`
	InitMsg    string       `json:"init_msg" yaml:"init_msg"`
	Migratable bool         `json:"migratable" yaml:"migratable"`
}

type executeContractReq struct {
	BaseReq rest.BaseReq `json:"base_req" yaml:"base_req"`
	Amount  sdk.Coins    `json:"coins" yaml:"coins"`
	ExecMsg string       `json:"exec_msg" yaml:"exec_msg"`
}

type migrateContractReq struct {
	BaseReq    rest.BaseReq `json:"base_req" yaml:"base_req"`
	MigrateMsg string       `json:"migrate_msg" yaml:"migrate_msg"`
	NewCodeID  uint64       `json:"new_code_id" yaml:"new_code_id"`
}

type updateContractOwnerReq struct {
	BaseReq  rest.BaseReq   `json:"base_req" yaml:"base_req"`
	NewOwner sdk.AccAddress `json:"new_owner" yaml:"new_owner"`
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
		wasmBytes := req.WasmBytes
		if wasmBytesLen := uint64(len(wasmBytes)); wasmBytesLen > types.EnforcedMaxContractSize {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "Binary size exceeds maximum limit")
			return
		}

		// gzip the wasm file
		if wasmUtils.IsWasm(wasmBytes) {
			wasmBytes, err = wasmUtils.GzipIt(wasmBytes)
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}
		} else if !wasmUtils.IsGzip(wasmBytes) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "Invalid input file, use wasm binary or zip")
			return
		}

		fromAddr, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// build and sign the transaction, then broadcast to Tendermint
		msg := types.NewMsgStoreCode(fromAddr, wasmBytes)
		if err = msg.ValidateBasic(); err != nil {
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

		initMsgBz := []byte(req.InitMsg)
		if !json.Valid(initMsgBz) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "msg must be a json string format")
			return
		}

		// limit the input size
		if initMsgLen := uint64(len(initMsgBz)); initMsgLen > types.EnforcedMaxContractMsgSize {
			rest.WriteErrorResponse(w, http.StatusBadRequest,
				fmt.Sprintf("init msg size exceeds the max size hard-cap (allowed:%d, actual: %d)",
					types.EnforcedMaxContractMsgSize, initMsgLen))
			return
		}

		msg := types.NewMsgInstantiateContract(fromAddr, codeID, initMsgBz, req.InitCoins, req.Migratable)
		if err = msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
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

		execMsgBz := []byte(req.ExecMsg)
		if !json.Valid(execMsgBz) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "msg must be a json string format")
			return
		}

		// limit the input size
		if execMsgLen := uint64(len(execMsgBz)); execMsgLen > types.EnforcedMaxContractMsgSize {
			rest.WriteErrorResponse(w, http.StatusBadRequest,
				fmt.Sprintf("exec msg size exceeds the max size hard-cap (allowed:%d, actual: %d)",
					types.EnforcedMaxContractMsgSize, execMsgLen))
			return
		}

		msg := types.NewMsgExecuteContract(fromAddr, contractAddress, execMsgBz, req.Amount)
		if err = msg.ValidateBasic(); err != nil {
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

func migrateContractHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req migrateContractReq
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

		migrateMsgBz := []byte(req.MigrateMsg)
		if !json.Valid(migrateMsgBz) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "msg must be a json string format")
			return
		}

		// limit the input size
		if migrateMsgLen := uint64(len(migrateMsgBz)); migrateMsgLen > types.EnforcedMaxContractMsgSize {
			rest.WriteErrorResponse(w, http.StatusBadRequest,
				fmt.Sprintf("migrate msg size exceeds the max size hard-cap (allowed:%d, actual: %d)",
					types.EnforcedMaxContractMsgSize, migrateMsgLen))
			return
		}

		msg := types.NewMsgMigrateContract(fromAddr, contractAddress, req.NewCodeID, migrateMsgBz)
		if err = msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

func updateOwnerContractHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req updateContractOwnerReq
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

		msg := types.NewMsgUpdateContractOwner(fromAddr, req.NewOwner, contractAddress)
		if err = msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}
