package wasmbinding

import (
	//	"fmt"

	//	sdk "github.com/cosmos/cosmos-sdk/types"

	//	"github.com/classic-terra/core/wasmbinding/bindings"
	marketkeeper "github.com/classic-terra/core/x/market/keeper"
	oraclekeeper "github.com/classic-terra/core/x/oracle/keeper"
	treasurykeeper "github.com/classic-terra/core/x/treasury/keeper"
)

type QueryPlugin struct {
	marketKeeper   *marketkeeper.Keeper
	oracleKeeper   *oraclekeeper.Keeper
	treasuryKeeper *treasurykeeper.Keeper
}

// NewQueryPlugin returns a reference to a new QueryPlugin.
func NewQueryPlugin(tmk *marketkeeper.Keeper, tok *oraclekeeper.Keeper, ttk *treasurykeeper.Keeper) *QueryPlugin {
	return &QueryPlugin{
		marketKeeper:   tmk,
		oracleKeeper:   tok,
		treasuryKeeper: ttk,
	}
}
