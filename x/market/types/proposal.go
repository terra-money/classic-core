package types

import (
	"fmt"
	"strings"

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

	return SeigniorageRoutes{Routes: pcp.Routes}.ValidateRoutes()
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
