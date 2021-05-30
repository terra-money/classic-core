package rest

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/multisig"
	"github.com/tendermint/tendermint/crypto/secp256k1"

	"github.com/terra-money/core/app"
)

func TestMultisigPubkey(t *testing.T) {
	cdc := app.MakeCodec()

	pubkey, err := cdc.MarshalBinaryBare(multisig.NewPubKeyMultisigThreshold(3,
		[]crypto.PubKey{
			secp256k1.GenPrivKey().PubKey(),
			secp256k1.GenPrivKey().PubKey(),
			secp256k1.GenPrivKey().PubKey(),
			secp256k1.GenPrivKey().PubKey(),
		}))

	require.NoError(t, err)

	fmt.Println(hex.EncodeToString(pubkey))

	pubkey, err = cdc.MarshalBinaryBare(multisig.NewPubKeyMultisigThreshold(2,
		[]crypto.PubKey{
			secp256k1.GenPrivKey().PubKey(),
			secp256k1.GenPrivKey().PubKey(),
		}))

	require.NoError(t, err)

	fmt.Println(hex.EncodeToString(pubkey))
}
