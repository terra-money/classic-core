package keeper

import (
	types "github.com/classic-terra/core/v2/x/dyncomm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// GetVotingPower calculates the voting power of a validator in percent
func (k Keeper) CalculateVotingPower(ctx sdk.Context, validator stakingtypes.Validator) (ret sdk.Dec) {
	totalPower := k.StakingKeeper.GetLastTotalPower(ctx).Int64()
	validatorPower := sdk.TokensToConsensusPower(
		validator.Tokens,
		k.StakingKeeper.PowerReduction(ctx),
	)
	return sdk.NewDec(validatorPower).QuoInt64(totalPower).MulInt64(100)
}

// CalculateDynCommission calculates the min commission according
// to StrathColes formula
func (k Keeper) CalculateDynCommission(ctx sdk.Context, validator stakingtypes.Validator) (ret sdk.Dec) {
	// The original parameters as defined
	// by Strath
	A := k.GetMaxZero(ctx)
	B := k.GetSlopeBase(ctx)
	C := k.GetSlopeVpImpact(ctx)
	D := k.GetCap(ctx).MulInt64(100)
	x := k.CalculateVotingPower(ctx, validator)
	factorA := x.Sub(A)
	quotient := x.Quo(C)
	factorB := quotient.Add(B)
	minComm := k.StakingKeeper.MinCommissionRate(ctx).MulInt64(100)

	y := factorA.Mul(factorB)
	if y.GT(D) {
		y = D
	}
	if minComm.GT(y) {
		y = minComm
	}
	return y.QuoInt64(100)
}

func (k Keeper) SetDynCommissionRate(ctx sdk.Context, validator string, rate sdk.Dec) {
	var preSetRate types.ValidatorCommissionRate
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetMinCommissionRatesKey(validator))
	targetRate := sdk.ZeroDec()

	if bz != nil {
		k.cdc.MustUnmarshal(bz, &preSetRate)
		targetRate = *preSetRate.TargetCommissionRate
	}
	bz = k.cdc.MustMarshal(
		&types.ValidatorCommissionRate{
			ValidatorAddress:     validator,
			MinCommissionRate:    &rate,
			TargetCommissionRate: &targetRate,
		},
	)
	store.Set(types.GetMinCommissionRatesKey(validator), bz)
}

func (k Keeper) SetTargetCommissionRate(ctx sdk.Context, validator string, rate sdk.Dec) {
	var preSetRate types.ValidatorCommissionRate
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetMinCommissionRatesKey(validator))
	minRate := sdk.ZeroDec()

	if bz != nil {
		k.cdc.MustUnmarshal(bz, &preSetRate)
		minRate = *preSetRate.MinCommissionRate
	}
	bz = k.cdc.MustMarshal(
		&types.ValidatorCommissionRate{
			ValidatorAddress:     validator,
			MinCommissionRate:    &minRate,
			TargetCommissionRate: &rate,
		},
	)
	store.Set(types.GetMinCommissionRatesKey(validator), bz)
}

func (k Keeper) GetDynCommissionRate(ctx sdk.Context, validator string) (rate sdk.Dec) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetMinCommissionRatesKey(validator))
	if bz == nil {
		return sdk.ZeroDec()
	}

	var validatorRate types.ValidatorCommissionRate
	k.cdc.MustUnmarshal(bz, &validatorRate)
	return *validatorRate.MinCommissionRate
}

func (k Keeper) GetTargetCommissionRate(ctx sdk.Context, validator string) (rate sdk.Dec) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetMinCommissionRatesKey(validator))
	if bz == nil {
		return sdk.ZeroDec()
	}

	var validatorRate types.ValidatorCommissionRate
	k.cdc.MustUnmarshal(bz, &validatorRate)
	return *validatorRate.TargetCommissionRate
}

// IterateDynCommissionRates iterates over dyn commission rates in the store
func (k Keeper) IterateDynCommissionRates(ctx sdk.Context, cb func(types.ValidatorCommissionRate) bool) {
	store := ctx.KVStore(k.storeKey)
	it := store.Iterator(nil, nil)
	defer it.Close()

	for ; it.Valid(); it.Next() {
		var entry types.ValidatorCommissionRate
		if err := entry.Unmarshal(it.Value()); err != nil {
			panic(err)
		}

		if cb(entry) {
			break
		}
	}
}

func (k Keeper) UpdateValidatorMinRates(ctx sdk.Context, validator stakingtypes.Validator) {
	var newRate sdk.Dec
	minRate := k.CalculateDynCommission(ctx, validator)
	newMaxRate := validator.Commission.MaxRate
	targetRate := k.GetTargetCommissionRate(ctx, validator.OperatorAddress)

	// assume the newRate will be the target rate ...
	newRate = targetRate

	// ... but enforce min rate
	if newRate.LT(minRate) {
		newRate = minRate
	}

	// new min rate pushes max rate
	if newMaxRate.LT(minRate) {
		newMaxRate = minRate
	}

	newValidator := validator
	newValidator.Commission = stakingtypes.NewCommissionWithTime(
		newRate,
		newMaxRate,
		validator.Commission.MaxChangeRate,
		validator.Commission.UpdateTime,
	)

	k.StakingKeeper.SetValidator(ctx, newValidator)
	k.SetDynCommissionRate(ctx, validator.OperatorAddress, minRate)

	// Debug
	targetRate = k.GetTargetCommissionRate(ctx, validator.OperatorAddress)
	ctx.Logger().Debug("dyncomm:", "val", validator.OperatorAddress, "min_rate", minRate, "new target_rate", targetRate)
}

func (k Keeper) UpdateAllBondedValidatorRates(ctx sdk.Context) (err error) {
	k.StakingKeeper.IterateValidators(ctx, func(index int64, validator stakingtypes.ValidatorI) (stop bool) {
		val := validator.(stakingtypes.Validator)

		if !val.IsBonded() {
			return false
		}

		k.UpdateValidatorMinRates(ctx, val)

		return false
	})

	return nil
}
