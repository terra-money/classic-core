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
	prefixClaim     = []byte("claim")

	paramStoreKeyParams = []byte("params")
)

func keyProgram(programID uint64) []byte {
	return []byte(fmt.Sprintf("%s:%d", prefixProgram, programID))
}

func keyVote(programID uint64, voterAddr sdk.AccAddress) []byte {
	return []byte(fmt.Sprintf("%s:%d:%s", prefixVote, programID, voterAddr))
}

func prefixCandQueueEndBlock(endBlock int64) []byte {
	return []byte(fmt.Sprintf("%s:%020d", prefixCandQueue, endBlock))
}

func keyCandidate(endBlock int64, programID uint64) []byte {
	return []byte(fmt.Sprintf("%s:%020d:%d", prefixCandQueue, endBlock, programID))
}

func keyClaim(recipient sdk.AccAddress) []byte {
	return []byte(fmt.Sprintf("%s:%s", prefixClaim, recipient))
}

func paramKeyTable() params.KeyTable {
	return params.NewKeyTable(
		paramStoreKeyParams, Params{},
	)
}
