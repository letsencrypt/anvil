// Code generated by protoc-gen-go.
// source: core/proto/core.proto
// DO NOT EDIT!

/*
Package proto is a generated protocol buffer package.

It is generated from these files:
	core/proto/core.proto

It has these top-level messages:
	Challenge
	ValidationRecord
	ProblemDetails
	Certificate
	Registration
	Authorization
	CertificateRequest
	Empty
*/
package proto

import proto1 "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto1.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto1.ProtoPackageIsVersion2 // please upgrade the proto package

type Challenge struct {
	Id                *int64              `protobuf:"varint,1,opt,name=id" json:"id,omitempty"`
	Type              *string             `protobuf:"bytes,2,opt,name=type" json:"type,omitempty"`
	Status            *string             `protobuf:"bytes,6,opt,name=status" json:"status,omitempty"`
	Uri               *string             `protobuf:"bytes,9,opt,name=uri" json:"uri,omitempty"`
	Token             *string             `protobuf:"bytes,3,opt,name=token" json:"token,omitempty"`
	KeyAuthorization  *string             `protobuf:"bytes,5,opt,name=keyAuthorization" json:"keyAuthorization,omitempty"`
	Validationrecords []*ValidationRecord `protobuf:"bytes,10,rep,name=validationrecords" json:"validationrecords,omitempty"`
	Error             *ProblemDetails     `protobuf:"bytes,7,opt,name=error" json:"error,omitempty"`
	XXX_unrecognized  []byte              `json:"-"`
}

func (m *Challenge) Reset()                    { *m = Challenge{} }
func (m *Challenge) String() string            { return proto1.CompactTextString(m) }
func (*Challenge) ProtoMessage()               {}
func (*Challenge) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *Challenge) GetId() int64 {
	if m != nil && m.Id != nil {
		return *m.Id
	}
	return 0
}

func (m *Challenge) GetType() string {
	if m != nil && m.Type != nil {
		return *m.Type
	}
	return ""
}

func (m *Challenge) GetStatus() string {
	if m != nil && m.Status != nil {
		return *m.Status
	}
	return ""
}

func (m *Challenge) GetUri() string {
	if m != nil && m.Uri != nil {
		return *m.Uri
	}
	return ""
}

func (m *Challenge) GetToken() string {
	if m != nil && m.Token != nil {
		return *m.Token
	}
	return ""
}

func (m *Challenge) GetKeyAuthorization() string {
	if m != nil && m.KeyAuthorization != nil {
		return *m.KeyAuthorization
	}
	return ""
}

func (m *Challenge) GetValidationrecords() []*ValidationRecord {
	if m != nil {
		return m.Validationrecords
	}
	return nil
}

func (m *Challenge) GetError() *ProblemDetails {
	if m != nil {
		return m.Error
	}
	return nil
}

type ValidationRecord struct {
	Hostname          *string  `protobuf:"bytes,1,opt,name=hostname" json:"hostname,omitempty"`
	Port              *string  `protobuf:"bytes,2,opt,name=port" json:"port,omitempty"`
	AddressesResolved [][]byte `protobuf:"bytes,3,rep,name=addressesResolved" json:"addressesResolved,omitempty"`
	AddressUsed       []byte   `protobuf:"bytes,4,opt,name=addressUsed" json:"addressUsed,omitempty"`
	Authorities       []string `protobuf:"bytes,5,rep,name=authorities" json:"authorities,omitempty"`
	Url               *string  `protobuf:"bytes,6,opt,name=url" json:"url,omitempty"`
	XXX_unrecognized  []byte   `json:"-"`
}

func (m *ValidationRecord) Reset()                    { *m = ValidationRecord{} }
func (m *ValidationRecord) String() string            { return proto1.CompactTextString(m) }
func (*ValidationRecord) ProtoMessage()               {}
func (*ValidationRecord) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *ValidationRecord) GetHostname() string {
	if m != nil && m.Hostname != nil {
		return *m.Hostname
	}
	return ""
}

func (m *ValidationRecord) GetPort() string {
	if m != nil && m.Port != nil {
		return *m.Port
	}
	return ""
}

func (m *ValidationRecord) GetAddressesResolved() [][]byte {
	if m != nil {
		return m.AddressesResolved
	}
	return nil
}

func (m *ValidationRecord) GetAddressUsed() []byte {
	if m != nil {
		return m.AddressUsed
	}
	return nil
}

func (m *ValidationRecord) GetAuthorities() []string {
	if m != nil {
		return m.Authorities
	}
	return nil
}

func (m *ValidationRecord) GetUrl() string {
	if m != nil && m.Url != nil {
		return *m.Url
	}
	return ""
}

type ProblemDetails struct {
	ProblemType      *string `protobuf:"bytes,1,opt,name=problemType" json:"problemType,omitempty"`
	Detail           *string `protobuf:"bytes,2,opt,name=detail" json:"detail,omitempty"`
	HttpStatus       *int32  `protobuf:"varint,3,opt,name=httpStatus" json:"httpStatus,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *ProblemDetails) Reset()                    { *m = ProblemDetails{} }
func (m *ProblemDetails) String() string            { return proto1.CompactTextString(m) }
func (*ProblemDetails) ProtoMessage()               {}
func (*ProblemDetails) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *ProblemDetails) GetProblemType() string {
	if m != nil && m.ProblemType != nil {
		return *m.ProblemType
	}
	return ""
}

func (m *ProblemDetails) GetDetail() string {
	if m != nil && m.Detail != nil {
		return *m.Detail
	}
	return ""
}

func (m *ProblemDetails) GetHttpStatus() int32 {
	if m != nil && m.HttpStatus != nil {
		return *m.HttpStatus
	}
	return 0
}

type Certificate struct {
	RegistrationID   *int64  `protobuf:"varint,1,opt,name=registrationID" json:"registrationID,omitempty"`
	Serial           *string `protobuf:"bytes,2,opt,name=serial" json:"serial,omitempty"`
	Digest           *string `protobuf:"bytes,3,opt,name=digest" json:"digest,omitempty"`
	Der              []byte  `protobuf:"bytes,4,opt,name=der" json:"der,omitempty"`
	Issued           *int64  `protobuf:"varint,5,opt,name=issued" json:"issued,omitempty"`
	Expires          *int64  `protobuf:"varint,6,opt,name=expires" json:"expires,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *Certificate) Reset()                    { *m = Certificate{} }
func (m *Certificate) String() string            { return proto1.CompactTextString(m) }
func (*Certificate) ProtoMessage()               {}
func (*Certificate) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *Certificate) GetRegistrationID() int64 {
	if m != nil && m.RegistrationID != nil {
		return *m.RegistrationID
	}
	return 0
}

func (m *Certificate) GetSerial() string {
	if m != nil && m.Serial != nil {
		return *m.Serial
	}
	return ""
}

func (m *Certificate) GetDigest() string {
	if m != nil && m.Digest != nil {
		return *m.Digest
	}
	return ""
}

func (m *Certificate) GetDer() []byte {
	if m != nil {
		return m.Der
	}
	return nil
}

func (m *Certificate) GetIssued() int64 {
	if m != nil && m.Issued != nil {
		return *m.Issued
	}
	return 0
}

func (m *Certificate) GetExpires() int64 {
	if m != nil && m.Expires != nil {
		return *m.Expires
	}
	return 0
}

type Registration struct {
	Id               *int64   `protobuf:"varint,1,opt,name=id" json:"id,omitempty"`
	Key              []byte   `protobuf:"bytes,2,opt,name=key" json:"key,omitempty"`
	Contact          []string `protobuf:"bytes,3,rep,name=contact" json:"contact,omitempty"`
	ContactsPresent  *bool    `protobuf:"varint,4,opt,name=contactsPresent" json:"contactsPresent,omitempty"`
	Agreement        *string  `protobuf:"bytes,5,opt,name=agreement" json:"agreement,omitempty"`
	InitialIP        []byte   `protobuf:"bytes,6,opt,name=initialIP" json:"initialIP,omitempty"`
	CreatedAt        *int64   `protobuf:"varint,7,opt,name=createdAt" json:"createdAt,omitempty"`
	Status           *string  `protobuf:"bytes,8,opt,name=status" json:"status,omitempty"`
	XXX_unrecognized []byte   `json:"-"`
}

func (m *Registration) Reset()                    { *m = Registration{} }
func (m *Registration) String() string            { return proto1.CompactTextString(m) }
func (*Registration) ProtoMessage()               {}
func (*Registration) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

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

func (m *Registration) GetContactsPresent() bool {
	if m != nil && m.ContactsPresent != nil {
		return *m.ContactsPresent
	}
	return false
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
	Id               *string      `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	Identifier       *string      `protobuf:"bytes,2,opt,name=identifier" json:"identifier,omitempty"`
	RegistrationID   *int64       `protobuf:"varint,3,opt,name=registrationID" json:"registrationID,omitempty"`
	Status           *string      `protobuf:"bytes,4,opt,name=status" json:"status,omitempty"`
	Expires          *int64       `protobuf:"varint,5,opt,name=expires" json:"expires,omitempty"`
	Challenges       []*Challenge `protobuf:"bytes,6,rep,name=challenges" json:"challenges,omitempty"`
	Combinations     []byte       `protobuf:"bytes,7,opt,name=combinations" json:"combinations,omitempty"`
	XXX_unrecognized []byte       `json:"-"`
}

func (m *Authorization) Reset()                    { *m = Authorization{} }
func (m *Authorization) String() string            { return proto1.CompactTextString(m) }
func (*Authorization) ProtoMessage()               {}
func (*Authorization) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

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

func (m *Authorization) GetChallenges() []*Challenge {
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

type CertificateRequest struct {
	Csr              []byte `protobuf:"bytes,1,opt,name=csr" json:"csr,omitempty"`
	RegID            *int64 `protobuf:"varint,2,opt,name=regID" json:"regID,omitempty"`
	XXX_unrecognized []byte `json:"-"`
}

func (m *CertificateRequest) Reset()                    { *m = CertificateRequest{} }
func (m *CertificateRequest) String() string            { return proto1.CompactTextString(m) }
func (*CertificateRequest) ProtoMessage()               {}
func (*CertificateRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{6} }

func (m *CertificateRequest) GetCsr() []byte {
	if m != nil {
		return m.Csr
	}
	return nil
}

func (m *CertificateRequest) GetRegID() int64 {
	if m != nil && m.RegID != nil {
		return *m.RegID
	}
	return 0
}

type Empty struct {
	XXX_unrecognized []byte `json:"-"`
}

func (m *Empty) Reset()                    { *m = Empty{} }
func (m *Empty) String() string            { return proto1.CompactTextString(m) }
func (*Empty) ProtoMessage()               {}
func (*Empty) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{7} }

func init() {
	proto1.RegisterType((*Challenge)(nil), "core.Challenge")
	proto1.RegisterType((*ValidationRecord)(nil), "core.ValidationRecord")
	proto1.RegisterType((*ProblemDetails)(nil), "core.ProblemDetails")
	proto1.RegisterType((*Certificate)(nil), "core.Certificate")
	proto1.RegisterType((*Registration)(nil), "core.Registration")
	proto1.RegisterType((*Authorization)(nil), "core.Authorization")
	proto1.RegisterType((*CertificateRequest)(nil), "core.CertificateRequest")
	proto1.RegisterType((*Empty)(nil), "core.Empty")
}

func init() { proto1.RegisterFile("core/proto/core.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 543 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x64, 0x93, 0xcf, 0x72, 0xd3, 0x30,
	0x10, 0xc6, 0xc7, 0x75, 0xdc, 0xd4, 0x1b, 0xb7, 0x49, 0x4c, 0x29, 0xe6, 0xd6, 0x71, 0x2f, 0x39,
	0xb5, 0xd0, 0x37, 0x28, 0x2d, 0x87, 0xdc, 0x32, 0xe1, 0xcf, 0x81, 0x9b, 0x6a, 0x2f, 0x89, 0xa6,
	0x8e, 0x65, 0xa4, 0x4d, 0x87, 0x70, 0xe6, 0x29, 0xb8, 0xf2, 0x28, 0xbc, 0x18, 0xab, 0xb5, 0x43,
	0x13, 0x7a, 0x93, 0x76, 0x25, 0xed, 0xf7, 0xfd, 0x76, 0x05, 0x2f, 0x0b, 0x63, 0xf1, 0xaa, 0xb1,
	0x86, 0xcc, 0x95, 0x5f, 0x5e, 0xca, 0x32, 0xed, 0xf9, 0x75, 0xfe, 0x27, 0x80, 0xf8, 0x76, 0xa9,
	0xaa, 0x0a, 0xeb, 0x05, 0xa6, 0x00, 0x07, 0xba, 0xcc, 0x82, 0xf3, 0x60, 0x12, 0xa6, 0x09, 0xf4,
	0x68, 0xd3, 0x60, 0x76, 0xc0, 0xbb, 0x38, 0x3d, 0x81, 0x43, 0x47, 0x8a, 0xd6, 0x2e, 0x3b, 0x94,
	0xfd, 0x00, 0xc2, 0xb5, 0xd5, 0x59, 0x2c, 0x9b, 0x63, 0x88, 0xc8, 0x3c, 0x60, 0x9d, 0x85, 0xb2,
	0xcd, 0x60, 0xf4, 0x80, 0x9b, 0x9b, 0x35, 0x2d, 0x8d, 0xd5, 0x3f, 0x14, 0x69, 0x53, 0x67, 0x91,
	0x64, 0xde, 0xc2, 0xf8, 0x51, 0x55, 0xba, 0x94, 0x98, 0x45, 0x56, 0x50, 0xba, 0x0c, 0xce, 0xc3,
	0xc9, 0xe0, 0xfa, 0xec, 0x52, 0xb4, 0x7d, 0xfe, 0x97, 0x9e, 0x4b, 0x3a, 0xbd, 0x80, 0x08, 0xad,
	0x35, 0x36, 0xeb, 0xf3, 0x0b, 0x83, 0xeb, 0xd3, 0xf6, 0xd8, 0xcc, 0x9a, 0xfb, 0x0a, 0x57, 0x77,
	0x48, 0x4a, 0x57, 0x2e, 0xff, 0x19, 0xc0, 0xe8, 0xd9, 0xcd, 0x11, 0x1c, 0x2d, 0x8d, 0xa3, 0x5a,
	0xad, 0x50, 0x2c, 0xc5, 0xde, 0x52, 0x63, 0x2c, 0x75, 0x96, 0x5e, 0xc3, 0x58, 0x95, 0xa5, 0x45,
	0xe7, 0xd0, 0xcd, 0xd1, 0x99, 0xea, 0x11, 0x4b, 0x76, 0x10, 0x4e, 0x92, 0xf4, 0x05, 0x0c, 0xba,
	0xd4, 0x27, 0xc7, 0xc1, 0x1e, 0x9f, 0x6f, 0x83, 0xad, 0x27, 0xd2, 0xe8, 0xd8, 0x51, 0xb8, 0xe5,
	0x50, 0xb5, 0x50, 0xf2, 0x29, 0x9c, 0xec, 0x0b, 0xf3, 0x77, 0x9a, 0x36, 0xf2, 0xd1, 0xb3, 0x0c,
	0xb6, 0x2c, 0x4b, 0xc9, 0x77, 0x42, 0x18, 0xfb, 0x92, 0xa8, 0xf9, 0xd0, 0xf2, 0xf5, 0x0c, 0xa3,
	0xdc, 0xc1, 0xe0, 0x16, 0x2d, 0xe9, 0xaf, 0xba, 0x50, 0x84, 0xe9, 0x19, 0x9c, 0x58, 0x5c, 0x68,
	0x47, 0x56, 0x1c, 0x4e, 0xef, 0xba, 0x26, 0xf9, 0xb6, 0xa0, 0xd5, 0xaa, 0x7a, 0x6a, 0x53, 0xa9,
	0x17, 0xe8, 0xa8, 0x6b, 0x05, 0xcb, 0x2b, 0xd1, 0x76, 0x06, 0x38, 0xa9, 0x9d, 0x5b, 0xb3, 0xa1,
	0x48, 0x2e, 0x0f, 0xa1, 0x8f, 0xdf, 0x1b, 0xcd, 0x36, 0x45, 0x7f, 0x98, 0xff, 0x0a, 0x20, 0x99,
	0xef, 0x94, 0xd9, 0x9b, 0x07, 0x7e, 0x8a, 0xbb, 0x2a, 0x75, 0x12, 0x7f, 0xb5, 0x30, 0x35, 0xa9,
	0x82, 0x84, 0x58, 0x9c, 0xbe, 0x82, 0x61, 0x17, 0x70, 0x33, 0x7e, 0x10, 0x6b, 0x92, 0xa2, 0x47,
	0xe9, 0x18, 0x62, 0xb5, 0xb0, 0x88, 0x2b, 0x1f, 0x6a, 0xa7, 0x80, 0x43, 0xba, 0x66, 0x88, 0xaa,
	0x9a, 0xce, 0xa4, 0x72, 0xe2, 0x43, 0x85, 0x45, 0x76, 0x5a, 0xde, 0x90, 0x74, 0x3a, 0xdc, 0x99,
	0xb8, 0x23, 0x81, 0xfb, 0x3b, 0x80, 0xe3, 0xbd, 0x99, 0xda, 0x51, 0x27, 0x0c, 0x75, 0xc9, 0x25,
	0x18, 0x18, 0xfb, 0x6d, 0x61, 0x3c, 0x87, 0x16, 0xfe, 0xf7, 0x72, 0x4f, 0xce, 0xed, 0x70, 0x68,
	0xc1, 0x5c, 0x00, 0x14, 0xdb, 0x3f, 0xe1, 0xd9, 0xf8, 0xf9, 0x1c, 0xb6, 0x83, 0xf7, 0xf4, 0x57,
	0x4e, 0x21, 0x29, 0xcc, 0xea, 0x5e, 0xd7, 0xf2, 0xb8, 0x13, 0xd5, 0x49, 0xfe, 0x06, 0xd2, 0x9d,
	0xbe, 0xcd, 0xf1, 0xdb, 0x9a, 0x9b, 0xe1, 0xd9, 0x15, 0xce, 0x8a, 0xd4, 0xc4, 0xff, 0x16, 0x96,
	0xc5, 0x6a, 0x0e, 0x04, 0x7a, 0x1f, 0xa2, 0xf7, 0xab, 0x86, 0x36, 0xef, 0xfa, 0x5f, 0x22, 0xf9,
	0x99, 0x7f, 0x03, 0x00, 0x00, 0xff, 0xff, 0x1b, 0xac, 0x4f, 0x1f, 0xb1, 0x03, 0x00, 0x00,
}
