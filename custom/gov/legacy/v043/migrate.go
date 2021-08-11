package v043

import (
	"fmt"

	proto "github.com/gogo/protobuf/proto"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	v036distr "github.com/cosmos/cosmos-sdk/x/distribution/legacy/v036"
	v040distr "github.com/cosmos/cosmos-sdk/x/distribution/types"
	v034gov "github.com/cosmos/cosmos-sdk/x/gov/legacy/v034"
	v036gov "github.com/cosmos/cosmos-sdk/x/gov/legacy/v036"
	v043gov "github.com/cosmos/cosmos-sdk/x/gov/types"
	v036params "github.com/cosmos/cosmos-sdk/x/params/legacy/v036"
	v040params "github.com/cosmos/cosmos-sdk/x/params/types/proposal"
	v038upgrade "github.com/cosmos/cosmos-sdk/x/upgrade/legacy/v038"
	v040upgrade "github.com/cosmos/cosmos-sdk/x/upgrade/types"

	v04treasury "github.com/terra-money/core/x/treasury/legacy/v04"
)

func migrateVoteOption(oldVoteOption v034gov.VoteOption) v043gov.VoteOption {
	switch oldVoteOption {
	case v034gov.OptionEmpty:
		return v043gov.OptionEmpty

	case v034gov.OptionYes:
		return v043gov.OptionYes

	case v034gov.OptionAbstain:
		return v043gov.OptionAbstain

	case v034gov.OptionNo:
		return v043gov.OptionNo

	case v034gov.OptionNoWithVeto:
		return v043gov.OptionNoWithVeto

	default:
		panic(fmt.Errorf("'%s' is not a valid vote option", oldVoteOption))
	}
}

func migrateProposalStatus(oldProposalStatus v034gov.ProposalStatus) v043gov.ProposalStatus {
	switch oldProposalStatus {

	case v034gov.StatusNil:
		return v043gov.StatusNil

	case v034gov.StatusDepositPeriod:
		return v043gov.StatusDepositPeriod

	case v034gov.StatusVotingPeriod:
		return v043gov.StatusVotingPeriod

	case v034gov.StatusPassed:
		return v043gov.StatusPassed

	case v034gov.StatusRejected:
		return v043gov.StatusRejected

	case v034gov.StatusFailed:
		return v043gov.StatusFailed

	default:
		panic(fmt.Errorf("'%s' is not a valid proposal status", oldProposalStatus))
	}
}

func migrateContent(oldContent v036gov.Content) *codectypes.Any {
	var protoProposal proto.Message

	switch oldContent := oldContent.(type) {
	case v036gov.TextProposal:
		{
			protoProposal = &v043gov.TextProposal{
				Title:       oldContent.Title,
				Description: oldContent.Description,
			}
			// Convert the content into Any.
			contentAny, err := codectypes.NewAnyWithValue(protoProposal)
			if err != nil {
				panic(err)
			}

			return contentAny
		}
	case v036distr.CommunityPoolSpendProposal:
		{
			protoProposal = &v040distr.CommunityPoolSpendProposal{
				Title:       oldContent.Title,
				Description: oldContent.Description,
				Recipient:   oldContent.Recipient.String(),
				Amount:      oldContent.Amount,
			}
		}
	case v038upgrade.CancelSoftwareUpgradeProposal:
		{
			protoProposal = &v040upgrade.CancelSoftwareUpgradeProposal{
				Description: oldContent.Description,
				Title:       oldContent.Title,
			}
		}
	case v038upgrade.SoftwareUpgradeProposal:
		{
			protoProposal = &v040upgrade.SoftwareUpgradeProposal{
				Description: oldContent.Description,
				Title:       oldContent.Title,
				Plan: v040upgrade.Plan{
					Name:   oldContent.Plan.Name,
					Height: oldContent.Plan.Height,
					Info:   oldContent.Plan.Info,
				},
			}
		}
	case v036params.ParameterChangeProposal:
		{
			newChanges := make([]v040params.ParamChange, len(oldContent.Changes))
			for i, oldChange := range oldContent.Changes {
				newChanges[i] = v040params.ParamChange{
					Subspace: oldChange.Subspace,
					Key:      oldChange.Key,
					Value:    oldChange.Value,
				}
			}

			protoProposal = &v040params.ParameterChangeProposal{
				Description: oldContent.Description,
				Title:       oldContent.Title,
				Changes:     newChanges,
			}
		}
	case v04treasury.TaxRateUpdateProposal:
		{
			// Change the legacy proposal to text proposal to keep the record
			protoProposal = &v043gov.TextProposal{
				Title:       oldContent.Title,
				Description: oldContent.Description,
			}

			// Convert the content into Any.
			contentAny, err := codectypes.NewAnyWithValue(protoProposal)
			if err != nil {
				panic(err)
			}

			return contentAny
		}
	case v04treasury.RewardWeightUpdateProposal:
		{
			// Change the legacy proposal to text proposal to keep the record
			protoProposal = &v043gov.TextProposal{
				Title:       oldContent.Title,
				Description: oldContent.Description,
			}

			// Convert the content into Any.
			contentAny, err := codectypes.NewAnyWithValue(protoProposal)
			if err != nil {
				panic(err)
			}

			return contentAny
		}
	default:
		panic(fmt.Errorf("%T is not a valid proposal content type", oldContent))
	}

	// Convert the content into Any.
	contentAny, err := codectypes.NewAnyWithValue(protoProposal)
	if err != nil {
		panic(err)
	}

	return contentAny
}

// Migrate accepts exported v0.36 x/gov genesis state and migrates it to
// v0.40 x/gov genesis state. The migration includes:
//
// - Convert vote option & proposal status from byte to enum.
// - Migrate proposal content to Any.
// - Convert addresses from bytes to bech32 strings.
// - Re-encode in v0.40 GenesisState.
func Migrate(oldGovState v036gov.GenesisState) *v043gov.GenesisState {
	newDeposits := make([]v043gov.Deposit, len(oldGovState.Deposits))
	for i, oldDeposit := range oldGovState.Deposits {
		newDeposits[i] = v043gov.Deposit{
			ProposalId: oldDeposit.ProposalID,
			Depositor:  oldDeposit.Depositor.String(),
			Amount:     oldDeposit.Amount,
		}
	}

	newVotes := make([]v043gov.Vote, len(oldGovState.Votes))
	for i, oldVote := range oldGovState.Votes {
		newVotes[i] = v043gov.Vote{
			ProposalId: oldVote.ProposalID,
			Voter:      oldVote.Voter.String(),
			Options:    []v043gov.WeightedVoteOption{{Option: migrateVoteOption(oldVote.Option), Weight: sdk.NewDec(1)}},
		}
	}

	newProposals := make([]v043gov.Proposal, len(oldGovState.Proposals))
	for i, oldProposal := range oldGovState.Proposals {
		newProposals[i] = v043gov.Proposal{
			ProposalId: oldProposal.ProposalID,
			Content:    migrateContent(oldProposal.Content),
			Status:     migrateProposalStatus(oldProposal.Status),
			FinalTallyResult: v043gov.TallyResult{
				Yes:        oldProposal.FinalTallyResult.Yes,
				Abstain:    oldProposal.FinalTallyResult.Abstain,
				No:         oldProposal.FinalTallyResult.No,
				NoWithVeto: oldProposal.FinalTallyResult.NoWithVeto,
			},
			SubmitTime:      oldProposal.SubmitTime,
			DepositEndTime:  oldProposal.DepositEndTime,
			TotalDeposit:    oldProposal.TotalDeposit,
			VotingStartTime: oldProposal.VotingStartTime,
			VotingEndTime:   oldProposal.VotingEndTime,
		}
	}

	return &v043gov.GenesisState{
		StartingProposalId: oldGovState.StartingProposalID,
		Deposits:           newDeposits,
		Votes:              newVotes,
		Proposals:          newProposals,
		DepositParams: v043gov.DepositParams{
			MinDeposit:       oldGovState.DepositParams.MinDeposit,
			MaxDepositPeriod: oldGovState.DepositParams.MaxDepositPeriod,
		},
		VotingParams: v043gov.VotingParams{
			VotingPeriod: oldGovState.VotingParams.VotingPeriod,
		},
		TallyParams: v043gov.TallyParams{
			Quorum:        oldGovState.TallyParams.Quorum,
			Threshold:     oldGovState.TallyParams.Threshold,
			VetoThreshold: oldGovState.TallyParams.Veto,
		},
	}
}
