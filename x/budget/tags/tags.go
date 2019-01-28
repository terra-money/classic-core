package tags

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Governance tags
var (
	ActionProgramWithdrawn = []byte("program-withdrawn")
	ActionProgramDropped   = []byte("program-dropped")
	ActionProgramPassed    = []byte("program-passed")
	ActionProgramRejected  = []byte("program-rejected")
	ActionProgramSubmitted = []byte("program-submitted")
	ActionProgramVote      = []byte("program-vote")

	Action            = sdk.TagAction
	Submitter         = "submitter"
	ProgramID         = "program-id"
	VotingPeriodStart = "voting-period-start"
	Executor          = "executor"
	Voter             = "voter"
	Option            = "option"
)
