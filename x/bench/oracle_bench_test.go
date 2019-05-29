package bench

import (
	"encoding/hex"
	"testing"

	"github.com/terra-project/core/types/assets"
	"github.com/terra-project/core/x/oracle"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func BenchmarkOracleFeedVotePerBlock(b *testing.B) {
	input := createTestInput()

	defaultOracleParams := oracle.DefaultParams()
	defaultOracleParams.VotePeriod = 1
	input.oracleKeeper.SetParams(input.ctx, defaultOracleParams)

	h := oracle.NewHandler(input.oracleKeeper)

	denoms := []string{
		assets.MicroSDRDenom,
		assets.MicroKRWDenom,
		assets.MicroUSDDenom,
		assets.MicroCNYDenom,
		assets.MicroJPYDenom,
		assets.MicroGBPDenom,
		assets.MicroEURDenom,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ctx := input.ctx.WithBlockHeight(int64(i))

		for j := 0; j < numOfValidators; j++ {
			for d := 0; d < len(denoms); d++ {
				salt := string(j)
				bz, err := oracle.VoteHash(salt, sdk.OneDec(), denoms[d], sdk.ValAddress(addrs[j]))
				if err != nil {
					panic(err)
				}

				voteMsg := oracle.NewMsgPriceFeed(hex.EncodeToString(bz), "", denoms[d], addrs[j], sdk.ValAddress(addrs[j]), sdk.ZeroDec())
				res := h(ctx, voteMsg)
				if !res.IsOK() {
					panic(res.Log)
				}

				voteMsg = oracle.NewMsgPriceFeed("", salt, denoms[d], addrs[j], sdk.ValAddress(addrs[j]), sdk.OneDec())
				res = h(ctx.WithBlockHeight(1), voteMsg)
				if !res.IsOK() {
					panic(res.Log)
				}

			}
		}

		oracle.EndBlocker(ctx, input.oracleKeeper)
	}
}
