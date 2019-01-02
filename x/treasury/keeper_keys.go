package treasury

import "fmt"

// nolint
var (
	KeyLunaTargetIssuance = []byte{0x01}
	KeyIncomePool         = []byte("income_pool")
	PrefixClaim           = []byte("claim")
	PrefixShare           = []byte("share")
)

// GetIncomePoolKey returns the key for the income pool
func GetIncomePoolKey() []byte {
	return KeyIncomePool
}

// GetShareKey is in format of prefixshare:shareID
func GetShareKey(shareID string) []byte {
	return []byte(fmt.Sprintf("%s:%s", PrefixShare, shareID))
}

// GetClaimsForSharePrefix is in format of prefixclaim:shareID
func GetClaimsForSharePrefix(shareID string) []byte {
	return []byte(fmt.Sprintf("%s:%s", PrefixClaim, shareID))
}

// GetClaimKey is in format of prefixclaim:shareID:claimID
func GetClaimKey(shareID string, claimID string) []byte {
	return []byte(fmt.Sprintf("%s:%s:%s", PrefixClaim, shareID, claimID))
}
