//noalias
package types

// Oracle module event types
const (
	EventTypeExchangeRateUpdate = "exchangerate_update"
	EventTypePrevote            = "prevote"
	EventTypeVote               = "vote"
	EventTypeFeedDeleate        = "feed_delegate"

	AttributeKeyDenom        = "denom"
	AttributeKeyVoter        = "voter"
	AttributeKeyExchangeRate = "exchangerate"
	AttributeKeyOperator     = "operator"
	AttributeKeyFeeder       = "feeder"

	AttributeValueCategory = ModuleName
)
