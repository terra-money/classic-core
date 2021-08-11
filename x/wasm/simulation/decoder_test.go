package simulation

import (
	"encoding/binary"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/tendermint/tendermint/crypto/ed25519"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/kv"

	"github.com/terra-money/core/x/wasm/keeper"
	"github.com/terra-money/core/x/wasm/types"
)

var (
	creatorPk    = ed25519.GenPrivKey().PubKey()
	contractPk   = ed25519.GenPrivKey().PubKey()
	creatorAddr  = sdk.AccAddress(creatorPk.Address())
	contractAddr = sdk.AccAddress(contractPk.Address())
)

func TestDecodeDistributionStore(t *testing.T) {
	cdc := keeper.MakeTestCodec(t)
	dec := NewDecodeStore(cdc)

	lastCodeIDbz := make([]byte, 8)
	lastInstanceIDbz := make([]byte, 8)
	binary.LittleEndian.PutUint64(lastCodeIDbz, 123)
	binary.LittleEndian.PutUint64(lastInstanceIDbz, 456)

	codeInfo := types.NewCodeInfo(1, []byte{1, 2, 3}, creatorAddr)
	contractInfo := types.NewContractInfo(1, contractAddr, creatorAddr, creatorAddr, []byte{4, 5, 6})
	emptyAdminContractInfo := types.NewContractInfo(1, contractAddr, creatorAddr, sdk.AccAddress{}, []byte{4, 5, 6})
	contractStore := []byte{7, 8, 9}

	kvPairs := kv.Pairs{
		Pairs: []kv.Pair{
			{Key: types.LastCodeIDKey, Value: lastCodeIDbz},
			{Key: types.LastInstanceIDKey, Value: lastInstanceIDbz},
			{Key: types.CodeKey, Value: cdc.MustMarshal(&codeInfo)},
			{Key: types.ContractInfoKey, Value: cdc.MustMarshal(&contractInfo)},
			{Key: append(types.ContractInfoKey, 0x1), Value: cdc.MustMarshal(&emptyAdminContractInfo)},
			{Key: types.ContractStoreKey, Value: contractStore},
			{Key: []byte{0x99}, Value: []byte{0x99}},
		},
	}

	tests := []struct {
		name        string
		expectedLog string
	}{
		{"LastCodeID", "lastCodeIDA: 123\nlastCodeIDB: 123"},
		{"LastInstanceID", "lastInstanceIDA: 456\nlastInstanceIDB: 456"},
		{"CodeInfo", fmt.Sprintf("%v\n%v", codeInfo, codeInfo)},
		{"ContractInfo", fmt.Sprintf("%v\n%v", contractInfo, contractInfo)},
		{"ContractInfo", fmt.Sprintf("%v\n%v", emptyAdminContractInfo, emptyAdminContractInfo)},
		{"ContractStore", fmt.Sprintf("%v\n%v", contractStore, contractStore)},
		{"other", ""},
	}

	for i, tt := range tests {
		i, tt := i, tt
		t.Run(tt.name, func(t *testing.T) {
			switch i {
			case len(tests) - 1:
				require.Panics(t, func() { dec(kvPairs.Pairs[i], kvPairs.Pairs[i]) }, tt.name)
			default:
				require.Equal(t, tt.expectedLog, dec(kvPairs.Pairs[i], kvPairs.Pairs[i]), tt.name)
			}
		})
	}
}
