// nolint
// autogenerated code using github.com/rigelrozanski/multitool
// aliases generated for the following subdirectories:
// ALIASGEN: github.com/terra-project/core/x/oracle/internal/types/
// ALIASGEN: github.com/terra-project/core/x/oracle/internal/keeper/
package oracle

import (
	"github.com/terra-project/core/x/oracle/internal/keeper"
	"github.com/terra-project/core/x/oracle/internal/types"
)

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
	QueryPrice             = types.QueryPrice
	QueryActives           = types.QueryActives
	QueryPrevotes          = types.QueryPrevotes
	QueryVotes             = types.QueryVotes
	QueryFeederDelegation  = types.QueryFeederDelegation
)

var (
	// functions aliases
	NewClaim                       = types.NewClaim
	RegisterCodec                  = types.RegisterCodec
	ErrInvalidHashLength           = types.ErrInvalidHashLength
	ErrUnknownDenomination         = types.ErrUnknownDenomination
	ErrInvalidPrice                = types.ErrInvalidPrice
	ErrVoterNotValidator           = types.ErrVoterNotValidator
	ErrVerificationFailed          = types.ErrVerificationFailed
	ErrNoPrevote                   = types.ErrNoPrevote
	ErrNoVote                      = types.ErrNoVote
	ErrNoVotingPermission          = types.ErrNoVotingPermission
	ErrNotRevealPeriod             = types.ErrNotRevealPeriod
	ErrInvalidSaltLength           = types.ErrInvalidSaltLength
	ErrInvalidMsgFormat            = types.ErrInvalidMsgFormat
	NewGenesisState                = types.NewGenesisState
	DefaultGenesisState            = types.DefaultGenesisState
	ValidateGenesis                = types.ValidateGenesis
	GetPrevoteKey                  = types.GetPrevoteKey
	GetVoteKey                     = types.GetVoteKey
	GetPriceKey                    = types.GetPriceKey
	GetFeederDelegationKey         = types.GetFeederDelegationKey
	NewMsgPricePrevote             = types.NewMsgPricePrevote
	NewMsgPriceVote                = types.NewMsgPriceVote
	NewMsgDelegateFeederPermission = types.NewMsgDelegateFeederPermission
	DefaultParams                  = types.DefaultParams
	NewQueryPriceParams            = types.NewQueryPriceParams
	NewQueryPrevotesParams         = types.NewQueryPrevotesParams
	NewQueryVotesParams            = types.NewQueryVotesParams
	NewQueryFeederDelegationParams = types.NewQueryFeederDelegationParams
	NewPricePrevote                = types.NewPricePrevote
	VoteHash                       = types.VoteHash
	NewPriceVote                   = types.NewPriceVote
	NewKeeper                      = keeper.NewKeeper
	ParamKeyTable                  = keeper.ParamKeyTable
	NewQuerier                     = keeper.NewQuerier

	// variable aliases
	ModuleCdc                             = types.ModuleCdc
	PrevoteKey                            = types.PrevoteKey
	VoteKey                               = types.VoteKey
	PriceKey                              = types.PriceKey
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
	DefaultSlashFraction                  = types.DefaultSlashFraction
	DefaultWhitelist                      = types.DefaultWhitelist
)

type (
	PriceBallot                 = types.PriceBallot
	Claim                       = types.Claim
	ClaimPool                   = types.ClaimPool
	DenomList                   = types.DenomList
	StakingKeeper               = types.StakingKeeper
	DistributionKeeper          = types.DistributionKeeper
	SupplyKeeper                = types.SupplyKeeper
	GenesisState                = types.GenesisState
	MsgPricePrevote             = types.MsgPricePrevote
	MsgPriceVote                = types.MsgPriceVote
	MsgDelegateFeederPermission = types.MsgDelegateFeederPermission
	Params                      = types.Params
	QueryPriceParams            = types.QueryPriceParams
	QueryPrevotesParams         = types.QueryPrevotesParams
	QueryVotesParams            = types.QueryVotesParams
	QueryFeederDelegationParams = types.QueryFeederDelegationParams
	PricePrevote                = types.PricePrevote
	PricePrevotes               = types.PricePrevotes
	PriceVote                   = types.PriceVote
	PriceVotes                  = types.PriceVotes
	Keeper                      = keeper.Keeper
)
