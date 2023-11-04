package dyncomm

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/classic-terra/core/v2/x/dyncomm/keeper"
	"github.com/classic-terra/core/v2/x/dyncomm/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// InitGenesis initializes default parameters
// and the keeper's address to pubkey map
func InitGenesis(ctx sdk.Context, keeper keeper.Keeper, data *types.GenesisState) {
	keeper.SetParams(ctx, data.Params)

	// iterate validators and set target rates
	keeper.StakingKeeper.IterateValidators(ctx, func(index int64, validator stakingtypes.ValidatorI) (stop bool) {
		val := validator.(stakingtypes.Validator)
		keeper.SetTargetCommissionRate(ctx, val.OperatorAddress, val.Commission.Rate)
		return false
	})

	err := keeper.UpdateAllBondedValidatorRates(ctx)
	if err != nil {
		panic("could not initialize genesis")
	}
}

// ExportGenesis writes the current store values
// to a genesis file, which can be imported again
// with InitGenesis
func ExportGenesis(ctx sdk.Context, keeper keeper.Keeper) (data *types.GenesisState) {
	params := keeper.GetParams(ctx)
	var rates []types.ValidatorCommissionRate

	// rates = append(rates)
	keeper.IterateDynCommissionRates(ctx, func(rate types.ValidatorCommissionRate) (stop bool) {
		rates = append(rates, rate)
		return false
	})

	return types.NewGenesisState(params, rates)
}
