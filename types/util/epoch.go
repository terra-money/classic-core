package util

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const blocksPerEpoch = 604800 // Approx. 1 week

// GetEpoch returns the current epoch, starting from 0
func GetEpoch(ctx sdk.Context) sdk.Int {
	curEpoch := ctx.BlockHeight() / blocksPerEpoch
	return sdk.NewInt(curEpoch)
}

func GetBlocksPerEpoch() int64 {
	return blocksPerEpoch
}
