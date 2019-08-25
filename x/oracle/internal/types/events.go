//noalias
package types

// Oracle module event types
const (
	EventTypePriceUpdate = "price_update"
	EventTypeSlash       = "slash"
	EventTypeLiveness    = "liveness"
	EventTypePrevote     = "prevote"
	EventTypeVote        = "vote"
	EventTypeFeedDeleate = "feed_delegate"

	AttributeKeyAddress     = "address"
	AttributeKeyHeight      = "height"
	AttributeKeyMissedVotes = "missed_votes"
	AttributeKeyDenom       = "denom"
	AttributeKeyVoter       = "voter"
	AttributeKeyPower       = "power"
	AttributeKeyPrice       = "price"
	AttributeKeyOperator    = "operator"
	AttributeKeyFeeder      = "feeder"

	AttributeValueCategory = ModuleName
)
