package types

import (
	"fmt"

	"gopkg.in/yaml.v2"

	core "github.com/terra-money/core/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
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

var _ params.ParamSet = &Params{}

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

// ParamSetPairs returns the parameter set pairs.
func (p *Params) ParamSetPairs() params.ParamSetPairs {
	return params.ParamSetPairs{
		params.NewParamSetPair(ParamStoreKeyVotePeriod, &p.VotePeriod, validateVotePeriod),
		params.NewParamSetPair(ParamStoreKeyVoteThreshold, &p.VoteThreshold, validateVoteThreshold),
		params.NewParamSetPair(ParamStoreKeyRewardBand, &p.RewardBand, validateRewardBand),
		params.NewParamSetPair(ParamStoreKeyRewardDistributionWindow, &p.RewardDistributionWindow, validateRewardDistributionWindow),
		params.NewParamSetPair(ParamStoreKeyWhitelist, &p.Whitelist, validateWhitelist),
		params.NewParamSetPair(ParamStoreKeySlashFraction, &p.SlashFraction, validateSlashFraction),
		params.NewParamSetPair(ParamStoreKeySlashWindow, &p.SlashWindow, validateSlashWindow),
		params.NewParamSetPair(ParamStoreKeyMinValidPerWindow, &p.MinValidPerWindow, validateMinValidPerWindow),
	}
}

// ParamKeyTable returns the parameter key table.
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable().RegisterParamSet(&Params{})
}

// String implements fmt.Stringer interface
func (p Params) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}

// ValidateBasic performs basic validation on oracle parameters.
func (p Params) ValidateBasic() error {
	if p.VotePeriod <= 0 {
		return fmt.Errorf("oracle parameter VotePeriod must be > 0, is %d", p.VotePeriod)
	}
	if p.VoteThreshold.LTE(sdk.NewDecWithPrec(33, 2)) {
		return fmt.Errorf("oracle parameter VoteTheshold must be greater than 33 percent")
	}

	if p.RewardBand.IsNegative() || p.RewardBand.GT(sdk.OneDec()) {
		return fmt.Errorf("oracle parameter RewardBand must be between [0, 1]")
	}

	if p.RewardDistributionWindow < p.VotePeriod {
		return fmt.Errorf("oracle parameter RewardDistributionWindow must be greater than or equal with votes period")
	}

	if p.SlashFraction.GT(sdk.OneDec()) || p.SlashFraction.IsNegative() {
		return fmt.Errorf("oracle parameter SlashRraction must be between [0, 1]")
	}

	if p.SlashWindow < p.VotePeriod {
		return fmt.Errorf("oracle parameter SlashWindow must be greater than or equal with votes period")
	}

	if p.MinValidPerWindow.GT(sdk.NewDecWithPrec(5, 1)) || p.MinValidPerWindow.IsNegative() {
		return fmt.Errorf("oracle parameter MinValidPerWindow must be between [0, 0.5]")
	}

	for _, denom := range p.Whitelist {
		if denom.TobinTax.LT(sdk.ZeroDec()) || denom.TobinTax.GT(sdk.OneDec()) {
			return fmt.Errorf("oracle parameter Whitelist Denom must have TobinTax between [0, 1]")
		}
		if len(denom.Name) == 0 {
			return fmt.Errorf("oracle parameter Whitelist Denom must have name")
		}
	}
	return nil
}

func validateVotePeriod(i interface{}) error {
	v, ok := i.(int64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v <= 0 {
		return fmt.Errorf("vote period must be positive: %d", v)
	}

	return nil
}

func validateVoteThreshold(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.LT(sdk.NewDecWithPrec(33, 2)) {
		return fmt.Errorf("vote threshold must be bigger than 33%%: %s", v)
	}

	if v.GT(sdk.OneDec()) {
		return fmt.Errorf("vote threshold too large: %s", v)
	}

	return nil
}

func validateRewardBand(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNegative() {
		return fmt.Errorf("reward band must be positive: %s", v)
	}

	if v.GT(sdk.OneDec()) {
		return fmt.Errorf("reward band is too large: %s", v)
	}

	return nil
}

func validateRewardDistributionWindow(i interface{}) error {
	v, ok := i.(int64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v <= 0 {
		return fmt.Errorf("reward distribution window must be positive: %d", v)
	}

	return nil
}

func validateWhitelist(i interface{}) error {
	v, ok := i.(DenomList)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	for _, d := range v {
		if d.TobinTax.LT(sdk.ZeroDec()) || d.TobinTax.GT(sdk.OneDec()) {
			return fmt.Errorf("oracle parameter Whitelist Denom must have TobinTax between [0, 1]")
		}
		if len(d.Name) == 0 {
			return fmt.Errorf("oracle parameter Whitelist Denom must have name")
		}
	}

	return nil
}

func validateSlashFraction(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNegative() {
		return fmt.Errorf("slash fraction must be positive: %s", v)
	}

	if v.GT(sdk.OneDec()) {
		return fmt.Errorf("slash fraction is too large: %s", v)
	}

	return nil
}

func validateSlashWindow(i interface{}) error {
	v, ok := i.(int64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v <= 0 {
		return fmt.Errorf("slash window must be positive: %d", v)
	}

	return nil
}

func validateMinValidPerWindow(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNegative() {
		return fmt.Errorf("min valid per window must be positive: %s", v)
	}

	if v.GT(sdk.NewDecWithPrec(5, 1)) {
		return fmt.Errorf("min valid per window is too large: %s", v)
	}

	return nil
}
