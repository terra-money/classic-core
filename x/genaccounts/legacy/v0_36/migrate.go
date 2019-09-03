// DONTCOVER
// nolint
package v0_36

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	v034distr "github.com/cosmos/cosmos-sdk/x/distribution/legacy/v0_34"
	v034staking "github.com/cosmos/cosmos-sdk/x/staking/legacy/v0_34"

	v034accounts "github.com/terra-project/core/x/genaccounts/legacy/v0_34"

	"github.com/tendermint/tendermint/crypto"
)

const (
	notBondedPoolName  = "not_bonded_tokens_pool"
	bondedPoolName     = "bonded_tokens_pool"
	feeCollectorName   = "fee_collector"
	oracleModuleName   = "oracle"
	marketModuleName   = "market"
	treasuryModuleName = "treasury"

	basic   = "basic"
	minter  = "minter"
	burner  = "burner"
	staking = "staking"
)

// Migrate accepts exported genesis state from v0.34 and migrates it to v0.36
// genesis state. It deletes the governance base accounts and creates the new module accounts.
// The remaining accounts are updated to the new GenesisAccount type from 0.36
func Migrate(
	oldGenState v034accounts.GenesisState, fees sdk.Coins, communityPool sdk.DecCoins,
	vals v034staking.Validators, ubds []v034staking.UnbondingDelegation,
	valOutRewards []v034distr.ValidatorOutstandingRewardsRecord, bondDenom, distrModuleName, govModuleName string,
) GenesisState {

	bondedAmt := sdk.ZeroInt()
	notBondedAmt := sdk.ZeroInt()

	// remove the two previous governance base accounts for deposits and burned
	// coins from rejected proposals add six new module accounts:
	// distribution, gov, mint, fee collector, bonded and not bonded pool
	var (
		newGenState   GenesisState
		govCoins      sdk.Coins
		extraAccounts = 8
	)

	for _, acc := range oldGenState {
		newGenState = append(
			newGenState,
			NewGenesisAccount(
				acc.Address, acc.Coins, acc.Sequence,
				acc.OriginalVesting, acc.DelegatedFree, acc.DelegatedVesting,
				acc.StartTime, acc.EndTime, acc.VestingSchedules, "", []string{},
			),
		)
	}

	// get staking module accounts coins
	for _, validator := range vals {
		switch validator.Status {
		case sdk.Bonded:
			bondedAmt = bondedAmt.Add(validator.Tokens)

		case sdk.Unbonding, sdk.Unbonded:
			notBondedAmt = notBondedAmt.Add(validator.Tokens)

		default:
			panic("invalid validator status")
		}
	}

	for _, ubd := range ubds {
		for _, entry := range ubd.Entries {
			notBondedAmt = notBondedAmt.Add(entry.Balance)
		}
	}

	bondedCoins := sdk.NewCoins(sdk.NewCoin(bondDenom, bondedAmt))
	notBondedCoins := sdk.NewCoins(sdk.NewCoin(bondDenom, notBondedAmt))

	// get distr module account coins
	var distrDecCoins sdk.DecCoins
	for _, reward := range valOutRewards {
		distrDecCoins = distrDecCoins.Add(reward.OutstandingRewards)
	}

	distrCoins, _ := distrDecCoins.Add(communityPool).TruncateDecimal()

	// get module account addresses
	feeCollectorAddr := sdk.AccAddress(crypto.AddressHash([]byte(feeCollectorName)))
	govAddr := sdk.AccAddress(crypto.AddressHash([]byte(govModuleName)))
	bondedAddr := sdk.AccAddress(crypto.AddressHash([]byte(bondedPoolName)))
	notBondedAddr := sdk.AccAddress(crypto.AddressHash([]byte(notBondedPoolName)))
	distrAddr := sdk.AccAddress(crypto.AddressHash([]byte(distrModuleName)))
	oracleAddr := sdk.AccAddress(crypto.AddressHash([]byte(oracleModuleName)))
	marketAddr := sdk.AccAddress(crypto.AddressHash([]byte(marketModuleName)))
	treasuryAddr := sdk.AccAddress(crypto.AddressHash([]byte(treasuryModuleName)))

	// create module genesis accounts
	feeCollectorModuleAcc := NewGenesisAccount(
		feeCollectorAddr, fees, 0,
		sdk.Coins{}, sdk.Coins{}, sdk.Coins{},
		0, 0, nil, feeCollectorName, []string{basic},
	)
	govModuleAcc := NewGenesisAccount(
		govAddr, govCoins, 0,
		sdk.Coins{}, sdk.Coins{}, sdk.Coins{},
		0, 0, nil, govModuleName, []string{burner},
	)
	distrModuleAcc := NewGenesisAccount(
		distrAddr, distrCoins, 0,
		sdk.Coins{}, sdk.Coins{}, sdk.Coins{},
		0, 0, nil, distrModuleName, []string{basic},
	)
	bondedModuleAcc := NewGenesisAccount(
		bondedAddr, bondedCoins, 0,
		sdk.Coins{}, sdk.Coins{}, sdk.Coins{},
		0, 0, nil, bondedPoolName, []string{burner, staking},
	)
	notBondedModuleAcc := NewGenesisAccount(
		notBondedAddr, notBondedCoins, 0,
		sdk.Coins{}, sdk.Coins{}, sdk.Coins{},
		0, 0, nil, notBondedPoolName, []string{burner, staking},
	)
	oracleModuleAcc := NewGenesisAccount(
		oracleAddr, sdk.Coins{}, 0,
		sdk.Coins{}, sdk.Coins{}, sdk.Coins{},
		0, 0, nil, oracleModuleName, []string{basic},
	)
	marketModuleAcc := NewGenesisAccount(
		marketAddr, sdk.Coins{}, 0,
		sdk.Coins{}, sdk.Coins{}, sdk.Coins{},
		0, 0, nil, marketModuleName, []string{minter, burner},
	)
	treasuryModuleAcc := NewGenesisAccount(
		treasuryAddr, sdk.Coins{}, 0,
		sdk.Coins{}, sdk.Coins{}, sdk.Coins{},
		0, 0, nil, treasuryModuleName, []string{minter},
	)

	newGenState = append(
		newGenState,
		[]GenesisAccount{
			feeCollectorModuleAcc, govModuleAcc, distrModuleAcc, bondedModuleAcc,
			notBondedModuleAcc, oracleModuleAcc, marketModuleAcc, treasuryModuleAcc,
		}...,
	)

	// verify the total number of accounts is correct
	if len(newGenState) != len(oldGenState)+extraAccounts {
		panic(
			fmt.Sprintf(
				"invalid total number of genesis accounts; got: %d, expected: %d",
				len(newGenState), len(oldGenState)+extraAccounts),
		)
	}

	return newGenState
}
