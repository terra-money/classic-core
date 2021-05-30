package main

import (
	"bufio"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/tendermint/tendermint/libs/cli"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keys"
	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authexported "github.com/cosmos/cosmos-sdk/x/auth/exported"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/genutil"

	"github.com/terra-money/core/x/auth/vesting"
)

const (
	flagClientHome = "home-client"
)

// AddGenesisAccountCmd returns add-genesis-account cobra Command.
func AddGenesisAccountCmd(ctx *server.Context, cdc *codec.Codec,
	defaultNodeHome, defaultClientHome string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-genesis-account [address_or_key_name] [coin][,[coin]] [denom|start|end|ratio][,[denom|start|end|ratio]]",
		Short: "Add genesis account to genesis.json",
		Long: strings.TrimSpace(
			fmt.Sprintf(`
Add genesis account or vesting account to genesis.json. 
It supports adding LazyGradedVestingAccount with args. 
'start' and 'end' must be specified with number of days from the genesis

Example:
$ %s add-genesis-account acc1 '10000000000uluna,1000000ukrw' 'uluna|30|60|0.1,ukrw|0|30|1'  

Or add normal account

Eaxmple:
$ %s add-genesis-account acc1 '10000000000uluna,1000000ukrw'
`, version.ServerName, version.ServerName),
		),
		Args: cobra.RangeArgs(2, 3),
		RunE: func(cmd *cobra.Command, args []string) error {
			config := ctx.Config
			config.SetRoot(viper.GetString(cli.HomeFlag))

			inBuf := bufio.NewReader(cmd.InOrStdin())
			addr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				// attempt to lookup address from Keybase if no address was provided
				kb, err := keys.NewKeyring(
					sdk.KeyringServiceName(),
					viper.GetString(flags.FlagKeyringBackend),
					viper.GetString(flagClientHome),
					inBuf,
				)
				if err != nil {
					return err
				}

				info, err := kb.Get(args[0])
				if err != nil {
					return fmt.Errorf("failed to get address from Keybase: %w", err)
				}

				addr = info.GetAddress()
			}

			coins, err := sdk.ParseCoins(args[1])
			if err != nil {
				return err
			}

			// retrieve the app state
			genFile := config.GenesisFile()
			appState, genDoc, err := genutil.GenesisStateFromGenFile(cdc, genFile)
			if err != nil {
				return err
			}

			var genAcc authexported.GenesisAccount
			if len(args) == 2 {
				// Build normal account
				acc := types.NewBaseAccountWithAddress(addr)
				acc.Coins = coins

				genAcc = &acc
				if err := genAcc.Validate(); err != nil {
					return err
				}
			} else {
				// Build vesting account
				genesisTime := genDoc.GenesisTime
				vestingSchedulesDenomMap := make(map[string]*vesting.VestingSchedule)
				unparsedSchedules := strings.Split(args[2], ",")
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

					lazySchedule := vesting.LazySchedule{
						StartTime: genesisTime.AddDate(0, 0, startDay).UTC().Unix(),
						EndTime:   genesisTime.AddDate(0, 0, endDay).UTC().Unix(),
						Ratio:     ratio,
					}

					if vs, ok := vestingSchedulesDenomMap[denom]; ok {
						vs.LazySchedules = append(vs.LazySchedules, lazySchedule)
					} else {
						vestingSchedulesDenomMap[denom] = &vesting.VestingSchedule{Denom: denom, LazySchedules: vesting.LazySchedules{lazySchedule}}
					}
				}

				vestingSchedules := vesting.VestingSchedules{}
				for denom, schedule := range vestingSchedulesDenomMap {
					schedule.Denom = denom

					vestingSchedules = append(vestingSchedules, *schedule)
				}

				acc := types.NewBaseAccountWithAddress(addr)
				acc.Coins = coins
				baseVestingAcc, err := vesting.NewBaseVestingAccount(&acc, acc.Coins, 0)
				if err != nil {
					return err
				}

				genAcc = vesting.NewLazyGradedVestingAccountRaw(baseVestingAcc, vestingSchedules)
				if err := genAcc.Validate(); err != nil {
					return err
				}
			}

			// Add genesis account to the app state
			var genesisState auth.GenesisState
			cdc.MustUnmarshalJSON(appState[auth.ModuleName], &genesisState)

			if genesisState.Accounts.Contains(addr) {
				return fmt.Errorf("cannot add account at existing address %v", addr)
			}

			genesisState.Accounts = append(genesisState.Accounts, genAcc)

			genesisStateBz := cdc.MustMarshalJSON(genesisState)
			appState[auth.ModuleName] = genesisStateBz

			appStateJSON, err := cdc.MarshalJSON(appState)
			if err != nil {
				return err
			}

			// export app state
			genDoc.AppState = appStateJSON

			return genutil.ExportGenesisFile(genDoc, genFile)
		},
	}

	cmd.Flags().String(cli.HomeFlag, defaultNodeHome, "node's home directory")
	cmd.Flags().String(flags.FlagKeyringBackend, flags.DefaultKeyringBackend, "Select keyring's backend (os|file|test)")
	cmd.Flags().String(flagClientHome, defaultClientHome, "client's home directory")
	return cmd
}
