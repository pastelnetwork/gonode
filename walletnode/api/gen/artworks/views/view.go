// Code generated by goa v3.4.3, DO NOT EDIT.
//
// artworks views
//
// Command:
// $ goa gen github.com/pastelnetwork/gonode/walletnode/api/design

package views

import (
	"unicode/utf8"

	goa "goa.design/goa/v3/pkg"
)

// RegisterResult is the viewed result type that is projected based on a view.
type RegisterResult struct {
	// Type to project
	Projected *RegisterResultView
	// View to render
	View string
}

// Task is the viewed result type that is projected based on a view.
type Task struct {
	// Type to project
	Projected *TaskView
	// View to render
	View string
}

// TaskCollection is the viewed result type that is projected based on a view.
type TaskCollection struct {
	// Type to project
	Projected TaskCollectionView
	// View to render
	View string
}

// Image is the viewed result type that is projected based on a view.
type Image struct {
	// Type to project
	Projected *ImageView
	// View to render
	View string
}

// RegisterResultView is a type that runs validations on a projected type.
type RegisterResultView struct {
	// Task ID of the registration process
	TaskID *string
}

// TaskView is a type that runs validations on a projected type.
type TaskView struct {
	// JOb ID of the registration process
	ID *string
	// Status of the registration process
	Status *string
	// List of states from the very beginning of the process
	States []*TaskStateView
	// txid
	Txid   *string
	Ticket *ArtworkTicketView
}

// TaskStateView is a type that runs validations on a projected type.
type TaskStateView struct {
	// Date of the status creation
	Date *string
	// Status of the registration process
	Status *string
}

// ArtworkTicketView is a type that runs validations on a projected type.
type ArtworkTicketView struct {
	// Name of the artwork
	Name *string
	// Description of the artwork
	Description *string
	// Keywords
	Keywords *string
	// Series name
	SeriesName *string
	// Number of copies issued
	IssuedCopies *int
	// Artwork creation video youtube URL
	YoutubeURL *string
	// Artist's PastelID
	ArtistPastelID *string
	// Passphrase of the artist's PastelID
	ArtistPastelIDPassphrase *string
	// Name of the artist
	ArtistName *string
	// Artist website URL
	ArtistWebsiteURL *string
	// Spendable address
	SpendableAddress *string
	// Used to find a suitable masternode with a fee equal or less
	MaximumFee *float64
	// Percentage the artist received in future sales. If set to 0% he only get
	// paids for the first sale on each copy of the NFT
	Royalty *float64
	// To donate 2% of the sale proceeds on every sale to TeamTrees which plants
	// trees
	Green               *bool
	ThumbnailCoordinate *ThumbnailcoordinateView
}

// ThumbnailcoordinateView is a type that runs validations on a projected type.
type ThumbnailcoordinateView struct {
	// X coordinate of the thumbnail's top left conner
	TopLeftX *int64
	// Y coordinate of the thumbnail's top left conner
	TopLeftY *int64
	// X coordinate of the thumbnail's bottom right conner
	BottomRightX *int64
	// Y coordinate of the thumbnail's bottom right conner
	BottomRightY *int64
}

// TaskCollectionView is a type that runs validations on a projected type.
type TaskCollectionView []*TaskView

// ImageView is a type that runs validations on a projected type.
type ImageView struct {
	// Uploaded image ID
	ImageID *string
	// Image expiration
	ExpiresIn *string
}

var (
	// RegisterResultMap is a map of attribute names in result type RegisterResult
	// indexed by view name.
	RegisterResultMap = map[string][]string{
		"default": []string{
			"task_id",
		},
	}
	// TaskMap is a map of attribute names in result type Task indexed by view name.
	TaskMap = map[string][]string{
		"tiny": []string{
			"id",
			"status",
			"txid",
			"ticket",
		},
		"default": []string{
			"id",
			"status",
			"states",
			"txid",
			"ticket",
		},
	}
	// TaskCollectionMap is a map of attribute names in result type TaskCollection
	// indexed by view name.
	TaskCollectionMap = map[string][]string{
		"tiny": []string{
			"id",
			"status",
			"txid",
			"ticket",
		},
		"default": []string{
			"id",
			"status",
			"states",
			"txid",
			"ticket",
		},
	}
	// ImageMap is a map of attribute names in result type Image indexed by view
	// name.
	ImageMap = map[string][]string{
		"default": []string{
			"image_id",
			"expires_in",
		},
	}
	// ThumbnailcoordinateMap is a map of attribute names in result type
	// Thumbnailcoordinate indexed by view name.
	ThumbnailcoordinateMap = map[string][]string{
		"default": []string{
			"top_left_x",
			"top_left_y",
			"bottom_right_x",
			"bottom_right_y",
		},
	}
)

// ValidateRegisterResult runs the validations defined on the viewed result
// type RegisterResult.
func ValidateRegisterResult(result *RegisterResult) (err error) {
	switch result.View {
	case "default", "":
		err = ValidateRegisterResultView(result.Projected)
	default:
		err = goa.InvalidEnumValueError("view", result.View, []interface{}{"default"})
	}
	return
}

// ValidateTask runs the validations defined on the viewed result type Task.
func ValidateTask(result *Task) (err error) {
	switch result.View {
	case "tiny":
		err = ValidateTaskViewTiny(result.Projected)
	case "default", "":
		err = ValidateTaskView(result.Projected)
	default:
		err = goa.InvalidEnumValueError("view", result.View, []interface{}{"tiny", "default"})
	}
	return
}

// ValidateTaskCollection runs the validations defined on the viewed result
// type TaskCollection.
func ValidateTaskCollection(result TaskCollection) (err error) {
	switch result.View {
	case "tiny":
		err = ValidateTaskCollectionViewTiny(result.Projected)
	case "default", "":
		err = ValidateTaskCollectionView(result.Projected)
	default:
		err = goa.InvalidEnumValueError("view", result.View, []interface{}{"tiny", "default"})
	}
	return
}

// ValidateImage runs the validations defined on the viewed result type Image.
func ValidateImage(result *Image) (err error) {
	switch result.View {
	case "default", "":
		err = ValidateImageView(result.Projected)
	default:
		err = goa.InvalidEnumValueError("view", result.View, []interface{}{"default"})
	}
	return
}

// ValidateRegisterResultView runs the validations defined on
// RegisterResultView using the "default" view.
func ValidateRegisterResultView(result *RegisterResultView) (err error) {
	if result.TaskID == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("task_id", "result"))
	}
	if result.TaskID != nil {
		if utf8.RuneCountInString(*result.TaskID) < 8 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("result.task_id", *result.TaskID, utf8.RuneCountInString(*result.TaskID), 8, true))
		}
	}
	if result.TaskID != nil {
		if utf8.RuneCountInString(*result.TaskID) > 8 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("result.task_id", *result.TaskID, utf8.RuneCountInString(*result.TaskID), 8, false))
		}
	}
	return
}

// ValidateTaskViewTiny runs the validations defined on TaskView using the
// "tiny" view.
func ValidateTaskViewTiny(result *TaskView) (err error) {
	if result.ID == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("id", "result"))
	}
	if result.Status == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("status", "result"))
	}
	if result.Ticket == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("ticket", "result"))
	}
	if result.ID != nil {
		if utf8.RuneCountInString(*result.ID) < 8 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("result.id", *result.ID, utf8.RuneCountInString(*result.ID), 8, true))
		}
	}
	if result.ID != nil {
		if utf8.RuneCountInString(*result.ID) > 8 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("result.id", *result.ID, utf8.RuneCountInString(*result.ID), 8, false))
		}
	}
	if result.Status != nil {
		if !(*result.Status == "Task Started" || *result.Status == "Connected" || *result.Status == "Image Probed" || *result.Status == "Image And Thumbnail Uploaded" || *result.Status == "Status Gen ReptorQ Symbols" || *result.Status == "Preburn Registration Fee" || *result.Status == "Ticket Accepted" || *result.Status == "Ticket Registered" || *result.Status == "Ticket Activated" || *result.Status == "Error Insufficient Fee" || *result.Status == "Error Fingerprints Dont Match" || *result.Status == "Error ThumbnailHashes Dont Match" || *result.Status == "Error GenRaptorQ Symbols Failed" || *result.Status == "Task Rejected" || *result.Status == "Task Completed") {
			err = goa.MergeErrors(err, goa.InvalidEnumValueError("result.status", *result.Status, []interface{}{"Task Started", "Connected", "Image Probed", "Image And Thumbnail Uploaded", "Status Gen ReptorQ Symbols", "Preburn Registration Fee", "Ticket Accepted", "Ticket Registered", "Ticket Activated", "Error Insufficient Fee", "Error Fingerprints Dont Match", "Error ThumbnailHashes Dont Match", "Error GenRaptorQ Symbols Failed", "Task Rejected", "Task Completed"}))
		}
	}
	if result.Txid != nil {
		if utf8.RuneCountInString(*result.Txid) < 64 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("result.txid", *result.Txid, utf8.RuneCountInString(*result.Txid), 64, true))
		}
	}
	if result.Txid != nil {
		if utf8.RuneCountInString(*result.Txid) > 64 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("result.txid", *result.Txid, utf8.RuneCountInString(*result.Txid), 64, false))
		}
	}
	if result.Ticket != nil {
		if err2 := ValidateArtworkTicketView(result.Ticket); err2 != nil {
			err = goa.MergeErrors(err, err2)
		}
	}
	return
}

// ValidateTaskView runs the validations defined on TaskView using the
// "default" view.
func ValidateTaskView(result *TaskView) (err error) {
	if result.ID == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("id", "result"))
	}
	if result.Status == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("status", "result"))
	}
	if result.Ticket == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("ticket", "result"))
	}
	if result.ID != nil {
		if utf8.RuneCountInString(*result.ID) < 8 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("result.id", *result.ID, utf8.RuneCountInString(*result.ID), 8, true))
		}
	}
	if result.ID != nil {
		if utf8.RuneCountInString(*result.ID) > 8 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("result.id", *result.ID, utf8.RuneCountInString(*result.ID), 8, false))
		}
	}
	if result.Status != nil {
		if !(*result.Status == "Task Started" || *result.Status == "Connected" || *result.Status == "Image Probed" || *result.Status == "Image And Thumbnail Uploaded" || *result.Status == "Status Gen ReptorQ Symbols" || *result.Status == "Preburn Registration Fee" || *result.Status == "Ticket Accepted" || *result.Status == "Ticket Registered" || *result.Status == "Ticket Activated" || *result.Status == "Error Insufficient Fee" || *result.Status == "Error Fingerprints Dont Match" || *result.Status == "Error ThumbnailHashes Dont Match" || *result.Status == "Error GenRaptorQ Symbols Failed" || *result.Status == "Task Rejected" || *result.Status == "Task Completed") {
			err = goa.MergeErrors(err, goa.InvalidEnumValueError("result.status", *result.Status, []interface{}{"Task Started", "Connected", "Image Probed", "Image And Thumbnail Uploaded", "Status Gen ReptorQ Symbols", "Preburn Registration Fee", "Ticket Accepted", "Ticket Registered", "Ticket Activated", "Error Insufficient Fee", "Error Fingerprints Dont Match", "Error ThumbnailHashes Dont Match", "Error GenRaptorQ Symbols Failed", "Task Rejected", "Task Completed"}))
		}
	}
	for _, e := range result.States {
		if e != nil {
			if err2 := ValidateTaskStateView(e); err2 != nil {
				err = goa.MergeErrors(err, err2)
			}
		}
	}
	if result.Txid != nil {
		if utf8.RuneCountInString(*result.Txid) < 64 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("result.txid", *result.Txid, utf8.RuneCountInString(*result.Txid), 64, true))
		}
	}
	if result.Txid != nil {
		if utf8.RuneCountInString(*result.Txid) > 64 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("result.txid", *result.Txid, utf8.RuneCountInString(*result.Txid), 64, false))
		}
	}
	if result.Ticket != nil {
		if err2 := ValidateArtworkTicketView(result.Ticket); err2 != nil {
			err = goa.MergeErrors(err, err2)
		}
	}
	return
}

// ValidateTaskStateView runs the validations defined on TaskStateView.
func ValidateTaskStateView(result *TaskStateView) (err error) {
	if result.Date == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("date", "result"))
	}
	if result.Status == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("status", "result"))
	}
	if result.Status != nil {
		if !(*result.Status == "Task Started" || *result.Status == "Connected" || *result.Status == "Image Probed" || *result.Status == "Image And Thumbnail Uploaded" || *result.Status == "Status Gen ReptorQ Symbols" || *result.Status == "Preburn Registration Fee" || *result.Status == "Ticket Accepted" || *result.Status == "Ticket Registered" || *result.Status == "Ticket Activated" || *result.Status == "Error Insufficient Fee" || *result.Status == "Error Fingerprints Dont Match" || *result.Status == "Error ThumbnailHashes Dont Match" || *result.Status == "Error GenRaptorQ Symbols Failed" || *result.Status == "Task Rejected" || *result.Status == "Task Completed") {
			err = goa.MergeErrors(err, goa.InvalidEnumValueError("result.status", *result.Status, []interface{}{"Task Started", "Connected", "Image Probed", "Image And Thumbnail Uploaded", "Status Gen ReptorQ Symbols", "Preburn Registration Fee", "Ticket Accepted", "Ticket Registered", "Ticket Activated", "Error Insufficient Fee", "Error Fingerprints Dont Match", "Error ThumbnailHashes Dont Match", "Error GenRaptorQ Symbols Failed", "Task Rejected", "Task Completed"}))
		}
	}
	return
}

// ValidateArtworkTicketView runs the validations defined on ArtworkTicketView.
func ValidateArtworkTicketView(result *ArtworkTicketView) (err error) {
	if result.ArtistName == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("artist_name", "result"))
	}
	if result.Name == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("name", "result"))
	}
	if result.IssuedCopies == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("issued_copies", "result"))
	}
	if result.ArtistPastelID == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("artist_pastelid", "result"))
	}
	if result.ArtistPastelIDPassphrase == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("artist_pastelid_passphrase", "result"))
	}
	if result.SpendableAddress == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("spendable_address", "result"))
	}
	if result.MaximumFee == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("maximum_fee", "result"))
	}
	if result.Name != nil {
		if utf8.RuneCountInString(*result.Name) > 256 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("result.name", *result.Name, utf8.RuneCountInString(*result.Name), 256, false))
		}
	}
	if result.Description != nil {
		if utf8.RuneCountInString(*result.Description) > 1024 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("result.description", *result.Description, utf8.RuneCountInString(*result.Description), 1024, false))
		}
	}
	if result.Keywords != nil {
		if utf8.RuneCountInString(*result.Keywords) > 256 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("result.keywords", *result.Keywords, utf8.RuneCountInString(*result.Keywords), 256, false))
		}
	}
	if result.SeriesName != nil {
		if utf8.RuneCountInString(*result.SeriesName) > 256 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("result.series_name", *result.SeriesName, utf8.RuneCountInString(*result.SeriesName), 256, false))
		}
	}
	if result.IssuedCopies != nil {
		if *result.IssuedCopies < 1 {
			err = goa.MergeErrors(err, goa.InvalidRangeError("result.issued_copies", *result.IssuedCopies, 1, true))
		}
	}
	if result.IssuedCopies != nil {
		if *result.IssuedCopies > 1000 {
			err = goa.MergeErrors(err, goa.InvalidRangeError("result.issued_copies", *result.IssuedCopies, 1000, false))
		}
	}
	if result.YoutubeURL != nil {
		if utf8.RuneCountInString(*result.YoutubeURL) > 128 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("result.youtube_url", *result.YoutubeURL, utf8.RuneCountInString(*result.YoutubeURL), 128, false))
		}
	}
	if result.ArtistPastelID != nil {
		err = goa.MergeErrors(err, goa.ValidatePattern("result.artist_pastelid", *result.ArtistPastelID, "^[a-zA-Z0-9]+$"))
	}
	if result.ArtistPastelID != nil {
		if utf8.RuneCountInString(*result.ArtistPastelID) < 86 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("result.artist_pastelid", *result.ArtistPastelID, utf8.RuneCountInString(*result.ArtistPastelID), 86, true))
		}
	}
	if result.ArtistPastelID != nil {
		if utf8.RuneCountInString(*result.ArtistPastelID) > 86 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("result.artist_pastelid", *result.ArtistPastelID, utf8.RuneCountInString(*result.ArtistPastelID), 86, false))
		}
	}
	if result.ArtistName != nil {
		if utf8.RuneCountInString(*result.ArtistName) > 256 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("result.artist_name", *result.ArtistName, utf8.RuneCountInString(*result.ArtistName), 256, false))
		}
	}
	if result.ArtistWebsiteURL != nil {
		if utf8.RuneCountInString(*result.ArtistWebsiteURL) > 256 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("result.artist_website_url", *result.ArtistWebsiteURL, utf8.RuneCountInString(*result.ArtistWebsiteURL), 256, false))
		}
	}
	if result.SpendableAddress != nil {
		err = goa.MergeErrors(err, goa.ValidatePattern("result.spendable_address", *result.SpendableAddress, "^[a-zA-Z0-9]+$"))
	}
	if result.SpendableAddress != nil {
		if utf8.RuneCountInString(*result.SpendableAddress) < 35 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("result.spendable_address", *result.SpendableAddress, utf8.RuneCountInString(*result.SpendableAddress), 35, true))
		}
	}
	if result.SpendableAddress != nil {
		if utf8.RuneCountInString(*result.SpendableAddress) > 35 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("result.spendable_address", *result.SpendableAddress, utf8.RuneCountInString(*result.SpendableAddress), 35, false))
		}
	}
	if result.MaximumFee != nil {
		if *result.MaximumFee < 1e-05 {
			err = goa.MergeErrors(err, goa.InvalidRangeError("result.maximum_fee", *result.MaximumFee, 1e-05, true))
		}
	}
	if result.Royalty != nil {
		if *result.Royalty < 0 {
			err = goa.MergeErrors(err, goa.InvalidRangeError("result.royalty", *result.Royalty, 0, true))
		}
	}
	if result.Royalty != nil {
		if *result.Royalty > 100 {
			err = goa.MergeErrors(err, goa.InvalidRangeError("result.royalty", *result.Royalty, 100, false))
		}
	}
	if result.ThumbnailCoordinate != nil {
		if err2 := ValidateThumbnailcoordinateView(result.ThumbnailCoordinate); err2 != nil {
			err = goa.MergeErrors(err, err2)
		}
	}
	return
}

// ValidateThumbnailcoordinateView runs the validations defined on
// ThumbnailcoordinateView using the "default" view.
func ValidateThumbnailcoordinateView(result *ThumbnailcoordinateView) (err error) {
	if result.TopLeftX == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("top_left_x", "result"))
	}
	if result.TopLeftY == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("top_left_y", "result"))
	}
	if result.BottomRightX == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("bottom_right_x", "result"))
	}
	if result.BottomRightY == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("bottom_right_y", "result"))
	}
	return
}

// ValidateTaskCollectionViewTiny runs the validations defined on
// TaskCollectionView using the "tiny" view.
func ValidateTaskCollectionViewTiny(result TaskCollectionView) (err error) {
	for _, item := range result {
		if err2 := ValidateTaskViewTiny(item); err2 != nil {
			err = goa.MergeErrors(err, err2)
		}
	}
	return
}

// ValidateTaskCollectionView runs the validations defined on
// TaskCollectionView using the "default" view.
func ValidateTaskCollectionView(result TaskCollectionView) (err error) {
	for _, item := range result {
		if err2 := ValidateTaskView(item); err2 != nil {
			err = goa.MergeErrors(err, err2)
		}
	}
	return
}

// ValidateImageView runs the validations defined on ImageView using the
// "default" view.
func ValidateImageView(result *ImageView) (err error) {
	if result.ImageID == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("image_id", "result"))
	}
	if result.ExpiresIn == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("expires_in", "result"))
	}
	if result.ImageID != nil {
		if utf8.RuneCountInString(*result.ImageID) < 8 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("result.image_id", *result.ImageID, utf8.RuneCountInString(*result.ImageID), 8, true))
		}
	}
	if result.ImageID != nil {
		if utf8.RuneCountInString(*result.ImageID) > 8 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("result.image_id", *result.ImageID, utf8.RuneCountInString(*result.ImageID), 8, false))
		}
	}
	if result.ExpiresIn != nil {
		err = goa.MergeErrors(err, goa.ValidateFormat("result.expires_in", *result.ExpiresIn, goa.FormatDateTime))
	}
	return
}
