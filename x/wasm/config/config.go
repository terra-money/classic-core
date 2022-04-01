package config

import (
	"github.com/spf13/cast"

	servertypes "github.com/cosmos/cosmos-sdk/server/types"
)

// config default values
const (
	DefaultContractQueryGasLimit   = uint64(3000000)
	DefaultContractDebugMode       = false
	DefaultContractMemoryCacheSize = uint32(100)
)

// DBDir used to store wasm data to
var DBDir = "data/wasm"

// Config is the extra config required for wasm
type Config struct {
	// The maximum gas amount can be spent for contract query
	// The external query will invoke contract vm on wasm module,
	// so we need to restrict the max usage to prevent DoS attack
	ContractQueryGasLimit uint64 `mapstructure:"contract-query-gas-limit"`

	// The flag to specify whether print contract logs or not
	ContractDebugMode bool `mapstructure:"contract-debug-mode"`

	// The WASM VM memory cache size in MiB not bytes
	ContractMemoryCacheSize uint32 `mapstructure:"contract-memory-cache-size"`
}

// DefaultConfig returns the default settings for WasmConfig
func DefaultConfig() *Config {
	return &Config{
		ContractQueryGasLimit:   DefaultContractQueryGasLimit,
		ContractDebugMode:       DefaultContractDebugMode,
		ContractMemoryCacheSize: DefaultContractMemoryCacheSize,
	}
}

// GetConfig load config values from the app options
func GetConfig(appOpts servertypes.AppOptions) *Config {
	return &Config{
		ContractQueryGasLimit:   cast.ToUint64(appOpts.Get("wasm.contract-query-gas-limit")),
		ContractDebugMode:       cast.ToBool(appOpts.Get("wasm.contract-debug-mode")),
		ContractMemoryCacheSize: cast.ToUint32(appOpts.Get("wasm.contract-memory-cache-size")),
	}
}

// DefaultConfigTemplate default config template for wasm module
const DefaultConfigTemplate = `
[wasm]
# The maximum gas amount can be spent for contract query.
# The contract query will invoke contract execution vm,
# so we need to restrict the max usage to prevent DoS attack
contract-query-gas-limit = "{{ .WASMConfig.ContractQueryGasLimit }}"

# The flag to specify whether print contract logs or not
contract-debug-mode = "{{ .WASMConfig.ContractDebugMode }}"

# The WASM VM memory cache size in MiB not bytes
contract-memory-cache-size = "{{ .WASMConfig.ContractMemoryCacheSize }}"
`
