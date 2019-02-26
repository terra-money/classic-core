package treasury

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/x/params"
)

const (
	//nolint default paramspace for treasury keeper
	DefaultParamspace = "treasury"
)

// nolint
var (
	KeyRewardWeight = []byte("reward_weight")
	KeyIncomePool   = []byte("income_pool")

	PrefixIssuance = []byte("issuance")
	PrefixClaim    = []byte("claim")

	ParamStoreKeyParams = []byte("params")
)

// KeyClaim is in format of prefixclaim:claimType:claimID
func KeyClaim(claimID string) []byte {
	return []byte(fmt.Sprintf("%s:%s", PrefixClaim, claimID))
}

// KeyIssuance is in format of PrefixIssuance:denom
func KeyIssuance(denom string) []byte {
	return []byte(fmt.Sprintf("%s:%s", PrefixIssuance, denom))
}

// ParamKeyTable for treasury module
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable(
		ParamStoreKeyParams, Params{},
	)
}
