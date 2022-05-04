package stargatelayer

import (
	wasmvmtypes "github.com/CosmWasm/wasmvm/types"

	proto "github.com/gogo/protobuf/proto"
)

// NormalizeResponse normalize response to ensure deterministic
func NormalizeResponse(binding interface{}, bz []byte) ([]byte, error) {
	// all values are proto message
	message, ok := binding.(proto.Message)
	if !ok {
		return nil, wasmvmtypes.Unknown{}
	}

	// unmarshal binary into stargate response data structure
	err := proto.Unmarshal(bz, message)
	if err != nil {
		return nil, wasmvmtypes.Unknown{}
	}

	// build new deterministic response
	bz, err = proto.Marshal(message)
	if err != nil {
		return nil, wasmvmtypes.Unknown{}
	}

	// clear proto message
	message.Reset()

	return bz, nil
}
