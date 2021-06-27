package artworkregister

type ImageThumbnail struct {
	TopLeftX     int64 `json:"top_left_x"`
	TopLeftY     int64 `json:"top_left_y"`
	BottomRightX int64 `json:"bottom_right_x"`
	BottomRightY int64 `json:"bottom_right_y"`
}
