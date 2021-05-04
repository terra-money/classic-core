// DONTCOVER
// nolint
package v04

import (
	"fmt"
	"strings"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	v036gov "github.com/cosmos/cosmos-sdk/x/gov/legacy/v036"
)

const (
	// ModuleName nonlint
	ModuleName = "treasury"

	// RouterKey is the message route for distribution
	RouterKey = ModuleName

	// ProposalTypeTaxRateUpdate defines the type for a TaxRateUpdateProposal
	ProposalTypeTaxRateUpdate = "TaxRateUpdate"

	// ProposalTypeRewardWeightUpdate defines the type for a RewardWeightUpdateProposal
	ProposalTypeRewardWeightUpdate = "RewardWeightUpdate"

	// Block units
	BlocksPerMinute = uint64(10)
	BlocksPerHour   = BlocksPerMinute * 60
	BlocksPerDay    = BlocksPerHour * 24
	BlocksPerWeek   = BlocksPerDay * 7
)

type (
	// GenesisState - all treasury state that must be provided at genesis
	GenesisState struct {
		Params               Params             `json:"params" yaml:"params"` // market params
		TaxRate              sdk.Dec            `json:"tax_rate" yaml:"tax_rate"`
		RewardWeight         sdk.Dec            `json:"reward_weight" yaml:"reward_weight"`
		TaxCaps              map[string]sdk.Int `json:"tax_caps" yaml:"tax_caps"`
		TaxProceed           sdk.Coins          `json:"tax_proceed" yaml:"tax_proceed"`
		EpochInitialIssuance sdk.Coins          `json:"epoch_initial_issuance" yaml:"epoch_initial_issuance"`
		CumulativeHeight     int64              `json:"cumulated_height" yaml:"cumulated_height"`
		TRs                  []sdk.Dec          `json:"TRs" yaml:"TRs"`
		SRs                  []sdk.Dec          `json:"SRs" yaml:"SRs"`
		TSLs                 []sdk.Int          `json:"TSLs" yaml:"TSLs"`
	}

	// Params treasury parameters
	Params struct {
		TaxPolicy               PolicyConstraints `json:"tax_policy" yaml:"tax_policy"`
		RewardPolicy            PolicyConstraints `json:"reward_policy" yaml:"reward_policy"`
		SeigniorageBurdenTarget sdk.Dec           `json:"seigniorage_burden_target" yaml:"seigniorage_burden_target"`
		MiningIncrement         sdk.Dec           `json:"mining_increment" yaml:"mining_increment"`
		WindowShort             int64             `json:"window_short" yaml:"window_short"`
		WindowLong              int64             `json:"window_long" yaml:"window_long"`
		WindowProbation         int64             `json:"window_probation" yaml:"window_probation"`
	}

	// PolicyConstraints wraps constraints around updating a key Treasury variable
	PolicyConstraints struct {
		RateMin       sdk.Dec  `json:"rate_min"`
		RateMax       sdk.Dec  `json:"rate_max"`
		Cap           sdk.Coin `json:"cap"`
		ChangeRateMax sdk.Dec  `json:"change_max"`
	}

	// TaxRateUpdateProposal updates treasury tax-rate
	TaxRateUpdateProposal struct {
		Title       string  `json:"title" yaml:"title"`             // Title of the Proposal
		Description string  `json:"description" yaml:"description"` // Description of the Proposal
		TaxRate     sdk.Dec `json:"tax_rate" yaml:"tax_rate"`       // target TaxRate
	}

	// RewardWeightUpdateProposal update treasury tax-rate
	RewardWeightUpdateProposal struct {
		Title        string  `json:"title" yaml:"title"`                 // Title of the Proposal
		Description  string  `json:"description" yaml:"description"`     // Description of the Proposal
		RewardWeight sdk.Dec `json:"reward_weight" yaml:"reward_weight"` // target RewardWeight
	}
)

var _ v036gov.Content = TaxRateUpdateProposal{}
var _ v036gov.Content = RewardWeightUpdateProposal{}

// GetTitle returns the title of an TaxRateUpdateProposal.
func (p TaxRateUpdateProposal) GetTitle() string { return p.Title }

// GetDescription returns the description of an TaxRateUpdateProposal.
func (p TaxRateUpdateProposal) GetDescription() string { return p.Description }

// ProposalRoute returns the routing key of an TaxRateUpdateProposal.
func (TaxRateUpdateProposal) ProposalRoute() string { return RouterKey }

// ProposalType returns the type of an TaxRateUpdateProposal.
func (p TaxRateUpdateProposal) ProposalType() string { return ProposalTypeTaxRateUpdate }

// ValidateBasic runs basic stateless validity checks
func (p TaxRateUpdateProposal) ValidateBasic() error {
	err := v036gov.ValidateAbstract(p)
	if err != nil {
		return err
	}

	if p.TaxRate.IsNegative() || p.TaxRate.GT(sdk.OneDec()) {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "Invalid tax-rate: "+p.TaxRate.String())
	}

	return nil
}

// String implements the Stringer interface.
func (p TaxRateUpdateProposal) String() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf(`Community Pool Spend Proposal:
  Title:        %s
  Description:  %s
  TaxRate:      %s
`, p.Title, p.Description, p.TaxRate))
	return b.String()
}

// GetTitle returns the title of an RewardWeightUpdateProposal.
func (p RewardWeightUpdateProposal) GetTitle() string { return p.Title }

// GetDescription returns the description of an RewardWeightUpdateProposal.
func (p RewardWeightUpdateProposal) GetDescription() string { return p.Description }

// ProposalRoute returns the routing key of an RewardWeightUpdateProposal.
func (RewardWeightUpdateProposal) ProposalRoute() string { return RouterKey }

// ProposalType returns the type of an RewardWeightUpdateProposal.
func (p RewardWeightUpdateProposal) ProposalType() string { return ProposalTypeRewardWeightUpdate }

// ValidateBasic runs basic stateless validity checks
func (p RewardWeightUpdateProposal) ValidateBasic() error {
	err := v036gov.ValidateAbstract(p)
	if err != nil {
		return err
	}

	if p.RewardWeight.IsNegative() || p.RewardWeight.GT(sdk.OneDec()) {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "Invalid reward-weight: "+p.RewardWeight.String())
	}

	return nil
}

// String implements the Stringer interface.
func (p RewardWeightUpdateProposal) String() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf(`Community Pool Spend Proposal:
  Title:        %s
  Description:  %s
  RewardWeight: %s
`, p.Title, p.Description, p.RewardWeight))
	return b.String()
}

// RegisterLegacyAminoCodec nolint
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(TaxRateUpdateProposal{}, "treasury/TaxRateUpdateProposal", nil)
	cdc.RegisterConcrete(RewardWeightUpdateProposal{}, "treasury/RewardWeightUpdateProposal", nil)
}
