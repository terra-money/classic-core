package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto/ed25519"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	authvesttypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"

	core "github.com/terra-money/core/types"
	"github.com/terra-money/core/x/vesting/types"
)

var (
	pk1   = ed25519.GenPrivKey().PubKey()
	pk2   = ed25519.GenPrivKey().PubKey()
	addr1 = sdk.ValAddress(pk1.Address())
	addr2 = sdk.ValAddress(pk2.Address())
)

// require invalid vesting account fails validation
func TestValidateGenesisInvalidAccounts(t *testing.T) {
	acc1 := authtypes.NewBaseAccountWithAddress(sdk.AccAddress(addr1))
	coins := sdk.NewCoins(sdk.NewInt64Coin(core.MicroLunaDenom, 150))
	baseVestingAcc := authvesttypes.NewBaseVestingAccount(acc1, coins, 0)

	// invalid delegated vesting
	baseVestingAcc.DelegatedVesting = coins.Add(coins...)

	acc2 := authtypes.NewBaseAccountWithAddress(sdk.AccAddress(addr2))

	genAccs := make([]authtypes.GenesisAccount, 2)
	genAccs[0] = baseVestingAcc
	genAccs[1] = acc2

	require.Error(t, authtypes.ValidateGenAccounts(genAccs))
	baseVestingAcc.DelegatedVesting = coins
	genAccs[0] = baseVestingAcc
	require.NoError(t, authtypes.ValidateGenAccounts(genAccs))

	// invalid vesting time
	genAccs[0] = types.NewLazyGradedVestingAccountRaw(baseVestingAcc, types.VestingSchedules{types.VestingSchedule{core.MicroLunaDenom, types.Schedules{types.Schedule{1654668078, 1554668078, sdk.OneDec()}}}})
	require.Error(t, authtypes.ValidateGenAccounts(genAccs))
}
