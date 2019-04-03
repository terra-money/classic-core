package budget

import (
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// no-lint
type Vote MsgVoteProgram

func (v Vote) String() string {
	return MsgVoteProgram(v).String()
}

// NewMsgVoteProgram creates a MsgVoteProgram instance
func NewVote(programID uint64, option bool, voter sdk.AccAddress) Vote {
	return Vote(
		MsgVoteProgram{
			ProgramID: programID,
			Option:    option,
			Voter:     voter,
		})
}

// Votes is a collection of Vote
type Votes []Vote

func (v Votes) String() (out string) {
	for _, val := range v {
		out += val.String() + "\n"
	}
	return strings.TrimSpace(out)
}
