package types

import (
	"encoding/json"

	wasmTypes "github.com/CosmWasm/go-cosmwasm/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// WasmQuerier - query registration interface for other modules
type WasmQuerier interface {
	Query(ctx sdk.Context, request wasmTypes.QueryRequest) ([]byte, sdk.Error)
	QueryCustom(ctx sdk.Context, data json.RawMessage) ([]byte, sdk.Error)
}

// Querier - wasm query handler
type Querier struct {
	Ctx      sdk.Context
	Queriers map[string]WasmQuerier
}

// WasmCustomQuery - wasm custom query
type WasmCustomQuery struct {
	Route     string          `json:"route"`
	QueryData json.RawMessage `json:"query_data"`
}

var _ wasmTypes.Querier = Querier{}

// Routes of pre-determined wasm querier
const (
	WasmQueryRouteBank    = "bank"
	WasmQueryRouteStaking = "staking"
	WasmQueryRouteWasm    = "wasm"
)

// Query - interface for wasmTypes.Querier
func (q Querier) Query(request wasmTypes.QueryRequest) ([]byte, error) {
	switch {
	case request.Bank != nil:
		if querier, ok := q.Queriers[WasmQueryRouteBank]; ok {
			return querier.Query(q.Ctx, request)
		} else {
			return nil, ErrNoRegisteredQuerier(WasmQueryRouteBank)
		}
	case request.Custom != nil:
		var customQuery WasmCustomQuery
		err := json.Unmarshal(request.Custom, &customQuery)
		if err != nil {
			return nil, sdk.ErrInternal(err.Error())
		}

		if querier, ok := q.Queriers[customQuery.Route]; ok {
			return querier.QueryCustom(q.Ctx, customQuery.QueryData)
		} else {
			return nil, ErrNoRegisteredQuerier(customQuery.Route)
		}
	case request.Staking != nil:
		if querier, ok := q.Queriers[WasmQueryRouteStaking]; ok {
			return querier.Query(q.Ctx, request)
		} else {
			return nil, ErrNoRegisteredQuerier(WasmQueryRouteStaking)
		}
	case request.Wasm != nil:
		if querier, ok := q.Queriers[WasmQueryRouteWasm]; ok {
			return querier.Query(q.Ctx, request)
		} else {
			return nil, ErrNoRegisteredQuerier(WasmQueryRouteWasm)
		}
	}

	return nil, wasmTypes.Unknown{}
}
