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
