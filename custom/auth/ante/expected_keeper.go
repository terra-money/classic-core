package ante

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"

	feesharettypes "github.com/classic-terra/core/x/feeshare/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

// TreasuryKeeper for tax charging & recording
type TreasuryKeeper interface {
	RecordEpochTaxProceeds(ctx sdk.Context, delta sdk.Coins)
	GetTaxRate(ctx sdk.Context) (taxRate sdk.Dec)
	GetTaxCap(ctx sdk.Context, denom string) (taxCap sdk.Int)
	GetBurnSplitRate(ctx sdk.Context) sdk.Dec
	HasBurnTaxExemptionAddress(ctx sdk.Context, addresses ...string) bool
	GetMinInitialDepositRatio(ctx sdk.Context) sdk.Dec
}

// OracleKeeper for feeder validation
type OracleKeeper interface {
	ValidateFeeder(ctx sdk.Context, feederAddr sdk.AccAddress, validatorAddr sdk.ValAddress) error
}

// BankKeeper defines the contract needed for supply related APIs (noalias)
type BankKeeper interface {
	SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	SendCoinsFromModuleToModule(ctx sdk.Context, senderModule string, recipientModule string, amt sdk.Coins) error
	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
}

type FeeShareKeeper interface {
	GetParams(ctx sdk.Context) (params feesharettypes.Params)
	GetFeeShare(ctx sdk.Context, contract sdk.Address) (feesharettypes.FeeShare, bool)
}

type DistrKeeper interface {
	FundCommunityPool(ctx sdk.Context, amount sdk.Coins, sender sdk.AccAddress) error
	GetFeePool(ctx sdk.Context) distributiontypes.FeePool
}

type GovKeeper interface {
	GetDepositParams(ctx sdk.Context) govtypes.DepositParams
}
