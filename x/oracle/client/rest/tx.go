package rest

import (
	"fmt"
	"net/http"

	"github.com/terra-project/core/x/oracle/internal/types"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"

	"github.com/gorilla/mux"
)

func resgisterTxRoute(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc(fmt.Sprintf("/oracle/denoms/{%s}/prevotes", RestDenom), submitPrevoteHandlerFunction(cliCtx)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/oracle/denoms/{%s}/votes", RestDenom), submitVoteHandlerFunction(cliCtx)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/oracle/voters/{%s}/feeder", RestVoter), submitDelegateHandlerFunction(cliCtx)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/oracle/voters/{%s}/aggregate_prevote", RestVoter), submitAggregatePrevoteHandlerFunction(cliCtx)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/oracle/voters/{%s}/aggregate_vote", RestVoter), submitAggregateVoteHandlerFunction(cliCtx)).Methods("POST")
}

// PrevoteReq ...
type PrevoteReq struct {
	BaseReq rest.BaseReq `json:"base_req"`

	Hash         string  `json:"hash"`
	ExchangeRate sdk.Dec `json:"exchange_rate"`
	Salt         string  `json:"salt"`

	Validator sdk.ValAddress `json:"validator"`
}

func submitPrevoteHandlerFunction(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		denom := vars[RestDenom]

		var req PrevoteReq
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

		// Default validator is self address
		valAddress := req.Validator
		if len(req.Validator) == 0 {
			valAddress = sdk.ValAddress(fromAddress)
		}

		var hash types.VoteHash

		// If hash is not given, then retrieve hash from exchange_rate and salt
		if len(req.Hash) == 0 && (!req.ExchangeRate.Equal(sdk.ZeroDec()) && len(req.Salt) > 0) {
			hash = types.GetVoteHash(req.Salt, req.ExchangeRate, denom, valAddress)
		} else {
			hash, err = types.VoteHashFromHexString(req.Hash)
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}
		}

		// create the message
		msg := types.NewMsgExchangeRatePrevote(hash, denom, fromAddress, valAddress)
		err = msg.ValidateBasic()
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

//VoteReq ...
type VoteReq struct {
	BaseReq rest.BaseReq `json:"base_req"`

	ExchangeRate sdk.Dec `json:"exchange_rate"`
	Salt         string  `json:"salt"`

	Validator sdk.ValAddress `json:"validator"`
}

func submitVoteHandlerFunction(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		denom := vars[RestDenom]

		var req VoteReq
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

		// Default validator is self address
		valAddress := req.Validator
		if len(req.Validator) == 0 {
			valAddress = sdk.ValAddress(fromAddress)
		}

		// create the message
		msg := types.NewMsgExchangeRateVote(req.ExchangeRate, req.Salt, denom, fromAddress, valAddress)
		err = msg.ValidateBasic()
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

// DelegateReq is request body to set feeder of validator
type DelegateReq struct {
	BaseReq rest.BaseReq `json:"base_req"`
	Feeder  string       `json:"feeder"`
}

func submitDelegateHandlerFunction(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		voter := vars[RestVoter]

		// Get voter validator address
		valAddress, err := sdk.ValAddressFromBech32(voter)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		var req DelegateReq
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

		// Bytes comparison, so do not require type conversion
		if !valAddress.Equals(fromAddress) {
			err := fmt.Errorf("[%v] can not change [%v] delegation", fromAddress, valAddress)
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		feeder, err := sdk.AccAddressFromBech32(req.Feeder)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// create the message
		msg := types.NewMsgDelegateFeedConsent(valAddress, feeder)
		err = msg.ValidateBasic()
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

// AggregatePrevoteReq ...
type AggregatePrevoteReq struct {
	BaseReq rest.BaseReq `json:"base_req"`

	Hash          string `json:"hash"`
	ExchangeRates string `json:"exchange_rates"`
	Salt          string `json:"salt"`
}

func submitAggregatePrevoteHandlerFunction(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		voter := vars[RestVoter]

		valAddress, err := sdk.ValAddressFromBech32(voter)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		var req AggregatePrevoteReq
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

		var hash types.AggregateVoteHash

		// If hash is not given, then retrieve hash from exchange_rate and salt
		if len(req.Hash) == 0 && (len(req.ExchangeRates) > 0 && len(req.Salt) > 0) {
			_, err := types.ParseExchangeRateTuples(req.ExchangeRates)
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}

			hash = types.GetAggregateVoteHash(req.Salt, req.ExchangeRates, valAddress)
		} else {
			hash, err = types.AggregateVoteHashFromHexString(req.Hash)
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}
		}

		// create the message
		msg := types.NewMsgAggregateExchangeRatePrevote(hash, fromAddress, valAddress)
		err = msg.ValidateBasic()
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

// AggregateVoteReq ...
type AggregateVoteReq struct {
	BaseReq rest.BaseReq `json:"base_req"`

	ExchangeRates string `json:"exchange_rates"`
	Salt          string `json:"salt"`
}

func submitAggregateVoteHandlerFunction(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		voter := vars[RestVoter]

		valAddress, err := sdk.ValAddressFromBech32(voter)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		var req AggregateVoteReq
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

		// Check validation of tuples
		_, err = types.ParseExchangeRateTuples(req.ExchangeRates)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// create the message
		msg := types.NewMsgAggregateExchangeRateVote(req.Salt, req.ExchangeRates, fromAddress, valAddress)
		err = msg.ValidateBasic()
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}
