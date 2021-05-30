package rest

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	genutilrest "github.com/cosmos/cosmos-sdk/x/genutil/client/rest"

	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/types"
	"github.com/terra-money/core/client/lcd"
)

// TxQueryMaxHeightRange maximum allowed height range for /txs query
const TxQueryMaxHeightRange = 100

// QueryTxsRequestHandlerFn implements a REST handler that searches for transactions.
// Genesis transactions are returned if the height parameter is set to zero,
// otherwise the transactions are searched for by events.
func QueryTxsRequestHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			rest.WriteErrorResponse(
				w, http.StatusBadRequest,
				fmt.Sprintf("failed to parse query parameters: %s", err),
			)
			return
		}

		// Check public node option
		isPublicOpen := viper.GetBool(lcd.FlagPublic)

		// if the height query param is set to zero, query for genesis transactions
		heightStr := r.FormValue("height")
		if heightStr != "" {
			if height, err := strconv.ParseInt(heightStr, 10, 64); err == nil && height == 0 {
				if isPublicOpen {
					rest.WriteErrorResponse(
						w, http.StatusBadRequest,
						fmt.Sprintf("query genesis txs is not allowed for the public node"),
					)
				} else {
					genutilrest.QueryGenesisTxs(cliCtx, w)
				}

				return
			}
		}

		// enforce tx.height query parameter
		if isPublicOpen {
			if err := validateTxHeightRange(r); err != nil {
				rest.WriteErrorResponse(
					w, http.StatusBadRequest,
					err.Error(),
				)
				return
			}
		}

		var (
			events      []string
			txs         []sdk.TxResponse
			page, limit int
		)

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		if len(r.Form) == 0 {
			rest.PostProcessResponseBare(w, cliCtx, txs)
			return
		}

		events, page, limit, err = rest.ParseHTTPArgs(r)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		searchResult, err := utils.QueryTxsByEvents(cliCtx, events, page, limit)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		rest.PostProcessResponseBare(w, cliCtx, searchResult)
	}
}

func validateTxHeightRange(r *http.Request) error {
	txHeightStr := r.FormValue(types.TxHeightKey)
	txMinHeightStr := r.FormValue(rest.TxMinHeightKey)
	txMaxHeightStr := r.FormValue(rest.TxMaxHeightKey)

	if txHeightStr == "" && (txMinHeightStr == "" || txMaxHeightStr == "") {
		return fmt.Errorf(
			"it is not allowed to query txs without %s or (%s && %s) options. please refer {URL}/swagger-ui",
			types.TxHeightKey, rest.TxMaxHeightKey, rest.TxMinHeightKey)
	}

	if txMinHeightStr != "" && txMaxHeightStr != "" {
		txMinHeight, err := strconv.ParseInt(txMinHeightStr, 10, 64)
		if err != nil {
			return err
		}

		txMaxHeight, err := strconv.ParseInt(txMaxHeightStr, 10, 64)
		if err != nil {
			return err
		}

		if txMaxHeight < txMinHeight {
			return fmt.Errorf("%s must be bigger than %s", rest.TxMaxHeightKey, rest.TxMinHeightKey)
		}

		if txMaxHeight-txMinHeight > TxQueryMaxHeightRange {
			return fmt.Errorf("tx height range must be smaller than %d", TxQueryMaxHeightRange)
		}
	}

	return nil
}
