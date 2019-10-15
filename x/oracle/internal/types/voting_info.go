package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// VotingInfo defines the voting info for a validator
type VotingInfo struct {
	Address            sdk.ValAddress `json:"address" yaml:"address"`                           // validator consensus address
	StartHeight        int64          `json:"start_height" yaml:"start_height"`                 // height at which validator was first a candidate OR was unjailed
	IndexOffset        int64          `json:"index_offset" yaml:"index_offset"`                 // index offset into signed block bit array
	MissedVotesCounter int64          `json:"missed_votes_counter" yaml:"missed_votes_counter"` // missed blocks counter (to avoid scanning the array every time)
}

// NewVotingInfo creates a new NewVotingInfo instance
func NewVotingInfo(
	valAddr sdk.ValAddress, startHeight,
	indexOffset, missedvotesCounter int64,
) VotingInfo {

	return VotingInfo{
		Address:            valAddr,
		StartHeight:        startHeight,
		IndexOffset:        indexOffset,
		MissedVotesCounter: missedvotesCounter,
	}
}

// String implements fmt.Stringer interface
func (i VotingInfo) String() string {
	return fmt.Sprintf(`Validator Signing Info:
  Address:               %s
  Start Height:          %d
  Index Offset:          %d
  Missed Votes Counter: %d`,
		i.Address, i.StartHeight,
		i.IndexOffset, i.MissedVotesCounter)
}
