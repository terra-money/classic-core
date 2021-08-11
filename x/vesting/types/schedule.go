package types

import (
	"errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

//-----------------------------------------------------------------------------
// Schedule

// NewSchedule returns new Schedule instance
func NewSchedule(startTime, endTime int64, ratio sdk.Dec) Schedule {
	return Schedule{
		StartTime: startTime,
		EndTime:   endTime,
		Ratio:     ratio,
	}
}

// GetStartTime returns start time
func (s Schedule) GetStartTime() int64 {
	return s.StartTime
}

// GetEndTime returns end time
func (s Schedule) GetEndTime() int64 {
	return s.EndTime
}

// GetRatio returns ratio
func (s Schedule) GetRatio() sdk.Dec {
	return s.Ratio
}

// Validate checks that the lazy schedule is valid.
func (s Schedule) Validate() error {

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

// Schedules stores all lazy schedules
type Schedules []Schedule

//-----------------------------------------------------------------------------
// Vesting Schedule

// NewVestingSchedule creates a new vesting lazy schedule instance.
func NewVestingSchedule(denom string, schedules Schedules) VestingSchedule {
	return VestingSchedule{
		Denom:     denom,
		Schedules: schedules,
	}
}

// GetVestedRatio returns the ratio of tokens that have vested by blockTime.
func (vs VestingSchedule) GetVestedRatio(blockTime int64) sdk.Dec {
	sumRatio := sdk.ZeroDec()
	for _, lazySchedule := range vs.Schedules {
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

// GetDenom returns the denom of vesting schedule
func (vs VestingSchedule) GetDenom() string {
	return vs.Denom
}

// Validate checks that the vesting lazy schedule is valid.
func (vs VestingSchedule) Validate() error {
	sumRatio := sdk.ZeroDec()
	for _, lazySchedule := range vs.Schedules {
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

// VestingSchedules stores all vesting schedules passed as part of a LazyGradedVestingAccount
type VestingSchedules []VestingSchedule
