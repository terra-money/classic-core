package v040

import (
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	kmultisig "github.com/cosmos/cosmos-sdk/crypto/keys/multisig"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	v039auth "github.com/cosmos/cosmos-sdk/x/auth/legacy/v039"
	v040auth "github.com/cosmos/cosmos-sdk/x/auth/types"
	v040authvesting "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	v40mint "github.com/cosmos/cosmos-sdk/x/mint/types"

	v039authcustom "github.com/terra-money/core/custom/auth/legacy/v039"
	v40treasury "github.com/terra-money/core/x/treasury/types"
	v040vesting "github.com/terra-money/core/x/vesting/types"
)

// convertBaseAccount converts a 0.39 BaseAccount to a 0.40 BaseAccount.
func convertBaseAccount(old *v039auth.BaseAccount) *v040auth.BaseAccount {
	var any *codectypes.Any

	// If the old genesis had a pubkey, we pack it inside an Any. Or else, we
	// just leave it nil.
	if old.PubKey != nil {
		var pk cryptotypes.PubKey
		if mpk, ok := old.PubKey.(*v039authcustom.LegacyAminoPubKey); ok {
			pk = kmultisig.NewLegacyAminoPubKey(int(mpk.Threshold.Int64()), mpk.PubKeys)
		} else {
			pk = old.PubKey
		}

		var err error
		any, err = codectypes.NewAnyWithValue(pk)
		if err != nil {
			panic(err)
		}
	}

	return &v040auth.BaseAccount{
		Address:       old.Address.String(),
		PubKey:        any,
		AccountNumber: old.AccountNumber,
		Sequence:      old.Sequence,
	}
}

// convertBaseVestingAccount converts a 0.39 BaseVestingAccount to a 0.40 BaseVestingAccount.
func convertBaseVestingAccount(old *v039auth.BaseVestingAccount) *v040authvesting.BaseVestingAccount {
	baseAccount := convertBaseAccount(old.BaseAccount)

	return &v040authvesting.BaseVestingAccount{
		BaseAccount:      baseAccount,
		OriginalVesting:  old.OriginalVesting,
		DelegatedFree:    old.DelegatedFree,
		DelegatedVesting: old.DelegatedVesting,
		EndTime:          old.EndTime,
	}
}

// Migrate accepts exported x/auth genesis state from v0.38/v0.39 and migrates
// it to v0.40 x/auth genesis state. The migration includes:
//
// - Removing coins from account encoding.
// - Re-encode in v0.40 GenesisState.
func Migrate(authGenState v039auth.GenesisState) *v040auth.GenesisState {
	mintModuleAddress := v040auth.NewModuleAddress(v40mint.ModuleName)
	treasuryModuleAddress := v040auth.NewModuleAddress(v40treasury.ModuleName)

	// Convert v0.39 accounts to v0.40 ones.
	var v040Accounts = make([]v040auth.GenesisAccount, len(authGenState.Accounts))
	for i, v039Account := range authGenState.Accounts {
		switch v039Account := v039Account.(type) {
		case *v039auth.BaseAccount:
			{
				// columbus chain has mint address as normal account
				// so need to changes this to module account
				if !v039Account.GetAddress().Equals(mintModuleAddress) {
					v040Accounts[i] = convertBaseAccount(v039Account)
				} else {
					v040Accounts[i] = &v040auth.ModuleAccount{
						BaseAccount: convertBaseAccount(v039Account),
						Name:        v40mint.ModuleName,
						Permissions: []string{v040auth.Minter},
					}
				}
			}
		case *v039auth.ModuleAccount:
			{
				// burn permission must be added to treasury module
				permissions := v039Account.Permissions
				if v039Account.GetAddress().Equals(treasuryModuleAddress) {
					permissions = append(permissions, v040auth.Burner)
				}

				v040Accounts[i] = &v040auth.ModuleAccount{
					BaseAccount: convertBaseAccount(v039Account.BaseAccount),
					Name:        v039Account.Name,
					Permissions: permissions,
				}
			}
		case *v039auth.BaseVestingAccount:
			{
				v040Accounts[i] = convertBaseVestingAccount(v039Account)
			}
		case *v039authcustom.LazyGradedVestingAccount:
			{
				vestingSchedules := make([]v040vesting.VestingSchedule, len(v039Account.VestingSchedules))
				for j, vestingSchedule := range v039Account.VestingSchedules {
					schedules := make([]v040vesting.Schedule, len(vestingSchedule.LazySchedules))

					sumRatio := sdk.ZeroDec()
					for k, schedule := range vestingSchedule.LazySchedules {
						schedules[k] = v040vesting.Schedule{
							StartTime: schedule.StartTime,
							EndTime:   schedule.EndTime,
							Ratio:     schedule.Ratio,
						}

						sumRatio = sumRatio.Add(schedule.Ratio)
					}

					// Correct rounding error
					diff := sdk.OneDec().Sub(sumRatio)
					if !diff.IsZero() {
						lastIndex := len(vestingSchedule.LazySchedules) - 1
						schedules[lastIndex].Ratio = schedules[lastIndex].Ratio.Add(diff)
					}

					vestingSchedules[j] = v040vesting.VestingSchedule{
						Denom:     vestingSchedule.Denom,
						Schedules: schedules,
					}
				}
				v040Accounts[i] = &v040vesting.LazyGradedVestingAccount{
					BaseVestingAccount: convertBaseVestingAccount(v039Account.BaseVestingAccount),
					VestingSchedules:   vestingSchedules,
				}
			}
		default:
			panic(sdkerrors.Wrapf(sdkerrors.ErrInvalidType, "got invalid type %T", v039Account))
		}

	}

	// Convert v0.40 accounts into Anys.
	anys := make([]*codectypes.Any, len(v040Accounts))
	for i, v040Account := range v040Accounts {
		any, err := codectypes.NewAnyWithValue(v040Account)
		if err != nil {
			panic(err)
		}

		anys[i] = any
	}

	return &v040auth.GenesisState{
		Params: v040auth.Params{
			MaxMemoCharacters:      authGenState.Params.MaxMemoCharacters,
			TxSigLimit:             authGenState.Params.TxSigLimit,
			TxSizeCostPerByte:      authGenState.Params.TxSizeCostPerByte,
			SigVerifyCostED25519:   authGenState.Params.SigVerifyCostED25519,
			SigVerifyCostSecp256k1: authGenState.Params.SigVerifyCostSecp256k1,
		},
		Accounts: anys,
	}
}
