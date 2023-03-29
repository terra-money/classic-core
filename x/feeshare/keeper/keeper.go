package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/tendermint/tendermint/libs/log"

	revtypes "github.com/classic-terra/core/x/feeshare/types"
	wasmkeeper "github.com/classic-terra/core/x/wasm/keeper"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// Keeper of this module maintains collections of feeshares for contracts
// registered to receive transaction fees.
type Keeper struct {
	storeKey   storetypes.StoreKey
	cdc        codec.BinaryCodec
	paramstore paramtypes.Subspace

	bankKeeper revtypes.BankKeeper
	wasmKeeper wasmkeeper.Keeper

	feeCollectorName string
}

// NewKeeper creates new instances of the fees Keeper
func NewKeeper(
	storeKey storetypes.StoreKey,
	cdc codec.BinaryCodec,
	ps paramtypes.Subspace,
	bk revtypes.BankKeeper,
	wk wasmkeeper.Keeper,
	feeCollector string,
) Keeper {
	// set KeyTable if it has not already been set
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(revtypes.ParamKeyTable())
	}

	return Keeper{
		storeKey:         storeKey,
		cdc:              cdc,
		paramstore:       ps,
		bankKeeper:       bk,
		wasmKeeper:       wk,
		feeCollectorName: feeCollector,
	}
}

// SendCoinsFromAccountToFeeCollector transfers amt to the fee collector account.
func (k Keeper) SendCoinsFromAccountToFeeCollector(ctx sdk.Context, senderAddr sdk.AccAddress, amt sdk.Coins) error {
	if senderAddr.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "senderAddr address cannot be empty")
	}
	return k.bankKeeper.SendCoinsFromAccountToModule(ctx, senderAddr, k.feeCollectorName, amt)
}

// SendCoinsFromFeeCollectorToAccount transfers amt from the fee collector account to the recipient.
func (k Keeper) SendCoinsFromFeeCollectorToAccount(ctx sdk.Context, recipientAddr sdk.AccAddress, amt sdk.Coins) error {
	if recipientAddr.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "recipient address cannot be empty")
	}
	return k.bankKeeper.SendCoinsFromModuleToAccount(ctx, k.feeCollectorName, recipientAddr, amt)
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", revtypes.ModuleName))
}
