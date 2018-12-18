package oracle

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewHandler returns a handler for "bank" type messages.
func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case PriceFeedMsg:
			return k.handlePriceFeedMsg(ctx, msg)
		default:
			errMsg := "Unrecognized bank Msg type: %s" + msg.Type()
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// NewEndBlocker checks proposals and generates a EndBlocker
//func NewEndBlocker(k Keeper) sdk.EndBlocker {
//	return func(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
//		newTags := sdk.NewTags()
//
//		// Update price elects
//		if ctx.BlockHeight()%k.timeout == 0 {
//
//			newTags.AppendTag("action", []byte("price_update"))
//			for _, whiteDenom := range denomWhiteList {
//				elect, err := k.Elect(ctx, whiteDenom)
//				if err == nil {
//					newTags.AppendTag(whiteDenom, []byte(elect.Price.String()))
//				}
//			}
//		}
//
//		return abci.ResponseEndBlock{
//			Tags: newTags,
//		}
//	}
//}

func EndBlocker(ctx sdk.Context, k Keeper) sdk.Tags {
	resTags := sdk.NewTags()

	// Update price elects
	if ctx.BlockHeight()%k.timeout == 0 {

		resTags.AppendTag("action", []byte("price_update"))
		for _, whiteDenom := range denomWhiteList {
			elect, err := k.Elect(ctx, whiteDenom)
			if err == nil {
				resTags.AppendTag(whiteDenom, []byte(elect.Price.String()))
			}
		}
	}

	return resTags
}

// Handle is used by other modules to handle Msg
func (keeper Keeper) handlePriceFeedMsg(ctx sdk.Context, pfm PriceFeedMsg) sdk.Result {
	valset := keeper.valset

	// Check the feeder is a validater
	val := valset.Validator(ctx, sdk.ValAddress(pfm.Feeder.Bytes()))
	if val == nil {
		return ErrNotValidator(DefaultCodespace, pfm.Feeder).Result()
	}

	priceVote := NewPriceVote(pfm, val.GetPower())
	keeper.AddVote(ctx, priceVote)

	return sdk.Result{}
}
