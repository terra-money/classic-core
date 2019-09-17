// nolint
package slashing

import (
	"github.com/cosmos/cosmos-sdk/x/slashing"
)

const (
	DefaultCodespace            = slashing.DefaultCodespace
	CodeInvalidValidator        = slashing.CodeInvalidValidator
	CodeValidatorJailed         = slashing.CodeValidatorJailed
	CodeValidatorNotJailed      = slashing.CodeValidatorNotJailed
	CodeMissingSelfDelegation   = slashing.CodeMissingSelfDelegation
	CodeSelfDelegationTooLow    = slashing.CodeSelfDelegationTooLow
	CodeMissingSigningInfo      = slashing.CodeMissingSigningInfo
	ModuleName                  = slashing.ModuleName
	StoreKey                    = slashing.StoreKey
	RouterKey                   = slashing.RouterKey
	QuerierRoute                = slashing.QuerierRoute
	QueryParameters             = slashing.QueryParameters
	QuerySigningInfo            = slashing.QuerySigningInfo
	QuerySigningInfos           = slashing.QuerySigningInfos
	DefaultParamspace           = slashing.DefaultParamspace
	DefaultMaxEvidenceAge       = slashing.DefaultMaxEvidenceAge
	DefaultSignedBlocksWindow   = slashing.DefaultSignedBlocksWindow
	DefaultDowntimeJailDuration = slashing.DefaultDowntimeJailDuration
)

var (
	// functions aliases
	ErrNoValidatorForAddress                 = slashing.ErrNoValidatorForAddress
	ErrBadValidatorAddr                      = slashing.ErrBadValidatorAddr
	ErrValidatorJailed                       = slashing.ErrValidatorJailed
	ErrValidatorNotJailed                    = slashing.ErrValidatorNotJailed
	ErrMissingSelfDelegation                 = slashing.ErrMissingSelfDelegation
	ErrSelfDelegationTooLowToUnjail          = slashing.ErrSelfDelegationTooLowToUnjail
	ErrNoSigningInfoFound                    = slashing.ErrNoSigningInfoFound
	NewGenesisState                          = slashing.NewGenesisState
	DefaultGenesisState                      = slashing.DefaultGenesisState
	ValidateGenesis                          = slashing.ValidateGenesis
	GetValidatorSigningInfoKey               = slashing.GetValidatorSigningInfoKey
	GetValidatorSigningInfoAddress           = slashing.GetValidatorSigningInfoAddress
	GetValidatorMissedBlockBitArrayPrefixKey = slashing.GetValidatorMissedBlockBitArrayPrefixKey
	GetValidatorMissedBlockBitArrayKey       = slashing.GetValidatorMissedBlockBitArrayKey
	GetAddrPubkeyRelationKey                 = slashing.GetAddrPubkeyRelationKey
	NewMsgUnjail                             = slashing.NewMsgUnjail
	ParamKeyTable                            = slashing.ParamKeyTable
	NewParams                                = slashing.NewParams
	DefaultParams                            = slashing.DefaultParams
	NewQuerySigningInfoParams                = slashing.NewQuerySigningInfoParams
	NewQuerySigningInfosParams               = slashing.NewQuerySigningInfosParams
	NewValidatorSigningInfo                  = slashing.NewValidatorSigningInfo
	NewCosmosAppModule                       = slashing.NewAppModule
	NewKeeper                                = slashing.NewKeeper

	// variable aliases
	CosmosModuleCdc                 = slashing.ModuleCdc
	ValidatorSigningInfoKey         = slashing.ValidatorSigningInfoKey
	ValidatorMissedBlockBitArrayKey = slashing.ValidatorMissedBlockBitArrayKey
	AddrPubkeyRelationKey           = slashing.AddrPubkeyRelationKey
	DoubleSignJailEndTime           = slashing.DoubleSignJailEndTime
	DefaultMinSignedPerWindow       = slashing.DefaultMinSignedPerWindow
	DefaultSlashFractionDoubleSign  = slashing.DefaultSlashFractionDoubleSign
	DefaultSlashFractionDowntime    = slashing.DefaultSlashFractionDowntime
	KeyMaxEvidenceAge               = slashing.KeyMaxEvidenceAge
	KeySignedBlocksWindow           = slashing.KeySignedBlocksWindow
	KeyMinSignedPerWindow           = slashing.KeyMinSignedPerWindow
	KeyDowntimeJailDuration         = slashing.KeyDowntimeJailDuration
	KeySlashFractionDoubleSign      = slashing.KeySlashFractionDoubleSign
	KeySlashFractionDowntime        = slashing.KeySlashFractionDowntime
)

type (
	CodeType                = slashing.CodeType
	GenesisState            = slashing.GenesisState
	MissedBlock             = slashing.MissedBlock
	MsgUnjail               = slashing.MsgUnjail
	Params                  = slashing.Params
	QuerySigningInfoParams  = slashing.QuerySigningInfoParams
	QuerySigningInfosParams = slashing.QuerySigningInfosParams
	ValidatorSigningInfo    = slashing.ValidatorSigningInfo
	Keeper                  = slashing.Keeper
	CosmosAppModule         = slashing.AppModule
	CosmosAppModuleBasic    = slashing.AppModuleBasic
)
