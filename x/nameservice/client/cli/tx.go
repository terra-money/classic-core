package cli

import (
	"fmt"
	"github.com/spf13/cobra"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"

	"github.com/terra-project/core/x/nameservice/internal/types"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd(cdc *codec.Codec) *cobra.Command {
	marketTxCmd := &cobra.Command{
		Use:                        "nameservice",
		Aliases:                    []string{"ns"},
		Short:                      "Nameservice transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	marketTxCmd.AddCommand(client.PostCommands(
		GetOpenAuctionCmd(cdc),
		GetBidAuctionCmd(cdc),
		GetRenewRegistryCmd(cdc),
		GetUpdateOwnerCmd(cdc),
		GetRevealBidCmd(cdc),
		GetRegisterSubNameCmd(cdc),
		GetUnregisterSubNameCmd(cdc),
	)...)

	return marketTxCmd
}

// GetOpenAuctionCmd will create and send a MsgOpenAuction
func GetOpenAuctionCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "open-auction [name]",
		Aliases: []string{"open"},
		Args:    cobra.ExactArgs(1),
		Short:   "Create an auction for the name",
		Long: strings.TrimSpace(`
Create an auction for the name. 

$ terracli ns open "wallet.terra"
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			name := types.Name(args[0])
			if err := name.Validate(); err != nil {
				return err
			}

			if levels := name.Levels(); levels != 2 {
				return fmt.Errorf("must submit by the second level name")
			}

			fromAddress := cliCtx.GetFromAddress()

			// build and sign the transaction, then broadcast to Tendermint
			msg := types.NewMsgOpenAuction(name, fromAddress)
			err := msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	return cmd
}

// GetBidAuctionCmd will create and send a MsgBidAuction
func GetBidAuctionCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "bid-auction [name] [amount] [deposit] [salt]",
		Aliases: []string{"bid"},
		Args:    cobra.ExactArgs(4),
		Short:   "Submit bid for an auction of the name",
		Long: strings.TrimSpace(`
Submit bid for an auction of the name. To hide the bid amount, 
it creates a hash containing the amount and salt. The amount actually 
transferred is deposit, and this value must be greater than 
the bid amount.

$ terracli ns bid-auction "wallet.terra" "100uluna" "1000uluna" "salt"
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			name := types.Name(args[0])
			if err := name.Validate(); err != nil {
				return err
			}

			if levels := name.Levels(); levels != 2 {
				return fmt.Errorf("must submit by the second level name")
			}

			bidAmount, err := sdk.ParseCoin(args[1])
			if err != nil {
				return err
			}

			deposit, err := sdk.ParseCoin(args[2])
			if err != nil {
				return err
			}

			if deposit.IsLT(bidAmount) {
				return fmt.Errorf("deposit must be bigger than bid amount")
			}
			salt := args[3]

			fromAddress := cliCtx.GetFromAddress()

			bidHash := types.GetBidHash(salt, name, bidAmount, fromAddress)

			// build and sign the transaction, then broadcast to Tendermint
			msg := types.NewMsgBidAuction(name, bidHash.String(), deposit, fromAddress)
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	return cmd
}

// GetRevealBidCmd will create and send a MsgRevealBid
func GetRevealBidCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "reveal-bid [name] [amount] [salt]",
		Aliases: []string{"reveal"},
		Args:    cobra.ExactArgs(3),
		Short:   "Submit a message to reveal a bid",
		Long: strings.TrimSpace(`
Submit a message to reveal a bid whose value is hidden by hash. 

$ terracli ns reveal-bid "wallet.terra" "100uluna" "salt"
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			name := types.Name(args[0])
			if err := name.Validate(); err != nil {
				return err
			}

			if levels := name.Levels(); levels != 2 {
				return fmt.Errorf("must submit by the second level name")
			}

			bidAmount, err := sdk.ParseCoin(args[1])
			if err != nil {
				return err
			}

			salt := args[2]

			fromAddress := cliCtx.GetFromAddress()

			// build and sign the transaction, then broadcast to Tendermint
			msg := types.NewMsgRevealBid(name, salt, bidAmount, fromAddress)
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	return cmd
}

// GetRenewRegistryCmd will create and send a MsgRenewRegistry
func GetRenewRegistryCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "renew [name] [amount]",
		Args:  cobra.ExactArgs(2),
		Short: "Extend the name expiration period by paying fee",
		Long: strings.TrimSpace(`
Extend the name expiration period by paying fee. Only the registry 
owner can extend. Any coins registered in the oracle can be paid, 
and the extended expiration period is automatically calculated 
based on the amount paid.

$ terracli ns renew "wallet.terra" "100uluna,20000ukrw"
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			name := types.Name(args[0])
			if err := name.Validate(); err != nil {
				return err
			}

			if levels := name.Levels(); levels != 2 {
				return fmt.Errorf("must submit by the second level name")
			}

			fees, err := sdk.ParseCoins(args[1])
			if err != nil {
				return err
			}

			fromAddress := cliCtx.GetFromAddress()

			// build and sign the transaction, then broadcast to Tendermint
			msg := types.NewMsgRenewRegistry(name, fees, fromAddress)
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	return cmd
}

// GetUpdateOwnerCmd will create and send a MsgUpdateOwner
func GetUpdateOwnerCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-owner [name] [owner-addr]",
		Args:  cobra.ExactArgs(2),
		Short: "Change the owner for the name",
		Long: strings.TrimSpace(`
Change the owner address for the name. Only the registry 
owner is allowed to execute this msg.

$ terracli ns update-owner "wallet.terra" terra~
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			name := types.Name(args[0])
			if err := name.Validate(); err != nil {
				return err
			}

			if levels := name.Levels(); levels != 2 {
				return fmt.Errorf("must submit by the second level name")
			}

			addr, err := sdk.AccAddressFromBech32(args[1])
			if err != nil {
				return err
			}

			fromAddress := cliCtx.GetFromAddress()

			// build and sign the transaction, then broadcast to Tendermint
			msg := types.NewMsgUpdateOwner(name, addr, fromAddress)
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	return cmd
}

// GetRegisterSubNameCmd will create and send a MsgRegisterSubName
func GetRegisterSubNameCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "register [name] [address]",
		Args:  cobra.ExactArgs(2),
		Short: "Register (sub-name, address) entry in the name registry",
		Long: strings.TrimSpace(`
Register (sub-name, address) in the name registry. Only the registry 
owner can register. Each account address can be registered only once.

$ terracli ns register "acc1.wallet.terra" terra~
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			name := types.Name(args[0])
			if err := name.Validate(); err != nil {
				return err
			}

			if levels := name.Levels(); levels != 2 && levels != 3 {
				return fmt.Errorf("must submit by the second or third level name")
			}

			addr, err := sdk.AccAddressFromBech32(args[1])
			if err != nil {
				return err
			}

			fromAddress := cliCtx.GetFromAddress()

			// build and sign the transaction, then broadcast to Tendermint
			msg := types.NewMsgRegisterSubName(name, addr, fromAddress)
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	return cmd
}

// GetUnregisterSubNameCmd will create and send a MsgUnregisterSubName
func GetUnregisterSubNameCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "unregister [name]",
		Args:  cobra.ExactArgs(1),
		Short: "Unregister sub-name from the name registry",
		Long: strings.TrimSpace(`
Unregister sub-name from the name registry. Only the registry 
owner can unregister.

$ terracli ns unregister "acc1.wallet.terra"
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			name := types.Name(args[0])
			if err := name.Validate(); err != nil {
				return err
			}

			if levels := name.Levels(); levels != 2 && levels != 3 {
				return fmt.Errorf("must submit by the second or third level name")
			}

			fromAddress := cliCtx.GetFromAddress()

			// build and sign the transaction, then broadcast to Tendermint
			msg := types.NewMsgUnregisterSubName(name, fromAddress)
			err := msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	return cmd
}
