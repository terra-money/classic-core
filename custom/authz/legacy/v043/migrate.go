package v043

import (
	"fmt"

	proto "github.com/gogo/protobuf/proto"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"

	v05market "github.com/terra-money/core/x/market/types"
	v04msgauth "github.com/terra-money/core/x/msgauth/legacy/v04"

	sdk "github.com/cosmos/cosmos-sdk/types"
	v043authz "github.com/cosmos/cosmos-sdk/x/authz"
	v043bank "github.com/cosmos/cosmos-sdk/x/bank/types"
	v043gov "github.com/cosmos/cosmos-sdk/x/gov/types"
)

func migrateAuthorization(oldAuthorization v04msgauth.Authorization) *codectypes.Any {
	var protoAuthorization proto.Message

	switch oldAuthorization := oldAuthorization.(type) {
	case v04msgauth.GenericAuthorization:
		{
			var msgTypeURL string
			if oldAuthorization.GrantMsgType == "swap" {
				msgTypeURL = sdk.MsgTypeURL(&v05market.MsgSwap{})
			} else if oldAuthorization.GrantMsgType == "vote" {
				msgTypeURL = sdk.MsgTypeURL(&v043gov.MsgVote{})
			} else {
				panic(fmt.Errorf("%T is not a valid generic authorization msg type", oldAuthorization.GrantMsgType))
			}

			protoAuthorization = &v043authz.GenericAuthorization{
				Msg: msgTypeURL,
			}
		}
	case v04msgauth.SendAuthorization:
		{
			protoAuthorization = &v043bank.SendAuthorization{
				SpendLimit: oldAuthorization.SpendLimit,
			}
		}

	default:
		panic(fmt.Errorf("%T is not a valid authorization type", oldAuthorization))
	}

	// Convert the Authorization into Any.
	authorizationAny, err := codectypes.NewAnyWithValue(protoAuthorization)
	if err != nil {
		panic(err)
	}

	return authorizationAny
}

// Migrate accepts exported v0.4 x/msgauth genesis state and migrates it to
// cosmos-sdk@v0.43 x/authz genesis state. The migration includes:
//
// - Convert vote option & proposal status from byte to enum.
// - Migrate proposal content to Any.
// - Convert addresses from bytes to bech32 strings.
// - Re-encode in v0.43 GenesisState.
func Migrate(msgauthGenState v04msgauth.GenesisState) *v043authz.GenesisState {
	entries := make([]v043authz.GrantAuthorization, len(msgauthGenState.AuthorizationEntries))
	for i, e := range msgauthGenState.AuthorizationEntries {
		entries[i] = v043authz.GrantAuthorization{
			Granter:       e.Granter.String(),
			Grantee:       e.Grantee.String(),
			Authorization: migrateAuthorization(e.Authorization),
			Expiration:    e.Expiration,
		}
	}

	return &v043authz.GenesisState{
		Authorization: entries,
	}
}
