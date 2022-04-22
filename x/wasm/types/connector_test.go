package types

import (
	"fmt"
	"testing"

	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/ed25519"

	sdk "github.com/cosmos/cosmos-sdk/types"
	core "github.com/terra-money/core/types"
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

	event := buildEvent(EventTypeFromContract, addr, wasmvmtypes.EventAttributes{
		{Key: "key1", Value: "value1"},
		{Key: "key2", Value: "value2"},
		{Key: "key3", Value: "value3"},
	})
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

func TestParseToCoins(t *testing.T) {
	coins := wasmvmtypes.Coins{wasmvmtypes.NewCoin(1234, core.MicroLunaDenom)}
	_, err := ParseToCoins(coins)
	require.NoError(t, err)
}

func TestEncodeSdkCoins(t *testing.T) {
	sdkCoin := sdk.NewCoin(core.MicroLunaDenom, sdk.NewInt(15000))
	encodedCoin := EncodeSdkCoin(sdkCoin)
	sdkCoins := sdk.NewCoins(sdkCoin)
	encodedCoins := EncodeSdkCoins(sdkCoins)
	require.Equal(t, encodedCoin, encodedCoins[0])
}

func keyPubAddr() (crypto.PrivKey, crypto.PubKey, sdk.AccAddress) {
	key := ed25519.GenPrivKey()
	pub := key.PubKey()
	addr := sdk.AccAddress(pub.Address())
	return key, pub, addr
}
