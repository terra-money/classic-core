package oracle

import (
	"fmt"
	"terra/types/assets"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Params oracle parameters
type Params struct {
	Whitelist     []string `json:"whitelist"`      // type of coin to mint
	VotePeriod    sdk.Int  `json:"vote_period"`    // maximum annual change in inflation rate
	VoteThreshold sdk.Dec  `json:"vote_threshold"` // maximum inflation rate
}

// NewParams creates a new param instance
func NewParams(whitelist []string, votePeriod sdk.Int, voteThreshold sdk.Dec) Params {
	return Params{
		Whitelist:     whitelist,
		VotePeriod:    votePeriod,
		VoteThreshold: voteThreshold,
	}
}

// DefaultParams creates default oracle module parameters
func DefaultParams() Params {
	return NewParams(
		assets.GetAllDenoms(),
		sdk.NewInt(1000000),
		sdk.NewDecWithPrec(66, 2), // 66%
	)
}

func validateParams(params Params) error {
	if len(params.Whitelist) == 0 {
		return fmt.Errorf("oracle parameter whitelist should not be nil")
	}
	if params.VotePeriod.LT(sdk.ZeroInt()) {
		return fmt.Errorf("oracle parameter VotePeriod must be > 0, is %s", params.VotePeriod.String())
	}
	if params.VoteThreshold.LT(sdk.NewDecWithPrec(33, 2)) {
		return fmt.Errorf("oracle parameter VoteThreshold must be greater than 33%")
	}
	return nil
}
