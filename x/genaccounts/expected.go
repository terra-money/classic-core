package genaccounts

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/exported"
)

// AccountKeeper defines expected account keeper
type AccountKeeper interface {
	NewAccount(sdk.Context, exported.Account) exported.Account
	SetAccount(sdk.Context, exported.Account)
	IterateAccounts(ctx sdk.Context, process func(exported.Account) (stop bool))
}
