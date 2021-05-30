package keeper

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"

	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	core "github.com/terra-money/core/types"
)

type StakingInitMsg struct {
	Name      string         `json:"name"`
	Symbol    string         `json:"symbol"`
	Decimals  uint8          `json:"decimals"`
	Validator sdk.ValAddress `json:"validator"`
	ExitTax   sdk.Dec        `json:"exit_tax"`
	// MinWithdrawal is uint128 encoded as a string (use sdk.Int?)
	MinWithdrawal string `json:"min_withdrawal"`
}

// StakingHandleMsg is used to encode handle messages
type StakingHandleMsg struct {
	Transfer *transferPayload `json:"transfer,omitempty"`
	Bond     *struct{}        `json:"bond,omitempty"`
	Unbond   *unbondPayload   `json:"unbond,omitempty"`
	Claim    *struct{}        `json:"claim,omitempty"`
	Reinvest *struct{}        `json:"reinvest,omitempty"`
	Change   *ownerPayload    `json:"change_owner,omitempty"`
}

type transferPayload struct {
	Recipient sdk.Address `json:"recipient"`
	// uint128 encoded as string
	Amount string `json:"amount"`
}

type unbondPayload struct {
	// uint128 encoded as string
	Amount string `json:"amount"`
}

// StakingQueryMsg is used to encode query messages
type StakingQueryMsg struct {
	Balance    *addressQuery `json:"balance,omitempty"`
	Claims     *addressQuery `json:"claims,omitempty"`
	TokenInfo  *struct{}     `json:"token_info,omitempty"`
	Investment *struct{}     `json:"investment,omitempty"`
}

type addressQuery struct {
	Address sdk.AccAddress `json:"address"`
}

type BalanceResponse struct {
	Balance string `json:"balance,omitempty"`
}

type ClaimsResponse struct {
	Claims string `json:"claims,omitempty"`
}

type TokenInfoResponse struct {
	Name     string `json:"name"`
	Symbol   string `json:"symbol"`
	Decimals uint8  `json:"decimals"`
}

type InvestmentResponse struct {
	TokenSupply  string         `json:"token_supply"`
	StakedTokens sdk.Coin       `json:"staked_tokens"`
	NominalValue sdk.Dec        `json:"nominal_value"`
	Owner        sdk.AccAddress `json:"owner"`
	Validator    sdk.ValAddress `json:"validator"`
	ExitTax      sdk.Dec        `json:"exit_tax"`
	// MinWithdrawal is uint128 encoded as a string (use sdk.Int?)
	MinWithdrawal string `json:"min_withdrawal"`
}

func TestInitializeStaking(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "wasm")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)
	viper.Set(flags.FlagHome, tempDir)

	input := CreateTestInput(t)
	ctx, accKeeper, stakingKeeper, keeper := input.Ctx, input.AccKeeper, input.StakingKeeper, input.WasmKeeper

	valAddr := addValidator(ctx, stakingKeeper, accKeeper, sdk.NewInt64Coin(core.MicroLunaDenom, 1234567))
	ctx = nextBlock(ctx, stakingKeeper)
	v, found := stakingKeeper.GetValidator(ctx, valAddr)
	assert.True(t, found)
	assert.Equal(t, v.GetDelegatorShares(), sdk.NewDec(1234567))

	deposit := sdk.NewCoins(sdk.NewInt64Coin(sdk.DefaultBondDenom, 100000), sdk.NewInt64Coin(core.MicroLunaDenom, 500000))
	creatorAddr := createFakeFundedAccount(ctx, accKeeper, deposit)

	// upload staking derivates code
	stakingCode, err := ioutil.ReadFile("./testdata/staking.wasm")
	require.NoError(t, err)
	stakingID, err := keeper.StoreCode(ctx, creatorAddr, stakingCode)
	require.NoError(t, err)
	require.Equal(t, uint64(1), stakingID)

	// register to a valid address
	initMsg := StakingInitMsg{
		Name:          "Staking Derivatives",
		Symbol:        "DRV",
		Decimals:      0,
		Validator:     valAddr,
		ExitTax:       sdk.MustNewDecFromStr("0.10"),
		MinWithdrawal: "100",
	}
	initBz, err := json.Marshal(&initMsg)
	require.NoError(t, err)

	stakingAddr, err := keeper.InstantiateContract(ctx, stakingID, creatorAddr, initBz, nil, true)
	require.NoError(t, err)
	require.NotEmpty(t, stakingAddr)

	// nothing spent here
	checkAccount(t, ctx, accKeeper, creatorAddr, deposit)

	// try to register with a validator not on the list and it fails
	_, _, bob := keyPubAddr()
	badInitMsg := StakingInitMsg{
		Name:          "Missing Validator",
		Symbol:        "MISS",
		Decimals:      0,
		Validator:     sdk.ValAddress(bob),
		ExitTax:       sdk.MustNewDecFromStr("0.10"),
		MinWithdrawal: "100",
	}
	badBz, err := json.Marshal(&badInitMsg)
	require.NoError(t, err)

	_, err = keeper.InstantiateContract(ctx, stakingID, creatorAddr, badBz, nil, true)
	require.Error(t, err)

	// no changes to bonding shares
	val, _ := stakingKeeper.GetValidator(ctx, valAddr)
	assert.Equal(t, val.GetDelegatorShares(), sdk.NewDec(1234567))
}

// InitInfo nolint
type InitInfo struct {
	valAddr      sdk.ValAddress
	contractAddr sdk.AccAddress
	creatorAddr  sdk.AccAddress
}

func initializeStaking(t *testing.T, input TestInput) InitInfo {
	ctx, accKeeper, stakingKeeper, keeper := input.Ctx, input.AccKeeper, input.StakingKeeper, input.WasmKeeper

	valAddr := addValidator(ctx, stakingKeeper, accKeeper, sdk.NewInt64Coin(core.MicroLunaDenom, 1000000))
	ctx = nextBlock(ctx, stakingKeeper)

	v, found := stakingKeeper.GetValidator(ctx, valAddr)
	assert.True(t, found)
	assert.Equal(t, v.GetDelegatorShares(), sdk.NewDec(1000000))
	assert.Equal(t, v.Status, sdk.Bonded)

	deposit := sdk.NewCoins(sdk.NewInt64Coin(sdk.DefaultBondDenom, 100000), sdk.NewInt64Coin(core.MicroLunaDenom, 500000))
	creatorAddr := createFakeFundedAccount(ctx, accKeeper, deposit)

	// upload staking derivates code
	stakingCode, err := ioutil.ReadFile("./testdata/staking.wasm")
	require.NoError(t, err)
	stakingID, err := keeper.StoreCode(ctx, creatorAddr, stakingCode)
	require.NoError(t, err)
	require.Equal(t, uint64(1), stakingID)

	// register to a valid address
	initMsg := StakingInitMsg{
		Name:          "Staking Derivatives",
		Symbol:        "DRV",
		Decimals:      0,
		Validator:     valAddr,
		ExitTax:       sdk.MustNewDecFromStr("0.10"),
		MinWithdrawal: "100",
	}
	initBz, err := json.Marshal(&initMsg)
	require.NoError(t, err)

	stakingAddr, err := keeper.InstantiateContract(ctx, stakingID, creatorAddr, initBz, nil, true)
	require.NoError(t, err)
	require.NotEmpty(t, stakingAddr)

	return InitInfo{valAddr, stakingAddr, creatorAddr}
}

func TestBonding(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "wasm")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)
	viper.Set(flags.FlagHome, tempDir)

	input := CreateTestInput(t)
	initInfo := initializeStaking(t, input)

	ctx, accKeeper, stakingKeeper, keeper := input.Ctx, input.AccKeeper, input.StakingKeeper, input.WasmKeeper
	valAddr, contractAddr := initInfo.valAddr, initInfo.contractAddr

	// initial checks of bonding state
	val, found := stakingKeeper.GetValidator(ctx, valAddr)
	require.True(t, found)
	initPower := val.GetDelegatorShares()

	// bob has 160k, putting 80k into the contract
	full := sdk.NewCoins(sdk.NewInt64Coin(core.MicroLunaDenom, 160000))
	funds := sdk.NewCoins(sdk.NewInt64Coin(core.MicroLunaDenom, 80000))
	bob := createFakeFundedAccount(ctx, accKeeper, full)

	// check contract state before
	assertBalance(t, ctx, keeper, contractAddr, bob, "0")
	assertClaims(t, ctx, keeper, contractAddr, bob, "0")
	assertSupply(t, ctx, keeper, contractAddr, "0", sdk.NewInt64Coin(core.MicroLunaDenom, 0))

	bond := StakingHandleMsg{
		Bond: &struct{}{},
	}
	bondBz, err := json.Marshal(bond)
	require.NoError(t, err)
	_, err = keeper.ExecuteContract(ctx, contractAddr, bob, bondBz, funds)
	require.NoError(t, err)

	// check some account values - the money is on neither account (cuz it is bonded)
	checkAccount(t, ctx, accKeeper, contractAddr, sdk.Coins{})
	checkAccount(t, ctx, accKeeper, bob, funds)

	// make sure the proper number of tokens have been bonded
	val, _ = stakingKeeper.GetValidator(ctx, valAddr)
	finalPower := val.GetDelegatorShares()
	assert.Equal(t, sdk.NewInt(80000), finalPower.Sub(initPower).TruncateInt())

	// check the delegation itself
	d, found := stakingKeeper.GetDelegation(ctx, contractAddr, valAddr)
	require.True(t, found)
	assert.Equal(t, d.Shares, sdk.MustNewDecFromStr("80000"))

	// check we have the desired balance
	assertBalance(t, ctx, keeper, contractAddr, bob, "80000")
	assertClaims(t, ctx, keeper, contractAddr, bob, "0")
	assertSupply(t, ctx, keeper, contractAddr, "80000", sdk.NewInt64Coin(core.MicroLunaDenom, 80000))
}

func TestUnbonding(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "wasm")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)
	viper.Set(flags.FlagHome, tempDir)

	input := CreateTestInput(t)
	initInfo := initializeStaking(t, input)

	ctx, accKeeper, stakingKeeper, keeper := input.Ctx, input.AccKeeper, input.StakingKeeper, input.WasmKeeper
	valAddr, contractAddr, creatorAddr := initInfo.valAddr, initInfo.contractAddr, initInfo.creatorAddr

	// initial checks of bonding state
	val, found := stakingKeeper.GetValidator(ctx, valAddr)
	require.True(t, found)
	initPower := val.GetDelegatorShares()

	// bob has 160k, putting 80k into the contract
	full := sdk.NewCoins(sdk.NewInt64Coin(core.MicroLunaDenom, 160000))
	funds := sdk.NewCoins(sdk.NewInt64Coin(core.MicroLunaDenom, 80000))
	bob := createFakeFundedAccount(ctx, accKeeper, full)

	bond := StakingHandleMsg{
		Bond: &struct{}{},
	}
	bondBz, err := json.Marshal(bond)
	require.NoError(t, err)
	_, err = keeper.ExecuteContract(ctx, contractAddr, bob, bondBz, funds)
	require.NoError(t, err)

	// update height a bit
	ctx = nextBlock(ctx, stakingKeeper)

	// now unbond 30k - note that 3k (10%) goes to the owner as a tax, 27k unbonded and available as claims
	unbond := StakingHandleMsg{
		Unbond: &unbondPayload{
			Amount: "30000",
		},
	}
	unbondBz, err := json.Marshal(unbond)
	require.NoError(t, err)
	_, err = keeper.ExecuteContract(ctx, contractAddr, bob, unbondBz, nil)
	require.NoError(t, err)

	// check some account values - the money is on neither account (cuz it is bonded)
	// Note: why is this immediate? just test setup?
	checkAccount(t, ctx, accKeeper, contractAddr, sdk.Coins{})
	checkAccount(t, ctx, accKeeper, bob, funds)

	// make sure the proper number of tokens have been bonded (80k - 27k = 53k)
	val, _ = stakingKeeper.GetValidator(ctx, valAddr)
	finalPower := val.GetDelegatorShares()
	assert.Equal(t, sdk.NewInt(53000), finalPower.Sub(initPower).TruncateInt(), finalPower.String())

	// check the delegation itself
	d, found := stakingKeeper.GetDelegation(ctx, contractAddr, valAddr)
	require.True(t, found)
	assert.Equal(t, d.Shares, sdk.MustNewDecFromStr("53000"))

	// check there is unbonding in progress
	un, found := stakingKeeper.GetUnbondingDelegation(ctx, contractAddr, valAddr)
	require.True(t, found)
	require.Equal(t, 1, len(un.Entries))
	assert.Equal(t, "27000", un.Entries[0].Balance.String())

	// check we have the desired balance
	assertBalance(t, ctx, keeper, contractAddr, bob, "50000")
	assertBalance(t, ctx, keeper, contractAddr, creatorAddr, "3000")
	assertClaims(t, ctx, keeper, contractAddr, bob, "27000")
	assertSupply(t, ctx, keeper, contractAddr, "53000", sdk.NewInt64Coin(core.MicroLunaDenom, 53000))
}

func TestReinvest(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "wasm")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)
	viper.Set(flags.FlagHome, tempDir)

	input := CreateTestInput(t)
	initInfo := initializeStaking(t, input)

	ctx, accKeeper, stakingKeeper, keeper := input.Ctx, input.AccKeeper, input.StakingKeeper, input.WasmKeeper
	valAddr, contractAddr, creatorAddr := initInfo.valAddr, initInfo.contractAddr, initInfo.creatorAddr

	// initial checks of bonding state
	val, found := stakingKeeper.GetValidator(ctx, valAddr)
	require.True(t, found)
	initPower := val.GetDelegatorShares()
	assert.Equal(t, val.Tokens, sdk.NewInt(1000000), "%s", val.Tokens)

	// full is 2x funds, 1x goes to the contract, other stays on his wallet
	full := sdk.NewCoins(sdk.NewInt64Coin(core.MicroLunaDenom, 400000))
	funds := sdk.NewCoins(sdk.NewInt64Coin(core.MicroLunaDenom, 200000))
	bob := createFakeFundedAccount(ctx, accKeeper, full)

	// we will stake 200k to a validator with 1M self-bond
	// this means we should get 1/6 of the rewards
	bond := StakingHandleMsg{
		Bond: &struct{}{},
	}
	bondBz, err := json.Marshal(bond)
	require.NoError(t, err)
	_, err = keeper.ExecuteContract(ctx, contractAddr, bob, bondBz, funds)
	require.NoError(t, err)

	// update height a bit to solidify the delegation
	ctx = nextBlock(ctx, stakingKeeper)
	// we get 1/6, our share should be 40k minus 10% commission = 36k
	setValidatorRewards(ctx, stakingKeeper, input.DistrKeeper, valAddr, "240000")

	// this should withdraw our outstanding 40k of rewards and reinvest them in the same delegation
	reinvest := StakingHandleMsg{
		Reinvest: &struct{}{},
	}
	reinvestBz, err := json.Marshal(reinvest)
	require.NoError(t, err)
	_, err = keeper.ExecuteContract(ctx, contractAddr, bob, reinvestBz, nil)
	require.NoError(t, err)

	// check some account values - the money is on neither account (cuz it is bonded)
	// Note: why is this immediate? just test setup?
	checkAccount(t, ctx, accKeeper, contractAddr, sdk.Coins{})
	checkAccount(t, ctx, accKeeper, bob, funds)

	// check the delegation itself
	d, found := stakingKeeper.GetDelegation(ctx, contractAddr, valAddr)
	require.True(t, found)
	// we started with 200k and added 36k
	assert.Equal(t, d.Shares, sdk.MustNewDecFromStr("236000"))

	// make sure the proper number of tokens have been bonded (80k + 40k = 120k)
	val, _ = stakingKeeper.GetValidator(ctx, valAddr)
	finalPower := val.GetDelegatorShares()
	assert.Equal(t, sdk.NewInt(236000), finalPower.Sub(initPower).TruncateInt(), finalPower.String())

	// check there is no unbonding in progress
	un, found := stakingKeeper.GetUnbondingDelegation(ctx, contractAddr, valAddr)
	assert.False(t, found, "%#v", un)

	// check we have the desired balance
	assertBalance(t, ctx, keeper, contractAddr, bob, "200000")
	assertBalance(t, ctx, keeper, contractAddr, creatorAddr, "0")
	assertClaims(t, ctx, keeper, contractAddr, bob, "0")
	assertSupply(t, ctx, keeper, contractAddr, "200000", sdk.NewInt64Coin(core.MicroLunaDenom, 236000))
}

// adds a few validators and returns a list of validators that are registered
func addValidator(ctx sdk.Context, stakingKeeper staking.Keeper, accountKeeper auth.AccountKeeper, value sdk.Coin) sdk.ValAddress {
	_, pub, accAddr := keyPubAddr()

	addr := sdk.ValAddress(accAddr)

	owner := createFakeFundedAccount(ctx, accountKeeper, sdk.Coins{value})

	msg := staking.MsgCreateValidator{
		Description: types.Description{
			Moniker: "Validator power",
		},
		Commission: types.CommissionRates{
			Rate:          sdk.MustNewDecFromStr("0.1"),
			MaxRate:       sdk.MustNewDecFromStr("0.2"),
			MaxChangeRate: sdk.MustNewDecFromStr("0.01"),
		},
		MinSelfDelegation: sdk.OneInt(),
		DelegatorAddress:  owner,
		ValidatorAddress:  addr,
		PubKey:            pub,
		Value:             value,
	}

	h := staking.NewHandler(stakingKeeper)
	_, err := h(ctx, msg)
	if err != nil {
		panic(err)
	}
	return addr
}

// this will commit the current set, update the block height and set historic info
// basically, letting two blocks pass
func nextBlock(ctx sdk.Context, stakingKeeper staking.Keeper) sdk.Context {
	staking.EndBlocker(ctx, stakingKeeper)
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)
	staking.BeginBlocker(ctx, stakingKeeper)
	return ctx
}

func setValidatorRewards(ctx sdk.Context, stakingKeeper staking.Keeper, distKeeper distribution.Keeper, valAddr sdk.ValAddress, reward string) {
	// allocate some rewards
	vali := stakingKeeper.Validator(ctx, valAddr)
	amount, err := sdk.NewDecFromStr(reward)
	if err != nil {
		panic(err)
	}
	payout := sdk.DecCoins{{Denom: core.MicroLunaDenom, Amount: amount}}
	distKeeper.AllocateTokensToValidator(ctx, vali, payout)
}

func assertBalance(t *testing.T, ctx sdk.Context, keeper Keeper, contract sdk.AccAddress, addr sdk.AccAddress, expected string) {
	query := StakingQueryMsg{
		Balance: &addressQuery{
			Address: addr,
		},
	}
	queryBz, err := json.Marshal(query)
	require.NoError(t, err)
	res, err := keeper.queryToContract(ctx, contract, queryBz)
	require.NoError(t, err)
	var balance BalanceResponse
	err = json.Unmarshal(res, &balance)
	require.NoError(t, err)
	assert.Equal(t, expected, balance.Balance)
}

func assertClaims(t *testing.T, ctx sdk.Context, keeper Keeper, contract sdk.AccAddress, addr sdk.AccAddress, expected string) {
	query := StakingQueryMsg{
		Claims: &addressQuery{
			Address: addr,
		},
	}
	queryBz, err := json.Marshal(query)
	require.NoError(t, err)
	res, err := keeper.queryToContract(ctx, contract, queryBz)
	require.NoError(t, err)
	var claims ClaimsResponse
	err = json.Unmarshal(res, &claims)
	require.NoError(t, err)
	assert.Equal(t, expected, claims.Claims)
}

func assertSupply(t *testing.T, ctx sdk.Context, keeper Keeper, contract sdk.AccAddress, expectedIssued string, expectedBonded sdk.Coin) {
	query := StakingQueryMsg{Investment: &struct{}{}}
	queryBz, err := json.Marshal(query)
	require.NoError(t, err)
	res, err := keeper.queryToContract(ctx, contract, queryBz)
	require.NoError(t, err)
	var invest InvestmentResponse
	err = json.Unmarshal(res, &invest)
	require.NoError(t, err)
	assert.Equal(t, expectedIssued, invest.TokenSupply)
	assert.Equal(t, expectedBonded, invest.StakedTokens)
}
