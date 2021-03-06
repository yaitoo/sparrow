// Code generated by protoc-gen-go. DO NOT EDIT.
// source: pkg/eheader.proto

package pkg

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

type EHeader struct {
	// User defined
	Status uint32 `protobuf:"varint,1,opt,name=Status,proto3" json:"Status,omitempty"`
	// String format: User defined (ASCII String).
	PacketID string `protobuf:"bytes,2,opt,name=PacketID,proto3" json:"PacketID,omitempty"`
	// String format: UUID (ASCII String).
	TrackID string `protobuf:"bytes,3,opt,name=TrackID,proto3" json:"TrackID,omitempty"`
	// Payload Protobuf Length.
	PayloadLength uint32 `protobuf:"varint,4,opt,name=PayloadLength,proto3" json:"PayloadLength,omitempty"`
	// Route compression.
	Routing string `protobuf:"bytes,5,opt,name=Routing,proto3" json:"Routing,omitempty"`
	// Client Version.
	Version string `protobuf:"bytes,6,opt,name=Version,proto3" json:"Version,omitempty"`
	// Parameters
	Parameters string `protobuf:"bytes,7,opt,name=Parameters,proto3" json:"Parameters,omitempty"`
	// QueryStrings
	QueryStrings         string   `protobuf:"bytes,8,opt,name=QueryStrings,proto3" json:"QueryStrings,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *EHeader) Reset()         { *m = EHeader{} }
func (m *EHeader) String() string { return proto.CompactTextString(m) }
func (*EHeader) ProtoMessage()    {}
func (*EHeader) Descriptor() ([]byte, []int) {
	return fileDescriptor_f65f94cfcf3fa7c0, []int{0}
}

func (m *EHeader) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_EHeader.Unmarshal(m, b)
}
func (m *EHeader) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_EHeader.Marshal(b, m, deterministic)
}
func (m *EHeader) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EHeader.Merge(m, src)
}
func (m *EHeader) XXX_Size() int {
	return xxx_messageInfo_EHeader.Size(m)
}
func (m *EHeader) XXX_DiscardUnknown() {
	xxx_messageInfo_EHeader.DiscardUnknown(m)
}

var xxx_messageInfo_EHeader proto.InternalMessageInfo

func (m *EHeader) GetStatus() uint32 {
	if m != nil {
		return m.Status
	}
	return 0
}

func (m *EHeader) GetPacketID() string {
	if m != nil {
		return m.PacketID
	}
	return ""
}

func (m *EHeader) GetTrackID() string {
	if m != nil {
		return m.TrackID
	}
	return ""
}

func (m *EHeader) GetPayloadLength() uint32 {
	if m != nil {
		return m.PayloadLength
	}
	return 0
}

func (m *EHeader) GetRouting() string {
	if m != nil {
		return m.Routing
	}
	return ""
}

func (m *EHeader) GetVersion() string {
	if m != nil {
		return m.Version
	}
	return ""
}

func (m *EHeader) GetParameters() string {
	if m != nil {
		return m.Parameters
	}
	return ""
}

func (m *EHeader) GetQueryStrings() string {
	if m != nil {
		return m.QueryStrings
	}
	return ""
}

func init() {
	proto.RegisterType((*EHeader)(nil), "pkg.EHeader")
}

func init() { proto.RegisterFile("pkg/eheader.proto", fileDescriptor_f65f94cfcf3fa7c0) }

var fileDescriptor_f65f94cfcf3fa7c0 = []byte{
	// 211 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x54, 0x90, 0x31, 0x4f, 0xc3, 0x30,
	0x10, 0x46, 0x65, 0x0a, 0x49, 0x39, 0xd1, 0x01, 0x0f, 0xe8, 0xc4, 0x80, 0xaa, 0x8a, 0xa1, 0x13,
	0x0c, 0xfc, 0x85, 0x22, 0x51, 0x89, 0x21, 0xa4, 0x88, 0xfd, 0x68, 0x4f, 0x6e, 0x64, 0xb0, 0xa3,
	0xf3, 0x65, 0xe8, 0x0f, 0x67, 0x47, 0x75, 0x02, 0x6a, 0xc6, 0xf7, 0x9e, 0x3f, 0x0f, 0x07, 0xd7,
	0xad, 0x77, 0x8f, 0xbc, 0x67, 0xda, 0xb1, 0x3c, 0xb4, 0x12, 0x35, 0xda, 0x49, 0xeb, 0xdd, 0xe2,
	0xc7, 0x40, 0xf9, 0xfc, 0x92, 0xb5, 0xbd, 0x81, 0x62, 0xa3, 0xa4, 0x5d, 0x42, 0x33, 0x37, 0xcb,
	0x59, 0x3d, 0x90, 0xbd, 0x85, 0x69, 0x45, 0x5b, 0xcf, 0xba, 0x5e, 0xe1, 0xd9, 0xdc, 0x2c, 0x2f,
	0xeb, 0x7f, 0xb6, 0x08, 0xe5, 0xbb, 0xd0, 0xd6, 0xaf, 0x57, 0x38, 0xc9, 0xe9, 0x0f, 0xed, 0x3d,
	0xcc, 0x2a, 0x3a, 0x7c, 0x45, 0xda, 0xbd, 0x72, 0x70, 0xba, 0xc7, 0xf3, 0xfc, 0xe9, 0x58, 0x1e,
	0xf7, 0x75, 0xec, 0xb4, 0x09, 0x0e, 0x2f, 0xfa, 0xfd, 0x80, 0xc7, 0xf2, 0xc1, 0x92, 0x9a, 0x18,
	0xb0, 0xe8, 0xcb, 0x80, 0xf6, 0x0e, 0xa0, 0x22, 0xa1, 0x6f, 0x56, 0x96, 0x84, 0x65, 0x8e, 0x27,
	0xc6, 0x2e, 0xe0, 0xea, 0xad, 0x63, 0x39, 0x6c, 0x54, 0x9a, 0xe0, 0x12, 0x4e, 0xf3, 0x8b, 0x91,
	0xfb, 0x2c, 0xf2, 0x0d, 0x9e, 0x7e, 0x03, 0x00, 0x00, 0xff, 0xff, 0x9f, 0x60, 0x7e, 0x13, 0x18,
	0x01, 0x00, 0x00,
}
