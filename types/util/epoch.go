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

	BlocksPerEpoch = BlocksPerWeek
)

// GetEpoch returns the current epoch, starting from 0
func GetEpoch(ctx sdk.Context) sdk.Int {
	curEpoch := ctx.BlockHeight() / BlocksPerEpoch
	return sdk.NewInt(curEpoch)
}

// IsPeriodLastBlock returns true if we are at the last block of the period
func IsPeriodLastBlock(ctx sdk.Context, blocksPerPeriod int64) bool {
	return (ctx.BlockHeight()+1)%blocksPerPeriod == 0
}
