package types

import (
	"encoding/json"

	wasmvmtypes "github.com/CosmWasm/wasmvm/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// WasmQuerierInterface - query registration interface for other modules
type WasmQuerierInterface interface {
	Query(ctx sdk.Context, request wasmvmtypes.QueryRequest) ([]byte, error)
	QueryCustom(ctx sdk.Context, data json.RawMessage) ([]byte, error)
}

// IBCWasmQuerierInterface - query registration interface for ibc querier
type IBCWasmQuerierInterface interface {
	Query(ctx sdk.Context, contractAddr sdk.AccAddress, request wasmvmtypes.QueryRequest) ([]byte, error)
}

// StargateWasmQuerierInterface - query registration interface for stargate querier
type StargateWasmQuerierInterface interface {
	Query(ctx sdk.Context, request wasmvmtypes.QueryRequest) ([]byte, error)
}

// Querier - wasm query handler
type Querier struct {
	Ctx             sdk.Context
	ContractAddr    sdk.AccAddress
	Queriers        map[string]WasmQuerierInterface
	StargateQuerier StargateWasmQuerierInterface
	IBCQuerier      IBCWasmQuerierInterface
}

// NewWasmQuerier return wasm querier
func NewWasmQuerier() Querier {
	return Querier{
		Queriers: make(map[string]WasmQuerierInterface),
	}
}

// WasmCustomQuery - wasm custom query
type WasmCustomQuery struct {
	Route     string          `json:"route"`
	QueryData json.RawMessage `json:"query_data"`
}

var _ wasmvmtypes.Querier = Querier{}

// Routes of pre-determined wasm querier
const (
	WasmQueryRouteBank     = "bank"
	WasmQueryRouteStaking  = "staking"
	WasmQueryRouteMarket   = "market"
	WasmQueryRouteOracle   = "oracle"
	WasmQueryRouteTreasury = "treasury"
	WasmQueryRouteWasm     = "wasm"
)

// WithCtx returns new querier with context
func (q Querier) WithCtx(ctx sdk.Context) Querier {
	q.Ctx = ctx
	return q
}

// WithContractAddr returns new querier with contractAddr
func (q Querier) WithContractAddr(contractAddr sdk.AccAddress) Querier {
	q.ContractAddr = contractAddr
	return q
}

// GasConsumed consume gas in the current context
func (q Querier) GasConsumed() uint64 {
	return q.Ctx.GasMeter().GasConsumed()
}

// Query - interface for wasmvmtypes.Querier
func (q Querier) Query(request wasmvmtypes.QueryRequest, gasLimit uint64) ([]byte, error) {
	// set a limit for a ctx
	// gasLimit passed from the go-cosmwasm part, so need to divide it with gas multiplier
	ctx := q.Ctx.WithGasMeter(sdk.NewGasMeter(gasLimit / GasMultiplier))

	// make sure we charge the higher level context even on panic
	defer func() {
		q.Ctx.GasMeter().ConsumeGas(ctx.GasMeter().GasConsumed(), "contract sub-query")
	}()

	// do the query
	switch {
	case request.Bank != nil:
		if querier, ok := q.Queriers[WasmQueryRouteBank]; ok {
			return querier.Query(ctx, request)
		}

		return nil, sdkerrors.Wrap(ErrNoRegisteredQuerier, WasmQueryRouteBank)
	case request.Custom != nil:
		var customQuery WasmCustomQuery
		err := json.Unmarshal(request.Custom, &customQuery)
		if err != nil {
			return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
		}

		if querier, ok := q.Queriers[customQuery.Route]; ok {
			return querier.QueryCustom(ctx, customQuery.QueryData)
		}

		return nil, sdkerrors.Wrap(ErrNoRegisteredQuerier, customQuery.Route)
	case request.Staking != nil:
		if querier, ok := q.Queriers[WasmQueryRouteStaking]; ok {
			return querier.Query(ctx, request)
		}

		return nil, sdkerrors.Wrap(ErrNoRegisteredQuerier, WasmQueryRouteStaking)

	case request.Wasm != nil:
		if querier, ok := q.Queriers[WasmQueryRouteWasm]; ok {
			return querier.Query(ctx, request)
		}

		return nil, sdkerrors.Wrap(ErrNoRegisteredQuerier, WasmQueryRouteWasm)
	case request.Stargate != nil:
		if q.StargateQuerier != nil {
			return q.StargateQuerier.Query(ctx, request)
		}

		return nil, sdkerrors.Wrap(ErrNoRegisteredQuerier, "stargate")
	case request.IBC != nil:
		if q.IBCQuerier != nil {
			return q.IBCQuerier.Query(ctx, q.ContractAddr, request)
		}

		return nil, sdkerrors.Wrap(ErrNoRegisteredQuerier, "IBC")
	}

	return nil, wasmvmtypes.Unknown{}
}
