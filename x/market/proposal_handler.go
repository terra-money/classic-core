package market

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/terra-money/core/x/market/keeper"
	"github.com/terra-money/core/x/market/types"
)

// NewSeigniorageRouteChangeProposalHandler creates a new governance Handler for a SeigniorageRouteChangeProposal
func NewSeigniorageRouteChangeProposalHandler(k keeper.Keeper) govtypes.Handler {
	return func(ctx sdk.Context, content govtypes.Content) error {
		switch c := content.(type) {
		case *types.SeigniorageRouteChangeProposal:
			return handleSeigniorageRouteChangeProposal(ctx, k, c)

		default:
			return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized param proposal content type: %T", c)
		}
	}
}

func handleSeigniorageRouteChangeProposal(ctx sdk.Context, k keeper.Keeper, p *types.SeigniorageRouteChangeProposal) error {
	k.SetSeigniorageRoutes(ctx, p.Routes)
	return nil
}
