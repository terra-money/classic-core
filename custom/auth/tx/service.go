package tx

import (
	"context"

	gogogrpc "github.com/gogo/protobuf/grpc"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	customante "github.com/terra-money/core/custom/auth/ante"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ ServiceServer = txServer{}

// txServer is the server for the protobuf Tx service.
type txServer struct {
	treasuryKeeper customante.TreasuryKeeper
}

// NewTxServer creates a new Tx service server.
func NewTxServer(treasuryKeeper customante.TreasuryKeeper) ServiceServer {
	return txServer{
		treasuryKeeper: treasuryKeeper,
	}
}

// ComputeTax implements the ServiceServer.ComputeTax RPC method.
func (ts txServer) ComputeTax(c context.Context, req *ComputeTaxRequest) (*ComputeTaxResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be nil")
	}

	taxAmount := customante.FilterMsgAndComputeTax(ctx, ts.treasuryKeeper, req.Tx.GetMsgs()...)
	return &ComputeTaxResponse{
		TaxAmount: taxAmount,
	}, nil
}

// RegisterTxService registers the tx service on the gRPC router.
func RegisterTxService(
	qrt gogogrpc.Server,
	treasuryKeeper customante.TreasuryKeeper,
) {
	RegisterServiceServer(
		qrt,
		NewTxServer(treasuryKeeper),
	)
}

// RegisterGRPCGatewayRoutes mounts the tx service's GRPC-gateway routes on the
// given Mux.
func RegisterGRPCGatewayRoutes(clientConn gogogrpc.ClientConn, mux *runtime.ServeMux) {
	_ = RegisterServiceHandlerClient(context.Background(), mux, NewServiceClient(clientConn))
}

var _ codectypes.UnpackInterfacesMessage = ComputeTaxRequest{}

// UnpackInterfaces implements the UnpackInterfacesMessage interface.
func (m ComputeTaxRequest) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	return m.Tx.UnpackInterfaces(unpacker)
}
