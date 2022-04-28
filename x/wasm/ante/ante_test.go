package ante_test

import (
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/libs/log"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	dbm "github.com/tendermint/tm-db"

	"github.com/terra-money/core/x/wasm/ante"
	"github.com/terra-money/core/x/wasm/types"
)

func TestCountTxDecorator(t *testing.T) {
	keyWasm := sdk.NewKVStoreKey(types.StoreKey)
	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(keyWasm, sdk.StoreTypeIAVL, db)
	require.NoError(t, ms.LoadLatestVersion())
	const myCurrentBlockHeight = 100

	specs := map[string]struct {
		setupDB        func(t *testing.T, ctx sdk.Context)
		simulate       bool
		nextAssertAnte func(ctx sdk.Context, tx sdk.Tx, simulate bool) (sdk.Context, error)
		expErr         bool
	}{
		"no initial counter set": {
			setupDB: func(t *testing.T, ctx sdk.Context) {},
			nextAssertAnte: func(ctx sdk.Context, tx sdk.Tx, simulate bool) (sdk.Context, error) {
				gotCounter, ok := types.TXCounter(ctx)
				require.True(t, ok)
				assert.Equal(t, uint32(0), gotCounter)
				// and stored +1
				bz := ctx.MultiStore().GetKVStore(keyWasm).Get(types.TXCounterKey)
				assert.Equal(t, []byte{0, 0, 0, 0, 0, 0, 0, myCurrentBlockHeight, 0, 0, 0, 1}, bz)
				return ctx, nil
			},
		},
		"persistent counter incremented - big endian": {
			setupDB: func(t *testing.T, ctx sdk.Context) {
				bz := []byte{0, 0, 0, 0, 0, 0, 0, myCurrentBlockHeight, 1, 0, 0, 2}
				ctx.MultiStore().GetKVStore(keyWasm).Set(types.TXCounterKey, bz)
			},
			nextAssertAnte: func(ctx sdk.Context, tx sdk.Tx, simulate bool) (sdk.Context, error) {
				gotCounter, ok := types.TXCounter(ctx)
				require.True(t, ok)
				assert.Equal(t, uint32(1<<24+2), gotCounter)
				// and stored +1
				bz := ctx.MultiStore().GetKVStore(keyWasm).Get(types.TXCounterKey)
				assert.Equal(t, []byte{0, 0, 0, 0, 0, 0, 0, myCurrentBlockHeight, 1, 0, 0, 3}, bz)
				return ctx, nil
			},
		},
		"old height counter replaced": {
			setupDB: func(t *testing.T, ctx sdk.Context) {
				previousHeight := byte(myCurrentBlockHeight - 1)
				bz := []byte{0, 0, 0, 0, 0, 0, 0, previousHeight, 0, 0, 0, 1}
				ctx.MultiStore().GetKVStore(keyWasm).Set(types.TXCounterKey, bz)
			},
			nextAssertAnte: func(ctx sdk.Context, tx sdk.Tx, simulate bool) (sdk.Context, error) {
				gotCounter, ok := types.TXCounter(ctx)
				require.True(t, ok)
				assert.Equal(t, uint32(0), gotCounter)
				// and stored +1
				bz := ctx.MultiStore().GetKVStore(keyWasm).Get(types.TXCounterKey)
				assert.Equal(t, []byte{0, 0, 0, 0, 0, 0, 0, myCurrentBlockHeight, 0, 0, 0, 1}, bz)
				return ctx, nil
			},
		},
		"simulation not persisted": {
			setupDB: func(t *testing.T, ctx sdk.Context) {
			},
			simulate: true,
			nextAssertAnte: func(ctx sdk.Context, tx sdk.Tx, simulate bool) (sdk.Context, error) {
				_, ok := types.TXCounter(ctx)
				assert.False(t, ok)
				require.True(t, simulate)
				// and not stored
				assert.False(t, ctx.MultiStore().GetKVStore(keyWasm).Has(types.TXCounterKey))
				return ctx, nil
			},
		},
	}
	for name, spec := range specs {
		t.Run(name, func(t *testing.T) {
			ctx := sdk.NewContext(ms.CacheMultiStore(), tmproto.Header{
				Height: myCurrentBlockHeight,
				Time:   time.Date(2021, time.September, 27, 12, 0, 0, 0, time.UTC),
			}, false, log.NewNopLogger())

			spec.setupDB(t, ctx)
			var anyTx sdk.Tx

			// when
			ante := ante.NewCountTXDecorator(keyWasm)
			_, gotErr := ante.AnteHandle(ctx, anyTx, spec.simulate, spec.nextAssertAnte)
			if spec.expErr {
				require.Error(t, gotErr)
				return
			}
			require.NoError(t, gotErr)
		})
	}
}
