// nolint
package gov

import (
	"github.com/cosmos/cosmos-sdk/x/gov"
)

const (
	MaxDescriptionLength         = gov.MaxDescriptionLength
	MaxTitleLength               = gov.MaxTitleLength
	DefaultCodespace             = gov.DefaultCodespace
	CodeUnknownProposal          = gov.CodeUnknownProposal
	CodeInactiveProposal         = gov.CodeInactiveProposal
	CodeAlreadyActiveProposal    = gov.CodeAlreadyActiveProposal
	CodeAlreadyFinishedProposal  = gov.CodeAlreadyFinishedProposal
	CodeAddressNotStaked         = gov.CodeAddressNotStaked
	CodeInvalidContent           = gov.CodeInvalidContent
	CodeInvalidProposalType      = gov.CodeInvalidProposalType
	CodeInvalidVote              = gov.CodeInvalidVote
	CodeInvalidGenesis           = gov.CodeInvalidGenesis
	CodeInvalidProposalStatus    = gov.CodeInvalidProposalStatus
	CodeProposalHandlerNotExists = gov.CodeProposalHandlerNotExists
	ModuleName                   = gov.ModuleName
	StoreKey                     = gov.StoreKey
	RouterKey                    = gov.RouterKey
	QuerierRoute                 = gov.QuerierRoute
	DefaultParamspace            = gov.DefaultParamspace
	TypeMsgDeposit               = gov.TypeMsgDeposit
	TypeMsgVote                  = gov.TypeMsgVote
	TypeMsgSubmitProposal        = gov.TypeMsgSubmitProposal
	StatusNil                    = gov.StatusNil
	StatusDepositPeriod          = gov.StatusDepositPeriod
	StatusVotingPeriod           = gov.StatusVotingPeriod
	StatusPassed                 = gov.StatusPassed
	StatusRejected               = gov.StatusRejected
	StatusFailed                 = gov.StatusFailed
	ProposalTypeText             = gov.ProposalTypeText
	ProposalTypeSoftwareUpgrade  = gov.ProposalTypeSoftwareUpgrade
	QueryParams                  = gov.QueryParams
	QueryProposals               = gov.QueryProposals
	QueryProposal                = gov.QueryProposal
	QueryDeposits                = gov.QueryDeposits
	QueryDeposit                 = gov.QueryDeposit
	QueryVotes                   = gov.QueryVotes
	QueryVote                    = gov.QueryVote
	QueryTally                   = gov.QueryTally
	ParamDeposit                 = gov.ParamDeposit
	ParamVoting                  = gov.ParamVoting
	ParamTallying                = gov.ParamTallying
	OptionEmpty                  = gov.OptionEmpty
	OptionYes                    = gov.OptionYes
	OptionAbstain                = gov.OptionAbstain
	OptionNo                     = gov.OptionNo
	OptionNoWithVeto             = gov.OptionNoWithVeto
)

var (
	// functions aliases
	ValidateAbstract              = gov.ValidateAbstract
	NewDeposit                    = gov.NewDeposit
	ErrUnknownProposal            = gov.ErrUnknownProposal
	ErrInactiveProposal           = gov.ErrInactiveProposal
	ErrAlreadyActiveProposal      = gov.ErrAlreadyActiveProposal
	ErrAlreadyFinishedProposal    = gov.ErrAlreadyFinishedProposal
	ErrAddressNotStaked           = gov.ErrAddressNotStaked
	ErrInvalidProposalContent     = gov.ErrInvalidProposalContent
	ErrInvalidProposalType        = gov.ErrInvalidProposalType
	ErrInvalidVote                = gov.ErrInvalidVote
	ErrInvalidGenesis             = gov.ErrInvalidGenesis
	ErrNoProposalHandlerExists    = gov.ErrNoProposalHandlerExists
	ProposalKey                   = gov.ProposalKey
	ActiveProposalByTimeKey       = gov.ActiveProposalByTimeKey
	ActiveProposalQueueKey        = gov.ActiveProposalQueueKey
	InactiveProposalByTimeKey     = gov.InactiveProposalByTimeKey
	InactiveProposalQueueKey      = gov.InactiveProposalQueueKey
	DefaultGenesisState           = gov.DefaultGenesisState
	DepositsKey                   = gov.DepositsKey
	DepositKey                    = gov.DepositKey
	VotesKey                      = gov.VotesKey
	VoteKey                       = gov.VoteKey
	SplitProposalKey              = gov.SplitProposalKey
	SplitActiveProposalQueueKey   = gov.SplitActiveProposalQueueKey
	SplitInactiveProposalQueueKey = gov.SplitInactiveProposalQueueKey
	SplitKeyDeposit               = gov.SplitKeyDeposit
	SplitKeyVote                  = gov.SplitKeyVote
	NewMsgSubmitProposal          = gov.NewMsgSubmitProposal
	NewMsgDeposit                 = gov.NewMsgDeposit
	NewMsgVote                    = gov.NewMsgVote
	ParamKeyTable                 = gov.ParamKeyTable
	NewDepositParams              = gov.NewDepositParams
	NewTallyParams                = gov.NewTallyParams
	NewVotingParams               = gov.NewVotingParams
	NewParams                     = gov.NewParams
	NewProposal                   = gov.NewProposal
	ProposalStatusFromString      = gov.ProposalStatusFromString
	ValidProposalStatus           = gov.ValidProposalStatus
	NewTallyResult                = gov.NewTallyResult
	NewTallyResultFromMap         = gov.NewTallyResultFromMap
	EmptyTallyResult              = gov.EmptyTallyResult
	NewTextProposal               = gov.NewTextProposal
	NewSoftwareUpgradeProposal    = gov.NewSoftwareUpgradeProposal
	RegisterProposalType          = gov.RegisterProposalType
	ContentFromProposalType       = gov.ContentFromProposalType
	IsValidProposalType           = gov.IsValidProposalType
	ProposalHandler               = gov.ProposalHandler
	NewQueryProposalParams        = gov.NewQueryProposalParams
	NewQueryDepositParams         = gov.NewQueryDepositParams
	NewQueryVoteParams            = gov.NewQueryVoteParams
	NewQueryProposalsParams       = gov.NewQueryProposalsParams
	NewVote                       = gov.NewVote
	VoteOptionFromString          = gov.VoteOptionFromString
	ValidVoteOption               = gov.ValidVoteOption
	NewKeeper                     = gov.NewKeeper
	NewRouter                     = gov.NewRouter
	NewCosmosAppModule            = gov.NewAppModule
	NewCosmosAppModuleBasic       = gov.NewAppModuleBasic

	// variable aliases
	CosmosModuleCdc             = gov.ModuleCdc
	ProposalsKeyPrefix          = gov.ProposalsKeyPrefix
	ActiveProposalQueuePrefix   = gov.ActiveProposalQueuePrefix
	InactiveProposalQueuePrefix = gov.InactiveProposalQueuePrefix
	ProposalIDKey               = gov.ProposalIDKey
	DepositsKeyPrefix           = gov.DepositsKeyPrefix
	VotesKeyPrefix              = gov.VotesKeyPrefix
	ParamStoreKeyDepositParams  = gov.ParamStoreKeyDepositParams
	ParamStoreKeyVotingParams   = gov.ParamStoreKeyVotingParams
	ParamStoreKeyTallyParams    = gov.ParamStoreKeyTallyParams
)

type (
	Content                 = gov.Content
	Keeper                  = gov.Keeper
	Handler                 = gov.Handler
	Deposit                 = gov.Deposit
	Deposits                = gov.Deposits
	MsgSubmitProposal       = gov.MsgSubmitProposal
	MsgDeposit              = gov.MsgDeposit
	MsgVote                 = gov.MsgVote
	DepositParams           = gov.DepositParams
	TallyParams             = gov.TallyParams
	VotingParams            = gov.VotingParams
	Params                  = gov.Params
	Proposal                = gov.Proposal
	Proposals               = gov.Proposals
	ProposalQueue           = gov.ProposalQueue
	ProposalStatus          = gov.ProposalStatus
	TallyResult             = gov.TallyResult
	TextProposal            = gov.TextProposal
	SoftwareUpgradeProposal = gov.SoftwareUpgradeProposal
	QueryProposalParams     = gov.QueryProposalParams
	QueryDepositParams      = gov.QueryDepositParams
	QueryVoteParams         = gov.QueryVoteParams
	QueryProposalsParams    = gov.QueryProposalsParams
	Vote                    = gov.Vote
	Votes                   = gov.Votes
	VoteOption              = gov.VoteOption
	CosmosAppModule         = gov.AppModule
	CosmosAppModuleBasic    = gov.AppModuleBasic
)
