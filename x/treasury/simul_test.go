package treasury

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"math"
	"os"
	"strconv"
	"testing"

	"github.com/terra-project/core/types/assets"
	"github.com/terra-project/core/types/mock"
	"github.com/terra-project/core/types/util"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type SimulObj struct {
	t        int64
	TV       float64
	M        float64
	LS       float64
	MR       float64
	MRL      float64
	MRL_MA4  float64
	MRL_MA52 float64
	SB       float64
	f        float64
	w        float64
}

func loadSimulData(fileName string) (simulDatas []SimulObj) {
	csvFile, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}

	reader := csv.NewReader(bufio.NewReader(csvFile))
	_, err = reader.Read()
	if err == io.EOF {
		return
	} else if err != nil {
		panic(err)
	}

	for true {
		line, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}

		t, err := strconv.ParseInt(line[0], 10, 64)
		if err != nil {
			panic(err)
		}
		TV, err := strconv.ParseFloat(line[1], 64)
		if err != nil {
			panic(err)
		}
		M, err := strconv.ParseFloat(line[2], 64)
		if err != nil {
			panic(err)
		}
		LS, err := strconv.ParseFloat(line[3], 64)
		if err != nil {
			panic(err)
		}
		MR, err := strconv.ParseFloat(line[4], 64)
		if err != nil {
			panic(err)
		}
		MRL, err := strconv.ParseFloat(line[5], 64)
		if err != nil {
			panic(err)
		}
		MRL_MA4, err := strconv.ParseFloat(line[6], 64)
		if err != nil {
			panic(err)
		}
		MRL_MA52, err := strconv.ParseFloat(line[7], 64)
		if err != nil {
			panic(err)
		}
		SB, err := strconv.ParseFloat(line[8], 64)
		if err != nil {
			panic(err)
		}
		f, err := strconv.ParseFloat(line[9], 64)
		if err != nil {
			panic(err)
		}
		w, err := strconv.ParseFloat(line[10], 64)
		if err != nil {
			panic(err)
		}

		simulDatas = append(simulDatas, SimulObj{
			t:        t,
			TV:       TV,
			M:        M,
			LS:       LS,
			MR:       MR,
			MRL:      MRL,
			MRL_MA4:  MRL_MA4,
			MRL_MA52: MRL_MA52,
			SB:       SB,
			f:        f,
			w:        w,
		})
	}

	return
}

func TestSimulation(t *testing.T) {
	simulDatas := loadSimulData("mr-test-1.csv")

	input := createTestInput(t)

	valset := mock.NewMockValSet()
	validator := mock.NewMockValidator(sdk.ValAddress(addrs[0].Bytes()), mLunaAmt)
	valset.Validators = append(valset.Validators, validator)
	input.treasuryKeeper.valset = &valset
	params := input.treasuryKeeper.GetParams(input.ctx)
	params.WindowProbation = sdk.ZeroInt()
	input.treasuryKeeper.SetParams(input.ctx, params)

	for _, data := range simulDatas {

		currentEpoch := data.t
		currentBlockHeight := currentEpoch * util.GetBlocksPerEpoch()
		input.ctx = input.ctx.WithBlockHeight(currentBlockHeight)

		// Data of Current Epoch
		dataTV := sdk.NewDecWithPrec(int64(data.TV*math.Pow10(10)), 10)
		dataF := sdk.NewDecWithPrec(int64(data.f*math.Pow10(10)), 10)
		dataW := sdk.NewDecWithPrec(int64(data.w*math.Pow10(10)), 10)
		dataLS := sdk.NewDecWithPrec(int64(data.LS*math.Pow10(10)), 10)
		dataMRL := sdk.NewDecWithPrec(int64(data.MRL*math.Pow10(10)), 10)

		// Update current epoch data
		input.treasuryKeeper.RecordTaxProceeds(input.ctx, sdk.NewCoins(sdk.NewCoin("msdr", dataTV.MulInt64(assets.MicroUnit).Mul(dataF).TruncateInt())))
		validator.Power = dataLS.TruncateInt()

		f := input.treasuryKeeper.GetTaxRate(input.ctx, sdk.NewInt(currentEpoch))
		w := input.treasuryKeeper.GetRewardWeight(input.ctx, sdk.NewInt(currentEpoch))
		mrl := MRL(input.ctx, input.treasuryKeeper, sdk.NewInt(currentEpoch))

		// Check Invariant
		if !f.Equal(dataF) || !w.Equal(dataW) || !mrl.Equal(dataMRL) {
			fmt.Printf("target: (%v, %v, %v), result: (%v, %v, %v)\n", dataF, dataW, dataMRL, f, w, mrl)
		}

		input.ctx = input.ctx.WithBlockHeight(currentBlockHeight + util.GetBlocksPerEpoch() - 1)

		EndBlocker(input.ctx, input.treasuryKeeper)
	}

}
