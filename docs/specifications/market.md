# Market

The Market module facilitates the atomic swaps among terra currencies and luna.

## Overview

The market module adopts uni-swap modle to facilitate swaps between all terra currencies that have an active exchange rate with Luna registered with the Oracle module and Luna.

* A user can swap SDT \(TerraSDR\) and UST \(TerraUSD\) at the exchange rate registered with the oracle. For example, if Luna&lt;&gt;SDT exchange rate returned by `GetLunaExchangeRate` by the oracle is 10, and Luna&lt;&gt;KRT exchange rate is 10,000, a swapping 1 SDT will return 1000 KRT.
* A user cap swap any of the Terra currencies for Luna at the oracle exchange rate. Using the same exchange rates in the above example, a user can swap 1 SDT for 0.1 Luna, or 0.1 Luna for 1 SDT.

## Safety mechanisms for Luna swaps

* A Terra liquidity pool (`BasePool`) determines the maximum amount of Terra that can be issued or burned in any 24 hour period. Swap transactions after the cap has been hit pay most of the swaped coins for spread fees. This is to prevent excessive volatility in Luna supply which can lead to divesting attacks \(a large increase in Terra supply putting the peg at risk\) or consensus attacks \(a large increase in Luna supply being staked can lead to a consensus attack on the blockchain\). Luna’s market capitalization is expected to be larger than Terra’s, meaning that a cap relative to Terra supply serves as an effective cap on Luna’s supply.
* The system only charges Tobin Tax (0.3%) for the Terra to Terra swap without constant product spread.
* A spread is enforced on swaps involving Luna, currently between 2-100%.

  ```text
  // Compute a spread, which is initialiy MinSpread and grows reciprocally to 1
  // Swap PoolA -> PoolB (Both pools are in SDR units)
  CP = BasePool * BasePool = PoolA * PoolB
  swapAmt = (PoolB - CP / (PoolA + offerAmt))
          = PoolB * offerAmt / (PoolA + offerAmt)

  // Since both pools are the same unit, askAmt should be the same as offerAmt 
  // including the spread caused by constant product.
  spread = 1 - swapAmt / offerAmt
  ```

  where `MinSwapSpread` is the minimum luna swap spreads charged respectively. The spread starts at the minimum and reciprocally increases to the max(100%) spread as the current terra supply approximates the daily supply cap in either direction.

## Swap procedure

```go
// MsgSwap contains a swap request
type MsgSwap struct {
    Trader    sdk.AccAddress `json:"trader"`     // Address of the trader
    OfferCoin sdk.Coin       `json:"offer_coin"` // Coin being offered
    AskDenom  string         `json:"ask_denom"`  // Denom of the coin to swap to
}
```

The trader can submit a `MsgSwap` transaction with the amount / denomination of the coin to be swapped, the "offer", and the denomination of the coins to be swapped into, the "ask".

If the trader's `Account` has insufficient balance to execute the swap, the swap transaction fails. Upon successful completion of swaps involving Luna, a portion of the coins to be credited to the user's account is withheld as the spread fee.

## Spread rewards

The spread fee charged in swaps involving Luna is burned to maintain the constant product.

## Parameters

```go
// Params market parameters
type Params struct {
	PoolUpdateInterval       int64   `json:"pool_update_interval" yaml:"pool_update_interval"`                // reset interval of BasePool
	DailyTerraLiquidityRatio sdk.Dec `json:"daily_terra_liquidity_ratio" yaml:"daily_terra_liquidity_ratio"`  // daily % inflation or deflation cap on Terra
	MinSpread                sdk.Dec `json:"min_spread" yaml:"min_spread"`                                    // minimum spread for swaps involving Luna
	TobinTax                 sdk.Dec `json:"tobin_tax" yaml:"tobin_tax"`                                      // a tax on Terra<>Terra swap
}
```

