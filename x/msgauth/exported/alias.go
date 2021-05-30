// nolint
package exported

import (
	"time"

	"github.com/terra-money/core/x/msgauth/internal/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Keeper interface {
	//DispatchActions executes the provided messages via authorization grants from the message signer to the grantee
	DispatchActions(ctx sdk.Context, grantee sdk.AccAddress, msgs []sdk.Msg) sdk.Result

	// Grant grants the provided authorization to the grantee on the granter's account with the provided expiration time
	// If there is an existing authorization grant for the same sdk.Msg type, this grant overwrites that.
	Grant(ctx sdk.Context, grantee sdk.AccAddress, granter sdk.AccAddress, authorization types.Authorization, expiration time.Time)

	// Revoke removes any authorization for the provided message type granted to the grantee by the granter.
	Revoke(ctx sdk.Context, grantee sdk.AccAddress, granter sdk.AccAddress, msgType sdk.Msg)

	// GetAuthorizationGrant Returns any Authorization (or nil), with the expiration time,
	// granted to the grantee by the granter for the provided msg type.
	GetAuthorizationGrant(ctx sdk.Context, grantee sdk.AccAddress, granter sdk.AccAddress, msgType sdk.Msg) (grant types.AuthorizationGrant, found bool)
}

var RegisterMsgAuthTypeCodec = types.RegisterMsgAuthTypeCodec

type (
	MsgGrantAuthorization  = types.MsgGrantAuthorization
	MsgRevokeAuthorization = types.MsgRevokeAuthorization
	MsgExecAuthorized      = types.MsgExecAuthorized
)
