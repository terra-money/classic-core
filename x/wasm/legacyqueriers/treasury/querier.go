package treasury

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/encoding"
	"google.golang.org/grpc/encoding/proto"
)

var protoCodec = encoding.GetCodec(proto.Name)

// QueryLegacyTreasury return empty response for backward compatibility
func QueryLegacyTreasury(path string) (bz []byte, err error) {
	switch path {
	case "/terra.treasury.v1beta1.Query/TaxRate":
		bz, err = protoCodec.Marshal(&QueryTaxRateResponse{TaxRate: sdk.ZeroDec()})
	case "/terra.treasury.v1beta1.Query/TaxCap":
		bz, err = protoCodec.Marshal(&QueryTaxCapResponse{TaxCap: sdk.ZeroInt()})
	case "/terra.treasury.v1beta1.Query/TaxCaps":
		var taxCaps []QueryTaxCapsResponseItem
		bz, err = protoCodec.Marshal(&QueryTaxCapsResponse{TaxCaps: taxCaps})
	case "/terra.treasury.v1beta1.Query/RewardWeight":
		bz, err = protoCodec.Marshal(&QueryRewardWeightResponse{RewardWeight: sdk.ZeroDec()})
	case "/terra.treasury.v1beta1.Query/SeigniorageProceeds":
		bz, err = protoCodec.Marshal(&QuerySeigniorageProceedsResponse{SeigniorageProceeds: sdk.ZeroInt()})
	case "/terra.treasury.v1beta1.Query/TaxProceeds":
		bz, err = protoCodec.Marshal(&QueryTaxProceedsResponse{TaxProceeds: sdk.Coins{}})
	case "/terra.treasury.v1beta1.Query/Indicators":
		bz, err = protoCodec.Marshal(&QueryIndicatorsResponse{
			TRLYear:  sdk.ZeroDec(),
			TRLMonth: sdk.ZeroDec(),
		})
	case "/terra.treasury.v1beta1.Query/Params":
		bz, err = protoCodec.Marshal(&QueryParamsResponse{Params: DefaultParams()})
	}

	return bz, err
}
