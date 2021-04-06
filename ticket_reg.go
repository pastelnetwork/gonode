package pastel

type RegTicket struct {
	Version        int       `json:"version"`
	AuthorPastelID string    `json:"author_pastel_id"`
	BlockNum       int       `json:"block_num"`
	ImageHash      string    `json:"image_hash"`
	TotalCopies    int       `json:"total_copies"`
	Ticket         ArtTicket `json:"art_ticket"`
	Reserved       string    `json:"reserved"`
}
