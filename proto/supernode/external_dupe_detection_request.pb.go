// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.18.1
// source: external_dupe_detection_request.proto

package supernode

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

var File_external_dupe_detection_request_proto protoreflect.FileDescriptor

var file_external_dupe_detection_request_proto_rawDesc = []byte{
	0x0a, 0x25, 0x65, 0x78, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x5f, 0x64, 0x75, 0x70, 0x65, 0x5f,
	0x64, 0x65, 0x74, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x09, 0x73, 0x75, 0x70, 0x65, 0x72, 0x6e, 0x6f,
	0x64, 0x65, 0x1a, 0x16, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x5f, 0x61, 0x72, 0x74,
	0x77, 0x6f, 0x72, 0x6b, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x32, 0xc0, 0x01, 0x0a, 0x15, 0x45,
	0x78, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x44, 0x75, 0x70, 0x65, 0x44, 0x65, 0x74, 0x65, 0x63,
	0x74, 0x69, 0x6f, 0x6e, 0x12, 0x41, 0x0a, 0x07, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x12,
	0x19, 0x2e, 0x73, 0x75, 0x70, 0x65, 0x72, 0x6e, 0x6f, 0x64, 0x65, 0x2e, 0x53, 0x65, 0x73, 0x73,
	0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x17, 0x2e, 0x73, 0x75, 0x70,
	0x65, 0x72, 0x6e, 0x6f, 0x64, 0x65, 0x2e, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x52, 0x65,
	0x70, 0x6c, 0x79, 0x28, 0x01, 0x30, 0x01, 0x12, 0x64, 0x0a, 0x16, 0x53, 0x65, 0x6e, 0x64, 0x45,
	0x44, 0x44, 0x54, 0x69, 0x63, 0x6b, 0x65, 0x74, 0x53, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72,
	0x65, 0x12, 0x25, 0x2e, 0x73, 0x75, 0x70, 0x65, 0x72, 0x6e, 0x6f, 0x64, 0x65, 0x2e, 0x53, 0x65,
	0x6e, 0x64, 0x54, 0x69, 0x63, 0x6b, 0x65, 0x74, 0x53, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72,
	0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x23, 0x2e, 0x73, 0x75, 0x70, 0x65, 0x72,
	0x6e, 0x6f, 0x64, 0x65, 0x2e, 0x53, 0x65, 0x6e, 0x64, 0x54, 0x69, 0x63, 0x6b, 0x65, 0x74, 0x53,
	0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x42, 0x31, 0x5a,
	0x2f, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x70, 0x61, 0x73, 0x74,
	0x65, 0x6c, 0x6e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x2f, 0x67, 0x6f, 0x6e, 0x6f, 0x64, 0x65,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x73, 0x75, 0x70, 0x65, 0x72, 0x6e, 0x6f, 0x64, 0x65,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var file_external_dupe_detection_request_proto_goTypes = []interface{}{
	(*SessionRequest)(nil),             // 0: supernode.SessionRequest
	(*SendTicketSignatureRequest)(nil), // 1: supernode.SendTicketSignatureRequest
	(*SessionReply)(nil),               // 2: supernode.SessionReply
	(*SendTicketSignatureReply)(nil),   // 3: supernode.SendTicketSignatureReply
}
var file_external_dupe_detection_request_proto_depIdxs = []int32{
	0, // 0: supernode.ExternalDupeDetection.Session:input_type -> supernode.SessionRequest
	1, // 1: supernode.ExternalDupeDetection.SendEDDTicketSignature:input_type -> supernode.SendTicketSignatureRequest
	2, // 2: supernode.ExternalDupeDetection.Session:output_type -> supernode.SessionReply
	3, // 3: supernode.ExternalDupeDetection.SendEDDTicketSignature:output_type -> supernode.SendTicketSignatureReply
	2, // [2:4] is the sub-list for method output_type
	0, // [0:2] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_external_dupe_detection_request_proto_init() }
func file_external_dupe_detection_request_proto_init() {
	if File_external_dupe_detection_request_proto != nil {
		return
	}
	file_register_artwork_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_external_dupe_detection_request_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_external_dupe_detection_request_proto_goTypes,
		DependencyIndexes: file_external_dupe_detection_request_proto_depIdxs,
	}.Build()
	File_external_dupe_detection_request_proto = out.File
	file_external_dupe_detection_request_proto_rawDesc = nil
	file_external_dupe_detection_request_proto_goTypes = nil
	file_external_dupe_detection_request_proto_depIdxs = nil
}