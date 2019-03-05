package treasury

// import (
// 	"github.com/cosmos/cosmos-sdk/codec"
// 	sdk "github.com/cosmos/cosmos-sdk/types"
// 	"github.com/cosmos/cosmos-sdk/x/staking/types"
// 	abci "github.com/tendermint/tendermint/abci/types"
// )

// // query endpoints supported by the treasury Querier
// const (
// 	QueryMiningRewardWeight = "mingReward"
// 	QueryTaxRate            = "taxRate"
// 	QueryIncomePool         = "incomePool"
// 	QueryOutstandingClaims  = "outstandingClaims"
// 	QueryParameters         = "parameters"
// )

// // creates a querier for staking REST endpoints
// func NewQuerier(k Keeper, cdc *codec.Codec) sdk.Querier {
// 	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
// 		switch path[0] {
// 		case QueryMiningRewardWeight:
// 			return queryMiningRewardWeight(ctx, cdc, k)
// 		case QueryTaxRate:
// 			return queryValidator(ctx, cdc, req, k)
// 		case QueryIncomePool:
// 			return queryValidatorDelegations(ctx, cdc, req, k)
// 		case QueryOutstandingClaims:
// 			return queryValidatorUnbondingDelegations(ctx, cdc, req, k)
// 		case QueryParameters:
// 			return queryParameters(ctx, cdc, k)
// 		default:
// 			return nil, sdk.ErrUnknownRequest("unknown treasury query endpoint")
// 		}
// 	}
// }

// // defines the params for the following queries:
// // - 'custom/staking/delegation'
// // - 'custom/staking/unbondingDelegation'
// // - 'custom/staking/delegatorValidator'
// type QueryBondsParams struct {
// 	DelegatorAddr sdk.AccAddress
// 	ValidatorAddr sdk.ValAddress
// }

// func NewQueryBondsParams(delegatorAddr sdk.AccAddress, validatorAddr sdk.ValAddress) QueryBondsParams {
// 	return QueryBondsParams{
// 		DelegatorAddr: delegatorAddr,
// 		ValidatorAddr: validatorAddr,
// 	}
// }

// func queryMiningRewardWeight(ctx sdk.Context, cdc *codec.Codec, req abci.RequestQuery, k keep.Keeper) (res []byte, err sdk.Error) {
// 	var params QueryValidatorParams

// 	errRes := cdc.UnmarshalJSON(req.Data, &params)
// 	if errRes != nil {
// 		return []byte{}, sdk.ErrUnknownAddress("")
// 	}

// 	validator, found := k.GetValidator(ctx, params.ValidatorAddr)
// 	if !found {
// 		return []byte{}, types.ErrNoValidatorFound(types.DefaultCodespace)
// 	}

// 	res, errRes = codec.MarshalJSONIndent(cdc, validator)
// 	if errRes != nil {
// 		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", errRes.Error()))
// 	}
// 	return res, nil
// }
