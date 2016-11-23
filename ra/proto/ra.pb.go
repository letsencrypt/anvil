// Code generated by protoc-gen-go.
// source: ra/proto/ra.proto
// DO NOT EDIT!

/*
Package proto is a generated protocol buffer package.

It is generated from these files:
	ra/proto/ra.proto

It has these top-level messages:
	Registration
	Authorization
	NewAuthorizationRequest
	NewCertificateRequest
	UpdateRegistrationRequest
	UpdateAuthorizationRequest
	RevokeCertificateWithRegRequest
	AdministrativelyRevokeCertificateRequest
*/
package proto

import proto1 "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import core "github.com/letsencrypt/boulder/core/proto"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto1.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto1.ProtoPackageIsVersion2 // please upgrade the proto package

type Registration struct {
	Id               *int64   `protobuf:"varint,1,opt,name=id" json:"id,omitempty"`
	Key              []byte   `protobuf:"bytes,2,opt,name=key" json:"key,omitempty"`
	Contact          []string `protobuf:"bytes,3,rep,name=contact" json:"contact,omitempty"`
	Agreement        *string  `protobuf:"bytes,4,opt,name=agreement" json:"agreement,omitempty"`
	InitialIP        []byte   `protobuf:"bytes,5,opt,name=initialIP" json:"initialIP,omitempty"`
	CreatedAt        *int64   `protobuf:"varint,6,opt,name=createdAt" json:"createdAt,omitempty"`
	Status           *string  `protobuf:"bytes,7,opt,name=status" json:"status,omitempty"`
	XXX_unrecognized []byte   `json:"-"`
}

func (m *Registration) Reset()                    { *m = Registration{} }
func (m *Registration) String() string            { return proto1.CompactTextString(m) }
func (*Registration) ProtoMessage()               {}
func (*Registration) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *Registration) GetId() int64 {
	if m != nil && m.Id != nil {
		return *m.Id
	}
	return 0
}

func (m *Registration) GetKey() []byte {
	if m != nil {
		return m.Key
	}
	return nil
}

func (m *Registration) GetContact() []string {
	if m != nil {
		return m.Contact
	}
	return nil
}

func (m *Registration) GetAgreement() string {
	if m != nil && m.Agreement != nil {
		return *m.Agreement
	}
	return ""
}

func (m *Registration) GetInitialIP() []byte {
	if m != nil {
		return m.InitialIP
	}
	return nil
}

func (m *Registration) GetCreatedAt() int64 {
	if m != nil && m.CreatedAt != nil {
		return *m.CreatedAt
	}
	return 0
}

func (m *Registration) GetStatus() string {
	if m != nil && m.Status != nil {
		return *m.Status
	}
	return ""
}

type Authorization struct {
	Id               *string           `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	Identifier       *string           `protobuf:"bytes,2,opt,name=identifier" json:"identifier,omitempty"`
	RegistrationID   *int64            `protobuf:"varint,3,opt,name=registrationID" json:"registrationID,omitempty"`
	Status           *string           `protobuf:"bytes,4,opt,name=status" json:"status,omitempty"`
	Expires          *int64            `protobuf:"varint,5,opt,name=expires" json:"expires,omitempty"`
	Challenges       []*core.Challenge `protobuf:"bytes,6,rep,name=challenges" json:"challenges,omitempty"`
	Combinations     []byte            `protobuf:"bytes,7,opt,name=combinations" json:"combinations,omitempty"`
	XXX_unrecognized []byte            `json:"-"`
}

func (m *Authorization) Reset()                    { *m = Authorization{} }
func (m *Authorization) String() string            { return proto1.CompactTextString(m) }
func (*Authorization) ProtoMessage()               {}
func (*Authorization) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *Authorization) GetId() string {
	if m != nil && m.Id != nil {
		return *m.Id
	}
	return ""
}

func (m *Authorization) GetIdentifier() string {
	if m != nil && m.Identifier != nil {
		return *m.Identifier
	}
	return ""
}

func (m *Authorization) GetRegistrationID() int64 {
	if m != nil && m.RegistrationID != nil {
		return *m.RegistrationID
	}
	return 0
}

func (m *Authorization) GetStatus() string {
	if m != nil && m.Status != nil {
		return *m.Status
	}
	return ""
}

func (m *Authorization) GetExpires() int64 {
	if m != nil && m.Expires != nil {
		return *m.Expires
	}
	return 0
}

func (m *Authorization) GetChallenges() []*core.Challenge {
	if m != nil {
		return m.Challenges
	}
	return nil
}

func (m *Authorization) GetCombinations() []byte {
	if m != nil {
		return m.Combinations
	}
	return nil
}

type NewAuthorizationRequest struct {
	Authz            *Authorization `protobuf:"bytes,1,opt,name=authz" json:"authz,omitempty"`
	RegID            *int64         `protobuf:"varint,2,opt,name=regID" json:"regID,omitempty"`
	XXX_unrecognized []byte         `json:"-"`
}

func (m *NewAuthorizationRequest) Reset()                    { *m = NewAuthorizationRequest{} }
func (m *NewAuthorizationRequest) String() string            { return proto1.CompactTextString(m) }
func (*NewAuthorizationRequest) ProtoMessage()               {}
func (*NewAuthorizationRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *NewAuthorizationRequest) GetAuthz() *Authorization {
	if m != nil {
		return m.Authz
	}
	return nil
}

func (m *NewAuthorizationRequest) GetRegID() int64 {
	if m != nil && m.RegID != nil {
		return *m.RegID
	}
	return 0
}

type NewCertificateRequest struct {
	Csr              []byte `protobuf:"bytes,1,opt,name=csr" json:"csr,omitempty"`
	RegID            *int64 `protobuf:"varint,2,opt,name=regID" json:"regID,omitempty"`
	XXX_unrecognized []byte `json:"-"`
}

func (m *NewCertificateRequest) Reset()                    { *m = NewCertificateRequest{} }
func (m *NewCertificateRequest) String() string            { return proto1.CompactTextString(m) }
func (*NewCertificateRequest) ProtoMessage()               {}
func (*NewCertificateRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *NewCertificateRequest) GetCsr() []byte {
	if m != nil {
		return m.Csr
	}
	return nil
}

func (m *NewCertificateRequest) GetRegID() int64 {
	if m != nil && m.RegID != nil {
		return *m.RegID
	}
	return 0
}

type UpdateRegistrationRequest struct {
	Base             *Registration `protobuf:"bytes,1,opt,name=base" json:"base,omitempty"`
	Update           *Registration `protobuf:"bytes,2,opt,name=update" json:"update,omitempty"`
	XXX_unrecognized []byte        `json:"-"`
}

func (m *UpdateRegistrationRequest) Reset()                    { *m = UpdateRegistrationRequest{} }
func (m *UpdateRegistrationRequest) String() string            { return proto1.CompactTextString(m) }
func (*UpdateRegistrationRequest) ProtoMessage()               {}
func (*UpdateRegistrationRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func (m *UpdateRegistrationRequest) GetBase() *Registration {
	if m != nil {
		return m.Base
	}
	return nil
}

func (m *UpdateRegistrationRequest) GetUpdate() *Registration {
	if m != nil {
		return m.Update
	}
	return nil
}

type UpdateAuthorizationRequest struct {
	Authz            *Authorization  `protobuf:"bytes,1,opt,name=authz" json:"authz,omitempty"`
	ChallengeIndex   *int64          `protobuf:"varint,2,opt,name=challengeIndex" json:"challengeIndex,omitempty"`
	Response         *core.Challenge `protobuf:"bytes,3,opt,name=response" json:"response,omitempty"`
	XXX_unrecognized []byte          `json:"-"`
}

func (m *UpdateAuthorizationRequest) Reset()                    { *m = UpdateAuthorizationRequest{} }
func (m *UpdateAuthorizationRequest) String() string            { return proto1.CompactTextString(m) }
func (*UpdateAuthorizationRequest) ProtoMessage()               {}
func (*UpdateAuthorizationRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

func (m *UpdateAuthorizationRequest) GetAuthz() *Authorization {
	if m != nil {
		return m.Authz
	}
	return nil
}

func (m *UpdateAuthorizationRequest) GetChallengeIndex() int64 {
	if m != nil && m.ChallengeIndex != nil {
		return *m.ChallengeIndex
	}
	return 0
}

func (m *UpdateAuthorizationRequest) GetResponse() *core.Challenge {
	if m != nil {
		return m.Response
	}
	return nil
}

type RevokeCertificateWithRegRequest struct {
	Cert             []byte `protobuf:"bytes,1,opt,name=cert" json:"cert,omitempty"`
	Code             *int64 `protobuf:"varint,2,opt,name=code" json:"code,omitempty"`
	RegID            *int64 `protobuf:"varint,3,opt,name=regID" json:"regID,omitempty"`
	XXX_unrecognized []byte `json:"-"`
}

func (m *RevokeCertificateWithRegRequest) Reset()                    { *m = RevokeCertificateWithRegRequest{} }
func (m *RevokeCertificateWithRegRequest) String() string            { return proto1.CompactTextString(m) }
func (*RevokeCertificateWithRegRequest) ProtoMessage()               {}
func (*RevokeCertificateWithRegRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{6} }

func (m *RevokeCertificateWithRegRequest) GetCert() []byte {
	if m != nil {
		return m.Cert
	}
	return nil
}

func (m *RevokeCertificateWithRegRequest) GetCode() int64 {
	if m != nil && m.Code != nil {
		return *m.Code
	}
	return 0
}

func (m *RevokeCertificateWithRegRequest) GetRegID() int64 {
	if m != nil && m.RegID != nil {
		return *m.RegID
	}
	return 0
}

type AdministrativelyRevokeCertificateRequest struct {
	Cert             []byte  `protobuf:"bytes,1,opt,name=cert" json:"cert,omitempty"`
	Code             *int64  `protobuf:"varint,2,opt,name=code" json:"code,omitempty"`
	AdminName        *string `protobuf:"bytes,3,opt,name=adminName" json:"adminName,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *AdministrativelyRevokeCertificateRequest) Reset() {
	*m = AdministrativelyRevokeCertificateRequest{}
}
func (m *AdministrativelyRevokeCertificateRequest) String() string { return proto1.CompactTextString(m) }
func (*AdministrativelyRevokeCertificateRequest) ProtoMessage()    {}
func (*AdministrativelyRevokeCertificateRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor0, []int{7}
}

func (m *AdministrativelyRevokeCertificateRequest) GetCert() []byte {
	if m != nil {
		return m.Cert
	}
	return nil
}

func (m *AdministrativelyRevokeCertificateRequest) GetCode() int64 {
	if m != nil && m.Code != nil {
		return *m.Code
	}
	return 0
}

func (m *AdministrativelyRevokeCertificateRequest) GetAdminName() string {
	if m != nil && m.AdminName != nil {
		return *m.AdminName
	}
	return ""
}

func init() {
	proto1.RegisterType((*Registration)(nil), "ra.Registration")
	proto1.RegisterType((*Authorization)(nil), "ra.Authorization")
	proto1.RegisterType((*NewAuthorizationRequest)(nil), "ra.NewAuthorizationRequest")
	proto1.RegisterType((*NewCertificateRequest)(nil), "ra.NewCertificateRequest")
	proto1.RegisterType((*UpdateRegistrationRequest)(nil), "ra.UpdateRegistrationRequest")
	proto1.RegisterType((*UpdateAuthorizationRequest)(nil), "ra.UpdateAuthorizationRequest")
	proto1.RegisterType((*RevokeCertificateWithRegRequest)(nil), "ra.RevokeCertificateWithRegRequest")
	proto1.RegisterType((*AdministrativelyRevokeCertificateRequest)(nil), "ra.AdministrativelyRevokeCertificateRequest")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion3

// Client API for RegistrationAuthority service

type RegistrationAuthorityClient interface {
	NewRegistration(ctx context.Context, in *Registration, opts ...grpc.CallOption) (*Registration, error)
	NewAuthorization(ctx context.Context, in *NewAuthorizationRequest, opts ...grpc.CallOption) (*Authorization, error)
	NewCertificate(ctx context.Context, in *NewCertificateRequest, opts ...grpc.CallOption) (*core.Certificate, error)
	UpdateRegistration(ctx context.Context, in *UpdateRegistrationRequest, opts ...grpc.CallOption) (*Registration, error)
	UpdateAuthorization(ctx context.Context, in *UpdateAuthorizationRequest, opts ...grpc.CallOption) (*Authorization, error)
	RevokeCertificateWithReg(ctx context.Context, in *RevokeCertificateWithRegRequest, opts ...grpc.CallOption) (*core.Empty, error)
	DeactivateRegistration(ctx context.Context, in *Registration, opts ...grpc.CallOption) (*core.Empty, error)
	DeactivateAuthorization(ctx context.Context, in *Authorization, opts ...grpc.CallOption) (*core.Empty, error)
	AdministrativelyRevokeCertificate(ctx context.Context, in *AdministrativelyRevokeCertificateRequest, opts ...grpc.CallOption) (*core.Empty, error)
}

type registrationAuthorityClient struct {
	cc *grpc.ClientConn
}

func NewRegistrationAuthorityClient(cc *grpc.ClientConn) RegistrationAuthorityClient {
	return &registrationAuthorityClient{cc}
}

func (c *registrationAuthorityClient) NewRegistration(ctx context.Context, in *Registration, opts ...grpc.CallOption) (*Registration, error) {
	out := new(Registration)
	err := grpc.Invoke(ctx, "/ra.RegistrationAuthority/NewRegistration", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *registrationAuthorityClient) NewAuthorization(ctx context.Context, in *NewAuthorizationRequest, opts ...grpc.CallOption) (*Authorization, error) {
	out := new(Authorization)
	err := grpc.Invoke(ctx, "/ra.RegistrationAuthority/NewAuthorization", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *registrationAuthorityClient) NewCertificate(ctx context.Context, in *NewCertificateRequest, opts ...grpc.CallOption) (*core.Certificate, error) {
	out := new(core.Certificate)
	err := grpc.Invoke(ctx, "/ra.RegistrationAuthority/NewCertificate", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *registrationAuthorityClient) UpdateRegistration(ctx context.Context, in *UpdateRegistrationRequest, opts ...grpc.CallOption) (*Registration, error) {
	out := new(Registration)
	err := grpc.Invoke(ctx, "/ra.RegistrationAuthority/UpdateRegistration", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *registrationAuthorityClient) UpdateAuthorization(ctx context.Context, in *UpdateAuthorizationRequest, opts ...grpc.CallOption) (*Authorization, error) {
	out := new(Authorization)
	err := grpc.Invoke(ctx, "/ra.RegistrationAuthority/UpdateAuthorization", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *registrationAuthorityClient) RevokeCertificateWithReg(ctx context.Context, in *RevokeCertificateWithRegRequest, opts ...grpc.CallOption) (*core.Empty, error) {
	out := new(core.Empty)
	err := grpc.Invoke(ctx, "/ra.RegistrationAuthority/RevokeCertificateWithReg", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *registrationAuthorityClient) DeactivateRegistration(ctx context.Context, in *Registration, opts ...grpc.CallOption) (*core.Empty, error) {
	out := new(core.Empty)
	err := grpc.Invoke(ctx, "/ra.RegistrationAuthority/DeactivateRegistration", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *registrationAuthorityClient) DeactivateAuthorization(ctx context.Context, in *Authorization, opts ...grpc.CallOption) (*core.Empty, error) {
	out := new(core.Empty)
	err := grpc.Invoke(ctx, "/ra.RegistrationAuthority/DeactivateAuthorization", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *registrationAuthorityClient) AdministrativelyRevokeCertificate(ctx context.Context, in *AdministrativelyRevokeCertificateRequest, opts ...grpc.CallOption) (*core.Empty, error) {
	out := new(core.Empty)
	err := grpc.Invoke(ctx, "/ra.RegistrationAuthority/AdministrativelyRevokeCertificate", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for RegistrationAuthority service

type RegistrationAuthorityServer interface {
	NewRegistration(context.Context, *Registration) (*Registration, error)
	NewAuthorization(context.Context, *NewAuthorizationRequest) (*Authorization, error)
	NewCertificate(context.Context, *NewCertificateRequest) (*core.Certificate, error)
	UpdateRegistration(context.Context, *UpdateRegistrationRequest) (*Registration, error)
	UpdateAuthorization(context.Context, *UpdateAuthorizationRequest) (*Authorization, error)
	RevokeCertificateWithReg(context.Context, *RevokeCertificateWithRegRequest) (*core.Empty, error)
	DeactivateRegistration(context.Context, *Registration) (*core.Empty, error)
	DeactivateAuthorization(context.Context, *Authorization) (*core.Empty, error)
	AdministrativelyRevokeCertificate(context.Context, *AdministrativelyRevokeCertificateRequest) (*core.Empty, error)
}

func RegisterRegistrationAuthorityServer(s *grpc.Server, srv RegistrationAuthorityServer) {
	s.RegisterService(&_RegistrationAuthority_serviceDesc, srv)
}

func _RegistrationAuthority_NewRegistration_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Registration)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RegistrationAuthorityServer).NewRegistration(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ra.RegistrationAuthority/NewRegistration",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RegistrationAuthorityServer).NewRegistration(ctx, req.(*Registration))
	}
	return interceptor(ctx, in, info, handler)
}

func _RegistrationAuthority_NewAuthorization_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NewAuthorizationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RegistrationAuthorityServer).NewAuthorization(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ra.RegistrationAuthority/NewAuthorization",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RegistrationAuthorityServer).NewAuthorization(ctx, req.(*NewAuthorizationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RegistrationAuthority_NewCertificate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NewCertificateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RegistrationAuthorityServer).NewCertificate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ra.RegistrationAuthority/NewCertificate",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RegistrationAuthorityServer).NewCertificate(ctx, req.(*NewCertificateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RegistrationAuthority_UpdateRegistration_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateRegistrationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RegistrationAuthorityServer).UpdateRegistration(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ra.RegistrationAuthority/UpdateRegistration",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RegistrationAuthorityServer).UpdateRegistration(ctx, req.(*UpdateRegistrationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RegistrationAuthority_UpdateAuthorization_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateAuthorizationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RegistrationAuthorityServer).UpdateAuthorization(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ra.RegistrationAuthority/UpdateAuthorization",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RegistrationAuthorityServer).UpdateAuthorization(ctx, req.(*UpdateAuthorizationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RegistrationAuthority_RevokeCertificateWithReg_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RevokeCertificateWithRegRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RegistrationAuthorityServer).RevokeCertificateWithReg(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ra.RegistrationAuthority/RevokeCertificateWithReg",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RegistrationAuthorityServer).RevokeCertificateWithReg(ctx, req.(*RevokeCertificateWithRegRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RegistrationAuthority_DeactivateRegistration_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Registration)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RegistrationAuthorityServer).DeactivateRegistration(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ra.RegistrationAuthority/DeactivateRegistration",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RegistrationAuthorityServer).DeactivateRegistration(ctx, req.(*Registration))
	}
	return interceptor(ctx, in, info, handler)
}

func _RegistrationAuthority_DeactivateAuthorization_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Authorization)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RegistrationAuthorityServer).DeactivateAuthorization(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ra.RegistrationAuthority/DeactivateAuthorization",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RegistrationAuthorityServer).DeactivateAuthorization(ctx, req.(*Authorization))
	}
	return interceptor(ctx, in, info, handler)
}

func _RegistrationAuthority_AdministrativelyRevokeCertificate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AdministrativelyRevokeCertificateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RegistrationAuthorityServer).AdministrativelyRevokeCertificate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ra.RegistrationAuthority/AdministrativelyRevokeCertificate",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RegistrationAuthorityServer).AdministrativelyRevokeCertificate(ctx, req.(*AdministrativelyRevokeCertificateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _RegistrationAuthority_serviceDesc = grpc.ServiceDesc{
	ServiceName: "ra.RegistrationAuthority",
	HandlerType: (*RegistrationAuthorityServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "NewRegistration",
			Handler:    _RegistrationAuthority_NewRegistration_Handler,
		},
		{
			MethodName: "NewAuthorization",
			Handler:    _RegistrationAuthority_NewAuthorization_Handler,
		},
		{
			MethodName: "NewCertificate",
			Handler:    _RegistrationAuthority_NewCertificate_Handler,
		},
		{
			MethodName: "UpdateRegistration",
			Handler:    _RegistrationAuthority_UpdateRegistration_Handler,
		},
		{
			MethodName: "UpdateAuthorization",
			Handler:    _RegistrationAuthority_UpdateAuthorization_Handler,
		},
		{
			MethodName: "RevokeCertificateWithReg",
			Handler:    _RegistrationAuthority_RevokeCertificateWithReg_Handler,
		},
		{
			MethodName: "DeactivateRegistration",
			Handler:    _RegistrationAuthority_DeactivateRegistration_Handler,
		},
		{
			MethodName: "DeactivateAuthorization",
			Handler:    _RegistrationAuthority_DeactivateAuthorization_Handler,
		},
		{
			MethodName: "AdministrativelyRevokeCertificate",
			Handler:    _RegistrationAuthority_AdministrativelyRevokeCertificate_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: fileDescriptor0,
}

func init() { proto1.RegisterFile("ra/proto/ra.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 615 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x9c, 0x54, 0xc1, 0x6e, 0xd3, 0x40,
	0x14, 0x8c, 0xeb, 0x26, 0xc1, 0x2f, 0x6e, 0xd2, 0x3c, 0x68, 0xeb, 0x1a, 0x51, 0x5c, 0xf7, 0x92,
	0x03, 0x4a, 0xa5, 0x72, 0x40, 0x88, 0x0b, 0xa5, 0x41, 0x28, 0x15, 0x8a, 0x50, 0x24, 0x84, 0x40,
	0x42, 0x62, 0x6b, 0x3f, 0x92, 0x55, 0x13, 0x3b, 0xac, 0x37, 0x6d, 0xd3, 0x1b, 0x5f, 0xc1, 0x07,
	0xf0, 0xa3, 0xc8, 0x6b, 0xa7, 0xa9, 0xed, 0x04, 0x22, 0x6e, 0xeb, 0xdd, 0x9d, 0x79, 0x33, 0xb3,
	0x23, 0x43, 0x53, 0xb0, 0xe3, 0x89, 0x08, 0x65, 0x78, 0x2c, 0x58, 0x5b, 0x2d, 0x70, 0x43, 0x30,
	0x7b, 0xc7, 0x0b, 0x05, 0xa5, 0x07, 0xf1, 0x32, 0x39, 0x72, 0x7f, 0x6a, 0x60, 0xf6, 0x69, 0xc0,
	0x23, 0x29, 0x98, 0xe4, 0x61, 0x80, 0x00, 0x1b, 0xdc, 0xb7, 0x34, 0x47, 0x6b, 0xe9, 0x58, 0x03,
	0xfd, 0x92, 0x66, 0xd6, 0x86, 0xa3, 0xb5, 0x4c, 0x6c, 0x40, 0xd5, 0x0b, 0x03, 0xc9, 0x3c, 0x69,
	0xe9, 0x8e, 0xde, 0x32, 0xb0, 0x09, 0x06, 0x1b, 0x08, 0xa2, 0x31, 0x05, 0xd2, 0xda, 0x74, 0xb4,
	0x64, 0x8b, 0x07, 0x5c, 0x72, 0x36, 0xea, 0x7e, 0xb0, 0xca, 0x0a, 0xd6, 0x04, 0xc3, 0x13, 0xc4,
	0x24, 0xf9, 0xa7, 0xd2, 0xaa, 0x28, 0xda, 0x3a, 0x54, 0x22, 0xc9, 0xe4, 0x34, 0xb2, 0xaa, 0x31,
	0xca, 0xfd, 0xad, 0xc1, 0xd6, 0xe9, 0x54, 0x0e, 0x43, 0xc1, 0x6f, 0xf3, 0x22, 0x0c, 0x44, 0x00,
	0xee, 0x53, 0x20, 0xf9, 0x77, 0x4e, 0x42, 0x69, 0x31, 0x70, 0x17, 0xea, 0xe2, 0x9e, 0xe8, 0x6e,
	0xc7, 0xd2, 0x73, 0xcc, 0x89, 0x9e, 0x06, 0x54, 0xe9, 0x66, 0xc2, 0x05, 0x45, 0x4a, 0x8d, 0x8e,
	0x47, 0x00, 0xde, 0x90, 0x8d, 0x46, 0x14, 0x0c, 0x28, 0xb2, 0x2a, 0x8e, 0xde, 0xaa, 0x9d, 0x34,
	0xda, 0x2a, 0x8f, 0xb3, 0xf9, 0x3e, 0x3e, 0x02, 0xd3, 0x0b, 0xc7, 0x17, 0x3c, 0x50, 0xe4, 0x89,
	0x4a, 0xd3, 0x3d, 0x87, 0xbd, 0x1e, 0x5d, 0x67, 0x74, 0xf6, 0xe9, 0xc7, 0x94, 0x22, 0x89, 0x0e,
	0x94, 0xd9, 0x54, 0x0e, 0x6f, 0x95, 0xe2, 0xda, 0x49, 0xb3, 0x2d, 0x58, 0x3b, 0x6b, 0x68, 0x0b,
	0xca, 0x82, 0x06, 0xdd, 0x8e, 0xd2, 0xaf, 0xbb, 0xcf, 0x61, 0xa7, 0x47, 0xd7, 0x67, 0x24, 0x62,
	0x57, 0x1e, 0x93, 0x34, 0x67, 0xaa, 0x81, 0xee, 0x45, 0x42, 0xf1, 0x98, 0x79, 0xd0, 0x57, 0xd8,
	0xff, 0x38, 0xf1, 0xd5, 0xe5, 0x85, 0xf5, 0x39, 0xf0, 0x00, 0x36, 0x2f, 0x58, 0x44, 0xa9, 0x82,
	0xed, 0x58, 0x41, 0xe6, 0x59, 0x1d, 0xa8, 0x4c, 0x15, 0x58, 0x91, 0x2d, 0xb9, 0xe1, 0xce, 0xc0,
	0x4e, 0xe8, 0xff, 0xd3, 0xe2, 0x2e, 0xd4, 0xef, 0xa2, 0xed, 0x06, 0x3e, 0xdd, 0x24, 0xb2, 0xf1,
	0x10, 0x1e, 0x08, 0x8a, 0x26, 0x61, 0x10, 0x91, 0x7a, 0xa5, 0x62, 0xe0, 0xee, 0x7b, 0x78, 0xda,
	0xa7, 0xab, 0xf0, 0x92, 0xee, 0x25, 0xf2, 0x89, 0xcb, 0x61, 0x9f, 0x06, 0xf3, 0xf9, 0x26, 0x6c,
	0x7a, 0x24, 0x64, 0x9a, 0x4c, 0xfc, 0x15, 0xfa, 0x94, 0x4e, 0xb8, 0xcb, 0x49, 0x95, 0xc0, 0xfd,
	0x0c, 0xad, 0x53, 0x7f, 0xcc, 0x83, 0xd4, 0xda, 0x15, 0x8d, 0x66, 0x05, 0xf6, 0x75, 0x68, 0xe3,
	0x7e, 0xc7, 0x3c, 0x3d, 0x36, 0x4e, 0x94, 0x1b, 0x27, 0xbf, 0xca, 0xb0, 0x73, 0x3f, 0xb4, 0x34,
	0x01, 0x39, 0xc3, 0x17, 0xd0, 0xe8, 0xd1, 0x75, 0x26, 0xf2, 0x42, 0xc4, 0x76, 0x31, 0xf4, 0x12,
	0x76, 0x60, 0x3b, 0x5f, 0x2b, 0x7c, 0x1c, 0xdf, 0x5b, 0x51, 0x36, 0xbb, 0x18, 0xbd, 0x5b, 0xc2,
	0xd7, 0x50, 0xcf, 0x16, 0x0a, 0xf7, 0x53, 0x8e, 0xa2, 0x69, 0xbb, 0x99, 0xe6, 0xbf, 0x38, 0x71,
	0x4b, 0xf8, 0x0e, 0xb0, 0xd8, 0x2e, 0x7c, 0x12, 0xb3, 0xac, 0x6c, 0xdd, 0x52, 0x43, 0xe7, 0xf0,
	0x70, 0x49, 0x8f, 0xf0, 0x60, 0xc1, 0xb4, 0xbe, 0xad, 0x1e, 0x58, 0xab, 0x8a, 0x81, 0x47, 0xc9,
	0xec, 0xbf, 0xd6, 0xc6, 0xae, 0x25, 0x56, 0xdf, 0x8e, 0x27, 0x72, 0xe6, 0x96, 0xf0, 0x25, 0xec,
	0x76, 0x88, 0x79, 0x92, 0x5f, 0xe5, 0x8d, 0x16, 0x1f, 0x2b, 0x07, 0x7d, 0x05, 0x7b, 0x0b, 0x68,
	0xd6, 0x5a, 0x51, 0x7a, 0x1e, 0xfc, 0x0d, 0x0e, 0xff, 0x59, 0x49, 0x7c, 0xa6, 0x68, 0xd6, 0x6c,
	0x6e, 0x6e, 0xc2, 0x9b, 0xea, 0x97, 0xb2, 0xfa, 0xa1, 0xff, 0x09, 0x00, 0x00, 0xff, 0xff, 0x2c,
	0x26, 0x73, 0x85, 0xff, 0x05, 0x00, 0x00,
}
