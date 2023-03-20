package wasm

import (
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	wasmvmtypes "github.com/CosmWasm/wasmvm/types"

	"github.com/classic-terra/core/x/treasury/keeper"
	"github.com/classic-terra/core/x/treasury/types"
	wasm "github.com/classic-terra/core/x/wasm/exported"
)

var _ wasm.WasmQuerierInterface = Querier{}

// Querier - staking query interface for wasm contract
type Querier struct {
	keeper keeper.Keeper
}

// NewQuerier return bank wasm query interface
func NewWasmQuerier(keeper keeper.Keeper) Querier {
	return Querier{keeper}
}

// Query - implement query function
func (Querier) Query(_ sdk.Context, _ wasmvmtypes.QueryRequest) ([]byte, error) {
	return nil, nil
}

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
func (querier Querier) QueryCustom(ctx sdk.Context, data json.RawMessage) ([]byte, error) {
	var query CosmosQuery
	err := json.Unmarshal(data, &query)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	var bz []byte

	switch {
	case query.TaxRate != nil:
		rate := querier.keeper.GetTaxRate(ctx)
		bz, err = json.Marshal(TaxRateQueryResponse{Rate: rate.String()})
	case query.TaxCap != nil:
		cap := querier.keeper.GetTaxCap(ctx, query.TaxCap.Denom)
		bz, err = json.Marshal(TaxCapQueryResponse{Cap: cap.String()})
	default:
		return nil, sdkerrors.ErrInvalidRequest
	}

	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}
