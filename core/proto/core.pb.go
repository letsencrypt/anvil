// Code generated by protoc-gen-go. DO NOT EDIT.
// source: core/proto/core.proto

package proto

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
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

type Challenge struct {
	Id                   *int64              `protobuf:"varint,1,opt,name=id" json:"id,omitempty"`
	Type                 *string             `protobuf:"bytes,2,opt,name=type" json:"type,omitempty"`
	Status               *string             `protobuf:"bytes,6,opt,name=status" json:"status,omitempty"`
	Uri                  *string             `protobuf:"bytes,9,opt,name=uri" json:"uri,omitempty"`
	Token                *string             `protobuf:"bytes,3,opt,name=token" json:"token,omitempty"`
	KeyAuthorization     *string             `protobuf:"bytes,5,opt,name=keyAuthorization" json:"keyAuthorization,omitempty"`
	Validationrecords    []*ValidationRecord `protobuf:"bytes,10,rep,name=validationrecords" json:"validationrecords,omitempty"`
	Error                *ProblemDetails     `protobuf:"bytes,7,opt,name=error" json:"error,omitempty"`
	XXX_NoUnkeyedLiteral struct{}            `json:"-"`
	XXX_unrecognized     []byte              `json:"-"`
	XXX_sizecache        int32               `json:"-"`
}

func (m *Challenge) Reset()         { *m = Challenge{} }
func (m *Challenge) String() string { return proto.CompactTextString(m) }
func (*Challenge) ProtoMessage()    {}
func (*Challenge) Descriptor() ([]byte, []int) {
	return fileDescriptor_80ea9561f1d738ba, []int{0}
}

func (m *Challenge) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Challenge.Unmarshal(m, b)
}
func (m *Challenge) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Challenge.Marshal(b, m, deterministic)
}
func (m *Challenge) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Challenge.Merge(m, src)
}
func (m *Challenge) XXX_Size() int {
	return xxx_messageInfo_Challenge.Size(m)
}
func (m *Challenge) XXX_DiscardUnknown() {
	xxx_messageInfo_Challenge.DiscardUnknown(m)
}

var xxx_messageInfo_Challenge proto.InternalMessageInfo

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
	// A list of addresses tried before the address used (see
	// core/objects.go and the comment on the ValidationRecord structure
	// definition for more information.
	AddressesTried       [][]byte `protobuf:"bytes,7,rep,name=addressesTried" json:"addressesTried,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ValidationRecord) Reset()         { *m = ValidationRecord{} }
func (m *ValidationRecord) String() string { return proto.CompactTextString(m) }
func (*ValidationRecord) ProtoMessage()    {}
func (*ValidationRecord) Descriptor() ([]byte, []int) {
	return fileDescriptor_80ea9561f1d738ba, []int{1}
}

func (m *ValidationRecord) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ValidationRecord.Unmarshal(m, b)
}
func (m *ValidationRecord) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ValidationRecord.Marshal(b, m, deterministic)
}
func (m *ValidationRecord) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ValidationRecord.Merge(m, src)
}
func (m *ValidationRecord) XXX_Size() int {
	return xxx_messageInfo_ValidationRecord.Size(m)
}
func (m *ValidationRecord) XXX_DiscardUnknown() {
	xxx_messageInfo_ValidationRecord.DiscardUnknown(m)
}

var xxx_messageInfo_ValidationRecord proto.InternalMessageInfo

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

func (m *ValidationRecord) GetAddressesTried() [][]byte {
	if m != nil {
		return m.AddressesTried
	}
	return nil
}

type ProblemDetails struct {
	ProblemType          *string  `protobuf:"bytes,1,opt,name=problemType" json:"problemType,omitempty"`
	Detail               *string  `protobuf:"bytes,2,opt,name=detail" json:"detail,omitempty"`
	HttpStatus           *int32   `protobuf:"varint,3,opt,name=httpStatus" json:"httpStatus,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ProblemDetails) Reset()         { *m = ProblemDetails{} }
func (m *ProblemDetails) String() string { return proto.CompactTextString(m) }
func (*ProblemDetails) ProtoMessage()    {}
func (*ProblemDetails) Descriptor() ([]byte, []int) {
	return fileDescriptor_80ea9561f1d738ba, []int{2}
}

func (m *ProblemDetails) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ProblemDetails.Unmarshal(m, b)
}
func (m *ProblemDetails) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ProblemDetails.Marshal(b, m, deterministic)
}
func (m *ProblemDetails) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ProblemDetails.Merge(m, src)
}
func (m *ProblemDetails) XXX_Size() int {
	return xxx_messageInfo_ProblemDetails.Size(m)
}
func (m *ProblemDetails) XXX_DiscardUnknown() {
	xxx_messageInfo_ProblemDetails.DiscardUnknown(m)
}

var xxx_messageInfo_ProblemDetails proto.InternalMessageInfo

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
	RegistrationID       *int64   `protobuf:"varint,1,opt,name=registrationID" json:"registrationID,omitempty"`
	Serial               *string  `protobuf:"bytes,2,opt,name=serial" json:"serial,omitempty"`
	Digest               *string  `protobuf:"bytes,3,opt,name=digest" json:"digest,omitempty"`
	Der                  []byte   `protobuf:"bytes,4,opt,name=der" json:"der,omitempty"`
	Issued               *int64   `protobuf:"varint,5,opt,name=issued" json:"issued,omitempty"`
	Expires              *int64   `protobuf:"varint,6,opt,name=expires" json:"expires,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Certificate) Reset()         { *m = Certificate{} }
func (m *Certificate) String() string { return proto.CompactTextString(m) }
func (*Certificate) ProtoMessage()    {}
func (*Certificate) Descriptor() ([]byte, []int) {
	return fileDescriptor_80ea9561f1d738ba, []int{3}
}

func (m *Certificate) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Certificate.Unmarshal(m, b)
}
func (m *Certificate) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Certificate.Marshal(b, m, deterministic)
}
func (m *Certificate) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Certificate.Merge(m, src)
}
func (m *Certificate) XXX_Size() int {
	return xxx_messageInfo_Certificate.Size(m)
}
func (m *Certificate) XXX_DiscardUnknown() {
	xxx_messageInfo_Certificate.DiscardUnknown(m)
}

var xxx_messageInfo_Certificate proto.InternalMessageInfo

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
	Id                   *int64   `protobuf:"varint,1,opt,name=id" json:"id,omitempty"`
	Key                  []byte   `protobuf:"bytes,2,opt,name=key" json:"key,omitempty"`
	Contact              []string `protobuf:"bytes,3,rep,name=contact" json:"contact,omitempty"`
	ContactsPresent      *bool    `protobuf:"varint,4,opt,name=contactsPresent" json:"contactsPresent,omitempty"`
	Agreement            *string  `protobuf:"bytes,5,opt,name=agreement" json:"agreement,omitempty"`
	InitialIP            []byte   `protobuf:"bytes,6,opt,name=initialIP" json:"initialIP,omitempty"`
	CreatedAt            *int64   `protobuf:"varint,7,opt,name=createdAt" json:"createdAt,omitempty"`
	Status               *string  `protobuf:"bytes,8,opt,name=status" json:"status,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Registration) Reset()         { *m = Registration{} }
func (m *Registration) String() string { return proto.CompactTextString(m) }
func (*Registration) ProtoMessage()    {}
func (*Registration) Descriptor() ([]byte, []int) {
	return fileDescriptor_80ea9561f1d738ba, []int{4}
}

func (m *Registration) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Registration.Unmarshal(m, b)
}
func (m *Registration) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Registration.Marshal(b, m, deterministic)
}
func (m *Registration) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Registration.Merge(m, src)
}
func (m *Registration) XXX_Size() int {
	return xxx_messageInfo_Registration.Size(m)
}
func (m *Registration) XXX_DiscardUnknown() {
	xxx_messageInfo_Registration.DiscardUnknown(m)
}

var xxx_messageInfo_Registration proto.InternalMessageInfo

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
	Id                   *string      `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	Identifier           *string      `protobuf:"bytes,2,opt,name=identifier" json:"identifier,omitempty"`
	RegistrationID       *int64       `protobuf:"varint,3,opt,name=registrationID" json:"registrationID,omitempty"`
	Status               *string      `protobuf:"bytes,4,opt,name=status" json:"status,omitempty"`
	Expires              *int64       `protobuf:"varint,5,opt,name=expires" json:"expires,omitempty"`
	Challenges           []*Challenge `protobuf:"bytes,6,rep,name=challenges" json:"challenges,omitempty"`
	Combinations         []byte       `protobuf:"bytes,7,opt,name=combinations" json:"combinations,omitempty"`
	V2                   *bool        `protobuf:"varint,8,opt,name=v2" json:"v2,omitempty"`
	XXX_NoUnkeyedLiteral struct{}     `json:"-"`
	XXX_unrecognized     []byte       `json:"-"`
	XXX_sizecache        int32        `json:"-"`
}

func (m *Authorization) Reset()         { *m = Authorization{} }
func (m *Authorization) String() string { return proto.CompactTextString(m) }
func (*Authorization) ProtoMessage()    {}
func (*Authorization) Descriptor() ([]byte, []int) {
	return fileDescriptor_80ea9561f1d738ba, []int{5}
}

func (m *Authorization) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Authorization.Unmarshal(m, b)
}
func (m *Authorization) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Authorization.Marshal(b, m, deterministic)
}
func (m *Authorization) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Authorization.Merge(m, src)
}
func (m *Authorization) XXX_Size() int {
	return xxx_messageInfo_Authorization.Size(m)
}
func (m *Authorization) XXX_DiscardUnknown() {
	xxx_messageInfo_Authorization.DiscardUnknown(m)
}

var xxx_messageInfo_Authorization proto.InternalMessageInfo

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

func (m *Authorization) GetV2() bool {
	if m != nil && m.V2 != nil {
		return *m.V2
	}
	return false
}

type Order struct {
	Id                   *int64          `protobuf:"varint,1,opt,name=id" json:"id,omitempty"`
	RegistrationID       *int64          `protobuf:"varint,2,opt,name=registrationID" json:"registrationID,omitempty"`
	Expires              *int64          `protobuf:"varint,3,opt,name=expires" json:"expires,omitempty"`
	Error                *ProblemDetails `protobuf:"bytes,4,opt,name=error" json:"error,omitempty"`
	CertificateSerial    *string         `protobuf:"bytes,5,opt,name=certificateSerial" json:"certificateSerial,omitempty"`
	Authorizations       []string        `protobuf:"bytes,6,rep,name=authorizations" json:"authorizations,omitempty"`
	Status               *string         `protobuf:"bytes,7,opt,name=status" json:"status,omitempty"`
	Names                []string        `protobuf:"bytes,8,rep,name=names" json:"names,omitempty"`
	BeganProcessing      *bool           `protobuf:"varint,9,opt,name=beganProcessing" json:"beganProcessing,omitempty"`
	Created              *int64          `protobuf:"varint,10,opt,name=created" json:"created,omitempty"`
	XXX_NoUnkeyedLiteral struct{}        `json:"-"`
	XXX_unrecognized     []byte          `json:"-"`
	XXX_sizecache        int32           `json:"-"`
}

func (m *Order) Reset()         { *m = Order{} }
func (m *Order) String() string { return proto.CompactTextString(m) }
func (*Order) ProtoMessage()    {}
func (*Order) Descriptor() ([]byte, []int) {
	return fileDescriptor_80ea9561f1d738ba, []int{6}
}

func (m *Order) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Order.Unmarshal(m, b)
}
func (m *Order) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Order.Marshal(b, m, deterministic)
}
func (m *Order) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Order.Merge(m, src)
}
func (m *Order) XXX_Size() int {
	return xxx_messageInfo_Order.Size(m)
}
func (m *Order) XXX_DiscardUnknown() {
	xxx_messageInfo_Order.DiscardUnknown(m)
}

var xxx_messageInfo_Order proto.InternalMessageInfo

func (m *Order) GetId() int64 {
	if m != nil && m.Id != nil {
		return *m.Id
	}
	return 0
}

func (m *Order) GetRegistrationID() int64 {
	if m != nil && m.RegistrationID != nil {
		return *m.RegistrationID
	}
	return 0
}

func (m *Order) GetExpires() int64 {
	if m != nil && m.Expires != nil {
		return *m.Expires
	}
	return 0
}

func (m *Order) GetError() *ProblemDetails {
	if m != nil {
		return m.Error
	}
	return nil
}

func (m *Order) GetCertificateSerial() string {
	if m != nil && m.CertificateSerial != nil {
		return *m.CertificateSerial
	}
	return ""
}

func (m *Order) GetAuthorizations() []string {
	if m != nil {
		return m.Authorizations
	}
	return nil
}

func (m *Order) GetStatus() string {
	if m != nil && m.Status != nil {
		return *m.Status
	}
	return ""
}

func (m *Order) GetNames() []string {
	if m != nil {
		return m.Names
	}
	return nil
}

func (m *Order) GetBeganProcessing() bool {
	if m != nil && m.BeganProcessing != nil {
		return *m.BeganProcessing
	}
	return false
}

func (m *Order) GetCreated() int64 {
	if m != nil && m.Created != nil {
		return *m.Created
	}
	return 0
}

type Empty struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Empty) Reset()         { *m = Empty{} }
func (m *Empty) String() string { return proto.CompactTextString(m) }
func (*Empty) ProtoMessage()    {}
func (*Empty) Descriptor() ([]byte, []int) {
	return fileDescriptor_80ea9561f1d738ba, []int{7}
}

func (m *Empty) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Empty.Unmarshal(m, b)
}
func (m *Empty) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Empty.Marshal(b, m, deterministic)
}
func (m *Empty) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Empty.Merge(m, src)
}
func (m *Empty) XXX_Size() int {
	return xxx_messageInfo_Empty.Size(m)
}
func (m *Empty) XXX_DiscardUnknown() {
	xxx_messageInfo_Empty.DiscardUnknown(m)
}

var xxx_messageInfo_Empty proto.InternalMessageInfo

func init() {
	proto.RegisterType((*Challenge)(nil), "core.Challenge")
	proto.RegisterType((*ValidationRecord)(nil), "core.ValidationRecord")
	proto.RegisterType((*ProblemDetails)(nil), "core.ProblemDetails")
	proto.RegisterType((*Certificate)(nil), "core.Certificate")
	proto.RegisterType((*Registration)(nil), "core.Registration")
	proto.RegisterType((*Authorization)(nil), "core.Authorization")
	proto.RegisterType((*Order)(nil), "core.Order")
	proto.RegisterType((*Empty)(nil), "core.Empty")
}

func init() { proto.RegisterFile("core/proto/core.proto", fileDescriptor_80ea9561f1d738ba) }

var fileDescriptor_80ea9561f1d738ba = []byte{
	// 736 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x55, 0x41, 0x6e, 0xdb, 0x3a,
	0x10, 0x85, 0x2c, 0x2b, 0xb6, 0xc6, 0xfe, 0x89, 0x43, 0xe4, 0x07, 0xc2, 0xc7, 0x47, 0x20, 0x68,
	0xf1, 0x21, 0x04, 0x1f, 0x09, 0xe0, 0x1b, 0xa4, 0x49, 0x17, 0x59, 0xd5, 0x60, 0xd2, 0x2e, 0xba,
	0x53, 0xa4, 0xa9, 0xcd, 0x46, 0x16, 0x0d, 0x92, 0x36, 0xea, 0xde, 0xa1, 0x07, 0xe9, 0x61, 0x7a,
	0x8d, 0x2e, 0x7b, 0x84, 0xa2, 0xe0, 0x50, 0xb6, 0x65, 0x39, 0x45, 0x77, 0x33, 0x8f, 0x94, 0x39,
	0xf3, 0xde, 0x9b, 0x31, 0xfc, 0x9d, 0x4b, 0x85, 0xd7, 0x0b, 0x25, 0x8d, 0xbc, 0xb6, 0xe1, 0x15,
	0x85, 0xac, 0x6b, 0xe3, 0xe4, 0x4b, 0x07, 0xc2, 0xdb, 0x59, 0x56, 0x96, 0x58, 0x4d, 0x91, 0x1d,
	0x43, 0x47, 0x14, 0x91, 0x17, 0x7b, 0xa9, 0xcf, 0x3b, 0xa2, 0x60, 0x0c, 0xba, 0x66, 0xbd, 0xc0,
	0xa8, 0x13, 0x7b, 0x69, 0xc8, 0x29, 0x66, 0xe7, 0x70, 0xa4, 0x4d, 0x66, 0x96, 0x3a, 0x3a, 0x22,
	0xb4, 0xce, 0xd8, 0x08, 0xfc, 0xa5, 0x12, 0x51, 0x48, 0xa0, 0x0d, 0xd9, 0x19, 0x04, 0x46, 0x3e,
	0x63, 0x15, 0xf9, 0x84, 0xb9, 0x84, 0x5d, 0xc2, 0xe8, 0x19, 0xd7, 0x37, 0x4b, 0x33, 0x93, 0x4a,
	0x7c, 0xce, 0x8c, 0x90, 0x55, 0x14, 0xd0, 0x85, 0x03, 0x9c, 0xdd, 0xc1, 0xe9, 0x2a, 0x2b, 0x45,
	0x41, 0x99, 0xc2, 0x5c, 0xaa, 0x42, 0x47, 0x10, 0xfb, 0xe9, 0x60, 0x7c, 0x7e, 0x45, 0xbd, 0xbc,
	0xdb, 0x1e, 0x73, 0x3a, 0xe6, 0x87, 0x1f, 0xb0, 0x4b, 0x08, 0x50, 0x29, 0xa9, 0xa2, 0x5e, 0xec,
	0xa5, 0x83, 0xf1, 0x99, 0xfb, 0x72, 0xa2, 0xe4, 0x53, 0x89, 0xf3, 0x3b, 0x34, 0x99, 0x28, 0x35,
	0x77, 0x57, 0x92, 0x1f, 0x1e, 0x8c, 0xda, 0xbf, 0xc9, 0xfe, 0x81, 0xfe, 0x4c, 0x6a, 0x53, 0x65,
	0x73, 0x24, 0x72, 0x42, 0xbe, 0xcd, 0x2d, 0x45, 0x0b, 0xa9, 0xcc, 0x86, 0x22, 0x1b, 0xb3, 0xff,
	0xe1, 0x34, 0x2b, 0x0a, 0x85, 0x5a, 0xa3, 0xe6, 0xa8, 0x65, 0xb9, 0xc2, 0x22, 0xf2, 0x63, 0x3f,
	0x1d, 0xf2, 0xc3, 0x03, 0x16, 0xc3, 0xa0, 0x06, 0xdf, 0x6a, 0x2c, 0xa2, 0x6e, 0xec, 0xa5, 0x43,
	0xde, 0x84, 0xe8, 0x86, 0xe3, 0xc5, 0x08, 0xd4, 0x51, 0x10, 0xfb, 0x69, 0xc8, 0x9b, 0x90, 0x23,
	0xbf, 0xac, 0x15, 0xb1, 0x21, 0xfb, 0x0f, 0x8e, 0xb7, 0x4f, 0x3d, 0x2a, 0x81, 0x45, 0xd4, 0xa3,
	0x02, 0x5a, 0x68, 0xf2, 0x11, 0x8e, 0xf7, 0x99, 0xb0, 0xaf, 0x2d, 0x1c, 0xf2, 0x68, 0xb5, 0x77,
	0x0d, 0x37, 0x21, 0x6b, 0x81, 0x82, 0x2e, 0xd7, 0x5d, 0xd7, 0x19, 0xbb, 0x00, 0x98, 0x19, 0xb3,
	0x78, 0x70, 0xf6, 0xb0, 0xaa, 0x07, 0xbc, 0x81, 0x24, 0x5f, 0x3d, 0x18, 0xdc, 0xa2, 0x32, 0xe2,
	0x83, 0xc8, 0x33, 0x83, 0xb6, 0x46, 0x85, 0x53, 0xa1, 0x8d, 0x22, 0xb6, 0xef, 0xef, 0x6a, 0xeb,
	0xb5, 0x50, 0xb2, 0x1c, 0x2a, 0x91, 0x6d, 0xdf, 0x73, 0x19, 0xd5, 0x21, 0xa6, 0xa8, 0x4d, 0xed,
	0xb0, 0x3a, 0xb3, 0x6c, 0x14, 0xa8, 0x6a, 0x26, 0x6d, 0x68, 0x6f, 0x0a, 0xad, 0x97, 0x58, 0x90,
	0xd5, 0x7c, 0x5e, 0x67, 0x2c, 0x82, 0x1e, 0x7e, 0x5a, 0x08, 0x85, 0xce, 0xcd, 0x3e, 0xdf, 0xa4,
	0xc9, 0x77, 0x0f, 0x86, 0xbc, 0x51, 0xc6, 0xc1, 0x6c, 0x8c, 0xc0, 0x7f, 0xc6, 0x35, 0x55, 0x34,
	0xe4, 0x36, 0xb4, 0x3f, 0x96, 0xcb, 0xca, 0x64, 0xb9, 0x21, 0xb1, 0x43, 0xbe, 0x49, 0x59, 0x0a,
	0x27, 0x75, 0xa8, 0x27, 0x0a, 0x35, 0x56, 0x86, 0x8a, 0xeb, 0xf3, 0x36, 0xcc, 0xfe, 0x85, 0x30,
	0x9b, 0x2a, 0xc4, 0xb9, 0xbd, 0xe3, 0xc6, 0x62, 0x07, 0xd8, 0x53, 0x51, 0x09, 0x23, 0xb2, 0xf2,
	0x7e, 0x42, 0x05, 0x0f, 0xf9, 0x0e, 0xb0, 0xa7, 0xb9, 0xc2, 0xcc, 0x60, 0x71, 0x63, 0xc8, 0xeb,
	0x3e, 0xdf, 0x01, 0x8d, 0xb9, 0xed, 0x37, 0xe7, 0x36, 0xf9, 0xe9, 0xc1, 0x5f, 0xfb, 0x53, 0xb7,
	0xeb, 0x34, 0xa4, 0x4e, 0x2f, 0x00, 0x44, 0x81, 0x95, 0x95, 0x0d, 0x55, 0x2d, 0x41, 0x03, 0x79,
	0x41, 0x46, 0xff, 0xb7, 0x32, 0xba, 0x0a, 0xba, 0x7b, 0x9b, 0xa3, 0x21, 0x42, 0xb0, 0x27, 0x02,
	0xbb, 0x06, 0xc8, 0x37, 0xcb, 0xc9, 0x2a, 0x64, 0x07, 0xff, 0xc4, 0x8d, 0xef, 0x76, 0x69, 0xf1,
	0xc6, 0x15, 0x96, 0xc0, 0x30, 0x97, 0xf3, 0x27, 0x51, 0xd1, 0x9b, 0x9a, 0x58, 0x18, 0xf2, 0x3d,
	0xcc, 0xb6, 0xb7, 0x1a, 0x13, 0x09, 0x7d, 0xde, 0x59, 0x8d, 0x93, 0x6f, 0x1d, 0x08, 0xde, 0x28,
	0xeb, 0x92, 0xb6, 0xc4, 0x87, 0x8d, 0x75, 0x5e, 0x6c, 0xac, 0xd1, 0x80, 0xbf, 0xdf, 0xc0, 0x76,
	0xf5, 0x74, 0xff, 0xb8, 0x7a, 0xec, 0xd6, 0xc8, 0x77, 0xc3, 0xf1, 0xe0, 0x0c, 0xef, 0x2c, 0x70,
	0x78, 0x40, 0xf3, 0xdd, 0x54, 0xcd, 0xd1, 0x13, 0xf2, 0x16, 0xda, 0x20, 0xbd, 0xb7, 0x47, 0xfa,
	0x19, 0x04, 0x76, 0x7f, 0x59, 0x37, 0xd8, 0xcf, 0x5c, 0x62, 0x8d, 0xfa, 0x84, 0xd3, 0xac, 0x9a,
	0x28, 0x99, 0xa3, 0xd6, 0xa2, 0x9a, 0xd2, 0x42, 0xef, 0xf3, 0x36, 0x4c, 0x66, 0x77, 0xde, 0x8a,
	0xc0, 0xf5, 0x5c, 0xa7, 0x49, 0x0f, 0x82, 0xd7, 0xf3, 0x85, 0x59, 0xbf, 0xea, 0xbd, 0x0f, 0xe8,
	0xaf, 0xe6, 0x57, 0x00, 0x00, 0x00, 0xff, 0xff, 0xbe, 0x8c, 0x5a, 0x8e, 0x82, 0x06, 0x00, 0x00,
}
