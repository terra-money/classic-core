package nameservice

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/terra-project/core/x/nameservice/internal/keeper"
	"testing"
)

func setup(t *testing.T) (keeper.TestInput, sdk.Handler) {
	input := keeper.CreateTestInput(t)
	params := DefaultParams()
	params.MinNameLength = 2
	input.NameserviceKeeper.SetParams(input.Ctx, params)
	h := NewHandler(input.NameserviceKeeper)
	return input, h
}
