package cli

import (
	flag "github.com/spf13/pflag"
)

const (
	flagAddressValidator  = "validator"
	flagAddressDelegator  = "delegator"
	flagStartHeight       = "start"
	flagEndHeight         = "end"
	flagOnlyFromValidator = "only-from-validator"
	flagIsValidator       = "is-validator"
	flagComission         = "commission"
	flagWithdrawTo        = "withdraw-to"
)

var (
	fsValidator = flag.NewFlagSet("", flag.ContinueOnError)
	fsDelegator = flag.NewFlagSet("", flag.ContinueOnError)
)

func init() {
	fsValidator.String(flagAddressValidator, "", "The Bech32 address of the validator")
	fsDelegator.String(flagAddressDelegator, "", "The Bech32 address of the delegator")
}
