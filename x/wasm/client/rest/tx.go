package rest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"

	feeutils "github.com/terra-project/core/custom/auth/client/utils"
	wasmUtils "github.com/terra-project/core/x/wasm/client/utils"
	"github.com/terra-project/core/x/wasm/types"
)

func registerTxRoutes(clientCtx client.Context, r *mux.Router) {
	r.HandleFunc("/wasm/codes", storeCodeHandlerFn(clientCtx)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/wasm/codes/{%s}", RestCodeID), instantiateContractHandlerFn(clientCtx)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/wasm/contracts/{%s}", RestContractAddress), executeContractHandlerFn(clientCtx)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/wasm/contract/{%s}/migrate", RestContractAddress), migrateContractHandlerFn(clientCtx)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/wasm/contract/{%s}/owner", RestContractAddress), updateOwnerContractHandlerFn(clientCtx)).Methods("POST")
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

func storeCodeHandlerFn(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req storeCodeReq
		if !rest.ReadRESTReq(w, r, clientCtx.LegacyAmino, &req) {
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
			if rest.CheckBadRequestError(w, err) {
				return
			}
		} else if !wasmUtils.IsGzip(wasmBytes) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "Invalid input file, use wasm binary or zip")
			return
		}

		fromAddr, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if rest.CheckBadRequestError(w, err) {
			return
		}

		// build and sign the transaction, then broadcast to Tendermint
		msg := types.NewMsgStoreCode(fromAddr, wasmBytes)
		if err = msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		tx.WriteGeneratedTxResponse(clientCtx, w, req.BaseReq, msg)
	}
}

func instantiateContractHandlerFn(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req instantiateContractReq
		if !rest.ReadRESTReq(w, r, clientCtx.LegacyAmino, &req) {
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
		if rest.CheckBadRequestError(w, err) {
			return
		}

		fromAddr, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if rest.CheckBadRequestError(w, err) {
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

func executeContractHandlerFn(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req executeContractReq
		if !rest.ReadRESTReq(w, r, clientCtx.LegacyAmino, &req) {
			return
		}
		vars := mux.Vars(r)
		contractAddr := vars[RestContractAddress]

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		contractAddress, err := sdk.AccAddressFromBech32(contractAddr)
		if rest.CheckBadRequestError(w, err) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		fromAddr, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if rest.CheckBadRequestError(w, err) {
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

func migrateContractHandlerFn(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req migrateContractReq
		if !rest.ReadRESTReq(w, r, clientCtx.LegacyAmino, &req) {
			return
		}
		vars := mux.Vars(r)
		contractAddr := vars[RestContractAddress]

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		contractAddress, err := sdk.AccAddressFromBech32(contractAddr)
		if rest.CheckBadRequestError(w, err) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		fromAddr, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if rest.CheckBadRequestError(w, err) {
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
		if rest.CheckBadRequestError(w, msg.ValidateBasic()) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		tx.WriteGeneratedTxResponse(clientCtx, w, req.BaseReq, msg)
	}
}

func updateOwnerContractHandlerFn(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req updateContractOwnerReq
		if !rest.ReadRESTReq(w, r, clientCtx.LegacyAmino, &req) {
			return
		}
		vars := mux.Vars(r)
		contractAddr := vars[RestContractAddress]

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		contractAddress, err := sdk.AccAddressFromBech32(contractAddr)
		if rest.CheckBadRequestError(w, err) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		fromAddr, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if rest.CheckBadRequestError(w, err) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		msg := types.NewMsgUpdateContractOwner(fromAddr, req.NewOwner, contractAddress)
		if rest.CheckBadRequestError(w, msg.ValidateBasic()) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		tx.WriteGeneratedTxResponse(clientCtx, w, req.BaseReq, msg)
	}
}
