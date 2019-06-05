# Oracle

The Oracle module forms a consensus on the exchange rate of Luna with respect to various fiat currencies, in order to facilitate swaps amongst the stablecoins mirroring those currencies as well as ensuring their price-stability. 

## Overview

The objective of the oracle module is to get accurate exchange rates of Luna with various fiat currencies such that the system can facilitate fair exchanges of Terra stablecoins and Luna. Should the system fail to gain an accurate understanding of Luna, a small set of arbitrageurs could profit at the cost of the system. 

In order to get fair exchange rates, the oracle operates in the following way: 

- Over a `VotePeriod`, validators submit price feed votes, and the weighted-by-Luna-stake median of the votes is tallied to be the correct price of Luna in the subsequent period. Winners of the ballot, i.e. voters that have managed to vote within a small band around the weighted median, get rewarded by the fees that are collected in Market swap operations. The VotePeriod is kept extremely tight (currently 1 minute) to minimize the risk of price drift. 

- A spread is charged for transactions involving Luna (2 ~ 10%), and the fees collected here is used to compensate ballot winners. Oracle swap rewards are doled out the end of every `VotePeriod`, and rewards the validator and its delegations in accordance with the logic of the distribution module. 

- In order to minimize the risk of oracle frontrunning (people waiting to see where the price consensus forms and voting to receive the reward at the very end of the  `VotePeriod`), oracle votes are demarcated into two: prevotes and votes. 
```
Period  |  P1 |  P2 |  P3 |  ...    |
Prevote |  O  |  O  |  O  |  ...    |
        |-----\-----\-----\-----    |
Vote    |     |  O  |  O  |  ...    |
```
In the prevote stage, a validator should submit the hash of the part of real vote msg to prove the validator is not just copying other validators price vote. In vote phrase, the validator should reveal the real price by submitting `MsgPriceVote` with the salt.

The submission order has to be kept in (vote -> prevote) order. If an prevote comes early, it will replace previous prevote so next vote, which reveals the proof for previous prevote, will be failed.

- If an insufficient amount of votes have been received for a currency, below `VoteThreshold`, its exchange rate is deleted from the store, and no swaps can be made with it. 


## Vote procedure

### Submit a prevote

```golang
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

The `MsgPricePrevote` is just the submission of the leading 20 bytes of the SHA256 hex string run over a string containing the metadata of the actual `MsgPriceVote` to follow in the next period. The string is of the format: `salt:price:denom:voter`. 

Effectively this scheme forces the voter to commit to a firm price submission before knowing the votes of others, and thereby reduces centralization and free-rider risk in the oracle. 

### Submit a vote

```golang
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

The `MsgPriceVote` contains the actual price vote. the `Salt` parameter must match the salt used to create the prevote, otherwise the voter cannot be rewarded. 

## Parameters

```golang
// Params oracle parameters
type Params struct {
	VotePeriod       int64   `json:"vote_period"`        // voting period in block height; tallys and reward claim period
	VoteThreshold    sdk.Dec `json:"vote_threshold"`     // minimum stake power threshold to update price
	OracleRewardBand sdk.Dec `json:"oracle_reward_band"` // band around the oracle weighted median to reward
}
```
