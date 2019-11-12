// Code generated by protoc-gen-go. DO NOT EDIT.
// source: ca/proto/ca.proto

package proto

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	proto1 "github.com/letsencrypt/boulder/core/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type IssueCertificateRequest struct {
	Csr                  []byte   `protobuf:"bytes,1,opt,name=csr" json:"csr,omitempty"`
	RegistrationID       *int64   `protobuf:"varint,2,opt,name=registrationID" json:"registrationID,omitempty"`
	OrderID              *int64   `protobuf:"varint,3,opt,name=orderID" json:"orderID,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *IssueCertificateRequest) Reset()         { *m = IssueCertificateRequest{} }
func (m *IssueCertificateRequest) String() string { return proto.CompactTextString(m) }
func (*IssueCertificateRequest) ProtoMessage()    {}
func (*IssueCertificateRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_8f9fdc2529716820, []int{0}
}

func (m *IssueCertificateRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_IssueCertificateRequest.Unmarshal(m, b)
}
func (m *IssueCertificateRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_IssueCertificateRequest.Marshal(b, m, deterministic)
}
func (m *IssueCertificateRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_IssueCertificateRequest.Merge(m, src)
}
func (m *IssueCertificateRequest) XXX_Size() int {
	return xxx_messageInfo_IssueCertificateRequest.Size(m)
}
func (m *IssueCertificateRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_IssueCertificateRequest.DiscardUnknown(m)
}

var xxx_messageInfo_IssueCertificateRequest proto.InternalMessageInfo

func (m *IssueCertificateRequest) GetCsr() []byte {
	if m != nil {
		return m.Csr
	}
	return nil
}

func (m *IssueCertificateRequest) GetRegistrationID() int64 {
	if m != nil && m.RegistrationID != nil {
		return *m.RegistrationID
	}
	return 0
}

func (m *IssueCertificateRequest) GetOrderID() int64 {
	if m != nil && m.OrderID != nil {
		return *m.OrderID
	}
	return 0
}

type IssuePrecertificateResponse struct {
	DER                  []byte   `protobuf:"bytes,1,opt,name=DER" json:"DER,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *IssuePrecertificateResponse) Reset()         { *m = IssuePrecertificateResponse{} }
func (m *IssuePrecertificateResponse) String() string { return proto.CompactTextString(m) }
func (*IssuePrecertificateResponse) ProtoMessage()    {}
func (*IssuePrecertificateResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_8f9fdc2529716820, []int{1}
}

func (m *IssuePrecertificateResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_IssuePrecertificateResponse.Unmarshal(m, b)
}
func (m *IssuePrecertificateResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_IssuePrecertificateResponse.Marshal(b, m, deterministic)
}
func (m *IssuePrecertificateResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_IssuePrecertificateResponse.Merge(m, src)
}
func (m *IssuePrecertificateResponse) XXX_Size() int {
	return xxx_messageInfo_IssuePrecertificateResponse.Size(m)
}
func (m *IssuePrecertificateResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_IssuePrecertificateResponse.DiscardUnknown(m)
}

var xxx_messageInfo_IssuePrecertificateResponse proto.InternalMessageInfo

func (m *IssuePrecertificateResponse) GetDER() []byte {
	if m != nil {
		return m.DER
	}
	return nil
}

type IssueCertificateForPrecertificateRequest struct {
	DER                  []byte   `protobuf:"bytes,1,opt,name=DER" json:"DER,omitempty"`
	SCTs                 [][]byte `protobuf:"bytes,2,rep,name=SCTs" json:"SCTs,omitempty"`
	RegistrationID       *int64   `protobuf:"varint,3,opt,name=registrationID" json:"registrationID,omitempty"`
	OrderID              *int64   `protobuf:"varint,4,opt,name=orderID" json:"orderID,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *IssueCertificateForPrecertificateRequest) Reset() {
	*m = IssueCertificateForPrecertificateRequest{}
}
func (m *IssueCertificateForPrecertificateRequest) String() string { return proto.CompactTextString(m) }
func (*IssueCertificateForPrecertificateRequest) ProtoMessage()    {}
func (*IssueCertificateForPrecertificateRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_8f9fdc2529716820, []int{2}
}

func (m *IssueCertificateForPrecertificateRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_IssueCertificateForPrecertificateRequest.Unmarshal(m, b)
}
func (m *IssueCertificateForPrecertificateRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_IssueCertificateForPrecertificateRequest.Marshal(b, m, deterministic)
}
func (m *IssueCertificateForPrecertificateRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_IssueCertificateForPrecertificateRequest.Merge(m, src)
}
func (m *IssueCertificateForPrecertificateRequest) XXX_Size() int {
	return xxx_messageInfo_IssueCertificateForPrecertificateRequest.Size(m)
}
func (m *IssueCertificateForPrecertificateRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_IssueCertificateForPrecertificateRequest.DiscardUnknown(m)
}

var xxx_messageInfo_IssueCertificateForPrecertificateRequest proto.InternalMessageInfo

func (m *IssueCertificateForPrecertificateRequest) GetDER() []byte {
	if m != nil {
		return m.DER
	}
	return nil
}

func (m *IssueCertificateForPrecertificateRequest) GetSCTs() [][]byte {
	if m != nil {
		return m.SCTs
	}
	return nil
}

func (m *IssueCertificateForPrecertificateRequest) GetRegistrationID() int64 {
	if m != nil && m.RegistrationID != nil {
		return *m.RegistrationID
	}
	return 0
}

func (m *IssueCertificateForPrecertificateRequest) GetOrderID() int64 {
	if m != nil && m.OrderID != nil {
		return *m.OrderID
	}
	return 0
}

type GenerateOCSPRequest struct {
	CertDER   []byte  `protobuf:"bytes,1,opt,name=certDER" json:"certDER,omitempty"`
	Status    *string `protobuf:"bytes,2,opt,name=status" json:"status,omitempty"`
	Reason    *int32  `protobuf:"varint,3,opt,name=reason" json:"reason,omitempty"`
	RevokedAt *int64  `protobuf:"varint,4,opt,name=revokedAt" json:"revokedAt,omitempty"`
	// If serial is set than certDER must not be set and issuerID must be set
	Serial               *string  `protobuf:"bytes,5,opt,name=serial" json:"serial,omitempty"`
	IssuerID             *int64   `protobuf:"varint,6,opt,name=issuerID" json:"issuerID,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GenerateOCSPRequest) Reset()         { *m = GenerateOCSPRequest{} }
func (m *GenerateOCSPRequest) String() string { return proto.CompactTextString(m) }
func (*GenerateOCSPRequest) ProtoMessage()    {}
func (*GenerateOCSPRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_8f9fdc2529716820, []int{3}
}

func (m *GenerateOCSPRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GenerateOCSPRequest.Unmarshal(m, b)
}
func (m *GenerateOCSPRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GenerateOCSPRequest.Marshal(b, m, deterministic)
}
func (m *GenerateOCSPRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GenerateOCSPRequest.Merge(m, src)
}
func (m *GenerateOCSPRequest) XXX_Size() int {
	return xxx_messageInfo_GenerateOCSPRequest.Size(m)
}
func (m *GenerateOCSPRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_GenerateOCSPRequest.DiscardUnknown(m)
}

var xxx_messageInfo_GenerateOCSPRequest proto.InternalMessageInfo

func (m *GenerateOCSPRequest) GetCertDER() []byte {
	if m != nil {
		return m.CertDER
	}
	return nil
}

func (m *GenerateOCSPRequest) GetStatus() string {
	if m != nil && m.Status != nil {
		return *m.Status
	}
	return ""
}

func (m *GenerateOCSPRequest) GetReason() int32 {
	if m != nil && m.Reason != nil {
		return *m.Reason
	}
	return 0
}

func (m *GenerateOCSPRequest) GetRevokedAt() int64 {
	if m != nil && m.RevokedAt != nil {
		return *m.RevokedAt
	}
	return 0
}

func (m *GenerateOCSPRequest) GetSerial() string {
	if m != nil && m.Serial != nil {
		return *m.Serial
	}
	return ""
}

func (m *GenerateOCSPRequest) GetIssuerID() int64 {
	if m != nil && m.IssuerID != nil {
		return *m.IssuerID
	}
	return 0
}

type OCSPResponse struct {
	Response             []byte   `protobuf:"bytes,1,opt,name=response" json:"response,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *OCSPResponse) Reset()         { *m = OCSPResponse{} }
func (m *OCSPResponse) String() string { return proto.CompactTextString(m) }
func (*OCSPResponse) ProtoMessage()    {}
func (*OCSPResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_8f9fdc2529716820, []int{4}
}

func (m *OCSPResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_OCSPResponse.Unmarshal(m, b)
}
func (m *OCSPResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_OCSPResponse.Marshal(b, m, deterministic)
}
func (m *OCSPResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_OCSPResponse.Merge(m, src)
}
func (m *OCSPResponse) XXX_Size() int {
	return xxx_messageInfo_OCSPResponse.Size(m)
}
func (m *OCSPResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_OCSPResponse.DiscardUnknown(m)
}

var xxx_messageInfo_OCSPResponse proto.InternalMessageInfo

func (m *OCSPResponse) GetResponse() []byte {
	if m != nil {
		return m.Response
	}
	return nil
}

func init() {
	proto.RegisterType((*IssueCertificateRequest)(nil), "ca.IssueCertificateRequest")
	proto.RegisterType((*IssuePrecertificateResponse)(nil), "ca.IssuePrecertificateResponse")
	proto.RegisterType((*IssueCertificateForPrecertificateRequest)(nil), "ca.IssueCertificateForPrecertificateRequest")
	proto.RegisterType((*GenerateOCSPRequest)(nil), "ca.GenerateOCSPRequest")
	proto.RegisterType((*OCSPResponse)(nil), "ca.OCSPResponse")
}

func init() { proto.RegisterFile("ca/proto/ca.proto", fileDescriptor_8f9fdc2529716820) }

var fileDescriptor_8f9fdc2529716820 = []byte{
	// 415 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xa4, 0x53, 0xcd, 0x6e, 0xd3, 0x40,
	0x10, 0x8e, 0xed, 0xa6, 0x69, 0x47, 0x06, 0xb5, 0x5b, 0xa0, 0x96, 0x8b, 0x44, 0xf0, 0x01, 0x59,
	0x08, 0x39, 0x52, 0xaf, 0x9c, 0x4a, 0x5c, 0x50, 0x24, 0x24, 0xaa, 0x2d, 0x5c, 0xb8, 0xad, 0xb6,
	0x53, 0xb0, 0x00, 0x6f, 0x99, 0x5d, 0x23, 0xf1, 0x1a, 0x7d, 0x13, 0xde, 0x10, 0xed, 0x66, 0xed,
	0x38, 0x91, 0xa3, 0x1c, 0x7a, 0x9b, 0x6f, 0x66, 0x67, 0xbe, 0x6f, 0x7e, 0x16, 0x8e, 0xa5, 0x98,
	0xdd, 0x91, 0x32, 0x6a, 0x26, 0x45, 0xe1, 0x0c, 0x16, 0x4a, 0x91, 0x3e, 0x95, 0x8a, 0xb0, 0x0d,
	0x28, 0xc2, 0x65, 0x28, 0xfb, 0x05, 0xa7, 0x0b, 0xad, 0x1b, 0x9c, 0x23, 0x99, 0xea, 0xb6, 0x92,
	0xc2, 0x20, 0xc7, 0xdf, 0x0d, 0x6a, 0xc3, 0x8e, 0x20, 0x92, 0x9a, 0x92, 0x60, 0x1a, 0xe4, 0x31,
	0xb7, 0x26, 0x7b, 0x05, 0x8f, 0x09, 0xbf, 0x55, 0xda, 0x90, 0x30, 0x95, 0xaa, 0x17, 0x65, 0x12,
	0x4e, 0x83, 0x3c, 0xe2, 0x1b, 0x5e, 0x96, 0xc0, 0x44, 0xd1, 0x0d, 0xd2, 0xa2, 0x4c, 0x22, 0xf7,
	0xa0, 0x85, 0xd9, 0x0c, 0xce, 0x1c, 0xdd, 0x15, 0xa1, 0xec, 0x33, 0xea, 0x3b, 0x55, 0x6b, 0xb4,
	0x94, 0xe5, 0x25, 0x6f, 0x29, 0xcb, 0x4b, 0x9e, 0xdd, 0x07, 0x90, 0x6f, 0x0a, 0x7c, 0xaf, 0x68,
	0x33, 0xbf, 0x53, 0xbc, 0x9e, 0xce, 0x18, 0xec, 0x5d, 0xcf, 0x3f, 0xeb, 0x24, 0x9c, 0x46, 0x79,
	0xcc, 0x9d, 0x3d, 0xd0, 0x45, 0xb4, 0xab, 0x8b, 0xbd, 0xf5, 0x2e, 0xfe, 0x05, 0x70, 0xf2, 0x01,
	0x6b, 0x24, 0x61, 0xf0, 0xd3, 0xfc, 0xfa, 0xaa, 0xe5, 0x4f, 0x60, 0x62, 0x55, 0xad, 0x34, 0xb4,
	0x90, 0x3d, 0x83, 0x7d, 0x6d, 0x84, 0x69, 0xb4, 0x9b, 0xd8, 0x21, 0xf7, 0xc8, 0xfa, 0x09, 0x85,
	0x56, 0xb5, 0xd3, 0x30, 0xe6, 0x1e, 0xb1, 0xe7, 0x70, 0x48, 0xf8, 0x47, 0xfd, 0xc0, 0x9b, 0x0b,
	0xe3, 0xd9, 0x57, 0x0e, 0x57, 0x0d, 0xa9, 0x12, 0x3f, 0x93, 0xb1, 0xaf, 0xe6, 0x10, 0x4b, 0xe1,
	0xa0, 0xb2, 0xb3, 0xb2, 0x92, 0xf7, 0x5d, 0x52, 0x87, 0xb3, 0xd7, 0x10, 0x2f, 0xa5, 0xfa, 0x51,
	0xa7, 0x70, 0x40, 0xde, 0xf6, 0x62, 0x3b, 0x7c, 0x7e, 0x1f, 0xc2, 0x93, 0xde, 0xbc, 0x2f, 0x1a,
	0xf3, 0x5d, 0x51, 0x65, 0xfe, 0xb2, 0x2f, 0x70, 0x32, 0xb0, 0x3e, 0x76, 0x56, 0x48, 0x51, 0x6c,
	0x39, 0xa3, 0xf4, 0x45, 0x17, 0x1c, 0x5e, 0x7a, 0x36, 0x62, 0xb7, 0xf0, 0x72, 0xe7, 0x8e, 0xd9,
	0x9b, 0x21, 0x92, 0x6d, 0xa7, 0x90, 0x1e, 0x17, 0xee, 0xc8, 0x7b, 0x4f, 0xb3, 0x11, 0x7b, 0x0b,
	0x71, 0x7f, 0x6d, 0xec, 0xd4, 0x96, 0x1c, 0x58, 0x64, 0x7a, 0x64, 0x03, 0xfd, 0x71, 0x65, 0xa3,
	0xf3, 0x8f, 0xf0, 0xc8, 0x7a, 0xfc, 0x73, 0x45, 0x0f, 0xaa, 0xf6, 0x6e, 0xf2, 0x75, 0xec, 0x3e,
	0xe0, 0xff, 0x00, 0x00, 0x00, 0xff, 0xff, 0xa3, 0xb6, 0x36, 0xac, 0xaf, 0x03, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// CertificateAuthorityClient is the client API for CertificateAuthority service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type CertificateAuthorityClient interface {
	IssuePrecertificate(ctx context.Context, in *IssueCertificateRequest, opts ...grpc.CallOption) (*IssuePrecertificateResponse, error)
	IssueCertificateForPrecertificate(ctx context.Context, in *IssueCertificateForPrecertificateRequest, opts ...grpc.CallOption) (*proto1.Certificate, error)
	GenerateOCSP(ctx context.Context, in *GenerateOCSPRequest, opts ...grpc.CallOption) (*OCSPResponse, error)
}

type certificateAuthorityClient struct {
	cc *grpc.ClientConn
}

func NewCertificateAuthorityClient(cc *grpc.ClientConn) CertificateAuthorityClient {
	return &certificateAuthorityClient{cc}
}

func (c *certificateAuthorityClient) IssuePrecertificate(ctx context.Context, in *IssueCertificateRequest, opts ...grpc.CallOption) (*IssuePrecertificateResponse, error) {
	out := new(IssuePrecertificateResponse)
	err := c.cc.Invoke(ctx, "/ca.CertificateAuthority/IssuePrecertificate", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *certificateAuthorityClient) IssueCertificateForPrecertificate(ctx context.Context, in *IssueCertificateForPrecertificateRequest, opts ...grpc.CallOption) (*proto1.Certificate, error) {
	out := new(proto1.Certificate)
	err := c.cc.Invoke(ctx, "/ca.CertificateAuthority/IssueCertificateForPrecertificate", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *certificateAuthorityClient) GenerateOCSP(ctx context.Context, in *GenerateOCSPRequest, opts ...grpc.CallOption) (*OCSPResponse, error) {
	out := new(OCSPResponse)
	err := c.cc.Invoke(ctx, "/ca.CertificateAuthority/GenerateOCSP", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CertificateAuthorityServer is the server API for CertificateAuthority service.
type CertificateAuthorityServer interface {
	IssuePrecertificate(context.Context, *IssueCertificateRequest) (*IssuePrecertificateResponse, error)
	IssueCertificateForPrecertificate(context.Context, *IssueCertificateForPrecertificateRequest) (*proto1.Certificate, error)
	GenerateOCSP(context.Context, *GenerateOCSPRequest) (*OCSPResponse, error)
}

// UnimplementedCertificateAuthorityServer can be embedded to have forward compatible implementations.
type UnimplementedCertificateAuthorityServer struct {
}

func (*UnimplementedCertificateAuthorityServer) IssuePrecertificate(ctx context.Context, req *IssueCertificateRequest) (*IssuePrecertificateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method IssuePrecertificate not implemented")
}
func (*UnimplementedCertificateAuthorityServer) IssueCertificateForPrecertificate(ctx context.Context, req *IssueCertificateForPrecertificateRequest) (*proto1.Certificate, error) {
	return nil, status.Errorf(codes.Unimplemented, "method IssueCertificateForPrecertificate not implemented")
}
func (*UnimplementedCertificateAuthorityServer) GenerateOCSP(ctx context.Context, req *GenerateOCSPRequest) (*OCSPResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GenerateOCSP not implemented")
}

func RegisterCertificateAuthorityServer(s *grpc.Server, srv CertificateAuthorityServer) {
	s.RegisterService(&_CertificateAuthority_serviceDesc, srv)
}

func _CertificateAuthority_IssuePrecertificate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(IssueCertificateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CertificateAuthorityServer).IssuePrecertificate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ca.CertificateAuthority/IssuePrecertificate",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CertificateAuthorityServer).IssuePrecertificate(ctx, req.(*IssueCertificateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CertificateAuthority_IssueCertificateForPrecertificate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(IssueCertificateForPrecertificateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CertificateAuthorityServer).IssueCertificateForPrecertificate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ca.CertificateAuthority/IssueCertificateForPrecertificate",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CertificateAuthorityServer).IssueCertificateForPrecertificate(ctx, req.(*IssueCertificateForPrecertificateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CertificateAuthority_GenerateOCSP_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GenerateOCSPRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CertificateAuthorityServer).GenerateOCSP(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ca.CertificateAuthority/GenerateOCSP",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CertificateAuthorityServer).GenerateOCSP(ctx, req.(*GenerateOCSPRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _CertificateAuthority_serviceDesc = grpc.ServiceDesc{
	ServiceName: "ca.CertificateAuthority",
	HandlerType: (*CertificateAuthorityServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "IssuePrecertificate",
			Handler:    _CertificateAuthority_IssuePrecertificate_Handler,
		},
		{
			MethodName: "IssueCertificateForPrecertificate",
			Handler:    _CertificateAuthority_IssueCertificateForPrecertificate_Handler,
		},
		{
			MethodName: "GenerateOCSP",
			Handler:    _CertificateAuthority_GenerateOCSP_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "ca/proto/ca.proto",
}

// OCSPGeneratorClient is the client API for OCSPGenerator service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type OCSPGeneratorClient interface {
	GenerateOCSP(ctx context.Context, in *GenerateOCSPRequest, opts ...grpc.CallOption) (*OCSPResponse, error)
}

type oCSPGeneratorClient struct {
	cc *grpc.ClientConn
}

func NewOCSPGeneratorClient(cc *grpc.ClientConn) OCSPGeneratorClient {
	return &oCSPGeneratorClient{cc}
}

func (c *oCSPGeneratorClient) GenerateOCSP(ctx context.Context, in *GenerateOCSPRequest, opts ...grpc.CallOption) (*OCSPResponse, error) {
	out := new(OCSPResponse)
	err := c.cc.Invoke(ctx, "/ca.OCSPGenerator/GenerateOCSP", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// OCSPGeneratorServer is the server API for OCSPGenerator service.
type OCSPGeneratorServer interface {
	GenerateOCSP(context.Context, *GenerateOCSPRequest) (*OCSPResponse, error)
}

// UnimplementedOCSPGeneratorServer can be embedded to have forward compatible implementations.
type UnimplementedOCSPGeneratorServer struct {
}

func (*UnimplementedOCSPGeneratorServer) GenerateOCSP(ctx context.Context, req *GenerateOCSPRequest) (*OCSPResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GenerateOCSP not implemented")
}

func RegisterOCSPGeneratorServer(s *grpc.Server, srv OCSPGeneratorServer) {
	s.RegisterService(&_OCSPGenerator_serviceDesc, srv)
}

func _OCSPGenerator_GenerateOCSP_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GenerateOCSPRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OCSPGeneratorServer).GenerateOCSP(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ca.OCSPGenerator/GenerateOCSP",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OCSPGeneratorServer).GenerateOCSP(ctx, req.(*GenerateOCSPRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _OCSPGenerator_serviceDesc = grpc.ServiceDesc{
	ServiceName: "ca.OCSPGenerator",
	HandlerType: (*OCSPGeneratorServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GenerateOCSP",
			Handler:    _OCSPGenerator_GenerateOCSP_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "ca/proto/ca.proto",
}
