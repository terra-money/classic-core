package cli

import (
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/terra-project/core/x/treasury/types"
)

// GetCmdSubmitTaxRateUpdateProposal implements the command to submit a tax-rate-update proposal
func GetCmdSubmitTaxRateUpdateProposal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tax-rate-update [proposal-file]",
		Args:  cobra.ExactArgs(1),
		Short: "Submit a tax rate update proposal",
		Long: `Submit a tax rate update proposal along with an initial deposit.
The proposal details must be supplied via a JSON file.

Example:
$ terrad tx treasury submit-proposal tax-rate-update <path/to/proposal.json> --from=<key_or_address>

Where proposal.json contains:

{
  "title": "Update Tax Rate",
  "description": "Lets update tax rate to 1.5%%",
  "tax_rate": "0.015",
  "deposit": [
    {
      "denom": "stake",
      "amount": "10000"
    }
  ]
}
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			proposal, err := ParseTaxRateUpdateProposalWithDeposit(clientCtx.JSONCodec, args[0])
			if err != nil {
				return err
			}

			from := clientCtx.GetFromAddress()
			content := types.NewTaxRateUpdateProposal(proposal.Title, proposal.Description, proposal.TaxRate)

			deposit, err := sdk.ParseCoinsNormalized(proposal.Deposit)
			if err != nil {
				return err
			}

			msg, err := govtypes.NewMsgSubmitProposal(content, deposit, from)
			if err != nil {
				return err
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	return cmd
}

// GetCmdSubmitRewardWeightUpdateProposal implements the command to submit a reward-weight-update proposal
func GetCmdSubmitRewardWeightUpdateProposal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "reward-weight-update [proposal-file]",
		Args:  cobra.ExactArgs(1),
		Short: "Submit a reward weight update proposal",
		Long: `Submit a reward weight update proposal along with an initial deposit.
The proposal details must be supplied via a JSON file.

Example:
$ terrad tx treasury submit-proposal reward-weight-update <path/to/proposal.json> --from=<key_or_address>

Where proposal.json contains:

{
  "title": "Update Reward Weight",
  "description": "Lets update reward weight to 1.5%%",
  "reward_weight": "0.015",
  "deposit": [
    {
      "denom": "stake",
      "amount": "10000"
    }
  ]
}
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			proposal, err := ParseRewardWeightUpdateProposalWithDeposit(clientCtx.JSONCodec, args[0])
			if err != nil {
				return err
			}

			from := clientCtx.GetFromAddress()
			content := types.NewRewardWeightUpdateProposal(proposal.Title, proposal.Description, proposal.RewardWeight)

			deposit, err := sdk.ParseCoinsNormalized(proposal.Deposit)
			if err != nil {
				return err
			}

			msg, err := govtypes.NewMsgSubmitProposal(content, deposit, from)
			if err != nil {
				return err
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	return cmd
}
