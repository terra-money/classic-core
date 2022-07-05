package ante

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	treasury "github.com/terra-money/core/x/treasury/types"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
)

// BurnTaxFeeDecorator will immediately burn the collected Tax
type BurnTaxFeeDecorator struct {
	treasuryKeeper TreasuryKeeper
	bankKeeper BankKeeper
}

// NewBurnTaxFeeDecorator returns new tax fee decorator instance
func NewBurnTaxFeeDecorator(treasuryKeeper TreasuryKeeper, bankKeeper BankKeeper) BurnTaxFeeDecorator {
	return BurnTaxFeeDecorator{
		treasuryKeeper: treasuryKeeper,
		bankKeeper: bankKeeper,
	}
}

// AnteHandle handles msg tax fee checking
func (btfd BurnTaxFeeDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
	type isTaxKey string
	k := isTaxKey("tax")

	feeTx, ok := tx.(sdk.FeeTx)
	if !ok {
		return ctx, sdkerrors.Wrap(sdkerrors.ErrTxDecode, "Tx must be a FeeTx")
	}

	msgs := feeTx.GetMsgs()

	//At this point we have already run the DeductFees AnteHandler and taken the fees from the sending account
	//Now we remove the taxes from the gas reward and immediately burn it

	if !simulate {
		// Compute taxes again.  Slightly redundant
		ctx.Logger().Info(fmt.Sprintf("Value of tax param is %s", ctx.Value("tax")))
		if ctx.Value(k) == true {
			taxes := FilterMsgAndComputeTax(ctx, btfd.treasuryKeeper, msgs...)

			ctx.Logger().Info(fmt.Sprintf("Burning the Tax %s", taxes))
			// Record tax proceeds
			if !taxes.IsZero() {
				ctx.Logger().Info(fmt.Sprintf("Burning the Tax %s", taxes))
				btfd.bankKeeper.SendCoinsFromModuleToModule(ctx, types.FeeCollectorName, treasury.BurnModuleName, taxes)
				if err != nil {
					return ctx, sdkerrors.Wrapf(sdkerrors.ErrInsufficientFunds, err.Error())
				}
			}
		}
	}

	return next(ctx, tx, simulate)
}

