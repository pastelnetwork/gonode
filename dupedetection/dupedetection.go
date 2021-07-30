//go:generate mockery --name=Client

package dupedetection

import "context"

type DupeDetectionResult struct {
	DupeDetectionSystemVer  string              `json:"dupe_detection_system_version"`
	ImageHash               string              `json:"hash_of_candidate_image_file,omitempty"`
	IsLikelyDupe            float64             `json:"is_likely_dupe,omitempty"`
	PastelRarenessScore     float64             `json:"overall_average_rareness_score"`
	InternetRarenessScore   float64             `json:"is_rare_on_internet"`
	MatchesFoundOnFirstPage int                 `json:"matches_found_on_first_page"`
	NumberOfResultPages     int                 `json:"number_of_pages_of_results"`
	FirstMatchURL           string              `json:"url_of_first_match_in_page"`
	OpenNSFWScore           float64             `json:"open_nsfw_score"`
	AlternateNSFWScores     AlternateNSFWScores `json:"alternative_nsfw_scores"`
	FingerPrints            string              `json:"image_fingerprint_of_candidate_image_file"`
}

type AlternateNSFWScores struct {
	Drawings float64 `json:"drawings"`
	Hentai   float32 `json:"hentai"`
	Neutral  float32 `json:"neutral"`
	Porn     float32 `json:"porn"`
	Sexy     float32 `json:"sexy"`
}

type Client interface {
	// Generate returns fingerprints and ranks for a given image
	Generate(ctx context.Context, path string) (*DupeDetectionResult, error)
}
