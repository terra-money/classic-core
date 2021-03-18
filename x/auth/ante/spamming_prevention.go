package ante

import (
	"sync"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/spf13/viper"

	oracleexported "github.com/terra-project/core/x/oracle/exported"
)

// FlagTxGasHardLimit defines the hard cap to prevent tx spamming attack
const FlagTxGasHardLimit = "tx-gas-hard-limit"

// SpammingPreventionDecorator will check if the transaction's gas is smaller than
// configured hard cap
type SpammingPreventionDecorator struct {
	oracleKeeper     OracleKeeper
	oraclePrevoteMap map[string]int64
	oracleVoteMap    map[string]int64
	mu               *sync.Mutex
}

// NewSpammingPreventionDecorator returns new spamming prevention decorator instance
func NewSpammingPreventionDecorator(oracleKeeper OracleKeeper) SpammingPreventionDecorator {
	return SpammingPreventionDecorator{
		oracleKeeper:     oracleKeeper,
		oraclePrevoteMap: make(map[string]int64),
		oracleVoteMap:    make(map[string]int64),
	}
}

// AnteHandle handles msg tax fee checking
func (spd SpammingPreventionDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
	if ctx.IsCheckTx() {
		feeTx, ok := tx.(FeeTx)
		if !ok {
			return ctx, sdkerrors.Wrap(sdkerrors.ErrTxDecode, "Tx must be a FeeTx")
		}

		gas := feeTx.GetGas()
		gasHardLimit := viper.GetUint64(FlagTxGasHardLimit)
		if gas > gasHardLimit {
			return ctx, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "Tx cannot spend more than %d gas", gasHardLimit)
		}

		err := spd.checkOracleSpamming(ctx, feeTx.GetMsgs())
		if err != nil {
			return ctx, err
		}
	}

	return next(ctx, tx, simulate)
}

func (spd SpammingPreventionDecorator) checkOracleSpamming(ctx sdk.Context, msgs []sdk.Msg) error {
	spd.mu.Lock()
	defer spd.mu.Unlock()

	curHeight := ctx.BlockHeight()
	for _, msg := range msgs {
		switch msg := msg.(type) {
		case oracleexported.MsgAggregateExchangeRatePrevote:
			err := spd.oracleKeeper.ValidateFeeder(ctx, msg.Feeder, msg.Validator, true)
			if err != nil {
				return err
			}

			valAddrStr := msg.Validator.String()
			if lastSubmittedHeight, ok := spd.oraclePrevoteMap[valAddrStr]; ok && lastSubmittedHeight == curHeight {
				return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "the validator has already been submitted prevote at the current height")
			}

			spd.oraclePrevoteMap[valAddrStr] = curHeight
			continue
		case oracleexported.MsgAggregateExchangeRateVote:
			err := spd.oracleKeeper.ValidateFeeder(ctx, msg.Feeder, msg.Validator, true)
			if err != nil {
				return err
			}

			valAddrStr := msg.Validator.String()
			if lastSubmittedHeight, ok := spd.oracleVoteMap[valAddrStr]; ok && lastSubmittedHeight == curHeight {
				return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "the validator has already been submitted prevote at the current height")
			}

			spd.oracleVoteMap[valAddrStr] = curHeight
			continue
		default:
			return nil
		}
	}

	return nil
}
