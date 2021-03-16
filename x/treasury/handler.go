package treasury

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/terra-project/core/x/treasury/keeper"
	"github.com/terra-project/core/x/treasury/types"
)

// NewTreasuryPolicyUpdateHandler custom gov proposal handler
func NewTreasuryPolicyUpdateHandler(k keeper.Keeper) govtypes.Handler {
	return func(ctx sdk.Context, content govtypes.Content) error {
		switch c := content.(type) {
		case *types.TaxRateUpdateProposal:
			return keeper.HandleTaxRateUpdateProposal(ctx, k, c)
		case *types.RewardWeightUpdateProposal:
			return keeper.HandleRewardWeightUpdateProposal(ctx, k, c)

		default:
			return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized treasury proposal content type: %T", c)
		}
	}
}
