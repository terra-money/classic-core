package v040

import (
	"testing"

	"github.com/stretchr/testify/require"
	v039authcustom "github.com/terra-money/core/custom/auth/legacy/v039"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/crypto/types/multisig"
	sdk "github.com/cosmos/cosmos-sdk/types"
	v038auth "github.com/cosmos/cosmos-sdk/x/auth/legacy/v038"
)

func TestMultisigPubkeyMigration(t *testing.T) {
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount("terra", "terrapub")

	v04Codec := codec.NewLegacyAmino()
	v039authcustom.RegisterLegacyAminoCodec(v04Codec)
	jsonAccount := `{"type":"core/LazyGradedVestingAccount","value":{"address":"terra1dp0taj85ruc299rkdvzp4z5pfg6z6swaed74e6","coins":[],"public_key":{"type":"tendermint/PubKeyMultisigThreshold","value":{"threshold":"2","pubkeys":[{"type":"tendermint/PubKeySecp256k1","value":"AyETa9Y9ihObzeRPWMP0MBAa0Mqune3I+5KonOCPTtkv"},{"type":"tendermint/PubKeySecp256k1","value":"AzzLltyI4MzxLpcmS1vfpXJeAk/sgS1eVYmvXgFpGRtg"},{"type":"tendermint/PubKeySecp256k1","value":"AnZjvWmye3JPEL95xRcGeFRf4o8pHDK0dkZjf6B9D4FA"}]}},"account_number":"317776","sequence":"61","original_vesting":[{"denom":"usdr","amount":"1000000000000000"}],"delegated_free":[{"denom":"uluna","amount":"57620630000008"}],"delegated_vesting":[],"end_time":"0","vesting_schedules":[{"denom":"usdr","schedules":[{"start_time":"1556085600","end_time":"1558677600","ratio":"0.100000000000000000"},{"start_time":"1587708000","end_time":"1590300000","ratio":"0.100000000000000000"},{"start_time":"1619244000","end_time":"1621836000","ratio":"0.100000000000000000"},{"start_time":"1650780000","end_time":"1653372000","ratio":"0.100000000000000000"},{"start_time":"1682316000","end_time":"1684908000","ratio":"0.100000000000000000"},{"start_time":"1713938400","end_time":"1716530400","ratio":"0.100000000000000000"},{"start_time":"1745474400","end_time":"1748066400","ratio":"0.100000000000000000"},{"start_time":"1777010400","end_time":"1779602400","ratio":"0.100000000000000000"},{"start_time":"1808546400","end_time":"1811138400","ratio":"0.100000000000000000"},{"start_time":"1840168800","end_time":"1842760800","ratio":"0.100000000000000000"}]}]}}`

	var oldAccount v038auth.GenesisAccount
	v04Codec.MustUnmarshalJSON([]byte(jsonAccount), &oldAccount)

	account := convertBaseVestingAccount(oldAccount.(*v039authcustom.LazyGradedVestingAccount).BaseVestingAccount)

	pk, ok := account.GetPubKey().(multisig.PubKey)
	require.True(t, ok)
	require.Equal(t, 3, len(pk.GetPubKeys()))
}
