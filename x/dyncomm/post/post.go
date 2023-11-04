package post

import (
	dyncommkeeper "github.com/classic-terra/core/v2/x/dyncomm/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// DyncommDecorator does post runMsg store
// modifications for dyncomm module
type DyncommDecorator struct {
	dyncommKeeper dyncommkeeper.Keeper
}

func NewDyncommPostDecorator(dk dyncommkeeper.Keeper) DyncommDecorator {
	return DyncommDecorator{
		dyncommKeeper: dk,
	}
}

func (dd DyncommDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (sdk.Context, error) {
	if simulate {
		return next(ctx, tx, simulate)
	}

	if ctx.IsCheckTx() {
		return next(ctx, tx, simulate)
	}

	msgs := tx.GetMsgs()
	dd.FilterMsgsAndProcessMsgs(ctx, msgs...)

	return next(ctx, tx, simulate)
}

func (dd DyncommDecorator) FilterMsgsAndProcessMsgs(ctx sdk.Context, msgs ...sdk.Msg) {
	for _, msg := range msgs {
		switch msg.(type) {
		case *stakingtypes.MsgEditValidator:
			dd.ProcessEditValidator(ctx, msg)
		case *stakingtypes.MsgCreateValidator:
			dd.ProcessCreateValidator(ctx, msg)
		default:
			continue
		}
	}
}

func (dd DyncommDecorator) ProcessEditValidator(ctx sdk.Context, msg sdk.Msg) {
	msgEditValidator := msg.(*stakingtypes.MsgEditValidator)

	// no update of CommissionRate provided
	if msgEditValidator.CommissionRate == nil {
		return
	}

	// post handler runs after successfully
	// calling runMsgs -> we can set state changes here!
	newIntendedRate := msgEditValidator.CommissionRate
	dd.dyncommKeeper.SetTargetCommissionRate(ctx, msgEditValidator.ValidatorAddress, *newIntendedRate)
}

func (dd DyncommDecorator) ProcessCreateValidator(ctx sdk.Context, msg sdk.Msg) {
	// post handler runs after successfully
	// calling runMsgs -> we can set state changes here!
	msgCreateValidator := msg.(*stakingtypes.MsgCreateValidator)
	newIntendedRate := msgCreateValidator.Commission.Rate
	dd.dyncommKeeper.SetTargetCommissionRate(ctx, msgCreateValidator.ValidatorAddress, newIntendedRate)
}
