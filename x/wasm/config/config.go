package config

import (
	"github.com/spf13/cast"

	servertypes "github.com/cosmos/cosmos-sdk/server/types"
)

// config default values
const (
	DefaultContractQueryGasLimit  = uint64(3000000)
	DefaultContractDebugMode      = false
	DefaultWriteVMMemoryCacheSize = uint32(500)
	DefaultReadVMMemoryCacheSize  = uint32(300)
	DefaultNumReadVM              = uint32(1)
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

	// The write WASM VM memory cache size in MiB not bytes
	WriteVMMemoryCacheSize uint32 `mapstructure:"write-vm-memory-cache-size"`

	// The read WASM VM memory cache size in MiB not bytes
	ReadVMMemoryCacheSize uint32 `mapstructure:"read-vm-memory-cache-size"`

	// The number of read WASM VMs
	NumReadVMs uint32 `mapstructure:"num-read-vms"`
}

// DefaultConfig returns the default settings for WasmConfig
func DefaultConfig() *Config {
	return &Config{
		ContractQueryGasLimit:  DefaultContractQueryGasLimit,
		ContractDebugMode:      DefaultContractDebugMode,
		WriteVMMemoryCacheSize: DefaultWriteVMMemoryCacheSize,
		ReadVMMemoryCacheSize:  DefaultReadVMMemoryCacheSize,
		NumReadVMs:             DefaultNumReadVM,
	}
}

// GetConfig load config values from the app options
func GetConfig(appOpts servertypes.AppOptions) *Config {
	return &Config{
		ContractQueryGasLimit:  cast.ToUint64(appOpts.Get("wasm.contract-query-gas-limit")),
		ContractDebugMode:      cast.ToBool(appOpts.Get("wasm.contract-debug-mode")),
		WriteVMMemoryCacheSize: cast.ToUint32(appOpts.Get("wasm.write-vm-memory-cache-size")),
		ReadVMMemoryCacheSize:  cast.ToUint32(appOpts.Get("wasm.read-vm-memory-cache-size")),
		NumReadVMs:             cast.ToUint32(appOpts.Get("wasm.num-read-vms")),
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

# The write WASM VM memory cache size in MiB not bytes
write-vm-memory-cache-size = "{{ .WASMConfig.WriteVMMemoryCacheSize }}"

# The read WASM VM memory cache size in MiB not bytes
read-vm-memory-cache-size = "{{ .WASMConfig.ReadVMMemoryCacheSize }}"

# The number of read WASM VMs
num-read-vms = "{{ .WASMConfig.NumReadVMs }}"
`
