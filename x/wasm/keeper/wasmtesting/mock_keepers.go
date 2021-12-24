package wasmtesting

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/terra-money/core/x/wasm/types"
)

var _ types.ICS20TransferPortSource = &MockIBCTransferKeeper{}

// MockIBCTransferKeeper nolint
type MockIBCTransferKeeper struct {
	GetPortFn func(ctx sdk.Context) string
}

// GetPort no lint
func (m MockIBCTransferKeeper) GetPort(ctx sdk.Context) string {
	if m.GetPortFn == nil {
		panic("not expected to be called")
	}
	return m.GetPortFn(ctx)
}
