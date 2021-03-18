package ante

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/spf13/viper"

	core "github.com/terra-project/core/types"
)

// FlagTxGasHardLimit defines the hard cap to prevent tx spamming attack
const FlagTxGasHardLimit = "tx-gas-hard-limit"
const transactionGasHardCap = 30000000

// SpammingPreventionDecorator will check if the transaction's gas is smaller than
// configured hard cap
type SpammingPreventionDecorator struct {
}

// NewSpammingPreventionDecorator returns new spamming prevention decorator instance
func NewSpammingPreventionDecorator() SpammingPreventionDecorator {
	return SpammingPreventionDecorator{}
}

// AnteHandle handles msg tax fee checking
func (spd SpammingPreventionDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
	feeTx, ok := tx.(FeeTx)
	if !ok {
		return ctx, sdkerrors.Wrap(sdkerrors.ErrTxDecode, "Tx must be a FeeTx")
	}

	gas := feeTx.GetGas()
	if ctx.IsCheckTx() {
		gasHardLimit := viper.GetUint64(FlagTxGasHardLimit)
		if gas > gasHardLimit {
			return ctx, sdkerrors.Wrapf(sdkerrors.ErrOutOfGas, "Tx cannot spend more than %d gas", gasHardLimit)
		}
	}

	if !core.IsWaitingForSoftfork(ctx, 2) {
		if gas > transactionGasHardCap {
			return ctx, sdkerrors.Wrap(sdkerrors.ErrOutOfGas, "Tx exceed max gas usage")
		}
	}

	return next(ctx, tx, simulate)
}
