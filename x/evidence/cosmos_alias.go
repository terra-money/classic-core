package evidence

import (
	"github.com/cosmos/cosmos-sdk/x/evidence"
)

// nolint

const (
	ModuleName               = evidence.ModuleName
	StoreKey                 = evidence.StoreKey
	RouterKey                = evidence.RouterKey
	QuerierRoute             = evidence.QuerierRoute
	DefaultParamspace        = evidence.DefaultParamspace
	QueryEvidence            = evidence.QueryEvidence
	QueryAllEvidence         = evidence.QueryAllEvidence
	QueryParameters          = evidence.QueryParameters
	TypeMsgSubmitEvidence    = evidence.TypeMsgSubmitEvidence
	EventTypeSubmitEvidence  = evidence.EventTypeSubmitEvidence
	AttributeValueCategory   = evidence.AttributeValueCategory
	AttributeKeyEvidenceHash = evidence.AttributeKeyEvidenceHash
	DefaultMaxEvidenceAge    = evidence.DefaultMaxEvidenceAge
)

var (
	NewKeeper              = evidence.NewKeeper
	NewQuerier             = evidence.NewQuerier
	NewMsgSubitEvidence    = evidence.NewMsgSubmitEvidence
	NewRouter              = evidence.NewRouter
	NewQueryEvidenceParams = evidence.NewQueryEvidenceParams
	NewQueryAllEvideParams = evidence.NewQueryAllEvidenceParams
	NewGenesisState        = evidence.NewGenesisState

	DefaultGenesisState          = evidence.DefaultGenesisState
	ConvertDuplicateVoteEvidence = evidence.ConvertDuplicateVoteEvidence
	KeyMaxEvidenceAge            = evidence.KeyMaxEvidenceAge
	DoubleSignJailEndTime        = evidence.DoubleSignJailEndTime
	ParamKeyTable                = evidence.ParamKeyTable

	NewCosmosAppModule = evidence.NewAppModule
	NewCosmosAppModsic = evidence.NewAppModuleBasic
	CosmosModuleCdc    = evidence.ModuleCdc
)

type (
	Keeper = evidence.Keeper

	GenesisState      = evidence.GenesisState
	MsgSubmitEvidence = evidence.MsgSubmitEvidence
	Handler           = evidence.Handler
	Router            = evidence.Router
	Equivocation      = evidence.Equivocation

	CosmosAppModule      = evidence.AppModule
	CosmosAppModuleBasic = evidence.AppModuleBasic
)
