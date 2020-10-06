package config

import "strings"

// config default values
const (
	defaultContractQueryGasLimit = uint64(3000000)
)

// config flags for wasm module
const (
	FlagContractQueryGasLimit    = "contract-query-gas-limit"
	FlagContractLoggingWhitelist = "contract-logging-whitelist"
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
}

// Config defines the server's top level configuration
type Config struct {
	BaseConfig `mapstructure:",squash"`
	loggingAll bool
}

// DefaultConfig returns the default settings for WasmConfig
func DefaultConfig() *Config {
	return &Config{
		BaseConfig: BaseConfig{
			ContractQueryGasLimit:    defaultContractQueryGasLimit,
			ContractLoggingWhitelist: "",
		},
		loggingAll: false,
	}
}

// LoggingAll return whitelist config is set to "*"
func (config Config) LoggingAll() bool {
	return config.loggingAll
}

// WhitelistToMap return logging whitelist map
func (config *Config) WhitelistToMap() (loggingWhitelist map[string]bool) {
	loggingWhitelist = make(map[string]bool)

	if config.ContractLoggingWhitelist != "*" {
		for _, addr := range strings.Split(config.ContractLoggingWhitelist, ",") {
			loggingWhitelist[addr] = true
		}

		config.loggingAll = false
	} else {
		config.loggingAll = true
	}

	return
}
