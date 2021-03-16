package types

import (
	wasmTypes "github.com/CosmWasm/go-cosmwasm/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewCodeInfo fills a new Contract struct
func NewCodeInfo(codeID uint64, codeHash []byte, creator sdk.AccAddress) CodeInfo {
	return CodeInfo{
		CodeID:   codeID,
		CodeHash: codeHash,
		Creator:  creator.String(),
	}
}

// NewContractInfo creates a new instance of a given WASM contract info
func NewContractInfo(codeID uint64, address, owner sdk.AccAddress, initMsg []byte, migratable bool) ContractInfo {
	return ContractInfo{
		Address:    address.String(),
		CodeID:     codeID,
		Owner:      owner.String(),
		InitMsg:    initMsg,
		Migratable: migratable,
	}
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
