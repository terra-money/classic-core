package types

import (
	"bytes"
)

// GenesisState - all oracle state that must be provided at genesis
type GenesisState struct {
	Params      Params                  `json:"params" yaml:"params"`
	VotingInfos map[string]VotingInfo   `json:"voting_infos" yaml:"voting_infos"`
	MissedVotes map[string][]MissedVote `json:"missed_votes" yaml:"missed_votes"`
}

// NewGenesisState creates a new GenesisState object
func NewGenesisState(
	params Params, votingInfo map[string]VotingInfo, MissedVotes map[string][]MissedVote,
) GenesisState {

	return GenesisState{
		Params:      params,
		VotingInfos: votingInfo,
		MissedVotes: MissedVotes,
	}
}

// MissedVote
type MissedVote struct {
	Index  int64 `json:"index" yaml:"index"`
	Missed bool  `json:"missed" yaml:"missed"`
}

// NewMissedVote creates a new MissedVote instance
func NewMissedVote(index int64, missed bool) MissedVote {
	return MissedVote{
		Index:  index,
		Missed: missed,
	}
}

// DefaultGenesisState - default GenesisState used by columbus-2
func DefaultGenesisState() GenesisState {
	return GenesisState{
		Params:      DefaultParams(),
		VotingInfos: make(map[string]VotingInfo),
		MissedVotes: make(map[string][]MissedVote),
	}
}

// ValidateGenesis validates the oracle genesis parameters
func ValidateGenesis(data GenesisState) error {
	return data.Params.Validate()
}

// Checks whether 2 GenesisState structs are equivalent.
func (data GenesisState) Equal(data2 GenesisState) bool {
	b1 := ModuleCdc.MustMarshalBinaryBare(data)
	b2 := ModuleCdc.MustMarshalBinaryBare(data2)
	return bytes.Equal(b1, b2)
}

// Returns if a GenesisState is empty or has data in it
func (data GenesisState) IsEmpty() bool {
	emptyGenState := GenesisState{}
	return data.Equal(emptyGenState)
}
