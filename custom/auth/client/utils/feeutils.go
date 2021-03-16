package utils

import (
	"fmt"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/cosmos/cosmos-sdk/x/auth/legacy/legacytx"
	"github.com/spf13/pflag"

	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	core "github.com/terra-project/core/types"

	marketexported "github.com/terra-project/core/x/market/exported"
	msgauthexported "github.com/terra-project/core/x/msgauth/exported"
	treasuryexported "github.com/terra-project/core/x/treasury/exported"
	wasmexported "github.com/terra-project/core/x/wasm/exported"
)

type (
	// EstimateFeeReq defines a tx fee estimation request.
	EstimateFeeReq struct {
		BaseReq rest.BaseReq `json:"base_req" yaml:"base_req"`
		Msg     []sdk.Msg    `json:"msgs" yaml:"msgs"`
	}

	// EstimateFeeResp defines a tx fee estimation response
	EstimateFeeResp struct {
		Fee legacytx.StdFee `json:"fee" yaml:"fee"`
	}
)

// ComputeFeesWithBaseReq returns fee amount with given stdTx.
func ComputeFeesWithBaseReq(
	clientCtx client.Context, br rest.BaseReq, msgs ...sdk.Msg) (*legacytx.StdFee, error) {

	gasSetting, err := flags.ParseGasSetting(br.Gas)
	if err != nil {
		return nil, err
	}

	gasAdj, err := ParseFloat64(br.GasAdjustment, flags.DefaultGasAdjustment)
	if err != nil {
		return nil, err
	}

	gas := gasSetting.Gas
	if gasSetting.Simulate {
		txf := tx.Factory{}.
			WithFees(br.Fees.String()).
			WithGasPrices(br.GasPrices.String()).
			WithGas(gasSetting.Gas).
			WithGasAdjustment(gasAdj).
			WithMemo(br.Memo).
			WithChainID(br.ChainID).
			WithSimulateAndExecute(br.Simulate).
			WithTxConfig(clientCtx.TxConfig).
			WithTimeoutHeight(br.TimeoutHeight)

		_, adj, err := tx.CalculateGas(clientCtx.QueryWithData, txf, msgs...)

		if err != nil {
			return nil, err
		}

		gas = adj
	}

	// Computes taxes of the msgs
	taxes, err := FilterMsgAndComputeTax(clientCtx, msgs...)
	if err != nil {
		return nil, err
	}

	fees := br.Fees.Add(taxes...)
	gasPrices := br.GasPrices

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

	return &legacytx.StdFee{
		Amount: fees,
		Gas:    gas,
	}, nil
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

// ComputeFeesWithCmd returns fee amount with cli options.
func ComputeFeesWithCmd(
	clientCtx client.Context, flagSet *pflag.FlagSet, msgs ...sdk.Msg) (*legacytx.StdFee, error) {
	txf := tx.NewFactoryCLI(clientCtx, flagSet)

	gas := txf.Gas()
	if txf.SimulateAndExecute() {
		_, adj, err := tx.CalculateGas(clientCtx.QueryWithData, txf, msgs...)
		if err != nil {
			return nil, err
		}

		gas = adj
	}

	// Computes taxes of the msgs
	taxes, err := FilterMsgAndComputeTax(clientCtx, msgs...)
	if err != nil {
		return nil, err
	}

	fees := txf.Fees().Add(taxes...)
	gasPrices := txf.GasPrices()

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

	return &legacytx.StdFee{
		Amount: fees,
		Gas:    gas,
	}, nil
}

// FilterMsgAndComputeTax computes the stability tax on MsgSend and MsgMultiSend.
func FilterMsgAndComputeTax(clientCtx client.Context, msgs ...sdk.Msg) (taxes sdk.Coins, err error) {
	taxRate, err := queryTaxRate(clientCtx)
	if err != nil {
		return nil, err
	}

	for _, msg := range msgs {
		switch msg := msg.(type) {
		case *banktypes.MsgSend:
			tax, err := computeTax(clientCtx, taxRate, msg.Amount)
			if err != nil {
				return nil, err
			}

			taxes = taxes.Add(tax...)

		case *banktypes.MsgMultiSend:
			for _, input := range msg.Inputs {
				tax, err := computeTax(clientCtx, taxRate, input.Coins)
				if err != nil {
					return nil, err
				}

				taxes = taxes.Add(tax...)
			}

		case *msgauthexported.MsgExecAuthorized:
			tax, err := FilterMsgAndComputeTax(clientCtx, msg.GetMsgs()...)
			if err != nil {
				return nil, err
			}

			taxes = taxes.Add(tax...)

		case *marketexported.MsgSwapSend:
			tax, err := computeTax(clientCtx, taxRate, sdk.NewCoins(msg.OfferCoin))
			if err != nil {
				return nil, err
			}

			taxes = taxes.Add(tax...)

		case *wasmexported.MsgInstantiateContract:
			tax, err := computeTax(clientCtx, taxRate, msg.InitCoins)
			if err != nil {
				return nil, err
			}

			taxes = taxes.Add(tax...)

		case *wasmexported.MsgExecuteContract:
			tax, err := computeTax(clientCtx, taxRate, msg.Coins)
			if err != nil {
				return nil, err
			}

			taxes = taxes.Add(tax...)
		}
	}

	return
}

// computes the stability tax according to tax-rate and tax-cap
func computeTax(clientCtx client.Context, taxRate sdk.Dec, principal sdk.Coins) (taxes sdk.Coins, err error) {

	for _, coin := range principal {

		if coin.Denom == core.MicroLunaDenom {
			continue
		}

		taxCap, err := queryTaxCap(clientCtx, coin.Denom)
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

func queryTaxRate(clientCtx client.Context) (sdk.Dec, error) {

	// Query tax-rate
	res, _, err := clientCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", treasuryexported.QuerierRoute, treasuryexported.QueryTaxRate), nil)
	if err != nil {
		return sdk.Dec{}, err
	}

	var taxRate sdk.Dec
	clientCtx.LegacyAmino.MustUnmarshalJSON(res, &taxRate)
	return taxRate, nil
}

func queryTaxCap(clientCtx client.Context, denom string) (sdk.Int, error) {
	// Query tax-cap

	params := treasuryexported.NewQueryTaxCapParams(denom)
	bz := clientCtx.LegacyAmino.MustMarshalJSON(params)
	res, _, err := clientCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", treasuryexported.QuerierRoute, treasuryexported.QueryTaxCap), bz)
	if err != nil {
		return sdk.Int{}, err
	}

	var taxCap sdk.Int
	clientCtx.LegacyAmino.MustUnmarshalJSON(res, &taxCap)

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
