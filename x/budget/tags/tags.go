package tags

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Governance tags
var (
	ActionProgramWithdrawn = "program-withdrawn"
	ActionProgramDropped   = "program-dropped"
	ActionProgramPassed    = "program-passed"
	ActionProgramRejected  = "program-rejected"
	ActionProgramSubmitted = "program-submitted"
	ActionProgramVote      = "program-vote"

	Action            = sdk.TagAction
	Submitter         = "submitter"
	ProgramID         = "program-id"
	VotingPeriodStart = "voting-period-start"
	Executor          = "executor"
	Voter             = "voter"
	Option            = "option"
)
