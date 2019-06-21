package plan

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"

	"github.com/terra-project/core/types"
	"github.com/terra-project/core/types/assets"
	"github.com/terra-project/core/x/oracle"
)

var (
	preseedAddresses = [...]string{
		"terra1p54hc4yy2ajg67j645dn73w3378j6k05v52cnk",
		"terra1weyggr6rnggy7vdnp79pxr7vskndegcnswh8gl",
		"terra1jd2yth3zl2vez34q4m0wxngaazwy0e8csxuucd",
		"terra1upg95nlwkfkrq4hhjrn3k9s6ud0aqx36gwnlsn",
		"terra1f2z4q5kdelhfk7xq3xmxzlhp6ntrtzu659pl0s",
	}
	seedAddresses = [...]string{
		"terra1m8lvnh4zju4zcjh34rjhyyuyjk4m79536jf2tm",
		"terra1y9n2ywyu5dahtxar6k4z4jz97ynt8km4catuk6",
		"terra13cvpljjkm2huadv38dfx8jyucfda0cs6jp3sku",
		"terra1kgj2zc0qxltr5vsgycmxp7x337q5s9lvxlh694",
		"terra1zjkr4t3dglt8mzykhtvrq5wltapstlmpdly87u",
		"terra1y4umfuqfg76t8mfcff6zzx7elvy93jtp4xcdvw",
	}

	// TagUpdate230000 is tag key for update 230000
	TagUpdate230000 = "update_230000"
)

const (
	// seconds per 30 days
	secondsPerMonth = 30 * 24 * 60 * 60

	genesisUnixTime = 1556085600
)

// Update230000 update vesting schedule and oracle param
func Update230000(ctx sdk.Context, accKeeper auth.AccountKeeper, oracleKeeper oracle.Keeper) bool {

	// check update height
	if ctx.BlockHeight() != 230000 {
		return false
	}

	// update vesting schedule
	accKeeper.IterateAccounts(ctx, func(acc auth.Account) (stop bool) {
		stop = false

		vacc, ok := acc.(auth.VestingAccount)

		if !ok {
			return
		}

		gvacc, ok := vacc.(types.GradedVestingAccount)
		if ok {

			var lazyVestingSchedules []types.LazyVestingSchedule

			isPreseedAccount := false
			for _, addr := range preseedAddresses {
				if addr == gvacc.GetAddress().String() {
					lazyVestingSchedules = updatePreseedSchedules(gvacc)
					isPreseedAccount = true
					break
				}
			}

			isSeedAccount := false
			for _, addr := range seedAddresses {
				if addr == gvacc.GetAddress().String() {
					lazyVestingSchedules = updateSeedSchedules(gvacc)
					isSeedAccount = true
					break
				}
			}

			if !isPreseedAccount && !isSeedAccount {
				// update to LazyGradedVestingAccount
				vestingSchedules := gvacc.GetVestingSchedules()

				for _, vs := range vestingSchedules {
					var lazySchedules []types.LazySchedule
					for _, s := range vs.Schedules {
						lazySchedules = append(lazySchedules, types.NewLazySchedule(s.GetCliff(), s.GetCliff()+secondsPerMonth, s.GetRatio()))
					}

					lazyVestingSchedule := types.NewLazyVestingSchedule(vs.GetDenom(), lazySchedules)
					lazyVestingSchedules = append(lazyVestingSchedules, lazyVestingSchedule)
				}
			}

			for _, lvs := range lazyVestingSchedules {
				if !lvs.IsValid() {
					panic(fmt.Sprintf("not valid schdule: %v\n %v", gvacc, lvs))
				}
			}

			baseAccount := &auth.BaseAccount{
				Address:       gvacc.GetAddress(),
				PubKey:        gvacc.GetPubKey(),
				Coins:         gvacc.GetCoins().Sort(),
				AccountNumber: gvacc.GetAccountNumber(),
				Sequence:      gvacc.GetSequence(),
			}

			baseVestingAcc := &auth.BaseVestingAccount{
				BaseAccount:      baseAccount,
				OriginalVesting:  gvacc.GetOriginalVesting(),
				DelegatedFree:    gvacc.GetDelegatedFree(),
				DelegatedVesting: gvacc.GetDelegatedVesting(),
				EndTime:          gvacc.GetEndTime(),
			}

			lazyAccount := types.BaseLazyGradedVestingAccount{
				BaseVestingAccount:   baseVestingAcc,
				LazyVestingSchedules: lazyVestingSchedules,
			}

			accKeeper.SetAccount(ctx, lazyAccount)
		}

		return
	})

	// update oracle reward band param
	oracleParams := oracleKeeper.GetParams(ctx)
	oracleParams.OracleRewardBand = sdk.NewDecWithPrec(2, 2) // 2%
	oracleKeeper.SetParams(ctx, oracleParams)

	return true
}

// preseed account does not have any other schedules
func updatePreseedSchedules(gvacc types.GradedVestingAccount) []types.LazyVestingSchedule {

	vestingSchedules := gvacc.GetVestingSchedules()
	if len(vestingSchedules) != 1 || vestingSchedules[0].GetDenom() != assets.MicroLunaDenom {
		panic(fmt.Sprintf("Invalid Preseed Account: %v", gvacc))
	}

	vestingSchedule := vestingSchedules[0]
	if len(vestingSchedule.Schedules) != 4 {
		panic(fmt.Sprintf("Invalid Preseed Account: %v", gvacc))
	}

	// strict preseed account check
	for _, s := range vestingSchedule.Schedules {
		if s.GetCliff() == 1558677600 && s.GetRatio().Equal(sdk.NewDecWithPrec(10, 2)) {
			continue
		} else if s.GetCliff() == 1561356000 && s.GetRatio().Equal(sdk.NewDecWithPrec(10, 2)) {
			continue
		} else if s.GetCliff() == 1563948000 && s.GetRatio().Equal(sdk.NewDecWithPrec(10, 2)) {
			continue
		} else if s.GetCliff() == 1587708000 && s.GetRatio().Equal(sdk.NewDecWithPrec(70, 2)) {
			continue
		} else {
			panic(fmt.Sprintf("Invalid Preseed Account: %v", gvacc))
		}
	}

	var lazyVestingSchedules []types.LazyVestingSchedule
	var lazyVestingSchedule types.LazyVestingSchedule
	var lazySchedules []types.LazySchedule

	genesisTime := time.Unix(genesisUnixTime, 0)
	lazySchedules = append(lazySchedules,
		types.NewLazySchedule(
			genesisTime.AddDate(0, 1, 0).Unix(), genesisTime.AddDate(0, 2, 0).Unix(), sdk.NewDecWithPrec(10, 2)),
		types.NewLazySchedule(
			genesisTime.AddDate(0, 2, 0).Unix(), genesisTime.AddDate(0, 12, 0).Unix(), sdk.NewDecWithPrec(27, 2)),
		types.NewLazySchedule(
			genesisTime.AddDate(0, 12, 0).Unix(), genesisTime.AddDate(0, 17, 0).Unix(), sdk.NewDecWithPrec(48, 2)),
		types.NewLazySchedule(
			genesisTime.AddDate(0, 17, 0).Unix(), genesisTime.AddDate(0, 18, 0).Unix(), sdk.NewDecWithPrec(15, 2)),
	)

	lazyVestingSchedule = types.NewLazyVestingSchedule(assets.MicroLunaDenom, lazySchedules)
	lazyVestingSchedules = append(lazyVestingSchedules, lazyVestingSchedule)

	return lazyVestingSchedules
}

// only terra1y9n2ywyu5dahtxar6k4z4jz97ynt8km4catuk6 account has difference vesting schedule
// except that all seed account does not have any other schedules
func updateSeedSchedules(gvacc types.GradedVestingAccount) []types.LazyVestingSchedule {
	vestingSchedules := gvacc.GetVestingSchedules()
	if len(vestingSchedules) != 1 || vestingSchedules[0].GetDenom() != assets.MicroLunaDenom {
		panic(fmt.Sprintf("Invalid Seed Account: %v", gvacc))
	}

	vestingSchedule := vestingSchedules[0]
	ratio := sdk.OneDec()

	// strict seed account check
	if gvacc.GetAddress().String() == "terra1y9n2ywyu5dahtxar6k4z4jz97ynt8km4catuk6" {
		ratio = sdk.NewDecWithPrec(467, 3)

		if len(vestingSchedule.Schedules) != 5 {
			panic(fmt.Sprintf("Invalid Seed Account: %v", gvacc))
		}

		for _, s := range vestingSchedule.Schedules {
			if s.GetCliff() == 1558677600 && s.GetRatio().Equal(sdk.NewDecWithPrec(47, 3)) {
				continue
			} else if s.GetCliff() == 1561356000 && s.GetRatio().Equal(sdk.NewDecWithPrec(47, 3)) {
				continue
			} else if s.GetCliff() == 1563948000 && s.GetRatio().Equal(sdk.NewDecWithPrec(47, 3)) {
				continue
			} else if s.GetCliff() == 1582524000 && s.GetRatio().Equal(sdk.NewDecWithPrec(326, 3)) {
				continue
			} else if s.GetCliff() == 1603519200 && s.GetRatio().Equal(sdk.NewDecWithPrec(533, 3)) {
				continue
			} else {
				panic(fmt.Sprintf("Invalid Seed Account: %v", gvacc))
			}
		}
	} else {
		if len(vestingSchedule.Schedules) != 4 {
			panic(fmt.Sprintf("Invalid Seed Account: %v", gvacc))
		}

		for _, s := range vestingSchedule.Schedules {
			if s.GetCliff() == 1558677600 && s.GetRatio().Equal(sdk.NewDecWithPrec(10, 2)) {
				continue
			} else if s.GetCliff() == 1561356000 && s.GetRatio().Equal(sdk.NewDecWithPrec(10, 2)) {
				continue
			} else if s.GetCliff() == 1563948000 && s.GetRatio().Equal(sdk.NewDecWithPrec(10, 2)) {
				continue
			} else if (s.GetCliff() == 1582524000 || s.GetCliff() == 1579845600) && s.GetRatio().Equal(sdk.NewDecWithPrec(70, 2)) {
				continue
			} else {
				panic(fmt.Sprintf("Invalid Seed Account: %v", gvacc))
			}
		}
	}

	var lazyVestingSchedules []types.LazyVestingSchedule
	var lazyVestingSchedule types.LazyVestingSchedule
	var lazySchedules []types.LazySchedule

	genesisTime := time.Unix(genesisUnixTime, 0)
	lazySchedules = append(lazySchedules,
		types.NewLazySchedule(
			genesisTime.AddDate(0, 1, 0).Unix(), genesisTime.AddDate(0, 2, 0).Unix(), ratio.Mul(sdk.NewDecWithPrec(10, 2))),
		types.NewLazySchedule(
			genesisTime.AddDate(0, 2, 0).Unix(), genesisTime.AddDate(0, 10, 0).Unix(), ratio.Mul(sdk.NewDecWithPrec(30, 2))),
		types.NewLazySchedule(
			genesisTime.AddDate(0, 10, 0).Unix(), genesisTime.AddDate(0, 13, 0).Unix(), ratio.Mul(sdk.NewDecWithPrec(60, 2))),
	)

	if !ratio.Equal(sdk.OneDec()) {
		lazySchedules = append(lazySchedules,
			types.NewLazySchedule(
				genesisTime.AddDate(0, 18, 0).Unix(), genesisTime.AddDate(0, 19, 0).Unix(), sdk.OneDec().Sub(ratio)),
		)
	}

	lazyVestingSchedule = types.NewLazyVestingSchedule(assets.MicroLunaDenom, lazySchedules)
	lazyVestingSchedules = append(lazyVestingSchedules, lazyVestingSchedule)

	return lazyVestingSchedules
}
