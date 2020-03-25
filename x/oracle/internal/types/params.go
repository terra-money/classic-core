package types

import (
	"fmt"

	core "github.com/terra-project/core/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params/subspace"
)

// DefaultParamspace defines default space for oracle params
const DefaultParamspace = ModuleName

// Parameter keys
var (
	ParamStoreKeyVotePeriod               = []byte("voteperiod")
	ParamStoreKeyVoteThreshold            = []byte("votethreshold")
	ParamStoreKeyRewardBand               = []byte("rewardband")
	ParamStoreKeyRewardDistributionWindow = []byte("rewarddistributionwindow")
	ParamStoreKeyWhitelist                = []byte("whitelist")
	ParamStoreKeySlashFraction            = []byte("slashfraction")
	ParamStoreKeySlashWindow              = []byte("slashwindow")
	ParamStoreKeyMinValidPerWindow        = []byte("minvalidperwindow")
)

// Default parameter values
const (
	DefaultVotePeriod               = core.BlocksPerMinute / 2 // 30 seconds
	DefaultSlashWindow              = core.BlocksPerWeek       // window for a week
	DefaultRewardDistributionWindow = core.BlocksPerYear       // window for a year
)

// Default parameter values
var (
	DefaultVoteThreshold = sdk.NewDecWithPrec(50, 2) // 50%
	DefaultRewardBand    = sdk.NewDecWithPrec(2, 2)  // 2% (-1, 1)
	DefaultTobinTax      = sdk.NewDecWithPrec(25, 4) // 0.25%
	DefaultWhitelist     = DenomList{
		{Name: core.MicroKRWDenom, TobinTax: DefaultTobinTax},
		{Name: core.MicroSDRDenom, TobinTax: DefaultTobinTax},
		{Name: core.MicroUSDDenom, TobinTax: DefaultTobinTax},
		{Name: core.MicroMNTDenom, TobinTax: DefaultTobinTax.MulInt64(8)}}
	DefaultSlashFraction     = sdk.NewDecWithPrec(1, 4) // 0.01%
	DefaultMinValidPerWindow = sdk.NewDecWithPrec(5, 2) // 5%
)

var _ subspace.ParamSet = &Params{}

// Params oracle parameters
type Params struct {
	VotePeriod               int64     `json:"vote_period" yaml:"vote_period"`                               // the number of blocks during which voting takes place.
	VoteThreshold            sdk.Dec   `json:"vote_threshold" yaml:"vote_threshold"`                         // the minimum percentage of votes that must be received for a ballot to pass.
	RewardBand               sdk.Dec   `json:"reward_band" yaml:"reward_band"`                               // the ratio of allowable exchange rate error that can be rewarded.
	RewardDistributionWindow int64     `json:"reward_distribution_window" yaml:"reward_distribution_window"` // the number of blocks during which seigniorage reward comes in and then is distributed.
	Whitelist                DenomList `json:"whitelist" yaml:"whitelist"`                                   // the denom list that can be activated,
	SlashFraction            sdk.Dec   `json:"slash_fraction" yaml:"slash_fraction"`                         // the ratio of penalty on bonded tokens
	SlashWindow              int64     `json:"slash_window" yaml:"slash_window"`                             // the number of blocks for slashing tallying
	MinValidPerWindow        sdk.Dec   `json:"min_valid_per_window" yaml:"min_valid_per_window"`             // the ratio of minimum valid oracle votes per slash window to avoid slashing
}

// DefaultParams creates default oracle module parameters
func DefaultParams() Params {
	return Params{
		VotePeriod:               DefaultVotePeriod,
		VoteThreshold:            DefaultVoteThreshold,
		RewardBand:               DefaultRewardBand,
		RewardDistributionWindow: DefaultRewardDistributionWindow,
		Whitelist:                DefaultWhitelist,
		SlashFraction:            DefaultSlashFraction,
		SlashWindow:              DefaultSlashWindow,
		MinValidPerWindow:        DefaultMinValidPerWindow,
	}
}

// Validate validates a set of params
func (params Params) Validate() error {
	if params.VotePeriod <= 0 {
		return fmt.Errorf("oracle parameter VotePeriod must be > 0, is %d", params.VotePeriod)
	}
	if params.VoteThreshold.LTE(sdk.NewDecWithPrec(33, 2)) {
		return fmt.Errorf("oracle parameter VoteTheshold must be greater than 33 percent")
	}
	if params.RewardBand.IsNegative() || params.RewardBand.GT(sdk.OneDec()) {
		return fmt.Errorf("oracle parameter RewardBand must be between [0, 1]")
	}
	if params.RewardDistributionWindow < params.VotePeriod {
		return fmt.Errorf("oracle parameter RewardDistributionWindow must be greater than or equal with votes period")
	}
	if params.SlashFraction.GT(sdk.OneDec()) || params.SlashFraction.IsNegative() {
		return fmt.Errorf("oracle parameter SlashRraction must be between [0, 1]")
	}
	if params.SlashWindow < params.VotePeriod {
		return fmt.Errorf("oracle parameter SlashWindow must be greater than or equal with votes period")
	}
	if params.MinValidPerWindow.GT(sdk.NewDecWithPrec(5, 1)) || params.MinValidPerWindow.IsNegative() {
		return fmt.Errorf("oracle parameter MinValidPerWindow must be between [0, 0.5]")
	}
	for _, denom := range params.Whitelist {
		if denom.TobinTax.LT(sdk.ZeroDec()) || denom.TobinTax.GT(sdk.OneDec()) {
			return fmt.Errorf("oracle parameter Whitelist Denom must have TobinTax between [0, 1]")
		}
		if len(denom.Name) == 0 {
			return fmt.Errorf("oracle parameter Whitelist Denom must have name")
		}
	}
	return nil
}

// ParamSetPairs implements the ParamSet interface and returns all the key/value pairs
// pairs of oracle module's parameters.
func (params *Params) ParamSetPairs() subspace.ParamSetPairs {
	return subspace.ParamSetPairs{
		{Key: ParamStoreKeyVotePeriod, Value: &params.VotePeriod},
		{Key: ParamStoreKeyVoteThreshold, Value: &params.VoteThreshold},
		{Key: ParamStoreKeyRewardBand, Value: &params.RewardBand},
		{Key: ParamStoreKeyRewardDistributionWindow, Value: &params.RewardDistributionWindow},
		{Key: ParamStoreKeyWhitelist, Value: &params.Whitelist},
		{Key: ParamStoreKeySlashFraction, Value: &params.SlashFraction},
		{Key: ParamStoreKeySlashWindow, Value: &params.SlashWindow},
		{Key: ParamStoreKeyMinValidPerWindow, Value: &params.MinValidPerWindow},
	}
}

// String implements fmt.Stringer interface
func (params Params) String() string {
	return fmt.Sprintf(`Oracle Params:
	VotePeriod:                  %d
	VoteThreshold:               %s
	RewardBand:                  %s
	RewardDistributionWindow:    %d
	Whitelist                    %s
	SlashFraction                %s
	SlashWindow                  %d
	MinValidPerWindow            %s
	`, params.VotePeriod, params.VoteThreshold, params.RewardBand,
		params.RewardDistributionWindow, params.Whitelist,
		params.SlashFraction, params.SlashWindow, params.MinValidPerWindow)
}
