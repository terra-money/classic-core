package cli

import (
	flag "github.com/spf13/pflag"
)

const (
	flagDenom        = "denom"
	flagTargetPrice  = "targetprice"
	flagCurrentPrice = "currentprice"
	flagVoterAddress = "address"
)

// common flagsets to add to various functions
var (
	fsDenom        = flag.NewFlagSet("", flag.ContinueOnError)
	fsTargetPrice  = flag.NewFlagSet("", flag.ContinueOnError)
	fsCurrentPrice = flag.NewFlagSet("", flag.ContinueOnError)
	//fsVoterAddress = flag.NewFlagSet("", flag.ContinueOnError)
)

func init() {
	fsDenom.String(flagDenom, "", "Denomination of the coin")
	fsTargetPrice.Float32(flagTargetPrice, 0.0, "Target price in Luna for the coin")
	fsCurrentPrice.Float32(flagCurrentPrice, 0.0, "Current price in Luna for the coin")
	//fsVoterAddress.String(flagVoterAddress, "", "Bech32 validator address of the voter")
}
