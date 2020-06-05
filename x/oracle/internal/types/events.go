//noalias
package types

// Oracle module event types
const (
	EventTypeExchangeRateUpdate      = "exchange_rate_update"
	EventTypeCrossExchangeRateUpdate = "cross_exchange_rate_update"
	EventTypePrevote                 = "prevote"
	EventTypeVote                    = "vote"
	EventTypeFeedDelegate            = "feed_delegate"
	EventTypeAggregatePrevote        = "aggregate_prevote"
	EventTypeAggregateVote           = "aggregate_vote"

	AttributeKeyDenom             = "denom"
	AttributeKeyDenom1            = "denom1"
	AttributeKeyDenom2            = "denom2"
	AttributeKeyVoter             = "voter"
	AttributeKeyExchangeRate      = "exchange_rate"
	AttributeKeyExchangeRates     = "exchange_rates"
	AttributeKeyCrossExchangeRate = "cross_exchange_rate"
	AttributeKeyOperator          = "operator"
	AttributeKeyFeeder            = "feeder"

	AttributeValueCategory = ModuleName
)
