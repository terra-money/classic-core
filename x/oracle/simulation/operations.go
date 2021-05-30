package simulation

// DONTCOVER

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/simapp/helpers"
	simappparams "github.com/cosmos/cosmos-sdk/simapp/params"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	core "github.com/terra-money/core/types"
	"github.com/terra-money/core/x/oracle/internal/keeper"
	"github.com/terra-money/core/x/oracle/internal/types"
)

// Simulation operation weights constants
const (
	OpWeightMsgExchangeRatePrevote = "op_weight_msg_exchange_rate_prevote"
	OpWeightMsgExchangeRateVote    = "op_weight_msg_exchange_rate_vote"
	OpWeightMsgDelegateFeedConsent = "op_weight_msg_exchange_feed_consent"

	salt = "1234"
)

var (
	whitelist                      = []string{core.MicroKRWDenom, core.MicroUSDDenom, core.MicroSDRDenom, core.MicroMNTDenom}
	voteHashMap map[string]sdk.Dec = make(map[string]sdk.Dec)
)

// WeightedOperations returns all the operations from the module with their respective weights
func WeightedOperations(
	appParams simulation.AppParams,
	cdc *codec.Codec,
	ak authkeeper.AccountKeeper,
	k keeper.Keeper) simulation.WeightedOperations {
	var (
		weightMsgExchangeRatePrevote int
		weightMsgExchangeRateVote    int
		weightMsgDelegateFeedConsent int
	)
	appParams.GetOrGenerate(cdc, OpWeightMsgExchangeRatePrevote, &weightMsgExchangeRatePrevote, nil,
		func(_ *rand.Rand) {
			weightMsgExchangeRatePrevote = simappparams.DefaultWeightMsgSend * 2
		},
	)

	appParams.GetOrGenerate(cdc, OpWeightMsgExchangeRateVote, &weightMsgExchangeRateVote, nil,
		func(_ *rand.Rand) {
			weightMsgExchangeRateVote = simappparams.DefaultWeightMsgSend * 2
		},
	)

	appParams.GetOrGenerate(cdc, OpWeightMsgDelegateFeedConsent, &weightMsgDelegateFeedConsent, nil,
		func(_ *rand.Rand) {
			weightMsgDelegateFeedConsent = simappparams.DefaultWeightMsgSetWithdrawAddress
		},
	)

	return simulation.WeightedOperations{
		simulation.NewWeightedOperation(
			weightMsgExchangeRatePrevote,
			SimulateMsgExchangeRatePrevote(ak, k),
		),
		simulation.NewWeightedOperation(
			weightMsgExchangeRateVote,
			SimulateMsgExchangeRateVote(ak, k),
		),
		simulation.NewWeightedOperation(
			weightMsgDelegateFeedConsent,
			SimulateMsgDelegateFeedConsent(ak, k),
		),
	}
}

// SimulateMsgExchangeRatePrevote generates a MsgExchangeRatePrevote with random values.
// nolint: funlen
func SimulateMsgExchangeRatePrevote(ak authkeeper.AccountKeeper, k keeper.Keeper) simulation.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simulation.Account, chainID string,
	) (simulation.OperationMsg, []simulation.FutureOperation, error) {

		simAccount, _ := simulation.RandomAcc(r, accs)
		address := sdk.ValAddress(simAccount.Address)

		// ensure the validator exists
		val := k.StakingKeeper.Validator(ctx, address)
		power := k.StakingKeeper.GetLastValidatorPower(ctx, address)
		if val == nil || !val.IsBonded() || power == 0 {
			return simulation.NoOpMsg(types.ModuleName), nil, nil
		}

		denom := whitelist[simulation.RandIntBetween(r, 0, len(whitelist))]
		price := sdk.NewDecWithPrec(int64(simulation.RandIntBetween(r, 1, 10000)), int64(1))
		voteHash := types.GetVoteHash(salt, price, denom, address)

		feederAddr := k.GetOracleDelegate(ctx, address)
		feederSimAccount, _ := simulation.FindAccount(accs, feederAddr)
		feederAccount := ak.GetAccount(ctx, feederAddr)

		fees, err := simulation.RandomFees(r, ctx, feederAccount.SpendableCoins(ctx.BlockTime()))
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}

		msg := types.NewMsgExchangeRatePrevote(voteHash, denom, feederAddr, address)

		tx := helpers.GenTx(
			[]sdk.Msg{msg},
			fees,
			helpers.DefaultGenTxGas,
			chainID,
			[]uint64{feederAccount.GetAccountNumber()},
			[]uint64{feederAccount.GetSequence()},
			feederSimAccount.PrivKey,
		)

		_, _, err = app.Deliver(tx)
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}

		voteHashMap[denom+address.String()] = price

		return simulation.NewOperationMsg(msg, true, ""), nil, nil
	}
}

// SimulateMsgExchangeRateVote generates a MsgExchangeRateVote with random values.
// nolint: funlen
func SimulateMsgExchangeRateVote(ak authkeeper.AccountKeeper, k keeper.Keeper) simulation.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simulation.Account, chainID string,
	) (simulation.OperationMsg, []simulation.FutureOperation, error) {

		simAccount, _ := simulation.RandomAcc(r, accs)
		address := sdk.ValAddress(simAccount.Address)

		// ensure the validator exists
		val := k.StakingKeeper.Validator(ctx, address)
		power := k.StakingKeeper.GetLastValidatorPower(ctx, address)
		if val == nil || !val.IsBonded() || power == 0 {
			return simulation.NoOpMsg(types.ModuleName), nil, nil
		}

		// ensure vote hash exists
		denom := whitelist[simulation.RandIntBetween(r, 0, len(whitelist))]
		price, ok := voteHashMap[denom+address.String()]
		if !ok {
			return simulation.NoOpMsg(types.ModuleName), nil, nil
		}

		// get prevote
		prevote, err := k.GetExchangeRatePrevote(ctx, denom, address)
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, nil
		}

		params := k.GetParams(ctx)
		if (ctx.BlockHeight()/params.VotePeriod)-(prevote.SubmitBlock/params.VotePeriod) != 1 {
			return simulation.NoOpMsg(types.ModuleName), nil, nil
		}

		feederAddr := k.GetOracleDelegate(ctx, address)
		feederSimAccount, _ := simulation.FindAccount(accs, feederAddr)
		feederAccount := ak.GetAccount(ctx, feederAddr)

		fees, err := simulation.RandomFees(r, ctx, feederAccount.SpendableCoins(ctx.BlockTime()))
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}

		msg := types.NewMsgExchangeRateVote(price, salt, denom, feederAddr, address)

		tx := helpers.GenTx(
			[]sdk.Msg{msg},
			fees,
			helpers.DefaultGenTxGas,
			chainID,
			[]uint64{feederAccount.GetAccountNumber()},
			[]uint64{feederAccount.GetSequence()},
			feederSimAccount.PrivKey,
		)

		_, _, err = app.Deliver(tx)
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}

		return simulation.NewOperationMsg(msg, true, ""), nil, nil
	}
}

// SimulateMsgDelegateFeedConsent generates a MsgDelegateFeedConsent with random values.
// nolint: funlen
func SimulateMsgDelegateFeedConsent(ak authkeeper.AccountKeeper, k keeper.Keeper) simulation.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simulation.Account, chainID string,
	) (simulation.OperationMsg, []simulation.FutureOperation, error) {

		simAccount, _ := simulation.RandomAcc(r, accs)
		delegateAccount, _ := simulation.RandomAcc(r, accs)
		valAddress := sdk.ValAddress(simAccount.Address)
		delegateValAddress := sdk.ValAddress(delegateAccount.Address)
		account := ak.GetAccount(ctx, simAccount.Address)

		// ensure the validator exists
		val := k.StakingKeeper.Validator(ctx, valAddress)
		power := k.StakingKeeper.GetLastValidatorPower(ctx, valAddress)
		if val == nil || !val.IsBonded() || power == 0 {
			return simulation.NoOpMsg(types.ModuleName), nil, nil
		}

		// ensure the target address is not a validator
		val2 := k.StakingKeeper.Validator(ctx, delegateValAddress)
		if val2 != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, nil
		}

		fees, err := simulation.RandomFees(r, ctx, account.SpendableCoins(ctx.BlockTime()))
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}

		msg := types.NewMsgDelegateFeedConsent(valAddress, delegateAccount.Address)

		tx := helpers.GenTx(
			[]sdk.Msg{msg},
			fees,
			helpers.DefaultGenTxGas,
			chainID,
			[]uint64{account.GetAccountNumber()},
			[]uint64{account.GetSequence()},
			simAccount.PrivKey,
		)

		_, _, err = app.Deliver(tx)
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}

		return simulation.NewOperationMsg(msg, true, ""), nil, nil
	}
}
