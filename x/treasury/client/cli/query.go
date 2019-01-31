package cli

import (
	"fmt"
	"terra/x/treasury"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	flagShareID = "shareID"
)

// GetCmdQueryAssets implements the query price command.
func GetCmdQueryAssets(storeName string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "assets",
		Short: "Query the current size of the Treasury asssets in Terra",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			bz, err := cliCtx.QueryStore(treasury.GetIncomePoolKey(), storeName)
			if err != nil {
				return err
			} else if len(bz) == 0 {
				panic("No income pool found")
			}

			res := sdk.Coins{}
			cdc.MustUnmarshalBinaryLengthPrefixed(bz, &res)

			fmt.Println(res[0].Amount)

			return nil
		},
	}

	return cmd
}

// GetCmdQueryShare implements the query price command.
func GetCmdQueryShare(storeName string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "share [shareID]",
		Short: "Query the share corresponding to [shareID] and fetch attendant claims. Share ID is one of 'oracle', 'debt', and 'budget'",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			shareId := viper.GetString(flagShareID)
			key := treasury.GetShareKey(shareId)
			bz, err := cliCtx.QueryStore(key, storeName)
			if err != nil {
				return err
			} else if len(bz) == 0 {
				return fmt.Errorf("No share found with id %s", shareId)
			}

			var share treasury.Share
			cdc.MustUnmarshalBinaryLengthPrefixed(bz, &share)

			fmt.Printf("Share id: %s weight: %f ", shareId, share.GetWeight())

			// fetch claims
			claimKey := treasury.GetClaimsForSharePrefix(shareId)
			claimBytes, err := cliCtx.QueryStore(claimKey, storeName)
			if err != nil {
				return err
			}

			var matchingClaims []treasury.Claim
			err = cdc.UnmarshalJSON(claimBytes, &matchingClaims)
			if err != nil {
				return err
			}

			if len(matchingClaims) == 0 {
				fmt.Println("No matching claims found")
				return nil
			}

			for _, claim := range matchingClaims {
				fmt.Printf("  %f - %s\n", claim.GetWeight(), claim.ID())
			}

			return nil
		},
	}

	cmd.Flags().String(flagShareID, "", "id of the share to query")
	cmd.MarkFlagRequired(flagShareID)

	return cmd
}
