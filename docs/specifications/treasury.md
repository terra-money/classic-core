# Treasury

The Treasury module is the "central bank" of the Terra economy. It monitors changes in the macroeconomic variables, and adjusts Terra monetary polices accordingly.

## Observed macroecnomic variables

The treasury observes two main variables:

* Periodic Tax returns: stability fee income that has been made in a given period, namely `WindowShort`. In tandem to this, the treasury also monitors the rolling average of stability fee revenues in a longer timeframe, `WindowLong`, to be able to compare the performance of tax income in a shorter time window compared to a longer one.
* Periodic Terra seigniorage burn: the amount of Terra seigniorage that has been burned \(total periodic seigniorage \* mining reward weight\) in a given period. Similar to tax returns, the treasury measures the index in light of a short and long time window.

Tax income and seigniorage burn combined makes up the total mining rewards for Luna.

## Monetary policy tools

The treasury module has two monetary policy levers in its toolkit. The tax rate, by which it can increase fees coming in from Terra transactions, and and the mining reward weight, which is the portion of seigniorage that is burned to reward miners via scarcity. Every `WindowLong`, it re-evaluates each lever to stabilize unit staking returns for Luna, thereby optimizing for stable cash flows from Terra staking.

### Tax rate

```go
// t(t+1) = t(t) * (TL_year(t) + INC) / TL_month(t)
func (k Keeper) UpdateTaxPolicy(ctx sdk.Context) (newTaxRate sdk.Dec) {
	params := k.GetParams(ctx)

	oldTaxRate := k.GetTaxRate(ctx)
	inc := params.MiningIncrement
	tlYear := k.rollingAverageIndicator(ctx, params.WindowLong, types.TRLKey)
	tlMonth := k.rollingAverageIndicator(ctx, params.WindowShort, types.TRLKey)

	// No revenues, hike as much as possible.
	if tlMonth.Equal(sdk.ZeroDec()) {
		newTaxRate = params.TaxPolicy.RateMax
	} else {
		newTaxRate = oldTaxRate.Mul(tlYear.Mul(inc)).Quo(tlMonth)
	}

	newTaxRate = params.TaxPolicy.Clamp(oldTaxRate, newTaxRate)

	// Set the new tax rate to the store
	k.SetTaxRate(ctx, newTaxRate)
	return
}
```

At the point of evaluation, the treasury hikes up tax rates when tax revenues in a shorter time window is performing poorly in comparison to the longer term tax revenue average. It lowers tax rates when short term tax revenues are outperforming the longer term index.

### Reward weight

```go
// w(t+1) = w(t)*SB_target/SB_rolling(t)
func (k Keeper) UpdateRewardPolicy(ctx sdk.Context) (newRewardWeight sdk.Dec) {
	params := k.GetParams(ctx)

	oldWeight := k.GetRewardWeight(ctx)
	sbTarget := params.SeigniorageBurdenTarget

	seigniorageSum := k.sumIndicator(ctx, params.WindowShort, types.SRKey)
	totalSum := k.sumIndicator(ctx, params.WindowShort, types.MRKey)

	// No revenues; hike as much as possible
	if totalSum.Equal(sdk.ZeroDec()) || seigniorageSum.Equal(sdk.ZeroDec()) {
		newRewardWeight = params.RewardPolicy.RateMax
	} else {
		// Seigniorage burden out of total rewards
		sb := seigniorageSum.Quo(totalSum)
		newRewardWeight = oldWeight.Mul(sbTarget.Quo(sb))
	}

	newRewardWeight = params.RewardPolicy.Clamp(oldWeight, newRewardWeight)

	// Set the new reward weight
	k.SetRewardWeight(ctx, newRewardWeight)
	return
}
```

The treasury mirrors the tax rate when adjusting the mining reward weight. It observes the overall burden seigniorage burn needs to bear in the overall reward profile, `SeigniorageBurdenTarget`, and hikes up rates accordingly as tax rates rise. In order to make sure that unit mining rewards do not stay stagnant, the treasury adds a `MiningIncrement` to each policy update, such that mining rewards increase steadily over time.

### Policy contraints

```go
// PolicyConstraints wraps constraints around updating a key Treasury variable
type PolicyConstraints struct {
    RateMin       sdk.Dec  `json:"rate_min"`
    RateMax       sdk.Dec  `json:"rate_max"`
    Cap           sdk.Coin `json:"cap"`
    ChangeRateMax sdk.Dec  `json:"change_max"`
}

// Clamp constrains a policy variable update within the policy constraints
func (pc PolicyConstraints) Clamp(prevRate sdk.Dec, newRate sdk.Dec) (clampedRate sdk.Dec) {
    if newRate.LT(pc.RateMin) {
        newRate = pc.RateMin
    } else if newRate.GT(pc.RateMax) {
        newRate = pc.RateMax
    }

    delta := newRate.Sub(prevRate)
    if newRate.GT(prevRate) {
        if delta.GT(pc.ChangeRateMax) {
            newRate = prevRate.Add(pc.ChangeRateMax)
        }
    } else {
        if delta.Abs().GT(pc.ChangeRateMax) {
            newRate = prevRate.Sub(pc.ChangeRateMax)
        }
    }
    return newRate
}
```

Both tax rate and seigniorage burn weight updates are limited by `PolicyConstraint`, which specifies the floor, ceiling, and the max periodic changes for each variable.

## Parameters

```go
// Params treasury parameters
type Params struct {
    TaxPolicy    PolicyConstraints `json:"tax_policy"`
    RewardPolicy PolicyConstraints `json:"reward_policy"`

    SeigniorageBurdenTarget sdk.Dec `json:"seigniorage_burden_target"`
    MiningIncrement         sdk.Dec `json:"mining_increment"`

    WindowShort     sdk.Int `json:"window_short"`
    WindowLong      sdk.Int `json:"window_long"`
    WindowProbation sdk.Int `json:"window_probation"`
}
```

