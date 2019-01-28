package feeder

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/utils"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtxb "github.com/cosmos/cosmos-sdk/x/auth/client/txbuilder"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/libs/log"
	"math"
	"os"
	"terra/x/oracle"
	"time"

	// Import statik for light client stuff
	_ "github.com/cosmos/cosmos-sdk/client/lcd/statik"
)

const (
	flagUpdateInterval  = "update-interval"
	flagDataSourceFile  = "data-source-file"
	flagDataSourceURL   = "data-source-url"
	flagDataSourceFixed = "data-source-json"
)

// FeedDaemon represents the oracle feeder daemon
type FeedDaemon struct {
	CliCtx context.CLIContext
	Cdc    *codec.Codec

	source Source
	log    log.Logger
}

// NewFeedDaemon creates a new rest server instance
func NewFeedDaemon(cdc *codec.Codec) (*FeedDaemon, error) {
	cliCtx := context.NewCLIContext().WithCodec(cdc).WithAccountDecoder(cdc)

	logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout)).With("module", "feed-server")

	var src Source
	if viper.IsSet(flagDataSourceURL) {
		src = URLSource{viper.GetString(flagDataSourceURL)}

	} else if viper.IsSet(flagDataSourceFile) {
		src = FileSource{viper.GetString(flagDataSourceFile)}

	} else if viper.IsSet(flagDataSourceFixed) {
		var err error

		src, err = CreateJsonSource(viper.GetString(flagDataSourceFixed))
		if err != nil {
			return nil, err
		}

	} else {
		return nil, fmt.Errorf("One of the sources should have to be set.")
	}

	return &FeedDaemon{
		CliCtx: cliCtx,
		Cdc:    cdc,

		source: src,
		log:    logger,
	}, nil
}

func votePrice(price Price, cdc *codec.Codec) error {

	txBldr := authtxb.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
	cliCtx := context.NewCLIContext().WithCodec(cdc).WithAccountDecoder(cdc)

	voterAddress, err := cliCtx.GetFromAddress()
	if err != nil {
		return err
	}

	target := sdk.NewDecWithPrec(int64(math.Round(price.TargetPrice*100)), 2)
	current := sdk.NewDecWithPrec(int64(math.Round(price.CurrentPrice*100)), 2)

	// build and sign the transaction, then broadcast to Tendermint
	msg := oracle.NewPriceFeedMsg(price.Denom, target, current, voterAddress)
	if cliCtx.GenerateOnly {
		return utils.PrintUnsignedStdTx(os.Stdout, txBldr, cliCtx, []sdk.Msg{msg}, false)
	}

	return utils.CompleteAndBroadcastTxCli(txBldr, cliCtx, []sdk.Msg{msg})

}

func (fd *FeedDaemon) voteAll(cdc *codec.Codec) error {
	prices, err := fd.source.getData()
	if err != nil {
		return err
	}

	for _, price := range prices {
		err := votePrice(price, cdc)

		if err != nil {
			fd.log.Error(err.Error())
		}
	}

	return nil
}

// Start starts the rest server
func (fd *FeedDaemon) Start() (err error) {

	interval := viper.GetDuration(flagUpdateInterval)

	//
	done := make(chan bool)
	shutdown := make(chan bool)

	server.TrapSignal(func() {
		close(shutdown)
	})

	go func() {
	loop:
		for {
			select {
			case <-time.After(interval):
				err := fd.voteAll(fd.CliCtx.Codec)
				if err != nil {
					fd.log.Error(err.Error())
				}

			case <-shutdown: // triggered on the stop signal
				break loop   // exit
			}
		}
	}()

	<-done

	return nil
}

// ServeCommand will start a Gaia Lite REST service as a blocking process. It
// takes a codec to create a FeedDaemon object and a function to register all
// necessary routes.
func ServeCommand(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "feed-daemon",
		Short: "Start Feeder, a Oracle feeding daemon",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			fd, err := NewFeedDaemon(cdc)
			if err != nil {
				return err
			}

			err = fd.Start()

			return err
		},
	}

	//
	cmd.Flags().String(client.FlagListenAddr, "tcp://localhost:1317", "The address for the server to listen on")
	cmd.Flags().String(client.FlagChainID, "", "Chain ID of Tendermint node")
	cmd.Flags().String(client.FlagNode, "tcp://localhost:26657", "Address of the node to connect to")
	cmd.Flags().Bool(client.FlagTrustNode, false, "Trust connected full node (don't verify proofs for responses)")

	viper.BindPFlag(client.FlagTrustNode, cmd.Flags().Lookup(client.FlagTrustNode))
	viper.BindPFlag(client.FlagChainID, cmd.Flags().Lookup(client.FlagChainID))
	viper.BindPFlag(client.FlagNode, cmd.Flags().Lookup(client.FlagNode))

	return cmd
}
