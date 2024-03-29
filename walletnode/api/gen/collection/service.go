// Code generated by goa v3.15.0, DO NOT EDIT.
//
// collection service
//
// Command:
// $ goa gen github.com/pastelnetwork/gonode/walletnode/api/design

package collection

import (
	"context"

	collectionviews "github.com/pastelnetwork/gonode/walletnode/api/gen/collection/views"
	goa "goa.design/goa/v3/pkg"
	"goa.design/goa/v3/security"
)

// OpenAPI Collection service
type Service interface {
	// Streams the state of the registration process.
	RegisterCollection(context.Context, *RegisterCollectionPayload) (res *RegisterCollectionResponse, err error)
	// Streams the state of the registration process.
	RegisterTaskState(context.Context, *RegisterTaskStatePayload, RegisterTaskStateServerStream) (err error)
	// Gets the history of the task's states.
	GetTaskHistory(context.Context, *GetTaskHistoryPayload) (res []*TaskHistory, err error)
}

// Auther defines the authorization functions to be implemented by the service.
type Auther interface {
	// APIKeyAuth implements the authorization logic for the APIKey security scheme.
	APIKeyAuth(ctx context.Context, key string, schema *security.APIKeyScheme) (context.Context, error)
}

// APIName is the name of the API as defined in the design.
const APIName = "walletnode"

// APIVersion is the version of the API as defined in the design.
const APIVersion = "1.0"

// ServiceName is the name of the service as defined in the design. This is the
// same value that is set in the endpoint request contexts under the ServiceKey
// key.
const ServiceName = "collection"

// MethodNames lists the service method names as defined in the design. These
// are the same values that are set in the endpoint request contexts under the
// MethodKey key.
var MethodNames = [3]string{"registerCollection", "registerTaskState", "getTaskHistory"}

// RegisterTaskStateServerStream is the interface a "registerTaskState"
// endpoint server stream must satisfy.
type RegisterTaskStateServerStream interface {
	// Send streams instances of "TaskState".
	Send(*TaskState) error
	// Close closes the stream.
	Close() error
}

// RegisterTaskStateClientStream is the interface a "registerTaskState"
// endpoint client stream must satisfy.
type RegisterTaskStateClientStream interface {
	// Recv reads instances of "TaskState" from the stream.
	Recv() (*TaskState, error)
}

type Details struct {
	// details regarding the status
	Message *string
	// important fields regarding status history
	Fields map[string]any
}

// GetTaskHistoryPayload is the payload type of the collection service
// getTaskHistory method.
type GetTaskHistoryPayload struct {
	// Task ID of the registration process
	TaskID string
}

// RegisterCollectionPayload is the payload type of the collection service
// registerCollection method.
type RegisterCollectionPayload struct {
	// name of the collection
	CollectionName string
	// type of items, store by collection
	ItemType string
	// list of authorized contributors
	ListOfPastelidsOfAuthorizedContributors []string
	// max no of entries in the collection
	MaxCollectionEntries int
	// no of days to finalize collection
	NoOfDaysToFinalizeCollection int
	// item copy count in the collection
	CollectionItemCopyCount int
	// royalty fee
	Royalty float64
	// green
	Green bool
	// max open nfsw score sense and nft items can have
	MaxPermittedOpenNsfwScore float64
	// min similarity for 1st entry to have
	MinimumSimilarityScoreToFirstEntryInCollection float64
	// App PastelID
	AppPastelID string
	// Spendable address
	SpendableAddress string
	// Passphrase of the owner's PastelID
	Key *string
}

// RegisterCollectionResponse is the result type of the collection service
// registerCollection method.
type RegisterCollectionResponse struct {
	// Uploaded file ID
	TaskID string
}

// RegisterTaskStatePayload is the payload type of the collection service
// registerTaskState method.
type RegisterTaskStatePayload struct {
	// Task ID of the registration process
	TaskID string
}

type TaskHistory struct {
	// Timestamp of the status creation
	Timestamp *string
	// past status string
	Status string
	// message string (if any)
	Message *string
	// details of the status
	Details *Details
}

// TaskState is the result type of the collection service registerTaskState
// method.
type TaskState struct {
	// Date of the status creation
	Date string
	// Status of the registration process
	Status string
}

// MakeUnAuthorized builds a goa.ServiceError from an error.
func MakeUnAuthorized(err error) *goa.ServiceError {
	return goa.NewServiceError(err, "UnAuthorized", false, false, false)
}

// MakeBadRequest builds a goa.ServiceError from an error.
func MakeBadRequest(err error) *goa.ServiceError {
	return goa.NewServiceError(err, "BadRequest", false, false, false)
}

// MakeNotFound builds a goa.ServiceError from an error.
func MakeNotFound(err error) *goa.ServiceError {
	return goa.NewServiceError(err, "NotFound", false, false, false)
}

// MakeInternalServerError builds a goa.ServiceError from an error.
func MakeInternalServerError(err error) *goa.ServiceError {
	return goa.NewServiceError(err, "InternalServerError", false, false, false)
}

// NewRegisterCollectionResponse initializes result type
// RegisterCollectionResponse from viewed result type
// RegisterCollectionResponse.
func NewRegisterCollectionResponse(vres *collectionviews.RegisterCollectionResponse) *RegisterCollectionResponse {
	return newRegisterCollectionResponse(vres.Projected)
}

// NewViewedRegisterCollectionResponse initializes viewed result type
// RegisterCollectionResponse from result type RegisterCollectionResponse using
// the given view.
func NewViewedRegisterCollectionResponse(res *RegisterCollectionResponse, view string) *collectionviews.RegisterCollectionResponse {
	p := newRegisterCollectionResponseView(res)
	return &collectionviews.RegisterCollectionResponse{Projected: p, View: "default"}
}

// newRegisterCollectionResponse converts projected type
// RegisterCollectionResponse to service type RegisterCollectionResponse.
func newRegisterCollectionResponse(vres *collectionviews.RegisterCollectionResponseView) *RegisterCollectionResponse {
	res := &RegisterCollectionResponse{}
	if vres.TaskID != nil {
		res.TaskID = *vres.TaskID
	}
	return res
}

// newRegisterCollectionResponseView projects result type
// RegisterCollectionResponse to projected type RegisterCollectionResponseView
// using the "default" view.
func newRegisterCollectionResponseView(res *RegisterCollectionResponse) *collectionviews.RegisterCollectionResponseView {
	vres := &collectionviews.RegisterCollectionResponseView{
		TaskID: &res.TaskID,
	}
	return vres
}
