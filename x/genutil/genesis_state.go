package genutil

import (
	"errors"
	"fmt"

	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/terra-project/core/x/auth"
)

// validate GenTx transactions
func ValidateGenesis(genesisState GenesisState) error {
	for i, genTx := range genesisState.GenTxs {
		var tx auth.StdTx
		if err := ModuleCdc.UnmarshalJSON(genTx, &tx); err != nil {
			return err
		}

		msgs := tx.GetMsgs()
		if len(msgs) != 1 {
			return errors.New(
				"must provide genesis StdTx with exactly 1 CreateValidator message")
		}

		// TODO abstract back to staking
		if _, ok := msgs[0].(staking.MsgCreateValidator); !ok {
			return fmt.Errorf(
				"Genesis transaction %v does not contain a MsgCreateValidator", i)
		}
	}
	return nil
}
