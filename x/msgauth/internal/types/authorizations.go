package types

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// Authorization represents the interface of various Authorization instances
type Authorization interface {
	MsgType() string
	Accept(msg sdk.Msg, block abci.Header) (allow bool, updated Authorization, delete bool)
}

// AuthorizationGrant represent the stored grant instance in the keeper store
type AuthorizationGrant struct {
	Authorization Authorization `json:"authorization"`

	Expiration time.Time `json:"expiration"`
}

// NewAuthorizationGrant returns new AuthroizationGrant instance
func NewAuthorizationGrant(authorization Authorization, expiration time.Time) AuthorizationGrant {
	return AuthorizationGrant{Authorization: authorization, Expiration: expiration}
}

// GGMPair is struct that just has a granter-grantee-msgtype pair with no other data.
// It is intended to be used as a marshalable pointer. For example, a GGPair can be used to construct the
// key to getting an Grant from state.
type GGMPair struct {
	GranterAddress sdk.AccAddress
	GranteeAddress sdk.AccAddress
	MsgType        string
}
