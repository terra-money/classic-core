package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	supplyexported "github.com/cosmos/cosmos-sdk/x/supply/exported"
)

// expected supply keeper
type SupplyKeeper interface {
	GetSupply(ctx sdk.Context) (supply supplyexported.SupplyI)
	MintCoins(ctx sdk.Context, name string, amt sdk.Coins) sdk.Error
	SendCoinsFromModuleToModule(ctx sdk.Context, senderModule string, recipientModule string, amt sdk.Coins) sdk.Error
}

// expected market keeper
type MarketKeeper interface {
	ComputeInternalSwap(ctx sdk.Context, offerCoin sdk.DecCoin, askDenom string) (sdk.DecCoin, sdk.Error)
}

// expected keeper for staking module
type StakingKeeper interface {
	TotalBondedTokens(sdk.Context) sdk.Int // total bonded tokens within the validator set
}

// expected keeper for distribution module
type DistributionKeeper interface {
	GetFeePool(ctx sdk.Context) (feePool distrtypes.FeePool)
	SetFeePool(ctx sdk.Context, feePool distrtypes.FeePool)
}
