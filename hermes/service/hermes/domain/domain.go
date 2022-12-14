package domain

const (
	// ScoreReportLimit is the socre limit at which we report to CNode
	ScoreReportLimit = 3
)

// DDFingerprints is dd & fp domain representation
type DDFingerprints struct {
	Sha256HashOfArtImageFile           string    `json:"sha256_hash_of_art_image_file,omitempty"`
	PathToArtImageFile                 string    `json:"path_to_art_image_file,omitempty"`
	ImageFingerprintVector             []float64 `json:"new_model_image_fingerprint_vector,omitempty"`
	DatetimeFingerprintAddedToDatabase string    `json:"datetime_fingerprint_added_to_database,omitempty"`
	ImageThumbnailAsBase64             string    `json:"thumbnail_of_image,omitempty"`
	RequestType                        string    `json:"request_type,omitempty"`
	IDString                           string    `json:"open_api_subset_id_string,omitempty"`
}

// SnScore is domain representation of Sn score record
type SnScore struct {
	TxID      string
	PastelID  string
	IPAddress string
	Score     int
}
