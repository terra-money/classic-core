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

func TestParseEvents_Empty(t *testing.T) {
	_, _, addr := keyPubAddr()
	events, err := ParseEvents(addr, wasmvmtypes.EventAttributes{}, wasmvmtypes.Events{})
	require.NoError(t, err)
	require.Empty(t, events)
}

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

func TestParseEvents_WithoutAttributes(t *testing.T) {
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

func TestBuildEvent_EmptyEvent(t *testing.T) {
	_, _, addr := keyPubAddr()
	event := buildEvent(EventTypeExecuteContract, addr, wasmvmtypes.EventAttributes{})
	require.Nil(t, event)
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

func TestParseToCoin_InvalidAmount(t *testing.T) {
	_, err := ParseToCoin(wasmvmtypes.Coin{
		Denom:  "denom",
		Amount: "invalid",
	})

	require.Error(t, err)
}

func TestParseToCoins(t *testing.T) {
	coins := wasmvmtypes.Coins{wasmvmtypes.NewCoin(1234, core.MicroLunaDenom)}
	_, err := ParseToCoins(coins)
	require.NoError(t, err)
}

func TestParseToCoins_InvalidAmount(t *testing.T) {
	coins := wasmvmtypes.Coins{wasmvmtypes.Coin{Amount: "invalid", Denom: "denom"}}
	_, err := ParseToCoins(coins)
	require.Error(t, err)
}

func TestEncodeSdkCoins(t *testing.T) {
	sdkCoin := sdk.NewCoin(core.MicroLunaDenom, sdk.NewInt(15000))
	encodedCoin := EncodeSdkCoin(sdkCoin)
	sdkCoins := sdk.NewCoins(sdkCoin)
	encodedCoins := EncodeSdkCoins(sdkCoins)
	require.Equal(t, encodedCoin, encodedCoins[0])
}

func TestEncodeSdkEvents(t *testing.T) {
	events := EncodeSdkEvents(sdk.Events{sdk.NewEvent(
		EventTypeWasmPrefix,
		sdk.NewAttribute("k1", "v1"),
		sdk.NewAttribute("k2", "v2"),
		sdk.NewAttribute("k3", "v3"),
		sdk.NewAttribute("k4", "v4"),
	), sdk.NewEvent(
		EventTypeFromContract,
		sdk.NewAttribute("k1", "v1"),
		sdk.NewAttribute("k2", "v2"),
		sdk.NewAttribute("k3", "v3"),
		sdk.NewAttribute("k4", "v4"),
	)})

	require.Equal(t, wasmvmtypes.Events{wasmvmtypes.Event{
		Type: EventTypeWasmPrefix,
		Attributes: []wasmvmtypes.EventAttribute{
			{
				Key:   "k1",
				Value: "v1",
			},
			{
				Key:   "k2",
				Value: "v2",
			},
			{
				Key:   "k3",
				Value: "v3",
			},
			{
				Key:   "k4",
				Value: "v4",
			},
		},
	}}, events)
}

func TestConvertWasmIBCTimeoutHeightToCosmosHeight(t *testing.T) {
	emptyHeight := ConvertWasmIBCTimeoutHeightToCosmosHeight(nil)
	require.Zero(t, emptyHeight.RevisionNumber)
	require.Zero(t, emptyHeight.RevisionHeight)

	revision := uint64(10)
	height := uint64(20)
	ibcTimeoutBlock := wasmvmtypes.IBCTimeoutBlock{
		Revision: revision,
		Height:   height,
	}
	ibcHeight := ConvertWasmIBCTimeoutHeightToCosmosHeight(&ibcTimeoutBlock)
	require.Equal(t, revision, ibcHeight.RevisionNumber)
	require.Equal(t, height, ibcHeight.RevisionHeight)
}

func keyPubAddr() (crypto.PrivKey, crypto.PubKey, sdk.AccAddress) {
	key := ed25519.GenPrivKey()
	pub := key.PubKey()
	addr := sdk.AccAddress(pub.Address())
	return key, pub, addr
}
