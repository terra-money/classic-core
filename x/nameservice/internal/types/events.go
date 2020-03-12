//noalias
package types

// Nameservice module event types
const (
	EventTypeOpen       = "open"
	EventTypeBid        = "bid"
	EventTypeReveal     = "reveal"
	EventTypeRegister   = "register"
	EventTypeUnregister = "unregister"
	EventTypeRenew      = "renew"

	AttributeKeyName      = "name"
	AttributeKeyDeposit   = "deposit"
	AttributeKeyBidder    = "bidder"
	AttributeKeyAmount    = "amount"
	AttributeKeyAddress   = "address"
	AttributeKeyOrganizer = "organizer"
	AttributeKeyEndTime   = "end_time"
	AttributeKeyOwner     = "owner"
	AttributeKeyNewOwner  = "new_owner"
	AttributeKeyFee       = "fee"

	AttributeValueCategory = ModuleName
)
