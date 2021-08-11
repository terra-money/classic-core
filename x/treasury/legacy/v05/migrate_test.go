package v05_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/terra-money/core/app"
	core "github.com/terra-money/core/types"
	v04treasury "github.com/terra-money/core/x/treasury/legacy/v04"
	v05treasury "github.com/terra-money/core/x/treasury/legacy/v05"
)

func TestMigrate(t *testing.T) {
	sdk.GetConfig().SetBech32PrefixForAccount(core.Bech32PrefixAccAddr, core.Bech32PrefixAccPub)

	encodingConfig := app.MakeEncodingConfig()
	clientCtx := client.Context{}.
		WithInterfaceRegistry(encodingConfig.InterfaceRegistry).
		WithTxConfig(encodingConfig.TxConfig).
		WithLegacyAmino(encodingConfig.Amino).
		WithCodec(encodingConfig.Marshaler)

	treasuryGenState := v04treasury.GenesisState{
		TaxRate:      sdk.NewDecWithPrec(2, 2),
		RewardWeight: sdk.NewDecWithPrec(5, 2),
		TaxCaps: map[string]sdk.Int{
			core.MicroLunaDenom: sdk.NewInt(1),
			core.MicroSDRDenom:  sdk.NewInt(100),
		},
		TaxProceed: sdk.NewCoins(
			sdk.NewCoin(core.MicroLunaDenom, sdk.NewInt(100)),
			sdk.NewCoin(core.MicroSDRDenom, sdk.NewInt(1000)),
			sdk.NewCoin(core.MicroUSDDenom, sdk.NewInt(5000)),
		),
		EpochInitialIssuance: sdk.NewCoins(
			sdk.NewCoin(core.MicroLunaDenom, sdk.NewInt(100)),
			sdk.NewCoin(core.MicroSDRDenom, sdk.NewInt(1000)),
			sdk.NewCoin(core.MicroUSDDenom, sdk.NewInt(5000)),
		),
		CumulativeHeight: 100,
		TRs: []sdk.Dec{
			sdk.NewDec(100),
			sdk.NewDec(200),
			sdk.NewDec(300),
		},
		SRs: []sdk.Dec{
			sdk.NewDec(100),
			sdk.NewDec(200),
			sdk.NewDec(300),
		},
		TSLs: []sdk.Int{
			sdk.NewInt(100),
			sdk.NewInt(200),
			sdk.NewInt(300),
		},
		Params: v04treasury.Params{
			TaxPolicy: v04treasury.PolicyConstraints{
				RateMin:       sdk.NewDecWithPrec(1, 2),
				RateMax:       sdk.NewDecWithPrec(10, 2),
				Cap:           sdk.NewCoin(core.MicroSDRDenom, sdk.NewInt(1000000)),
				ChangeRateMax: sdk.NewDecWithPrec(5, 3),
			},
			RewardPolicy: v04treasury.PolicyConstraints{
				RateMin:       sdk.NewDecWithPrec(1, 2),
				RateMax:       sdk.NewDecWithPrec(10, 2),
				Cap:           sdk.NewCoin("unused", sdk.ZeroInt()),
				ChangeRateMax: sdk.NewDecWithPrec(5, 3),
			},
			SeigniorageBurdenTarget: sdk.NewDecWithPrec(67, 2),
			MiningIncrement:         sdk.NewDecWithPrec(107, 2),
			WindowShort:             4,
			WindowLong:              52,
			WindowProbation:         18,
		},
	}

	migrated := v05treasury.Migrate(treasuryGenState)

	bz, err := clientCtx.Codec.MarshalJSON(migrated)
	require.NoError(t, err)

	// Indent the JSON bz correctly.
	var jsonObj map[string]interface{}
	err = json.Unmarshal(bz, &jsonObj)
	require.NoError(t, err)
	indentedBz, err := json.MarshalIndent(jsonObj, "", "\t")
	require.NoError(t, err)

	// Make sure about:
	// - EpochState has correct JSON.
	expected := `{
	"epoch_initial_issuance": [
		{
			"amount": "100",
			"denom": "uluna"
		},
		{
			"amount": "1000",
			"denom": "usdr"
		},
		{
			"amount": "5000",
			"denom": "uusd"
		}
	],
	"epoch_states": [
		{
			"epoch": "0",
			"seigniorage_reward": "100.000000000000000000",
			"tax_reward": "100.000000000000000000",
			"total_staked_luna": "100"
		},
		{
			"epoch": "1",
			"seigniorage_reward": "200.000000000000000000",
			"tax_reward": "200.000000000000000000",
			"total_staked_luna": "200"
		},
		{
			"epoch": "2",
			"seigniorage_reward": "300.000000000000000000",
			"tax_reward": "300.000000000000000000",
			"total_staked_luna": "300"
		}
	],
	"params": {
		"mining_increment": "1.070000000000000000",
		"reward_policy": {
			"cap": {
				"amount": "0",
				"denom": "unused"
			},
			"change_rate_max": "0.000000000000000000",
			"rate_max": "1.000000000000000000",
			"rate_min": "0.000000000000000000"
		},
		"seigniorage_burden_target": "0.670000000000000000",
		"tax_policy": {
			"cap": {
				"amount": "1000000",
				"denom": "usdr"
			},
			"change_rate_max": "0.005000000000000000",
			"rate_max": "0.100000000000000000",
			"rate_min": "0.010000000000000000"
		},
		"window_long": "52",
		"window_probation": "18",
		"window_short": "4"
	},
	"reward_weight": "1.000000000000000000",
	"tax_caps": [
		{
			"denom": "uluna",
			"tax_cap": "1"
		},
		{
			"denom": "usdr",
			"tax_cap": "100"
		}
	],
	"tax_proceeds": [
		{
			"amount": "100",
			"denom": "uluna"
		},
		{
			"amount": "1000",
			"denom": "usdr"
		},
		{
			"amount": "5000",
			"denom": "uusd"
		}
	],
	"tax_rate": "0.020000000000000000"
}`

	assert.JSONEq(t, expected, string(indentedBz))
}
