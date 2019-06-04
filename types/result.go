package types

import (
	sdktypes "github.com/cosmos/cosmos-sdk/types"
)

// TxSearchReponse defines tx search response structure
type TxSearchResponse struct {
	Txs        []sdktypes.TxResponse `json:"txs"`
	TotalCount int                   `json:"total_count"`
}

// NewTxSearchResponse returns a TxSearchResponse object
func NewTxSearchResponse(txs []sdktypes.TxResponse, totalCount int) TxSearchResponse {
	return TxSearchResponse{
		Txs:        txs,
		TotalCount: totalCount,
	}
}
