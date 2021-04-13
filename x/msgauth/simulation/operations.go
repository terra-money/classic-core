package simulation

//DONTCOVER

import (
	"math/rand"
	"strings"
	"time"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/simapp/helpers"
	simappparams "github.com/cosmos/cosmos-sdk/simapp/params"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"github.com/terra-project/core/x/msgauth/keeper"
	"github.com/terra-project/core/x/msgauth/types"
)

// Simulation operation weights constants
const (
	OpWeightMsgGrantAuthorization = "op_weight_msg_grant_authorization"
	OpWeightRevokeAuthorization   = "op_weight_msg_revoke_authorization"
	OpWeightExecAuthorized        = "op_weight_msg_execute_authorized"
)

// WeightedOperations returns all the operations from the module with their respective weights
func WeightedOperations(
	appParams simtypes.AppParams,
	cdc codec.JSONMarshaler,
	ak types.AccountKeeper,
	bk types.BankKeeper,
	k keeper.Keeper,
) simulation.WeightedOperations {

	var (
		weightMsgGrantAuthorization int
		weightRevokeAuthorization   int
		weightExecAuthorized        int
	)

	appParams.GetOrGenerate(cdc, OpWeightMsgGrantAuthorization, &weightMsgGrantAuthorization, nil,
		func(_ *rand.Rand) {
			weightMsgGrantAuthorization = simappparams.DefaultWeightMsgDelegate
		},
	)

	appParams.GetOrGenerate(cdc, OpWeightRevokeAuthorization, &weightRevokeAuthorization, nil,
		func(_ *rand.Rand) {
			weightRevokeAuthorization = simappparams.DefaultWeightMsgUndelegate
		},
	)

	appParams.GetOrGenerate(cdc, OpWeightExecAuthorized, &weightExecAuthorized, nil,
		func(_ *rand.Rand) {
			weightExecAuthorized = simappparams.DefaultWeightMsgSend
		},
	)

	return simulation.WeightedOperations{
		simulation.NewWeightedOperation(
			weightMsgGrantAuthorization,
			SimulateMsgGrantAuthorization(ak, bk, k),
		),
		simulation.NewWeightedOperation(
			weightRevokeAuthorization,
			SimulateMsgRevokeAuthorization(ak, bk, k),
		),
		simulation.NewWeightedOperation(
			weightExecAuthorized,
			SimulateMsgExecuteAuthorized(ak, bk, k),
		),
	}
}

// SimulateMsgGrantAuthorization generates a MsgGrantAuthorization with random values.
// nolint: funlen
func SimulateMsgGrantAuthorization(
	ak types.AccountKeeper,
	bk types.BankKeeper,
	k keeper.Keeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {

		granter, _ := simtypes.RandomAcc(r, accs)
		grantee, _ := simtypes.RandomAcc(r, accs)
		if granter.Address.Equals(grantee.Address) {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgGrantAuthorization, "unable to grant to self"), nil, nil
		}

		if _, hasGrant := k.GetGrant(ctx, granter.Address, grantee.Address, banktypes.TypeMsgSend); hasGrant {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgGrantAuthorization, "grant already exists"), nil, nil
		}

		account := ak.GetAccount(ctx, granter.Address)
		spendableCoins := bk.SpendableCoins(ctx, granter.Address)
		fees, err := simtypes.RandomFees(r, ctx, spendableCoins)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgGrantAuthorization, "unable to generate fees"), nil, err
		}

		msg, err := types.NewMsgGrantAuthorization(granter.Address, grantee.Address,
			types.NewSendAuthorization(spendableCoins.Sub(fees)), time.Hour)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgGrantAuthorization, "unable to generate grant msg"), nil, err
		}

		txGen := simappparams.MakeTestEncodingConfig().TxConfig
		tx, err := helpers.GenTx(
			txGen,
			[]sdk.Msg{msg},
			fees,
			helpers.DefaultGenTxGas,
			chainID,
			[]uint64{account.GetAccountNumber()},
			[]uint64{account.GetSequence()},
			granter.PrivKey,
		)

		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "unable to generate mock tx"), nil, err
		}

		_, _, err = app.Deliver(txGen.TxEncoder(), tx)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "unable to deliver tx"), nil, err
		}

		return simtypes.NewOperationMsg(msg, true, ""), nil, nil
	}
}

// SimulateMsgRevokeAuthorization generates a MsgRevokeAuthorization with random values.
// nolint: funlen
func SimulateMsgRevokeAuthorization(
	ak types.AccountKeeper,
	bk types.BankKeeper,
	k keeper.Keeper,
) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {

		hasGrant := false
		var targetGrant types.AuthorizationGrant
		var granterAddr sdk.AccAddress
		var granteeAddr sdk.AccAddress
		k.IterateGrants(ctx, func(granter, grantee sdk.AccAddress, grant types.AuthorizationGrant) bool {
			targetGrant = grant
			granterAddr = granter
			granteeAddr = grantee
			hasGrant = true
			return true
		})

		if !hasGrant {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgExecAuthorized, "no grant exist"), nil, nil
		}

		granter, _ := simtypes.FindAccount(accs, granterAddr)
		account := ak.GetAccount(ctx, granter.Address)

		spendableCoins := bk.SpendableCoins(ctx, granterAddr)
		fees, err := simtypes.RandomFees(r, ctx, spendableCoins)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgRevokeAuthorization, "failed to generate fees"), nil, err
		}

		msg := types.NewMsgRevokeAuthorization(granterAddr, granteeAddr, targetGrant.GetAuthorization().MsgType())

		txGen := simappparams.MakeTestEncodingConfig().TxConfig
		tx, err := helpers.GenTx(
			txGen,
			[]sdk.Msg{msg},
			fees,
			helpers.DefaultGenTxGas,
			chainID,
			[]uint64{account.GetAccountNumber()},
			[]uint64{account.GetSequence()},
			granter.PrivKey,
		)

		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "unable to generate mock tx"), nil, err
		}

		_, _, err = app.Deliver(txGen.TxEncoder(), tx)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "unable to deliver tx"), nil, err
		}

		return simtypes.NewOperationMsg(msg, true, ""), nil, nil
	}
}

// SimulateMsgExecuteAuthorized generates a MsgExecuteAuthorized with random values.
// nolint: funlen
func SimulateMsgExecuteAuthorized(
	ak types.AccountKeeper,
	bk types.BankKeeper,
	k keeper.Keeper,
) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		hasGrant := false
		var targetGrant types.AuthorizationGrant
		var granterAddr sdk.AccAddress
		var granteeAddr sdk.AccAddress
		k.IterateGrants(ctx, func(granter, grantee sdk.AccAddress, grant types.AuthorizationGrant) bool {
			targetGrant = grant
			granterAddr = granter
			granteeAddr = grantee
			hasGrant = true
			return true
		})

		if !hasGrant {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgExecAuthorized, "no grant exists"), nil, nil
		}

		grantee, _ := simtypes.FindAccount(accs, granteeAddr)
		granteeAccount := ak.GetAccount(ctx, granteeAddr)

		granterSpendableCoins := bk.SpendableCoins(ctx, granterAddr)
		if granterSpendableCoins.Empty() {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgExecAuthorized, "no spendable coins"), nil, nil
		}

		granteeSpendableCoins := bk.SpendableCoins(ctx, granteeAddr)
		fees, err := simtypes.RandomFees(r, ctx, granteeSpendableCoins)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgExecAuthorized, "failed to generate fees"), nil, err
		}

		coins := simtypes.RandSubsetCoins(r, granterSpendableCoins)
		if err := bk.SendEnabledCoins(ctx, coins...); err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgExecAuthorized, err.Error()), nil, nil
		}

		execMsg := banktypes.NewMsgSend(
			granterAddr,
			granteeAddr,
			coins,
		)

		msg, err := types.NewMsgExecAuthorized(grantee.Address, []sdk.Msg{execMsg})
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgExecAuthorized, "failed to create msg"), nil, err
		}

		allow, _, _ := targetGrant.GetAuthorization().Accept(execMsg, ctx.BlockHeader())
		if !allow {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgExecAuthorized, "not allowed"), nil, nil
		}

		txGen := simappparams.MakeTestEncodingConfig().TxConfig
		tx, err := helpers.GenTx(
			txGen,
			[]sdk.Msg{msg},
			fees,
			helpers.DefaultGenTxGas,
			chainID,
			[]uint64{granteeAccount.GetAccountNumber()},
			[]uint64{granteeAccount.GetSequence()},
			grantee.PrivKey,
		)

		if err != nil {
			if strings.Contains(err.Error(), "insufficient fee") {
				return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgExecAuthorized, "skip low fee due to tax"), nil, nil
			}
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "unable to generate mock tx"), nil, err
		}

		_, _, err = app.Deliver(txGen.TxEncoder(), tx)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "unable to deliver tx"), nil, err
		}

		return simtypes.NewOperationMsg(msg, true, ""), nil, nil
	}
}
