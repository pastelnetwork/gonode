//go:generate mockery --name=Client

package dupedetection

import "context"

// DupeDetection is the dupe detection result which will be sent to the caller
type DupeDetection struct {
	DupeDetectionSystemVer  string              `json:"dupe_detection_system_version"`
	ImageHash               string              `json:"hash_of_candidate_image_file"`
	PastelRarenessScore     float64             `json:"overall_average_rareness_score"`
	InternetRarenessScore   float64             `json:"is_rare_on_internet"`
	MatchesFoundOnFirstPage int                 `json:"matches_found_on_first_page"`
	NumberOfResultPages     int                 `json:"number_of_pages_of_results"`
	FirstMatchURL           string              `json:"url_of_first_match_in_page"`
	OpenNSFWScore           float64             `json:"open_nsfw_score"`
	AlternateNSFWScores     AlternateNSFWScores `json:"alternative_nsfw_scores"`
	ImageHashes             ImageHashes         `json:"image_hashes"`
	Fingerprints            string              `json:"image_fingerprint_of_candidate_image_file"`
}

// AlternateNSFWScores represents alternate NSFW scores in the output of dupe detection service
type AlternateNSFWScores struct {
	Drawings float64 `json:"drawings"`
	Hentai   float64 `json:"hentai"`
	Neutral  float64 `json:"neutral"`
	Porn     float64 `json:"porn"`
	Sexy     float64 `json:"sexy"`
}

// ImageHashes represents image hashes in the output of dupe detection service
type ImageHashes struct {
	PerceptualHash string `json:"perceptual_hash"`
	AverageHash    string `json:"average_hash"`
	DifferenceHash string `json:"difference_hash"`
}

// Client represents the interface to communicate with dupe-detection service
type Client interface {
	// Generate returns fingerprints and ranks for a given image
	Generate(ctx context.Context, img []byte, format string) (*DupeDetection, error)
}
