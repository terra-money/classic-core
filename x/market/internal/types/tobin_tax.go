package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// TobinTax - struct to store tobin tax for the specific denom with high volatility
type TobinTax struct {
	Denom   string  `json:"denom" yaml:"denom"`
	TaxRate sdk.Dec `json:"tax_rate" yaml:"tax_rate"`
}

// String implements fmt.Stringer interface
func (tt TobinTax) String() string {
	return fmt.Sprintf(`TobinTax
	Denom:      %s, 
	TaxRate:    %s`,
		tt.Denom, tt.TaxRate)
}

// TobinTaxList is convience wrapper to handle TobinTax array
type TobinTaxList []TobinTax

// String implements fmt.Stringer interface
func (ttl TobinTaxList) String() (out string) {
	out = ""
	for _, tt := range ttl {
		out += tt.String() + "\n"
	}

	return
}
