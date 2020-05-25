package config

const (
	defaultLRUCacheSize          = uint64(0)
	defaultContractQueryGasLimit = uint64(3000000)
)

// config flags for wasm module
const (
	FlagContractQueryGasLimit = "contract-query-gas-limit"
	FlagCacheSize             = "lru-size"
)

// BaseConfig is the extra config required for wasm
type BaseConfig struct {
	// The maximum gas amount can be spent for contract query
	// The external query will invoke contract vm on wasm module,
	// so we need to restrict the max usage to prevent DoS attack
	ContractQueryGasLimit uint64 `mapstructure:"contract-query-gas-limit"`

	// Storing instances in the LRU will have no effect on
	// the results (still deterministic), but should lower
	// execution time at the cost of increased memory usage.
	// We cannot pick universal parameters for this, so
	// we should allow node operators to set it.
	CacheSize uint64 `mapstructure:"lru-size"`
}

// Config defines the server's top level configuration
type Config struct {
	BaseConfig `mapstructure:",squash"`
}

// DefaultConfig returns the default settings for WasmConfig
func DefaultConfig() *Config {
	return &Config{
		BaseConfig: BaseConfig{
			ContractQueryGasLimit: defaultContractQueryGasLimit,
			CacheSize:             defaultLRUCacheSize,
		},
	}
}
