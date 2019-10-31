package oracle

import (
	"github.com/terra-project/core/x/oracle/internal/keeper"
	"github.com/terra-project/core/x/oracle/internal/types"
)

//nolint
const (
	DefaultCodespace       = types.DefaultCodespace
	CodeUnknownDenom       = types.CodeUnknownDenom
	CodeInvalidPrice       = types.CodeInvalidPrice
	CodeVoterNotValidator  = types.CodeVoterNotValidator
	CodeInvalidVote        = types.CodeInvalidVote
	CodeNoVotingPermission = types.CodeNoVotingPermission
	CodeInvalidHashLength  = types.CodeInvalidHashLength
	CodeInvalidPrevote     = types.CodeInvalidPrevote
	CodeVerificationFailed = types.CodeVerificationFailed
	CodeNotRevealPeriod    = types.CodeNotRevealPeriod
	CodeInvalidSaltLength  = types.CodeInvalidSaltLength
	CodeInvalidMsgFormat   = types.CodeInvalidMsgFormat
	ModuleName             = types.ModuleName
	StoreKey               = types.StoreKey
	RouterKey              = types.RouterKey
	QuerierRoute           = types.QuerierRoute
	DefaultParamspace      = types.DefaultParamspace
	DefaultVotePeriod      = types.DefaultVotePeriod
	DefaultVotesWindow     = types.DefaultVotesWindow
	QueryParameters        = types.QueryParameters
	QueryExchangeRate      = types.QueryExchangeRate
	QueryActives           = types.QueryActives
	QueryPrevotes          = types.QueryPrevotes
	QueryVotes             = types.QueryVotes
	QueryFeederDelegation  = types.QueryFeederDelegation
)

//nolint
var (
	// functions aliases
	NewClaim                       = types.NewClaim
	RegisterCodec                  = types.RegisterCodec
	ErrInvalidHashLength           = types.ErrInvalidHashLength
	ErrUnknownDenomination         = types.ErrUnknownDenomination
	ErrInvalidPrice                = types.ErrInvalidPrice
	ErrVerificationFailed          = types.ErrVerificationFailed
	ErrNoPrevote                   = types.ErrNoPrevote
	ErrNoVote                      = types.ErrNoVote
	ErrNoVotingPermission          = types.ErrNoVotingPermission
	ErrNotRevealPeriod             = types.ErrNotRevealPeriod
	ErrInvalidSaltLength           = types.ErrInvalidSaltLength
	NewGenesisState                = types.NewGenesisState
	DefaultGenesisState            = types.DefaultGenesisState
	ValidateGenesis                = types.ValidateGenesis
	GetPrevoteKey                  = types.GetPrevoteKey
	GetVoteKey                     = types.GetVoteKey
	GetPriceKey                    = types.GetPriceKey
	GetFeederDelegationKey         = types.GetFeederDelegationKey
	NewMsgPrevote                  = types.NewMsgPrevote
	NewMsgVote                     = types.NewMsgVote
	NewMsgDelegateConsent = types.NewMsgDelegateConsent
	DefaultParams                  = types.DefaultParams
	NewQueryExchangeRateParams     = types.NewQueryExchangeRateParams
	NewQueryPrevotesParams         = types.NewQueryPrevotesParams
	NewQueryVotesParams            = types.NewQueryVotesParams
	NewQueryFeederDelegationParams = types.NewQueryFeederDelegationParams
	NewPrevote                     = types.NewPrevote
	VoteHash                       = types.VoteHash
	NewVote                        = types.NewVote
	NewKeeper                      = keeper.NewKeeper
	ParamKeyTable                  = keeper.ParamKeyTable
	NewQuerier                     = keeper.NewQuerier

	// variable aliases
	ModuleCdc                             = types.ModuleCdc
	PrevoteKey                            = types.PrevoteKey
	VoteKey                               = types.VoteKey
	PriceKey                              = types.ExchangeRateKey
	FeederDelegationKey                   = types.FeederDelegationKey
	ParamStoreKeyVotePeriod               = types.ParamStoreKeyVotePeriod
	ParamStoreKeyVoteThreshold            = types.ParamStoreKeyVoteThreshold
	ParamStoreKeyRewardBand               = types.ParamStoreKeyRewardBand
	ParamStoreKeyRewardDistributionPeriod = types.ParamStoreKeyRewardDistributionPeriod
	ParamStoreKeyWhitelist                = types.ParamStoreKeyWhitelist
	DefaultVoteThreshold                  = types.DefaultVoteThreshold
	DefaultRewardBand                     = types.DefaultRewardBand
	DefaultRewardDistributionPeriod       = types.DefaultRewardDistributionPeriod
	DefaultMinValidVotesPerWindow         = types.DefaultMinValidVotesPerWindow
	DefaultWhitelist                      = types.DefaultWhitelist
)

//nolint
type (
	ExchangeRateBallot          = types.ExchangeRateBallot
	Claim                       = types.Claim
	ClaimPool                   = types.ClaimPool
	DenomList                   = types.DenomList
	StakingKeeper               = types.StakingKeeper
	DistributionKeeper          = types.DistributionKeeper
	SupplyKeeper                = types.SupplyKeeper
	GenesisState                = types.GenesisState
	MsgPrevote                  = types.MsgPrevote
	MsgVote                     = types.MsgVote
	MsgDelegateConsent = types.MsgDelegateConsent
	Params                      = types.Params
	QueryExchangeRateParams     = types.QueryExchangeRateParams
	QueryPrevotesParams         = types.QueryPrevotesParams
	QueryVotesParams            = types.QueryVotesParams
	QueryFeederDelegationParams = types.QueryFeederDelegationParams
	Prevote                     = types.Prevote
	Prevotes                    = types.Prevotes
	Vote                        = types.Vote
	Votes                       = types.Votes
	Keeper                      = keeper.Keeper
)
