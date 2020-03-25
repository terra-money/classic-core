package rest

import (
	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/context"
)

const (
	RestCodeID          = "code_id"
	RestContractAddress = "contract_address"
	RestKey             = "key"
	RestSubkey          = "subkey"
	RestMsg             = "msg"
)

// RegisterRoutes registers staking-related REST handlers to a router
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router) {
	registerQueryRoutes(cliCtx, r)
	registerTxRoutes(cliCtx, r)
}
