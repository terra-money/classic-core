package types

import (
	wasmvmtypes "github.com/CosmWasm/wasmvm/types"

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

// NewEnv initializes the environment for a contract instance
func NewEnv(ctx sdk.Context, contractAddr sdk.AccAddress) wasmvmtypes.Env {
	env := wasmvmtypes.Env{
		Block: wasmvmtypes.BlockInfo{
			Height:    uint64(ctx.BlockHeight()),
			Time:      uint64(ctx.BlockTime().Unix()),
			TimeNanos: uint64(ctx.BlockTime().Nanosecond()),
			ChainID:   ctx.ChainID(),
		},
		Contract: wasmvmtypes.ContractInfo{
			Address: contractAddr.String(),
		},
	}
	return env
}

// NewInfo initializes the MessageInfo for a contract instance
func NewInfo(creator sdk.AccAddress, deposit sdk.Coins) wasmvmtypes.MessageInfo {
	return wasmvmtypes.MessageInfo{
		Sender: creator.String(),
		Funds:  NewWasmCoins(deposit),
	}
}

// NewWasmCoins translates between Cosmos SDK coins and Wasm coins
func NewWasmCoins(cosmosCoins sdk.Coins) (wasmCoins []wasmvmtypes.Coin) {
	for _, coin := range cosmosCoins {
		wasmCoin := wasmvmtypes.Coin{
			Denom:  coin.Denom,
			Amount: coin.Amount.String(),
		}
		wasmCoins = append(wasmCoins, wasmCoin)
	}
	return wasmCoins
}
