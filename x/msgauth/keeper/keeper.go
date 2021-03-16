package keeper

import (
	"bytes"
	"fmt"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/gogo/protobuf/proto"

	"github.com/terra-project/core/x/msgauth/types"
)

// Keeper of the msgauth store
type Keeper struct {
	storeKey sdk.StoreKey
	cdc      codec.BinaryMarshaler
	router   sdk.Router
}

// NewKeeper constructs a message authorization Keeper
func NewKeeper(
	cdc codec.BinaryMarshaler,
	storeKey sdk.StoreKey,
	router sdk.Router,
) Keeper {
	return Keeper{
		cdc:      cdc,
		storeKey: storeKey,
		router:   router,
	}
}

// GetGrant returns grant between granter and grantee for the given msg type
func (k Keeper) GetGrant(ctx sdk.Context, granterAddr sdk.AccAddress, granteeAddr sdk.AccAddress, msgType string) (grant types.AuthorizationGrant, found bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetGrantKey(granterAddr, granteeAddr, msgType))
	if bz == nil {
		return grant, false
	}

	k.cdc.MustUnmarshalBinaryBare(bz, &grant)
	return grant, true
}

// GetGrants returns all the grants between granter and grantee
func (k Keeper) GetGrants(ctx sdk.Context, granterAddr sdk.AccAddress, granteeAddr sdk.AccAddress) (grants types.AuthorizationGrants) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.GetGrantsKey(granterAddr, granteeAddr))
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var authorizationGrant types.AuthorizationGrant
		k.cdc.MustUnmarshalBinaryBare(iter.Value(), &authorizationGrant)
		grants = append(grants, authorizationGrant)
	}

	return grants
}

// GetAllGrants returns all the grants of a granter
func (k Keeper) GetAllGrants(ctx sdk.Context, granterAddr sdk.AccAddress) (grants types.AuthorizationGrants) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.GetAllGrantsKey(granterAddr))
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var authorizationGrant types.AuthorizationGrant
		k.cdc.MustUnmarshalBinaryBare(iter.Value(), &authorizationGrant)
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

			authorization := grant.Authorization.GetCachedValue().(types.AuthorizationI)
			allow, updated, del := authorization.Accept(msg, ctx.BlockHeader())
			if !allow {
				return sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "authorization not found")
			}

			if del {
				k.RevokeGrant(ctx, granterAddr, granteeAddr, msg.Type())
				k.RevokeFromGrantQueue(ctx, granterAddr, granteeAddr, msg.Type(), grant.Expiration)
			} else if updated != nil {
				protoMsg, ok := updated.(proto.Message)
				if !ok {
					return fmt.Errorf("%T does not implement proto.Message", authorization)
				}

				any, err := codectypes.NewAnyWithValue(protoMsg)
				if err != nil {
					return err
				}

				grant.Authorization = any
				k.SetGrant(ctx, granterAddr, granteeAddr, msg.Type(), grant)
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

		var events sdk.Events
		for _, event := range res.Events {
			events = append(events, sdk.Event(event))
		}

		ctx.EventManager().EmitEvents(events)
	}

	return nil
}

// SetGrant method grants the provided authorization to the grantee on the granter's account with the provided expiration
// time. If there is an existing authorization grant for the same `sdk.Msg` type, this grant
// overwrites that.
func (k Keeper) SetGrant(
	ctx sdk.Context, granterAddr, granteeAddr sdk.AccAddress,
	msgType string, grant types.AuthorizationGrant) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryBare(&grant)

	store.Set(types.GetGrantKey(granterAddr, granteeAddr, msgType), bz)
}

// RevokeGrant removes method revokes any authorization for the provided message type granted to the grantee by the granter.
func (k Keeper) RevokeGrant(ctx sdk.Context, granterAddr, granteeAddr sdk.AccAddress, msgType string) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetGrantKey(granterAddr, granteeAddr, msgType))
}

// IterateGrants iterates over all authorization grants
func (k Keeper) IterateGrants(ctx sdk.Context,
	handler func(granterAddr, granteeAddr sdk.AccAddress, grant types.AuthorizationGrant) bool) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.GrantKey)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var grant types.AuthorizationGrant
		granterAddr, granteeAddr := types.ExtractAddressesFromGrantKey(iter.Key())
		k.cdc.MustUnmarshalBinaryBare(iter.Value(), &grant)
		if handler(granterAddr, granteeAddr, grant) {
			break
		}
	}
}

// grant queue timeslice operations

// GetGrantQueueTimeSlice gets a specific grant queue timeslice. A timeslice is a slice of GGMPair
// corresponding to grants that expire at a certain time.
func (k Keeper) GetGrantQueueTimeSlice(ctx sdk.Context, timestamp time.Time) (ggmPairs types.GGMPairs) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetGrantTimeKey(timestamp))
	if bz == nil {
		return types.GGMPairs{}
	}
	k.cdc.MustUnmarshalBinaryBare(bz, &ggmPairs)
	return ggmPairs
}

// SetGrantQueueTimeSlice sets a specific grant queue timeslice.
func (k Keeper) SetGrantQueueTimeSlice(ctx sdk.Context, timestamp time.Time, keys types.GGMPairs) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryBare(&keys)
	store.Set(types.GetGrantTimeKey(timestamp), bz)
}

// InsertGrantQueue inserts an grant to the appropriate timeslice in the grant queue
func (k Keeper) InsertGrantQueue(ctx sdk.Context, granterAddr,
	granteeAddr sdk.AccAddress, msgType string, completionTime time.Time) {

	timeSlice := k.GetGrantQueueTimeSlice(ctx, completionTime)
	ggmPair := types.GGMPair{GranterAddress: granterAddr.String(), GranteeAddress: granteeAddr.String(), MsgType: msgType}
	if len(timeSlice.Pairs) == 0 {
		k.SetGrantQueueTimeSlice(ctx, completionTime, types.GGMPairs{Pairs: []types.GGMPair{ggmPair}})
	} else {
		timeSlice.Pairs = append(timeSlice.Pairs, ggmPair)
		k.SetGrantQueueTimeSlice(ctx, completionTime, timeSlice)
	}
}

// RevokeFromGrantQueue removes grant data from the timeslice queue
func (k Keeper) RevokeFromGrantQueue(ctx sdk.Context, granterAddr,
	granteeAddr sdk.AccAddress, msgType string, completionTime time.Time) {
	timeSlice := k.GetGrantQueueTimeSlice(ctx, completionTime)

	granterAddrStr := granterAddr.String()
	granteeAddrStr := granteeAddr.String()

	for idx, ggmPair := range timeSlice.Pairs {
		if ggmPair.GranterAddress == granterAddrStr &&
			ggmPair.GranteeAddress == granteeAddrStr &&
			ggmPair.MsgType == msgType {

			lastIdx := len(timeSlice.Pairs) - 1
			timeSlice.Pairs[idx] = timeSlice.Pairs[lastIdx]
			timeSlice.Pairs = timeSlice.Pairs[:lastIdx]

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
func (k Keeper) DequeueAllMatureGrantQueue(ctx sdk.Context) (matureGrants types.GGMPairs) {
	store := ctx.KVStore(k.storeKey)
	// gets an iterator for all timeslices from time 0 until the current Blockheader time
	grantTimesliceIterator := k.GrantQueueIterator(ctx, ctx.BlockHeader().Time)
	for ; grantTimesliceIterator.Valid(); grantTimesliceIterator.Next() {
		timeslice := types.GGMPairs{}
		value := grantTimesliceIterator.Value()
		k.cdc.MustUnmarshalBinaryBare(value, &timeslice)
		matureGrants.Pairs = append(matureGrants.Pairs, timeslice.Pairs...)
		store.Delete(grantTimesliceIterator.Key())
	}
	return matureGrants
}
