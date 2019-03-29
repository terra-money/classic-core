package budget

import (
	"fmt"

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
	return []byte(fmt.Sprintf("%s:%d", prefixProgram, programID))
}

func keyVote(programID uint64, voterAddr sdk.AccAddress) []byte {
	return []byte(fmt.Sprintf("%s:%d:%s", prefixVote, programID, voterAddr))
}

func prefixVoteForProgram(programID uint64) []byte {
	return []byte(fmt.Sprintf("%s:%d", prefixVote, programID))
}

func prefixCandQueueEndBlock(endBlock int64) []byte {
	return []byte(fmt.Sprintf("%s:%d", prefixCandQueue, endBlock))
}

func keyCandidate(endBlock int64, programID uint64) []byte {
	return []byte(fmt.Sprintf("%s:%d:%d", prefixCandQueue, endBlock, programID))
}

func paramKeyTable() params.KeyTable {
	return params.NewKeyTable(
		paramStoreKeyParams, Params{},
	)
}
