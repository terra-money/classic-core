package types

import (
	"fmt"
	"terra/x/treasury"

	"github.com/cosmos/cosmos-sdk/x/bank"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
)

const (
	costGetCoins      sdk.Gas = 10
	costHasCoins      sdk.Gas = 10
	costSetCoins      sdk.Gas = 100
	costSubtractCoins sdk.Gas = 10
	costAddCoins      sdk.Gas = 10
)

// TaxKeeper defines a module interface that facilitates the transfer of coins
// between accounts.
type TaxKeeper interface {
	bank.Keeper
	// SetCoins(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) sdk.Error
	// SubtractCoins(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) (sdk.Coins, sdk.Tags, sdk.Error)
	// AddCoins(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) (sdk.Coins, sdk.Tags, sdk.Error)
}

//var _ TaxKeeper = (*BaseTaxKeeper)(nil)

// BaseTaxKeeper manages transfers between accounts. It implements the Keeper
// interface.
type BaseTaxKeeper struct {
	am auth.AccountKeeper
	tk treasury.Keeper
	fk auth.FeeCollectionKeeper
}

// NewBaseTaxKeeper returns a new BaseKeeper
func NewBaseTaxKeeper(am auth.AccountKeeper, tk treasury.Keeper, fk auth.FeeCollectionKeeper) BaseTaxKeeper {
	return BaseTaxKeeper{am: am, tk: tk, fk: fk}
}

// // GetCoins returns the coins at the addr.
// func (keeper BaseTaxKeeper) GetCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins {
// 	return getCoins(ctx, keeper.am, addr)
// }

// // SetCoins sets the coins at the addr.
// func (keeper BaseTaxKeeper) SetCoins(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) sdk.Error {
// 	return setCoins(ctx, keeper.am, addr, amt)
// }

// // HasCoins returns whether or not an account has at least amt coins.
// func (keeper BaseTaxKeeper) HasCoins(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) bool {
// 	return hasCoins(ctx, keeper.am, addr, amt)
// }

// // SubtractCoins subtracts amt from the coins at the addr.
// func (keeper BaseTaxKeeper) SubtractCoins(
// 	ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins,
// ) (sdk.Coins, sdk.Tags, sdk.Error) {

// 	return subtractCoins(ctx, keeper.am, addr, amt)
// }

// // AddCoins adds amt to the coins at the addr.
// func (keeper BaseTaxKeeper) AddCoins(
// 	ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins,
// ) (sdk.Coins, sdk.Tags, sdk.Error) {

// 	return addCoins(ctx, keeper.am, addr, amt)
// }

// SendCoins moves coins from one account to another
func (keeper BaseTaxKeeper) SendCoins(
	ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins,
) (sdk.Tags, sdk.Error) {

	taxes := sdk.Coins{}
	for _, coin := range amt {
		taxRate := keeper.tk.GetTaxRate(ctx, coin.Denom)
		taxDue := sdk.NewDecFromInt(coin.Amount).Mul(taxRate).RoundInt()

		taxes = append(taxes, sdk.NewCoin(coin.Denom, taxDue))
	}

	_, taxTags, err := subtractCoins(ctx, keeper.am, fromAddr, taxes)
	if err != nil {
		return nil, err
	}

	_, subTags, err := subtractCoins(ctx, keeper.am, fromAddr, amt)
	if err != nil {
		return nil, err
	}

	_, addTags, err := addCoins(ctx, keeper.am, toAddr, amt)
	if err != nil {
		return nil, err
	}

	taxTags = taxTags.AppendTags(subTags)

	return taxTags.AppendTags(addTags), nil
}

// InputOutputCoins handles a list of inputs and outputs
func (keeper BaseTaxKeeper) InputOutputCoins(ctx sdk.Context, inputs []bank.Input, outputs []bank.Output) (sdk.Tags, sdk.Error) {
	allTags := sdk.EmptyTags()

	for _, in := range inputs {
		_, tags, err := subtractCoins(ctx, keeper.am, in.Address, in.Coins)
		if err != nil {
			return nil, err
		}
		allTags = allTags.AppendTags(tags)
	}

	for _, out := range outputs {

		taxes := sdk.Coins{}
		for _, coin := range out.Coins {
			taxRate := keeper.tk.GetTaxRate(ctx, coin.Denom)
			taxDue := sdk.NewDecFromInt(coin.Amount).Mul(taxRate).RoundInt()

			taxes = append(taxes, sdk.NewCoin(coin.Denom, taxDue))
		}

		_, taxTags, err := subtractCoins(ctx, keeper.am, out.Address, taxes)
		if err != nil {
			return nil, err
		}
		allTags = allTags.AppendTags(taxTags)

		_, tags, err := addCoins(ctx, keeper.am, out.Address, out.Coins)
		if err != nil {
			return nil, err
		}
		allTags = allTags.AppendTags(tags)
	}

	return allTags, nil
}

func getCoins(ctx sdk.Context, am auth.AccountKeeper, addr sdk.AccAddress) sdk.Coins {
	ctx.GasMeter().ConsumeGas(costGetCoins, "getCoins")
	acc := am.GetAccount(ctx, addr)
	if acc == nil {
		return sdk.Coins{}
	}
	return acc.GetCoins()
}

func setCoins(ctx sdk.Context, am auth.AccountKeeper, addr sdk.AccAddress, amt sdk.Coins) sdk.Error {
	ctx.GasMeter().ConsumeGas(costSetCoins, "setCoins")
	acc := am.GetAccount(ctx, addr)
	if acc == nil {
		acc = am.NewAccountWithAddress(ctx, addr)
	}
	err := acc.SetCoins(amt)
	if err != nil {
		// Handle w/ #870
		panic(err)
	}
	am.SetAccount(ctx, acc)
	return nil
}

// HasCoins returns whether or not an account has at least amt coins.
func hasCoins(ctx sdk.Context, am auth.AccountKeeper, addr sdk.AccAddress, amt sdk.Coins) bool {
	ctx.GasMeter().ConsumeGas(costHasCoins, "hasCoins")
	return getCoins(ctx, am, addr).IsAllGTE(amt)
}

// SubtractCoins subtracts amt from the coins at the addr.
func subtractCoins(ctx sdk.Context, am auth.AccountKeeper, addr sdk.AccAddress, amt sdk.Coins) (sdk.Coins, sdk.Tags, sdk.Error) {
	ctx.GasMeter().ConsumeGas(costSubtractCoins, "subtractCoins")
	oldCoins := getCoins(ctx, am, addr)
	newCoins := oldCoins.Minus(amt)
	if !newCoins.IsNotNegative() {
		return amt, nil, sdk.ErrInsufficientCoins(fmt.Sprintf("%s < %s", oldCoins, amt))
	}
	err := setCoins(ctx, am, addr, newCoins)
	tags := sdk.NewTags("sender", []byte(addr.String()))
	return newCoins, tags, err
}

// AddCoins adds amt to the coins at the addr.
func addCoins(ctx sdk.Context, am auth.AccountKeeper, addr sdk.AccAddress, amt sdk.Coins) (sdk.Coins, sdk.Tags, sdk.Error) {
	ctx.GasMeter().ConsumeGas(costAddCoins, "addCoins")
	oldCoins := getCoins(ctx, am, addr)
	newCoins := oldCoins.Plus(amt)
	if !newCoins.IsNotNegative() {
		return amt, nil, sdk.ErrInsufficientCoins(fmt.Sprintf("%s < %s", oldCoins, amt))
	}
	err := setCoins(ctx, am, addr, newCoins)
	tags := sdk.NewTags("recipient", []byte(addr.String()))
	return newCoins, tags, err
}
