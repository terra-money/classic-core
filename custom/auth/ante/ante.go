package ante

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	cosmosante "github.com/cosmos/cosmos-sdk/x/auth/ante"
	"github.com/cosmos/cosmos-sdk/x/auth/keeper"
	"github.com/cosmos/cosmos-sdk/x/auth/signing"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
)

// NewAnteHandler returns an AnteHandler that checks and increments sequence
// numbers, checks signatures & account numbers, and deducts fees from the first
// signer.
func NewAnteHandler(
	ak keeper.AccountKeeper, bankKeeper types.BankKeeper, oracleKeeper OracleKeeper, treasuryKeeper TreasuryKeeper,
	sigGasConsumer cosmosante.SignatureVerificationGasConsumer, signModeHandler signing.SignModeHandler) sdk.AnteHandler {
	return sdk.ChainAnteDecorators(
		cosmosante.NewSetUpContextDecorator(), // outermost AnteDecorator. SetUpContext must be called first
		cosmosante.NewRejectExtensionOptionsDecorator(),
		NewSpammingPreventionDecorator(oracleKeeper), // spamming prevention
		NewTaxFeeDecorator(treasuryKeeper),           // mempool gas fee validation & record tax proceeds
		cosmosante.NewValidateBasicDecorator(),
		cosmosante.TxTimeoutHeightDecorator{},
		cosmosante.NewValidateMemoDecorator(ak),
		cosmosante.NewConsumeGasForTxSizeDecorator(ak),
		cosmosante.NewRejectFeeGranterDecorator(),
		cosmosante.NewSetPubKeyDecorator(ak), // SetPubKeyDecorator must be called before all signature verification decorators
		cosmosante.NewValidateSigCountDecorator(ak),
		cosmosante.NewDeductFeeDecorator(ak, bankKeeper),
		cosmosante.NewSigGasConsumeDecorator(ak, sigGasConsumer),
		cosmosante.NewSigVerificationDecorator(ak, signModeHandler),
		cosmosante.NewIncrementSequenceDecorator(ak),
	)
}
