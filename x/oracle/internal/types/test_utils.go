// nolint:deadcode unused noalias
package types

import (
	"math"
	"math/rand"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking/exported"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/secp256k1"
)

const oracleDecPrecision = 6

// GenerateRandomTestCase nolint
func GenerateRandomTestCase() (rates []float64, valValAddrs []sdk.ValAddress, stakingKeeper DummyStakingKeeper) {
	valValAddrs = []sdk.ValAddress{}
	mockValidators := []MockValidator{}

	base := math.Pow10(oracleDecPrecision)

	rand.Seed(int64(time.Now().Nanosecond()))
	numInputs := 10 + (rand.Int() % 100)
	for i := 0; i < numInputs; i++ {
		rate := float64(int64(rand.Float64()*base)) / base
		rates = append(rates, rate)

		pubKey := secp256k1.GenPrivKey().PubKey()
		valValAddr := sdk.ValAddress(pubKey.Address())
		valValAddrs = append(valValAddrs, valValAddr)

		power := rand.Int63()%1000 + 1
		mockValidator := NewMockValidator(valValAddr, power)
		mockValidators = append(mockValidators, mockValidator)
	}

	stakingKeeper = NewDummyStakingKeeper(mockValidators)

	return
}

var _ StakingKeeper = DummyStakingKeeper{}

// DummyStakingKeeper dummy staking keeper to test ballot
type DummyStakingKeeper struct {
	validators []MockValidator
}

// NewDummyStakingKeeper returns new DummyStakingKeeper instance
func NewDummyStakingKeeper(validators []MockValidator) DummyStakingKeeper {
	return DummyStakingKeeper{
		validators: validators,
	}
}

// Validators nolint
func (sk DummyStakingKeeper) Validators() []MockValidator {
	return sk.validators
}

// Validator nolint
func (sk DummyStakingKeeper) Validator(ctx sdk.Context, address sdk.ValAddress) exported.ValidatorI {
	for _, validator := range sk.validators {
		if validator.GetOperator().Equals(address) {
			return validator
		}
	}

	return nil
}

// TotalBondedTokens nolint
func (DummyStakingKeeper) TotalBondedTokens(_ sdk.Context) sdk.Int {
	return sdk.ZeroInt()
}

// Slash nolint
func (DummyStakingKeeper) Slash(sdk.Context, sdk.ConsAddress, int64, int64, sdk.Dec) {}

// IterateValidators nolint
func (DummyStakingKeeper) IterateValidators(sdk.Context, func(index int64, validator exported.ValidatorI) (stop bool)) {
}

// Jail nolint
func (DummyStakingKeeper) Jail(sdk.Context, sdk.ConsAddress) {
}

type MockValidator struct {
	power    int64
	operator sdk.ValAddress
}

var _ exported.ValidatorI = MockValidator{}

func (MockValidator) IsJailed() bool                                  { return false }
func (MockValidator) GetMoniker() string                              { return "" }
func (MockValidator) GetStatus() sdk.BondStatus                       { return sdk.Bonded }
func (MockValidator) IsBonded() bool                                  { return true }
func (MockValidator) IsUnbonded() bool                                { return false }
func (MockValidator) IsUnbonding() bool                               { return false }
func (v MockValidator) GetOperator() sdk.ValAddress                   { return v.operator }
func (MockValidator) GetConsPubKey() crypto.PubKey                    { return nil }
func (MockValidator) GetConsAddr() sdk.ConsAddress                    { return nil }
func (v MockValidator) GetTokens() sdk.Int                            { return sdk.TokensFromConsensusPower(v.power) }
func (v MockValidator) GetBondedTokens() sdk.Int                      { return sdk.TokensFromConsensusPower(v.power) }
func (v MockValidator) GetConsensusPower() int64                      { return v.power }
func (v MockValidator) GetCommission() sdk.Dec                        { return sdk.ZeroDec() }
func (v MockValidator) GetMinSelfDelegation() sdk.Int                 { return sdk.OneInt() }
func (v MockValidator) GetDelegatorShares() sdk.Dec                   { return sdk.NewDec(v.power) }
func (v MockValidator) TokensFromShares(sdk.Dec) sdk.Dec              { return sdk.ZeroDec() }
func (v MockValidator) TokensFromSharesTruncated(sdk.Dec) sdk.Dec     { return sdk.ZeroDec() }
func (v MockValidator) TokensFromSharesRoundUp(sdk.Dec) sdk.Dec       { return sdk.ZeroDec() }
func (v MockValidator) SharesFromTokens(amt sdk.Int) (sdk.Dec, error) { return sdk.ZeroDec(), nil }
func (v MockValidator) SharesFromTokensTruncated(amt sdk.Int) (sdk.Dec, error) {
	return sdk.ZeroDec(), nil
}
func (v MockValidator) SetPower(power int64) { v.power = power }

func NewMockValidator(valAddr sdk.ValAddress, power int64) MockValidator {
	return MockValidator{
		power:    power,
		operator: valAddr,
	}
}
