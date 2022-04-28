package types

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"

	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
)

func TestRegisterGasCost(t *testing.T) {
	require.Equal(t, registerCost, RegisterContractCosts())
}

func TestOverflow_ToWasmVM(t *testing.T) {
	require.Panics(t,
		func() { ToWasmVMGas(uint64(0xff_ff_ff_ff_ff_ff_ff_ff)) },
	)

	require.NotPanics(t, func() {
		ToWasmVMGas(uint64(0xff_ff_ff_ff))
	})
}

func TestFromWasmVMGas(t *testing.T) {
	res := FromWasmVMGas(uint64(0xff_ff_ff_ff))
	require.Equal(t, uint64(0xff_ff_ff_ff)/GasMultiplier, res)
}

func TestCompileCosts(t *testing.T) {
	bzLength := 10

	cost := CompileCosts(bzLength)
	require.Equal(t, sdk.Gas(10*compileCostPerByte), cost)
}

func TestInstantiateContractCosts(t *testing.T) {
	msgLength := 10

	cost := InstantiateContractCosts(msgLength)
	require.Equal(t, sdk.Gas(instantiateCost+uint64(msgLength)*contractMessageDataCostPerByte), cost)
}

func TestReplyCosts(t *testing.T) {
	eventsNum := uint64(10)
	dataLength := uint64(10)

	typeLength := uint64(10)
	attributesNum := uint64(10)
	keyLength := uint64(10)
	valueLength := uint64(10)

	attributes := make(wasmvmtypes.EventAttributes, attributesNum)
	for i := 0; i < int(attributesNum); i++ {
		attributes[i] = wasmvmtypes.EventAttribute{
			Key:   strings.Repeat("a", int(keyLength)),
			Value: strings.Repeat("a", int(valueLength)),
		}
	}

	events := make(wasmvmtypes.Events, eventsNum)
	for i := 0; i < int(eventsNum); i++ {
		events[i] = wasmvmtypes.Event{
			Type:       strings.Repeat("a", int(typeLength)),
			Attributes: attributes,
		}
	}

	reply := wasmvmtypes.Reply{
		ID: 10,
		Result: wasmvmtypes.SubcallResult{
			Ok: &wasmvmtypes.SubcallResponse{
				Events: events,
				Data:   make([]byte, dataLength),
			},
		},
	}

	cost := ReplyCosts(reply)

	totalAttributesNum := eventsNum * attributesNum
	require.Equal(t,
		instantiateCost+
			(typeLength*eventsNum+dataLength)*contractMessageDataCostPerByte+
			((keyLength+valueLength)*eventAttributeDataCostPerByte+eventAttributeCost)*totalAttributesNum,
		cost)
}

func TestEmptyEventCost(t *testing.T) {
	require.Equal(t, sdk.Gas(0), eventTypeCosts(""))
}

func TestEmptyAttributeCosts(t *testing.T) {
	require.Equal(t, sdk.Gas(0), eventAttributeCosts([]wasmvmtypes.EventAttribute{}))
}

func TestEventCosts(t *testing.T) {
	eventsNum := uint64(10)

	typeLength := uint64(10)
	attributesNum := uint64(10)
	keyLength := uint64(10)
	valueLength := uint64(10)

	attributes := make(wasmvmtypes.EventAttributes, attributesNum)
	for i := 0; i < int(attributesNum); i++ {
		attributes[i] = wasmvmtypes.EventAttribute{
			Key:   strings.Repeat("a", int(keyLength)),
			Value: strings.Repeat("a", int(valueLength)),
		}
	}

	events := make(wasmvmtypes.Events, eventsNum)
	for i := 0; i < int(eventsNum); i++ {
		events[i] = wasmvmtypes.Event{
			Type:       strings.Repeat("a", int(typeLength)),
			Attributes: attributes,
		}
	}

	cost := EventCosts(attributes, events)
	require.Equal(t,
		((keyLength+valueLength)*eventAttributeDataCostPerByte+eventAttributeCost)*attributesNum+
			(typeLength*eventAttributeDataCostPerByte+ // type cost
				((keyLength+valueLength)*eventAttributeDataCostPerByte+eventAttributeCost)*attributesNum+ // event attribute cost
				customEventCost)*eventsNum,
		cost)
}
