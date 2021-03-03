package cli

import (
	"io/ioutil"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type (

	// TaxRateUpdateProposalJSON defines a TaxRateUpdateProposal with a deposit
	TaxRateUpdateProposalJSON struct {
		Title       string    `json:"title" yaml:"title"`
		Description string    `json:"description" yaml:"description"`
		TaxRate     sdk.Dec   `json:"tax_rate" yaml:"tax_rate"`
		Deposit     sdk.Coins `json:"deposit" yaml:"deposit"`
	}

	// RewardWeightUpdateProposalJSON defines a RewardWeightUpdateProposal with a deposit
	RewardWeightUpdateProposalJSON struct {
		Title        string    `json:"title" yaml:"title"`
		Description  string    `json:"description" yaml:"description"`
		RewardWeight sdk.Dec   `json:"reward_weight" yaml:"reward_weight"`
		Deposit      sdk.Coins `json:"deposit" yaml:"deposit"`
	}
)

// ParseTaxRateUpdateProposalJSON reads and parses a TaxRateUpdateProposalJSON from a file.
func ParseTaxRateUpdateProposalJSON(cdc *codec.Codec, proposalFile string) (TaxRateUpdateProposalJSON, error) {
	proposal := TaxRateUpdateProposalJSON{}

	contents, err := ioutil.ReadFile(proposalFile)
	if err != nil {
		return proposal, err
	}

	if err := cdc.UnmarshalJSON(contents, &proposal); err != nil {
		return proposal, err
	}

	return proposal, nil
}

// ParseRewardWeightUpdateProposalJSON reads and parses a RewardWeightUpdateProposalJSON from a file.
func ParseRewardWeightUpdateProposalJSON(cdc *codec.Codec, proposalFile string) (RewardWeightUpdateProposalJSON, error) {
	proposal := RewardWeightUpdateProposalJSON{}

	contents, err := ioutil.ReadFile(proposalFile)
	if err != nil {
		return proposal, err
	}

	if err := cdc.UnmarshalJSON(contents, &proposal); err != nil {
		return proposal, err
	}

	return proposal, nil
}
