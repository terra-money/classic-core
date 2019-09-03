// nolint
package bank

import (
	"github.com/cosmos/cosmos-sdk/x/bank"
)

const (
	DefaultCodespace         = bank.DefaultCodespace
	CodeSendDisabled         = bank.CodeSendDisabled
	CodeInvalidInputsOutputs = bank.CodeInvalidInputsOutputs
	ModuleName               = bank.ModuleName
	RouterKey                = bank.RouterKey
	QuerierRoute             = bank.QuerierRoute
	DefaultParamspace        = bank.DefaultParamspace
)

var (
	// functions aliases
	ErrNoInputs                     = bank.ErrNoInputs
	ErrNoOutputs                    = bank.ErrNoOutputs
	ErrInputOutputMismatch          = bank.ErrInputOutputMismatch
	ErrSendDisabled                 = bank.ErrSendDisabled
	NewBaseKeeper                   = bank.NewBaseKeeper
	NewInput                        = bank.NewInput
	NewOutput                       = bank.NewOutput
	ParamKeyTable                   = bank.ParamKeyTable
	NewCosmosAppModule              = bank.NewAppModule
	SimulateMsgSend                 = bank.SimulateMsgSend
	SimulateSingleInputMsgMultiSend = bank.SimulateSingleInputMsgMultiSend

	// variable aliases
	ParamStoreKeySendEnabled = bank.ParamStoreKeySendEnabled
	CosmosModuleCdc          = bank.ModuleCdc
)

type (
	BaseKeeper           = bank.BaseKeeper // ibc module depends on this
	Keeper               = bank.Keeper
	MsgSend              = bank.MsgSend
	MsgMultiSend         = bank.MsgMultiSend
	Input                = bank.Input
	Output               = bank.Output
	CosmosAppModule      = bank.AppModule
	CosmosAppModuleBasic = bank.AppModuleBasic
)
