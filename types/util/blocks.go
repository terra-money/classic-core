package util

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// nolint
const (
	BlocksPerMinute = int64(10)
	BlocksPerHour   = BlocksPerMinute * 60
	BlocksPerDay    = BlocksPerHour * 24
	BlocksPerWeek   = BlocksPerDay * 7
	BlocksPerMonth  = BlocksPerDay * 30
	BlocksPerYear   = BlocksPerDay * 365
)

// IsPeriodLastBlock returns true if we are at the last block of the period
func IsPeriodLastBlock(ctx sdk.Context, blocksPerPeriod int64) bool {
	return (ctx.BlockHeight()+1)%blocksPerPeriod == 0
}
