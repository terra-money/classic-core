package rest

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"

	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/types"
	"github.com/terra-project/core/client/lcd"
)

// QueryTxsHandlerFn implements a REST handler that searches for transactions.
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

		// if the height query param is set to zero, query for genesis transactions
		heightStr := r.FormValue("height")
		if heightStr != "" {
			if height, err := strconv.ParseInt(heightStr, 10, 64); err == nil && height == 0 {
				rest.WriteErrorResponse(
					w, http.StatusBadRequest,
					fmt.Sprintf("query genesis txs is not allowed for the public node"),
				)
				return
			}
		}

		txHeightStr := r.FormValue(types.TxHeightKey)

		// enforce tx.height query parameter
		isPublicOpen := viper.GetBool(lcd.FlagPublic)
		if isPublicOpen {
			if txHeightStr == "" {
				rest.WriteErrorResponse(
					w, http.StatusBadRequest,
					fmt.Sprint("it is not allowed to query txs without tx.height option. please refer {URL}/swagger-ui"),
				)
				return
			}
		}

		// parse tx.height query parameter
		var txHeightEvents []string
		if _, err := strconv.ParseInt(txHeightStr, 10, 64); len(txHeightStr) != 0 && err != nil {
			// remove query parameter to prevent duplicated handling
			delete(r.Form, types.TxHeightKey)

			txHeightEvents, err = parseHeightRange(txHeightStr)
			if err != nil {
				rest.WriteErrorResponse(
					w, http.StatusBadRequest,
					fmt.Sprintf("failed to parse %s: %s", types.TxHeightKey, err.Error()),
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

		//append tx height events
		if txHeightEvents != nil && len(txHeightEvents) > 0 {
			events = append(events, txHeightEvents...)
		}

		searchResult, err := utils.QueryTxsByEvents(cliCtx, events, page, limit)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		rest.PostProcessResponseBare(w, cliCtx, searchResult)
	}
}

func parseHeightRange(txHeightStr string) (events []string, err error) {
	queries := strings.Split(txHeightStr, ",")
	if len(queries) != 2 {
		err = fmt.Errorf("invalid tx.height options: %s", queries)
		return
	}

	var upperBound int64
	var lowerBound int64
	for _, query := range queries {
		switch {
		case strings.HasPrefix(query, "GTE"):
			if h, err2 := strconv.ParseInt(query[3:], 10, 64); err2 == nil {
				lowerBound = h
				events = append(events, fmt.Sprintf("%s>=%d", types.TxHeightKey, h))
			} else {
				err = fmt.Errorf("failed to parse integer: %s", query[3:])
				return
			}
		case strings.HasPrefix(query, "GT"):
			if h, err2 := strconv.ParseInt(query[2:], 10, 64); err2 == nil {
				lowerBound = h + 1
				events = append(events, fmt.Sprintf("%s>%d", types.TxHeightKey, h))
			} else {
				err = fmt.Errorf("failed to parse integer: %s", query[2:])
				return
			}
		case strings.HasPrefix(query, "LTE"):
			if h, err2 := strconv.ParseInt(query[3:], 10, 64); err2 == nil {
				upperBound = h
				events = append(events, fmt.Sprintf("%s<=%d", types.TxHeightKey, h))
			} else {
				err = fmt.Errorf("failed to parse integer: %s", query[3:])
				return
			}
		case strings.HasPrefix(query, "LT"):
			if h, err2 := strconv.ParseInt(query[2:], 10, 64); err2 == nil {
				upperBound = h - 1
				events = append(events, fmt.Sprintf("%s<%d", types.TxHeightKey, h))
			} else {
				err = fmt.Errorf("failed to parse integer: %s", query[2:])
				return
			}
		default:
			err = fmt.Errorf("invalid operator: %s", query)
			return
		}
	}

	boundDiff := upperBound - lowerBound
	if boundDiff > 100 {
		err = fmt.Errorf("max allowed tx.height range gap is 100: %d", boundDiff)
		return
	} else if boundDiff <= 0 {
		err = fmt.Errorf("tx.height range gap should be positive: %d", boundDiff)
		return
	}

	return
}
