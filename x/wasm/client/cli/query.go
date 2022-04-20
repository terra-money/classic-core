package cli

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/terra-money/core/x/wasm/types"
)

// GetQueryCmd returns the cli query commands for wasm   module
func GetQueryCmd() *cobra.Command {
	queryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the wasm module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	queryCmd.AddCommand(
		GetCmdQueryByteCode(),
		GetCmdQueryCodeInfo(),
		GetCmdGetContractInfo(),
		GetCmdGetContractStore(),
		GetCmdGetRawStore(),
		GetCmdQueryParams(),
	)
	return queryCmd
}

// GetCmdQueryCodeInfo is for querying code information
func GetCmdQueryCodeInfo() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "code [code-id]",
		Short: "query code information",
		Long:  "query code information",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			codeID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			res, err := queryClient.CodeInfo(context.Background(), &types.QueryCodeInfoRequest{CodeId: codeID})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetCmdQueryByteCode returns the bytecode for a given contract
func GetCmdQueryByteCode() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bytecode [code-id] [output-filename]",
		Short: "Downloads wasm bytecode for given code id",
		Long:  "Downloads wasm bytecode for given code id",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			codeID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			res, err := queryClient.ByteCode(context.Background(), &types.QueryByteCodeRequest{CodeId: codeID})
			if err != nil {
				return err
			}

			fmt.Printf("Downloading wasm code to %s\n", args[1])
			return ioutil.WriteFile(args[1], res.ByteCode, 0600)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetCmdGetContractInfo gets details about a given contract
func GetCmdGetContractInfo() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "contract [contract-address]",
		Short: "Prints out metadata of a contract given its address",
		Long:  "Prints out metadata of a contract given its address",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			addr := args[0]
			_, err = sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			res, err := queryClient.ContractInfo(context.Background(), &types.QueryContractInfoRequest{
				ContractAddress: addr,
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetCmdGetContractStore send query msg to a given contract
func GetCmdGetContractStore() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "contract-store [bech32-address] [msg]",
		Short: "Query contract store of the address with query data and prints the returned result",
		Long:  "Query contract store of the address with query data and prints the returned result",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			addr := args[0]
			_, err = sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			msg := args[1]
			msgBz := []byte(msg)
			if !json.Valid(msgBz) {
				return errors.New("msg must be a json string format")
			}

			res, err := queryClient.ContractStore(context.Background(), &types.QueryContractStoreRequest{
				ContractAddress: addr,
				QueryMsg:        msgBz,
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetCmdGetRawStore dumps full internal state of a given contract
func GetCmdGetRawStore() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "raw-store [bech32-address] [base64-raw-key]",
		Short: "Prints out raw store of a contract",
		Long:  "Prints out raw store of a contract",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			addr := args[0]
			_, err = sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			key := args[1]
			keyBz, err := base64.StdEncoding.DecodeString(key)
			if err != nil {
				return err
			}

			res, err := queryClient.RawStore(context.Background(), &types.QueryRawStoreRequest{
				ContractAddress: addr,
				Key:             keyBz,
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetCmdQueryParams implements the query params command.
func GetCmdQueryParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "params",
		Args:  cobra.NoArgs,
		Short: "Query the current wasm params",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.Params(context.Background(), &types.QueryParamsRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}
