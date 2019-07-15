package tx

import (
	"fmt"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/client/utils"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authtxb "github.com/cosmos/cosmos-sdk/x/auth/client/txbuilder"
	"github.com/cosmos/cosmos-sdk/x/bank"

	"github.com/terra-project/core/types/assets"
	"github.com/terra-project/core/x/treasury"
)

// ComputeFeesWithStdTx returns fee amount with given stdTx.
func ComputeFeesWithStdTx(
	cliCtx context.CLIContext,
	cdc *codec.Codec,
	tx auth.StdTx,
	gasAdjustment float64,
	gasPrices sdk.DecCoins) (fees sdk.Coins, gas uint64, err error) {

	gas = tx.Fee.Gas
	sim := (gas == 0)

	if sim {
		tx.Signatures = []auth.StdSignature{{}}
		txBytes, err := utils.GetTxEncoder(cdc)(tx)
		if err != nil {
			return nil, 0, err
		}

		_, adj, err := utils.CalculateGas(cliCtx.Query, cliCtx.Codec, txBytes, gasAdjustment)

		if err != nil {
			return nil, 0, err
		}

		gas = adj
	}

	// Computes taxes of the msgs
	taxes, err := filterMsgAndComputeTax(cliCtx, cdc, tx.Msgs)
	if err != nil {
		return nil, 0, err
	}

	fees = fees.Add(taxes)

	if !gasPrices.IsZero() {
		glDec := sdk.NewDec(int64(gas))

		// Derive the fees based on the provided gas prices, where
		// fee = ceil(gasPrice * gasLimit).
		gasFees := make(sdk.Coins, len(gasPrices))
		for i, gp := range gasPrices {
			fee := gp.Amount.Mul(glDec)
			gasFees[i] = sdk.NewCoin(gp.Denom, fee.Ceil().RoundInt())
		}

		fees = fees.Add(gasFees.Sort())
	}

	return
}

// ComputeReqParams no-lint
type ComputeReqParams struct {
	Memo          string
	ChainID       string
	AccountNumber uint64
	Sequence      uint64
	GasPrices     sdk.DecCoins
	Gas           string
	GasAdjustment string

	Msgs []sdk.Msg
}

// ComputeFee returns fee amount with given transfer, gas, gas prices, and fees amount.
func ComputeFees(
	cliCtx context.CLIContext,
	cdc *codec.Codec,
	req ComputeReqParams) (fees sdk.Coins, gas uint64, err error) {

	gasPrices := req.GasPrices
	gasAdj, err := parseFloat64(req.GasAdjustment, client.DefaultGasAdjustment)
	if err != nil {
		return nil, 0, err
	}

	if req.Gas == "0" {
		req.Gas = client.GasFlagAuto
	}

	sim, gas, err := client.ParseGas(req.Gas)
	txBldr := authtxb.NewTxBuilder(
		utils.GetTxEncoder(cdc), req.AccountNumber, req.Sequence, gas, gasAdj,
		sim, req.ChainID, req.Memo, sdk.Coins{}, req.GasPrices,
	)

	kb, _ := keys.NewKeyBaseFromHomeFlag()
	txBldr = txBldr.WithKeybase(kb)

	if sim {
		txBldr, err = utils.PrepareTxBuilder(txBldr, cliCtx)
		if err != nil {
			return nil, 0, err
		}

		txBldr, err = utils.EnrichWithGas(txBldr, cliCtx, req.Msgs)
		if err != nil {
			return nil, 0, err
		}

		gas = txBldr.Gas()
	}

	// Computes taxes of the msgs
	taxes, err := filterMsgAndComputeTax(cliCtx, cdc, req.Msgs)
	if err != nil {
		return nil, 0, err
	}

	fees = fees.Add(taxes)

	if !gasPrices.IsZero() {
		glDec := sdk.NewDec(int64(gas))

		// Derive the fees based on the provided gas prices, where
		// fee = ceil(gasPrice * gasLimit).
		gasFees := make(sdk.Coins, len(gasPrices))
		for i, gp := range gasPrices {
			fee := gp.Amount.Mul(glDec)
			gasFees[i] = sdk.NewCoin(gp.Denom, fee.Ceil().RoundInt())
		}

		fees = fees.Add(gasFees.Sort())
	}

	return
}

// filterMsgAndComputeTax computes the stability tax on MsgSend and MsgMultiSend.
func filterMsgAndComputeTax(cliCtx context.CLIContext, cdc *codec.Codec, msgs []sdk.Msg) (taxes sdk.Coins, err error) {
	taxRate, err := queryTaxRate(cliCtx, cdc)
	if err != nil {
		return nil, err
	}

	for _, msg := range msgs {
		switch msg := msg.(type) {
		case bank.MsgSend:
			tax, err := computeTax(cliCtx, cdc, taxRate, msg.Amount)
			if err != nil {
				return nil, err
			}

			taxes = taxes.Add(tax)

		case bank.MsgMultiSend:
			for _, input := range msg.Inputs {
				tax, err := computeTax(cliCtx, cdc, taxRate, input.Coins)
				if err != nil {
					return nil, err
				}

				taxes = taxes.Add(tax)
			}
		}
	}

	return
}

// computes the stability tax according to tax-rate and tax-cap
func computeTax(cliCtx context.CLIContext, cdc *codec.Codec, taxRate sdk.Dec, principal sdk.Coins) (taxes sdk.Coins, err error) {

	for _, coin := range principal {

		if coin.Denom == assets.MicroLunaDenom {
			continue
		}

		taxCap, err := queryTaxCap(cliCtx, cdc, coin.Denom)
		if err != nil {
			return nil, err
		}

		taxDue := sdk.NewDecFromInt(coin.Amount).Mul(taxRate).TruncateInt()

		// If tax due is greater than the tax cap, cap!
		if taxDue.GT(taxCap) {
			taxDue = taxCap
		}

		if taxDue.Equal(sdk.ZeroInt()) {
			continue
		}

		taxes = taxes.Add(sdk.NewCoins(sdk.NewCoin(coin.Denom, taxDue)))
	}

	return
}

func queryTaxRate(cliCtx context.CLIContext, cdc *codec.Codec) (sdk.Dec, error) {
	// Query current-epoch
	res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", treasury.QuerierRoute, treasury.QueryCurrentEpoch), nil)
	if err != nil {
		return sdk.Dec{}, err
	}

	var epochResponse treasury.QueryCurrentEpochResponse
	cdc.MustUnmarshalJSON(res, &epochResponse)
	epoch := epochResponse.CurrentEpoch

	// Query tax-rate
	res, err = cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s", treasury.QuerierRoute, treasury.QueryTaxRate, epoch.String()), nil)
	if err != nil {
		return sdk.Dec{}, err
	}

	var taxRateResponse treasury.QueryTaxRateResponse
	cdc.MustUnmarshalJSON(res, &taxRateResponse)
	taxRate := taxRateResponse.TaxRate

	return taxRate, nil
}

func queryTaxCap(cliCtx context.CLIContext, cdc *codec.Codec, denom string) (sdk.Int, error) {
	// Query tax-cap
	res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s", treasury.QuerierRoute, treasury.QueryTaxCap, denom), nil)
	if err != nil {
		return sdk.Int{}, err
	}

	var taxCapResponse treasury.QueryTaxCapResponse
	cdc.MustUnmarshalJSON(res, &taxCapResponse)
	taxCap := taxCapResponse.TaxCap

	return taxCap, nil
}

func parseFloat64(s string, defaultIfEmpty float64) (n float64, err error) {
	if len(s) == 0 {
		return defaultIfEmpty, nil
	}

	n, err = strconv.ParseFloat(s, 64)

	return
}
