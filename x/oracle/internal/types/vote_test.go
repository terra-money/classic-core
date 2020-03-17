package types

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestParseExchangeRateTuples(t *testing.T) {
	valid := "123.0uluna,123.123ukrw"
	_, err := ParseExchangeRateTuples(valid)
	require.NoError(t, err)

	duplicatedDenom := "100.0uluna,123.123ukrw,121233.123ukrw"
	_, err = ParseExchangeRateTuples(duplicatedDenom)
	require.Error(t, err)

	invalidCoin := "123.123"
	_, err = ParseExchangeRateTuples(invalidCoin)
	require.Error(t, err)

	invalidCoinWithValid := "123.0uluna,123.1"
	_, err = ParseExchangeRateTuples(invalidCoinWithValid)
	require.Error(t, err)
}
