package keeper

import (
	"io/ioutil"
	"testing"

	"github.com/classic-terra/core/x/wasm/config"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

// BenchmarkVerification benchmarks secp256k1 verification which is 1000 gas based on cpu time.
//
// Just this function is copied from
// https://github.com/cosmos/cosmos-sdk/blob/90e9370bd80d9a3d41f7203ddb71166865561569/crypto/keys/internal/benchmarking/bench.go#L48-L62
// And thus under the GO license (BSD style)
// go test -benchmem -run=^$ -bench ^BenchmarkGasNormalization$ github.com/classic-terra/core/x/wasm/keeper
func BenchmarkGasNormalization(b *testing.B) {
	priv := secp256k1.GenPrivKey()
	pub := priv.PubKey()

	// use a short message, so this time doesn't get dominated by hashing.
	message := []byte("Hello, world!")
	signature, err := priv.Sign(message)
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pub.VerifySignature(message, signature)
	}
}

// By comparing the timing for queries on pinned vs unpinned, the difference gives us the overhead of
// instantiating an unpinned contract. That value can be used to determine a reasonable gas price
// for the InstantiationCost
// go test -benchmem -run=^$ -bench ^BenchmarkInstantiationOverhead$ github.com/classic-terra/core/x/wasm/keeper
func BenchmarkInstantiationOverhead(b *testing.B) {
	specs := map[string]struct {
		pinned bool
	}{
		"unpinned, memory db": {
			pinned: false,
		},
		"pinned, memory db": {
			pinned: true,
		},
	}
	for name, spec := range specs {
		b.Run(name, func(b *testing.B) {
			// explicitly set ContractMemoryCacheSize to zero to disable InMemoryCache
			wasmConfig := config.DefaultConfig()
			wasmConfig.ContractMemoryCacheSize = 0
			wasmConfig.ContractDebugMode = true
			input := CreateTestInput(b, wasmConfig)

			example := InstantiateHackatomExampleContract(b, input)
			if spec.pinned {
				require.NoError(b, input.WasmKeeper.pinCode(input.Ctx, example.CodeID))
			}
			input.Ctx = input.Ctx.WithGasMeter(sdk.NewGasMeter(100_000_000_000))
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, err := input.WasmKeeper.queryToContract(input.Ctx, example.Contract, []byte(`{"verifier":{}}`))
				require.NoError(b, err)
			}
		})
	}
}

// Calculate the time it takes to compile some wasm code the first time.
// This will help us adjust pricing for StoreCode
// go test -benchmem -run=^$ -bench ^BenchmarkCompilation$ github.com/classic-terra/core/x/wasm/keeper
// run stat -f%z x/wasm/keeper/testdata/hackatom.wasm to get byte size of wasm file
func BenchmarkCompilation(b *testing.B) {
	specs := map[string]struct {
		wasmFile string
	}{
		// 214650 bytes
		"hackatom": {
			wasmFile: "./testdata/hackatom.wasm",
		},
		// 125239
		"burner": {
			wasmFile: "./testdata/burner.wasm",
		},
		// 222488
		"maker": {
			wasmFile: "./testdata/maker.wasm",
		},
	}

	for name, spec := range specs {
		b.Run(name, func(b *testing.B) {
			input := CreateTestInput(b, config.DefaultConfig())

			// print out code size for comparisons
			code, err := ioutil.ReadFile(spec.wasmFile)
			require.NoError(b, err)
			b.Logf("\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b(size: %d)  ", len(code))

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = StoreExampleContract(b, input, spec.wasmFile)
			}
		})
	}
}
