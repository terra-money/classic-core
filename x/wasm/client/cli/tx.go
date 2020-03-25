package cli

import (
	"fmt"
	"io/ioutil"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"

	wasmUtils "github.com/terra-project/core/x/wasm/client/utils"
	"github.com/terra-project/core/x/wasm/internal/types"
)

const (
	flagTo     = "to"
	flagAmount = "amount"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd(cdc *codec.Codec) *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Wasm transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	txCmd.AddCommand(client.PostCommands(
		StoreCodeCmd(cdc),
		InstantiateContractCmd(cdc),
		ExecuteContractCmd(cdc),
	)...)
	return txCmd
}

// StoreCodeCmd will upload code to be reused.
func StoreCodeCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "store [from_key_or_address] [wasm-file]",
		Short: "Upload a wasm binary",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContextWithFrom(args[0]).WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			// parse coins trying to be sent
			wasm, err := ioutil.ReadFile(args[1])
			if err != nil {
				return err
			}

			// limit the input size
			if len(wasm) > types.MaxWasmSize {
				return fmt.Errorf("input size exceeds the max size hard-cap (allowed:%d, actual: %d)",
					types.MaxWasmSize, len(wasm))
			}

			// gzip the wasm file
			if wasmUtils.IsWasm(wasm) {
				wasm, err = wasmUtils.GzipIt(wasm)

				if err != nil {
					return err
				}
			} else if !wasmUtils.IsGzip(wasm) {
				return fmt.Errorf("invalid input file. Use wasm binary or gzip")
			}

			// build and sign the transaction, then broadcast to Tendermint
			msg := types.MsgStoreCode{
				Sender:       cliCtx.GetFromAddress(),
				WASMByteCode: wasm,
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
	return cmd
}

// InstantiateContractCmd will instantiate a contract from previously uploaded code.
func InstantiateContractCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "instantiate [from_key_or_address] [code-id-int64] [json-encoded-args] [coins]",
		Short: "Instantiate a wasm contract",
		Args:  cobra.RangeArgs(3, 4),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithFrom(args[0]).WithCodec(cdc)

			// get the id of the code to instantiate
			codeID, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}

			initMsg := args[2]

			var coins sdk.Coins
			if len(args) == 4 {
				coins, err = sdk.ParseCoins(args[3])
				if err != nil {
					return err
				}
			}

			// build and sign the transaction, then broadcast to Tendermint
			msg := types.MsgInstantiateContract{
				Sender:    cliCtx.GetFromAddress(),
				CodeID:    codeID,
				InitCoins: coins,
				InitMsg:   []byte(initMsg),
			}
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	return cmd
}

// ExecuteContractCmd will instantiate a contract from previously uploaded code.
func ExecuteContractCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "execute [from_key_or_address] [contract_addr_bech32] [json_encoded_args] [coins]",
		Short: "Execute a command on a wasm contract",
		Args:  cobra.RangeArgs(3, 4),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithFrom(args[0]).WithCodec(cdc)

			// get the id of the code to instantiate
			contractAddr, err := sdk.AccAddressFromBech32(args[1])
			if err != nil {
				return err
			}

			execMsg := args[2]

			var coins sdk.Coins
			if len(args) == 4 {
				coins, err = sdk.ParseCoins(args[3])
				if err != nil {
					return err
				}
			}

			// build and sign the transaction, then broadcast to Tendermint
			msg := types.MsgExecuteContract{
				Sender:   cliCtx.GetFromAddress(),
				Contract: contractAddr,
				Coins:    coins,
				Msg:      []byte(execMsg),
			}
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	return cmd
}
