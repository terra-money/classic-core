package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/terra-project/core/x/supply"
)

// TreasuryKeeper is expected keeper for treasury
type TreasuryKeeper interface {
	GetTaxRate(ctx sdk.Context) (rate sdk.Dec)
	GetTaxCap(ctx sdk.Context, denom string) (taxCap sdk.Int)
	RecordEpochTaxProceeds(ctx sdk.Context, delta sdk.Coins)
}

// SupplyKeeper defines the expected supply Keeper (noalias)
type SupplyKeeper interface {
	SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) sdk.Error
	GetModuleAccount(ctx sdk.Context, moduleName string) supply.ModuleAccountI
	GetModuleAddress(moduleName string) sdk.AccAddress
}
