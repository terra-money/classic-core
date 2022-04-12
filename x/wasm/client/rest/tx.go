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

	wasmUtils "github.com/terra-money/core/x/wasm/client/utils"
	"github.com/terra-money/core/x/wasm/types"
)

func registerTxRoutes(clientCtx client.Context, r *mux.Router) {
	r.HandleFunc("/wasm/codes", storeCodeHandlerFn(clientCtx)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/wasm/codes/{%s}", RestCodeID), instantiateContractHandlerFn(clientCtx)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/wasm/codes/{%s}/migrate", RestCodeID), migrateCodeHandlerFn(clientCtx)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/wasm/contracts/{%s}", RestContractAddress), executeContractHandlerFn(clientCtx)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/wasm/contract/{%s}/migrate", RestContractAddress), migrateContractHandlerFn(clientCtx)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/wasm/contract/{%s}/admin/update", RestContractAddress), updateContractAdminHandlerFn(clientCtx)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/wasm/contract/{%s}/admin/clear", RestContractAddress), clearContractAdminHandlerFn(clientCtx)).Methods("POST")
}

type storeCodeReq struct {
	BaseReq   rest.BaseReq `json:"base_req" yaml:"base_req"`
	WasmBytes []byte       `json:"wasm_bytes"`
}

type migrateCodeReq struct {
	BaseReq   rest.BaseReq `json:"base_req" yaml:"base_req"`
	WasmBytes []byte       `json:"wasm_bytes"`
}

type instantiateContractReq struct {
	BaseReq   rest.BaseReq `json:"base_req" yaml:"base_req"`
	InitCoins sdk.Coins    `json:"init_coins" yaml:"init_coins"`
	InitMsg   string       `json:"init_msg" yaml:"init_msg"`
	Admin     string       `json:"admin" yaml:"admin"`
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

type updateContractAdminReq struct {
	BaseReq  rest.BaseReq `json:"base_req" yaml:"base_req"`
	NewAdmin string       `json:"new_admin" yaml:"new_admin"`
}

type clearContractAdminReq struct {
	BaseReq rest.BaseReq `json:"base_req" yaml:"base_req"`
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

func migrateCodeHandlerFn(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req migrateCodeReq
		if !rest.ReadRESTReq(w, r, clientCtx.LegacyAmino, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		vars := mux.Vars(r)
		strCodeID := vars[RestCodeID]

		// get the id of the code to migrate
		codeID, err := strconv.ParseUint(strCodeID, 10, 64)
		if rest.CheckBadRequestError(w, err) {
			return
		}

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
		msg := types.NewMsgMigrateCode(codeID, fromAddr, wasmBytes)
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

		adminAddr := sdk.AccAddress{}
		if len(req.Admin) != 0 {
			adminAddr, err = sdk.AccAddressFromBech32(req.Admin)
			if rest.CheckBadRequestError(w, err) {
				rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}
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

		msg := types.NewMsgInstantiateContract(fromAddr, adminAddr, codeID, initMsgBz, req.InitCoins)
		if rest.CheckBadRequestError(w, msg.ValidateBasic()) {
			return
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

func updateContractAdminHandlerFn(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req updateContractAdminReq
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

		newAdminAddr, err := sdk.AccAddressFromBech32(req.NewAdmin)
		if rest.CheckBadRequestError(w, err) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		msg := types.NewMsgUpdateContractAdmin(fromAddr, newAdminAddr, contractAddress)
		if rest.CheckBadRequestError(w, msg.ValidateBasic()) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		tx.WriteGeneratedTxResponse(clientCtx, w, req.BaseReq, msg)
	}
}

func clearContractAdminHandlerFn(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req clearContractAdminReq
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

		msg := types.NewMsgClearContractAdmin(fromAddr, contractAddress)
		if rest.CheckBadRequestError(w, msg.ValidateBasic()) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		tx.WriteGeneratedTxResponse(clientCtx, w, req.BaseReq, msg)
	}
}
