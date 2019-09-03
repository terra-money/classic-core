package rest

import (
	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/context"
)

const (
	RestDenom = "denom"
	RestEpoch = "epoch"
)

// RegisterRoutes registers oracle-related REST handlers to a router
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router) {
	// resgisterTxRoute(cliCtx, r)
	registerQueryRoute(cliCtx, r)
}
