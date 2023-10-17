// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v4.24.4
// source: register_collection_wn.proto

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

type SendCollectionTicketForSignatureRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	CollectionTicket []byte `protobuf:"bytes,1,opt,name=collection_ticket,json=collectionTicket,proto3" json:"collection_ticket,omitempty"`
	CreatorSignature []byte `protobuf:"bytes,2,opt,name=creator_signature,json=creatorSignature,proto3" json:"creator_signature,omitempty"`
	BurnTxid         string `protobuf:"bytes,3,opt,name=burn_txid,json=burnTxid,proto3" json:"burn_txid,omitempty"`
}

func (x *SendCollectionTicketForSignatureRequest) Reset() {
	*x = SendCollectionTicketForSignatureRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_register_collection_wn_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SendCollectionTicketForSignatureRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SendCollectionTicketForSignatureRequest) ProtoMessage() {}

func (x *SendCollectionTicketForSignatureRequest) ProtoReflect() protoreflect.Message {
	mi := &file_register_collection_wn_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SendCollectionTicketForSignatureRequest.ProtoReflect.Descriptor instead.
func (*SendCollectionTicketForSignatureRequest) Descriptor() ([]byte, []int) {
	return file_register_collection_wn_proto_rawDescGZIP(), []int{0}
}

func (x *SendCollectionTicketForSignatureRequest) GetCollectionTicket() []byte {
	if x != nil {
		return x.CollectionTicket
	}
	return nil
}

func (x *SendCollectionTicketForSignatureRequest) GetCreatorSignature() []byte {
	if x != nil {
		return x.CreatorSignature
	}
	return nil
}

func (x *SendCollectionTicketForSignatureRequest) GetBurnTxid() string {
	if x != nil {
		return x.BurnTxid
	}
	return ""
}

type SendCollectionTicketForSignatureResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	CollectionRegTxid string `protobuf:"bytes,1,opt,name=collection_reg_txid,json=collectionRegTxid,proto3" json:"collection_reg_txid,omitempty"`
}

func (x *SendCollectionTicketForSignatureResponse) Reset() {
	*x = SendCollectionTicketForSignatureResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_register_collection_wn_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SendCollectionTicketForSignatureResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SendCollectionTicketForSignatureResponse) ProtoMessage() {}

func (x *SendCollectionTicketForSignatureResponse) ProtoReflect() protoreflect.Message {
	mi := &file_register_collection_wn_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SendCollectionTicketForSignatureResponse.ProtoReflect.Descriptor instead.
func (*SendCollectionTicketForSignatureResponse) Descriptor() ([]byte, []int) {
	return file_register_collection_wn_proto_rawDescGZIP(), []int{1}
}

func (x *SendCollectionTicketForSignatureResponse) GetCollectionRegTxid() string {
	if x != nil {
		return x.CollectionRegTxid
	}
	return ""
}

var File_register_collection_wn_proto protoreflect.FileDescriptor

var file_register_collection_wn_proto_rawDesc = []byte{
	0x0a, 0x1c, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x5f, 0x63, 0x6f, 0x6c, 0x6c, 0x65,
	0x63, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x77, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0a,
	0x77, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x6e, 0x6f, 0x64, 0x65, 0x1a, 0x0f, 0x63, 0x6f, 0x6d, 0x6d,
	0x6f, 0x6e, 0x5f, 0x77, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xa0, 0x01, 0x0a, 0x27,
	0x53, 0x65, 0x6e, 0x64, 0x43, 0x6f, 0x6c, 0x6c, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x54, 0x69,
	0x63, 0x6b, 0x65, 0x74, 0x46, 0x6f, 0x72, 0x53, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x2b, 0x0a, 0x11, 0x63, 0x6f, 0x6c, 0x6c, 0x65,
	0x63, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x74, 0x69, 0x63, 0x6b, 0x65, 0x74, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0c, 0x52, 0x10, 0x63, 0x6f, 0x6c, 0x6c, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x54, 0x69,
	0x63, 0x6b, 0x65, 0x74, 0x12, 0x2b, 0x0a, 0x11, 0x63, 0x72, 0x65, 0x61, 0x74, 0x6f, 0x72, 0x5f,
	0x73, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52,
	0x10, 0x63, 0x72, 0x65, 0x61, 0x74, 0x6f, 0x72, 0x53, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72,
	0x65, 0x12, 0x1b, 0x0a, 0x09, 0x62, 0x75, 0x72, 0x6e, 0x5f, 0x74, 0x78, 0x69, 0x64, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x62, 0x75, 0x72, 0x6e, 0x54, 0x78, 0x69, 0x64, 0x22, 0x5a,
	0x0a, 0x28, 0x53, 0x65, 0x6e, 0x64, 0x43, 0x6f, 0x6c, 0x6c, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e,
	0x54, 0x69, 0x63, 0x6b, 0x65, 0x74, 0x46, 0x6f, 0x72, 0x53, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75,
	0x72, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x2e, 0x0a, 0x13, 0x63, 0x6f,
	0x6c, 0x6c, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x72, 0x65, 0x67, 0x5f, 0x74, 0x78, 0x69,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x11, 0x63, 0x6f, 0x6c, 0x6c, 0x65, 0x63, 0x74,
	0x69, 0x6f, 0x6e, 0x52, 0x65, 0x67, 0x54, 0x78, 0x69, 0x64, 0x32, 0x91, 0x04, 0x0a, 0x12, 0x52,
	0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x43, 0x6f, 0x6c, 0x6c, 0x65, 0x63, 0x74, 0x69, 0x6f,
	0x6e, 0x12, 0x43, 0x0a, 0x07, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x1a, 0x2e, 0x77,
	0x61, 0x6c, 0x6c, 0x65, 0x74, 0x6e, 0x6f, 0x64, 0x65, 0x2e, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f,
	0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x18, 0x2e, 0x77, 0x61, 0x6c, 0x6c, 0x65,
	0x74, 0x6e, 0x6f, 0x64, 0x65, 0x2e, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x70,
	0x6c, 0x79, 0x28, 0x01, 0x30, 0x01, 0x12, 0x51, 0x0a, 0x0d, 0x41, 0x63, 0x63, 0x65, 0x70, 0x74,
	0x65, 0x64, 0x4e, 0x6f, 0x64, 0x65, 0x73, 0x12, 0x20, 0x2e, 0x77, 0x61, 0x6c, 0x6c, 0x65, 0x74,
	0x6e, 0x6f, 0x64, 0x65, 0x2e, 0x41, 0x63, 0x63, 0x65, 0x70, 0x74, 0x65, 0x64, 0x4e, 0x6f, 0x64,
	0x65, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1e, 0x2e, 0x77, 0x61, 0x6c, 0x6c,
	0x65, 0x74, 0x6e, 0x6f, 0x64, 0x65, 0x2e, 0x41, 0x63, 0x63, 0x65, 0x70, 0x74, 0x65, 0x64, 0x4e,
	0x6f, 0x64, 0x65, 0x73, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x45, 0x0a, 0x09, 0x43, 0x6f, 0x6e,
	0x6e, 0x65, 0x63, 0x74, 0x54, 0x6f, 0x12, 0x1c, 0x2e, 0x77, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x6e,
	0x6f, 0x64, 0x65, 0x2e, 0x43, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x54, 0x6f, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x1a, 0x2e, 0x77, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x6e, 0x6f, 0x64,
	0x65, 0x2e, 0x43, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x54, 0x6f, 0x52, 0x65, 0x70, 0x6c, 0x79,
	0x12, 0x45, 0x0a, 0x09, 0x4d, 0x65, 0x73, 0x68, 0x4e, 0x6f, 0x64, 0x65, 0x73, 0x12, 0x1c, 0x2e,
	0x77, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x6e, 0x6f, 0x64, 0x65, 0x2e, 0x4d, 0x65, 0x73, 0x68, 0x4e,
	0x6f, 0x64, 0x65, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1a, 0x2e, 0x77, 0x61,
	0x6c, 0x6c, 0x65, 0x74, 0x6e, 0x6f, 0x64, 0x65, 0x2e, 0x4d, 0x65, 0x73, 0x68, 0x4e, 0x6f, 0x64,
	0x65, 0x73, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x8d, 0x01, 0x0a, 0x20, 0x53, 0x65, 0x6e, 0x64,
	0x43, 0x6f, 0x6c, 0x6c, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x54, 0x69, 0x63, 0x6b, 0x65, 0x74,
	0x46, 0x6f, 0x72, 0x53, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x12, 0x33, 0x2e, 0x77,
	0x61, 0x6c, 0x6c, 0x65, 0x74, 0x6e, 0x6f, 0x64, 0x65, 0x2e, 0x53, 0x65, 0x6e, 0x64, 0x43, 0x6f,
	0x6c, 0x6c, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x54, 0x69, 0x63, 0x6b, 0x65, 0x74, 0x46, 0x6f,
	0x72, 0x53, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x34, 0x2e, 0x77, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x6e, 0x6f, 0x64, 0x65, 0x2e, 0x53,
	0x65, 0x6e, 0x64, 0x43, 0x6f, 0x6c, 0x6c, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x54, 0x69, 0x63,
	0x6b, 0x65, 0x74, 0x46, 0x6f, 0x72, 0x53, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x45, 0x0a, 0x09, 0x47, 0x65, 0x74, 0x54, 0x6f,
	0x70, 0x4d, 0x4e, 0x73, 0x12, 0x1c, 0x2e, 0x77, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x6e, 0x6f, 0x64,
	0x65, 0x2e, 0x47, 0x65, 0x74, 0x54, 0x6f, 0x70, 0x4d, 0x4e, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x1a, 0x2e, 0x77, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x6e, 0x6f, 0x64, 0x65, 0x2e,
	0x47, 0x65, 0x74, 0x54, 0x6f, 0x70, 0x4d, 0x4e, 0x73, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x42, 0x32,
	0x5a, 0x30, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x70, 0x61, 0x73,
	0x74, 0x65, 0x6c, 0x6e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x2f, 0x67, 0x6f, 0x6e, 0x6f, 0x64,
	0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x77, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x6e, 0x6f,
	0x64, 0x65, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_register_collection_wn_proto_rawDescOnce sync.Once
	file_register_collection_wn_proto_rawDescData = file_register_collection_wn_proto_rawDesc
)

func file_register_collection_wn_proto_rawDescGZIP() []byte {
	file_register_collection_wn_proto_rawDescOnce.Do(func() {
		file_register_collection_wn_proto_rawDescData = protoimpl.X.CompressGZIP(file_register_collection_wn_proto_rawDescData)
	})
	return file_register_collection_wn_proto_rawDescData
}

var file_register_collection_wn_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_register_collection_wn_proto_goTypes = []interface{}{
	(*SendCollectionTicketForSignatureRequest)(nil),  // 0: walletnode.SendCollectionTicketForSignatureRequest
	(*SendCollectionTicketForSignatureResponse)(nil), // 1: walletnode.SendCollectionTicketForSignatureResponse
	(*SessionRequest)(nil),                           // 2: walletnode.SessionRequest
	(*AcceptedNodesRequest)(nil),                     // 3: walletnode.AcceptedNodesRequest
	(*ConnectToRequest)(nil),                         // 4: walletnode.ConnectToRequest
	(*MeshNodesRequest)(nil),                         // 5: walletnode.MeshNodesRequest
	(*GetTopMNsRequest)(nil),                         // 6: walletnode.GetTopMNsRequest
	(*SessionReply)(nil),                             // 7: walletnode.SessionReply
	(*AcceptedNodesReply)(nil),                       // 8: walletnode.AcceptedNodesReply
	(*ConnectToReply)(nil),                           // 9: walletnode.ConnectToReply
	(*MeshNodesReply)(nil),                           // 10: walletnode.MeshNodesReply
	(*GetTopMNsReply)(nil),                           // 11: walletnode.GetTopMNsReply
}
var file_register_collection_wn_proto_depIdxs = []int32{
	2,  // 0: walletnode.RegisterCollection.Session:input_type -> walletnode.SessionRequest
	3,  // 1: walletnode.RegisterCollection.AcceptedNodes:input_type -> walletnode.AcceptedNodesRequest
	4,  // 2: walletnode.RegisterCollection.ConnectTo:input_type -> walletnode.ConnectToRequest
	5,  // 3: walletnode.RegisterCollection.MeshNodes:input_type -> walletnode.MeshNodesRequest
	0,  // 4: walletnode.RegisterCollection.SendCollectionTicketForSignature:input_type -> walletnode.SendCollectionTicketForSignatureRequest
	6,  // 5: walletnode.RegisterCollection.GetTopMNs:input_type -> walletnode.GetTopMNsRequest
	7,  // 6: walletnode.RegisterCollection.Session:output_type -> walletnode.SessionReply
	8,  // 7: walletnode.RegisterCollection.AcceptedNodes:output_type -> walletnode.AcceptedNodesReply
	9,  // 8: walletnode.RegisterCollection.ConnectTo:output_type -> walletnode.ConnectToReply
	10, // 9: walletnode.RegisterCollection.MeshNodes:output_type -> walletnode.MeshNodesReply
	1,  // 10: walletnode.RegisterCollection.SendCollectionTicketForSignature:output_type -> walletnode.SendCollectionTicketForSignatureResponse
	11, // 11: walletnode.RegisterCollection.GetTopMNs:output_type -> walletnode.GetTopMNsReply
	6,  // [6:12] is the sub-list for method output_type
	0,  // [0:6] is the sub-list for method input_type
	0,  // [0:0] is the sub-list for extension type_name
	0,  // [0:0] is the sub-list for extension extendee
	0,  // [0:0] is the sub-list for field type_name
}

func init() { file_register_collection_wn_proto_init() }
func file_register_collection_wn_proto_init() {
	if File_register_collection_wn_proto != nil {
		return
	}
	file_common_wn_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_register_collection_wn_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SendCollectionTicketForSignatureRequest); i {
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
		file_register_collection_wn_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SendCollectionTicketForSignatureResponse); i {
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
			RawDescriptor: file_register_collection_wn_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_register_collection_wn_proto_goTypes,
		DependencyIndexes: file_register_collection_wn_proto_depIdxs,
		MessageInfos:      file_register_collection_wn_proto_msgTypes,
	}.Build()
	File_register_collection_wn_proto = out.File
	file_register_collection_wn_proto_rawDesc = nil
	file_register_collection_wn_proto_goTypes = nil
	file_register_collection_wn_proto_depIdxs = nil
}
