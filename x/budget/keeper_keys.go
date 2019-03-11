package budget

import (
	"bytes"
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
)

// nolint
var (
	KeyDelimiter         = []byte(":")
	PrefixProgram        = []byte("program")
	PrefixVote           = []byte("vote")
	KeyNextProgramID     = []byte("new-program-id")
	PrefixCandidateQueue = []byte("candidate-queue")
	ParamStoreKeyParams  = []byte("params")
	DefaultParamspace    = "budget"
)

// KeyProgram creates a key of the form "Programs"|{state}|{ProgramID}
func KeyProgram(programID uint64) []byte {
	return []byte(fmt.Sprintf("%s:%d", PrefixProgram, programID))
}

// Key for getting a specific vote from the store
func KeyVote(programID uint64, voterAddr sdk.AccAddress) []byte {
	return []byte(fmt.Sprintf("%s:%d:%s", PrefixVote, programID, voterAddr))
}

// Key for getting a specific vote from the store
func PrefixVoteForProgram(programID uint64) []byte {
	return []byte(fmt.Sprintf("%s:%d", PrefixVote, programID))
}

// Returns the key for a programID in the activeprogramQueue
func PrefixCandidateQueueTime(endTime time.Time) []byte {
	return bytes.Join([][]byte{
		PrefixCandidateQueue,
		sdk.FormatTimeBytes(endTime),
	}, KeyDelimiter)
}

// Returns the key for a programID in the activeprogramQueue
func KeyCandidate(endTime time.Time, programID uint64) []byte {
	return bytes.Join([][]byte{
		PrefixCandidateQueue,
		sdk.FormatTimeBytes(endTime),
		sdk.Uint64ToBigEndian(programID),
	}, KeyDelimiter)
}

// ParamTable for budget module
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable(
		ParamStoreKeyParams, Params{},
	)
}
