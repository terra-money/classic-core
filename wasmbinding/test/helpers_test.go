package wasmbinding_test

import (
	"os"
	"testing"

	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	apptesting "github.com/classic-terra/core/app/testing"
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

func (s *WasmTestSuite) InstantiateContract(addr sdk.AccAddress, contractDir string) sdk.AccAddress {
	wasmKeeper := s.App.WasmKeeper

	codeId := s.storeReflectCode(addr, contractDir)

	cInfo := wasmKeeper.GetCodeInfo(s.Ctx, codeId)
	s.Require().NotNil(cInfo)

	contractAddr := s.instantiateContract(addr, codeId)

	// check if contract is instantiated
	info := wasmKeeper.GetContractInfo(s.Ctx, contractAddr)
	s.Require().NotNil(info)

	return contractAddr
}

func (s *WasmTestSuite) storeReflectCode(addr sdk.AccAddress, contractDir string) uint64 {
	wasmCode, err := os.ReadFile(contractDir)
	s.Require().NoError(err)

	codeId, _, err := wasmkeeper.NewDefaultPermissionKeeper(s.App.WasmKeeper).Create(s.Ctx, addr, wasmCode, &wasmtypes.AllowEverybody)
	s.Require().NoError(err)

	return codeId
}

func (s *WasmTestSuite) instantiateContract(funder sdk.AccAddress, codeId uint64) sdk.AccAddress {
	initMsgBz := []byte("{}")
	contractKeeper := wasmkeeper.NewDefaultPermissionKeeper(s.App.WasmKeeper)
	addr, _, err := contractKeeper.Instantiate(s.Ctx, codeId, funder, funder, initMsgBz, nil)
	s.Require().NoError(err)

	return addr
}
