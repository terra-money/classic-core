package rest

import (
	"github.com/cosmos/cosmos-sdk/client"
	clientrest "github.com/cosmos/cosmos-sdk/client/rest"

	"github.com/gorilla/mux"
)

// RestDenom is the wildcard part of the request path
const RestDenom = "denom"

// RegisterRoutes registers market-related REST handlers to a router
func RegisterRoutes(cliCtx client.Context, rtr *mux.Router) {
	r := clientrest.WithHTTPDeprecationHeaders(rtr)

	registerQueryRoutes(cliCtx, r)
	registerTxHandlers(cliCtx, r)
}
