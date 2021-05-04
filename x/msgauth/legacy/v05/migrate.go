package v05

import (
	"fmt"

	proto "github.com/gogo/protobuf/proto"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"

	v04msgauth "github.com/terra-project/core/x/msgauth/legacy/v04"
	v05msgauth "github.com/terra-project/core/x/msgauth/types"
)

func migrateAuthorization(oldAuthorization v04msgauth.Authorization) *codectypes.Any {
	var protoAuthorization proto.Message

	switch oldAuthorization := oldAuthorization.(type) {
	case v04msgauth.GenericAuthorization:
		{
			protoAuthorization = &v05msgauth.GenericAuthorization{
				GrantMsgType: oldAuthorization.GrantMsgType,
			}
		}
	case v04msgauth.SendAuthorization:
		{
			protoAuthorization = &v05msgauth.SendAuthorization{
				SpendLimit: oldAuthorization.SpendLimit,
			}
		}

	default:
		panic(fmt.Errorf("%T is not a valid proposal content type", oldAuthorization))
	}

	// Convert the Authorization into Any.
	authorizationAny, err := codectypes.NewAnyWithValue(protoAuthorization)
	if err != nil {
		panic(err)
	}

	return authorizationAny
}

// Migrate accepts exported v0.4 x/msgauth genesis state and migrates it to
// v0.5 x/msgauth genesis state. The migration includes:
//
// - Convert vote option & proposal status from byte to enum.
// - Migrate proposal content to Any.
// - Convert addresses from bytes to bech32 strings.
// - Re-encode in v0.40 GenesisState.
func Migrate(msgauthGenState v04msgauth.GenesisState) *v05msgauth.GenesisState {
	entries := make([]v05msgauth.AuthorizationEntry, len(msgauthGenState.AuthorizationEntries))
	for i, e := range msgauthGenState.AuthorizationEntries {
		entries[i] = v05msgauth.AuthorizationEntry{
			Granter:       e.Granter.String(),
			Grantee:       e.Grantee.String(),
			Authorization: migrateAuthorization(e.Authorization),
			Expiration:    e.Expiration,
		}
	}

	return &v05msgauth.GenesisState{
		AuthorizationEntries: entries,
	}
}
