package oracle

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// InitGenesis initialize default parameters
// and the keeper's address to pubkey map
func InitGenesis(ctx sdk.Context, keeper Keeper, data GenesisState) {

	for addr, info := range data.VotingInfos {
		address, err := sdk.ValAddressFromBech32(addr)
		if err != nil {
			panic(err)
		}
		keeper.SetVotingInfo(ctx, address, info)
	}

	for addr, array := range data.MissedVotes {
		address, err := sdk.ValAddressFromBech32(addr)
		if err != nil {
			panic(err)
		}
		for _, missed := range array {
			keeper.SetMissedVoteBitArray(ctx, address, missed.Index, missed.Missed)
		}
	}

	keeper.SetParams(ctx, data.Params)
}

// ExportGenesis writes the current store values
// to a genesis file, which can be imported again
// with InitGenesis
func ExportGenesis(ctx sdk.Context, keeper Keeper) (data GenesisState) {
	params := keeper.GetParams(ctx)
	votingInfos := make(map[string]VotingInfo)
	missedVotes := make(map[string][]MissedVote)
	keeper.IterateVotingInfos(ctx, func(info VotingInfo) (stop bool) {
		bechAddr := info.Address.String()

		votingInfos[bechAddr] = info
		localMissedBlocks := []MissedVote{}

		keeper.IterateMissedVoteBitArray(ctx, info.Address, func(index int64, missed bool) (stop bool) {
			localMissedBlocks = append(localMissedBlocks, NewMissedVote(index, missed))
			return false
		})
		missedVotes[bechAddr] = localMissedBlocks

		return false
	})

	return NewGenesisState(params, votingInfos, missedVotes)
}
