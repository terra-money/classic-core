package rest

import (
	"math"
	"net/http"
	"terra/x/oracle"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/utils"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gorilla/mux"
)

//nolint
const (
	RestVoteDenom = "denom"
	RestVoter     = "voteraddress"
	storeName     = "oracle"
)

// RegisterRoutes - Central function to define routes that get registered by the main application
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *codec.Codec) {
	// GET /vote/{denom}/{voteraddress}
	r.HandleFunc("oracle/vote/{denom}/{voteraddress}",
		GetVoteHandlerFunction(cdc, storeName, cliCtx)).Methods("GET")
	// POST /vote/{denom}
	r.HandleFunc("oracle/vote/{denom}",
		SubmitVoteHandlerFunction(cdc, cliCtx)).Methods("POST")
	// GET /elect/{denom}
	r.HandleFunc("oracle/elect/{denom}",
		GetElectHandlerFunction(cdc, storeName, cliCtx)).Methods("GET")
}

//nolint
type VoteReq struct {
	BaseReq      utils.BaseReq `json:"base_req"`
	Price        float64       `json:"price"`
	Denom        string        `json:"denom"`
	VoterAddress string        `json:"voter_address"`
}

// GetVoteHandlerFunction handles the request to get the currently unelected outstanding price oracle vote
func GetVoteHandlerFunction(cdc *codec.Codec, storeName string, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		denom := vars[RestVoteDenom]
		voter := vars[RestVoter]

		voterAcc, err := cliCtx.GetAccount([]byte(voter))
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		res, err := cliCtx.QueryStore(oracle.GetVoteKey(denom, voterAcc.GetAddress()), storeName)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		utils.PostProcessResponse(w, cdc, res, cliCtx.Indent)
	}
}

// SubmitVoteHandlerFunction handles a POST vote request
func SubmitVoteHandlerFunction(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req VoteReq
		err := utils.ReadRESTReq(w, r, cdc, &req)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		price := sdk.NewDecWithPrec(int64(math.Round(req.Price*100)), 2)

		voterAcc, err := cliCtx.GetAccount([]byte(req.VoterAddress))
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// create the message
		msg := oracle.NewPriceFeedMsg(req.Denom, price, voterAcc.GetAddress())
		err = msg.ValidateBasic()
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.CompleteAndBroadcastTxREST(w, r, cliCtx, baseReq, []sdk.Msg{msg}, cdc)
	}
}

// GetElectHandlerFunction handles the GET request for the currently valid Price elect
func GetElectHandlerFunction(cdc *codec.Codec, storeName string, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		denom := vars[RestVoteDenom]

		res, err := cliCtx.QueryStore(oracle.GetElectKey(denom), storeName)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		utils.PostProcessResponse(w, cdc, res, cliCtx.Indent)
	}
}
