// Code generated by protoc-gen-go. DO NOT EDIT.
// source: register_nft_wn.proto

package walletnode

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

type SendSignedNFTTicketRequest struct {
	NftTicket            []byte             `protobuf:"bytes,1,opt,name=nft_ticket,json=nftTicket,proto3" json:"nft_ticket,omitempty"`
	CreatorSignature     []byte             `protobuf:"bytes,2,opt,name=creator_signature,json=creatorSignature,proto3" json:"creator_signature,omitempty"`
	Label                string             `protobuf:"bytes,3,opt,name=label,proto3" json:"label,omitempty"`
	EncodeParameters     *EncoderParameters `protobuf:"bytes,5,opt,name=encode_parameters,json=encodeParameters,proto3" json:"encode_parameters,omitempty"`
	DdFpFiles            []byte             `protobuf:"bytes,6,opt,name=dd_fp_files,json=ddFpFiles,proto3" json:"dd_fp_files,omitempty"`
	RqFiles              []byte             `protobuf:"bytes,7,opt,name=rq_files,json=rqFiles,proto3" json:"rq_files,omitempty"`
	XXX_NoUnkeyedLiteral struct{}           `json:"-"`
	XXX_unrecognized     []byte             `json:"-"`
	XXX_sizecache        int32              `json:"-"`
}

func (m *SendSignedNFTTicketRequest) Reset()         { *m = SendSignedNFTTicketRequest{} }
func (m *SendSignedNFTTicketRequest) String() string { return proto.CompactTextString(m) }
func (*SendSignedNFTTicketRequest) ProtoMessage()    {}
func (*SendSignedNFTTicketRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_47ba950f5e714b46, []int{0}
}

func (m *SendSignedNFTTicketRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SendSignedNFTTicketRequest.Unmarshal(m, b)
}
func (m *SendSignedNFTTicketRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SendSignedNFTTicketRequest.Marshal(b, m, deterministic)
}
func (m *SendSignedNFTTicketRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SendSignedNFTTicketRequest.Merge(m, src)
}
func (m *SendSignedNFTTicketRequest) XXX_Size() int {
	return xxx_messageInfo_SendSignedNFTTicketRequest.Size(m)
}
func (m *SendSignedNFTTicketRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_SendSignedNFTTicketRequest.DiscardUnknown(m)
}

var xxx_messageInfo_SendSignedNFTTicketRequest proto.InternalMessageInfo

func (m *SendSignedNFTTicketRequest) GetNftTicket() []byte {
	if m != nil {
		return m.NftTicket
	}
	return nil
}

func (m *SendSignedNFTTicketRequest) GetCreatorSignature() []byte {
	if m != nil {
		return m.CreatorSignature
	}
	return nil
}

func (m *SendSignedNFTTicketRequest) GetLabel() string {
	if m != nil {
		return m.Label
	}
	return ""
}

func (m *SendSignedNFTTicketRequest) GetEncodeParameters() *EncoderParameters {
	if m != nil {
		return m.EncodeParameters
	}
	return nil
}

func (m *SendSignedNFTTicketRequest) GetDdFpFiles() []byte {
	if m != nil {
		return m.DdFpFiles
	}
	return nil
}

func (m *SendSignedNFTTicketRequest) GetRqFiles() []byte {
	if m != nil {
		return m.RqFiles
	}
	return nil
}

type SendSignedNFTTicketReply struct {
	RegistrationFee      int64    `protobuf:"varint,1,opt,name=registration_fee,json=registrationFee,proto3" json:"registration_fee,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SendSignedNFTTicketReply) Reset()         { *m = SendSignedNFTTicketReply{} }
func (m *SendSignedNFTTicketReply) String() string { return proto.CompactTextString(m) }
func (*SendSignedNFTTicketReply) ProtoMessage()    {}
func (*SendSignedNFTTicketReply) Descriptor() ([]byte, []int) {
	return fileDescriptor_47ba950f5e714b46, []int{1}
}

func (m *SendSignedNFTTicketReply) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SendSignedNFTTicketReply.Unmarshal(m, b)
}
func (m *SendSignedNFTTicketReply) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SendSignedNFTTicketReply.Marshal(b, m, deterministic)
}
func (m *SendSignedNFTTicketReply) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SendSignedNFTTicketReply.Merge(m, src)
}
func (m *SendSignedNFTTicketReply) XXX_Size() int {
	return xxx_messageInfo_SendSignedNFTTicketReply.Size(m)
}
func (m *SendSignedNFTTicketReply) XXX_DiscardUnknown() {
	xxx_messageInfo_SendSignedNFTTicketReply.DiscardUnknown(m)
}

var xxx_messageInfo_SendSignedNFTTicketReply proto.InternalMessageInfo

func (m *SendSignedNFTTicketReply) GetRegistrationFee() int64 {
	if m != nil {
		return m.RegistrationFee
	}
	return 0
}

type SendPreBurntFeeTxidRequest struct {
	Txid                 string   `protobuf:"bytes,1,opt,name=txid,proto3" json:"txid,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SendPreBurntFeeTxidRequest) Reset()         { *m = SendPreBurntFeeTxidRequest{} }
func (m *SendPreBurntFeeTxidRequest) String() string { return proto.CompactTextString(m) }
func (*SendPreBurntFeeTxidRequest) ProtoMessage()    {}
func (*SendPreBurntFeeTxidRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_47ba950f5e714b46, []int{2}
}

func (m *SendPreBurntFeeTxidRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SendPreBurntFeeTxidRequest.Unmarshal(m, b)
}
func (m *SendPreBurntFeeTxidRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SendPreBurntFeeTxidRequest.Marshal(b, m, deterministic)
}
func (m *SendPreBurntFeeTxidRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SendPreBurntFeeTxidRequest.Merge(m, src)
}
func (m *SendPreBurntFeeTxidRequest) XXX_Size() int {
	return xxx_messageInfo_SendPreBurntFeeTxidRequest.Size(m)
}
func (m *SendPreBurntFeeTxidRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_SendPreBurntFeeTxidRequest.DiscardUnknown(m)
}

var xxx_messageInfo_SendPreBurntFeeTxidRequest proto.InternalMessageInfo

func (m *SendPreBurntFeeTxidRequest) GetTxid() string {
	if m != nil {
		return m.Txid
	}
	return ""
}

type SendPreBurntFeeTxidReply struct {
	NFTRegTxid           string   `protobuf:"bytes,1,opt,name=NFT_reg_txid,json=NFTRegTxid,proto3" json:"NFT_reg_txid,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SendPreBurntFeeTxidReply) Reset()         { *m = SendPreBurntFeeTxidReply{} }
func (m *SendPreBurntFeeTxidReply) String() string { return proto.CompactTextString(m) }
func (*SendPreBurntFeeTxidReply) ProtoMessage()    {}
func (*SendPreBurntFeeTxidReply) Descriptor() ([]byte, []int) {
	return fileDescriptor_47ba950f5e714b46, []int{3}
}

func (m *SendPreBurntFeeTxidReply) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SendPreBurntFeeTxidReply.Unmarshal(m, b)
}
func (m *SendPreBurntFeeTxidReply) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SendPreBurntFeeTxidReply.Marshal(b, m, deterministic)
}
func (m *SendPreBurntFeeTxidReply) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SendPreBurntFeeTxidReply.Merge(m, src)
}
func (m *SendPreBurntFeeTxidReply) XXX_Size() int {
	return xxx_messageInfo_SendPreBurntFeeTxidReply.Size(m)
}
func (m *SendPreBurntFeeTxidReply) XXX_DiscardUnknown() {
	xxx_messageInfo_SendPreBurntFeeTxidReply.DiscardUnknown(m)
}

var xxx_messageInfo_SendPreBurntFeeTxidReply proto.InternalMessageInfo

func (m *SendPreBurntFeeTxidReply) GetNFTRegTxid() string {
	if m != nil {
		return m.NFTRegTxid
	}
	return ""
}

type SendTicketRequest struct {
	Ticket               []byte   `protobuf:"bytes,1,opt,name=ticket,proto3" json:"ticket,omitempty"`
	TicketSignature      string   `protobuf:"bytes,2,opt,name=ticket_signature,json=ticketSignature,proto3" json:"ticket_signature,omitempty"`
	Fgpt                 string   `protobuf:"bytes,3,opt,name=fgpt,proto3" json:"fgpt,omitempty"`
	FgptSignature        string   `protobuf:"bytes,4,opt,name=fgpt_signature,json=fgptSignature,proto3" json:"fgpt_signature,omitempty"`
	FeeTxid              string   `protobuf:"bytes,5,opt,name=fee_txid,json=feeTxid,proto3" json:"fee_txid,omitempty"`
	Thumbnail            []byte   `protobuf:"bytes,6,opt,name=thumbnail,proto3" json:"thumbnail,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SendTicketRequest) Reset()         { *m = SendTicketRequest{} }
func (m *SendTicketRequest) String() string { return proto.CompactTextString(m) }
func (*SendTicketRequest) ProtoMessage()    {}
func (*SendTicketRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_47ba950f5e714b46, []int{4}
}

func (m *SendTicketRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SendTicketRequest.Unmarshal(m, b)
}
func (m *SendTicketRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SendTicketRequest.Marshal(b, m, deterministic)
}
func (m *SendTicketRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SendTicketRequest.Merge(m, src)
}
func (m *SendTicketRequest) XXX_Size() int {
	return xxx_messageInfo_SendTicketRequest.Size(m)
}
func (m *SendTicketRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_SendTicketRequest.DiscardUnknown(m)
}

var xxx_messageInfo_SendTicketRequest proto.InternalMessageInfo

func (m *SendTicketRequest) GetTicket() []byte {
	if m != nil {
		return m.Ticket
	}
	return nil
}

func (m *SendTicketRequest) GetTicketSignature() string {
	if m != nil {
		return m.TicketSignature
	}
	return ""
}

func (m *SendTicketRequest) GetFgpt() string {
	if m != nil {
		return m.Fgpt
	}
	return ""
}

func (m *SendTicketRequest) GetFgptSignature() string {
	if m != nil {
		return m.FgptSignature
	}
	return ""
}

func (m *SendTicketRequest) GetFeeTxid() string {
	if m != nil {
		return m.FeeTxid
	}
	return ""
}

func (m *SendTicketRequest) GetThumbnail() []byte {
	if m != nil {
		return m.Thumbnail
	}
	return nil
}

type SendTicketReply struct {
	TicketTxid           string   `protobuf:"bytes,1,opt,name=ticket_txid,json=ticketTxid,proto3" json:"ticket_txid,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SendTicketReply) Reset()         { *m = SendTicketReply{} }
func (m *SendTicketReply) String() string { return proto.CompactTextString(m) }
func (*SendTicketReply) ProtoMessage()    {}
func (*SendTicketReply) Descriptor() ([]byte, []int) {
	return fileDescriptor_47ba950f5e714b46, []int{5}
}

func (m *SendTicketReply) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SendTicketReply.Unmarshal(m, b)
}
func (m *SendTicketReply) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SendTicketReply.Marshal(b, m, deterministic)
}
func (m *SendTicketReply) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SendTicketReply.Merge(m, src)
}
func (m *SendTicketReply) XXX_Size() int {
	return xxx_messageInfo_SendTicketReply.Size(m)
}
func (m *SendTicketReply) XXX_DiscardUnknown() {
	xxx_messageInfo_SendTicketReply.DiscardUnknown(m)
}

var xxx_messageInfo_SendTicketReply proto.InternalMessageInfo

func (m *SendTicketReply) GetTicketTxid() string {
	if m != nil {
		return m.TicketTxid
	}
	return ""
}

type UploadImageRequest struct {
	// Types that are valid to be assigned to Payload:
	//	*UploadImageRequest_ImagePiece
	//	*UploadImageRequest_MetaData_
	Payload              isUploadImageRequest_Payload `protobuf_oneof:"payload"`
	XXX_NoUnkeyedLiteral struct{}                     `json:"-"`
	XXX_unrecognized     []byte                       `json:"-"`
	XXX_sizecache        int32                        `json:"-"`
}

func (m *UploadImageRequest) Reset()         { *m = UploadImageRequest{} }
func (m *UploadImageRequest) String() string { return proto.CompactTextString(m) }
func (*UploadImageRequest) ProtoMessage()    {}
func (*UploadImageRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_47ba950f5e714b46, []int{6}
}

func (m *UploadImageRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UploadImageRequest.Unmarshal(m, b)
}
func (m *UploadImageRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UploadImageRequest.Marshal(b, m, deterministic)
}
func (m *UploadImageRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UploadImageRequest.Merge(m, src)
}
func (m *UploadImageRequest) XXX_Size() int {
	return xxx_messageInfo_UploadImageRequest.Size(m)
}
func (m *UploadImageRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_UploadImageRequest.DiscardUnknown(m)
}

var xxx_messageInfo_UploadImageRequest proto.InternalMessageInfo

type isUploadImageRequest_Payload interface {
	isUploadImageRequest_Payload()
}

type UploadImageRequest_ImagePiece struct {
	ImagePiece []byte `protobuf:"bytes,1,opt,name=image_piece,json=imagePiece,proto3,oneof"`
}

type UploadImageRequest_MetaData_ struct {
	MetaData *UploadImageRequest_MetaData `protobuf:"bytes,2,opt,name=meta_data,json=metaData,proto3,oneof"`
}

func (*UploadImageRequest_ImagePiece) isUploadImageRequest_Payload() {}

func (*UploadImageRequest_MetaData_) isUploadImageRequest_Payload() {}

func (m *UploadImageRequest) GetPayload() isUploadImageRequest_Payload {
	if m != nil {
		return m.Payload
	}
	return nil
}

func (m *UploadImageRequest) GetImagePiece() []byte {
	if x, ok := m.GetPayload().(*UploadImageRequest_ImagePiece); ok {
		return x.ImagePiece
	}
	return nil
}

func (m *UploadImageRequest) GetMetaData() *UploadImageRequest_MetaData {
	if x, ok := m.GetPayload().(*UploadImageRequest_MetaData_); ok {
		return x.MetaData
	}
	return nil
}

// XXX_OneofWrappers is for the internal use of the proto package.
func (*UploadImageRequest) XXX_OneofWrappers() []interface{} {
	return []interface{}{
		(*UploadImageRequest_ImagePiece)(nil),
		(*UploadImageRequest_MetaData_)(nil),
	}
}

type UploadImageRequest_Coordinate struct {
	TopLeftX             int64    `protobuf:"varint,1,opt,name=top_left_x,json=topLeftX,proto3" json:"top_left_x,omitempty"`
	TopLeftY             int64    `protobuf:"varint,2,opt,name=top_left_y,json=topLeftY,proto3" json:"top_left_y,omitempty"`
	BottomRightX         int64    `protobuf:"varint,3,opt,name=bottom_right_x,json=bottomRightX,proto3" json:"bottom_right_x,omitempty"`
	BottomRightY         int64    `protobuf:"varint,4,opt,name=bottom_right_y,json=bottomRightY,proto3" json:"bottom_right_y,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *UploadImageRequest_Coordinate) Reset()         { *m = UploadImageRequest_Coordinate{} }
func (m *UploadImageRequest_Coordinate) String() string { return proto.CompactTextString(m) }
func (*UploadImageRequest_Coordinate) ProtoMessage()    {}
func (*UploadImageRequest_Coordinate) Descriptor() ([]byte, []int) {
	return fileDescriptor_47ba950f5e714b46, []int{6, 0}
}

func (m *UploadImageRequest_Coordinate) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UploadImageRequest_Coordinate.Unmarshal(m, b)
}
func (m *UploadImageRequest_Coordinate) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UploadImageRequest_Coordinate.Marshal(b, m, deterministic)
}
func (m *UploadImageRequest_Coordinate) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UploadImageRequest_Coordinate.Merge(m, src)
}
func (m *UploadImageRequest_Coordinate) XXX_Size() int {
	return xxx_messageInfo_UploadImageRequest_Coordinate.Size(m)
}
func (m *UploadImageRequest_Coordinate) XXX_DiscardUnknown() {
	xxx_messageInfo_UploadImageRequest_Coordinate.DiscardUnknown(m)
}

var xxx_messageInfo_UploadImageRequest_Coordinate proto.InternalMessageInfo

func (m *UploadImageRequest_Coordinate) GetTopLeftX() int64 {
	if m != nil {
		return m.TopLeftX
	}
	return 0
}

func (m *UploadImageRequest_Coordinate) GetTopLeftY() int64 {
	if m != nil {
		return m.TopLeftY
	}
	return 0
}

func (m *UploadImageRequest_Coordinate) GetBottomRightX() int64 {
	if m != nil {
		return m.BottomRightX
	}
	return 0
}

func (m *UploadImageRequest_Coordinate) GetBottomRightY() int64 {
	if m != nil {
		return m.BottomRightY
	}
	return 0
}

type UploadImageRequest_MetaData struct {
	// size of the image
	Size int64 `protobuf:"varint,1,opt,name=size,proto3" json:"size,omitempty"`
	// valid image format such as jpeg, png ...
	Format string `protobuf:"bytes,2,opt,name=format,proto3" json:"format,omitempty"`
	// sha3-256 hash of the image
	Hash []byte `protobuf:"bytes,3,opt,name=hash,proto3" json:"hash,omitempty"`
	// thumbnail of the image
	Thumbnail            *UploadImageRequest_Coordinate `protobuf:"bytes,4,opt,name=thumbnail,proto3" json:"thumbnail,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                       `json:"-"`
	XXX_unrecognized     []byte                         `json:"-"`
	XXX_sizecache        int32                          `json:"-"`
}

func (m *UploadImageRequest_MetaData) Reset()         { *m = UploadImageRequest_MetaData{} }
func (m *UploadImageRequest_MetaData) String() string { return proto.CompactTextString(m) }
func (*UploadImageRequest_MetaData) ProtoMessage()    {}
func (*UploadImageRequest_MetaData) Descriptor() ([]byte, []int) {
	return fileDescriptor_47ba950f5e714b46, []int{6, 1}
}

func (m *UploadImageRequest_MetaData) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UploadImageRequest_MetaData.Unmarshal(m, b)
}
func (m *UploadImageRequest_MetaData) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UploadImageRequest_MetaData.Marshal(b, m, deterministic)
}
func (m *UploadImageRequest_MetaData) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UploadImageRequest_MetaData.Merge(m, src)
}
func (m *UploadImageRequest_MetaData) XXX_Size() int {
	return xxx_messageInfo_UploadImageRequest_MetaData.Size(m)
}
func (m *UploadImageRequest_MetaData) XXX_DiscardUnknown() {
	xxx_messageInfo_UploadImageRequest_MetaData.DiscardUnknown(m)
}

var xxx_messageInfo_UploadImageRequest_MetaData proto.InternalMessageInfo

func (m *UploadImageRequest_MetaData) GetSize() int64 {
	if m != nil {
		return m.Size
	}
	return 0
}

func (m *UploadImageRequest_MetaData) GetFormat() string {
	if m != nil {
		return m.Format
	}
	return ""
}

func (m *UploadImageRequest_MetaData) GetHash() []byte {
	if m != nil {
		return m.Hash
	}
	return nil
}

func (m *UploadImageRequest_MetaData) GetThumbnail() *UploadImageRequest_Coordinate {
	if m != nil {
		return m.Thumbnail
	}
	return nil
}

type UploadImageReply struct {
	PreviewThumbnailHash []byte   `protobuf:"bytes,1,opt,name=preview_thumbnail_hash,json=previewThumbnailHash,proto3" json:"preview_thumbnail_hash,omitempty"`
	MediumThumbnailHash  []byte   `protobuf:"bytes,2,opt,name=medium_thumbnail_hash,json=mediumThumbnailHash,proto3" json:"medium_thumbnail_hash,omitempty"`
	SmallThumbnailHash   []byte   `protobuf:"bytes,3,opt,name=small_thumbnail_hash,json=smallThumbnailHash,proto3" json:"small_thumbnail_hash,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *UploadImageReply) Reset()         { *m = UploadImageReply{} }
func (m *UploadImageReply) String() string { return proto.CompactTextString(m) }
func (*UploadImageReply) ProtoMessage()    {}
func (*UploadImageReply) Descriptor() ([]byte, []int) {
	return fileDescriptor_47ba950f5e714b46, []int{7}
}

func (m *UploadImageReply) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UploadImageReply.Unmarshal(m, b)
}
func (m *UploadImageReply) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UploadImageReply.Marshal(b, m, deterministic)
}
func (m *UploadImageReply) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UploadImageReply.Merge(m, src)
}
func (m *UploadImageReply) XXX_Size() int {
	return xxx_messageInfo_UploadImageReply.Size(m)
}
func (m *UploadImageReply) XXX_DiscardUnknown() {
	xxx_messageInfo_UploadImageReply.DiscardUnknown(m)
}

var xxx_messageInfo_UploadImageReply proto.InternalMessageInfo

func (m *UploadImageReply) GetPreviewThumbnailHash() []byte {
	if m != nil {
		return m.PreviewThumbnailHash
	}
	return nil
}

func (m *UploadImageReply) GetMediumThumbnailHash() []byte {
	if m != nil {
		return m.MediumThumbnailHash
	}
	return nil
}

func (m *UploadImageReply) GetSmallThumbnailHash() []byte {
	if m != nil {
		return m.SmallThumbnailHash
	}
	return nil
}

func init() {
	proto.RegisterType((*SendSignedNFTTicketRequest)(nil), "walletnode.SendSignedNFTTicketRequest")
	proto.RegisterType((*SendSignedNFTTicketReply)(nil), "walletnode.SendSignedNFTTicketReply")
	proto.RegisterType((*SendPreBurntFeeTxidRequest)(nil), "walletnode.SendPreBurntFeeTxidRequest")
	proto.RegisterType((*SendPreBurntFeeTxidReply)(nil), "walletnode.SendPreBurntFeeTxidReply")
	proto.RegisterType((*SendTicketRequest)(nil), "walletnode.SendTicketRequest")
	proto.RegisterType((*SendTicketReply)(nil), "walletnode.SendTicketReply")
	proto.RegisterType((*UploadImageRequest)(nil), "walletnode.UploadImageRequest")
	proto.RegisterType((*UploadImageRequest_Coordinate)(nil), "walletnode.UploadImageRequest.Coordinate")
	proto.RegisterType((*UploadImageRequest_MetaData)(nil), "walletnode.UploadImageRequest.MetaData")
	proto.RegisterType((*UploadImageReply)(nil), "walletnode.UploadImageReply")
}

func init() { proto.RegisterFile("register_nft_wn.proto", fileDescriptor_47ba950f5e714b46) }

var fileDescriptor_47ba950f5e714b46 = []byte{
	// 964 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x56, 0xdd, 0x72, 0xdb, 0x44,
	0x14, 0xae, 0x9a, 0x1f, 0xdb, 0xc7, 0x6e, 0x93, 0x6c, 0xd3, 0xe2, 0x9a, 0x34, 0x18, 0x4f, 0x81,
	0x74, 0x98, 0x71, 0x32, 0x86, 0x4b, 0x6e, 0xc8, 0x8f, 0x1b, 0x32, 0xc4, 0x13, 0x14, 0x33, 0xb4,
	0xdc, 0x68, 0xd6, 0xd2, 0x91, 0xbc, 0x53, 0x49, 0xab, 0xac, 0xd6, 0x24, 0xe6, 0x9e, 0xe1, 0x01,
	0xe0, 0x39, 0xb8, 0xe5, 0x1d, 0x78, 0x2a, 0x66, 0x77, 0x25, 0x5b, 0xb2, 0x9d, 0xf4, 0xca, 0xab,
	0xef, 0xfb, 0xce, 0xd9, 0xb3, 0xdf, 0xd9, 0x1f, 0xc3, 0x73, 0x81, 0x01, 0x4b, 0x25, 0x0a, 0x27,
	0xf6, 0xa5, 0x73, 0x1b, 0x77, 0x13, 0xc1, 0x25, 0x27, 0x70, 0x4b, 0xc3, 0x10, 0x65, 0xcc, 0x3d,
	0x6c, 0x6d, 0xb9, 0x3c, 0x8a, 0x78, 0x3c, 0x23, 0x3b, 0x7f, 0x3c, 0x86, 0xd6, 0x35, 0xc6, 0xde,
	0x35, 0x0b, 0x62, 0xf4, 0x06, 0xfd, 0xe1, 0x90, 0xb9, 0x1f, 0x50, 0xda, 0x78, 0x33, 0xc1, 0x54,
	0x92, 0x57, 0x00, 0x2a, 0x97, 0xd4, 0x60, 0xd3, 0x6a, 0x5b, 0x07, 0x0d, 0xbb, 0x16, 0xfb, 0xd2,
	0xa8, 0xc8, 0xd7, 0xb0, 0xe3, 0x0a, 0xa4, 0x92, 0x0b, 0x27, 0x65, 0x41, 0x4c, 0xe5, 0x44, 0x60,
	0xf3, 0xb1, 0x56, 0x6d, 0x67, 0xc4, 0x75, 0x8e, 0x93, 0x5d, 0xd8, 0x08, 0xe9, 0x08, 0xc3, 0xe6,
	0x5a, 0xdb, 0x3a, 0xa8, 0xd9, 0xe6, 0x83, 0x5c, 0xc0, 0x0e, 0xc6, 0x2e, 0xf7, 0xd0, 0x49, 0xa8,
	0xa0, 0x11, 0x4a, 0x14, 0x69, 0x73, 0xa3, 0x6d, 0x1d, 0xd4, 0x7b, 0xaf, 0xba, 0xf3, 0xca, 0xbb,
	0x67, 0x5a, 0x24, 0xae, 0x66, 0x22, 0x7b, 0xdb, 0xc4, 0xcd, 0x11, 0xb2, 0x0f, 0x75, 0xcf, 0x73,
	0xfc, 0xc4, 0xf1, 0x59, 0x88, 0x69, 0x73, 0xd3, 0x94, 0xeb, 0x79, 0xfd, 0xa4, 0xaf, 0x00, 0xf2,
	0x12, 0xaa, 0xe2, 0x26, 0x23, 0x2b, 0x9a, 0xac, 0x88, 0x1b, 0x4d, 0x75, 0xce, 0xa0, 0xb9, 0xd2,
	0x86, 0x24, 0x9c, 0x92, 0x37, 0xb0, 0x6d, 0x9c, 0x15, 0x54, 0x32, 0x1e, 0x3b, 0x3e, 0xa2, 0xb6,
	0x62, 0xcd, 0xde, 0x2a, 0xe2, 0x7d, 0xc4, 0xce, 0x91, 0x71, 0xf3, 0x4a, 0xe0, 0xf1, 0x44, 0xc4,
	0xb2, 0x8f, 0x38, 0xbc, 0x63, 0x5e, 0xee, 0x26, 0x81, 0x75, 0x79, 0xc7, 0x3c, 0x1d, 0x5c, 0xb3,
	0xf5, 0xb8, 0xf3, 0x9d, 0x99, 0x78, 0x29, 0x42, 0x4d, 0xdc, 0x86, 0xc6, 0xa0, 0x3f, 0x74, 0x04,
	0x06, 0x4e, 0x21, 0x0e, 0x06, 0xfd, 0xa1, 0x8d, 0x81, 0x92, 0x75, 0xfe, 0xb3, 0x60, 0x47, 0x85,
	0x97, 0xbb, 0xf6, 0x02, 0x36, 0x4b, 0x1d, 0xcb, 0xbe, 0xd4, 0x42, 0xcc, 0x68, 0xa1, 0x5b, 0x35,
	0x7b, 0xcb, 0xe0, 0xf3, 0x66, 0x11, 0x58, 0xf7, 0x83, 0x44, 0x66, 0xbd, 0xd2, 0x63, 0xf2, 0x05,
	0x3c, 0x55, 0xbf, 0x85, 0xe0, 0x75, 0xcd, 0x3e, 0x51, 0xe8, 0x3c, 0xf4, 0x25, 0x54, 0x7d, 0x44,
	0x53, 0xf1, 0x86, 0x16, 0x54, 0x7c, 0xb3, 0x2a, 0xb2, 0x07, 0x35, 0x39, 0x9e, 0x44, 0xa3, 0x98,
	0xb2, 0x30, 0x6f, 0xcf, 0x0c, 0xe8, 0xf4, 0x60, 0xab, 0xb8, 0x16, 0xe5, 0xc0, 0x67, 0x50, 0xcf,
	0x2a, 0x2e, 0x1a, 0x60, 0x20, 0x6d, 0xc0, 0xbf, 0x6b, 0x40, 0x7e, 0x4e, 0x42, 0x4e, 0xbd, 0x1f,
	0x22, 0x1a, 0x60, 0xee, 0xc0, 0xe7, 0x50, 0x67, 0xea, 0xdb, 0x49, 0x18, 0xba, 0xa6, 0x5b, 0x8d,
	0xf3, 0x47, 0x36, 0x68, 0xf0, 0x4a, 0x61, 0xa4, 0x0f, 0xb5, 0x08, 0x25, 0x75, 0x3c, 0x2a, 0xa9,
	0x76, 0xa1, 0xde, 0xfb, 0xaa, 0xb8, 0xe1, 0x96, 0xb3, 0x76, 0x2f, 0x51, 0xd2, 0x53, 0x2a, 0xe9,
	0xf9, 0x23, 0xbb, 0x1a, 0x65, 0xe3, 0xd6, 0xdf, 0x16, 0xc0, 0x09, 0xe7, 0xc2, 0x63, 0x31, 0x95,
	0x48, 0xf6, 0x00, 0x24, 0x4f, 0x9c, 0x10, 0x7d, 0xe9, 0xdc, 0x65, 0xdb, 0xa4, 0x2a, 0x79, 0xf2,
	0x23, 0xfa, 0xf2, 0x5d, 0x89, 0x9d, 0xea, 0x59, 0xe7, 0xec, 0x7b, 0xf2, 0x1a, 0x9e, 0x8e, 0xb8,
	0x94, 0x3c, 0x72, 0x04, 0x0b, 0xc6, 0x2a, 0x7e, 0x4d, 0x2b, 0x1a, 0x06, 0xb5, 0x15, 0xf8, 0x6e,
	0x49, 0x35, 0xd5, 0x6d, 0x28, 0xab, 0xde, 0xb7, 0xfe, 0xb2, 0xa0, 0x9a, 0xd7, 0xab, 0xba, 0x99,
	0xb2, 0xdf, 0xf3, 0x5d, 0xab, 0xc7, 0x6a, 0x93, 0xf8, 0x5c, 0x44, 0x54, 0x66, 0x5b, 0x20, 0xfb,
	0x52, 0xda, 0x31, 0x4d, 0xc7, 0x7a, 0xea, 0x86, 0xad, 0xc7, 0xe4, 0x6d, 0xb1, 0x6f, 0xeb, 0xda,
	0xab, 0x37, 0x1f, 0xf1, 0x6a, 0x6e, 0x49, 0xa1, 0xc5, 0xc7, 0x35, 0xa8, 0x24, 0x74, 0xaa, 0xc4,
	0x9d, 0x7f, 0x2c, 0xd8, 0x2e, 0xc5, 0xa9, 0x7e, 0x7f, 0x0b, 0x2f, 0x12, 0x81, 0xbf, 0x31, 0xbc,
	0x75, 0x66, 0x41, 0x8e, 0x2e, 0xc7, 0xec, 0xe4, 0xdd, 0x8c, 0x1d, 0xe6, 0xe4, 0xb9, 0x2a, 0xaf,
	0x07, 0xcf, 0x23, 0xf4, 0xd8, 0x24, 0x5a, 0x0c, 0x32, 0x57, 0xd1, 0x33, 0x43, 0x96, 0x63, 0x8e,
	0x60, 0x37, 0x8d, 0x68, 0x18, 0x2e, 0x86, 0x98, 0x65, 0x13, 0xcd, 0x95, 0x22, 0x7a, 0x7f, 0x56,
	0xa0, 0x6e, 0x67, 0x37, 0xec, 0xc0, 0x97, 0xe4, 0x04, 0x2a, 0xd7, 0x98, 0xa6, 0x8c, 0xc7, 0xa4,
	0x55, 0x34, 0x23, 0x03, 0x33, 0x23, 0x5a, 0xcd, 0x95, 0x5c, 0x12, 0x4e, 0x0f, 0xac, 0x23, 0x8b,
	0xfc, 0x04, 0x4f, 0xbe, 0x77, 0x5d, 0x4c, 0x24, 0x7a, 0x03, 0xee, 0x61, 0x4a, 0xda, 0x45, 0x79,
	0x89, 0xca, 0x13, 0xee, 0x3f, 0xa0, 0x50, 0x1e, 0x9e, 0x41, 0xed, 0x84, 0xc7, 0x31, 0xba, 0x72,
	0xc8, 0xc9, 0x5e, 0x51, 0x3c, 0x83, 0xf3, 0x54, 0xad, 0x7b, 0xd8, 0x2c, 0xcd, 0x25, 0xa6, 0x63,
	0x53, 0x55, 0x29, 0xcd, 0x0c, 0x5e, 0x99, 0xa6, 0xc0, 0xaa, 0x34, 0xbf, 0x98, 0x43, 0x6d, 0x63,
	0xa0, 0x76, 0xa3, 0x3a, 0x6c, 0xa4, 0x53, 0x76, 0xa4, 0x44, 0xe6, 0x29, 0xdb, 0x0f, 0x6a, 0x54,
	0xe2, 0x0b, 0x80, 0x2b, 0xc1, 0x47, 0xa8, 0x77, 0x0f, 0x29, 0xbd, 0x15, 0x73, 0x3c, 0x4f, 0xf7,
	0xe9, 0x7d, 0xb4, 0xee, 0x03, 0x71, 0xe1, 0xd9, 0x8a, 0xdb, 0x9f, 0x7c, 0xb9, 0x58, 0xc4, 0xea,
	0x57, 0xb2, 0xf5, 0xfa, 0xa3, 0x3a, 0x55, 0x70, 0x36, 0xc9, 0xc2, 0x4d, 0xbf, 0x3c, 0xc9, 0xea,
	0xc7, 0x63, 0x79, 0x92, 0x95, 0x4f, 0xc6, 0x39, 0xc0, 0xfc, 0x0e, 0x2d, 0xbb, 0xb2, 0xf4, 0x4e,
	0x94, 0x5d, 0x59, 0xbc, 0x7a, 0x2f, 0xa1, 0x5e, 0x38, 0x9e, 0x64, 0xff, 0xe1, 0xf3, 0xde, 0xda,
	0xbb, 0x97, 0x37, 0x16, 0x5f, 0xc0, 0xce, 0x5b, 0x94, 0xa7, 0xa7, 0xea, 0x3e, 0x1a, 0xd1, 0x14,
	0xf5, 0x21, 0x2c, 0x05, 0x29, 0xfa, 0x58, 0xc1, 0x79, 0xca, 0x4f, 0x8a, 0x6c, 0x4e, 0x25, 0xe1,
	0xf4, 0xb8, 0xf7, 0xeb, 0x51, 0xc0, 0xe4, 0x78, 0x32, 0xea, 0xba, 0x3c, 0x3a, 0x4c, 0x68, 0x2a,
	0x31, 0x8c, 0x51, 0xde, 0x72, 0xf1, 0xe1, 0x30, 0xe0, 0x4a, 0x7e, 0xa8, 0xff, 0xdd, 0x1c, 0xce,
	0xe3, 0x47, 0x9b, 0x1a, 0xf9, 0xe6, 0xff, 0x00, 0x00, 0x00, 0xff, 0xff, 0xc1, 0xcc, 0xc4, 0x4e,
	0x25, 0x09, 0x00, 0x00,
}
