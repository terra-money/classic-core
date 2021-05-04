package v04

import (
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// ModuleName nolint
	ModuleName = "msgauth"
)

type (
	// Authorization represents the interface of various Authorization instances
	Authorization interface{}

	// GenericAuthorization grants the permission to execute any transaction of the provided
	// msg type without restrictions
	GenericAuthorization struct {
		// GrantMsgType is the type of Msg this capability grant allows
		GrantMsgType string `json:"grant_msg_type"`
	}

	//SendAuthorization grants the permission to execute send transaction
	SendAuthorization struct {
		// SpendLimit specifies the maximum amount of tokens that can be spent
		// by this authorization and will be updated as tokens are spent. If it is
		// empty, there is no spend limit and any amount of coins can be spent.
		SpendLimit sdk.Coins `json:"spend_limit"`
	}

	// AuthorizationEntry hold each authorization information
	AuthorizationEntry struct {
		Granter       sdk.AccAddress `json:"granter" yaml:"granter"`
		Grantee       sdk.AccAddress `json:"grantee" yaml:"grantee"`
		Authorization Authorization  `json:"authorization" yaml:"authorization"`
		Expiration    time.Time      `json:"expiration" yaml:"expiration"`
	}

	// GenesisState is the struct representation of the export genesis
	GenesisState struct {
		AuthorizationEntries []AuthorizationEntry `json:"authorization_entries" yaml:"authorization_entries"`
	}
)

// RegisterLegacyAminoCodec nolint
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterInterface((*Authorization)(nil), nil)
	cdc.RegisterConcrete(GenericAuthorization{}, "msgauth/GenericAuthorization", nil)
	cdc.RegisterConcrete(SendAuthorization{}, "msgauth/SendAuthorization", nil)
}
