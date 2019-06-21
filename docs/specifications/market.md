# Market

The Market module facilitates the atomic swaps among terra currencies and luna.

## Overview

The market module facilitates swaps between all terra currencies that have an active exchange rate with Luna registered with the Oracle module and Luna.

* A user can swap SDT \(TerraSDR\) and UST \(TerraUSD\) at the exchange rate registered with the oracle. For example, if Luna&lt;&gt;SDT exchange rate returned by `GetLunaSwapRate` by the oracle is 10, and Luna&lt;&gt;KRT exchange rate is 10,000, a swapping 1 SDT will return 1000 KRT.
* A user cap swap any of the Terra currencies for Luna at the oracle exchange rate. Using the same exchange rates in the above example, a user can swap 1 SDT for 0.1 Luna, or 0.1 Luna for 1 SDT.

## Safety mechanisms for Luna swaps

* A daily Luna supply change cap is enforced, such that Luna supply can inflate or deflate only up to the cap in any given 24 hour period. Swap transactions after the cap has been hit fails. This is to prevent excessive volatility in Luna supply which can lead to divesting attacks \(a large increase in Terra supply putting the peg at risk\) or consensus attacks \(a large increase in Luna supply being staked can lead to a consensus attack on the blockchain\).
* A spread is enforced on swaps involving Luna, currently between 2-10%.

  ```text
  // Compute a spread, which is initialiy MinSwapSpread and grows linearly to MaxSwapSpread with delta
  spread = MinSwapSpread + dailyDelta / maxDelta *  (MaxSwapSpread - MinSwapSpread)
  ```

  where `MinSwapSpread` and `MaxSwapSpread` is the minimum and maximum luna swap spreads charged respectively. The spread starts at the minimum and linearly increases to the max spread as the current luna supply approximates the daily supply cap in either direction.

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

The spread fee charged in swaps involving Luna is distributed to the `SwapFeePool` in the oracle to be distributed to the oracle voters that voted close to the elected price at the end of every oracle `VotePeriod`.

## Parameters

```go
// Params market parameters
type Params struct {
    DailyLunaDeltaCap sdk.Dec `json:"daily_luna_delta_limit"` // daily % inflation or deflation cap on Luna
    MinSwapSpread     sdk.Dec `json:"min_swap_spread"`        // minimum spread for swaps involving Luna
    MaxSwapSpread     sdk.Dec `json:"max_swap_spread"`        // maximum spread for swaps involving Luna
}
```

