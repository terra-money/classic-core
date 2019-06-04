package oracle

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
)

var (
	prefixPrevote          = []byte("prevote")
	prefixVote             = []byte("vote")
	prefixPrice            = []byte("price")
	prefixDropCounter      = []byte("drop")
	paramStoreKeyParams    = []byte("params")
	prefixFeederDelegation = []byte("feederdelegation")
	prefixClaim   = []byte("claim")

	keySwapFeePool = []byte("swapfeepool")
)

func keyPrevote(denom string, voter sdk.ValAddress) []byte {
	return []byte(fmt.Sprintf("%s:%s:%s", prefixPrevote, denom, voter))
}

func keyVote(denom string, voter sdk.ValAddress) []byte {
	return []byte(fmt.Sprintf("%s:%s:%s", prefixVote, denom, voter))
}

func keyPrice(denom string) []byte {
	return []byte(fmt.Sprintf("%s:%s", prefixPrice, denom))
}

func keyClaim(recipient sdk.AccAddress) []byte {
	return []byte(fmt.Sprintf("%s:%s", prefixClaim, recipient))
}

func keyDropCounter(denom string) []byte {
	return []byte(fmt.Sprintf("%s:%s", prefixDropCounter, denom))
}

func paramKeyTable() params.KeyTable {
	return params.NewKeyTable(
		paramStoreKeyParams, Params{},
	)
}

func keyFeederDelegation(delegate sdk.ValAddress) []byte {
	return []byte(fmt.Sprintf("%s:%s", prefixFeederDelegation, delegate))
}
