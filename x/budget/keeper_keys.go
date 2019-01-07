package budget

import (
	"bytes"
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	KeyDelimiter     = []byte(":")
	PrefixProgram    = []byte("program")
	PrefixVote       = []byte("vote")
	KeyNextProgramID = []byte("newProgramID")

	// PrefixActiveProgramQueue   = []byte("activeProgramQueue")
	PrefixInactiveProgramQueue = []byte("inactiveProgramQueue")
)

// GenerateProgramKey creates a key of the form "Programs"|{state}|{ProgramID}
func KeyProgram(programID int64) []byte {
	return []byte(fmt.Sprintf("%s:%d", KeyProgram, programID))
}

// Key for getting a specific vote from the store
func KeyVote(programID uint64, voterAddr sdk.AccAddress) []byte {
	return []byte(fmt.Sprintf("%s:%d:%d", PrefixVote, programID, voterAddr))
}

// // Returns the key for a programID in the activeprogramQueue
// func PrefixActiveProgramQueueTime() []byte {
// 	return PrefixActiveProgramQueue
// }

// // Returns the key for a programID in the activeprogramQueue
// func KeyActiveProgramQueueProgram(programID uint64) []byte {
// 	return bytes.Join([][]byte{
// 		PrefixActiveProgramQueue,
// 		sdk.Uint64ToBigEndian(programID),
// 	}, KeyDelimiter)
// }

// Returns the key for a programID in the activeprogramQueue
func PrefixInactiveProgramQueueTime(endTime time.Time) []byte {
	return bytes.Join([][]byte{
		PrefixInactiveProgramQueue,
		sdk.FormatTimeBytes(endTime),
	}, KeyDelimiter)
}

// Returns the key for a programID in the activeprogramQueue
func KeyInactiveProgramQueueProgram(endTime time.Time, programID uint64) []byte {
	return bytes.Join([][]byte{
		PrefixInactiveProgramQueue,
		sdk.FormatTimeBytes(endTime),
		sdk.Uint64ToBigEndian(programID),
	}, KeyDelimiter)
}
