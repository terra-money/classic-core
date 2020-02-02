// nolint
package bank

import (
	"github.com/cosmos/cosmos-sdk/x/bank"
)

const (
	QueryBalance       = bank.QueryBalance
	ModuleName         = bank.ModuleName
	QuerierRoute       = bank.QuerierRoute
	RouterKey          = bank.RouterKey
	DefaultParamspace  = bank.DefaultParamspace
	DefaultSendEnabled = bank.DefaultSendEnabled
)

var (
	RegisterInvariants          = bank.RegisterInvariants
	NonnegativeBalanceInvariant = bank.NonnegativeBalanceInvariant
	NewBaseKeeper               = bank.NewBaseKeeper
	NewBaseSendKeeper           = bank.NewBaseSendKeeper
	NewBaseViewKeeper           = bank.NewBaseViewKeeper
	NewQuerier                  = bank.NewQuerier
	ErrNoInputs                 = bank.ErrNoInputs
	ErrNoOutputs                = bank.ErrNoOutputs
	ErrInputOutputMismatch      = bank.ErrInputOutputMismatch
	ErrSendDisabled             = bank.ErrSendDisabled
	NewGenesisState             = bank.NewGenesisState
	DefaultGenesisState         = bank.DefaultGenesisState
	ValidateGenesis             = bank.ValidateGenesis
	NewMsgSend                  = bank.NewMsgSend
	NewMsgMultiSend             = bank.NewMsgMultiSend
	NewInput                    = bank.NewInput
	NewOutput                   = bank.NewOutput
	ValidateInputsOutputs       = bank.ValidateInputsOutputs
	ParamKeyTable               = bank.ParamKeyTable
	NewQueryBalanceParams       = bank.NewQueryBalanceParams
	ParamStoreKeySendEnabled    = bank.ParamStoreKeySendEnabled
	NewCosmosAppModule          = bank.NewAppModule

	CosmosModuleCdc = bank.ModuleCdc
)

type (
	Keeper             = bank.Keeper
	BaseKeeper         = bank.BaseKeeper
	SendKeeper         = bank.SendKeeper
	BaseSendKeeper     = bank.BaseSendKeeper
	ViewKeeper         = bank.ViewKeeper
	BaseViewKeeper     = bank.BaseViewKeeper
	GenesisState       = bank.GenesisState
	MsgSend            = bank.MsgSend
	MsgMultiSend       = bank.MsgMultiSend
	Input              = bank.Input
	Output             = bank.Output
	QueryBalanceParams = bank.QueryBalanceParams

	CosmosAppModule      = bank.AppModule
	CosmosAppModuleBasic = bank.AppModuleBasic
)
