package types

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/stretchr/testify/require"

	"github.com/tendermint/tendermint/crypto/ed25519"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	core "github.com/terra-money/core/types"

	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
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

func TestNewEnv(t *testing.T) {
	ctx := sdk.NewContext(nil, tmproto.Header{
		Height: 100,
		Time:   time.Now(),
	}, false, nil)

	require.NotPanics(t, func() {
		_ = NewEnv(ctx, sdk.AccAddress{})
		_ = NewEnv(WithTXCounter(ctx, 100), sdk.AccAddress{})
	})
}

func TestNewInfo(t *testing.T) {
	_, _, addr := keyPubAddr()

	deposit := sdk.NewCoins(sdk.NewCoin(core.MicroLunaDenom, sdk.NewInt(10)))
	info := NewInfo(addr, deposit)
	require.Equal(t, addr.String(), info.Sender)
	require.Equal(t, wasmvmtypes.Coins{wasmvmtypes.NewCoin(10, core.MicroLunaDenom)}, info.Funds)
}

func TestGenerateContractAddress(t *testing.T) {
	addr := GenerateContractAddress(1, 1)
	require.Equal(t,
		[]byte{0x3b, 0x1a, 0x74, 0x85, 0xc6, 0x16, 0x2c, 0x58, 0x83, 0xee, 0x45, 0xfb, 0x2d, 0x74, 0x77, 0xa8, 0x7d, 0x8a, 0x4c, 0xe5},
		addr.Bytes(),
	)
}
