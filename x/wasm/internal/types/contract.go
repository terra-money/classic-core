package types

import (
	"encoding/hex"
	"fmt"

	wasmTypes "github.com/CosmWasm/go-cosmwasm/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	core "github.com/terra-project/core/types"
)

// Model is a struct that holds a KV pair
type Model struct {
	Key   core.HexBytes `json:"key"`
	Value core.HexBytes `json:"value"`
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
	CodeHash core.HexBytes  `json:"code_hash"`
	Creator  sdk.AccAddress `json:"creator"`
}

// String implements fmt.Stringer interface
func (ci CodeInfo) String() string {
	return fmt.Sprintf(`CodeInfo
	CodeHash:    %s, 
	Creator:     %s`,
		ci.CodeHash, ci.Creator)
}

// NewCodeInfo fills a new Contract struct
func NewCodeInfo(codeHash []byte, creator sdk.AccAddress) CodeInfo {
	return CodeInfo{
		CodeHash: codeHash,
		Creator:  creator,
	}
}

// ContractInfo stores a WASM contract instance
type ContractInfo struct {
	CodeID  uint64         `json:"code_id"`
	Address sdk.AccAddress `json:"address"`
	Creator sdk.AccAddress `json:"creator"`
	InitMsg core.HexBytes  `json:"init_msg"`
}

// NewContractInfo creates a new instance of a given WASM contract info
func NewContractInfo(codeID uint64, address sdk.AccAddress, creator sdk.AccAddress, initMsg []byte) ContractInfo {
	return ContractInfo{
		CodeID:  codeID,
		Address: address,
		Creator: creator,
		InitMsg: initMsg,
	}
}

// String implements fmt.Stringer interface
func (ci ContractInfo) String() string {
	return fmt.Sprintf(`ContractInfo
	CodeID:     %d, 
	Creator:    %s,
	InitMsg:    %s`,
		ci.CodeID, ci.Creator, hex.EncodeToString(ci.InitMsg))
}

// NewWasmAPIParams initializes params for a contract instance
func NewWasmAPIParams(ctx sdk.Context, creator sdk.AccAddress, deposit sdk.Coins, contractAddr sdk.AccAddress) wasmTypes.Env {
	return wasmTypes.Env{
		Block: wasmTypes.BlockInfo{
			Height:  uint64(ctx.BlockHeight()),
			Time:    uint64(ctx.BlockTime().Unix()),
			ChainID: ctx.ChainID(),
		},
		Message: wasmTypes.MessageInfo{
			Sender:    wasmTypes.CanonicalAddress(creator),
			SentFunds: NewWasmCoins(deposit),
		},
		Contract: wasmTypes.ContractInfo{
			Address: wasmTypes.CanonicalAddress(contractAddr),
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
