package genaccounts

import (
	"errors"
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/terra-project/core/x/auth"
	"github.com/terra-project/core/x/supply"
)

// GenesisAccount is a struct for account initialization used exclusively during genesis
type GenesisAccount struct {
	Address       sdk.AccAddress `json:"address" yaml:"address"`
	Coins         sdk.Coins      `json:"coins" yaml:"coins"`
	Sequence      uint64         `json:"sequence_number" yaml:"sequence_number"`
	AccountNumber uint64         `json:"account_number" yaml:"account_number"`

	// vesting account fields
	OriginalVesting  sdk.Coins `json:"original_vesting" yaml:"original_vesting"`   // total vesting coins upon initialization
	DelegatedFree    sdk.Coins `json:"delegated_free" yaml:"delegated_free"`       // delegated vested coins at time of delegation
	DelegatedVesting sdk.Coins `json:"delegated_vesting" yaml:"delegated_vesting"` // delegated vesting coins at time of delegation
	StartTime        int64     `json:"start_time" yaml:"start_time"`               // vesting start time (UNIX Epoch time)
	EndTime          int64     `json:"end_time" yaml:"end_time"`                   // vesting end time (UNIX Epoch time)

	// for lazy vesting account
	VestingSchedules []auth.VestingSchedule `json:"vesting_schedules" yaml:"vesting_schedules"` // lazy vesting schedules

	// module account fields
	ModuleName        string   `json:"module_name" yaml:"module_name"`               // name of the module account
	ModulePermissions []string `json:"module_permissions" yaml:"module_permissions"` // permissions of module account
}

// Validate checks for errors on the vesting and module account parameters
func (ga GenesisAccount) Validate() error {
	if !ga.OriginalVesting.IsZero() {
		if ga.OriginalVesting.IsAnyGT(ga.Coins) {
			return errors.New("vesting amount cannot be greater than total amount")
		}

		if len(ga.VestingSchedules) > 0 { // lazy vesting account
			for _, schedule := range ga.VestingSchedules {
				if !schedule.IsValid() {
					return errors.New("invalid lazy vesting schedule")
				}
			}
		} else if ga.StartTime >= ga.EndTime { // or normal vesting account

			return errors.New("vesting start-time cannot be before end-time")
		}
	}

	// don't allow blank (i.e just whitespaces) on the module name
	if ga.ModuleName != "" && strings.TrimSpace(ga.ModuleName) == "" {
		return errors.New("module account name cannot be blank")
	}

	return nil
}

// NewGenesisAccountRaw creates a new GenesisAccount object
func NewGenesisAccountRaw(address sdk.AccAddress, coins,
	vestingAmount sdk.Coins, vestingStartTime, vestingEndTime int64,
	lazyVestingSchedules []auth.VestingSchedule, module string, permissions ...string) GenesisAccount {

	return GenesisAccount{
		Address:           address,
		Coins:             coins,
		Sequence:          0,
		AccountNumber:     0, // ignored set by the account keeper during InitGenesis
		OriginalVesting:   vestingAmount,
		DelegatedFree:     sdk.Coins{}, // ignored
		DelegatedVesting:  sdk.Coins{}, // ignored
		StartTime:         vestingStartTime,
		EndTime:           vestingEndTime,
		VestingSchedules:  lazyVestingSchedules,
		ModuleName:        module,
		ModulePermissions: permissions,
	}
}

// NewGenesisAccount creates a GenesisAccount instance from a BaseAccount.
func NewGenesisAccount(acc *auth.BaseAccount) GenesisAccount {
	return GenesisAccount{
		Address:       acc.Address,
		Coins:         acc.Coins,
		AccountNumber: acc.AccountNumber,
		Sequence:      acc.Sequence,
	}
}

// NewGenesisAccountI creates a GenesisAccount instance from an Account interface.
func NewGenesisAccountI(acc auth.Account) (GenesisAccount, error) {
	gacc := GenesisAccount{
		Address:       acc.GetAddress(),
		Coins:         acc.GetCoins(),
		AccountNumber: acc.GetAccountNumber(),
		Sequence:      acc.GetSequence(),
	}

	if err := gacc.Validate(); err != nil {
		return gacc, err
	}

	switch acc := acc.(type) {
	case auth.LazyGradedVestingAccount:
		gacc.OriginalVesting = acc.GetOriginalVesting()
		gacc.DelegatedFree = acc.GetDelegatedFree()
		gacc.DelegatedVesting = acc.GetDelegatedVesting()
		gacc.VestingSchedules = acc.GetVestingSchedules()
	case auth.VestingAccount:
		gacc.OriginalVesting = acc.GetOriginalVesting()
		gacc.DelegatedFree = acc.GetDelegatedFree()
		gacc.DelegatedVesting = acc.GetDelegatedVesting()
		gacc.StartTime = acc.GetStartTime()
		gacc.EndTime = acc.GetEndTime()
	case supply.ModuleAccountI:
		gacc.ModuleName = acc.GetName()
		gacc.ModulePermissions = acc.GetPermissions()
	}

	return gacc, nil
}

// ToAccount converts a GenesisAccount to an Account interface
func (ga *GenesisAccount) ToAccount() auth.Account {
	bacc := auth.NewBaseAccount(ga.Address, ga.Coins.Sort(), nil, ga.AccountNumber, ga.Sequence)

	// vesting accounts
	if !ga.OriginalVesting.IsZero() {
		baseVestingAcc := auth.NewBaseVestingAccount(
			bacc, ga.OriginalVesting, ga.DelegatedFree,
			ga.DelegatedVesting, ga.EndTime,
		)

		switch {
		case len(ga.VestingSchedules) != 0:
			return auth.NewBaseLazyGradedVestingAccountRaw(baseVestingAcc, ga.VestingSchedules)
		case ga.StartTime != 0 && ga.EndTime != 0:
			return auth.NewContinuousVestingAccountRaw(baseVestingAcc, ga.StartTime)
		case ga.EndTime != 0:
			return auth.NewDelayedVestingAccountRaw(baseVestingAcc)
		default:
			panic(fmt.Sprintf("invalid genesis vesting account: %+v", ga))
		}
	}

	// module accounts
	if ga.ModuleName != "" {
		return supply.NewModuleAccount(bacc, ga.ModuleName, ga.ModulePermissions...)
	}

	return bacc
}

//___________________________________
type GenesisAccounts []GenesisAccount

// genesis accounts contain an address
func (gaccs GenesisAccounts) Contains(acc sdk.AccAddress) bool {
	for _, gacc := range gaccs {
		if gacc.Address.Equals(acc) {
			return true
		}
	}
	return false
}
