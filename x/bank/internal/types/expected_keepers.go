package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	authexported "github.com/cosmos/cosmos-sdk/x/auth/exported"
	supplyexported "github.com/cosmos/cosmos-sdk/x/supply/exported"
)

// AccountKeeper defines the account contract that must be fulfilled when
// creating a x/bank keeper.
type AccountKeeper interface {
	NewAccountWithAddress(ctx sdk.Context, addr sdk.AccAddress) authexported.Account

	GetAccount(ctx sdk.Context, addr sdk.AccAddress) authexported.Account
	GetAllAccounts(ctx sdk.Context) []authexported.Account
	SetAccount(ctx sdk.Context, acc authexported.Account)

	IterateAccounts(ctx sdk.Context, process func(authexported.Account) bool)
}

type SupplyKeeper interface {
	BurnCoins(ctx sdk.Context, name string, amt sdk.Coins) error
	GetModuleAccount(ctx sdk.Context, moduleName string) supplyexported.ModuleAccountI
}
