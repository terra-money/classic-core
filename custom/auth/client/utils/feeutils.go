package utils

import (
	"strconv"

	"github.com/spf13/pflag"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/cosmos/cosmos-sdk/x/auth/legacy/legacytx"
)

type (
	// EstimateFeeReq defines a tx fee estimation request.
	EstimateFeeReq struct {
		BaseReq rest.BaseReq `json:"base_req" yaml:"base_req"`
		Msgs    []sdk.Msg    `json:"msgs" yaml:"msgs"`
	}

	// EstimateFeeResp defines a tx fee estimation response
	EstimateFeeResp struct {
		Fee legacytx.StdFee `json:"fee" yaml:"fee"`
	}
)

var _ codectypes.UnpackInterfacesMessage = EstimateFeeReq{}

// UnpackInterfaces implements the UnpackInterfacesMessage interface.
func (m EstimateFeeReq) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	for _, m := range m.Msgs {
		err := codectypes.UnpackInterfaces(m, unpacker)
		if err != nil {
			return err
		}
	}

	return nil
}

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
			WithAccountNumber(br.AccountNumber).
			WithSequence(br.Sequence).
			WithMemo(br.Memo).
			WithChainID(br.ChainID).
			WithSimulateAndExecute(br.Simulate || gasSetting.Simulate).
			WithTxConfig(clientCtx.TxConfig).
			WithTimeoutHeight(br.TimeoutHeight).
			WithAccountRetriever(clientCtx.AccountRetriever)

		// Prepare AccountNumber & SequenceNumber when not given
		clientCtx.FromAddress, err = sdk.AccAddressFromBech32(br.From)
		if err != nil {
			return nil, err
		}

		txf, err := prepareFactory(clientCtx, txf)
		if err != nil {
			return nil, err
		}

		_, adj, err := tx.CalculateGas(clientCtx, txf, msgs...)
		if err != nil {
			return nil, err
		}

		gas = adj
	}

	fees := br.Fees
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
		txf, err := prepareFactory(clientCtx, txf)
		if err != nil {
			return nil, err
		}

		_, adj, err := tx.CalculateGas(clientCtx, txf, msgs...)
		if err != nil {
			return nil, err
		}

		gas = adj
	}

	fees := txf.Fees()
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

// ParseFloat64 parses string to float64
func ParseFloat64(s string, defaultIfEmpty float64) (n float64, err error) {
	if len(s) == 0 {
		return defaultIfEmpty, nil
	}

	n, err = strconv.ParseFloat(s, 64)

	return
}

// prepareFactory ensures the account defined by ctx.GetFromAddress() exists and
// if the account number and/or the account sequence number are zero (not set),
// they will be queried for and set on the provided Factory. A new Factory with
// the updated fields will be returned.
func prepareFactory(clientCtx client.Context, txf tx.Factory) (tx.Factory, error) {
	from := clientCtx.GetFromAddress()

	if err := txf.AccountRetriever().EnsureExists(clientCtx, from); err != nil {
		return txf, err
	}

	initNum, initSeq := txf.AccountNumber(), txf.Sequence()
	if initNum == 0 || initSeq == 0 {
		num, seq, err := txf.AccountRetriever().GetAccountNumberSequence(clientCtx, from)
		if err != nil {
			return txf, err
		}

		if initNum == 0 {
			txf = txf.WithAccountNumber(num)
		}

		if initSeq == 0 {
			txf = txf.WithSequence(seq)
		}
	}

	return txf, nil
}
