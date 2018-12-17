package market

import (
	"terra/x/oracle"
	"terra/x/treasury"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
)

//nolint
type Keeper struct {
	storeKey  sdk.StoreKey      // Key to our module's store
	codespace sdk.CodespaceType // Reserves space for error codes
	cdc       *codec.Codec      // Codec to encore/decode structs

	bk bank.Keeper   // Read & write terra & luna balance
	ok oracle.Keeper // Read terra & luna prices
	tk treasury.Keeper
}

// NewKeeper crates a new keeper with write and read access
func NewKeeper(
	// cdc *amino.Codec,
	// marketKey sdk.StoreKey,
	bk bank.Keeper,
	ok oracle.Keeper,
	tk treasury.Keeper,
	// codespace sdk.CodespaceType,
) Keeper {
	return Keeper{
		// storeKey:  marketKey,
		// cdc:       cdc,
		bk: bk,
		ok: ok,
		tk: tk,
		// codespace: codespace,
	}
}
