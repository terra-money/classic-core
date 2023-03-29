<!--
order: 4
-->

# Transactions

This section defines the `sdk.Msg` concrete types that result in the state transitions defined on the previous section.

## `MsgRegisterFeeShare`

Defines a transaction signed by a developer to register a contract for transaction fee distribution. The sender must be an EOA that corresponds to the contract deployer address.

```go
type MsgRegisterFeeShare struct {
  // contract_address in bech32 format
  ContractAddress string `protobuf:"bytes,1,opt,name=contract_address,json=contractAddress,proto3" json:"contract_address,omitempty"`
  // deployer_address is the bech32 address of message sender. It must be the
  // same the contract's admin address
  DeployerAddress string `protobuf:"bytes,2,opt,name=deployer_address,json=deployerAddress,proto3" json:"deployer_address,omitempty"`
  // withdrawer_address is the bech32 address of account receiving the
  // transaction fees
  WithdrawerAddress string `protobuf:"bytes,3,opt,name=withdrawer_address,json=withdrawerAddress,proto3" json:"withdrawer_address,omitempty"`
}
```

The message content stateless validation fails if:

- Contract bech32 address is invalid
- Deployer bech32 address is invalid
- Withdraw bech32 address is invalid

### `MsgUpdateFeeShare`

Defines a transaction signed by a developer to update the withdraw address of a contract registered for transaction fee distribution. The sender must be the admin of the contract.

```go
type MsgUpdateFeeShare struct {
  // contract_address in bech32 format
  ContractAddress string `protobuf:"bytes,1,opt,name=contract_address,json=contractAddress,proto3" json:"contract_address,omitempty"`
  // deployer_address is the bech32 address of message sender. It must be the
  // same the contract's admin address
  DeployerAddress string `protobuf:"bytes,2,opt,name=deployer_address,json=deployerAddress,proto3" json:"deployer_address,omitempty"`
  // withdrawer_address is the bech32 address of account receiving the
  // transaction fees
  WithdrawerAddress string `protobuf:"bytes,3,opt,name=withdrawer_address,json=withdrawerAddress,proto3" json:"withdrawer_address,omitempty"`
}
```

The message content stateless validation fails if:

- Contract bech32 address is invalid
- Deployer bech32 address is invalid
- Withdraw bech32 address is invalid

### `MsgCancelFeeShare`

Defines a transaction signed by a developer to remove the information for a registered contract. Transaction fees will no longer be distributed to the developer for this smart contract. The sender must be an admin that corresponds to the contract.

```go
type MsgCancelFeeShare struct {
  // contract_address in bech32 format
  ContractAddress string `protobuf:"bytes,1,opt,name=contract_address,json=contractAddress,proto3" json:"contract_address,omitempty"`
  // deployer_address is the bech32 address of message sender. It must be the
  // same the contract's admin address
  DeployerAddress string `protobuf:"bytes,2,opt,name=deployer_address,json=deployerAddress,proto3" json:"deployer_address,omitempty"`
}
```

The message content stateless validation fails if:

- Contract bech32 address is invalid
- Contract bech32 address is zero
- Deployer bech32 address is invalid
