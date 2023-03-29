<!--
order: 6
-->

# Events

The `x/feeshare` module emits the following events:

## Register Fee Split

| Type                 | Attribute Key          | Attribute Value           |
| :------------------- | :--------------------- | :------------------------ |
| `register_feeshare`  | `"contract"`            | `{msg.ContractAddress}`   |
| `register_feeshare`  | `"sender"`              | `{msg.DeployerAddress}`   |
| `register_feeshare`  | `"withdrawer_address"`  | `{msg.WithdrawerAddress}` |

## Update Fee Split

| Type               | Attribute Key          | Attribute Value           |
| :----------------- | :--------------------- | :------------------------ |
| `update_feeshare`  | `"contract"`            | `{msg.ContractAddress}`   |
| `update_feeshare`  | `"sender"`              | `{msg.DeployerAddress}`   |
| `update_feeshare`  | `"withdrawer_address"`  | `{msg.WithdrawerAddress}` |

## Cancel Fee Split

| Type               | Attribute Key | Attribute Value         |
| :----------------- | :------------ | :---------------------- |
| `cancel_feeshare`  | `"contract"`   | `{msg.ContractAddress}` |
| `cancel_feeshare`  | `"sender"`     | `{msg.DeployerAddress}` |
