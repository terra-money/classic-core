// nolint
package distribution

import (
	distr "github.com/cosmos/cosmos-sdk/x/distribution"
)

const (
	DefaultParamspace                = distr.DefaultParamspace
	DefaultCodespace                 = distr.DefaultCodespace
	CodeInvalidInput                 = distr.CodeInvalidInput
	CodeNoDistributionInfo           = distr.CodeNoDistributionInfo
	CodeNoValidatorCommission        = distr.CodeNoValidatorCommission
	CodeSetWithdrawAddrDisabled      = distr.CodeSetWithdrawAddrDisabled
	ModuleName                       = distr.ModuleName
	StoreKey                         = distr.StoreKey
	RouterKey                        = distr.RouterKey
	QuerierRoute                     = distr.QuerierRoute
	ProposalTypeCommunityPoolSpend   = distr.ProposalTypeCommunityPoolSpend
	QueryParams                      = distr.QueryParams
	QueryValidatorOutstandingRewards = distr.QueryValidatorOutstandingRewards
	QueryValidatorCommission         = distr.QueryValidatorCommission
	QueryValidatorSlashes            = distr.QueryValidatorSlashes
	QueryDelegationRewards           = distr.QueryDelegationRewards
	QueryDelegatorTotalRewards       = distr.QueryDelegatorTotalRewards
	QueryDelegatorValidators         = distr.QueryDelegatorValidators
	QueryWithdrawAddr                = distr.QueryWithdrawAddr
	QueryCommunityPool               = distr.QueryCommunityPool
	ParamCommunityTax                = distr.ParamCommunityTax
	ParamBaseProposerReward          = distr.ParamBaseProposerReward
	ParamBonusProposerReward         = distr.ParamBonusProposerReward
	ParamWithdrawAddrEnabled         = distr.ParamWithdrawAddrEnabled
)

var (
	// functions aliases
	RegisterInvariants                         = distr.RegisterInvariants
	AllInvariants                              = distr.AllInvariants
	NonNegativeOutstandingInvariant            = distr.NonNegativeOutstandingInvariant
	CanWithdrawInvariant                       = distr.CanWithdrawInvariant
	ReferenceCountInvariant                    = distr.ReferenceCountInvariant
	ModuleAccountInvariant                     = distr.ModuleAccountInvariant
	NewKeeper                                  = distr.NewKeeper
	GetValidatorOutstandingRewardsAddress      = distr.GetValidatorOutstandingRewardsAddress
	GetDelegatorWithdrawInfoAddress            = distr.GetDelegatorWithdrawInfoAddress
	GetDelegatorStartingInfoAddresses          = distr.GetDelegatorStartingInfoAddresses
	GetValidatorHistoricalRewardsAddressPeriod = distr.GetValidatorHistoricalRewardsAddressPeriod
	GetValidatorCurrentRewardsAddress          = distr.GetValidatorCurrentRewardsAddress
	GetValidatorAccumulatedCommissionAddress   = distr.GetValidatorAccumulatedCommissionAddress
	GetValidatorSlashEventAddressHeight        = distr.GetValidatorSlashEventAddressHeight
	GetValidatorOutstandingRewardsKey          = distr.GetValidatorOutstandingRewardsKey
	GetDelegatorWithdrawAddrKey                = distr.GetDelegatorWithdrawAddrKey
	GetDelegatorStartingInfoKey                = distr.GetDelegatorStartingInfoKey
	GetValidatorHistoricalRewardsPrefix        = distr.GetValidatorHistoricalRewardsPrefix
	GetValidatorHistoricalRewardsKey           = distr.GetValidatorHistoricalRewardsKey
	GetValidatorCurrentRewardsKey              = distr.GetValidatorCurrentRewardsKey
	GetValidatorAccumulatedCommissionKey       = distr.GetValidatorAccumulatedCommissionKey
	GetValidatorSlashEventPrefix               = distr.GetValidatorSlashEventPrefix
	GetValidatorSlashEventKey                  = distr.GetValidatorSlashEventKey
	GetValidatorSlashEventKeyPrefix            = distr.GetValidatorSlashEventKeyPrefix
	ParamKeyTable                              = distr.ParamKeyTable
	HandleCommunityPoolSpendProposal           = distr.HandleCommunityPoolSpendProposal
	NewQuerier                                 = distr.NewQuerier
	MakeTestCodec                              = distr.MakeTestCodec
	CreateTestInputDefault                     = distr.CreateTestInputDefault
	CreateTestInputAdvanced                    = distr.CreateTestInputAdvanced
	NewDelegatorStartingInfo                   = distr.NewDelegatorStartingInfo
	ErrNilDelegatorAddr                        = distr.ErrNilDelegatorAddr
	ErrNilWithdrawAddr                         = distr.ErrNilWithdrawAddr
	ErrNilValidatorAddr                        = distr.ErrNilValidatorAddr
	ErrNoDelegationDistInfo                    = distr.ErrNoDelegationDistInfo
	ErrNoValidatorDistInfo                     = distr.ErrNoValidatorDistInfo
	ErrNoValidatorCommission                   = distr.ErrNoValidatorCommission
	ErrSetWithdrawAddrDisabled                 = distr.ErrSetWithdrawAddrDisabled
	ErrBadDistribution                         = distr.ErrBadDistribution
	ErrInvalidProposalAmount                   = distr.ErrInvalidProposalAmount
	ErrEmptyProposalRecipient                  = distr.ErrEmptyProposalRecipient
	InitialFeePool                             = distr.InitialFeePool
	NewMsgSetWithdrawAddress                   = distr.NewMsgSetWithdrawAddress
	NewMsgWithdrawDelegatorReward              = distr.NewMsgWithdrawDelegatorReward
	NewMsgWithdrawValidatorCommission          = distr.NewMsgWithdrawValidatorCommission
	NewCommunityPoolSpendProposal              = distr.NewCommunityPoolSpendProposal
	NewQueryValidatorOutstandingRewardsParams  = distr.NewQueryValidatorOutstandingRewardsParams
	NewQueryValidatorCommissionParams          = distr.NewQueryValidatorCommissionParams
	NewQueryValidatorSlashesParams             = distr.NewQueryValidatorSlashesParams
	NewQueryDelegationRewardsParams            = distr.NewQueryDelegationRewardsParams
	NewQueryDelegatorParams                    = distr.NewQueryDelegatorParams
	NewQueryDelegatorWithdrawAddrParams        = distr.NewQueryDelegatorWithdrawAddrParams
	NewQueryDelegatorTotalRewardsResponse      = distr.NewQueryDelegatorTotalRewardsResponse
	NewDelegationDelegatorReward               = distr.NewDelegationDelegatorReward
	NewValidatorHistoricalRewards              = distr.NewValidatorHistoricalRewards
	NewValidatorCurrentRewards                 = distr.NewValidatorCurrentRewards
	InitialValidatorAccumulatedCommission      = distr.InitialValidatorAccumulatedCommission
	NewValidatorSlashEvent                     = distr.NewValidatorSlashEvent
	NewCommunityPoolSpendProposalHandler       = distr.NewCommunityPoolSpendProposalHandler
	NewCosmosAppModule                         = distr.NewAppModule

	// variable aliases
	FeePoolKey                           = distr.FeePoolKey
	ProposerKey                          = distr.ProposerKey
	ValidatorOutstandingRewardsPrefix    = distr.ValidatorOutstandingRewardsPrefix
	DelegatorWithdrawAddrPrefix          = distr.DelegatorWithdrawAddrPrefix
	DelegatorStartingInfoPrefix          = distr.DelegatorStartingInfoPrefix
	ValidatorHistoricalRewardsPrefix     = distr.ValidatorHistoricalRewardsPrefix
	ValidatorCurrentRewardsPrefix        = distr.ValidatorCurrentRewardsPrefix
	ValidatorAccumulatedCommissionPrefix = distr.ValidatorAccumulatedCommissionPrefix
	ValidatorSlashEventPrefix            = distr.ValidatorSlashEventPrefix
	ParamStoreKeyCommunityTax            = distr.ParamStoreKeyCommunityTax
	ParamStoreKeyBaseProposerReward      = distr.ParamStoreKeyBaseProposerReward
	ParamStoreKeyBonusProposerReward     = distr.ParamStoreKeyBonusProposerReward
	ParamStoreKeyWithdrawAddrEnabled     = distr.ParamStoreKeyWithdrawAddrEnabled
	TestAddrs                            = distr.TestAddrs
	EventTypeRewards                     = distr.EventTypeRewards
	EventTypeCommission                  = distr.EventTypeCommission
	AttributeValueCategory               = distr.AttributeValueCategory
	AttributeKeyValidator                = distr.AttributeKeyValidator
	CosmosModuleCdc                      = distr.ModuleCdc
)

type (
	Hooks                                  = distr.Hooks
	Keeper                                 = distr.Keeper
	DelegatorStartingInfo                  = distr.DelegatorStartingInfo
	CodeType                               = distr.CodeType
	FeePool                                = distr.FeePool
	DelegatorWithdrawInfo                  = distr.DelegatorWithdrawInfo
	ValidatorOutstandingRewardsRecord      = distr.ValidatorOutstandingRewardsRecord
	ValidatorAccumulatedCommissionRecord   = distr.ValidatorAccumulatedCommissionRecord
	ValidatorHistoricalRewardsRecord       = distr.ValidatorHistoricalRewardsRecord
	ValidatorCurrentRewardsRecord          = distr.ValidatorCurrentRewardsRecord
	DelegatorStartingInfoRecord            = distr.DelegatorStartingInfoRecord
	ValidatorSlashEventRecord              = distr.ValidatorSlashEventRecord
	GenesisState                           = distr.GenesisState
	MsgSetWithdrawAddress                  = distr.MsgSetWithdrawAddress
	MsgWithdrawDelegatorReward             = distr.MsgWithdrawDelegatorReward
	MsgWithdrawValidatorCommission         = distr.MsgWithdrawValidatorCommission
	CommunityPoolSpendProposal             = distr.CommunityPoolSpendProposal
	QueryValidatorOutstandingRewardsParams = distr.QueryValidatorOutstandingRewardsParams
	QueryValidatorCommissionParams         = distr.QueryValidatorCommissionParams
	QueryValidatorSlashesParams            = distr.QueryValidatorSlashesParams
	QueryDelegationRewardsParams           = distr.QueryDelegationRewardsParams
	QueryDelegatorParams                   = distr.QueryDelegatorParams
	QueryDelegatorWithdrawAddrParams       = distr.QueryDelegatorWithdrawAddrParams
	QueryDelegatorTotalRewardsResponse     = distr.QueryDelegatorTotalRewardsResponse
	DelegationDelegatorReward              = distr.DelegationDelegatorReward
	ValidatorHistoricalRewards             = distr.ValidatorHistoricalRewards
	ValidatorCurrentRewards                = distr.ValidatorCurrentRewards
	ValidatorAccumulatedCommission         = distr.ValidatorAccumulatedCommission
	ValidatorSlashEvent                    = distr.ValidatorSlashEvent
	ValidatorSlashEvents                   = distr.ValidatorSlashEvents
	ValidatorOutstandingRewards            = distr.ValidatorOutstandingRewards
	CosmosAppModule                        = distr.AppModule
	CosmosAppModuleBasic                   = distr.AppModuleBasic
)
