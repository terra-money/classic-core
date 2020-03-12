package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/terra-project/core/types"
	"sort"
)

// RenewalFee - holds required fee amount for the specific name length
type RenewalFee struct {
	Length int      `json:"length" yaml:"length"`
	Amount sdk.Coin `json:"amount" yaml:"amount"`
}

// String implements fmt.Stringer interface
func (fee RenewalFee) String() string {
	return fmt.Sprintf(`RenewalFee
	Length: %d,
	Amount: %s`,
		fee.Length, fee.Amount)
}

// RenewalFees gives default operation for renewal fee handling
type RenewalFees []RenewalFee

// RenewalFeeForLength return proper fee amount for the given length
func (fees RenewalFees) RenewalFeeForLength(n int) sdk.Coin {
	if len(fees) == 0 {
		return sdk.NewCoin(types.MicroSDRDenom, sdk.ZeroInt())
	}

	sort.Sort(fees)

	for _, fee := range fees {
		// return amount of equal or bigger length item
		if fee.Length >= n {
			return fee.Amount
		}
	}

	// return last item
	return fees[len(fees)-1].Amount
}

// String implements fmt.Stringer interface
func (fees RenewalFees) String() (out string) {
	for _, fee := range fees {
		out += fee.String() + "\n"
	}
	return
}

// Len implements sort.Interface
func (fees RenewalFees) Len() int {
	return len(fees)
}

// Less reports whether the element with
// index i should sort before the element with index j.
func (fees RenewalFees) Less(i, j int) bool {
	return fees[i].Length <= fees[j].Length
}

// Swap implements sort.Interface.
func (fees RenewalFees) Swap(i, j int) {
	fees[i], fees[j] = fees[j], fees[i]
}
