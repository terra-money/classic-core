package types

import (
	"encoding/hex"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto/tmhash"
)

// ensure Msg interface compliance at compile time
var (
	_ sdk.Msg = &MsgOpenAuction{}
	_ sdk.Msg = &MsgBidAuction{}
	_ sdk.Msg = &MsgRevealBid{}
	_ sdk.Msg = &MsgRenewRegistry{}
	_ sdk.Msg = &MsgRegisterSubName{}
	_ sdk.Msg = &MsgUnregisterSubName{}
)

// MsgOpenAuction - struct for opening the name auction .
type MsgOpenAuction struct {
	Name      Name           `json:"name" yaml:"name"`
	Organizer sdk.AccAddress `json:"organizer" yaml:"organizer"`
}

// NewMsgOpenAuction creates a MsgOpenAuction instance
func NewMsgOpenAuction(name Name, organizer sdk.AccAddress) MsgOpenAuction {
	return MsgOpenAuction{
		Name:      name,
		Organizer: organizer,
	}
}

// Route implements sdk.Msg
func (msg MsgOpenAuction) Route() string { return RouterKey }

// Type implements sdk.Msg
func (msg MsgOpenAuction) Type() string { return "openauction" }

// GetSignBytes implements sdk.Msg
func (msg MsgOpenAuction) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners implements sdk.Msg
func (msg MsgOpenAuction) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Organizer}
}

// ValidateBasic Implements sdk.Msg
func (msg MsgOpenAuction) ValidateBasic() sdk.Error {
	if err := msg.Name.Validate(); err != nil {
		return ErrInvalidName(DefaultCodespace, msg.Name, err.Error())
	}

	if msg.Name.Levels() != 2 {
		return ErrInvalidName(DefaultCodespace, msg.Name, "only second level name is accepted for auction")
	}

	if msg.Organizer.Empty() {
		return sdk.ErrInvalidAddress(msg.Organizer.String())
	}

	return nil
}

// String implements fmt.Stringer interface
func (msg MsgOpenAuction) String() string {
	return fmt.Sprintf(`MsgOpenAuction
	Name:         %s,
	Organizer:    %s`,
		msg.Name, msg.Organizer)
}

// MsgBidAuction - struct for bidding the name auction
// The purpose of hash is to hide bidding amount
// which is formatted as hex string in first 20 bytes of SHA256("salt:name:amount:bidder")
type MsgBidAuction struct {
	Name    Name           `json:"name" yaml:"name"`
	Hash    string         `json:"hash" yaml:"hash"`
	Deposit sdk.Coin       `json:"deposit" yaml:"deposit"`
	Bidder  sdk.AccAddress `json:"bidder" yaml:"bidder"`
}

// NewMsgBidAuction creates a MsgBidAuction instance
func NewMsgBidAuction(name Name, hash string, deposit sdk.Coin, bidder sdk.AccAddress) MsgBidAuction {
	return MsgBidAuction{
		Name:    name,
		Hash:    hash,
		Deposit: deposit,
		Bidder:  bidder,
	}
}

// Route implements sdk.Msg
func (msg MsgBidAuction) Route() string { return RouterKey }

// Type implements sdk.Msg
func (msg MsgBidAuction) Type() string { return "bidauction" }

// GetSignBytes implements sdk.Msg
func (msg MsgBidAuction) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners implements sdk.Msg
func (msg MsgBidAuction) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Bidder}
}

// ValidateBasic Implements sdk.Msg
func (msg MsgBidAuction) ValidateBasic() sdk.Error {

	if err := msg.Name.Validate(); err != nil {
		return ErrInvalidName(DefaultCodespace, msg.Name, err.Error())
	}

	if msg.Name.Levels() != 2 {
		return ErrInvalidName(DefaultCodespace, msg.Name, "only second level name is accepted for auction")
	}

	if bz, err := hex.DecodeString(msg.Hash); err != nil {
		return sdk.ErrInternal(fmt.Sprintf("Failed to decode hash: %s", err.Error()))
	} else if len(bz) != tmhash.TruncatedSize {
		return ErrInvalidHashLength(DefaultCodespace, len(bz))
	}

	if msg.Bidder.Empty() {
		return sdk.ErrInvalidAddress(msg.Bidder.String())
	}

	if !msg.Deposit.IsValid() {
		return sdk.ErrInvalidCoins(msg.Deposit.String())
	}

	return nil
}

// String implements fmt.Stringer interface
func (msg MsgBidAuction) String() string {
	return fmt.Sprintf(`MsgBidAuction
	Name:    %s,
	Hash:    %s,
    Deposit: %s,
    Bidder:  %s`,
		msg.Name, msg.Hash, msg.Deposit, msg.Bidder)
}

// MsgRevealBid - struct for revealing the proof of the bidding msg
type MsgRevealBid struct {
	Name   Name           `json:"name" yaml:"name"`
	Salt   string         `json:"salt" yaml:"salt"`
	Amount sdk.Coin       `json:"amount" yaml:"amount"`
	Bidder sdk.AccAddress `json:"bidder" yaml:"bidder"`
}

// NewMsgRevealBid creates a MsgRevealBid instance
func NewMsgRevealBid(name Name, salt string, amount sdk.Coin, bidder sdk.AccAddress) MsgRevealBid {
	return MsgRevealBid{
		Name:   name,
		Salt:   salt,
		Amount: amount,
		Bidder: bidder,
	}
}

// Route implements sdk.Msg
func (msg MsgRevealBid) Route() string { return RouterKey }

// Type implements sdk.Msg
func (msg MsgRevealBid) Type() string { return "revealbid" }

// GetSignBytes implements sdk.Msg
func (msg MsgRevealBid) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners implements sdk.Msg
func (msg MsgRevealBid) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Bidder}
}

// ValidateBasic Implements sdk.Msg
func (msg MsgRevealBid) ValidateBasic() sdk.Error {

	if err := msg.Name.Validate(); err != nil {
		return ErrInvalidName(DefaultCodespace, msg.Name, err.Error())
	}

	if msg.Name.Levels() != 2 {
		return ErrInvalidName(DefaultCodespace, msg.Name, "only second level name is accepted for auction")
	}

	if len(msg.Salt) > 4 || len(msg.Salt) < 1 {
		return ErrInvalidSaltLength(DefaultCodespace, len(msg.Salt))
	}

	if msg.Bidder.Empty() {
		return sdk.ErrInvalidAddress(msg.Bidder.String())
	}

	if !msg.Amount.IsValid() {
		return sdk.ErrInvalidCoins(msg.Amount.String())
	}

	return nil
}

// String implements fmt.Stringer interface
func (msg MsgRevealBid) String() string {
	return fmt.Sprintf(`MsgRevealBid
	Name:    %s,
	Salt:    %s,
    Amount: %s,
    Bidder:  %s`,
		msg.Name, msg.Salt, msg.Amount, msg.Bidder)
}

// MsgRenewRegistry - struct for renewal of the name registry
// The renewal tax can be vary according to the length of name
type MsgRenewRegistry struct {
	Name  Name           `json:"name" yaml:"name"`
	Fee   sdk.Coins      `json:"fee" yaml:"fee"`
	Owner sdk.AccAddress `json:"owner" yaml:"owner"`
}

// NewMsgRenewRegistry creates a MsgRenewRegistry instance
func NewMsgRenewRegistry(name Name, fee sdk.Coins, owner sdk.AccAddress) MsgRenewRegistry {
	return MsgRenewRegistry{
		Name:  name,
		Fee:   fee,
		Owner: owner,
	}
}

// Route implements sdk.Msg
func (msg MsgRenewRegistry) Route() string { return RouterKey }

// Type implements sdk.Msg
func (msg MsgRenewRegistry) Type() string { return "renewregistry" }

// GetSignBytes implements sdk.Msg
func (msg MsgRenewRegistry) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners implements sdk.Msg
func (msg MsgRenewRegistry) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}

// ValidateBasic Implements sdk.Msg
func (msg MsgRenewRegistry) ValidateBasic() sdk.Error {

	if err := msg.Name.Validate(); err != nil {
		return ErrInvalidName(DefaultCodespace, msg.Name, err.Error())
	}

	if msg.Name.Levels() != 2 {
		return ErrInvalidName(DefaultCodespace, msg.Name, "only second level name is accepted for registry")
	}

	if !msg.Fee.IsValid() {
		return sdk.ErrInvalidCoins(msg.Fee.String())
	}

	if msg.Owner.Empty() {
		return sdk.ErrInvalidAddress(msg.Owner.String())
	}

	return nil
}

// String implements fmt.Stringer interface
func (msg MsgRenewRegistry) String() string {
	return fmt.Sprintf(`MsgRenewRegistry
	Name:    %s,
	Fee:     %s,
    Owner:   %s`,
		msg.Name, msg.Fee, msg.Owner)
}

// MsgUpdateOwner - struct for unregistering new sub name
type MsgUpdateOwner struct {
	Name     Name           `json:"name" yaml:"name"`
	NewOwner sdk.AccAddress `json:"new_owner" yaml:"new_owner"`
	Owner    sdk.AccAddress `json:"owner" yaml:"owner"`
}

// NewMsgUpdateOwner creates a MsgUpdateOwner instance
func NewMsgUpdateOwner(name Name, newOwner, owner sdk.AccAddress) MsgUpdateOwner {
	return MsgUpdateOwner{
		Name:     name,
		NewOwner: newOwner,
		Owner:    owner,
	}
}

// Route implements sdk.Msg
func (msg MsgUpdateOwner) Route() string { return RouterKey }

// Type implements sdk.Msg
func (msg MsgUpdateOwner) Type() string { return "updateowner" }

// GetSignBytes implements sdk.Msg
func (msg MsgUpdateOwner) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners implements sdk.Msg
func (msg MsgUpdateOwner) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}

// ValidateBasic Implements sdk.Msg
func (msg MsgUpdateOwner) ValidateBasic() sdk.Error {

	if err := msg.Name.Validate(); err != nil {
		return ErrInvalidName(DefaultCodespace, msg.Name, err.Error())
	}

	if levels := msg.Name.Levels(); !(levels == 2) {
		return ErrInvalidName(DefaultCodespace, msg.Name, "only second level name is accepted for register")
	}

	if msg.Owner.Empty() {
		return sdk.ErrInvalidAddress(msg.Owner.String())
	}

	if msg.NewOwner.Empty() {
		return sdk.ErrInvalidAddress(msg.NewOwner.String())
	}

	return nil
}

// String implements fmt.Stringer interface
func (msg MsgUpdateOwner) String() string {
	return fmt.Sprintf(`MsgUpdateOwner
	Name:        %s,
    Owner:       %s,
    NewOwner:    %s`,
		msg.Name, msg.Owner, msg.NewOwner)
}

// MsgRegisterSubName - struct for registering new sub name
type MsgRegisterSubName struct {
	Name    Name           `json:"name" yaml:"name"`
	Address sdk.AccAddress `json:"address" yaml:"address"`
	Owner   sdk.AccAddress `json:"owner" yaml:"owner"`
}

// NewMsgRegisterSubName creates a MsgRegisterSubName instance
func NewMsgRegisterSubName(name Name, address sdk.AccAddress, owner sdk.AccAddress) MsgRegisterSubName {
	return MsgRegisterSubName{
		Name:    name,
		Address: address,
		Owner:   owner,
	}
}

// Route implements sdk.Msg
func (msg MsgRegisterSubName) Route() string { return RouterKey }

// Type implements sdk.Msg
func (msg MsgRegisterSubName) Type() string { return "registersubname" }

// GetSignBytes implements sdk.Msg
func (msg MsgRegisterSubName) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners implements sdk.Msg
func (msg MsgRegisterSubName) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}

// ValidateBasic Implements sdk.Msg
func (msg MsgRegisterSubName) ValidateBasic() sdk.Error {
	if err := msg.Name.Validate(); err != nil {
		return ErrInvalidName(DefaultCodespace, msg.Name, err.Error())
	}

	if levels := msg.Name.Levels(); !(levels == 2 || levels == 3) {
		return ErrInvalidName(DefaultCodespace, msg.Name, "only the second or third level name is accepted for register")
	}

	if msg.Address.Empty() {
		return sdk.ErrInvalidAddress(msg.Address.String())
	}

	if msg.Owner.Empty() {
		return sdk.ErrInvalidAddress(msg.Owner.String())
	}

	return nil
}

// String implements fmt.Stringer interface
func (msg MsgRegisterSubName) String() string {
	return fmt.Sprintf(`MsgRegisterSubName
	Name:     %s,
    Address:  %s,
    Owner:    %s`,
		msg.Name, msg.Address, msg.Owner)
}

// MsgUnregisterSubName - struct for unregistering new sub name
type MsgUnregisterSubName struct {
	Name  Name           `json:"name" yaml:"name"`
	Owner sdk.AccAddress `json:"owner" yaml:"owner"`
}

// NewMsgUnregisterSubName creates a MsgUnregisterSubName instance
func NewMsgUnregisterSubName(name Name, owner sdk.AccAddress) MsgUnregisterSubName {
	return MsgUnregisterSubName{
		Name:  name,
		Owner: owner,
	}
}

// Route implements sdk.Msg
func (msg MsgUnregisterSubName) Route() string { return RouterKey }

// Type implements sdk.Msg
func (msg MsgUnregisterSubName) Type() string { return "unregistersubname" }

// GetSignBytes implements sdk.Msg
func (msg MsgUnregisterSubName) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners implements sdk.Msg
func (msg MsgUnregisterSubName) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}

// ValidateBasic Implements sdk.Msg
func (msg MsgUnregisterSubName) ValidateBasic() sdk.Error {

	if err := msg.Name.Validate(); err != nil {
		return ErrInvalidName(DefaultCodespace, msg.Name, err.Error())
	}

	if levels := msg.Name.Levels(); !(levels == 2 || levels == 3) {
		return ErrInvalidName(DefaultCodespace, msg.Name, "only the second or third level name is accepted for unregister")
	}

	if msg.Owner.Empty() {
		return sdk.ErrInvalidAddress(msg.Owner.String())
	}

	return nil
}

// String implements fmt.Stringer interface
func (msg MsgUnregisterSubName) String() string {
	return fmt.Sprintf(`MsgUnregisterSubName
	Name:        %s,
    Owner:       %s`,
		msg.Name, msg.Owner)
}
