package keeper

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerror "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/terra-money/core/x/wasm/types"
)

type Recurse struct {
	Depth    uint32         `json:"depth"`
	Work     uint32         `json:"work"`
	Contract sdk.AccAddress `json:"contract"`
}

type recurseWrapper struct {
	Recurse Recurse `json:"recurse"`
}

func buildQuery(t *testing.T, msg Recurse) []byte {
	wrapper := recurseWrapper{Recurse: msg}
	bz, err := json.Marshal(wrapper)
	require.NoError(t, err)
	return bz
}

type recurseResponse struct {
	Hashed []byte `json:"hashed"`
}

// wasmQuerierWithCounter - wasm query interface for wasm contract
type wasmQuerierWithCounter struct {
	real WasmQuerier
}

// newWasmQuerier returns wasm querier
func newWasmQuerierWithCounter(keeper Keeper) wasmQuerierWithCounter {
	return wasmQuerierWithCounter{real: NewWasmQuerier(keeper)}
}

// Query increase counter and execute real querier
func (querier wasmQuerierWithCounter) Query(ctx sdk.Context, request wasmvmtypes.QueryRequest) ([]byte, error) {
	totalWasmQueryCounter++
	return querier.real.Query(ctx, request)
}

// QueryCustom implements custom query interface
func (wasmQuerierWithCounter) QueryCustom(sdk.Context, json.RawMessage) ([]byte, error) {
	return nil, nil
}

// number os wasm queries called from a contract
var totalWasmQueryCounter int

func initRecurseContract(t *testing.T) (contract sdk.AccAddress, creator sdk.AccAddress, ctx sdk.Context, keeper Keeper, cdc *codec.LegacyAmino) {
	input := CreateTestInput(t)
	ctx, cdc, accKeeper, bankKeeper, keeper := input.Ctx, input.Cdc, input.AccKeeper, input.BankKeeper, input.WasmKeeper
	keeper.RegisterQueriers(map[string]types.WasmQuerierInterface{
		types.WasmQueryRouteWasm: newWasmQuerierWithCounter(keeper),
	}, nil)

	deposit := sdk.NewCoins(sdk.NewInt64Coin("denom", 100000))
	creator = createFakeFundedAccount(ctx, accKeeper, bankKeeper, deposit.Add(deposit...))

	// store the code
	wasmCode, err := ioutil.ReadFile("./testdata/hackatom.wasm")
	require.NoError(t, err)
	codeID, err := keeper.StoreCode(ctx, creator, wasmCode)
	require.NoError(t, err)

	// instantiate the contract
	_, _, bob := keyPubAddr()
	_, _, fred := keyPubAddr()
	initMsg := HackatomExampleInitMsg{
		Verifier:    fred,
		Beneficiary: bob,
	}
	initMsgBz, err := json.Marshal(initMsg)
	require.NoError(t, err)
	contractAddr, _, err := keeper.InstantiateContract(ctx, codeID, creator, sdk.AccAddress{}, initMsgBz, deposit)
	require.NoError(t, err)

	return contractAddr, creator, ctx, keeper, cdc
}

func TestGasCostOnQuery(t *testing.T) {
	GasNoWork := types.InstantiateContractCosts(0) + 3_509
	// Note: about 100 SDK gas (10k wasmVM gas) for each round of sha256
	GasWork50 := GasNoWork + 5_662 // this is a little shy of 50k gas - to keep an eye on the limit

	const (
		GasReturnUnhashed uint64 = 198
		GasReturnHashed   uint64 = 171
	)

	cases := map[string]struct {
		gasLimit    uint64
		msg         Recurse
		expectedGas uint64
	}{
		"no recursion, no work": {
			gasLimit:    400_000,
			msg:         Recurse{},
			expectedGas: GasNoWork,
		},
		"no recursion, some work": {
			gasLimit: 400_000,
			msg: Recurse{
				Work: 50, // 50 rounds of sha256 inside the contract
			},
			expectedGas: GasWork50,
		},
		"recursion 1, no work": {
			gasLimit: 400_000,
			msg: Recurse{
				Depth: 1,
			},
			expectedGas: 2*GasNoWork + GasReturnUnhashed,
		},
		"recursion 1, some work": {
			gasLimit: 400_000,
			msg: Recurse{
				Depth: 1,
				Work:  50,
			},
			expectedGas: 2*GasWork50 + GasReturnHashed + 1,
		},
		"recursion 4, some work": {
			gasLimit: 400_000,
			msg: Recurse{
				Depth: 4,
				Work:  50,
			},
			// NOTE: +6 for rounding issues
			expectedGas: 5*GasWork50 + 4*GasReturnHashed + 1,
		},
	}

	contractAddr, creator, ctx, keeper, _ := initRecurseContract(t)

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			// external limit has no effect (we get a panic if this is enforced)
			keeper.wasmConfig.ContractQueryGasLimit = 1000

			// make sure we set a limit before calling
			ctx = ctx.WithGasMeter(sdk.NewGasMeter(tc.gasLimit))
			require.Equal(t, uint64(0), ctx.GasMeter().GasConsumed())

			// do the query
			recurse := tc.msg
			recurse.Contract = contractAddr
			msg := buildQuery(t, recurse)
			data, err := keeper.queryToContract(ctx, contractAddr, msg)
			require.NoError(t, err)

			// check the gas is what we expected
			assert.Equal(t, tc.expectedGas, ctx.GasMeter().GasConsumed())

			// assert result is 32 byte sha256 hash (if hashed), or contractAddr if not
			var resp recurseResponse
			err = json.Unmarshal(data, &resp)
			require.NoError(t, err)
			if recurse.Work == 0 {
				assert.Equal(t, len(resp.Hashed), len(creator.String()))
			} else {
				assert.Equal(t, len(resp.Hashed), 32)
			}
		})
	}
}

func TestGasOnExternalQuery(t *testing.T) {
	GasNoWork := types.InstantiateContractCosts(0) + 3_509
	// Note: about 100 SDK gas (10k wasmVM gas) for each round of sha256
	GasWork50 := GasNoWork + 5_662 // this is a little shy of 50k gas - to keep an eye on the limit

	cases := map[string]struct {
		gasLimit       uint64
		msg            Recurse
		expectOutOfGas bool
	}{
		"no recursion, plenty gas": {
			gasLimit: 400_000,
			msg: Recurse{
				Work: 50, // 50 rounds of sha256 inside the contract
			},
		},
		"recursion 4, plenty gas": {
			// this uses 244708 gas
			gasLimit: 400_000,
			msg: Recurse{
				Depth: 4,
				Work:  50,
			},
		},
		"no recursion, external gas limit": {
			gasLimit: 5000, // this is not enough
			msg: Recurse{
				Work: 50,
			},
			expectOutOfGas: true,
		},
		"recursion 4, external gas limit": {
			// this uses 244708 gas but give less
			gasLimit: 4 * GasWork50,
			msg: Recurse{
				Depth: 4,
				Work:  50,
			},
			expectOutOfGas: true,
		},
	}

	contractAddr, _, ctx, keeper, cdc := initRecurseContract(t)

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			// set the external gas limit (normally from config file)
			keeper.wasmConfig.ContractQueryGasLimit = tc.gasLimit
			querier := NewLegacyQuerier(keeper, cdc)

			recurse := tc.msg
			recurse.Contract = contractAddr
			msg := buildQuery(t, recurse)

			// do the query
			bz, err := cdc.MarshalJSON(types.NewQueryContractParams(contractAddr, msg))
			require.NoError(t, err)

			if tc.expectOutOfGas {
				_, err = querier(ctx, []string{types.QueryContractStore}, abci.RequestQuery{Data: []byte(bz)})
				require.Error(t, err)
				require.Contains(t, err.Error(), sdkerror.ErrOutOfGas.Error())
			} else {
				// otherwise, make sure we get a good success
				_, err = querier(ctx, []string{types.QueryContractStore}, abci.RequestQuery{Data: []byte(bz)})
				require.NoError(t, err)
			}
		})
	}
}

func TestLimitRecursiveQueryGas(t *testing.T) {
	// The point of this test from https://github.com/CosmWasm/cosmwasm/issues/456
	// Basically, if I burn 90% of gas in CPU loop, then query out (to my self)
	// the sub-query will have all the original gas (minus the 40k instance charge)
	// and can burn 90% and call a sub-contract again...
	// This attack would allow us to use far more than the provided gas before
	// eventually hitting an OutOfGas panic.

	GasNoWork := types.InstantiateContractCosts(0) + 3_509
	GasWork2k := GasNoWork + 228_931

	// This is overhead for calling into a sub-contract
	const GasReturnHashed uint64 = 176

	cases := map[string]struct {
		gasLimit                  uint64
		msg                       Recurse
		expectQueriesFromContract int
		expectedGas               uint64
		expectOutOfGas            bool
	}{
		"no recursion, lots of work": {
			gasLimit: 4_000_000,
			msg: Recurse{
				Depth: 0,
				Work:  2000,
			},
			expectQueriesFromContract: 0,
			expectedGas:               GasWork2k,
		},
		"recursion 5, lots of work": {
			gasLimit: 4_000_000,
			msg: Recurse{
				Depth: 5,
				Work:  2000,
			},
			expectQueriesFromContract: 5,
			// NOTE: +2 for rounding issues
			expectedGas: GasWork2k + 5*(GasWork2k+GasReturnHashed) + 1,
		},
		// this is where we expect an error...
		// it has enough gas to run 4 times and die on the 5th (4th time dispatching to sub-contract)
		// however, if we don't charge the cpu gas before sub-dispatching, we can recurse over 20 times
		// TODO: figure out how to asset how deep it went
		"deep recursion, should die on 5th level": {
			gasLimit: 1_200_000,
			msg: Recurse{
				Depth: 50,
				Work:  2000,
			},
			expectQueriesFromContract: 4,
			expectOutOfGas:            true,
		},
	}

	contractAddr, _, ctx, keeper, _ := initRecurseContract(t)

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			// reset the counter before test
			totalWasmQueryCounter = 0

			// make sure we set a limit before calling
			ctx = ctx.WithGasMeter(sdk.NewGasMeter(tc.gasLimit))
			require.Equal(t, uint64(0), ctx.GasMeter().GasConsumed())

			// prepare the query
			recurse := tc.msg
			recurse.Contract = contractAddr
			msg := buildQuery(t, recurse)

			// if we expect out of gas, make sure this panics
			if tc.expectOutOfGas {
				require.Panics(t, func() {
					_, err := keeper.queryToContract(ctx, contractAddr, msg)
					t.Logf("Got error not panic: %#v", err)
				})
				assert.Equal(t, tc.expectQueriesFromContract, totalWasmQueryCounter)
				return
			}

			// otherwise, we expect a successful call
			_, err := keeper.queryToContract(ctx, contractAddr, msg)
			require.NoError(t, err)
			assert.Equal(t, tc.expectedGas, ctx.GasMeter().GasConsumed())

			assert.Equal(t, tc.expectQueriesFromContract, totalWasmQueryCounter)
		})
	}
}

func TestLimitRecursiveQueryDepth(t *testing.T) {
	contractAddr, _, ctx, keeper, _ := initRecurseContract(t)

	// exceed max query depth
	msg := buildQuery(t, Recurse{
		Depth: types.ContractMaxQueryDepth,
	})

	_, err := keeper.queryToContract(ctx, contractAddr, msg)
	require.Error(t, err)

	msg = buildQuery(t, Recurse{
		Depth: types.ContractMaxQueryDepth - 1, // need to include first query
	})
	_, err = keeper.queryToContract(ctx, contractAddr, msg)
	require.NoError(t, err)
}
