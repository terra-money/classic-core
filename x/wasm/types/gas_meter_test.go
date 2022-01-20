package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOverflow_ToWasmVM(t *testing.T) {
	require.Panics(t,
		func() { ToWasmVMGas(uint64(0x7f_ff_ff_ff_ff_ff_ff_ff)) },
	)

	require.NotPanics(t, func() {
		ToWasmVMGas(uint64(0x7f_ff_ff_ff_ff_ff_ff))
	})
}
