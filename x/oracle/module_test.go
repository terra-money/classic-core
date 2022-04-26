package oracle_test

import (
	"encoding/json"
	"testing"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/server/api"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/terra-money/core/x/oracle"
	test_util "github.com/terra-money/core/x/oracle/keeper"
	"github.com/terra-money/core/x/oracle/types"
)

func TestMarketModule(t *testing.T) {
	input := test_util.CreateTestInput(t)
	ctx, keeper, accKeeper, bankKeeper := input.Ctx, input.OracleKeeper, input.AccountKeeper, input.BankKeeper
	appModuleBasic := oracle.AppModuleBasic{}
	appCodec := test_util.MakeTestCodec(t)
	interfaceRegistry := codectypes.NewInterfaceRegistry()
	amino := codec.NewLegacyAmino()

	marshaler := codec.NewProtoCodec(interfaceRegistry)
	clientCtx := client.Context{}
	server := api.New(clientCtx, log.NewNopLogger())

	wantDefaultGenesis := map[string]json.RawMessage{
		"mockAppModuleBasic": json.RawMessage(``),
	}

	require.Equal(t, types.ModuleName, appModuleBasic.Name())
	err := appModuleBasic.ValidateGenesis(marshaler, nil, wantDefaultGenesis["mockAppModuleBasic"])
	require.Error(t, err)
	require.NotPanics(t, func() {
		appModuleBasic.RegisterLegacyAminoCodec(amino)
		appModuleBasic.RegisterInterfaces(interfaceRegistry)
		appModuleBasic.DefaultGenesis(marshaler)
		appModuleBasic.RegisterRESTRoutes(clientCtx, server.Router)
		appModuleBasic.RegisterGRPCGatewayRoutes(clientCtx, server.GRPCGatewayRouter)
		appModuleBasic.GetQueryCmd()
		appModuleBasic.GetTxCmd()
	})

	appModule := oracle.NewAppModule(appCodec, keeper, accKeeper, bankKeeper)

	require.Equal(t, types.ModuleName, appModule.Name())
	require.NotPanics(t, func() {
		appModule.Route()
		appModule.QuerierRoute()
		appModule.LegacyQuerierHandler(amino)
		appModule.ConsensusVersion()
		appModule.BeginBlock(ctx, abci.RequestBeginBlock{})
		appModule.EndBlock(ctx, abci.RequestEndBlock{})
	})

}
