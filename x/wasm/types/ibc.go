package types

import (
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const portIDPrefix = "wasm."

// PortIDForContract build port ID from a contract address
func PortIDForContract(addr sdk.AccAddress) string {
	return portIDPrefix + addr.String()
}

// ContractFromPortID extract a contract address from port ID
func ContractFromPortID(portID string) (sdk.AccAddress, error) {
	if !strings.HasPrefix(portID, portIDPrefix) {
		return nil, sdkerrors.Wrapf(ErrInvalid, "without prefix")
	}
	return sdk.AccAddressFromBech32(portID[len(portIDPrefix):])
}
