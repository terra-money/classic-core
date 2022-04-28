package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPortIDForContract(t *testing.T) {
	_, _, addr := keyPubAddr()
	portID := PortIDForContract(addr)
	require.Equal(t, portIDPrefix+addr.String(), portID)

	addr2, err := ContractFromPortID(portID)
	require.NoError(t, err)
	require.Equal(t, addr, addr2)
}

func TestContractFromPortID_InvalidPortID(t *testing.T) {
	_, _, addr := keyPubAddr()

	_, err := ContractFromPortID("invalid" + addr.String())
	require.Error(t, err)
}
