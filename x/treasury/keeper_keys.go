package treasury

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
)

const (
	//nolint default paramspace for treasury keeper
	DefaultParamspace = "treasury"
)

// nolint
var (
	KeyRewardWeight     = []byte("reward_weight")
	PrefixClaim         = []byte("claim")
	ParamStoreKeyParams = []byte("params")

	keyTaxRate        = []byte("tax_rate")
	prefixTaxProceeds = []byte("tax_proceeds")
	prefixTaxCap      = []byte("tax_cap")
	prefixIssuance    = []byte("issuance")
)

// KeyClaim is in format of prefixclaim:claimType:claimID
func KeyClaim(claimID string) []byte {
	return []byte(fmt.Sprintf("%s:%s", PrefixClaim, claimID))
}

func keyTaxProceeds(epoch sdk.Int) []byte {
	return []byte(fmt.Sprintf("%s:%s", prefixTaxProceeds, epoch))
}

func keyTaxCap(denom string) []byte {
	return []byte(fmt.Sprintf("%s:%s", prefixTaxCap, denom))
}

func keyIssuance(denom string, epoch sdk.Int) []byte {
	return []byte(fmt.Sprintf("%s:%s:%s", prefixIssuance, denom, epoch))
}

// ParamKeyTable for treasury module
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable(
		ParamStoreKeyParams, Params{},
	)
}
