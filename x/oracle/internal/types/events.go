//noalias
package types

// Oracle module event types
const (
	EventTypePriceUpdate = "price_update"
	EventTypePrevote     = "prevote"
	EventTypeVote        = "vote"
	EventTypeFeedDeleate = "feed_delegate"

	AttributeKeyDenom    = "denom"
	AttributeKeyVoter    = "voter"
	AttributeKeyPrice    = "price"
	AttributeKeyOperator = "operator"
	AttributeKeyFeeder   = "feeder"

	AttributeValueCategory = ModuleName
)
