# Pay

The pay module is the base transactional layer of the Terra blockchain: it allows assets to be sent from one `Account` to another. It reads the currently effective `tax-rate` and `tax-cap` parameters from the `treasury` module to enforce a stability layer fee. 

## Overview

### Send 

```golang
// MsgSend - high level transaction of the coin module
type MsgSend struct {
	FromAddress sdk.AccAddress `json:"from_address"`
	ToAddress   sdk.AccAddress `json:"to_address"`
	Amount      sdk.Coins      `json:"amount"`
}
```

The pay module can be used to send coins from one `Account` to another. A `MsgSend` is constructed to facilitate the transfer. If the balance of coins in the `Account` is insufficient or the recipient `Account` does not exist, the transaction fails. 

### Multisend via batching  

```golang
// MsgMultiSend - high level transaction of the coin module
type MsgMultiSend struct {
	Inputs  []Input  `json:"inputs"`
	Outputs []Output `json:"outputs"`
}
```

The pay module can be used to send multiple transactions at once. `Inputs` contains the incoming transactions, and `Outputs` contains the outgoing transactions. The coin balance of the `Inputs` and the `Outputs` must match exactly. Batching transactions via multisend has the benefit of conserving network bandwidth and gas fees. 

If any of the `Accounts` fails, then taxes and fees already paid through the transaction is not refunded. 

## Fees

### Gas fees

As with any other transaction, `MsgSend` and `MsgMultiSend` has to pay a gas fee the size of which depends on validator's preferences (each validator sets his own min-gas-fees) and the complexity of the transaction. [Notes on gas and fees](../guide/users.md#a-note-on-gas-and-fees) has a more detailed explanation of how gas is computed. Important detail to note here is that gas fees are specified by the sender when the transaction is outbound. 

### Stability fees

Further to the gas fee, the pay module charges a stability fee that is a percentage of the transaction's value. It reads the `tax-rate` and `tax-cap` parameters from the treasury module to compute the amount of stability tax that needs to be charged. 

- `tax-rate`: an sdk.Dec object specifying what % of send transactions must be paid in stability fees
- `tax-cap`: a cap unique to each currency specifying the absolute cap that can be charged in stability fees from a given transaction. 

For an example `MsgSend` transaction of 1000 usdr tokens,

```
stability fee = min(1000 * tax_rate, tax_cap(usdr))
```

For a `MsgMultiSend` transaction, a stability fee is charged from every outbound transaction. 

Unlike with the gas fee which needs to be specified by the sender, the stability fee is automatically deducted from the sender's `Account`. 
