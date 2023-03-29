package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerror "github.com/cosmos/cosmos-sdk/types/errors"
)

// NewFeeShare returns an instance of FeeShare.
func NewFeeShare(contract sdk.Address, deployer, withdrawer sdk.AccAddress) FeeShare {
	return FeeShare{
		ContractAddress:   contract.String(),
		DeployerAddress:   deployer.String(),
		WithdrawerAddress: withdrawer.String(),
	}
}

// GetContractAddr returns the contract address
func (fs FeeShare) GetContractAddr() sdk.Address {
	contract, err := sdk.AccAddressFromBech32(fs.ContractAddress)
	if err != nil {
		return nil
	}
	return contract
}

// GetDeployerAddr returns the contract deployer address
func (fs FeeShare) GetDeployerAddr() sdk.AccAddress {
	contract, err := sdk.AccAddressFromBech32(fs.DeployerAddress)
	if err != nil {
		return nil
	}
	return contract
}

// GetWithdrawerAddr returns the account address to where the funds proceeding
// from the fees will be received.
func (fs FeeShare) GetWithdrawerAddr() sdk.AccAddress {
	contract, err := sdk.AccAddressFromBech32(fs.WithdrawerAddress)
	if err != nil {
		return nil
	}
	return contract
}

// Validate performs a stateless validation of a FeeShare
func (fs FeeShare) Validate() error {
	if _, err := sdk.AccAddressFromBech32(fs.ContractAddress); err != nil {
		return err
	}

	if _, err := sdk.AccAddressFromBech32(fs.DeployerAddress); err != nil {
		return err
	}

	if fs.WithdrawerAddress == "" {
		return sdkerror.Wrap(sdkerror.ErrInvalidAddress, "withdrawer address cannot be empty")
	}

	if _, err := sdk.AccAddressFromBech32(fs.WithdrawerAddress); err != nil {
		return err
	}

	return nil
}
