package v043_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/terra-money/core/app"
	v043authz "github.com/terra-money/core/custom/authz/legacy/v043"
	core "github.com/terra-money/core/types"
	v04msgauth "github.com/terra-money/core/x/msgauth/legacy/v04"
)

func TestMigrate(t *testing.T) {
	sdk.GetConfig().SetBech32PrefixForAccount(core.Bech32PrefixAccAddr, core.Bech32PrefixAccPub)

	encodingConfig := app.MakeEncodingConfig()
	clientCtx := client.Context{}.
		WithInterfaceRegistry(encodingConfig.InterfaceRegistry).
		WithTxConfig(encodingConfig.TxConfig).
		WithLegacyAmino(encodingConfig.Amino).
		WithCodec(encodingConfig.Marshaler)

	granter, err := sdk.AccAddressFromBech32("terra13vs2znvhdcy948ejsh7p8p22j8l4n4y07062qq")
	require.NoError(t, err)

	grantee, err := sdk.AccAddressFromBech32("terra1mx72uukvzqtzhc6gde7shrjqfu5srk22v7gmww")
	require.NoError(t, err)

	msgauthGenState := v04msgauth.GenesisState{
		AuthorizationEntries: []v04msgauth.AuthorizationEntry{
			{
				Granter: granter,
				Grantee: grantee,
				Authorization: v04msgauth.GenericAuthorization{
					GrantMsgType: "vote",
				},
			},
			{
				Granter: granter,
				Grantee: grantee,
				Authorization: v04msgauth.SendAuthorization{
					SpendLimit: sdk.Coins{
						{
							Denom:  core.MicroUSDDenom,
							Amount: sdk.NewInt(100),
						},
					},
				},
			},
		},
	}

	migrated := v043authz.Migrate(msgauthGenState)

	bz, err := clientCtx.Codec.MarshalJSON(migrated)
	require.NoError(t, err)

	// Indent the JSON bz correctly.
	var jsonObj map[string]interface{}
	err = json.Unmarshal(bz, &jsonObj)
	require.NoError(t, err)
	indentedBz, err := json.MarshalIndent(jsonObj, "", "\t")
	require.NoError(t, err)

	// Make sure about:
	// - GenericAuthorization has correct JSON.
	// - SendAuthorization has correct JSON.
	expected := `{
	"authorization": [
		{
			"authorization": {
				"@type": "/cosmos.authz.v1beta1.GenericAuthorization",
				"msg": "/cosmos.gov.v1beta1.MsgVote"
			},
			"expiration": "0001-01-01T00:00:00Z",
			"grantee": "terra1mx72uukvzqtzhc6gde7shrjqfu5srk22v7gmww",
			"granter": "terra13vs2znvhdcy948ejsh7p8p22j8l4n4y07062qq"
		},
		{
			"authorization": {
				"@type": "/cosmos.bank.v1beta1.SendAuthorization",
				"spend_limit": [
					{
						"amount": "100",
						"denom": "uusd"
					}
				]
			},
			"expiration": "0001-01-01T00:00:00Z",
			"grantee": "terra1mx72uukvzqtzhc6gde7shrjqfu5srk22v7gmww",
			"granter": "terra13vs2znvhdcy948ejsh7p8p22j8l4n4y07062qq"
		}
	]
}`

	require.Equal(t, expected, string(indentedBz))
}
