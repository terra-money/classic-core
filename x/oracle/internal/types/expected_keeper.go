package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingexported "github.com/cosmos/cosmos-sdk/x/staking/exported"
	supplyexported "github.com/cosmos/cosmos-sdk/x/supply/exported"
)

// StakingKeeper is expected keeper for staking module
type StakingKeeper interface {
	Validator(ctx sdk.Context, address sdk.ValAddress) stakingexported.ValidatorI // get validator by operator address; nil when validator not found
	TotalBondedTokens(sdk.Context) sdk.Int                                        // total bonded tokens within the validator set
}

// DistributionKeeper is expected keeper for distribution module
type DistributionKeeper interface {
	AllocateTokensToValidator(ctx sdk.Context, val stakingexported.ValidatorI, tokens sdk.DecCoins)
}

// SupplyKeeper is expected keeper for supply module
type SupplyKeeper interface {
	GetModuleAccount(ctx sdk.Context, moduleName string) supplyexported.ModuleAccountI
	SetModuleAccount(sdk.Context, supplyexported.ModuleAccountI)
	GetSupply(ctx sdk.Context) (supply supplyexported.SupplyI)
	SetSupply(ctx sdk.Context, supply supplyexported.SupplyI)
	SendCoinsFromModuleToModule(ctx sdk.Context, senderModule string, recipientModule string, amt sdk.Coins) sdk.Error
}
