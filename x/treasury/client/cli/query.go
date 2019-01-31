package cli

import (
	"fmt"
	"terra/x/treasury"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetCmdQueryAssets implements the query price command.
func GetCmdQueryAssets(storeName string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "assets",
		Short: "Query the current size of the Treasury asssets in Terra",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			bz, err := cliCtx.QueryStore(treasury.KeyIncomePool, storeName)
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
