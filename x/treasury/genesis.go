package treasury

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GenesisState - all distribution state that must be provided at genesis
type GenesisState struct {
	OracleShare  Share `json:"oracle_share"`  // oracle share
	DebtShare    Share `json:"debt_share"`    // debt share
	SubsidyShare Share `json:"subsidy_share"` // subsidy share
}

// NewGenesisState - new treasury genesis state instance
func NewGenesisState(oracleWeight sdk.Dec, debtWeight sdk.Dec, subsidyWeight sdk.Dec) GenesisState {
	return GenesisState{
		OracleShare:  NewBaseShare(OracleShareID, oracleWeight),
		DebtShare:    NewBaseShare(DebtShareID, debtWeight),
		SubsidyShare: NewBaseShare(SubsidyShareID, subsidyWeight),
	}
}

// DefaultGenesisState - get raw genesis raw message for testing
func DefaultGenesisState() GenesisState {
	return GenesisState{
		OracleShare:  NewBaseShare(OracleShareID, sdk.NewDecWithPrec(10, 2)),
		DebtShare:    NewBaseShare(DebtShareID, sdk.ZeroDec()),
		SubsidyShare: NewBaseShare(SubsidyShareID, sdk.NewDecWithPrec(90, 2)),
	}
}

// InitGenesis creates the new treasury genesis
func InitGenesis(ctx sdk.Context, keeper Keeper, data GenesisState) {
	keeper.ResetShares(ctx, []Share{
		data.OracleShare, data.DebtShare, data.SubsidyShare,
	})
}

// ExportGenesis returns a GenesisState for a given context and keeper. The
// GenesisState will contain the pool, and validator/delegator distribution info's
func ExportGenesis(ctx sdk.Context, keeper Keeper) GenesisState {
	oracleShare := keeper.GetShare(ctx, OracleShareID)
	debtShare := keeper.GetShare(ctx, DebtShareID)
	subsidyShare := keeper.GetShare(ctx, SubsidyShareID)
	return NewGenesisState(oracleShare.GetWeight(), debtShare.GetWeight(), subsidyShare.GetWeight())
}

// ValidateGenesis validates the provided oracle genesis state to ensure the
// expected invariants holds. (i.e. params in correct bounds, no duplicate validators)
func ValidateGenesis(data GenesisState) error {
	oracleWeight := data.OracleShare.GetWeight()
	debtWeight := data.DebtShare.GetWeight()
	subsidyWeight := data.SubsidyShare.GetWeight()

	sane := oracleWeight.Add(debtWeight).Add(subsidyWeight).Equal(sdk.OneDec())
	if !sane {
		return fmt.Errorf("Share weights should sum to one")
	}

	return nil
}
