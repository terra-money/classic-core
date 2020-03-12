package nameservice

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"testing"
	"time"

	"github.com/terra-project/core/x/nameservice/internal/keeper"
)

func TestIterate_BidAuction(t *testing.T) {
	input := keeper.CreateTestInput(t)

	params := input.NameserviceKeeper.GetParams(input.Ctx)
	blockTime := input.Ctx.BlockTime()

	validName := Name("wallet.terra")
	validName2 := Name("harvest.terra")
	validName3 := Name("chai.terra")
	validNameHash, _ := validName.NameHash()
	validNameHash2, _ := validName2.NameHash()
	validNameHash3, _ := validName3.NameHash()
	auction := NewAuction(validName, AuctionStatusBid, blockTime)
	auction2 := NewAuction(validName2, AuctionStatusBid, blockTime)
	auction3 := NewAuction(validName3, AuctionStatusBid, blockTime.Add(time.Second))

	input.NameserviceKeeper.SetAuction(input.Ctx, validNameHash, auction)
	input.NameserviceKeeper.SetAuction(input.Ctx, validNameHash2, auction2)
	input.NameserviceKeeper.SetAuction(input.Ctx, validNameHash3, auction3)

	input.NameserviceKeeper.InsertBidAuctionQueue(input.Ctx, validNameHash, auction.EndTime)
	input.NameserviceKeeper.InsertBidAuctionQueue(input.Ctx, validNameHash2, auction2.EndTime)
	input.NameserviceKeeper.InsertBidAuctionQueue(input.Ctx, validNameHash3, auction3.EndTime)

	// only auction 1 has bid
	input.NameserviceKeeper.SetBid(input.Ctx, validNameHash, NewBid([]byte{}, sdk.NewInt64Coin("foo", 123), keeper.Addrs[0]))

	EndBlocker(input.Ctx, input.NameserviceKeeper)

	// auction 1 have bid, so move to reveal period
	auction, err := input.NameserviceKeeper.GetAuction(input.Ctx, validNameHash)
	require.NoError(t, err)
	require.Equal(t, AuctionStatusReveal, auction.Status)
	require.Equal(t, blockTime.Add(params.RevealPeriod), auction.EndTime)

	// auction 2 have no bid, so close it
	_, err = input.NameserviceKeeper.GetAuction(input.Ctx, validNameHash2)
	require.Error(t, err)

	// auction 3 have left end time
	auction3, err = input.NameserviceKeeper.GetAuction(input.Ctx, validNameHash3)
	require.NoError(t, err)
	require.Equal(t, AuctionStatusBid, auction3.Status)
}

func TestIterate_RevealAuction(t *testing.T) {
	input := keeper.CreateTestInput(t)

	params := input.NameserviceKeeper.GetParams(input.Ctx)
	blockTime := input.Ctx.BlockTime()

	validName := Name("wallet.terra")
	validName2 := Name("harvest.terra")
	validName3 := Name("chai.terra")
	validNameHash, _ := validName.NameHash()
	validNameHash2, _ := validName2.NameHash()
	validNameHash3, _ := validName3.NameHash()
	auction := NewAuction(validName, AuctionStatusReveal, blockTime)
	auction.TopBidAmount = sdk.NewCoins(params.MinDeposit)
	auction.TopBidder = keeper.Addrs[0]

	auction2 := NewAuction(validName2, AuctionStatusReveal, blockTime)
	auction3 := NewAuction(validName3, AuctionStatusReveal, blockTime.Add(time.Second))

	input.NameserviceKeeper.SetAuction(input.Ctx, validNameHash, auction)
	input.NameserviceKeeper.SetAuction(input.Ctx, validNameHash2, auction2)
	input.NameserviceKeeper.SetAuction(input.Ctx, validNameHash3, auction3)

	input.NameserviceKeeper.InsertRevealAuctionQueue(input.Ctx, validNameHash, auction.EndTime)
	input.NameserviceKeeper.InsertRevealAuctionQueue(input.Ctx, validNameHash2, auction2.EndTime)
	input.NameserviceKeeper.InsertRevealAuctionQueue(input.Ctx, validNameHash3, auction3.EndTime)

	// lazy bidder (slashing)
	input.NameserviceKeeper.SetBid(input.Ctx, validNameHash, NewBid([]byte{}, params.MinDeposit, keeper.Addrs[1]))

	err := input.SupplyKeeper.SendCoinsFromAccountToModule(input.Ctx, keeper.Addrs[0], ModuleName, sdk.NewCoins(params.MinDeposit))
	require.NoError(t, err)

	err = input.SupplyKeeper.SendCoinsFromAccountToModule(input.Ctx, keeper.Addrs[1], ModuleName, sdk.NewCoins(params.MinDeposit))
	require.NoError(t, err)

	EndBlocker(input.Ctx, input.NameserviceKeeper)

	// auction 1 should be deleted and new registry for valid name 1 should be exists
	_, err = input.NameserviceKeeper.GetAuction(input.Ctx, validNameHash)
	require.Error(t, err)

	registry, err := input.NameserviceKeeper.GetRegistry(input.Ctx, validNameHash)
	require.NoError(t, err)
	require.Equal(t, NewRegistry(validName, keeper.Addrs[0], blockTime.Add(params.RenewalInterval)), registry)

	// lazy bidder slashing
	_, err = input.NameserviceKeeper.GetBid(input.Ctx, validNameHash, keeper.Addrs[1])
	require.Error(t, err)
	require.Equal(t, keeper.InitCoins.Sub(sdk.NewCoins(params.MinDeposit)), input.BankKeeper.GetCoins(input.Ctx, keeper.Addrs[1]))

	// auction 2 should be deleted and no registry will be added
	_, err = input.NameserviceKeeper.GetAuction(input.Ctx, validNameHash2)
	require.Error(t, err)

	_, err = input.NameserviceKeeper.GetRegistry(input.Ctx, validNameHash2)
	require.Error(t, err)

	// auction 3 should be exists
	_, err = input.NameserviceKeeper.GetAuction(input.Ctx, validNameHash3)
	require.NoError(t, err)
}

func TestIterate_ActiveRegistry(t *testing.T) {
	input := keeper.CreateTestInput(t)

	blockTime := input.Ctx.BlockTime()

	validName := Name("wallet.terra")
	validName2 := Name("harvest.terra")
	validNameHash, _ := validName.NameHash()
	validNameHash2, _ := validName2.NameHash()

	registry := NewRegistry(validName, keeper.Addrs[0], blockTime)
	registry2 := NewRegistry(validName2, keeper.Addrs[0], blockTime.Add(time.Second))

	input.NameserviceKeeper.SetRegistry(input.Ctx, validNameHash, registry)
	input.NameserviceKeeper.SetRegistry(input.Ctx, validNameHash2, registry2)

	input.NameserviceKeeper.InsertActiveRegistryQueue(input.Ctx, validNameHash, blockTime)
	input.NameserviceKeeper.InsertActiveRegistryQueue(input.Ctx, validNameHash, blockTime.Add(time.Second))

	input.NameserviceKeeper.SetResolve(input.Ctx, validNameHash, []byte{}, keeper.Addrs[0])
	input.NameserviceKeeper.SetReverseResolve(input.Ctx, keeper.Addrs[0], validNameHash)

	EndBlocker(input.Ctx, input.NameserviceKeeper)

	// registry 1 must be expired
	_, err := input.NameserviceKeeper.GetRegistry(input.Ctx, validNameHash)
	require.Error(t, err)

	_, err = input.NameserviceKeeper.GetResolve(input.Ctx, validNameHash, []byte{})
	require.Error(t, err)

	_, err = input.NameserviceKeeper.GetReverseResolve(input.Ctx, keeper.Addrs[0])
	require.Error(t, err)

	// registry 2 must be stay in the store
	_, err = input.NameserviceKeeper.GetRegistry(input.Ctx, validNameHash2)
	require.NoError(t, err)
}
