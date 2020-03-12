package nameservice

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/terra-project/core/x/nameservice/internal/keeper"
)

func TestExportInitGenesis(t *testing.T) {
	input := keeper.CreateTestInput(t)
	name := Name("wallet.terra")
	nameHash, childNameHash := name.NameHash()
	endTime := time.Now().UTC()
	// insert auction
	input.NameserviceKeeper.SetAuction(input.Ctx, nameHash, NewAuction(name, AuctionStatusBid, endTime))
	input.NameserviceKeeper.InsertBidAuctionQueue(input.Ctx, nameHash, endTime)

	// insert bid
	bidAmount := sdk.NewInt64Coin("foo", 123)
	bidHash := GetBidHash("salt", name, bidAmount, keeper.Addrs[0])
	input.NameserviceKeeper.SetBid(input.Ctx, nameHash, NewBid(bidHash, bidAmount, keeper.Addrs[0]))

	// insert registry
	input.NameserviceKeeper.SetRegistry(input.Ctx, nameHash, NewRegistry(name, keeper.Addrs[0], endTime))
	input.NameserviceKeeper.InsertActiveRegistryQueue(input.Ctx, nameHash, endTime)

	// insert resolve
	input.NameserviceKeeper.SetResolve(input.Ctx, nameHash, childNameHash, keeper.Addrs[0])
	input.NameserviceKeeper.SetReverseResolve(input.Ctx, keeper.Addrs[0], nameHash)

	genesis := ExportGenesis(input.Ctx, input.NameserviceKeeper)

	newInput := keeper.CreateTestInput(t)
	InitGenesis(newInput.Ctx, newInput.NameserviceKeeper, genesis)
	newGenesis := ExportGenesis(newInput.Ctx, newInput.NameserviceKeeper)

	require.Equal(t, genesis, newGenesis)
}
