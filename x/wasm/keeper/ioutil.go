package keeper

import (
	"bytes"
	"compress/gzip"
	"io"
	"io/ioutil"

	"github.com/terra-money/core/x/wasm/types"
)

// magic bytes to identify gzip.
// See https://www.ietf.org/rfc/rfc1952.txt
// and https://github.com/golang/go/blob/master/src/net/http/sniff.go#L186
var gzipIdent = []byte("\x1F\x8B\x08")

// uncompress returns gzip uncompressed content or given src when not gzip.
func (k Keeper) uncompress(src []byte, maxContractSize uint64) ([]byte, error) {
	switch n := uint64(len(src)); {
	case n < 3:
		return src, nil
	case n > maxContractSize:
		return nil, types.ErrExceedMaxContractSize
	}

	if !bytes.Equal(gzipIdent, src[0:3]) {
		return src, nil
	}

	zr, err := gzip.NewReader(bytes.NewReader(src))
	if err != nil {
		return nil, err
	}
	zr.Multistream(false)
	defer zr.Close()

	return ioutil.ReadAll(LimitReader(zr, int64(maxContractSize)))
}

// LimitReader returns LimitedReader
func LimitReader(r io.Reader, n int64) io.Reader {
	return &LimitedReader{r: &io.LimitedReader{R: r, N: n}}
}

// LimitedReader is a Reader that reads from r
// but stops with types.ErrExceedMaxContractSize after n bytes.
type LimitedReader struct {
	r *io.LimitedReader
}

func (l *LimitedReader) Read(p []byte) (n int, err error) {
	if l.r.N <= 0 {
		return 0, types.ErrExceedMaxContractSize
	}
	return l.r.Read(p)
}
