package types

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestSeigniorageRouteChangeProposal(t *testing.T) {
	pc1 := NewSeigniorageRoute([]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19}, sdk.NewDecFromIntWithPrec(sdk.NewInt(10), 6))
	pc2 := NewSeigniorageRoute([]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 20}, sdk.NewDecFromIntWithPrec(sdk.NewInt(300), 6))
	pcp := NewSeigniorageRouteChangeProposal("test title", "test description", []SeigniorageRoute{pc1, pc2})

	require.Equal(t, "test title", pcp.GetTitle())
	require.Equal(t, "test description", pcp.GetDescription())
	require.Equal(t, RouterKey, pcp.ProposalRoute())
	require.Equal(t, ProposalTypeChange, pcp.ProposalType())
	require.Nil(t, pcp.ValidateBasic())

	pc3 := NewSeigniorageRoute(sdk.AccAddress{}, sdk.NewDecFromIntWithPrec(sdk.NewInt(10), 6))
	pcp = NewSeigniorageRouteChangeProposal("test title", "test description", []SeigniorageRoute{pc3})
	require.Error(t, pcp.ValidateBasic())

	pc4 := NewSeigniorageRoute([]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19}, sdk.ZeroDec())
	pcp = NewSeigniorageRouteChangeProposal("test title", "test description", []SeigniorageRoute{pc4})
	require.Error(t, pcp.ValidateBasic())

	pc5 := NewSeigniorageRoute([]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19}, sdk.OneDec())
	pcp = NewSeigniorageRouteChangeProposal("test title", "test description", []SeigniorageRoute{pc5})
	require.Error(t, pcp.ValidateBasic())

	pc6 := NewSeigniorageRoute([]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19}, sdk.OneDec())
	pc7 := NewSeigniorageRoute([]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19}, sdk.OneDec())
	pcp = NewSeigniorageRouteChangeProposal("test title", "test description", []SeigniorageRoute{pc6, pc7})
	require.Error(t, pcp.ValidateBasic())
}
