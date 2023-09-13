// Code generated by goa v3.13.0, DO NOT EDIT.
//
// nft views
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

// ImageRes is the viewed result type that is projected based on a view.
type ImageRes struct {
	// Type to project
	Projected *ImageResView
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
	Ticket *NftRegisterPayloadView
}

// TaskStateView is a type that runs validations on a projected type.
type TaskStateView struct {
	// Date of the status creation
	Date *string
	// Status of the registration process
	Status *string
}

// NftRegisterPayloadView is a type that runs validations on a projected type.
type NftRegisterPayloadView struct {
	// Name of the NFT
	Name *string
	// Description of the NFT
	Description *string
	// Keywords
	Keywords *string
	// Series name
	SeriesName *string
	// Number of copies issued
	IssuedCopies *int
	// NFT creation video youtube URL
	YoutubeURL *string
	// Creator's PastelID
	CreatorPastelID *string
	// Name of the NFT creator
	CreatorName *string
	// NFT creator website URL
	CreatorWebsiteURL *string
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
	// To make it publicly accessible
	MakePubliclyAccessible *bool
	// Act Collection TxID to add given ticket in collection
	CollectionActTxid *string
	// OpenAPI GroupID string
	OpenAPIGroupID *string
	// Passphrase of the owner's PastelID
	Key *string
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

// ImageResView is a type that runs validations on a projected type.
type ImageResView struct {
	// Uploaded image ID
	ImageID *string
	// Image expiration
	ExpiresIn *string
	// Estimated fee
	EstimatedFee *float64
}

var (
	// RegisterResultMap is a map indexing the attribute names of RegisterResult by
	// view name.
	RegisterResultMap = map[string][]string{
		"default": {
			"task_id",
		},
	}
	// TaskMap is a map indexing the attribute names of Task by view name.
	TaskMap = map[string][]string{
		"tiny": {
			"id",
			"status",
			"txid",
			"ticket",
		},
		"default": {
			"id",
			"status",
			"states",
			"txid",
			"ticket",
		},
	}
	// TaskCollectionMap is a map indexing the attribute names of TaskCollection by
	// view name.
	TaskCollectionMap = map[string][]string{
		"tiny": {
			"id",
			"status",
			"txid",
			"ticket",
		},
		"default": {
			"id",
			"status",
			"states",
			"txid",
			"ticket",
		},
	}
	// ImageResMap is a map indexing the attribute names of ImageRes by view name.
	ImageResMap = map[string][]string{
		"default": {
			"image_id",
			"expires_in",
			"estimated_fee",
		},
	}
	// ThumbnailcoordinateMap is a map indexing the attribute names of
	// Thumbnailcoordinate by view name.
	ThumbnailcoordinateMap = map[string][]string{
		"default": {
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
		err = goa.InvalidEnumValueError("view", result.View, []any{"default"})
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
		err = goa.InvalidEnumValueError("view", result.View, []any{"tiny", "default"})
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
		err = goa.InvalidEnumValueError("view", result.View, []any{"tiny", "default"})
	}
	return
}

// ValidateImageRes runs the validations defined on the viewed result type
// ImageRes.
func ValidateImageRes(result *ImageRes) (err error) {
	switch result.View {
	case "default", "":
		err = ValidateImageResView(result.Projected)
	default:
		err = goa.InvalidEnumValueError("view", result.View, []any{"default"})
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
		if !(*result.Status == "Task Started" || *result.Status == "Connected" || *result.Status == "Validating Duplicate Reg Tickets" || *result.Status == "Validating Burn Txn" || *result.Status == "Burn Txn Validated" || *result.Status == "Image Probed" || *result.Status == "Image And Thumbnail Uploaded" || *result.Status == "Status Gen ReptorQ Symbols" || *result.Status == "Preburn Registration Fee" || *result.Status == "Downloaded" || *result.Status == "Request Accepted" || *result.Status == "Request Registered" || *result.Status == "Request Activated" || *result.Status == "Error Setting up mesh of supernodes" || *result.Status == "Error Sending Reg Metadata" || *result.Status == "Error Uploading Image" || *result.Status == "Error Converting Image to Bytes" || *result.Status == "Error Encoding Image" || *result.Status == "Error Creating Ticket" || *result.Status == "Error Signing Ticket" || *result.Status == "Error Uploading Ticket" || *result.Status == "Error Activating Ticket" || *result.Status == "Error Probing Image" || *result.Status == "Error checking dd-server availability before probe image" || *result.Status == "Error Generating DD and Fingerprint IDs" || *result.Status == "Error comparing suitable storage fee with task request maximum fee" || *result.Status == "Error balance not sufficient" || *result.Status == "Error getting hash of the image" || *result.Status == "Error sending signed ticket to SNs" || *result.Status == "Error checking balance" || *result.Status == "Error burning reg fee to get reg ticket id" || *result.Status == "Error validating reg ticket txn id" || *result.Status == "Error validating activate ticket txn id" || *result.Status == "Error Insufficient Fee" || *result.Status == "Error Signatures Dont Match" || *result.Status == "Error Fingerprints Dont Match" || *result.Status == "Error ThumbnailHashes Dont Match" || *result.Status == "Error GenRaptorQ Symbols Failed" || *result.Status == "Error File Don't Match" || *result.Status == "Error Not Enough SuperNode" || *result.Status == "Error Find Responding SNs" || *result.Status == "Error Not Enough Downloaded Filed" || *result.Status == "Error Download Failed" || *result.Status == "Error Invalid Burn TxID" || *result.Status == "Task Failed" || *result.Status == "Task Rejected" || *result.Status == "Task Completed") {
			err = goa.MergeErrors(err, goa.InvalidEnumValueError("result.status", *result.Status, []any{"Task Started", "Connected", "Validating Duplicate Reg Tickets", "Validating Burn Txn", "Burn Txn Validated", "Image Probed", "Image And Thumbnail Uploaded", "Status Gen ReptorQ Symbols", "Preburn Registration Fee", "Downloaded", "Request Accepted", "Request Registered", "Request Activated", "Error Setting up mesh of supernodes", "Error Sending Reg Metadata", "Error Uploading Image", "Error Converting Image to Bytes", "Error Encoding Image", "Error Creating Ticket", "Error Signing Ticket", "Error Uploading Ticket", "Error Activating Ticket", "Error Probing Image", "Error checking dd-server availability before probe image", "Error Generating DD and Fingerprint IDs", "Error comparing suitable storage fee with task request maximum fee", "Error balance not sufficient", "Error getting hash of the image", "Error sending signed ticket to SNs", "Error checking balance", "Error burning reg fee to get reg ticket id", "Error validating reg ticket txn id", "Error validating activate ticket txn id", "Error Insufficient Fee", "Error Signatures Dont Match", "Error Fingerprints Dont Match", "Error ThumbnailHashes Dont Match", "Error GenRaptorQ Symbols Failed", "Error File Don't Match", "Error Not Enough SuperNode", "Error Find Responding SNs", "Error Not Enough Downloaded Filed", "Error Download Failed", "Error Invalid Burn TxID", "Task Failed", "Task Rejected", "Task Completed"}))
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
		if err2 := ValidateNftRegisterPayloadView(result.Ticket); err2 != nil {
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
		if !(*result.Status == "Task Started" || *result.Status == "Connected" || *result.Status == "Validating Duplicate Reg Tickets" || *result.Status == "Validating Burn Txn" || *result.Status == "Burn Txn Validated" || *result.Status == "Image Probed" || *result.Status == "Image And Thumbnail Uploaded" || *result.Status == "Status Gen ReptorQ Symbols" || *result.Status == "Preburn Registration Fee" || *result.Status == "Downloaded" || *result.Status == "Request Accepted" || *result.Status == "Request Registered" || *result.Status == "Request Activated" || *result.Status == "Error Setting up mesh of supernodes" || *result.Status == "Error Sending Reg Metadata" || *result.Status == "Error Uploading Image" || *result.Status == "Error Converting Image to Bytes" || *result.Status == "Error Encoding Image" || *result.Status == "Error Creating Ticket" || *result.Status == "Error Signing Ticket" || *result.Status == "Error Uploading Ticket" || *result.Status == "Error Activating Ticket" || *result.Status == "Error Probing Image" || *result.Status == "Error checking dd-server availability before probe image" || *result.Status == "Error Generating DD and Fingerprint IDs" || *result.Status == "Error comparing suitable storage fee with task request maximum fee" || *result.Status == "Error balance not sufficient" || *result.Status == "Error getting hash of the image" || *result.Status == "Error sending signed ticket to SNs" || *result.Status == "Error checking balance" || *result.Status == "Error burning reg fee to get reg ticket id" || *result.Status == "Error validating reg ticket txn id" || *result.Status == "Error validating activate ticket txn id" || *result.Status == "Error Insufficient Fee" || *result.Status == "Error Signatures Dont Match" || *result.Status == "Error Fingerprints Dont Match" || *result.Status == "Error ThumbnailHashes Dont Match" || *result.Status == "Error GenRaptorQ Symbols Failed" || *result.Status == "Error File Don't Match" || *result.Status == "Error Not Enough SuperNode" || *result.Status == "Error Find Responding SNs" || *result.Status == "Error Not Enough Downloaded Filed" || *result.Status == "Error Download Failed" || *result.Status == "Error Invalid Burn TxID" || *result.Status == "Task Failed" || *result.Status == "Task Rejected" || *result.Status == "Task Completed") {
			err = goa.MergeErrors(err, goa.InvalidEnumValueError("result.status", *result.Status, []any{"Task Started", "Connected", "Validating Duplicate Reg Tickets", "Validating Burn Txn", "Burn Txn Validated", "Image Probed", "Image And Thumbnail Uploaded", "Status Gen ReptorQ Symbols", "Preburn Registration Fee", "Downloaded", "Request Accepted", "Request Registered", "Request Activated", "Error Setting up mesh of supernodes", "Error Sending Reg Metadata", "Error Uploading Image", "Error Converting Image to Bytes", "Error Encoding Image", "Error Creating Ticket", "Error Signing Ticket", "Error Uploading Ticket", "Error Activating Ticket", "Error Probing Image", "Error checking dd-server availability before probe image", "Error Generating DD and Fingerprint IDs", "Error comparing suitable storage fee with task request maximum fee", "Error balance not sufficient", "Error getting hash of the image", "Error sending signed ticket to SNs", "Error checking balance", "Error burning reg fee to get reg ticket id", "Error validating reg ticket txn id", "Error validating activate ticket txn id", "Error Insufficient Fee", "Error Signatures Dont Match", "Error Fingerprints Dont Match", "Error ThumbnailHashes Dont Match", "Error GenRaptorQ Symbols Failed", "Error File Don't Match", "Error Not Enough SuperNode", "Error Find Responding SNs", "Error Not Enough Downloaded Filed", "Error Download Failed", "Error Invalid Burn TxID", "Task Failed", "Task Rejected", "Task Completed"}))
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
		if err2 := ValidateNftRegisterPayloadView(result.Ticket); err2 != nil {
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
		if !(*result.Status == "Task Started" || *result.Status == "Connected" || *result.Status == "Validating Duplicate Reg Tickets" || *result.Status == "Validating Burn Txn" || *result.Status == "Burn Txn Validated" || *result.Status == "Image Probed" || *result.Status == "Image And Thumbnail Uploaded" || *result.Status == "Status Gen ReptorQ Symbols" || *result.Status == "Preburn Registration Fee" || *result.Status == "Downloaded" || *result.Status == "Request Accepted" || *result.Status == "Request Registered" || *result.Status == "Request Activated" || *result.Status == "Error Setting up mesh of supernodes" || *result.Status == "Error Sending Reg Metadata" || *result.Status == "Error Uploading Image" || *result.Status == "Error Converting Image to Bytes" || *result.Status == "Error Encoding Image" || *result.Status == "Error Creating Ticket" || *result.Status == "Error Signing Ticket" || *result.Status == "Error Uploading Ticket" || *result.Status == "Error Activating Ticket" || *result.Status == "Error Probing Image" || *result.Status == "Error checking dd-server availability before probe image" || *result.Status == "Error Generating DD and Fingerprint IDs" || *result.Status == "Error comparing suitable storage fee with task request maximum fee" || *result.Status == "Error balance not sufficient" || *result.Status == "Error getting hash of the image" || *result.Status == "Error sending signed ticket to SNs" || *result.Status == "Error checking balance" || *result.Status == "Error burning reg fee to get reg ticket id" || *result.Status == "Error validating reg ticket txn id" || *result.Status == "Error validating activate ticket txn id" || *result.Status == "Error Insufficient Fee" || *result.Status == "Error Signatures Dont Match" || *result.Status == "Error Fingerprints Dont Match" || *result.Status == "Error ThumbnailHashes Dont Match" || *result.Status == "Error GenRaptorQ Symbols Failed" || *result.Status == "Error File Don't Match" || *result.Status == "Error Not Enough SuperNode" || *result.Status == "Error Find Responding SNs" || *result.Status == "Error Not Enough Downloaded Filed" || *result.Status == "Error Download Failed" || *result.Status == "Error Invalid Burn TxID" || *result.Status == "Task Failed" || *result.Status == "Task Rejected" || *result.Status == "Task Completed") {
			err = goa.MergeErrors(err, goa.InvalidEnumValueError("result.status", *result.Status, []any{"Task Started", "Connected", "Validating Duplicate Reg Tickets", "Validating Burn Txn", "Burn Txn Validated", "Image Probed", "Image And Thumbnail Uploaded", "Status Gen ReptorQ Symbols", "Preburn Registration Fee", "Downloaded", "Request Accepted", "Request Registered", "Request Activated", "Error Setting up mesh of supernodes", "Error Sending Reg Metadata", "Error Uploading Image", "Error Converting Image to Bytes", "Error Encoding Image", "Error Creating Ticket", "Error Signing Ticket", "Error Uploading Ticket", "Error Activating Ticket", "Error Probing Image", "Error checking dd-server availability before probe image", "Error Generating DD and Fingerprint IDs", "Error comparing suitable storage fee with task request maximum fee", "Error balance not sufficient", "Error getting hash of the image", "Error sending signed ticket to SNs", "Error checking balance", "Error burning reg fee to get reg ticket id", "Error validating reg ticket txn id", "Error validating activate ticket txn id", "Error Insufficient Fee", "Error Signatures Dont Match", "Error Fingerprints Dont Match", "Error ThumbnailHashes Dont Match", "Error GenRaptorQ Symbols Failed", "Error File Don't Match", "Error Not Enough SuperNode", "Error Find Responding SNs", "Error Not Enough Downloaded Filed", "Error Download Failed", "Error Invalid Burn TxID", "Task Failed", "Task Rejected", "Task Completed"}))
		}
	}
	return
}

// ValidateNftRegisterPayloadView runs the validations defined on
// NftRegisterPayloadView.
func ValidateNftRegisterPayloadView(result *NftRegisterPayloadView) (err error) {
	if result.CreatorName == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("creator_name", "result"))
	}
	if result.Name == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("name", "result"))
	}
	if result.CreatorPastelID == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("creator_pastelid", "result"))
	}
	if result.SpendableAddress == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("spendable_address", "result"))
	}
	if result.MaximumFee == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("maximum_fee", "result"))
	}
	if result.Key == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("key", "result"))
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
		if *result.IssuedCopies > 1000 {
			err = goa.MergeErrors(err, goa.InvalidRangeError("result.issued_copies", *result.IssuedCopies, 1000, false))
		}
	}
	if result.YoutubeURL != nil {
		if utf8.RuneCountInString(*result.YoutubeURL) > 128 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("result.youtube_url", *result.YoutubeURL, utf8.RuneCountInString(*result.YoutubeURL), 128, false))
		}
	}
	if result.CreatorPastelID != nil {
		err = goa.MergeErrors(err, goa.ValidatePattern("result.creator_pastelid", *result.CreatorPastelID, "^[a-zA-Z0-9]+$"))
	}
	if result.CreatorPastelID != nil {
		if utf8.RuneCountInString(*result.CreatorPastelID) < 86 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("result.creator_pastelid", *result.CreatorPastelID, utf8.RuneCountInString(*result.CreatorPastelID), 86, true))
		}
	}
	if result.CreatorPastelID != nil {
		if utf8.RuneCountInString(*result.CreatorPastelID) > 86 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("result.creator_pastelid", *result.CreatorPastelID, utf8.RuneCountInString(*result.CreatorPastelID), 86, false))
		}
	}
	if result.CreatorName != nil {
		if utf8.RuneCountInString(*result.CreatorName) > 256 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("result.creator_name", *result.CreatorName, utf8.RuneCountInString(*result.CreatorName), 256, false))
		}
	}
	if result.CreatorWebsiteURL != nil {
		if utf8.RuneCountInString(*result.CreatorWebsiteURL) > 256 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("result.creator_website_url", *result.CreatorWebsiteURL, utf8.RuneCountInString(*result.CreatorWebsiteURL), 256, false))
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
		if *result.Royalty > 20 {
			err = goa.MergeErrors(err, goa.InvalidRangeError("result.royalty", *result.Royalty, 20, false))
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

// ValidateImageResView runs the validations defined on ImageResView using the
// "default" view.
func ValidateImageResView(result *ImageResView) (err error) {
	if result.ImageID == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("image_id", "result"))
	}
	if result.ExpiresIn == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("expires_in", "result"))
	}
	if result.EstimatedFee == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("estimated_fee", "result"))
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
	if result.EstimatedFee != nil {
		if *result.EstimatedFee < 1e-05 {
			err = goa.MergeErrors(err, goa.InvalidRangeError("result.estimated_fee", *result.EstimatedFee, 1e-05, true))
		}
	}
	return
}
