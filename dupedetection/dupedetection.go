//go:generate mockery --name=Client

package dupedetection

import "context"

type DupeDetectionResult struct {
	PastelRarenessScore   float64
	InternetRarenessScore float64
	OpenNSFWScore         float64
	AlternateNSFWScores   AlternateNSFWScores
	FingerPrints          string
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
