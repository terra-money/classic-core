package utils

import (
	"context"

	"github.com/spf13/pflag"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/migrations/legacytx"
	"github.com/cosmos/cosmos-sdk/x/authz"

	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	wasmexported "github.com/CosmWasm/wasmd/x/wasm"
	marketexported "github.com/classic-terra/core/v2/x/market/exported"
	treasuryexported "github.com/classic-terra/core/v2/x/treasury/exported"
)

type (
	// EstimateFeeResp defines a tx fee estimation response
	EstimateFeeResp struct {
		Fee legacytx.StdFee `json:"fee" yaml:"fee"`
	}
)

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
	clientCtx client.Context, flagSet *pflag.FlagSet, msgs ...sdk.Msg,
) (*legacytx.StdFee, error) {
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

	// Computes taxes of the msgs
	taxes, err := FilterMsgAndComputeTax(clientCtx, msgs...)
	if err != nil {
		return nil, err
	}

	fees := txf.Fees().Add(taxes...)
	gasPrices := txf.GasPrices()

	if !gasPrices.IsZero() {
		glDec := sdk.NewDec(int64(gas))
		adjustment := sdk.NewDecWithPrec(int64(txf.GasAdjustment())*100, 2)

		if adjustment.LT(sdk.OneDec()) {
			adjustment = sdk.OneDec()
		}

		// Derive the fees based on the provided gas prices, where
		// fee = ceil(gasPrice * gasLimit).
		gasFees := make(sdk.Coins, len(gasPrices))
		for i, gp := range gasPrices {
			fee := gp.Amount.Mul(glDec).Mul(adjustment)
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

		case *authz.MsgExec:
			messages, err := msg.GetMessages()
			if err != nil {
				panic(err)
			}

			tax, err := FilterMsgAndComputeTax(clientCtx, messages...)
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
			tax, err := computeTax(clientCtx, taxRate, msg.Funds)
			if err != nil {
				return nil, err
			}

			taxes = taxes.Add(tax...)

		case *wasmexported.MsgInstantiateContract2:
			tax, err := computeTax(clientCtx, taxRate, msg.Funds)
			if err != nil {
				return nil, err
			}

			taxes = taxes.Add(tax...)

		case *wasmexported.MsgExecuteContract:
			tax, err := computeTax(clientCtx, taxRate, msg.Funds)
			if err != nil {
				return nil, err
			}

			taxes = taxes.Add(tax...)
		}
	}

	return taxes, nil
}

// computes the stability tax according to tax-rate and tax-cap
func computeTax(clientCtx client.Context, taxRate sdk.Dec, principal sdk.Coins) (taxes sdk.Coins, err error) {
	for _, coin := range principal {

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
	queryClient := treasuryexported.NewQueryClient(clientCtx)

	res, err := queryClient.TaxRate(context.Background(), &treasuryexported.QueryTaxRateRequest{})
	if err != nil {
		return sdk.ZeroDec(), err
	}
	return res.TaxRate, err
}

func queryTaxCap(clientCtx client.Context, denom string) (sdk.Int, error) {
	queryClient := treasuryexported.NewQueryClient(clientCtx)

	res, err := queryClient.TaxCap(context.Background(), &treasuryexported.QueryTaxCapRequest{Denom: denom})
	if err != nil {
		return sdk.NewInt(0), err
	}
	return res.TaxCap, err
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
