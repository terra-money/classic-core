package cli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	distr "github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/distribution/types"

	"github.com/terra-project/core/x/distribution/client/common"
)

// GetCmdQueryParams implements the query params command.
func GetCmdQueryParams(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "params",
		Args:  cobra.NoArgs,
		Short: "Query distribution params",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			params, err := common.QueryParams(cliCtx, queryRoute)
			if err != nil {
				return err
			}
			return cliCtx.PrintOutput(params)
		},
	}
}

// GetCmdQueryValidatorOutstandingRewards implements the query validator outstanding rewards command.
func GetCmdQueryValidatorOutstandingRewards(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validator-outstanding-rewards --validator [validator-address]",
		Args:  cobra.NoArgs,
		Short: "Query distribution outstanding (un-withdrawn) rewards",
		Long: strings.TrimSpace(`Query distribution outstanding (un-withdrawn) rewards for a validator and all their delegations:

$ terracli query dist validator-outstanding-rewards --validator terravaloper1gghjut3ccd8ay0zduzj64hwre2fxs9ldmqhffj
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			valAddrStr := viper.GetString(flagAddressValidator)
			valAddr, err := sdk.ValAddressFromBech32(valAddrStr)
			if err != nil {
				return err
			}

			bz := cdc.MustMarshalJSON(distr.NewQueryValidatorOutstandingRewardsParams(valAddr))

			route := fmt.Sprintf("custom/%s/validator_outstanding_rewards", queryRoute)
			res, err := cliCtx.QueryWithData(route, bz)
			if err != nil {
				return err
			}

			var outstandingRewards types.ValidatorOutstandingRewards
			cdc.MustUnmarshalJSON(res, &outstandingRewards)
			return cliCtx.PrintOutput(outstandingRewards)
		},
	}

	cmd.Flags().AddFlagSet(fsValidator)

	cmd.MarkFlagRequired(flagAddressValidator)

	return cmd
}

// GetCmdQueryValidatorCommission implements the query validator commission command.
func GetCmdQueryValidatorCommission(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validator-commission --validator [validator]",
		Args:  cobra.NoArgs,
		Short: "Query distribution validator commission",
		Long: strings.TrimSpace(`Query validator commission rewards from delegators to a validator:

$ terracli query distr commission --validator terravaloper1gghjut3ccd8ay0zduzj64hwre2fxs9ldmqhffj
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			valAddrStr := viper.GetString(flagAddressValidator)
			valAddr, err := sdk.ValAddressFromBech32(valAddrStr)
			if err != nil {
				return err
			}

			res, err := common.QueryValidatorCommission(cliCtx, cdc, queryRoute, valAddr)
			if err != nil {
				return err
			}

			var valCom types.ValidatorAccumulatedCommission
			cdc.MustUnmarshalJSON(res, &valCom)
			return cliCtx.PrintOutput(valCom)
		},
	}

	cmd.Flags().AddFlagSet(fsValidator)

	cmd.MarkFlagRequired(flagAddressValidator)

	return cmd
}

// GetCmdQueryValidatorSlashes implements the query validator slashes command.
func GetCmdQueryValidatorSlashes(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validator-slashes --validator [validator] --start [start-height] --end [end-height]",
		Args:  cobra.NoArgs,
		Short: "Query distribution validator slashes",
		Long: strings.TrimSpace(`Query all slashes of a validator for a given block range:

$ terracli query distr slashes --validator terravaloper1gghjut3ccd8ay0zduzj64hwre2fxs9ldmqhffj --start 0  --end 100
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			valAddrStr := viper.GetString(flagAddressValidator)
			valAddr, err := sdk.ValAddressFromBech32(valAddrStr)
			if err != nil {
				return err
			}

			startHeightStr := viper.GetString(flagStartHeight)
			startHeight, err := strconv.ParseUint(startHeightStr, 10, 64)
			if err != nil {
				return fmt.Errorf("start-height %s not a valid uint, please input a valid start-height", startHeightStr)
			}

			endHeightStr := viper.GetString(flagEndHeight)
			endHeight, err := strconv.ParseUint(endHeightStr, 10, 64)
			if err != nil {
				return fmt.Errorf("end-height %s not a valid uint, please input a valid end-height", endHeightStr)
			}

			params := distr.NewQueryValidatorSlashesParams(valAddr, startHeight, endHeight)
			bz, err := cdc.MarshalJSON(params)
			if err != nil {
				return err
			}

			res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/validator_slashes", queryRoute), bz)
			if err != nil {
				return err
			}

			var slashes types.ValidatorSlashEvents
			cdc.MustUnmarshalJSON(res, &slashes)
			return cliCtx.PrintOutput(slashes)
		},
	}

	cmd.Flags().AddFlagSet(fsValidator)
	cmd.Flags().String(flagStartHeight, "", "The start height of given query")
	cmd.Flags().String(flagEndHeight, "", "The end height of given query")

	cmd.MarkFlagRequired(flagAddressValidator)
	cmd.MarkFlagRequired(flagStartHeight)
	cmd.MarkFlagRequired(flagEndHeight)

	return cmd
}

// GetCmdQueryDelegatorRewards implements the query delegator rewards command.
func GetCmdQueryDelegatorRewards(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rewards --delegator [delegator-addr] -validator [validator-addr]",
		Args:  cobra.NoArgs,
		Short: "Query all distribution delegator rewards or rewards from a particular validator",
		Long: strings.TrimSpace(`Query all rewards earned by a delegator, optionally restrict to rewards from a single validator:

$ terracli query distr rewards --delegator terra1gghjut3ccd8ay0zduzj64hwre2fxs9ld75ru9p
$ terracli query distr rewards --delegator terra1gghjut3ccd8ay0zduzj64hwre2fxs9ld75ru9p --validator terravaloper1gghjut3ccd8ay0zduzj64hwre2fxs9ldmqhffj
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			var resp []byte
			var err error

			delAddrStr := viper.GetString(flagAddressDelegator)
			valAddrStr := viper.GetString(flagAddressValidator)

			if len(valAddrStr) == 0 {
				resp, err = common.QueryDelegatorTotalRewards(cliCtx, cdc, queryRoute, delAddrStr)
			} else {
				resp, err = common.QueryDelegationRewards(cliCtx, cdc, queryRoute, delAddrStr, valAddrStr)
			}

			if err != nil {
				return err
			}

			var result sdk.DecCoins
			cdc.MustUnmarshalJSON(resp, &result)
			return cliCtx.PrintOutput(result)
		},
	}

	cmd.Flags().AddFlagSet(fsValidator)
	cmd.Flags().AddFlagSet(fsDelegator)

	cmd.MarkFlagRequired(flagAddressDelegator)

	return cmd
}
