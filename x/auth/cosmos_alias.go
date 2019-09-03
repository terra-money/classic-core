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
	NewBaseAccount                 = auth.NewBaseAccount
	ProtoBaseAccount               = auth.ProtoBaseAccount
	NewBaseAccountWithAddress      = auth.NewBaseAccountWithAddress
	NewBaseVestingAccount          = auth.NewBaseVestingAccount
	NewContinuousVestingAccountRaw = auth.NewContinuousVestingAccountRaw
	NewContinuousVestingAccount    = auth.NewContinuousVestingAccount
	NewDelayedVestingAccountRaw    = auth.NewDelayedVestingAccountRaw
	NewDelayedVestingAccount       = auth.NewDelayedVestingAccount
	NewGenesisState                = auth.NewGenesisState
	ValidateGenesis                = auth.ValidateGenesis
	AddressStoreKey                = auth.AddressStoreKey
	NewParams                      = auth.NewParams
	ParamKeyTable                  = auth.ParamKeyTable
	DefaultParams                  = auth.DefaultParams
	NewQueryAccountParams          = auth.NewQueryAccountParams
	NewStdTx                       = auth.NewStdTx
	CountSubKeys                   = auth.CountSubKeys
	NewStdFee                      = auth.NewStdFee
	StdSignBytes                   = auth.StdSignBytes
	DefaultTxDecoder               = auth.DefaultTxDecoder
	DefaultTxEncoder               = auth.DefaultTxEncoder
	NewTxBuilder                   = auth.NewTxBuilder
	NewTxBuilderFromCLI            = auth.NewTxBuilderFromCLI
	MakeSignature                  = auth.MakeSignature
	NewAccountRetriever            = auth.NewAccountRetriever
	NewAccountKeeper               = auth.NewAccountKeeper

	// variable aliases
	AddressStoreKeyPrefix     = auth.AddressStoreKeyPrefix
	GlobalAccountNumberKey    = auth.GlobalAccountNumberKey
	KeyMaxMemoCharacters      = auth.KeyMaxMemoCharacters
	KeyTxSigLimit             = auth.KeyTxSigLimit
	KeyTxSizeCostPerByte      = auth.KeyTxSizeCostPerByte
	KeySigVerifyCostED25519   = auth.KeySigVerifyCostED25519
	KeySigVerifyCostSecp256k1 = auth.KeySigVerifyCostSecp256k1
	CosmosModuleCdc           = auth.ModuleCdc
	NewCosmosAppModule        = auth.NewAppModule
)

type (
	Account                  = auth.Account
	VestingAccount           = auth.VestingAccount
	BaseAccount              = auth.BaseAccount
	BaseVestingAccount       = auth.BaseVestingAccount
	ContinuousVestingAccount = auth.ContinuousVestingAccount
	DelayedVestingAccount    = auth.DelayedVestingAccount
	GenesisState             = auth.GenesisState
	Params                   = auth.Params
	QueryAccountParams       = auth.QueryAccountParams
	StdSignMsg               = auth.StdSignMsg
	StdTx                    = auth.StdTx
	StdFee                   = auth.StdFee
	StdSignDoc               = auth.StdSignDoc
	StdSignature             = auth.StdSignature
	TxBuilder                = auth.TxBuilder
	AccountKeeper            = auth.AccountKeeper
	CosmosAppModule          = auth.AppModule
	CosmosAppModuleBasic     = auth.AppModuleBasic
)
