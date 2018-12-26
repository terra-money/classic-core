package rest

import (
	"net/http"
	"terra/x/oracle"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/utils"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/gorilla/mux"
)

//nolint
const (
	RestVoteDenom = "denom"
	RestVoter     = "voteraddress"
	storeName     = "oracle"
)

func registerQueryRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *codec.Codec) {
	r.HandleFunc(
		"/oracle/vote/{denom}/{voteraddress}",
		getVoteHandlerFunction(cdc, storeName, cliCtx),
	).Methods("GET")

	r.HandleFunc("/oracle/current/{denom}",
		getCurrentHandlerFunction(cdc, storeName, cliCtx),
	).Methods("GET")

	r.HandleFunc("/oracle/target/{denom}",
		getTargetHandlerFunction(cdc, storeName, cliCtx),
	).Methods("GET")

	r.HandleFunc("/oracle/whitelist",
		getWhitelistHandlerFunction(cdc, storeName, cliCtx),
	).Methods("GET")
}

//VoteReq ...
type VoteReq struct {
	BaseReq      utils.BaseReq `json:"base_req"`
	TargetPrice  float64       `json:"target_price"`
	CurrentPrice float64       `json:"current_price"`
	Denom        string        `json:"denom"`
	VoterAddress string        `json:"voter_address"`
}

// GetVoteHandlerFunction handles the request to get the currently unelected outstanding price oracle vote
func getVoteHandlerFunction(cdc *codec.Codec, storeName string, cliCtx context.CLIContext) http.HandlerFunc {
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

// Retrieves the current effective price in Luna for the asset
func getCurrentHandlerFunction(cdc *codec.Codec, storeName string, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		denom := vars[RestVoteDenom]

		res, err := cliCtx.QueryStore(oracle.GetObservedPriceKey(denom), storeName)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		utils.PostProcessResponse(w, cdc, res, cliCtx.Indent)
	}
}

// Retrieves the target peg in Luna for the asset
func getTargetHandlerFunction(cdc *codec.Codec, storeName string, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		denom := vars[RestVoteDenom]

		res, err := cliCtx.QueryStore(oracle.GetTargetPriceKey(denom), storeName)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		utils.PostProcessResponse(w, cdc, res, cliCtx.Indent)
	}
}

// Retrieves the oracle whitelist
func getWhitelistHandlerFunction(cdc *codec.Codec, storeName string, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		res, err := cliCtx.QueryStore(oracle.KeyWhitelist, storeName)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		utils.PostProcessResponse(w, cdc, res, cliCtx.Indent)
	}
}
