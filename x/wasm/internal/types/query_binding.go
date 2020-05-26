package types

import (
	"encoding/json"

	wasmTypes "github.com/CosmWasm/go-cosmwasm/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// WasmQuerierInterface - query registration interface for other modules
type WasmQuerierInterface interface {
	Query(ctx sdk.Context, request wasmTypes.QueryRequest) ([]byte, error)
	QueryCustom(ctx sdk.Context, data json.RawMessage) ([]byte, error)
}

// Querier - wasm query handler
type Querier struct {
	Ctx      sdk.Context
	Queriers map[string]WasmQuerierInterface
}

// NewQuerier return wasm querier
func NewQuerier() Querier {
	return Querier{
		Queriers: make(map[string]WasmQuerierInterface),
	}
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
	WasmQueryRouteMarket  = "market"
	WasmQueryRouteWasm    = "wasm"
)

// WithCtx returns new querier with context
func (q Querier) WithCtx(ctx sdk.Context) Querier {
	q.Ctx = ctx
	return q
}

// Query - interface for wasmTypes.Querier
func (q Querier) Query(request wasmTypes.QueryRequest) ([]byte, error) {
	switch {
	case request.Bank != nil:
		if querier, ok := q.Queriers[WasmQueryRouteBank]; ok {
			return querier.Query(q.Ctx, request)
		} else {
			return nil, sdkerrors.Wrap(ErrNoRegisteredQuerier, WasmQueryRouteBank)
		}
	case request.Custom != nil:
		var customQuery WasmCustomQuery
		err := json.Unmarshal(request.Custom, &customQuery)
		if err != nil {
			return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
		}

		if querier, ok := q.Queriers[customQuery.Route]; ok {
			return querier.QueryCustom(q.Ctx, customQuery.QueryData)
		} else {
			return nil, sdkerrors.Wrap(ErrNoRegisteredQuerier, customQuery.Route)
		}
	case request.Staking != nil:
		if querier, ok := q.Queriers[WasmQueryRouteStaking]; ok {
			return querier.Query(q.Ctx, request)
		} else {
			return nil, sdkerrors.Wrap(ErrNoRegisteredQuerier, WasmQueryRouteStaking)
		}
	case request.Wasm != nil:
		if querier, ok := q.Queriers[WasmQueryRouteWasm]; ok {
			return querier.Query(q.Ctx, request)
		} else {
			return nil, sdkerrors.Wrap(ErrNoRegisteredQuerier, WasmQueryRouteWasm)
		}
	}

	return nil, wasmTypes.Unknown{}
}
