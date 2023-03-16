package pastel

import (
	"bytes"
	"encoding/json"
	"strconv"

	"github.com/DataDog/zstd"
	"github.com/btcsuite/btcutil/base58"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/utils"
)

// SeparatorByte is the separator in dd_and_fingerprints.signature i.e. '.'
const SeparatorByte byte = 46

// DDAndFingerprints represents  the dd & fingerprints data returned by dd-service & SN
type DDAndFingerprints struct {
	BlockHash                                   string  `json:"pastel_block_hash_when_request_submitted"`
	BlockHeight                                 string  `json:"pastel_block_height_when_request_submitted"`
	TimestampOfRequest                          string  `json:"utc_timestamp_when_request_submitted"`
	SubmitterPastelID                           string  `json:"pastel_id_of_submitter"`
	SN1PastelID                                 string  `json:"pastel_id_of_registering_supernode_1"`
	SN2PastelID                                 string  `json:"pastel_id_of_registering_supernode_2"`
	SN3PastelID                                 string  `json:"pastel_id_of_registering_supernode_3"`
	IsOpenAPIRequest                            bool    `json:"is_pastel_openapi_request"`
	OpenAPISubsetID                             string  `json:"open_api_subset_id_string"`
	DupeDetectionSystemVersion                  string  `json:"dupe_detection_system_version"`
	IsLikelyDupe                                bool    `json:"is_likely_dupe"`
	IsRareOnInternet                            bool    `json:"is_rare_on_internet"`
	OverallRarenessScore                        float32 `json:"overall_rareness_score "`
	PctOfTop10MostSimilarWithDupeProbAbove25pct float32 `json:"pct_of_top_10_most_similar_with_dupe_prob_above_25pct"`
	PctOfTop10MostSimilarWithDupeProbAbove33pct float32 `json:"pct_of_top_10_most_similar_with_dupe_prob_above_33pct"`
	PctOfTop10MostSimilarWithDupeProbAbove50pct float32 `json:"pct_of_top_10_most_similar_with_dupe_prob_above_50pct"`

	OpenNSFWScore                              float32                `json:"open_nsfw_score"`
	ImageFingerprintOfCandidateImageFile       []float32              `json:"image_fingerprint_of_candidate_image_file"`
	HashOfCandidateImageFile                   string                 `json:"hash_of_candidate_image_file"`
	RarenessScoresTableJSONCompressedB64       string                 `json:"rareness_scores_table_json_compressed_b64"`
	CollectionNameString                       string                 `json:"collection_name_string,omitempty"`
	OpenAPIGroupIDString                       string                 `json:"open_api_group_id_string,omitempty"`
	GroupRarenessScore                         float32                `json:"group_rareness_score,omitempty"`
	CandidateImageThumbnailWebpAsBase64String  string                 `json:"candidate_image_thumbnail_webp_as_base64_string,omitempty"`
	DoesNotImpactTheFollowingCollectionStrings string                 `json:"does_not_impact_the_following_collection_strings,omitempty"`
	IsInvalidSenseRequest                      bool                   `json:"is_invalid_sense_request,omitempty"`
	InvalidSenseRequestReason                  string                 `json:"invalid_sense_request_reason,omitempty"`
	SimilarityScoreToFirstEntryInCollection    float32                `json:"similarity_score_to_first_entry_in_collection,omitempty"`
	CPProbability                              float32                `json:"cp_probability,omitempty"`
	ChildProbability                           float32                `json:"child_probability,omitempty"`
	ImageFilePath                              string                 `json:"image_file_path,omitempty"`
	InternetRareness                           *InternetRareness      `json:"internet_rareness"`
	AlternativeNSFWScores                      *AlternativeNSFWScores `json:"alternative_nsfw_scores"`
}

// InternetRareness defines internet rareness scores
type InternetRareness struct {
	RareOnInternetSummaryTableAsJSONCompressedB64    string `json:"rare_on_internet_summary_table_as_json_compressed_b64"`
	RareOnInternetGraphJSONCompressedB64             string `json:"rare_on_internet_graph_json_compressed_b64"`
	AlternativeRareOnInternetDictAsJSONCompressedB64 string `json:"alternative_rare_on_internet_dict_as_json_compressed_b64"`
	MinNumberOfExactMatchesInPage                    uint32 `json:"min_number_of_exact_matches_in_page"`
	EarliestAvailableDateOfInternetResults           string `json:"earliest_available_date_of_internet_results"`
}

// AlternativeNSFWScores represents the not-safe-for-work of an image
type AlternativeNSFWScores struct {
	Drawings float32 `json:"drawings"`
	Hentai   float32 `json:"hentai"`
	Neutral  float32 `json:"neutral"`
	Porn     float32 `json:"porn"`
	Sexy     float32 `json:"sexy"`
}

// Validate validates output
func (s DDAndFingerprints) Validate() error {
	if s.BlockHash == "" {
		return errors.Errorf("BlockHash is empty")
	}

	if s.BlockHeight == "" {
		return errors.Errorf("BlockHeight is empty")
	}

	if s.TimestampOfRequest == "" {
		return errors.Errorf("TimestampOfRequest is empty")
	}

	if s.SubmitterPastelID == "" {
		return errors.Errorf("SubmitterPastelID is empty")
	}

	if s.SN1PastelID == "" {
		return errors.Errorf("SN1PastelID is empty")
	}

	if s.SN2PastelID == "" {
		return errors.Errorf("SN2PastelID is empty")
	}

	if s.SN3PastelID == "" {
		return errors.Errorf("SN3PastelID is empty")
	}

	if s.HashOfCandidateImageFile == "" {
		return errors.Errorf("HashOfCandidateImageFile is empty")
	}

	if len(s.ImageFingerprintOfCandidateImageFile) == 0 {
		return errors.Errorf("ImageFingerprintOfCandidateImageFile is empty")
	}

	if s.InternetRareness == nil || s.AlternativeNSFWScores == nil {
		return errors.Errorf("InternetRareness or AlternativeNSFWScores is nil")
	}

	return nil
}

// ExtractCompressSignedDDAndFingerprints  decompresses & decodes
// the probe image reply which is: Base64URL(compress(dd_and_fingerprints.signature))
// and returns :
//- DDAndFingerprints structure data
//- JSON bytes of DDAndFingerprints struct (due to float value, so marshal might different on different platform)
//- signature, signature returned is base64 string of signature
func ExtractCompressSignedDDAndFingerprints(compressed []byte) (*DDAndFingerprints, []byte, []byte, error) {
	// Decompress compressedSignedDDAndFingerprints
	decompressedReply, err := zstd.Decompress(nil, compressed)
	if err != nil {
		return nil, nil, nil, errors.Errorf("decompress: %w", err)
	}

	var ddData, signature []byte
	for i := len(decompressedReply) - 1; i >= 0; i-- {
		if decompressedReply[i] == SeparatorByte {
			ddData = decompressedReply[:i]
			if i+1 >= len(decompressedReply) {
				return nil, nil, nil, errors.New("invalid probe image reply")
			}

			signature = decompressedReply[i+1:]
			break
		}
	}

	ddDataJSON, err := utils.B64Decode(ddData)
	if err != nil {
		return nil, nil, nil, errors.Errorf("decode %s failed: %w", string(decompressedReply), err)
	}

	dDataAndFingerprints := &DDAndFingerprints{}
	if err := json.Unmarshal(ddDataJSON, dDataAndFingerprints); err != nil {
		return nil, nil, nil, errors.Errorf("json decode dd_and_fingerprints: %w", err)
	}

	return dDataAndFingerprints, ddDataJSON, signature, nil
}

// ToCompressSignedDDAndFingerprints converts dd_and_fingerptints data to compress((b64(dd_and_fp).(sig)))
// https://pastel.wiki/en/Architecture/Workflows/NewArtRegistration
// Step 4.A.4
func ToCompressSignedDDAndFingerprints(ddData *DDAndFingerprints, signature []byte) ([]byte, error) {
	ddDataJSON, err := json.Marshal(ddData)

	if err != nil {
		return nil, errors.Errorf("marshal: %w", err)
	}

	res := utils.B64Encode(ddDataJSON)

	res = append(res, SeparatorByte)
	res = append(res, signature...)

	// Compress it
	compressed, err := zstd.CompressLevel(nil, res, 22)
	if err != nil {
		return nil, errors.Errorf("compress fingerprint data: %w", err)
	}

	return compressed, nil
}

// GetIDFiles generates ID Files for dd_and_fingerprints files and rq_id files
// file is b64 encoded file appended with signatures and compressed, ic is the initial counter
// and max is the number of ids to generate
func GetIDFiles(file []byte, ic uint32, max uint32) (ids []string, files [][]byte, err error) {
	var idFiles [][]byte
	for i := uint32(0); i < max; i++ {
		var buffer bytes.Buffer
		counter := ic + i

		buffer.Write(file)
		buffer.WriteByte(SeparatorByte)
		buffer.WriteString(strconv.Itoa(int(counter)))

		compressedData, err := zstd.CompressLevel(nil, buffer.Bytes(), 22)
		if err != nil {
			return ids, idFiles, errors.Errorf("compress identifiers file: %w", err)
		}

		idFiles = append(idFiles, compressedData)

		hash, err := utils.Sha3256hash(compressedData)
		if err != nil {
			return ids, idFiles, errors.Errorf("sha3-256-hash error getting an id file: %w", err)
		}

		ids = append(ids, base58.Encode(hash))
	}

	return ids, idFiles, nil
}
