package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	supplyexported "github.com/cosmos/cosmos-sdk/x/supply/exported"

	"github.com/terra-money/core/x/oracle/internal/types"
)

// GetOracleAccount returns oracle ModuleAccount
func (k Keeper) GetOracleAccount(ctx sdk.Context) supplyexported.ModuleAccountI {
	return k.supplyKeeper.GetModuleAccount(ctx, types.ModuleName)
}

// GetRewardPool retrieves the balance of the oracle module account
func (k Keeper) GetRewardPool(ctx sdk.Context) sdk.Coins {
	acc := k.supplyKeeper.GetModuleAccount(ctx, types.ModuleName)
	return acc.GetCoins()
}
