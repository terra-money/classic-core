package mint

import (
	"github.com/terra-project/core/types/assets"
	"github.com/terra-project/core/types/util"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/staking"
)

// StoreKey is string representation of the store key for mint
const StoreKey = "mint"

// Keeper is an instance of the Mint keeper module.
// Adds / subtracts balances from accounts and maintains a global state
// of issuance of currencies on the Terra network.
type Keeper struct {
	cdc *codec.Codec
	key sdk.StoreKey
	sk  staking.Keeper
	bk  bank.Keeper
	ak  auth.AccountKeeper
}

// NewKeeper creates a new instance of the mint module.
func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, sk staking.Keeper, bk bank.Keeper, ak auth.AccountKeeper) Keeper {
	return Keeper{
		cdc: cdc,
		key: key,
		sk:  sk,
		bk:  bk,
		ak:  ak,
	}
}

// Mint credits {coin} to the {recipient} account, and reflects the increase in issuance
func (k Keeper) Mint(ctx sdk.Context, recipient sdk.AccAddress, coin sdk.Coin) (err sdk.Error) {

	_, _, err = k.bk.AddCoins(ctx, recipient, sdk.Coins{coin})
	if err != nil {
		return err
	}

	if coin.Denom == assets.MicroLunaDenom {
		pool := k.sk.GetPool(ctx)
		pool.NotBondedTokens = pool.NotBondedTokens.Add(coin.Amount)
		k.sk.SetPool(ctx, pool)
	}

	return k.ChangeIssuance(ctx, coin.Denom, coin.Amount)
}

// Burn deducts {coin} from the {payer} account, and reflects the decrease in issuance
func (k Keeper) Burn(ctx sdk.Context, payer sdk.AccAddress, coin sdk.Coin) (err sdk.Error) {
	_, _, err = k.bk.SubtractCoins(ctx, payer, sdk.Coins{coin})
	if err != nil {
		return err
	}

	if coin.Denom == assets.MicroLunaDenom {
		pool := k.sk.GetPool(ctx)
		pool.NotBondedTokens = pool.NotBondedTokens.Sub(coin.Amount)
		k.sk.SetPool(ctx, pool)
	}

	return k.ChangeIssuance(ctx, coin.Denom, coin.Amount.Neg())
}

// ChangeIssuance updates the issuance to reflect
func (k Keeper) ChangeIssuance(ctx sdk.Context, denom string, delta sdk.Int) (err sdk.Error) {
	store := ctx.KVStore(k.key)
	curDay := sdk.NewInt(ctx.BlockHeight() / util.BlocksPerDay)

	// If genesis issuance is not on disk, GetIssuance will do a fresh read of account balances
	// and the change in issuance should be reported automatically.
	if !store.Has(keyIssuance(denom, sdk.ZeroInt())) {
		k.GetIssuance(ctx, denom, curDay)
		return
	}

	curIssuance := k.GetIssuance(ctx, denom, curDay)
	newIssuance := curIssuance.Add(delta)

	if newIssuance.IsNegative() {
		err = sdk.ErrInternal("Issuance should never fall below 0")
	} else {
		bz := k.cdc.MustMarshalBinaryLengthPrefixed(newIssuance)
		store.Set(keyIssuance(denom, curDay), bz)
	}

	return
}

// GetIssuance fetches the total issuance count of the coin matching {denom}. If the {day} applies
// to a previous period, fetches the last stored snapshot issuance of the coin. For virgin calls,
// iterates through the accountkeeper and computes the genesis issuance.
func (k Keeper) GetIssuance(ctx sdk.Context, denom string, day sdk.Int) (issuance sdk.Int) {
	store := ctx.KVStore(k.key)

	if bz := store.Get(keyIssuance(denom, day)); bz != nil {
		k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &issuance)
	} else {
		// Genesis epoch; nothing exists in store so we must read it
		// from accountkeeper
		if day.LTE(sdk.ZeroInt()) {
			issuance = sdk.ZeroInt()
			countIssuance := func(acc auth.Account) (stop bool) {
				issuance = issuance.Add(acc.GetCoins().AmountOf(denom))
				return false
			}
			k.ak.IterateAccounts(ctx, countIssuance)
		} else {
			// Fetch the issuance snapshot of the previous epoch
			issuance = k.GetIssuance(ctx, denom, day.Sub(sdk.OneInt()))
		}

		// Set issuance to the store
		store := ctx.KVStore(k.key)
		bz := k.cdc.MustMarshalBinaryLengthPrefixed(issuance)
		store.Set(keyIssuance(denom, day), bz)
	}

	return
}

// PeekEpochSeigniorage retrieves the size of the seigniorage pool at epoch
func (k Keeper) PeekEpochSeigniorage(ctx sdk.Context, epoch sdk.Int) (epochSeigniorage sdk.Int) {

	daysPerEpoch := util.BlocksPerEpoch / util.BlocksPerDay
	epochLastDay := epoch.Add(sdk.OneInt()).MulRaw(daysPerEpoch).Sub(sdk.OneInt())

	//fmt.Println(epochLastDay)
	today := sdk.NewInt(ctx.BlockHeight() / util.BlocksPerDay)
	if epochLastDay.GT(today) {
		epochLastDay = today
	}

	prevEpochLastDay := epochLastDay.SubRaw(daysPerEpoch)
	if prevEpochLastDay.IsNegative() {
		prevEpochLastDay = sdk.ZeroInt()
	}

	prevEpochIssuance := k.GetIssuance(ctx, assets.MicroLunaDenom, prevEpochLastDay)
	epochIssuance := k.GetIssuance(ctx, assets.MicroLunaDenom, epochLastDay)
	epochSeigniorage = prevEpochIssuance.Sub(epochIssuance)

	if epochSeigniorage.LT(sdk.ZeroInt()) {
		epochSeigniorage = sdk.ZeroInt()
	}

	return
}
