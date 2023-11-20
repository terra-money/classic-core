package dyncomm

import (
	"time"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/classic-terra/core/v2/x/dyncomm/keeper"
	"github.com/classic-terra/core/v2/x/dyncomm/types"

	core "github.com/classic-terra/core/v2/types"
)

// EndBlocker is called at the end of every block
func EndBlocker(ctx sdk.Context, k keeper.Keeper) {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyEndBlocker)

	if !core.IsPeriodLastBlock(ctx, core.BlocksPerWeek) {
		return
	}

	ctx.Logger().Info("End Epoch - Calculation of Dyncomm is due")
	k.UpdateAllBondedValidatorRates(ctx)
}
