package wasm

import (
	"bytes"

	"github.com/terra-money/core/x/wasm/keeper"
	"github.com/terra-money/core/x/wasm/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// InitGenesis sets wasm information for genesis.
//
// CONTRACT: all types of accounts must have been already initialized/created
func InitGenesis(ctx sdk.Context, keeper keeper.Keeper, data *types.GenesisState) {
	keeper.SetParams(ctx, data.Params)
	keeper.SetLastCodeID(ctx, data.LastCodeID)
	keeper.SetLastInstanceID(ctx, data.LastInstanceID)

	for _, code := range data.Codes {
		// To cope with CosmWasm version update,
		// we temporarily allow empty code bytes
		// TODO - remove after columbus-5 update
		if len(code.CodeBytes) != 0 {
			codeHash, err := keeper.CompileCode(ctx, code.CodeBytes)
			if err != nil {
				panic(err)
			}

			if !bytes.Equal(codeHash, code.CodeInfo.CodeHash) {
				panic("CodeHash is not same")
			}
		}

		keeper.SetCodeInfo(ctx, code.CodeInfo.CodeID, code.CodeInfo)
	}

	for _, contract := range data.Contracts {
		contractAddr, err := sdk.AccAddressFromBech32(contract.ContractInfo.Address)
		if err != nil {
			panic(err)
		}

		keeper.SetContractInfo(ctx, contractAddr, contract.ContractInfo)
		keeper.SetContractStore(ctx, contractAddr, contract.ContractStore)
	}
}

// ExportGenesis returns a GenesisState for a given context and keeper.
func ExportGenesis(ctx sdk.Context, keeper keeper.Keeper) *types.GenesisState {
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
		codeInfo, err := keeper.GetCodeInfo(ctx, i)
		if err != nil {
			panic(err)
		}

		var bytecode []byte
		if len(codeInfo.CodeHash) != 0 {
			bytecode, err = keeper.GetByteCode(ctx, i)
			if err != nil {
				panic(err)
			}
		}

		codes = append(codes, types.Code{
			CodeInfo:  codeInfo,
			CodeBytes: bytecode,
		})
	}

	keeper.IterateContractInfo(ctx, func(contract types.ContractInfo) bool {
		contractAddr, err := sdk.AccAddressFromBech32(contract.Address)
		if err != nil {
			panic(err)
		}

		contractStateIterator := keeper.GetContractStoreIterator(ctx, contractAddr)
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
