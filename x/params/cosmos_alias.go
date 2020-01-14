// nolint
package params

import (
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/params/types"
)

const (
	StoreKey             = params.StoreKey
	TStoreKey            = params.TStoreKey
	TestParamStore       = params.TestParamStore
	DefaultCodespace     = params.DefaultCodespace
	CodeUnknownSubspace  = params.CodeUnknownSubspace
	CodeSettingParameter = params.CodeSettingParameter
	CodeEmptyData        = params.CodeEmptyData
	ModuleName           = params.ModuleName
	RouterKey            = params.RouterKey
	ProposalTypeChange   = params.ProposalTypeChange
)

var (
	// functions aliases
	NewSubspace                   = params.NewSubspace
	NewKeyTable                   = params.NewKeyTable
	DefaultTestComponents         = params.DefaultTestComponents
	ErrUnknownSubspace            = params.ErrUnknownSubspace
	ErrSettingParameter           = params.ErrSettingParameter
	ErrEmptyChanges               = params.ErrEmptyChanges
	ErrEmptySubspace              = params.ErrEmptySubspace
	ErrEmptyKey                   = params.ErrEmptyKey
	ErrEmptyValue                 = params.ErrEmptyValue
	NewParameterChangeProposal    = params.NewParameterChangeProposal
	NewParamChange                = params.NewParamChange
	NewParamChangeWithSubkey      = params.NewParamChangeWithSubkey
	ValidateChanges               = params.ValidateChanges
	NewKeeper                     = params.NewKeeper
	NewParamChangeProposalHandler = params.NewParamChangeProposalHandler

	// variables aliases
	CosmosModuleCdc = types.ModuleCdc
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
)
