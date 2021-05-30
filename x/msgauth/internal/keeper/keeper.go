package keeper

import (
	"bytes"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/terra-money/core/x/msgauth/internal/types"
)

// Keeper of the msgauth store
type Keeper struct {
	cdc             *codec.Codec
	storeKey        sdk.StoreKey
	router          sdk.Router
	allowedMsgTypes []string
}

// NewKeeper constructs a message authorisation Keeper
func NewKeeper(cdc *codec.Codec, storeKey sdk.StoreKey, router sdk.Router, allowedMsgTypes ...string) Keeper {
	return Keeper{
		cdc:             cdc,
		storeKey:        storeKey,
		router:          router,
		allowedMsgTypes: allowedMsgTypes,
	}
}

// IsGrantable returns the flag that the given msg type is grantable or not
func (k Keeper) IsGrantable(msgType string) bool {
	for _, mt := range k.allowedMsgTypes {
		if mt == msgType {
			return true
		}
	}

	return false
}

// GetGrant returns grant between granter and grantee for the given msg type
func (k Keeper) GetGrant(ctx sdk.Context, granterAddr sdk.AccAddress, granteeAddr sdk.AccAddress, msgType string) (grant types.AuthorizationGrant, found bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetGrantKey(granterAddr, granteeAddr, msgType))
	if bz == nil {
		return grant, false
	}

	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &grant)
	return grant, true
}

// GetGrants returns all the grants between granter and grantee
func (k Keeper) GetGrants(ctx sdk.Context, granterAddr sdk.AccAddress, granteeAddr sdk.AccAddress) (grants []types.AuthorizationGrant) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.GetGrantKey(granterAddr, granteeAddr, ""))
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var authorizationGrant types.AuthorizationGrant
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iter.Value(), &authorizationGrant)
		grants = append(grants, authorizationGrant)
	}

	return grants
}

// DispatchActions attempts to execute the provided messages via authorization
// grants from the message signer to the grantee.
func (k Keeper) DispatchActions(ctx sdk.Context, granteeAddr sdk.AccAddress, msgs []sdk.Msg) error {
	for _, msg := range msgs {
		signers := msg.GetSigners()
		if len(signers) != 1 {
			return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "authorization can be given to msg with only one signer")
		}
		granterAddr := signers[0]
		if !bytes.Equal(granterAddr, granteeAddr) {
			grant, found := k.GetGrant(ctx, granterAddr, granteeAddr, msg.Type())
			if !found {
				return sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "authorization not found")
			}

			allow, updated, del := grant.Authorization.Accept(msg, ctx.BlockHeader())
			if !allow {
				return sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "authorization not found")
			}

			if del {
				k.RevokeGrant(ctx, granterAddr, granteeAddr, msg.Type())
				k.RevokeFromGrantQueue(ctx, granterAddr, granteeAddr, msg.Type(), grant.Expiration)
			} else if updated != nil {
				grant.Authorization = updated
				k.SetGrant(ctx, granterAddr, granteeAddr, grant)
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

// SetGrant method grants the provided authorization to the grantee on the granter's account with the provided expiration
// time. If there is an existing authorization grant for the same `sdk.Msg` type, this grant
// overwrites that.
func (k Keeper) SetGrant(ctx sdk.Context, granterAddr sdk.AccAddress, granteeAddr sdk.AccAddress, grant types.AuthorizationGrant) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(grant)
	store.Set(types.GetGrantKey(granterAddr, granteeAddr, grant.Authorization.MsgType()), bz)
}

// RevokeGrant removes method revokes any authorization for the provided message type granted to the grantee by the granter.
func (k Keeper) RevokeGrant(ctx sdk.Context, granterAddr sdk.AccAddress, granteeAddr sdk.AccAddress, msgType string) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetGrantKey(granterAddr, granteeAddr, msgType))
}

// IterateGrants iterates over all authorization grants
func (k Keeper) IterateGrants(ctx sdk.Context,
	handler func(granterAddr sdk.AccAddress, granteeAddr sdk.AccAddress, grant types.AuthorizationGrant) bool) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.GrantKey)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var grant types.AuthorizationGrant
		granterAddr, granteeAddr := types.ExtractAddressesFromGrantKey(iter.Key())
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iter.Value(), &grant)
		if handler(granterAddr, granteeAddr, grant) {
			break
		}
	}
}

// grant queue timeslice operations

// GetGrantQueueTimeSlice gets a specific grant queue timeslice. A timeslice is a slice of GGMPair
// corresponding to grants that expire at a certain time.
func (k Keeper) GetGrantQueueTimeSlice(ctx sdk.Context, timestamp time.Time) (ggmPairs []types.GGMPair) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetGrantTimeKey(timestamp))
	if bz == nil {
		return []types.GGMPair{}
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &ggmPairs)
	return ggmPairs
}

// SetGrantQueueTimeSlice sets a specific grant queue timeslice.
func (k Keeper) SetGrantQueueTimeSlice(ctx sdk.Context, timestamp time.Time, keys []types.GGMPair) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(keys)
	store.Set(types.GetGrantTimeKey(timestamp), bz)
}

// InsertGrantQueue inserts an grant to the appropriate timeslice in the grant queue
func (k Keeper) InsertGrantQueue(ctx sdk.Context, granterAddr,
	granteeAddr sdk.AccAddress, msgType string, completionTime time.Time) {

	timeSlice := k.GetGrantQueueTimeSlice(ctx, completionTime)
	ggmPair := types.GGMPair{GranterAddress: granterAddr, GranteeAddress: granteeAddr, MsgType: msgType}
	if len(timeSlice) == 0 {
		k.SetGrantQueueTimeSlice(ctx, completionTime, []types.GGMPair{ggmPair})
	} else {
		timeSlice = append(timeSlice, ggmPair)
		k.SetGrantQueueTimeSlice(ctx, completionTime, timeSlice)
	}
}

// RevokeFromGrantQueue removes grant data from the timeslice queue
func (k Keeper) RevokeFromGrantQueue(ctx sdk.Context, granterAddr,
	granteeAddr sdk.AccAddress, msgType string, completionTime time.Time) {
	timeSlice := k.GetGrantQueueTimeSlice(ctx, completionTime)
	for idx, ggmPair := range timeSlice {
		if ggmPair.GranterAddress.Equals(granterAddr) &&
			ggmPair.GranteeAddress.Equals(granteeAddr) &&
			ggmPair.MsgType == msgType {

			lastIdx := len(timeSlice) - 1
			timeSlice[idx] = timeSlice[lastIdx]
			timeSlice = timeSlice[:lastIdx]

			k.SetGrantQueueTimeSlice(ctx, completionTime, timeSlice)
			return
		}
	}
}

// GrantQueueIterator returns all the grant queue timeslices from time 0 until endTime
func (k Keeper) GrantQueueIterator(ctx sdk.Context, endTime time.Time) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return store.Iterator(types.GrantQueueKey,
		sdk.InclusiveEndBytes(types.GetGrantTimeKey(endTime)))
}

// DequeueAllMatureGrantQueue returns a concatenated list of all the timeslices inclusively previous to
// current block time, and deletes the timeslices from the queue
func (k Keeper) DequeueAllMatureGrantQueue(ctx sdk.Context) (matureGrants []types.GGMPair) {
	store := ctx.KVStore(k.storeKey)
	// gets an iterator for all timeslices from time 0 until the current Blockheader time
	grantTimesliceIterator := k.GrantQueueIterator(ctx, ctx.BlockHeader().Time)
	for ; grantTimesliceIterator.Valid(); grantTimesliceIterator.Next() {
		timeslice := []types.GGMPair{}
		value := grantTimesliceIterator.Value()
		k.cdc.MustUnmarshalBinaryLengthPrefixed(value, &timeslice)
		matureGrants = append(matureGrants, timeslice...)
		store.Delete(grantTimesliceIterator.Key())
	}
	return matureGrants
}
