package rest

import (
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/gorilla/mux"
)

// REST Variable names
// nolint
const (
	RestDenom = "denom"
	RestDay   = "day"
	RestEpoch = "epoch"
)

// RegisterRoutes - Central function to define routes that get registered by the main application
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *codec.Codec) {
	registerQueryRoutes(cliCtx, r, cdc)
}
