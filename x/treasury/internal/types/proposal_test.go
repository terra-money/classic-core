package types

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestTaxRateUpdateProposal(t *testing.T) {
	// invalid title
	proposal := NewTaxRateUpdateProposal("", "description", sdk.NewDec(1))
	require.Error(t, proposal.ValidateBasic())

	// invalid descrription
	proposal = NewTaxRateUpdateProposal("title", "", sdk.NewDec(1))
	require.Error(t, proposal.ValidateBasic())

	// invalid tax-rate
	proposal = NewTaxRateUpdateProposal("title", "description", sdk.NewDec(2))
	require.Error(t, proposal.ValidateBasic())

	proposal = NewTaxRateUpdateProposal("title", "description", sdk.NewDec(-1))
	require.Error(t, proposal.ValidateBasic())

	proposal = NewTaxRateUpdateProposal("title", "description", sdk.NewDecWithPrec(1, 1))
	require.NoError(t, proposal.ValidateBasic())
}

func TestRewardWeightUpdateProposal(t *testing.T) {
	// invalid title
	proposal := NewRewardWeightUpdateProposal("", "description", sdk.NewDec(1))
	require.Error(t, proposal.ValidateBasic())

	// invalid descrription
	proposal = NewRewardWeightUpdateProposal("title", "", sdk.NewDec(1))
	require.Error(t, proposal.ValidateBasic())

	// invalid reward-weight
	proposal = NewRewardWeightUpdateProposal("title", "description", sdk.NewDec(2))
	require.Error(t, proposal.ValidateBasic())

	proposal = NewRewardWeightUpdateProposal("title", "description", sdk.NewDec(-1))
	require.Error(t, proposal.ValidateBasic())

	proposal = NewRewardWeightUpdateProposal("title", "description", sdk.NewDecWithPrec(1, 1))
	require.NoError(t, proposal.ValidateBasic())
}
