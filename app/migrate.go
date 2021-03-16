package app

// TODO change this script to columbus-4 to columbus-5 migration script
//This file implements a genesis migration from cosmoshub-3 to cosmoshub-4. It migrates state from the modules in cosmoshub-3.
//This file also implements setting an initial height from an upgrade.

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bank "github.com/cosmos/cosmos-sdk/x/bank/types"
	captypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	evtypes "github.com/cosmos/cosmos-sdk/x/evidence/types"
	"github.com/cosmos/cosmos-sdk/x/genutil/client/cli"
	"github.com/cosmos/cosmos-sdk/x/genutil/types"
	ibcxfertypes "github.com/cosmos/cosmos-sdk/x/ibc/applications/transfer/types"
	host "github.com/cosmos/cosmos-sdk/x/ibc/core/24-host"
	"github.com/cosmos/cosmos-sdk/x/ibc/core/exported"
	ibccoretypes "github.com/cosmos/cosmos-sdk/x/ibc/core/types"
	staking "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	tmjson "github.com/tendermint/tendermint/libs/json"
	tmtypes "github.com/tendermint/tendermint/types"
)

const (
	flagGenesisTime     = "genesis-time"
	flagInitialHeight   = "initial-height"
	flagReplacementKeys = "replacement-cons-keys"
	flagNoProp29        = "no-prop-29"
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

			firstMigration := "v0.38"
			importGenesis := args[0]

			jsonBlob, err := ioutil.ReadFile(importGenesis)

			if err != nil {
				return errors.Wrap(err, "failed to read provided genesis file")
			}

			jsonBlob, err = migrateTendermintGenesis(jsonBlob)

			if err != nil {
				return errors.Wrap(err, "failed to migration from 0.32 Tendermint params to 0.34 parms")
			}

			genDoc, err := tmtypes.GenesisDocFromJSON(jsonBlob)
			if err != nil {
				return errors.Wrapf(err, "failed to read genesis document from file %s", importGenesis)
			}

			var initialState types.AppMap
			if err := json.Unmarshal(genDoc.AppState, &initialState); err != nil {
				return errors.Wrap(err, "failed to JSON unmarshal initial genesis state")
			}

			migrationFunc := cli.GetMigrationCallback(firstMigration)
			if migrationFunc == nil {
				return fmt.Errorf("unknown migration function for version: %s", firstMigration)
			}

			// TODO: handler error from migrationFunc call
			newGenState := migrationFunc(initialState, clientCtx)

			secondMigration := "v0.39"

			migrationFunc = cli.GetMigrationCallback(secondMigration)
			if migrationFunc == nil {
				return fmt.Errorf("unknown migration function for version: %s", secondMigration)
			}

			// TODO: handler error from migrationFunc call
			newGenState = migrationFunc(newGenState, clientCtx)

			thirdMigration := "v0.40"

			migrationFunc = cli.GetMigrationCallback(thirdMigration)
			if migrationFunc == nil {
				return fmt.Errorf("unknown migration function for version: %s", thirdMigration)
			}

			// TODO: handler error from migrationFunc call
			newGenState = migrationFunc(newGenState, clientCtx)

			var bankGenesis bank.GenesisState

			clientCtx.JSONMarshaler.MustUnmarshalJSON(newGenState[bank.ModuleName], &bankGenesis)

			bankGenesis.DenomMetadata = []bank.Metadata{
				{
					Description: "The native staking token of the Cosmos Hub.",
					DenomUnits: []*bank.DenomUnit{
						{Denom: "uatom", Exponent: uint32(0), Aliases: []string{"microatom"}},
						{Denom: "matom", Exponent: uint32(3), Aliases: []string{"milliatom"}},
						{Denom: "atom", Exponent: uint32(6), Aliases: []string{}},
					},
					Base:    "uatom",
					Display: "atom",
				},
			}
			newGenState[bank.ModuleName] = clientCtx.JSONMarshaler.MustMarshalJSON(&bankGenesis)

			var stakingGenesis staking.GenesisState

			clientCtx.JSONMarshaler.MustUnmarshalJSON(newGenState[staking.ModuleName], &stakingGenesis)

			ibcTransferGenesis := ibcxfertypes.DefaultGenesisState()
			ibcCoreGenesis := ibccoretypes.DefaultGenesisState()
			capGenesis := captypes.DefaultGenesis()
			evGenesis := evtypes.DefaultGenesisState()

			ibcTransferGenesis.Params.ReceiveEnabled = false
			ibcTransferGenesis.Params.SendEnabled = false

			ibcCoreGenesis.ClientGenesis.Params.AllowedClients = []string{exported.Tendermint}
			stakingGenesis.Params.HistoricalEntries = 10000

			newGenState[ibcxfertypes.ModuleName] = clientCtx.JSONMarshaler.MustMarshalJSON(ibcTransferGenesis)
			newGenState[host.ModuleName] = clientCtx.JSONMarshaler.MustMarshalJSON(ibcCoreGenesis)
			newGenState[captypes.ModuleName] = clientCtx.JSONMarshaler.MustMarshalJSON(capGenesis)
			newGenState[evtypes.ModuleName] = clientCtx.JSONMarshaler.MustMarshalJSON(evGenesis)
			newGenState[staking.ModuleName] = clientCtx.JSONMarshaler.MustMarshalJSON(&stakingGenesis)

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
	cmd.Flags().String(flagReplacementKeys, "", "Proviide a JSON file to replace the consensus keys of validators")
	cmd.Flags().String(flags.FlagChainID, "", "override chain_id with this flag")
	cmd.Flags().Bool(flagNoProp29, false, "Do not implement fund recovery from prop29")

	return cmd
}

// MigrateTendermintGenesis makes sure a later version of Tendermint can parse
// a JSON blob exported by an older version of Tendermint.
func migrateTendermintGenesis(jsonBlob []byte) ([]byte, error) {
	var jsonObj map[string]interface{}
	err := json.Unmarshal(jsonBlob, &jsonObj)
	if err != nil {
		return nil, err
	}

	consensusParams, ok := jsonObj["consensus_params"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("exported json does not contain consensus_params field")
	}
	evidenceParams, ok := consensusParams["evidence"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("exported json does not contain consensus_params.evidence field")

	}

	evidenceParams["max_age_num_blocks"] = evidenceParams["max_age"]
	delete(evidenceParams, "max_age")

	evidenceParams["max_age_duration"] = "172800000000000"
	evidenceParams["max_bytes"] = "50000"

	jsonBlob, err = json.Marshal(jsonObj)

	if err != nil {
		return nil, errors.Wrapf(err, "Error resserializing JSON blob after tendermint migrations")
	}

	return jsonBlob, nil
}
