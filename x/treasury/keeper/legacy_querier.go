package keeper

import (
	"math"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	core "github.com/terra-money/core/types"
	"github.com/terra-money/core/x/treasury/types"
)

// NewLegacyQuerier is the module level router for state queries
func NewLegacyQuerier(k Keeper, legacyQuerierCdc *codec.LegacyAmino) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err error) {
		switch path[0] {
		case types.QueryTaxRate:
			return queryTaxRate(ctx, k, legacyQuerierCdc)
		case types.QueryTaxCap:
			return queryTaxCap(ctx, req, k, legacyQuerierCdc)
		case types.QueryTaxCaps:
			return queryTaxCaps(ctx, k, legacyQuerierCdc)
		case types.QueryRewardWeight:
			return queryRewardWeight(ctx, k, legacyQuerierCdc)
		case types.QuerySeigniorageProceeds:
			return querySeigniorageProceeds(ctx, k, legacyQuerierCdc)
		case types.QueryTaxProceeds:
			return queryTaxProceeds(ctx, k, legacyQuerierCdc)
		case types.QueryParameters:
			return queryParameters(ctx, k, legacyQuerierCdc)
		case types.QueryIndicators:
			return queryIndicators(ctx, k, legacyQuerierCdc)
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unknown %s query endpoint: %s", types.ModuleName, path[0])
		}
	}
}

func queryIndicators(ctx sdk.Context, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	// Compute Total Staked Luna (TSL)
	TSL := k.stakingKeeper.TotalBondedTokens(ctx)

	// Compute Tax Rewards (TR)
	taxRewards := sdk.NewDecCoinsFromCoins(k.PeekEpochTaxProceeds(ctx)...)
	TR := k.alignCoins(ctx, taxRewards, core.MicroSDRDenom)

	epoch := k.GetEpoch(ctx)
	var res types.IndicatorQueryResponse
	if epoch == 0 {
		res = types.IndicatorQueryResponse{
			TRLYear:  TR.QuoInt(TSL),
			TRLMonth: TR.QuoInt(TSL),
		}
	} else {
		params := k.GetParams(ctx)
		previousEpochCtx := ctx.WithBlockHeight(ctx.BlockHeight() - int64(core.BlocksPerWeek))
		trlYear := k.rollingAverageIndicator(previousEpochCtx, int64(params.WindowLong-1), TRL)
		trlMonth := k.rollingAverageIndicator(previousEpochCtx, int64(params.WindowShort-1), TRL)

		computedEpochForYear := int64(math.Min(float64(params.WindowLong-1), float64(epoch)))
		computedEpochForMonty := int64(math.Min(float64(params.WindowShort-1), float64(epoch)))

		trlYear = trlYear.MulInt64(computedEpochForYear).Add(TR.QuoInt(TSL)).QuoInt64(computedEpochForYear + 1)
		trlMonth = trlMonth.MulInt64(computedEpochForMonty).Add(TR.QuoInt(TSL)).QuoInt64(computedEpochForMonty + 1)

		res = types.IndicatorQueryResponse{
			TRLYear:  trlYear,
			TRLMonth: trlMonth,
		}
	}

	bz, err := codec.MarshalJSONIndent(legacyQuerierCdc, res)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}

func queryTaxRate(ctx sdk.Context, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	taxRate := k.GetTaxRate(ctx)
	bz, err := codec.MarshalJSONIndent(legacyQuerierCdc, taxRate)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}

func queryTaxCap(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	var params types.QueryTaxCapParams
	err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	taxCap := k.GetTaxCap(ctx, params.Denom)
	bz, err := codec.MarshalJSONIndent(legacyQuerierCdc, taxCap)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}

func queryTaxCaps(ctx sdk.Context, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	var taxCaps types.TaxCapsQueryResponse
	k.IterateTaxCap(ctx, func(denom string, taxCap sdk.Int) bool {
		taxCaps = append(taxCaps, types.TaxCapsResponseItem{
			Denom:  denom,
			TaxCap: taxCap,
		})
		return false
	})

	bz, err := codec.MarshalJSONIndent(legacyQuerierCdc, taxCaps)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}

func queryRewardWeight(ctx sdk.Context, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	taxRate := k.GetRewardWeight(ctx)
	bz, err := codec.MarshalJSONIndent(legacyQuerierCdc, taxRate)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}

func querySeigniorageProceeds(ctx sdk.Context, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	seigniorage := k.PeekEpochSeigniorage(ctx)
	bz, err := codec.MarshalJSONIndent(legacyQuerierCdc, seigniorage)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return bz, nil
}

func queryTaxProceeds(ctx sdk.Context, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	proceeds := k.PeekEpochTaxProceeds(ctx)
	bz, err := codec.MarshalJSONIndent(legacyQuerierCdc, proceeds)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return bz, nil
}

func queryParameters(ctx sdk.Context, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	bz, err := codec.MarshalJSONIndent(legacyQuerierCdc, k.GetParams(ctx))
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return bz, nil
}
