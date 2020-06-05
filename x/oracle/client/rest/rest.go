package rest

import (
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/gorilla/mux"
)

//nolint
const (
	RestDenom  = "denom"
	RestDenom1 = "denom1"
	RestDenom2 = "denom2"
	RestVoter  = "voter"
)

// RegisterRoutes registers oracle-related REST handlers to a router
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router) {
	resgisterTxRoute(cliCtx, r)
	registerQueryRoute(cliCtx, r)
}
