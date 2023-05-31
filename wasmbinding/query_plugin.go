package wasmbinding

import (
	"encoding/json"

	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/classic-terra/core/wasmbinding/bindings"
	marketkeeper "github.com/classic-terra/core/x/market/keeper"
	markettypes "github.com/classic-terra/core/x/market/types"
)

// TaxCapQueryResponse - tax cap query response for wasm module
type TaxCapQueryResponse struct {
	// uint64 string, eg "1000000"
	Cap string `json:"cap"`
}

// CustomQuerier dispatches custom CosmWasm bindings queries.
func CustomQuerier(qp *QueryPlugin) func(ctx sdk.Context, request json.RawMessage) ([]byte, error) {
	return func(ctx sdk.Context, request json.RawMessage) ([]byte, error) {
		var contractQuery bindings.TerraQuery
		if err := json.Unmarshal(request, &contractQuery); err != nil {
			return nil, sdkerrors.Wrap(err, "terra query")
		}

		switch {
		case contractQuery.Swap != nil:
			q := marketkeeper.NewQuerier(*qp.marketKeeper)
			res, err := q.Swap(sdk.WrapSDKContext(ctx), &markettypes.QuerySwapRequest{
				OfferCoin: contractQuery.Swap.OfferCoin.String(),
				AskDenom:  contractQuery.Swap.AskDenom,
			})
			if err != nil {
				return nil, err
			}

			bz, err := json.Marshal(bindings.SwapQueryResponse{Receive: ConvertSdkCoinToWasmCoin(res.ReturnCoin)})
			if err != nil {
				return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
			}

			return bz, nil

		case contractQuery.ExchangeRates != nil:
			// LUNA / BASE_DENOM
			baseDenomExchangeRate, err := qp.oracleKeeper.GetLunaExchangeRate(ctx, contractQuery.ExchangeRates.BaseDenom)
			if err != nil {
				return nil, err
			}

			var items []bindings.ExchangeRateItem
			for _, quoteDenom := range contractQuery.ExchangeRates.QuoteDenoms {
				// LUNA / QUOTE_DENOM
				quoteDenomExchangeRate, err := qp.oracleKeeper.GetLunaExchangeRate(ctx, quoteDenom)
				if err != nil {
					continue
				}

				// (LUNA / QUOTE_DENOM) / (BASE_DENOM / LUNA) = BASE_DENOM / QUOTE_DENOM
				items = append(items, bindings.ExchangeRateItem{
					ExchangeRate: quoteDenomExchangeRate.Quo(baseDenomExchangeRate).String(),
					QuoteDenom:   quoteDenom,
				})
			}

			bz, err := json.Marshal(bindings.ExchangeRatesQueryResponse{
				BaseDenom:     contractQuery.ExchangeRates.BaseDenom,
				ExchangeRates: items,
			})
			if err != nil {
				return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
			}

			return bz, nil

		case contractQuery.TaxRate != nil:
			rate := qp.treasuryKeeper.GetTaxRate(ctx)
			bz, err := json.Marshal(bindings.TaxRateQueryResponse{Rate: rate.String()})
			if err != nil {
				return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
			}

			return bz, nil

		case contractQuery.TaxCap != nil:
			cap := qp.treasuryKeeper.GetTaxCap(ctx, contractQuery.TaxCap.Denom)
			bz, err := json.Marshal(TaxCapQueryResponse{Cap: cap.String()})
			if err != nil {
				return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
			}

			return bz, nil

		default:
			return nil, wasmvmtypes.UnsupportedRequest{Kind: "unknown terra query variant"}
		}
	}
}
