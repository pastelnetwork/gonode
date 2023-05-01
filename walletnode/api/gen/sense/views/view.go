// Code generated by goa v3.7.6, DO NOT EDIT.
//
// sense views
//
// Command:
// $ goa gen github.com/pastelnetwork/gonode/walletnode/api/design

package views

import (
	"unicode/utf8"

	goa "goa.design/goa/v3/pkg"
)

// Image is the viewed result type that is projected based on a view.
type Image struct {
	// Type to project
	Projected *ImageView
	// View to render
	View string
}

// StartProcessingResult is the viewed result type that is projected based on a
// view.
type StartProcessingResult struct {
	// Type to project
	Projected *StartProcessingResultView
	// View to render
	View string
}

// ImageView is a type that runs validations on a projected type.
type ImageView struct {
	// Uploaded image ID
	ImageID *string
	// Image expiration
	ExpiresIn *string
	// Estimated fee
	TotalEstimatedFee *float64
	// The amount that's required to be preburned
	RequiredPreburnAmount *float64
}

// StartProcessingResultView is a type that runs validations on a projected
// type.
type StartProcessingResultView struct {
	// Task ID of processing task
	TaskID *string
}

var (
	// ImageMap is a map indexing the attribute names of Image by view name.
	ImageMap = map[string][]string{
		"default": {
			"image_id",
			"expires_in",
			"total_estimated_fee",
			"required_preburn_amount",
		},
	}
	// StartProcessingResultMap is a map indexing the attribute names of
	// StartProcessingResult by view name.
	StartProcessingResultMap = map[string][]string{
		"default": {
			"task_id",
		},
	}
)

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

// ValidateStartProcessingResult runs the validations defined on the viewed
// result type StartProcessingResult.
func ValidateStartProcessingResult(result *StartProcessingResult) (err error) {
	switch result.View {
	case "default", "":
		err = ValidateStartProcessingResultView(result.Projected)
	default:
		err = goa.InvalidEnumValueError("view", result.View, []interface{}{"default"})
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
	if result.TotalEstimatedFee == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("total_estimated_fee", "result"))
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
	if result.TotalEstimatedFee != nil {
		if *result.TotalEstimatedFee < 1e-05 {
			err = goa.MergeErrors(err, goa.InvalidRangeError("result.total_estimated_fee", *result.TotalEstimatedFee, 1e-05, true))
		}
	}
	if result.RequiredPreburnAmount != nil {
		if *result.RequiredPreburnAmount < 1e-05 {
			err = goa.MergeErrors(err, goa.InvalidRangeError("result.required_preburn_amount", *result.RequiredPreburnAmount, 1e-05, true))
		}
	}
	return
}

// ValidateStartProcessingResultView runs the validations defined on
// StartProcessingResultView using the "default" view.
func ValidateStartProcessingResultView(result *StartProcessingResultView) (err error) {
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
