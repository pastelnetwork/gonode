// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0-devel
// 	protoc        v3.14.0
// source: common_wn.proto

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

type SessionRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	IsPrimary bool `protobuf:"varint,1,opt,name=is_primary,json=isPrimary,proto3" json:"is_primary,omitempty"`
}

func (x *SessionRequest) Reset() {
	*x = SessionRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_common_wn_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SessionRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SessionRequest) ProtoMessage() {}

func (x *SessionRequest) ProtoReflect() protoreflect.Message {
	mi := &file_common_wn_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SessionRequest.ProtoReflect.Descriptor instead.
func (*SessionRequest) Descriptor() ([]byte, []int) {
	return file_common_wn_proto_rawDescGZIP(), []int{0}
}

func (x *SessionRequest) GetIsPrimary() bool {
	if x != nil {
		return x.IsPrimary
	}
	return false
}

type SessionReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SessID string `protobuf:"bytes,1,opt,name=sessID,proto3" json:"sessID,omitempty"`
}

func (x *SessionReply) Reset() {
	*x = SessionReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_common_wn_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SessionReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SessionReply) ProtoMessage() {}

func (x *SessionReply) ProtoReflect() protoreflect.Message {
	mi := &file_common_wn_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SessionReply.ProtoReflect.Descriptor instead.
func (*SessionReply) Descriptor() ([]byte, []int) {
	return file_common_wn_proto_rawDescGZIP(), []int{1}
}

func (x *SessionReply) GetSessID() string {
	if x != nil {
		return x.SessID
	}
	return ""
}

type AcceptedNodesRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *AcceptedNodesRequest) Reset() {
	*x = AcceptedNodesRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_common_wn_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AcceptedNodesRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AcceptedNodesRequest) ProtoMessage() {}

func (x *AcceptedNodesRequest) ProtoReflect() protoreflect.Message {
	mi := &file_common_wn_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AcceptedNodesRequest.ProtoReflect.Descriptor instead.
func (*AcceptedNodesRequest) Descriptor() ([]byte, []int) {
	return file_common_wn_proto_rawDescGZIP(), []int{2}
}

type AcceptedNodesReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Peers []*AcceptedNodesReply_Peer `protobuf:"bytes,1,rep,name=peers,proto3" json:"peers,omitempty"`
}

func (x *AcceptedNodesReply) Reset() {
	*x = AcceptedNodesReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_common_wn_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AcceptedNodesReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AcceptedNodesReply) ProtoMessage() {}

func (x *AcceptedNodesReply) ProtoReflect() protoreflect.Message {
	mi := &file_common_wn_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AcceptedNodesReply.ProtoReflect.Descriptor instead.
func (*AcceptedNodesReply) Descriptor() ([]byte, []int) {
	return file_common_wn_proto_rawDescGZIP(), []int{3}
}

func (x *AcceptedNodesReply) GetPeers() []*AcceptedNodesReply_Peer {
	if x != nil {
		return x.Peers
	}
	return nil
}

type ConnectToRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SessID string `protobuf:"bytes,1,opt,name=sessID,proto3" json:"sessID,omitempty"`
	NodeID string `protobuf:"bytes,2,opt,name=nodeID,proto3" json:"nodeID,omitempty"`
}

func (x *ConnectToRequest) Reset() {
	*x = ConnectToRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_common_wn_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ConnectToRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ConnectToRequest) ProtoMessage() {}

func (x *ConnectToRequest) ProtoReflect() protoreflect.Message {
	mi := &file_common_wn_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ConnectToRequest.ProtoReflect.Descriptor instead.
func (*ConnectToRequest) Descriptor() ([]byte, []int) {
	return file_common_wn_proto_rawDescGZIP(), []int{4}
}

func (x *ConnectToRequest) GetSessID() string {
	if x != nil {
		return x.SessID
	}
	return ""
}

func (x *ConnectToRequest) GetNodeID() string {
	if x != nil {
		return x.NodeID
	}
	return ""
}

type ConnectToReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *ConnectToReply) Reset() {
	*x = ConnectToReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_common_wn_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ConnectToReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ConnectToReply) ProtoMessage() {}

func (x *ConnectToReply) ProtoReflect() protoreflect.Message {
	mi := &file_common_wn_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ConnectToReply.ProtoReflect.Descriptor instead.
func (*ConnectToReply) Descriptor() ([]byte, []int) {
	return file_common_wn_proto_rawDescGZIP(), []int{5}
}

type MeshNodesRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Nodes []*MeshNodesRequest_Node `protobuf:"bytes,1,rep,name=nodes,proto3" json:"nodes,omitempty"`
}

func (x *MeshNodesRequest) Reset() {
	*x = MeshNodesRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_common_wn_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MeshNodesRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MeshNodesRequest) ProtoMessage() {}

func (x *MeshNodesRequest) ProtoReflect() protoreflect.Message {
	mi := &file_common_wn_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MeshNodesRequest.ProtoReflect.Descriptor instead.
func (*MeshNodesRequest) Descriptor() ([]byte, []int) {
	return file_common_wn_proto_rawDescGZIP(), []int{6}
}

func (x *MeshNodesRequest) GetNodes() []*MeshNodesRequest_Node {
	if x != nil {
		return x.Nodes
	}
	return nil
}

type MeshNodesReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *MeshNodesReply) Reset() {
	*x = MeshNodesReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_common_wn_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MeshNodesReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MeshNodesReply) ProtoMessage() {}

func (x *MeshNodesReply) ProtoReflect() protoreflect.Message {
	mi := &file_common_wn_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MeshNodesReply.ProtoReflect.Descriptor instead.
func (*MeshNodesReply) Descriptor() ([]byte, []int) {
	return file_common_wn_proto_rawDescGZIP(), []int{7}
}

type SendRegMetadataRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	CreatorPastelID string `protobuf:"bytes,1,opt,name=creatorPastelID,proto3" json:"creatorPastelID,omitempty"`
	BlockHash       string `protobuf:"bytes,2,opt,name=blockHash,proto3" json:"blockHash,omitempty"`
	BurnTxid        string `protobuf:"bytes,3,opt,name=burn_txid,json=burnTxid,proto3" json:"burn_txid,omitempty"`
}

func (x *SendRegMetadataRequest) Reset() {
	*x = SendRegMetadataRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_common_wn_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SendRegMetadataRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SendRegMetadataRequest) ProtoMessage() {}

func (x *SendRegMetadataRequest) ProtoReflect() protoreflect.Message {
	mi := &file_common_wn_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SendRegMetadataRequest.ProtoReflect.Descriptor instead.
func (*SendRegMetadataRequest) Descriptor() ([]byte, []int) {
	return file_common_wn_proto_rawDescGZIP(), []int{8}
}

func (x *SendRegMetadataRequest) GetCreatorPastelID() string {
	if x != nil {
		return x.CreatorPastelID
	}
	return ""
}

func (x *SendRegMetadataRequest) GetBlockHash() string {
	if x != nil {
		return x.BlockHash
	}
	return ""
}

func (x *SendRegMetadataRequest) GetBurnTxid() string {
	if x != nil {
		return x.BurnTxid
	}
	return ""
}

type SendRegMetadataReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *SendRegMetadataReply) Reset() {
	*x = SendRegMetadataReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_common_wn_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SendRegMetadataReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SendRegMetadataReply) ProtoMessage() {}

func (x *SendRegMetadataReply) ProtoReflect() protoreflect.Message {
	mi := &file_common_wn_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SendRegMetadataReply.ProtoReflect.Descriptor instead.
func (*SendRegMetadataReply) Descriptor() ([]byte, []int) {
	return file_common_wn_proto_rawDescGZIP(), []int{9}
}

type ProbeImageRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Payload []byte `protobuf:"bytes,1,opt,name=payload,proto3" json:"payload,omitempty"`
}

func (x *ProbeImageRequest) Reset() {
	*x = ProbeImageRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_common_wn_proto_msgTypes[10]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ProbeImageRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProbeImageRequest) ProtoMessage() {}

func (x *ProbeImageRequest) ProtoReflect() protoreflect.Message {
	mi := &file_common_wn_proto_msgTypes[10]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProbeImageRequest.ProtoReflect.Descriptor instead.
func (*ProbeImageRequest) Descriptor() ([]byte, []int) {
	return file_common_wn_proto_rawDescGZIP(), []int{10}
}

func (x *ProbeImageRequest) GetPayload() []byte {
	if x != nil {
		return x.Payload
	}
	return nil
}

type ProbeImageReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	CompressedSignedDDAndFingerprints []byte `protobuf:"bytes,1,opt,name=compressedSignedDDAndFingerprints,proto3" json:"compressedSignedDDAndFingerprints,omitempty"`
	IsValidBurnTxid                   bool   `protobuf:"varint,2,opt,name=is_valid_burn_txid,json=isValidBurnTxid,proto3" json:"is_valid_burn_txid,omitempty"`
}

func (x *ProbeImageReply) Reset() {
	*x = ProbeImageReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_common_wn_proto_msgTypes[11]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ProbeImageReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProbeImageReply) ProtoMessage() {}

func (x *ProbeImageReply) ProtoReflect() protoreflect.Message {
	mi := &file_common_wn_proto_msgTypes[11]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProbeImageReply.ProtoReflect.Descriptor instead.
func (*ProbeImageReply) Descriptor() ([]byte, []int) {
	return file_common_wn_proto_rawDescGZIP(), []int{11}
}

func (x *ProbeImageReply) GetCompressedSignedDDAndFingerprints() []byte {
	if x != nil {
		return x.CompressedSignedDDAndFingerprints
	}
	return nil
}

func (x *ProbeImageReply) GetIsValidBurnTxid() bool {
	if x != nil {
		return x.IsValidBurnTxid
	}
	return false
}

type AcceptedNodesReply_Peer struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	NodeID string `protobuf:"bytes,1,opt,name=nodeID,proto3" json:"nodeID,omitempty"`
}

func (x *AcceptedNodesReply_Peer) Reset() {
	*x = AcceptedNodesReply_Peer{}
	if protoimpl.UnsafeEnabled {
		mi := &file_common_wn_proto_msgTypes[12]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AcceptedNodesReply_Peer) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AcceptedNodesReply_Peer) ProtoMessage() {}

func (x *AcceptedNodesReply_Peer) ProtoReflect() protoreflect.Message {
	mi := &file_common_wn_proto_msgTypes[12]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AcceptedNodesReply_Peer.ProtoReflect.Descriptor instead.
func (*AcceptedNodesReply_Peer) Descriptor() ([]byte, []int) {
	return file_common_wn_proto_rawDescGZIP(), []int{3, 0}
}

func (x *AcceptedNodesReply_Peer) GetNodeID() string {
	if x != nil {
		return x.NodeID
	}
	return ""
}

type MeshNodesRequest_Node struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SessID string `protobuf:"bytes,1,opt,name=sessID,proto3" json:"sessID,omitempty"`
	NodeID string `protobuf:"bytes,2,opt,name=nodeID,proto3" json:"nodeID,omitempty"`
}

func (x *MeshNodesRequest_Node) Reset() {
	*x = MeshNodesRequest_Node{}
	if protoimpl.UnsafeEnabled {
		mi := &file_common_wn_proto_msgTypes[13]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MeshNodesRequest_Node) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MeshNodesRequest_Node) ProtoMessage() {}

func (x *MeshNodesRequest_Node) ProtoReflect() protoreflect.Message {
	mi := &file_common_wn_proto_msgTypes[13]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MeshNodesRequest_Node.ProtoReflect.Descriptor instead.
func (*MeshNodesRequest_Node) Descriptor() ([]byte, []int) {
	return file_common_wn_proto_rawDescGZIP(), []int{6, 0}
}

func (x *MeshNodesRequest_Node) GetSessID() string {
	if x != nil {
		return x.SessID
	}
	return ""
}

func (x *MeshNodesRequest_Node) GetNodeID() string {
	if x != nil {
		return x.NodeID
	}
	return ""
}

var File_common_wn_proto protoreflect.FileDescriptor

var file_common_wn_proto_rawDesc = []byte{
	0x0a, 0x0f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x5f, 0x77, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x0a, 0x77, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x6e, 0x6f, 0x64, 0x65, 0x22, 0x2f, 0x0a,
	0x0e, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x1d, 0x0a, 0x0a, 0x69, 0x73, 0x5f, 0x70, 0x72, 0x69, 0x6d, 0x61, 0x72, 0x79, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x08, 0x52, 0x09, 0x69, 0x73, 0x50, 0x72, 0x69, 0x6d, 0x61, 0x72, 0x79, 0x22, 0x26,
	0x0a, 0x0c, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x16,
	0x0a, 0x06, 0x73, 0x65, 0x73, 0x73, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06,
	0x73, 0x65, 0x73, 0x73, 0x49, 0x44, 0x22, 0x16, 0x0a, 0x14, 0x41, 0x63, 0x63, 0x65, 0x70, 0x74,
	0x65, 0x64, 0x4e, 0x6f, 0x64, 0x65, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x6f,
	0x0a, 0x12, 0x41, 0x63, 0x63, 0x65, 0x70, 0x74, 0x65, 0x64, 0x4e, 0x6f, 0x64, 0x65, 0x73, 0x52,
	0x65, 0x70, 0x6c, 0x79, 0x12, 0x39, 0x0a, 0x05, 0x70, 0x65, 0x65, 0x72, 0x73, 0x18, 0x01, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x23, 0x2e, 0x77, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x6e, 0x6f, 0x64, 0x65,
	0x2e, 0x41, 0x63, 0x63, 0x65, 0x70, 0x74, 0x65, 0x64, 0x4e, 0x6f, 0x64, 0x65, 0x73, 0x52, 0x65,
	0x70, 0x6c, 0x79, 0x2e, 0x50, 0x65, 0x65, 0x72, 0x52, 0x05, 0x70, 0x65, 0x65, 0x72, 0x73, 0x1a,
	0x1e, 0x0a, 0x04, 0x50, 0x65, 0x65, 0x72, 0x12, 0x16, 0x0a, 0x06, 0x6e, 0x6f, 0x64, 0x65, 0x49,
	0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x6e, 0x6f, 0x64, 0x65, 0x49, 0x44, 0x22,
	0x42, 0x0a, 0x10, 0x43, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x54, 0x6f, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x65, 0x73, 0x73, 0x49, 0x44, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x65, 0x73, 0x73, 0x49, 0x44, 0x12, 0x16, 0x0a, 0x06, 0x6e,
	0x6f, 0x64, 0x65, 0x49, 0x44, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x6e, 0x6f, 0x64,
	0x65, 0x49, 0x44, 0x22, 0x10, 0x0a, 0x0e, 0x43, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x54, 0x6f,
	0x52, 0x65, 0x70, 0x6c, 0x79, 0x22, 0x83, 0x01, 0x0a, 0x10, 0x4d, 0x65, 0x73, 0x68, 0x4e, 0x6f,
	0x64, 0x65, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x37, 0x0a, 0x05, 0x6e, 0x6f,
	0x64, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x21, 0x2e, 0x77, 0x61, 0x6c, 0x6c,
	0x65, 0x74, 0x6e, 0x6f, 0x64, 0x65, 0x2e, 0x4d, 0x65, 0x73, 0x68, 0x4e, 0x6f, 0x64, 0x65, 0x73,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x2e, 0x4e, 0x6f, 0x64, 0x65, 0x52, 0x05, 0x6e, 0x6f,
	0x64, 0x65, 0x73, 0x1a, 0x36, 0x0a, 0x04, 0x4e, 0x6f, 0x64, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x73,
	0x65, 0x73, 0x73, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x65, 0x73,
	0x73, 0x49, 0x44, 0x12, 0x16, 0x0a, 0x06, 0x6e, 0x6f, 0x64, 0x65, 0x49, 0x44, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x06, 0x6e, 0x6f, 0x64, 0x65, 0x49, 0x44, 0x22, 0x10, 0x0a, 0x0e, 0x4d,
	0x65, 0x73, 0x68, 0x4e, 0x6f, 0x64, 0x65, 0x73, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x22, 0x7d, 0x0a,
	0x16, 0x53, 0x65, 0x6e, 0x64, 0x52, 0x65, 0x67, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x28, 0x0a, 0x0f, 0x63, 0x72, 0x65, 0x61, 0x74,
	0x6f, 0x72, 0x50, 0x61, 0x73, 0x74, 0x65, 0x6c, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x0f, 0x63, 0x72, 0x65, 0x61, 0x74, 0x6f, 0x72, 0x50, 0x61, 0x73, 0x74, 0x65, 0x6c, 0x49,
	0x44, 0x12, 0x1c, 0x0a, 0x09, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x48, 0x61, 0x73, 0x68, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x48, 0x61, 0x73, 0x68, 0x12,
	0x1b, 0x0a, 0x09, 0x62, 0x75, 0x72, 0x6e, 0x5f, 0x74, 0x78, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x08, 0x62, 0x75, 0x72, 0x6e, 0x54, 0x78, 0x69, 0x64, 0x22, 0x16, 0x0a, 0x14,
	0x53, 0x65, 0x6e, 0x64, 0x52, 0x65, 0x67, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x52,
	0x65, 0x70, 0x6c, 0x79, 0x22, 0x2d, 0x0a, 0x11, 0x50, 0x72, 0x6f, 0x62, 0x65, 0x49, 0x6d, 0x61,
	0x67, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x18, 0x0a, 0x07, 0x70, 0x61, 0x79,
	0x6c, 0x6f, 0x61, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x07, 0x70, 0x61, 0x79, 0x6c,
	0x6f, 0x61, 0x64, 0x22, 0x8c, 0x01, 0x0a, 0x0f, 0x50, 0x72, 0x6f, 0x62, 0x65, 0x49, 0x6d, 0x61,
	0x67, 0x65, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x4c, 0x0a, 0x21, 0x63, 0x6f, 0x6d, 0x70, 0x72,
	0x65, 0x73, 0x73, 0x65, 0x64, 0x53, 0x69, 0x67, 0x6e, 0x65, 0x64, 0x44, 0x44, 0x41, 0x6e, 0x64,
	0x46, 0x69, 0x6e, 0x67, 0x65, 0x72, 0x70, 0x72, 0x69, 0x6e, 0x74, 0x73, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0c, 0x52, 0x21, 0x63, 0x6f, 0x6d, 0x70, 0x72, 0x65, 0x73, 0x73, 0x65, 0x64, 0x53, 0x69,
	0x67, 0x6e, 0x65, 0x64, 0x44, 0x44, 0x41, 0x6e, 0x64, 0x46, 0x69, 0x6e, 0x67, 0x65, 0x72, 0x70,
	0x72, 0x69, 0x6e, 0x74, 0x73, 0x12, 0x2b, 0x0a, 0x12, 0x69, 0x73, 0x5f, 0x76, 0x61, 0x6c, 0x69,
	0x64, 0x5f, 0x62, 0x75, 0x72, 0x6e, 0x5f, 0x74, 0x78, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x08, 0x52, 0x0f, 0x69, 0x73, 0x56, 0x61, 0x6c, 0x69, 0x64, 0x42, 0x75, 0x72, 0x6e, 0x54, 0x78,
	0x69, 0x64, 0x42, 0x32, 0x5a, 0x30, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d,
	0x2f, 0x70, 0x61, 0x73, 0x74, 0x65, 0x6c, 0x6e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x2f, 0x67,
	0x6f, 0x6e, 0x6f, 0x64, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x77, 0x61, 0x6c, 0x6c,
	0x65, 0x74, 0x6e, 0x6f, 0x64, 0x65, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_common_wn_proto_rawDescOnce sync.Once
	file_common_wn_proto_rawDescData = file_common_wn_proto_rawDesc
)

func file_common_wn_proto_rawDescGZIP() []byte {
	file_common_wn_proto_rawDescOnce.Do(func() {
		file_common_wn_proto_rawDescData = protoimpl.X.CompressGZIP(file_common_wn_proto_rawDescData)
	})
	return file_common_wn_proto_rawDescData
}

var file_common_wn_proto_msgTypes = make([]protoimpl.MessageInfo, 14)
var file_common_wn_proto_goTypes = []interface{}{
	(*SessionRequest)(nil),          // 0: walletnode.SessionRequest
	(*SessionReply)(nil),            // 1: walletnode.SessionReply
	(*AcceptedNodesRequest)(nil),    // 2: walletnode.AcceptedNodesRequest
	(*AcceptedNodesReply)(nil),      // 3: walletnode.AcceptedNodesReply
	(*ConnectToRequest)(nil),        // 4: walletnode.ConnectToRequest
	(*ConnectToReply)(nil),          // 5: walletnode.ConnectToReply
	(*MeshNodesRequest)(nil),        // 6: walletnode.MeshNodesRequest
	(*MeshNodesReply)(nil),          // 7: walletnode.MeshNodesReply
	(*SendRegMetadataRequest)(nil),  // 8: walletnode.SendRegMetadataRequest
	(*SendRegMetadataReply)(nil),    // 9: walletnode.SendRegMetadataReply
	(*ProbeImageRequest)(nil),       // 10: walletnode.ProbeImageRequest
	(*ProbeImageReply)(nil),         // 11: walletnode.ProbeImageReply
	(*AcceptedNodesReply_Peer)(nil), // 12: walletnode.AcceptedNodesReply.Peer
	(*MeshNodesRequest_Node)(nil),   // 13: walletnode.MeshNodesRequest.Node
}
var file_common_wn_proto_depIdxs = []int32{
	12, // 0: walletnode.AcceptedNodesReply.peers:type_name -> walletnode.AcceptedNodesReply.Peer
	13, // 1: walletnode.MeshNodesRequest.nodes:type_name -> walletnode.MeshNodesRequest.Node
	2,  // [2:2] is the sub-list for method output_type
	2,  // [2:2] is the sub-list for method input_type
	2,  // [2:2] is the sub-list for extension type_name
	2,  // [2:2] is the sub-list for extension extendee
	0,  // [0:2] is the sub-list for field type_name
}

func init() { file_common_wn_proto_init() }
func file_common_wn_proto_init() {
	if File_common_wn_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_common_wn_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SessionRequest); i {
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
		file_common_wn_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SessionReply); i {
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
		file_common_wn_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AcceptedNodesRequest); i {
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
		file_common_wn_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AcceptedNodesReply); i {
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
		file_common_wn_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ConnectToRequest); i {
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
		file_common_wn_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ConnectToReply); i {
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
		file_common_wn_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MeshNodesRequest); i {
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
		file_common_wn_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MeshNodesReply); i {
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
		file_common_wn_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SendRegMetadataRequest); i {
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
		file_common_wn_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SendRegMetadataReply); i {
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
		file_common_wn_proto_msgTypes[10].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ProbeImageRequest); i {
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
		file_common_wn_proto_msgTypes[11].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ProbeImageReply); i {
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
		file_common_wn_proto_msgTypes[12].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AcceptedNodesReply_Peer); i {
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
		file_common_wn_proto_msgTypes[13].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MeshNodesRequest_Node); i {
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
			RawDescriptor: file_common_wn_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   14,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_common_wn_proto_goTypes,
		DependencyIndexes: file_common_wn_proto_depIdxs,
		MessageInfos:      file_common_wn_proto_msgTypes,
	}.Build()
	File_common_wn_proto = out.File
	file_common_wn_proto_rawDesc = nil
	file_common_wn_proto_goTypes = nil
	file_common_wn_proto_depIdxs = nil
}
