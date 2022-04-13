package client

import (
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
	"github.com/terra-money/core/x/market/client/cli"
	"github.com/terra-money/core/x/market/client/rest"
)

// ProposalHandler is the param change proposal handler.
var ProposalHandler = govclient.NewProposalHandler(cli.NewSubmitSeigniorageRouteChangeTxCmd, rest.ProposalRESTHandler)
