// nolint
package params

import (
	"github.com/cosmos/cosmos-sdk/x/params"
)

const (
	StoreKey         = params.StoreKey
	TStoreKey        = params.TStoreKey
	ModuleNam        = params.ModuleName
	RouterKey        = params.RouterKey
	ProposalTeChange = params.ProposalTypeChange
)

var (
	// functions aliases
	NewParamSetPair               = params.NewParamSetPair
	NewSubspace                   = params.NewSubspace
	NewKeyTable                   = params.NewKeyTable
	NewKeeper                     = params.NewKeeper
	NewParamChangeProposalHandler = params.NewParamChangeProposalHandler
	ErrUnknownSubspace            = params.ErrUnknownSubspace
	ErrSettingParameter           = params.ErrSettingParameter
	ErrEmptyChanges               = params.ErrEmptyChanges
	ErrEmptySubspace              = params.ErrEmptySubspace
	ErrEmptyKey                   = params.ErrEmptyKey
	ErrEmptyValue                 = params.ErrEmptyValue
	NewParameterChangeProposal    = params.NewParameterChangeProposal
	NewParamChange                = params.NewParamChange
	ValidateChanges               = params.ValidateChanges

	// variables aliases
	CosmosModuleCdc    = params.ModuleCdc
	NewCosmosAppModule = params.NewAppModule
)

type (
	ParamSetPair            = params.ParamSetPair
	ParamSetPairs           = params.ParamSetPairs
	ParamSet                = params.ParamSet
	Subspace                = params.Subspace
	ReadOnlySubspace        = params.ReadOnlySubspace
	KeyTable                = params.KeyTable
	ParameterChangeProposal = params.ParameterChangeProposal
	ParamChange             = params.ParamChange
	Keeper                  = params.Keeper
	CosmosAppModuleBasic    = params.AppModuleBasic
	CosmosAppModule         = params.AppModule
)
