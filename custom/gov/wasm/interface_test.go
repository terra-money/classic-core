package wasm

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto/secp256k1"

	wasmvmtypes "github.com/CosmWasm/wasmvm/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

func TestEncoding(t *testing.T) {
	addrs := []sdk.AccAddress{
		sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address()),
		sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address()),
		sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address()),
	}

	cases := map[string]struct {
		sender sdk.AccAddress
		input  wasmvmtypes.CosmosMsg
		// set if valid
		output sdk.Msg
		// set if invalid
		isError bool
	}{
		"yes vote": {
			sender: addrs[0],
			input: wasmvmtypes.CosmosMsg{
				Gov: &wasmvmtypes.GovMsg{
					Vote: &wasmvmtypes.VoteMsg{
						ProposalId: 1,
						Vote:       wasmvmtypes.Yes,
					},
				},
			},
			output: &govtypes.MsgVote{
				Voter:      addrs[0].String(),
				ProposalId: 1,
				Option:     govtypes.OptionYes,
			},
		},
		"no vote": {
			sender: addrs[0],
			input: wasmvmtypes.CosmosMsg{
				Gov: &wasmvmtypes.GovMsg{
					Vote: &wasmvmtypes.VoteMsg{
						ProposalId: 1,
						Vote:       wasmvmtypes.No,
					},
				},
			},
			output: &govtypes.MsgVote{
				Voter:      addrs[0].String(),
				ProposalId: 1,
				Option:     govtypes.OptionNo,
			},
		},
		"no_with_veto vote": {
			sender: addrs[0],
			input: wasmvmtypes.CosmosMsg{
				Gov: &wasmvmtypes.GovMsg{
					Vote: &wasmvmtypes.VoteMsg{
						ProposalId: 1,
						Vote:       wasmvmtypes.NoWithVeto,
					},
				},
			},
			output: &govtypes.MsgVote{
				Voter:      addrs[0].String(),
				ProposalId: 1,
				Option:     govtypes.OptionNoWithVeto,
			},
		},
		"abstain vote": {
			sender: addrs[0],
			input: wasmvmtypes.CosmosMsg{
				Gov: &wasmvmtypes.GovMsg{
					Vote: &wasmvmtypes.VoteMsg{
						ProposalId: 1,
						Vote:       wasmvmtypes.Abstain,
					},
				},
			},
			output: &govtypes.MsgVote{
				Voter:      addrs[0].String(),
				ProposalId: 1,
				Option:     govtypes.OptionAbstain,
			},
		},
	}

	parser := NewWasmMsgParser()
	for name, tc := range cases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			res, err := parser.Parse(tc.sender, tc.input)
			if tc.isError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.output, res)
			}
		})
	}

}
