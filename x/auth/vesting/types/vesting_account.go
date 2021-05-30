package types

import (
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authexported "github.com/cosmos/cosmos-sdk/x/auth/exported"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	vestexported "github.com/cosmos/cosmos-sdk/x/auth/vesting/exported"
	vesttypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	"github.com/tendermint/tendermint/crypto"

	customauthtypes "github.com/terra-money/core/x/auth/internal/types"

	"gopkg.in/yaml.v2"
)

// for pretty purpose
type vestingAccountYAML struct {
	Address          sdk.AccAddress `json:"address" yaml:"address"`
	Coins            sdk.Coins      `json:"coins" yaml:"coins"`
	PubKey           string         `json:"public_key" yaml:"public_key"`
	AccountNumber    uint64         `json:"account_number" yaml:"account_number"`
	Sequence         uint64         `json:"sequence" yaml:"sequence"`
	OriginalVesting  sdk.Coins      `json:"original_vesting" yaml:"original_vesting"`
	DelegatedFree    sdk.Coins      `json:"delegated_free" yaml:"delegated_free"`
	DelegatedVesting sdk.Coins      `json:"delegated_vesting" yaml:"delegated_vesting"`
	EndTime          int64          `json:"end_time" yaml:"end_time"`

	// custom fields based on concrete vesting type which can be omitted
	VestingSchedules VestingSchedules `json:"vesting_schedules,omitempty" yaml:"vesting_schedules,omitempty"`
}

// To prevent stack overflow
type vestingAccountJSON struct {
	Address          sdk.AccAddress `json:"address" yaml:"address"`
	Coins            sdk.Coins      `json:"coins" yaml:"coins"`
	PubKey           crypto.PubKey  `json:"public_key" yaml:"public_key"`
	AccountNumber    uint64         `json:"account_number" yaml:"account_number"`
	Sequence         uint64         `json:"sequence" yaml:"sequence"`
	OriginalVesting  sdk.Coins      `json:"original_vesting" yaml:"original_vesting"`
	DelegatedFree    sdk.Coins      `json:"delegated_free" yaml:"delegated_free"`
	DelegatedVesting sdk.Coins      `json:"delegated_vesting" yaml:"delegated_vesting"`
	EndTime          int64          `json:"end_time" yaml:"end_time"`

	// custom fields based on concrete vesting type which can be omitted
	VestingSchedules VestingSchedules `json:"vesting_schedules,omitempty" yaml:"vesting_schedules,omitempty"`
}

//-----------------------------------------------------------------------------
// Lazy Graded Vesting Account

// LazyGradedVestingAccount implements the LazyGradedVestingAccount interface. It vests all
// coins according to a predefined schedule.
var _ vestexported.VestingAccount = (*LazyGradedVestingAccount)(nil)
var _ authexported.GenesisAccount = (*LazyGradedVestingAccount)(nil)

// Register the vesting account types on the auth module codec
func init() {
	customauthtypes.RegisterAccountTypeCodec(&vesttypes.BaseVestingAccount{}, "core/BaseVestingAccount")
	customauthtypes.RegisterAccountTypeCodec(&LazyGradedVestingAccount{}, "core/LazyGradedVestingAccount")
}

// LazyGradedVestingAccount implements the VestingAccount interface. It vests tokens according to
// a predefined set of vesting schedule.
type LazyGradedVestingAccount struct {
	*vesttypes.BaseVestingAccount

	VestingSchedules VestingSchedules `json:"vesting_schedules"`
}

// NewLazyGradedVestingAccountRaw creates a new LazyGradedVestingAccount object from BaseVestingAccount
func NewLazyGradedVestingAccountRaw(baseVestingAcc *vesttypes.BaseVestingAccount, lazyVestingSchedules VestingSchedules) *LazyGradedVestingAccount {
	return &LazyGradedVestingAccount{
		BaseVestingAccount: baseVestingAcc,
		VestingSchedules:   lazyVestingSchedules,
	}
}

// NewLazyGradedVestingAccount returns a new LazyGradedVestingAccount
func NewLazyGradedVestingAccount(baseAcc *authtypes.BaseAccount, lazyVestingSchedules VestingSchedules) *LazyGradedVestingAccount {
	baseVestingAcc := &vesttypes.BaseVestingAccount{
		BaseAccount:     baseAcc,
		OriginalVesting: baseAcc.Coins,
		EndTime:         0,
	}

	return &LazyGradedVestingAccount{baseVestingAcc, lazyVestingSchedules}
}

// GetVestingSchedules returns the VestingSchedules of the graded lazy vesting account
func (lgva LazyGradedVestingAccount) GetVestingSchedules() VestingSchedules {
	return lgva.VestingSchedules
}

// GetVestingSchedule returns the VestingSchedule of the given denom
func (lgva LazyGradedVestingAccount) GetVestingSchedule(denom string) (VestingSchedule, bool) {
	for _, vs := range lgva.VestingSchedules {
		if vs.Denom == denom {
			return vs, true
		}
	}

	return VestingSchedule{}, false
}

// GetVestedCoins returns the total amount of vested coins for a graded vesting
// account. All coins are vested continuously once the schedule's StartTime has elapsed until EndTime.
func (lgva LazyGradedVestingAccount) GetVestedCoins(blockTime time.Time) sdk.Coins {
	var vestedCoins sdk.Coins
	for _, ovc := range lgva.OriginalVesting {
		if vestingSchedule, exists := lgva.GetVestingSchedule(ovc.Denom); exists {
			vestedRatio := vestingSchedule.GetVestedRatio(blockTime.Unix())
			vestedAmt := ovc.Amount.ToDec().Mul(vestedRatio).RoundInt()
			if vestedAmt.Equal(sdk.ZeroInt()) {
				continue
			}
			vestedCoins = append(vestedCoins, sdk.NewCoin(ovc.Denom, vestedAmt))
		} else {
			vestedCoins = append(vestedCoins, sdk.NewCoin(ovc.Denom, ovc.Amount))
		}
	}

	return vestedCoins
}

// GetVestingCoins returns the total number of vesting coins for a graded
// vesting account.
func (lgva LazyGradedVestingAccount) GetVestingCoins(blockTime time.Time) sdk.Coins {
	return lgva.OriginalVesting.Sub(lgva.GetVestedCoins(blockTime))
}

// SpendableCoins returns the total number of spendable coins for a graded
// vesting account.
func (lgva LazyGradedVestingAccount) SpendableCoins(blockTime time.Time) sdk.Coins {
	return lgva.spendableCoins(lgva.GetVestingCoins(blockTime))
}

// TrackDelegation tracks a desired delegation amount by setting the appropriate
// values for the amount of delegated vesting, delegated free, and reducing the
// overall amount of base coins.
func (lgva *LazyGradedVestingAccount) TrackDelegation(blockTime time.Time, amount sdk.Coins) {
	lgva.trackDelegation(lgva.GetVestingCoins(blockTime), amount)
}

// GetStartTime returns zero since a lazy graded vesting account has no start time.
func (lgva LazyGradedVestingAccount) GetStartTime() int64 {
	return 0
}

// GetEndTime returns zero since a lazy graded vesting account has no end time.
func (lgva LazyGradedVestingAccount) GetEndTime() int64 {
	return 0
}

// Validate checks for errors on the account fields
func (lgva LazyGradedVestingAccount) Validate() error {
	for _, vestingSchedule := range lgva.GetVestingSchedules() {
		if err := vestingSchedule.Validate(); err != nil {
			return err
		}
	}

	return lgva.BaseVestingAccount.Validate()
}

// String implements fmt.Stringer interface
func (lgva LazyGradedVestingAccount) String() string {
	out, _ := lgva.MarshalYAML()
	return out.(string)
}

// MarshalYAML returns the YAML representation of a LazyGradedVestingAccount.
func (lgva LazyGradedVestingAccount) MarshalYAML() (interface{}, error) {
	alias := vestingAccountYAML{
		Address:          lgva.Address,
		Coins:            lgva.Coins,
		AccountNumber:    lgva.AccountNumber,
		Sequence:         lgva.Sequence,
		OriginalVesting:  lgva.OriginalVesting,
		DelegatedFree:    lgva.DelegatedFree,
		DelegatedVesting: lgva.DelegatedVesting,
		EndTime:          lgva.EndTime,
		VestingSchedules: lgva.VestingSchedules,
	}

	if lgva.PubKey != nil {
		pks, err := sdk.Bech32ifyPubKey(sdk.Bech32PubKeyTypeAccPub, lgva.GetPubKey())
		if err != nil {
			return nil, err
		}

		alias.PubKey = pks
	}

	bz, err := yaml.Marshal(alias)
	if err != nil {
		return nil, err
	}

	return string(bz), err
}

// MarshalJSON returns the JSON representation of a LazyGradedVestingAccount.
func (lgva LazyGradedVestingAccount) MarshalJSON() ([]byte, error) {
	alias := vestingAccountJSON{
		Address:          lgva.Address,
		Coins:            lgva.Coins,
		PubKey:           lgva.GetPubKey(),
		AccountNumber:    lgva.AccountNumber,
		Sequence:         lgva.Sequence,
		OriginalVesting:  lgva.OriginalVesting,
		DelegatedFree:    lgva.DelegatedFree,
		DelegatedVesting: lgva.DelegatedVesting,
		EndTime:          lgva.EndTime,
		VestingSchedules: lgva.VestingSchedules,
	}

	return codec.Cdc.MarshalJSON(alias)
}

// UnmarshalJSON unmarshals raw JSON bytes into a LazyGradedVestingAccount.
func (lgva *LazyGradedVestingAccount) UnmarshalJSON(bz []byte) error {
	var alias vestingAccountJSON
	if err := codec.Cdc.UnmarshalJSON(bz, &alias); err != nil {
		return err
	}

	lgva.BaseVestingAccount = &vesttypes.BaseVestingAccount{
		BaseAccount:      authtypes.NewBaseAccount(alias.Address, alias.Coins, alias.PubKey, alias.AccountNumber, alias.Sequence),
		OriginalVesting:  alias.OriginalVesting,
		DelegatedFree:    alias.DelegatedFree,
		DelegatedVesting: alias.DelegatedVesting,
		EndTime:          alias.EndTime,
	}

	lgva.VestingSchedules = alias.VestingSchedules

	return nil
}

// spendableCoins returns all the spendable coins for a vesting account given a
// set of vesting coins.
//
// CONTRACT: The account's coins, delegated vesting coins, vestingCoins must be
// sorted.
func (lgva LazyGradedVestingAccount) spendableCoins(vestingCoins sdk.Coins) sdk.Coins {
	var spendableCoins sdk.Coins
	bc := lgva.GetCoins()

	for _, coin := range bc {
		// zip/lineup all coins by their denomination to provide O(n) time
		baseAmt := coin.Amount
		vestingAmt := vestingCoins.AmountOf(coin.Denom)
		delVestingAmt := lgva.DelegatedVesting.AmountOf(coin.Denom)

		// compute min((BC + DV) - V, BC) per the specification
		min := sdk.MinInt(baseAmt.Add(delVestingAmt).Sub(vestingAmt), baseAmt)
		spendableCoin := sdk.NewCoin(coin.Denom, min)

		if !spendableCoin.IsZero() {
			spendableCoins = spendableCoins.Add(spendableCoin)
		}
	}

	return spendableCoins
}

// trackDelegation tracks a delegation amount for any given vesting account type
// given the amount of coins currently vesting. It returns the resulting base
// coins.
//
// CONTRACT: The account's coins, delegation coins, vesting coins, and delegated
// vesting coins must be sorted.
func (lgva *LazyGradedVestingAccount) trackDelegation(vestingCoins, amount sdk.Coins) {
	bc := lgva.GetCoins()

	for _, coin := range amount {
		// zip/lineup all coins by their denomination to provide O(n) time

		baseAmt := bc.AmountOf(coin.Denom)
		vestingAmt := vestingCoins.AmountOf(coin.Denom)
		delVestingAmt := lgva.DelegatedVesting.AmountOf(coin.Denom)

		// Panic if the delegation amount is zero or if the base coins does not
		// exceed the desired delegation amount.
		if coin.Amount.IsZero() || baseAmt.LT(coin.Amount) {
			panic("delegation attempt with zero coins or insufficient funds")
		}

		// compute x and y per the specification, where:
		// X := min(max(V - DV, 0), D)
		// Y := D - X
		x := sdk.MinInt(sdk.MaxInt(vestingAmt.Sub(delVestingAmt), sdk.ZeroInt()), coin.Amount)
		y := coin.Amount.Sub(x)

		if !x.IsZero() {
			xCoin := sdk.NewCoin(coin.Denom, x)
			lgva.DelegatedVesting = lgva.DelegatedVesting.Add(xCoin)
		}

		if !y.IsZero() {
			yCoin := sdk.NewCoin(coin.Denom, y)
			lgva.DelegatedFree = lgva.DelegatedFree.Add(yCoin)
		}
	}
}
