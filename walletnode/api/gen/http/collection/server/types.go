// Code generated by goa v3.15.0, DO NOT EDIT.
//
// collection HTTP server types
//
// Command:
// $ goa gen github.com/pastelnetwork/gonode/walletnode/api/design

package server

import (
	"unicode/utf8"

	collection "github.com/pastelnetwork/gonode/walletnode/api/gen/collection"
	collectionviews "github.com/pastelnetwork/gonode/walletnode/api/gen/collection/views"
	goa "goa.design/goa/v3/pkg"
)

// RegisterCollectionRequestBody is the type of the "collection" service
// "registerCollection" endpoint HTTP request body.
type RegisterCollectionRequestBody struct {
	// name of the collection
	CollectionName *string `form:"collection_name,omitempty" json:"collection_name,omitempty" xml:"collection_name,omitempty"`
	// type of items, store by collection
	ItemType *string `form:"item_type,omitempty" json:"item_type,omitempty" xml:"item_type,omitempty"`
	// list of authorized contributors
	ListOfPastelidsOfAuthorizedContributors []string `form:"list_of_pastelids_of_authorized_contributors,omitempty" json:"list_of_pastelids_of_authorized_contributors,omitempty" xml:"list_of_pastelids_of_authorized_contributors,omitempty"`
	// max no of entries in the collection
	MaxCollectionEntries *int `form:"max_collection_entries,omitempty" json:"max_collection_entries,omitempty" xml:"max_collection_entries,omitempty"`
	// no of days to finalize collection
	NoOfDaysToFinalizeCollection *int `form:"no_of_days_to_finalize_collection,omitempty" json:"no_of_days_to_finalize_collection,omitempty" xml:"no_of_days_to_finalize_collection,omitempty"`
	// item copy count in the collection
	CollectionItemCopyCount *int `form:"collection_item_copy_count,omitempty" json:"collection_item_copy_count,omitempty" xml:"collection_item_copy_count,omitempty"`
	// royalty fee
	Royalty *float64 `form:"royalty,omitempty" json:"royalty,omitempty" xml:"royalty,omitempty"`
	// green
	Green *bool `form:"green,omitempty" json:"green,omitempty" xml:"green,omitempty"`
	// max open nfsw score sense and nft items can have
	MaxPermittedOpenNsfwScore *float64 `form:"max_permitted_open_nsfw_score,omitempty" json:"max_permitted_open_nsfw_score,omitempty" xml:"max_permitted_open_nsfw_score,omitempty"`
	// min similarity for 1st entry to have
	MinimumSimilarityScoreToFirstEntryInCollection *float64 `form:"minimum_similarity_score_to_first_entry_in_collection,omitempty" json:"minimum_similarity_score_to_first_entry_in_collection,omitempty" xml:"minimum_similarity_score_to_first_entry_in_collection,omitempty"`
	// App PastelID
	AppPastelID *string `form:"app_pastelid,omitempty" json:"app_pastelid,omitempty" xml:"app_pastelid,omitempty"`
	// Spendable address
	SpendableAddress *string `form:"spendable_address,omitempty" json:"spendable_address,omitempty" xml:"spendable_address,omitempty"`
}

// RegisterCollectionResponseBody is the type of the "collection" service
// "registerCollection" endpoint HTTP response body.
type RegisterCollectionResponseBody struct {
	// Uploaded file ID
	TaskID string `form:"task_id" json:"task_id" xml:"task_id"`
}

// RegisterTaskStateResponseBody is the type of the "collection" service
// "registerTaskState" endpoint HTTP response body.
type RegisterTaskStateResponseBody struct {
	// Date of the status creation
	Date string `form:"date" json:"date" xml:"date"`
	// Status of the registration process
	Status string `form:"status" json:"status" xml:"status"`
}

// GetTaskHistoryResponseBody is the type of the "collection" service
// "getTaskHistory" endpoint HTTP response body.
type GetTaskHistoryResponseBody []*TaskHistoryResponse

// RegisterCollectionUnAuthorizedResponseBody is the type of the "collection"
// service "registerCollection" endpoint HTTP response body for the
// "UnAuthorized" error.
type RegisterCollectionUnAuthorizedResponseBody struct {
	// Name is the name of this class of errors.
	Name string `form:"name" json:"name" xml:"name"`
	// ID is a unique identifier for this particular occurrence of the problem.
	ID string `form:"id" json:"id" xml:"id"`
	// Message is a human-readable explanation specific to this occurrence of the
	// problem.
	Message string `form:"message" json:"message" xml:"message"`
	// Is the error temporary?
	Temporary bool `form:"temporary" json:"temporary" xml:"temporary"`
	// Is the error a timeout?
	Timeout bool `form:"timeout" json:"timeout" xml:"timeout"`
	// Is the error a server-side fault?
	Fault bool `form:"fault" json:"fault" xml:"fault"`
}

// RegisterCollectionBadRequestResponseBody is the type of the "collection"
// service "registerCollection" endpoint HTTP response body for the
// "BadRequest" error.
type RegisterCollectionBadRequestResponseBody struct {
	// Name is the name of this class of errors.
	Name string `form:"name" json:"name" xml:"name"`
	// ID is a unique identifier for this particular occurrence of the problem.
	ID string `form:"id" json:"id" xml:"id"`
	// Message is a human-readable explanation specific to this occurrence of the
	// problem.
	Message string `form:"message" json:"message" xml:"message"`
	// Is the error temporary?
	Temporary bool `form:"temporary" json:"temporary" xml:"temporary"`
	// Is the error a timeout?
	Timeout bool `form:"timeout" json:"timeout" xml:"timeout"`
	// Is the error a server-side fault?
	Fault bool `form:"fault" json:"fault" xml:"fault"`
}

// RegisterCollectionNotFoundResponseBody is the type of the "collection"
// service "registerCollection" endpoint HTTP response body for the "NotFound"
// error.
type RegisterCollectionNotFoundResponseBody struct {
	// Name is the name of this class of errors.
	Name string `form:"name" json:"name" xml:"name"`
	// ID is a unique identifier for this particular occurrence of the problem.
	ID string `form:"id" json:"id" xml:"id"`
	// Message is a human-readable explanation specific to this occurrence of the
	// problem.
	Message string `form:"message" json:"message" xml:"message"`
	// Is the error temporary?
	Temporary bool `form:"temporary" json:"temporary" xml:"temporary"`
	// Is the error a timeout?
	Timeout bool `form:"timeout" json:"timeout" xml:"timeout"`
	// Is the error a server-side fault?
	Fault bool `form:"fault" json:"fault" xml:"fault"`
}

// RegisterCollectionInternalServerErrorResponseBody is the type of the
// "collection" service "registerCollection" endpoint HTTP response body for
// the "InternalServerError" error.
type RegisterCollectionInternalServerErrorResponseBody struct {
	// Name is the name of this class of errors.
	Name string `form:"name" json:"name" xml:"name"`
	// ID is a unique identifier for this particular occurrence of the problem.
	ID string `form:"id" json:"id" xml:"id"`
	// Message is a human-readable explanation specific to this occurrence of the
	// problem.
	Message string `form:"message" json:"message" xml:"message"`
	// Is the error temporary?
	Temporary bool `form:"temporary" json:"temporary" xml:"temporary"`
	// Is the error a timeout?
	Timeout bool `form:"timeout" json:"timeout" xml:"timeout"`
	// Is the error a server-side fault?
	Fault bool `form:"fault" json:"fault" xml:"fault"`
}

// RegisterTaskStateNotFoundResponseBody is the type of the "collection"
// service "registerTaskState" endpoint HTTP response body for the "NotFound"
// error.
type RegisterTaskStateNotFoundResponseBody struct {
	// Name is the name of this class of errors.
	Name string `form:"name" json:"name" xml:"name"`
	// ID is a unique identifier for this particular occurrence of the problem.
	ID string `form:"id" json:"id" xml:"id"`
	// Message is a human-readable explanation specific to this occurrence of the
	// problem.
	Message string `form:"message" json:"message" xml:"message"`
	// Is the error temporary?
	Temporary bool `form:"temporary" json:"temporary" xml:"temporary"`
	// Is the error a timeout?
	Timeout bool `form:"timeout" json:"timeout" xml:"timeout"`
	// Is the error a server-side fault?
	Fault bool `form:"fault" json:"fault" xml:"fault"`
}

// RegisterTaskStateInternalServerErrorResponseBody is the type of the
// "collection" service "registerTaskState" endpoint HTTP response body for the
// "InternalServerError" error.
type RegisterTaskStateInternalServerErrorResponseBody struct {
	// Name is the name of this class of errors.
	Name string `form:"name" json:"name" xml:"name"`
	// ID is a unique identifier for this particular occurrence of the problem.
	ID string `form:"id" json:"id" xml:"id"`
	// Message is a human-readable explanation specific to this occurrence of the
	// problem.
	Message string `form:"message" json:"message" xml:"message"`
	// Is the error temporary?
	Temporary bool `form:"temporary" json:"temporary" xml:"temporary"`
	// Is the error a timeout?
	Timeout bool `form:"timeout" json:"timeout" xml:"timeout"`
	// Is the error a server-side fault?
	Fault bool `form:"fault" json:"fault" xml:"fault"`
}

// GetTaskHistoryNotFoundResponseBody is the type of the "collection" service
// "getTaskHistory" endpoint HTTP response body for the "NotFound" error.
type GetTaskHistoryNotFoundResponseBody struct {
	// Name is the name of this class of errors.
	Name string `form:"name" json:"name" xml:"name"`
	// ID is a unique identifier for this particular occurrence of the problem.
	ID string `form:"id" json:"id" xml:"id"`
	// Message is a human-readable explanation specific to this occurrence of the
	// problem.
	Message string `form:"message" json:"message" xml:"message"`
	// Is the error temporary?
	Temporary bool `form:"temporary" json:"temporary" xml:"temporary"`
	// Is the error a timeout?
	Timeout bool `form:"timeout" json:"timeout" xml:"timeout"`
	// Is the error a server-side fault?
	Fault bool `form:"fault" json:"fault" xml:"fault"`
}

// GetTaskHistoryInternalServerErrorResponseBody is the type of the
// "collection" service "getTaskHistory" endpoint HTTP response body for the
// "InternalServerError" error.
type GetTaskHistoryInternalServerErrorResponseBody struct {
	// Name is the name of this class of errors.
	Name string `form:"name" json:"name" xml:"name"`
	// ID is a unique identifier for this particular occurrence of the problem.
	ID string `form:"id" json:"id" xml:"id"`
	// Message is a human-readable explanation specific to this occurrence of the
	// problem.
	Message string `form:"message" json:"message" xml:"message"`
	// Is the error temporary?
	Temporary bool `form:"temporary" json:"temporary" xml:"temporary"`
	// Is the error a timeout?
	Timeout bool `form:"timeout" json:"timeout" xml:"timeout"`
	// Is the error a server-side fault?
	Fault bool `form:"fault" json:"fault" xml:"fault"`
}

// TaskHistoryResponse is used to define fields on response body types.
type TaskHistoryResponse struct {
	// Timestamp of the status creation
	Timestamp *string `form:"timestamp,omitempty" json:"timestamp,omitempty" xml:"timestamp,omitempty"`
	// past status string
	Status string `form:"status" json:"status" xml:"status"`
	// message string (if any)
	Message *string `form:"message,omitempty" json:"message,omitempty" xml:"message,omitempty"`
	// details of the status
	Details *DetailsResponse `form:"details,omitempty" json:"details,omitempty" xml:"details,omitempty"`
}

// DetailsResponse is used to define fields on response body types.
type DetailsResponse struct {
	// details regarding the status
	Message *string `form:"message,omitempty" json:"message,omitempty" xml:"message,omitempty"`
	// important fields regarding status history
	Fields map[string]any `form:"fields,omitempty" json:"fields,omitempty" xml:"fields,omitempty"`
}

// NewRegisterCollectionResponseBody builds the HTTP response body from the
// result of the "registerCollection" endpoint of the "collection" service.
func NewRegisterCollectionResponseBody(res *collectionviews.RegisterCollectionResponseView) *RegisterCollectionResponseBody {
	body := &RegisterCollectionResponseBody{
		TaskID: *res.TaskID,
	}
	return body
}

// NewRegisterTaskStateResponseBody builds the HTTP response body from the
// result of the "registerTaskState" endpoint of the "collection" service.
func NewRegisterTaskStateResponseBody(res *collection.TaskState) *RegisterTaskStateResponseBody {
	body := &RegisterTaskStateResponseBody{
		Date:   res.Date,
		Status: res.Status,
	}
	return body
}

// NewGetTaskHistoryResponseBody builds the HTTP response body from the result
// of the "getTaskHistory" endpoint of the "collection" service.
func NewGetTaskHistoryResponseBody(res []*collection.TaskHistory) GetTaskHistoryResponseBody {
	body := make([]*TaskHistoryResponse, len(res))
	for i, val := range res {
		body[i] = marshalCollectionTaskHistoryToTaskHistoryResponse(val)
	}
	return body
}

// NewRegisterCollectionUnAuthorizedResponseBody builds the HTTP response body
// from the result of the "registerCollection" endpoint of the "collection"
// service.
func NewRegisterCollectionUnAuthorizedResponseBody(res *goa.ServiceError) *RegisterCollectionUnAuthorizedResponseBody {
	body := &RegisterCollectionUnAuthorizedResponseBody{
		Name:      res.Name,
		ID:        res.ID,
		Message:   res.Message,
		Temporary: res.Temporary,
		Timeout:   res.Timeout,
		Fault:     res.Fault,
	}
	return body
}

// NewRegisterCollectionBadRequestResponseBody builds the HTTP response body
// from the result of the "registerCollection" endpoint of the "collection"
// service.
func NewRegisterCollectionBadRequestResponseBody(res *goa.ServiceError) *RegisterCollectionBadRequestResponseBody {
	body := &RegisterCollectionBadRequestResponseBody{
		Name:      res.Name,
		ID:        res.ID,
		Message:   res.Message,
		Temporary: res.Temporary,
		Timeout:   res.Timeout,
		Fault:     res.Fault,
	}
	return body
}

// NewRegisterCollectionNotFoundResponseBody builds the HTTP response body from
// the result of the "registerCollection" endpoint of the "collection" service.
func NewRegisterCollectionNotFoundResponseBody(res *goa.ServiceError) *RegisterCollectionNotFoundResponseBody {
	body := &RegisterCollectionNotFoundResponseBody{
		Name:      res.Name,
		ID:        res.ID,
		Message:   res.Message,
		Temporary: res.Temporary,
		Timeout:   res.Timeout,
		Fault:     res.Fault,
	}
	return body
}

// NewRegisterCollectionInternalServerErrorResponseBody builds the HTTP
// response body from the result of the "registerCollection" endpoint of the
// "collection" service.
func NewRegisterCollectionInternalServerErrorResponseBody(res *goa.ServiceError) *RegisterCollectionInternalServerErrorResponseBody {
	body := &RegisterCollectionInternalServerErrorResponseBody{
		Name:      res.Name,
		ID:        res.ID,
		Message:   res.Message,
		Temporary: res.Temporary,
		Timeout:   res.Timeout,
		Fault:     res.Fault,
	}
	return body
}

// NewRegisterTaskStateNotFoundResponseBody builds the HTTP response body from
// the result of the "registerTaskState" endpoint of the "collection" service.
func NewRegisterTaskStateNotFoundResponseBody(res *goa.ServiceError) *RegisterTaskStateNotFoundResponseBody {
	body := &RegisterTaskStateNotFoundResponseBody{
		Name:      res.Name,
		ID:        res.ID,
		Message:   res.Message,
		Temporary: res.Temporary,
		Timeout:   res.Timeout,
		Fault:     res.Fault,
	}
	return body
}

// NewRegisterTaskStateInternalServerErrorResponseBody builds the HTTP response
// body from the result of the "registerTaskState" endpoint of the "collection"
// service.
func NewRegisterTaskStateInternalServerErrorResponseBody(res *goa.ServiceError) *RegisterTaskStateInternalServerErrorResponseBody {
	body := &RegisterTaskStateInternalServerErrorResponseBody{
		Name:      res.Name,
		ID:        res.ID,
		Message:   res.Message,
		Temporary: res.Temporary,
		Timeout:   res.Timeout,
		Fault:     res.Fault,
	}
	return body
}

// NewGetTaskHistoryNotFoundResponseBody builds the HTTP response body from the
// result of the "getTaskHistory" endpoint of the "collection" service.
func NewGetTaskHistoryNotFoundResponseBody(res *goa.ServiceError) *GetTaskHistoryNotFoundResponseBody {
	body := &GetTaskHistoryNotFoundResponseBody{
		Name:      res.Name,
		ID:        res.ID,
		Message:   res.Message,
		Temporary: res.Temporary,
		Timeout:   res.Timeout,
		Fault:     res.Fault,
	}
	return body
}

// NewGetTaskHistoryInternalServerErrorResponseBody builds the HTTP response
// body from the result of the "getTaskHistory" endpoint of the "collection"
// service.
func NewGetTaskHistoryInternalServerErrorResponseBody(res *goa.ServiceError) *GetTaskHistoryInternalServerErrorResponseBody {
	body := &GetTaskHistoryInternalServerErrorResponseBody{
		Name:      res.Name,
		ID:        res.ID,
		Message:   res.Message,
		Temporary: res.Temporary,
		Timeout:   res.Timeout,
		Fault:     res.Fault,
	}
	return body
}

// NewRegisterCollectionPayload builds a collection service registerCollection
// endpoint payload.
func NewRegisterCollectionPayload(body *RegisterCollectionRequestBody, key *string) *collection.RegisterCollectionPayload {
	v := &collection.RegisterCollectionPayload{
		CollectionName:            *body.CollectionName,
		ItemType:                  *body.ItemType,
		MaxCollectionEntries:      *body.MaxCollectionEntries,
		MaxPermittedOpenNsfwScore: *body.MaxPermittedOpenNsfwScore,
		MinimumSimilarityScoreToFirstEntryInCollection: *body.MinimumSimilarityScoreToFirstEntryInCollection,
		AppPastelID:      *body.AppPastelID,
		SpendableAddress: *body.SpendableAddress,
	}
	if body.NoOfDaysToFinalizeCollection != nil {
		v.NoOfDaysToFinalizeCollection = *body.NoOfDaysToFinalizeCollection
	}
	if body.CollectionItemCopyCount != nil {
		v.CollectionItemCopyCount = *body.CollectionItemCopyCount
	}
	if body.Royalty != nil {
		v.Royalty = *body.Royalty
	}
	if body.Green != nil {
		v.Green = *body.Green
	}
	v.ListOfPastelidsOfAuthorizedContributors = make([]string, len(body.ListOfPastelidsOfAuthorizedContributors))
	for i, val := range body.ListOfPastelidsOfAuthorizedContributors {
		v.ListOfPastelidsOfAuthorizedContributors[i] = val
	}
	if body.NoOfDaysToFinalizeCollection == nil {
		v.NoOfDaysToFinalizeCollection = 7
	}
	if body.CollectionItemCopyCount == nil {
		v.CollectionItemCopyCount = 1
	}
	if body.Royalty == nil {
		v.Royalty = 0
	}
	if body.Green == nil {
		v.Green = false
	}
	v.Key = key

	return v
}

// NewRegisterTaskStatePayload builds a collection service registerTaskState
// endpoint payload.
func NewRegisterTaskStatePayload(taskID string) *collection.RegisterTaskStatePayload {
	v := &collection.RegisterTaskStatePayload{}
	v.TaskID = taskID

	return v
}

// NewGetTaskHistoryPayload builds a collection service getTaskHistory endpoint
// payload.
func NewGetTaskHistoryPayload(taskID string) *collection.GetTaskHistoryPayload {
	v := &collection.GetTaskHistoryPayload{}
	v.TaskID = taskID

	return v
}

// ValidateRegisterCollectionRequestBody runs the validations defined on
// RegisterCollectionRequestBody
func ValidateRegisterCollectionRequestBody(body *RegisterCollectionRequestBody) (err error) {
	if body.CollectionName == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("collection_name", "body"))
	}
	if body.ItemType == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("item_type", "body"))
	}
	if body.ListOfPastelidsOfAuthorizedContributors == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("list_of_pastelids_of_authorized_contributors", "body"))
	}
	if body.MaxCollectionEntries == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("max_collection_entries", "body"))
	}
	if body.MaxPermittedOpenNsfwScore == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("max_permitted_open_nsfw_score", "body"))
	}
	if body.MinimumSimilarityScoreToFirstEntryInCollection == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("minimum_similarity_score_to_first_entry_in_collection", "body"))
	}
	if body.AppPastelID == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("app_pastelid", "body"))
	}
	if body.SpendableAddress == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("spendable_address", "body"))
	}
	if body.ItemType != nil {
		if !(*body.ItemType == "sense" || *body.ItemType == "nft") {
			err = goa.MergeErrors(err, goa.InvalidEnumValueError("body.item_type", *body.ItemType, []any{"sense", "nft"}))
		}
	}
	if body.MaxCollectionEntries != nil {
		if *body.MaxCollectionEntries < 1 {
			err = goa.MergeErrors(err, goa.InvalidRangeError("body.max_collection_entries", *body.MaxCollectionEntries, 1, true))
		}
	}
	if body.MaxCollectionEntries != nil {
		if *body.MaxCollectionEntries > 10000 {
			err = goa.MergeErrors(err, goa.InvalidRangeError("body.max_collection_entries", *body.MaxCollectionEntries, 10000, false))
		}
	}
	if body.NoOfDaysToFinalizeCollection != nil {
		if *body.NoOfDaysToFinalizeCollection < 1 {
			err = goa.MergeErrors(err, goa.InvalidRangeError("body.no_of_days_to_finalize_collection", *body.NoOfDaysToFinalizeCollection, 1, true))
		}
	}
	if body.NoOfDaysToFinalizeCollection != nil {
		if *body.NoOfDaysToFinalizeCollection > 7 {
			err = goa.MergeErrors(err, goa.InvalidRangeError("body.no_of_days_to_finalize_collection", *body.NoOfDaysToFinalizeCollection, 7, false))
		}
	}
	if body.CollectionItemCopyCount != nil {
		if *body.CollectionItemCopyCount < 1 {
			err = goa.MergeErrors(err, goa.InvalidRangeError("body.collection_item_copy_count", *body.CollectionItemCopyCount, 1, true))
		}
	}
	if body.CollectionItemCopyCount != nil {
		if *body.CollectionItemCopyCount > 1000 {
			err = goa.MergeErrors(err, goa.InvalidRangeError("body.collection_item_copy_count", *body.CollectionItemCopyCount, 1000, false))
		}
	}
	if body.Royalty != nil {
		if *body.Royalty < 0 {
			err = goa.MergeErrors(err, goa.InvalidRangeError("body.royalty", *body.Royalty, 0, true))
		}
	}
	if body.Royalty != nil {
		if *body.Royalty > 20 {
			err = goa.MergeErrors(err, goa.InvalidRangeError("body.royalty", *body.Royalty, 20, false))
		}
	}
	if body.MaxPermittedOpenNsfwScore != nil {
		if *body.MaxPermittedOpenNsfwScore < 0 {
			err = goa.MergeErrors(err, goa.InvalidRangeError("body.max_permitted_open_nsfw_score", *body.MaxPermittedOpenNsfwScore, 0, true))
		}
	}
	if body.MaxPermittedOpenNsfwScore != nil {
		if *body.MaxPermittedOpenNsfwScore > 1 {
			err = goa.MergeErrors(err, goa.InvalidRangeError("body.max_permitted_open_nsfw_score", *body.MaxPermittedOpenNsfwScore, 1, false))
		}
	}
	if body.MinimumSimilarityScoreToFirstEntryInCollection != nil {
		if *body.MinimumSimilarityScoreToFirstEntryInCollection < 0 {
			err = goa.MergeErrors(err, goa.InvalidRangeError("body.minimum_similarity_score_to_first_entry_in_collection", *body.MinimumSimilarityScoreToFirstEntryInCollection, 0, true))
		}
	}
	if body.MinimumSimilarityScoreToFirstEntryInCollection != nil {
		if *body.MinimumSimilarityScoreToFirstEntryInCollection > 1 {
			err = goa.MergeErrors(err, goa.InvalidRangeError("body.minimum_similarity_score_to_first_entry_in_collection", *body.MinimumSimilarityScoreToFirstEntryInCollection, 1, false))
		}
	}
	if body.AppPastelID != nil {
		err = goa.MergeErrors(err, goa.ValidatePattern("body.app_pastelid", *body.AppPastelID, "^[a-zA-Z0-9]+$"))
	}
	if body.AppPastelID != nil {
		if utf8.RuneCountInString(*body.AppPastelID) < 86 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("body.app_pastelid", *body.AppPastelID, utf8.RuneCountInString(*body.AppPastelID), 86, true))
		}
	}
	if body.AppPastelID != nil {
		if utf8.RuneCountInString(*body.AppPastelID) > 86 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("body.app_pastelid", *body.AppPastelID, utf8.RuneCountInString(*body.AppPastelID), 86, false))
		}
	}
	if body.SpendableAddress != nil {
		err = goa.MergeErrors(err, goa.ValidatePattern("body.spendable_address", *body.SpendableAddress, "^[a-zA-Z0-9]+$"))
	}
	if body.SpendableAddress != nil {
		if utf8.RuneCountInString(*body.SpendableAddress) < 35 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("body.spendable_address", *body.SpendableAddress, utf8.RuneCountInString(*body.SpendableAddress), 35, true))
		}
	}
	if body.SpendableAddress != nil {
		if utf8.RuneCountInString(*body.SpendableAddress) > 35 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("body.spendable_address", *body.SpendableAddress, utf8.RuneCountInString(*body.SpendableAddress), 35, false))
		}
	}
	return
}
