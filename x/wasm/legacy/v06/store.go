package v06

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"

	v05 "github.com/terra-money/core/x/wasm/legacy/v05"
	v06 "github.com/terra-money/core/x/wasm/types"
)

// migrateContractInfo migrate old ContractInfo to include
// empty IBCPort.
func migrateContractInfo(store sdk.KVStore, prefixBz []byte, cdc codec.BinaryCodec) error {
	prefixStore := prefix.NewStore(store, prefixBz)
	iter := prefixStore.Iterator(nil, nil)
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		var oldContractInfo v05.ContractInfo
		if err := cdc.Unmarshal(iter.Value(), &oldContractInfo); err != nil {
			return err
		}

		contractInfo := v06.ContractInfo{
			Address:   oldContractInfo.Address,
			Creator:   oldContractInfo.Creator,
			Admin:     oldContractInfo.Admin,
			CodeID:    oldContractInfo.CodeID,
			InitMsg:   oldContractInfo.InitMsg,
			IBCPortID: "",
		}

		bz, err := cdc.Marshal(&contractInfo)
		if err != nil {
			return err
		}

		prefixStore.Set(iter.Key(), bz)
	}

	return nil
}

// MigrateStore performs in-place store migrations from v0.5 to v0.6. The
// migration includes:
//
// - Add empty IBCPort to all ContractInfos.
func MigrateStore(ctx sdk.Context, storeKey sdk.StoreKey, cdc codec.BinaryCodec) error {
	store := ctx.KVStore(storeKey)
	return migrateContractInfo(store, v05.ContractInfoKey, cdc)
}
