# Mint

The Mint module tracks the issuance of various terra currencies.

## Overview

### Mint

```golang
// Mint credits {coin} to the {recipient} account, and reflects the increase in issuance
func (k Keeper) Mint(ctx sdk.Context, recipient sdk.AccAddress, coin sdk.Coin) (err sdk.Error) {

	_, _, err = k.bk.AddCoins(ctx, recipient, sdk.Coins{coin})
	if err != nil {
		return err
	}

	...

	return k.ChangeIssuance(ctx, coin.Denom, coin.Amount)
}
```

Every time Mint is called, the issuance of the coin is incremented in the state.  

### Burn 

```golang
// Burn deducts {coin} from the {payer} account, and reflects the decrease in issuance
func (k Keeper) Burn(ctx sdk.Context, payer sdk.AccAddress, coin sdk.Coin) (err sdk.Error) {
	_, _, err = k.bk.SubtractCoins(ctx, payer, sdk.Coins{coin})
	if err != nil {
		return err
	}

	...

	return k.ChangeIssuance(ctx, coin.Denom, coin.Amount.Neg())
}
```

Every time Burn is called, the issuance of the coin is decremented in the state.  

### GetIssuance

```golang
func (k Keeper) GetIssuance(ctx sdk.Context, denom string, day sdk.Int) (issuance sdk.Int)
```

GetIssuance fetches the total issuance count of the coin matching `denom` for the `day`. If the `day` applies to a previous period, fetches the last stored snapshot issuance of the coin. For virgin calls, iterates through the accountkeeper and computes the genesis issuance.

For day 0, seigniorage is not recorded as mint mint starts its issuance memory from day 0. 
