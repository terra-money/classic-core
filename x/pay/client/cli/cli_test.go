package cli

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/terra-project/core/testutil"
	"github.com/terra-project/core/types/assets"
)

func TestSendTx(t *testing.T) {
	cdc, rootCmd, txCmd, _ := testutil.PrepareCmdTest()

	txCmd.AddCommand(SendTxCmd(cdc))

	// normal case all parameter given
	_, err := testutil.ExecuteCommand(
		rootCmd,
		`tx`,
		`send`,
		`--from=terra1wg2mlrxdmnnkkykgqg4znky86nyrtc45q336yv`,
		`--coins=1000uluna`,
		`--to=terra1wg2mlrxdmnnkkykgqg4znky86nyrtc45q336yv`,
		`--generate-only`,
		`--offline`,
		`--chain-id=columbus`,
	)

	require.Nil(t, err)
}

func TestSetManualFees(t *testing.T) {
	coins := sdk.NewCoins(sdk.NewCoin(assets.MicroLunaDenom, sdk.NewInt(1000)), sdk.NewCoin(assets.MicroKRWDenom, sdk.NewInt(10000)))
	_, err := sdk.ParseCoins(coins.String())
	require.NoError(t, err)
}
