package rest

import (
	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client"
)

// Rest Flags
const (
	RestCodeID          = "code_id"
	RestContractAddress = "contract_address"
)

// RegisterRoutes registers staking-related REST handlers to a router
func RegisterRoutes(clientCtx client.Context, rtr *mux.Router) {
	registerQueryRoutes(clientCtx, rtr)
	registerTxRoutes(clientCtx, rtr)
}
