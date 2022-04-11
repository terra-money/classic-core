package treasury

import (
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	wasmvmtypes "github.com/CosmWasm/wasmvm/types"

	wasm "github.com/terra-money/core/x/wasm/exported"
)

var _ wasm.WasmQuerierInterface = WasmQuerier{}

// WasmQuerier - staking query interface for wasm contract
type WasmQuerier struct{}

// NewWasmQuerier return bank wasm query interface
func NewWasmQuerier() WasmQuerier {
	return WasmQuerier{}
}

// Query - implement query function
func (WasmQuerier) Query(_ sdk.Context, _ wasmvmtypes.QueryRequest) ([]byte, error) {
	return nil, nil
}

// CosmosQuery contains various treasury queries
type CosmosQuery struct {
	TaxRate *struct{}          `json:"tax_rate,omitempty"`
	TaxCap  *QueryTaxCapParams `json:"tax_cap,omitempty"`
}

// QueryTaxCapParams for query
// - 'custom/treasury/taxRate
type QueryTaxCapParams struct {
	Denom string `json:"denom"`
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
		rate := sdk.ZeroDec()
		bz, err = json.Marshal(TaxRateQueryResponse{Rate: rate.String()})
	} else if query.TaxCap != nil {
		cap := sdk.ZeroInt()
		bz, err = json.Marshal(TaxCapQueryResponse{Cap: cap.String()})
	} else {
		return nil, sdkerrors.ErrInvalidRequest
	}

	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}
