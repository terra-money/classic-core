package util

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// nolint
const (
	BlocksPerMinute = int64(12)
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

// GetBlocksPerEpoch gets the number of blocks processed in one epoch
func GetBlocksPerEpoch() int64 {
	return BlocksPerEpoch
}

// IsEpochLastBlock checks whether we are at the last block of the current epoch
func IsEpochLastBlock(ctx sdk.Context) bool {
	return (ctx.BlockHeight()+1)%BlocksPerEpoch == 0
}
