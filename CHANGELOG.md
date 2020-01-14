## 0.3.2

### Improvements
#### [\#313](https://github.com/terra-project/core/pull/313) upgrade SDK
* Bump SDK version to [v0.37.5](https://github.com/cosmos/cosmos-sdk/releases/tag/v0.37.5)
* Tendermint version to [v0.32.8](https://github.com/tendermint/tendermint/releases/tag/v0.32.8)
#### [\#312](https://github.com/terra-project/core/pull/312) upgrade golangci-lint version to v1.22.2

## 0.3.1

### Bug Fixes
#### [\#303](https://github.com/terra-project/core/pull/303) fix estimate fee endpoint for multiple signature tx
#### [\#304](https://github.com/terra-project/core/pull/304) genesis scrpit update

### Improvements
#### [\#301](https://github.com/terra-project/core/pull/301) README update
#### [\#305](https://github.com/terra-project/core/pull/305) swagger update
#### [\#306](https://github.com/terra-project/core/pull/306) circleci update for goreleaser

## 0.3.0
### Breaking Changes
#### [\#265](https://github.com/terra-project/core/pull/265) Oracle refactor & Oracle slashing
##### Slashing
A validator get slashed `SlashFraction`% if the one perform any of the following violations in `SlashWindow - minValidPerWindow` voteperiods over a window of `SlashWindow` voteperiods:

1. A vote is missing for any of the denom in the whitelist. Oracle voters looking to abstain must still submit a "vote of no confidence", which has 0 for the luna ExchangeRate field of the prevote.
2. A submitted vote is more than max(`RewardBand`, standard deviation) from the elected median
3. Oracle voters who submit abstain vote with other invalid votes will also get slashed.

##### Codec Changes
* `oracle/MsgDelegateFeederPermission` => `oracle/MsgDelegateFeedConsent`

##### New End Points
* `/oracle/voters/{%s}/miss` return the # of vote periods missed in this oracle slash window.

##### Path Changes
* `/oracle/denoms/{denom}/price` => `/oracle/denoms/{denom}/exchange_rate`

##### Request Body Changes

* POST /oracle/denoms/{denom}/prevotes
* POST /oracle/denoms/{denom}/votes
```
Price sdk.Dec `json:"price"`
```
has been changed to 
```
ExchangeRate sdk.Dec `json:"exchange_rate"`
```

#### [\#256](https://github.com/terra-project/core/pull/256) Oracle endpoints improvement
##### New EndPoints
```
/oracle/voters/{validator}/votes
/oracle/voters/{validator}/prevotes
/oracle/denoms/prices
```

#### [\#250](https://github.com/terra-project/core/pull/250) Oracle whitelist & Reward distribution update
* Create a whitelist param that stores an array of denoms that are whitelisted by the protocol. 
* Edit the oracle `Reward Pool of a VotePeriod = oracle module account / (n vote periods)`. 
* Oracle module account is whitelisted in the bank module such that users can donate funds to the oracle module account

#### [\#234](https://github.com/terra-project/core/pull/234) Adopt gov module
`distribution` module already contains `community-pool-spend` proposal suitable for `budget` so budget module is removed. There are two custom governance proposals from `treasury` module; `tax-rate-update` & `reward-weight-update` proposals. 

##### New EndPoints
```
(GET)/gov/proposals
(GET)/gov/proposals/{proposalId}
(GET)/gov/proposals/{proposalId}/proposer
(GET)/gov/proposals/{proposalId}/deposits
(GET)/gov/proposals/{proposalId}/deposits/{depositor}
(GET/POST)/gov/proposals/{proposalId}/votes
(GET)/gov/proposals/{proposalId}/votes/{voter}
(GET)/gov/proposals/{proposalId}/tally
(GET)/gov/parameters/deposit
(GET)/gov/parameters/tallying
(GET)/gov/parameters/voting
(POST)/gov/proposals/tax_rate_update
(POST)/gov/proposals/reward_weight_update
(POST)/gov/proposals/param_change
(POST)/gov/proposals/community_pool_spend
```

#### [\#233](https://github.com/terra-project/core/pull/233) Swap constant product
As proposed [here](https://agora.terra.money/uploads/short-url/92QHxFtEmWUEwf9kWTminuobwpM.pdf), apply constant product to swap feature. 

##### Compute Pools
```
// Both LUNA and TERRA pools are using SDR units.
cp = basePool*basePool
terraPool = (basePool + terraDelta)
lunaPool = cp/terraPool
```

##### LUNA to TERRA swap
```
// offerAmt must be SDR units
newLunaPool = lunaPool + offerAmt
newTerraPool = cp / newLunaPool
returnAmt = newTerraPool - terraPool

// Swap return SDR Amt to TERRA
returnLunaAmt = market.swap(returnAmt, "LUNA")
```

##### TERRA to LUNA swap
```
// offerAmt must be SDR units
newTerraPool = terraPool + offerAmt
newLunaPool = cp / newTerraPool
returnAmt = newLunaPool - lunaPool

// Swap return SDR Amt to proper TERRA
returnTerraAmt = market.swap(returnAmt, "TERRA")
```

##### TERRA to TERRA swap

Apply only fixed tobin-tax without computing and changing pools


##### New EndPoints
```
/market/terra_pool_delta
```

#### [\#231](https://github.com/terra-project/core/pull/231) Bump SDK to v0.37.x
##### REST end points, which are changed
All REST responses now wrap the original resource/result. The response
will contain two fields: height and result.
```
/market/params => /market/parameters
/oracle/params => /oracle/parameters
/treasury/tax-rate => /treasury/tax_rate
/treasury/tax-rate/{epoch} => /treasury/tax_rate/{epoch}
/treasury/tax-cap => /treasury/tax_cap
/treasury/tax-cap/{denom} => /treasury/tax_cap/{denom}
/treasury/reward-weight => /treasury/reward_weight
/treasury/reward-weight/{epoch} => /treasury/reward_weight/{epoch}
/treasury/tax-proceeds => /treasury/tax_proceeds
/treasury/tax-proceeds/{epoch} => /treasury/tax_proceeds/{epoch}
/treasury/seigniorage-proceeds => /treasury/seigniorage_proceeds
/treasury/seigniorage-proceeds/{epoch} => /treasury/seigniorage_proceeds/{epoch}
/treasury/current-epoch => /treasury/current_epoch
/treasury/params => /treasury/parameters
```
##### REST end points, which response object key is removed
```
/treasury/current_epoch
/treasury/seigniorage_proceeds/{epoch}
/treasury/seigniorage_proceeds
/treasury/tax_proceeds/{epoch}
/treasury/tax_proceeds
/treasury/reward_weight/{epoch}
/treasury/reward_weight
/treasury/tax_cap/{denom}
/treasury/tax_rate/{epoch}
/treasury/tax_rate
/oracle/denoms/actives
/oracle/denoms/{denom}/price
/oracle/denoms/{denom}/prevotes/{voter}
/oracle/denoms/{denom}/prevotes
/oracle/denoms/{denom}/votes/{voter}
/oracle/denoms/{denom}/votes
```
##### New REST endpoints
```
/supply/total
/supply/total/{denomination}
/market/last_day_issuance
/oracle/voters/{%s}/voting_info
/oracle/voting_infos
/treasury/historical_issuance/{epoch}
```
##### Codec changes
```
auth/Account => core/Account
auth/StdTx => core/StdTx
pay/MsgSend => bank/MsgSend
pay/MsgMultiSend => bank/MsgMultiSend
```
##### Other Changes
* GradedVestingAccount is fully removed
* LazyGradedVestingAccount's vesting_lazy_schedules is changed to vesting_schedules
* Improve the UX of fee and tax. The sender have to specify fees containing tax amount with following methods.
  1. Use `/bank/accounts/{address}/transfers` without fees. `terracli` will compute tax amount and fill fees field containing both gas & tax fee.
  2. Use `/txs/estimate_fee` to estimate fees of StdTX, and replace StdTx.Fee.Amount to estimated fee
  3. Compute tax with `/treasury/tax_rate` & `/treasury/tax_cap` add computed tax with original gas fee

## 0.2.4
### Bug fixes
#### [\#196](https://github.com/terra-project/core/pull/196) peek epoch seigniorage
Change PeekEpochSeigniorage to compute seigniorage by subtracting current issuance from previous issuance

#### [\#198](https://github.com/terra-project/core/pull/198) Use next block for treasury tax and reward update
updateTaxPolicy and updateRewardPolicy are updating new tax-rate and reward-weight with current ctx. The ctx height is the last block of current epoch, but treasury should update next epoch's tax-rate and reward-weight at the last block of current epoch.
In updateTaxPolicy and updateRewardPolicy, change ctx input of keeper setter to ctx with next epoch height.

### Features
#### [\#193](https://github.com/terra-project/core/pull/193) Recover old hd path
Added `--old-hd-path` option to `$terracli keys add` command for recovering old bip44 path(for atom)
##### Example
```
$ terracli keys add tmp --recover --old-hd-path
Enter a passphrase to encrypt your key to disk:
Repeat the passphrase:
> Enter your bip39 mnemonic
candy hint hamster cute inquiry bright industry decide assist wedding carpet fiber arm menu machine lottery type alert fan march argue adapt recycle stomach

NAME:	TYPE:	ADDRESS:					PUBKEY:
tmp   local	terra1gaczd45crhwfa4x05k9747cuxwfmnduvmtyefs	terrapub1addwnpepqv6tse2pyag9ts5vy6dk4h3qh7xc9qhat4jx449n6nrfve3jhzldz3f3l7p
```

## 0.2.3
- [\#187](https://github.com/terra-project/core/pull/187): Change all time instance timezone to UTC to remove gap in time calculation

### Changes
#### [\#187](https://github.com/terra-project/core/pull/187) Bugfix/fix-time-zone
In update_230000.go, we change genesis time derivation from 
```
genesisTime := time.Unix(genesisUnixTime, 0)
```
to
```
genesisTime := time.Unix(genesisUnixTime, 0).UTC()
```

## 0.2.2

- [\#185](https://github.com/terra-project/core/pull/185): Improve oracle specs
- [\#184](https://github.com/terra-project/core/pull/184): Fix `terracli` docs
- [\#183](https://github.com/terra-project/core/pull/183): Change all GradedVestingAccounts to LazyGradedVestingAccounts.
- [\#179](https://github.com/terra-project/core/pull/179): Conform querier responses to be returned in JSON format 
- [\#178](https://github.com/terra-project/core/pull/178): Change BIP44 PATH to 330

### Changes
#### [\#185](https://github.com/terra-project/core/pull/185) Oracle `MsgFeederDelegatePermission` specs
Added docs for using `MsgFeederDelegatePermission` to oracle specs

#### [\#185](https://github.com/terra-project/core/pull/185) Oracle price vote denom error fix 
Oracle specs now specify micro units `uluna` and `uusd` for correct denominations for price prevotes and votes 

#### [\#184](https://github.com/terra-project/core/pull/184) Minor terracli fix 

#### [\#183](https://github.com/terra-project/core/pull/183) Oracle param update
```
OracleRewardBand: 1% => 2%
```

#### [\#183](https://github.com/terra-project/core/pull/183) Market param update
```
DailyLunaDeltaCap: 0.5% => 0.1%
```

#### [\#183](https://github.com/terra-project/core/pull/183) LazyGradedVestingAccount

* Spread out the cliffs for presale investors, with varying degrees of severity (details [\#180](https://github.com/terra-project/core/issues/180))

#### [\#179](https://github.com/terra-project/core/pull/179) Align Querier responses to JSON

* Querier was returning misaligned formats for return values, now aligned to JSON format

#### [\#178](https://github.com/terra-project/core/pull/178) Correctly use 330 as the coin type field in BIP 44 PATH

* We were previously using the Cosmos coin type field for the BIP44 path. Changed to Terra's own 330. 


## 0.2.1

- [\#166](https://github.com/terra-project/core/pull/166): Newly added parameters were not being added to the columbus-2 genesis.json file. Fixed. 

## 0.2.0

### Bug Fixes

* [\#140](https://github.com/terra-project/core/pull/140) Fix export bug.

* [\#140](https://github.com/terra-project/core/pull/140) Client querier bug fix (distr outstanding rewards)

* [\#140](https://github.com/terra-project/core/pull/140) Fix budget module to delete all votes when submitter withdraws the program and to use DeleteVotesForProgram to delete all votes for a program.

### Improvements
#### [\#140](https://github.com/terra-project/core/pull/140) Msg Types

```
cosmos-sdk/MsgSend => pay/MsgSend
cosmos-sdk/MsgMultiSend => pay/MsgMultiSend

cosmos-sdk/MsgCreateValidator => staking/MsgCreateValidator
cosmos-sdk/MsgEditValidator => staking/MsgEditValidator
cosmos-sdk/MsgDelegate => staking/MsgDelegate
cosmos-sdk/MsgUndelegate => staking/MsgUndelegate
cosmos-sdk/MsgBeginRedelegate => staking/MsgBeginRedelegate

cosmos-sdk/MsgWithdrawDelegationReward => distribution/MsgWithdrawDelegationReward
cosmos-sdk/MsgWithdrawValidatorCommission => distribution/MsgWithdrawValidatorCommission
cosmos-sdk/MsgModifyWithdrawAddress => distribution/MsgModifyWithdrawAddress

cosmos-sdk/MsgUnjail => slashing/MsgUnjail
```
  
#### [\#140](https://github.com/terra-project/core/pull/140) Oracle updates prevoting/voting
MsgPriceFeed is split into ```MsgPricePrevote``` and ```MsgPriceVote```
```
Period  |  P1 |  P2 |  P3 |  ...    |
Prevote |  O  |  O  |  O  |  ...    |
        |-----\-----\-----\-----    |
Vote    |     |  O  |  O  |  ...    |
```
In prevote stage, a validator should submit the hash of the part of real vote msg to prove the validator is not just copying other validators price vote. In vote phrase, the validator should reveal the real price by submitting MsgPriceVote with ```salt```.

The submission order has to be kept in (vote -> prevote) order. If an prevote comes early, it will replace previous prevote so next vote, which reveals the proof for previous prevote, will be failed.

#### [\#148](https://github.com/terra-project/core/pull/148) Oracle voting right delegation
By using the oracle/MsgDelegateFeederPermission a validator can assign the right to vote to another account at any time. The validator account will preserve its right to vote at any time.

#### [\#140](https://github.com/terra-project/core/pull/140) & [\#148](https://github.com/terra-project/core/pull/148) Rest Interface Update
##### Change rest interface url
```
"/distribution/parameters" => "/distribution/params"
"/staking/parameters" => "/staking/params"
```

##### Change request body
```
Send request body
From:
type SendReq struct {
	BaseReq rest.BaseReq `json:"base_req"`
	Amount   sdk.Coins    `json:"coins"`
}

To:
type SendReq struct {
	BaseReq rest.BaseReq `json:"base_req"`
	Coins   sdk.Coins    `json:"coins"`
}
```

##### New rest interfaces
```
(GET/POST) "/oracle/denoms/{%s}/votes"
(GET/POST) "/oracle/denoms/{%s}/prevotes"
(GET/POST) "/oracle/voters/{%s}/feeder"
(GET/POST) "/market/swap"
(GET/POST) "/market/params"
```

#### [\#140](https://github.com/terra-project/core/pull/140) Add transaction logs for tax and swap amount
##### Send Tx
Add **tax** log to send transaction for recording real amount which a transaction pay.
Ex) txs/B515331BF9EA9A92AD59A85D593E5A2B170E3D297C59E85DDA9FA6FF33790E9B
```
{
  "logs": [
    {
      "msg_index": 0,
      "success": true,
      "log": "{\"tax\":\"400uluna\"}"
    }
  ]
}
```

##### Swap Tx
Add **swap_coin** log to swap transaction for recording the amount of swapped coin along with offered coin
```
{
  "logs": [
    {
      "msg_index": 0,
      "success": true,
      "log": "{\"swap_coin\":\"400ukrw\"}"
    }
  ]
}
```

#### [\#150](https://github.com/terra-project/core/pull/150) Market Swap protections

##### Add bidirectional Luna supply change cap on market swaps. 
A daily trading cap (luna supply change cap) protects excessive luna volatility. Capping Luna deflation prevents divesting attacks (attacker swaps large amount into terra to avoid slippage) and consensus attacks by limiting access to staking tokens. Early parameters are 2% - 10% on both sides of the trade. 

##### Add bidirectional Luna spread fees on market swaps 
To protect against short term price deviations between the open market and the on-chain oracle, we now charge a 2-10% spread on swaps that involve luna. 

##### Change oracle reward scheme from monthly seigniorage to validators to minute distribution
Swap spreads are distributed to oracle ballot winners on the oracle VotePeriod; this vastly shortens distribution periods. Also, all stakeholders receive oracle rewards (includes delegators).

##### Swaps halt immediately after an illiquid oracle vote
Previously we facilitated swaps for 10 VotePeriods after the last valid oracle ballot. We now stop swaps immediately to prevent arbitrage attacks from price drift.


### Parameter Changes 

#### [\#150](https://github.com/terra-project/core/pull/150) Change MiningRewardWeight.Max from 20% to 90%. This is to reduce volatility in fees at network infancy. 

#### Changed BlocksPerMinute from 12 to 5 to more accurately reflect Columbus block times.
