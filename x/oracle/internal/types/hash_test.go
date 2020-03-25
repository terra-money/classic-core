package types

import (
	"encoding/hex"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/mock"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"

	"testing"
)

func TestVoteHash(t *testing.T) {
	_, addrs := mock.GeneratePrivKeyAddressPairs(1)

	voteHash := GetVoteHash("salt", sdk.OneDec(), "denom", sdk.ValAddress(addrs[0]))
	hexStr := hex.EncodeToString(voteHash)
	voteHashRes, err := VoteHashFromHexString(hexStr)
	require.NoError(t, err)
	require.Equal(t, voteHash, voteHashRes)
	require.True(t, voteHash.Equal(voteHash))
	require.True(t, VoteHash([]byte{}).Empty())

	got, _ := yaml.Marshal(&voteHash)
	require.Equal(t, voteHash.String()+"\n", string(got))

	res := VoteHash{}
	testMarshal(t, &voteHash, &res, voteHash.MarshalJSON, (&res).UnmarshalJSON)
	testMarshal(t, &voteHash, &res, voteHash.Marshal, (&res).Unmarshal)
}

func TestAggregateVoteHash(t *testing.T) {
	_, addr := mock.GeneratePrivKeyAddressPairs(1)
	aggregateVoteHash := GetAggregateVoteHash("salt", "100ukrw,200uusd", sdk.ValAddress(addr[0]))
	hexStr := hex.EncodeToString(aggregateVoteHash)
	aggregateVoteHashRes, err := AggregateVoteHashFromHexString(hexStr)
	require.NoError(t, err)
	require.Equal(t, aggregateVoteHash, aggregateVoteHashRes)
	require.True(t, aggregateVoteHash.Equal(aggregateVoteHash))
	require.True(t, AggregateVoteHash([]byte{}).Empty())

	got, _ := yaml.Marshal(&aggregateVoteHash)
	require.Equal(t, aggregateVoteHash.String()+"\n", string(got))

	res := AggregateVoteHash{}
	testMarshal(t, &aggregateVoteHash, &res, aggregateVoteHash.MarshalJSON, (&res).UnmarshalJSON)
	testMarshal(t, &aggregateVoteHash, &res, aggregateVoteHash.Marshal, (&res).Unmarshal)
}

func testMarshal(t *testing.T, original interface{}, res interface{}, marshal func() ([]byte, error), unmarshal func([]byte) error) {
	bz, err := marshal()
	require.Nil(t, err)
	err = unmarshal(bz)
	require.Nil(t, err)
	require.Equal(t, original, res)
}
