package wasm

import (
	"bytes"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/terra-money/core/x/wasm/internal/types"
	// authexported "github.com/cosmos/cosmos-sdk/x/auth/exported"
	// "github.com/terra-money/core/x/wasm/internal/types"
)

// InitGenesis sets wasm information for genesis.
//
// CONTRACT: all types of accounts must have been already initialized/created
func InitGenesis(ctx sdk.Context, keeper Keeper, data types.GenesisState) {
	keeper.SetParams(ctx, data.Params)
	keeper.SetLastCodeID(ctx, data.LastCodeID)
	keeper.SetLastInstanceID(ctx, data.LastInstanceID)

	for _, code := range data.Codes {
		codeHash, err := keeper.CompileCode(ctx, code.CodesBytes)
		if err != nil {
			panic(err)
		}

		if !bytes.Equal(codeHash, code.CodeInfo.CodeHash.Bytes()) {
			panic("CodeHash is not same")
		}

		keeper.SetCodeInfo(ctx, code.CodeInfo.CodeID, code.CodeInfo)
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

	lastCodeID, err := keeper.GetLastCodeID(ctx)
	if err != nil {
		panic(err)
	}

	lastInstanceID, err := keeper.GetLastInstanceID(ctx)
	if err != nil {
		panic(err)
	}

	for i := uint64(1); i <= lastCodeID; i++ {
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

	return types.NewGenesisState(params, lastCodeID, lastInstanceID, codes, contracts)
}
