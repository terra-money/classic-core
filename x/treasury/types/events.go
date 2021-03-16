package types

// Treasury module event types
const (
	EventTypePolicyUpdate       = "policy_update"
	EventTypeTaxRateUpdate      = "tax_rate_update"
	EventTypeRewardWeightUpdate = "reward_weight_update"

	AttributeKeyTaxRate      = "tax_rate"
	AttributeKeyRewardWeight = "reward_weight"
	AttributeKeyTaxCap       = "tax_cap"

	AttributeValueCategory = ModuleName
)
