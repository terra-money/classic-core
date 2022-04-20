package types

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestSeigniorageRoutes_ValidateRoutes(t *testing.T) {
	addr1 := sdk.AccAddress([]byte("addr1_______________")).String()
	addr2 := sdk.AccAddress([]byte("addr2_______________")).String()
	testCases := []struct {
		name   string
		routes SeigniorageRoutes
		expErr error
	}{
		{
			"empty address",
			SeigniorageRoutes{
				Routes: []SeigniorageRoute{{
					Address: "",
					Weight:  sdk.NewDecWithPrec(1, 1),
				},
				},
			},
			ErrInvalidAddress,
		},
		{
			"invalid address",
			SeigniorageRoutes{
				Routes: []SeigniorageRoute{{
					Address: "osmosis1sznj93ytuxwxh27ufk0amx547p3jr374c63zzm",
					Weight:  sdk.NewDecWithPrec(1, 1),
				},
				}},
			ErrInvalidAddress,
		},
		{
			"duplicate addresses",
			SeigniorageRoutes{
				Routes: []SeigniorageRoute{
					{
						Address: addr1,
						Weight:  sdk.NewDecWithPrec(1, 1),
					},
					{
						Address: addr1,
						Weight:  sdk.NewDecWithPrec(2, 1),
					},
				}},
			ErrDuplicateRoute,
		},
		{
			"negative weight",
			SeigniorageRoutes{
				Routes: []SeigniorageRoute{
					{
						Address: addr1,
						Weight:  sdk.NewDecWithPrec(-1, 1),
					},
				}},
			ErrInvalidWeight,
		},
		{
			"zero weight",
			SeigniorageRoutes{
				Routes: []SeigniorageRoute{
					{
						Address: addr1,
						Weight:  sdk.ZeroDec(),
					},
				}},
			ErrInvalidWeight,
		},
		{
			"one weight",
			SeigniorageRoutes{
				Routes: []SeigniorageRoute{
					{
						Address: addr1,
						Weight:  sdk.OneDec(),
					},
				}},
			ErrInvalidWeightsSum,
		},
		{
			"weight sum exceeding one",
			SeigniorageRoutes{
				Routes: []SeigniorageRoute{
					{
						Address: addr1,
						Weight:  sdk.NewDecWithPrec(5, 1),
					},
					{
						Address: addr2,
						Weight:  sdk.NewDecWithPrec(5, 1),
					},
				}},
			ErrInvalidWeightsSum,
		},
		{
			"valid routes",
			SeigniorageRoutes{
				Routes: []SeigniorageRoute{
					{
						Address: addr1,
						Weight:  sdk.NewDecWithPrec(1, 1),
					},
					{
						Address: addr2,
						Weight:  sdk.NewDecWithPrec(2, 1),
					},
				}},
			nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.routes.ValidateRoutes()
			if tc.expErr != nil {
				require.ErrorIs(t, err, tc.expErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
