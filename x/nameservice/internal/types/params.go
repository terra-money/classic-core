package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/params/subspace"
	core "github.com/terra-project/core/types"
	"time"
)

// DefaultParamspace defines default space for oracle params
const DefaultParamspace = ModuleName

// Parameter keys
var (
	ParamStoreKeyBidPeriod       = []byte("bidperiod")
	ParamStoreKeyRevealPeriod    = []byte("revealperiod")
	ParamStoreKeyGracePeriod     = []byte("graceperiod")
	ParamStoreKeyRenewalInterval = []byte("renewalinterval")
	ParamStoreKeyMinDeposit      = []byte("mindeposit")
	ParamStoreKeyRootName        = []byte("rootname")
	ParamStoreKeyMinNameLength   = []byte("minnamelength")
	ParamStoreKeyRenewalFees     = []byte("renewalfees")
)

// Default parameter values
const (
	DefaultBidPeriod       = time.Hour * 24 * 2   // 2 days
	DefaultRevealPeriod    = time.Hour * 24 * 3   // 3 days
	DefaultGracePeriod     = time.Hour * 24 * 5   // 5 days
	DefaultRenewalInterval = time.Hour * 24 * 365 // 365 days
	DefaultRootName        = "terra"
	DefaultMinNameLength   = 3
)

// Default parameter values
var (
	DefaultMinDeposit  = sdk.NewInt64Coin(core.MicroLunaDenom, 512*core.MicroUnit) // 512 LUNA
	DefaultRenewalFees = RenewalFees{
		{3, sdk.NewInt64Coin(core.MicroSDRDenom, 400*core.MicroUnit)},
		{4, sdk.NewInt64Coin(core.MicroSDRDenom, 100*core.MicroUnit)},
		{5, sdk.NewInt64Coin(core.MicroSDRDenom, 5*core.MicroUnit)},
	}
)

var _ subspace.ParamSet = &Params{}

// Params nameservice parameters
type Params struct {
	BidPeriod       time.Duration `json:"bid_period" yaml:"bid_period"`
	RevealPeriod    time.Duration `json:"reveal_period" yaml:"reveal_period"`
	RenewalInterval time.Duration `json:"renewal_interval" yaml:"renewal_interval"`
	MinDeposit      sdk.Coin      `json:"min_deposit" yaml:"min_deposit"`
	RootName        Name          `json:"root_name" yaml:"root_name"`
	MinNameLength   int           `json:"min_name_length" yaml:"min_name_length"`
	RenewalFees     RenewalFees   `json:"renewal_fees" yaml:"renewal_fees"`
}

// DefaultParams creates default oracle module parameters
func DefaultParams() Params {
	return Params{
		BidPeriod:       DefaultBidPeriod,
		RevealPeriod:    DefaultRevealPeriod,
		RenewalInterval: DefaultRenewalInterval,
		MinDeposit:      DefaultMinDeposit,
		RootName:        DefaultRootName,
		MinNameLength:   DefaultMinNameLength,
		RenewalFees:     DefaultRenewalFees,
	}
}

// Validate validates a set of params
func (params Params) Validate() error {
	if params.BidPeriod.Minutes() < 10 {
		return fmt.Errorf("nameservice parameter BidPeriod must be bigger than 10 minutes")
	}

	if params.RevealPeriod.Minutes() < 10 {
		return fmt.Errorf("nameservice parameter RevealPeriod must be bigger than 10 minutes")
	}

	if params.RenewalInterval.Minutes() < 10 {
		return fmt.Errorf("nameservice parameter RenewalInterval must be bigger than 10 minutes")
	}

	if !params.MinDeposit.IsValid() {
		return fmt.Errorf("nameservice parameter MinDeposit must be valid coin")
	}

	if err := params.RootName.Validate(); err != nil {
		return fmt.Errorf("nameservice parameter RootName is invalid; %s", err.Error())
	}

	if params.MinNameLength <= 0 {
		return fmt.Errorf("nameservice parameter MinNameLength must be bigger than 0")
	}

	uniqueMap := make(map[int]bool)
	for _, fee := range params.RenewalFees {
		if _, ok := uniqueMap[fee.Length]; ok {
			return fmt.Errorf("nameservice parameter RenewalFees has duplicated items")
		}
		if fee.Length < 0 {
			return fmt.Errorf("nameservice parameter RenewalFees can't hold negative Length")
		}
		if !fee.Amount.IsValid() {
			return fmt.Errorf("nameservice parameter RenewalFees contains invalid coins")
		}

		uniqueMap[fee.Length] = true
	}

	return nil
}

// ParamKeyTable for namespace module
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable().RegisterParamSet(&Params{})
}

// ParamSetPairs implements the ParamSet interface and returns all the key/value pairs
// pairs of oracle module's parameters.
func (params *Params) ParamSetPairs() subspace.ParamSetPairs {
	return subspace.ParamSetPairs{
		{Key: ParamStoreKeyBidPeriod, Value: &params.BidPeriod},
		{Key: ParamStoreKeyRevealPeriod, Value: &params.RevealPeriod},
		{Key: ParamStoreKeyRenewalInterval, Value: &params.RenewalInterval},
		{Key: ParamStoreKeyMinDeposit, Value: &params.MinDeposit},
		{Key: ParamStoreKeyRootName, Value: &params.RootName},
		{Key: ParamStoreKeyMinNameLength, Value: &params.MinNameLength},
		{Key: ParamStoreKeyRenewalFees, Value: &params.RenewalFees},
	}
}

// String implements fmt.Stringer interface
func (params Params) String() string {
	return fmt.Sprintf(`Nameservice Params:
	BidPeriod:       %d
	RevealInterval:  %d
	RenewalInterval: %d
	MinDeposit       %s
    RootName         %s
    MinNameLength    %d
    RenewalFees      %s
	`, params.BidPeriod, params.RevealPeriod,
		params.RenewalInterval, params.MinDeposit, params.RootName,
		params.MinNameLength, params.RenewalFees)
}
