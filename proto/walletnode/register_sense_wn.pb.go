// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v3.12.4
// source: register_sense_wn.proto

package walletnode

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type SendSignedSenseTicketRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ActionTicket     []byte `protobuf:"bytes,1,opt,name=action_ticket,json=actionTicket,proto3" json:"action_ticket,omitempty"`
	CreatorSignature []byte `protobuf:"bytes,2,opt,name=creator_signature,json=creatorSignature,proto3" json:"creator_signature,omitempty"`
	DdFpFiles        []byte `protobuf:"bytes,3,opt,name=dd_fp_files,json=ddFpFiles,proto3" json:"dd_fp_files,omitempty"`
}

func (x *SendSignedSenseTicketRequest) Reset() {
	*x = SendSignedSenseTicketRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_register_sense_wn_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SendSignedSenseTicketRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SendSignedSenseTicketRequest) ProtoMessage() {}

func (x *SendSignedSenseTicketRequest) ProtoReflect() protoreflect.Message {
	mi := &file_register_sense_wn_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SendSignedSenseTicketRequest.ProtoReflect.Descriptor instead.
func (*SendSignedSenseTicketRequest) Descriptor() ([]byte, []int) {
	return file_register_sense_wn_proto_rawDescGZIP(), []int{0}
}

func (x *SendSignedSenseTicketRequest) GetActionTicket() []byte {
	if x != nil {
		return x.ActionTicket
	}
	return nil
}

func (x *SendSignedSenseTicketRequest) GetCreatorSignature() []byte {
	if x != nil {
		return x.CreatorSignature
	}
	return nil
}

func (x *SendSignedSenseTicketRequest) GetDdFpFiles() []byte {
	if x != nil {
		return x.DdFpFiles
	}
	return nil
}

var File_register_sense_wn_proto protoreflect.FileDescriptor

var file_register_sense_wn_proto_rawDesc = []byte{
	0x0a, 0x17, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x5f, 0x73, 0x65, 0x6e, 0x73, 0x65,
	0x5f, 0x77, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0a, 0x77, 0x61, 0x6c, 0x6c, 0x65,
	0x74, 0x6e, 0x6f, 0x64, 0x65, 0x1a, 0x0f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x5f, 0x77, 0x6e,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x90, 0x01, 0x0a, 0x1c, 0x53, 0x65, 0x6e, 0x64, 0x53,
	0x69, 0x67, 0x6e, 0x65, 0x64, 0x53, 0x65, 0x6e, 0x73, 0x65, 0x54, 0x69, 0x63, 0x6b, 0x65, 0x74,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x23, 0x0a, 0x0d, 0x61, 0x63, 0x74, 0x69, 0x6f,
	0x6e, 0x5f, 0x74, 0x69, 0x63, 0x6b, 0x65, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0c,
	0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x54, 0x69, 0x63, 0x6b, 0x65, 0x74, 0x12, 0x2b, 0x0a, 0x11,
	0x63, 0x72, 0x65, 0x61, 0x74, 0x6f, 0x72, 0x5f, 0x73, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x10, 0x63, 0x72, 0x65, 0x61, 0x74, 0x6f, 0x72,
	0x53, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x12, 0x1e, 0x0a, 0x0b, 0x64, 0x64, 0x5f,
	0x66, 0x70, 0x5f, 0x66, 0x69, 0x6c, 0x65, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x09,
	0x64, 0x64, 0x46, 0x70, 0x46, 0x69, 0x6c, 0x65, 0x73, 0x32, 0x83, 0x07, 0x0a, 0x0d, 0x52, 0x65,
	0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x53, 0x65, 0x6e, 0x73, 0x65, 0x12, 0x43, 0x0a, 0x07, 0x53,
	0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x1a, 0x2e, 0x77, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x6e,
	0x6f, 0x64, 0x65, 0x2e, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x18, 0x2e, 0x77, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x6e, 0x6f, 0x64, 0x65, 0x2e,
	0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x28, 0x01, 0x30, 0x01,
	0x12, 0x51, 0x0a, 0x0d, 0x41, 0x63, 0x63, 0x65, 0x70, 0x74, 0x65, 0x64, 0x4e, 0x6f, 0x64, 0x65,
	0x73, 0x12, 0x20, 0x2e, 0x77, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x6e, 0x6f, 0x64, 0x65, 0x2e, 0x41,
	0x63, 0x63, 0x65, 0x70, 0x74, 0x65, 0x64, 0x4e, 0x6f, 0x64, 0x65, 0x73, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x1e, 0x2e, 0x77, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x6e, 0x6f, 0x64, 0x65,
	0x2e, 0x41, 0x63, 0x63, 0x65, 0x70, 0x74, 0x65, 0x64, 0x4e, 0x6f, 0x64, 0x65, 0x73, 0x52, 0x65,
	0x70, 0x6c, 0x79, 0x12, 0x45, 0x0a, 0x09, 0x43, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x54, 0x6f,
	0x12, 0x1c, 0x2e, 0x77, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x6e, 0x6f, 0x64, 0x65, 0x2e, 0x43, 0x6f,
	0x6e, 0x6e, 0x65, 0x63, 0x74, 0x54, 0x6f, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1a,
	0x2e, 0x77, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x6e, 0x6f, 0x64, 0x65, 0x2e, 0x43, 0x6f, 0x6e, 0x6e,
	0x65, 0x63, 0x74, 0x54, 0x6f, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x45, 0x0a, 0x09, 0x4d, 0x65,
	0x73, 0x68, 0x4e, 0x6f, 0x64, 0x65, 0x73, 0x12, 0x1c, 0x2e, 0x77, 0x61, 0x6c, 0x6c, 0x65, 0x74,
	0x6e, 0x6f, 0x64, 0x65, 0x2e, 0x4d, 0x65, 0x73, 0x68, 0x4e, 0x6f, 0x64, 0x65, 0x73, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1a, 0x2e, 0x77, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x6e, 0x6f,
	0x64, 0x65, 0x2e, 0x4d, 0x65, 0x73, 0x68, 0x4e, 0x6f, 0x64, 0x65, 0x73, 0x52, 0x65, 0x70, 0x6c,
	0x79, 0x12, 0x57, 0x0a, 0x0f, 0x53, 0x65, 0x6e, 0x64, 0x52, 0x65, 0x67, 0x4d, 0x65, 0x74, 0x61,
	0x64, 0x61, 0x74, 0x61, 0x12, 0x22, 0x2e, 0x77, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x6e, 0x6f, 0x64,
	0x65, 0x2e, 0x53, 0x65, 0x6e, 0x64, 0x52, 0x65, 0x67, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74,
	0x61, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x20, 0x2e, 0x77, 0x61, 0x6c, 0x6c, 0x65,
	0x74, 0x6e, 0x6f, 0x64, 0x65, 0x2e, 0x53, 0x65, 0x6e, 0x64, 0x52, 0x65, 0x67, 0x4d, 0x65, 0x74,
	0x61, 0x64, 0x61, 0x74, 0x61, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x4a, 0x0a, 0x0a, 0x50, 0x72,
	0x6f, 0x62, 0x65, 0x49, 0x6d, 0x61, 0x67, 0x65, 0x12, 0x1d, 0x2e, 0x77, 0x61, 0x6c, 0x6c, 0x65,
	0x74, 0x6e, 0x6f, 0x64, 0x65, 0x2e, 0x50, 0x72, 0x6f, 0x62, 0x65, 0x49, 0x6d, 0x61, 0x67, 0x65,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1b, 0x2e, 0x77, 0x61, 0x6c, 0x6c, 0x65, 0x74,
	0x6e, 0x6f, 0x64, 0x65, 0x2e, 0x50, 0x72, 0x6f, 0x62, 0x65, 0x49, 0x6d, 0x61, 0x67, 0x65, 0x52,
	0x65, 0x70, 0x6c, 0x79, 0x28, 0x01, 0x12, 0x6b, 0x0a, 0x16, 0x53, 0x65, 0x6e, 0x64, 0x53, 0x69,
	0x67, 0x6e, 0x65, 0x64, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x54, 0x69, 0x63, 0x6b, 0x65, 0x74,
	0x12, 0x28, 0x2e, 0x77, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x6e, 0x6f, 0x64, 0x65, 0x2e, 0x53, 0x65,
	0x6e, 0x64, 0x53, 0x69, 0x67, 0x6e, 0x65, 0x64, 0x53, 0x65, 0x6e, 0x73, 0x65, 0x54, 0x69, 0x63,
	0x6b, 0x65, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x27, 0x2e, 0x77, 0x61, 0x6c,
	0x6c, 0x65, 0x74, 0x6e, 0x6f, 0x64, 0x65, 0x2e, 0x53, 0x65, 0x6e, 0x64, 0x53, 0x69, 0x67, 0x6e,
	0x65, 0x64, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x54, 0x69, 0x63, 0x6b, 0x65, 0x74, 0x52, 0x65,
	0x70, 0x6c, 0x79, 0x12, 0x51, 0x0a, 0x0d, 0x53, 0x65, 0x6e, 0x64, 0x41, 0x63, 0x74, 0x69, 0x6f,
	0x6e, 0x41, 0x63, 0x74, 0x12, 0x20, 0x2e, 0x77, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x6e, 0x6f, 0x64,
	0x65, 0x2e, 0x53, 0x65, 0x6e, 0x64, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x41, 0x63, 0x74, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1e, 0x2e, 0x77, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x6e,
	0x6f, 0x64, 0x65, 0x2e, 0x53, 0x65, 0x6e, 0x64, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x41, 0x63,
	0x74, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x4a, 0x0a, 0x11, 0x47, 0x65, 0x74, 0x44, 0x44, 0x44,
	0x61, 0x74, 0x61, 0x62, 0x61, 0x73, 0x65, 0x48, 0x61, 0x73, 0x68, 0x12, 0x1c, 0x2e, 0x77, 0x61,
	0x6c, 0x6c, 0x65, 0x74, 0x6e, 0x6f, 0x64, 0x65, 0x2e, 0x47, 0x65, 0x74, 0x44, 0x42, 0x48, 0x61,
	0x73, 0x68, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x17, 0x2e, 0x77, 0x61, 0x6c, 0x6c,
	0x65, 0x74, 0x6e, 0x6f, 0x64, 0x65, 0x2e, 0x44, 0x42, 0x48, 0x61, 0x73, 0x68, 0x52, 0x65, 0x70,
	0x6c, 0x79, 0x12, 0x54, 0x0a, 0x10, 0x47, 0x65, 0x74, 0x44, 0x44, 0x53, 0x65, 0x72, 0x76, 0x65,
	0x72, 0x53, 0x74, 0x61, 0x74, 0x73, 0x12, 0x20, 0x2e, 0x77, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x6e,
	0x6f, 0x64, 0x65, 0x2e, 0x44, 0x44, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x53, 0x74, 0x61, 0x74,
	0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1e, 0x2e, 0x77, 0x61, 0x6c, 0x6c, 0x65,
	0x74, 0x6e, 0x6f, 0x64, 0x65, 0x2e, 0x44, 0x44, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x53, 0x74,
	0x61, 0x74, 0x73, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x45, 0x0a, 0x09, 0x47, 0x65, 0x74, 0x54,
	0x6f, 0x70, 0x4d, 0x4e, 0x73, 0x12, 0x1c, 0x2e, 0x77, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x6e, 0x6f,
	0x64, 0x65, 0x2e, 0x47, 0x65, 0x74, 0x54, 0x6f, 0x70, 0x4d, 0x4e, 0x73, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x1a, 0x2e, 0x77, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x6e, 0x6f, 0x64, 0x65,
	0x2e, 0x47, 0x65, 0x74, 0x54, 0x6f, 0x70, 0x4d, 0x4e, 0x73, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x42,
	0x32, 0x5a, 0x30, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x70, 0x61,
	0x73, 0x74, 0x65, 0x6c, 0x6e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x2f, 0x67, 0x6f, 0x6e, 0x6f,
	0x64, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x77, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x6e,
	0x6f, 0x64, 0x65, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_register_sense_wn_proto_rawDescOnce sync.Once
	file_register_sense_wn_proto_rawDescData = file_register_sense_wn_proto_rawDesc
)

func file_register_sense_wn_proto_rawDescGZIP() []byte {
	file_register_sense_wn_proto_rawDescOnce.Do(func() {
		file_register_sense_wn_proto_rawDescData = protoimpl.X.CompressGZIP(file_register_sense_wn_proto_rawDescData)
	})
	return file_register_sense_wn_proto_rawDescData
}

var file_register_sense_wn_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_register_sense_wn_proto_goTypes = []interface{}{
	(*SendSignedSenseTicketRequest)(nil), // 0: walletnode.SendSignedSenseTicketRequest
	(*SessionRequest)(nil),               // 1: walletnode.SessionRequest
	(*AcceptedNodesRequest)(nil),         // 2: walletnode.AcceptedNodesRequest
	(*ConnectToRequest)(nil),             // 3: walletnode.ConnectToRequest
	(*MeshNodesRequest)(nil),             // 4: walletnode.MeshNodesRequest
	(*SendRegMetadataRequest)(nil),       // 5: walletnode.SendRegMetadataRequest
	(*ProbeImageRequest)(nil),            // 6: walletnode.ProbeImageRequest
	(*SendActionActRequest)(nil),         // 7: walletnode.SendActionActRequest
	(*GetDBHashRequest)(nil),             // 8: walletnode.GetDBHashRequest
	(*DDServerStatsRequest)(nil),         // 9: walletnode.DDServerStatsRequest
	(*GetTopMNsRequest)(nil),             // 10: walletnode.GetTopMNsRequest
	(*SessionReply)(nil),                 // 11: walletnode.SessionReply
	(*AcceptedNodesReply)(nil),           // 12: walletnode.AcceptedNodesReply
	(*ConnectToReply)(nil),               // 13: walletnode.ConnectToReply
	(*MeshNodesReply)(nil),               // 14: walletnode.MeshNodesReply
	(*SendRegMetadataReply)(nil),         // 15: walletnode.SendRegMetadataReply
	(*ProbeImageReply)(nil),              // 16: walletnode.ProbeImageReply
	(*SendSignedActionTicketReply)(nil),  // 17: walletnode.SendSignedActionTicketReply
	(*SendActionActReply)(nil),           // 18: walletnode.SendActionActReply
	(*DBHashReply)(nil),                  // 19: walletnode.DBHashReply
	(*DDServerStatsReply)(nil),           // 20: walletnode.DDServerStatsReply
	(*GetTopMNsReply)(nil),               // 21: walletnode.GetTopMNsReply
}
var file_register_sense_wn_proto_depIdxs = []int32{
	1,  // 0: walletnode.RegisterSense.Session:input_type -> walletnode.SessionRequest
	2,  // 1: walletnode.RegisterSense.AcceptedNodes:input_type -> walletnode.AcceptedNodesRequest
	3,  // 2: walletnode.RegisterSense.ConnectTo:input_type -> walletnode.ConnectToRequest
	4,  // 3: walletnode.RegisterSense.MeshNodes:input_type -> walletnode.MeshNodesRequest
	5,  // 4: walletnode.RegisterSense.SendRegMetadata:input_type -> walletnode.SendRegMetadataRequest
	6,  // 5: walletnode.RegisterSense.ProbeImage:input_type -> walletnode.ProbeImageRequest
	0,  // 6: walletnode.RegisterSense.SendSignedActionTicket:input_type -> walletnode.SendSignedSenseTicketRequest
	7,  // 7: walletnode.RegisterSense.SendActionAct:input_type -> walletnode.SendActionActRequest
	8,  // 8: walletnode.RegisterSense.GetDDDatabaseHash:input_type -> walletnode.GetDBHashRequest
	9,  // 9: walletnode.RegisterSense.GetDDServerStats:input_type -> walletnode.DDServerStatsRequest
	10, // 10: walletnode.RegisterSense.GetTopMNs:input_type -> walletnode.GetTopMNsRequest
	11, // 11: walletnode.RegisterSense.Session:output_type -> walletnode.SessionReply
	12, // 12: walletnode.RegisterSense.AcceptedNodes:output_type -> walletnode.AcceptedNodesReply
	13, // 13: walletnode.RegisterSense.ConnectTo:output_type -> walletnode.ConnectToReply
	14, // 14: walletnode.RegisterSense.MeshNodes:output_type -> walletnode.MeshNodesReply
	15, // 15: walletnode.RegisterSense.SendRegMetadata:output_type -> walletnode.SendRegMetadataReply
	16, // 16: walletnode.RegisterSense.ProbeImage:output_type -> walletnode.ProbeImageReply
	17, // 17: walletnode.RegisterSense.SendSignedActionTicket:output_type -> walletnode.SendSignedActionTicketReply
	18, // 18: walletnode.RegisterSense.SendActionAct:output_type -> walletnode.SendActionActReply
	19, // 19: walletnode.RegisterSense.GetDDDatabaseHash:output_type -> walletnode.DBHashReply
	20, // 20: walletnode.RegisterSense.GetDDServerStats:output_type -> walletnode.DDServerStatsReply
	21, // 21: walletnode.RegisterSense.GetTopMNs:output_type -> walletnode.GetTopMNsReply
	11, // [11:22] is the sub-list for method output_type
	0,  // [0:11] is the sub-list for method input_type
	0,  // [0:0] is the sub-list for extension type_name
	0,  // [0:0] is the sub-list for extension extendee
	0,  // [0:0] is the sub-list for field type_name
}

func init() { file_register_sense_wn_proto_init() }
func file_register_sense_wn_proto_init() {
	if File_register_sense_wn_proto != nil {
		return
	}
	file_common_wn_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_register_sense_wn_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SendSignedSenseTicketRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_register_sense_wn_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_register_sense_wn_proto_goTypes,
		DependencyIndexes: file_register_sense_wn_proto_depIdxs,
		MessageInfos:      file_register_sense_wn_proto_msgTypes,
	}.Build()
	File_register_sense_wn_proto = out.File
	file_register_sense_wn_proto_rawDesc = nil
	file_register_sense_wn_proto_goTypes = nil
	file_register_sense_wn_proto_depIdxs = nil
}
