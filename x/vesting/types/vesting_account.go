package types

import (
	fmt "fmt"
	"time"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	vestexported "github.com/cosmos/cosmos-sdk/x/auth/vesting/exported"
	vesttypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"

	"gopkg.in/yaml.v2"
)

// for pretty purpose
type vestingAccountYAML struct {
	Address          sdk.AccAddress `json:"address" yaml:"address"`
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

//-----------------------------------------------------------------------------
// Lazy Graded Vesting Account

// LazyGradedVestingAccount implements the LazyGradedVestingAccount interface. It vests all
// coins according to a predefined schedule.
var _ vestexported.VestingAccount = (*LazyGradedVestingAccount)(nil)
var _ authtypes.GenesisAccount = (*LazyGradedVestingAccount)(nil)

// NewLazyGradedVestingAccountRaw creates a new LazyGradedVestingAccount object from BaseVestingAccount
func NewLazyGradedVestingAccountRaw(baseVestingAcc *vesttypes.BaseVestingAccount, lazyVestingSchedules VestingSchedules) *LazyGradedVestingAccount {
	return &LazyGradedVestingAccount{
		BaseVestingAccount: baseVestingAcc,
		VestingSchedules:   lazyVestingSchedules,
	}
}

// NewLazyGradedVestingAccount returns a new LazyGradedVestingAccount
func NewLazyGradedVestingAccount(baseAcc *authtypes.BaseAccount, originalVesting sdk.Coins, lazyVestingSchedules VestingSchedules) *LazyGradedVestingAccount {
	baseVestingAcc := &vesttypes.BaseVestingAccount{
		BaseAccount:      baseAcc,
		OriginalVesting:  originalVesting,
		DelegatedFree:    sdk.NewCoins(),
		DelegatedVesting: sdk.NewCoins(),
		EndTime:          0,
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

// LockedCoins returns the set of coins that are not spendable (i.e. locked).
func (lgva LazyGradedVestingAccount) LockedCoins(blockTime time.Time) sdk.Coins {
	return lgva.BaseVestingAccount.LockedCoinsFromVesting(lgva.GetVestingCoins(blockTime))
}

// TrackDelegation tracks a delegation amount for any given vesting account type
// given the amount of coins currently vesting. It returns the resulting base
// coins.
//
// CONTRACT: The account's coins, delegation coins, vesting coins, and delegated
// vesting coins must be sorted.
func (lgva *LazyGradedVestingAccount) TrackDelegation(blockTime time.Time, balance, amount sdk.Coins) {
	lgva.BaseVestingAccount.TrackDelegation(balance, lgva.GetVestingCoins(blockTime), amount)
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
	denomMap := make(map[string]bool)
	for _, vestingSchedule := range lgva.GetVestingSchedules() {
		if _, ok := denomMap[vestingSchedule.Denom]; ok {
			return fmt.Errorf("cannot have multiple vesting schedules for %s", vestingSchedule.Denom)
		}

		if err := vestingSchedule.Validate(); err != nil {
			return err
		}

		denomMap[vestingSchedule.Denom] = true
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
	accAddr, err := sdk.AccAddressFromBech32(lgva.Address)
	if err != nil {
		return nil, err
	}

	out := vestingAccountYAML{
		Address:          accAddr,
		AccountNumber:    lgva.AccountNumber,
		PubKey:           getPKString(lgva),
		Sequence:         lgva.Sequence,
		OriginalVesting:  lgva.OriginalVesting,
		DelegatedFree:    lgva.DelegatedFree,
		DelegatedVesting: lgva.DelegatedVesting,
		EndTime:          lgva.EndTime,
		VestingSchedules: lgva.VestingSchedules,
	}

	return marshalYaml(out)
}

type getPK interface {
	GetPubKey() cryptotypes.PubKey
}

func getPKString(g getPK) string {
	if pk := g.GetPubKey(); pk != nil {
		return pk.String()
	}
	return ""
}

func marshalYaml(i interface{}) (interface{}, error) {
	bz, err := yaml.Marshal(i)
	if err != nil {
		return nil, err
	}
	return string(bz), nil
}
