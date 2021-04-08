package cli

import (
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authclient "github.com/cosmos/cosmos-sdk/x/auth/client"

	feeutils "github.com/terra-project/core/custom/auth/client/utils"
	"github.com/terra-project/core/x/msgauth/types"
)

// FlagPeriod is flag to specify grant period
const FlagPeriod = "period"

// GetTxCmd returns the transaction commands for this module
func GetTxCmd() *cobra.Command {
	AuthorizationTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Msg authorization transactions subcommands",
		Long:                       "Authorize and revoke access to execute transactions on behalf of your address",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	AuthorizationTxCmd.AddCommand(
		GetCmdGrantAuthorization(),
		GetCmdRevokeAuthorization(),
		GetCmdSendAs(),
	)

	return AuthorizationTxCmd
}

func GetCmdGrantAuthorization() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "grant [grantee-address] [msg-type] [limit]",
		Short: "Grant authorization of a specific msg type to an address",
		Long: strings.TrimSpace(`
Grant authorization of a specific msg type to an address 
to let the address execute a transaction on your behalf,

$ terracli tx msgauth grant terra... send 1000000uluna,10000000ukrw --from [granter]

Or, you can just give authorization of other msg types

$ terracli tx msgauth grant terra... swap --from [granter]
				`),
		Args: cobra.RangeArgs(2, 3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			granter := clientCtx.FromAddress
			grantee, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			msgType := args[1]

			var authorization types.AuthorizationI
			if msgType == (types.SendAuthorization{}.MsgType()) {
				limit, err := sdk.ParseCoinsNormalized(args[2])
				if err != nil {
					return err
				}

				authorization = types.NewSendAuthorization(limit)
			} else {
				authorization = types.NewGenericAuthorization(msgType)
			}

			period := time.Duration(viper.GetInt64(FlagPeriod)) * time.Second

			msg, err := types.NewMsgGrantAuthorization(granter, grantee, authorization, period)
			if err != nil {
				return err
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)

		},
	}

	cmd.Flags().Int64(FlagPeriod, int64(3600*24*365), "The second unit of time duration which the authorization is active for the user; Default is a year")
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func GetCmdRevokeAuthorization() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "revoke [grantee_address] [msg_type]",
		Short: "Revoke authorization",
		Long: strings.TrimSpace(`
Revoke authorization from an address for a msg type,

$ terracli msgauth revoke terra... send --from [granter]
`),
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			granter := clientCtx.FromAddress
			grantee, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			msgAuthorized := args[1]

			msg := types.NewMsgRevokeAuthorization(granter, grantee, msgAuthorized)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func GetCmdSendAs() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "send-as [granter] [tx_json] --from [grantee]",
		Short: "Execute tx on behalf of granter account",
		Long: strings.TrimSpace(`
Execute tx on behalf of granter account,

$ terracli msgauth send-as terra... ./tx.json --from [grantee]

tx.json should be format of StdTx
`),
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			// Generate transaction factory for gas simulation
			txf := tx.NewFactoryCLI(clientCtx, cmd.Flags())

			grantee := clientCtx.FromAddress

			stdTx, err := authclient.ReadTxFromFile(clientCtx, args[1])
			if err != nil {
				return err
			}

			msg, err := types.NewMsgExecAuthorized(grantee, stdTx.GetMsgs())
			if err != nil {
				return err
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			if !clientCtx.GenerateOnly && txf.Fees().IsZero() {
				// estimate tax and gas
				stdFee, err := feeutils.ComputeFeesWithCmd(clientCtx, cmd.Flags(), msg)

				if err != nil {
					return err
				}

				// override gas and fees
				txf = txf.
					WithFees(stdFee.Amount.String()).
					WithGas(stdFee.Gas).
					WithSimulateAndExecute(false).
					WithGasPrices("")
			}

			// build and sign the transaction, then broadcast to Tendermint
			return tx.GenerateOrBroadcastTxWithFactory(clientCtx, txf, msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}
