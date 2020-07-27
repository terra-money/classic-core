package types

// msgauth module events
const (
	EventGrantAuthorization   = "grant_authorization"
	EventRevokeAuthorization  = "revoke_authorization"
	EventExecuteAuthorization = "execute_authorization"

	AttributeKeyGrantType      = "grant_type"
	AttributeKeyGranteeAddress = "grantee"
	AttributeKeyGranterAddress = "granter"

	AttributeValueCategory = ModuleName
)
