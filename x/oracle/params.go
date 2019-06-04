package oracle

import (
	"fmt"

	"github.com/terra-project/core/types/util"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Params oracle parameters
type Params struct {
	VotePeriod       int64   `json:"vote_period"`        // voting period in block height; tallys and reward claim period
	VoteThreshold    sdk.Dec `json:"vote_threshold"`     // minimum stake power threshold to update price
	OracleRewardBand sdk.Dec `json:"oracle_reward_band"` // band around the oracle weighted median to reward
}

// NewParams creates a new param instance
func NewParams(votePeriod int64, voteThreshold sdk.Dec, oracleRewardBand sdk.Dec) Params {
	return Params{
		VotePeriod:       votePeriod,
		VoteThreshold:    voteThreshold,
		OracleRewardBand: oracleRewardBand,
	}
}

// DefaultParams creates default oracle module parameters
func DefaultParams() Params {
	return NewParams(
		util.BlocksPerMinute,      // 1 minute
		sdk.NewDecWithPrec(50, 2), // 50%
		sdk.NewDecWithPrec(1, 2),  // 1%
	)
}

func validateParams(params Params) error {
	if params.VotePeriod <= 0 {
		return fmt.Errorf("oracle parameter VotePeriod must be > 0, is %d", params.VotePeriod)
	}
	if params.VoteThreshold.LTE(sdk.NewDecWithPrec(33, 2)) {
		return fmt.Errorf("oracle parameter VoteTheshold must be greater than 33 percent")
	}
	if params.OracleRewardBand.IsNegative() {
		return fmt.Errorf("oracle parameter OracleRewardBand must be positive")
	}
	return nil
}

func (params Params) String() string {
	return fmt.Sprintf(`Oracle Params:
  VotePeriod:          %d
  VoteThreshold:       %s
  OracleRewardBand:    %s
  `, params.VotePeriod, params.VoteThreshold, params.OracleRewardBand)
}
