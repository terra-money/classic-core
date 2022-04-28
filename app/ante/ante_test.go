package ante_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/terra-money/core/app"
	"github.com/terra-money/core/app/ante"

	sdk "github.com/cosmos/cosmos-sdk/types"
	cosmosante "github.com/cosmos/cosmos-sdk/x/auth/ante"
	"github.com/cosmos/cosmos-sdk/x/auth/types"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func TestNewAnteHandler(t *testing.T) {
	_, err := ante.NewAnteHandler(ante.HandlerOptions{})
	require.Error(t, err)

	_, err = ante.NewAnteHandler(ante.HandlerOptions{
		HandlerOptions: cosmosante.HandlerOptions{
			AccountKeeper: dummyAccountKeeper{},
		},
	})
	require.Error(t, err)

	_, err = ante.NewAnteHandler(ante.HandlerOptions{
		HandlerOptions: cosmosante.HandlerOptions{
			AccountKeeper: dummyAccountKeeper{},
			BankKeeper:    dummyBankKeeper{},
		},
	})
	require.Error(t, err)

	_, err = ante.NewAnteHandler(ante.HandlerOptions{
		HandlerOptions: cosmosante.HandlerOptions{
			AccountKeeper: dummyAccountKeeper{},
			BankKeeper:    dummyBankKeeper{},
		},
		OracleKeeper: dummyOracleKeeper{},
	})
	require.Error(t, err)

	encodingConfig := app.MakeEncodingConfig()
	signModeHandler := encodingConfig.TxConfig.SignModeHandler()
	_, err = ante.NewAnteHandler(ante.HandlerOptions{
		HandlerOptions: cosmosante.HandlerOptions{
			AccountKeeper:   dummyAccountKeeper{},
			BankKeeper:      dummyBankKeeper{},
			SignModeHandler: signModeHandler,
		},
		OracleKeeper: dummyOracleKeeper{},
	})
	require.Error(t, err)

	_, err = ante.NewAnteHandler(ante.HandlerOptions{
		HandlerOptions: cosmosante.HandlerOptions{
			AccountKeeper:   dummyAccountKeeper{},
			BankKeeper:      dummyBankKeeper{},
			SignModeHandler: signModeHandler,
		},
		OracleKeeper:      dummyOracleKeeper{},
		TXCounterStoreKey: sdk.NewKVStoreKey("wasm"),
	})
	require.NoError(t, err)

}

type dummyAccountKeeper struct{}

// GetParams nolint
func (_ dummyAccountKeeper) GetParams(_ sdk.Context) (params types.Params) {
	return types.DefaultParams()
}

// GetAccount nolint
func (_ dummyAccountKeeper) GetAccount(_ sdk.Context, addr sdk.AccAddress) types.AccountI {
	return types.NewBaseAccountWithAddress(addr)
}

// SetAccount nolint
func (_ dummyAccountKeeper) SetAccount(_ sdk.Context, _ types.AccountI) {}

// GetModuleAddress nolint
func (_ dummyAccountKeeper) GetModuleAddress(moduleName string) sdk.AccAddress {
	return types.NewModuleAddress(moduleName)
}

type dummyBankKeeper struct{}

// SendCoinsFromAccountToModule nolint
func (_ dummyBankKeeper) SendCoinsFromAccountToModule(_ sdk.Context, _ sdk.AccAddress, _ string, _ sdk.Coins) error {
	return nil
}

type dummyOracleKeeper struct {
	feeders map[string]string
}

func (ok dummyOracleKeeper) ValidateFeeder(ctx sdk.Context, feederAddr sdk.AccAddress, validatorAddr sdk.ValAddress) error {
	if val, ok := ok.feeders[validatorAddr.String()]; ok && val == feederAddr.String() {
		return nil
	}

	return sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "cannot ensure feeder right")
}
