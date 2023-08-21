package ante

import (
	treasury "github.com/classic-terra/core/v2/x/treasury/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	cosmosante "github.com/cosmos/cosmos-sdk/x/auth/ante"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
)

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
				distributionDeltaCoins := sdk.NewCoins()

				for _, taxCoin := range taxes {
					splitcoinAmount := burnSplitRate.MulInt(taxCoin.Amount).RoundInt()
					distributionDeltaCoins = distributionDeltaCoins.Add(sdk.NewCoin(taxCoin.Denom, splitcoinAmount))
				}

				taxes = taxes.Sub(distributionDeltaCoins...)
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
