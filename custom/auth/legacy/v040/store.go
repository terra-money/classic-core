// Package v040 creates in-place store migrations for fixing
// multisig pubkey migration problem
// ref: https://github.com/terra-money/core/issues/562
package v040

import (
	"github.com/cosmos/cosmos-sdk/crypto/types/multisig"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
)

// Set MultiSig account PubKey as nil
func migrateMultiSigAccount(account types.AccountI) (types.AccountI, error) {
	_, ok := account.GetPubKey().(multisig.PubKey)
	if !ok {
		return nil, nil
	}

	_ = account.SetPubKey(nil)
	return account.(types.AccountI), nil
}

// MigrateAccount migrates multisig account's PubKey as nil to restore mistakenly set PubKey
// References: https://github.com/terra-money/core/issues/562
//
func MigrateAccount(account types.AccountI) (types.AccountI, error) {
	return migrateMultiSigAccount(account)
}
