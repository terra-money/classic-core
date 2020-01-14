package tags

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Oracle tags
var (
	ActionPriceUpdate  = "price-update"  // normal cases
	ActionTallyDropped = "tally-dropped" // emitted when price update is illiquid

	Action = sdk.TagAction
	Denom  = "denom"
	Voter  = "voter"
	Power  = "power"
	Price  = "price"

	Operator     = "operator"
	FeedDelegate = "feed_delegate"
)
