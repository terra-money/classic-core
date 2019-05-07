package cli

import (
	flag "github.com/spf13/pflag"

	"github.com/cosmos/cosmos-sdk/x/staking/types"
)

// nolint
const (
	flagAddressValidator    = "validator"
	flagAddressValidatorSrc = "addr-validator-source"
	flagAddressValidatorDst = "addr-validator-dest"
	flagAddressDelegator    = "delegator"
	flagPubKey              = "pubkey"
	flagAmount              = "amount"
	flagSharesAmount        = "shares-amount"
	flagSharesFraction      = "shares-fraction"

	flagMoniker  = "moniker"
	flagIdentity = "identity"
	flagWebsite  = "website"
	flagDetails  = "details"

	flagCommissionRate          = "commission-rate"
	flagCommissionMaxRate       = "commission-max-rate"
	flagCommissionMaxChangeRate = "commission-max-change-rate"

	flagMinSelfDelegation = "min-self-delegation"

	flagGenesisFormat = "genesis-format"
	flagNodeID        = "node-id"
	flagIP            = "ip"

	flagOffline = "offline"
)

// common flagsets to add to various functions
var (
	fsPk                = flag.NewFlagSet("", flag.ContinueOnError)
	fsAmount            = flag.NewFlagSet("", flag.ContinueOnError)
	fsShares            = flag.NewFlagSet("", flag.ContinueOnError)
	fsDescriptionCreate = flag.NewFlagSet("", flag.ContinueOnError)
	fsCommissionCreate  = flag.NewFlagSet("", flag.ContinueOnError)
	fsCommissionUpdate  = flag.NewFlagSet("", flag.ContinueOnError)
	fsMinSelfDelegation = flag.NewFlagSet("", flag.ContinueOnError)
	fsDescriptionEdit   = flag.NewFlagSet("", flag.ContinueOnError)
	fsValidator         = flag.NewFlagSet("", flag.ContinueOnError)
	fsDelegator         = flag.NewFlagSet("", flag.ContinueOnError)
	fsRedelegation      = flag.NewFlagSet("", flag.ContinueOnError)
)

func init() {
	fsPk.String(flagPubKey, "", "The Bech32 encoded PubKey of the validator")
	fsAmount.String(flagAmount, "", "Amount of coins to bond")
	fsShares.String(flagSharesAmount, "", "Amount of source-shares to either unbond or redelegate as a positive integer or decimal")
	fsShares.String(flagSharesFraction, "", "Fraction of source-shares to either unbond or redelegate as a positive integer or decimal >0 and <=1")
	fsDescriptionCreate.String(flagMoniker, "", "The validator's name")
	fsDescriptionCreate.String(flagIdentity, "", "The optional identity signature (ex. UPort or Keybase)")
	fsDescriptionCreate.String(flagWebsite, "", "The validator's (optional) website")
	fsDescriptionCreate.String(flagDetails, "", "The validator's (optional) details")
	fsCommissionUpdate.String(flagCommissionRate, "", "The new commission rate percentage")
	fsCommissionCreate.String(flagCommissionRate, "", "The initial commission rate percentage")
	fsCommissionCreate.String(flagCommissionMaxRate, "", "The maximum commission rate percentage")
	fsCommissionCreate.String(flagCommissionMaxChangeRate, "", "The maximum commission change rate percentage (per day)")
	fsMinSelfDelegation.String(flagMinSelfDelegation, "", "The minimum self delegation required on the validator")
	fsDescriptionEdit.String(flagMoniker, types.DoNotModifyDesc, "The validator's name")
	fsDescriptionEdit.String(flagIdentity, types.DoNotModifyDesc, "The (optional) identity signature (ex. UPort or Keybase)")
	fsDescriptionEdit.String(flagWebsite, types.DoNotModifyDesc, "The validator's (optional) website")
	fsDescriptionEdit.String(flagDetails, types.DoNotModifyDesc, "The validator's (optional) details")
	fsValidator.String(flagAddressValidator, "", "The Bech32 address of the validator")
	fsDelegator.String(flagAddressDelegator, "", "The Bech32 address of the delegator")
	fsRedelegation.String(flagAddressValidatorSrc, "", "The Bech32 address of the source validator")
	fsRedelegation.String(flagAddressValidatorDst, "", "The Bech32 address of the destination validator")
}
