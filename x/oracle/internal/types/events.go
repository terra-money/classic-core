//noalias
package types

// Oracle module event types
const (
	EventTypeExchangeRateUpdate = "exchange_rate_update"
	EventTypePrevote            = "prevote"
	EventTypeVote               = "vote"
	EventTypeFeedDeleate        = "feed_delegate"

	AttributeKeyDenom        = "denom"
	AttributeKeyVoter        = "voter"
	AttributeKeyExchangeRate = "exchange_rate"
	AttributeKeyOperator     = "operator"
	AttributeKeyFeeder       = "feeder"

	AttributeValueCategory = ModuleName
)
