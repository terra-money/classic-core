package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"terra/version"

	"github.com/tendermint/tendermint/p2p"

	terraInit "terra/cmd/init"

	"github.com/cosmos/cosmos-sdk/baseapp"

	"terra/app"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/cli"
	"github.com/tendermint/tendermint/libs/common"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
	tmtypes "github.com/tendermint/tendermint/types"
)

const (
	flagClientHome = "home-client"
)

func main() {
	cdc := app.MakeCodec()

	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(sdk.Bech32PrefixAccAddr, sdk.Bech32PrefixAccPub)
	config.SetBech32PrefixForValidator(sdk.Bech32PrefixValAddr, sdk.Bech32PrefixValPub)
	config.SetBech32PrefixForConsensusNode(sdk.Bech32PrefixConsAddr, sdk.Bech32PrefixConsPub)
	config.Seal()

	ctx := server.NewDefaultContext()
	cobra.EnableCommandSorting = false
	rootCmd := &cobra.Command{
		Use:               "terrad",
		Short:             "Terra Daemon (server)",
		PersistentPreRunE: server.PersistentPreRunEFn(ctx),
	}
	appInit := app.TerraAppInit()
	rootCmd.AddCommand(terraInit.InitCmd(ctx, cdc, appInit))
	rootCmd.AddCommand(terraInit.CollectGenTxsCmd(ctx, cdc))
	rootCmd.AddCommand(terraInit.TestnetFilesCmd(ctx, cdc, server.AppInit{}))
	rootCmd.AddCommand(terraInit.GenTxCmd(ctx, cdc))

	// preoccupy the version command (that will be added in server.AddCommands)
	rootCmd.AddCommand(version.VersionCmd)

	server.AddCommands(ctx, cdc, rootCmd, appInit,
		newApp, exportAppStateAndTMValidators)

	// prepare and add flags
	executor := cli.PrepareBaseCmd(rootCmd, "TE", app.DefaultNodeHome)
	err := executor.Execute()
	if err != nil {
		// handle with #870
		panic(err)
	}
}

// func main() {
// 	cdc := app.MakeCodec()
// 	ctx := server.NewDefaultContext()

// 	rootCmd := &cobra.Command{
// 		Use:               "terrad",
// 		Short:             "Terra Daemon (server)",
// 		PersistentPreRunE: server.PersistentPreRunEFn(ctx),
// 	}

// 	appInit := server.AppInit{
// 		AppGenState: TerraAppGenState,
// 	}
// 	rootCmd.AddCommand(InitCmd(ctx, cdc, appInit))

// 	server.AddCommands(ctx, cdc, rootCmd, appInit,
// 		newApp, exportAppStateAndTMValidators)

// 	// prepare and add flags
// 	rootDir := os.ExpandEnv("$HOME/.terrad")
// 	executor := cli.PrepareBaseCmd(rootCmd, "BC", rootDir)

// 	err := executor.Execute()
// 	if err != nil {
// 		// Note: Handle with #870
// 		panic(err)
// 	}
// }

// create the genesis app state
func TerraAppGenState(cdc *codec.Codec, genDoc tmtypes.GenesisDoc, appGenTxs []json.RawMessage) (
	appState json.RawMessage, err error) {

	if len(appGenTxs) != 1 {
		err = errors.New("must provide a single genesis transaction")
		return
	}

	var tx server.SimpleGenTx
	err = cdc.UnmarshalJSON(appGenTxs[0], &tx)
	if err != nil {
		return
	}

	appState = json.RawMessage(fmt.Sprintf(`{
		"accounts": [
			{
			"address": "%s",
			"coins": [
			  {
				"denom": "terra",
				"amount": "1000"
			  },
			  {
				"denom": "luna",
				"amount": "1000000000"
			  }
			],
			"sequence_number": "0",
			"account_number": "0"
		  }    
		],
		"auth": {
		  "collected_fees": null
		},
		"stake": {
		  "pool": {
			"loose_tokens": "1000000000.0000000000",
			"bonded_tokens": "0.0000000000"
		  },
		  "params": {
			"unbonding_time": "259200000000000",
			"max_validators": 100,
			"bond_denom": "luna"
		  },
		  "intra_tx_counter": 0,
		  "last_total_power": "0",
		  "validators": null,
		  "bonds": null,
		  "unbonding_delegations": null,
		  "redelegations": null
		},
		"mint": {
		  "minter": {
			"inflation_last_time": "1970-01-01T00:00:00Z",
			"inflation": "1000.00000000"
		  },
		  "params": {
			"mint_denom": "luna",
			"inflation_rate_change": "0.1000000000",
			"inflation_max": "1000.00000000",
			"inflation_min": "1000.00000000",
			"goal_bonded": "0.6700000000"
		  }
		},
		"distr": {
		  "fee_pool": {
			"val_accum": {
			  "update_height": "0",
			  "accum": "0.0000000000"
			},
			"val_pool": null,
			"community_pool": null
		  },
		  "community_tax": "0.0200000000",
		  "base_proposer_reward": "0.200000000",
		  "bonus_proposer_reward": "0.0400000000",
		  "validator_dist_infos": null,
		  "delegator_dist_infos": null,
		  "delegator_withdraw_infos": null,
		  "previous_proposer": "cosmosvalcons1m46yrx"
		},
		"gov": {
		  "starting_proposal_id": "1",
		  "deposits": null,
		  "votes": null,
		  "proposals": null,
		  "deposit_params": {
			"min_deposit": [
			  {
				"denom": "luna",
				"amount": "10"
			  }
			],
			"max_deposit_period": "86400000000000"
		  },
		  "voting_params": {
			"voting_period": "86400000000000"
		  },
		  "tally_params": {
			"threshold": "0.5000000000",
			"veto": "0.3340000000",
			"governance_penalty": "0.0100000000"
		  }
		},
		"slashing": {
		  "Params": {
			"max-evidence-age": "120000000000",
			"signed-blocks-window": "5000",
			"min-signed-per-window": "0.5000000000",
			"double-sign-unbond-duration": "300000000000",
			"downtime-unbond-duration": "600000000000",
			"slash-fraction-double-sign": "0.200000000",
			"slash-fraction-downtime": "0.0100000000"
		  },
		  "SigningInfos": {},
		  "MissedBlocks": {},
		  "SlashingPeriods": null
		}
	  }`, tx.Addr))
	return
}

// get cmd to initialize all files for tendermint and application
// nolint: errcheck
func InitCmd(ctx *server.Context, cdc *codec.Codec, appInit server.AppInit) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize genesis config, priv-validator file, and p2p-node file",
		Args:  cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {

			config := ctx.Config
			config.SetRoot(viper.GetString(cli.HomeFlag))
			chainID := viper.GetString(client.FlagChainID)
			if chainID == "" {
				chainID = fmt.Sprintf("test-chain-%v", common.RandStr(6))
			}

			nodeKey, err := p2p.LoadOrGenNodeKey(config.NodeKeyFile())
			if err != nil {
				return err
			}
			nodeID := string(nodeKey.ID())

			pk := terraInit.ReadOrCreatePrivValidator(config.PrivValidatorFile())
			genTx, appMessage, validator, err := server.SimpleAppGenTx(cdc, pk)
			if err != nil {
				return err
			}

			appState, err := appInit.AppGenState(
				cdc, tmtypes.GenesisDoc{}, []json.RawMessage{genTx})
			if err != nil {
				return err
			}
			appStateJSON, err := cdc.MarshalJSON(appState)
			if err != nil {
				return err
			}

			toPrint := struct {
				ChainID    string          `json:"chain_id"`
				NodeID     string          `json:"node_id"`
				AppMessage json.RawMessage `json:"app_message"`
			}{
				chainID,
				nodeID,
				appMessage,
			}
			out, err := codec.MarshalJSONIndent(cdc, toPrint)
			if err != nil {
				return err
			}
			fmt.Fprintf(os.Stderr, "%s\n", string(out))
			return terraInit.ExportGenesisFile(config.GenesisFile(), chainID,
				[]tmtypes.GenesisValidator{validator}, appStateJSON)
		},
	}

	cmd.Flags().String(cli.HomeFlag, app.DefaultNodeHome, "node's home directory")
	cmd.Flags().String(flagClientHome, app.DefaultCLIHome, "client's home directory")
	cmd.Flags().String(client.FlagChainID, "",
		"genesis file chain-id, if left blank will be randomly created")
	cmd.Flags().String(client.FlagName, "", "validator's moniker")
	cmd.MarkFlagRequired(client.FlagName)
	return cmd
}

func newApp(logger log.Logger, db dbm.DB, storeTracer io.Writer) abci.Application {
	return app.NewTerraApp(logger, db, storeTracer,
		baseapp.SetPruning(viper.GetString("pruning")),
		baseapp.SetMinimumFees(viper.GetString("minimum_fees")),
	)
}

func exportAppStateAndTMValidators(logger log.Logger, db dbm.DB, storeTracer io.Writer) (
	json.RawMessage, []tmtypes.GenesisValidator, error) {
	bapp := app.NewTerraApp(logger, db, storeTracer)
	return bapp.ExportAppStateAndValidators()
}
