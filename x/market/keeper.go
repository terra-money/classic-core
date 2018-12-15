package market

import (
	"terra/x/oracle"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	amino "github.com/tendermint/go-amino"
)

//nolint
type Keeper struct {
	storeKey  sdk.StoreKey      // Key to our module's store
	codespace sdk.CodespaceType // Reserves space for error codes
	cdc       *codec.Codec      // Codec to encore/decode structs

	bk bank.Keeper   // Read & write terra & luna balance
	ok oracle.Keeper // Read terra & luna prices
}

// NewKeeper crates a new keeper with write and read access
func NewKeeper(
	cdc *amino.Codec,
	marketKey sdk.StoreKey,
	bk bank.Keeper,
	ok oracle.Keeper,
	codespace sdk.CodespaceType,
) Keeper {
	return Keeper{
		storeKey:  marketKey,
		cdc:       cdc,
		bk:        bk,
		ok:        ok,
		codespace: codespace,
	}
}

// GetCoinSupply retrieves the current total issuance of the coin
func (mk Keeper) GetCoinSupply(ctx sdk.Context, denom string) (res sdk.Int) {
	store := ctx.KVStore(mk.storeKey)
	bz := store.Get(GetCoinSupplyKey(denom))
	if bz == nil {
		return
	}
	mk.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &res)
	return
}

// SetCoinSupply records the current total supply of the coin
func (mk Keeper) SetCoinSupply(ctx sdk.Context, denom string, issuance sdk.Int) {
	store := ctx.KVStore(mk.storeKey)
	bz := mk.cdc.MustMarshalBinaryLengthPrefixed(issuance)
	store.Set(GetCoinSupplyKey(denom), bz)
}
