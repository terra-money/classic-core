// nolint
package exported

import "github.com/terra-money/core/x/wasm/types"

var (
	EncodeSdkCoin  = types.EncodeSdkCoin
	EncodeSdkCoins = types.EncodeSdkCoins
	ParseToCoin    = types.ParseToCoin
	ParseToCoins   = types.ParseToCoins

	ErrInvalidMsg = types.ErrInvalidMsg
)

type (
	WasmMsgParserInterface = types.WasmMsgParserInterface
	WasmQuerierInterface   = types.WasmQuerierInterface
	MsgInstantiateContract = types.MsgInstantiateContract
	MsgExecuteContract     = types.MsgExecuteContract
	MsgStoreCode           = types.MsgStoreCode
)
