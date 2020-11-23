package rest

import (
	"time"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"
)

//nolint
const (
	RestGranter = "granter"
	RestGrantee = "grantee"
	RestMsgType = "msg_type"
)

// RegisterRoutes register routes for querier and tx broadcast
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router) {
	registerQueryRoutes(cliCtx, r)
	registerTxRoutes(cliCtx, r)
}

// GrantRequest defines the properties of a grant request's body.
type GrantRequest struct {
	BaseReq rest.BaseReq  `json:"base_req" yaml:"base_req"`
	Period  time.Duration `json:"period"`
	Limit   sdk.Coins     `json:"limit,omitempty"`
}

// RevokeRequest defines the properties of a revoke request's body.
type RevokeRequest struct {
	BaseReq rest.BaseReq `json:"base_req" yaml:"base_req"`
}

// ExecuteRequest defines the properties of a execute request's body.
type ExecuteRequest struct {
	BaseReq rest.BaseReq `json:"base_req" yaml:"base_req"`
	Msgs    []sdk.Msg    `json:"msgs" yaml:"msgs"`
}
