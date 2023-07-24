package wasmbinding

import (
	"encoding/json"
	"fmt"

	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/classic-terra/core/v2/wasmbinding/bindings"
	marketkeeper "github.com/classic-terra/core/v2/x/market/keeper"
	markettypes "github.com/classic-terra/core/v2/x/market/types"
)

// TaxCapQueryResponse - tax cap query response for wasm module
type TaxCapQueryResponse struct {
	// uint64 string, eg "1000000"
	Cap string `json:"cap"`
}

// StargateQuerier dispatches whitelisted stargate queries
func StargateQuerier(queryRouter baseapp.GRPCQueryRouter, cdc codec.Codec) func(ctx sdk.Context, request *wasmvmtypes.StargateQuery) ([]byte, error) {
	return func(ctx sdk.Context, request *wasmvmtypes.StargateQuery) ([]byte, error) {
		protoResponseType, err := GetWhitelistedQuery(request.Path)
		if err != nil {
			return nil, err
		}

		route := queryRouter.Route(request.Path)
		if route == nil {
			return nil, wasmvmtypes.UnsupportedRequest{Kind: fmt.Sprintf("No route to query '%s'", request.Path)}
		}

		res, err := route(ctx, abci.RequestQuery{
			Data: request.Data,
			Path: request.Path,
		})
		if err != nil {
			return nil, err
		}

		bz, err := ConvertProtoToJSONMarshal(protoResponseType, res.Value, cdc)
		if err != nil {
			return nil, err
		}

		return bz, nil
	}
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
			taxRate := qp.treasuryKeeper.GetTaxRate(ctx)
			bz, err := json.Marshal(bindings.TaxRateQueryResponse{Rate: taxRate.String()})
			if err != nil {
				return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
			}

			return bz, nil

		case contractQuery.TaxCap != nil:
			taxCap := qp.treasuryKeeper.GetTaxCap(ctx, contractQuery.TaxCap.Denom)
			bz, err := json.Marshal(TaxCapQueryResponse{Cap: taxCap.String()})
			if err != nil {
				return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
			}

			return bz, nil

		default:
			return nil, wasmvmtypes.UnsupportedRequest{Kind: "unknown terra query variant"}
		}
	}
}
