package wasm

import (
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	wasmTypes "github.com/CosmWasm/go-cosmwasm/types"

	"github.com/terra-project/core/x/treasury/internal/keeper"
	"github.com/terra-project/core/x/treasury/internal/types"
	wasm "github.com/terra-project/core/x/wasm/exported"
)

var _ wasm.WasmQuerierInterface = WasmQuerier{}

// WasmQuerier - staking query interface for wasm contract
type WasmQuerier struct {
	keeper keeper.Keeper
}

// NewWasmQuerier return bank wasm query interface
func NewWasmQuerier(keeper keeper.Keeper) WasmQuerier {
	return WasmQuerier{keeper}
}

// Query - implement query function
func (WasmQuerier) Query(_ sdk.Context, _ wasmTypes.QueryRequest) ([]byte, error) { return nil, nil }

// CosmosQuery contains various treasury queries
type CosmosQuery struct {
	TaxRate *struct{}                `json:"tax_rate,omitempty"`
	TaxCap  *types.QueryTaxCapParams `json:"tax_cap,omitempty"`
}

// TaxRateQueryResponse - tax rate query response for wasm module
type TaxRateQueryResponse struct {
	// decimal string, eg "0.02"
	Rate string `json:"rate"`
}

// TaxCapQueryResponse - tax cap query response for wasm module
type TaxCapQueryResponse struct {
	// uint64 string, eg "1000000"
	Cap string `json:"cap"`
}

// QueryCustom implements custom query interface
func (querier WasmQuerier) QueryCustom(ctx sdk.Context, data json.RawMessage) ([]byte, error) {
	var query CosmosQuery
	err := json.Unmarshal(data, &query)

	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	var bz []byte

	if query.TaxRate != nil {
		rate := querier.keeper.GetTaxRate(ctx)
		bz, err = json.Marshal(TaxRateQueryResponse{Rate: rate.String()})
	} else if query.TaxCap != nil {
		cap := querier.keeper.GetTaxCap(ctx, query.TaxCap.Denom)
		bz, err = json.Marshal(TaxCapQueryResponse{Cap: cap.String()})
	} else {
		return nil, sdkerrors.ErrInvalidRequest
	}

	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}
