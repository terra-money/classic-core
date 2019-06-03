## 0.2.0

### Bug Fixes

* [\#140](https://github.com/terra-project/core/pull/140) Fix export bug.

* [\#140](https://github.com/terra-project/core/pull/140) client querier bug fix (distr outstanding rewards)

### Improvements
#### [\#140](https://github.com/terra-project/core/pull/140) Msg Types

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