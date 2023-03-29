<!--
order: 7
-->

# Parameters

The fee Split module contains the following parameters:

| Key                        | Type        | Default Value    |
| :------------------------- | :---------- | :--------------- |
| `EnableFeeShare`           | bool        | `true`           |
| `DeveloperShares`          | sdk.Dec     | `50%`            |
| `AllowedDenoms`            | []string{}  | `[]string(nil)`  |

## Enable FeeShare Module

The `EnableFeeShare` parameter toggles all state transitions in the module. When the parameter is disabled, it will prevent any transaction fees from being distributed to contract deplorers and it will disallow contract registrations, updates or cancellations.

### Developer Shares Amount

The `DeveloperShares` parameter is the percentage of transaction fees that are sent to the contract deplorers.

### Allowed Denominations

The `AllowedDenoms` parameter is used to specify which fees coins will be paid to contract developers. If this is empty, all fees paid will be split. If not, only fees specified here will be paid out to the withdrawal address.
