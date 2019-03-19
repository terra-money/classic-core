package rest

import (
	"fmt"
	"math"
	"net/http"
	"terra/x/oracle"

	"github.com/cosmos/cosmos-sdk/client/context"
	clientrest "github.com/cosmos/cosmos-sdk/client/rest"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"

	"github.com/gorilla/mux"
)

//nolint
const (
	RestVoteDenom  = "denom"
	RestVoter      = "voter"
	RestPrice      = "price"
	RestParamsType = "params"

	queryRoute = "oracle"
)

// RegisterRoutes registers staking-related REST handlers to a router
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *codec.Codec) {
	// POST method
	r.HandleFunc("/oracle/vote", submitVoteHandlerFunction(cdc, cliCtx)).Methods("POST")

	// GET method
	r.HandleFunc(fmt.Sprintf("/oracle/votes/{%s}/{%s}", RestVoteDenom, RestVoter), queryVotesHandlerFunction(cdc, cliCtx)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/oracle/price/{%s}", RestPrice), queryPriceHandlerFunction(cdc, cliCtx)).Methods("GET")
	r.HandleFunc("/oracle/active", queryActiveHandlerFunction(cdc, cliCtx)).Methods("GET")

	r.HandleFunc("/oracle/params", queryParamsHandlerFn(cdc, cliCtx)).Methods("GET")
}

//VoteReq ...
type VoteReq struct {
	BaseReq      rest.BaseReq `json:"base_req"`
	Price        float64      `json:"price"`
	Denom        string       `json:"denom"`
	VoterAddress string       `json:"voter_address"`
}

func submitVoteHandlerFunction(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req VoteReq
		if !rest.ReadRESTReq(w, r, cdc, &req) {
			return
		}

		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		fromAddress := cliCtx.GetFromAddress()
		if req.VoterAddress != fromAddress.String() {
			rest.WriteErrorResponse(w, http.StatusUnauthorized, "Must use own address")
			return
		}

		price := sdk.NewDecWithPrec(int64(math.Round(req.Price*100)), 2)

		// create the message
		msg := oracle.NewMsgPriceFeed(req.Denom, price, fromAddress)
		err := msg.ValidateBasic()
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		if req.BaseReq.GenerateOnly {
			clientrest.WriteGenerateStdTxResponse(w, cdc, cliCtx, req.BaseReq, []sdk.Msg{msg})
			return
		}

		clientrest.CompleteAndBroadcastTxREST(w, cliCtx, req.BaseReq, []sdk.Msg{msg}, cdc)
	}
}

func queryVotesHandlerFunction(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		denom := vars[RestVoteDenom]
		voter := vars[RestVoter]

		var voterAddress sdk.AccAddress
		params := oracle.NewQueryVoteParams(voterAddress, denom)

		if len(voter) != 0 {
			voterAddress, err := sdk.AccAddressFromBech32(voter)
			if err != nil {
				return
			}
			params.Voter = voterAddress
		}

		bz, err := cdc.MarshalJSON(params)
		if err != nil {
			return
		}

		res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", queryRoute, oracle.QueryVotes), bz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		rest.PostProcessResponse(w, cdc, res, cliCtx.Indent)
	}
}

func queryPriceHandlerFunction(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		denom := vars[RestVoteDenom]

		res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s", queryRoute, oracle.QueryPrice, denom), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cdc, res, cliCtx.Indent)
	}
}

func queryActiveHandlerFunction(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", queryRoute, oracle.QueryActive), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cdc, res, cliCtx.Indent)
	}
}

func queryParamsHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", queryRoute, oracle.QueryParams), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cdc, res, cliCtx.Indent)
	}
}
