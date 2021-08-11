package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
)

// Constant gas parameters
const (
	GasMultiplier = uint64(100) // Please note that all gas prices returned to the wasmVM engine should have this multiplied

	compileCostPerByte             = uint64(2)       // sdk gas cost per bytes
	instantiateCost                = uint64(40_000)  // sdk gas cost for executing wasmVM engine
	registerCost                   = uint64(160_000) // sdk gas cost for creating contract
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
	return compileCostPerByte * uint64(byteLength)
}

// InstantiateContractCosts costs when interacting with a wasm contract
func InstantiateContractCosts(msgLen int) sdk.Gas {
	dataCosts := sdk.Gas(msgLen) * contractMessageDataCostPerByte
	return instantiateCost + dataCosts
}

// RegisterContractCosts costs when registering a new contract to the store
func RegisterContractCosts() sdk.Gas {
	return registerCost
}

// ReplyCosts costs to to handle a message reply
func ReplyCosts(reply wasmvmtypes.Reply) sdk.Gas {
	msgLen := len(reply.Result.Err)

	var eventGas sdk.Gas
	if reply.Result.Ok != nil {
		msgLen += len(reply.Result.Ok.Data)

		var attrs []wasmvmtypes.EventAttribute
		for _, e := range reply.Result.Ok.Events {
			msgLen += len(e.Type)
			attrs = append(attrs, e.Attributes...)
		}

		eventGas += eventAttributeCosts(attrs)
	}

	return eventGas + InstantiateContractCosts(msgLen)
}

// EventCosts costs to persist an event
func EventCosts(attrs []wasmvmtypes.EventAttribute, events wasmvmtypes.Events) sdk.Gas {
	gas := eventAttributeCosts(attrs)
	for _, e := range events {
		gas += customEventCost + sdk.Gas(len(e.Type))*eventAttributeDataCostPerByte
		gas += eventAttributeCosts(e.Attributes)
	}
	return gas
}

func eventAttributeCosts(attrs []wasmvmtypes.EventAttribute) sdk.Gas {
	if len(attrs) == 0 {
		return 0
	}

	var storedBytes uint64
	for _, l := range attrs {
		storedBytes += uint64(len(l.Key) + len(l.Value))
	}

	// total Length * costs + attribute count * costs
	r := sdk.NewIntFromUint64(eventAttributeDataCostPerByte).Mul(sdk.NewIntFromUint64(storedBytes)).
		Add(sdk.NewIntFromUint64(eventAttributeCost).Mul(sdk.NewIntFromUint64(uint64(len(attrs)))))
	if !r.IsUint64() {
		panic(sdk.ErrorOutOfGas{Descriptor: "overflow"})
	}

	return r.Uint64()
}

// ToWasmVMGas converts from sdk gas to wasmvm gas
func ToWasmVMGas(source sdk.Gas) uint64 {
	x := source * GasMultiplier
	if x < source {
		panic(sdk.ErrorOutOfGas{Descriptor: "overflow"})
	}
	return x
}

// FromWasmVMGas converts from wasmvm gas to sdk gas
func FromWasmVMGas(source uint64) sdk.Gas {
	return source / GasMultiplier
}
