package cli

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
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
	txCmd.AddCommand(flags.PostCommands(
		StoreCodeCmd(cdc),
		InstantiateContractCmd(cdc),
		ExecuteContractCmd(cdc),
	)...)
	return txCmd
}

// StoreCodeCmd will upload code to be reused.
func StoreCodeCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "store [wasm-file]",
		Short: "Upload a wasm binary",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

			fromAddr := cliCtx.GetFromAddress()
			if fromAddr.Empty() {
				return fmt.Errorf("must specify flag --from")
			}

			// parse coins trying to be sent
			wasm, err := ioutil.ReadFile(args[0])
			if err != nil {
				return err
			}

			// limit the input size
			if wasmLen := uint64(len(wasm)); wasmLen > types.EnforcedMaxContractSize {
				return fmt.Errorf("wasm code size exceeds the max size hard-cap (allowed:%d, actual: %d)",
					types.EnforcedMaxContractSize, wasmLen)
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
				Sender:       fromAddr,
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
		Use:   "instantiate [code-id-int64] [json-encoded-args] [coins]",
		Short: "Instantiate a wasm contract",
		Long: `
Instantiate a wasm contract of the code which has the given id

$ terracli instantiate 1 '{"arbiter": "terra~~"}'

You can also instantiate it with funds

$ terracli instantiate 1 '{"arbiter": "terra~~"}' "1000000uluna"
`,
		Args: cobra.RangeArgs(2, 3),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

			fromAddr := cliCtx.GetFromAddress()
			if fromAddr.Empty() {
				return fmt.Errorf("must specify flag --from")
			}

			// get the id of the code to instantiate
			codeID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			initMsgBz := []byte(args[1])

			// limit the input size
			if initMsgLen := uint64(len(initMsgBz)); initMsgLen > types.EnforcedMaxContractMsgSize {
				return fmt.Errorf("init msg size exceeds the max size hard-cap (allowed:%d, actual: %d)",
					types.EnforcedMaxContractMsgSize, initMsgLen)
			}

			var coins sdk.Coins
			if len(args) == 3 {
				coins, err = sdk.ParseCoins(args[2])
				if err != nil {
					return err
				}
			}

			// build and sign the transaction, then broadcast to Tendermint
			msg := types.MsgInstantiateContract{
				Sender:    fromAddr,
				CodeID:    codeID,
				InitCoins: coins,
				InitMsg:   initMsgBz,
			}
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	return cmd
}

// ExecuteContractCmd will instantiate a contract from previously uploaded code.
func ExecuteContractCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "execute [contract_addr_bech32] [json_encoded_args] [coins]",
		Short: "Execute a command on a wasm contract",
		Args:  cobra.RangeArgs(2, 3),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

			fromAddr := cliCtx.GetFromAddress()
			if fromAddr.Empty() {
				return fmt.Errorf("must specify flag --from")
			}

			// get the id of the code to instantiate
			contractAddr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			execMsgBz := []byte(args[1])

			// limit the input size
			if execMsgLen := uint64(len(execMsgBz)); execMsgLen > types.EnforcedMaxContractMsgSize {
				return fmt.Errorf("exec msg size exceeds the max size hard-cap (allowed:%d, actual: %d)",
					types.EnforcedMaxContractMsgSize, execMsgLen)
			}

			var coins sdk.Coins
			if len(args) == 3 {
				coins, err = sdk.ParseCoins(args[2])
				if err != nil {
					return err
				}
			}

			// build and sign the transaction, then broadcast to Tendermint
			msg := types.MsgExecuteContract{
				Sender:   fromAddr,
				Contract: contractAddr,
				Coins:    coins,
				Msg:      execMsgBz,
			}
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	return cmd
}
