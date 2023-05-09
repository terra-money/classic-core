package keeper

import (
	"github.com/classic-terra/core/v2/x/treasury/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func HandleAddBurnTaxExemptionAddressProposal(ctx sdk.Context, k Keeper, p *types.AddBurnTaxExemptionAddressProposal) error {
	for _, address := range p.Addresses {
		k.AddBurnTaxExemptionAddress(ctx, address)
	}

	return nil
}

func HandleRemoveBurnTaxExemptionAddressProposal(ctx sdk.Context, k Keeper, p *types.RemoveBurnTaxExemptionAddressProposal) error {
	for _, address := range p.Addresses {
		err := k.RemoveBurnTaxExemptionAddress(ctx, address)
		if err != nil {
			return err
		}
	}

	return nil
}
