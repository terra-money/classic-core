package msgauth

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// InitGenesis register all exported authorization entries
func InitGenesis(ctx sdk.Context, keeper Keeper, data GenesisState) {
	for _, entry := range data.AuthorizationEntries {
		keeper.SetGrant(ctx, entry.Granter, entry.Grantee, AuthorizationGrant{
			Authorization: entry.Authorization,
			Expiration:    entry.Expiration,
		})

		keeper.InsertGrantQueue(ctx, entry.Granter, entry.Grantee,
			entry.Authorization.MsgType(), entry.Expiration)
	}
}

// ExportGenesis writes the current store values
// to a genesis file, which can be imported again
// with InitGenesis
func ExportGenesis(ctx sdk.Context, keeper Keeper) (data GenesisState) {
	var entries []AuthorizationEntry
	keeper.IterateGrants(ctx, func(granter, grantee sdk.AccAddress, grant AuthorizationGrant) bool {
		entries = append(entries, AuthorizationEntry{
			Granter:       granter,
			Grantee:       grantee,
			Expiration:    grant.Expiration,
			Authorization: grant.Authorization,
		})
		return false
	})

	return NewGenesisState(entries)
}
