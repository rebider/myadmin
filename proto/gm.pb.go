// Code generated by protoc-gen-go. DO NOT EDIT.
// source: gm.proto

/*
Package gm is a generated protocol buffer package.

It is generated from these files:
	gm.proto

It has these top-level messages:
	MSetDisableTos
	MSetDisableToc
	MSendMailTos
	MSendMailToc
*/
package gm

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type MSetDisableToc_ENUM int32

const (
	MSetDisableToc_success MSetDisableToc_ENUM = 1
	MSetDisableToc_fail    MSetDisableToc_ENUM = 2
)

var MSetDisableToc_ENUM_name = map[int32]string{
	1: "success",
	2: "fail",
}
var MSetDisableToc_ENUM_value = map[string]int32{
	"success": 1,
	"fail":    2,
}

func (x MSetDisableToc_ENUM) Enum() *MSetDisableToc_ENUM {
	p := new(MSetDisableToc_ENUM)
	*p = x
	return p
}
func (x MSetDisableToc_ENUM) String() string {
	return proto.EnumName(MSetDisableToc_ENUM_name, int32(x))
}
func (x *MSetDisableToc_ENUM) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(MSetDisableToc_ENUM_value, data, "MSetDisableToc_ENUM")
	if err != nil {
		return err
	}
	*x = MSetDisableToc_ENUM(value)
	return nil
}
func (MSetDisableToc_ENUM) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{1, 0} }

type MSendMailToc_ENUM int32

const (
	MSendMailToc_success MSendMailToc_ENUM = 1
	MSendMailToc_fail    MSendMailToc_ENUM = 2
)

var MSendMailToc_ENUM_name = map[int32]string{
	1: "success",
	2: "fail",
}
var MSendMailToc_ENUM_value = map[string]int32{
	"success": 1,
	"fail":    2,
}

func (x MSendMailToc_ENUM) Enum() *MSendMailToc_ENUM {
	p := new(MSendMailToc_ENUM)
	*p = x
	return p
}
func (x MSendMailToc_ENUM) String() string {
	return proto.EnumName(MSendMailToc_ENUM_name, int32(x))
}
func (x *MSendMailToc_ENUM) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(MSendMailToc_ENUM_value, data, "MSendMailToc_ENUM")
	if err != nil {
		return err
	}
	*x = MSendMailToc_ENUM(value)
	return nil
}
func (MSendMailToc_ENUM) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{3, 0} }

// 设置封禁
type MSetDisableTos struct {
	Token            *string `protobuf:"bytes,1,req,name=token" json:"token,omitempty"`
	Type             *int32  `protobuf:"varint,2,req,name=type" json:"type,omitempty"`
	PlayerId         *int32  `protobuf:"varint,3,req,name=player_id,json=playerId" json:"player_id,omitempty"`
	Sec              *int32  `protobuf:"varint,4,req,name=sec" json:"sec,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *MSetDisableTos) Reset()                    { *m = MSetDisableTos{} }
func (m *MSetDisableTos) String() string            { return proto.CompactTextString(m) }
func (*MSetDisableTos) ProtoMessage()               {}
func (*MSetDisableTos) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *MSetDisableTos) GetToken() string {
	if m != nil && m.Token != nil {
		return *m.Token
	}
	return ""
}

func (m *MSetDisableTos) GetType() int32 {
	if m != nil && m.Type != nil {
		return *m.Type
	}
	return 0
}

func (m *MSetDisableTos) GetPlayerId() int32 {
	if m != nil && m.PlayerId != nil {
		return *m.PlayerId
	}
	return 0
}

func (m *MSetDisableTos) GetSec() int32 {
	if m != nil && m.Sec != nil {
		return *m.Sec
	}
	return 0
}

// 设置封禁
type MSetDisableToc struct {
	Result           *MSetDisableToc_ENUM `protobuf:"varint,1,req,name=result,enum=MSetDisableToc_ENUM" json:"result,omitempty"`
	Type             *int32               `protobuf:"varint,2,req,name=type" json:"type,omitempty"`
	XXX_unrecognized []byte               `json:"-"`
}

func (m *MSetDisableToc) Reset()                    { *m = MSetDisableToc{} }
func (m *MSetDisableToc) String() string            { return proto.CompactTextString(m) }
func (*MSetDisableToc) ProtoMessage()               {}
func (*MSetDisableToc) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *MSetDisableToc) GetResult() MSetDisableToc_ENUM {
	if m != nil && m.Result != nil {
		return *m.Result
	}
	return MSetDisableToc_success
}

func (m *MSetDisableToc) GetType() int32 {
	if m != nil && m.Type != nil {
		return *m.Type
	}
	return 0
}

// 发送邮件
type MSendMailTos struct {
	Token            *string             `protobuf:"bytes,1,req,name=token" json:"token,omitempty"`
	Title            *string             `protobuf:"bytes,2,req,name=title" json:"title,omitempty"`
	Content          *string             `protobuf:"bytes,3,req,name=content" json:"content,omitempty"`
	PropList         []*MSendMailTosProp `protobuf:"bytes,4,rep,name=propList" json:"propList,omitempty"`
	PlayerNameList   *string             `protobuf:"bytes,5,req,name=playerNameList" json:"playerNameList,omitempty"`
	XXX_unrecognized []byte              `json:"-"`
}

func (m *MSendMailTos) Reset()                    { *m = MSendMailTos{} }
func (m *MSendMailTos) String() string            { return proto.CompactTextString(m) }
func (*MSendMailTos) ProtoMessage()               {}
func (*MSendMailTos) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *MSendMailTos) GetToken() string {
	if m != nil && m.Token != nil {
		return *m.Token
	}
	return ""
}

func (m *MSendMailTos) GetTitle() string {
	if m != nil && m.Title != nil {
		return *m.Title
	}
	return ""
}

func (m *MSendMailTos) GetContent() string {
	if m != nil && m.Content != nil {
		return *m.Content
	}
	return ""
}

func (m *MSendMailTos) GetPropList() []*MSendMailTosProp {
	if m != nil {
		return m.PropList
	}
	return nil
}

func (m *MSendMailTos) GetPlayerNameList() string {
	if m != nil && m.PlayerNameList != nil {
		return *m.PlayerNameList
	}
	return ""
}

type MSendMailTosProp struct {
	PropType         *int32 `protobuf:"varint,1,req,name=propType" json:"propType,omitempty"`
	PropId           *int32 `protobuf:"varint,2,req,name=propId" json:"propId,omitempty"`
	PropNum          *int32 `protobuf:"varint,3,req,name=propNum" json:"propNum,omitempty"`
	XXX_unrecognized []byte `json:"-"`
}

func (m *MSendMailTosProp) Reset()                    { *m = MSendMailTosProp{} }
func (m *MSendMailTosProp) String() string            { return proto.CompactTextString(m) }
func (*MSendMailTosProp) ProtoMessage()               {}
func (*MSendMailTosProp) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2, 0} }

func (m *MSendMailTosProp) GetPropType() int32 {
	if m != nil && m.PropType != nil {
		return *m.PropType
	}
	return 0
}

func (m *MSendMailTosProp) GetPropId() int32 {
	if m != nil && m.PropId != nil {
		return *m.PropId
	}
	return 0
}

func (m *MSendMailTosProp) GetPropNum() int32 {
	if m != nil && m.PropNum != nil {
		return *m.PropNum
	}
	return 0
}

// 发送邮件
type MSendMailToc struct {
	Result           *MSendMailToc_ENUM `protobuf:"varint,1,req,name=result,enum=MSendMailToc_ENUM" json:"result,omitempty"`
	XXX_unrecognized []byte             `json:"-"`
}

func (m *MSendMailToc) Reset()                    { *m = MSendMailToc{} }
func (m *MSendMailToc) String() string            { return proto.CompactTextString(m) }
func (*MSendMailToc) ProtoMessage()               {}
func (*MSendMailToc) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *MSendMailToc) GetResult() MSendMailToc_ENUM {
	if m != nil && m.Result != nil {
		return *m.Result
	}
	return MSendMailToc_success
}

func init() {
	proto.RegisterType((*MSetDisableTos)(nil), "m_set_disable_tos")
	proto.RegisterType((*MSetDisableToc)(nil), "m_set_disable_toc")
	proto.RegisterType((*MSendMailTos)(nil), "m_send_mail_tos")
	proto.RegisterType((*MSendMailTosProp)(nil), "m_send_mail_tos.prop")
	proto.RegisterType((*MSendMailToc)(nil), "m_send_mail_toc")
	proto.RegisterEnum("MSetDisableToc_ENUM", MSetDisableToc_ENUM_name, MSetDisableToc_ENUM_value)
	proto.RegisterEnum("MSendMailToc_ENUM", MSendMailToc_ENUM_name, MSendMailToc_ENUM_value)
}

func init() { proto.RegisterFile("gm.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 318 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x51, 0x5d, 0x4b, 0xc3, 0x30,
	0x14, 0x65, 0x5d, 0xb7, 0x75, 0x77, 0x30, 0xe7, 0xc5, 0x8f, 0x32, 0x11, 0x46, 0x1f, 0x64, 0x2f,
	0x56, 0xdc, 0x7f, 0xf0, 0x61, 0xa0, 0x7b, 0x28, 0xf3, 0xb9, 0xd4, 0x34, 0x4a, 0x30, 0x6d, 0x42,
	0x93, 0x21, 0xfb, 0x07, 0xfe, 0x6c, 0xb9, 0x69, 0x26, 0xd8, 0x89, 0xbe, 0xdd, 0x73, 0x72, 0x72,
	0x4f, 0xce, 0x09, 0x44, 0x6f, 0x55, 0xaa, 0x1b, 0x65, 0x55, 0x22, 0xe1, 0xb4, 0xca, 0x0d, 0xb7,
	0x79, 0x29, 0x4c, 0xf1, 0x22, 0x79, 0x6e, 0x95, 0xc1, 0x33, 0x18, 0x58, 0xf5, 0xce, 0xeb, 0xb8,
	0xb7, 0x08, 0x96, 0xe3, 0xac, 0x05, 0x88, 0x10, 0xda, 0xbd, 0xe6, 0x71, 0xb0, 0x08, 0x96, 0x83,
	0xcc, 0xcd, 0x78, 0x05, 0x63, 0x2d, 0x8b, 0x3d, 0x6f, 0x72, 0x51, 0xc6, 0x7d, 0x77, 0x10, 0xb5,
	0xc4, 0xba, 0xc4, 0x19, 0xf4, 0x0d, 0x67, 0x71, 0xe8, 0x68, 0x1a, 0x93, 0x8f, 0x63, 0x37, 0x86,
	0x77, 0x30, 0x6c, 0xb8, 0xd9, 0x49, 0xeb, 0xec, 0xa6, 0xab, 0xcb, 0xf4, 0x48, 0x93, 0x3e, 0x6c,
	0x9e, 0x9f, 0x32, 0x2f, 0xfb, 0xed, 0x21, 0xc9, 0x35, 0x84, 0xa4, 0xc1, 0x09, 0x8c, 0xcc, 0x8e,
	0x31, 0x6e, 0xcc, 0xac, 0x87, 0x11, 0x84, 0xaf, 0x85, 0x90, 0xb3, 0x20, 0xf9, 0x0c, 0xe0, 0x84,
	0xb6, 0xd6, 0x65, 0x5e, 0x15, 0x42, 0xfe, 0x91, 0x92, 0x58, 0x61, 0x65, 0xbb, 0x9d, 0x58, 0x02,
	0x18, 0xc3, 0x88, 0xa9, 0xda, 0xf2, 0xda, 0xba, 0x94, 0xe3, 0xec, 0x00, 0xf1, 0x1e, 0x22, 0xdd,
	0x28, 0xfd, 0x28, 0x8c, 0x8d, 0xc3, 0x45, 0x7f, 0x39, 0x59, 0x9d, 0xa7, 0x1d, 0x27, 0xaa, 0x5a,
	0x67, 0xdf, 0x32, 0xbc, 0x81, 0x69, 0xdb, 0xd1, 0xa6, 0xa8, 0xb8, 0xbb, 0x38, 0x70, 0x3b, 0x3b,
	0xec, 0x7c, 0x0b, 0x21, 0xdd, 0xc1, 0x79, 0x6b, 0xb1, 0xa5, 0xcc, 0x3d, 0xdf, 0xb1, 0xc7, 0x78,
	0x01, 0x43, 0x9a, 0xd7, 0xa5, 0x6f, 0xc3, 0x23, 0x7a, 0x30, 0x4d, 0x9b, 0x5d, 0xe5, 0xbf, 0xe5,
	0x00, 0x93, 0xbc, 0xdb, 0x04, 0xc3, 0xdb, 0xce, 0x0f, 0x74, 0x13, 0xfc, 0xec, 0xff, 0x9f, 0xae,
	0xbf, 0x02, 0x00, 0x00, 0xff, 0xff, 0xa3, 0x1c, 0x35, 0x40, 0x5d, 0x02, 0x00, 0x00,
}