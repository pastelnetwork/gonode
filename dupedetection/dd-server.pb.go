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
	ImageFilepath                         string   `protobuf:"bytes,1,opt,name=image_filepath,json=imageFilepath,proto3" json:"image_filepath,omitempty"`
	PastelBlockHashWhenRequestSubmitted   string   `protobuf:"bytes,2,opt,name=pastel_block_hash_when_request_submitted,json=pastelBlockHashWhenRequestSubmitted,proto3" json:"pastel_block_hash_when_request_submitted,omitempty"`
	PastelBlockHeightWhenRequestSubmitted string   `protobuf:"bytes,3,opt,name=pastel_block_height_when_request_submitted,json=pastelBlockHeightWhenRequestSubmitted,proto3" json:"pastel_block_height_when_request_submitted,omitempty"`
	UtcTimestampWhenRequestSubmitted      string   `protobuf:"bytes,4,opt,name=utc_timestamp_when_request_submitted,json=utcTimestampWhenRequestSubmitted,proto3" json:"utc_timestamp_when_request_submitted,omitempty"`
	PastelIdOfSubmitter                   string   `protobuf:"bytes,5,opt,name=pastel_id_of_submitter,json=pastelIdOfSubmitter,proto3" json:"pastel_id_of_submitter,omitempty"`
	PastelIdOfRegisteringSupernode_1      string   `protobuf:"bytes,6,opt,name=pastel_id_of_registering_supernode_1,json=pastelIdOfRegisteringSupernode1,proto3" json:"pastel_id_of_registering_supernode_1,omitempty"`
	PastelIdOfRegisteringSupernode_2      string   `protobuf:"bytes,7,opt,name=pastel_id_of_registering_supernode_2,json=pastelIdOfRegisteringSupernode2,proto3" json:"pastel_id_of_registering_supernode_2,omitempty"`
	PastelIdOfRegisteringSupernode_3      string   `protobuf:"bytes,8,opt,name=pastel_id_of_registering_supernode_3,json=pastelIdOfRegisteringSupernode3,proto3" json:"pastel_id_of_registering_supernode_3,omitempty"`
	IsPastelOpenapiRequest                bool     `protobuf:"varint,9,opt,name=is_pastel_openapi_request,json=isPastelOpenapiRequest,proto3" json:"is_pastel_openapi_request,omitempty"`
	OpenApiSubsetIdString                 string   `protobuf:"bytes,10,opt,name=open_api_subset_id_string,json=openApiSubsetIdString,proto3" json:"open_api_subset_id_string,omitempty"`
	OpenApiGroupIdString                  string   `protobuf:"bytes,11,opt,name=open_api_group_id_string,json=openApiGroupIdString,proto3" json:"open_api_group_id_string,omitempty"`
	CollectionNameString                  string   `protobuf:"bytes,12,opt,name=collection_name_string,json=collectionNameString,proto3" json:"collection_name_string,omitempty"`
	XXX_NoUnkeyedLiteral                  struct{} `json:"-"`
	XXX_unrecognized                      []byte   `json:"-"`
	XXX_sizecache                         int32    `json:"-"`
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

func (m *RarenessScoreRequest) GetImageFilepath() string {
	if m != nil {
		return m.ImageFilepath
	}
	return ""
}

func (m *RarenessScoreRequest) GetPastelBlockHashWhenRequestSubmitted() string {
	if m != nil {
		return m.PastelBlockHashWhenRequestSubmitted
	}
	return ""
}

func (m *RarenessScoreRequest) GetPastelBlockHeightWhenRequestSubmitted() string {
	if m != nil {
		return m.PastelBlockHeightWhenRequestSubmitted
	}
	return ""
}

func (m *RarenessScoreRequest) GetUtcTimestampWhenRequestSubmitted() string {
	if m != nil {
		return m.UtcTimestampWhenRequestSubmitted
	}
	return ""
}

func (m *RarenessScoreRequest) GetPastelIdOfSubmitter() string {
	if m != nil {
		return m.PastelIdOfSubmitter
	}
	return ""
}

func (m *RarenessScoreRequest) GetPastelIdOfRegisteringSupernode_1() string {
	if m != nil {
		return m.PastelIdOfRegisteringSupernode_1
	}
	return ""
}

func (m *RarenessScoreRequest) GetPastelIdOfRegisteringSupernode_2() string {
	if m != nil {
		return m.PastelIdOfRegisteringSupernode_2
	}
	return ""
}

func (m *RarenessScoreRequest) GetPastelIdOfRegisteringSupernode_3() string {
	if m != nil {
		return m.PastelIdOfRegisteringSupernode_3
	}
	return ""
}

func (m *RarenessScoreRequest) GetIsPastelOpenapiRequest() bool {
	if m != nil {
		return m.IsPastelOpenapiRequest
	}
	return false
}

func (m *RarenessScoreRequest) GetOpenApiSubsetIdString() string {
	if m != nil {
		return m.OpenApiSubsetIdString
	}
	return ""
}

func (m *RarenessScoreRequest) GetOpenApiGroupIdString() string {
	if m != nil {
		return m.OpenApiGroupIdString
	}
	return ""
}

func (m *RarenessScoreRequest) GetCollectionNameString() string {
	if m != nil {
		return m.CollectionNameString
	}
	return ""
}

type ImageRarenessScoreReply struct {
	PastelBlockHashWhenRequestSubmitted           string            `protobuf:"bytes,1,opt,name=pastel_block_hash_when_request_submitted,json=pastelBlockHashWhenRequestSubmitted,proto3" json:"pastel_block_hash_when_request_submitted,omitempty"`
	PastelBlockHeightWhenRequestSubmitted         string            `protobuf:"bytes,2,opt,name=pastel_block_height_when_request_submitted,json=pastelBlockHeightWhenRequestSubmitted,proto3" json:"pastel_block_height_when_request_submitted,omitempty"`
	UtcTimestampWhenRequestSubmitted              string            `protobuf:"bytes,3,opt,name=utc_timestamp_when_request_submitted,json=utcTimestampWhenRequestSubmitted,proto3" json:"utc_timestamp_when_request_submitted,omitempty"`
	PastelIdOfSubmitter                           string            `protobuf:"bytes,4,opt,name=pastel_id_of_submitter,json=pastelIdOfSubmitter,proto3" json:"pastel_id_of_submitter,omitempty"`
	PastelIdOfRegisteringSupernode_1              string            `protobuf:"bytes,5,opt,name=pastel_id_of_registering_supernode_1,json=pastelIdOfRegisteringSupernode1,proto3" json:"pastel_id_of_registering_supernode_1,omitempty"`
	PastelIdOfRegisteringSupernode_2              string            `protobuf:"bytes,6,opt,name=pastel_id_of_registering_supernode_2,json=pastelIdOfRegisteringSupernode2,proto3" json:"pastel_id_of_registering_supernode_2,omitempty"`
	PastelIdOfRegisteringSupernode_3              string            `protobuf:"bytes,7,opt,name=pastel_id_of_registering_supernode_3,json=pastelIdOfRegisteringSupernode3,proto3" json:"pastel_id_of_registering_supernode_3,omitempty"`
	IsPastelOpenapiRequest                        bool              `protobuf:"varint,8,opt,name=is_pastel_openapi_request,json=isPastelOpenapiRequest,proto3" json:"is_pastel_openapi_request,omitempty"`
	OpenApiSubsetIdString                         string            `protobuf:"bytes,9,opt,name=open_api_subset_id_string,json=openApiSubsetIdString,proto3" json:"open_api_subset_id_string,omitempty"`
	DupeDetectionSystemVersion                    string            `protobuf:"bytes,10,opt,name=dupe_detection_system_version,json=dupeDetectionSystemVersion,proto3" json:"dupe_detection_system_version,omitempty"`
	IsLikelyDupe                                  bool              `protobuf:"varint,11,opt,name=is_likely_dupe,json=isLikelyDupe,proto3" json:"is_likely_dupe,omitempty"`
	IsRareOnInternet                              bool              `protobuf:"varint,12,opt,name=is_rare_on_internet,json=isRareOnInternet,proto3" json:"is_rare_on_internet,omitempty"`
	OverallRarenessScore                          float32           `protobuf:"fixed32,13,opt,name=overall_rareness_score,json=overallRarenessScore,proto3" json:"overall_rareness_score,omitempty"`
	PctOfTop_10MostSimilarWithDupeProbAbove_25Pct float32           `protobuf:"fixed32,14,opt,name=pct_of_top_10_most_similar_with_dupe_prob_above_25pct,json=pctOfTop10MostSimilarWithDupeProbAbove25pct,proto3" json:"pct_of_top_10_most_similar_with_dupe_prob_above_25pct,omitempty"`
	PctOfTop_10MostSimilarWithDupeProbAbove_33Pct float32           `protobuf:"fixed32,15,opt,name=pct_of_top_10_most_similar_with_dupe_prob_above_33pct,json=pctOfTop10MostSimilarWithDupeProbAbove33pct,proto3" json:"pct_of_top_10_most_similar_with_dupe_prob_above_33pct,omitempty"`
	PctOfTop_10MostSimilarWithDupeProbAbove_50Pct float32           `protobuf:"fixed32,16,opt,name=pct_of_top_10_most_similar_with_dupe_prob_above_50pct,json=pctOfTop10MostSimilarWithDupeProbAbove50pct,proto3" json:"pct_of_top_10_most_similar_with_dupe_prob_above_50pct,omitempty"`
	RarenessScoresTableJsonCompressedB64          string            `protobuf:"bytes,17,opt,name=rareness_scores_table_json_compressed_b64,json=rarenessScoresTableJsonCompressedB64,proto3" json:"rareness_scores_table_json_compressed_b64,omitempty"`
	InternetRareness                              *InternetRareness `protobuf:"bytes,18,opt,name=internet_rareness,json=internetRareness,proto3" json:"internet_rareness,omitempty"`
	OpenNsfwScore                                 float32           `protobuf:"fixed32,19,opt,name=open_nsfw_score,json=openNsfwScore,proto3" json:"open_nsfw_score,omitempty"`
	AlternativeNsfwScores                         *AltNsfwScores    `protobuf:"bytes,20,opt,name=alternative_nsfw_scores,json=alternativeNsfwScores,proto3" json:"alternative_nsfw_scores,omitempty"`
	ImageFingerprintOfCandidateImageFile          []float64         `protobuf:"fixed64,21,rep,packed,name=image_fingerprint_of_candidate_image_file,json=imageFingerprintOfCandidateImageFile,proto3" json:"image_fingerprint_of_candidate_image_file,omitempty"`
	CollectionNameString                          string            `protobuf:"bytes,22,opt,name=collection_name_string,json=collectionNameString,proto3" json:"collection_name_string,omitempty"`
	HashOfCandidateImageFile                      string            `protobuf:"bytes,23,opt,name=hash_of_candidate_image_file,json=hashOfCandidateImageFile,proto3" json:"hash_of_candidate_image_file,omitempty"`
	OpenApiGroupIdString                          string            `protobuf:"bytes,24,opt,name=open_api_group_id_string,json=openApiGroupIdString,proto3" json:"open_api_group_id_string,omitempty"`
	GroupRarenessScore                            float32           `protobuf:"fixed32,25,opt,name=group_rareness_score,json=groupRarenessScore,proto3" json:"group_rareness_score,omitempty"`
	CandidateImageThumbnailWebpAsBase64String     string            `protobuf:"bytes,26,opt,name=candidate_image_thumbnail_webp_as_base64_string,json=candidateImageThumbnailWebpAsBase64String,proto3" json:"candidate_image_thumbnail_webp_as_base64_string,omitempty"`
	DoesNotImpactTheFollowingCollectionStrings    string            `protobuf:"bytes,27,opt,name=does_not_impact_the_following_collection_strings,json=doesNotImpactTheFollowingCollectionStrings,proto3" json:"does_not_impact_the_following_collection_strings,omitempty"`
	IsInvalidSenseRequest                         bool              `protobuf:"varint,28,opt,name=is_invalid_sense_request,json=isInvalidSenseRequest,proto3" json:"is_invalid_sense_request,omitempty"`
	InvalidSenseRequestReason                     string            `protobuf:"bytes,29,opt,name=invalid_sense_request_reason,json=invalidSenseRequestReason,proto3" json:"invalid_sense_request_reason,omitempty"`
	SimilarityScoreToFirstEntryInCollection       float32           `protobuf:"fixed32,30,opt,name=similarity_score_to_first_entry_in_collection,json=similarityScoreToFirstEntryInCollection,proto3" json:"similarity_score_to_first_entry_in_collection,omitempty"`
	CpProbability                                 float32           `protobuf:"fixed32,31,opt,name=cp_probability,json=cpProbability,proto3" json:"cp_probability,omitempty"`
	ChildProbability                              float32           `protobuf:"fixed32,32,opt,name=child_probability,json=childProbability,proto3" json:"child_probability,omitempty"`
	ImageFilePath                                 string            `protobuf:"bytes,33,opt,name=image_file_path,json=imageFilePath,proto3" json:"image_file_path,omitempty"`
	XXX_NoUnkeyedLiteral                          struct{}          `json:"-"`
	XXX_unrecognized                              []byte            `json:"-"`
	XXX_sizecache                                 int32             `json:"-"`
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

func (m *ImageRarenessScoreReply) GetPastelBlockHashWhenRequestSubmitted() string {
	if m != nil {
		return m.PastelBlockHashWhenRequestSubmitted
	}
	return ""
}

func (m *ImageRarenessScoreReply) GetPastelBlockHeightWhenRequestSubmitted() string {
	if m != nil {
		return m.PastelBlockHeightWhenRequestSubmitted
	}
	return ""
}

func (m *ImageRarenessScoreReply) GetUtcTimestampWhenRequestSubmitted() string {
	if m != nil {
		return m.UtcTimestampWhenRequestSubmitted
	}
	return ""
}

func (m *ImageRarenessScoreReply) GetPastelIdOfSubmitter() string {
	if m != nil {
		return m.PastelIdOfSubmitter
	}
	return ""
}

func (m *ImageRarenessScoreReply) GetPastelIdOfRegisteringSupernode_1() string {
	if m != nil {
		return m.PastelIdOfRegisteringSupernode_1
	}
	return ""
}

func (m *ImageRarenessScoreReply) GetPastelIdOfRegisteringSupernode_2() string {
	if m != nil {
		return m.PastelIdOfRegisteringSupernode_2
	}
	return ""
}

func (m *ImageRarenessScoreReply) GetPastelIdOfRegisteringSupernode_3() string {
	if m != nil {
		return m.PastelIdOfRegisteringSupernode_3
	}
	return ""
}

func (m *ImageRarenessScoreReply) GetIsPastelOpenapiRequest() bool {
	if m != nil {
		return m.IsPastelOpenapiRequest
	}
	return false
}

func (m *ImageRarenessScoreReply) GetOpenApiSubsetIdString() string {
	if m != nil {
		return m.OpenApiSubsetIdString
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

func (m *ImageRarenessScoreReply) GetOverallRarenessScore() float32 {
	if m != nil {
		return m.OverallRarenessScore
	}
	return 0
}

func (m *ImageRarenessScoreReply) GetPctOfTop_10MostSimilarWithDupeProbAbove_25Pct() float32 {
	if m != nil {
		return m.PctOfTop_10MostSimilarWithDupeProbAbove_25Pct
	}
	return 0
}

func (m *ImageRarenessScoreReply) GetPctOfTop_10MostSimilarWithDupeProbAbove_33Pct() float32 {
	if m != nil {
		return m.PctOfTop_10MostSimilarWithDupeProbAbove_33Pct
	}
	return 0
}

func (m *ImageRarenessScoreReply) GetPctOfTop_10MostSimilarWithDupeProbAbove_50Pct() float32 {
	if m != nil {
		return m.PctOfTop_10MostSimilarWithDupeProbAbove_50Pct
	}
	return 0
}

func (m *ImageRarenessScoreReply) GetRarenessScoresTableJsonCompressedB64() string {
	if m != nil {
		return m.RarenessScoresTableJsonCompressedB64
	}
	return ""
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

func (m *ImageRarenessScoreReply) GetImageFingerprintOfCandidateImageFile() []float64 {
	if m != nil {
		return m.ImageFingerprintOfCandidateImageFile
	}
	return nil
}

func (m *ImageRarenessScoreReply) GetCollectionNameString() string {
	if m != nil {
		return m.CollectionNameString
	}
	return ""
}

func (m *ImageRarenessScoreReply) GetHashOfCandidateImageFile() string {
	if m != nil {
		return m.HashOfCandidateImageFile
	}
	return ""
}

func (m *ImageRarenessScoreReply) GetOpenApiGroupIdString() string {
	if m != nil {
		return m.OpenApiGroupIdString
	}
	return ""
}

func (m *ImageRarenessScoreReply) GetGroupRarenessScore() float32 {
	if m != nil {
		return m.GroupRarenessScore
	}
	return 0
}

func (m *ImageRarenessScoreReply) GetCandidateImageThumbnailWebpAsBase64String() string {
	if m != nil {
		return m.CandidateImageThumbnailWebpAsBase64String
	}
	return ""
}

func (m *ImageRarenessScoreReply) GetDoesNotImpactTheFollowingCollectionStrings() string {
	if m != nil {
		return m.DoesNotImpactTheFollowingCollectionStrings
	}
	return ""
}

func (m *ImageRarenessScoreReply) GetIsInvalidSenseRequest() bool {
	if m != nil {
		return m.IsInvalidSenseRequest
	}
	return false
}

func (m *ImageRarenessScoreReply) GetInvalidSenseRequestReason() string {
	if m != nil {
		return m.InvalidSenseRequestReason
	}
	return ""
}

func (m *ImageRarenessScoreReply) GetSimilarityScoreToFirstEntryInCollection() float32 {
	if m != nil {
		return m.SimilarityScoreToFirstEntryInCollection
	}
	return 0
}

func (m *ImageRarenessScoreReply) GetCpProbability() float32 {
	if m != nil {
		return m.CpProbability
	}
	return 0
}

func (m *ImageRarenessScoreReply) GetChildProbability() float32 {
	if m != nil {
		return m.ChildProbability
	}
	return 0
}

func (m *ImageRarenessScoreReply) GetImageFilePath() string {
	if m != nil {
		return m.ImageFilePath
	}
	return ""
}

type InternetRareness struct {
	RareOnInternetSummaryTableAsJsonCompressedB64    string   `protobuf:"bytes,1,opt,name=rare_on_internet_summary_table_as_json_compressed_b64,json=rareOnInternetSummaryTableAsJsonCompressedB64,proto3" json:"rare_on_internet_summary_table_as_json_compressed_b64,omitempty"`
	RareOnInternetGraphJsonCompressedB64             string   `protobuf:"bytes,2,opt,name=rare_on_internet_graph_json_compressed_b64,json=rareOnInternetGraphJsonCompressedB64,proto3" json:"rare_on_internet_graph_json_compressed_b64,omitempty"`
	AlternativeRareOnInternetDictAsJsonCompressedB64 string   `protobuf:"bytes,3,opt,name=alternative_rare_on_internet_dict_as_json_compressed_b64,json=alternativeRareOnInternetDictAsJsonCompressedB64,proto3" json:"alternative_rare_on_internet_dict_as_json_compressed_b64,omitempty"`
	MinNumberOfExactMatchesInPage                    uint32   `protobuf:"varint,4,opt,name=min_number_of_exact_matches_in_page,json=minNumberOfExactMatchesInPage,proto3" json:"min_number_of_exact_matches_in_page,omitempty"`
	EarliestAvailableDateOfInternetResults           string   `protobuf:"bytes,5,opt,name=earliest_available_date_of_internet_results,json=earliestAvailableDateOfInternetResults,proto3" json:"earliest_available_date_of_internet_results,omitempty"`
	XXX_NoUnkeyedLiteral                             struct{} `json:"-"`
	XXX_unrecognized                                 []byte   `json:"-"`
	XXX_sizecache                                    int32    `json:"-"`
}

func (m *InternetRareness) Reset()         { *m = InternetRareness{} }
func (m *InternetRareness) String() string { return proto.CompactTextString(m) }
func (*InternetRareness) ProtoMessage()    {}
func (*InternetRareness) Descriptor() ([]byte, []int) {
	return fileDescriptor_1ddb247232398ef5, []int{2}
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

func (m *InternetRareness) GetRareOnInternetSummaryTableAsJsonCompressedB64() string {
	if m != nil {
		return m.RareOnInternetSummaryTableAsJsonCompressedB64
	}
	return ""
}

func (m *InternetRareness) GetRareOnInternetGraphJsonCompressedB64() string {
	if m != nil {
		return m.RareOnInternetGraphJsonCompressedB64
	}
	return ""
}

func (m *InternetRareness) GetAlternativeRareOnInternetDictAsJsonCompressedB64() string {
	if m != nil {
		return m.AlternativeRareOnInternetDictAsJsonCompressedB64
	}
	return ""
}

func (m *InternetRareness) GetMinNumberOfExactMatchesInPage() uint32 {
	if m != nil {
		return m.MinNumberOfExactMatchesInPage
	}
	return 0
}

func (m *InternetRareness) GetEarliestAvailableDateOfInternetResults() string {
	if m != nil {
		return m.EarliestAvailableDateOfInternetResults
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
	return fileDescriptor_1ddb247232398ef5, []int{3}
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

func init() {
	proto.RegisterType((*RarenessScoreRequest)(nil), "dupedetection.RarenessScoreRequest")
	proto.RegisterType((*ImageRarenessScoreReply)(nil), "dupedetection.ImageRarenessScoreReply")
	proto.RegisterType((*InternetRareness)(nil), "dupedetection.InternetRareness")
	proto.RegisterType((*AltNsfwScores)(nil), "dupedetection.AltNsfwScores")
}

func init() { proto.RegisterFile("dd-server.proto", fileDescriptor_1ddb247232398ef5) }

var fileDescriptor_1ddb247232398ef5 = []byte{
	// 1427 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xb4, 0x57, 0x7b, 0x53, 0x14, 0x47,
	0x10, 0xaf, 0x3b, 0x14, 0x71, 0x14, 0xc1, 0xe5, 0xe1, 0x42, 0x30, 0x12, 0x44, 0x82, 0x5a, 0xbc,
	0x91, 0x98, 0x7f, 0x92, 0x02, 0x51, 0x73, 0x96, 0x02, 0xb5, 0x47, 0x42, 0x1e, 0x55, 0x99, 0x9a,
	0xdd, 0xed, 0xbb, 0x1d, 0xdd, 0x9d, 0xd9, 0xcc, 0xcc, 0x71, 0xde, 0x07, 0xc8, 0x17, 0x4b, 0xe5,
	0x3b, 0xe5, 0xdf, 0xd4, 0xf4, 0xee, 0xde, 0x8b, 0x83, 0xc2, 0xd7, 0x7f, 0x77, 0xd3, 0xbf, 0x5f,
	0xf7, 0x3c, 0xfa, 0xd7, 0xdd, 0x4b, 0xc6, 0xc2, 0x70, 0x45, 0x83, 0x3a, 0x05, 0xb5, 0x9a, 0x2a,
	0x69, 0xa4, 0x33, 0x1a, 0x36, 0x52, 0x08, 0xc1, 0x40, 0x60, 0xb8, 0x14, 0x0b, 0xff, 0x0e, 0x93,
	0x49, 0x8f, 0x29, 0x10, 0xa0, 0x75, 0x35, 0x90, 0x0a, 0x3c, 0xf8, 0xab, 0x01, 0xda, 0x38, 0x0f,
	0xc8, 0x2d, 0x9e, 0xb0, 0x3a, 0xd0, 0x1a, 0x8f, 0x21, 0x65, 0x26, 0x72, 0x4b, 0xf3, 0xa5, 0xe5,
	0xeb, 0xde, 0x28, 0xae, 0xbe, 0xc8, 0x17, 0x9d, 0x9f, 0xc9, 0x72, 0xca, 0xb4, 0x81, 0x98, 0xfa,
	0xb1, 0x0c, 0xde, 0xd1, 0x88, 0xe9, 0x88, 0x36, 0x23, 0x10, 0x54, 0x65, 0x8e, 0xa8, 0x6e, 0xf8,
	0x09, 0x37, 0x06, 0x42, 0xb7, 0x8c, 0x0e, 0xee, 0x67, 0xf8, 0x3d, 0x0b, 0xff, 0x89, 0xe9, 0xe8,
	0x24, 0x02, 0x91, 0x07, 0xad, 0x16, 0x50, 0xe7, 0x37, 0xf2, 0xa8, 0xd7, 0x2d, 0xf0, 0x7a, 0x64,
	0xce, 0x73, 0x3c, 0x84, 0x8e, 0x1f, 0x74, 0x3b, 0x46, 0xfc, 0x40, 0xd7, 0x07, 0x64, 0xb1, 0x61,
	0x02, 0x6a, 0x78, 0x02, 0xda, 0xb0, 0x24, 0x3d, 0xcf, 0xe9, 0x15, 0x74, 0x3a, 0xdf, 0x30, 0xc1,
	0x71, 0x01, 0x1d, 0xe8, 0x6f, 0x8b, 0x4c, 0xe7, 0x5b, 0xe5, 0x21, 0x95, 0xb5, 0xb6, 0x07, 0xe5,
	0x5e, 0x45, 0x0f, 0x13, 0x99, 0xb5, 0x12, 0x1e, 0xd6, 0x0a, 0x92, 0x72, 0xde, 0x90, 0xc5, 0x1e,
	0x92, 0x82, 0x3a, 0xd7, 0x06, 0x14, 0x17, 0x75, 0xaa, 0x1b, 0x29, 0x28, 0x21, 0x43, 0xa0, 0x1b,
	0xee, 0x30, 0xba, 0xb8, 0xd7, 0x71, 0xe1, 0x75, 0x80, 0xd5, 0x02, 0xb7, 0x71, 0x49, 0x77, 0x9b,
	0xee, 0xb5, 0xcb, 0xb8, 0xdb, 0xbc, 0xa4, 0xbb, 0x2d, 0x77, 0xe4, 0x32, 0xee, 0xb6, 0x9c, 0xef,
	0xc9, 0x0c, 0xd7, 0x34, 0xf7, 0x28, 0x53, 0x10, 0x2c, 0xe5, 0xc5, 0x85, 0xbb, 0xd7, 0xe7, 0x4b,
	0xcb, 0x23, 0xde, 0x34, 0xd7, 0x47, 0x68, 0x3f, 0xcc, 0xcc, 0x45, 0x16, 0x3e, 0x25, 0x33, 0x96,
	0x40, 0x2d, 0x43, 0x37, 0x7c, 0x0d, 0xc6, 0x6e, 0x49, 0x1b, 0x1b, 0xc1, 0x25, 0x18, 0x7e, 0xca,
	0x02, 0x76, 0x53, 0x5e, 0x45, 0x73, 0x25, 0xac, 0xa2, 0xd1, 0xd9, 0x21, 0x6e, 0x9b, 0x59, 0x57,
	0xb2, 0x91, 0x76, 0x11, 0x6f, 0x20, 0x71, 0x32, 0x27, 0xbe, 0xb4, 0xd6, 0x36, 0x6f, 0x9b, 0x4c,
	0x07, 0x32, 0x8e, 0x33, 0x79, 0x50, 0xc1, 0x12, 0x28, 0x58, 0x37, 0x33, 0x56, 0xc7, 0x7a, 0xc0,
	0x12, 0xc8, 0x58, 0x0b, 0xff, 0x38, 0xe4, 0x4e, 0xc5, 0x0a, 0xa3, 0x4f, 0x4b, 0x69, 0xdc, 0xfa,
	0x20, 0x89, 0x94, 0xbe, 0x94, 0x44, 0xca, 0x5f, 0x42, 0x22, 0x43, 0x9f, 0x2c, 0x91, 0x2b, 0x9f,
	0x2e, 0x91, 0xab, 0x9f, 0x57, 0x22, 0xc3, 0x9f, 0x57, 0x22, 0xd7, 0x3e, 0x83, 0x44, 0x46, 0x3e,
	0x5e, 0x22, 0xd7, 0x2f, 0x92, 0xc8, 0x2e, 0xb9, 0x6b, 0x9b, 0x01, 0x6d, 0x77, 0x03, 0xaa, 0x5b,
	0xda, 0x40, 0x42, 0x4f, 0x41, 0x69, 0x2e, 0x45, 0x2e, 0xb0, 0x59, 0x0b, 0xda, 0x2f, 0x30, 0x55,
	0x84, 0xfc, 0x92, 0x21, 0x9c, 0x45, 0x72, 0x8b, 0x6b, 0x1a, 0xf3, 0x77, 0x10, 0xb7, 0xa8, 0xc5,
	0xa1, 0xb6, 0x46, 0xbc, 0x9b, 0x5c, 0xbf, 0xc6, 0xc5, 0xfd, 0x46, 0x0a, 0xce, 0x0a, 0x99, 0xe0,
	0x9a, 0x2a, 0xa6, 0x80, 0x4a, 0x41, 0xb9, 0x30, 0xa0, 0x04, 0x18, 0x14, 0xd4, 0x88, 0x37, 0xce,
	0xb5, 0x15, 0xcd, 0xa1, 0xa8, 0xe4, 0xeb, 0x56, 0x82, 0xf2, 0x14, 0x14, 0x8b, 0x63, 0xe4, 0x58,
	0x39, 0x51, 0x6d, 0xf5, 0xe4, 0x8e, 0xce, 0x97, 0x96, 0xcb, 0xde, 0x64, 0x6e, 0xed, 0xd1, 0x9a,
	0xf3, 0x96, 0x3c, 0x49, 0x03, 0x63, 0xdf, 0xc2, 0xc8, 0x94, 0x6e, 0xac, 0xd3, 0x44, 0xda, 0x64,
	0xe5, 0x09, 0x8f, 0x99, 0xa2, 0x4d, 0x6e, 0x22, 0xdc, 0x20, 0x4d, 0x95, 0xf4, 0x29, 0xf3, 0xe5,
	0x29, 0xd0, 0x4d, 0xcb, 0x70, 0x6f, 0xa1, 0xd3, 0xc7, 0x69, 0x60, 0x0e, 0x6b, 0xc7, 0x32, 0xdd,
	0x58, 0x7f, 0x23, 0xb5, 0xa9, 0x66, 0xbc, 0x13, 0x6e, 0x22, 0x7b, 0x84, 0x23, 0x25, 0xfd, 0x5d,
	0xcb, 0x41, 0xca, 0xc7, 0xc4, 0xda, 0xda, 0xb2, 0xb1, 0xc6, 0x3e, 0x24, 0x16, 0x52, 0x3e, 0x26,
	0xd6, 0x93, 0x75, 0x1b, 0x6b, 0xfc, 0x43, 0x62, 0x21, 0xc5, 0x39, 0x21, 0x0f, 0x7b, 0x6f, 0x5c,
	0x53, 0xc3, 0xfc, 0x18, 0xe8, 0x5b, 0x2d, 0x05, 0x0d, 0x64, 0x92, 0x2a, 0xd0, 0x1a, 0x42, 0xea,
	0xef, 0x6c, 0xbb, 0xb7, 0x31, 0x3b, 0x16, 0x55, 0xf7, 0x2b, 0xe8, 0x63, 0x0b, 0x7f, 0xa5, 0xa5,
	0x78, 0xd6, 0x06, 0xef, 0xed, 0x6c, 0x3b, 0xaf, 0xc9, 0xed, 0xe2, 0xd9, 0xdb, 0x6f, 0xea, 0x3a,
	0xf3, 0xa5, 0xe5, 0x1b, 0x9b, 0xf7, 0x56, 0x7b, 0x26, 0x92, 0xd5, 0x22, 0x0d, 0x8a, 0xd7, 0xf5,
	0xc6, 0x79, 0xdf, 0x8a, 0xb3, 0x44, 0xc6, 0x30, 0xe5, 0x85, 0xae, 0x35, 0xf3, 0xcc, 0x98, 0xc0,
	0xc3, 0x8e, 0xda, 0xe5, 0x03, 0x5d, 0x6b, 0x66, 0x29, 0x71, 0x4c, 0xee, 0xb0, 0xd8, 0x72, 0x99,
	0xe1, 0xa7, 0xd0, 0x05, 0xd7, 0xee, 0x24, 0xc6, 0x9e, 0xeb, 0x8b, 0xbd, 0x1b, 0x9b, 0x36, 0x5b,
	0x7b, 0x53, 0x5d, 0xe4, 0xce, 0xb2, 0xbd, 0xa4, 0x62, 0x32, 0x12, 0x75, 0x50, 0xa9, 0xe2, 0x02,
	0x9f, 0x27, 0x60, 0x22, 0xe4, 0x21, 0x33, 0x40, 0x3b, 0x83, 0x93, 0x3b, 0x35, 0x3f, 0xb4, 0x5c,
	0xf2, 0x16, 0xf3, 0xa1, 0xa9, 0x8d, 0x3f, 0xac, 0x3d, 0x2b, 0xd0, 0x95, 0x62, 0x9e, 0xba, 0xa0,
	0xf5, 0x4c, 0x9f, 0xdf, 0x7a, 0x9c, 0x1f, 0xc8, 0x1c, 0x76, 0x94, 0xf3, 0x76, 0x70, 0x07, 0xb9,
	0xae, 0xc5, 0x0c, 0x8c, 0x7a, 0x51, 0xa3, 0x74, 0x2f, 0x68, 0x94, 0xeb, 0x64, 0x32, 0x83, 0xf7,
	0x69, 0x74, 0x06, 0x5f, 0xc2, 0x41, 0x5b, 0xaf, 0x42, 0x7d, 0xb2, 0xd6, 0xbf, 0x43, 0x13, 0x35,
	0x12, 0x5f, 0x30, 0x1e, 0xd3, 0x26, 0xf8, 0x29, 0x65, 0x9a, 0xfa, 0x4c, 0xc3, 0xce, 0x76, 0xb1,
	0x81, 0x59, 0xdc, 0xc0, 0xc3, 0xa0, 0x67, 0xdb, 0xc7, 0x05, 0xe9, 0x04, 0xfc, 0x74, 0x57, 0xef,
	0x21, 0x23, 0xdf, 0x55, 0x48, 0xd6, 0x43, 0x09, 0x9a, 0x0a, 0x69, 0x28, 0x4f, 0x52, 0x16, 0x18,
	0x6a, 0x22, 0xa0, 0x35, 0x19, 0xc7, 0xb2, 0x69, 0x8b, 0x73, 0xd7, 0x0d, 0x67, 0x31, 0xb4, 0xfb,
	0x15, 0x06, 0x79, 0x64, 0x79, 0x07, 0xd2, 0x54, 0x90, 0x75, 0x1c, 0xc1, 0x8b, 0x82, 0xf3, 0xac,
	0x4d, 0xc9, 0x82, 0x68, 0xe7, 0x3b, 0xe2, 0x72, 0x4d, 0xb9, 0x38, 0x65, 0xb1, 0xbd, 0x2c, 0x10,
	0x1a, 0xda, 0xd5, 0x7a, 0x0e, 0xab, 0xda, 0x14, 0xd7, 0x95, 0xcc, 0x5c, 0xb5, 0xd6, 0xa2, 0x58,
	0xff, 0x48, 0xe6, 0x06, 0xb2, 0xa8, 0x02, 0xa6, 0xa5, 0x70, 0xef, 0xe2, 0x56, 0x66, 0xf8, 0x59,
	0xaa, 0x87, 0x00, 0xe7, 0x4f, 0xb2, 0x92, 0x4b, 0x9f, 0x9b, 0x56, 0x76, 0xe3, 0xd4, 0x48, 0x5a,
	0xe3, 0x4a, 0x1b, 0x0a, 0xc2, 0xa8, 0x16, 0xe5, 0xa2, 0xeb, 0x90, 0xee, 0xd7, 0xf8, 0x1c, 0xdf,
	0x76, 0x48, 0xf8, 0x16, 0xc7, 0xf2, 0x85, 0x65, 0x3c, 0xb7, 0x84, 0x8a, 0xe8, 0x1c, 0xd0, 0x8e,
	0xfd, 0x41, 0x8a, 0xb5, 0x84, 0xf9, 0x3c, 0xe6, 0xa6, 0xe5, 0xde, 0xcb, 0x94, 0x15, 0xa4, 0x47,
	0x9d, 0x45, 0xe7, 0x31, 0xb9, 0x1d, 0x44, 0x3c, 0x0e, 0x7b, 0x90, 0xf3, 0x88, 0x1c, 0x47, 0x43,
	0x37, 0x78, 0x89, 0x8c, 0x75, 0xf2, 0x91, 0xe2, 0xb7, 0xc4, 0x37, 0x7d, 0xdf, 0x12, 0x47, 0xcc,
	0x44, 0x0b, 0xff, 0x0d, 0x91, 0xf1, 0x7e, 0xf5, 0x3b, 0x31, 0x79, 0xd2, 0xdf, 0x38, 0xa8, 0x6e,
	0x24, 0x09, 0x53, 0xad, 0xbc, 0x36, 0x31, 0x3d, 0xb0, 0x3c, 0x65, 0xa3, 0xd4, 0x8a, 0xea, 0xe9,
	0x2d, 0xd5, 0x8c, 0x89, 0x55, 0x6a, 0x57, 0x9f, 0xad, 0x53, 0xbf, 0x92, 0x47, 0x67, 0xa2, 0xd5,
	0x15, 0x4b, 0xa3, 0x81, 0x21, 0xca, 0x9d, 0x0a, 0xd8, 0x09, 0xf1, 0xd2, 0xc2, 0xcf, 0x7a, 0x56,
	0xe4, 0x69, 0x77, 0x2d, 0x3a, 0x13, 0x25, 0xe4, 0x81, 0x39, 0xef, 0x28, 0xd9, 0x9c, 0xb5, 0xde,
	0xc5, 0xef, 0xed, 0x98, 0xfb, 0x3c, 0x30, 0x83, 0x4e, 0xf3, 0x8a, 0xdc, 0x4f, 0xb8, 0xa0, 0xa2,
	0x91, 0xf8, 0xa0, 0x6c, 0x81, 0x80, 0xf7, 0x56, 0x10, 0x09, 0x33, 0x41, 0x04, 0x36, 0x7f, 0x69,
	0xca, 0xea, 0x80, 0x43, 0xd8, 0xa8, 0x77, 0x37, 0xe1, 0xe2, 0x00, 0x91, 0x87, 0xb5, 0xe7, 0x16,
	0xf7, 0x26, 0x83, 0x55, 0xc4, 0x11, 0xab, 0x83, 0xf3, 0x07, 0x79, 0x0c, 0x4c, 0xc5, 0xdc, 0x26,
	0x2b, 0x3b, 0x65, 0x3c, 0xc6, 0xbb, 0x47, 0x29, 0xcb, 0x5a, 0xe7, 0x18, 0x0a, 0x74, 0x23, 0x36,
	0x3a, 0x9f, 0xca, 0x96, 0x0a, 0xca, 0x6e, 0xc1, 0xd8, 0x67, 0x06, 0x0e, 0x6b, 0xed, 0x57, 0xce,
	0xd0, 0x0b, 0x7f, 0x97, 0xc8, 0x68, 0x4f, 0xed, 0x75, 0x66, 0xc9, 0x48, 0xa8, 0x58, 0x13, 0xf5,
	0x59, 0xc2, 0xbc, 0x6a, 0xff, 0x77, 0xa6, 0xc9, 0x70, 0x04, 0xc2, 0x30, 0x8e, 0x0f, 0x50, 0xf6,
	0xf2, 0x7f, 0x8e, 0x4b, 0xae, 0x09, 0x68, 0x18, 0xc5, 0x62, 0xbc, 0xb1, 0xb2, 0x57, 0xfc, 0x75,
	0x1c, 0x72, 0x25, 0x95, 0x4a, 0xe0, 0x49, 0xcb, 0x1e, 0xfe, 0xb6, 0x6b, 0x1a, 0xde, 0xb7, 0x70,
	0xa7, 0x65, 0x0f, 0x7f, 0x6f, 0xbe, 0x27, 0x13, 0xfb, 0x3d, 0xc3, 0x0e, 0x7e, 0x39, 0x3b, 0x8c,
	0x38, 0x67, 0x87, 0x7b, 0xe7, 0x7e, 0x5f, 0xf3, 0x18, 0xf4, 0x19, 0x3d, 0xbb, 0xd4, 0xdf, 0xdd,
	0x06, 0x7f, 0x24, 0xec, 0xad, 0xfd, 0xbe, 0x52, 0xe7, 0x26, 0x6a, 0xf8, 0xab, 0x81, 0x4c, 0xd6,
	0xb2, 0x41, 0x50, 0x80, 0x69, 0x4a, 0xf5, 0x6e, 0xad, 0x2e, 0xed, 0x90, 0xb8, 0xd6, 0xe3, 0xc8,
	0x1f, 0xc6, 0xcf, 0xf9, 0xad, 0xff, 0x03, 0x00, 0x00, 0xff, 0xff, 0xaf, 0xe2, 0x03, 0xfc, 0xe1,
	0x0f, 0x00, 0x00,
}
