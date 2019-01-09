package cli

import (
	"fmt"
	"terra/x/oracle"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// GetCmdQueryPrice implements the query price command.
func GetCmdQueryPrice(storeName string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "price [denom]",
		Short: "Query the current price of [denom] asset",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			denom := viper.GetString(flagDenom)
			key := oracle.KeyObservedPrice(denom)
			bz, err := cliCtx.QueryStore(key, storeName)
			if err != nil {
				return err
			} else if len(bz) == 0 {
				return fmt.Errorf("No price found with denom %s", denom)
			}

			res := oracle.PriceVote{}
			cdc.MustUnmarshalBinaryLengthPrefixed(bz, &res)

			fmt.Println(res.FeedMsg.ObservedPrice.String())

			return nil
		},
	}

	cmd.Flags().AddFlagSet(fsDenom)

	return cmd
}

// GetCmdQueryTarget implements the query price command.
func GetCmdQueryTarget(storeName string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "target [denom]",
		Short: "Query the target price of [denom] asset",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			denom := viper.GetString(flagDenom)
			key := oracle.KeyTargetPrice(denom)
			bz, err := cliCtx.QueryStore(key, storeName)
			if err != nil {
				return err
			} else if len(bz) == 0 {
				return fmt.Errorf("No price found with denom %s", denom)
			}

			res := oracle.PriceVote{}
			cdc.MustUnmarshalBinaryLengthPrefixed(bz, &res)

			fmt.Println(res.FeedMsg.TargetPrice.String())

			return nil
		},
	}

	cmd.Flags().AddFlagSet(fsDenom)
	return cmd
}

// GetCmdQueryVote implements the query price command.
func GetCmdQueryVote(storeName string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "vote [denom]",
		Short: "Query the most recent oracle vote for [denom] asset",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			denom := viper.GetString(flagDenom)
			voterAddress, err := cliCtx.GetFromAddress()
			if err != nil {
				return err
			}

			key := oracle.KeyVote(denom, voterAddress)
			bz, err := cliCtx.QueryStore(key, storeName)
			if err != nil {
				return err
			} else if len(bz) == 0 {
				return fmt.Errorf("No vote found by voter %s with denom %s", voterAddress, denom)
			}

			res := oracle.PriceVote{}
			cdc.MustUnmarshalBinaryLengthPrefixed(bz, &res)

			// parse out the validator
			output, err := codec.MarshalJSONIndent(cdc, res)
			if err != nil {
				return err
			}

			fmt.Println(string(output))
			return nil
		},
	}

	cmd.Flags().AddFlagSet(fsDenom)
	cmd.MarkFlagRequired(flagDenom)

	return cmd
}

// GetCmdQueryWhitelist implements the query price command.
func GetCmdQueryWhitelist(storeName string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "whitelist",
		Short: "Query the whitelist of Terra oracle assets",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			key := oracle.KeyWhitelist
			bz, err := cliCtx.QueryStore(key, storeName)
			if err != nil {
				return err
			}

			res := []string{}
			cdc.MustUnmarshalBinaryLengthPrefixed(bz, &res)

			for _, denom := range res {
				fmt.Println(denom)
			}

			return nil
		},
	}
	return cmd
}
