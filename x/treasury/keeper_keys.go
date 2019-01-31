package treasury

import "fmt"

// nolint
var (
	KeyLunaTargetIssuance = []byte("luna_target")
	KeyIncomePool         = []byte("income_pool")
	PrefixClaim           = []byte("claim")
	KeyClaimsTally        = []byte("claims_tally")
)

// GetClaimKey is in format of prefixclaim:shareID:claimID
func GetClaimKey(claimID string) []byte {
	return []byte(fmt.Sprintf("%s:%s", PrefixClaim, claimID))
}
