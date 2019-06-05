package market

import sdk "github.com/cosmos/cosmos-sdk/types"

// expected oracle keeper
type OracleKeeper interface {
	AddSwapFeePool(ctx sdk.Context, fees sdk.Coins)
	GetLunaSwapRate(ctx sdk.Context, denom string) (price sdk.Dec, err sdk.Error)
}

// expected mint keeper
type MintKeeper interface {
	Mint(ctx sdk.Context, recipient sdk.AccAddress, coin sdk.Coin) (err sdk.Error)
	Burn(ctx sdk.Context, payer sdk.AccAddress, coin sdk.Coin) (err sdk.Error)
	GetIssuance(ctx sdk.Context, denom string, day sdk.Int) (issuance sdk.Int)
}
