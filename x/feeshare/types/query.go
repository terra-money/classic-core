package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// ValidateBasic runs stateless checks on the query requests
func (q QueryFeeShareRequest) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(q.ContractAddress); err != nil {
		return sdkerrors.Wrapf(err, "invalid contract address %s", q.ContractAddress)
	}

	return nil
}

// ValidateBasic runs stateless checks on the query requests
func (q QueryDeployerFeeSharesRequest) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(q.DeployerAddress); err != nil {
		return sdkerrors.Wrapf(err, "invalid deployer address %s", q.DeployerAddress)
	}

	return nil
}

// ValidateBasic runs stateless checks on the query requests
func (q QueryWithdrawerFeeSharesRequest) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(q.WithdrawerAddress); err != nil {
		return sdkerrors.Wrapf(err, "invalid withdraw address %s", q.WithdrawerAddress)
	}

	return nil
}
