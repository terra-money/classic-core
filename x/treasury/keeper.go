package treasury

import (
	"terra/types/assets"
	"terra/types/tax"
	"terra/types/util"

	"github.com/cosmos/cosmos-sdk/x/auth"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Keeper of the treasury store
type Keeper struct {
	key sdk.StoreKey
	cdc *codec.Codec

	tk tax.Keeper
	fk auth.FeeCollectionKeeper
}

// NewKeeper constructs a new keeper
func NewKeeper(key sdk.StoreKey, cdc *codec.Codec,
	taxKeeper tax.Keeper, fk auth.FeeCollectionKeeper) Keeper {
	return Keeper{
		key: key,
		cdc: cdc,
		tk:  taxKeeper,
		fk:  fk,
	}
}

func (keeper Keeper) PayTax(ctx sdk.Context, revenue sdk.Coins) {
	keeper.fk.AddCollectedFees(ctx, revenue)
}

func (keeper Keeper) PayMintIncome(ctx sdk.Context, revenue sdk.Coins) {
	if revenue[0].Denom == assets.LunaDenom {
		keeper.deposit(ctx, revenue)
	}
	// Burn otherwise
}

func (keeper Keeper) AddClaim(ctx sdk.Context, claim Claim) {
	prevClaim := util.Get(
		keeper.key,
		keeper.cdc,
		ctx,
		GetClaimKey(claim.Account),
	)

	if prevClaim != nil {
		prevClaim := prevClaim.(Claim)
		claim.Weight = prevClaim.Weight.Add(claim.Weight)
	}

	util.Set(
		keeper.key,
		keeper.cdc,
		ctx,
		GetClaimKey(claim.Account),
		claim,
	)
}

//---------------------------------------
//---------------------------------------
//---------------------------------------

func (keeper Keeper) deposit(ctx sdk.Context, funds sdk.Coins) {
	incomePool := util.Get(
		keeper.key,
		keeper.cdc,
		ctx,
		KeyIncomePool,
	).(sdk.Coins)

	incomePool = incomePool.Plus(funds)

	util.Set(
		keeper.key,
		keeper.cdc,
		ctx,
		KeyIncomePool,
		incomePool,
	)
}

func (keeper Keeper) withdraw(ctx sdk.Context, funds sdk.Coins) {
	incomePool := util.Get(
		keeper.key,
		keeper.cdc,
		ctx,
		KeyIncomePool,
	).(sdk.Coins)

	incomePool = incomePool.Minus(funds)

	util.Set(
		keeper.key,
		keeper.cdc,
		ctx,
		KeyIncomePool,
		incomePool,
	)
}

func (keeper Keeper) getClaims(ctx sdk.Context) (res []Claim) {
	claims := util.Collect(
		keeper.key,
		keeper.cdc,
		ctx,
		PrefixClaim,
	)

	for _, c := range claims {
		res = append(res, c.(Claim))
	}
	return
}

func (keeper Keeper) clearClaims(ctx sdk.Context) {
	util.Clear(
		keeper.key,
		ctx,
		PrefixClaim,
	)
}
