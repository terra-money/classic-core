# Budget

The Budget module governs how a portion of Terra seigniorage can be deployed via a distributed governance of Terra validators to drive Terra's adoption. 

## Overview 

A portion of Terra's growth (seigniorage) is routed to budget programs continuously. Therefore, long-lasting institutions (such as an ecosystem development fund, a bug bounty program) is more suitable for the budget rather than one-off proposals. 

At the end of every treasury update cycle, a portion of seigniorage collected minus the amount burned for mining rewards (1 - `MiningRewardWeight`) is routed to the budget to be distributed among programs. 

Each active program is associated with a weight, which is the sum of voting staking power in support minus against (yes votes - no votes). At the end of the budget `VotePeriod`, the seigniorage routed from the treasury is disbursed pro-rata to the program weights. 

Though we expect budget rewards to be quite random close to genesis, we expect that in time budget programs that offer the highest returns to the community and sets a high bar for transparency will rise above the pack. 

## Budget program 

```golang
// Program defines the basic properties of a staking Program
type Program struct {
	ProgramID   uint64         `json:"program_id"`  // ID of the Program
	Title       string         `json:"title"`       // Title of the Program
	Description string         `json:"description"` // Description of the Program
	Submitter   sdk.AccAddress `json:"submitter"`   // Validator address of the proposer
	Executor    sdk.AccAddress `json:"executor"`    // Account address of the executor
	SubmitBlock int64          `json:"submit_time"` // Block height from which the Program is open for votations
}
```

The budget program contains simple metadata about the program, such as title, description, submitter, and executor. 

In order to submit a budget program for consideration, a `MsgSubmitProgram` must be submitted, which will require a small deposit to be paid to prevent spamming. 

In order to withdraw a budget program that is still being considered or in the active set, the Submitter can send a `MsgWithdrawProgram`, which will remove the program from the store and refund the deposit. 

To vote on programs, either in the candidate or active set, the validator must submit a `MsgVoteProgram` with a binary option in support or against.

The validator is not obligated to vote on any budget programs (for now). 

## Program states 

### Candidate state

Programs that are newly submitted and satisfies the condition `SubmitBlock + VotePeriod > ctx.BlockHeight()` are in the candidate state. When the `VotePeriod` has expired since the submitted block, votes are tallied on the program, and if the program's weight is greater than the `ActiveThreshold` it is transitioned to the active state. Otherwise, it is simply dropped from the store and the submit deposit is burned. 

### Withdrawn state

Programs that are withdrawn while still in the candidate / active state are withdrawn, and the submit deposit is returned to the submitter. Only the submitter may send a `MsgWithdrawProgram` transaction. 

### Active state 

Programs that are in the active state receive budget subsidies. At each `VotePeriod`, their weights are readjusted to reflect the votes of validators. If an active program's weight falls below `LegacyThreshold`, it enters a legacied state and is deleted from the store. 

### Legacied state

Active programs that fell out of favor.


## Parameters 

```golang
// Params budget parameters
type Params struct {
	ActiveThreshold sdk.Dec  `json:"active_threshold"` // threshold of vote that will transition a program open -> active budget queue
	LegacyThreshold sdk.Dec  `json:"legacy_threshold"` // threshold of vote that will transition a program active -> legacy budget queue
	VotePeriod      int64    `json:"vote_period"`      // vote period
	Deposit         sdk.Coin `json:"deposit"`          // Minimum deposit in TerraSDR
}
```