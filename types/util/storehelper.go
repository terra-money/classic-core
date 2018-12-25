package util

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func Get(dbKey sdk.StoreKey, cdc *codec.Codec, ctx sdk.Context, storeKey []byte) (res Value) {
	store := ctx.KVStore(dbKey)
	bz := store.Get(storeKey)
	if bz == nil {
		return
	}
	cdc.MustUnmarshalBinaryLengthPrefixed(bz, &res)
	return
}

func Collect(dbKey sdk.StoreKey, cdc *codec.Codec, ctx sdk.Context, storePrefix []byte) (res []Value) {
	store := ctx.KVStore(dbKey)
	iter := sdk.KVStorePrefixIterator(store, storePrefix)
	for ; iter.Valid(); iter.Next() {
		var v Value
		cdc.MustUnmarshalBinaryLengthPrefixed(iter.Value(), &v)
		res = append(res, v)
	}
	iter.Close()
	return
}

func Set(dbKey sdk.StoreKey, cdc *codec.Codec, ctx sdk.Context, storeKey []byte, value Value) {
	store := ctx.KVStore(dbKey)
	bz := cdc.MustMarshalBinaryLengthPrefixed(value)
	store.Set(storeKey, bz)
}

func Delete(dbKey sdk.StoreKey, ctx sdk.Context, storeKey []byte) {
	store := ctx.KVStore(dbKey)
	store.Delete(storeKey)
}

func Clear(dbKey sdk.StoreKey, ctx sdk.Context, storePrefix []byte) {
	store := ctx.KVStore(dbKey)
	iter := sdk.KVStorePrefixIterator(store, storePrefix)
	for ; iter.Valid(); iter.Next() {
		store.Delete(iter.Key())
	}
	iter.Close()
}

type Value interface {
}
