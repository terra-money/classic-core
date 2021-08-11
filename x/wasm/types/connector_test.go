package types

import (
	"fmt"
	"testing"

	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/ed25519"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestParseEvents(t *testing.T) {
	_, _, addr := keyPubAddr()

	events, err := ParseEvents(addr, wasmvmtypes.EventAttributes{
		{Key: "key1", Value: "value1"},
		{Key: "key2", Value: "value2"},
		{Key: "key3", Value: "value3"},
	}, wasmvmtypes.Events{
		{
			Type: "type1",
			Attributes: wasmvmtypes.EventAttributes{
				{Key: "key1", Value: "value1"},
				{Key: "key2", Value: "value2"},
				{Key: "key3", Value: "value3"},
			},
		},
		{
			Type: "type2",
			Attributes: wasmvmtypes.EventAttributes{
				{Key: "key1", Value: "value1"},
				{Key: "key2", Value: "value2"},
				{Key: "key3", Value: "value3"},
			},
		},
	})
	require.NoError(t, err)
	require.Equal(t, sdk.Events{sdk.NewEvent(
		EventTypeWasmPrefix,
		[]sdk.Attribute{
			{Key: AttributeKeyContractAddress, Value: addr.String()},
			{Key: "key1", Value: "value1"},
			{Key: "key2", Value: "value2"},
			{Key: "key3", Value: "value3"},
		}...,
	), sdk.NewEvent(
		EventTypeFromContract,
		[]sdk.Attribute{
			{Key: AttributeKeyContractAddress, Value: addr.String()},
			{Key: "key1", Value: "value1"},
			{Key: "key2", Value: "value2"},
			{Key: "key3", Value: "value3"},
		}...,
	), sdk.NewEvent(
		fmt.Sprintf("%s-type1", EventTypeWasmPrefix),
		[]sdk.Attribute{
			{Key: AttributeKeyContractAddress, Value: addr.String()},
			{Key: "key1", Value: "value1"},
			{Key: "key2", Value: "value2"},
			{Key: "key3", Value: "value3"},
		}...,
	), sdk.NewEvent(
		fmt.Sprintf("%s-type2", EventTypeWasmPrefix),
		[]sdk.Attribute{
			{Key: AttributeKeyContractAddress, Value: addr.String()},
			{Key: "key1", Value: "value1"},
			{Key: "key2", Value: "value2"},
			{Key: "key3", Value: "value3"},
		}...,
	)}, events)
}

func TestParseEventsWithoutAttributes(t *testing.T) {
	_, _, addr := keyPubAddr()

	events, err := ParseEvents(addr, wasmvmtypes.EventAttributes{}, wasmvmtypes.Events{
		{
			Type: "type1",
			Attributes: wasmvmtypes.EventAttributes{
				{Key: "key1", Value: "value1"},
				{Key: "key2", Value: "value2"},
				{Key: "key3", Value: "value3"},
			},
		},
		{
			Type: "type2",
			Attributes: wasmvmtypes.EventAttributes{
				{Key: "key1", Value: "value1"},
				{Key: "key2", Value: "value2"},
				{Key: "key3", Value: "value3"},
			},
		},
	})
	require.NoError(t, err)
	require.Equal(t, sdk.Events{sdk.NewEvent(
		fmt.Sprintf("%s-type1", EventTypeWasmPrefix),
		[]sdk.Attribute{
			{Key: AttributeKeyContractAddress, Value: addr.String()},
			{Key: "key1", Value: "value1"},
			{Key: "key2", Value: "value2"},
			{Key: "key3", Value: "value3"},
		}...,
	), sdk.NewEvent(
		fmt.Sprintf("%s-type2", EventTypeWasmPrefix),
		[]sdk.Attribute{
			{Key: AttributeKeyContractAddress, Value: addr.String()},
			{Key: "key1", Value: "value1"},
			{Key: "key2", Value: "value2"},
			{Key: "key3", Value: "value3"},
		}...,
	)}, events)
}

func TestBuildEvent(t *testing.T) {
	_, _, addr := keyPubAddr()

	event, err := buildEvent(EventTypeFromContract, addr, wasmvmtypes.EventAttributes{
		{Key: "key1", Value: "value1"},
		{Key: "key2", Value: "value2"},
		{Key: "key3", Value: "value3"},
	})
	require.NoError(t, err)
	require.Equal(t, sdk.NewEvent(
		EventTypeFromContract,
		[]sdk.Attribute{
			{Key: AttributeKeyContractAddress, Value: addr.String()},
			{Key: "key1", Value: "value1"},
			{Key: "key2", Value: "value2"},
			{Key: "key3", Value: "value3"},
		}...,
	), *event)
}

func keyPubAddr() (crypto.PrivKey, crypto.PubKey, sdk.AccAddress) {
	key := ed25519.GenPrivKey()
	pub := key.PubKey()
	addr := sdk.AccAddress(pub.Address())
	return key, pub, addr
}
