package treasury

import sdk "github.com/cosmos/cosmos-sdk/types"

// nolint
var (
	KeyLunaTargetIssuance = []byte{0x01}
	KeyIncomePool         = []byte{0x02}
	PrefixClaim           = []byte{0x03}
)

// GetClaimKey is in format of prefixclaim||accaddress
func GetClaimKey(account sdk.AccAddress) []byte {
	return append(PrefixClaim, []byte(account.String())...)
}
