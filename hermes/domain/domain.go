package domain

// CollectionState represents the state of collection
type CollectionState string

const (
	// ScoreReportLimit is the socre limit at which we report to CNode
	ScoreReportLimit = 3
	//UndefinedCollectionState is the state of collection when current no collection entries > max no of collection entries
	UndefinedCollectionState CollectionState = "undefined"
	// InProcessCollectionState is the state of collection when current no collection entries < max no of collection entries
	InProcessCollectionState CollectionState = "in_process"
	// FinalizedCollectionState is the state of collection when current no collection entries = max no of collection entries
	FinalizedCollectionState CollectionState = "finalized"
)

// DDFingerprints is dd & fp domain representation
type DDFingerprints struct {
	Sha256HashOfArtImageFile                   string    `json:"sha256_hash_of_art_image_file,omitempty"`
	PathToArtImageFile                         string    `json:"path_to_art_image_file,omitempty"`
	ImageFingerprintVector                     []float64 `json:"new_model_image_fingerprint_vector,omitempty"`
	DatetimeFingerprintAddedToDatabase         string    `json:"datetime_fingerprint_added_to_database,omitempty"`
	ImageThumbnailAsBase64                     string    `json:"thumbnail_of_image,omitempty"`
	RequestType                                string    `json:"request_type,omitempty"`
	OpenAPIGroupIDString                       string    `json:"open_api_group_id_string,omitempty"`
	CollectionNameString                       string    `json:"collection_name_string,omitempty"`
	DoesNotImpactTheFollowingCollectionsString string    `json:"does_not_impact_the_following_collection_strings,omitempty"`
	RegTXID                                    string    `json:"registration_ticket_txid"`
	TxIDTimestamp                              int64     `json:"txid_timestamp"`
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

// Collection is the domain representation collection table
type Collection struct {
	CollectionTicketTXID                           string          `json:"collection_ticket_txid"`
	CollectionName                                 string          `json:"collection_name_string"`
	CollectionTicketActivationBlockHeight          int             `json:"collection_ticket_activation_block_height"`
	CollectionFinalAllowedBlockHeight              int             `json:"collection_final_allowed_block_height"`
	MaxPermittedOpenNSFWScore                      float64         `json:"max_permitted_open_nsfw_score"`
	MinimumSimilarityScoreToFirstEntryInCollection float64         `json:"minimum_similarity_score_to_first_entry_in_collection"`
	CollectionState                                CollectionState `json:"collection_state"`
	DatetimeCollectionStateUpdated                 string          `json:"datetime_collection_state_updated"`
}

// String returns string representation of CollectionState
func (c CollectionState) String() string {
	switch c {
	case InProcessCollectionState:
		return "in_process"
	case FinalizedCollectionState:
		return "finalized"
	}

	return "undefined"
}

// NonImpactedCollection is the domain representation of the record from does not impact collection table
type NonImpactedCollection struct {
	ID                       int    `json:"id"`
	CollectionName           string `json:"collection_name"`
	Sha256HashOfArtImageFile string `json:"sha256_hash_of_art_image_file,omitempty"`
}

// NonImpactedCollections is the array type for holding non-impacted collection
type NonImpactedCollections []*NonImpactedCollection

// PastelBlock is the domain representation of pastel_blocks table
type PastelBlock struct {
	BlockHash          string `json:"block_hash"`
	BlockHeight        int32  `json:"block_height"`
	DatetimeBlockAdded string `json:"datetime_block_added"`
}

// IsValid checks if the pastel block is valid by matching the height and hash passed as an argument
func (p *PastelBlock) IsValid(blockHeight int32, hash string) bool {
	return p.BlockHash == hash && p.BlockHeight == blockHeight
}
