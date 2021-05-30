package simulation

import (
	"encoding/binary"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/tendermint/tendermint/crypto/ed25519"
	tmkv "github.com/tendermint/tendermint/libs/kv"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/terra-money/core/x/wasm/internal/types"
)

var (
	creatorPk    = ed25519.GenPrivKey().PubKey()
	contractPk   = ed25519.GenPrivKey().PubKey()
	creatorAddr  = sdk.AccAddress(creatorPk.Address())
	contractAddr = sdk.AccAddress(contractPk.Address())
)

func makeTestCodec() (cdc *codec.Codec) {
	cdc = codec.New()
	sdk.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	types.RegisterCodec(cdc)
	return
}

func TestDecodeDistributionStore(t *testing.T) {
	cdc := makeTestCodec()

	lastCodeIDbz := make([]byte, 8)
	lastInstanceIDbz := make([]byte, 8)
	binary.LittleEndian.PutUint64(lastCodeIDbz, 123)
	binary.LittleEndian.PutUint64(lastInstanceIDbz, 456)

	codeInfo := types.NewCodeInfo(1, []byte{1, 2, 3}, creatorAddr)
	contractInfo := types.NewContractInfo(1, contractAddr, creatorAddr, []byte{4, 5, 6}, true)
	contractStore := []byte{7, 8, 9}

	kvPairs := tmkv.Pairs{
		tmkv.Pair{Key: types.LastCodeIDKey, Value: lastCodeIDbz},
		tmkv.Pair{Key: types.LastInstanceIDKey, Value: lastInstanceIDbz},
		tmkv.Pair{Key: types.CodeKey, Value: cdc.MustMarshalBinaryLengthPrefixed(codeInfo)},
		tmkv.Pair{Key: types.ContractInfoKey, Value: cdc.MustMarshalBinaryLengthPrefixed(contractInfo)},
		tmkv.Pair{Key: types.ContractStoreKey, Value: cdc.MustMarshalBinaryLengthPrefixed(contractStore)},
		tmkv.Pair{Key: []byte{0x99}, Value: []byte{0x99}},
	}

	tests := []struct {
		name        string
		expectedLog string
	}{
		{"LastCodeID", "lastCodeIDA: 123\nlastCodeIDB: 123"},
		{"LastInstanceID", "lastInstanceIDA: 456\nlastInstanceIDB: 456"},
		{"CodeInfo", fmt.Sprintf("%v\n%v", codeInfo, codeInfo)},
		{"ContractInfo", fmt.Sprintf("%v\n%v", contractInfo, contractInfo)},
		{"ContractStore", fmt.Sprintf("%v\n%v", contractStore, contractStore)},
		{"other", ""},
	}

	for i, tt := range tests {
		i, tt := i, tt
		t.Run(tt.name, func(t *testing.T) {
			switch i {
			case len(tests) - 1:
				require.Panics(t, func() { DecodeStore(cdc, kvPairs[i], kvPairs[i]) }, tt.name)
			default:
				require.Equal(t, tt.expectedLog, DecodeStore(cdc, kvPairs[i], kvPairs[i]), tt.name)
			}
		})
	}
}
