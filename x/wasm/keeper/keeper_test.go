package keeper

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/terra-money/core/x/wasm/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestNewKeeper(t *testing.T) {
	input := CreateTestInput(t)
	keeper := input.WasmKeeper
	require.NotNil(t, keeper)
}

func TestCodeInfo(t *testing.T) {
	input := CreateTestInput(t)
	ctx, keeper := input.Ctx, input.WasmKeeper

	codeID := uint64(1)
	_, _, creatorAddr := keyPubAddr()
	expected := types.NewCodeInfo(codeID, []byte{1, 2, 3}, creatorAddr)
	keeper.SetCodeInfo(ctx, 1, expected)

	as, err := keeper.GetCodeInfo(ctx, codeID)
	require.NoError(t, err)
	require.Equal(t, expected, as)
}

func TestContractInfo(t *testing.T) {
	input := CreateTestInput(t)
	ctx, keeper := input.Ctx, input.WasmKeeper

	_, _, alice := keyPubAddr()
	_, _, bob := keyPubAddr()

	codeID := uint64(1)
	instanceID := uint64(1)
	_, _, creatorAddr := keyPubAddr()
	contractAddr := types.GenerateContractAddress(codeID, instanceID)

	initMsg := HackatomExampleInitMsg{
		Verifier:    alice,
		Beneficiary: bob,
	}
	initMsgBz, err := json.Marshal(initMsg)
	require.NoError(t, err)

	expected := types.NewContractInfo(codeID, contractAddr, creatorAddr, sdk.AccAddress{}, initMsgBz)
	keeper.SetContractInfo(ctx, contractAddr, expected)

	as, err := keeper.GetContractInfo(ctx, contractAddr)
	require.NoError(t, err)
	require.Equal(t, expected, as)

	keeper.IterateContractInfo(ctx, func(contractInfo types.ContractInfo) bool {
		require.Equal(t, expected, contractInfo)
		return false
	})
}

func TestContractStore(t *testing.T) {
	models := []types.Model{
		{
			Key:   []byte("a"),
			Value: []byte("aa"),
		},
		{
			Key:   []byte("b"),
			Value: []byte("bb"),
		},
		{
			Key:   []byte("c"),
			Value: []byte("cc"),
		},
	}

	input := CreateTestInput(t)
	ctx, keeper := input.Ctx, input.WasmKeeper

	_, _, contractAddr := keyPubAddr()
	keeper.SetContractStore(ctx, contractAddr, models)

	i := 0
	for iter := keeper.GetContractStoreIterator(ctx, contractAddr); iter.Valid(); iter.Next() {
		require.Equal(t, models[i], types.Model{
			Key:   iter.Key(),
			Value: iter.Value(),
		})

		i++
	}
}
