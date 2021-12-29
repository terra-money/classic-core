package v05

import (
	"gopkg.in/yaml.v2"
)

const (
	// ModuleName is the name of the wasm module
	ModuleName = "wasm"

	defaultMaxContractSize    = uint64(600 * 1024) // 600 KB
	defaultMaxContractGas     = uint64(20_000_000) // 20,000,000
	defaultMaxContractMsgSize = uint64(4 * 1024)   // 4KB
)

// String implements fmt.Stringer interface
func (p Params) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}
