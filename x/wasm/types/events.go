// noalias
package types

// Wasm module event types
const (
	EventTypeStoreCode           = "store_code"
	EventTypeMigrateCode         = "migrate_code"
	EventTypeInstantiateContract = "instantiate_contract"
	EventTypeExecuteContract     = "execute_contract"
	EventTypeMigrateContract     = "migrate_contract"
	EventTypeUpdateContractAdmin = "update_contract_admin"
	EventTypeClearContractAdmin  = "clear_contract_admin"
	EventTypeWasmPrefix          = "wasm"

	// Deprecated
	EventTypeFromContract = "from_contract"

	AttributeKeySender          = "sender"
	AttributeKeyCodeID          = "code_id"
	AttributeKeyContractAddress = "contract_address"
	AttributeKeyContractID      = "contract_id"
	AttributeKeyAdmin           = "admin"
	AttributeKeyCreator         = "creator"

	AttributeValueCategory = ModuleName
)
