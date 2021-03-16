package msgauth

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/terra-project/core/x/msgauth/keeper"
	"github.com/terra-project/core/x/msgauth/types"
)

// InitGenesis register all exported authorization entries
func InitGenesis(ctx sdk.Context, keeper keeper.Keeper, data *types.GenesisState) {
	for _, entry := range data.AuthorizationEntries {
		granter, err := sdk.AccAddressFromBech32(entry.Granter)
		if err != nil {
			panic(fmt.Errorf("Invalid granter address %s", entry.Granter))
		}

		grantee, err := sdk.AccAddressFromBech32(entry.Grantee)
		if err != nil {
			panic(fmt.Errorf("Invalid grantee address %s", entry.Grantee))
		}

		authorization := entry.GetAuthorization()
		grant, err := types.NewAuthorizationGrant(authorization, entry.Expiration)
		if err != nil {
			panic(err)
		}

		keeper.SetGrant(ctx, granter, grantee, authorization.MsgType(), grant)
		keeper.InsertGrantQueue(ctx, granter, grantee, authorization.MsgType(), entry.Expiration)
	}
}

// ExportGenesis writes the current store values
// to a genesis file, which can be imported again
// with InitGenesis
func ExportGenesis(ctx sdk.Context, keeper keeper.Keeper) (data *types.GenesisState) {
	var entries []types.AuthorizationEntry
	keeper.IterateGrants(ctx, func(granter, grantee sdk.AccAddress, grant types.AuthorizationGrant) bool {
		entries = append(entries, types.AuthorizationEntry{
			Granter:       granter.String(),
			Grantee:       grantee.String(),
			Expiration:    grant.Expiration,
			Authorization: grant.Authorization,
		})
		return false
	})

	return types.NewGenesisState(entries)
}
