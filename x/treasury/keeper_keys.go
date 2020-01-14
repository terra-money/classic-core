package treasury

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
)

// nolint
var (
	prefixRewardWeight  = []byte("reward_weight")
	PrefixClaim         = []byte("claim")
	paramStoreKeyParams = []byte("params")

	prefixTaxRate     = []byte("tax_rate")
	prefixTaxProceeds = []byte("tax_proceeds")
	prefixTaxCap      = []byte("tax_cap")
	prefixIssuance    = []byte("issuance")
)

func keyTaxRate(epoch sdk.Int) []byte {
	return []byte(fmt.Sprintf("%s:%s", prefixTaxRate, epoch))
}

func keyRewardWeight(epoch sdk.Int) []byte {
	return []byte(fmt.Sprintf("%s:%s", prefixRewardWeight, epoch))
}

func keyClaim(claimID string) []byte {
	return []byte(fmt.Sprintf("%s:%s", PrefixClaim, claimID))
}

func keyTaxProceeds(epoch sdk.Int) []byte {
	return []byte(fmt.Sprintf("%s:%s", prefixTaxProceeds, epoch))
}

func keyTaxCap(denom string) []byte {
	return []byte(fmt.Sprintf("%s:%s", prefixTaxCap, denom))
}

func paramKeyTable() params.KeyTable {
	return params.NewKeyTable(
		paramStoreKeyParams, Params{},
	)
}
