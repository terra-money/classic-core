package nameservice

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/terra-project/core/x/market"
	"github.com/terra-project/core/x/nameservice/internal/keeper"
	"github.com/terra-project/core/x/nameservice/internal/types"
	"testing"
)

func TestMarketFilters(t *testing.T) {
	input, h := setup(t)

	// Case 1: non-oracle message being sent fails
	bankMsg := bank.MsgSend{}
	res := h(input.Ctx, bankMsg)
	require.False(t, res.IsOK())

	// Case 2: Normal MsgSwap submission goes through
	prevoteMsg := NewMsgOpenAuction("valid.terra", keeper.Addrs[0])
	res = h(input.Ctx, prevoteMsg)
	require.True(t, res.IsOK())
}

func TestHandle_MsgOpenAuction(t *testing.T) {
	input, h := setup(t)

	params := input.NameserviceKeeper.GetParams(input.Ctx)

	validName := Name("wallet.terra")
	validNameHash, _ := validName.NameHash()
	expectedEndTime := input.Ctx.BlockTime().Add(params.BidPeriod)

	// normal auction open
	msg := NewMsgOpenAuction(validName, keeper.Addrs[0])
	res := h(input.Ctx, msg)
	require.True(t, res.IsOK())

	auction, err := input.NameserviceKeeper.GetAuction(input.Ctx, validNameHash)
	require.NoError(t, err)
	require.Equal(t, NewAuction(validName, AuctionStatusBid, expectedEndTime), auction)
	input.NameserviceKeeper.IterateBidAuctionQueue(input.Ctx, expectedEndTime, func(nameHash NameHash, auction Auction) bool {
		require.Equal(t, validNameHash, nameHash)
		return true
	})

	// invalid root name
	invalidRootName := Name("wallet.luna")
	msg = NewMsgOpenAuction(invalidRootName, keeper.Addrs[0])
	res = h(input.Ctx, msg)
	require.False(t, res.IsOK())
	require.Equal(t, res.Code, CodeInvalidRootName)

	// invalid length name
	invalidLengthName := Name("t.terra")
	msg = NewMsgOpenAuction(invalidLengthName, keeper.Addrs[0])
	res = h(input.Ctx, msg)
	require.False(t, res.IsOK())
	require.Equal(t, res.Code, CodeInvalidNameLength)

	// auction exists error
	msg = NewMsgOpenAuction(validName, keeper.Addrs[0])
	res = h(input.Ctx, msg)
	require.False(t, res.IsOK())
	require.Equal(t, res.Code, CodeAuctionExists)

	// registry exists error
	validName2 := Name("harvest.terra")
	validNameHash2, _ := validName2.NameHash()
	input.NameserviceKeeper.SetRegistry(input.Ctx, validNameHash2, NewRegistry(validName2, keeper.Addrs[0], expectedEndTime))
	msg = NewMsgOpenAuction(validName2, keeper.Addrs[0])
	res = h(input.Ctx, msg)
	require.False(t, res.IsOK())
	require.Equal(t, res.Code, types.CodeNameAlreadyTaken)
}

func TestHandle_MsgBidAuction(t *testing.T) {
	input, h := setup(t)

	params := input.NameserviceKeeper.GetParams(input.Ctx)

	salt := "salt"
	validName := Name("wallet.terra")
	validNameHash, _ := validName.NameHash()
	bidAmount := params.MinDeposit
	deposit := params.MinDeposit.Add(bidAmount)
	endTime := input.Ctx.BlockTime().Add(params.BidPeriod)

	bidHash := GetBidHash(salt, validName, bidAmount, keeper.Addrs[0])

	// register auction
	auction := NewAuction(validName, AuctionStatusBid, endTime)
	input.NameserviceKeeper.SetAuction(input.Ctx, validNameHash, NewAuction(validName, AuctionStatusBid, endTime))

	// valid bid
	msg := NewMsgBidAuction(validName, bidHash, deposit, keeper.Addrs[0])
	res := h(input.Ctx, msg)
	require.True(t, res.IsOK())

	bid, err := input.NameserviceKeeper.GetBid(input.Ctx, validNameHash, keeper.Addrs[0])
	require.NoError(t, err)
	require.Equal(t, NewBid(bidHash, deposit, keeper.Addrs[0]), bid)
	require.Equal(t, keeper.InitCoins.Sub(sdk.NewCoins(deposit)), input.BankKeeper.GetCoins(input.Ctx, keeper.Addrs[0]))
	require.Equal(t, sdk.NewCoins(deposit), input.SupplyKeeper.GetModuleAccount(input.Ctx, ModuleName).GetCoins())

	// different root name gives different name hash
	invalidName := Name("wallet.luna")
	msg = NewMsgBidAuction(invalidName, bidHash, deposit, keeper.Addrs[0])
	res = h(input.Ctx, msg)
	require.False(t, res.IsOK())
	require.Equal(t, CodeAuctionNotExists, res.Code)

	// auction not exists error
	validName2 := Name("wallet2.terra")
	msg = NewMsgBidAuction(validName2, bidHash, deposit, keeper.Addrs[0])
	res = h(input.Ctx, msg)
	require.False(t, res.IsOK())
	require.Equal(t, CodeAuctionNotExists, res.Code)

	// bid already exists error
	msg = NewMsgBidAuction(validName, bidHash, deposit, keeper.Addrs[0])
	res = h(input.Ctx, msg)
	require.False(t, res.IsOK())
	require.Equal(t, CodeBidAlreadyExists, res.Code)

	// auction is not bid status error
	auction.Status = AuctionStatusReveal
	input.NameserviceKeeper.SetAuction(input.Ctx, validNameHash, auction)
	input.NameserviceKeeper.DeleteBid(input.Ctx, validNameHash, keeper.Addrs[0])
	msg = NewMsgBidAuction(validName, bidHash, deposit, keeper.Addrs[0])
	res = h(input.Ctx, msg)
	require.False(t, res.IsOK())
	require.Equal(t, CodeAuctionNotBidStatus, res.Code)

	// min deposit error
	deposit.Amount = sdk.OneInt()
	auction.Status = AuctionStatusBid
	input.NameserviceKeeper.SetAuction(input.Ctx, validNameHash, auction)
	msg = NewMsgBidAuction(validName, bidHash, deposit, keeper.Addrs[0])
	res = h(input.Ctx, msg)
	require.False(t, res.IsOK())
	require.Equal(t, types.CodeLowDeposit, res.Code)
}

func TestHandle_MsgRevealBid(t *testing.T) {
	input, h := setup(t)

	params := input.NameserviceKeeper.GetParams(input.Ctx)

	salt := "salt"
	validName := Name("wallet.terra")
	validNameHash, _ := validName.NameHash()
	bidAmount := params.MinDeposit
	deposit := bidAmount.Add(bidAmount)
	endTime := input.Ctx.BlockTime().Add(params.BidPeriod)

	rawBidHash := GetBidHash(salt, validName, bidAmount, keeper.Addrs[0])

	input.NameserviceKeeper.SetAuction(input.Ctx, validNameHash, NewAuction(validName, AuctionStatusReveal, endTime))
	input.NameserviceKeeper.SetBid(input.Ctx, validNameHash, NewBid(rawBidHash, deposit, keeper.Addrs[0]))
	err := input.SupplyKeeper.SendCoinsFromAccountToModule(input.Ctx, keeper.Addrs[0], ModuleName, sdk.NewCoins(deposit))
	require.NoError(t, err)

	// valid reveal
	msg := NewMsgRevealBid(validName, salt, bidAmount, keeper.Addrs[0])
	res := h(input.Ctx, msg)

	require.True(t, res.IsOK())
	require.Equal(t, keeper.InitCoins.Sub(sdk.NewCoins(bidAmount)), input.BankKeeper.GetCoins(input.Ctx, keeper.Addrs[0]))
	require.Equal(t, sdk.NewCoins(bidAmount), input.SupplyKeeper.GetModuleAccount(input.Ctx, ModuleName).GetCoins())

	auction, err := input.NameserviceKeeper.GetAuction(input.Ctx, validNameHash)
	require.NoError(t, err)
	require.Equal(t, keeper.Addrs[0], auction.TopBidder)
	require.Equal(t, sdk.NewCoins(bidAmount), auction.TopBidAmount)

	// different root name gives different name hash
	invalidName := Name("wallet.luna")
	msg = NewMsgRevealBid(invalidName, salt, bidAmount, keeper.Addrs[0])
	res = h(input.Ctx, msg)
	require.False(t, res.IsOK())
	require.Equal(t, CodeAuctionNotExists, res.Code)

	// auction not exists error
	validName2 := Name("wallet2.terra")
	msg = NewMsgRevealBid(validName2, salt, bidAmount, keeper.Addrs[0])
	res = h(input.Ctx, msg)
	require.False(t, res.IsOK())
	require.Equal(t, CodeAuctionNotExists, res.Code)

	// bid not exists error
	msg = NewMsgRevealBid(validName, salt, bidAmount, keeper.Addrs[1])
	res = h(input.Ctx, msg)
	require.False(t, res.IsOK())
	require.Equal(t, CodeBidNotExists, res.Code)

	// hash validation error
	salt2 := "s2"
	input.NameserviceKeeper.SetBid(input.Ctx, validNameHash, NewBid(rawBidHash, deposit, keeper.Addrs[0]))

	msg = NewMsgRevealBid(validName, salt2, bidAmount, keeper.Addrs[0])
	res = h(input.Ctx, msg)
	require.False(t, res.IsOK())
	require.Equal(t, CodeVerificationFailed, res.Code)

	// deposit amount is smaller than bid amount
	bidAmount2 := bidAmount.Add(deposit)
	rawBidHash2 := GetBidHash(salt, validName, bidAmount2, keeper.Addrs[1])
	input.NameserviceKeeper.SetBid(input.Ctx, validNameHash, NewBid(rawBidHash2, deposit, keeper.Addrs[1]))

	msg = NewMsgRevealBid(validName, salt, bidAmount2, keeper.Addrs[1])
	res = h(input.Ctx, msg)
	require.False(t, res.IsOK())
	require.Equal(t, CodeDepositSmallerThanBidAmount, res.Code)

	// valid reveal from addr 1 (win top bidder)
	deposit2 := bidAmount2

	input.NameserviceKeeper.SetBid(input.Ctx, validNameHash, NewBid(rawBidHash2, deposit2, keeper.Addrs[1]))
	err = input.SupplyKeeper.SendCoinsFromAccountToModule(input.Ctx, keeper.Addrs[1], ModuleName, sdk.NewCoins(deposit2))
	require.NoError(t, err)

	msg = NewMsgRevealBid(validName, salt, bidAmount2, keeper.Addrs[1])
	res = h(input.Ctx, msg)
	require.True(t, res.IsOK())

	auction, err = input.NameserviceKeeper.GetAuction(input.Ctx, validNameHash)
	require.NoError(t, err)
	require.Equal(t, keeper.Addrs[1], auction.TopBidder)
	require.Equal(t, sdk.NewCoins(bidAmount2), auction.TopBidAmount)
	require.Equal(t, keeper.InitCoins, input.BankKeeper.GetCoins(input.Ctx, keeper.Addrs[0]))
	require.Equal(t, keeper.InitCoins.Sub(sdk.NewCoins(bidAmount2)), input.BankKeeper.GetCoins(input.Ctx, keeper.Addrs[1]))

	// valid reveal from addr 2 (lose)
	rawBidHash3 := GetBidHash(salt, validName, bidAmount, keeper.Addrs[2])
	input.NameserviceKeeper.SetBid(input.Ctx, validNameHash, NewBid(rawBidHash3, deposit, keeper.Addrs[2]))
	err = input.SupplyKeeper.SendCoinsFromAccountToModule(input.Ctx, keeper.Addrs[2], ModuleName, sdk.NewCoins(deposit))
	require.NoError(t, err)

	msg = NewMsgRevealBid(validName, salt, bidAmount, keeper.Addrs[2])
	res = h(input.Ctx, msg)
	require.True(t, res.IsOK())

	auction, err = input.NameserviceKeeper.GetAuction(input.Ctx, validNameHash)
	require.NoError(t, err)
	require.Equal(t, keeper.Addrs[1], auction.TopBidder)
	require.Equal(t, sdk.NewCoins(bidAmount2), auction.TopBidAmount)
	require.Equal(t, keeper.InitCoins, input.BankKeeper.GetCoins(input.Ctx, keeper.Addrs[2]))
	require.Equal(t, keeper.InitCoins.Sub(sdk.NewCoins(bidAmount2)), input.BankKeeper.GetCoins(input.Ctx, keeper.Addrs[1]))
}

func TestHandle_MsgRenewRegistry(t *testing.T) {
	input, h := setup(t)

	params := input.NameserviceKeeper.GetParams(input.Ctx)

	validName := Name("wallet.terra")
	validNameHash, _ := validName.NameHash()
	endTime := input.Ctx.BlockTime().Add(params.BidPeriod)
	fee := sdk.NewCoins(params.RenewalFees.RenewalFeeForLength(5))

	registry := NewRegistry(validName, keeper.Addrs[0], endTime)
	input.NameserviceKeeper.SetRegistry(input.Ctx, validNameHash, registry)

	// valid renew
	msg := NewMsgRenewRegistry(validName, fee, keeper.Addrs[0])
	res := h(input.Ctx, msg)
	require.True(t, res.IsOK())
	require.Equal(t, keeper.InitCoins.Sub(fee), input.BankKeeper.GetCoins(input.Ctx, keeper.Addrs[0]))

	registry, err := input.NameserviceKeeper.GetRegistry(input.Ctx, validNameHash)
	require.NoError(t, err)
	require.Equal(t, endTime.Add(params.RenewalInterval), registry.EndTime)
	input.NameserviceKeeper.IterateActiveRegistryQueue(input.Ctx, endTime, func(_ NameHash, _ Registry) bool {
		assert.Fail(t, "this entry must be deleted after renew registry")
		return true
	})
	input.NameserviceKeeper.IterateActiveRegistryQueue(input.Ctx, endTime.Add(params.RenewalInterval), func(nameHash NameHash, _ Registry) bool {
		require.Equal(t, validNameHash, nameHash)
		return true
	})

	// different root name gives different name hash
	invalidName := Name("wallet.luna")
	msg = NewMsgRenewRegistry(invalidName, fee, keeper.Addrs[0])
	res = h(input.Ctx, msg)
	require.False(t, res.IsOK())
	require.Equal(t, CodeRegistryNotExists, res.Code)

	// registry not exists
	validName2 := Name("harvest.terra")
	msg = NewMsgRenewRegistry(validName2, fee, keeper.Addrs[0])
	res = h(input.Ctx, msg)
	require.False(t, res.IsOK())
	require.Equal(t, CodeRegistryNotExists, res.Code)

	// permission error
	msg = NewMsgRenewRegistry(validName, fee, keeper.Addrs[1])
	res = h(input.Ctx, msg)
	require.False(t, res.IsOK())
	require.Equal(t, sdk.CodeUnauthorized, res.Code)

	// no registered swap rate error
	fee2 := fee
	fee2[0].Denom = "foo"

	msg = NewMsgRenewRegistry(validName, fee2, keeper.Addrs[0])
	res = h(input.Ctx, msg)
	require.False(t, res.IsOK())
	require.Equal(t, market.CodeNoEffectivePrice, res.Code)
}

func TestHandle_MsgUpdateOwner(t *testing.T) {
	input, h := setup(t)

	params := input.NameserviceKeeper.GetParams(input.Ctx)

	validName := Name("wallet.terra")
	validNameHash, _ := validName.NameHash()
	endTime := input.Ctx.BlockTime().Add(params.BidPeriod)

	registry := NewRegistry(validName, keeper.Addrs[0], endTime)
	input.NameserviceKeeper.SetRegistry(input.Ctx, validNameHash, registry)

	// valid update
	msg := NewMsgUpdateOwner(validName, keeper.Addrs[1], keeper.Addrs[0])
	res := h(input.Ctx, msg)
	require.True(t, res.IsOK())

	registry, err := input.NameserviceKeeper.GetRegistry(input.Ctx, validNameHash)
	require.NoError(t, err)
	require.Equal(t, keeper.Addrs[1], registry.Owner)

	// different root name gives different name hash
	invalidName := Name("wallet.luna")
	msg = NewMsgUpdateOwner(invalidName, keeper.Addrs[1], keeper.Addrs[0])
	res = h(input.Ctx, msg)
	require.False(t, res.IsOK())
	require.Equal(t, CodeRegistryNotExists, res.Code)

	// registry not exists
	validName2 := Name("harvest.terra")
	msg = NewMsgUpdateOwner(validName2, keeper.Addrs[1], keeper.Addrs[0])
	res = h(input.Ctx, msg)
	require.False(t, res.IsOK())
	require.Equal(t, CodeRegistryNotExists, res.Code)

	// permission error
	msg = NewMsgUpdateOwner(validName, keeper.Addrs[1], keeper.Addrs[0])
	res = h(input.Ctx, msg)
	require.False(t, res.IsOK())
	require.Equal(t, sdk.CodeUnauthorized, res.Code)
}

func TestHandle_MsgRegisterSubName(t *testing.T) {
	input, h := setup(t)

	params := input.NameserviceKeeper.GetParams(input.Ctx)

	validName := Name("wallet.terra")
	validNameHash, _ := validName.NameHash()
	endTime := input.Ctx.BlockTime().Add(params.RenewalInterval)

	registry := NewRegistry(validName, keeper.Addrs[0], endTime)
	input.NameserviceKeeper.SetRegistry(input.Ctx, validNameHash, registry)

	// valid sub-name register
	msg := NewMsgRegisterSubName(validName, keeper.Addrs[1], keeper.Addrs[0])
	res := h(input.Ctx, msg)
	require.True(t, res.IsOK())

	resolve, err := input.NameserviceKeeper.GetResolve(input.Ctx, validNameHash, []byte{})
	require.NoError(t, err)
	require.Equal(t, keeper.Addrs[1], resolve)

	reverse, err := input.NameserviceKeeper.GetReverseResolve(input.Ctx, keeper.Addrs[1])
	require.NoError(t, err)
	require.Equal(t, validNameHash, reverse)

	// different root name gives different name hash
	invalidName := Name("wallet.luna")
	msg = NewMsgRegisterSubName(invalidName, keeper.Addrs[2], keeper.Addrs[0])
	res = h(input.Ctx, msg)
	require.False(t, res.IsOK())
	require.Equal(t, CodeRegistryNotExists, res.Code)

	// registry not exists
	validName2 := Name("harvest.terra")
	msg = NewMsgRegisterSubName(validName2, keeper.Addrs[1], keeper.Addrs[0])
	res = h(input.Ctx, msg)
	require.False(t, res.IsOK())
	require.Equal(t, CodeRegistryNotExists, res.Code)

	// unauthorized
	validName3 := Name("user1.wallet.terra")
	msg = NewMsgRegisterSubName(validName3, keeper.Addrs[2], keeper.Addrs[1])
	res = h(input.Ctx, msg)
	require.False(t, res.IsOK())
	require.Equal(t, sdk.CodeUnauthorized, res.Code)

	// address already in use
	msg = NewMsgRegisterSubName(validName3, keeper.Addrs[1], keeper.Addrs[0])
	res = h(input.Ctx, msg)
	require.False(t, res.IsOK())
	require.Equal(t, CodeAddressAlreadyRegistered, res.Code)

	// name already in use
	msg = NewMsgRegisterSubName(validName, keeper.Addrs[2], keeper.Addrs[0])
	res = h(input.Ctx, msg)
	require.False(t, res.IsOK())
	require.Equal(t, CodeNameAlreadyTaken, res.Code)
}

func TestHandle_MsgUnregisterSubName(t *testing.T) {
	input, h := setup(t)

	params := input.NameserviceKeeper.GetParams(input.Ctx)

	validName := Name("wallet.terra")
	validNameHash, _ := validName.NameHash()
	endTime := input.Ctx.BlockTime().Add(params.RenewalInterval)

	registry := NewRegistry(validName, keeper.Addrs[0], endTime)
	input.NameserviceKeeper.SetRegistry(input.Ctx, validNameHash, registry)
	input.NameserviceKeeper.SetResolve(input.Ctx, validNameHash, []byte{}, keeper.Addrs[0])
	input.NameserviceKeeper.SetReverseResolve(input.Ctx, keeper.Addrs[0], validNameHash)

	// valid sub-name unregister
	msg := NewMsgUnregisterSubName(validName, keeper.Addrs[0])
	res := h(input.Ctx, msg)
	require.True(t, res.IsOK())

	_, err := input.NameserviceKeeper.GetResolve(input.Ctx, validNameHash, []byte{})
	require.Error(t, err)

	_, err = input.NameserviceKeeper.GetReverseResolve(input.Ctx, keeper.Addrs[0])
	require.Error(t, err)

	// register it again
	input.NameserviceKeeper.SetResolve(input.Ctx, validNameHash, []byte{}, keeper.Addrs[0])
	input.NameserviceKeeper.SetReverseResolve(input.Ctx, keeper.Addrs[0], validNameHash)

	// different root name gives different name hash
	invalidName := Name("wallet.luna")
	msg = NewMsgUnregisterSubName(invalidName, keeper.Addrs[0])
	res = h(input.Ctx, msg)
	require.False(t, res.IsOK())
	require.Equal(t, CodeRegistryNotExists, res.Code)

	// registry not exists
	validName2 := Name("harvest.terra")
	msg = NewMsgUnregisterSubName(validName2, keeper.Addrs[0])
	res = h(input.Ctx, msg)
	require.False(t, res.IsOK())
	require.Equal(t, CodeRegistryNotExists, res.Code)

	// unauthorized
	validName3 := Name("user1.wallet.terra")
	msg = NewMsgUnregisterSubName(validName3, keeper.Addrs[1])
	res = h(input.Ctx, msg)
	require.False(t, res.IsOK())
	require.Equal(t, sdk.CodeUnauthorized, res.Code)

	// name not in use
	msg = NewMsgUnregisterSubName(validName3, keeper.Addrs[0])
	res = h(input.Ctx, msg)
	require.False(t, res.IsOK())
	require.Equal(t, CodeResolveNotExists, res.Code)
}
