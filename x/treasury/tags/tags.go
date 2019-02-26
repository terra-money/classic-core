package tags

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Treasury tags
var (
	ActionSettle            = "settle"
	ActionReward            = "reward"
	ActionTaxUpdate         = "tax-update"
	ActionMinerRewardUpdate = "miner-reward-update"

	Action      = sdk.TagAction
	Denom       = "denom"
	Amount      = "amount"
	Tax         = "tax"
	MinerReward = "miner-weight"
	Oracle      = "oracle-reward"
	Budget      = "budget"
)
