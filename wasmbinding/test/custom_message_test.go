package wasmbinding_test

import (
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"

	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
	core "github.com/classic-terra/core/v2/types"
	"github.com/classic-terra/core/v2/wasmbinding/bindings"
	markettypes "github.com/classic-terra/core/v2/x/market/types"
)

// go test -v -run ^TestSwap$ github.com/classic-terra/core/v2/wasmbinding/test
// oracle rate: 1 uluna = 1.7 usdr
// 1000 uluna from trader goes to contract
// 1666 usdr (after 2% tax) is swapped into which goes back to contract
func (s *WasmTestSuite) Swap(contractPath string, executeFunc func(contract sdk.AccAddress, sender sdk.AccAddress, msg bindings.TerraMsg, funds sdk.Coin) error) {
	s.SetupTest()
	actor := s.RandomAccountAddress()

	// fund
	s.FundAcc(actor, sdk.NewCoins(sdk.NewInt64Coin(core.MicroLunaDenom, 1000000000)))

	// instantiate reflect contract
	contractAddr := s.InstantiateContract(actor, contractPath)
	s.Require().NotEmpty(contractAddr)

	// setup swap environment
	// Set Oracle Price
	lunaPriceInSDR := sdk.NewDecWithPrec(17, 1)
	s.App.OracleKeeper.SetLunaExchangeRate(s.Ctx, core.MicroSDRDenom, lunaPriceInSDR)

	actorBeforeSwap := s.App.BankKeeper.GetAllBalances(s.Ctx, actor)
	contractBeforeSwap := s.App.BankKeeper.GetAllBalances(s.Ctx, contractAddr)

	// Calculate expected swapped SDR
	expectedSwappedSDR := sdk.NewDec(1000).Mul(lunaPriceInSDR)
	tax := markettypes.DefaultMinStabilitySpread.Mul(expectedSwappedSDR)
	expectedSwappedSDR = expectedSwappedSDR.Sub(tax)

	// execute custom Msg
	msg := bindings.TerraMsg{
		Swap: &bindings.Swap{
			OfferCoin: sdk.NewCoin(core.MicroLunaDenom, sdk.NewInt(1000)),
			AskDenom:  core.MicroSDRDenom,
		},
	}

	err := executeFunc(contractAddr, actor, msg, sdk.NewCoin(core.MicroLunaDenom, sdk.NewInt(1000)))
	s.Require().NoError(err)

	// check result after swap
	actorAfterSwap := s.App.BankKeeper.GetAllBalances(s.Ctx, actor)
	contractAfterSwap := s.App.BankKeeper.GetAllBalances(s.Ctx, contractAddr)

	s.Require().Equal(actorBeforeSwap.AmountOf(core.MicroLunaDenom).Sub(sdk.NewInt(1000)), actorAfterSwap.AmountOf(core.MicroLunaDenom))
	s.Require().Equal(contractBeforeSwap.AmountOf(core.MicroSDRDenom).Add(expectedSwappedSDR.TruncateInt()), contractAfterSwap.AmountOf(core.MicroSDRDenom))
}

// go test -v -run ^TestSwapSend$ github.com/classic-terra/core/v2/wasmbinding/test
// oracle rate: 1 uluna = 1.7 usdr
// 1000 uluna from trader goes to contract
// 1666 usdr (after 2% tax) is swapped into which goes back to contract
// 1666 usdr is sent to trader
func (s *WasmTestSuite) SwapSend(contractPath string, executeFunc func(contract sdk.AccAddress, sender sdk.AccAddress, msg bindings.TerraMsg, funds sdk.Coin) error) {
	s.SetupTest()
	actor := s.RandomAccountAddress()

	// fund
	s.FundAcc(actor, sdk.NewCoins(sdk.NewInt64Coin(core.MicroLunaDenom, 1000000000)))

	// instantiate reflect contract
	contractAddr := s.InstantiateContract(actor, contractPath)
	s.Require().NotEmpty(contractAddr)

	// setup swap environment
	// Set Oracle Price
	lunaPriceInSDR := sdk.NewDecWithPrec(17, 1)
	s.App.OracleKeeper.SetLunaExchangeRate(s.Ctx, core.MicroSDRDenom, lunaPriceInSDR)

	actorBeforeSwap := s.App.BankKeeper.GetAllBalances(s.Ctx, actor)

	// Calculate expected swapped SDR
	expectedSwappedSDR := sdk.NewDec(1000).Mul(lunaPriceInSDR)
	tax := markettypes.DefaultMinStabilitySpread.Mul(expectedSwappedSDR)
	expectedSwappedSDR = expectedSwappedSDR.Sub(tax)

	// execute custom Msg
	msg := bindings.TerraMsg{
		SwapSend: &bindings.SwapSend{
			ToAddress: actor.String(),
			OfferCoin: sdk.NewCoin(core.MicroLunaDenom, sdk.NewInt(1000)),
			AskDenom:  core.MicroSDRDenom,
		},
	}

	err := executeFunc(contractAddr, actor, msg, sdk.NewCoin(core.MicroLunaDenom, sdk.NewInt(1000)))
	s.Require().NoError(err)

	// check result after swap
	actorAfterSwap := s.App.BankKeeper.GetAllBalances(s.Ctx, actor)
	expectedActorAfterSwap := actorBeforeSwap.Sub(sdk.NewCoins(sdk.NewInt64Coin(core.MicroLunaDenom, 1000)))
	expectedActorAfterSwap = expectedActorAfterSwap.Add(sdk.NewCoin(core.MicroSDRDenom, expectedSwappedSDR.TruncateInt()))

	s.Require().Equal(expectedActorAfterSwap, actorAfterSwap)
}

type ReflectExec struct {
	ReflectMsg    *ReflectMsgs    `json:"reflect_msg,omitempty"`
	ReflectSubMsg *ReflectSubMsgs `json:"reflect_sub_msg,omitempty"`
}

type ReflectMsgs struct {
	Msgs []wasmvmtypes.CosmosMsg `json:"msgs"`
}

type ReflectSubMsgs struct {
	Msgs []wasmvmtypes.SubMsg `json:"msgs"`
}

func (s *WasmTestSuite) executeCustom(contract sdk.AccAddress, sender sdk.AccAddress, msg bindings.TerraMsg, funds sdk.Coin) error {
	customBz, err := json.Marshal(msg)
	s.Require().NoError(err)
	reflectMsg := ReflectExec{
		ReflectMsg: &ReflectMsgs{
			Msgs: []wasmvmtypes.CosmosMsg{{
				Custom: customBz,
			}},
		},
	}
	reflectBz, err := json.Marshal(reflectMsg)
	s.Require().NoError(err)

	// no funds sent if amount is 0
	var coins sdk.Coins
	if !funds.Amount.IsNil() {
		coins = sdk.Coins{funds}
	}

	contractKeeper := wasmkeeper.NewDefaultPermissionKeeper(s.App.WasmKeeper)
	_, err = contractKeeper.Execute(s.Ctx, contract, sender, reflectBz, coins)
	return err
}

type customSwap struct {
	Swap *bindings.Swap `json:"swap"`
}

type customSwapSend struct {
	SwapSend *bindings.SwapSend `json:"swap_send"`
}

func (s *WasmTestSuite) executeOldBindings(contract sdk.AccAddress, sender sdk.AccAddress, msg bindings.TerraMsg, funds sdk.Coin) error {
	var reflectBz []byte
	switch {
	case msg.Swap != nil:
		customSwap := customSwap{
			Swap: msg.Swap,
		}
		var err error
		reflectBz, err = json.Marshal(customSwap)
		s.Require().NoError(err)
	case msg.SwapSend != nil:
		customSwapSend := customSwapSend{
			SwapSend: msg.SwapSend,
		}
		var err error
		reflectBz, err = json.Marshal(customSwapSend)
		s.Require().NoError(err)
	}

	// no funds sent if amount is 0
	var coins sdk.Coins
	if !funds.Amount.IsNil() {
		coins = sdk.Coins{funds}
	}

	contractKeeper := wasmkeeper.NewDefaultPermissionKeeper(s.App.WasmKeeper)
	_, err := contractKeeper.Execute(s.Ctx, contract, sender, reflectBz, coins)
	return err
}
