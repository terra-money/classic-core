package v05

import (
	v04wasm "github.com/terra-money/core/x/wasm/legacy/v04"
)

// Migrate accepts exported v0.4 x/wasm and
// migrates it to v0.5 x/wasm genesis state. The migration includes:
//
// - Add new params for event and data size limit to x/wasm genesis state.
// - Change code bytes and code hash to empty bytes
// - Re-encode in v0.5 GenesisState.
func Migrate(
	wasmGenState v04wasm.GenesisState,
) *GenesisState {
	// CosmWasm version is not compatible, so remove code bytes and code hash
	codes := make([]Code, len(wasmGenState.Codes))
	for i, c := range wasmGenState.Codes {
		codes[i] = Code{
			CodeInfo: CodeInfo{
				CodeID:   c.CodeInfo.CodeID,
				CodeHash: []byte{},
				Creator:  c.CodeInfo.Creator.String(),
			},
			CodeBytes: []byte{},
		}
	}

	contracts := make([]Contract, len(wasmGenState.Contracts))
	for i, c := range wasmGenState.Contracts {
		models := make([]Model, len(c.ContractStore))
		for j, m := range c.ContractStore {
			models[j] = Model{
				Key:   m.Key,
				Value: m.Value,
			}
		}

		adminAddr := ""
		if c.ContractInfo.Migratable {
			adminAddr = c.ContractInfo.Owner.String()
		}

		contracts[i] = Contract{
			ContractInfo: ContractInfo{
				CodeID:  c.ContractInfo.CodeID,
				Address: c.ContractInfo.Address.String(),
				Creator: c.ContractInfo.Owner.String(),
				Admin:   adminAddr,
				InitMsg: c.ContractInfo.InitMsg,
			},
			ContractStore: models,
		}
	}

	return &GenesisState{
		Params: Params{
			MaxContractSize:    defaultMaxContractSize,
			MaxContractMsgSize: defaultMaxContractMsgSize,
			MaxContractGas:     defaultMaxContractGas,
		},
		Codes:          codes,
		Contracts:      contracts,
		LastCodeID:     wasmGenState.LastCodeID,
		LastInstanceID: wasmGenState.LastInstanceID,
	}
}
