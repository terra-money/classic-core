package types

import sdk "github.com/cosmos/cosmos-sdk/types"

// query endpoints supported by the wasm Querier
const (
	QueryGetByteCode     = "bytecode"
	QueryGetCodeInfo     = "codeInfo"
	QueryGetContractInfo = "contractInfo"
	QueryRawStore        = "rawStore"
	QueryContractStore   = "contractStore"
	QueryParameters      = "parameters"
)

// QueryCodeIDParams defines the params for the following queries:
// - 'custom/wasm/codeInfo
// - 'custom/wasm/bytecode
type QueryCodeIDParams struct {
	CodeID uint64
}

// NewQueryCodeIDParams returns QueryCodeIDParams instance
func NewQueryCodeIDParams(codeID uint64) QueryCodeIDParams {
	return QueryCodeIDParams{codeID}
}

// QueryContractAddressParams defines the params for the following queries:
// - 'custom/wasm/contractInfo
type QueryContractAddressParams struct {
	ContractAddress sdk.AccAddress
}

// NewQueryContractAddressParams returns QueryContractAddressParams instance
func NewQueryContractAddressParams(contractAddress sdk.AccAddress) QueryContractAddressParams {
	return QueryContractAddressParams{contractAddress}
}

// QueryRawStoreParams defines the params for the following queries:
// - 'custom/wasm/rawStore'
type QueryRawStoreParams struct {
	ContractAddress sdk.AccAddress
	Key             []byte
}

// NewQueryRawStoreParams returns QueryRawStoreParams instance
func NewQueryRawStoreParams(contractAddress sdk.AccAddress, key []byte) QueryRawStoreParams {
	return QueryRawStoreParams{contractAddress, key}
}

// QueryContractParams defines the params for the following queries:
// - 'custom/wasm/contractStore'
type QueryContractParams struct {
	ContractAddress sdk.AccAddress
	Msg             []byte
}

// NewQueryContractParams returns QueryContractParams instance
func NewQueryContractParams(contractAddress sdk.AccAddress, msg []byte) QueryContractParams {
	return QueryContractParams{contractAddress, msg}
}
