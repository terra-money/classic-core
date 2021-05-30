package ante

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/bank"

	core "github.com/terra-money/core/types"
	marketexported "github.com/terra-money/core/x/market/exported"
	msgauthexported "github.com/terra-money/core/x/msgauth/exported"
	oracleexported "github.com/terra-money/core/x/oracle/exported"
	wasmexported "github.com/terra-money/core/x/wasm/exported"
)

// FeeTx defines the interface to be implemented by Tx to use the FeeDecorators
type FeeTx interface {
	sdk.Tx
	GetGas() uint64
	GetFee() sdk.Coins
	FeePayer() sdk.AccAddress
}

// TaxFeeDecorator will check if the transaction's fee is at least as large
// as tax + the local validator's minimum gasFee (defined in validator config)
// and record tax proceeds to treasury module to track tax proceeds.
// If fee is too low, decorator returns error and tx is rejected from mempool.
// Note this only applies when ctx.CheckTx = true
// If fee is high enough or not CheckTx, then call next AnteHandler
// CONTRACT: Tx must implement FeeTx to use MempoolFeeDecorator
type TaxFeeDecorator struct {
	treasuryKeeper TreasuryKeeper
}

// NewTaxFeeDecorator returns new tax fee decorator instance
func NewTaxFeeDecorator(treasuryKeeper TreasuryKeeper) TaxFeeDecorator {
	return TaxFeeDecorator{
		treasuryKeeper: treasuryKeeper,
	}
}

// AnteHandle handles msg tax fee checking
func (tfd TaxFeeDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
	feeTx, ok := tx.(FeeTx)
	if !ok {
		return ctx, sdkerrors.Wrap(sdkerrors.ErrTxDecode, "Tx must be a FeeTx")
	}

	feeCoins := feeTx.GetFee()
	gas := feeTx.GetGas()

	if !simulate {
		// Compute taxes
		taxes := FilterMsgAndComputeTax(ctx, tfd.treasuryKeeper, feeTx.GetMsgs())

		// Mempool fee validation
		if ctx.IsCheckTx() && !(isOracleTx(ctx, feeTx.GetMsgs()) && gas <= 1000000) {
			if err := EnsureSufficientMempoolFees(ctx, gas, feeCoins, taxes); err != nil {
				return ctx, sdkerrors.Wrapf(sdkerrors.ErrInsufficientFee, err.Error())
			}
		}

		// Ensure paid fee is enough to cover taxes
		if _, hasNeg := feeCoins.SafeSub(taxes); hasNeg {
			return ctx, sdkerrors.Wrapf(sdkerrors.ErrInsufficientFee, "insufficient fees; got: %s required: %s", feeCoins, taxes)
		}

		// Record tax proceeds
		if !taxes.IsZero() {
			tfd.treasuryKeeper.RecordEpochTaxProceeds(ctx, taxes)
		}
	}

	return next(ctx, tx, simulate)
}

// EnsureSufficientMempoolFees verifies that the given transaction has supplied
// enough fees(gas + stability) to cover a proposer's minimum fees. A result object is returned
// indicating success or failure.
//
// Contract: This should only be called during CheckTx as it cannot be part of
// consensus.
func EnsureSufficientMempoolFees(ctx sdk.Context, gas uint64, feeCoins sdk.Coins, taxes sdk.Coins) error {
	requiredFees := sdk.Coins{}
	minGasPrices := ctx.MinGasPrices()
	if !minGasPrices.IsZero() {
		requiredFees = make(sdk.Coins, len(minGasPrices))

		// Determine the required fees by multiplying each required minimum gas
		// price by the gas limit, where fee = ceil(minGasPrice * gasLimit).
		glDec := sdk.NewDec(int64(gas))
		for i, gp := range minGasPrices {
			fee := gp.Amount.Mul(glDec)
			requiredFees[i] = sdk.NewCoin(gp.Denom, fee.Ceil().RoundInt())
		}
	}

	// Before checking gas prices, remove taxed from fee
	var hasNeg bool
	if feeCoins, hasNeg = feeCoins.SafeSub(taxes); hasNeg {
		return fmt.Errorf("insufficient fees; got: %q, required: %q = %q(gas) +%q(stability)", feeCoins.Add(taxes...), requiredFees.Add(taxes...), requiredFees, taxes)
	}

	if !requiredFees.IsZero() && !feeCoins.IsAnyGTE(requiredFees) {
		return fmt.Errorf("insufficient fees; got: %q, required: %q = %q(gas) +%q(stability)", feeCoins.Add(taxes...), requiredFees.Add(taxes...), requiredFees, taxes)
	}

	return nil
}

// FilterMsgAndComputeTax computes the stability tax on MsgSend and MsgMultiSend.
func FilterMsgAndComputeTax(ctx sdk.Context, tk TreasuryKeeper, msgs []sdk.Msg) sdk.Coins {
	taxes := sdk.Coins{}
	for _, msg := range msgs {
		switch msg := msg.(type) {
		case bank.MsgSend:
			taxes = taxes.Add(computeTax(ctx, tk, msg.Amount)...)

		case bank.MsgMultiSend:
			for _, input := range msg.Inputs {
				taxes = taxes.Add(computeTax(ctx, tk, input.Coins)...)
			}

		case marketexported.MsgSwapSend:
			taxes = taxes.Add(computeTax(ctx, tk, sdk.NewCoins(msg.OfferCoin))...)

		case wasmexported.MsgInstantiateContract:
			taxes = taxes.Add(computeTax(ctx, tk, msg.InitCoins)...)

		case wasmexported.MsgExecuteContract:
			taxes = taxes.Add(computeTax(ctx, tk, msg.Coins)...)

		case msgauthexported.MsgExecAuthorized:
			taxes = taxes.Add(FilterMsgAndComputeTax(ctx, tk, msg.Msgs)...)
		}
	}

	return taxes
}

// computes the stability tax according to tax-rate and tax-cap
func computeTax(ctx sdk.Context, tk TreasuryKeeper, principal sdk.Coins) sdk.Coins {
	taxRate := tk.GetTaxRate(ctx)
	if taxRate.Equal(sdk.ZeroDec()) {
		return sdk.Coins{}
	}

	taxes := sdk.Coins{}
	for _, coin := range principal {
		if coin.Denom == core.MicroLunaDenom || coin.Denom == sdk.DefaultBondDenom {
			continue
		}

		taxDue := sdk.NewDecFromInt(coin.Amount).Mul(taxRate).TruncateInt()

		// If tax due is greater than the tax cap, cap!
		taxCap := tk.GetTaxCap(ctx, coin.Denom)
		if taxDue.GT(taxCap) {
			taxDue = taxCap
		}

		if taxDue.Equal(sdk.ZeroInt()) {
			continue
		}

		taxes = taxes.Add(sdk.NewCoin(coin.Denom, taxDue))
	}

	return taxes
}

func isOracleTx(ctx sdk.Context, msgs []sdk.Msg) bool {
	for _, msg := range msgs {
		switch msg.(type) {
		case oracleexported.MsgAggregateExchangeRatePrevote:
			continue
		case oracleexported.MsgAggregateExchangeRateVote:
			continue
		default:
			return false
		}
	}

	return true
}
