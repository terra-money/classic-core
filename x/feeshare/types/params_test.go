package types

import (
	"testing"

	appparams "github.com/classic-terra/core/v2/app/params"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/stretchr/testify/require"
)

func TestParamKeyTable(t *testing.T) {
	require.IsType(t, paramtypes.KeyTable{}, ParamKeyTable())
	require.NotEmpty(t, ParamKeyTable())
}

func TestParamSetPairs(t *testing.T) {
	params := DefaultParams()
	require.NotEmpty(t, params.ParamSetPairs())
}

func TestParamsValidate(t *testing.T) {
	devShares := sdk.NewDecWithPrec(60, 2)
	acceptedDenoms := []string{appparams.BondDenom}

	testCases := []struct {
		name     string
		params   Params
		expError bool
	}{
		{"default", DefaultParams(), false},
		{
			"valid: enabled",
			NewParams(true, devShares, acceptedDenoms),
			false,
		},
		{
			"valid: disabled",
			NewParams(false, devShares, acceptedDenoms),
			false,
		},
		{
			"valid: 100% devs",
			Params{true, sdk.NewDecFromInt(sdk.NewInt(1)), acceptedDenoms},
			false,
		},
		{
			"empty",
			Params{},
			true,
		},
		{
			"invalid: share > 1",
			Params{true, sdk.NewDecFromInt(sdk.NewInt(2)), acceptedDenoms},
			true,
		},
		{
			"invalid: share < 0",
			Params{true, sdk.NewDecFromInt(sdk.NewInt(-1)), acceptedDenoms},
			true,
		},
		{
			"valid: all denoms allowed",
			Params{true, sdk.NewDecFromInt(sdk.NewInt(-1)), []string{}},
			true,
		},
	}
	for _, tc := range testCases {
		err := tc.params.Validate()

		if tc.expError {
			require.Error(t, err, tc.name)
		} else {
			require.NoError(t, err, tc.name)
		}
	}
}

func TestParamsValidateShares(t *testing.T) {
	testCases := []struct {
		name     string
		value    interface{}
		expError bool
	}{
		{"default", DefaultDeveloperShares, false},
		{"valid", sdk.NewDecFromInt(sdk.NewInt(1)), false},
		{"invalid - wrong type - bool", false, true},
		{"invalid - wrong type - string", "", true},
		{"invalid - wrong type - int64", int64(123), true},
		{"invalid - wrong type - math.Int", sdk.NewInt(1), true},
		{"invalid - is nil", nil, true},
		{"invalid - is negative", sdk.NewDecFromInt(sdk.NewInt(-1)), true},
		{"invalid - is > 1", sdk.NewDecFromInt(sdk.NewInt(2)), true},
	}
	for _, tc := range testCases {
		err := validateShares(tc.value)

		if tc.expError {
			require.Error(t, err, tc.name)
		} else {
			require.NoError(t, err, tc.name)
		}
	}
}

func TestParamsValidateBool(t *testing.T) {
	err := validateBool(DefaultEnableFeeShare)
	require.NoError(t, err)
	err = validateBool(true)
	require.NoError(t, err)
	err = validateBool(false)
	require.NoError(t, err)
	err = validateBool("")
	require.Error(t, err)
	err = validateBool(int64(123))
	require.Error(t, err)
}
