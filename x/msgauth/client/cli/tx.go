package cli

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authclient "github.com/cosmos/cosmos-sdk/x/auth/client/utils"

	feeutils "github.com/terra-money/core/x/auth/client/utils"
	"github.com/terra-money/core/x/msgauth/internal/types"
)

// FlagPeriod is flag to specify grant period
const FlagPeriod = "period"

// GetTxCmd returns the transaction commands for this module
func GetTxCmd(storeKey string, cdc *codec.Codec) *cobra.Command {
	AuthorizationTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Msg authorization transactions subcommands",
		Long:                       "Authorize and revoke access to execute transactions on behalf of your address",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	AuthorizationTxCmd.AddCommand(flags.PostCommands(
		GetCmdGrantAuthorization(cdc),
		GetCmdRevokeAuthorization(cdc),
		GetCmdSendAs(cdc),
	)...)

	return AuthorizationTxCmd
}

func GetCmdGrantAuthorization(cdc *codec.Codec) *cobra.Command {
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
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(authclient.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)

			granter := cliCtx.FromAddress
			grantee, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			msgType := args[1]

			var authorization types.Authorization
			if msgType == (types.SendAuthorization{}.MsgType()) {
				limit, err := sdk.ParseCoins(args[2])
				if err != nil {
					return err
				}

				authorization = types.NewSendAuthorization(limit)
			} else {
				authorization = types.NewGenericAuthorization(msgType)
			}

			period := time.Duration(viper.GetInt64(FlagPeriod)) * time.Second

			msg := types.NewMsgGrantAuthorization(granter, grantee, authorization, period)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return authclient.CompleteAndBroadcastTxCLI(txBldr, cliCtx, []sdk.Msg{msg})

		},
	}

	cmd.Flags().Int64(FlagPeriod, int64(3600*24*365), "The second unit of time duration which the authorization is active for the user; Default is a year")

	return cmd
}

func GetCmdRevokeAuthorization(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "revoke [grantee_address] [msg_type]",
		Short: "Revoke authorization",
		Long: strings.TrimSpace(`
Revoke authorization from an address for a msg type,

$ terracli msgauth revoke terra... send --from [granter]
`),
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(authclient.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)

			granter := cliCtx.FromAddress
			grantee, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			msgAuthorized := args[1]

			msg := types.NewMsgRevokeAuthorization(granter, grantee, msgAuthorized)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return authclient.CompleteAndBroadcastTxCLI(txBldr, cliCtx, []sdk.Msg{msg})
		},
	}
	return cmd
}

func GetCmdSendAs(cdc *codec.Codec) *cobra.Command {
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
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(authclient.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)

			grantee := cliCtx.FromAddress

			var stdTx auth.StdTx
			bz, err := ioutil.ReadFile(args[1])
			if err != nil {
				return err
			}

			err = cdc.UnmarshalJSON(bz, &stdTx)
			if err != nil {
				return err
			}

			msg := types.NewMsgExecAuthorized(grantee, stdTx.Msgs)

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

			return authclient.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
	return cmd
}
