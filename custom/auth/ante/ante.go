package ante

import (
	ibcante "github.com/cosmos/ibc-go/v6/modules/core/ante"
	ibckeeper "github.com/cosmos/ibc-go/v6/modules/core/keeper"

	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	"github.com/cosmos/cosmos-sdk/x/auth/signing"
	distributionkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	govkeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
)

// HandlerOptions are the options required for constructing a default SDK AnteHandler.
type HandlerOptions struct {
	AccountKeeper          ante.AccountKeeper
	BankKeeper             BankKeeper
	ExtensionOptionChecker ante.ExtensionOptionChecker
	FeegrantKeeper         ante.FeegrantKeeper
	OracleKeeper           OracleKeeper
	TreasuryKeeper         TreasuryKeeper
	SignModeHandler        signing.SignModeHandler
	SigGasConsumer         ante.SignatureVerificationGasConsumer
	TxFeeChecker           ante.TxFeeChecker
	IBCKeeper              ibckeeper.Keeper
	DistributionKeeper     distributionkeeper.Keeper
	GovKeeper              govkeeper.Keeper
	WasmConfig             *wasmtypes.WasmConfig
	TXCounterStoreKey      storetypes.StoreKey
}

// NewAnteHandler returns an AnteHandler that checks and increments sequence
// numbers, checks signatures & account numbers, and deducts fees from the first
// signer.
func NewAnteHandler(options HandlerOptions) (sdk.AnteHandler, error) {
	if options.AccountKeeper == nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrLogic, "account keeper is required for ante builder")
	}

	if options.BankKeeper == nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrLogic, "bank keeper is required for ante builder")
	}

	if options.OracleKeeper == nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrLogic, "oracle keeper is required for ante builder")
	}

	if options.TreasuryKeeper == nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrLogic, "treasury keeper is required for ante builder")
	}

	if options.SignModeHandler == nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrLogic, "sign mode handler is required for ante builder")
	}

	if options.WasmConfig == nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrLogic, "wasm config is required for ante builder")
	}

	if options.TXCounterStoreKey == nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrLogic, "tx counter key is required for ante builder")
	}

	return sdk.ChainAnteDecorators(
		ante.NewSetUpContextDecorator(), // outermost AnteDecorator. SetUpContext must be called first
		wasmkeeper.NewLimitSimulationGasDecorator(options.WasmConfig.SimulationGasLimit),
		wasmkeeper.NewCountTXDecorator(options.TXCounterStoreKey),
		ante.NewExtensionOptionsDecorator(options.ExtensionOptionChecker),
		ante.NewValidateBasicDecorator(),
		ante.NewTxTimeoutHeightDecorator(),
		ante.NewValidateMemoDecorator(options.AccountKeeper),
		// SpammingPreventionDecorator prevents spamming oracle vote tx attempts at same height
		NewSpammingPreventionDecorator(options.OracleKeeper),
		// MinInitialDepositDecorator prevents submitting governance proposal low initial deposit
		NewMinInitialDepositDecorator(options.GovKeeper, options.TreasuryKeeper),
		ante.NewConsumeGasForTxSizeDecorator(options.AccountKeeper),
		NewFeeDecorator(options.AccountKeeper, options.BankKeeper, options.FeegrantKeeper, options.TreasuryKeeper),
		ante.NewSetPubKeyDecorator(options.AccountKeeper), // SetPubKeyDecorator must be called before all signature verification decorators
		ante.NewValidateSigCountDecorator(options.AccountKeeper),
		ante.NewSigGasConsumeDecorator(options.AccountKeeper, options.SigGasConsumer),
		ante.NewSigVerificationDecorator(options.AccountKeeper, options.SignModeHandler),
		ante.NewIncrementSequenceDecorator(options.AccountKeeper),
		ibcante.NewRedundantRelayDecorator(&options.IBCKeeper),
	), nil
}
