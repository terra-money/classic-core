package cli

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/terra-money/core/x/wasm/client/utils"
	"github.com/terra-money/core/x/wasm/internal/types"
)

const flagRaw = "raw"

// GetQueryCmd returns the cli query commands for wasm   module
func GetQueryCmd(queryRoute string, cdc *codec.Codec) *cobra.Command {
	queryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the wasm module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	queryCmd.AddCommand(flags.GetCommands(
		GetCmdQueryByteCode(queryRoute, cdc),
		GetCmdQueryCodeInfo(queryRoute, cdc),
		GetCmdGetContractInfo(queryRoute, cdc),
		GetCmdGetContractStore(queryRoute, cdc),
		GetCmdGetRawStore(queryRoute, cdc),
		GetCmdQueryParams(queryRoute, cdc),
	)...)
	return queryCmd
}

// GetCmdQueryCodeInfo is for querying code information
func GetCmdQueryCodeInfo(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "code [code-id]",
		Short: "query code information",
		Long:  "query code information",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			codeID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			params := types.NewQueryCodeIDParams(codeID)
			bz, err := cliCtx.Codec.MarshalJSON(params)
			if err != nil {
				return err
			}

			route := fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryGetCodeInfo)
			res, _, err := cliCtx.QueryWithData(route, bz)
			if err != nil {
				return err
			}

			var codeInfo types.CodeInfo
			cdc.MustUnmarshalJSON(res, &codeInfo)
			return cliCtx.PrintOutput(codeInfo)
		},
	}
}

// GetCmdQueryByteCode returns the bytecode for a given contract
func GetCmdQueryByteCode(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "bytecode [code-id] [output-filename]",
		Short: "Downloads wasm bytecode for given code id",
		Long:  "Downloads wasm bytecode for given code id",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			codeID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			params := types.NewQueryCodeIDParams(codeID)
			bz, err := cliCtx.Codec.MarshalJSON(params)
			if err != nil {
				return err
			}

			route := fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryGetByteCode)
			res, _, err := cliCtx.QueryWithData(route, bz)
			if err != nil {
				return err
			}

			if len(res) == 0 {
				return fmt.Errorf("contract not found")
			}
			var bytecode []byte
			err = json.Unmarshal(res, &bytecode)
			if err != nil {
				return err
			}

			if len(bytecode) == 0 {
				return fmt.Errorf("contract not found")
			}

			fmt.Printf("Downloading wasm code to %s\n", args[1])
			return ioutil.WriteFile(args[1], bytecode, 0600)
		},
	}
}

// GetCmdGetContractInfo gets details about a given contract
func GetCmdGetContractInfo(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "contract [contract-address]",
		Short: "Prints out metadata of a contract given its address",
		Long:  "Prints out metadata of a contract given its address",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			addr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			params := types.NewQueryContractAddressParams(addr)
			bz, err := cliCtx.Codec.MarshalJSON(params)
			if err != nil {
				return err
			}

			route := fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryGetContractInfo)
			res, _, err := cliCtx.QueryWithData(route, bz)
			if err != nil {
				return err
			}

			var contractInfo types.ContractInfo
			cdc.MustUnmarshalJSON(res, &contractInfo)
			return cliCtx.PrintOutput(contractInfo)
		},
	}
}

// GetCmdGetContractStore send query msg to a given contract
func GetCmdGetContractStore(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "contract-store [bech32-address] [msg]",
		Short: "Query contract store of the address with query data and prints the returned result",
		Long:  "Query contract store of the address with query data and prints the returned result",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			addr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			msg := args[1]
			msgBz := []byte(msg)
			if !json.Valid(msgBz) {
				return errors.New("msg must be a json string format")
			}

			params := types.NewQueryContractParams(addr, msgBz)
			bz, err := cliCtx.Codec.MarshalJSON(params)
			if err != nil {
				return err
			}

			route := fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryContractStore)
			res, _, err := cliCtx.QueryWithData(route, bz)
			if err != nil {
				return err
			}

			fmt.Println(string(res))
			return nil
		},
	}
}

// GetCmdGetRawStore dumps full internal state of a given contract
func GetCmdGetRawStore(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "raw-store [bech32-address] [key] [subkey]",
		Short: "Prints out raw store of a contract",
		Long:  "Prints out raw store of a contract",
		Args:  cobra.RangeArgs(2, 3),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			addr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			// need to extend key with prefix of its length
			key := args[1]
			subkey := ""
			if len(args) == 3 {
				subkey = args[2]
			}

			keyBz := append(utils.EncodeKey(key), []byte(subkey)...)
			params := types.NewQueryRawStoreParams(addr, keyBz)
			bz, err := cliCtx.Codec.MarshalJSON(params)
			if err != nil {
				return err
			}

			route := fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryRawStore)
			res, _, err := cliCtx.QueryWithData(route, bz)
			if err != nil {
				return err
			}

			model := types.Model{
				Key:   keyBz,
				Value: res,
			}

			return cliCtx.PrintOutput(model)
		},
	}
}

// GetCmdQueryParams implements the query params command.
func GetCmdQueryParams(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "params",
		Args:  cobra.NoArgs,
		Short: "Query the current wasm params",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryParameters), nil)
			if err != nil {
				return err
			}

			var params types.Params
			cdc.MustUnmarshalJSON(res, &params)
			return cliCtx.PrintOutput(params)
		},
	}

	return cmd
}
