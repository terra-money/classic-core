package cli

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"

	feeutils "github.com/terra-money/core/x/auth/client/utils"
	wasmUtils "github.com/terra-money/core/x/wasm/client/utils"
	"github.com/terra-money/core/x/wasm/internal/types"
)

const (
	flagTo         = "to"
	flagAmount     = "amount"
	flagMigratable = "migratable"
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
		MigrateContractCmd(cdc),
		UpdateContractOwnerCmd(cdc),
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
			wasmBytes, err := ioutil.ReadFile(args[0])
			if err != nil {
				return err
			}

			// limit the input size
			if wasmLen := uint64(len(wasmBytes)); wasmLen > types.EnforcedMaxContractSize {
				return fmt.Errorf("wasm code size exceeds the max size hard-cap (allowed:%d, actual: %d)",
					types.EnforcedMaxContractSize, wasmLen)
			}

			// gzip the wasm file
			if wasmUtils.IsWasm(wasmBytes) {
				wasmBytes, err = wasmUtils.GzipIt(wasmBytes)

				if err != nil {
					return err
				}
			} else if !wasmUtils.IsGzip(wasmBytes) {
				return fmt.Errorf("invalid input file. Use wasm binary or gzip")
			}

			// build and sign the transaction, then broadcast to Tendermint
			msg := types.NewMsgStoreCode(fromAddr, wasmBytes)
			if err := msg.ValidateBasic(); err != nil {
				return err
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

			migratable := viper.GetBool(flagMigratable)

			// get the id of the code to instantiate
			codeID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			initMsgBz := []byte(args[1])
			if !json.Valid(initMsgBz) {
				return errors.New("msg must be a json string format")
			}

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
			msg := types.NewMsgInstantiateContract(fromAddr, codeID, initMsgBz, coins, migratable)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			if !cliCtx.GenerateOnly && txBldr.Fees().IsZero() {
				// extimate tax and gas
				fees, gas, err := feeutils.ComputeFees(cliCtx, feeutils.ComputeReqParams{
					Memo:          txBldr.Memo(),
					ChainID:       txBldr.ChainID(),
					AccountNumber: txBldr.AccountNumber(),
					Sequence:      txBldr.Sequence(),
					GasPrices:     txBldr.GasPrices(),
					Gas:           fmt.Sprintf("%d", txBldr.Gas()),
					GasAdjustment: fmt.Sprintf("%f", txBldr.GasAdjustment()),
					Msgs:          []sdk.Msg{msg},
				})

				if err != nil {
					return err
				}

				// override gas and fees
				txBldr = auth.NewTxBuilder(txBldr.TxEncoder(), txBldr.AccountNumber(), txBldr.Sequence(),
					gas, txBldr.GasAdjustment(), false, txBldr.ChainID(), txBldr.Memo(), fees, sdk.DecCoins{})
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd.Flags().Bool(flagMigratable, false, "setting the flag will make the contract migratable")
	return cmd
}

// ExecuteContractCmd will instantiate a contract from previously uploaded code.
func ExecuteContractCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "execute [contract-addr-bech32] [json-encoded-args] [coins]",
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
			if !json.Valid(execMsgBz) {
				return errors.New("msg must be a json string format")
			}

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
			msg := types.NewMsgExecuteContract(fromAddr, contractAddr, execMsgBz, coins)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			if !cliCtx.GenerateOnly && txBldr.Fees().IsZero() {
				// extimate tax and gas
				fees, gas, err := feeutils.ComputeFees(cliCtx, feeutils.ComputeReqParams{
					Memo:          txBldr.Memo(),
					ChainID:       txBldr.ChainID(),
					AccountNumber: txBldr.AccountNumber(),
					Sequence:      txBldr.Sequence(),
					GasPrices:     txBldr.GasPrices(),
					Gas:           fmt.Sprintf("%d", txBldr.Gas()),
					GasAdjustment: fmt.Sprintf("%f", txBldr.GasAdjustment()),
					Msgs:          []sdk.Msg{msg},
				})

				if err != nil {
					return err
				}

				// override gas and fees
				txBldr = auth.NewTxBuilder(txBldr.TxEncoder(), txBldr.AccountNumber(), txBldr.Sequence(),
					gas, txBldr.GasAdjustment(), false, txBldr.ChainID(), txBldr.Memo(), fees, sdk.DecCoins{})
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	return cmd
}

// MigrateContractCmd will instantiate a contract from previously uploaded code.
func MigrateContractCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "migrate [contract-addr-bech32] [new-code-id] [json-encoded-args]",
		Short: "Migrate a contract to new code base",
		Long: strings.TrimSpace(`
Migrate a contract to new code by calling migrate function of 
the new code.

$ terracli tx wasm migrate terra... 10 '{"verifier": "terra..."}'
		`),
		Args: cobra.ExactArgs(3),
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

			newCodeID, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}

			migrateMsgBz := []byte(args[2])

			// limit the input size
			if migrateMsgLen := uint64(len(migrateMsgBz)); migrateMsgLen > types.EnforcedMaxContractMsgSize {
				return fmt.Errorf("migrate msg size exceeds the max size hard-cap (allowed:%d, actual: %d)",
					types.EnforcedMaxContractMsgSize, migrateMsgLen)
			}

			// build and sign the transaction, then broadcast to Tendermint
			msg := types.NewMsgMigrateContract(fromAddr, contractAddr, newCodeID, migrateMsgBz)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	return cmd
}

// UpdateContractOwnerCmd will instantiate a contract from previously uploaded code.
func UpdateContractOwnerCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-owner [contract-addr-bech32] [new-owner]",
		Short: "update a contract owner",
		Long: strings.TrimSpace(`
Update a contract owner to a new address

$ terracli tx wasm update-owner terra... terra...
		`),
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

			fromAddr := cliCtx.GetFromAddress()
			if fromAddr.Empty() {
				return fmt.Errorf("must specify flag --from")
			}

			contractAddr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			newOwnerAddr, err := sdk.AccAddressFromBech32(args[1])
			if err != nil {
				return err
			}

			// build and sign the transaction, then broadcast to Tendermint
			msg := types.NewMsgUpdateContractOwner(fromAddr, newOwnerAddr, contractAddr)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	return cmd
}
