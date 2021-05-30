package keeper

import (
	"bytes"
	"compress/gzip"
	"errors"
	"io"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/terra-money/core/x/wasm/types"
)

func TestUncompress(t *testing.T) {
	input := CreateTestInput(t)
	ctx, keeper := input.Ctx, input.WasmKeeper

	wasmRaw, err := ioutil.ReadFile("./testdata/hackatom.wasm")
	require.NoError(t, err)

	wasmGzipped, err := ioutil.ReadFile("./testdata/hackatom.wasm.gzip")
	require.NoError(t, err)

	specs := map[string]struct {
		src       []byte
		expError  error
		expResult []byte
	}{
		"handle wasm uncompressed": {
			src:       wasmRaw,
			expResult: wasmRaw,
		},
		"handle wasm compressed": {
			src:       wasmGzipped,
			expResult: wasmRaw,
		},
		"handle nil slice": {
			src:       nil,
			expResult: nil,
		},
		"handle short unidentified": {
			src:       []byte{0x1, 0x2},
			expResult: []byte{0x1, 0x2},
		},
		"handle big input slice": {
			src:      []byte(strings.Repeat("a", int(keeper.MaxContractSize(ctx)+1))),
			expError: types.ErrExceedMaxContractSize,
		},
		"handle gzip identifier only": {
			src:      gzipIdent,
			expError: io.ErrUnexpectedEOF,
		},
		"handle broken gzip": {
			src:      append(gzipIdent, byte(0x1)),
			expError: io.ErrUnexpectedEOF,
		},
		"handle incomplete gzip": {
			src:      wasmGzipped[:len(wasmGzipped)-5],
			expError: io.ErrUnexpectedEOF,
		},
		"handle big gzip output": {
			src:      asGzip(strings.Repeat("a", int(keeper.MaxContractSize(ctx)+1))),
			expError: io.ErrUnexpectedEOF,
		},
		"handle other big gzip output": {
			src:      asGzip(strings.Repeat("a", 2*int(keeper.MaxContractSize(ctx)))),
			expError: io.ErrUnexpectedEOF,
		},
	}

	limit := input.WasmKeeper.MaxContractSize(input.Ctx)
	for msg, spec := range specs {
		t.Run(msg, func(t *testing.T) {
			r, err := keeper.uncompress(spec.src, limit)
			require.True(t, errors.Is(spec.expError, err), "exp %+v got %+v", spec.expError, err)
			if spec.expError != nil {
				return
			}
			assert.Equal(t, spec.expResult, r)
		})
	}

}

func asGzip(src string) []byte {
	var buf bytes.Buffer
	if _, err := io.Copy(gzip.NewWriter(&buf), strings.NewReader(src)); err != nil {
		panic(err)
	}
	return buf.Bytes()
}
