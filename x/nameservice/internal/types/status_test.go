package types

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestAuctionStatus(t *testing.T) {
	items := []string{"Bid", "Reveal"}

	for _, item := range items {
		auctionStatus, err := AuctionStatusFromString(item)
		require.NoError(t, err)

		bz, err := auctionStatus.Marshal()
		require.NoError(t, err)

		status := AuctionStatus(0xff)
		err = status.Unmarshal(bz)
		require.NoError(t, err)

		require.Equal(t, auctionStatus, status)

		bz2, err := auctionStatus.MarshalJSON()
		require.NoError(t, err)

		status2 := AuctionStatus(0xff)
		err = status2.UnmarshalJSON(bz2)
		require.NoError(t, err)

		require.Equal(t, auctionStatus, status2)
		require.Equal(t, item, fmt.Sprintf("%s", auctionStatus))
	}

	_, err := AuctionStatusFromString("INVALID")
	require.Error(t, err)
}
