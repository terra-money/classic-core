package budget

import (
	"bytes"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
)

// nolint
var (
	keyDelimiter     = []byte(":")
	keyNextProgramID = []byte("new-program-id")

	prefixProgram   = []byte("program")
	prefixVote      = []byte("vote")
	prefixCandQueue = []byte("candidate-queue")

	paramStoreKeyParams = []byte("params")
)

func keyProgram(programID uint64) []byte {
	return bytes.Join([][]byte{
		prefixProgram,
		sdk.Uint64ToBigEndian(programID),
	}, keyDelimiter)
}

func keyVote(programID uint64, voterAddr sdk.AccAddress) []byte {
	return bytes.Join([][]byte{
		prefixVote,
		sdk.Uint64ToBigEndian(programID),
		voterAddr,
	}, keyDelimiter)
}

func prefixVoteForProgram(programID uint64) []byte {
	return bytes.Join([][]byte{
		prefixVote,
		sdk.Uint64ToBigEndian(programID),
	}, keyDelimiter)
}

func prefixCandQueueEndBlock(endBlock int64) []byte {
	return bytes.Join([][]byte{
		prefixCandQueue,
		sdk.Uint64ToBigEndian(uint64(endBlock)),
	}, keyDelimiter)
}

func keyCandidate(endBlock int64, programID uint64) []byte {
	return bytes.Join([][]byte{
		prefixCandQueue,
		sdk.Uint64ToBigEndian(uint64(endBlock)),
		sdk.Uint64ToBigEndian(programID),
	}, keyDelimiter)
}

func paramKeyTable() params.KeyTable {
	return params.NewKeyTable(
		paramStoreKeyParams, Params{},
	)
}
