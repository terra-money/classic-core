package treasury

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	core "github.com/terra-money/core/types"
)

// Default parameter values
var (
	DefaultTaxPolicy = PolicyConstraints{
		RateMin:       sdk.ZeroDec(),
		RateMax:       sdk.ZeroDec(),
		Cap:           sdk.NewCoin(core.MicroSDRDenom, sdk.ZeroInt()),
		ChangeRateMax: sdk.ZeroDec(),
	}
	DefaultRewardPolicy = PolicyConstraints{
		RateMin:       sdk.ZeroDec(),
		RateMax:       sdk.ZeroDec(),
		Cap:           sdk.NewCoin("unused", sdk.ZeroInt()),
		ChangeRateMax: sdk.ZeroDec(),
	}
	DefaultSeigniorageBurdenTarget = sdk.ZeroDec()
	DefaultMiningIncrement         = sdk.ZeroDec()
	DefaultWindowShort             = uint64(0)
	DefaultWindowLong              = uint64(0)
	DefaultWindowProbation         = uint64(0)
	DefaultTaxRate                 = sdk.ZeroDec()
	DefaultRewardWeight            = sdk.ZeroDec()
)

// DefaultParams creates default treasury module parameters
func DefaultParams() Params {
	return Params{
		TaxPolicy:               DefaultTaxPolicy,
		RewardPolicy:            DefaultRewardPolicy,
		SeigniorageBurdenTarget: DefaultSeigniorageBurdenTarget,
		MiningIncrement:         DefaultMiningIncrement,
		WindowShort:             DefaultWindowShort,
		WindowLong:              DefaultWindowLong,
		WindowProbation:         DefaultWindowProbation,
	}
}
