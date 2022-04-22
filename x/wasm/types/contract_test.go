package types

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto/ed25519"
)

var (
	creatorPk    = ed25519.GenPrivKey().PubKey()
	contractPk   = ed25519.GenPrivKey().PubKey()
	creatorAddr  = sdk.AccAddress(creatorPk.Address())
	contractAddr = sdk.AccAddress(contractPk.Address())
)

func TestNewCodeInfo(t *testing.T) {
	codeInfo := CodeInfo{
		CodeID:   1,
		CodeHash: []byte{1, 2, 3},
		Creator:  creatorAddr.String(),
	}
	require.Equal(t, codeInfo, NewCodeInfo(1, []byte{1, 2, 3}, creatorAddr))
}

func TestNewContractInfo(t *testing.T) {
	contractInfo := ContractInfo{
		Address: contractAddr.String(),
		CodeID:  1,
		Creator: creatorAddr.String(),
		Admin:   creatorAddr.String(),
		InitMsg: []byte{},
	}
	require.Equal(t, contractInfo, NewContractInfo(1, contractAddr, creatorAddr, creatorAddr, []byte{}))
}
