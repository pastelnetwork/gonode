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

	OpenAPIGroupIDString                       string `json:"open_api_group_id_string,omitempty"`
	CollectionNameString                       string `json:"collection_name_string,omitempty"`
	DoesNotImpactTheFollowingCollectionsString string `json:"does_not_impact_the_following_collection_strings,omitempty"`
}

// SnScore is domain representation of Sn score record
type SnScore struct {
	TxID      string
	PastelID  string
	IPAddress string
	Score     int
}

// MasterNodeConf is the domain representation of masternode.conf file
type MasterNodeConf struct {
	MnAddress  string `json:"mnAddress"`
	MnPrivKey  string `json:"mnPrivKey"`
	Txid       string `json:"txid"`
	OutIndex   string `json:"outIndex"`
	ExtAddress string `json:"extAddress"`
	ExtP2P     string `json:"extP2P"`
	ExtCfg     string `json:"extCfg"`
	ExtKey     string `json:"extKey"`
}
