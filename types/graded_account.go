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
type Schedule struct {
	Cliff int64   `json:"cliff"`
	Ratio sdk.Dec `json:"ratio"`
}

func NewSchedule(cliff int64, ratio sdk.Dec) Schedule {
	return Schedule{
		Cliff: cliff,
		Ratio: ratio,
	}
}

// GetRatio returns cliff
func (s Schedule) GetCliff() int64 {
	return s.Cliff
}

// GetRatio returns ratio
func (s Schedule) GetRatio() sdk.Dec {
	return s.Ratio
}

// String implements the fmt.Stringer interface
func (s Schedule) String() string {
	return fmt.Sprintf(`Schedule:
	Cliff: %v,
	Ratio: %v`,
		s.Cliff, s.Ratio)
}

// IsValid checks that the schedule is valid.
func (s Schedule) IsValid() bool {

	cliff := s.GetCliff()
	ratio := s.GetRatio()

	return cliff >= 0 && ratio.GT(sdk.ZeroDec())
}

//-----------------------------------------------------------------------------
// Vesting Schedule

// VestingSchedule maps the ratio of tokens that becomes vested by blocktime (in nanoseconds) from genesis.
// The sum of values in the Schedule should sum to 1.0.
// CONTRACT: assumes that entries are
type VestingSchedule struct {
	Denom     string     `json:"denom"`
	Schedules []Schedule `json:"schedules"` // maps blocktime to percentage vested. Should sum to 1.
}

// NewVestingSchedule creates a new vesting schedule instance.
func NewVestingSchedule(denom string, schedules []Schedule) VestingSchedule {
	return VestingSchedule{
		Denom:     denom,
		Schedules: schedules,
	}
}

// GetVestedRatio returns the ratio of tokens that have vested by blockTime.
func (vs VestingSchedule) GetVestedRatio(blockTime int64) sdk.Dec {
	sumRatio := sdk.ZeroDec()
	for _, schedule := range vs.Schedules {
		cliff := schedule.GetCliff()
		ratio := schedule.GetRatio()

		if blockTime >= cliff {
			sumRatio = sumRatio.Add(ratio)
		}
	}
	return sumRatio
}

// GetDenom returns the denom of vesting schedule
func (vs VestingSchedule) GetDenom() string {
	return vs.Denom
}

// IsValid checks that the vestingschedule is valid.
func (vs VestingSchedule) IsValid() bool {
	sumRatio := sdk.ZeroDec()
	for _, schedule := range vs.Schedules {

		if !schedule.IsValid() {
			return false
		}

		sumRatio = sumRatio.Add(schedule.GetRatio())
	}

	return sumRatio.Equal(sdk.OneDec())
}

// String implements the fmt.Stringer interface
func (vs VestingSchedule) String() string {
	return fmt.Sprintf(`VestingSchedule:
	Denom: %v,
	Schedules: %v`,
		vs.Denom, vs.Schedules)
}

//-----------------------------------------------------------------------------
// Graded Vesting Account

// GradedVestingAccount implements the VestingAccount interface. It vests all
// coins according to a predefined schedule.
var _ auth.VestingAccount = (*GradedVestingAccount)(nil)

// GradedVestingAccount implements the VestingAccount interface. It vests tokens according to
// a predefined set of vesting cliffs.
type GradedVestingAccount struct {
	*auth.BaseVestingAccount

	VestingSchedules []VestingSchedule `json:"vesting_schedules"`
}

// NewGradedVestingAccount returns a GradedVestingAccount
func NewGradedVestingAccount(baseAcc *auth.BaseAccount, vestingSchedules []VestingSchedule) *GradedVestingAccount {
	baseVestingAcc := &auth.BaseVestingAccount{
		BaseAccount:     baseAcc,
		OriginalVesting: baseAcc.Coins,
		EndTime:         0,
	}

	return &GradedVestingAccount{baseVestingAcc, vestingSchedules}
}

// GetSchedules returns the VestingSchedules of the graded vesting account
func (gva GradedVestingAccount) GetVestingSchedules() []VestingSchedule {
	return gva.VestingSchedules
}

// GetVestingSchedule returns the VestingSchedule of the given denom
func (gva GradedVestingAccount) GetVestingSchedule(denom string) (VestingSchedule, bool) {
	for _, vs := range gva.VestingSchedules {
		if vs.Denom == denom {
			return vs, true
		}
	}

	return VestingSchedule{}, false
}

// GetVestedCoins returns the total amount of vested coins for a graded vesting
// account. All coins are only vested once the schedule has elapsed.
func (gva GradedVestingAccount) GetVestedCoins(blockTime time.Time) sdk.Coins {
	var vestedCoins sdk.Coins
	for _, ovc := range gva.OriginalVesting {
		if vestingSchedule, exists := gva.GetVestingSchedule(ovc.Denom); exists {
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
func (gva GradedVestingAccount) GetVestingCoins(blockTime time.Time) sdk.Coins {
	return gva.OriginalVesting.Sub(gva.GetVestedCoins(blockTime))
}

// SpendableCoins returns the total number of spendable coins for a graded
// vesting account.
func (gva GradedVestingAccount) SpendableCoins(blockTime time.Time) sdk.Coins {
	return gva.spendableCoins(gva.GetVestingCoins(blockTime))
}

// TrackDelegation tracks a desired delegation amount by setting the appropriate
// values for the amount of delegated vesting, delegated free, and reducing the
// overall amount of base coins.
func (gva *GradedVestingAccount) TrackDelegation(blockTime time.Time, amount sdk.Coins) {
	gva.trackDelegation(gva.GetVestingCoins(blockTime), amount)
}

// GetStartTime returns zero since a graded vesting account has no start time.
func (gva *GradedVestingAccount) GetStartTime() int64 {
	return 0
}

// GetEndTime returns the time when vesting ends for a graded vesting account.
func (gva *GradedVestingAccount) GetEndTime() int64 {
	return 0
}

func (gva GradedVestingAccount) String() string {
	var pubkey string

	if gva.PubKey != nil {
		pubkey = sdk.MustBech32ifyAccPub(gva.PubKey)
	}

	return fmt.Sprintf(`Graded Vesting Account:
  Address:          %s
  Pubkey:           %s
  Coins:            %s
  AccountNumber:    %d
  Sequence:         %d
  OriginalVesting:  %s
  DelegatedFree:    %s
  DelegatedVesting: %s
  VestingSchedules:        %v `,
		gva.Address, pubkey, gva.Coins, gva.AccountNumber, gva.Sequence,
		gva.OriginalVesting, gva.DelegatedFree, gva.DelegatedVesting,
		gva.VestingSchedules,
	)
}

// spendableCoins returns all the spendable coins for a vesting account given a
// set of vesting coins.
//
// CONTRACT: The account's coins, delegated vesting coins, vestingCoins must be
// sorted.
func (gva GradedVestingAccount) spendableCoins(vestingCoins sdk.Coins) sdk.Coins {
	var spendableCoins sdk.Coins
	bc := gva.GetCoins()

	for _, coin := range bc {
		// zip/lineup all coins by their denomination to provide O(n) time
		baseAmt := coin.Amount
		vestingAmt := vestingCoins.AmountOf(coin.Denom)
		delVestingAmt := gva.DelegatedVesting.AmountOf(coin.Denom)

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
func (gva *GradedVestingAccount) trackDelegation(vestingCoins, amount sdk.Coins) {
	bc := gva.GetCoins()

	for _, coin := range amount {
		// zip/lineup all coins by their denomination to provide O(n) time

		baseAmt := bc.AmountOf(coin.Denom)
		vestingAmt := vestingCoins.AmountOf(coin.Denom)
		delVestingAmt := gva.DelegatedVesting.AmountOf(coin.Denom)

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
			gva.DelegatedVesting = gva.DelegatedVesting.Add(sdk.Coins{xCoin})
		}

		if !y.IsZero() {
			yCoin := sdk.NewCoin(coin.Denom, y)
			gva.DelegatedFree = gva.DelegatedFree.Add(sdk.Coins{yCoin})
		}

		gva.Coins = gva.Coins.Sub(sdk.Coins{coin})
	}
}
