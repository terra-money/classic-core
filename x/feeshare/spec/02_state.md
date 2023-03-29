<!--
order: 2
-->

# State

## State Objects

The `x/feeshare` module keeps the following objects in the state:

| State Object          | Description                           | Key                                                               | Value              | Store |
| :-------------------- | :------------------------------------ | :---------------------------------------------------------------- | :----------------- | :---- |
| `FeeShare`            | Fee split bytecode                    | `[]byte{1} + []byte(contract_address)`                            | `[]byte{feeshare}` | KV    |
| `DeployerFeeShares`   | Contract by deployer address bytecode | `[]byte{2} + []byte(deployer_address) + []byte(contract_address)` | `[]byte{1}`        | KV    |
| `WithdrawerFeeShares` | Contract by withdraw address bytecode | `[]byte{3} + []byte(withdraw_address) + []byte(contract_address)` | `[]byte{1}`        | KV    |

### FeeShare

A FeeShare defines an instance that organizes fee distribution conditions for
the owner of a given smart contract

```go
type FeeShare struct {
  // contract_address is the bech32 address of a registered contract in string form
  ContractAddress string `protobuf:"bytes,1,opt,name=contract_address,json=contractAddress,proto3" json:"contract_address,omitempty"`
  // deployer_address is the bech32 address of message sender. It must be the
  // same as the contracts admin address.
  DeployerAddress string `protobuf:"bytes,2,opt,name=deployer_address,json=deployerAddress,proto3" json:"deployer_address,omitempty"`
  // withdrawer_address is the bech32 address of account receiving the
  // transaction fees.
  WithdrawerAddress string `protobuf:"bytes,3,opt,name=withdrawer_address,json=withdrawerAddress,proto3" json:"withdrawer_address,omitempty"`
}
```

### ContractAddress

`ContractAddress` defines the contract address that has been registered for fee distribution.

### DeployerAddress

A `DeployerAddress` is the admin address for a registered contract.

### WithdrawerAddress

The `WithdrawerAddress` is the address that receives transaction fees for a registered contract.

## Genesis State

The `x/feeshare` module's `GenesisState` defines the state necessary for initializing the chain from a previously exported height. It contains the module parameters and the fee share for registered contracts:

```go
// GenesisState defines the module's genesis state.
type GenesisState struct {
  // module parameters
  Params Params `protobuf:"bytes,1,opt,name=params,proto3" json:"params"`
  // active registered contracts for fee distribution
  FeeShares []FeeShare `protobuf:"bytes,2,rep,name=feeshares,json=feeshares,proto3" json:"feeshares"`
}
```
