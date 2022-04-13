package utils

import (
	"io/ioutil"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/terra-money/core/x/market/types"
)

type (
	// SeigniorageRoutesJSON defines a slice of SeigniorageRouteJSON objects which can be
	// converted to a slice of SeigniorageRoute objects.
	SeigniorageRoutesJSON []SeigniorageRouteJSON

	// SeigniorageRouteJSON defines a parameter change used in JSON input. This
	// allows values to be specified in raw JSON instead of being string encoded.
	SeigniorageRouteJSON struct {
		Address string `json:"address" yaml:"address"`
		Weight  string `json:"weight" yaml:"weight"`
	}

	// SeigniorageRouteChangeProposalJSON defines a ParameterChangeProposal with a deposit used
	// to parse parameter change proposals from a JSON file.
	SeigniorageRouteChangeProposalJSON struct {
		Title       string                `json:"title" yaml:"title"`
		Description string                `json:"description" yaml:"description"`
		Routes      SeigniorageRoutesJSON `json:"routes" yaml:"routes"`
		Deposit     string                `json:"deposit" yaml:"deposit"`
	}

	// SeigniorageRouteChangeProposalReq defines a parameter change proposal request body.
	SeigniorageRouteChangeProposalReq struct {
		BaseReq rest.BaseReq `json:"base_req" yaml:"base_req"`

		Title       string                `json:"title" yaml:"title"`
		Description string                `json:"description" yaml:"description"`
		Routes      SeigniorageRoutesJSON `json:"routes" yaml:"routes"`
		Proposer    sdk.AccAddress        `json:"proposer" yaml:"proposer"`
		Deposit     sdk.Coins             `json:"deposit" yaml:"deposit"`
	}
)

func NewSeigniorageRouteJSON(address, weight string) SeigniorageRouteJSON {
	return SeigniorageRouteJSON{address, weight}
}

// ToSeigniorageRoute converts a SeigniorageRouteJSON object to SeigniorageRoute.
func (pcj SeigniorageRouteJSON) ToSeigniorageRoute() (*types.SeigniorageRoute, error) {
	address, err := sdk.AccAddressFromBech32(pcj.Address)
	if err != nil {
		return nil, err
	}

	weight, err := sdk.NewDecFromStr(pcj.Weight)
	if err != nil {
		return nil, err
	}

	route := types.NewSeigniorageRoute(address, weight)
	return &route, nil
}

// ToSeigniorageRoutes converts a slice of SeigniorageRouteJSON objects to a slice of
// SeigniorageRoute.
func (pcj SeigniorageRoutesJSON) ToSeigniorageRoutes() ([]types.SeigniorageRoute, error) {
	res := make([]types.SeigniorageRoute, len(pcj))
	for i, pc := range pcj {
		route, err := pc.ToSeigniorageRoute()
		if err != nil {
			return nil, err
		}

		res[i] = *route
	}
	return res, nil
}

// ParseSeigniorageRouteChangeProposalJSON reads and parses a SeigniorageRouteChangeProposalJSON from
// file.
func ParseSeigniorageRouteChangeProposalJSON(cdc *codec.LegacyAmino, proposalFile string) (SeigniorageRouteChangeProposalJSON, error) {
	proposal := SeigniorageRouteChangeProposalJSON{}

	contents, err := ioutil.ReadFile(proposalFile)
	if err != nil {
		return proposal, err
	}

	if err := cdc.UnmarshalJSON(contents, &proposal); err != nil {
		return proposal, err
	}

	return proposal, nil
}
