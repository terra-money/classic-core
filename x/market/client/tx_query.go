package client

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/gov"
	"github.com/cosmos/cosmos-sdk/x/gov/tags"
)

const (
	defaultPage  = 1
	defaultLimit = 30 // should be consistent with tendermint/tendermint/rpc/core/pipe.go:19
)


type QueryHistoryParams struct {
	TraderAddress sdk.AccAddress
	AskDenom      string
	OfferDenom    string
}

// QueryVotesByTxQuery will query for votes via a direct txs tags query. It
// will fetch and build votes directly from the returned txs and return a JSON
// marshalled result or any error that occurred.
//
// NOTE: SearchTxs is used to facilitate the txs query which does not currently
// support configurable pagination.
func QueryHistoryByTxQuery(
	cdc *codec.Codec, cliCtx context.CLIContext, params QueryHistoryParams,
) ([]byte, error) {

	tags := []string{
		fmt.Sprintf("%s='%s'", tags.Action, tags.ActionSwap)
	}

	if len(params.AskDenom) != 0 {
		tags = append(tags, fmt.Sprintf("%s='%s'", tags.AskDenom, []byte(params.AskDenom)))
	}

	if len(params.OfferDenom) != 0 {
		tags = append(tags, fmt.Sprintf("%s='%s'", params.OfferDenom, []byte(params.OfferDenom)))
	}

	// NOTE: SearchTxs is used to facilitate the txs query which does not currently
	// support configurable pagination.
	infos, err := tx.SearchTxs(cliCtx, cdc, tags, defaultPage, defaultLimit)
	if err != nil {
		return nil, err
	}

	var swaps []market.MsgSwap

	for _, info := range infos {
		for _, msg := range info.Tx.GetMsgs() {
			if msg.Type() == "swap" {
				swapMsg := msg.(market.MsgSwap)

				swaps = append(swaps, swapMsg)
			}
		}
	}

	if cliCtx.Indent {
		return cdc.MarshalJSONIndent(votes, "", "  ")
	}

	return cdc.MarshalJSON(votes)
}
