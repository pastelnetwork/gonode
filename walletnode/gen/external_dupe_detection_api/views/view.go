// Code generated by goa v3.5.2, DO NOT EDIT.
//
// external_dupe_detection_api views
//
// Command:
// $ goa gen github.com/pastelnetwork/gonode/walletnode/api/design

package views

import (
	"unicode/utf8"

	goa "goa.design/goa/v3/pkg"
)

// Externaldupedetetionapiinitiatesubmissionresult is the viewed result type
// that is projected based on a view.
type Externaldupedetetionapiinitiatesubmissionresult struct {
	// Type to project
	Projected *ExternaldupedetetionapiinitiatesubmissionresultView
	// View to render
	View string
}

// ExternaldupedetetionapiinitiatesubmissionresultView is a type that runs
// validations on a projected type.
type ExternaldupedetetionapiinitiatesubmissionresultView struct {
	// SHA3-256 format, unique external dupe detection request id
	RequestID *string
}

var (
	// ExternaldupedetetionapiinitiatesubmissionresultMap is a map indexing the
	// attribute names of Externaldupedetetionapiinitiatesubmissionresult by view
	// name.
	ExternaldupedetetionapiinitiatesubmissionresultMap = map[string][]string{
		"default": {
			"request_id",
		},
	}
)

// ValidateExternaldupedetetionapiinitiatesubmissionresult runs the validations
// defined on the viewed result type
// Externaldupedetetionapiinitiatesubmissionresult.
func ValidateExternaldupedetetionapiinitiatesubmissionresult(result *Externaldupedetetionapiinitiatesubmissionresult) (err error) {
	switch result.View {
	case "default", "":
		err = ValidateExternaldupedetetionapiinitiatesubmissionresultView(result.Projected)
	default:
		err = goa.InvalidEnumValueError("view", result.View, []interface{}{"default"})
	}
	return
}

// ValidateExternaldupedetetionapiinitiatesubmissionresultView runs the
// validations defined on ExternaldupedetetionapiinitiatesubmissionresultView
// using the "default" view.
func ValidateExternaldupedetetionapiinitiatesubmissionresultView(result *ExternaldupedetetionapiinitiatesubmissionresultView) (err error) {
	if result.RequestID == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("request_id", "result"))
	}
	if result.RequestID != nil {
		if utf8.RuneCountInString(*result.RequestID) < 64 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("result.request_id", *result.RequestID, utf8.RuneCountInString(*result.RequestID), 64, true))
		}
	}
	if result.RequestID != nil {
		if utf8.RuneCountInString(*result.RequestID) > 64 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("result.request_id", *result.RequestID, utf8.RuneCountInString(*result.RequestID), 64, false))
		}
	}
	return
}
