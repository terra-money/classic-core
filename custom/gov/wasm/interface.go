package wasm

import (
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/gov/types"

	wasmvmtypes "github.com/CosmWasm/wasmvm/types"

	wasm "github.com/terra-money/core/x/wasm/exported"
)

var _ wasm.WasmMsgParserInterface = WasmMsgParser{}

// WasmMsgParser - wasm msg parser for staking msgs
type WasmMsgParser struct{}

// NewWasmMsgParser returns bank wasm msg parser
func NewWasmMsgParser() WasmMsgParser {
	return WasmMsgParser{}
}

// Parse implements wasm staking msg parser
func (WasmMsgParser) Parse(contractAddr sdk.AccAddress, wasmMsg wasmvmtypes.CosmosMsg) (sdk.Msg, error) {
	msg := wasmMsg.Gov

	var option types.VoteOption
	switch msg.Vote.Vote {
	case wasmvmtypes.Yes:
		option = types.OptionYes
	case wasmvmtypes.No:
		option = types.OptionNo
	case wasmvmtypes.NoWithVeto:
		option = types.OptionNoWithVeto
	case wasmvmtypes.Abstain:
		option = types.OptionAbstain
	}

	cosmosMsg := &types.MsgVote{
		ProposalId: msg.Vote.ProposalId,
		Voter:      contractAddr.String(),
		Option:     option,
	}

	return cosmosMsg, cosmosMsg.ValidateBasic()
}

// ParseCustom implements custom parser
func (WasmMsgParser) ParseCustom(_ sdk.AccAddress, _ json.RawMessage) (sdk.Msg, error) {
	return nil, nil
}
