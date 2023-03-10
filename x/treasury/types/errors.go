package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var ErrNoSuchBurnTaxExemptionAddress = sdkerrors.Register(ModuleName, 1, "no such address in extemption list")
