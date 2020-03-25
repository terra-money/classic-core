package keeper

import (
	"fmt"
	"path/filepath"

	wasm "github.com/confio/go-cosmwasm"
	"github.com/spf13/viper"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"

	"github.com/terra-project/core/x/params"
	"github.com/terra-project/core/x/wasm/internal/types"
)

// Keeper will have a reference to Wasmer with it's own data directory.
type Keeper struct {
	cdc      *codec.Codec
	storeKey sdk.StoreKey

	paramSpace params.Subspace
	codespace  sdk.CodespaceType

	accountKeeper auth.AccountKeeper
	bankKeeper    bank.Keeper

	router sdk.Router

	wasmer        wasm.Wasmer

	queryGasLimit uint64
	cacheSize     uint64
}

// NewKeeper creates a new contract Keeper instance
func NewKeeper(cdc *codec.Codec, storeKey sdk.StoreKey, paramspace params.Subspace, accountKeeper auth.AccountKeeper, bankKeeper bank.Keeper, router sdk.Router, wasmConfig types.WasmConfig) Keeper {
	homeDir := viper.GetString(flags.FlagHome)
	wasmer, err := wasm.NewWasmer(filepath.Join(homeDir, "wasm"), 0)

	if err != nil {
		panic(err)
	}

	return Keeper{
		storeKey:      storeKey,
		cdc:           cdc,
		paramSpace:    paramspace.WithKeyTable(ParamKeyTable()),
		wasmer:        *wasmer,
		accountKeeper: accountKeeper,
		bankKeeper:    bankKeeper,
		router:        router,
		queryGasLimit: wasmConfig.ContractQueryGasLimit,
		cacheSize:     wasmConfig.CacheSize,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// GetCodeInfo returns CodeInfo for the given codeID
func (k Keeper) GetCodeInfo(ctx sdk.Context, codeID uint64) (codeInfo types.CodeInfo, err sdk.Error) {
	store := ctx.KVStore(k.storeKey)
	codeInfoBz := store.Get(types.GetCodeInfoKey(codeID))
	if codeInfoBz == nil {
		return types.CodeInfo{}, types.ErrNotFound(fmt.Sprintf("codeID %d", codeID))
	}
	k.cdc.MustUnmarshalBinaryBare(codeInfoBz, &codeInfo)
	return
}

// SetCodeInfo stores CodeInfo for the given codeID
func (k Keeper) SetCodeInfo(ctx sdk.Context, codeID uint64, codeInfo types.CodeInfo) {
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshalBinaryBare(codeInfo)
	store.Set(types.GetCodeInfoKey(codeID), b)
}

// GetContractInfo returns contract info of the given address
func (k Keeper) GetContractInfo(ctx sdk.Context, contractAddress sdk.AccAddress) (contractInfo types.ContractInfo, err sdk.Error) {
	store := ctx.KVStore(k.storeKey)
	contractBz := store.Get(types.GetContractInfoKey(contractAddress))
	if contractBz == nil {
		return types.ContractInfo{}, types.ErrNotFound(fmt.Sprintf("constractInfo %s", contractAddress.String()))
	}
	k.cdc.MustUnmarshalBinaryBare(contractBz, &contractInfo)
	return contractInfo, nil
}

// SetContractInfo stores ContractInfo for the given contractAddress
func (k Keeper) SetContractInfo(ctx sdk.Context, contractAddress sdk.AccAddress, codeInfo types.ContractInfo) {
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshalBinaryBare(codeInfo)
	store.Set(types.GetContractInfoKey(contractAddress), b)
}

// IterateContractInfo iterates all contract infos
func (k Keeper) IterateContractInfo(ctx sdk.Context, cb func(types.ContractInfo) bool) {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.ContractInfoKey)
	iter := prefixStore.Iterator(nil, nil)
	for ; iter.Valid(); iter.Next() {
		var contract types.ContractInfo
		k.cdc.MustUnmarshalBinaryBare(iter.Value(), &contract)
		// cb returns true to stop early
		if cb(contract) {
			break
		}
	}
}

// GetContractStoreIterator returns iterator for a contract store
func (k Keeper) GetContractStoreIterator(ctx sdk.Context, contractAddress sdk.AccAddress) sdk.Iterator {
	prefixStoreKey := types.GetContractStoreKey(contractAddress)
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), prefixStoreKey)
	return prefixStore.Iterator(nil, nil)
}

// SetContractStore records all the Models on the contract store
func (k Keeper) SetContractStore(ctx sdk.Context, contractAddress sdk.AccAddress, models []types.Model) {
	prefixStoreKey := types.GetContractStoreKey(contractAddress)
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), prefixStoreKey)
	for _, model := range models {
		prefixStore.Set(model.Key, model.Value)
	}
}

// GetByteCode returns ByteCode of the given CodeHash
func (k Keeper) GetByteCode(ctx sdk.Context, codeID uint64) ([]byte, sdk.Error) {
	codeInfo, sdkErr := k.GetCodeInfo(ctx, codeID)
	if sdkErr != nil {
		return nil, sdkErr
	}

	byteCode, err := k.wasmer.GetCode(codeInfo.CodeHash)
	if err != nil {
		return nil, sdk.ErrInternal(err.Error())
	}
	return byteCode, nil
}
