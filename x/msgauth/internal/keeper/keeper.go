package keeper

import (
	"bytes"
	"time"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/terra-project/core/x/msgauth/internal/types"
)

type Keeper struct {
	cdc      *codec.Codec
	storeKey sdk.StoreKey
	router   baseapp.Router
}

// NewKeeper constructs a message authorisation Keeper
func NewKeeper(cdc *codec.Codec, storeKey sdk.StoreKey, router baseapp.Router) Keeper {
	return Keeper{
		cdc:      cdc,
		storeKey: storeKey,
		router:   router,
	}
}

func (k Keeper) getAuthorizationGrant(ctx sdk.Context, actor []byte) (grant types.AuthorizationGrant, found bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(actor)
	if bz == nil {
		return grant, false
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &grant)
	return grant, true
}

func (k Keeper) update(ctx sdk.Context, grantee sdk.AccAddress, granter sdk.AccAddress, updated types.Authorization) {
	actor := types.GetAuthorizationKey(grantee, granter, updated.MsgType())
	grant, found := k.getAuthorizationGrant(ctx, actor)
	if !found {
		return
	}
	grant.Authorization = updated
	store := ctx.KVStore(k.storeKey)
	store.Set(actor, k.cdc.MustMarshalBinaryLengthPrefixed(grant))
}

// DispatchActions attempts to execute the provided messages via authorization
// grants from the message signer to the grantee.
func (k Keeper) DispatchActions(ctx sdk.Context, grantee sdk.AccAddress, msgs []sdk.Msg) error {
	for _, msg := range msgs {
		signers := msg.GetSigners()
		if len(signers) != 1 {
			return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "authorization can be given to msg with only one signer")
		}
		granter := signers[0]
		if !bytes.Equal(granter, grantee) {
			authorization, _ := k.GetAuthorization(ctx, grantee, granter, msg.Type())
			if authorization == nil {
				return sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "authorization not found")
			}
			allow, updated, del := authorization.Accept(msg, ctx.BlockHeader())
			if !allow {
				return sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "authorization not found")
			}
			if del {
				k.Revoke(ctx, grantee, granter, msg.Type())
			} else if updated != nil {
				k.update(ctx, grantee, granter, updated)
			}
		}

		handler := k.router.Route(ctx, msg.Route())
		if handler == nil {
			return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized message route: %s", msg.Route())
		}

		res, err := handler(ctx, msg)
		if err != nil {
			return sdkerrors.Wrapf(err, "failed to execute message; message %s", msg.Type())
		}

		ctx.EventManager().EmitEvents(res.Events)
	}

	return nil
}

// Grant method grants the provided authorization to the grantee on the granter's account with the provided expiration
// time. If there is an existing authorization grant for the same `sdk.Msg` type, this grant
// overwrites that.
func (k Keeper) Grant(ctx sdk.Context, grantee sdk.AccAddress, granter sdk.AccAddress, authorization types.Authorization, expiration time.Time) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(types.AuthorizationGrant{Authorization: authorization, Expiration: expiration.Unix()})
	actor := types.GetAuthorizationKey(grantee, granter, authorization.MsgType())
	store.Set(actor, bz)
}

// Revoke method revokes any authorization for the provided message type granted to the grantee by the granter.
func (k Keeper) Revoke(ctx sdk.Context, grantee sdk.AccAddress, granter sdk.AccAddress, msgType string) error {
	store := ctx.KVStore(k.storeKey)
	actor := types.GetAuthorizationKey(grantee, granter, msgType)
	_, found := k.getAuthorizationGrant(ctx, actor)
	if !found {
		return sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "authorization not found")
	}
	store.Delete(actor)

	return nil
}

// GetAuthorization Returns any `Authorization` (or `nil`), with the expiration time,
// granted to the grantee by the granter for the provided msg type.
func (k Keeper) GetAuthorization(ctx sdk.Context, grantee sdk.AccAddress, granter sdk.AccAddress, msgType string) (cap types.Authorization, expiration int64) {
	grant, found := k.getAuthorizationGrant(ctx, types.GetAuthorizationKey(grantee, granter, msgType))
	if !found {
		return nil, 0
	}

	if grant.Expiration != 0 && grant.Expiration < (ctx.BlockHeader().Time.Unix()) {
		k.Revoke(ctx, grantee, granter, msgType)
		return nil, 0
	}

	return grant.Authorization, grant.Expiration
}

// IterateAuthorization iterates over all authorization grants
func (k Keeper) IterateAuthorization(ctx sdk.Context,
	handler func(grantee sdk.AccAddress, granter sdk.AccAddress, authorizationGrant types.AuthorizationGrant) bool) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.AuthorizationKey)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var authorizationGrant types.AuthorizationGrant
		granteeAddr, granterAddr := types.ExtractAddressesFromAuthorizationKey(iter.Key())
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iter.Value(), &authorizationGrant)
		if handler(granteeAddr, granterAddr, authorizationGrant) {
			break
		}
	}
}
