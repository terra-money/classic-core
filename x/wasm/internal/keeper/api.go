package keeper

import (
	"fmt"
	cosmwasm "github.com/confio/go-cosmwasm"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func humanAddress(canon []byte) (string, error) {
	if len(canon) != sdk.AddrLen {
		return "", fmt.Errorf("Expected %d byte address", sdk.AddrLen)
	}
	return sdk.AccAddress(canon).String(), nil
}

func canonicalAddress(human string) ([]byte, error) {
	return sdk.AccAddressFromBech32(human)
}

var cosmwasmAPI = cosmwasm.GoAPI{
	HumanAddress:     humanAddress,
	CanonicalAddress: canonicalAddress,
}
