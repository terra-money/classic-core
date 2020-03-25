package wasm

import (
	"bytes"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/terra-project/core/x/wasm/internal/types"
	// authexported "github.com/cosmos/cosmos-sdk/x/auth/exported"
	// "github.com/terra-project/core/x/wasm/internal/types"
)

// InitGenesis sets wasm information for genesis.
//
// CONTRACT: all types of accounts must have been already initialized/created
func InitGenesis(ctx sdk.Context, keeper Keeper, data types.GenesisState) {
	keeper.SetParams(ctx, data.Params)

	for _, code := range data.Codes {
		newCodeID, err := keeper.StoreCode(ctx, code.CodeInfo.Creator, code.CodesBytes)
		if err != nil {
			panic(err)
		}
		newInfo, err := keeper.GetCodeInfo(ctx, newCodeID)
		if err != nil {
			panic(err)
		}

		if !bytes.Equal(code.CodeInfo.CodeHash, newInfo.CodeHash) {
			panic("code hashes not same")
		}
	}

	for _, contract := range data.Contracts {
		keeper.SetContractInfo(ctx, contract.ContractInfo.Address, contract.ContractInfo)
		keeper.SetContractStore(ctx, contract.ContractInfo.Address, contract.ContractStore)
	}

}

// ExportGenesis returns a GenesisState for a given context and keeper.
func ExportGenesis(ctx sdk.Context, keeper Keeper) types.GenesisState {
	var codes []types.Code
	var contracts []types.Contract

	maxCodeID := keeper.GetNextCodeID(ctx)
	for i := uint64(1); i < maxCodeID; i++ {
		bytecode, err := keeper.GetByteCode(ctx, i)
		if err != nil {
			panic(err)
		}

		codeInfo, err := keeper.GetCodeInfo(ctx, i)
		if err != nil {
			panic(err)
		}

		codes = append(codes, types.Code{
			CodeInfo:   codeInfo,
			CodesBytes: bytecode,
		})
	}

	keeper.IterateContractInfo(ctx, func(contract types.ContractInfo) bool {
		contractStateIterator := keeper.GetContractStoreIterator(ctx, contract.Address)
		var models []types.Model
		for ; contractStateIterator.Valid(); contractStateIterator.Next() {
			m := types.Model{
				Key:   contractStateIterator.Key(),
				Value: contractStateIterator.Value(),
			}
			models = append(models, m)
		}

		contracts = append(contracts, types.Contract{
			ContractInfo:  contract,
			ContractStore: models,
		})

		return false
	})

	params := keeper.GetParams(ctx)

	return types.NewGenesisState(params, codes, contracts)
}
