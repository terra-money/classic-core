package common

import (
	"encoding/json"
	"fmt"
)

// Convenience struct for CLI output
type PrettyParams struct {
	BaseProposerReward  json.RawMessage `json:"base_proposer_reward"`
	BonusProposerReward json.RawMessage `json:"bonus_proposer_reward"`
	WithdrawAddrEnabled json.RawMessage `json:"withdraw_addr_enabled"`
}

// Construct a new PrettyParams
func NewPrettyParams(baseProposerReward json.RawMessage, bonusProposerReward json.RawMessage, withdrawAddrEnabled json.RawMessage) PrettyParams {
	return PrettyParams{
		BaseProposerReward:  baseProposerReward,
		BonusProposerReward: bonusProposerReward,
		WithdrawAddrEnabled: withdrawAddrEnabled,
	}
}

func (pp PrettyParams) String() string {
	return fmt.Sprintf(`Distribution Params:
  Base Proposer Reward:   %s
  Bonus Proposer Reward:  %s
  Withdraw Addr Enabled:  %s`,
		pp.BaseProposerReward, pp.BonusProposerReward, pp.WithdrawAddrEnabled)

}
