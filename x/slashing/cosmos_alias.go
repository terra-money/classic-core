// nolint
package slashing

import (
	"github.com/cosmos/cosmos-sdk/x/slashing"
)

const (
	ModuleName                  = slashing.ModuleName
	StoreKey                    = slashing.StoreKey
	RouterKey                   = slashing.RouterKey
	QuerierRoute                = slashing.QuerierRoute
	QueryParameters             = slashing.QueryParameters
	QuerySigningInfo            = slashing.QuerySigningInfo
	QuerySigningInfos           = slashing.QuerySigningInfos
	DefaultParamspace           = slashing.DefaultParamspace
	DefaultSignedBlocksWindow   = slashing.DefaultSignedBlocksWindow
	DefaultDowntimeJailDuration = slashing.DefaultDowntimeJailDuration
)

var (
	// functions aliases
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
	DefaultMinSignedPerWindow       = slashing.DefaultMinSignedPerWindow
	DefaultSlashFractionDoubleSign  = slashing.DefaultSlashFractionDoubleSign
	DefaultSlashFractionDowntime    = slashing.DefaultSlashFractionDowntime
	KeySignedBlocksWindow           = slashing.KeySignedBlocksWindow
	KeyMinSignedPerWindow           = slashing.KeyMinSignedPerWindow
	KeyDowntimeJailDuration         = slashing.KeyDowntimeJailDuration
	KeySlashFractionDoubleSign      = slashing.KeySlashFractionDoubleSign
	KeySlashFractionDowntime        = slashing.KeySlashFractionDowntime
)

type (
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
