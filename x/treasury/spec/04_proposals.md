## Proposals

The Treasury module defines special proposals which allow the [Tax Rate](./02_state.md#TaxRate) and [Reward Weight](./02_state.md#RewardWeight) values in the `KVStore` to be voted on and changed accordingly, subject to the [policy constraints](./03_end_block.md#PolicyConstraints) imposed by `pc.Clamp()`.

### TaxRateUpdateProposal

```go
type TaxRateUpdateProposal struct {
	Title       string  // Title of the Proposal
	Description string  // Description of the Proposal
	TaxRate     sdk.Dec // target TaxRate
}
```

::: details JSON Example

```json
{
  "type": "treasury/TaxRateUpdateProposal",
  "value": {
    "title": "proposal title",
    "description": "proposal description",
    "tax_rate": "0.001000000000000000"
  }
}
```

### RewardWeightUpdateProposal

```go
type RewardWeightUpdateProposal struct {
	Title        string  // Title of the Proposal
	Description  string  // Description of the Proposal
	RewardWeight sdk.Dec // target RewardWeight
}
```

::: details JSON Example

```json
{
  "type": "treasury/RewardWeightUpdateProposal",
  "value": {
    "title": "proposal title",
    "description": "proposal description",
    "reward_weight": "0.001000000000000000"
  }
}
```