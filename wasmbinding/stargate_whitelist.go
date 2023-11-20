package wasmbinding

import (
	"fmt"
	"sync"

	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
	markettypes "github.com/classic-terra/core/v2/x/market/types"
	oracletypes "github.com/classic-terra/core/v2/x/oracle/types"
	treasurytypes "github.com/classic-terra/core/v2/x/treasury/types"
	"github.com/cosmos/cosmos-sdk/codec"
)

// stargateWhitelist keeps whitelist and its deterministic
// response binding for stargate queries.
//
// The query can be multi-thread, so we have to use
// thread safe sync.Map.
var stargateWhitelist sync.Map

func init() {
	// market
	setWhitelistedQuery("/terra.market.v1beta1.Query/Swap", &markettypes.QuerySwapResponse{})

	// treasury
	setWhitelistedQuery("/terra.treasury.v1beta1.Query/TaxCap", &treasurytypes.QueryTaxCapResponse{})
	setWhitelistedQuery("/terra.treasury.v1beta1.Query/TaxRate", &treasurytypes.QueryTaxRateResponse{})

	// oracle
	setWhitelistedQuery("/terra.oracle.v1beta1.Query/ExchangeRate", &oracletypes.QueryExchangeRateResponse{})
}

// GetWhitelistedQuery returns the whitelisted query at the provided path.
// If the query does not exist, or it was setup wrong by the chain, this returns an error.
func GetWhitelistedQuery(queryPath string) (codec.ProtoMarshaler, error) {
	protoResponseAny, isWhitelisted := stargateWhitelist.Load(queryPath)
	if !isWhitelisted {
		return nil, wasmvmtypes.UnsupportedRequest{Kind: fmt.Sprintf("'%s' path is not allowed from the contract", queryPath)}
	}
	protoResponseType, ok := protoResponseAny.(codec.ProtoMarshaler)
	if !ok {
		return nil, wasmvmtypes.Unknown{}
	}
	return protoResponseType, nil
}

func setWhitelistedQuery(queryPath string, protoType codec.ProtoMarshaler) {
	stargateWhitelist.Store(queryPath, protoType)
}
