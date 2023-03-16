package ante

import (
	"fmt"
	core "github.com/classic-terra/core/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govkeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

// MinInitialDeposit Decorator will check Initial Deposits for MsgSubmitProposal
type MinInitialDepositDecorator struct {
	govKeeper      govkeeper.Keeper
	treasuryKeeper TreasuryKeeper
}

// NewMinInitialDeposit returns new min initial deposit decorator instance
func NewMinInitialDepositDecorator(govKeeper govkeeper.Keeper, treasuryKeeper TreasuryKeeper) MinInitialDepositDecorator {
	return MinInitialDepositDecorator{
		govKeeper:      govKeeper,
		treasuryKeeper: treasuryKeeper,
	}
}

// IsMsgSubmitProposal checks whether the input msg is a MsgSubmitProposal
func IsMsgSubmitProposal(msg sdk.Msg) bool {
	_, ok := msg.(*govtypes.MsgSubmitProposal)
	return ok
}

// HandleCheckMinInitialDeposit
func HandleCheckMinInitialDeposit(ctx sdk.Context, msg sdk.Msg, govKeeper govkeeper.Keeper, treasuryKeeper TreasuryKeeper) (err error) {
	submitPropMsg, ok := msg.(*govtypes.MsgSubmitProposal)
	if !ok {
		return fmt.Errorf("Could not dereference msg as MsgSubmitProposal")
	}

	minDeposit := govKeeper.GetDepositParams(ctx).MinDeposit
	requiredAmount := sdk.NewDecFromInt(minDeposit.AmountOf(core.MicroLunaDenom)).Mul(treasuryKeeper.GetMinInitialDepositRatio(ctx)).TruncateInt()

	requiredDepositCoins := sdk.NewCoins(
		sdk.NewCoin(core.MicroLunaDenom, requiredAmount),
	)
	initialDepositCoins := submitPropMsg.GetInitialDeposit()

	if !initialDepositCoins.IsAllGTE(requiredDepositCoins) {
		return fmt.Errorf("Not enough initial deposit provided. Expected %q; got %q", requiredDepositCoins, initialDepositCoins)
	}

	return nil
}

// AnteHandle handles checking MsgSubmitProposal
func (midd MinInitialDepositDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
	if simulate {
		return next(ctx, tx, simulate)
	}

	msgs := tx.GetMsgs()
	for _, msg := range msgs {

		if !IsMsgSubmitProposal(msg) {
			continue
		}

		err := HandleCheckMinInitialDeposit(ctx, msg, midd.govKeeper, midd.treasuryKeeper)
		if err != nil {
			return ctx, sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, err.Error())
		}

	}

	return next(ctx, tx, simulate)
}
