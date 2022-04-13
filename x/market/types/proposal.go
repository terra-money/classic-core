package types

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	customgovtypes "github.com/terra-money/core/custom/gov/types"
)

const (
	// ProposalTypeChange defines the type for a SeigniorageRouteChangeProposal
	ProposalTypeChange = "SeigniorageRouteChange"
)

// Assert SeigniorageRouteChangeProposal implements govtypes.Content at compile-time
var _ govtypes.Content = &SeigniorageRouteChangeProposal{}

func init() {
	govtypes.RegisterProposalType(ProposalTypeChange)
	customgovtypes.RegisterProposalTypeCodec(&SeigniorageRouteChangeProposal{}, "market/SeigniorageRouteChangeProposal")
}

func NewSeigniorageRouteChangeProposal(title, description string, changes []SeigniorageRoute) *SeigniorageRouteChangeProposal {
	return &SeigniorageRouteChangeProposal{title, description, changes}
}

// GetTitle returns the title of a parameter change proposal.
func (pcp *SeigniorageRouteChangeProposal) GetTitle() string { return pcp.Title }

// GetDescription returns the description of a parameter change proposal.
func (pcp *SeigniorageRouteChangeProposal) GetDescription() string { return pcp.Description }

// ProposalRoute returns the routing key of a parameter change proposal.
func (pcp *SeigniorageRouteChangeProposal) ProposalRoute() string { return RouterKey }

// ProposalType returns the type of a parameter change proposal.
func (pcp *SeigniorageRouteChangeProposal) ProposalType() string { return ProposalTypeChange }

// ValidateBasic validates the parameter change proposal
func (pcp *SeigniorageRouteChangeProposal) ValidateBasic() error {
	err := govtypes.ValidateAbstract(pcp)
	if err != nil {
		return err
	}

	return ValidateChanges(pcp.Routes)
}

// String implements the Stringer interface.
func (pcp SeigniorageRouteChangeProposal) String() string {
	var b strings.Builder

	b.WriteString(fmt.Sprintf(`Seigniorage Route Change Proposal:
  Title:       %s
  Description: %s
  Changes:
`, pcp.Title, pcp.Description))

	for _, pc := range pcp.Routes {
		b.WriteString(fmt.Sprintf(`    Seigniorage Route:
      Address: %s
      Weight:      %s
`, pc.Address, pc.Weight))
	}

	return b.String()
}

// ValidateChanges performs basic validation checks over a set of SeigniorageRoute. It
// returns an error if any SeigniorageRoute is invalid.
func ValidateChanges(changes []SeigniorageRoute) error {
	if len(changes) == 0 {
		return ErrEmptyChanges
	}

	weightSum := sdk.ZeroDec()
	addrMap := map[string]bool{}
	for _, pc := range changes {
		if len(pc.Address) == 0 {
			return ErrEmptyAddress
		}

		// each weight must be bigger than zero
		if pc.Weight.IsZero() {
			return ErrZeroWeight
		}

		// check duplicated address
		if _, exists := addrMap[pc.Address]; exists {
			return ErrDuplicateRoute
		}

		weightSum = weightSum.Add(pc.Weight)
		addrMap[pc.Address] = true
	}

	// the sum of weights must be smaller than one
	if weightSum.GTE(sdk.OneDec()) {
		return ErrInvalidWeightSum
	}

	return nil
}
