package utils

import (
	"fmt"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"

	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	core "github.com/terra-money/core/types"

	marketexported "github.com/terra-money/core/x/market/exported"
	msgauthexported "github.com/terra-money/core/x/msgauth/exported"
	treasuryexported "github.com/terra-money/core/x/treasury/exported"
	wasmexported "github.com/terra-money/core/x/wasm/exported"
)

type (
	// EstimateFeeReq defines a tx encoding request.
	EstimateFeeReq struct {
		Tx            auth.StdTx   `json:"tx"`
		GasAdjustment string       `json:"gas_adjustment"`
		GasPrices     sdk.DecCoins `json:"gas_prices"`
	}

	// EstimateFeeResp defines a tx encoding response.
	EstimateFeeResp struct {
		Fees sdk.Coins `json:"fees"`
		Gas  uint64    `json:"gas"`
	}
)

// String implements fmt.Stringer interface
func (r EstimateFeeResp) String() string {
	return fmt.Sprintf(`EstimateFeeResp
	fees: %s,
	gas:  %d`,
		r.Fees, r.Gas)
}

// ComputeFeesWithStdTx returns fee amount with given stdTx.
func ComputeFeesWithStdTx(
	cliCtx context.CLIContext,
	tx auth.StdTx,
	gasAdjustment float64,
	gasPrices sdk.DecCoins) (fees sdk.Coins, gas uint64, err error) {

	gas = tx.Fee.Gas
	sim := (gas == 0)

	if sim {
		tx.Signatures = []auth.StdSignature{}

		signers := make(map[string]bool)
		for _, msg := range tx.Msgs {
			for _, signer := range msg.GetSigners() {
				if _, ok := signers[signer.String()]; !ok {
					signers[signer.String()] = true
					tx.Signatures = append(tx.Signatures, auth.StdSignature{})
				}
			}
		}

		txBytes, err := utils.GetTxEncoder(cliCtx.Codec)(tx)
		if err != nil {
			return nil, 0, err
		}

		_, adj, err := utils.CalculateGas(cliCtx.QueryWithData, cliCtx.Codec, txBytes, gasAdjustment)

		if err != nil {
			return nil, 0, err
		}

		gas = adj
	}

	// Computes taxes of the msgs
	taxes, err := filterMsgAndComputeTax(cliCtx, tx.Msgs)
	if err != nil {
		return nil, 0, err
	}

	fees = fees.Add(taxes...)

	if !gasPrices.IsZero() {
		glDec := sdk.NewDec(int64(gas))

		// Derive the fees based on the provided gas prices, where
		// fee = ceil(gasPrice * gasLimit).
		gasFees := make(sdk.Coins, len(gasPrices))
		for i, gp := range gasPrices {
			fee := gp.Amount.Mul(glDec)
			gasFees[i] = sdk.NewCoin(gp.Denom, fee.Ceil().RoundInt())
		}

		fees = fees.Add(gasFees.Sort()...)
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

// ComputeFees returns fee amount with given transfer, gas, gas prices, and fees amount.
func ComputeFees(
	cliCtx context.CLIContext,
	req ComputeReqParams) (fees sdk.Coins, gas uint64, err error) {

	gasPrices := req.GasPrices
	gasAdj, err := ParseFloat64(req.GasAdjustment, flags.DefaultGasAdjustment)
	if err != nil {
		return nil, 0, err
	}

	if req.Gas == "0" {
		req.Gas = flags.GasFlagAuto
	}

	sim, gas, err := flags.ParseGas(req.Gas)
	txBldr := auth.NewTxBuilder(
		utils.GetTxEncoder(cliCtx.Codec), req.AccountNumber, req.Sequence, flags.DefaultGasLimit, gasAdj,
		sim, req.ChainID, req.Memo, sdk.Coins{}, req.GasPrices,
	)

	if sim {
		txBldr, err = utils.EnrichWithGas(txBldr, cliCtx, req.Msgs)
		if err != nil {
			return nil, 0, err
		}

		gas = txBldr.Gas()
	}

	// Computes taxes of the msgs
	taxes, err := filterMsgAndComputeTax(cliCtx, req.Msgs)
	if err != nil {
		return nil, 0, err
	}

	fees = fees.Add(taxes...)

	if !gasPrices.IsZero() {
		glDec := sdk.NewDec(int64(gas))

		// Derive the fees based on the provided gas prices, where
		// fee = ceil(gasPrice * gasLimit).
		gasFees := make(sdk.Coins, len(gasPrices))
		for i, gp := range gasPrices {
			fee := gp.Amount.Mul(glDec)
			gasFees[i] = sdk.NewCoin(gp.Denom, fee.Ceil().RoundInt())
		}

		fees = fees.Add(gasFees.Sort()...)
	}

	return
}

// filterMsgAndComputeTax computes the stability tax on MsgSend and MsgMultiSend.
func filterMsgAndComputeTax(cliCtx context.CLIContext, msgs []sdk.Msg) (taxes sdk.Coins, err error) {
	taxRate, err := queryTaxRate(cliCtx)
	if err != nil {
		return nil, err
	}

	for _, msg := range msgs {
		switch msg := msg.(type) {
		case bank.MsgSend:
			tax, err := computeTax(cliCtx, taxRate, msg.Amount)
			if err != nil {
				return nil, err
			}

			taxes = taxes.Add(tax...)

		case bank.MsgMultiSend:
			for _, input := range msg.Inputs {
				tax, err := computeTax(cliCtx, taxRate, input.Coins)
				if err != nil {
					return nil, err
				}

				taxes = taxes.Add(tax...)
			}

		case msgauthexported.MsgExecAuthorized:
			tax, err := filterMsgAndComputeTax(cliCtx, msg.Msgs)
			if err != nil {
				return nil, err
			}

			taxes = taxes.Add(tax...)

		case marketexported.MsgSwapSend:
			tax, err := computeTax(cliCtx, taxRate, sdk.NewCoins(msg.OfferCoin))
			if err != nil {
				return nil, err
			}

			taxes = taxes.Add(tax...)

		case wasmexported.MsgInstantiateContract:
			tax, err := computeTax(cliCtx, taxRate, msg.InitCoins)
			if err != nil {
				return nil, err
			}

			taxes = taxes.Add(tax...)

		case wasmexported.MsgExecuteContract:
			tax, err := computeTax(cliCtx, taxRate, msg.Coins)
			if err != nil {
				return nil, err
			}

			taxes = taxes.Add(tax...)
		}
	}

	return
}

// computes the stability tax according to tax-rate and tax-cap
func computeTax(cliCtx context.CLIContext, taxRate sdk.Dec, principal sdk.Coins) (taxes sdk.Coins, err error) {

	for _, coin := range principal {

		if coin.Denom == core.MicroLunaDenom {
			continue
		}

		taxCap, err := queryTaxCap(cliCtx, coin.Denom)
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

		taxes = taxes.Add(sdk.NewCoin(coin.Denom, taxDue))
	}

	return
}

func queryTaxRate(cliCtx context.CLIContext) (sdk.Dec, error) {

	// Query tax-rate
	res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", treasuryexported.QuerierRoute, treasuryexported.QueryTaxRate), nil)
	if err != nil {
		return sdk.Dec{}, err
	}

	var taxRate sdk.Dec
	cliCtx.Codec.MustUnmarshalJSON(res, &taxRate)
	return taxRate, nil
}

func queryTaxCap(cliCtx context.CLIContext, denom string) (sdk.Int, error) {
	// Query tax-cap

	params := treasuryexported.NewQueryTaxCapParams(denom)
	bz := cliCtx.Codec.MustMarshalJSON(params)
	res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", treasuryexported.QuerierRoute, treasuryexported.QueryTaxCap), bz)
	if err != nil {
		return sdk.Int{}, err
	}

	var taxCap sdk.Int
	cliCtx.Codec.MustUnmarshalJSON(res, &taxCap)

	return taxCap, nil
}

// ParseFloat64 parses string to float64
func ParseFloat64(s string, defaultIfEmpty float64) (n float64, err error) {
	if len(s) == 0 {
		return defaultIfEmpty, nil
	}

	n, err = strconv.ParseFloat(s, 64)

	return
}
