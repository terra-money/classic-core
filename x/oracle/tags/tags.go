package tags

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Oracle tags
var (
	ActionPriceUpdate   = "price-update"
	ActionVoteSubmitted = "vote-submitted"
	ActionTallyDropped  = "tally-dropped"
	ActionWhitelist     = "blacklist"
	ActionBlacklist     = "blacklist"

	Action = sdk.TagAction
	Denom  = "denom"
	Voter  = "voter"
	Power  = "power"
	Price  = "price"
)
