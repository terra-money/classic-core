package cli

import (
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	amino "github.com/tendermint/go-amino"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/utils"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtxb "github.com/cosmos/cosmos-sdk/x/auth/client/txbuilder"
	"github.com/cosmos/cosmos-sdk/x/distribution/types"

	"github.com/terra-project/core/x/distribution/client/common"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd(storeKey string, cdc *amino.Codec) *cobra.Command {
	distTxCmd := &cobra.Command{
		Use:   "dist",
		Args:  cobra.NoArgs,
		Short: "Distribution transactions subcommands",
	}

	distTxCmd.AddCommand(client.PostCommands(
		GetCmdWithdrawRewards(cdc),
		GetCmdSetWithdrawAddr(cdc),
	)...)

	return distTxCmd
}

// command to withdraw rewards
func GetCmdWithdrawRewards(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "withdraw-rewards --validator [validator-addr]",
		Args:  cobra.NoArgs,
		Short: "witdraw rewards from a given delegation address, and optionally withdraw validator commission if the delegation address given is a validator operator",
		Long: strings.TrimSpace(`witdraw rewards from a given delegation address, and optionally withdraw validator commission if the delegation address given is a validator operator:

$ terracli tx distr withdraw-rewards --validator terravaloper1gghjut3ccd8ay0zduzj64hwre2fxs9ldmqhffj --from mykey
$ terracli tx distr withdraw-rewards --validator terravaloper1gghjut3ccd8ay0zduzj64hwre2fxs9ldmqhffj --from mykey --commission
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := authtxb.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().
				WithCodec(cdc).
				WithAccountDecoder(cdc)

			delAddr := cliCtx.GetFromAddress()

			valAddrStr := viper.GetString(flagAddressValidator)
			valAddr, err := sdk.ValAddressFromBech32(valAddrStr)
			if err != nil {
				return err
			}

			msgs := []sdk.Msg{types.NewMsgWithdrawDelegatorReward(delAddr, valAddr)}
			if viper.GetBool(flagComission) {
				msgs = append(msgs, types.NewMsgWithdrawValidatorCommission(valAddr))
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, msgs, false)
		},
	}

	cmd.Flags().AddFlagSet(fsValidator)
	cmd.Flags().Bool(flagComission, false, "also withdraw validator's commission")

	cmd.MarkFlagRequired(flagAddressValidator)
	return cmd
}

// command to withdraw all rewards
func GetCmdWithdrawAllRewards(cdc *codec.Codec, queryRoute string) *cobra.Command {
	return &cobra.Command{
		Use:   "withdraw-all-rewards",
		Args:  cobra.NoArgs,
		Short: "withdraw all delegations rewards for a delegator",
		Long: strings.TrimSpace(`Withdraw all rewards for a single delegator:

$ terracli tx distr withdraw-all-rewards --from mykey
`),
		RunE: func(cmd *cobra.Command, args []string) error {

			txBldr := authtxb.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().
				WithCodec(cdc).
				WithAccountDecoder(cdc)

			delAddr := cliCtx.GetFromAddress()
			msgs, err := common.WithdrawAllDelegatorRewards(cliCtx, cdc, queryRoute, delAddr)
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, msgs, false)
		},
	}
}

// command to replace a delegator's withdrawal address
func GetCmdSetWithdrawAddr(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-withdraw-addr --withdraw-to [withdraw-addr]",
		Args:  cobra.NoArgs,
		Short: "change the default withdraw address for rewards associated with an address",
		Long: strings.TrimSpace(`Set the withdraw address for rewards associated with a delegator address:

$ terracli tx set-withdraw-addr --withdraw-to terra1gghjut3ccd8ay0zduzj64hwre2fxs9ld75ru9p --from mykey
`),
		RunE: func(cmd *cobra.Command, args []string) error {

			txBldr := authtxb.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().
				WithCodec(cdc).
				WithAccountDecoder(cdc)

			delAddr := cliCtx.GetFromAddress()

			withdrawAddrToStr := viper.GetString(flagWithdrawTo)
			withdrawAddrTo, err := sdk.AccAddressFromBech32(withdrawAddrToStr)
			if err != nil {
				return err
			}

			msg := types.NewMsgSetWithdrawAddress(delAddr, withdrawAddrTo)
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg}, false)
		},
	}

	cmd.Flags().String(flagWithdrawTo, "", "Target address to withdraw")

	cmd.MarkFlagRequired(flagWithdrawTo)

	return cmd
}
