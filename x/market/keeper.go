package market

import (
	"terra/x/oracle"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	amino "github.com/tendermint/go-amino"
)

//nolint
type Keeper struct {
	storeKey  sdk.StoreKey      // Key to our module's store
	codespace sdk.CodespaceType // Reserves space for error codes
	cdc       *codec.Codec      // Codec to encore/decode structs

	bk bank.Keeper              // Read & write terra & luna balance
	ok oracle.Keeper            // Read terra & luna prices
	fk auth.FeeCollectionKeeper // Set & get terra & luna fees
}

// NewKeeper crates a new keeper with write and read access
func NewKeeper(cdc *amino.Codec, marketKey sdk.StoreKey, bk bank.Keeper,
	ok oracle.Keeper, fk auth.FeeCollectionKeeper, codespace sdk.CodespaceType,
	rp ReserveParams, genesisBalance sdk.Coins, initTerraFee sdk.Rat) Keeper {
	k = Keeper{
		storeKey:  marketKey,
		cdc:       cdc,
		bk:        bk,
		ok:        ok,
		fk:        fk,
		codespace: codespace,
	}

	// Starting tx fee must not be zero
	if initTerraFee.IsZero {
		panic(error)
	}

	k.SetIssuanceMeta(genesisBalance)
	k.SetReserveParams(rp)
	k.SetTerraFee(initTerraFee)

	return k
}

// GetReserveParams retrieves the parameters for the reserve
func (mk Keeper) GetReserveParams(ctx sdk.Context) ReserveParams {
	store := ctx.KVStore(mk.key)
	bz := store.Get(GetReserveParamsKey())
	if bz == nil {
		panic(error)
	}

	rp := &(ReserveParams)
	fck.cdc.MustUnmarshalBinary(bz, rp)
	return *rp
}

// SetReserveParams sets parameters for the reserve
func (mk Keeper) SetReserveParams(ctx sdk.Context, rp ReserveParams) {
	store := ctx.KVStore(mk.storeKey)
	bz := k.cdc.MustMarshalBinary(rp)
	store.Set(GetReserveParamsKey(), bz)
}

// GetTerraFee gets the currently effective Terra tx fee
func (mk Keeper) GetTerraFee(ctx sdk.Context) sdk.Rat {
	store := ctx.KVStore(mk.storeKey)
	bz := store.Get(KeyTerraFee)
	if bz == nil {
		panic(error)
	}

	tf := &(sdk.Rat)
	fck.cdc.MustUnmarshalBinary(bz, tf)
	return tf
}

// SetTerraFee records the currently effective Terra tx fee
func (mk Keeper) SetTerraFee(ctx sdk.Context, tf sdk.Rat) {
	store := ctx.KVStore(mk.storeKey)
	bz := k.cdc.MustMarshalBinary(tf)
	store.Set(KeyTerraFee, bz)
}

// GetIssuanceMeta gets the current mint params of Terra and Luna
func (mk Keeper) GetIssuanceMeta(ctx sdk.Context) sdk.Coins {
	store := ctx.KVStore(mk.storeKey)
	bz := store.Get(KeyIssuanceMeta)
	if bz == nil {
		panic(error)
	}

	im := &(sdk.Coins)
	fck.cdc.MustUnmarshalBinary(bz, im)
	return *im
}

// SetIssuanceMeta sets the current mint params of Terra and Luna
func (mk Keeper) SetIssuanceMeta(ctx sdk.Context, im sdk.Coins) {
	store := ctx.KVStore(mk.storeKey)
	bz := k.cdc.MustMarshalBinary(im)
	store.Set(KeyIssuanceMeta, bz)
}
