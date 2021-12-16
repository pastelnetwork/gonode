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

// separatorByte is the separator in dd_and_fingerprints.signature i.e. '.'
var SeparatorByte byte = 46

// DDAndFingerprints represents  the dd & fingerprints data returned by dd-service & SN
type DDAndFingerprints struct {
	Block                      string                 `json:"block"`
	Principal                  string                 `json:"principal"`
	DupeDetectionSystemVersion string                 `json:"dupe_detection_system_version"`
	Fingerprints               Fingerprint            `json:"fingerprints"`
	PastelRarenessScore        float32                `json:"pastel_rareness_score"`
	IsLikelyDupe               bool                   `json:"is_likely_dupe"`
	InternetRarenessScore      float32                `json:"internet_rareness_score"`
	IsRareOnInternet           bool                   `json:"is_rare_on_internet"`
	MatchesFoundOnFirstPage    uint32                 `json:"matches_found_on_first_page"`
	NumberOfPagesOfResults     uint32                 `json:"number_of_pages_of_results"`
	URLOfFirstMatchInPage      string                 `json:"url_of_first_match_in_page"`
	OpenNSFWScore              float32                `json:"open_nsfw_score"`
	PerceptualHashOverlapCount uint32                 `json:"perceptual_hash_overlap_count"`
	AlternateNSFWScores        *AlternativeNSFWScore  `json:"alternate_nsfw_scores"`
	PerceptualImageHashes      *PerceptualImageHashes `json:"perceptual_image_hashes"`
	FingerprintsStat           *DDFingerprintsStat    `json:"fingerprints_stat"`
	Maxes                      *DDMaxes               `json:"maxes"`
	Percentile                 *DDPercentiles         `json:"percentile"`
	Score                      *DDScores              `json:"score"`

	// DD-Server does not directly return these comporessed fingerprints
	// We generate and assign to this field to avoid repeated operations
	ZstdCompressedFingerprint []byte `json:"zstd_compressed_fingerprint,omitempty"`
}

// AlternativeNSFWScore represents the not-safe-for-work of an image
type AlternativeNSFWScore struct {
	Drawing float32 `json:"drawing"`
	Hentai  float32 `json:"hentai"`
	Neutral float32 `json:"neutral"`
	Porn    float32 `json:"porn"`
	Sexy    float32 `json:"sexy"`
}

type DDPercentiles struct {
	PearsonTop1BpsPercentile             float32 `json:"pearson_top_1_bps_percentile"`
	SpearmanTop1BpsPercentile            float32 `json:"spearman_top_1_bps_percentile"`
	KendallTop1BpsPercentile             float32 `json:"kendall_top_1_bps_percentile"`
	HoeffdingTop10BpsPercentile          float32 `json:"hoeffding_top_10_bps_percentile"`
	MutualInformationTop100BpsPercentile float32 `json:"mutual_information_top_100_bps_percentile"`
	HsicTop100BpsPercentile              float32 `json:"hsic_top_100_bps_percentile"`
	XgbimportanceTop100BpsPercentile     float32 `json:"xgbimportance_top_100_bps_percentile"`
}

type DDFingerprintsStat struct {
	NumberOfFingerprintsRequiringFurtherTesting1 uint32 `json:"number_of_fingerprints_requiring_further_testing_1"`
	NumberOfFingerprintsRequiringFurtherTesting2 uint32 `json:"number_of_fingerprints_requiring_further_testing_2"`
	NumberOfFingerprintsRequiringFurtherTesting3 uint32 `json:"number_of_fingerprints_requiring_further_testing_3"`
	NumberOfFingerprintsRequiringFurtherTesting4 uint32 `json:"number_of_fingerprints_requiring_further_testing_4"`
	NumberOfFingerprintsRequiringFurtherTesting5 uint32 `json:"number_of_fingerprints_requiring_further_testing_5"`
	NumberOfFingerprintsRequiringFurtherTesting6 uint32 `json:"number_of_fingerprints_requiring_further_testing_6"`
	NumberOfFingerprintsOfSuspectedDupes         uint32 `json:"number_of_fingerprints_of_suspected_dupes"`
}

type DDMaxes struct {
	PearsonMax           float32 `json:"pearson_max"`
	SpearmanMax          float32 `json:"spearman_max"`
	KendallMax           float32 `json:"kendall_max"`
	HoeffdingMax         float32 `json:"hoeffding_max"`
	MutualInformationMax float32 `json:"mutual_information_max"`
	HsicMax              float32 `json:"hsic_max"`
	XgbimportanceMax     float32 `json:"xgbimportance_max"`
}

type DDScores struct {
	CombinedRarenessScore         float32 `json:"combined_rareness_score"`
	XgboostPredictedRarenessScore float32 `json:"xgboost_predicted_rareness_score"`
	NnPredictedRarenessScore      float32 `json:"nn_predicted_rareness_score"`
	OverallAverageRarenessScore   float32 `json:"overall_average_rareness_score"`
}

// AlternateNSFWScores represents alternate NSFW scores from dupe detection service
type AlternateNSFWScores struct {
	Drawing float32 `json:"drawing"`
	Hentai  float32 `json:"hentai"`
	Neutral float32 `json:"neutral"`
	Porn    float32 `json:"porn"`
	Sexy    float32 `json:"sexy"`
}

// PerceptualImageHashes represents image hashes from dupe detection service
type PerceptualImageHashes struct {
	PDQHash        string `json:"pdq_hash"`
	PerceptualHash string `json:"perceptual_hash"`
	AverageHash    string `json:"average_hash"`
	DifferenceHash string `json:"difference_hash"`
	NeuralHash     string `json:"neuralhash_hash"`
}

// GetDDFingerprintsAndSigFromProbeImageReply decodes and decompresses
// the probe image reply which is: Base64URL(compress(dd_and_fingerprints.signature))
// and returns dd_and_fingerprints data and signature
func GetDDFingerprintsAndSigFromProbeImageReply(probeImageReply []byte) (*DDAndFingerprints, []byte, error) {
	decompressedReply, err := zstd.Decompress(nil, probeImageReply)
	if err != nil {
		return nil, nil, errors.Errorf("decompress probe image reply: %w", err)
	}

	var ddData, signature []byte
	for i := len(decompressedReply) - 1; i >= 0; i-- {
		if decompressedReply[i] == SeparatorByte {
			ddData = decompressedReply[:i]
			if i+1 >= len(decompressedReply) {
				return nil, nil, errors.New("invalid probe image reply")
			}

			signature = decompressedReply[i+1:]
			break
		}
	}

	ddDataJson, err := utils.B64Decode(ddData)
	if err != nil {
		return nil, nil, errors.Errorf("b64 decode: %w", err)
	}

	ddDataAndFingerprints := &DDAndFingerprints{}
	if err := json.Unmarshal(ddDataJson, ddDataAndFingerprints); err != nil {
		return nil, nil, errors.Errorf("json decode dd_and_fingerprints: %w", err)
	}

	return ddDataAndFingerprints, signature, nil
}

// GetProbeImageReply returns Base64URL(compress(dd_and_fingerprints.signature))
func GetProbeImageReply(ddData *DDAndFingerprints, signature []byte) ([]byte, error) {
	ddDataJson, err := json.Marshal(ddData)
	if err != nil {
		return nil, errors.Errorf("GetProbeImageReply: failed to marshal dd-data: %w", err)
	}

	res := utils.B64Encode(ddDataJson)

	res = append(res, SeparatorByte)
	res = append(res, signature...)

	compressed, err := zstd.CompressLevel(nil, res, 22)
	if err != nil {
		return nil, errors.Errorf("GetProbeImageReply: compress: %w", err)
	}

	return compressed, nil
}

// GetIDFiles is supposed to generates ID Files for dd_and_fingerprints files and rq_id files
// file is b64 encoded file appended with signatures and compressed, ic is the initial counter
// and max is the number of ids to generate
func GetIDFiles(file []byte, ic uint32, max uint32) (ids []string, err error) {
	for i := uint32(0); i < max; i++ {
		var buffer bytes.Buffer
		counter := ic + i

		buffer.Write(file)
		buffer.WriteByte(SeparatorByte)
		buffer.WriteString(strconv.Itoa(int(counter)))

		compressedData, err := zstd.CompressLevel(nil, buffer.Bytes(), 22)
		if err != nil {
			return ids, errors.Errorf("compress identifiers file: %w", err)
		}

		hash, err := utils.Sha3256hash(compressedData)
		if err != nil {
			return ids, errors.Errorf("sha3256hash: %w", err)
		}

		ids = append(ids, base58.Encode(hash))
	}

	return ids, nil
}
