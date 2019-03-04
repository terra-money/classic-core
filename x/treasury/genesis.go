package treasury

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GenesisState - all treasury state that must be provided at genesis
type GenesisState struct {
	Params Params `json:"params"` // treasury params
	//GenesisIssuance map[string]sdk.Int `json:"genesis_issuance"` // genesis issuance for coins
}

func NewGenesisState(params Params /*, genesisIssuance map[string]sdk.Int*/) GenesisState {
	return GenesisState{
		Params: params,
		//GenesisIssuance: genesisIssuance,
	}
}

// get raw genesis raw message for testing
func DefaultGenesisState() GenesisState {
	return GenesisState{
		Params: DefaultParams(),
		// GenesisIssuance: map[string]sdk.Int{
		// 	assets.LunaDenom: sdk.NewInt(10 ^ 9),
		// },
	}
}

// new oracle genesis
func InitGenesis(ctx sdk.Context, keeper Keeper, data GenesisState) {
	// for key, value := range data.GenesisIssuance {
	// 	keeper.SetIssuance(ctx, key, value)
	// }

	keeper.SetParams(ctx, data.Params)
}

// ExportGenesis returns a GenesisState for a given context and keeper. The
// GenesisState will contain the pool, and validator/delegator distribution info's
func ExportGenesis(ctx sdk.Context, k Keeper) GenesisState {
	params := k.GetParams(ctx)
	// store := ctx.KVStore(k.key)
	// iter := sdk.KVStorePrefixIterator(store, PrefixIssuance)

	// genData := map[string]sdk.Int{}
	// for ; iter.Valid(); iter.Next() {

	// 	var denom string
	// 	var issuance sdk.Int
	// 	k.cdc.MustUnmarshalBinaryLengthPrefixed(iter.Key(), &denom)
	// 	k.cdc.MustUnmarshalBinaryLengthPrefixed(iter.Value(), &issuance)

	// 	genData[denom] = issuance
	// }
	// iter.Close()

	return NewGenesisState(params)
}

// ValidateGenesis validates the provided oracle genesis state to ensure the
// expected invariants holds. (i.e. params in correct bounds, no duplicate validators)
func ValidateGenesis(data GenesisState) error {
	// for denom, issuance := range data.GenesisIssuance {
	// 	if issuance.LT(sdk.ZeroInt()) {
	// 		return fmt.Errorf("Genesis issuance cannot be negative for %s", denom)
	// 	}
	// }

	return validateParams(data.Params)
}
