package rest

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"

	"github.com/terra-project/core/x/nameservice/internal/types"
)

func registerTxRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc(fmt.Sprintf("/nameservice/names/{%s}/auction", RestName), submitOpenAuctionHandlerFn(cliCtx)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/nameservice/names/{%s}/auction/bids", RestName), submitBidAuctionHandlerFn(cliCtx)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/nameservice/names/{%s}/auction/reveals", RestName), submitRevealBidHandlerFn(cliCtx)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/nameservice/names/{%s}/registry/renew", RestName), submitRenewRegistryHandlerFn(cliCtx)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/nameservice/names/{%s}/registry/owner", RestName), submitUpdateOwnerHandlerFn(cliCtx)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/nameservice/names/{%s}/register", RestName), submitRegisterHandlerFn(cliCtx)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/nameservice/names/{%s}/unregister", RestName), submitUnregisterHandlerFn(cliCtx)).Methods("POST")
}

// OpenAuctionReq ...
type OpenAuctionReq struct {
	BaseReq rest.BaseReq `json:"base_req"`
}

func submitOpenAuctionHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		nameStr := vars[RestName]

		var req OpenAuctionReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()

		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		fromAddress, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		name := types.Name(nameStr)
		if err := name.Validate(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		if levels := name.Levels(); levels != 2 {
			err := fmt.Errorf("must submit by the second level name")
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// create the message
		msg := types.NewMsgOpenAuction(name, fromAddress)
		err = msg.ValidateBasic()
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

// BidAuctionReq ...
type BidAuctionReq struct {
	BaseReq rest.BaseReq `json:"base_req"`

	Salt    string   `json:"salt"`
	Amount  sdk.Coin `json:"amount"`
	Deposit sdk.Coin `json:"deposit"`
}

func submitBidAuctionHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		nameStr := vars[RestName]

		var req BidAuctionReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()

		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		fromAddress, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		name := types.Name(nameStr)
		if err := name.Validate(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		if levels := name.Levels(); levels != 2 {
			err := fmt.Errorf("must submit by the second level name")
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		bidHash := types.GetBidHash(req.Salt, name, req.Amount, fromAddress)

		// create the message
		msg := types.NewMsgBidAuction(name, bidHash, req.Deposit, fromAddress)
		err = msg.ValidateBasic()
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

// RevealBidReq ...
type RevealBidReq struct {
	BaseReq rest.BaseReq `json:"base_req"`

	Salt   string   `json:"salt"`
	Amount sdk.Coin `json:"amount"`
}

func submitRevealBidHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		nameStr := vars[RestName]

		var req RevealBidReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()

		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		fromAddress, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		name := types.Name(nameStr)
		if err := name.Validate(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		if levels := name.Levels(); levels != 2 {
			err := fmt.Errorf("must submit by the second level name")
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// create the message
		msg := types.NewMsgRevealBid(name, req.Salt, req.Amount, fromAddress)
		err = msg.ValidateBasic()
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

// RenewRegistryReq ...
type RenewRegistryReq struct {
	BaseReq rest.BaseReq `json:"base_req"`
	Amount  sdk.Coins    `json:"amount"`
}

func submitRenewRegistryHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		nameStr := vars[RestName]

		var req RenewRegistryReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()

		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		fromAddress, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		name := types.Name(nameStr)
		if err := name.Validate(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		if levels := name.Levels(); levels != 2 {
			err := fmt.Errorf("must submit by the second level name")
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// create the message
		msg := types.NewMsgRenewRegistry(name, req.Amount, fromAddress)
		err = msg.ValidateBasic()
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

// UpdateOwnerReq ...
type UpdateOwnerReq struct {
	BaseReq   rest.BaseReq   `json:"base_req"`
	OwnerAddr sdk.AccAddress `json:"owner_addr"`
}

func submitUpdateOwnerHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		nameStr := vars[RestName]

		var req UpdateOwnerReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()

		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		fromAddress, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		name := types.Name(nameStr)
		if err := name.Validate(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		if levels := name.Levels(); levels != 2 {
			err := fmt.Errorf("must submit by the second level name")
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// create the message
		msg := types.NewMsgUpdateOwner(name, req.OwnerAddr, fromAddress)
		err = msg.ValidateBasic()
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

// RegisterReq ...
type RegisterReq struct {
	BaseReq rest.BaseReq   `json:"base_req"`
	Address sdk.AccAddress `json:"address"`
}

func submitRegisterHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		nameStr := vars[RestName]

		var req RegisterReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()

		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		fromAddress, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		name := types.Name(nameStr)
		if err := name.Validate(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		if levels := name.Levels(); levels != 2 && levels != 3 {
			err := fmt.Errorf("must submit by the second or third level name")
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// create the message
		msg := types.NewMsgRegisterSubName(name, req.Address, fromAddress)
		err = msg.ValidateBasic()
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

// UnregisterReq ...
type UnregisterReq struct {
	BaseReq rest.BaseReq `json:"base_req"`
}

func submitUnregisterHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		nameStr := vars[RestName]

		var req UnregisterReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()

		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		fromAddress, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		name := types.Name(nameStr)
		if err := name.Validate(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		if levels := name.Levels(); levels != 2 && levels != 3 {
			err := fmt.Errorf("must submit by the second or third level name")
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// create the message
		msg := types.NewMsgUnregisterSubName(name, fromAddress)
		err = msg.ValidateBasic()
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}
