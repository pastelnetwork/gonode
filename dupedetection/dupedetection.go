//go:generate mockery --name=Client

package dupedetection

import "context"

// DupeDetection is the dupe detection result which will be sent to the caller
type DupeDetection struct {
	DupeDetectionSystemVer  string
	PastelRarenessScore     float64
	InternetRarenessScore   float64
	MatchesFoundOnFirstPage int
	NumberOfResultPages     int
	FirstMatchURL           string
	OpenNSFWScore           float64
	AlternateNSFWScores     AlternateNSFWScores
	ImageHashes             ImageHashes
	FingerPrints            []float64
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

type Client interface {
	// Generate returns fingerprints and ranks for a given image
	Generate(ctx context.Context, path string) (*DupeDetection, error)
}
