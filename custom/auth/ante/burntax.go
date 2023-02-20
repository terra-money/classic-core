package ante

import (
	treasury "github.com/classic-terra/core/x/treasury/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	cosmosante "github.com/cosmos/cosmos-sdk/x/auth/ante"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
)

// TaxPowerUpgradeHeight is when taxes are allowed to go into effect
// This will still need a parameter change proposal, but can be activated
// anytime after this height
const TaxPowerUpgradeHeight = 9346889

// BurnTaxFeeDecorator will immediately burn the collected Tax
type BurnTaxFeeDecorator struct {
	accountKeeper  cosmosante.AccountKeeper
	treasuryKeeper TreasuryKeeper
	bankKeeper     BankKeeper
	distrKeeper    DistrKeeper
}

// NewBurnTaxFeeDecorator returns new tax fee decorator instance
func NewBurnTaxFeeDecorator(accountKeeper cosmosante.AccountKeeper, treasuryKeeper TreasuryKeeper, bankKeeper BankKeeper, distrKeeper DistrKeeper) BurnTaxFeeDecorator {
	return BurnTaxFeeDecorator{
		accountKeeper:  accountKeeper,
		treasuryKeeper: treasuryKeeper,
		bankKeeper:     bankKeeper,
		distrKeeper:    distrKeeper,
	}
}

// AnteHandle handles msg tax fee checking
func (btfd BurnTaxFeeDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
	// Do not proceed if you are below this block height
	currHeight := ctx.BlockHeight()
	if currHeight < TaxPowerUpgradeHeight {
		return next(ctx, tx, simulate)
	}

	feeTx, ok := tx.(sdk.FeeTx)
	if !ok {
		return ctx, sdkerrors.Wrap(sdkerrors.ErrTxDecode, "Tx must be a FeeTx")
	}

	msgs := feeTx.GetMsgs()

	// At this point we have already run the DeductFees AnteHandler and taken the fees from the sending account
	// Now we remove the taxes from the gas reward and immediately burn it
	if !simulate {
		// Compute taxes again.
		taxes := FilterMsgAndComputeTax(ctx, btfd.treasuryKeeper, msgs...)

		// Record tax proceeds
		if !taxes.IsZero() {
			burnSplitRate := btfd.treasuryKeeper.GetBurnSplitRate(ctx)

			if burnSplitRate.IsPositive() {
				communityDeltaCoins := sdk.NewCoins()

				for _, taxCoin := range taxes {
					splitcoinAmount := burnSplitRate.MulInt(taxCoin.Amount).RoundInt()
					communityDeltaCoins = communityDeltaCoins.Add(sdk.NewCoin(taxCoin.Denom, splitcoinAmount))
				}

				taxes = taxes.Sub(communityDeltaCoins)

				if err = btfd.distrKeeper.FundCommunityPool(
					ctx,
					communityDeltaCoins,
					btfd.accountKeeper.GetModuleAddress(types.FeeCollectorName),
				); err != nil {
					return ctx, sdkerrors.Wrapf(sdkerrors.ErrInsufficientFunds, err.Error())
				}
			}

			if !taxes.IsZero() {
				if err = btfd.bankKeeper.SendCoinsFromModuleToModule(
					ctx,
					types.FeeCollectorName,
					treasury.BurnModuleName,
					taxes,
				); err != nil {
					return ctx, sdkerrors.Wrapf(sdkerrors.ErrInsufficientFunds, err.Error())
				}
			}
		}
	}

	return next(ctx, tx, simulate)
}
