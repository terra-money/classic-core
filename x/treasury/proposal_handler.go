package treasury

import (
	"github.com/classic-terra/core/v2/x/treasury/keeper"
	"github.com/classic-terra/core/v2/x/treasury/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

func NewProposalHandler(k keeper.Keeper) govtypes.Handler {
	return func(ctx sdk.Context, content govtypes.Content) error {
		switch c := content.(type) {
		case *types.AddBurnTaxExemptionAddressProposal:
			return handleAddBurnTaxExemptionAddressProposal(ctx, k, c)
		case *types.RemoveBurnTaxExemptionAddressProposal:
			return handleRemoveBurnTaxExemptionAddressProposal(ctx, k, c)
		default:
			return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized treasury proposal content type: %T", c)
		}
	}
}

func handleAddBurnTaxExemptionAddressProposal(ctx sdk.Context, k keeper.Keeper, p *types.AddBurnTaxExemptionAddressProposal) error {
	return keeper.HandleAddBurnTaxExemptionAddressProposal(ctx, k, p)
}

func handleRemoveBurnTaxExemptionAddressProposal(ctx sdk.Context, k keeper.Keeper, p *types.RemoveBurnTaxExemptionAddressProposal) error {
	return keeper.HandleRemoveBurnTaxExemptionAddressProposal(ctx, k, p)
}
