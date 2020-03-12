package nameservice

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/terra-project/core/x/nameservice/internal/types"
)

// NewHandler returns a handler for "oracle" type messages.
func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case MsgOpenAuction:
			return handleMsgOpenAuction(ctx, k, msg)
		case MsgBidAuction:
			return handleMsgBidAuction(ctx, k, msg)
		case MsgRevealBid:
			return handleMsgRevealBid(ctx, k, msg)
		case MsgRenewRegistry:
			return handleMsgRenewRegistry(ctx, k, msg)
		case MsgUpdateOwner:
			return handleMsgUpdateOwner(ctx, k, msg)
		case MsgRegisterSubName:
			return handleMsgRegisterSubName(ctx, k, msg)
		case MsgUnregisterSubName:
			return handleMsgUnregisterSubName(ctx, k, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized nameservice message type: %T", msg)
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// handleMsgOpenAuction handles a MsgOpenAuction
func handleMsgOpenAuction(ctx sdk.Context, keeper Keeper, msg MsgOpenAuction) sdk.Result {
	rootName, parentName, _ := msg.Name.Split()

	// check root name is same with params
	if paramRootName := keeper.RootName(ctx); paramRootName != rootName {
		return ErrInvalidRootName(keeper.Codespace(), paramRootName, rootName).Result()
	}

	// check name length
	if minNameLength := keeper.MinNameLength(ctx); len(parentName) < minNameLength {
		return ErrInvalidNameLength(keeper.Codespace(), minNameLength, len(parentName)).Result()
	}

	// do not validate name format because ante handler's msg.Validate() do name format validation
	nameHash, _ := msg.Name.NameHash()

	// check the name is not registered
	if _, err := keeper.GetRegistry(ctx, nameHash); err == nil {
		return ErrNameAlreadyTaken(keeper.Codespace()).Result()
	}

	// ensure there is no opened auction for this name
	if _, err := keeper.GetAuction(ctx, nameHash); err == nil {
		return ErrAuctionExists(keeper.Codespace()).Result()
	}

	curTime := ctx.BlockTime()
	bidEndTime := curTime.Add(keeper.BidPeriod(ctx))

	keeper.SetAuction(ctx, nameHash, NewAuction(msg.Name, AuctionStatusBid, bidEndTime))
	keeper.InsertBidAuctionQueue(ctx, nameHash, bidEndTime)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeOpen,
			sdk.NewAttribute(types.AttributeKeyName, msg.Name.String()),
			sdk.NewAttribute(types.AttributeKeyOrganizer, msg.Organizer.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
	})

	return sdk.Result{Events: ctx.EventManager().Events()}
}

// handleMsgBidAuction handles a MsgBidAuction
func handleMsgBidAuction(ctx sdk.Context, keeper Keeper, msg MsgBidAuction) sdk.Result {
	// do not validate name format because ante handler's msg.Validate() do name format validation
	// open auction handler also do min name length validation
	nameHash, _ := msg.Name.NameHash()

	// check the auction is exists
	auction, err := keeper.GetAuction(ctx, nameHash)
	if err != nil {
		return err.Result()
	}

	// check auction state is bid
	if auction.Status != types.AuctionStatusBid {
		return ErrAuctionNotBidStatus(keeper.Codespace(), auction.Status).Result()
	}

	// check the bidder did bid
	if _, err := keeper.GetBid(ctx, nameHash, msg.Bidder); err == nil {
		return ErrBidAlreadyExists(keeper.Codespace()).Result()
	}

	// check the deposit is greater than min deposit
	if minDeposit := keeper.MinDeposit(ctx); !sdk.NewCoins(msg.Deposit).IsAllGTE(sdk.NewCoins(minDeposit)) {
		return ErrLowDeposit(keeper.Codespace(), minDeposit).Result()
	}

	// transfer bidder deposit to module account
	err = keeper.SupplyKeeper.SendCoinsFromAccountToModule(ctx, msg.Bidder, ModuleName, sdk.NewCoins(msg.Deposit))
	if err != nil {
		return err.Result()
	}

	// parse bid hash
	bidHash, err2 := BidHashFromHexString(msg.Hash)
	if err2 != nil {
		return sdk.ErrUnknownRequest(err2.Error()).Result()
	}

	// store bid
	keeper.SetBid(ctx, nameHash, NewBid(bidHash, msg.Deposit, msg.Bidder))

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeBid,
			sdk.NewAttribute(types.AttributeKeyName, msg.Name.String()),
			sdk.NewAttribute(types.AttributeKeyBidder, msg.Bidder.String()),
			sdk.NewAttribute(types.AttributeKeyDeposit, msg.Deposit.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
	})

	return sdk.Result{Events: ctx.EventManager().Events()}
}

// handleMsgRevealBid handles a MsgRevealBid
func handleMsgRevealBid(ctx sdk.Context, keeper Keeper, msg MsgRevealBid) sdk.Result {
	nameHash, _ := msg.Name.NameHash()

	// check the auction is exists
	auction, err := keeper.GetAuction(ctx, nameHash)
	if err != nil {
		return err.Result()
	}

	// check auction state is reveal
	if auction.Status != types.AuctionStatusReveal {
		return ErrAuctionNotRevealStatus(keeper.Codespace(), auction.Status).Result()
	}

	// check bid is exists
	bid, err := keeper.GetBid(ctx, nameHash, msg.Bidder)
	if err != nil {
		return err.Result()
	}

	// check hash validation
	bidHash := GetBidHash(msg.Salt, msg.Name, msg.Amount, msg.Bidder)
	if !bid.Hash.Equal(bidHash) {
		return ErrVerificationFailed(keeper.Codespace(), bid.Hash, bidHash).Result()
	}

	// check deposit is not smaller than bid amount
	if !sdk.NewCoins(bid.Deposit).IsAllGTE(sdk.NewCoins(msg.Amount)) {
		return ErrDepositSmallerThanBidAmount(keeper.Codespace(), bid.Deposit).Result()
	}

	// check bid amount is greater than past top bidder
	if auction.TopBidder.Empty() || auction.TopBidAmount[0].IsLT(msg.Amount) {
		// new bidder is now top bidder
		// refund bid amount to past top bidder
		err = keeper.SupplyKeeper.SendCoinsFromModuleToAccount(ctx, ModuleName, auction.TopBidder, auction.TopBidAmount)
		if err != nil {
			return err.Result()
		}

		// update top bidder and top bid amount
		auction.TopBidder = msg.Bidder
		auction.TopBidAmount = sdk.NewCoins(msg.Amount)
		keeper.SetAuction(ctx, nameHash, auction)

		// refund left deposit to new top bidder except bid amount
		err = keeper.SupplyKeeper.SendCoinsFromModuleToAccount(ctx, ModuleName, msg.Bidder, sdk.NewCoins(bid.Deposit.Sub(msg.Amount)))
		if err != nil {
			return err.Result()
		}
	} else {
		// refund the deposit to lose bidder
		err = keeper.SupplyKeeper.SendCoinsFromModuleToAccount(ctx, ModuleName, bid.Bidder, sdk.NewCoins(bid.Deposit))
		if err != nil {
			return err.Result()
		}
	}

	// delete bid
	keeper.DeleteBid(ctx, nameHash, msg.Bidder)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeReveal,
			sdk.NewAttribute(types.AttributeKeyName, msg.Name.String()),
			sdk.NewAttribute(types.AttributeKeyBidder, msg.Bidder.String()),
			sdk.NewAttribute(types.AttributeKeyAmount, msg.Amount.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
	})

	return sdk.Result{Events: ctx.EventManager().Events()}
}

// handleMsgRenewRegistry handles a MsgRenewRegistry
func handleMsgRenewRegistry(ctx sdk.Context, keeper Keeper, msg MsgRenewRegistry) sdk.Result {
	_, parentName, _ := msg.Name.Split()
	nameHash, _ := msg.Name.NameHash()

	// check registry is exists
	registry, err := keeper.GetRegistry(ctx, nameHash)
	if err != nil {
		return err.Result()
	}

	// check permission
	if !registry.Owner.Equals(msg.Owner) {
		return sdk.ErrUnauthorized(fmt.Sprintf("registered owner is %s", registry.Owner)).Result()
	}

	// convert renewal fee to time duration
	// check swap validation
	extendedTime, err := keeper.ConvertRenewalFeeToTime(ctx, msg.Fee, len(parentName))
	if err != nil {
		return err.Result()
	}

	// transfer fee to module account
	err = keeper.SupplyKeeper.SendCoinsFromAccountToModule(ctx, msg.Owner, ModuleName, msg.Fee)
	if err != nil {
		return err.Result()
	}

	// TODO - how to handle fee
	err = keeper.SupplyKeeper.BurnCoins(ctx, ModuleName, msg.Fee)
	if err != nil {
		return err.Result()
	}

	// delete old item from active queue
	keeper.RemoveFromActiveRegistryQueue(ctx, nameHash, registry.EndTime)

	// update end time and insert to active queue again with update end time
	registry.EndTime = registry.EndTime.Add(extendedTime)
	keeper.SetRegistry(ctx, nameHash, registry)
	keeper.InsertActiveRegistryQueue(ctx, nameHash, registry.EndTime)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeRenew,
			sdk.NewAttribute(types.AttributeKeyName, msg.Name.String()),
			sdk.NewAttribute(types.AttributeKeyOwner, msg.Owner.String()),
			sdk.NewAttribute(types.AttributeKeyFee, msg.Fee.String()),
			sdk.NewAttribute(types.AttributeKeyEndTime, registry.EndTime.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
	})

	return sdk.Result{Events: ctx.EventManager().Events()}
}

// handleMsgUpdateOwner handles a MsgUpdateOwner
func handleMsgUpdateOwner(ctx sdk.Context, keeper Keeper, msg MsgUpdateOwner) sdk.Result {
	nameHash, _ := msg.Name.NameHash()

	// check registry is exists
	registry, err := keeper.GetRegistry(ctx, nameHash)
	if err != nil {
		return err.Result()
	}

	// check permission
	if !registry.Owner.Equals(msg.Owner) {
		return sdk.ErrUnauthorized(fmt.Sprintf("registered owner is %s", registry.Owner)).Result()
	}

	// update owner
	registry.Owner = msg.NewOwner
	keeper.SetRegistry(ctx, nameHash, registry)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeRenew,
			sdk.NewAttribute(types.AttributeKeyName, msg.Name.String()),
			sdk.NewAttribute(types.AttributeKeyOwner, msg.Owner.String()),
			sdk.NewAttribute(types.AttributeKeyNewOwner, msg.NewOwner.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
	})

	return sdk.Result{Events: ctx.EventManager().Events()}
}

// handleMsgRegisterSubName handles a MsgRegisterSubName
func handleMsgRegisterSubName(ctx sdk.Context, keeper Keeper, msg MsgRegisterSubName) sdk.Result {
	nameHash, childNameHash := msg.Name.NameHash()

	// check registry is exists
	registry, err := keeper.GetRegistry(ctx, nameHash)
	if err != nil {
		return err.Result()
	}

	// check permission
	if !registry.Owner.Equals(msg.Owner) {
		return sdk.ErrUnauthorized(fmt.Sprintf("registered owner is %s", registry.Owner)).Result()
	}

	// ensure is there any registered resolve for the address
	if _, err := keeper.GetReverseResolve(ctx, msg.Address); err == nil {
		return ErrAddressAlreadyRegistered(keeper.Codespace()).Result()
	}

	// ensure the name is not taken
	if _, err := keeper.GetResolve(ctx, nameHash, childNameHash); err == nil {
		return ErrNameAlreadyTaken(keeper.Codespace()).Result()
	}

	// set address & reverse resolve
	keeper.SetResolve(ctx, nameHash, childNameHash, msg.Address)
	keeper.SetReverseResolve(ctx, msg.Address, nameHash)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeRegister,
			sdk.NewAttribute(types.AttributeKeyName, msg.Name.String()),
			sdk.NewAttribute(types.AttributeKeyOwner, msg.Owner.String()),
			sdk.NewAttribute(types.AttributeKeyAddress, msg.Address.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
	})

	return sdk.Result{Events: ctx.EventManager().Events()}
}

// handleMsgUnregisterSubName handles a MsgUnregisterSubName
func handleMsgUnregisterSubName(ctx sdk.Context, keeper Keeper, msg MsgUnregisterSubName) sdk.Result {
	nameHash, childNameHash := msg.Name.NameHash()

	// check registry is exists
	registry, err := keeper.GetRegistry(ctx, nameHash)
	if err != nil {
		return err.Result()
	}

	// check permission
	if !registry.Owner.Equals(msg.Owner) {
		return sdk.ErrUnauthorized(fmt.Sprintf("registered owner is %s", registry.Owner)).Result()
	}

	// ensure the name is exists
	resolvedAddr, err := keeper.GetResolve(ctx, nameHash, childNameHash)
	if err != nil {
		return err.Result()
	}

	// delete resolve and reverse resolve
	keeper.DeleteResolve(ctx, nameHash, childNameHash)
	keeper.DeleteReverseResolve(ctx, resolvedAddr)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeUnregister,
			sdk.NewAttribute(types.AttributeKeyName, msg.Name.String()),
			sdk.NewAttribute(types.AttributeKeyOwner, msg.Owner.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
	})

	return sdk.Result{Events: ctx.EventManager().Events()}
}
