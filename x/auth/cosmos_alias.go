// nolint
package auth

import (
	"github.com/cosmos/cosmos-sdk/x/auth"
)

const (
	ModuleName                    = auth.ModuleName
	StoreKey                      = auth.StoreKey
	FeeCollectorName              = auth.FeeCollectorName
	QuerierRoute                  = auth.QuerierRoute
	DefaultParamspace             = auth.DefaultParamspace
	DefaultMaxMemoCharacters      = auth.DefaultMaxMemoCharacters
	DefaultTxSigLimit             = auth.DefaultTxSigLimit
	DefaultTxSizeCostPerByte      = auth.DefaultTxSizeCostPerByte
	DefaultSigVerifyCostED25519   = auth.DefaultSigVerifyCostED25519
	DefaultSigVerifyCostSecp256k1 = auth.DefaultSigVerifyCostSecp256k1
	QueryAccount                  = auth.QueryAccount
)

var (
	// functions aliases
	NewAnteHandler                    = auth.NewAnteHandler
	GetSignerAcc                      = auth.GetSignerAcc
	DefaultSigVerificationGasConsumer = auth.DefaultSigVerificationGasConsumer
	DeductFees                        = auth.DeductFees
	SetGasMeter                       = auth.SetGasMeter
	NewAccountKeeper                  = auth.NewAccountKeeper
	NewQuerier                        = auth.NewQuerier
	NewBaseAccount                    = auth.NewBaseAccount
	ProtoBaseAccount                  = auth.ProtoBaseAccount
	NewBaseAccountWithAddress         = auth.NewBaseAccountWithAddress
	NewAccountRetriever               = auth.NewAccountRetriever
	NewGenesisState                   = auth.NewGenesisState
	DefaultGenesisState               = auth.DefaultGenesisState
	ValidateGenesis                   = auth.ValidateGenesis
	SanitizeGenesisAccounts           = auth.SanitizeGenesisAccounts
	AddressStoreKey                   = auth.AddressStoreKey
	NewParams                         = auth.NewParams
	ParamKeyTable                     = auth.ParamKeyTable
	DefaultParams                     = auth.DefaultParams
	NewQueryAccountParams             = auth.NewQueryAccountParams
	NewStdTx                          = auth.NewStdTx
	CountSubKeys                      = auth.CountSubKeys
	NewStdFee                         = auth.NewStdFee
	StdSignBytes                      = auth.StdSignBytes
	DefaultTxDecoder                  = auth.DefaultTxDecoder
	DefaultTxEncoder                  = auth.DefaultTxEncoder
	NewTxBuilder                      = auth.NewTxBuilder
	NewTxBuilderFromCLI               = auth.NewTxBuilderFromCLI
	MakeSignature                     = auth.MakeSignature
	ValidateGenAccounts               = auth.ValidateGenAccounts
	GetGenesisStateFromAppState       = auth.GetGenesisStateFromAppState
	NewCosmosAppModule                = auth.NewAppModule

	// variable aliases
	CosmosModuleCdc           = auth.ModuleCdc
	AddressStoreKeyPrefix     = auth.AddressStoreKeyPrefix
	GlobalAccountNumberKey    = auth.GlobalAccountNumberKey
	KeyMaxMemoCharacters      = auth.KeyMaxMemoCharacters
	KeyTxSigLimit             = auth.KeyTxSigLimit
	KeyTxSizeCostPerByte      = auth.KeyTxSizeCostPerByte
	KeySigVerifyCostED25519   = auth.KeySigVerifyCostED25519
	KeySigVerifyCostSecp256k1 = auth.KeySigVerifyCostSecp256k1
)

type (
	SignatureVerificationGasConsumer = auth.SignatureVerificationGasConsumer
	AccountKeeper                    = auth.AccountKeeper
	BaseAccount                      = auth.BaseAccount
	NodeQuerier                      = auth.NodeQuerier
	AccountRetriever                 = auth.AccountRetriever
	GenesisState                     = auth.GenesisState
	Params                           = auth.Params
	QueryAccountParams               = auth.QueryAccountParams
	StdSignMsg                       = auth.StdSignMsg
	StdTx                            = auth.StdTx
	StdFee                           = auth.StdFee
	StdSignDoc                       = auth.StdSignDoc
	StdSignature                     = auth.StdSignature
	TxBuilder                        = auth.TxBuilder
	GenesisAccountIterator           = auth.GenesisAccountIterator

	CosmosAppModuleBasic = auth.AppModuleBasic
	CosmosAppModule      = auth.AppModule
)
