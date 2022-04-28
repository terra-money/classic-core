package types

import (
	"encoding/binary"

	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
	"github.com/tendermint/tendermint/crypto"

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
func NewContractInfo(codeID uint64, address, creator, admin sdk.AccAddress, initMsg []byte) ContractInfo {
	var adminAddr string
	if !admin.Empty() {
		adminAddr = admin.String()
	}

	return ContractInfo{
		Address: address.String(),
		CodeID:  codeID,
		Creator: creator.String(),
		Admin:   adminAddr,
		InitMsg: initMsg,
	}
}

// NewEnv initializes the environment for a contract instance
func NewEnv(ctx sdk.Context, contractAddr sdk.AccAddress) wasmvmtypes.Env {
	env := wasmvmtypes.Env{
		Block: wasmvmtypes.BlockInfo{
			Height:  uint64(ctx.BlockHeight()),
			Time:    uint64(ctx.BlockTime().UnixNano()),
			ChainID: ctx.ChainID(),
		},
		Contract: wasmvmtypes.ContractInfo{
			Address: contractAddr.String(),
		},
	}

	if txCounter, ok := TXCounter(ctx); ok {
		env.Transaction = &wasmvmtypes.TransactionInfo{Index: txCounter}
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

// GenerateContractAddress generates a contract address from codeID + instanceID
// and increases last instanceID
func GenerateContractAddress(codeID uint64, instanceID uint64) sdk.AccAddress {
	// NOTE: It is possible to get a duplicate address if either codeID or instanceID
	// overflow 32 bits. This is highly improbable, but something that could be refactored.
	contractID := codeID<<32 + instanceID
	return addrFromUint64(contractID)
}

func addrFromUint64(id uint64) sdk.AccAddress {
	addr := make([]byte, 20)
	addr[0] = 'C'
	binary.PutUvarint(addr[1:], id)
	return sdk.AccAddress(crypto.AddressHash(addr))
}
