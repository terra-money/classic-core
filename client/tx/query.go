package tx

import (
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"

	"github.com/terra-project/core/types"
)

// QueryTxsByTagsRequestHandlerFn implements a REST handler that searches for
// transactions by tags.
func QueryTxsByTagsRequestHandlerFn(cliCtx context.CLIContext, cdc *codec.Codec) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			tags        []string
			txs         []sdk.TxResponse
			page, limit int
		)

		err := r.ParseForm()
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest,
				sdk.AppendMsgToErr("could not parse query parameters", err.Error()))
			return
		}

		if len(r.Form) == 0 {
			rest.PostProcessResponse(w, cdc, txs, cliCtx.Indent)
			return
		}

		tags, page, limit, err = rest.ParseHTTPArgs(r)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		txs, totalCount, err := SearchTxs(cliCtx, cdc, tags, page, limit)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		rest.PostProcessResponse(w, cdc, types.NewTxSearchResponse(txs, totalCount), cliCtx.Indent)
	}
}
