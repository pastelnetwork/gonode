package pastel

type CascadeMultiVolumeTicket struct {
	ID string `json:"id"`
	HumanReadableIdentifier string            `json:"human_readable_identifier"`
    SHA3256HashOfOriginalFile string          `json:"sha3_256_hash_of_original_file"`
    Volumes          map[string]string `json:"volumes"`
}