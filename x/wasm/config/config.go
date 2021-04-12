package config

import (
	"github.com/spf13/cobra"
)

// config default values
const (
	DefaultContractQueryGasLimit   = uint64(3000000)
	DefaultContractDebugMode       = false
	DefaultContractMemoryCacheSize = uint32(500)
)

// config flags for wasm module
const (
	FlagContractQueryGasLimit   = "contract-query-gas-limit"
	FlagContractDebugMode       = "contract-debug-mode"
	FlagContractMemoryCacheSize = "contract-memory-cache-size"
)

// DBDir used to store wasm data to
var DBDir = "data/wasm"

// BaseConfig is the extra config required for wasm
type BaseConfig struct {
	// The maximum gas amount can be spent for contract query
	// The external query will invoke contract vm on wasm module,
	// so we need to restrict the max usage to prevent DoS attack
	ContractQueryGasLimit uint64 `mapstructure:"contract-query-gas-limit"`

	// Only The logs from the contracts, which are listed in
	// this array or instantiated from the address in this array,
	// are stored in the local storage. To keep all logs,
	// a node operator can set "*" (not recommended).
	ContractLoggingWhitelist string `mapstructure:"contract-logging-whitelist"`

	// The flag to specify whether print contract logs or not
	ContractDebugMode bool `mapstructure:"contract-debug-mode"`

	// The WASM VM memory cache size in MiB not bytes
	ContractMemoryCacheSize uint32 `mapstructure:"contract-memory-cache-size"`
}

// Config defines the server's top level configuration
type Config struct {
	BaseConfig `mapstructure:",squash"`
}

// DefaultConfig returns the default settings for WasmConfig
func DefaultConfig() *Config {
	return &Config{
		BaseConfig: BaseConfig{
			ContractQueryGasLimit:   DefaultContractQueryGasLimit,
			ContractDebugMode:       DefaultContractDebugMode,
			ContractMemoryCacheSize: DefaultContractMemoryCacheSize,
		},
	}
}

// AddModuleInitFlags implements servertypes.ModuleInitFlags interface.
func AddModuleInitFlags(startCmd *cobra.Command) {
	startCmd.Flags().Uint64(FlagContractQueryGasLimit, DefaultContractQueryGasLimit, "The maximum gas amount can be spent for contract query")
	startCmd.Flags().Bool(FlagContractDebugMode, DefaultContractDebugMode, "The flag to specify whether print contract logs or not")
	startCmd.Flags().Uint32(FlagContractMemoryCacheSize, DefaultContractMemoryCacheSize, "The WASM VM memory cache size in MiB not bytes")
}
