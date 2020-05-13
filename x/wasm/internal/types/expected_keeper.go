package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/exported"
)

// AccountKeeper - expected account keeper
type AccountKeeper interface {
	NewAccountWithAddress(ctx sdk.Context, addr sdk.AccAddress) exported.Account
	GetAccount(ctx sdk.Context, addr sdk.AccAddress) exported.Account
	SetAccount(ctx sdk.Context, acc exported.Account)
}

// BankKeeper - expected bank keeper
type BankKeeper interface {
	SendCoins(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) sdk.Error
}
