package types

import (
	fmt "fmt"
	"time"

	proto "github.com/gogo/protobuf/proto"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ codectypes.UnpackInterfacesMessage = AuthorizationGrant{}
var _ codectypes.UnpackInterfacesMessage = AuthorizationGrants{}

// AuthorizationI represents the interface of various Authorization instances
type AuthorizationI interface {
	MsgType() string
	Accept(msg sdk.Msg, blocktime tmproto.Header) (allow bool, updated AuthorizationI, delete bool)
}

// NewAuthorizationGrant returns new AuthorizationGrant instance
func NewAuthorizationGrant(authorization AuthorizationI, expiration time.Time) (AuthorizationGrant, error) {
	msg, ok := authorization.(proto.Message)
	if !ok {
		return AuthorizationGrant{}, fmt.Errorf("%T does not implement proto.Message", authorization)
	}

	any, err := codectypes.NewAnyWithValue(msg)
	if err != nil {
		return AuthorizationGrant{}, err
	}

	return AuthorizationGrant{Authorization: any, Expiration: expiration}, nil
}

// GetAuthorization returns the grant Authorization
func (g AuthorizationGrant) GetAuthorization() AuthorizationI {
	authorization, ok := g.Authorization.GetCachedValue().(AuthorizationI)
	if !ok {
		return nil
	}
	return authorization
}

// UnpackInterfaces implements UnpackInterfacesMessage.UnpackInterfaces
func (g AuthorizationGrant) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	var authorization AuthorizationI
	return unpacker.UnpackAny(g.Authorization, &authorization)
}

// AuthorizationGrants array of grant
type AuthorizationGrants []AuthorizationGrant

// Equal returns true if two slices (order-dependant) of grants are equal.
func (grants AuthorizationGrants) Equal(other AuthorizationGrants) bool {
	if len(grants) != len(other) {
		return false
	}

	for i, grant := range grants {
		if !grant.Equal(other[i]) {
			return false
		}
	}

	return true
}

// UnpackInterfaces implements UnpackInterfacesMessage.UnpackInterfaces
func (grants AuthorizationGrants) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	for _, grant := range grants {
		err := grant.UnpackInterfaces(unpacker)
		if err != nil {
			return err
		}
	}

	return nil
}

var grantableMsgTypes = map[string]struct{}{}

// RegisterGrantableMsgType registers a grantable msg type. It will panic if the type is
// already registered.
func RegisterGrantableMsgType(ty string) {
	if _, ok := grantableMsgTypes[ty]; ok {
		panic(fmt.Sprintf("already registered proposal type: %s", ty))
	}

	grantableMsgTypes[ty] = struct{}{}
}

// IsGrantableMsgType returns a boolean determining if the msg type is
// grantable.
//
// NOTE: Modules with their own proposal types must register them.
func IsGrantableMsgType(ty string) bool {
	_, ok := grantableMsgTypes[ty]
	return ok
}
