package v05_test

import (
	"encoding/base64"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/terra-money/core/app"
	core "github.com/terra-money/core/types"
	v04wasm "github.com/terra-money/core/x/wasm/legacy/v04"
	v05wasm "github.com/terra-money/core/x/wasm/legacy/v05"
)

func TestMigrate(t *testing.T) {
	sdk.GetConfig().SetBech32PrefixForAccount(core.Bech32PrefixAccAddr, core.Bech32PrefixAccPub)

	encodingConfig := app.MakeEncodingConfig()
	clientCtx := client.Context{}.
		WithInterfaceRegistry(encodingConfig.InterfaceRegistry).
		WithTxConfig(encodingConfig.TxConfig).
		WithLegacyAmino(encodingConfig.Amino).
		WithCodec(encodingConfig.Marshaler)

	addr, err := sdk.AccAddressFromBech32("terra1mx72uukvzqtzhc6gde7shrjqfu5srk22v7gmww")
	require.NoError(t, err)

	contractAddr, err := sdk.AccAddressFromBech32("terra13vs2znvhdcy948ejsh7p8p22j8l4n4y07062qq")
	require.NoError(t, err)

	bz, err := base64.StdEncoding.DecodeString("MTIz")
	require.NoError(t, err)

	bz2, err := base64.StdEncoding.DecodeString("NDU2")
	require.NoError(t, err)

	jsonBz, err := base64.StdEncoding.DecodeString("eyJrZXkiOiAidmFsdWUifQ==")
	require.NoError(t, err)

	wasmGenState := v04wasm.GenesisState{
		Codes: []v04wasm.Code{
			{
				CodeInfo: v04wasm.CodeInfo{
					CodeID:   1,
					CodeHash: bz,
					Creator:  addr,
				},
				CodesBytes: bz,
			},
			{
				CodeInfo: v04wasm.CodeInfo{
					CodeID:   2,
					CodeHash: bz,
					Creator:  addr,
				},
				CodesBytes: bz,
			},
		},
		Contracts: []v04wasm.Contract{
			{
				ContractInfo: v04wasm.ContractInfo{
					Address:    contractAddr,
					Owner:      addr,
					CodeID:     1,
					InitMsg:    jsonBz,
					Migratable: true,
				},
				ContractStore: []v04wasm.Model{
					{
						Key:   bz,
						Value: bz,
					},
					{
						Key:   bz2,
						Value: bz2,
					},
				},
			},
			{
				ContractInfo: v04wasm.ContractInfo{
					Address:    contractAddr,
					Owner:      addr,
					CodeID:     2,
					InitMsg:    jsonBz,
					Migratable: false,
				},
				ContractStore: []v04wasm.Model{
					{
						Key:   bz,
						Value: bz,
					},
					{
						Key:   bz2,
						Value: bz2,
					},
				},
			},
		},
		LastCodeID:     2,
		LastInstanceID: 2,
		Params: v04wasm.Params{
			MaxContractSize:    100,
			MaxContractGas:     10000,
			MaxContractMsgSize: 1024,
		},
	}

	migrated := v05wasm.Migrate(wasmGenState)

	bz, err = clientCtx.Codec.MarshalJSON(migrated)
	require.NoError(t, err)

	// Indent the JSON bz correctly.
	var jsonObj map[string]interface{}
	err = json.Unmarshal(bz, &jsonObj)
	require.NoError(t, err)
	indentedBz, err := json.MarshalIndent(jsonObj, "", "\t")
	require.NoError(t, err)

	// Make sure about:
	// ExchangeRateVotes removed
	// ExchangeRatePrevotes removed
	expected := `{
	"codes": [
		{
			"code_bytes": "",
			"code_info": {
				"code_hash": "",
				"code_id": "1",
				"creator": "terra1mx72uukvzqtzhc6gde7shrjqfu5srk22v7gmww"
			}
		},
		{
			"code_bytes": "",
			"code_info": {
				"code_hash": "",
				"code_id": "2",
				"creator": "terra1mx72uukvzqtzhc6gde7shrjqfu5srk22v7gmww"
			}
		}
	],
	"contracts": [
		{
			"contract_info": {
				"address": "terra13vs2znvhdcy948ejsh7p8p22j8l4n4y07062qq",
				"admin": "terra1mx72uukvzqtzhc6gde7shrjqfu5srk22v7gmww",
				"code_id": "1",
				"creator": "terra1mx72uukvzqtzhc6gde7shrjqfu5srk22v7gmww",
				"init_msg": {
					"key": "value"
				}
			},
			"contract_store": [
				{
					"key": "MTIz",
					"value": "MTIz"
				},
				{
					"key": "NDU2",
					"value": "NDU2"
				}
			]
		},
		{
			"contract_info": {
				"address": "terra13vs2znvhdcy948ejsh7p8p22j8l4n4y07062qq",
				"admin": "",
				"code_id": "2",
				"creator": "terra1mx72uukvzqtzhc6gde7shrjqfu5srk22v7gmww",
				"init_msg": {
					"key": "value"
				}
			},
			"contract_store": [
				{
					"key": "MTIz",
					"value": "MTIz"
				},
				{
					"key": "NDU2",
					"value": "NDU2"
				}
			]
		}
	],
	"last_code_id": "2",
	"last_instance_id": "2",
	"params": {
		"max_contract_gas": "20000000",
		"max_contract_msg_size": "4096",
		"max_contract_size": "614400"
	}
}`
	assert.JSONEq(t, expected, string(indentedBz))
}
