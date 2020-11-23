# Events

## Handlers

## MsgStoreCode

| Type       | Attribute Key | Attribute Value |
|------------|---------------|-----------------|
| store_code | codeID        | {codeID}        |
| store_code | sender        | {senderAddress} |  
| message    | module        | wasm            |
| message    | action        | store_code      |
| message    | sender        | {senderAddress} |

## MsgInstantiateContract

| Type                 | Attribute Key    | Attribute Value      |
|----------------------|------------------|----------------------|
| instantiate_contract | owner            | {ownerAddress}       |
| instantiate_contract | code_id          | {codeID}             |  
| instantiate_contract | contract_address | {contractAddress}    |  
| message              | module           | wasm                 |
| message              | action           | instantiate_contract |
| message              | sender           | {senderAddress}      |

## MsgExecuteContract

| Type             | Attribute Key    | Attribute Value   |
|------------------|------------------|-------------------|
| execute_contract | contract_address | {contractAddress} |
| execute_contract | sender           | {senderAddress}   |
| message          | module           | wasm              |
| message          | action           | execute_contract  |
| message          | sender           | {senderAddress}   |

## MsgMigrateContract

| Type             | Attribute Key    | Attribute Value   |
|------------------|------------------|-------------------|
| migrate_contract | code_id          | {codeID}          |
| migrate_contract | contract_address | {contractAddress} |
| message          | module           | wasm              |
| message          | action           | migrate_contract  |
| message          | sender           | {senderAddress}   |

## MsgUpdateContractOwner

| Type                  | Attribute Key    | Attribute Value        |
|-----------------------|------------------|------------------------|
| update_contract_owner | owner            | {ownerAddress}         |
| update_contract_owner | contract_address | {contractAddress}      |
| message               | module           | wasm                   |
| message               | action           | update_contract_owner  |
| message               | sender           | {senderAddress}        |