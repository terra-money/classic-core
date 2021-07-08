# Events

## Handlers

## MsgStoreCode

| Type       | Attribute Key | Attribute Value |
| ---------- | ------------- | --------------- |
| store_code | sender        | {senderAddress} |
| store_code | codeID        | {codeID}        |
| message    | module        | wasm            |
| message    | action        | store_code      |
| message    | sender        | {senderAddress} |

## MsgMigrateCode

| Type         | Attribute Key | Attribute Value |
| ------------ | ------------- | --------------- |
| migrate_code | sender        | {senderAddress} |
| migrate_code | codeID        | {codeID}        |
| message      | module        | wasm            |
| message      | action        | migrate_code    |
| message      | sender        | {senderAddress} |

## MsgInstantiateContract

| Type                 | Attribute Key    | Attribute Value      |
| -------------------- | ---------------- | -------------------- |
| instantiate_contract | creator          | {creatorAddress}     |
| instantiate_contract | admin            | {adminAddress}       |
| instantiate_contract | code_id          | {codeID}             |
| instantiate_contract | contract_address | {contractAddress}    |
| message              | module           | wasm                 |
| message              | action           | instantiate_contract |
| message              | sender           | {senderAddress}      |

## MsgExecuteContract

| Type             | Attribute Key    | Attribute Value   |
| ---------------- | ---------------- | ----------------- |
| execute_contract | sender           | {senderAddress}   |
| execute_contract | contract_address | {contractAddress} |
| wasm-*           | ...              | ...               |
| wasm             | ...              | ...               |
| from_contract    | ...              | ...               |
| message          | module           | wasm              |
| message          | action           | execute_contract  |
| message          | sender           | {senderAddress}   |

## MsgMigrateContract

| Type             | Attribute Key    | Attribute Value   |
| ---------------- | ---------------- | ----------------- |
| migrate_contract | code_id          | {codeID}          |
| migrate_contract | contract_address | {contractAddress} |
| message          | module           | wasm              |
| message          | action           | migrate_contract  |
| message          | sender           | {senderAddress}   |

## MsgUpdateContractAdmin

| Type                  | Attribute Key    | Attribute Value       |
| --------------------- | ---------------- | --------------------- |
| update_contract_admin | admin            | {adminAddress}        |
| update_contract_admin | contract_address | {contractAddress}     |
| message               | module           | wasm                  |
| message               | action           | update_contract_admin |
| message               | sender           | {senderAddress}       |

## MsgClearContractAdmin

| Type                 | Attribute Key    | Attribute Value      |
| -------------------- | ---------------- | -------------------- |
| clear_contract_admin | contract_address | {contractAddress}    |
| message              | module           | wasm                 |
| message              | action           | clear_contract_admin |
| message              | sender           | {senderAddress}      |
