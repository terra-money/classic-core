// noalias
package types

// Wasm module event types
const (
	EventTypeStoreCode           = "store_code"
	EventTypeInstantiateContract = "instantiate_contract"
	EventTypeFromContract        = "from_contract"
	EventTypeUpdateContractOwner = "update_contract_owner"

	AttributeKeySender          = "sender"
	AttributeKeyCodeID          = "code_id"
	AttributeKeyContractAddress = "contract_address"
	AttributeKeyContractID      = "contract_id"
	AttributeKeyOwner           = "owner"

	AttributeValueCategory = ModuleName
)
