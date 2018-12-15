package oracle

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GenesisState - all distribution state that must be provided at genesis
type GenesisState struct {
	VoteThreshold sdk.Dec    `json:"vote_threshold"`
	VotePeriod    sdk.Int    `json:"vote_period"`
	GenesisElects PriceVotes `json:"genesis_elects"`
}

// NewGenesisState generates a new oracle genesis state
func NewGenesisState(voteThreshold sdk.Dec, votePeriod sdk.Int, genElects PriceVotes) GenesisState {
	return GenesisState{
		VoteThreshold: voteThreshold,
		VotePeriod:    votePeriod,
		GenesisElects: genElects,
	}
}

// DefaultGenesisState get raw genesis raw message for testing
func DefaultGenesisState() GenesisState {

	return NewGenesisState(
		sdk.NewDecWithPrec(66, 2), // 66%
		sdk.NewInt(1000000),       // TODO: calibrate paramter
		PriceVotes{
			PriceVote{
				FeedMsg: PriceFeedMsg{
					Denom:        "sdr",
					TargetPrice:  sdk.NewDecWithPrec(1, 0),
					CurrentPrice: sdk.NewDecWithPrec(1, 0),
					Feeder:       nil,
				},
				Power: sdk.ZeroDec(),
			},
		},
	)
}

// Init store state from genesis data
func InitGenesis(ctx sdk.Context, keeper Keeper, data GenesisState) {
	keeper.SetThreshold(ctx, data.VoteThreshold)
	keeper.SetVotePeriod(ctx, data.VotePeriod)

	genesisWhitelist := Whitelist{}
	for _, vote := range data.GenesisElects {
		genesisWhitelist = append(genesisWhitelist, vote.FeedMsg.Denom)

		keeper.SetElect(ctx, vote)
	}
	keeper.SetWhitelist(ctx, genesisWhitelist)
}

// ExportGenesis returns a GenesisState for a given context and keeper
func ExportGenesis(ctx sdk.Context, keeper Keeper) GenesisState {
	threshold := keeper.GetThreshold(ctx)
	voteperiod := keeper.GetVotePeriod(ctx)

	electList := PriceVotes{}
	whitelist := keeper.GetWhitelist(ctx)
	for _, wDenom := range whitelist {
		elect := keeper.GetElect(ctx, wDenom)
		electList = append(electList, elect)
	}

	return NewGenesisState(threshold, voteperiod, electList)
}
