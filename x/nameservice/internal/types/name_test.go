package types

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestName_GetNameHash(t *testing.T) {
	validName := Name("chai.terra")
	invalidName := Name("terra")
	invalidName2 := Name("forth.third.second.terra")

	require.NotPanics(t, func() { validName.NameHash() })
	require.Panics(t, func() { invalidName.NameHash() })
	require.Panics(t, func() { invalidName2.NameHash() })
	require.Equal(t, validName.Levels(), 2)
	require.Equal(t, invalidName.Levels(), 1)
	require.NotPanics(t, func() {
		root, parent, child := validName.Split()
		require.Equal(t, "terra", root)
		require.Equal(t, "chai", parent)
		require.Equal(t, "", child)
	})
	require.Panics(t, func() { invalidName.Split() })
	require.Panics(t, func() { invalidName2.Split() })
}
