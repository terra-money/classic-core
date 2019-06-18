package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/cosmos-sdk/x/staking/types"
)

// GetCmdQueryValidator implements the validator query command.
func GetCmdQueryValidator(storeName string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validator --validator [validator-addr]",
		Args:  cobra.NoArgs,
		Short: "Query a validator",
		Long: strings.TrimSpace(`Query details about an individual validator:

$ terracli query staking validator --validator terravaloper1gghjut3ccd8ay0zduzj64hwre2fxs9ldmqhffj
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			addrStr := viper.GetString(flagAddressValidator)
			addr, err := sdk.ValAddressFromBech32(addrStr)
			if err != nil {
				return err
			}

			res, err := cliCtx.QueryStore(staking.GetValidatorKey(addr), storeName)
			if err != nil {
				return err
			}

			if len(res) == 0 {
				return fmt.Errorf("No validator found with address %s", addr)
			}

			return cliCtx.PrintOutput(types.MustUnmarshalValidator(cdc, res))
		},
	}

	cmd.Flags().AddFlagSet(fsValidator)

	cmd.MarkFlagRequired(flagAddressValidator)

	return cmd
}

// GetCmdQueryValidators implements the query all validators command.
func GetCmdQueryValidators(storeName string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "validators",
		Short: "Query for all validators",
		Args:  cobra.NoArgs,
		Long: strings.TrimSpace(`Query details about all validators on a network:

$ terracli query staking validators
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			resKVs, err := cliCtx.QuerySubspace(staking.ValidatorsKey, storeName)
			if err != nil {
				return err
			}

			var validators staking.Validators
			for _, kv := range resKVs {
				validators = append(validators, types.MustUnmarshalValidator(cdc, kv.Value))
			}

			return cliCtx.PrintOutput(validators)
		},
	}
}

// GetCmdQueryValidatorUnbondingDelegations implements the query all unbonding delegatations from a validator command.
func GetCmdQueryValidatorUnbondingDelegations(storeKey string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "unbonding-delegations-from --validator [validator-addr]",
		Short: "Query all unbonding delegatations from a validator",
		Long: strings.TrimSpace(`Query delegations that are unbonding _from_ a validator:

$ terracli query staking unbonding-delegations-from --validator terravaloper1gghjut3ccd8ay0zduzj64hwre2fxs9ldmqhffj
`),
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			addrStr := viper.GetString(flagAddressValidator)
			valAddr, err := sdk.ValAddressFromBech32(addrStr)
			if err != nil {
				return err
			}

			bz, err := cdc.MarshalJSON(staking.NewQueryValidatorParams(valAddr))
			if err != nil {
				return err
			}

			route := fmt.Sprintf("custom/%s/%s", storeKey, staking.QueryValidatorUnbondingDelegations)
			res, err := cliCtx.QueryWithData(route, bz)
			if err != nil {
				return err
			}

			var ubds staking.UnbondingDelegations
			cdc.MustUnmarshalJSON(res, &ubds)
			return cliCtx.PrintOutput(ubds)
		},
	}

	cmd.Flags().AddFlagSet(fsValidator)

	cmd.MarkFlagRequired(flagAddressValidator)

	return cmd
}

// GetCmdQueryValidatorRedelegations implements the query all redelegatations from a validator command.
func GetCmdQueryValidatorRedelegations(storeKey string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "redelegations-from --validator [validator-addr]",
		Short: "Query all outgoing redelegatations from a validator",
		Long: strings.TrimSpace(`Query delegations that are redelegating _from_ a validator:

$ terrali query staking redelegations-from --validator terravaloper1gghjut3ccd8ay0zduzj64hwre2fxs9ldmqhffj
`),
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			addrStr := viper.GetString(flagAddressValidator)
			valAddr, err := sdk.ValAddressFromBech32(addrStr)
			if err != nil {
				return err
			}

			bz, err := cdc.MarshalJSON(staking.NewQueryValidatorParams(valAddr))
			if err != nil {
				return err
			}

			route := fmt.Sprintf("custom/%s/%s", storeKey, staking.QueryValidatorRedelegations)
			res, err := cliCtx.QueryWithData(route, bz)
			if err != nil {
				return err
			}

			var reds staking.Redelegations
			cdc.MustUnmarshalJSON(res, &reds)
			return cliCtx.PrintOutput(reds)
		},
	}

	cmd.Flags().AddFlagSet(fsValidator)

	cmd.MarkFlagRequired(flagAddressValidator)

	return cmd
}

// GetCmdQueryDelegation the query delegation command.
func GetCmdQueryDelegation(storeName string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delegation --delegator [delegator-addr] --validator [validator-addr]",
		Short: "Query a delegation based on address and validator address",
		Long: strings.TrimSpace(`Query delegations for an individual delegator on an individual validator:

$ terracli query staking delegation --delegator terra1gghjut3ccd8ay0zduzj64hwre2fxs9ld75ru9p --validator terravaloper1gghjut3ccd8ay0zduzj64hwre2fxs9ldmqhffj
`),
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			valAddrStr := viper.GetString(flagAddressValidator)
			valAddr, err := sdk.ValAddressFromBech32(valAddrStr)
			if err != nil {
				return err
			}

			delAddrStr := viper.GetString(flagAddressDelegator)
			delAddr, err := sdk.AccAddressFromBech32(delAddrStr)
			if err != nil {
				return err
			}

			res, err := cliCtx.QueryStore(staking.GetDelegationKey(delAddr, valAddr), storeName)
			if err != nil {
				return err
			}

			delegation, err := types.UnmarshalDelegation(cdc, res)
			if err != nil {
				return err
			}

			return cliCtx.PrintOutput(delegation)
		},
	}

	cmd.Flags().AddFlagSet(fsValidator)
	cmd.Flags().AddFlagSet(fsDelegator)

	cmd.MarkFlagRequired(flagAddressValidator)
	cmd.MarkFlagRequired(flagAddressDelegator)

	return cmd
}

// GetCmdQueryDelegations implements the command to query all the delegations
// made from one delegator.
func GetCmdQueryDelegations(storeName string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delegations --delegator [delegator-addr]",
		Short: "Query all delegations made by one delegator",
		Long: strings.TrimSpace(`Query delegations for an individual delegator on all validators:

$ terracli query staking delegations --delegator terra1gghjut3ccd8ay0zduzj64hwre2fxs9ld75ru9p
`),
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			delAddrStr := viper.GetString(flagAddressValidator)
			delAddr, err := sdk.AccAddressFromBech32(delAddrStr)
			if err != nil {
				return err
			}

			resKVs, err := cliCtx.QuerySubspace(staking.GetDelegationsKey(delAddr), storeName)
			if err != nil {
				return err
			}

			var delegations staking.Delegations
			for _, kv := range resKVs {
				delegations = append(delegations, types.MustUnmarshalDelegation(cdc, kv.Value))
			}

			return cliCtx.PrintOutput(delegations)
		},
	}

	cmd.Flags().AddFlagSet(fsDelegator)

	cmd.MarkFlagRequired(flagAddressDelegator)

	return cmd
}

// GetCmdQueryValidatorDelegations implements the command to query all the
// delegations to a specific validator.
func GetCmdQueryValidatorDelegations(storeKey string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delegations-to --validator [validator-addr]",
		Short: "Query all delegations made to one validator",
		Long: strings.TrimSpace(`Query delegations on an individual validator:

$ terracli query staking delegations-to --validator terravaloper1gghjut3ccd8ay0zduzj64hwre2fxs9ldmqhffj
`),
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			valAddrStr := viper.GetString(flagAddressValidator)
			validatorAddr, err := sdk.ValAddressFromBech32(valAddrStr)
			if err != nil {
				return err
			}

			bz, err := cdc.MarshalJSON(staking.NewQueryValidatorParams(validatorAddr))
			if err != nil {
				return err
			}

			route := fmt.Sprintf("custom/%s/%s", storeKey, staking.QueryValidatorDelegations)
			res, err := cliCtx.QueryWithData(route, bz)
			if err != nil {
				return err
			}

			var dels staking.Delegations
			cdc.MustUnmarshalJSON(res, &dels)
			return cliCtx.PrintOutput(dels)
		},
	}

	cmd.Flags().AddFlagSet(fsValidator)

	cmd.MarkFlagRequired(flagAddressValidator)

	return cmd
}

// GetCmdQueryUnbondingDelegation implements the command to query a single
// unbonding-delegation record.
func GetCmdQueryUnbondingDelegation(storeName string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "unbonding-delegation --delegator [delegator-addr] --validator [validator-addr]",
		Args:  cobra.NoArgs,
		Short: "Query an unbonding-delegation record based on delegator and validator address",
		Long: strings.TrimSpace(`Query unbonding delegations for an individual delegator on an individual validator:

$ terracli query staking unbonding-delegation --delegator terra1gghjut3ccd8ay0zduzj64hwre2fxs9ld75ru9p --validator terravaloper1gghjut3ccd8ay0zduzj64hwre2fxs9ldmqhffj
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			valAddrStr := viper.GetString(flagAddressValidator)
			valAddr, err := sdk.ValAddressFromBech32(valAddrStr)
			if err != nil {
				return err
			}

			delAddrStr := viper.GetString(flagAddressDelegator)
			delAddr, err := sdk.AccAddressFromBech32(delAddrStr)
			if err != nil {
				return err
			}

			res, err := cliCtx.QueryStore(staking.GetUBDKey(delAddr, valAddr), storeName)
			if err != nil {
				return err
			}

			return cliCtx.PrintOutput(types.MustUnmarshalUBD(cdc, res))

		},
	}

	cmd.Flags().AddFlagSet(fsValidator)
	cmd.Flags().AddFlagSet(fsDelegator)

	cmd.MarkFlagRequired(flagAddressValidator)
	cmd.MarkFlagRequired(flagAddressDelegator)

	return cmd
}

// GetCmdQueryUnbondingDelegations implements the command to query all the
// unbonding-delegation records for a delegator.
func GetCmdQueryUnbondingDelegations(storeName string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "unbonding-delegations --delegator [delegator-addr]",
		Args:  cobra.NoArgs,
		Short: "Query all unbonding-delegations records for one delegator",
		Long: strings.TrimSpace(`Query unbonding delegations for an individual delegator:

$ terracli query staking unbonding-delegation --delegator terra1gghjut3ccd8ay0zduzj64hwre2fxs9ld75ru9p
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			delAddrStr := viper.GetString(flagAddressDelegator)
			delAddr, err := sdk.AccAddressFromBech32(delAddrStr)
			if err != nil {
				return err
			}

			resKVs, err := cliCtx.QuerySubspace(staking.GetUBDsKey(delAddr), storeName)
			if err != nil {
				return err
			}

			var ubds staking.UnbondingDelegations
			for _, kv := range resKVs {
				ubds = append(ubds, types.MustUnmarshalUBD(cdc, kv.Value))
			}

			return cliCtx.PrintOutput(ubds)
		},
	}

	cmd.Flags().AddFlagSet(fsDelegator)

	cmd.MarkFlagRequired(flagAddressDelegator)

	return cmd
}

// GetCmdQueryRedelegation implements the command to query a single
// redelegation record.
func GetCmdQueryRedelegation(storeName string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "redelegation --delegator [delegator-addr] --addr-validator-source [src-validator-addr] --addr-validator-dest [dst-validator-addr]",
		Args:  cobra.NoArgs,
		Short: "Query a redelegation record based on delegator and a source and destination validator address",
		Long: strings.TrimSpace(`Query a redelegation record  for an individual delegator between a source and destination validator:

$ terracli query staking redelegation --delegator terra1gghjut3ccd8ay0zduzj64hwre2fxs9ld75ru9p --addr-validator-source terravaloper1l2rsakp388kuv9k8qzq6lrm9taddae7fpx59wm --addr-validator-dest terravaloper1gghjut3ccd8ay0zduzj64hwre2fxs9ldmqhffj
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			valSrcAddrStr := viper.GetString(flagAddressValidatorSrc)
			valSrcAddr, err := sdk.ValAddressFromBech32(valSrcAddrStr)
			if err != nil {
				return err
			}

			valDstAddrStr := viper.GetString(flagAddressValidatorDst)
			valDstAddr, err := sdk.ValAddressFromBech32(valDstAddrStr)
			if err != nil {
				return err
			}

			delAddrStr := viper.GetString(flagAddressDelegator)
			delAddr, err := sdk.AccAddressFromBech32(delAddrStr)
			if err != nil {
				return err
			}

			res, err := cliCtx.QueryStore(staking.GetREDKey(delAddr, valSrcAddr, valDstAddr), storeName)
			if err != nil {
				return err
			}

			return cliCtx.PrintOutput(types.MustUnmarshalRED(cdc, res))
		},
	}

	cmd.Flags().AddFlagSet(fsDelegator)
	cmd.Flags().AddFlagSet(fsRedelegation)

	cmd.MarkFlagRequired(flagAddressDelegator)
	cmd.MarkFlagRequired(flagAddressValidatorDst)
	cmd.MarkFlagRequired(flagAddressValidatorSrc)

	return cmd
}

// GetCmdQueryRedelegations implements the command to query all the
// redelegation records for a delegator.
func GetCmdQueryRedelegations(storeName string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "redelegations --delegator [delegator-addr]",
		Args:  cobra.NoArgs,
		Short: "Query all redelegations records for one delegator",
		Long: strings.TrimSpace(`Query all redelegation records for an individual delegator:

$ terracli query staking redelegations --delegator terra1gghjut3ccd8ay0zduzj64hwre2fxs9ld75ru9p
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			delAddrStr := viper.GetString(flagAddressDelegator)
			delAddr, err := sdk.AccAddressFromBech32(delAddrStr)
			if err != nil {
				return err
			}

			resKVs, err := cliCtx.QuerySubspace(staking.GetREDsKey(delAddr), storeName)
			if err != nil {
				return err
			}

			var reds staking.Redelegations
			for _, kv := range resKVs {
				reds = append(reds, types.MustUnmarshalRED(cdc, kv.Value))
			}

			return cliCtx.PrintOutput(reds)
		},
	}

	cmd.Flags().AddFlagSet(fsDelegator)

	cmd.MarkFlagRequired(flagAddressDelegator)

	return cmd
}

// GetCmdQueryPool implements the pool query command.
func GetCmdQueryPool(storeName string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "pool",
		Args:  cobra.NoArgs,
		Short: "Query the current staking pool values",
		Long: strings.TrimSpace(`Query values for amounts stored in the staking pool:

$ terracli query staking pool
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, err := cliCtx.QueryStore(staking.PoolKey, storeName)
			if err != nil {
				return err
			}

			return cliCtx.PrintOutput(types.MustUnmarshalPool(cdc, res))
		},
	}
}

// GetCmdQueryPool implements the params query command.
func GetCmdQueryParams(storeName string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "params",
		Args:  cobra.NoArgs,
		Short: "Query the current staking parameters information",
		Long: strings.TrimSpace(`Query values set as staking parameters:

$ terracli query staking params
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			route := fmt.Sprintf("custom/%s/%s", storeName, staking.QueryParameters)
			bz, err := cliCtx.QueryWithData(route, nil)
			if err != nil {
				return err
			}

			var params staking.Params
			cdc.MustUnmarshalJSON(bz, &params)
			return cliCtx.PrintOutput(params)
		},
	}
}
