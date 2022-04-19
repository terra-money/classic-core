package keeper

import (
	"encoding/json"

	ibctransfertypes "github.com/cosmos/ibc-go/v3/modules/apps/transfer/types"
	channeltypes "github.com/cosmos/ibc-go/v3/modules/core/04-channel/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	wasmvmtypes "github.com/CosmWasm/wasmvm/types"

	"github.com/terra-money/core/x/wasm/types"
)

var _ types.IBCWasmQuerierInterface = IBCQuerier{}
var _ types.IBCWasmMsgParserInterface = IBCMsgParser{}

// IBCMsgParser - ibc msg parser for wasm msgs
type IBCMsgParser struct {
	portSource types.ICS20TransferPortSource
}

// NewIBCMsgParser returns ibc msg parser
func NewIBCMsgParser(portSource types.ICS20TransferPortSource) IBCMsgParser {
	return IBCMsgParser{portSource}
}

// Parse implements ibc msg parser
func (p IBCMsgParser) Parse(ctx sdk.Context, contractAddr sdk.AccAddress, wasmMsg wasmvmtypes.CosmosMsg) (sdk.Msg, error) {
	msg := wasmMsg.IBC

	if msg.CloseChannel != nil {
		cosmosMsg := channeltypes.NewMsgChannelCloseInit(
			types.PortIDForContract(contractAddr),
			msg.CloseChannel.ChannelID,
			contractAddr.String(),
		)

		return cosmosMsg, cosmosMsg.ValidateBasic()
	}

	if msg.Transfer != nil {
		token, err := types.ParseToCoin(msg.Transfer.Amount)
		if err != nil {
			return nil, err
		}

		cosmosMsg := ibctransfertypes.NewMsgTransfer(
			p.portSource.GetPort(ctx),
			msg.Transfer.ChannelID,
			token,
			contractAddr.String(),
			msg.Transfer.ToAddress,
			types.ConvertWasmIBCTimeoutHeightToCosmosHeight(msg.Transfer.Timeout.Block),
			msg.Transfer.Timeout.Timestamp,
		)

		return cosmosMsg, cosmosMsg.ValidateBasic()
	}

	return nil, sdkerrors.Wrap(types.ErrInvalidMsg, "Unknown variant of IBC")
}

// IBCQuerier - wasm query interface for wasm contract
type IBCQuerier struct {
	keeper        Keeper
	channelKeeper types.ChannelKeeper
}

// NewIBCQuerier returns wasm querier
func NewIBCQuerier(keeper Keeper, channelKeeper types.ChannelKeeper) IBCQuerier {
	return IBCQuerier{keeper, channelKeeper}
}

// Query - implement query function
func (querier IBCQuerier) Query(ctx sdk.Context, contractAddr sdk.AccAddress, request wasmvmtypes.QueryRequest) ([]byte, error) {
	if request.IBC.PortID != nil {
		contractInfo, err := querier.keeper.GetContractInfo(ctx, contractAddr)
		if err != nil {
			return nil, err
		}

		return json.Marshal(wasmvmtypes.PortIDResponse{
			PortID: contractInfo.IBCPortID,
		})
	}

	if request.IBC.Channel != nil {
		channelID := request.IBC.Channel.ChannelID
		portID := request.IBC.Channel.PortID
		if portID == "" {
			contractInfo, err := querier.keeper.GetContractInfo(ctx, contractAddr)
			if err != nil {
				return nil, err
			}

			portID = contractInfo.IBCPortID
		}

		got, found := querier.channelKeeper.GetChannel(ctx, portID, channelID)
		var channel *wasmvmtypes.IBCChannel

		// it must be in open state
		if found && got.State == channeltypes.OPEN {
			channel = &wasmvmtypes.IBCChannel{
				Endpoint: wasmvmtypes.IBCEndpoint{
					PortID:    portID,
					ChannelID: channelID,
				},
				CounterpartyEndpoint: wasmvmtypes.IBCEndpoint{
					PortID:    got.Counterparty.PortId,
					ChannelID: got.Counterparty.ChannelId,
				},
				Order:        got.Ordering.String(),
				Version:      got.Version,
				ConnectionID: got.ConnectionHops[0],
			}
		}

		return json.Marshal(wasmvmtypes.ChannelResponse{
			Channel: channel,
		})
	}

	if request.IBC.ListChannels != nil {
		portID := request.IBC.ListChannels.PortID
		channels := make(wasmvmtypes.IBCChannels, 0)
		querier.channelKeeper.IterateChannels(ctx, func(ch channeltypes.IdentifiedChannel) bool {
			// it must match the port and be in open state
			if (portID == "" || portID == ch.PortId) && ch.State == channeltypes.OPEN {
				newChan := wasmvmtypes.IBCChannel{
					Endpoint: wasmvmtypes.IBCEndpoint{
						PortID:    ch.PortId,
						ChannelID: ch.ChannelId,
					},
					CounterpartyEndpoint: wasmvmtypes.IBCEndpoint{
						PortID:    ch.Counterparty.PortId,
						ChannelID: ch.Counterparty.ChannelId,
					},
					Order:        ch.Ordering.String(),
					Version:      ch.Version,
					ConnectionID: ch.ConnectionHops[0],
				}
				channels = append(channels, newChan)
			}
			return false
		})

		return json.Marshal(wasmvmtypes.ListChannelsResponse{
			Channels: channels,
		})
	}

	return nil, wasmvmtypes.UnsupportedRequest{Kind: "unknown IBCQuery variant"}
}

// QueryCustom implements custom query interface
func (querier IBCQuerier) QueryCustom(ctx sdk.Context, data json.RawMessage) ([]byte, error) {
	return nil, wasmvmtypes.UnsupportedRequest{Kind: "unknown IBC variant"}
}
