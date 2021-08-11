package legacy

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	tmjson "github.com/tendermint/tendermint/libs/json"
	tmtypes "github.com/tendermint/tendermint/types"

	ibcxfertypes "github.com/cosmos/ibc-go/modules/apps/transfer/types"
	host "github.com/cosmos/ibc-go/modules/core/24-host"
	"github.com/cosmos/ibc-go/modules/core/exported"
	ibccoretypes "github.com/cosmos/ibc-go/modules/core/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	captypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	evtypes "github.com/cosmos/cosmos-sdk/x/evidence/types"
	"github.com/cosmos/cosmos-sdk/x/genutil/types"
	staking "github.com/cosmos/cosmos-sdk/x/staking/types"

	legacy05 "github.com/terra-money/core/app/legacy/v05"
	oracletypes "github.com/terra-money/core/x/oracle/types"
)

const (
	flagGenesisTime     = "genesis-time"
	flagInitialHeight   = "initial-height"
	flagReplacementKeys = "replacement-cons-keys"
)

// MigrateGenesisCmd returns a command to execute genesis state migration.
func MigrateGenesisCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "migrate [genesis-file]",
		Short: "Migrate genesis to a specified target version",
		Long: `Migrate the source genesis into the target version and print to STDOUT.

Example:
$ terrad migrate /path/to/genesis.json --chain-id=cosmoshub-4 --genesis-time=2019-04-22T17:00:00Z --initial-height=5000
`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			var err error

			importGenesis := args[0]

			jsonBlob, err := ioutil.ReadFile(importGenesis)
			if err != nil {
				return errors.Wrap(err, "failed to read provided genesis file")
			}

			genDoc, err := tmtypes.GenesisDocFromJSON(jsonBlob)
			if err != nil {
				return errors.Wrapf(err, "failed to read genesis document from file %s", importGenesis)
			}

			// increase block consensus params
			genDoc.ConsensusParams.Block.MaxBytes = int64(5_000_000)
			genDoc.ConsensusParams.Block.MaxGas = int64(1_000_000_000)

			// decrease evidence max bytes
			genDoc.ConsensusParams.Evidence.MaxBytes = int64(50000)

			var initialState types.AppMap
			if err := json.Unmarshal(genDoc.AppState, &initialState); err != nil {
				return errors.Wrap(err, "failed to JSON unmarshal initial genesis state")
			}

			// Migrate Terra specific state
			newGenState := legacy05.Migrate(initialState, clientCtx)

			var bankGenesis banktypes.GenesisState

			clientCtx.Codec.MustUnmarshalJSON(newGenState[banktypes.ModuleName], &bankGenesis)

			var oracleGenesis oracletypes.GenesisState
			clientCtx.Codec.MustUnmarshalJSON(newGenState[oracletypes.ModuleName], &oracleGenesis)

			// Register whitelist denom list
			denomMetadata := make([]banktypes.Metadata, len(oracleGenesis.Params.Whitelist)+1)
			denomMetadata[0] = banktypes.Metadata{
				Description: "The native staking token of the Terra Columbus.",
				DenomUnits: []*banktypes.DenomUnit{
					{Denom: "uluna", Exponent: uint32(0), Aliases: []string{"microluna"}},
					{Denom: "mluna", Exponent: uint32(3), Aliases: []string{"milliluna"}},
					{Denom: "luna", Exponent: uint32(6), Aliases: []string{}},
				},
				Base:    "uluna",
				Display: "luna",
				Name:    "LUNA",
				Symbol:  "LUNA",
			}

			for i, w := range oracleGenesis.Params.Whitelist {
				base := w.Name
				display := base[1:]
				denomMetadata[i+1] = banktypes.Metadata{
					Description: "The native stable token of the Terra Columbus.",
					DenomUnits: []*banktypes.DenomUnit{
						{Denom: "u" + display, Exponent: uint32(0), Aliases: []string{"micro" + display}},
						{Denom: "m" + display, Exponent: uint32(3), Aliases: []string{"milli" + display}},
						{Denom: display, Exponent: uint32(6), Aliases: []string{}},
					},
					Base:    base,
					Display: display,
					Name:    fmt.Sprintf("%s TERRA", strings.ToUpper(display)),
					Symbol:  fmt.Sprintf("%sT", strings.ToUpper(display[:len(display)-1])),
				}
			}

			bankGenesis.DenomMetadata = denomMetadata
			newGenState[banktypes.ModuleName] = clientCtx.Codec.MustMarshalJSON(&bankGenesis)

			var stakingGenesis staking.GenesisState

			clientCtx.Codec.MustUnmarshalJSON(newGenState[staking.ModuleName], &stakingGenesis)

			ibcTransferGenesis := ibcxfertypes.DefaultGenesisState()
			ibcCoreGenesis := ibccoretypes.DefaultGenesisState()
			capGenesis := captypes.DefaultGenesis()
			evGenesis := evtypes.DefaultGenesisState()

			ibcTransferGenesis.Params.ReceiveEnabled = false
			ibcTransferGenesis.Params.SendEnabled = false

			ibcCoreGenesis.ClientGenesis.Params.AllowedClients = []string{exported.Tendermint}
			stakingGenesis.Params.HistoricalEntries = 10000

			newGenState[ibcxfertypes.ModuleName] = clientCtx.Codec.MustMarshalJSON(ibcTransferGenesis)
			newGenState[host.ModuleName] = clientCtx.Codec.MustMarshalJSON(ibcCoreGenesis)
			newGenState[captypes.ModuleName] = clientCtx.Codec.MustMarshalJSON(capGenesis)
			newGenState[evtypes.ModuleName] = clientCtx.Codec.MustMarshalJSON(evGenesis)
			newGenState[staking.ModuleName] = clientCtx.Codec.MustMarshalJSON(&stakingGenesis)

			genDoc.AppState, err = json.Marshal(newGenState)
			if err != nil {
				return errors.Wrap(err, "failed to JSON marshal migrated genesis state")
			}

			genesisTime, _ := cmd.Flags().GetString(flagGenesisTime)
			if genesisTime != "" {
				var t time.Time

				err := t.UnmarshalText([]byte(genesisTime))
				if err != nil {
					return errors.Wrap(err, "failed to unmarshal genesis time")
				}

				genDoc.GenesisTime = t
			}

			chainID, _ := cmd.Flags().GetString(flags.FlagChainID)
			if chainID != "" {
				genDoc.ChainID = chainID
			}

			initialHeight, _ := cmd.Flags().GetInt(flagInitialHeight)

			genDoc.InitialHeight = int64(initialHeight)

			replacementKeys, _ := cmd.Flags().GetString(flagReplacementKeys)

			if replacementKeys != "" {
				genDoc = loadKeydataFromFile(clientCtx, replacementKeys, genDoc)
			}

			bz, err := tmjson.Marshal(genDoc)
			if err != nil {
				return errors.Wrap(err, "failed to marshal genesis doc")
			}

			sortedBz, err := sdk.SortJSON(bz)
			if err != nil {
				return errors.Wrap(err, "failed to sort JSON genesis doc")
			}

			fmt.Println(string(sortedBz))
			return nil
		},
	}

	cmd.Flags().String(flagGenesisTime, "", "override genesis_time with this flag")
	cmd.Flags().Int(flagInitialHeight, 0, "Set the starting height for the chain")
	cmd.Flags().String(flagReplacementKeys, "", "Provide a JSON file to replace the consensus keys of validators")
	cmd.Flags().String(flags.FlagChainID, "", "override chain_id with this flag")

	return cmd
}
