package keeper

import (
	"encoding/binary"
	"fmt"
	"path/filepath"

	wasm "github.com/CosmWasm/go-cosmwasm"
	"github.com/spf13/viper"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/terra-money/core/x/params"
	"github.com/terra-money/core/x/wasm/config"
	"github.com/terra-money/core/x/wasm/internal/types"
)

// Keeper will have a reference to Wasmer with it's own data directory.
type Keeper struct {
	cdc        *codec.Codec
	storeKey   sdk.StoreKey
	paramSpace params.Subspace

	accountKeeper  types.AccountKeeper
	bankKeeper     types.BankKeeper
	supplyKeeper   types.SupplyKeeper
	treasuryKeeper types.TreasuryKeeper

	router sdk.Router

	wasmer    wasm.Wasmer
	querier   types.Querier
	msgParser types.MsgParser

	// WASM config values
	wasmConfig       *config.Config
	loggingWhitelist map[string]bool
}

// NewKeeper creates a new contract Keeper instance
func NewKeeper(cdc *codec.Codec, storeKey sdk.StoreKey,
	paramspace params.Subspace, accountKeeper types.AccountKeeper,
	bankKeeper types.BankKeeper, supplyKeeper types.SupplyKeeper, treasuryKeeper types.TreasuryKeeper, router sdk.Router,
	supportedFeatures string,
	wasmConfig *config.Config) Keeper {
	homeDir := viper.GetString(flags.FlagHome)
	wasmer, err := wasm.NewWasmer(filepath.Join(homeDir, config.DBDir), supportedFeatures, 0)

	if err != nil {
		panic(err)
	}

	// set KeyTable if it has not already been set
	if !paramspace.HasKeyTable() {
		paramspace = paramspace.WithKeyTable(types.ParamKeyTable())
	}

	return Keeper{
		storeKey:         storeKey,
		cdc:              cdc,
		paramSpace:       paramspace,
		wasmer:           *wasmer,
		accountKeeper:    accountKeeper,
		bankKeeper:       bankKeeper,
		supplyKeeper:     supplyKeeper,
		treasuryKeeper:   treasuryKeeper,
		router:           router,
		wasmConfig:       wasmConfig,
		loggingWhitelist: wasmConfig.WhitelistToMap(),
		msgParser:        types.NewModuleMsgParser(),
		querier:          types.NewModuleQuerier(),
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// StoreConfig store wasm config to local config file
func (k Keeper) StoreConfig() {
	rootDir := viper.GetString(flags.FlagHome)
	wasmConfigFilePath := filepath.Join(rootDir, "config/wasm.toml")

	config.WriteConfigFile(wasmConfigFilePath, k.wasmConfig)
}

// GetLastCodeID return last code ID
func (k Keeper) GetLastCodeID(ctx sdk.Context) (uint64, error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.LastCodeIDKey)
	if bz == nil {
		return 0, sdkerrors.Wrap(types.ErrInvalidGenesis, "initial code ID hasn't been set")
	}

	return binary.BigEndian.Uint64(bz), nil
}

// SetLastCodeID set last code id
func (k Keeper) SetLastCodeID(ctx sdk.Context, id uint64) {
	store := ctx.KVStore(k.storeKey)
	bz := sdk.Uint64ToBigEndian(id)
	store.Set(types.LastCodeIDKey, bz)
}

// GetLastInstanceID return last instance ID
func (k Keeper) GetLastInstanceID(ctx sdk.Context) (uint64, error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.LastInstanceIDKey)
	if bz == nil {
		return 0, sdkerrors.Wrap(types.ErrInvalidGenesis, "initial instance ID hasn't been set")
	}

	return binary.BigEndian.Uint64(bz), nil
}

// SetLastInstanceID set last instance id
func (k Keeper) SetLastInstanceID(ctx sdk.Context, id uint64) {
	store := ctx.KVStore(k.storeKey)
	bz := sdk.Uint64ToBigEndian(id)
	store.Set(types.LastInstanceIDKey, bz)
}

// GetCodeInfo returns CodeInfo for the given codeID
func (k Keeper) GetCodeInfo(ctx sdk.Context, codeID uint64) (codeInfo types.CodeInfo, err error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetCodeInfoKey(codeID))
	if bz == nil {
		return types.CodeInfo{}, sdkerrors.Wrapf(types.ErrNotFound, "codeID %d", codeID)
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &codeInfo)
	return
}

// SetCodeInfo stores CodeInfo for the given codeID
func (k Keeper) SetCodeInfo(ctx sdk.Context, codeID uint64, codeInfo types.CodeInfo) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(codeInfo)
	store.Set(types.GetCodeInfoKey(codeID), bz)
}

// GetContractInfo returns contract info of the given address
func (k Keeper) GetContractInfo(ctx sdk.Context, contractAddress sdk.AccAddress) (contractInfo types.ContractInfo, err error) {
	store := ctx.KVStore(k.storeKey)
	contractBz := store.Get(types.GetContractInfoKey(contractAddress))
	if contractBz == nil {
		return types.ContractInfo{}, sdkerrors.Wrapf(types.ErrNotFound, "constractInfo %s", contractAddress.String())
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(contractBz, &contractInfo)
	return contractInfo, nil
}

// SetContractInfo stores ContractInfo for the given contractAddress
func (k Keeper) SetContractInfo(ctx sdk.Context, contractAddress sdk.AccAddress, codeInfo types.ContractInfo) {
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshalBinaryLengthPrefixed(codeInfo)
	store.Set(types.GetContractInfoKey(contractAddress), b)
}

// IterateContractInfo iterates all contract infos
func (k Keeper) IterateContractInfo(ctx sdk.Context, cb func(types.ContractInfo) bool) {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.ContractInfoKey)
	iter := prefixStore.Iterator(nil, nil)
	for ; iter.Valid(); iter.Next() {
		var contract types.ContractInfo
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iter.Value(), &contract)
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
func (k Keeper) GetByteCode(ctx sdk.Context, codeID uint64) ([]byte, error) {
	codeInfo, sdkErr := k.GetCodeInfo(ctx, codeID)
	if sdkErr != nil {
		return nil, sdkErr
	}

	byteCode, err := k.wasmer.GetCode(codeInfo.CodeHash.Bytes())
	if err != nil {
		return nil, err
	}
	return byteCode, nil
}

// RegisterMsgParsers register module msg parsers
func (k *Keeper) RegisterMsgParsers(parsers map[string]types.WasmMsgParserInterface) {
	for route, parser := range parsers {
		k.msgParser[route] = parser
	}
}

// RegisterQueriers register module queriers
func (k *Keeper) RegisterQueriers(queriers map[string]types.WasmQuerierInterface) {
	for route, querier := range queriers {
		k.querier.Queriers[route] = querier
	}
}
