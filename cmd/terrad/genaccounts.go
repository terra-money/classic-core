package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	authvesting "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"

	vestingtypes "github.com/terra-money/core/x/vesting/types"
)

const (
	flagVestingAmt       = "vesting-amount"
	flagVestingSchedules = "vesting-schedules"
)

// AddGenesisAccountCmd returns add-genesis-account cobra Command.
func AddGenesisAccountCmd(defaultNodeHome string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-genesis-account [address_or_key_name] [coin][,[coin]]",
		Short: "Add a genesis account to genesis.json",
		Long: `Add a genesis account to genesis.json. The provided account must specify
the account address or key name and a list of initial coins. If a key name is given,
the address will be looked up in the local Keybase. The list of initial tokens must
contain valid denominations. Accounts may optionally be supplied with vesting parameters.

It supports adding LazyGradedVestingAccount with args. 
'start' and 'end' must be specified with number of days from the genesis

Example:
$ terrad add-genesis-account acc1 '10000000000uluna,1000000ukrw' \
  --vesting-amount '10000000000uluna,1000000ukrw' \
  --vesting-schedules 'uluna|30|60|0.1,ukrw|0|30|1'
  
Or add normal account

Example:
$ terrad add-genesis-account acc1 '10000000000uluna,1000000ukrw'

`,
		Args: cobra.RangeArgs(2, 3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			cdc := clientCtx.Codec

			serverCtx := server.GetServerContextFromCmd(cmd)
			config := serverCtx.Config

			config.SetRoot(clientCtx.HomeDir)

			genFile := config.GenesisFile()
			appState, genDoc, err := genutiltypes.GenesisStateFromGenFile(genFile)
			if err != nil {
				return fmt.Errorf("failed to unmarshal genesis state: %w", err)
			}

			addr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				inBuf := bufio.NewReader(cmd.InOrStdin())
				keyringBackend, err := cmd.Flags().GetString(flags.FlagKeyringBackend)
				if err != nil {
					return err
				}

				// attempt to lookup address from Keybase if no address was provided
				kb, err := keyring.New(sdk.KeyringServiceName(), keyringBackend, clientCtx.HomeDir, inBuf)
				if err != nil {
					return err
				}

				info, err := kb.Key(args[0])
				if err != nil {
					return fmt.Errorf("failed to get address from Keybase: %w", err)
				}

				addr = info.GetAddress()
			}

			coins, err := sdk.ParseCoinsNormalized(args[1])
			if err != nil {
				return fmt.Errorf("failed to parse coins: %w", err)
			}

			vestingAmtStr, err := cmd.Flags().GetString(flagVestingAmt)
			if err != nil {
				return err
			}

			vestingSchedulesStr, err := cmd.Flags().GetString(flagVestingSchedules)
			if err != nil {
				return err
			}

			vestingAmt, err := sdk.ParseCoinsNormalized(vestingAmtStr)
			if err != nil {
				return fmt.Errorf("failed to parse vesting amount: %w", err)
			}

			// create concrete account type based on input parameters
			var genAccount authtypes.GenesisAccount

			balances := banktypes.Balance{Address: addr.String(), Coins: coins.Sort()}
			baseAccount := authtypes.NewBaseAccount(addr, nil, 0, 0)

			if !vestingAmt.IsZero() {
				// Build vesting account
				genesisTime := genDoc.GenesisTime
				vestingSchedulesDenomMap := make(map[string]*vestingtypes.VestingSchedule)
				unparsedSchedules := strings.Split(vestingSchedulesStr, ",")
				for _, unparsedSchedule := range unparsedSchedules {
					items := strings.Split(unparsedSchedule, "|")
					if len(items) != 4 {
						return errors.New("vesting schedule parse error")
					}

					denom := items[0]
					startDay, err := strconv.Atoi(items[1])
					if err != nil {
						return err
					}
					endDay, err := strconv.Atoi(items[2])
					if err != nil {
						return err
					}
					ratio, err := sdk.NewDecFromStr(items[3])
					if err != nil {
						return err
					}

					lazySchedule := vestingtypes.Schedule{
						StartTime: genesisTime.AddDate(0, 0, startDay).UTC().Unix(),
						EndTime:   genesisTime.AddDate(0, 0, endDay).UTC().Unix(),
						Ratio:     ratio,
					}

					if vs, ok := vestingSchedulesDenomMap[denom]; ok {
						vs.Schedules = append(vs.Schedules, lazySchedule)
					} else {
						vestingSchedulesDenomMap[denom] = &vestingtypes.VestingSchedule{Denom: denom, Schedules: vestingtypes.Schedules{lazySchedule}}
					}
				}

				vestingSchedules := vestingtypes.VestingSchedules{}
				for denom, schedule := range vestingSchedulesDenomMap {
					schedule.Denom = denom

					vestingSchedules = append(vestingSchedules, *schedule)
				}

				baseVestingAccount := authvesting.NewBaseVestingAccount(baseAccount, vestingAmt.Sort(), 0)

				for _, coin := range vestingAmt {
					if _, ok := vestingSchedulesDenomMap[coin.Denom]; !ok {
						return errors.New("must provide vesting schedule for all vesting amount")
					}
				}

				if (balances.Coins.IsZero() && !baseVestingAccount.OriginalVesting.IsZero()) ||
					baseVestingAccount.OriginalVesting.IsAnyGT(balances.Coins) {
					return errors.New("vesting amount cannot be greater than total amount")
				}

				genAccount = vestingtypes.NewLazyGradedVestingAccountRaw(baseVestingAccount, vestingSchedules)
			} else {
				genAccount = baseAccount
			}

			if err := genAccount.Validate(); err != nil {
				return fmt.Errorf("failed to validate new genesis account: %w", err)
			}

			authGenState := authtypes.GetGenesisStateFromAppState(cdc, appState)

			accs, err := authtypes.UnpackAccounts(authGenState.Accounts)
			if err != nil {
				return fmt.Errorf("failed to get accounts from any: %w", err)
			}

			if accs.Contains(addr) {
				return fmt.Errorf("cannot add account at existing address %s", addr)
			}

			// Add the new account to the set of genesis accounts and sanitize the
			// accounts afterwards.
			accs = append(accs, genAccount)
			accs = authtypes.SanitizeGenesisAccounts(accs)

			genAccs, err := authtypes.PackAccounts(accs)
			if err != nil {
				return fmt.Errorf("failed to convert accounts into any's: %w", err)
			}
			authGenState.Accounts = genAccs

			authGenStateBz, err := cdc.MarshalJSON(&authGenState)
			if err != nil {
				return fmt.Errorf("failed to marshal auth genesis state: %w", err)
			}

			appState[authtypes.ModuleName] = authGenStateBz

			bankGenState := banktypes.GetGenesisStateFromAppState(cdc, appState)
			bankGenState.Balances = append(bankGenState.Balances, balances)
			bankGenState.Balances = banktypes.SanitizeGenesisBalances(bankGenState.Balances)
			bankGenState.Supply = bankGenState.Supply.Add(balances.Coins...)

			bankGenStateBz, err := cdc.MarshalJSON(bankGenState)
			if err != nil {
				return fmt.Errorf("failed to marshal bank genesis state: %w", err)
			}

			appState[banktypes.ModuleName] = bankGenStateBz

			appStateJSON, err := json.Marshal(appState)
			if err != nil {
				return fmt.Errorf("failed to marshal application genesis state: %w", err)
			}

			genDoc.AppState = appStateJSON
			return genutil.ExportGenesisFile(genDoc, genFile)
		},
	}

	cmd.Flags().String(flags.FlagHome, defaultNodeHome, "The application home directory")
	cmd.Flags().String(flags.FlagKeyringBackend, flags.DefaultKeyringBackend, "Select keyring's backend (os|file|kwallet|pass|test)")
	cmd.Flags().String(flagVestingAmt, "", "amount of coins for vesting accounts")
	cmd.Flags().String(flagVestingSchedules, "", "comma separated vesting schedules [denom|start|end|ratio][,[denom|start|end|ratio]], where 'start' and 'end' is day unit")
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}
