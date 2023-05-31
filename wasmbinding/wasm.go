package wasmbinding

import (
	"github.com/CosmWasm/wasmd/x/wasm"

	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	//	"github.com/cosmos/cosmos-sdk/baseapp"
	//	"github.com/cosmos/cosmos-sdk/codec"

	marketkeeper "github.com/classic-terra/core/x/market/keeper"
	oraclekeeper "github.com/classic-terra/core/x/oracle/keeper"
	treasurykeeper "github.com/classic-terra/core/x/treasury/keeper"
	// tokenfactorykeeper "github.com/osmosis-labs/osmosis/v15/x/tokenfactory/keeper"
)

func RegisterCustomPlugins(
	marketKeeper *marketkeeper.Keeper,
	oracleKeeper *oraclekeeper.Keeper,
	treasuryKeeper *treasurykeeper.Keeper,
) []wasmkeeper.Option {
	wasmQueryPlugin := NewQueryPlugin(
		marketKeeper,
		oracleKeeper,
		treasuryKeeper,
	)

	queryPluginOpt := wasmkeeper.WithQueryPlugins(&wasmkeeper.QueryPlugins{
		Custom: CustomQuerier(wasmQueryPlugin),
	})
	messengerDecoratorOpt := wasmkeeper.WithMessageHandlerDecorator(
		CustomMessageDecorator(marketKeeper),
	)

	return []wasm.Option{
		queryPluginOpt,
		messengerDecoratorOpt,
	}
}
