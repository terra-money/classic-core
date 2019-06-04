package oracle

import sdk "github.com/cosmos/cosmos-sdk/types"

// expected coin keeper
type DistributionKeeper interface {
	AllocateTokensToValidator(ctx sdk.Context, val sdk.Validator, tokens sdk.DecCoins)
}

// expected fee keeper
type FeeCollectionKeeper interface {
	AddCollectedFees(ctx sdk.Context, coins sdk.Coins) sdk.Coins
}
