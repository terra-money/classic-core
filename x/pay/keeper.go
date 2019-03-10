package pay

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
)

// Keeper manages transfers between accounts. It implements the Keeper interface.
type Keeper struct {
	bank.Keeper

	key sdk.StoreKey
	cdc *codec.Codec

	ak auth.AccountKeeper
	fk auth.FeeCollectionKeeper
}

// NewKeeper returns a new Keeper
func NewKeeper(
	key sdk.StoreKey,
	cdc *codec.Codec,
	ak auth.AccountKeeper,
	fk auth.FeeCollectionKeeper) Keeper {
	return Keeper{
		key: key,
		cdc: cdc,
		ak:  ak,
		fk:  fk,
	}
}

// SubtractCoins subtracts amt from the coins at the addr.
func (keeper Keeper) SubtractCoins(
	ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins,
) (sdk.Coins, sdk.Tags, sdk.Error) {
	coins, tags, err := subtractCoins(ctx, keeper.ak, addr, amt)
	if err != nil {
		keeper.subtractIssuance(ctx, amt)
	}

	return coins, tags, nil
}

// AddCoins adds amt to the coins at the addr.
func (keeper Keeper) AddCoins(
	ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins,
) (sdk.Coins, sdk.Tags, sdk.Error) {
	coins, tags, err := addCoins(ctx, keeper.ak, addr, amt)
	if err != nil {
		keeper.addIssuance(ctx, amt)
	}

	return coins, tags, nil
}

// SendCoins moves coins from one account to another
func (keeper Keeper) SendCoins(
	ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins,
) (sdk.Tags, sdk.Error) {

	_, taxTags, taxErr := keeper.payTax(ctx, fromAddr, amt)
	if taxErr != nil {
		return nil, taxErr
	}

	_, subTags, err := subtractCoins(ctx, keeper.ak, fromAddr, amt)
	if err != nil {
		return nil, err
	}

	_, addTags, err := addCoins(ctx, keeper.ak, toAddr, amt)
	if err != nil {
		return nil, err
	}

	taxTags = taxTags.AppendTags(subTags)

	return taxTags.AppendTags(addTags), nil
}

// InputOutputCoins handles a list of inputs and outputs
func (keeper Keeper) InputOutputCoins(ctx sdk.Context, inputs []bank.Input, outputs []bank.Output) (sdk.Tags, sdk.Error) {
	allTags := sdk.EmptyTags()

	for _, in := range inputs {
		_, tags, err := subtractCoins(ctx, keeper.ak, in.Address, in.Coins)
		if err != nil {
			return nil, err
		}
		allTags = allTags.AppendTags(tags)
	}

	for _, out := range outputs {
		_, taxTags, taxErr := keeper.payTax(ctx, out.Address, out.Coins)
		if taxErr != nil {
			return nil, taxErr
		}
		allTags = allTags.AppendTags(taxTags)

		_, tags, err := addCoins(ctx, keeper.ak, out.Address, out.Coins)
		if err != nil {
			return nil, err
		}
		allTags = allTags.AppendTags(tags)
	}

	return allTags, nil
}

func getCoins(ctx sdk.Context, ak auth.AccountKeeper, addr sdk.AccAddress) sdk.Coins {
	acc := ak.GetAccount(ctx, addr)
	if acc == nil {
		return sdk.Coins{}
	}
	return acc.GetCoins()
}

func setCoins(ctx sdk.Context, ak auth.AccountKeeper, addr sdk.AccAddress, amt sdk.Coins) sdk.Error {
	acc := ak.GetAccount(ctx, addr)
	if acc == nil {
		acc = ak.NewAccountWithAddress(ctx, addr)
	}
	err := acc.SetCoins(amt)
	if err != nil {
		// Handle w/ #870
		panic(err)
	}
	ak.SetAccount(ctx, acc)
	return nil
}

// HasCoins returns whether or not an account has at least amt coins.
func hasCoins(ctx sdk.Context, k Keeper, addr sdk.AccAddress, amt sdk.Coins) bool {
	return getCoins(ctx, k.ak, addr).IsAllGTE(amt)
}

func getAccount(ctx sdk.Context, ak auth.AccountKeeper, addr sdk.AccAddress) auth.Account {
	return ak.GetAccount(ctx, addr)
}

func setAccount(ctx sdk.Context, k Keeper, acc auth.Account) {
	k.ak.SetAccount(ctx, acc)
}

// subtractCoins subtracts amt coins from an account with the given address addr.
//
// CONTRACT: If the account is a vesting account, the amount has to be spendable.
func subtractCoins(ctx sdk.Context, ak auth.AccountKeeper, addr sdk.AccAddress, amt sdk.Coins) (sdk.Coins, sdk.Tags, sdk.Error) {
	oldCoins, spendableCoins := sdk.Coins{}, sdk.Coins{}

	acc := getAccount(ctx, ak, addr)
	if acc != nil {
		oldCoins = acc.GetCoins()
		spendableCoins = acc.SpendableCoins(ctx.BlockHeader().Time)
	}

	// For non-vesting accounts, spendable coins will simply be the original coins.
	// So the check here is sufficient instead of subtracting from oldCoins.
	_, hasNeg := spendableCoins.SafeMinus(amt)
	if hasNeg {
		return amt, nil, sdk.ErrInsufficientCoins(
			fmt.Sprintf("insufficient account funds; %s < %s", spendableCoins, amt),
		)
	}

	newCoins := oldCoins.Minus(amt) // should not panic as spendable coins was already checked
	err := setCoins(ctx, ak, addr, newCoins)
	tags := sdk.NewTags(bank.TagKeySender, addr.String())

	return newCoins, tags, err
}

// AddCoins adds amt to the coins at the addr.
func addCoins(ctx sdk.Context, am auth.AccountKeeper, addr sdk.AccAddress, amt sdk.Coins) (sdk.Coins, sdk.Tags, sdk.Error) {
	oldCoins := getCoins(ctx, am, addr)
	newCoins := oldCoins.Plus(amt)

	if newCoins.IsAnyNegative() {
		return amt, nil, sdk.ErrInsufficientCoins(
			fmt.Sprintf("insufficient account funds; %s < %s", oldCoins, amt),
		)
	}

	err := setCoins(ctx, am, addr, newCoins)
	tags := sdk.NewTags(bank.TagKeyRecipient, addr.String())

	return newCoins, tags, err
}

// SendCoins moves coins from one account to another
func sendCoins(ctx sdk.Context, k Keeper, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) (sdk.Tags, sdk.Error) {
	// Safety check ensuring that when sending coins the keeper must maintain the
	// supply invariant.
	if !amt.IsValid() {
		return nil, sdk.ErrInvalidCoins(amt.String())
	}

	_, subTags, err := subtractCoins(ctx, k.ak, fromAddr, amt)
	if err != nil {
		return nil, err
	}

	_, addTags, err := addCoins(ctx, k.ak, toAddr, amt)
	if err != nil {
		return nil, err
	}

	return subTags.AppendTags(addTags), nil
}

// InputOutputCoins handles a list of inputs and outputs
// NOTE: Make sure to revert state changes from tx on error
func inputOutputCoins(ctx sdk.Context, k Keeper, inputs []bank.Input, outputs []bank.Output) (sdk.Tags, sdk.Error) {
	// Safety check ensuring that when sending coins the keeper must maintain the
	// supply invariant.
	if err := bank.ValidateInputsOutputs(inputs, outputs); err != nil {
		return nil, err
	}

	allTags := sdk.EmptyTags()

	for _, in := range inputs {
		_, tags, err := subtractCoins(ctx, k.ak, in.Address, in.Coins)
		if err != nil {
			return nil, err
		}
		allTags = allTags.AppendTags(tags)
	}

	for _, out := range outputs {
		_, tags, err := addCoins(ctx, k.ak, out.Address, out.Coins)
		if err != nil {
			return nil, err
		}
		allTags = allTags.AppendTags(tags)
	}

	return allTags, nil
}
