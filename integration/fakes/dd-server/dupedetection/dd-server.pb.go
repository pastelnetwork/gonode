// Code generated by protoc-gen-go. DO NOT EDIT.
// source: dd-server.proto

package dupedetection

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

type RarenessScoreRequest struct {
	Path                 string   `protobuf:"bytes,1,opt,name=path,proto3" json:"path,omitempty"`
	BlockHash            string   `protobuf:"bytes,2,opt,name=block_hash,json=blockHash,proto3" json:"block_hash,omitempty"`
	PastelId             string   `protobuf:"bytes,3,opt,name=pastel_id,json=pastelId,proto3" json:"pastel_id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RarenessScoreRequest) Reset()         { *m = RarenessScoreRequest{} }
func (m *RarenessScoreRequest) String() string { return proto.CompactTextString(m) }
func (*RarenessScoreRequest) ProtoMessage()    {}
func (*RarenessScoreRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_1ddb247232398ef5, []int{0}
}

func (m *RarenessScoreRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RarenessScoreRequest.Unmarshal(m, b)
}
func (m *RarenessScoreRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RarenessScoreRequest.Marshal(b, m, deterministic)
}
func (m *RarenessScoreRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RarenessScoreRequest.Merge(m, src)
}
func (m *RarenessScoreRequest) XXX_Size() int {
	return xxx_messageInfo_RarenessScoreRequest.Size(m)
}
func (m *RarenessScoreRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_RarenessScoreRequest.DiscardUnknown(m)
}

var xxx_messageInfo_RarenessScoreRequest proto.InternalMessageInfo

func (m *RarenessScoreRequest) GetPath() string {
	if m != nil {
		return m.Path
	}
	return ""
}

func (m *RarenessScoreRequest) GetBlockHash() string {
	if m != nil {
		return m.BlockHash
	}
	return ""
}

func (m *RarenessScoreRequest) GetPastelId() string {
	if m != nil {
		return m.PastelId
	}
	return ""
}

type ImageRarenessScoreReply struct {
	Block                                string                 `protobuf:"bytes,1,opt,name=block,proto3" json:"block,omitempty"`
	Principal                            string                 `protobuf:"bytes,2,opt,name=principal,proto3" json:"principal,omitempty"`
	DupeDetectionSystemVersion           string                 `protobuf:"bytes,3,opt,name=dupe_detection_system_version,json=dupeDetectionSystemVersion,proto3" json:"dupe_detection_system_version,omitempty"`
	IsLikelyDupe                         bool                   `protobuf:"varint,4,opt,name=is_likely_dupe,json=isLikelyDupe,proto3" json:"is_likely_dupe,omitempty"`
	IsRareOnInternet                     bool                   `protobuf:"varint,5,opt,name=is_rare_on_internet,json=isRareOnInternet,proto3" json:"is_rare_on_internet,omitempty"`
	RarenessScores                       *RarenessScores        `protobuf:"bytes,6,opt,name=rareness_scores,json=rarenessScores,proto3" json:"rareness_scores,omitempty"`
	InternetRareness                     *InternetRareness      `protobuf:"bytes,7,opt,name=internet_rareness,json=internetRareness,proto3" json:"internet_rareness,omitempty"`
	OpenNsfwScore                        float32                `protobuf:"fixed32,8,opt,name=open_nsfw_score,json=openNsfwScore,proto3" json:"open_nsfw_score,omitempty"`
	AlternativeNsfwScores                *AltNsfwScores         `protobuf:"bytes,9,opt,name=alternative_nsfw_scores,json=alternativeNsfwScores,proto3" json:"alternative_nsfw_scores,omitempty"`
	ImageFingerprintOfCandidateImageFile []float32              `protobuf:"fixed32,10,rep,packed,name=image_fingerprint_of_candidate_image_file,json=imageFingerprintOfCandidateImageFile,proto3" json:"image_fingerprint_of_candidate_image_file,omitempty"`
	FingerprintsStat                     *FingerprintsStat      `protobuf:"bytes,11,opt,name=fingerprints_stat,json=fingerprintsStat,proto3" json:"fingerprints_stat,omitempty"`
	HashOfCandidateImageFile             string                 `protobuf:"bytes,12,opt,name=hash_of_candidate_image_file,json=hashOfCandidateImageFile,proto3" json:"hash_of_candidate_image_file,omitempty"`
	PerceptualImageHashes                *PerceptualImageHashes `protobuf:"bytes,13,opt,name=perceptual_image_hashes,json=perceptualImageHashes,proto3" json:"perceptual_image_hashes,omitempty"`
	PerceptualHashOverlapCount           uint32                 `protobuf:"varint,14,opt,name=perceptual_hash_overlap_count,json=perceptualHashOverlapCount,proto3" json:"perceptual_hash_overlap_count,omitempty"`
	Maxes                                *Maxes                 `protobuf:"bytes,15,opt,name=maxes,proto3" json:"maxes,omitempty"`
	Percentile                           *Percentile            `protobuf:"bytes,16,opt,name=percentile,proto3" json:"percentile,omitempty"`
	XXX_NoUnkeyedLiteral                 struct{}               `json:"-"`
	XXX_unrecognized                     []byte                 `json:"-"`
	XXX_sizecache                        int32                  `json:"-"`
}

func (m *ImageRarenessScoreReply) Reset()         { *m = ImageRarenessScoreReply{} }
func (m *ImageRarenessScoreReply) String() string { return proto.CompactTextString(m) }
func (*ImageRarenessScoreReply) ProtoMessage()    {}
func (*ImageRarenessScoreReply) Descriptor() ([]byte, []int) {
	return fileDescriptor_1ddb247232398ef5, []int{1}
}

func (m *ImageRarenessScoreReply) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ImageRarenessScoreReply.Unmarshal(m, b)
}
func (m *ImageRarenessScoreReply) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ImageRarenessScoreReply.Marshal(b, m, deterministic)
}
func (m *ImageRarenessScoreReply) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ImageRarenessScoreReply.Merge(m, src)
}
func (m *ImageRarenessScoreReply) XXX_Size() int {
	return xxx_messageInfo_ImageRarenessScoreReply.Size(m)
}
func (m *ImageRarenessScoreReply) XXX_DiscardUnknown() {
	xxx_messageInfo_ImageRarenessScoreReply.DiscardUnknown(m)
}

var xxx_messageInfo_ImageRarenessScoreReply proto.InternalMessageInfo

func (m *ImageRarenessScoreReply) GetBlock() string {
	if m != nil {
		return m.Block
	}
	return ""
}

func (m *ImageRarenessScoreReply) GetPrincipal() string {
	if m != nil {
		return m.Principal
	}
	return ""
}

func (m *ImageRarenessScoreReply) GetDupeDetectionSystemVersion() string {
	if m != nil {
		return m.DupeDetectionSystemVersion
	}
	return ""
}

func (m *ImageRarenessScoreReply) GetIsLikelyDupe() bool {
	if m != nil {
		return m.IsLikelyDupe
	}
	return false
}

func (m *ImageRarenessScoreReply) GetIsRareOnInternet() bool {
	if m != nil {
		return m.IsRareOnInternet
	}
	return false
}

func (m *ImageRarenessScoreReply) GetRarenessScores() *RarenessScores {
	if m != nil {
		return m.RarenessScores
	}
	return nil
}

func (m *ImageRarenessScoreReply) GetInternetRareness() *InternetRareness {
	if m != nil {
		return m.InternetRareness
	}
	return nil
}

func (m *ImageRarenessScoreReply) GetOpenNsfwScore() float32 {
	if m != nil {
		return m.OpenNsfwScore
	}
	return 0
}

func (m *ImageRarenessScoreReply) GetAlternativeNsfwScores() *AltNsfwScores {
	if m != nil {
		return m.AlternativeNsfwScores
	}
	return nil
}

func (m *ImageRarenessScoreReply) GetImageFingerprintOfCandidateImageFile() []float32 {
	if m != nil {
		return m.ImageFingerprintOfCandidateImageFile
	}
	return nil
}

func (m *ImageRarenessScoreReply) GetFingerprintsStat() *FingerprintsStat {
	if m != nil {
		return m.FingerprintsStat
	}
	return nil
}

func (m *ImageRarenessScoreReply) GetHashOfCandidateImageFile() string {
	if m != nil {
		return m.HashOfCandidateImageFile
	}
	return ""
}

func (m *ImageRarenessScoreReply) GetPerceptualImageHashes() *PerceptualImageHashes {
	if m != nil {
		return m.PerceptualImageHashes
	}
	return nil
}

func (m *ImageRarenessScoreReply) GetPerceptualHashOverlapCount() uint32 {
	if m != nil {
		return m.PerceptualHashOverlapCount
	}
	return 0
}

func (m *ImageRarenessScoreReply) GetMaxes() *Maxes {
	if m != nil {
		return m.Maxes
	}
	return nil
}

func (m *ImageRarenessScoreReply) GetPercentile() *Percentile {
	if m != nil {
		return m.Percentile
	}
	return nil
}

type RarenessScores struct {
	CombinedRarenessScore         float32  `protobuf:"fixed32,1,opt,name=combined_rareness_score,json=combinedRarenessScore,proto3" json:"combined_rareness_score,omitempty"`
	XgboostPredictedRarenessScore float32  `protobuf:"fixed32,2,opt,name=xgboost_predicted_rareness_score,json=xgboostPredictedRarenessScore,proto3" json:"xgboost_predicted_rareness_score,omitempty"`
	NnPredictedRarenessScore      float32  `protobuf:"fixed32,3,opt,name=nn_predicted_rareness_score,json=nnPredictedRarenessScore,proto3" json:"nn_predicted_rareness_score,omitempty"`
	OverallAverageRarenessScore   float32  `protobuf:"fixed32,4,opt,name=overall_average_rareness_score,json=overallAverageRarenessScore,proto3" json:"overall_average_rareness_score,omitempty"`
	XXX_NoUnkeyedLiteral          struct{} `json:"-"`
	XXX_unrecognized              []byte   `json:"-"`
	XXX_sizecache                 int32    `json:"-"`
}

func (m *RarenessScores) Reset()         { *m = RarenessScores{} }
func (m *RarenessScores) String() string { return proto.CompactTextString(m) }
func (*RarenessScores) ProtoMessage()    {}
func (*RarenessScores) Descriptor() ([]byte, []int) {
	return fileDescriptor_1ddb247232398ef5, []int{2}
}

func (m *RarenessScores) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RarenessScores.Unmarshal(m, b)
}
func (m *RarenessScores) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RarenessScores.Marshal(b, m, deterministic)
}
func (m *RarenessScores) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RarenessScores.Merge(m, src)
}
func (m *RarenessScores) XXX_Size() int {
	return xxx_messageInfo_RarenessScores.Size(m)
}
func (m *RarenessScores) XXX_DiscardUnknown() {
	xxx_messageInfo_RarenessScores.DiscardUnknown(m)
}

var xxx_messageInfo_RarenessScores proto.InternalMessageInfo

func (m *RarenessScores) GetCombinedRarenessScore() float32 {
	if m != nil {
		return m.CombinedRarenessScore
	}
	return 0
}

func (m *RarenessScores) GetXgboostPredictedRarenessScore() float32 {
	if m != nil {
		return m.XgboostPredictedRarenessScore
	}
	return 0
}

func (m *RarenessScores) GetNnPredictedRarenessScore() float32 {
	if m != nil {
		return m.NnPredictedRarenessScore
	}
	return 0
}

func (m *RarenessScores) GetOverallAverageRarenessScore() float32 {
	if m != nil {
		return m.OverallAverageRarenessScore
	}
	return 0
}

type InternetRareness struct {
	MatchesFoundOnFirstPage uint32   `protobuf:"varint,1,opt,name=matches_found_on_first_page,json=matchesFoundOnFirstPage,proto3" json:"matches_found_on_first_page,omitempty"`
	NumberOfPagesOfResults  uint32   `protobuf:"varint,2,opt,name=number_of_pages_of_results,json=numberOfPagesOfResults,proto3" json:"number_of_pages_of_results,omitempty"`
	UrlOfFirstMatchInPage   string   `protobuf:"bytes,3,opt,name=url_of_first_match_in_page,json=urlOfFirstMatchInPage,proto3" json:"url_of_first_match_in_page,omitempty"`
	XXX_NoUnkeyedLiteral    struct{} `json:"-"`
	XXX_unrecognized        []byte   `json:"-"`
	XXX_sizecache           int32    `json:"-"`
}

func (m *InternetRareness) Reset()         { *m = InternetRareness{} }
func (m *InternetRareness) String() string { return proto.CompactTextString(m) }
func (*InternetRareness) ProtoMessage()    {}
func (*InternetRareness) Descriptor() ([]byte, []int) {
	return fileDescriptor_1ddb247232398ef5, []int{3}
}

func (m *InternetRareness) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_InternetRareness.Unmarshal(m, b)
}
func (m *InternetRareness) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_InternetRareness.Marshal(b, m, deterministic)
}
func (m *InternetRareness) XXX_Merge(src proto.Message) {
	xxx_messageInfo_InternetRareness.Merge(m, src)
}
func (m *InternetRareness) XXX_Size() int {
	return xxx_messageInfo_InternetRareness.Size(m)
}
func (m *InternetRareness) XXX_DiscardUnknown() {
	xxx_messageInfo_InternetRareness.DiscardUnknown(m)
}

var xxx_messageInfo_InternetRareness proto.InternalMessageInfo

func (m *InternetRareness) GetMatchesFoundOnFirstPage() uint32 {
	if m != nil {
		return m.MatchesFoundOnFirstPage
	}
	return 0
}

func (m *InternetRareness) GetNumberOfPagesOfResults() uint32 {
	if m != nil {
		return m.NumberOfPagesOfResults
	}
	return 0
}

func (m *InternetRareness) GetUrlOfFirstMatchInPage() string {
	if m != nil {
		return m.UrlOfFirstMatchInPage
	}
	return ""
}

type AltNsfwScores struct {
	Drawings             float32  `protobuf:"fixed32,1,opt,name=drawings,proto3" json:"drawings,omitempty"`
	Hentai               float32  `protobuf:"fixed32,2,opt,name=hentai,proto3" json:"hentai,omitempty"`
	Neutral              float32  `protobuf:"fixed32,3,opt,name=neutral,proto3" json:"neutral,omitempty"`
	Porn                 float32  `protobuf:"fixed32,4,opt,name=porn,proto3" json:"porn,omitempty"`
	Sexy                 float32  `protobuf:"fixed32,5,opt,name=sexy,proto3" json:"sexy,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *AltNsfwScores) Reset()         { *m = AltNsfwScores{} }
func (m *AltNsfwScores) String() string { return proto.CompactTextString(m) }
func (*AltNsfwScores) ProtoMessage()    {}
func (*AltNsfwScores) Descriptor() ([]byte, []int) {
	return fileDescriptor_1ddb247232398ef5, []int{4}
}

func (m *AltNsfwScores) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AltNsfwScores.Unmarshal(m, b)
}
func (m *AltNsfwScores) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AltNsfwScores.Marshal(b, m, deterministic)
}
func (m *AltNsfwScores) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AltNsfwScores.Merge(m, src)
}
func (m *AltNsfwScores) XXX_Size() int {
	return xxx_messageInfo_AltNsfwScores.Size(m)
}
func (m *AltNsfwScores) XXX_DiscardUnknown() {
	xxx_messageInfo_AltNsfwScores.DiscardUnknown(m)
}

var xxx_messageInfo_AltNsfwScores proto.InternalMessageInfo

func (m *AltNsfwScores) GetDrawings() float32 {
	if m != nil {
		return m.Drawings
	}
	return 0
}

func (m *AltNsfwScores) GetHentai() float32 {
	if m != nil {
		return m.Hentai
	}
	return 0
}

func (m *AltNsfwScores) GetNeutral() float32 {
	if m != nil {
		return m.Neutral
	}
	return 0
}

func (m *AltNsfwScores) GetPorn() float32 {
	if m != nil {
		return m.Porn
	}
	return 0
}

func (m *AltNsfwScores) GetSexy() float32 {
	if m != nil {
		return m.Sexy
	}
	return 0
}

type PerceptualImageHashes struct {
	PdqHash              string   `protobuf:"bytes,1,opt,name=pdq_hash,json=pdqHash,proto3" json:"pdq_hash,omitempty"`
	PerceptualHash       string   `protobuf:"bytes,2,opt,name=perceptual_hash,json=perceptualHash,proto3" json:"perceptual_hash,omitempty"`
	AverageHash          string   `protobuf:"bytes,3,opt,name=average_hash,json=averageHash,proto3" json:"average_hash,omitempty"`
	DifferenceHash       string   `protobuf:"bytes,4,opt,name=difference_hash,json=differenceHash,proto3" json:"difference_hash,omitempty"`
	NeuralhashHash       string   `protobuf:"bytes,5,opt,name=neuralhash_hash,json=neuralhashHash,proto3" json:"neuralhash_hash,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PerceptualImageHashes) Reset()         { *m = PerceptualImageHashes{} }
func (m *PerceptualImageHashes) String() string { return proto.CompactTextString(m) }
func (*PerceptualImageHashes) ProtoMessage()    {}
func (*PerceptualImageHashes) Descriptor() ([]byte, []int) {
	return fileDescriptor_1ddb247232398ef5, []int{5}
}

func (m *PerceptualImageHashes) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PerceptualImageHashes.Unmarshal(m, b)
}
func (m *PerceptualImageHashes) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PerceptualImageHashes.Marshal(b, m, deterministic)
}
func (m *PerceptualImageHashes) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PerceptualImageHashes.Merge(m, src)
}
func (m *PerceptualImageHashes) XXX_Size() int {
	return xxx_messageInfo_PerceptualImageHashes.Size(m)
}
func (m *PerceptualImageHashes) XXX_DiscardUnknown() {
	xxx_messageInfo_PerceptualImageHashes.DiscardUnknown(m)
}

var xxx_messageInfo_PerceptualImageHashes proto.InternalMessageInfo

func (m *PerceptualImageHashes) GetPdqHash() string {
	if m != nil {
		return m.PdqHash
	}
	return ""
}

func (m *PerceptualImageHashes) GetPerceptualHash() string {
	if m != nil {
		return m.PerceptualHash
	}
	return ""
}

func (m *PerceptualImageHashes) GetAverageHash() string {
	if m != nil {
		return m.AverageHash
	}
	return ""
}

func (m *PerceptualImageHashes) GetDifferenceHash() string {
	if m != nil {
		return m.DifferenceHash
	}
	return ""
}

func (m *PerceptualImageHashes) GetNeuralhashHash() string {
	if m != nil {
		return m.NeuralhashHash
	}
	return ""
}

type FingerprintsStat struct {
	NumberOfFingerprintsRequiringFurtherTesting_1 uint32   `protobuf:"varint,1,opt,name=number_of_fingerprints_requiring_further_testing_1,json=numberOfFingerprintsRequiringFurtherTesting1,proto3" json:"number_of_fingerprints_requiring_further_testing_1,omitempty"`
	NumberOfFingerprintsRequiringFurtherTesting_2 uint32   `protobuf:"varint,2,opt,name=number_of_fingerprints_requiring_further_testing_2,json=numberOfFingerprintsRequiringFurtherTesting2,proto3" json:"number_of_fingerprints_requiring_further_testing_2,omitempty"`
	NumberOfFingerprintsRequiringFurtherTesting_3 uint32   `protobuf:"varint,3,opt,name=number_of_fingerprints_requiring_further_testing_3,json=numberOfFingerprintsRequiringFurtherTesting3,proto3" json:"number_of_fingerprints_requiring_further_testing_3,omitempty"`
	NumberOfFingerprintsRequiringFurtherTesting_4 uint32   `protobuf:"varint,4,opt,name=number_of_fingerprints_requiring_further_testing_4,json=numberOfFingerprintsRequiringFurtherTesting4,proto3" json:"number_of_fingerprints_requiring_further_testing_4,omitempty"`
	NumberOfFingerprintsRequiringFurtherTesting_5 uint32   `protobuf:"varint,5,opt,name=number_of_fingerprints_requiring_further_testing_5,json=numberOfFingerprintsRequiringFurtherTesting5,proto3" json:"number_of_fingerprints_requiring_further_testing_5,omitempty"`
	NumberOfFingerprintsRequiringFurtherTesting_6 uint32   `protobuf:"varint,6,opt,name=number_of_fingerprints_requiring_further_testing_6,json=numberOfFingerprintsRequiringFurtherTesting6,proto3" json:"number_of_fingerprints_requiring_further_testing_6,omitempty"`
	NumberOfFingerprintsOfSuspectedDupes          uint32   `protobuf:"varint,7,opt,name=number_of_fingerprints_of_suspected_dupes,json=numberOfFingerprintsOfSuspectedDupes,proto3" json:"number_of_fingerprints_of_suspected_dupes,omitempty"`
	XXX_NoUnkeyedLiteral                          struct{} `json:"-"`
	XXX_unrecognized                              []byte   `json:"-"`
	XXX_sizecache                                 int32    `json:"-"`
}

func (m *FingerprintsStat) Reset()         { *m = FingerprintsStat{} }
func (m *FingerprintsStat) String() string { return proto.CompactTextString(m) }
func (*FingerprintsStat) ProtoMessage()    {}
func (*FingerprintsStat) Descriptor() ([]byte, []int) {
	return fileDescriptor_1ddb247232398ef5, []int{6}
}

func (m *FingerprintsStat) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_FingerprintsStat.Unmarshal(m, b)
}
func (m *FingerprintsStat) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_FingerprintsStat.Marshal(b, m, deterministic)
}
func (m *FingerprintsStat) XXX_Merge(src proto.Message) {
	xxx_messageInfo_FingerprintsStat.Merge(m, src)
}
func (m *FingerprintsStat) XXX_Size() int {
	return xxx_messageInfo_FingerprintsStat.Size(m)
}
func (m *FingerprintsStat) XXX_DiscardUnknown() {
	xxx_messageInfo_FingerprintsStat.DiscardUnknown(m)
}

var xxx_messageInfo_FingerprintsStat proto.InternalMessageInfo

func (m *FingerprintsStat) GetNumberOfFingerprintsRequiringFurtherTesting_1() uint32 {
	if m != nil {
		return m.NumberOfFingerprintsRequiringFurtherTesting_1
	}
	return 0
}

func (m *FingerprintsStat) GetNumberOfFingerprintsRequiringFurtherTesting_2() uint32 {
	if m != nil {
		return m.NumberOfFingerprintsRequiringFurtherTesting_2
	}
	return 0
}

func (m *FingerprintsStat) GetNumberOfFingerprintsRequiringFurtherTesting_3() uint32 {
	if m != nil {
		return m.NumberOfFingerprintsRequiringFurtherTesting_3
	}
	return 0
}

func (m *FingerprintsStat) GetNumberOfFingerprintsRequiringFurtherTesting_4() uint32 {
	if m != nil {
		return m.NumberOfFingerprintsRequiringFurtherTesting_4
	}
	return 0
}

func (m *FingerprintsStat) GetNumberOfFingerprintsRequiringFurtherTesting_5() uint32 {
	if m != nil {
		return m.NumberOfFingerprintsRequiringFurtherTesting_5
	}
	return 0
}

func (m *FingerprintsStat) GetNumberOfFingerprintsRequiringFurtherTesting_6() uint32 {
	if m != nil {
		return m.NumberOfFingerprintsRequiringFurtherTesting_6
	}
	return 0
}

func (m *FingerprintsStat) GetNumberOfFingerprintsOfSuspectedDupes() uint32 {
	if m != nil {
		return m.NumberOfFingerprintsOfSuspectedDupes
	}
	return 0
}

type Maxes struct {
	PearsonMax           float32  `protobuf:"fixed32,1,opt,name=pearson_max,json=pearsonMax,proto3" json:"pearson_max,omitempty"`
	SpearmanMax          float32  `protobuf:"fixed32,2,opt,name=spearman_max,json=spearmanMax,proto3" json:"spearman_max,omitempty"`
	KendallMax           float32  `protobuf:"fixed32,3,opt,name=kendall_max,json=kendallMax,proto3" json:"kendall_max,omitempty"`
	HoeffdingMax         float32  `protobuf:"fixed32,4,opt,name=hoeffding_max,json=hoeffdingMax,proto3" json:"hoeffding_max,omitempty"`
	MutualInformationMax float32  `protobuf:"fixed32,5,opt,name=mutual_information_max,json=mutualInformationMax,proto3" json:"mutual_information_max,omitempty"`
	HsicMax              float32  `protobuf:"fixed32,6,opt,name=hsic_max,json=hsicMax,proto3" json:"hsic_max,omitempty"`
	XgbimportanceMax     float32  `protobuf:"fixed32,7,opt,name=xgbimportance_max,json=xgbimportanceMax,proto3" json:"xgbimportance_max,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Maxes) Reset()         { *m = Maxes{} }
func (m *Maxes) String() string { return proto.CompactTextString(m) }
func (*Maxes) ProtoMessage()    {}
func (*Maxes) Descriptor() ([]byte, []int) {
	return fileDescriptor_1ddb247232398ef5, []int{7}
}

func (m *Maxes) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Maxes.Unmarshal(m, b)
}
func (m *Maxes) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Maxes.Marshal(b, m, deterministic)
}
func (m *Maxes) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Maxes.Merge(m, src)
}
func (m *Maxes) XXX_Size() int {
	return xxx_messageInfo_Maxes.Size(m)
}
func (m *Maxes) XXX_DiscardUnknown() {
	xxx_messageInfo_Maxes.DiscardUnknown(m)
}

var xxx_messageInfo_Maxes proto.InternalMessageInfo

func (m *Maxes) GetPearsonMax() float32 {
	if m != nil {
		return m.PearsonMax
	}
	return 0
}

func (m *Maxes) GetSpearmanMax() float32 {
	if m != nil {
		return m.SpearmanMax
	}
	return 0
}

func (m *Maxes) GetKendallMax() float32 {
	if m != nil {
		return m.KendallMax
	}
	return 0
}

func (m *Maxes) GetHoeffdingMax() float32 {
	if m != nil {
		return m.HoeffdingMax
	}
	return 0
}

func (m *Maxes) GetMutualInformationMax() float32 {
	if m != nil {
		return m.MutualInformationMax
	}
	return 0
}

func (m *Maxes) GetHsicMax() float32 {
	if m != nil {
		return m.HsicMax
	}
	return 0
}

func (m *Maxes) GetXgbimportanceMax() float32 {
	if m != nil {
		return m.XgbimportanceMax
	}
	return 0
}

type Percentile struct {
	PearsonTop_1BpsPercentile             float32  `protobuf:"fixed32,1,opt,name=pearson_top_1_bps_percentile,json=pearsonTop1BpsPercentile,proto3" json:"pearson_top_1_bps_percentile,omitempty"`
	SpearmanTop_1BpsPercentile            float32  `protobuf:"fixed32,2,opt,name=spearman_top_1_bps_percentile,json=spearmanTop1BpsPercentile,proto3" json:"spearman_top_1_bps_percentile,omitempty"`
	KendallTop_1BpsPercentile             float32  `protobuf:"fixed32,3,opt,name=kendall_top_1_bps_percentile,json=kendallTop1BpsPercentile,proto3" json:"kendall_top_1_bps_percentile,omitempty"`
	HoeffdingTop_10BpsPercentile          float32  `protobuf:"fixed32,4,opt,name=hoeffding_top_10_bps_percentile,json=hoeffdingTop10BpsPercentile,proto3" json:"hoeffding_top_10_bps_percentile,omitempty"`
	MutualInformationTop_100BpsPercentile float32  `protobuf:"fixed32,5,opt,name=mutual_information_top_100_bps_percentile,json=mutualInformationTop100BpsPercentile,proto3" json:"mutual_information_top_100_bps_percentile,omitempty"`
	HsicTop_100BpsPercentile              float32  `protobuf:"fixed32,6,opt,name=hsic_top_100_bps_percentile,json=hsicTop100BpsPercentile,proto3" json:"hsic_top_100_bps_percentile,omitempty"`
	XgbimportanceTop_100BpsPercentile     float32  `protobuf:"fixed32,7,opt,name=xgbimportance_top_100_bps_percentile,json=xgbimportanceTop100BpsPercentile,proto3" json:"xgbimportance_top_100_bps_percentile,omitempty"`
	XXX_NoUnkeyedLiteral                  struct{} `json:"-"`
	XXX_unrecognized                      []byte   `json:"-"`
	XXX_sizecache                         int32    `json:"-"`
}

func (m *Percentile) Reset()         { *m = Percentile{} }
func (m *Percentile) String() string { return proto.CompactTextString(m) }
func (*Percentile) ProtoMessage()    {}
func (*Percentile) Descriptor() ([]byte, []int) {
	return fileDescriptor_1ddb247232398ef5, []int{8}
}

func (m *Percentile) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Percentile.Unmarshal(m, b)
}
func (m *Percentile) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Percentile.Marshal(b, m, deterministic)
}
func (m *Percentile) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Percentile.Merge(m, src)
}
func (m *Percentile) XXX_Size() int {
	return xxx_messageInfo_Percentile.Size(m)
}
func (m *Percentile) XXX_DiscardUnknown() {
	xxx_messageInfo_Percentile.DiscardUnknown(m)
}

var xxx_messageInfo_Percentile proto.InternalMessageInfo

func (m *Percentile) GetPearsonTop_1BpsPercentile() float32 {
	if m != nil {
		return m.PearsonTop_1BpsPercentile
	}
	return 0
}

func (m *Percentile) GetSpearmanTop_1BpsPercentile() float32 {
	if m != nil {
		return m.SpearmanTop_1BpsPercentile
	}
	return 0
}

func (m *Percentile) GetKendallTop_1BpsPercentile() float32 {
	if m != nil {
		return m.KendallTop_1BpsPercentile
	}
	return 0
}

func (m *Percentile) GetHoeffdingTop_10BpsPercentile() float32 {
	if m != nil {
		return m.HoeffdingTop_10BpsPercentile
	}
	return 0
}

func (m *Percentile) GetMutualInformationTop_100BpsPercentile() float32 {
	if m != nil {
		return m.MutualInformationTop_100BpsPercentile
	}
	return 0
}

func (m *Percentile) GetHsicTop_100BpsPercentile() float32 {
	if m != nil {
		return m.HsicTop_100BpsPercentile
	}
	return 0
}

func (m *Percentile) GetXgbimportanceTop_100BpsPercentile() float32 {
	if m != nil {
		return m.XgbimportanceTop_100BpsPercentile
	}
	return 0
}

func init() {
	proto.RegisterType((*RarenessScoreRequest)(nil), "dupedetection.RarenessScoreRequest")
	proto.RegisterType((*ImageRarenessScoreReply)(nil), "dupedetection.ImageRarenessScoreReply")
	proto.RegisterType((*RarenessScores)(nil), "dupedetection.RarenessScores")
	proto.RegisterType((*InternetRareness)(nil), "dupedetection.InternetRareness")
	proto.RegisterType((*AltNsfwScores)(nil), "dupedetection.AltNsfwScores")
	proto.RegisterType((*PerceptualImageHashes)(nil), "dupedetection.PerceptualImageHashes")
	proto.RegisterType((*FingerprintsStat)(nil), "dupedetection.FingerprintsStat")
	proto.RegisterType((*Maxes)(nil), "dupedetection.Maxes")
	proto.RegisterType((*Percentile)(nil), "dupedetection.Percentile")
}

func init() { proto.RegisterFile("dd-server.proto", fileDescriptor_1ddb247232398ef5) }

var fileDescriptor_1ddb247232398ef5 = []byte{
	// 1388 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x57, 0xdd, 0x6e, 0xdb, 0x36,
	0x14, 0x86, 0x9d, 0xff, 0x93, 0x38, 0x49, 0xd9, 0xa4, 0x51, 0xfe, 0x16, 0xcf, 0x0d, 0xba, 0x74,
	0x5b, 0x93, 0x26, 0xfd, 0x01, 0x3a, 0x74, 0xc5, 0xda, 0x06, 0x59, 0x8d, 0x35, 0x75, 0xa1, 0x14,
	0x2b, 0x30, 0x0c, 0x20, 0x68, 0x8b, 0xb2, 0x85, 0xc8, 0x94, 0x42, 0x52, 0x89, 0xf3, 0x00, 0xbb,
	0xde, 0xd5, 0x9e, 0x69, 0xb7, 0x7b, 0x8b, 0xed, 0x1d, 0x76, 0x33, 0x1c, 0x52, 0xb2, 0x2d, 0xc5,
	0x19, 0x90, 0x5c, 0xd9, 0x3c, 0xe7, 0x3b, 0xdf, 0x47, 0x9e, 0x43, 0x1e, 0x8a, 0xb0, 0xe0, 0x79,
	0x8f, 0x14, 0x97, 0xe7, 0x5c, 0xee, 0xc6, 0x32, 0xd2, 0x11, 0xa9, 0x78, 0x49, 0xcc, 0x3d, 0xae,
	0x79, 0x4b, 0x07, 0x91, 0xa8, 0xf9, 0xb0, 0xe4, 0x32, 0xc9, 0x05, 0x57, 0xea, 0xa4, 0x15, 0x49,
	0xee, 0xf2, 0xb3, 0x84, 0x2b, 0x4d, 0x08, 0x8c, 0xc7, 0x4c, 0x77, 0x9c, 0x52, 0xb5, 0xb4, 0x33,
	0xe3, 0x9a, 0xff, 0x64, 0x13, 0xa0, 0x19, 0x46, 0xad, 0x53, 0xda, 0x61, 0xaa, 0xe3, 0x94, 0x8d,
	0x67, 0xc6, 0x58, 0xde, 0x31, 0xd5, 0x21, 0xeb, 0x30, 0x13, 0x33, 0xa5, 0x79, 0x48, 0x03, 0xcf,
	0x19, 0x33, 0xde, 0x69, 0x6b, 0xa8, 0x7b, 0xb5, 0x7f, 0xa6, 0x60, 0xa5, 0xde, 0x65, 0x6d, 0x5e,
	0x50, 0x8b, 0xc3, 0x4b, 0xb2, 0x04, 0x13, 0x86, 0x25, 0x15, 0xb3, 0x03, 0xb2, 0x01, 0x33, 0xb1,
	0x0c, 0x44, 0x2b, 0x88, 0x59, 0x98, 0x89, 0xf5, 0x0d, 0xe4, 0x35, 0x6c, 0xe2, 0x42, 0x68, 0x7f,
	0x25, 0x54, 0x5d, 0x2a, 0xcd, 0xbb, 0xf4, 0x9c, 0x4b, 0x15, 0x44, 0x22, 0x9d, 0xc0, 0x1a, 0x82,
	0x0e, 0x33, 0xcc, 0x89, 0x81, 0xfc, 0x6c, 0x11, 0x64, 0x1b, 0xe6, 0x03, 0x45, 0xc3, 0xe0, 0x94,
	0x87, 0x97, 0x14, 0x71, 0xce, 0x78, 0xb5, 0xb4, 0x33, 0xed, 0xce, 0x05, 0xea, 0xbd, 0x31, 0x1e,
	0x26, 0x31, 0x27, 0x8f, 0xe0, 0x6e, 0xa0, 0xa8, 0x64, 0x92, 0xd3, 0x48, 0xd0, 0x40, 0x68, 0x2e,
	0x05, 0xd7, 0xce, 0x84, 0x81, 0x2e, 0x06, 0x0a, 0xd7, 0xd3, 0x10, 0xf5, 0xd4, 0x4e, 0x8e, 0x60,
	0x41, 0xa6, 0x2b, 0xa4, 0x0a, 0x97, 0xa8, 0x9c, 0xc9, 0x6a, 0x69, 0x67, 0xf6, 0x60, 0x73, 0x37,
	0x97, 0xf8, 0xdd, 0x5c, 0x1e, 0x94, 0x3b, 0x2f, 0x73, 0x63, 0xf2, 0x1e, 0xee, 0x64, 0x5a, 0x34,
	0x73, 0x39, 0x53, 0x86, 0x69, 0xab, 0xc0, 0x94, 0x69, 0x67, 0x8c, 0xee, 0x62, 0x50, 0xb0, 0x90,
	0x07, 0xb0, 0x10, 0xc5, 0x5c, 0x50, 0xa1, 0xfc, 0x0b, 0x3b, 0x2d, 0x67, 0xba, 0x5a, 0xda, 0x29,
	0xbb, 0x15, 0x34, 0x7f, 0x50, 0xfe, 0x85, 0x91, 0x25, 0x9f, 0x60, 0x85, 0x85, 0x18, 0xcb, 0x74,
	0x70, 0xce, 0x87, 0xe0, 0xca, 0x99, 0x31, 0xda, 0x1b, 0x05, 0xed, 0xd7, 0xa1, 0xee, 0x47, 0x2b,
	0x77, 0x79, 0x28, 0x78, 0x60, 0x26, 0x9f, 0xe1, 0x61, 0x80, 0xa5, 0xa7, 0x7e, 0x20, 0xda, 0x5c,
	0x62, 0x11, 0x35, 0x8d, 0x7c, 0xda, 0x62, 0xc2, 0x0b, 0x3c, 0xa6, 0x39, 0xcd, 0xdc, 0x21, 0x77,
	0xa0, 0x3a, 0xb6, 0x53, 0x76, 0xb7, 0x8d, 0xe5, 0x68, 0x80, 0x6f, 0xf8, 0x6f, 0x33, 0x74, 0xdd,
	0xba, 0x42, 0x8e, 0x49, 0x1a, 0xa2, 0x54, 0x54, 0x69, 0xa6, 0x9d, 0xd9, 0x91, 0x49, 0x1a, 0xa2,
	0x52, 0x27, 0x9a, 0x69, 0x77, 0xd1, 0x2f, 0x58, 0xc8, 0x2b, 0xd8, 0xc0, 0x8d, 0x7d, 0xed, 0xcc,
	0xe6, 0xcc, 0x8e, 0x72, 0x10, 0x33, 0x72, 0x36, 0xbf, 0xc2, 0x4a, 0xcc, 0x65, 0x8b, 0xc7, 0x3a,
	0x61, 0x61, 0x1a, 0x88, 0x60, 0xae, 0x9c, 0x8a, 0x99, 0xd3, 0x76, 0x61, 0x4e, 0x1f, 0xfb, 0x68,
	0x43, 0xf2, 0xce, 0x60, 0xdd, 0xe5, 0x78, 0x94, 0x19, 0x37, 0xfc, 0x10, 0xbb, 0x9d, 0xe8, 0x39,
	0x97, 0x21, 0x8b, 0x69, 0x2b, 0x4a, 0x84, 0x76, 0xe6, 0xab, 0xa5, 0x9d, 0x8a, 0xbb, 0x36, 0x00,
	0x61, 0x60, 0xc3, 0x42, 0xde, 0x22, 0x82, 0x7c, 0x0d, 0x13, 0x5d, 0xd6, 0xe3, 0xca, 0x59, 0x30,
	0xd3, 0x59, 0x2a, 0x4c, 0xe7, 0x18, 0x7d, 0xae, 0x85, 0x90, 0x17, 0x00, 0x86, 0x49, 0x68, 0x5c,
	0xfa, 0xa2, 0x09, 0x58, 0x1d, 0x35, 0x7f, 0x03, 0x70, 0x87, 0xc0, 0xb5, 0x3f, 0xca, 0x30, 0x9f,
	0xdf, 0xdd, 0xe4, 0x39, 0xac, 0xb4, 0xa2, 0x6e, 0x33, 0x10, 0xdc, 0xa3, 0xf9, 0xe3, 0x61, 0xce,
	0x7c, 0xd9, 0x5d, 0xce, 0xdc, 0xb9, 0x40, 0xf2, 0x23, 0x54, 0x7b, 0xed, 0x66, 0x14, 0x29, 0x4d,
	0x63, 0xc9, 0xbd, 0xa0, 0xa5, 0xaf, 0x12, 0x94, 0x0d, 0xc1, 0x66, 0x8a, 0xfb, 0x98, 0xc1, 0xf2,
	0x44, 0xdf, 0xc3, 0xba, 0x10, 0xd7, 0x73, 0x8c, 0x19, 0x0e, 0x47, 0x88, 0x6b, 0xc2, 0xdf, 0xc2,
	0x17, 0x98, 0x6c, 0x16, 0x86, 0x94, 0xe1, 0x6f, 0x9b, 0x17, 0x19, 0xc6, 0x0d, 0xc3, 0x7a, 0x8a,
	0x7a, 0x6d, 0x41, 0x39, 0x92, 0xda, 0x9f, 0x25, 0x58, 0x2c, 0x9e, 0x55, 0xf2, 0x12, 0xd6, 0xbb,
	0x4c, 0xb7, 0x3a, 0x5c, 0x51, 0x3f, 0x4a, 0x84, 0x87, 0x4d, 0xc6, 0x0f, 0x24, 0x2e, 0x98, 0xb5,
	0x6d, 0x76, 0x2a, 0xee, 0x4a, 0x0a, 0x39, 0x42, 0x44, 0x43, 0x1c, 0xa1, 0xff, 0x23, 0x6b, 0x73,
	0xf2, 0x1d, 0xac, 0x89, 0xa4, 0xdb, 0xe4, 0x12, 0x37, 0x2d, 0x06, 0x28, 0xfc, 0x23, 0xb9, 0x4a,
	0x42, 0xad, 0x4c, 0x66, 0x2a, 0xee, 0x3d, 0x8b, 0x68, 0xf8, 0x18, 0xa1, 0x1a, 0xbe, 0x6b, 0xbd,
	0xe4, 0x05, 0xac, 0x25, 0x32, 0x44, 0xbc, 0xd5, 0x33, 0x1a, 0x34, 0x10, 0x56, 0xd8, 0xb6, 0xcf,
	0xe5, 0x44, 0x86, 0x0d, 0xdf, 0xe8, 0x1d, 0xa3, 0xbb, 0x2e, 0x90, 0xa4, 0xf6, 0x5b, 0x09, 0x2a,
	0xb9, 0x93, 0x4f, 0xd6, 0x60, 0xda, 0x93, 0xec, 0x22, 0x10, 0x6d, 0x95, 0x56, 0xb4, 0x3f, 0x26,
	0xf7, 0x60, 0xb2, 0xc3, 0x85, 0x66, 0x41, 0x5a, 0xaa, 0x74, 0x44, 0x1c, 0x98, 0x12, 0x3c, 0xd1,
	0x92, 0x85, 0x69, 0xfe, 0xb3, 0xa1, 0xb9, 0x7c, 0x22, 0x29, 0xd2, 0xa4, 0x9a, 0xff, 0x68, 0x53,
	0xbc, 0x77, 0x69, 0x1a, 0x6f, 0xd9, 0x35, 0xff, 0x6b, 0x7f, 0x95, 0x60, 0x79, 0xe4, 0x21, 0x22,
	0xab, 0x30, 0x1d, 0x7b, 0x67, 0xf6, 0xa2, 0xb2, 0xb7, 0xca, 0x54, 0xec, 0x9d, 0x99, 0x6b, 0xea,
	0x2b, 0x58, 0x28, 0x1c, 0xa4, 0xf4, 0x76, 0x99, 0xcf, 0x1f, 0x1d, 0xf2, 0x25, 0xcc, 0x65, 0xc5,
	0x36, 0x28, 0x9b, 0x92, 0xd9, 0xd4, 0x96, 0x71, 0x79, 0x81, 0xef, 0x73, 0xc9, 0x45, 0x2b, 0x45,
	0x8d, 0x5b, 0xae, 0x81, 0x39, 0x03, 0x0a, 0x9e, 0x48, 0x16, 0x9a, 0x83, 0x6b, 0x80, 0x13, 0x16,
	0x38, 0x30, 0x23, 0xb0, 0xf6, 0xf7, 0x04, 0x2c, 0x16, 0x7b, 0x15, 0xe9, 0xc0, 0xc1, 0xa0, 0xcc,
	0xb9, 0x8e, 0x27, 0xf9, 0x59, 0x12, 0xc8, 0x40, 0xb4, 0xa9, 0x9f, 0x48, 0xdd, 0xe1, 0x92, 0x6a,
	0xae, 0x34, 0x8e, 0xf7, 0xd3, 0xbd, 0xf3, 0x6d, 0x56, 0xfe, 0x61, 0x56, 0x37, 0x0b, 0x3b, 0xb2,
	0x51, 0x9f, 0x6c, 0xd0, 0xfe, 0xad, 0x94, 0x0e, 0xd2, 0x8d, 0x76, 0x13, 0xa5, 0x83, 0x5b, 0x29,
	0x3d, 0x31, 0x35, 0xb8, 0x99, 0xd2, 0x93, 0x5b, 0x29, 0x3d, 0x35, 0x75, 0xbc, 0x99, 0xd2, 0xd3,
	0x5b, 0x29, 0x3d, 0x33, 0x1b, 0xe1, 0x66, 0x4a, 0xcf, 0x6e, 0xa5, 0xf4, 0xdc, 0x7c, 0x89, 0xdc,
	0x4c, 0xe9, 0x39, 0x5e, 0xde, 0xd7, 0x28, 0x45, 0x3e, 0x55, 0x89, 0x8a, 0xb9, 0x69, 0xa8, 0x78,
	0x33, 0xd8, 0x0f, 0x94, 0x8a, 0xbb, 0x3d, 0x4a, 0xa0, 0xe1, 0x9f, 0x64, 0x60, 0xfc, 0xae, 0x52,
	0xb5, 0xdf, 0xcb, 0x30, 0x61, 0xae, 0x1c, 0xb2, 0x05, 0xb3, 0x31, 0x67, 0x52, 0x45, 0x82, 0x76,
	0x59, 0x2f, 0xed, 0x1f, 0x90, 0x9a, 0x8e, 0x59, 0x0f, 0x4f, 0xa2, 0xc2, 0x61, 0x97, 0x59, 0x84,
	0xed, 0x23, 0xb3, 0x99, 0x0d, 0x21, 0x5b, 0x30, 0x7b, 0xca, 0x85, 0x87, 0x1d, 0x1a, 0x11, 0xb6,
	0xa1, 0x40, 0x6a, 0x42, 0xc0, 0x7d, 0xa8, 0x74, 0x22, 0xee, 0xfb, 0x1e, 0xa6, 0x02, 0x21, 0xb6,
	0xb9, 0xcc, 0xf5, 0x8d, 0x08, 0x7a, 0x0a, 0xf7, 0xba, 0x89, 0xbd, 0xbe, 0x85, 0x1f, 0xc9, 0x2e,
	0x33, 0x5f, 0x96, 0x88, 0xb6, 0x6d, 0x67, 0xc9, 0x7a, 0xeb, 0x03, 0x27, 0x46, 0xad, 0xc2, 0x74,
	0x47, 0x05, 0x2d, 0x83, 0x9b, 0xb4, 0x9d, 0x0c, 0xc7, 0xe8, 0xfa, 0x06, 0xee, 0xf4, 0xda, 0xcd,
	0xa0, 0x1b, 0x47, 0x52, 0x33, 0xec, 0x11, 0x88, 0x99, 0x32, 0x98, 0xc5, 0x9c, 0xe3, 0x98, 0xf5,
	0x6a, 0xff, 0x8e, 0x01, 0x0c, 0xee, 0x54, 0xfc, 0x1e, 0xc9, 0xd2, 0xa2, 0xa3, 0x98, 0xee, 0xd3,
	0x66, 0xac, 0xe8, 0xd0, 0xa5, 0x6c, 0xf3, 0xe4, 0xa4, 0x98, 0x4f, 0x51, 0xbc, 0xff, 0x26, 0x56,
	0x43, 0xf1, 0x3f, 0xc0, 0x66, 0x3f, 0x6b, 0x23, 0x09, 0x6c, 0x1a, 0x57, 0x33, 0xd0, 0x55, 0x86,
	0x57, 0xb0, 0x91, 0x25, 0x75, 0x24, 0x41, 0x7a, 0x6d, 0xa6, 0x98, 0xab, 0xf1, 0x87, 0xb0, 0x35,
	0xc8, 0xb9, 0x61, 0x78, 0x5c, 0xa4, 0x48, 0xef, 0xcd, 0x3e, 0x0c, 0x49, 0x1e, 0xe7, 0x59, 0x3e,
	0xc3, 0xc3, 0x11, 0x45, 0xb1, 0x74, 0x57, 0xf8, 0x6c, 0x9d, 0xb6, 0xaf, 0xd4, 0xc9, 0xf0, 0x16,
	0x88, 0x5f, 0xc2, 0xba, 0xa9, 0xdb, 0x35, 0x54, 0xb6, 0x94, 0x2b, 0x08, 0x19, 0x15, 0xfd, 0x01,
	0xb6, 0xf3, 0xa5, 0xbd, 0x86, 0xc6, 0x56, 0xbb, 0x9a, 0xc3, 0x8e, 0xe0, 0x3b, 0xe8, 0xc1, 0xdd,
	0xc3, 0xdc, 0x63, 0xc5, 0xbc, 0xda, 0x08, 0x03, 0x72, 0xf5, 0xdd, 0x44, 0xee, 0xff, 0xdf, 0x6b,
	0x22, 0x7d, 0xc3, 0xad, 0x3d, 0x28, 0x3e, 0x14, 0x46, 0xbf, 0xbf, 0xde, 0xfc, 0xf4, 0x4b, 0xbd,
	0x1d, 0xe8, 0x4e, 0xd2, 0xdc, 0x6d, 0x45, 0xdd, 0x3d, 0xfb, 0x64, 0x13, 0x5c, 0x5f, 0x44, 0xf2,
	0x74, 0xaf, 0x1d, 0x89, 0xc8, 0xe3, 0x7b, 0xf8, 0xa2, 0x68, 0x4b, 0x93, 0xc9, 0x3d, 0x9f, 0x9d,
	0x72, 0xb5, 0xd7, 0x7f, 0x57, 0xee, 0xe5, 0x44, 0x9a, 0x93, 0xe6, 0x99, 0xf9, 0xe4, 0xbf, 0x00,
	0x00, 0x00, 0xff, 0xff, 0x9c, 0xed, 0x4f, 0xa9, 0x79, 0x0e, 0x00, 0x00,
}