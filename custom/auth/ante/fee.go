package ante

import (
	"fmt"
	"math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
)

// FeeDecorator deducts fees from the first signer of the tx
// If the first signer does not have the funds to pay for the fees, return with InsufficientFunds error
// Call next AnteHandler if fees successfully deducted
// CONTRACT: Tx must implement FeeTx interface to use DeductFeeDecorator
type FeeDecorator struct {
	accountKeeper  ante.AccountKeeper
	bankKeeper     BankKeeper
	feegrantKeeper ante.FeegrantKeeper
	treasuryKeeper TreasuryKeeper
}

func NewFeeDecorator(ak ante.AccountKeeper, bk BankKeeper, fk ante.FeegrantKeeper, tk TreasuryKeeper) FeeDecorator {
	return FeeDecorator{
		accountKeeper:  ak,
		bankKeeper:     bk,
		feegrantKeeper: fk,
		treasuryKeeper: tk,
	}
}

func (fd FeeDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (sdk.Context, error) {
	feeTx, ok := tx.(sdk.FeeTx)
	if !ok {
		return ctx, sdkerrors.Wrap(sdkerrors.ErrTxDecode, "Tx must be a FeeTx")
	}

	if !simulate && ctx.BlockHeight() > 0 && feeTx.GetGas() == 0 {
		return ctx, sdkerrors.Wrap(sdkerrors.ErrInvalidGasLimit, "must provide positive gas")
	}

	var (
		priority int64
		err      error
	)

	msgs := feeTx.GetMsgs()
	// Compute taxes
	taxes := FilterMsgAndComputeTax(ctx, fd.treasuryKeeper, msgs...)

	if !simulate {
		priority, err = fd.checkTxFee(ctx, tx, taxes)
		if err != nil {
			return ctx, err
		}
	}

	if err := fd.checkDeductFee(ctx, feeTx, taxes, simulate); err != nil {
		return ctx, err
	}

	newCtx := ctx.WithPriority(priority)

	return next(newCtx, tx, simulate)
}

func (fd FeeDecorator) checkDeductFee(ctx sdk.Context, feeTx sdk.FeeTx, taxes sdk.Coins, simulate bool) error {
	if addr := fd.accountKeeper.GetModuleAddress(types.FeeCollectorName); addr == nil {
		return fmt.Errorf("fee collector module account (%s) has not been set", types.FeeCollectorName)
	}

	fee := feeTx.GetFee()
	feePayer := feeTx.FeePayer()
	feeGranter := feeTx.FeeGranter()
	deductFeesFrom := feePayer

	// if feegranter set deduct fee from feegranter account.
	// this works with only when feegrant enabled.
	if feeGranter != nil {
		if fd.feegrantKeeper == nil {
			return sdkerrors.ErrInvalidRequest.Wrap("fee grants are not enabled")
		} else if !feeGranter.Equals(feePayer) {
			err := fd.feegrantKeeper.UseGrantedFees(ctx, feeGranter, feePayer, fee, feeTx.GetMsgs())
			if err != nil {
				return sdkerrors.Wrapf(err, "%s does not not allow to pay fees for %s", feeGranter, feePayer)
			}
		}

		deductFeesFrom = feeGranter
	}

	deductFeesFromAcc := fd.accountKeeper.GetAccount(ctx, deductFeesFrom)
	if deductFeesFromAcc == nil {
		return sdkerrors.ErrUnknownAddress.Wrapf("fee payer address: %s does not exist", deductFeesFrom)
	}

	// deduct the fees
	if !fee.IsZero() {
		err := DeductFees(fd.bankKeeper, ctx, deductFeesFromAcc, fee)
		if err != nil {
			return err
		}

		if !taxes.IsZero() && !simulate {
			err := fd.BurnTaxSplit(ctx, taxes)
			if err != nil {
				return err
			}

			// Record tax proceeds
			fd.treasuryKeeper.RecordEpochTaxProceeds(ctx, taxes)
		}
	}

	events := sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeTx,
			sdk.NewAttribute(sdk.AttributeKeyFee, fee.String()),
			sdk.NewAttribute(sdk.AttributeKeyFeePayer, deductFeesFrom.String()),
		),
	}
	ctx.EventManager().EmitEvents(events)

	return nil
}

// DeductFees deducts fees from the given account.
func DeductFees(bankKeeper types.BankKeeper, ctx sdk.Context, acc types.AccountI, fees sdk.Coins) error {
	if !fees.IsValid() {
		return sdkerrors.Wrapf(sdkerrors.ErrInsufficientFee, "invalid fee amount: %s", fees)
	}

	err := bankKeeper.SendCoinsFromAccountToModule(ctx, acc.GetAddress(), types.FeeCollectorName, fees)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInsufficientFunds, err.Error())
	}

	return nil
}

// checkTxFee implements the default fee logic, where the minimum price per
// unit of gas is fixed and set by each validator, can the tx priority is computed from the gas price.
// Transaction with only oracle messages will skip gas fee check and will have the most priority.
// It also checks enough fee for treasury tax
func (fd FeeDecorator) checkTxFee(ctx sdk.Context, tx sdk.Tx, taxes sdk.Coins) (int64, error) {
	feeTx, ok := tx.(sdk.FeeTx)
	if !ok {
		return 0, sdkerrors.Wrap(sdkerrors.ErrTxDecode, "Tx must be a FeeTx")
	}

	feeCoins := feeTx.GetFee()
	gas := feeTx.GetGas()
	msgs := feeTx.GetMsgs()
	isOracleTx := isOracleTx(msgs)

	// Ensure that the provided fees meet a minimum threshold for the validator,
	// if this is a CheckTx. This is only for local mempool purposes, and thus
	// is only ran on check tx.
	if ctx.IsCheckTx() && !isOracleTx {
		requiredGasFees := sdk.Coins{}
		minGasPrices := ctx.MinGasPrices()
		if !minGasPrices.IsZero() {
			requiredGasFees = make(sdk.Coins, len(minGasPrices))

			// Determine the required fees by multiplying each required minimum gas
			// price by the gas limit, where fee = ceil(minGasPrice * gasLimit).
			glDec := sdk.NewDec(int64(gas))
			for i, gp := range minGasPrices {
				fee := gp.Amount.Mul(glDec)
				requiredGasFees[i] = sdk.NewCoin(gp.Denom, fee.Ceil().RoundInt())
			}
		}

		requiredFees := requiredGasFees.Add(taxes...)

		// Check required fees
		if !requiredFees.IsZero() && !feeCoins.IsAnyGTE(requiredFees) {
			return 0, sdkerrors.Wrapf(sdkerrors.ErrInsufficientFee, "insufficient fees; got: %q, required: %q = %q(gas) + %q(stability)", feeCoins, requiredFees, requiredGasFees, taxes)
		}
	}

	priority := int64(math.MaxInt64)

	if !isOracleTx {
		priority = getTxPriority(feeCoins, int64(gas))
	}

	return priority, nil
}

// getTxPriority returns a naive tx priority based on the amount of the smallest denomination of the gas price
// provided in a transaction.
// NOTE: This implementation should be used with a great consideration as it opens potential attack vectors
// where txs with multiple coins could not be prioritize as expected.
func getTxPriority(fee sdk.Coins, gas int64) int64 {
	var priority int64
	for _, c := range fee {
		p := int64(math.MaxInt64)
		gasPrice := c.Amount.QuoRaw(gas)
		if gasPrice.IsInt64() {
			p = gasPrice.Int64()
		}
		if priority == 0 || p < priority {
			priority = p
		}
	}

	return priority
}
