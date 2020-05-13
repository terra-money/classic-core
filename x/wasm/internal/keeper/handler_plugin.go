package keeper

// import (
// 	"encoding/json"

// 	sdk "github.com/cosmos/cosmos-sdk/types"
// 	"github.com/cosmos/cosmos-sdk/x/bank"
// 	"github.com/cosmos/cosmos-sdk/x/distribution"
// 	"github.com/cosmos/cosmos-sdk/x/staking"

// 	wasmTypes "github.com/CosmWasm/go-cosmwasm/types"
// 	"github.com/terra-project/core/x/wasm/internal/types"
// )

// // MessageHandler - reponsible for handling contract output
// type MessageHandler struct {
// 	router   sdk.Router
// 	encoders MessageEncoders
// }

// // NewMessageHandler create msg handler
// func NewMessageHandler(router sdk.Router, customEncoders *MessageEncoders) MessageHandler {
// 	encoders := defaultEncoders()
// 	return MessageHandler{
// 		router:   router,
// 		encoders: encoders,
// 	}
// }

// // BankEncoder - cosmos-sdk bank module msg encoder
// type BankEncoder func(sender sdk.AccAddress, msg *wasmTypes.BankMsg) ([]sdk.Msg, sdk.Error)

// // CustomEncoder - custom msg encoder
// type CustomEncoder func(sender sdk.AccAddress, msg json.RawMessage) ([]sdk.Msg, sdk.Error)

// // StakingEncoder - cosmos-sdk staking module msg encoder
// type StakingEncoder func(sender sdk.AccAddress, msg *wasmTypes.StakingMsg) ([]sdk.Msg, sdk.Error)

// // WasmEncoder - wasm module encoder
// type WasmEncoder func(sender sdk.AccAddress, msg *wasmTypes.WasmMsg) ([]sdk.Msg, sdk.Error)

// // MessageEncoders - cosmwasm msg encoder set
// type MessageEncoders struct {
// 	Bank    BankEncoder
// 	Custom  CustomEncoder
// 	Staking StakingEncoder
// 	Wasm    WasmEncoder
// }

// // defaultEncoders - default msg encoder
// func defaultEncoders() MessageEncoders {
// 	return MessageEncoders{
// 		Bank:    encodeBankMsg,
// 		Custom:  noCustomMsg,
// 		Staking: encodeStakingMsg,
// 		Wasm:    encodeWasmMsg,
// 	}
// }

// func (e MessageEncoders) encode(contractAddr sdk.AccAddress, msg wasmTypes.CosmosMsg) ([]sdk.Msg, sdk.Error) {
// 	switch {
// 	case msg.Bank != nil:
// 		return e.Bank(contractAddr, msg.Bank)
// 	case msg.Custom != nil:
// 		return e.Custom(contractAddr, msg.Custom)
// 	case msg.Staking != nil:
// 		return e.Staking(contractAddr, msg.Staking)
// 	case msg.Wasm != nil:
// 		return e.Wasm(contractAddr, msg.Wasm)
// 	}
// 	return nil, sdk.ErrInternal("single msg cannot contain multiple msgs")
// }

// func encodeBankMsg(sender sdk.AccAddress, msg *wasmTypes.BankMsg) ([]sdk.Msg, sdk.Error) {
// 	if msg.Send == nil {
// 		return nil, types.ErrInvalidMsg("Unknown variant of Bank")
// 	}
// 	if len(msg.Send.Amount) == 0 {
// 		return nil, nil
// 	}
// 	fromAddr, stderr := sdk.AccAddressFromBech32(msg.Send.FromAddress)
// 	if stderr != nil {
// 		return nil, sdk.ErrInvalidAddress(msg.Send.FromAddress)
// 	}
// 	toAddr, stderr := sdk.AccAddressFromBech32(msg.Send.ToAddress)
// 	if stderr != nil {
// 		return nil, sdk.ErrInvalidAddress(msg.Send.ToAddress)
// 	}
// 	toSend, err := convertWasmCoinsToSdkCoins(msg.Send.Amount)
// 	if err != nil {
// 		return nil, err
// 	}
// 	sdkMsg := bank.MsgSend{
// 		FromAddress: fromAddr,
// 		ToAddress:   toAddr,
// 		Amount:      toSend,
// 	}
// 	return []sdk.Msg{sdkMsg}, nil
// }

// func noCustomMsg(sender sdk.AccAddress, msg json.RawMessage) ([]sdk.Msg, sdk.Error) {
// 	return nil, types.ErrInvalidMsg("Custom variant not supported")
// }

// func encodeStakingMsg(sender sdk.AccAddress, msg *wasmTypes.StakingMsg) ([]sdk.Msg, sdk.Error) {
// 	if msg.Delegate != nil {
// 		validator, err := sdk.ValAddressFromBech32(msg.Delegate.Validator)
// 		if err != nil {
// 			return nil, sdk.ErrInvalidAddress(msg.Delegate.Validator)
// 		}
// 		coin, sdkErr := convertWasmCoinToSdkCoin(msg.Delegate.Amount)
// 		if sdkErr != nil {
// 			return nil, sdkErr
// 		}
// 		sdkMsg := staking.MsgDelegate{
// 			DelegatorAddress: sender,
// 			ValidatorAddress: validator,
// 			Amount:           coin,
// 		}
// 		return []sdk.Msg{sdkMsg}, nil
// 	}
// 	if msg.Redelegate != nil {
// 		src, err := sdk.ValAddressFromBech32(msg.Redelegate.SrcValidator)
// 		if err != nil {
// 			return nil, sdk.ErrInvalidAddress(msg.Redelegate.SrcValidator)
// 		}
// 		dst, err := sdk.ValAddressFromBech32(msg.Redelegate.DstValidator)
// 		if err != nil {
// 			return nil, sdk.ErrInvalidAddress(msg.Redelegate.DstValidator)
// 		}
// 		coin, sdkErr := convertWasmCoinToSdkCoin(msg.Redelegate.Amount)
// 		if sdkErr != nil {
// 			return nil, sdkErr
// 		}
// 		sdkMsg := staking.MsgBeginRedelegate{
// 			DelegatorAddress:    sender,
// 			ValidatorSrcAddress: src,
// 			ValidatorDstAddress: dst,
// 			Amount:              coin,
// 		}
// 		return []sdk.Msg{sdkMsg}, nil
// 	}
// 	if msg.Undelegate != nil {
// 		validator, err := sdk.ValAddressFromBech32(msg.Undelegate.Validator)
// 		if err != nil {
// 			return nil, sdk.ErrInvalidAddress(msg.Undelegate.Validator)
// 		}
// 		coin, sdkErr := convertWasmCoinToSdkCoin(msg.Undelegate.Amount)
// 		if sdkErr != nil {
// 			return nil, sdkErr
// 		}
// 		sdkMsg := staking.MsgUndelegate{
// 			DelegatorAddress: sender,
// 			ValidatorAddress: validator,
// 			Amount:           coin,
// 		}
// 		return []sdk.Msg{sdkMsg}, nil
// 	}
// 	if msg.Withdraw != nil {
// 		var err error
// 		rcpt := sender
// 		if len(msg.Withdraw.Recipient) != 0 {
// 			rcpt, err = sdk.AccAddressFromBech32(msg.Withdraw.Recipient)
// 			if err != nil {
// 				return nil, sdk.ErrInvalidAddress(msg.Withdraw.Recipient)
// 			}
// 		}
// 		validator, err := sdk.ValAddressFromBech32(msg.Withdraw.Validator)
// 		if err != nil {
// 			return nil, sdk.ErrInvalidAddress(msg.Withdraw.Validator)
// 		}
// 		setMsg := distribution.MsgSetWithdrawAddress{
// 			DelegatorAddress: sender,
// 			WithdrawAddress:  rcpt,
// 		}
// 		withdrawMsg := distribution.MsgWithdrawDelegatorReward{
// 			DelegatorAddress: sender,
// 			ValidatorAddress: validator,
// 		}
// 		return []sdk.Msg{setMsg, withdrawMsg}, nil
// 	}
// 	return nil, types.ErrInvalidMsg("Unknown variant of Staking")
// }

// func encodeWasmMsg(sender sdk.AccAddress, msg *wasmTypes.WasmMsg) ([]sdk.Msg, sdk.Error) {
// 	if msg.Execute != nil {
// 		contractAddr, err := sdk.AccAddressFromBech32(msg.Execute.ContractAddr)
// 		if err != nil {
// 			return nil, sdk.ErrInvalidAddress(msg.Execute.ContractAddr)
// 		}
// 		coins, sdkErr := convertWasmCoinsToSdkCoins(msg.Execute.Send)
// 		if sdkErr != nil {
// 			return nil, sdkErr
// 		}

// 		sdkMsg := types.MsgExecuteContract{
// 			Sender:   sender,
// 			Contract: contractAddr,
// 			Msg:      msg.Execute.Msg,
// 			Coins:    coins,
// 		}
// 		return []sdk.Msg{sdkMsg}, nil
// 	}
// 	if msg.Instantiate != nil {
// 		coins, err := convertWasmCoinsToSdkCoins(msg.Instantiate.Send)
// 		if err != nil {
// 			return nil, err
// 		}

// 		sdkMsg := types.MsgInstantiateContract{
// 			Sender:    sender,
// 			CodeID:    msg.Instantiate.CodeID,
// 			InitMsg:   msg.Instantiate.Msg,
// 			InitCoins: coins,
// 		}
// 		return []sdk.Msg{sdkMsg}, nil
// 	}
// 	return nil, types.ErrInvalidMsg("Unknown variant of Wasm")
// }

// func (h MessageHandler) dispatch(ctx sdk.Context, contractAddr sdk.AccAddress, msg wasmTypes.CosmosMsg) sdk.Error {
// 	sdkMsgs, err := h.encoders.encode(contractAddr, msg)
// 	if err != nil {
// 		return err
// 	}
// 	for _, sdkMsg := range sdkMsgs {
// 		if err := h.handleSdkMessage(ctx, contractAddr, sdkMsg); err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }

// func (h MessageHandler) handleSdkMessage(ctx sdk.Context, contractAddr sdk.Address, msg sdk.Msg) sdk.Error {
// 	// make sure this account can send it
// 	for _, acct := range msg.GetSigners() {
// 		if !acct.Equals(contractAddr) {
// 			return sdk.ErrUnauthorized("contract doesn't have permission")
// 		}
// 	}

// 	// find the handler and execute it
// 	handler := h.router.Route(msg.Route())
// 	if handler == nil {
// 		return sdk.ErrUnknownRequest(msg.Route())
// 	}
// 	res := handler(ctx, msg)
// 	if !res.IsOK() {
// 		return sdk.NewError(res.Codespace, res.Code, res.Log)
// 	}
// 	// redispatch all events, (type sdk.EventTypeMessage will be filtered out in the handler)
// 	ctx.EventManager().EmitEvents(res.Events)

// 	return nil
// }

// func convertWasmCoinsToSdkCoins(coins []wasmTypes.Coin) (sdk.Coins, sdk.Error) {
// 	var toSend sdk.Coins
// 	for _, coin := range coins {
// 		c, err := convertWasmCoinToSdkCoin(coin)
// 		if err != nil {
// 			return nil, err
// 		}
// 		toSend = append(toSend, c)
// 	}
// 	return toSend, nil
// }

// func convertWasmCoinToSdkCoin(coin wasmTypes.Coin) (sdk.Coin, sdk.Error) {
// 	amount, ok := sdk.NewIntFromString(coin.Amount)
// 	if !ok {
// 		return sdk.Coin{}, sdk.ErrInvalidCoins(coin.Amount + coin.Denom)
// 	}
// 	return sdk.Coin{
// 		Denom:  coin.Denom,
// 		Amount: amount,
// 	}, nil
// }
