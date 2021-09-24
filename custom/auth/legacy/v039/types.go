package v039

// DONTCOVER
// nolint

import (
	"errors"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	kmultisig "github.com/cosmos/cosmos-sdk/crypto/keys/multisig"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	v038auth "github.com/cosmos/cosmos-sdk/x/auth/legacy/v038"
	v039auth "github.com/cosmos/cosmos-sdk/x/auth/legacy/v039"

	"github.com/tendermint/tendermint/crypto/sr25519"
)

type (
	vestingAccountJSON struct {
		Address          sdk.AccAddress     `json:"address" yaml:"address"`
		Coins            sdk.Coins          `json:"coins,omitempty" yaml:"coins"`
		PubKey           cryptotypes.PubKey `json:"public_key" yaml:"public_key"`
		AccountNumber    uint64             `json:"account_number" yaml:"account_number"`
		Sequence         uint64             `json:"sequence" yaml:"sequence"`
		OriginalVesting  sdk.Coins          `json:"original_vesting" yaml:"original_vesting"`
		DelegatedFree    sdk.Coins          `json:"delegated_free" yaml:"delegated_free"`
		DelegatedVesting sdk.Coins          `json:"delegated_vesting" yaml:"delegated_vesting"`
		EndTime          int64              `json:"end_time" yaml:"end_time"`

		// custom fields based on concrete vesting type which can be omitted
		VestingSchedules VestingSchedules `json:"vesting_schedules,omitempty" yaml:"vesting_schedules,omitempty"`
	}

	// LazyGradedVestingAccount nolint
	LazyGradedVestingAccount struct {
		*v039auth.BaseVestingAccount

		VestingSchedules VestingSchedules `json:"vesting_schedules"`
	}

	// LazySchedule nolint
	LazySchedule struct {
		StartTime int64   `json:"start_time"`
		EndTime   int64   `json:"end_time"`
		Ratio     sdk.Dec `json:"ratio"`
	}

	// LazySchedules nolint
	LazySchedules []LazySchedule

	// VestingSchedule nolint
	VestingSchedule struct {
		Denom         string        `json:"denom"`
		LazySchedules LazySchedules `json:"schedules"` // maps blocktime to percentage vested. Should sum to 1.
	}

	// VestingSchedules nolint
	VestingSchedules []VestingSchedule

	// LegacyAminoPubKey specifies a public key type
	// which nests multiple public keys and a threshold,
	// it uses legacy amino address rules.
	LegacyAminoPubKey struct {
		Threshold sdk.Int              `json:"threshold"`
		PubKeys   []cryptotypes.PubKey `json:"pubkeys"`
	}
)

var _ cryptotypes.PubKey = &LegacyAminoPubKey{}

// NewLazyGradedVestingAccountRaw nolint
func NewLazyGradedVestingAccountRaw(baseVestingAcc *v039auth.BaseVestingAccount, lazyVestingSchedules VestingSchedules) *LazyGradedVestingAccount {
	return &LazyGradedVestingAccount{
		BaseVestingAccount: baseVestingAcc,
		VestingSchedules:   lazyVestingSchedules,
	}
}

// Validate nolint
func (lgva LazyGradedVestingAccount) Validate() error {
	for _, vestingSchedule := range lgva.VestingSchedules {
		if err := vestingSchedule.Validate(); err != nil {
			return err
		}
	}

	return lgva.BaseVestingAccount.Validate()
}

// Validate nolint
func (s LazySchedule) Validate() error {
	startTime := s.StartTime
	endTime := s.EndTime
	ratio := s.Ratio

	if startTime < 0 {
		return errors.New("vesting start-time cannot be negative")
	}

	if endTime < startTime {
		return errors.New("vesting start-time cannot be before end-time")
	}

	if ratio.LTE(sdk.ZeroDec()) {
		return errors.New("vesting ratio cannot be smaller than or equal with zero")
	}

	return nil
}

// Validate nolint
func (vs VestingSchedule) Validate() error {
	sumRatio := sdk.ZeroDec()
	for _, lazySchedule := range vs.LazySchedules {

		if err := lazySchedule.Validate(); err != nil {
			return err
		}

		sumRatio = sumRatio.Add(lazySchedule.Ratio)
	}

	// add rounding to allow language specific calculation errors
	const fixedPointDecimals = 1000000000
	if !sumRatio.MulInt64(fixedPointDecimals).RoundInt().
		ToDec().QuoInt64(fixedPointDecimals).Equal(sdk.OneDec()) {
		return errors.New("vesting total ratio must be one")
	}

	return nil
}

// MarshalJSON returns the JSON representation of a LazyGradedVestingAccount.
func (lgva LazyGradedVestingAccount) MarshalJSON() ([]byte, error) {
	alias := vestingAccountJSON{
		Address:          lgva.Address,
		Coins:            lgva.Coins,
		PubKey:           lgva.PubKey,
		AccountNumber:    lgva.AccountNumber,
		Sequence:         lgva.Sequence,
		OriginalVesting:  lgva.OriginalVesting,
		DelegatedFree:    lgva.DelegatedFree,
		DelegatedVesting: lgva.DelegatedVesting,
		EndTime:          lgva.EndTime,
		VestingSchedules: lgva.VestingSchedules,
	}

	return internalCdc.MarshalJSON(alias)
}

// UnmarshalJSON unmarshals raw JSON bytes into a LazyGradedVestingAccount.
func (lgva *LazyGradedVestingAccount) UnmarshalJSON(bz []byte) error {
	var alias vestingAccountJSON
	if err := internalCdc.UnmarshalJSON(bz, &alias); err != nil {
		return err
	}

	lgva.BaseVestingAccount = &v039auth.BaseVestingAccount{
		BaseAccount:      v039auth.NewBaseAccount(alias.Address, alias.Coins, alias.PubKey, alias.AccountNumber, alias.Sequence),
		OriginalVesting:  alias.OriginalVesting,
		DelegatedFree:    alias.DelegatedFree,
		DelegatedVesting: alias.DelegatedVesting,
		EndTime:          alias.EndTime,
	}

	lgva.VestingSchedules = alias.VestingSchedules

	return nil
}

// Address implements cryptotypes.PubKey Address method
func (*LegacyAminoPubKey) Address() cryptotypes.Address {
	return nil
}

// Bytes no-lint
func (*LegacyAminoPubKey) Bytes() []byte { return nil }

// Equals no-lint
func (*LegacyAminoPubKey) Equals(key cryptotypes.PubKey) bool { return false }

// ProtoMessage no-lint
func (*LegacyAminoPubKey) ProtoMessage() {}

// Reset no-lint
func (*LegacyAminoPubKey) Reset() {}

// String no-lint
func (*LegacyAminoPubKey) String() string { return "not implemented" }

// Type no-lint
func (*LegacyAminoPubKey) Type() string { return "PubKeyMultisigThreshold" }

// VerifySignature no-lint
func (*LegacyAminoPubKey) VerifySignature(msg []byte, sig []byte) bool {
	panic("not implemented")
}

var internalCdc *codec.LegacyAmino

func init() {
	internalCdc = codec.NewLegacyAmino()
	RegisterLegacyAminoCodec(internalCdc)
}

// RegisterLegacyAminoCodec nonlint
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	registerCrypto(cdc)
	cdc.RegisterInterface((*v038auth.GenesisAccount)(nil), nil)
	cdc.RegisterInterface((*v038auth.Account)(nil), nil)
	cdc.RegisterConcrete(&v039auth.BaseAccount{}, "core/Account", nil)
	cdc.RegisterConcrete(&v039auth.BaseVestingAccount{}, "core/BaseVestingAccount", nil)
	cdc.RegisterConcrete(&LazyGradedVestingAccount{}, "core/LazyGradedVestingAccount", nil)
	cdc.RegisterConcrete(&v039auth.ModuleAccount{}, "supply/ModuleAccount", nil)
}

// registerCrypto registers all crypto dependency types with the provided Amino
// codec.
func registerCrypto(cdc *codec.LegacyAmino) {
	cdc.RegisterInterface((*cryptotypes.PubKey)(nil), nil)
	cdc.RegisterConcrete(sr25519.PubKey{},
		sr25519.PubKeyName, nil)
	cdc.RegisterConcrete(&ed25519.PubKey{},
		ed25519.PubKeyName, nil)
	cdc.RegisterConcrete(&secp256k1.PubKey{},
		secp256k1.PubKeyName, nil)
	cdc.RegisterConcrete(&LegacyAminoPubKey{},
		kmultisig.PubKeyAminoRoute, nil)

	cdc.RegisterInterface((*cryptotypes.PrivKey)(nil), nil)
	cdc.RegisterConcrete(sr25519.PrivKey{},
		sr25519.PrivKeyName, nil)
	cdc.RegisterConcrete(&ed25519.PrivKey{}, //nolint:staticcheck
		ed25519.PrivKeyName, nil)
	cdc.RegisterConcrete(&secp256k1.PrivKey{},
		secp256k1.PrivKeyName, nil)
}
