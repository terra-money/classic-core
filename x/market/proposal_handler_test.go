package market

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/terra-money/core/x/market/keeper"
	"github.com/terra-money/core/x/market/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func testProposal(routes ...types.SeigniorageRoute) *types.SeigniorageRouteChangeProposal {
	return types.NewSeigniorageRouteChangeProposal("title", "description", routes)
}

func TestProposalHandler(t *testing.T) {
	input, _ := setup(t)

	validAddr := keeper.Addrs[0]

	validWeight := sdk.NewDecWithPrec(1, 3)

	testCases := []struct {
		name     string
		proposal *types.SeigniorageRouteChangeProposal
		onHandle func()
		expErr   bool
	}{
		{
			"all fields",
			testProposal(types.NewSeigniorageRoute(validAddr, validWeight)),
			func() {

				routes := input.MarketKeeper.GetSeigniorageRoutes(input.Ctx)
				require.Equal(t, routes, []types.SeigniorageRoute{types.NewSeigniorageRoute(validAddr, validWeight)})
			},
			false,
		},
	}

	govHandler := NewSeigniorageRouteChangeProposalHandler(input.MarketKeeper)
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := govHandler(input.Ctx, tc.proposal)
			if tc.expErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				tc.onHandle()
			}
		})
	}
}
