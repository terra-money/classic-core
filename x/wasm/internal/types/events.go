// noalias
package types

// Wasm module event types
const (
	EventTypeStoreCode           = "store_code"
	EventTypeInstantiateContract = "instantiate_contract"
	EventTypeExecuteContract     = "execute_contract"

	AttributeKeySender          = "sender"
	AttributeKeyCodeID          = "code_id"
	AttributeKeyContractAddress = "contract_address"
	AttributeKeyContractID      = "contract_id"

	AttributeValueCategory = ModuleName
)
