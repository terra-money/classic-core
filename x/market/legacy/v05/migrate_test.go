package v05_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/terra-project/core/app"
	core "github.com/terra-project/core/types"
	v04market "github.com/terra-project/core/x/market/legacy/v04"
	v05market "github.com/terra-project/core/x/market/legacy/v05"
)

func TestMigrate(t *testing.T) {
	sdk.GetConfig().SetBech32PrefixForAccount(core.Bech32PrefixAccAddr, core.Bech32PrefixAccPub)

	encodingConfig := app.MakeEncodingConfig()
	clientCtx := client.Context{}.
		WithInterfaceRegistry(encodingConfig.InterfaceRegistry).
		WithTxConfig(encodingConfig.TxConfig).
		WithLegacyAmino(encodingConfig.Amino).
		WithJSONCodec(encodingConfig.Marshaler)

	marketGenState := v04market.GenesisState{
		TerraPoolDelta: sdk.ZeroDec(),
		Params: v04market.Params{
			BasePool:           sdk.NewDec(1000000),
			PoolRecoveryPeriod: int64(10000),
			MinStabilitySpread: sdk.NewDecWithPrec(2, 2),
		},
	}

	migrated := v05market.Migrate(marketGenState)

	bz, err := clientCtx.JSONCodec.MarshalJSON(migrated)
	require.NoError(t, err)

	// Indent the JSON bz correctly.
	var jsonObj map[string]interface{}
	err = json.Unmarshal(bz, &jsonObj)
	require.NoError(t, err)
	indentedBz, err := json.MarshalIndent(jsonObj, "", "\t")
	require.NoError(t, err)

	// Make sure about:
	// - BasePool to Mint & Burn pool
	expected := `{
	"burn_pool_delta": "0.000000000000000000",
	"mint_pool_delta": "0.000000000000000000",
	"params": {
		"burn_base_pool": "1000000.000000000000000000",
		"min_stability_spread": "0.020000000000000000",
		"mint_base_pool": "10000000.000000000000000000",
		"pool_recovery_period": "10000"
	}
}`

	require.Equal(t, expected, string(indentedBz))
}
