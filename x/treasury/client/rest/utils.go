package rest

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
)

type (
	// TaxRateUpdateProposalReq defines a tax-rate-update proposal request body.
	TaxRateUpdateProposalReq struct {
		BaseReq rest.BaseReq `json:"base_req" yaml:"base_req"`

		Title       string         `json:"title" yaml:"title"`
		Description string         `json:"description" yaml:"description"`
		TaxRate     sdk.Dec        `json:"tax_rate" yaml:"tax_rate"`
		Proposer    sdk.AccAddress `json:"proposer" yaml:"proposer"`
		Deposit     sdk.Coins      `json:"deposit" yaml:"deposit"`
	}

	// RewardWeightUpdateProposalReq defines a tax-rate-update proposal request body.
	RewardWeightUpdateProposalReq struct {
		BaseReq rest.BaseReq `json:"base_req" yaml:"base_req"`

		Title        string         `json:"title" yaml:"title"`
		Description  string         `json:"description" yaml:"description"`
		RewardWeight sdk.Dec        `json:"reward_weight" yaml:"reward_weight"`
		Proposer     sdk.AccAddress `json:"proposer" yaml:"proposer"`
		Deposit      sdk.Coins      `json:"deposit" yaml:"deposit"`
	}
)
