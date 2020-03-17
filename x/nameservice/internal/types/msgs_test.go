package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/mock"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func TestMsgOpenAuction_ValidateBasic(t *testing.T) {
	_, acc := mock.GeneratePrivKeyAddressPairs(1)

	require.NoError(t, NewMsgOpenAuction("valid.terra", acc[0]).ValidateBasic())
	require.Error(t, NewMsgOpenAuction("valid.dokwon.terra", acc[0]).ValidateBasic())

	// empty name
	require.Error(t, NewMsgOpenAuction("", acc[0]).ValidateBasic())

	// double dot name
	require.Error(t, NewMsgOpenAuction("invalid..terra", acc[0]).ValidateBasic())

	// too long label
	require.Error(t, NewMsgOpenAuction(Name(strings.Repeat("long", 16)+".terra"), acc[0]).ValidateBasic())

	// too long name
	require.Error(t, NewMsgOpenAuction(Name(strings.Repeat(strings.Repeat("long", 8)+".", 8)+"terra"), acc[0]).ValidateBasic())

	// empty address
	require.Error(t, NewMsgOpenAuction("valid.terra", sdk.AccAddress{}).ValidateBasic())
}

func TestMsgBidAuction_ValidateBasic(t *testing.T) {
	_, acc := mock.GeneratePrivKeyAddressPairs(1)

	validSalt := "salt"
	validName := Name("valid.terra")
	invalidName := Name("invalid.valid.terra")
	validCoin := sdk.NewInt64Coin("foo", 123)
	validHash := GetBidHash(validSalt, validName, validCoin, acc[0])

	require.NoError(t, NewMsgBidAuction(validName, validHash, validCoin, acc[0]).ValidateBasic())

	// invalid hash length
	require.Error(t, NewMsgBidAuction(validName, BidHash{123}, validCoin, acc[0]).ValidateBasic())

	// invalid deposit
	require.Error(t, NewMsgBidAuction(validName, validHash, sdk.Coin{}, acc[0]).ValidateBasic())

	// invalid bidder address
	require.Error(t, NewMsgBidAuction(validName, validHash, validCoin, sdk.AccAddress{}).ValidateBasic())

	// invalid name
	require.Error(t, NewMsgBidAuction("", validHash, validCoin, acc[0]).ValidateBasic())

	// invalid name
	require.Error(t, NewMsgBidAuction(invalidName, validHash, validCoin, acc[0]).ValidateBasic())

}

func TestMsgRevealBid_ValidateBasic(t *testing.T) {
	_, acc := mock.GeneratePrivKeyAddressPairs(1)

	validSalt := "salt"
	validName := Name("valid.terra")
	invalidName := Name("invalid.valid.terra")
	validCoin := sdk.NewInt64Coin("foo", 123)

	require.NoError(t, NewMsgRevealBid(validName, validSalt, validCoin, acc[0]).ValidateBasic())

	// invalid name
	require.Error(t, NewMsgRevealBid("", validSalt, validCoin, acc[0]).ValidateBasic())

	// invalid salt
	require.Error(t, NewMsgRevealBid(validName, "", validCoin, acc[0]).ValidateBasic())
	require.Error(t, NewMsgRevealBid(validName, "123456789", validCoin, acc[0]).ValidateBasic())

	// invalid amount
	require.Error(t, NewMsgRevealBid(validName, validSalt, sdk.Coin{}, acc[0]).ValidateBasic())

	// invalid address
	require.Error(t, NewMsgRevealBid(validName, validSalt, validCoin, sdk.AccAddress{}).ValidateBasic())

	// invalid name
	require.Error(t, NewMsgRevealBid(invalidName, validSalt, validCoin, acc[0]).ValidateBasic())
}

func TestMsgRenewRegistry_ValidateBasic(t *testing.T) {
	_, acc := mock.GeneratePrivKeyAddressPairs(1)

	validName := Name("valid.terra")
	validCoins := sdk.Coins{sdk.NewInt64Coin("foo", 12345)}
	require.NoError(t, NewMsgRenewRegistry(validName, validCoins, acc[0]).ValidateBasic())

	// invalid name
	require.Error(t, NewMsgRenewRegistry("", validCoins, acc[0]).ValidateBasic())

	// invalid coin
	require.Error(t, NewMsgRenewRegistry(validName, sdk.Coins{sdk.Coin{}}, acc[0]).ValidateBasic())

	// invalid address
	require.Error(t, NewMsgRenewRegistry(validName, validCoins, sdk.AccAddress{}).ValidateBasic())
}

func TestMsgUpdateOwner_ValidateBasic(t *testing.T) {
	_, acc := mock.GeneratePrivKeyAddressPairs(2)

	validName := Name("valid.terra")
	require.NoError(t, NewMsgUpdateOwner(validName, acc[0], acc[1]).ValidateBasic())

	// invalid name
	require.Error(t, NewMsgUpdateOwner("", acc[0], acc[1]).ValidateBasic())

	// invalid new owner address
	require.Error(t, NewMsgUpdateOwner(validName, sdk.AccAddress{}, acc[1]).ValidateBasic())

	// invalid owner address
	require.Error(t, NewMsgUpdateOwner(validName, acc[0], sdk.AccAddress{}).ValidateBasic())
}

func TestMsgRegisterSubName_ValidateBasic(t *testing.T) {
	_, acc := mock.GeneratePrivKeyAddressPairs(2)

	validName := Name("acc.wallet.terra")
	invalidName := Name("hi.acc.wallet.terra")

	// valid
	require.NoError(t, NewMsgRegisterSubName(validName, acc[1], acc[0]).ValidateBasic())

	// invalid name
	require.Error(t, NewMsgRegisterSubName("", acc[1], acc[0]).ValidateBasic())
	require.Error(t, NewMsgRegisterSubName(invalidName, acc[1], acc[0]).ValidateBasic())

	// invalid address
	require.Error(t, NewMsgRegisterSubName(validName, sdk.AccAddress{}, acc[0]).ValidateBasic())

	// invalid owner
	require.Error(t, NewMsgRegisterSubName(validName, acc[1], sdk.AccAddress{}).ValidateBasic())
}

func TestMsgUnregisterSubName_ValidateBasic(t *testing.T) {
	_, acc := mock.GeneratePrivKeyAddressPairs(1)

	validName := Name("acc.valid.terra")
	invalidName := Name("hi.acc.wallet.terra")

	// valid
	require.NoError(t, NewMsgUnregisterSubName(validName, acc[0]).ValidateBasic())

	// invalid name
	require.Error(t, NewMsgUnregisterSubName("", acc[0]).ValidateBasic())
	require.Error(t, NewMsgUnregisterSubName(invalidName, acc[0]).ValidateBasic())

	// invalid owner
	require.Error(t, NewMsgUnregisterSubName(validName, sdk.AccAddress{}).ValidateBasic())
}
