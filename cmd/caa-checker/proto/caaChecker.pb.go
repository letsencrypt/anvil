// Code generated by protoc-gen-go.
// source: caaChecker.proto
// DO NOT EDIT!

/*
Package caaChecker is a generated protocol buffer package.

It is generated from these files:
	caaChecker.proto

It has these top-level messages:
	Check
	Result
*/
package caaChecker

import proto "github.com/letsencrypt/boulder/Godeps/_workspace/src/github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "github.com/letsencrypt/boulder/Godeps/_workspace/src/golang.org/x/net/context"
	grpc "github.com/letsencrypt/boulder/Godeps/_workspace/src/google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
const _ = proto.ProtoPackageIsVersion1

type Check struct {
	Name             *string `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	IssuerDomain     *string `protobuf:"bytes,2,opt,name=issuerDomain" json:"issuerDomain,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *Check) Reset()                    { *m = Check{} }
func (m *Check) String() string            { return proto.CompactTextString(m) }
func (*Check) ProtoMessage()               {}
func (*Check) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *Check) GetName() string {
	if m != nil && m.Name != nil {
		return *m.Name
	}
	return ""
}

func (m *Check) GetIssuerDomain() string {
	if m != nil && m.IssuerDomain != nil {
		return *m.IssuerDomain
	}
	return ""
}

type Result struct {
	Present          *bool  `protobuf:"varint,1,opt,name=present" json:"present,omitempty"`
	Valid            *bool  `protobuf:"varint,2,opt,name=valid" json:"valid,omitempty"`
	XXX_unrecognized []byte `json:"-"`
}

func (m *Result) Reset()                    { *m = Result{} }
func (m *Result) String() string            { return proto.CompactTextString(m) }
func (*Result) ProtoMessage()               {}
func (*Result) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *Result) GetPresent() bool {
	if m != nil && m.Present != nil {
		return *m.Present
	}
	return false
}

func (m *Result) GetValid() bool {
	if m != nil && m.Valid != nil {
		return *m.Valid
	}
	return false
}

func init() {
	proto.RegisterType((*Check)(nil), "Check")
	proto.RegisterType((*Result)(nil), "Result")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion1

// Client API for CAAChecker service

type CAACheckerClient interface {
	ValidForIssuance(ctx context.Context, in *Check, opts ...grpc.CallOption) (*Result, error)
}

type cAACheckerClient struct {
	cc *grpc.ClientConn
}

func NewCAACheckerClient(cc *grpc.ClientConn) CAACheckerClient {
	return &cAACheckerClient{cc}
}

func (c *cAACheckerClient) ValidForIssuance(ctx context.Context, in *Check, opts ...grpc.CallOption) (*Result, error) {
	out := new(Result)
	err := grpc.Invoke(ctx, "/CAAChecker/ValidForIssuance", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for CAAChecker service

type CAACheckerServer interface {
	ValidForIssuance(context.Context, *Check) (*Result, error)
}

func RegisterCAACheckerServer(s *grpc.Server, srv CAACheckerServer) {
	s.RegisterService(&_CAAChecker_serviceDesc, srv)
}

func _CAAChecker_ValidForIssuance_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error) (interface{}, error) {
	in := new(Check)
	if err := dec(in); err != nil {
		return nil, err
	}
	out, err := srv.(CAACheckerServer).ValidForIssuance(ctx, in)
	if err != nil {
		return nil, err
	}
	return out, nil
}

var _CAAChecker_serviceDesc = grpc.ServiceDesc{
	ServiceName: "CAAChecker",
	HandlerType: (*CAACheckerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ValidForIssuance",
			Handler:    _CAAChecker_ValidForIssuance_Handler,
		},
	},
	Streams: []grpc.StreamDesc{},
}

var fileDescriptor0 = []byte{
	// 167 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xe2, 0x12, 0x48, 0x4e, 0x4c, 0x74,
	0xce, 0x48, 0x4d, 0xce, 0x4e, 0x2d, 0xd2, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x57, 0xb2, 0xe7, 0x62,
	0x05, 0x0b, 0x08, 0x09, 0x71, 0xb1, 0xe4, 0x25, 0xe6, 0xa6, 0x4a, 0x30, 0x2a, 0x30, 0x6a, 0x70,
	0x06, 0x81, 0xd9, 0x42, 0x4a, 0x5c, 0x3c, 0x99, 0xc5, 0xc5, 0xa5, 0xa9, 0x45, 0x2e, 0xf9, 0xb9,
	0x89, 0x99, 0x79, 0x12, 0x4c, 0x60, 0x39, 0x14, 0x31, 0x25, 0x0b, 0x2e, 0xb6, 0xa0, 0xd4, 0xe2,
	0xd2, 0x9c, 0x12, 0x21, 0x09, 0x2e, 0xf6, 0x82, 0xa2, 0xd4, 0xe2, 0xd4, 0xbc, 0x12, 0xb0, 0x21,
	0x1c, 0x41, 0x30, 0xae, 0x90, 0x08, 0x17, 0x6b, 0x59, 0x62, 0x4e, 0x66, 0x0a, 0xd8, 0x00, 0x8e,
	0x20, 0x08, 0xc7, 0xc8, 0x98, 0x8b, 0xcb, 0xd9, 0xd1, 0x11, 0xea, 0x1c, 0x21, 0x55, 0x2e, 0x81,
	0x30, 0x90, 0xb0, 0x5b, 0x7e, 0x91, 0x27, 0xd0, 0xfc, 0xc4, 0xbc, 0xe4, 0x54, 0x21, 0x36, 0x3d,
	0xb0, 0xac, 0x14, 0xbb, 0x1e, 0xc4, 0x0a, 0x25, 0x06, 0x40, 0x00, 0x00, 0x00, 0xff, 0xff, 0x17,
	0xdb, 0xa7, 0x3b, 0xc2, 0x00, 0x00, 0x00,
}
