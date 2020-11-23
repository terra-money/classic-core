<!--
order: 5
-->

# Events

The market module emits the following events:

## Handlers

### MsgGrantAuthorization

| Type                | Attribute Key | Attribute Value     |
|---------------------|---------------|---------------------|
| grant_authorization | grant_type    | {msgType}           |
| grant_authorization | granter       | {granterAddress}    |
| grant_authorization | grantee       | {granteeAddress}    |
| message             | module        | msgauth             |
| message             | action        | grant_authorization |
| message             | sender        | {senderAddress}     |

### MsgRevokeAuthorization

| Type                 | Attribute Key | Attribute Value      |
|----------------------|---------------|----------------------|
| revoke_authorization | grant_type    | {msgType}            |
| revoke_authorization | granter       | {granterAddress}     |
| revoke_authorization | grantee       | {granteeAddress}     |
| message              | module        | msgauth              |
| message              | action        | revoke_authorization |
| message              | sender        | {senderAddress}      |

### MsgExecAuthorized

| Type                  | Attribute Key   | Attribute Value       |
|-----------------------|-----------------|-----------------------|
| execute_authorization | grantee_address | {granteeAddress}      |
| message               | module          | msgauth               |
| message               | action          | execute_authorization |
| message               | sender          | {senderAddress}       |
