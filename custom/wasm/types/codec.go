package types

import (
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	"github.com/cosmos/cosmos-sdk/codec"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
)

// RegisterLegacyAminoCodec registers the account types and interface
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) { //nolint:staticcheck
	cdc.RegisterConcrete(&wasmtypes.MsgStoreCode{}, "wasm/MsgStoreCode", nil)
	cdc.RegisterConcrete(&wasmtypes.MsgInstantiateContract{}, "wasm/MsgInstantiateContract", nil)
	cdc.RegisterConcrete(&wasmtypes.MsgInstantiateContract2{}, "wasm/MsgInstantiateContract2", nil)
	cdc.RegisterConcrete(&wasmtypes.MsgExecuteContract{}, "wasm/MsgExecuteContract", nil)
	cdc.RegisterConcrete(&wasmtypes.MsgMigrateContract{}, "wasm/MsgMigrateContract", nil)
	cdc.RegisterConcrete(&wasmtypes.MsgUpdateAdmin{}, "wasm/MsgUpdateAdmin", nil)
	cdc.RegisterConcrete(&wasmtypes.MsgClearAdmin{}, "wasm/MsgClearAdmin", nil)

	cdc.RegisterConcrete(&wasmtypes.PinCodesProposal{}, "wasm/PinCodesProposal", nil)
	cdc.RegisterConcrete(&wasmtypes.UnpinCodesProposal{}, "wasm/UnpinCodesProposal", nil)
	cdc.RegisterConcrete(&wasmtypes.StoreCodeProposal{}, "wasm/StoreCodeProposal", nil)
	cdc.RegisterConcrete(&wasmtypes.InstantiateContractProposal{}, "wasm/InstantiateContractProposal", nil)
	cdc.RegisterConcrete(&wasmtypes.MigrateContractProposal{}, "wasm/MigrateContractProposal", nil)
	cdc.RegisterConcrete(&wasmtypes.SudoContractProposal{}, "wasm/SudoContractProposal", nil)
	cdc.RegisterConcrete(&wasmtypes.ExecuteContractProposal{}, "wasm/ExecuteContractProposal", nil)
	cdc.RegisterConcrete(&wasmtypes.UpdateAdminProposal{}, "wasm/UpdateAdminProposal", nil)
	cdc.RegisterConcrete(&wasmtypes.ClearAdminProposal{}, "wasm/ClearAdminProposal", nil)
	cdc.RegisterConcrete(&wasmtypes.UpdateInstantiateConfigProposal{}, "wasm/UpdateInstantiateConfigProposal", nil)

	cdc.RegisterInterface((*wasmtypes.ContractAuthzFilterX)(nil), nil)
	cdc.RegisterConcrete(&wasmtypes.AllowAllMessagesFilter{}, "wasm/AllowAllMessagesFilter", nil)
	cdc.RegisterConcrete(&wasmtypes.AcceptedMessageKeysFilter{}, "wasm/AcceptedMessageKeysFilter", nil)
	cdc.RegisterConcrete(&wasmtypes.AcceptedMessagesFilter{}, "wasm/AcceptedMessagesFilter", nil)

	cdc.RegisterInterface((*wasmtypes.ContractAuthzLimitX)(nil), nil)
	cdc.RegisterConcrete(&wasmtypes.MaxCallsLimit{}, "wasm/MaxCallsLimit", nil)
	cdc.RegisterConcrete(&wasmtypes.MaxFundsLimit{}, "wasm/MaxFundsLimit", nil)
	cdc.RegisterConcrete(&wasmtypes.CombinedLimit{}, "wasm/CombinedLimit", nil)

	cdc.RegisterConcrete(&wasmtypes.ContractExecutionAuthorization{}, "wasm/ContractExecutionAuthorization", nil)
	cdc.RegisterConcrete(&wasmtypes.ContractMigrationAuthorization{}, "wasm/ContractMigrationAuthorization", nil)
	cdc.RegisterConcrete(&wasmtypes.StoreAndInstantiateContractProposal{}, "wasm/StoreAndInstantiateContractProposal", nil)
}

var (
	amino = codec.NewLegacyAmino()

	// ModuleCdc references the global x/wasm module codec.

	ModuleCdc = codec.NewAminoCodec(amino)
)

func init() {
	RegisterLegacyAminoCodec(amino)
	cryptocodec.RegisterCrypto(amino)
	amino.Seal()
}
