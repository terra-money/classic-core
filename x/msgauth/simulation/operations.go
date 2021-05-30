package simulation

// DONTCOVER

import (
	"math/rand"
	"strings"
	"time"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/simapp/helpers"
	simappparams "github.com/cosmos/cosmos-sdk/simapp/params"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"github.com/terra-money/core/x/msgauth/internal/keeper"
	"github.com/terra-money/core/x/msgauth/internal/types"
)

// Simulation operation weights constants
const (
	OpWeightMsgGrantAuthorization = "op_weight_msg_grant_authorization"
	OpWeightRevokeAuthorization   = "op_weight_msg_revoke_authorization"
	OpWeightExecAuthorized        = "op_weight_msg_execute_authorized"
)

// WeightedOperations returns all the operations from the module with their respective weights
func WeightedOperations(
	appParams simulation.AppParams, cdc *codec.Codec, ak authkeeper.AccountKeeper, bk bank.Keeper, k keeper.Keeper,
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
			SimulateMsgGrantAuthorization(ak, k),
		),
		simulation.NewWeightedOperation(
			weightRevokeAuthorization,
			SimulateMsgRevokeAuthorization(ak, k),
		),
		simulation.NewWeightedOperation(
			weightExecAuthorized,
			SimulateMsgExecuteAuthorized(ak, bk, k),
		),
	}
}

// SimulateMsgGrantAuthorization generates a MsgGrantAuthorization with random values.
// nolint: funlen
func SimulateMsgGrantAuthorization(ak authkeeper.AccountKeeper, k keeper.Keeper) simulation.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simulation.Account, chainID string,
	) (simulation.OperationMsg, []simulation.FutureOperation, error) {

		granter, _ := simulation.RandomAcc(r, accs)
		grantee, _ := simulation.RandomAcc(r, accs)
		if granter.Address.Equals(grantee.Address) {
			return simulation.NoOpMsg(types.ModuleName), nil, nil
		}

		account := ak.GetAccount(ctx, granter.Address)

		spendableCoins := account.SpendableCoins(ctx.BlockTime())
		fees, err := simulation.RandomFees(r, ctx, spendableCoins)
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}

		msg := types.NewMsgGrantAuthorization(granter.Address, grantee.Address,
			types.NewSendAuthorization(spendableCoins.Sub(fees)), time.Hour)

		tx := helpers.GenTx(
			[]sdk.Msg{msg},
			fees,
			helpers.DefaultGenTxGas,
			chainID,
			[]uint64{account.GetAccountNumber()},
			[]uint64{account.GetSequence()},
			granter.PrivKey,
		)

		_, _, err = app.Deliver(tx)
		return simulation.NewOperationMsg(msg, true, ""), nil, err
	}
}

// SimulateMsgRevokeAuthorization generates a MsgRevokeAuthorization with random values.
// nolint: funlen
func SimulateMsgRevokeAuthorization(ak authkeeper.AccountKeeper, k keeper.Keeper) simulation.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simulation.Account, chainID string,
	) (simulation.OperationMsg, []simulation.FutureOperation, error) {

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
			return simulation.NoOpMsg(types.ModuleName), nil, nil
		}

		granter, _ := simulation.FindAccount(accs, granterAddr)
		account := ak.GetAccount(ctx, granter.Address)

		spendableCoins := account.SpendableCoins(ctx.BlockTime())
		fees, err := simulation.RandomFees(r, ctx, spendableCoins)
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}

		msg := types.NewMsgRevokeAuthorization(granterAddr, granteeAddr, targetGrant.Authorization.MsgType())

		tx := helpers.GenTx(
			[]sdk.Msg{msg},
			fees,
			helpers.DefaultGenTxGas,
			chainID,
			[]uint64{account.GetAccountNumber()},
			[]uint64{account.GetSequence()},
			granter.PrivKey,
		)

		_, _, err = app.Deliver(tx)
		return simulation.NewOperationMsg(msg, true, ""), nil, err
	}
}

// SimulateMsgExecuteAuthorized generates a MsgExecuteAuthorized with random values.
// nolint: funlen
func SimulateMsgExecuteAuthorized(ak authkeeper.AccountKeeper, bk bank.Keeper, k keeper.Keeper) simulation.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simulation.Account, chainID string,
	) (simulation.OperationMsg, []simulation.FutureOperation, error) {
		if !bk.GetSendEnabled(ctx) {
			return simulation.NoOpMsg(types.ModuleName), nil, nil
		}

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
			return simulation.NoOpMsg(types.ModuleName), nil, nil
		}

		grantee, _ := simulation.FindAccount(accs, granteeAddr)
		granterAccount := ak.GetAccount(ctx, granterAddr)
		granteeAccount := ak.GetAccount(ctx, granteeAddr)

		granterSpendableCoins := granterAccount.SpendableCoins(ctx.BlockTime())
		if granterSpendableCoins.Empty() {
			return simulation.NoOpMsg(types.ModuleName), nil, nil
		}

		granteeSpendableCoins := granteeAccount.SpendableCoins(ctx.BlockTime())
		fees, err := simulation.RandomFees(r, ctx, granteeSpendableCoins)
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}

		execMsg := bank.NewMsgSend(
			granterAddr,
			granteeAddr,
			simulation.RandSubsetCoins(r, granterSpendableCoins),
		)

		msg := types.NewMsgExecAuthorized(grantee.Address, []sdk.Msg{execMsg})

		allow, _, _ := targetGrant.Authorization.Accept(execMsg, ctx.BlockHeader())
		if !allow {
			return simulation.NoOpMsg(types.ModuleName), nil, nil
		}

		tx := helpers.GenTx(
			[]sdk.Msg{msg},
			fees,
			helpers.DefaultGenTxGas,
			chainID,
			[]uint64{granteeAccount.GetAccountNumber()},
			[]uint64{granteeAccount.GetSequence()},
			grantee.PrivKey,
		)

		_, _, err = app.Deliver(tx)
		if err != nil {
			if strings.Contains(err.Error(), "insufficient fee") {
				return simulation.NoOpMsg(types.ModuleName), nil, nil
			}

			return simulation.NoOpMsg(types.ModuleName), nil, err
		}

		return simulation.NewOperationMsg(msg, true, ""), nil, nil
	}
}
