package wasmbinding_test

import (
	"os"
	"testing"

	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	apptesting "github.com/classic-terra/core/v2/app/testing"
	"github.com/stretchr/testify/suite"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type WasmTestSuite struct {
	apptesting.KeeperTestHelper
}

func TestWasmTestSuite(t *testing.T) {
	suite.Run(t, new(WasmTestSuite))
}

func (s *WasmTestSuite) SetupTest() {
	s.Setup(s.T())
}

func (s *WasmTestSuite) InstantiateContract(addr sdk.AccAddress, contractPath string) sdk.AccAddress {
	wasmKeeper := s.App.WasmKeeper

	codeID := s.storeReflectCode(addr, contractPath)

	cInfo := wasmKeeper.GetCodeInfo(s.Ctx, codeID)
	s.Require().NotNil(cInfo)

	contractAddr := s.instantiateContract(addr, codeID)

	// check if contract is instantiated
	info := wasmKeeper.GetContractInfo(s.Ctx, contractAddr)
	s.Require().NotNil(info)

	return contractAddr
}

func (s *WasmTestSuite) storeReflectCode(addr sdk.AccAddress, contractPath string) uint64 {
	wasmCode, err := os.ReadFile(contractPath)
	s.Require().NoError(err)

	codeID, _, err := wasmkeeper.NewDefaultPermissionKeeper(s.App.WasmKeeper).Create(s.Ctx, addr, wasmCode, &wasmtypes.AllowEverybody)
	s.Require().NoError(err)

	return codeID
}

func (s *WasmTestSuite) instantiateContract(funder sdk.AccAddress, codeID uint64) sdk.AccAddress {
	initMsgBz := []byte("{}")
	contractKeeper := wasmkeeper.NewDefaultPermissionKeeper(s.App.WasmKeeper)
	addr, _, err := contractKeeper.Instantiate(s.Ctx, codeID, funder, funder, initMsgBz, "label", nil)
	s.Require().NoError(err)

	return addr
}
