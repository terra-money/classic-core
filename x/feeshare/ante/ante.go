package ante

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	feeshare "github.com/classic-terra/core/x/feeshare/types"
	wasmtypes "github.com/classic-terra/core/x/wasm/types"
)

// FeeSharePayoutDecorator Run his after we already deduct the fee from the account with
// the ante.NewDeductFeeDecorator() decorator. We pull funds from the FeeCollector ModuleAccount
type FeeSharePayoutDecorator struct {
	bankKeeper     BankKeeper
	feesharekeeper FeeShareKeeper
}

func NewFeeSharePayoutDecorator(bk BankKeeper, fs FeeShareKeeper) FeeSharePayoutDecorator {
	return FeeSharePayoutDecorator{
		bankKeeper:     bk,
		feesharekeeper: fs,
	}
}

func (fsd FeeSharePayoutDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
	feeTx, ok := tx.(sdk.FeeTx)
	if !ok {
		return ctx, sdkerrors.Wrap(sdkerrors.ErrTxDecode, "Tx must be a FeeTx")
	}

	err = FeeSharePayout(ctx, fsd.bankKeeper, feeTx.GetFee(), fsd.feesharekeeper, tx.GetMsgs())
	if err != nil {
		return ctx, sdkerrors.Wrapf(sdkerrors.ErrInsufficientFunds, err.Error())
	}

	return next(ctx, tx, simulate)
}

// FeePayLogic takes the total fees and splits them based on the governance params
// and the number of contracts we are executing on.
// This returns the amount of fees each contract developer should get.
// tested in ante_test.go
func FeePayLogic(fees sdk.Coins, govPercent sdk.Dec, numPairs int) sdk.Coins {
	var splitFees sdk.Coins
	for _, c := range fees.Sort() {
		rewardAmount := govPercent.MulInt(c.Amount).QuoInt64(int64(numPairs)).RoundInt()
		if !rewardAmount.IsZero() {
			splitFees = splitFees.Add(sdk.NewCoin(c.Denom, rewardAmount))
		}
	}
	return splitFees
}

// FeeSharePayout takes the total fees and redistributes 50% (or param set) to the contract developers
// provided they opted-in to payments.
func FeeSharePayout(ctx sdk.Context, bk BankKeeper, totalFees sdk.Coins, fsk FeeShareKeeper, msgs []sdk.Msg) error {
	params := fsk.GetParams(ctx)
	if !params.EnableFeeShare {
		return nil
	}

	// Get valid withdraw addresses from contracts
	toPay := make([]sdk.AccAddress, 0)
	for _, msg := range msgs {
		if _, ok := msg.(*wasmtypes.MsgExecuteContract); ok {
			contractAddr, err := sdk.AccAddressFromBech32(msg.(*wasmtypes.MsgExecuteContract).Contract)
			if err != nil {
				return err
			}

			shareData, _ := fsk.GetFeeShare(ctx, contractAddr)

			withdrawAddr := shareData.GetWithdrawerAddr()
			if withdrawAddr != nil && !withdrawAddr.Empty() {
				toPay = append(toPay, withdrawAddr)
			}
		}
	}

	// Do nothing if no one needs payment
	if len(toPay) == 0 {
		return nil
	}

	// Get only allowed governance fees to be paid (helps for taxes)
	var fees sdk.Coins
	if len(params.AllowedDenoms) == 0 {
		// If empty, we allow all denoms to be used as payment
		fees = totalFees
	} else {
		for _, fee := range totalFees.Sort() {
			for _, allowed := range params.AllowedDenoms {
				if fee.Denom == allowed {
					fees = fees.Add(fee)
				}
			}
		}
	}

	// FeeShare logic payouts for contracts
	numPairs := len(toPay)
	if numPairs > 0 {
		govPercent := params.DeveloperShares
		splitFees := FeePayLogic(fees, govPercent, numPairs)

		// pay fees evenly between all withdraw addresses
		for _, withdrawAddr := range toPay {
			err := bk.SendCoinsFromModuleToAccount(ctx, authtypes.FeeCollectorName, withdrawAddr, splitFees)
			if err != nil {
				return sdkerrors.Wrapf(feeshare.ErrFeeSharePayment, "failed to pay fees to contract developer: %s", err.Error())
			}
		}
	}

	return nil
}
