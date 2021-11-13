package app

import (
	"fmt"
	"io/ioutil"
	"sort"
	"time"

	core "github.com/terra-money/core/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	vestexported "github.com/cosmos/cosmos-sdk/x/auth/vesting/exported"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

func (app *TerraApp) trackingAll(ctx sdk.Context) {
	// Build validator token share map to calculate delegators staking tokens
	validators := stakingtypes.Validators(app.StakingKeeper.GetAllValidators(ctx))
	tokenShareRates := make(map[string]sdk.Dec)
	for _, validator := range validators {
		if validator.IsBonded() {
			tokenShareRates[validator.GetOperator().String()] = validator.GetBondedTokens().ToDec().Quo(validator.GetDelegatorShares())
		}
	}

	// Load oracle whitelist
	var denoms []string
	for _, denom := range app.OracleKeeper.Whitelist(ctx) {
		denoms = append(denoms, denom.Name)
	}

	denoms = append(denoms, app.StakingKeeper.BondDenom(ctx))

	// Minimum coins to be included in tracking
	minCoins := sdk.Coins{}
	accsPerDenom := map[string]ExportAccounts{}
	for _, denom := range denoms {
		minCoins = append(minCoins, sdk.NewCoin(denom, sdk.OneInt().MulRaw(core.MicroUnit)))
		accsPerDenom[denom] = ExportAccounts{}
	}

	minCoins = minCoins.Sort()
	vestingCoins := sdk.NewCoins()

	app.Logger().Info("Start Tracking Load Account")
	app.AccountKeeper.IterateAccounts(ctx, func(acc authtypes.AccountI) bool {

		// Skip module accounts from tracking
		if _, ok := acc.(authtypes.ModuleAccountI); ok {
			return false
		}

		// Record vesting accounts
		if vacc, ok := acc.(vestexported.VestingAccount); ok {
			vestingCoins = vestingCoins.Add(vacc.GetVestingCoins(ctx.BlockHeader().Time)...)
		}

		// Compute account balance
		balances := app.BankKeeper.GetAllBalances(ctx, acc.GetAddress())

		// Compute staking amount
		stakingAmt := sdk.ZeroInt()
		delegations := app.StakingKeeper.GetAllDelegatorDelegations(ctx, acc.GetAddress())
		undelegations := app.StakingKeeper.GetUnbondingDelegations(ctx, acc.GetAddress(), 100)
		for _, delegation := range delegations {
			valAddr := delegation.GetValidatorAddr().String()
			if tokenShareRate, ok := tokenShareRates[valAddr]; ok {
				delegationAmt := delegation.GetShares().Mul(tokenShareRate).TruncateInt()
				stakingAmt = stakingAmt.Add(delegationAmt)
			}
		}

		unbondingAmt := sdk.ZeroInt()
		for _, undelegation := range undelegations {
			undelegationAmt := sdk.ZeroInt()
			for _, entry := range undelegation.Entries {
				undelegationAmt = undelegationAmt.Add(entry.Balance)
			}

			unbondingAmt.Add(undelegationAmt)
		}

		// Add staking amount to account balance
		stakingCoins := sdk.NewCoins(sdk.NewCoin(app.StakingKeeper.BondDenom(ctx), stakingAmt.Add(unbondingAmt)))
		balances = balances.Add(stakingCoins...)

		// Check minimum coins
		for _, denom := range denoms {
			if amt := balances.AmountOf(denom); amt.GTE(sdk.NewInt(core.MicroUnit)) {
				accsPerDenom[denom] = append(accsPerDenom[denom], NewExportAccount(acc.GetAddress(), amt))
			}
		}

		return false
	})

	app.Logger().Info("End Tracking Load Account")

	go app.exportVestingSupply(ctx, vestingCoins)
	for _, denom := range denoms {
		go app.exportRanking(ctx, accsPerDenom[denom], denom)
	}

}

func (app *TerraApp) exportVestingSupply(ctx sdk.Context, vestingCoins sdk.Coins) {
	app.Logger().Info("Start Tracking Vesting Luna Supply")
	bz, err := codec.MarshalJSONIndent(app.legacyAmino, vestingCoins)
	if err != nil {
		app.Logger().Error(err.Error())
	}

	// nolint
	err = ioutil.WriteFile(fmt.Sprintf("/tmp/vesting-%s.json", time.Now().Format(time.RFC3339)), bz, 0644)
	if err != nil {
		app.Logger().Error(err.Error())
	}
	app.Logger().Info("End Tracking Vesting Luna Supply")
}

func (app *TerraApp) exportRanking(ctx sdk.Context, accs ExportAccounts, denom string) {
	app.Logger().Info(fmt.Sprintf("Start Wallet Balance Tracking for %s", denom))

	// sort descending order
	sort.Sort(accs)

	// nolint
	err := ioutil.WriteFile(fmt.Sprintf("/tmp/tracking-%s-%s.txt", denom, time.Now().Format(time.RFC3339)), []byte(accs.String()), 0644)
	if err != nil {
		app.Logger().Error(err.Error())
	}

	app.Logger().Info(fmt.Sprintf("End Wallet Balance Tracking for %s", denom))
}

// ExportAccount is ranking export account format
type ExportAccount struct {
	Address sdk.AccAddress `json:"address"`
	Amount  sdk.Int        `json:"amount"`
}

// NewExportAccount returns new ExportAccount instance
func NewExportAccount(address sdk.AccAddress, amount sdk.Int) ExportAccount {
	return ExportAccount{
		Address: address,
		Amount:  amount,
	}
}

// String - implement stringify interface
func (acc ExportAccount) String() (out string) {
	return fmt.Sprintf("%s,%s", acc.Address, acc.Amount)
}

// ExportAccounts simple wrapper to print ranking list
type ExportAccounts []ExportAccount

// Less - implement Sort interface
func (accs ExportAccounts) Len() int {
	return len(accs)
}

// Less - implement Sort interface descanding order
func (accs ExportAccounts) Less(i, j int) bool {
	return accs[i].Amount.GT(accs[j].Amount)
}

// Less - implement Sort interface
func (accs ExportAccounts) Swap(i, j int) { accs[i], accs[j] = accs[j], accs[i] }

// String - implement stringify interface
func (accs ExportAccounts) String() (out string) {
	out = ""
	for _, a := range accs {
		out += a.String() + "\n"
	}

	return
}
