package types

const (
	defaultLRUCacheSize          = uint64(0)
	defaultContractQueryGasLimit = uint64(3000000)
)

// WasmConfig is the extra config required for wasm
type WasmConfig struct {
	ContractQueryGasLimit uint64 `mapstructure:"contract_query_gas_limit"`
	CacheSize             uint64 `mapstructure:"lru_size"`
}

// DefaultWasmConfig returns the default settings for WasmConfig
func DefaultWasmConfig() WasmConfig {
	return WasmConfig{
		ContractQueryGasLimit: defaultContractQueryGasLimit,
		CacheSize:             defaultLRUCacheSize,
	}
}

// WasmWrapper allows us to use namespacing in the config file
// This is only used for parsing in the app, x/wasm expects WasmConfig
type WasmWrapper struct {
	Wasm WasmConfig `mapstructure:"wasm"`
}
