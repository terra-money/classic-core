package types

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/mock"
	"github.com/stretchr/testify/require"
)

func TestClaimPoolSort(t *testing.T) {
	_, addrs, _, _ := mock.CreateGenAccounts(2, sdk.Coins{})

	claim1 := NewClaim(1, sdk.ValAddress(addrs[0]))
	claim2 := NewClaim(2, sdk.ValAddress(addrs[0]))
	claim3 := NewClaim(3, sdk.ValAddress(addrs[1]))

	claimPool := ClaimPool{claim1, claim2, claim3}
	claimPool = claimPool.Sort()

	require.Equal(t, 2, len(claimPool))
	require.Equal(t, int64(3), claimPool[0].Weight)
}
