package types

import (
	"fmt"

	wasmTypes "github.com/CosmWasm/go-cosmwasm/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	core "github.com/terra-money/core/types"
)

// Model is a struct that holds a KV pair
type Model struct {
	Key   core.Base64Bytes `json:"key"`
	Value core.Base64Bytes `json:"value"`
}

// String implements fmt.Stringer interface
func (m Model) String() string {
	return fmt.Sprintf(`Model
	Key:   %s,
	Value: %s`,
		m.Key, m.Value)
}

// CodeInfo is data for the uploaded contract WASM code
type CodeInfo struct {
	CodeID   uint64           `json:"code_id"`
	CodeHash core.Base64Bytes `json:"code_hash"`
	Creator  sdk.AccAddress   `json:"creator"`
}

// String implements fmt.Stringer interface
func (ci CodeInfo) String() string {
	return fmt.Sprintf(`CodeInfo
	CodeID:      %d,
	CodeHash:    %s, 
	Creator:     %s`,
		ci.CodeID, ci.CodeHash, ci.Creator)
}

// NewCodeInfo fills a new Contract struct
func NewCodeInfo(codeID uint64, codeHash []byte, creator sdk.AccAddress) CodeInfo {
	return CodeInfo{
		CodeID:   codeID,
		CodeHash: codeHash,
		Creator:  creator,
	}
}

// ContractInfo stores a WASM contract instance
type ContractInfo struct {
	Address    sdk.AccAddress   `json:"address"`
	Owner      sdk.AccAddress   `json:"owner"`
	CodeID     uint64           `json:"code_id"`
	InitMsg    core.Base64Bytes `json:"init_msg"`
	Migratable bool             `json:"migratable"`
}

// NewContractInfo creates a new instance of a given WASM contract info
func NewContractInfo(codeID uint64, address, owner sdk.AccAddress, initMsg []byte, migratable bool) ContractInfo {
	return ContractInfo{
		Address:    address,
		CodeID:     codeID,
		Owner:      owner,
		InitMsg:    initMsg,
		Migratable: migratable,
	}
}

// String implements fmt.Stringer interface
func (ci ContractInfo) String() string {
	return fmt.Sprintf(`ContractInfo
	Address:    %s,
	CodeID:     %d, 
	Owner:      %s,
	InitMsg:    %s,
	Migratable  %v,
	`,
		ci.Address, ci.CodeID, ci.Owner, ci.InitMsg, ci.Migratable)
}

// NewWasmAPIParams initializes params for a contract instance
func NewWasmAPIParams(ctx sdk.Context, sender sdk.AccAddress, deposit sdk.Coins, contractAddr sdk.AccAddress) wasmTypes.Env {
	return wasmTypes.Env{
		Block: wasmTypes.BlockInfo{
			Height:  uint64(ctx.BlockHeight()),
			Time:    uint64(ctx.BlockTime().Unix()),
			ChainID: ctx.ChainID(),
		},
		Message: wasmTypes.MessageInfo{
			Sender:    sender.String(),
			SentFunds: NewWasmCoins(deposit),
		},
		Contract: wasmTypes.ContractInfo{
			Address: contractAddr.String(),
		},
	}
}

// NewWasmCoins translates between Cosmos SDK coins and Wasm coins
func NewWasmCoins(cosmosCoins sdk.Coins) (wasmCoins []wasmTypes.Coin) {
	for _, coin := range cosmosCoins {
		wasmCoin := wasmTypes.Coin{
			Denom:  coin.Denom,
			Amount: coin.Amount.String(),
		}
		wasmCoins = append(wasmCoins, wasmCoin)
	}
	return wasmCoins
}
