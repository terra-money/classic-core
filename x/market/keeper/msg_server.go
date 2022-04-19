package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	"github.com/terra-money/core/x/market/types"
	oracletypes "github.com/terra-money/core/x/oracle/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the market MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

func (k msgServer) Swap(goCtx context.Context, msg *types.MsgSwap) (*types.MsgSwapResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	addr, err := sdk.AccAddressFromBech32(msg.Trader)
	if err != nil {
		return nil, err
	}

	return k.handleSwapRequest(ctx, addr, addr, msg.OfferCoin, msg.AskDenom)
}

func (k msgServer) SwapSend(goCtx context.Context, msg *types.MsgSwapSend) (*types.MsgSwapSendResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	fromAddr, err := sdk.AccAddressFromBech32(msg.FromAddress)
	if err != nil {
		return nil, err
	}

	toAddr, err := sdk.AccAddressFromBech32(msg.ToAddress)
	if err != nil {
		return nil, err
	}

	res, err := k.handleSwapRequest(ctx, fromAddr, toAddr, msg.OfferCoin, msg.AskDenom)
	if err != nil {
		return nil, err
	}

	return &types.MsgSwapSendResponse{
		SwapCoin: res.SwapCoin,
		SwapFee:  res.SwapFee,
	}, nil
}

// handleMsgSwap handles the logic of a MsgSwap
// This function does not repeat checks that have already been performed in msg.ValidateBasic()
// Ex) assert(offerCoin.Denom != askDenom)
func (k msgServer) handleSwapRequest(ctx sdk.Context,
	trader sdk.AccAddress, receiver sdk.AccAddress,
	offerCoin sdk.Coin, askDenom string) (*types.MsgSwapResponse, error) {

	// Compute exchange rates between the ask and offer
	swapDecCoin, spread, err := k.ComputeSwap(ctx, offerCoin, askDenom)
	if err != nil {
		return nil, err
	}

	// Charge a spread if applicable; the spread is burned
	var feeDecCoin sdk.DecCoin
	if spread.IsPositive() {
		feeDecCoin = sdk.NewDecCoinFromDec(swapDecCoin.Denom, spread.Mul(swapDecCoin.Amount))
	} else {
		feeDecCoin = sdk.NewDecCoin(swapDecCoin.Denom, sdk.ZeroInt())
	}

	// Subtract fee from the swap coin
	swapDecCoin.Amount = swapDecCoin.Amount.Sub(feeDecCoin.Amount)

	// Update pool delta
	err = k.ApplySwapToPool(ctx, offerCoin, swapDecCoin)
	if err != nil {
		return nil, err
	}

	// Send offer coins to module account
	offerCoins := sdk.NewCoins(offerCoin)
	err = k.BankKeeper.SendCoinsFromAccountToModule(ctx, trader, types.ModuleName, offerCoins)
	if err != nil {
		return nil, err
	}

	// Burn offered coins and subtract from the trader's account
	err = k.BankKeeper.BurnCoins(ctx, types.ModuleName, offerCoins)
	if err != nil {
		return nil, err
	}

	// Mint asked coins and credit Trader's account
	swapCoin, decimalCoin := swapDecCoin.TruncateDecimal()
	feeDecCoin = feeDecCoin.Add(decimalCoin) // add truncated decimalCoin to swapFee
	feeCoin, _ := feeDecCoin.TruncateDecimal()

	mintCoins := sdk.NewCoins(swapCoin.Add(feeCoin))
	err = k.BankKeeper.MintCoins(ctx, types.ModuleName, mintCoins)
	if err != nil {
		return nil, err
	}

	// Send swap coin to the trader
	swapCoins := sdk.NewCoins(swapCoin)
	err = k.BankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, receiver, swapCoins)
	if err != nil {
		return nil, err
	}

	// Spend swap fee as validator rewards
	if feeCoin.IsPositive() {
		// Split feeCoin 50:50 between oracle voting rewards and block validation rewards
		blockValidationReward := feeCoin.Amount.QuoRaw(2)
		oracleVotesReward := feeCoin.Amount.Sub(blockValidationReward)

		// Send half to fee collector account
		if !blockValidationReward.IsZero() {
			blockValidationRewardCoins := sdk.Coins{{Amount: blockValidationReward, Denom: feeCoin.Denom}}
			err = k.BankKeeper.SendCoinsFromModuleToModule(ctx, types.ModuleName, authtypes.FeeCollectorName, blockValidationRewardCoins)
			if err != nil {
				return nil, err
			}
		}

		// Send left half to oracle account
		if !oracleVotesReward.IsZero() {
			oracleVotesRewardCoins := sdk.Coins{{Amount: oracleVotesReward, Denom: feeCoin.Denom}}
			err = k.BankKeeper.SendCoinsFromModuleToModule(ctx, types.ModuleName, oracletypes.ModuleName, oracleVotesRewardCoins)
			if err != nil {
				return nil, err
			}
		}
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventSwap,
			sdk.NewAttribute(types.AttributeKeyOffer, offerCoin.String()),
			sdk.NewAttribute(types.AttributeKeyTrader, trader.String()),
			sdk.NewAttribute(types.AttributeKeyRecipient, receiver.String()),
			sdk.NewAttribute(types.AttributeKeySwapCoin, swapCoin.String()),
			sdk.NewAttribute(types.AttributeKeySwapFee, feeCoin.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
	})

	return &types.MsgSwapResponse{
		SwapCoin: swapCoin,
		SwapFee:  feeCoin,
	}, nil
}
