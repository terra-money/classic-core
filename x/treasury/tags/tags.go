package tags

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Treasury tags
var (
	ActionSettle       = "settle"
	ActionPolicyUpdate = "policy-update"

	Action      = sdk.TagAction
	Denom       = "denom"
	Amount      = "amount"
	Rewardee    = "rewardee"
	Tax         = "tax"
	TaxCap      = "tax-cap"
	Class       = "class"
	MinerReward = "miner-weight"
	Oracle      = "oracle-reward"
	Budget      = "budget"
)
