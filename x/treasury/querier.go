package treasury

import (
	"strconv"
	"strings"

	"github.com/terra-project/core/types/util"

	"github.com/cosmos/cosmos-sdk/codec"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// query endpoints supported by the treasury Querier
const (
	QueryTaxRate             = "tax-rate"
	QueryTaxCap              = "tax-cap"
	QueryMiningRewardWeight  = "reward-weight"
	QuerySeigniorageProceeds = "seigniorage-proceeds"
	QueryCurrentEpoch        = "current-epoch"
	QueryParams              = "params"
	QueryIssuance            = "issuance"
	QueryTaxProceeds         = "tax-proceeds"
)

// NewQuerier is the module level router for state queries
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case QueryTaxRate:
			return queryTaxRate(ctx, path[1:], req, keeper)
		case QueryTaxCap:
			return queryTaxCap(ctx, path[1:], req, keeper)
		case QueryMiningRewardWeight:
			return queryMiningRewardWeight(ctx, path[1:], req, keeper)
		case QueryTaxProceeds:
			return queryTaxProceeds(ctx, path[1:], req, keeper)
		case QuerySeigniorageProceeds:
			return querySeigniorageProceeds(ctx, path[1:], req, keeper)
		case QueryIssuance:
			return queryIssuance(ctx, path[1:], req, keeper)
		case QueryCurrentEpoch:
			return queryCurrentEpoch(ctx, req, keeper)
		case QueryParams:
			return queryParams(ctx, req, keeper)
		default:
			return nil, sdk.ErrUnknownRequest("unknown treasury query endpoint")
		}
	}
}

// JSON response format
type QueryTaxRateResponse struct {
	TaxRate sdk.Dec `json:"tax_rate"`
}

func (r QueryTaxRateResponse) String() (out string) {
	out = r.TaxRate.String()
	return strings.TrimSpace(out)
}

// nolint: unparam
func queryTaxRate(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	epoch, ok := sdk.NewIntFromString(path[0])
	if !ok {
		return nil, sdk.ErrInternal("epoch parameter is not correctly formatted")
	}

	taxRate := keeper.GetTaxRate(ctx, epoch)
	bz, err := codec.MarshalJSONIndent(keeper.cdc, QueryTaxRateResponse{TaxRate: taxRate})
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}
	return bz, nil
}

// JSON response format
type QueryTaxCapResponse struct {
	TaxCap sdk.Int `json:"tax_cap"`
}

func (r QueryTaxCapResponse) String() (out string) {
	out = r.TaxCap.String()
	return strings.TrimSpace(out)
}

// nolint: unparam
func queryTaxCap(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	denom := path[0]
	taxCap := keeper.GetTaxCap(ctx, denom)

	bz, err := codec.MarshalJSONIndent(keeper.cdc, QueryTaxCapResponse{TaxCap: taxCap})
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}
	return bz, nil
}

// JSON response format
type QueryIssuanceResponse struct {
	Issuance sdk.Int `json:"issuance"`
}

func (r QueryIssuanceResponse) String() (out string) {
	out = r.Issuance.String()
	return strings.TrimSpace(out)
}

// nolint: unparam
func queryIssuance(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	denom := path[0]
	var dayStr string
	if len(path) == 2 {
		dayStr = path[1]
	}

	var day int64
	if len(dayStr) == 0 {
		day = ctx.BlockHeight() / util.BlocksPerDay
	} else {
		day, _ = strconv.ParseInt(dayStr, 10, 64)
	}

	issuance := keeper.mtk.GetIssuance(ctx, denom, sdk.NewInt(day))
	bz, err := codec.MarshalJSONIndent(keeper.cdc, QueryIssuanceResponse{Issuance: issuance})
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}
	return bz, nil
}

// JSON response format
type QueryMiningRewardWeightResponse struct {
	RewardWeight sdk.Dec `json:"reward_weight"`
}

func (r QueryMiningRewardWeightResponse) String() (out string) {
	out = r.RewardWeight.String()
	return strings.TrimSpace(out)
}

// nolint: unparam
func queryMiningRewardWeight(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	epoch, ok := sdk.NewIntFromString(path[0])
	if !ok {
		return nil, sdk.ErrInternal("epoch parameter is not correctly formatted")
	}

	rewardWeight := keeper.GetRewardWeight(ctx, epoch)
	bz, err := codec.MarshalJSONIndent(keeper.cdc, QueryMiningRewardWeightResponse{RewardWeight: rewardWeight})
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}
	return bz, nil
}

// JSON response format
type QueryTaxProceedsResponse struct {
	TaxProceeds sdk.Coins `json:"tax_proceeds"`
}

func (r QueryTaxProceedsResponse) String() (out string) {
	out = r.TaxProceeds.String()
	return strings.TrimSpace(out)
}

// nolint: unparam
func queryTaxProceeds(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	epoch, ok := sdk.NewIntFromString(path[0])
	if !ok {
		return nil, sdk.ErrInternal("epoch parameter is not correctly formatted")
	}

	pool := keeper.PeekTaxProceeds(ctx, epoch)
	bz, err := codec.MarshalJSONIndent(keeper.cdc, QueryTaxProceedsResponse{TaxProceeds: pool})
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}
	return bz, nil
}

// JSON response format
type QuerySeigniorageProceedsResponse struct {
	SeigniorageProceeds sdk.Int `json:"seigniorage_proceeds"`
}

func (r QuerySeigniorageProceedsResponse) String() (out string) {
	out = r.SeigniorageProceeds.String()
	return strings.TrimSpace(out)
}

// nolint: unparam
func querySeigniorageProceeds(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	epoch, ok := sdk.NewIntFromString(path[0])
	if !ok {
		return nil, sdk.ErrInternal("epoch parameter is not correctly formatted")
	}

	pool := keeper.mtk.PeekEpochSeigniorage(ctx, epoch)
	bz, err := codec.MarshalJSONIndent(keeper.cdc, QuerySeigniorageProceedsResponse{SeigniorageProceeds: pool})
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}
	return bz, nil
}

// JSON response format
type QueryCurrentEpochResponse struct {
	CurrentEpoch sdk.Int `json:"current_epoch"`
}

func (r QueryCurrentEpochResponse) String() (out string) {
	out = r.CurrentEpoch.String()
	return strings.TrimSpace(out)
}

func queryCurrentEpoch(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	curEpoch := util.GetEpoch(ctx)
	bz, err := codec.MarshalJSONIndent(keeper.cdc, QueryCurrentEpochResponse{CurrentEpoch: curEpoch})
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}
	return bz, nil
}

func queryParams(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	bz, err := codec.MarshalJSONIndent(keeper.cdc, keeper.GetParams(ctx))
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}
	return bz, nil
}
