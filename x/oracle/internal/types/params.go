package types

import (
	"fmt"

	core "github.com/terra-project/core/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params/subspace"
)

// DefaultParamspace
const DefaultParamspace = ModuleName

// Parameter keys
var (
	ParamStoreKeyVotePeriod             = []byte("voteperiod")
	ParamStoreKeyVoteThreshold          = []byte("votethreshold")
	ParamStoreKeyRewardBand             = []byte("rewardband")
	ParamStoreKeyRewardFraction         = []byte("rewardfraction")
	ParamStoreKeyVotesWindow            = []byte("voteswindow")
	ParamStoreKeyMinValidVotesPerWindow = []byte("minvalidvotesperwindow")
	ParamStoreKeySlashFraction          = []byte("slashfraction")
)

// Default parameter values
const (
	DefaultVotePeriod  = core.BlocksPerMinute // 1 minute
	DefaultVotesWindow = int64(1000)          // 1000 oracle period
)

// Default parameter values
var (
	DefaultVoteThreshold          = sdk.NewDecWithPrec(50, 2) // 50%
	DefaultRewardBand             = sdk.NewDecWithPrec(1, 2)  // 1%
	DefaultRewardFraction         = sdk.NewDecWithPrec(1, 2)  // 1%
	DefaultMinValidVotesPerWindow = sdk.NewDecWithPrec(5, 2)  // 5%
	DefaultSlashFraction          = sdk.NewDecWithPrec(1, 4)  // 0.01%
)

var _ subspace.ParamSet = &Params{}

// Params oracle parameters
type Params struct {
	VotePeriod             int64   `json:"vote_period" yaml:"vote_period"`
	VoteThreshold          sdk.Dec `json:"vote_threshold" yaml:"vote_threshold"`
	RewardBand             sdk.Dec `json:"reward_band" yaml:"reward_band"`
	VotesWindow            int64   `json:"votes_window" yaml:"votes_window"`
	MinValidVotesPerWindow sdk.Dec `json:"min_valid_votes_per_window" yaml:"min_valid_votes_per_window"`
	SlashFraction          sdk.Dec `json:"slash_fraction" yaml:"slash_fraction"`
	RewardFraction         sdk.Dec `json:"reward_fraction" yaml:"reward_fraction"`
}

// DefaultParams creates default oracle module parameters
func DefaultParams() Params {
	return Params{
		VotePeriod:             DefaultVotePeriod,
		VoteThreshold:          DefaultVoteThreshold,
		RewardBand:             DefaultRewardBand,
		RewardFraction:         DefaultRewardFraction,
		VotesWindow:            DefaultVotesWindow,
		MinValidVotesPerWindow: DefaultMinValidVotesPerWindow,
		SlashFraction:          DefaultSlashFraction,
	}
}

// validate a set of params
func (params Params) Validate() error {
	if params.VotePeriod <= 0 {
		return fmt.Errorf("oracle parameter VotePeriod must be > 0, is %d", params.VotePeriod)
	}
	if params.VoteThreshold.LTE(sdk.NewDecWithPrec(33, 2)) {
		return fmt.Errorf("oracle parameter VoteTheshold must be greater than 33 percent")
	}
	if params.RewardBand.IsNegative() {
		return fmt.Errorf("oracle parameter RewardBand must be positive")
	}
	if params.RewardFraction.IsNegative() {
		return fmt.Errorf("oracle parameter RewardBand must be positive")
	}
	if params.VotesWindow <= 10 {
		return fmt.Errorf("oracle parameter VotesWindow must be > 0, is %d", params.VotesWindow)
	}
	if params.SlashFraction.GT(sdk.NewDecWithPrec(1, 2)) || params.SlashFraction.IsNegative() {
		return fmt.Errorf("oracle parameter SlashFraction must be smaller or equal than 1 percent and positive")
	}
	if params.MinValidVotesPerWindow.IsNegative() || params.MinValidVotesPerWindow.GT(sdk.OneDec()) {
		return fmt.Errorf("Min valid votes per window should be less than or equal to one and greater than zero, is %s", params.MinValidVotesPerWindow.String())
	}
	return nil
}

// ParamSetPairs implements the ParamSet interface and returns all the key/value pairs
// pairs of oracle module's parameters.
// nolint
func (params *Params) ParamSetPairs() subspace.ParamSetPairs {
	return subspace.ParamSetPairs{
		{Key: ParamStoreKeyVotePeriod, Value: &params.VotePeriod},
		{Key: ParamStoreKeyVoteThreshold, Value: &params.VoteThreshold},
		{Key: ParamStoreKeyRewardBand, Value: &params.RewardBand},
		{Key: ParamStoreKeyRewardFraction, Value: &params.RewardFraction},
		{Key: ParamStoreKeyVotesWindow, Value: &params.VotesWindow},
		{Key: ParamStoreKeyMinValidVotesPerWindow, Value: &params.MinValidVotesPerWindow},
		{Key: ParamStoreKeySlashFraction, Value: &params.SlashFraction},
	}
}

// String implements fmt.Stringer interface
func (params Params) String() string {
	return fmt.Sprintf(`Treasury Params:
  VotePeriod:               %d
  VoteThreshold:            %s
	RewardBand:               %s
	RewardFraction:               %s
	VotesWindow:              %d
	MinValidVotesPerWindow:   %s
	SlashFraction:            %s
	`, params.VotePeriod, params.VoteThreshold, params.RewardBand, params.RewardFraction,
		params.VotesWindow, params.MinValidVotesPerWindow, params.SlashFraction)
}
