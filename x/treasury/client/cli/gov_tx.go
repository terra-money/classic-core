package cli

import (
	"fmt"
	"strings"

	"github.com/classic-terra/core/v2/x/treasury/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/gov/client/cli"
	govv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	"github.com/spf13/cobra"
)

func ProposalAddBurnTaxExemptionAddressCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-burn-tax-exemption-address [addresses] --title [text] --description [text]",
		Short: "Submit an add burn tax exemption address proposal",
		Long: fmt.Sprintf(`Submit a proposal to add addresses for burn tax exemption.
Example:
$ %s tx gov submit-proposal add-burn-tax-exemption-address terra1dczz24r33fwlj0q5ra7rcdryjpk9hxm8rwy39t,terra1qt8mrv72gtvmnca9z6ftzd7slqhaf8m60aa7ye --title "add burn tax exemption address" --description "add address to burn tax exemption list"
			`, version.AppName),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			addresses := strings.Split(args[0], ",")

			proposalTitle, err := cmd.Flags().GetString(cli.FlagTitle)
			if err != nil {
				return fmt.Errorf("proposal title: %s", err)
			}
			proposalDescr, err := cmd.Flags().GetString(cli.FlagDescription)
			if err != nil {
				return fmt.Errorf("proposal description: %s", err)
			}
			depositArg, err := cmd.Flags().GetString(cli.FlagDeposit)
			if err != nil {
				return err
			}
			deposit, err := sdk.ParseCoinsNormalized(depositArg)
			if err != nil {
				return err
			}

			content := types.AddBurnTaxExemptionAddressProposal{
				Title:       proposalTitle,
				Description: proposalDescr,
				Addresses:   addresses,
			}

			msg, err := govv1beta1.NewMsgSubmitProposal(&content, deposit, clientCtx.GetFromAddress())
			if err != nil {
				return err
			}
			if err = msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	// proposal flags
	cmd.Flags().String(cli.FlagTitle, "", "Title of proposal")
	cmd.Flags().String(cli.FlagDescription, "", "Description of proposal")
	cmd.Flags().String(cli.FlagDeposit, "", "Deposit of proposal")
	return cmd
}

func ProposalRemoveBurnTaxExemptionAddressCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove-burn-tax-exemption-address [addresses] --title [text] --description [text]",
		Short: "Submit a remove burn tax exemption address proposal",
		Long: fmt.Sprintf(`Submit a proposal to remove addresses from burn tax exemption.
Example:
$ %s tx gov submit-proposal remove-burn-tax-exemption-address terra1dczz24r33fwlj0q5ra7rcdryjpk9hxm8rwy39t,terra1qt8mrv72gtvmnca9z6ftzd7slqhaf8m60aa7ye --title "remove burn tax exemption address" --description "remove address from burn tax exemption list"
			`, version.AppName),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			addresses := strings.Split(args[0], ",")

			proposalTitle, err := cmd.Flags().GetString(cli.FlagTitle)
			if err != nil {
				return fmt.Errorf("proposal title: %s", err)
			}
			proposalDescr, err := cmd.Flags().GetString(cli.FlagDescription)
			if err != nil {
				return fmt.Errorf("proposal description: %s", err)
			}
			depositArg, err := cmd.Flags().GetString(cli.FlagDeposit)
			if err != nil {
				return err
			}
			deposit, err := sdk.ParseCoinsNormalized(depositArg)
			if err != nil {
				return err
			}

			content := types.RemoveBurnTaxExemptionAddressProposal{
				Title:       proposalTitle,
				Description: proposalDescr,
				Addresses:   addresses,
			}

			msg, err := govv1beta1.NewMsgSubmitProposal(&content, deposit, clientCtx.GetFromAddress())
			if err != nil {
				return err
			}
			if err = msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	// proposal flags
	cmd.Flags().String(cli.FlagTitle, "", "Title of proposal")
	cmd.Flags().String(cli.FlagDescription, "", "Description of proposal")
	cmd.Flags().String(cli.FlagDeposit, "", "Deposit of proposal")
	return cmd
}
