package budget

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// expected mint keeper
type MintKeeper interface {
	Mint(ctx sdk.Context, recipient sdk.AccAddress, coin sdk.Coin) (err sdk.Error)
	Burn(ctx sdk.Context, payer sdk.AccAddress, coin sdk.Coin) (err sdk.Error)
	PeekEpochSeigniorage(ctx sdk.Context, epoch sdk.Int) (epochSeigniorage sdk.Int)
}

// expected treasury keeper
type TreasuryKeeper interface {
	GetRewardWeight(ctx sdk.Context, epoch sdk.Int) (rewardWeight sdk.Dec)
	SetRewardWeight(ctx sdk.Context, weight sdk.Dec)
}

// expected market keeper
type MarketKeeper interface {
	GetSwapDecCoin(ctx sdk.Context, offerCoin sdk.DecCoin, askDenom string) (sdk.DecCoin, sdk.Error)
}
