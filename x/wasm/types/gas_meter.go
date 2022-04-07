package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
)

// Constant gas parameters
const (
	GasMultiplier = uint64(140_000_000) // Please note that all gas prices returned to the wasmVM engine should have this multiplied

	compileCostPerByte             = uint64(3)       // sdk gas cost per bytes
	instantiateCost                = uint64(60_000)  // sdk gas cost for executing wasmVM engine
	registerCost                   = uint64(140_000) // sdk gas cost for creating contract
	humanizeCost                   = uint64(5)       // sdk gas cost to convert canonical address to human address
	canonicalizeCost               = uint64(4)       // sdk gas cost to convert human address to canonical address
	deserializationCostPerByte     = uint64(1)       // sdk gas cost to deserialize data
	eventAttributeCost             = uint64(10)      // sdk gas cost to emit attribute
	customEventCost                = uint64(20)      // sdk gas cost to emit custom event
	eventAttributeDataCostPerByte  = uint64(1)       // sdk gas cost per bytes to emit event
	contractMessageDataCostPerByte = uint64(1)       // sdk gas cost per bytes to execute submessages

	// HumanizeWasmGasCost humanize cost in wasm gas unit
	HumanizeWasmGasCost = humanizeCost * GasMultiplier
	// CanonicalizeWasmGasCost canonicalize cost in wasm gas unit
	CanonicalizeWasmGasCost = canonicalizeCost * GasMultiplier
)

var (
	// JSONDeserializationWasmGasCost json deserialization cost in wasm gas unit
	JSONDeserializationWasmGasCost = wasmvmtypes.UFraction{
		Numerator:   deserializationCostPerByte * GasMultiplier,
		Denominator: 1,
	}
)

// CompileCosts costs to persist and "compile" a new wasm contract
func CompileCosts(byteLength int) sdk.Gas {
	return sdk.NewUint(compileCostPerByte).MulUint64(uint64(byteLength)).Uint64()
}

// InstantiateContractCosts costs when interacting with a wasm contract
func InstantiateContractCosts(msgLen int) sdk.Gas {
	dataCosts := sdk.NewUint(sdk.Gas(msgLen)).MulUint64(contractMessageDataCostPerByte)
	return dataCosts.AddUint64(instantiateCost).Uint64()
}

// RegisterContractCosts costs when registering a new contract to the store
func RegisterContractCosts() sdk.Gas {
	return registerCost
}

// ReplyCosts costs to to handle a message reply
func ReplyCosts(reply wasmvmtypes.Reply) sdk.Gas {
	msgLen := len(reply.Result.Err)

	eventGas := sdk.NewUint(0)
	if reply.Result.Ok != nil {
		msgLen += len(reply.Result.Ok.Data)

		var attrs []wasmvmtypes.EventAttribute
		for _, e := range reply.Result.Ok.Events {
			msgLen += len(e.Type)
			attrs = append(attrs, e.Attributes...)
		}

		eventGas = eventGas.AddUint64(eventAttributeCosts(attrs))
	}

	return eventGas.AddUint64(InstantiateContractCosts(msgLen)).Uint64()
}

// EventCosts costs to persist an event
func EventCosts(attrs []wasmvmtypes.EventAttribute, events wasmvmtypes.Events) sdk.Gas {
	gas := sdk.NewUint(eventAttributeCosts(attrs))

	for _, e := range events {
		typeCost := eventTypeCosts(e.Type)
		attributeCost := eventAttributeCosts(e.Attributes)
		eventCost := sdk.NewUint(customEventCost).AddUint64(typeCost).AddUint64(attributeCost)
		gas = gas.Add(eventCost)
	}

	return gas.Uint64()
}

func eventTypeCosts(eventType string) sdk.Gas {
	if len(eventType) == 0 {
		return 0
	}

	return sdk.NewUint(uint64(len(eventType))).MulUint64(eventAttributeDataCostPerByte).Uint64()
}

func eventAttributeCosts(attrs []wasmvmtypes.EventAttribute) sdk.Gas {
	if len(attrs) == 0 {
		return 0
	}

	storedBytes := sdk.NewUint(0)
	for _, l := range attrs {
		storedBytes = storedBytes.AddUint64(uint64(len(l.Key) + len(l.Value)))
	}

	// total Length * costs + attribute count * costs
	lengthCost := storedBytes.MulUint64(eventAttributeDataCostPerByte)
	attributeCost := sdk.NewUint(eventAttributeCost).MulUint64(uint64(len(attrs)))
	r := lengthCost.Add(attributeCost)

	return r.Uint64()
}

// ToWasmVMGas converts from sdk gas to wasmvm gas
func ToWasmVMGas(source sdk.Gas) uint64 {
	return sdk.NewUint(source).MulUint64(GasMultiplier).Uint64()
}

// FromWasmVMGas converts from wasmvm gas to sdk gas
func FromWasmVMGas(source uint64) sdk.Gas {
	return source / GasMultiplier
}
