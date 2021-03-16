package cli

import (
	"io/ioutil"

	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/terra-project/core/x/treasury/types"
)

// ParseTaxRateUpdateProposalWithDeposit reads and parses a TaxRateUpdateProposalJSON from a file.
func ParseTaxRateUpdateProposalWithDeposit(cdc codec.JSONMarshaler, proposalFile string) (types.TaxRateUpdateProposalWithDeposit, error) {
	proposal := types.TaxRateUpdateProposalWithDeposit{}

	contents, err := ioutil.ReadFile(proposalFile)
	if err != nil {
		return proposal, err
	}

	if err := cdc.UnmarshalJSON(contents, &proposal); err != nil {
		return proposal, err
	}

	return proposal, nil
}

// ParseRewardWeightUpdateProposalWithDeposit reads and parses a RewardWeightUpdateProposalJSON from a file.
func ParseRewardWeightUpdateProposalWithDeposit(cdc codec.JSONMarshaler, proposalFile string) (types.RewardWeightUpdateProposalWithDeposit, error) {
	proposal := types.RewardWeightUpdateProposalWithDeposit{}

	contents, err := ioutil.ReadFile(proposalFile)
	if err != nil {
		return proposal, err
	}

	if err := cdc.UnmarshalJSON(contents, &proposal); err != nil {
		return proposal, err
	}

	return proposal, nil
}
