package tax

import (
	"fmt"
	"terra/types/assets"

	"github.com/cosmos/cosmos-sdk/codec"
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

// Tax related variables
var (
	taxRateMin = sdk.ZeroDec()
	taxRateMax = sdk.NewDecWithPrec(2, 2) // 2%
)

// Keeper defines a module interface that facilitates the transfer of coins
// between accounts.
type Keeper interface {
	bank.Keeper
	GetIssuance(ctx sdk.Context, denom string) sdk.Int
	GetDebtRatio(ctx sdk.Context) sdk.Dec
}

var _ Keeper = (*BaseKeeper)(nil)

// BaseKeeper manages transfers between accounts. It implements the Keeper
// interface.
type BaseKeeper struct {
	key sdk.StoreKey
	cdc *codec.Codec

	am auth.AccountKeeper
	fk auth.FeeCollectionKeeper
}

// NewBaseKeeper returns a new BaseKeeper
func NewBaseKeeper(key sdk.StoreKey, cdc *codec.Codec, am auth.AccountKeeper, fk auth.FeeCollectionKeeper) BaseKeeper {
	return BaseKeeper{
		key: key,
		cdc: cdc,
		am:  am,
		fk:  fk,
	}
}

// GetCoins returns the coins at the addr.
func (keeper BaseKeeper) GetCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins {
	return getCoins(ctx, keeper.am, addr)
}

// SetCoins sets the coins at the addr.
func (keeper BaseKeeper) SetCoins(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) sdk.Error {
	return setCoins(ctx, keeper.am, addr, amt)
}

// HasCoins returns whether or not an account has at least amt coins.
func (keeper BaseKeeper) HasCoins(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) bool {
	return hasCoins(ctx, keeper.am, addr, amt)
}

// SubtractCoins subtracts amt from the coins at the addr.
func (keeper BaseKeeper) SubtractCoins(
	ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins,
) (sdk.Coins, sdk.Tags, sdk.Error) {

	return subtractCoins(ctx, keeper, addr, amt)
}

// AddCoins adds amt to the coins at the addr.
func (keeper BaseKeeper) AddCoins(
	ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins,
) (sdk.Coins, sdk.Tags, sdk.Error) {

	return addCoins(ctx, keeper, addr, amt)
}

// SendCoins moves coins from one account to another
func (keeper BaseKeeper) SendCoins(
	ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins,
) (sdk.Tags, sdk.Error) {

	taxes := keeper.calculateTaxes(ctx, amt)
	_, taxTags, err := subtractCoins(ctx, keeper, fromAddr, taxes)
	if err != nil {
		return nil, err
	}

	_, subTags, err := subtractCoins(ctx, keeper, fromAddr, amt)
	if err != nil {
		return nil, err
	}

	_, addTags, err := addCoins(ctx, keeper, toAddr, amt)
	if err != nil {
		return nil, err
	}

	taxTags = taxTags.AppendTags(subTags)

	return taxTags.AppendTags(addTags), nil
}

// InputOutputCoins handles a list of inputs and outputs
func (keeper BaseKeeper) InputOutputCoins(ctx sdk.Context, inputs []bank.Input, outputs []bank.Output) (sdk.Tags, sdk.Error) {
	allTags := sdk.EmptyTags()

	for _, in := range inputs {
		_, tags, err := subtractCoins(ctx, keeper, in.Address, in.Coins)
		if err != nil {
			return nil, err
		}
		allTags = allTags.AppendTags(tags)
	}

	for _, out := range outputs {

		_, tags, err := addCoins(ctx, keeper, out.Address, out.Coins)
		if err != nil {
			return nil, err
		}
		allTags = allTags.AppendTags(tags)

		taxes := keeper.calculateTaxes(ctx, out.Coins)
		_, taxTags, err := subtractCoins(ctx, keeper, out.Address, taxes)
		if err != nil {
			return nil, err
		}
		allTags = allTags.AppendTags(taxTags)
	}

	return allTags, nil
}

func (keeper BaseKeeper) calculateTaxes(ctx sdk.Context, principal sdk.Coins) sdk.Coins {
	taxes := sdk.Coins{}
	for _, coin := range principal {
		taxRate := keeper.getTaxRate(ctx, coin.Denom)
		taxDue := sdk.NewDecFromInt(coin.Amount).Mul(taxRate).RoundInt()

		taxes = append(taxes, sdk.NewCoin(coin.Denom, taxDue))
	}

	return taxes
}

// getTaxRate gets the currently effective tax rate for {denom}
func (keeper BaseKeeper) getTaxRate(ctx sdk.Context, denom string) (res sdk.Dec) {
	return taxRateMin.Add(taxRateMax.Sub(taxRateMin).Mul(keeper.GetDebtRatio(ctx)))
}

// setIssuance sets the total issuance of the coin with {denom}
func (keeper BaseKeeper) setIssuance(ctx sdk.Context, denom string, issuance sdk.Int) {
	store := ctx.KVStore(keeper.key)
	bz := keeper.cdc.MustMarshalBinaryLengthPrefixed(issuance)
	store.Set(GetCoinSupplyKey(denom), bz)
}

// GetIssuance retrieves the total issuance of the coin with {denom}
func (keeper BaseKeeper) GetIssuance(ctx sdk.Context, denom string) (res sdk.Int) {
	store := ctx.KVStore(keeper.key)
	bz := store.Get(GetCoinSupplyKey(denom))
	if bz == nil {
		res = sdk.ZeroInt()
		return
	}
	keeper.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &res)
	return
}

// GetDebtRatio gets the current debt of the system
func (keeper BaseKeeper) GetDebtRatio(ctx sdk.Context) sdk.Dec {
	lunaCurrentIssuance := keeper.GetIssuance(ctx, assets.LunaDenom)
	lunaTargetIssuance := sdk.NewInt(assets.LunaTargetIssuance)

	lunaDebt := lunaCurrentIssuance.Sub(lunaTargetIssuance)

	return sdk.NewDecFromInt(lunaDebt).Quo(sdk.NewDecFromInt(lunaCurrentIssuance))
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
func subtractCoins(ctx sdk.Context, keeper BaseKeeper, addr sdk.AccAddress, amt sdk.Coins) (sdk.Coins, sdk.Tags, sdk.Error) {
	ctx.GasMeter().ConsumeGas(costSubtractCoins, "subtractCoins")
	oldCoins := getCoins(ctx, keeper.am, addr)
	newCoins, hasNeg := oldCoins.SafeMinus(amt)
	if hasNeg {
		return amt, nil, sdk.ErrInsufficientCoins(fmt.Sprintf("%s : %s", oldCoins, amt))
	}

	// Update issuance
	for _, coin := range amt {
		issuance := keeper.GetIssuance(ctx, coin.Denom)
		issuance = issuance.Sub(coin.Amount)
		keeper.setIssuance(ctx, coin.Denom, issuance)
	}

	err := setCoins(ctx, keeper.am, addr, newCoins)
	tags := sdk.NewTags("sender", []byte(addr.String()))
	return newCoins, tags, err
}

// AddCoins adds amt to the coins at the addr.
func addCoins(ctx sdk.Context, keeper BaseKeeper, addr sdk.AccAddress, amt sdk.Coins) (sdk.Coins, sdk.Tags, sdk.Error) {
	ctx.GasMeter().ConsumeGas(costAddCoins, "addCoins")
	oldCoins := getCoins(ctx, keeper.am, addr)
	newCoins := oldCoins.Plus(amt)

	// Update issuance
	for _, coin := range amt {
		issuance := keeper.GetIssuance(ctx, coin.Denom)
		issuance = issuance.Add(coin.Amount)
		keeper.setIssuance(ctx, coin.Denom, issuance)
	}

	if !newCoins.IsNotNegative() {
		return amt, nil, sdk.ErrInsufficientCoins(fmt.Sprintf("%s < %s", oldCoins, amt))
	}
	err := setCoins(ctx, keeper.am, addr, newCoins)
	tags := sdk.NewTags("recipient", []byte(addr.String()))
	return newCoins, tags, err
}
