<!--
order: 3
-->

# State Transitions

The `x/feeshare` module allows for three types of state transitions: `RegisterFeeShare`, `UpdateFeeShare` and `CancelFeeShare`. The logic for distributing transaction fees is handled through the [Ante handler](/app/ante.go).

## Register Fee Share

A developer registers a contract for receiving transaction fees by defining the contract address and the withdrawal address for fees to be paid too. If this is not set, the developer can not get income from the contract. This is opt-in for tax purposes. When registering for fees to be paid, you MUST be the admin of said wasm contract. The withdrawal address can be the same as the contract's address if you so choose.

1. User submits a `RegisterFeeShare` to register a contract address, along with a withdrawal address that they would like to receive the fees to
2. Check if the following conditions pass:
    1. `x/feeshare` module is enabled via Governance
    2. the contract was not previously registered
    3. deployer has a valid account (it has done at least one transaction)
    4. the contract address exists
    5. the deployer signing the transaction is the admin of the contract
    6. the contract is already deployed
3. Store an instance of the provided share.

All transactions sent to the registered contract occurring after registration will have their fees distributed to the developer, according to the global `DeveloperShares` parameter in governance.

### Update Fee Split

A developer updates the withdraw address for a registered contract, defining the contract address and the new withdraw address.

1. The user submits a `UpdateFeeShare`
2. Check if the following conditions pass:
    1. `x/feeshare` module is enabled
    2. the contract is registered
    3. the signer of the transaction is the same as the contract admin per the WasmVM
3. Update the fee with the new withdrawal address.

After this update, the developer receives the fees on the new withdrawal address.

### Cancel Fee Split

A developer cancels receiving fees for a registered contract, defining the contract address.

1. The user submits a `CancelFeeShare`
2. Check if the following conditions pass:
    1. `x/feeshare` module is enabled
    2. the contract is registered
    3. the signer of the transaction is the same as the contract admin per the WasmVM
3. Remove share from storage

The developer no longer receives fees from transactions sent to this contract. All fees go to the community.
