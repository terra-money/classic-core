package rest

import (
	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client"
)

// Defines wildcard part of the request paths
const (
	RestDenom = "denom"
	RestEpoch = "epoch"
)

// RegisterRoutes registers oracle-related REST handlers to a router
func RegisterRoutes(clientCtx client.Context, r *mux.Router) {
	registerQueryRoute(clientCtx, r)
}
