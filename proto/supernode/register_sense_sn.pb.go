// Code generated by protoc-gen-go. DO NOT EDIT.
// source: register_sense_sn.proto

package supernode

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

func init() { proto.RegisterFile("register_sense_sn.proto", fileDescriptor_41ab8ec1337d0642) }

var fileDescriptor_41ab8ec1337d0642 = []byte{
	// 236 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x90, 0xcd, 0x4a, 0xc3, 0x40,
	0x14, 0x85, 0xa9, 0x0b, 0xc5, 0x01, 0x11, 0x66, 0x53, 0xad, 0x4b, 0x11, 0x04, 0x31, 0xe3, 0xcf,
	0x13, 0x54, 0x8a, 0x0f, 0xd0, 0xb8, 0x72, 0x53, 0xda, 0xe4, 0x38, 0x0e, 0x4d, 0xee, 0x1d, 0xe7,
	0xde, 0x20, 0x79, 0x30, 0xdf, 0x4f, 0x12, 0xda, 0x80, 0x45, 0xb1, 0xdb, 0x7b, 0xbe, 0xef, 0x1c,
	0x66, 0xcc, 0x38, 0xc1, 0x07, 0x51, 0xa4, 0x85, 0x80, 0x04, 0x0b, 0xa1, 0x2c, 0x26, 0x56, 0xb6,
	0xc7, 0xd2, 0x44, 0x24, 0xe2, 0x12, 0x93, 0xd3, 0x82, 0xeb, 0x9a, 0x69, 0xc8, 0x1e, 0xbe, 0x0e,
	0xcc, 0xc9, 0x7c, 0xe3, 0xe5, 0x9d, 0x66, 0xa7, 0xe6, 0x28, 0x87, 0x48, 0x60, 0xb2, 0xe7, 0xd9,
	0x60, 0x66, 0x9b, 0xdb, 0x1c, 0x1f, 0x0d, 0x44, 0x27, 0xe3, 0xdf, 0xa2, 0x58, 0xb5, 0xd7, 0xa3,
	0xbb, 0x91, 0x6d, 0xcd, 0x45, 0x0e, 0x2a, 0xf3, 0xe0, 0x09, 0xe5, 0x6c, 0x36, 0xa5, 0xf2, 0x39,
	0x90, 0x47, 0x8a, 0x29, 0x90, 0x8a, 0xbd, 0xfd, 0xe1, 0xfe, 0xc9, 0x6d, 0xa7, 0x6e, 0xf6, 0xc5,
	0x63, 0xd5, 0xda, 0x37, 0x73, 0xd6, 0x33, 0xdd, 0x53, 0x5e, 0x42, 0xb1, 0x86, 0x76, 0xf8, 0x52,
	0x9b, 0x04, 0x7b, 0xb5, 0x53, 0xb4, 0x93, 0x6f, 0xf7, 0x2e, 0xff, 0xc3, 0x62, 0xd5, 0x3e, 0xdd,
	0xbf, 0x3a, 0x1f, 0xf4, 0xbd, 0x59, 0x65, 0x05, 0xd7, 0x2e, 0x2e, 0x45, 0x51, 0x11, 0xf4, 0x93,
	0xd3, 0xda, 0x79, 0xee, 0x5c, 0xd7, 0x7f, 0xb0, 0x1b, 0xba, 0x56, 0x87, 0xfd, 0xe1, 0xf1, 0x3b,
	0x00, 0x00, 0xff, 0xff, 0xad, 0x84, 0xef, 0x53, 0xa8, 0x01, 0x00, 0x00,
}
