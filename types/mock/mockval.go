package mock

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type MockValidator struct {
	sdk.Validator

	Address sdk.ValAddress
	Power   sdk.Int
}

func NewMockValidator(address sdk.ValAddress, power sdk.Int) MockValidator {
	return MockValidator{
		Address: address,
		Power:   power,
	}
}

func (mv MockValidator) GetBondedTokens() sdk.Int {
	return mv.Power
}

func (mv MockValidator) GetOperator() sdk.ValAddress {
	return mv.Address
}

type MockValset struct {
	sdk.ValidatorSet

	Validators []MockValidator
}

func NewMockValSet() MockValset {
	return MockValset{
		Validators: []MockValidator{},
	}
}

func (mv MockValset) Validator(ctx sdk.Context, valAddress sdk.ValAddress) sdk.Validator {
	for _, val := range mv.Validators {
		if val.Address.Equals(valAddress) {
			return val
		}
	}
	return nil
}

func (mv MockValset) TotalBondedTokens(ctx sdk.Context) sdk.Int {
	rval := sdk.ZeroInt()
	for _, val := range mv.Validators {
		rval = rval.Add(val.Power)
	}
	return rval
}

func (mv MockValset) IterateBondedValidatorsByPower(ctx sdk.Context,
	handler func(index int64, validator sdk.Validator) (stop bool)) {
	for i, val := range mv.Validators {
		if handler(int64(i), val) {
			break
		}
	}
}
