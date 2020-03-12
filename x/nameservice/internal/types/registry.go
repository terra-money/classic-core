package types

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Registry - data struct to hold name information
type Registry struct {
	Name    Name           `json:"name" yaml:"name"`
	Owner   sdk.AccAddress `json:"owner" yaml:"owner"`
	EndTime time.Time      `json:"end_time" yaml:"end_time"`
}

// NewRegistry returns Registry instance
func NewRegistry(name Name, owner sdk.AccAddress, endTime time.Time) Registry {
	return Registry{
		Name:    name,
		Owner:   owner,
		EndTime: endTime,
	}
}

// String implements fmt.Stringer interface
func (r Registry) String() string {
	return fmt.Sprintf(`Registry
Name:    %s
Owner:   %s
EndTime: %s 
`, r.Name, r.Owner.String(), r.EndTime)
}
