package v05

import (
	v04wasm "github.com/terra-money/core/x/wasm/legacy/v04"
	v05wasm "github.com/terra-money/core/x/wasm/types"
)

// Migrate accepts exported v0.4 x/wasm and
// migrates it to v0.5 x/wasm genesis state. The migration includes:
//
// - Add new params for event and data size limit to x/wasm genesis state.
// - Change code bytes and code hash to empty bytes
// - Re-encode in v0.5 GenesisState.
func Migrate(
	wasmGenState v04wasm.GenesisState,
) *v05wasm.GenesisState {
	// CosmWasm version is not compatible, so remove code bytes and code hash
	codes := make([]v05wasm.Code, len(wasmGenState.Codes))
	for i, c := range wasmGenState.Codes {
		codes[i] = v05wasm.Code{
			CodeInfo: v05wasm.CodeInfo{
				CodeID:   c.CodeInfo.CodeID,
				CodeHash: []byte{},
				Creator:  c.CodeInfo.Creator.String(),
			},
			CodeBytes: []byte{},
		}
	}

	contracts := make([]v05wasm.Contract, len(wasmGenState.Contracts))
	for i, c := range wasmGenState.Contracts {
		models := make([]v05wasm.Model, len(c.ContractStore))
		for j, m := range c.ContractStore {
			models[j] = v05wasm.Model{
				Key:   m.Key,
				Value: m.Value,
			}
		}

		adminAddr := ""
		if c.ContractInfo.Migratable {
			adminAddr = c.ContractInfo.Owner.String()
		}

		contracts[i] = v05wasm.Contract{
			ContractInfo: v05wasm.ContractInfo{
				CodeID:  c.ContractInfo.CodeID,
				Address: c.ContractInfo.Address.String(),
				Creator: c.ContractInfo.Owner.String(),
				Admin:   adminAddr,
				InitMsg: c.ContractInfo.InitMsg,
			},
			ContractStore: models,
		}
	}

	return &v05wasm.GenesisState{
		Params: v05wasm.Params{
			MaxContractSize:    v05wasm.DefaultMaxContractSize,
			MaxContractMsgSize: v05wasm.DefaultMaxContractMsgSize,
			MaxContractGas:     v05wasm.DefaultMaxContractGas,
		},
		Codes:          codes,
		Contracts:      contracts,
		LastCodeID:     wasmGenState.LastCodeID,
		LastInstanceID: wasmGenState.LastInstanceID,
	}
}
