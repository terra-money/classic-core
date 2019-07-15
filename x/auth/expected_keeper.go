package auth

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// TreasuryKeeper is expected keeper for treasury
type TreasuryKeeper interface {
	GetTaxRate(ctx sdk.Context, epoch sdk.Int) (rate sdk.Dec)
	GetTaxCap(ctx sdk.Context, denom string) (taxCap sdk.Int)
	RecordTaxProceeds(ctx sdk.Context, delta sdk.Coins)
}
