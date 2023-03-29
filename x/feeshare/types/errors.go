package types

import (
	sdkerrrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// errors
var (
	ErrFeeShareDisabled              = sdkerrrors.Register(ModuleName, 1, "feeshare module is disabled by governance")
	ErrFeeShareAlreadyRegistered     = sdkerrrors.Register(ModuleName, 2, "feeshare already exists for given contract")
	ErrFeeShareNoContractDeployed    = sdkerrrors.Register(ModuleName, 3, "no contract deployed")
	ErrFeeShareContractNotRegistered = sdkerrrors.Register(ModuleName, 4, "no feeshare registered for contract")
	ErrFeeSharePayment               = sdkerrrors.Register(ModuleName, 5, "feeshare payment error")
)
