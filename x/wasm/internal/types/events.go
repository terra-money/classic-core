// noalias
package types

// Wasm module event types
const (
	EventTypeStoreCode           = "store_code"
	EventTypeInstantiateContract = "instantiate_contract"
	EventTypeExecuteContract     = "execute_contract"
	EventTypeMigrateContract     = "migrate_contract"
	EventTypeUpdateContractOwner = "update_contract_owner"
	EventTypeFromContract        = "from_contract"

	AttributeKeySender          = "sender"
	AttributeKeyCodeID          = "code_id"
	AttributeKeyContractAddress = "contract_address"
	AttributeKeyContractID      = "contract_id"
	AttributeKeyOwner           = "owner"

	AttributeValueCategory = ModuleName
)
