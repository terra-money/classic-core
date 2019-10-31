# Oracle

The Oracle module forms a consensus on the exchange rate of Luna with respect to various fiat currencies, in order to facilitate swaps amongst the stablecoins mirroring those currencies as well as ensuring their price-stability.

## Overview

The objective of the oracle module is to get accurate exchange rates of Luna with various fiat currencies such that the system can facilitate fair exchanges of Terra stablecoins and Luna. Should the system fail to gain an accurate understanding of Luna, a small set of arbitrageurs could profit at the cost of the entire system.

In order to get fair exchange rates, the oracle operates in the following way:

* Let P = {P1, P2, ...} be a time series split up by `params.VotePeriod`, currently 1 minute. In each P, validators must submit two votes: 
  * A `MsgPricePrevote`, containing the SHA256 hash of the exchange rate of Luna is with respect to a Terra peg. For example, in order to support swaps for Terra currencies pegged to KRW, USD, SDR, three prevotes must be submitted each containing the uluna&lt;&gt;ukrw, uluna&lt;&gt;uusd, and uluna&lt;&gt;usdr exchange rates. 
  * A `MsgPriceVote`, containing the salt used to create the hash for the prevote submitted in P-1.  
* At the end of each P, votes submitted are tallied. 
  * The submitted salt of each vote is used to verify consistency with the prevote submitted by the validator in P-1. If the validator has not submitted a prevote, or the SHA256 resulting from the salt does not match the hash from the prevote, the vote is dropped.
  * For each currency, if the total voting power of submitted votes exceeds 50%, a weighted median price of the vote is taken and is record on-chain as the effective exchange rate for Luna w.r.t. said currency for P+1.
  * Winners of the ballot for P-1, i.e. voters that have managed to vote within a small band around the weighted median, get rewarded by spread fees collected by swap operations during P. For spread rewards, see [this](market.md#spread-rewards).
* If an insufficient amount of votes have been received for a currency, below `VoteThreshold`, its exchange rate is deleted from the store, and no swaps can be made with it during P. 

```text
Period  |  P1 |  P2 |  P3 |  ...    |
Prevote |  O  |  O  |  O  |  ...    |
        |-----\-----\-----\-----    |
Vote    |     |  O  |  O  |  ...    |
```

Effectively this scheme forces the voter to commit to a firm price submission before knowing the votes of others, and thereby reduces centralization and free-rider risk in the oracle.

## Vote procedure

### Submit a prevote

```go
// MsgPricePrevote - struct for prevoting on the PriceVote.
// The purpose of prevote is to hide vote price with hash
// which is formatted as hex string in SHA256("salt:price:denom:voter")
type MsgPricePrevote struct {
    Hash      string         `json:"hash"` // hex string
    Denom     string         `json:"denom"`
    Feeder    sdk.AccAddress `json:"feeder"`
    Validator sdk.ValAddress `json:"validator"`
}
```

The `MsgPricePrevote` is just the submission of the leading 20 bytes of the SHA256 hex string run over a string containing the metadata of the actual `MsgPriceVote` to follow in the next period. The string is of the format: `salt:price:denom:voter`. Note that since in the subsequent `MsgPriceVote` the salt will have to be revealed, the salt used must be regenerated for each prevote submission.

`Denom` is the denomination of the currency for which the vote is being cast. For example, if the voter wishes to submit a prevote for the usd, then the correct `Denom` is `uusd`.

The price used in the hash must be the open market price of Luna, w.r.t. to the currency matching `Denom`. For example, if `Denom` is `uusd` and the going price for Luna is 1 USD, then "1" must be used as the price, as 1 `uluna` = 1 `uusd`. 

`Feeder` is used if the validator wishes to delegate oracle vote signing to a separate key to de-risk exposing their validator signing key.

`Validator` is the validator address of the original validator.

### Submit a vote

```go
// MsgPriceVote - struct for voting on the price of Luna denominated in various Terra assets.
// For example, if the validator believes that the effective price of Luna in USD is 10.39, that's
// what the price field would be, and if 1213.34 for KRW, same.
type MsgPriceVote struct {
    Price     sdk.Dec        `json:"price"` // the effective price of Luna in {Denom}
    Salt      string         `json:"salt"`
    Denom     string         `json:"denom"`
    Feeder    sdk.AccAddress `json:"feeder"`
    Validator sdk.ValAddress `json:"validator"`
}
```

The `MsgPriceVote` contains the actual price vote. The `Salt` parameter must match the salt used to create the prevote, otherwise the voter cannot be rewarded.


### Delegate voting rights to another key

Validators may also elect to delegate voting rights to another key to prevent the block signing key from being kept online. To do so, they must submit a `MsgDelegateFeederPermission`, delegating their oracle voting rights to a `FeedDelegate`, which in turn sign `MsgPricePrevote` and `MsgPriceVote` on behalf of the validator. 

{% hint style="info" %}
Make sure to populate the delegate address with some coins by which to pay fees.
{% endhint %}

```go
// MsgDelegateFeederPermission - struct for delegating oracle voting rights to another address.
type MsgDelegateFeederPermission struct {
	Operator     sdk.ValAddress `json:"operator"`
	FeedDelegate sdk.AccAddress `json:"feed_delegate"`
}
```

The `Operator` field contains the operator address of the validator. The `FeedDelegate` field is the address of the delegate account that will be submitting price related votes and prevotes on behalf of the `Operator`. 


## Parameters

```go
// Params oracle parameters
type Params struct {
    VotePeriod       int64   `json:"vote_period"`        // voting period in block height; tallys and reward claim period
    VoteThreshold    sdk.Dec `json:"vote_threshold"`     // minimum stake power threshold to update price
    OracleRewardBand sdk.Dec `json:"oracle_reward_band"` // band around the oracle weighted median to reward
}
```
