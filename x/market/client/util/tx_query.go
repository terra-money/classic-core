package util

import (
	"fmt"

	"terra/x/market"
	tags "terra/x/market/tags"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
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

// QueryHistoryByTxQuery will query for votes via a direct txs tags query. It
// will fetch and build votes directly from the returned txs and return a JSON
// marshalled result or any error that occurred.
//
// NOTE: SearchTxs is used to facilitate the txs query which does not currently
// support configurable pagination.
func QueryHistoryByTxQuery(
	cdc *codec.Codec, cliCtx context.CLIContext, params QueryHistoryParams,
) ([]byte, error) {

	queryTags := []string{
		fmt.Sprintf("%s='%s'", tags.Action, market.SwapMsg{}.Type()),
	}

	if len(params.AskDenom) != 0 {
		queryTags = append(queryTags, fmt.Sprintf("%s='%s'", tags.AskDenom, []byte(params.AskDenom)))
	}

	if len(params.OfferDenom) != 0 {
		queryTags = append(queryTags, fmt.Sprintf("%s='%s'", tags.OfferDenom, []byte(params.OfferDenom)))
	}

	if len(params.TraderAddress) != 0 {
		queryTags = append(queryTags, fmt.Sprintf("%s='%s'", tags.Trader, []byte(params.TraderAddress)))
	}

	// NOTE: SearchTxs is used to facilitate the txs query which does not currently
	// support configurable pagination.
	infos, err := tx.SearchTxs(cliCtx, cdc, queryTags, defaultPage, defaultLimit)
	if err != nil {
		return nil, err
	}

	var swaps []market.SwapMsg

	for _, info := range infos {
		for _, msg := range info.Tx.GetMsgs() {
			if msg.Type() == (market.SwapMsg{}).Type() {
				swapMsg := msg.(market.SwapMsg)

				swaps = append(swaps, swapMsg)
			}
		}
	}

	if cliCtx.Indent {
		return cdc.MarshalJSONIndent(swaps, "", "  ")
	}

	return cdc.MarshalJSON(swaps)
}
