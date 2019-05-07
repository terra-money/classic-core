package cli

import (
	"fmt"
	"strings"

	"github.com/cosmos/cosmos-sdk/x/auth"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/utils"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtxb "github.com/cosmos/cosmos-sdk/x/auth/client/txbuilder"
	"github.com/cosmos/cosmos-sdk/x/staking"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// GetCmdCreateValidator implements the create validator command handler.
// TODO: Add full description
func GetCmdCreateValidator(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-validator",
		Args:  cobra.NoArgs,
		Short: "create new validator initialized with a self-delegation to it",
		Long: `
		terracli tx staking create-validator \
  --amount=5000000uluna \
  --pubkey=$(terrad tendermint show-validator) \
  --moniker="choose a moniker" \
  --chain-id=<chain_id> \
  --from=<key_name> \
  --commission-rate="0.10" \
  --commission-max-rate="0.20" \
  --commission-max-change-rate="0.01" \
  --min-self-delegation="1"`,
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := authtxb.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().
				WithCodec(cdc).
				WithAccountDecoder(cdc)

			txBldr, msg, err := BuildCreateValidatorMsg(cliCtx, txBldr)
			if err != nil {
				return err
			}

			offline := viper.GetBool(flagOffline)

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg}, offline)
		},
	}

	cmd.Flags().AddFlagSet(fsPk)
	cmd.Flags().AddFlagSet(fsAmount)
	cmd.Flags().AddFlagSet(fsDescriptionCreate)
	cmd.Flags().AddFlagSet(fsCommissionCreate)
	cmd.Flags().AddFlagSet(fsMinSelfDelegation)

	cmd.Flags().String(flagIP, "", fmt.Sprintf("The node's public IP. It takes effect only when used in combination with --%s", client.FlagGenerateOnly))
	cmd.Flags().String(flagNodeID, "", "The node's ID")
	cmd.Flags().Bool(flagOffline, false, " Offline mode; Do not query a full node")

	cmd.MarkFlagRequired(client.FlagFrom)
	cmd.MarkFlagRequired(flagAmount)
	cmd.MarkFlagRequired(flagPubKey)
	cmd.MarkFlagRequired(flagMoniker)

	return cmd
}

// GetCmdEditValidator implements the create edit validator command.
// TODO: add full description
func GetCmdEditValidator(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "edit-validator",
		Args:  cobra.NoArgs,
		Short: "edit an existing validator account",
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := authtxb.NewTxBuilderFromCLI().WithTxEncoder(auth.DefaultTxEncoder(cdc))
			cliCtx := context.NewCLIContext().
				WithCodec(cdc).
				WithAccountDecoder(cdc)

			valAddr := cliCtx.GetFromAddress()
			description := staking.Description{
				Moniker:  viper.GetString(flagMoniker),
				Identity: viper.GetString(flagIdentity),
				Website:  viper.GetString(flagWebsite),
				Details:  viper.GetString(flagDetails),
			}

			var newRate *sdk.Dec

			commissionRate := viper.GetString(flagCommissionRate)
			if commissionRate != "" {
				rate, err := sdk.NewDecFromStr(commissionRate)
				if err != nil {
					return fmt.Errorf("invalid new commission rate: %v", err)
				}

				newRate = &rate
			}

			var newMinSelfDelegation *sdk.Int

			minSelfDelegationString := viper.GetString(flagMinSelfDelegation)
			if minSelfDelegationString != "" {
				msb, ok := sdk.NewIntFromString(minSelfDelegationString)
				if !ok {
					return fmt.Errorf(staking.ErrMinSelfDelegationInvalid(staking.DefaultCodespace).Error())
				}
				newMinSelfDelegation = &msb
			}

			msg := staking.NewMsgEditValidator(sdk.ValAddress(valAddr), description, newRate, newMinSelfDelegation)

			offline := viper.GetBool(flagOffline)

			// build and sign the transaction, then broadcast to Tendermint
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg}, offline)
		},
	}

	cmd.Flags().AddFlagSet(fsDescriptionEdit)
	cmd.Flags().AddFlagSet(fsCommissionUpdate)
	cmd.Flags().AddFlagSet(fsMinSelfDelegation)

	cmd.Flags().Bool(flagOffline, false, " Offline mode; Do not query a full node")

	return cmd
}

// GetCmdDelegate implements the delegate command.
func GetCmdDelegate(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delegate --validator [validator-addr] --amount [amount]",
		Args:  cobra.NoArgs,
		Short: "delegate liquid tokens to a validator",
		Long: strings.TrimSpace(`Delegate an amount of liquid coins to a validator from your wallet:

$ terracli tx staking delegate --validator terravaloper1l2rsakp388kuv9k8qzq6lrm9taddae7fpx59wm --amount 1000uluna --from mykey
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := authtxb.NewTxBuilderFromCLI().WithTxEncoder(auth.DefaultTxEncoder(cdc))
			cliCtx := context.NewCLIContext().
				WithCodec(cdc).
				WithAccountDecoder(cdc)

			amountStr := viper.GetString(flagAmount)
			amount, err := sdk.ParseCoin(amountStr)
			if err != nil {
				return err
			}

			delAddr := cliCtx.GetFromAddress()

			valAddrStr := viper.GetString(flagAddressValidator)
			valAddr, err := sdk.ValAddressFromBech32(valAddrStr)
			if err != nil {
				return err
			}

			msg := staking.NewMsgDelegate(delAddr, valAddr, amount)

			offline := viper.GetBool(flagOffline)

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg}, offline)
		},
	}

	cmd.Flags().AddFlagSet(fsAmount)
	cmd.Flags().AddFlagSet(fsValidator)

	cmd.MarkFlagRequired(flagAmount)
	cmd.MarkFlagRequired(flagAddressValidator)

	cmd.Flags().Bool(flagOffline, false, " Offline mode; Do not query a full node")

	return cmd
}

// GetCmdRedelegate the begin redelegation command.
func GetCmdRedelegate(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "redelegate --addr-validator-source [src-validator-addr] --addr-validator-dest [dst-validator-addr] --amount [amount]",
		Args:  cobra.NoArgs,
		Short: "redelegate illiquid tokens from one validator to another",
		Long: strings.TrimSpace(`Redelegate an amount of illiquid staking tokens from one validator to another:

$ terracli tx staking redelegate --addr-validator-source terravaloper1gghjut3ccd8ay0zduzj64hwre2fxs9ldmqhffj --addr-validator-dest terravaloper1l2rsakp388kuv9k8qzq6lrm9taddae7fpx59wm --amount 100uluna --from mykey
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := authtxb.NewTxBuilderFromCLI().WithTxEncoder(auth.DefaultTxEncoder(cdc))
			cliCtx := context.NewCLIContext().
				WithCodec(cdc).
				WithAccountDecoder(cdc)

			delAddr := cliCtx.GetFromAddress()

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

			amountStr := viper.GetString(flagAmount)
			amount, err := sdk.ParseCoin(amountStr)
			if err != nil {
				return err
			}

			msg := staking.NewMsgBeginRedelegate(delAddr, valSrcAddr, valDstAddr, amount)

			offline := viper.GetBool(flagOffline)

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg}, offline)
		},
	}

	cmd.Flags().AddFlagSet(fsAmount)
	cmd.Flags().AddFlagSet(fsRedelegation)

	cmd.MarkFlagRequired(flagAmount)
	cmd.MarkFlagRequired(flagAddressValidatorSrc)
	cmd.MarkFlagRequired(flagAddressValidatorDst)

	cmd.Flags().Bool(flagOffline, false, " Offline mode; Do not query a full node")

	return cmd
}

// GetCmdUnbond implements the unbond validator command.
func GetCmdUnbond(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "unbond --validator [validator-addr] --amount [amount]",
		Args:  cobra.NoArgs,
		Short: "unbond shares from a validator",
		Long: strings.TrimSpace(`Unbond an amount of bonded shares from a validator:

$ terracli tx staking unbond --validator terravaloper1gghjut3ccd8ay0zduzj64hwre2fxs9ldmqhffj --amount 100uluna --from mykey
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := authtxb.NewTxBuilderFromCLI().WithTxEncoder(auth.DefaultTxEncoder(cdc))
			cliCtx := context.NewCLIContext().
				WithCodec(cdc).
				WithAccountDecoder(cdc)

			delAddr := cliCtx.GetFromAddress()

			valAddrStr := viper.GetString(flagAddressValidator)
			valAddr, err := sdk.ValAddressFromBech32(valAddrStr)
			if err != nil {
				return err
			}

			amountStr := viper.GetString(flagAmount)
			amount, err := sdk.ParseCoin(amountStr)
			if err != nil {
				return err
			}

			msg := staking.NewMsgUndelegate(delAddr, valAddr, amount)

			offline := viper.GetBool(flagOffline)

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg}, offline)
		},
	}

	cmd.Flags().AddFlagSet(fsAmount)
	cmd.Flags().AddFlagSet(fsValidator)

	cmd.MarkFlagRequired(flagAmount)
	cmd.MarkFlagRequired(flagAddressValidator)

	cmd.Flags().Bool(flagOffline, false, " Offline mode; Do not query a full node")

	return cmd
}

// BuildCreateValidatorMsg makes a new MsgCreateValidator.
func BuildCreateValidatorMsg(cliCtx context.CLIContext, txBldr authtxb.TxBuilder) (authtxb.TxBuilder, sdk.Msg, error) {
	amounstStr := viper.GetString(flagAmount)
	amount, err := sdk.ParseCoin(amounstStr)
	if err != nil {
		return txBldr, nil, err
	}

	valAddr := cliCtx.GetFromAddress()
	pkStr := viper.GetString(flagPubKey)

	pk, err := sdk.GetConsPubKeyBech32(pkStr)
	if err != nil {
		return txBldr, nil, err
	}

	description := staking.NewDescription(
		viper.GetString(flagMoniker),
		viper.GetString(flagIdentity),
		viper.GetString(flagWebsite),
		viper.GetString(flagDetails),
	)

	// get the initial validator commission parameters
	rateStr := viper.GetString(flagCommissionRate)
	maxRateStr := viper.GetString(flagCommissionMaxRate)
	maxChangeRateStr := viper.GetString(flagCommissionMaxChangeRate)
	commissionMsg, err := buildCommissionMsg(rateStr, maxRateStr, maxChangeRateStr)
	if err != nil {
		return txBldr, nil, err
	}

	// get the initial validator min self delegation
	msbStr := viper.GetString(flagMinSelfDelegation)
	minSelfDelegation, ok := sdk.NewIntFromString(msbStr)
	if !ok {
		return txBldr, nil, fmt.Errorf(staking.ErrMinSelfDelegationInvalid(staking.DefaultCodespace).Error())
	}

	msg := staking.NewMsgCreateValidator(
		sdk.ValAddress(valAddr), pk, amount, description, commissionMsg, minSelfDelegation,
	)

	if viper.GetBool(client.FlagGenerateOnly) {
		ip := viper.GetString(flagIP)
		nodeID := viper.GetString(flagNodeID)
		if nodeID != "" && ip != "" {
			txBldr = txBldr.WithMemo(fmt.Sprintf("%s@%s:26656", nodeID, ip))
		}
	}

	return txBldr, msg, nil
}
