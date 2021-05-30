package cli

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/cosmos/cosmos-sdk/x/gov"

	"github.com/terra-money/core/x/treasury/internal/types"
)

// GetCmdSubmitTaxRateUpdateProposal implements the command to submit a tax-rate-update proposal
func GetCmdSubmitTaxRateUpdateProposal(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tax-rate-update [proposal-file]",
		Args:  cobra.ExactArgs(1),
		Short: "Submit a tax rate update proposal",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Submit a tax rate update proposal along with an initial deposit.
The proposal details must be supplied via a JSON file.

Example:
$ %s tx treasury submit-proposal tax-rate-update <path/to/proposal.json> --from=<key_or_address>

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
				version.ClientName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			proposal, err := ParseTaxRateUpdateProposalJSON(cdc, args[0])
			if err != nil {
				return err
			}

			from := cliCtx.GetFromAddress()
			content := types.NewTaxRateUpdateProposal(proposal.Title, proposal.Description, proposal.TaxRate)

			msg := gov.NewMsgSubmitProposal(content, proposal.Deposit, from)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	return cmd
}

// GetCmdSubmitRewardWeightUpdateProposal implements the command to submit a reward-weight-update proposal
func GetCmdSubmitRewardWeightUpdateProposal(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "reward-weight-update [proposal-file]",
		Args:  cobra.ExactArgs(1),
		Short: "Submit a reward weight update proposal",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Submit a reward weight update proposal along with an initial deposit.
The proposal details must be supplied via a JSON file.

Example:
$ %s tx treasury submit-proposal reward-weight-update <path/to/proposal.json> --from=<key_or_address>

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
				version.ClientName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			proposal, err := ParseRewardWeightUpdateProposalJSON(cdc, args[0])
			if err != nil {
				return err
			}

			from := cliCtx.GetFromAddress()
			content := types.NewRewardWeightUpdateProposal(proposal.Title, proposal.Description, proposal.RewardWeight)

			msg := gov.NewMsgSubmitProposal(content, proposal.Deposit, from)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	return cmd
}
