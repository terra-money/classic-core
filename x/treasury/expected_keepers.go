package treasury

import sdk "github.com/cosmos/cosmos-sdk/types"

// expected mint keeper
type MintKeeper interface {
	PeekEpochSeigniorage(ctx sdk.Context, epoch sdk.Int) (seignioragePool sdk.Int)
	Mint(ctx sdk.Context, recipient sdk.AccAddress, coin sdk.Coin) (err sdk.Error)
	GetIssuance(ctx sdk.Context, denom string, day sdk.Int) (issuance sdk.Int)
}

// expected market keeper
type MarketKeeper interface {
	GetSwapDecCoins(ctx sdk.Context, offerCoin sdk.DecCoin, askDenom string) (sdk.DecCoin, sdk.Error)
	GetSwapCoins(ctx sdk.Context, offerCoin sdk.Coin, askDenom string, isInternal bool) (sdk.Coin, sdk.Dec, sdk.Error)
}

// expected coin keeper
type DistributionKeeper interface {
	AllocateTokensToValidator(ctx sdk.Context, val sdk.Validator, tokens sdk.DecCoins)
}

// expected fee keeper
type FeeCollectionKeeper interface {
	AddCollectedFees(ctx sdk.Context, coins sdk.Coins) sdk.Coins
}
