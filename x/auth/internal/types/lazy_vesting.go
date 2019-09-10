package types

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
)

//-----------------------------------------------------------------------------
// Schedule

// Schedule no-lint
type LazySchedule struct {
	StartTime int64   `json:"start_time"`
	EndTime   int64   `json:"end_time"`
	Ratio     sdk.Dec `json:"ratio"`
}

func NewLazySchedule(startTime, endTime int64, ratio sdk.Dec) LazySchedule {
	return LazySchedule{
		StartTime: startTime,
		EndTime:   endTime,
		Ratio:     ratio,
	}
}

// GetStartTime returns start time
func (s LazySchedule) GetStartTime() int64 {
	return s.StartTime
}

// GetEndTime returns end time
func (s LazySchedule) GetEndTime() int64 {
	return s.EndTime
}

// GetRatio returns ratio
func (s LazySchedule) GetRatio() sdk.Dec {
	return s.Ratio
}

// String implements fmt.Stringer interface
func (s LazySchedule) String() string {
	return fmt.Sprintf(`LazySchedule:
	StartTime: %v,
	EndTime: %v,
	Ratio: %v`,
		s.StartTime, s.EndTime, s.Ratio)
}

// IsValid checks that the lazy schedule is valid.
func (s LazySchedule) IsValid() bool {

	startTime := s.GetStartTime()
	endTime := s.GetEndTime()
	ratio := s.GetRatio()

	return startTime >= 0 && endTime >= startTime && ratio.GT(sdk.ZeroDec())
}

//-----------------------------------------------------------------------------
// Vesting Lazy Schedule

// VestingSchedule maps the ratio of tokens that becomes vested by blocktime (in seconds) from genesis.
// The sum of values in the LazySchedule should sum to 1.0.
// CONTRACT: assumes that entries are
type VestingSchedule struct {
	Denom         string         `json:"denom"`
	LazySchedules []LazySchedule `json:"schedules"` // maps blocktime to percentage vested. Should sum to 1.
}

// NewVestingSchedule creates a new vesting lazy schedule instance.
func NewVestingSchedule(denom string, lazySchedules []LazySchedule) VestingSchedule {
	return VestingSchedule{
		Denom:         denom,
		LazySchedules: lazySchedules,
	}
}

// GetVestedRatio returns the ratio of tokens that have vested by blockTime.
func (vs VestingSchedule) GetVestedRatio(blockTime int64) sdk.Dec {
	sumRatio := sdk.ZeroDec()
	for _, lazySchedule := range vs.LazySchedules {
		startTime := lazySchedule.GetStartTime()
		endTime := lazySchedule.GetEndTime()
		ratio := lazySchedule.GetRatio()

		if blockTime < startTime {
			continue
		}

		if blockTime < endTime {
			ratio = ratio.MulInt64(blockTime - startTime).QuoInt64(endTime - startTime)
		}

		sumRatio = sumRatio.Add(ratio)

	}
	return sumRatio
}

// GetDenom returns the denom of vesting layz schedule
func (vs VestingSchedule) GetDenom() string {
	return vs.Denom
}

// IsValid checks that the vesting lazy schedule is valid.
func (vs VestingSchedule) IsValid() bool {
	sumRatio := sdk.ZeroDec()
	for _, lazySchedule := range vs.LazySchedules {

		if !lazySchedule.IsValid() {
			return false
		}

		sumRatio = sumRatio.Add(lazySchedule.GetRatio())
	}

	return sumRatio.Equal(sdk.OneDec())
}

// String implements fmt.Stringer interface
func (vs VestingSchedule) String() string {
	return fmt.Sprintf(`VestingSchedule:
	Denom: %v,
	LazySchedules: %v`,
		vs.Denom, vs.LazySchedules)
}

//-----------------------------------------------------------------------------
// Lazy Graded Vesting Account

// LazyGradedVestingAccount defines an account type that vests coins via a graded vesting lazy schedule.
type LazyGradedVestingAccount interface {
	auth.VestingAccount

	GetVestingSchedules() []VestingSchedule
	GetVestingSchedule(denom string) (VestingSchedule, bool)
}

// BaseLazyGradedVestingAccount implements the LazyGradedVestingAccount interface. It vests all
// coins according to a predefined schedule.
var _ LazyGradedVestingAccount = (*BaseLazyGradedVestingAccount)(nil)

// BaseLazyGradedVestingAccount implements the VestingAccount interface. It vests tokens according to
// a predefined set of vesting schedule.
type BaseLazyGradedVestingAccount struct {
	*auth.BaseVestingAccount

	VestingSchedules []VestingSchedule `json:"vesting_schedules"`
}

// NewBaseLazyGradedVestingAccount creates a new BaseLazyGradedVestingAccount object from BaseVestingAccount
func NewBaseLazyGradedVestingAccountRaw(baseVestingAcc *auth.BaseVestingAccount, lazyVestingSchedules []VestingSchedule) *BaseLazyGradedVestingAccount {
	return &BaseLazyGradedVestingAccount{baseVestingAcc, lazyVestingSchedules}
}

// NewBaseLazyGradedVestingAccount returns a new BaseLazyGradedVestingAccount
func NewBaseLazyGradedVestingAccount(baseAcc *auth.BaseAccount, lazyVestingSchedules []VestingSchedule) *BaseLazyGradedVestingAccount {
	baseVestingAcc := &auth.BaseVestingAccount{
		BaseAccount:     baseAcc,
		OriginalVesting: baseAcc.Coins,
		EndTime:         0,
	}

	return &BaseLazyGradedVestingAccount{baseVestingAcc, lazyVestingSchedules}
}

// GetVestingSchedules returns the VestingSchedules of the graded lazy vesting account
func (lgva BaseLazyGradedVestingAccount) GetVestingSchedules() []VestingSchedule {
	return lgva.VestingSchedules
}

// GetVestingSchedule returns the VestingSchedule of the given denom
func (lgva BaseLazyGradedVestingAccount) GetVestingSchedule(denom string) (VestingSchedule, bool) {
	for _, vs := range lgva.VestingSchedules {
		if vs.Denom == denom {
			return vs, true
		}
	}

	return VestingSchedule{}, false
}

// GetVestedCoins returns the total amount of vested coins for a graded vesting
// account. All coins are vested continuously once the schedule's StartTime has elapsed until EndTime.
func (lgva BaseLazyGradedVestingAccount) GetVestedCoins(blockTime time.Time) sdk.Coins {
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
func (lgva BaseLazyGradedVestingAccount) GetVestingCoins(blockTime time.Time) sdk.Coins {
	return lgva.OriginalVesting.Sub(lgva.GetVestedCoins(blockTime))
}

// SpendableCoins returns the total number of spendable coins for a graded
// vesting account.
func (lgva BaseLazyGradedVestingAccount) SpendableCoins(blockTime time.Time) sdk.Coins {
	return lgva.spendableCoins(lgva.GetVestingCoins(blockTime))
}

// TrackDelegation tracks a desired delegation amount by setting the appropriate
// values for the amount of delegated vesting, delegated free, and reducing the
// overall amount of base coins.
func (lgva *BaseLazyGradedVestingAccount) TrackDelegation(blockTime time.Time, amount sdk.Coins) {
	lgva.trackDelegation(lgva.GetVestingCoins(blockTime), amount)
}

// GetStartTime returns zero since a lazy graded vesting account has no start time.
func (lgva BaseLazyGradedVestingAccount) GetStartTime() int64 {
	return 0
}

// GetEndTime returns zero since a lazy graded vesting account has no end time.
func (lgva BaseLazyGradedVestingAccount) GetEndTime() int64 {
	return 0
}

// String implements fmt.Stringer interface
func (lgva BaseLazyGradedVestingAccount) String() string {
	var pubkey string

	if lgva.PubKey != nil {
		pubkey = sdk.MustBech32ifyAccPub(lgva.PubKey)
	}

	return fmt.Sprintf(`BaseLazyGradedVestingAccount:
  Address:          %s
  Pubkey:           %s
  Coins:            %s
  AccountNumber:    %d
  Sequence:         %d
  OriginalVesting:  %s
  DelegatedFree:    %s
  DelegatedVesting: %s
  VestingSchedules:        %v `,
		lgva.Address, pubkey, lgva.Coins, lgva.AccountNumber, lgva.Sequence,
		lgva.OriginalVesting, lgva.DelegatedFree, lgva.DelegatedVesting,
		lgva.VestingSchedules,
	)
}

// spendableCoins returns all the spendable coins for a vesting account given a
// set of vesting coins.
//
// CONTRACT: The account's coins, delegated vesting coins, vestingCoins must be
// sorted.
func (lgva BaseLazyGradedVestingAccount) spendableCoins(vestingCoins sdk.Coins) sdk.Coins {
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
			spendableCoins = spendableCoins.Add(sdk.Coins{spendableCoin})
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
func (lgva *BaseLazyGradedVestingAccount) trackDelegation(vestingCoins, amount sdk.Coins) {
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
			lgva.DelegatedVesting = lgva.DelegatedVesting.Add(sdk.Coins{xCoin})
		}

		if !y.IsZero() {
			yCoin := sdk.NewCoin(coin.Denom, y)
			lgva.DelegatedFree = lgva.DelegatedFree.Add(sdk.Coins{yCoin})
		}

		lgva.Coins = lgva.Coins.Sub(sdk.Coins{coin})
	}
}
