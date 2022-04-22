// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: terra/wasm/v1beta1/stargatelayer/wasm/query.proto

package wasm

import (
	encoding_json "encoding/json"
	fmt "fmt"
	_ "github.com/cosmos/cosmos-sdk/types"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
	types "github.com/terra-money/core/x/wasm/types"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	io "io"
	math "math"
	math_bits "math/bits"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

// QueryCodeInfoResponse is response type for the
// QueryyCodeInfo RPC method.
type QueryCodeInfoResponse struct {
	CodeInfo types.CodeInfo `protobuf:"bytes,1,opt,name=code_info,json=codeInfo,proto3" json:"code_info"`
}

func (m *QueryCodeInfoResponse) Reset()         { *m = QueryCodeInfoResponse{} }
func (m *QueryCodeInfoResponse) String() string { return proto.CompactTextString(m) }
func (*QueryCodeInfoResponse) ProtoMessage()    {}
func (*QueryCodeInfoResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_d777d3128e9826c9, []int{0}
}
func (m *QueryCodeInfoResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryCodeInfoResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryCodeInfoResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryCodeInfoResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryCodeInfoResponse.Merge(m, src)
}
func (m *QueryCodeInfoResponse) XXX_Size() int {
	return m.Size()
}
func (m *QueryCodeInfoResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryCodeInfoResponse.DiscardUnknown(m)
}

var xxx_messageInfo_QueryCodeInfoResponse proto.InternalMessageInfo

func (m *QueryCodeInfoResponse) GetCodeInfo() types.CodeInfo {
	if m != nil {
		return m.CodeInfo
	}
	return types.CodeInfo{}
}

// QueryByteCodeResponse is response type for the
// QueryyByteCode RPC method.
type QueryByteCodeResponse struct {
	ByteCode []byte `protobuf:"bytes,1,opt,name=byte_code,json=byteCode,proto3" json:"byte_code,omitempty"`
}

func (m *QueryByteCodeResponse) Reset()         { *m = QueryByteCodeResponse{} }
func (m *QueryByteCodeResponse) String() string { return proto.CompactTextString(m) }
func (*QueryByteCodeResponse) ProtoMessage()    {}
func (*QueryByteCodeResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_d777d3128e9826c9, []int{1}
}
func (m *QueryByteCodeResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryByteCodeResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryByteCodeResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryByteCodeResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryByteCodeResponse.Merge(m, src)
}
func (m *QueryByteCodeResponse) XXX_Size() int {
	return m.Size()
}
func (m *QueryByteCodeResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryByteCodeResponse.DiscardUnknown(m)
}

var xxx_messageInfo_QueryByteCodeResponse proto.InternalMessageInfo

func (m *QueryByteCodeResponse) GetByteCode() []byte {
	if m != nil {
		return m.ByteCode
	}
	return nil
}

// QueryContractInfoResponse is response type for the
// Query/ContractInfo RPC method.
type QueryContractInfoResponse struct {
	ContractInfo types.ContractInfo `protobuf:"bytes,1,opt,name=contract_info,json=contractInfo,proto3" json:"contract_info"`
}

func (m *QueryContractInfoResponse) Reset()         { *m = QueryContractInfoResponse{} }
func (m *QueryContractInfoResponse) String() string { return proto.CompactTextString(m) }
func (*QueryContractInfoResponse) ProtoMessage()    {}
func (*QueryContractInfoResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_d777d3128e9826c9, []int{2}
}
func (m *QueryContractInfoResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryContractInfoResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryContractInfoResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryContractInfoResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryContractInfoResponse.Merge(m, src)
}
func (m *QueryContractInfoResponse) XXX_Size() int {
	return m.Size()
}
func (m *QueryContractInfoResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryContractInfoResponse.DiscardUnknown(m)
}

var xxx_messageInfo_QueryContractInfoResponse proto.InternalMessageInfo

func (m *QueryContractInfoResponse) GetContractInfo() types.ContractInfo {
	if m != nil {
		return m.ContractInfo
	}
	return types.ContractInfo{}
}

// QueryContractStoreResponse is response type for the
// Query/ContractStore RPC method.
type QueryContractStoreResponse struct {
	QueryResult encoding_json.RawMessage `protobuf:"bytes,1,opt,name=query_result,json=queryResult,proto3,casttype=encoding/json.RawMessage" json:"query_result,omitempty"`
}

func (m *QueryContractStoreResponse) Reset()         { *m = QueryContractStoreResponse{} }
func (m *QueryContractStoreResponse) String() string { return proto.CompactTextString(m) }
func (*QueryContractStoreResponse) ProtoMessage()    {}
func (*QueryContractStoreResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_d777d3128e9826c9, []int{3}
}
func (m *QueryContractStoreResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryContractStoreResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryContractStoreResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryContractStoreResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryContractStoreResponse.Merge(m, src)
}
func (m *QueryContractStoreResponse) XXX_Size() int {
	return m.Size()
}
func (m *QueryContractStoreResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryContractStoreResponse.DiscardUnknown(m)
}

var xxx_messageInfo_QueryContractStoreResponse proto.InternalMessageInfo

func (m *QueryContractStoreResponse) GetQueryResult() encoding_json.RawMessage {
	if m != nil {
		return m.QueryResult
	}
	return nil
}

// QueryRawStoreResponse is response type for the
// Query/RawStore RPC method.
type QueryRawStoreResponse struct {
	// Data contains the raw store data
	Data []byte `protobuf:"bytes,1,opt,name=data,proto3" json:"data,omitempty"`
}

func (m *QueryRawStoreResponse) Reset()         { *m = QueryRawStoreResponse{} }
func (m *QueryRawStoreResponse) String() string { return proto.CompactTextString(m) }
func (*QueryRawStoreResponse) ProtoMessage()    {}
func (*QueryRawStoreResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_d777d3128e9826c9, []int{4}
}
func (m *QueryRawStoreResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryRawStoreResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryRawStoreResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryRawStoreResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryRawStoreResponse.Merge(m, src)
}
func (m *QueryRawStoreResponse) XXX_Size() int {
	return m.Size()
}
func (m *QueryRawStoreResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryRawStoreResponse.DiscardUnknown(m)
}

var xxx_messageInfo_QueryRawStoreResponse proto.InternalMessageInfo

func (m *QueryRawStoreResponse) GetData() []byte {
	if m != nil {
		return m.Data
	}
	return nil
}

// QueryParamsResponse is the response type for the Query/Params RPC method.
type QueryParamsResponse struct {
	// params defines the parameters of the module.
	Params types.Params `protobuf:"bytes,1,opt,name=params,proto3" json:"params"`
}

func (m *QueryParamsResponse) Reset()         { *m = QueryParamsResponse{} }
func (m *QueryParamsResponse) String() string { return proto.CompactTextString(m) }
func (*QueryParamsResponse) ProtoMessage()    {}
func (*QueryParamsResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_d777d3128e9826c9, []int{5}
}
func (m *QueryParamsResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryParamsResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryParamsResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryParamsResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryParamsResponse.Merge(m, src)
}
func (m *QueryParamsResponse) XXX_Size() int {
	return m.Size()
}
func (m *QueryParamsResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryParamsResponse.DiscardUnknown(m)
}

var xxx_messageInfo_QueryParamsResponse proto.InternalMessageInfo

func (m *QueryParamsResponse) GetParams() types.Params {
	if m != nil {
		return m.Params
	}
	return types.Params{}
}

func init() {
	proto.RegisterType((*QueryCodeInfoResponse)(nil), "terra.wasm.v1beta1.stargatelayer.wasm.QueryCodeInfoResponse")
	proto.RegisterType((*QueryByteCodeResponse)(nil), "terra.wasm.v1beta1.stargatelayer.wasm.QueryByteCodeResponse")
	proto.RegisterType((*QueryContractInfoResponse)(nil), "terra.wasm.v1beta1.stargatelayer.wasm.QueryContractInfoResponse")
	proto.RegisterType((*QueryContractStoreResponse)(nil), "terra.wasm.v1beta1.stargatelayer.wasm.QueryContractStoreResponse")
	proto.RegisterType((*QueryRawStoreResponse)(nil), "terra.wasm.v1beta1.stargatelayer.wasm.QueryRawStoreResponse")
	proto.RegisterType((*QueryParamsResponse)(nil), "terra.wasm.v1beta1.stargatelayer.wasm.QueryParamsResponse")
}

func init() {
	proto.RegisterFile("terra/wasm/v1beta1/stargatelayer/wasm/query.proto", fileDescriptor_d777d3128e9826c9)
}

var fileDescriptor_d777d3128e9826c9 = []byte{
	// 427 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x52, 0x4f, 0x6b, 0xd4, 0x40,
	0x14, 0xdf, 0x40, 0x29, 0xed, 0x74, 0xbd, 0xac, 0x0a, 0x35, 0xae, 0x69, 0x09, 0x08, 0x82, 0x98,
	0x61, 0xfd, 0x03, 0xde, 0x84, 0x78, 0x12, 0x91, 0x6a, 0xbc, 0x88, 0x20, 0xcb, 0xcb, 0xe4, 0x75,
	0x1a, 0xd9, 0xcc, 0x8b, 0x33, 0xb3, 0xae, 0xf9, 0x16, 0x7e, 0xac, 0x1e, 0x7b, 0xf4, 0x54, 0x64,
	0xf7, 0x5b, 0x78, 0x92, 0x4c, 0xc6, 0x75, 0x43, 0x73, 0x9b, 0xbc, 0xdf, 0x9f, 0xf7, 0xfb, 0x91,
	0xc7, 0x66, 0x16, 0xb5, 0x06, 0xbe, 0x02, 0x53, 0xf1, 0xef, 0xb3, 0x1c, 0x2d, 0xcc, 0xb8, 0xb1,
	0xa0, 0x25, 0x58, 0x5c, 0x40, 0x83, 0xba, 0x83, 0xbe, 0x2d, 0x51, 0x37, 0x49, 0xad, 0xc9, 0xd2,
	0xe4, 0xa1, 0x93, 0x24, 0xed, 0x3c, 0xf1, 0x92, 0xa4, 0x27, 0x71, 0x50, 0x78, 0x47, 0x92, 0x24,
	0xa7, 0xe0, 0xed, 0xab, 0x13, 0x87, 0x53, 0x49, 0x24, 0x17, 0xc8, 0xa1, 0x2e, 0x39, 0x28, 0x45,
	0x16, 0x6c, 0x49, 0xca, 0x78, 0xf4, 0xc1, 0x40, 0x1a, 0xb7, 0xa7, 0x83, 0x23, 0x41, 0xa6, 0x22,
	0xc3, 0x73, 0x30, 0xb8, 0xc5, 0x05, 0x95, 0xaa, 0xc3, 0xe3, 0x4f, 0xec, 0xee, 0x87, 0x36, 0xe8,
	0x6b, 0x2a, 0xf0, 0x8d, 0x3a, 0xa7, 0x0c, 0x4d, 0x4d, 0xca, 0xe0, 0xe4, 0x15, 0x3b, 0x14, 0x54,
	0xe0, 0xbc, 0x54, 0xe7, 0x74, 0x1c, 0x9c, 0x06, 0x8f, 0x8e, 0x9e, 0x4e, 0x93, 0x81, 0x1a, 0xff,
	0x84, 0xe9, 0xde, 0xe5, 0xf5, 0xc9, 0x28, 0x3b, 0x10, 0xfe, 0x3b, 0x7e, 0xee, 0x9d, 0xd3, 0xc6,
	0x62, 0x4b, 0xda, 0x3a, 0xdf, 0x67, 0x87, 0x79, 0x63, 0x71, 0xde, 0x32, 0x9d, 0xf3, 0x38, 0x3b,
	0xc8, 0x3d, 0x29, 0xbe, 0x60, 0xf7, 0x7c, 0x1e, 0x65, 0x35, 0x08, 0xdb, 0xcb, 0xf4, 0x96, 0xdd,
	0x12, 0x7e, 0xbe, 0x9b, 0xeb, 0x74, 0x38, 0xd7, 0x7f, 0x03, 0x9f, 0x6d, 0x2c, 0x76, 0x66, 0xf1,
	0x17, 0x16, 0xf6, 0x36, 0x7d, 0xb4, 0xa4, 0x71, 0xa7, 0xfe, 0xd8, 0xfd, 0xc0, 0xb9, 0x46, 0xb3,
	0x5c, 0xd8, 0x2e, 0x67, 0x3a, 0xfd, 0x73, 0x7d, 0x72, 0x8c, 0x4a, 0x50, 0x51, 0x2a, 0xc9, 0xbf,
	0x1a, 0x52, 0x49, 0x06, 0xab, 0x77, 0x68, 0x0c, 0x48, 0xcc, 0x8e, 0x9c, 0x22, 0x73, 0x82, 0xf8,
	0xb1, 0xaf, 0x9f, 0xc1, 0xaa, 0xef, 0x3c, 0x61, 0x7b, 0x05, 0x58, 0xf0, 0xcd, 0xdd, 0x3b, 0x3e,
	0x63, 0xb7, 0x1d, 0xf9, 0x3d, 0x68, 0xa8, 0xcc, 0x96, 0xfa, 0x92, 0xed, 0xd7, 0x6e, 0xe2, 0x8b,
	0x86, 0x43, 0x45, 0x3b, 0x8d, 0xaf, 0xe8, 0xf9, 0xe9, 0xd9, 0xe5, 0x3a, 0x0a, 0xae, 0xd6, 0x51,
	0xf0, 0x7b, 0x1d, 0x05, 0x3f, 0x37, 0xd1, 0xe8, 0x6a, 0x13, 0x8d, 0x7e, 0x6d, 0xa2, 0xd1, 0xe7,
	0x17, 0xb2, 0xb4, 0x17, 0xcb, 0x3c, 0x11, 0x54, 0x71, 0xe7, 0xf6, 0xa4, 0x22, 0x85, 0x0d, 0x17,
	0xa4, 0x91, 0xff, 0xe8, 0xee, 0xe8, 0xe6, 0x35, 0xe7, 0xfb, 0xee, 0x5c, 0x9e, 0xfd, 0x0d, 0x00,
	0x00, 0xff, 0xff, 0x1c, 0x92, 0xc3, 0xfb, 0xfd, 0x02, 0x00, 0x00,
}

func (m *QueryCodeInfoResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryCodeInfoResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryCodeInfoResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size, err := m.CodeInfo.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintQuery(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0xa
	return len(dAtA) - i, nil
}

func (m *QueryByteCodeResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryByteCodeResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryByteCodeResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.ByteCode) > 0 {
		i -= len(m.ByteCode)
		copy(dAtA[i:], m.ByteCode)
		i = encodeVarintQuery(dAtA, i, uint64(len(m.ByteCode)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *QueryContractInfoResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryContractInfoResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryContractInfoResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size, err := m.ContractInfo.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintQuery(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0xa
	return len(dAtA) - i, nil
}

func (m *QueryContractStoreResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryContractStoreResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryContractStoreResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.QueryResult) > 0 {
		i -= len(m.QueryResult)
		copy(dAtA[i:], m.QueryResult)
		i = encodeVarintQuery(dAtA, i, uint64(len(m.QueryResult)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *QueryRawStoreResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryRawStoreResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryRawStoreResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Data) > 0 {
		i -= len(m.Data)
		copy(dAtA[i:], m.Data)
		i = encodeVarintQuery(dAtA, i, uint64(len(m.Data)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *QueryParamsResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryParamsResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryParamsResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size, err := m.Params.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintQuery(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0xa
	return len(dAtA) - i, nil
}

func encodeVarintQuery(dAtA []byte, offset int, v uint64) int {
	offset -= sovQuery(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *QueryCodeInfoResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.CodeInfo.Size()
	n += 1 + l + sovQuery(uint64(l))
	return n
}

func (m *QueryByteCodeResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.ByteCode)
	if l > 0 {
		n += 1 + l + sovQuery(uint64(l))
	}
	return n
}

func (m *QueryContractInfoResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.ContractInfo.Size()
	n += 1 + l + sovQuery(uint64(l))
	return n
}

func (m *QueryContractStoreResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.QueryResult)
	if l > 0 {
		n += 1 + l + sovQuery(uint64(l))
	}
	return n
}

func (m *QueryRawStoreResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Data)
	if l > 0 {
		n += 1 + l + sovQuery(uint64(l))
	}
	return n
}

func (m *QueryParamsResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.Params.Size()
	n += 1 + l + sovQuery(uint64(l))
	return n
}

func sovQuery(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozQuery(x uint64) (n int) {
	return sovQuery(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *QueryCodeInfoResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: QueryCodeInfoResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryCodeInfoResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field CodeInfo", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.CodeInfo.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *QueryByteCodeResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: QueryByteCodeResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryByteCodeResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ByteCode", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ByteCode = append(m.ByteCode[:0], dAtA[iNdEx:postIndex]...)
			if m.ByteCode == nil {
				m.ByteCode = []byte{}
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *QueryContractInfoResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: QueryContractInfoResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryContractInfoResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ContractInfo", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.ContractInfo.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *QueryContractStoreResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: QueryContractStoreResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryContractStoreResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field QueryResult", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.QueryResult = append(m.QueryResult[:0], dAtA[iNdEx:postIndex]...)
			if m.QueryResult == nil {
				m.QueryResult = []byte{}
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *QueryRawStoreResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: QueryRawStoreResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryRawStoreResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Data", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Data = append(m.Data[:0], dAtA[iNdEx:postIndex]...)
			if m.Data == nil {
				m.Data = []byte{}
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *QueryParamsResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: QueryParamsResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryParamsResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Params", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Params.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipQuery(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowQuery
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthQuery
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupQuery
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthQuery
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthQuery        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowQuery          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupQuery = fmt.Errorf("proto: unexpected end of group")
)
