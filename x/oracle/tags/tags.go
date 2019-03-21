package tags

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Oracle tags
var (
	ActionPriceUpdate  = "price-update"  // normal cases
	ActionTallyDropped = "tally-dropped" // emitted when price update is illiquid
	ActionWhitelist    = "whitelist"     // emitted on virgin listing
	ActionBlacklist    = "blacklist"     // emitted on delisting

	Action = sdk.TagAction
	Denom  = "denom"
	Voter  = "voter"
	Power  = "power"
	Price  = "price"
)
