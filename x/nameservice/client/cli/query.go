package cli

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/terra-project/core/x/nameservice/internal/types"
)

const flagStatus = "status"

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	nameserviceQueryCmd := &cobra.Command{
		Use:                        "nameservice",
		Aliases:                    []string{"ns"},
		Short:                      "Querying commands for the nameservice module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	nameserviceQueryCmd.AddCommand(client.GetCommands(
		GetCmdQueryAuctions(cdc),
		GetCmdQueryBids(cdc),
		GetCmdQueryRegistry(cdc),
		GetCmdQueryResolve(cdc),
		GetCmdQueryReverse(cdc),
		GetCmdQueryParams(cdc))...)

	return nameserviceQueryCmd
}

// GetCmdQueryAuctions implements the query auction command.
func GetCmdQueryAuctions(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "auctions [name]",
		Args:  cobra.RangeArgs(0, 1),
		Short: "Query auctions",
		Long: `Query auctions with name.
$ terracli query ns auctions "chai.terra"

It is possible to query without name argument, then it returns all active auctions.
$ terracli query ns auctions

Or, can filter active auctions with status
$ terracli query ns auctions --status Bid
	`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			var bz []byte
			if len(args) == 1 {
				name := types.Name(args[0])
				if name.Levels() != 2 {
					return fmt.Errorf("must submit by the second level name")
				}

				nameHash, _ := name.NameHash()
				params := types.QueryAuctionsParams{NameHash: nameHash}
				bz = cliCtx.Codec.MustMarshalJSON(params)
			} else {
				statusStr := viper.GetString(flagStatus)
				status, err := types.AuctionStatusFromString(statusStr)
				if err != nil {
					return err
				}

				params := types.QueryAuctionsParams{Status: status}
				bz = cliCtx.Codec.MustMarshalJSON(params)
			}

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryAuctions), bz)
			if err != nil {
				return err
			}

			var auctions types.Auctions
			cdc.MustUnmarshalJSON(res, &auctions)
			return cliCtx.PrintOutput(auctions)
		},
	}

	cmd.Flags().String(flagStatus, types.AuctionStatusNil.String(), "auction status must be one of (Nil, Bid, Reveal)")
	return cmd
}

// GetCmdQueryBids implements the query bid command.
func GetCmdQueryBids(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bids [name] [bidderAddr]",
		Args:  cobra.RangeArgs(1, 2),
		Short: "Query bids",
		Long: `Query bids of an auction.
$ terracli query ns bids "chai.terra"

It is possible to query with bidder address
$ terracli query ns bids "chai.terra" terra~~
	`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			name := types.Name(args[0])
			if name.Levels() != 2 {
				return fmt.Errorf("must submit by the second level name")
			}

			nameHash, _ := name.NameHash()

			var bz []byte
			if len(args) == 2 {
				addr, err := sdk.AccAddressFromBech32(args[1])
				if err != nil {
					return err
				}

				params := types.QueryBidsParams{NameHash: nameHash, Bidder: addr}
				bz = cliCtx.Codec.MustMarshalJSON(params)
			} else {
				params := types.QueryBidsParams{NameHash: nameHash, Bidder: nil}
				bz = cliCtx.Codec.MustMarshalJSON(params)
			}

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryBids), bz)
			if err != nil {
				return err
			}

			var bids types.Bids
			cdc.MustUnmarshalJSON(res, &bids)
			return cliCtx.PrintOutput(bids)
		},
	}

	return cmd
}

// GetCmdQueryRegistry implements the query registry command.
func GetCmdQueryRegistry(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "registry [name]",
		Args:  cobra.ExactArgs(1),
		Short: "Query registry",
		Long: `Query registry of the name.
$ terracli query ns registry "chai.terra"
	`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			name := types.Name(args[0])
			if name.Levels() != 2 {
				return fmt.Errorf("must submit by the second level name")
			}

			nameHash, _ := name.NameHash()

			params := types.QueryRegistryParams{NameHash: nameHash}
			bz := cliCtx.Codec.MustMarshalJSON(params)

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryRegistry), bz)
			if err != nil {
				return err
			}

			var registry types.Registry
			cdc.MustUnmarshalJSON(res, &registry)
			return cliCtx.PrintOutput(registry)
		},
	}

	return cmd
}

// GetCmdQueryResolve implements the query resolve command.
func GetCmdQueryResolve(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "resolve [name]",
		Args:  cobra.ExactArgs(1),
		Short: "Resolve given name to account address",
		Long: `Resolve second level name to account address
$ terracli query ns resolve "chai.terra"

It is possible to resolve third level name
$ terracli query ns resolve "account.chai.terra"
	`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			name := types.Name(args[0])
			if levels := name.Levels(); levels != 2 && levels != 3 {
				return fmt.Errorf("must submit by the second or third level name")
			}

			nameHash, childNameHash := name.NameHash()
			params := types.QueryResolveParams{NameHash: nameHash, ChildNameHash: childNameHash}
			bz := cliCtx.Codec.MustMarshalJSON(params)

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryResolve), bz)
			if err != nil {
				return err
			}

			var addr sdk.AccAddress
			cdc.MustUnmarshalJSON(res, &addr)
			return cliCtx.PrintOutput(addr)
		},
	}

	return cmd
}

// GetCmdQueryReverse implements the query reverse command.
func GetCmdQueryReverse(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "reverse [addr]",
		Args:  cobra.ExactArgs(1),
		Short: "Query a registry with the address",
		Long: `Query a registry that contains the address as an entry.
$ terracli query ns reverse terra~
	`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			addr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			params := types.QueryReverseParams{Address: addr}
			bz := cliCtx.Codec.MustMarshalJSON(params)

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryReverse), bz)
			if err != nil {
				return err
			}

			var registry types.Registry
			cdc.MustUnmarshalJSON(res, &registry)
			return cliCtx.PrintOutput(registry)
		},
	}

	return cmd
}

// GetCmdQueryParams implements the query params command.
func GetCmdQueryParams(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "params",
		Args:  cobra.NoArgs,
		Short: "Query the current market params",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryParameters), nil)
			if err != nil {
				return err
			}

			var params types.Params
			cdc.MustUnmarshalJSON(res, &params)
			return cliCtx.PrintOutput(params)
		},
	}

	return cmd
}
