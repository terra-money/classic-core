package types

import (
	"errors"
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

//-----------------------------------------------------------------------------
// Schedule

// LazySchedule defines a vesting schedule which is used for LazyGradedVestingAccount
type LazySchedule struct {
	StartTime int64   `json:"start_time"`
	EndTime   int64   `json:"end_time"`
	Ratio     sdk.Dec `json:"ratio"`
}

// NewLazySchedule returns new LazySchedule instance
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

// Validate checks that the lazy schedule is valid.
func (s LazySchedule) Validate() error {

	startTime := s.GetStartTime()
	endTime := s.GetEndTime()
	ratio := s.GetRatio()

	if startTime < 0 {
		return errors.New("vesting start-time cannot be negative")
	}

	if endTime < startTime {
		return errors.New("vesting start-time cannot be before end-time")
	}

	if ratio.LTE(sdk.ZeroDec()) {
		return errors.New("vesting ratio cannot be smaller than or equal with zero")
	}

	return nil
}

// LazySchedules stores all lazy schedules
type LazySchedules []LazySchedule

// String implements stringer interface
func (vs LazySchedules) String() string {
	lazySchedulesListString := make([]string, len(vs))
	for _, lazySchedule := range vs {
		lazySchedulesListString = append(lazySchedulesListString, lazySchedule.String())
	}
	return strings.TrimSpace(fmt.Sprintf(`Lazy Schedules:
		%s`, strings.Join(lazySchedulesListString, ", ")))
}

//-----------------------------------------------------------------------------
// Vesting Lazy Schedule

// VestingSchedule maps the ratio of tokens that becomes vested by blocktime (in seconds) from genesis.
// The sum of values in the LazySchedule should sum to 1.0.
// CONTRACT: assumes that entries are
type VestingSchedule struct {
	Denom         string        `json:"denom"`
	LazySchedules LazySchedules `json:"schedules"` // maps blocktime to percentage vested. Should sum to 1.
}

// NewVestingSchedule creates a new vesting lazy schedule instance.
func NewVestingSchedule(denom string, lazySchedules LazySchedules) VestingSchedule {
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

// Validate checks that the vesting lazy schedule is valid.
func (vs VestingSchedule) Validate() error {
	sumRatio := sdk.ZeroDec()
	for _, lazySchedule := range vs.LazySchedules {

		if err := lazySchedule.Validate(); err != nil {
			return err
		}

		sumRatio = sumRatio.Add(lazySchedule.GetRatio())
	}

	// add rounding to allow language specific calculation errors
	const fixedPointDecimals = 1000000000
	if !sumRatio.MulInt64(fixedPointDecimals).RoundInt().
		ToDec().QuoInt64(fixedPointDecimals).Equal(sdk.OneDec()) {
		return errors.New("vesting total ratio must be one")
	}

	return nil
}

// String implements fmt.Stringer interface
func (vs VestingSchedule) String() string {
	return fmt.Sprintf(`VestingSchedule:
	Denom: %v,
	LazySchedules: %v`,
		vs.Denom, vs.LazySchedules)
}

// VestingSchedules stores all vesting schedules passed as part of a LazyGradedVestingAccount
type VestingSchedules []VestingSchedule

// String implements stringer interface
func (vs VestingSchedules) String() string {
	vestingSchedulesListString := make([]string, len(vs))
	for _, vestingSchedule := range vs {
		vestingSchedulesListString = append(vestingSchedulesListString, vestingSchedule.String())
	}
	return strings.TrimSpace(fmt.Sprintf(`Vesting Schedules:
		%s`, strings.Join(vestingSchedulesListString, ", ")))
}
