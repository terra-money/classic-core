package msgauth

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// InitGenesis register all exported authorization entries
func InitGenesis(ctx sdk.Context, keeper Keeper, data GenesisState) {
	for _, entry := data.AuthorizationEntries {
		keeper.Grant(ctx, entry.Grantee, entry.Granter, entry.Authorization, entry.Expiration)
	}
}

// ExportGenesis writes the current store values
// to a genesis file, which can be imported again
// with InitGenesis
func ExportGenesis(ctx sdk.Context, keeper Keeper) (data GenesisState) {
	var entries []AuthorizationEntry
	keeper.IterateAuthorization(ctx, func(grantee, granter sdk.AccAddress, authorizationGrant AuthorizationGrant) bool {
		append(entries, AuthorizationEntry{
			Granter: granter,
			Grantee: grantee,
			Expiration: authorizationGrant.Expiration,
			Authorization: authorizationGrant.Authorization,
		})
		return false
	})

	return NewGenesisState(entries)
}
