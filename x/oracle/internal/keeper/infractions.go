package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/terra-project/core/x/oracle/internal/types"
)

// HandleWrongVotes handles a wrong votes, must be called once per validator per voting period.
func (k Keeper) HandleBallotAttendees(ctx sdk.Context, ballotAttendees map[string]bool) {
	for addr, valid := range ballotAttendees {
		valAddr, err := sdk.ValAddressFromBech32(addr)
		if err != nil {
			panic(err) // NOTE never occurs
		}

		k.handleBallotAttendee(ctx, valAddr, valid)
	}
}

func (k Keeper) handleBallotAttendee(ctx sdk.Context, valAddr sdk.ValAddress, valid bool) {
	logger := k.Logger(ctx)
	height := ctx.BlockHeight()

	// fetch voting info
	votingInfo, found := k.getVotingInfo(ctx, valAddr)
	if !found {
		panic(fmt.Sprintf("Expected signing info for validator %s but not found", valAddr))
	}

	// this is a relative index, so it counts blocks the validator *should* have signed
	// will use the 0-value default signing info if not present, except for start height
	index := votingInfo.IndexOffset % k.VotesWindow(ctx)
	votingInfo.IndexOffset++

	// Update signed block bit array & counter
	// This counter just tracks the sum of the bit array
	// That way we avoid needing to read/write the whole array each time
	previous := k.GetMissedVoteBitArray(ctx, valAddr, index)
	missed := !valid
	switch {
	case !previous && missed:
		// Array value has changed from not missed to missed, increment counter
		k.SetMissedVoteBitArray(ctx, valAddr, index, true)
		votingInfo.MissedVotesCounter++
	case previous && !missed:
		// Array value has changed from missed to not missed, decrement counter
		k.SetMissedVoteBitArray(ctx, valAddr, index, false)
		votingInfo.MissedVotesCounter--
	default:
		// Array value at this index has not changed, no need to update counter
	}

	if missed {
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeLiveness,
				sdk.NewAttribute(types.AttributeKeyAddress, valAddr.String()),
				sdk.NewAttribute(types.AttributeKeyMissedVotes, fmt.Sprintf("%d", votingInfo.MissedVotesCounter)),
				sdk.NewAttribute(types.AttributeKeyHeight, fmt.Sprintf("%d", height)),
			),
		)

		logger.Info(
			fmt.Sprintf("Wrong voting validator %s at height %d, %d missed, threshold %d", valAddr, height, votingInfo.MissedVotesCounter, k.MinValidVotesPerWindow(ctx)))
	}

	minHeight := votingInfo.StartHeight + k.VotesWindow(ctx)
	maxMissed := k.VotesWindow(ctx) - k.MinValidVotesPerWindow(ctx)

	// if we are past the minimum height and the validator has missed too many blocks, punish them
	if height > minHeight && votingInfo.MissedVotesCounter > maxMissed {
		validator := k.StakingKeeper.Validator(ctx, valAddr)
		if validator != nil && !validator.IsJailed() {
			power := validator.GetConsensusPower()

			// Downtime confirmed: slash and jail the validator
			logger.Info(fmt.Sprintf("Validator %s past min height of %d and below valid blocks threshold of %d",
				valAddr, minHeight, k.MinValidVotesPerWindow(ctx)))

			// We need to retrieve the stake distribution which signed the block, so we subtract ValidatorUpdateDelay from the evidence height,
			// and subtract an additional 1 since this is the LastCommit.
			// Note that this *can* result in a negative "distributionHeight" up to -ValidatorUpdateDelay-1,
			// i.e. at the end of the pre-genesis block (none) = at the beginning of the genesis block.
			// That's fine since this is just used to filter unbonding delegations & redelegations.
			distributionHeight := height - sdk.ValidatorUpdateDelay - 1

			ctx.EventManager().EmitEvent(
				sdk.NewEvent(
					types.EventTypeSlash,
					sdk.NewAttribute(types.AttributeKeyAddress, valAddr.String()),
					sdk.NewAttribute(types.AttributeKeyPower, fmt.Sprintf("%d", power)),
				),
			)

			k.StakingKeeper.Slash(ctx, validator.GetConsAddr(), distributionHeight, power, k.SlashFraction(ctx))

			// We need to reset the counter & array so that the validator won't be immediately slashed for downtime upon rebonding.
			votingInfo.MissedVotesCounter = 0
			votingInfo.IndexOffset = 0
			k.clearMissedVoteBitArray(ctx, valAddr)
		} else {
			// Validator was (a) not found or (b) already jailed, don't slash
			logger.Info(
				fmt.Sprintf("Validator %s would have been slashed for invalid oracle votes, but was either not found in store or already jailed", valAddr),
			)
		}
	}

	// Set the updated signing info
	k.SetVotingInfo(ctx, valAddr, votingInfo)
}
