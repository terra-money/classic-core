package types

import (
	"encoding/hex"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/mock"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"
	"testing"
)

func TestNameHash(t *testing.T) {
	nameHash := GetNameHash("name")
	hexStr := hex.EncodeToString(nameHash)
	nameHashRes, err := NameHashFromHexString(hexStr)
	require.NoError(t, err)
	require.Equal(t, nameHash, nameHashRes)
	require.True(t, nameHash.Equal(nameHash))
	require.True(t, NameHash([]byte{}).Empty())

	got, _ := yaml.Marshal(&nameHash)
	require.Equal(t, nameHash.String()+"\n", string(got))

	res := NameHash{}
	testMarshal(t, &nameHash, &res, nameHash.MarshalJSON, (&res).UnmarshalJSON)
	testMarshal(t, &nameHash, &res, nameHash.Marshal, (&res).Unmarshal)
}

func TestBidHash(t *testing.T) {
	_, addr := mock.GeneratePrivKeyAddressPairs(1)
	bidHash := GetBidHash("salt", "name", sdk.NewInt64Coin("foo", 123), addr[0])
	hexStr := hex.EncodeToString(bidHash)
	bidHashRes, err := BidHashFromHexString(hexStr)
	require.NoError(t, err)
	require.Equal(t, bidHash, bidHashRes)
	require.True(t, bidHash.Equal(bidHash))
	require.True(t, BidHash([]byte{}).Empty())

	got, _ := yaml.Marshal(&bidHash)
	require.Equal(t, bidHash.String()+"\n", string(got))

	res := BidHash{}
	testMarshal(t, &bidHash, &res, bidHash.MarshalJSON, (&res).UnmarshalJSON)
	testMarshal(t, &bidHash, &res, bidHash.Marshal, (&res).Unmarshal)
}

func testMarshal(t *testing.T, original interface{}, res interface{}, marshal func() ([]byte, error), unmarshal func([]byte) error) {
	bz, err := marshal()
	require.Nil(t, err)
	err = unmarshal(bz)
	require.Nil(t, err)
	require.Equal(t, original, res)
}
