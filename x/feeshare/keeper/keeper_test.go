package keeper_test

// import (
// 	"fmt"
// 	"testing"
// 	"time"

// 	"github.com/cosmos/cosmos-sdk/baseapp"
// 	sdk "github.com/cosmos/cosmos-sdk/types"
// 	"github.com/stretchr/testify/suite"
// 	tmrand "github.com/tendermint/tendermint/libs/rand"
// 	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

// 	"github.com/classic-terra/core/v2/app"
// 	apphelpers "github.com/classic-terra/core/v2/app/helpers"
// 	"github.com/classic-terra/core/v2/x/feeshare/keeper"
// 	"github.com/classic-terra/core/v2/x/feeshare/types"
// 	wasmtypes "github.com/classic-terra/core/v2/x/wasm/types"
// 	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
// )

// // BankKeeper defines the expected interface needed to retrieve account balances.
// type BankKeeper interface {
// 	MintCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) error
// 	SendCoins(ctx sdk.Context, fromAddr, toAddr sdk.AccAddress, amt sdk.Coins) error
// 	SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
// 	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
// 	SendCoinsFromModuleToModule(ctx sdk.Context, senderModule, recipientModule string, amt sdk.Coins) error
// }\

// type IntegrationTestSuite struct {
// 	suite.Suite

// 	ctx               sdk.Context
// 	app               *app.TerraApp
// 	bankKeeper        BankKeeper
// 	accountKeeper     types.AccountKeeper
// 	queryClient       types.QueryClient
// 	feeShareMsgServer types.MsgServer
// 	wasmMsgServer     wasmtypes.MsgServer
// }

// func (s *IntegrationTestSuite) SetupTest() {
// 	isCheckTx := false
// 	s.app = apphelpers.Setup(s.T(), isCheckTx, 1)

// 	s.ctx = s.app.BaseApp.NewContext(isCheckTx, tmproto.Header{
// 		ChainID: fmt.Sprintf("test-chain-%s", tmrand.Str(4)),
// 		Height:  9,
// 		Time:    time.Now().UTC(),
// 	})

// 	queryHelper := baseapp.NewQueryServerTestHelper(s.ctx, s.app.InterfaceRegistry())
// 	types.RegisterQueryServer(queryHelper, keeper.NewQuerier(s.app.FeeShareKeeper))

// 	s.queryClient = types.NewQueryClient(queryHelper)
// 	s.bankKeeper = s.app.BankKeeper
// 	s.accountKeeper = s.app.AccountKeeper
// 	s.feeShareMsgServer = s.app.FeeShareKeeper
// }

// func (s *IntegrationTestSuite) FundAccount(ctx sdk.Context, addr sdk.AccAddress, amounts sdk.Coins) error {
// 	if err := s.bankKeeper.MintCoins(ctx, minttypes.ModuleName, amounts); err != nil {
// 		return err
// 	}

// 	return s.bankKeeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, addr, amounts)
// }

// func TestKeeperTestSuite(t *testing.T) {
// 	suite.Run(t, new(IntegrationTestSuite))
// }
