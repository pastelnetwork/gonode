// Code generated by goa v3.14.0, DO NOT EDIT.
//
// nft service
//
// Command:
// $ goa gen github.com/pastelnetwork/gonode/walletnode/api/design

package nft

import (
	"context"

	nftviews "github.com/pastelnetwork/gonode/walletnode/api/gen/nft/views"
	goa "goa.design/goa/v3/pkg"
	"goa.design/goa/v3/security"
)

// Pastel NFT
type Service interface {
	// Runs a new registration process for the new NFT.
	Register(context.Context, *RegisterPayload) (res *RegisterResult, err error)
	// Streams the state of the registration process.
	RegisterTaskState(context.Context, *RegisterTaskStatePayload, RegisterTaskStateServerStream) (err error)
	// Gets the history of the task's states.
	GetTaskHistory(context.Context, *GetTaskHistoryPayload) (res []*TaskHistory, err error)
	// Returns a single task.
	RegisterTask(context.Context, *RegisterTaskPayload) (res *Task, err error)
	// List of all tasks.
	RegisterTasks(context.Context) (res TaskCollection, err error)
	// Upload the image that is used when registering a new NFT.
	UploadImage(context.Context, *UploadImagePayload) (res *ImageRes, err error)
	// Streams the search result for NFT
	NftSearch(context.Context, *NftSearchPayload, NftSearchServerStream) (err error)
	// Gets the NFT detail
	NftGet(context.Context, *NftGetPayload) (res *NftDetail, err error)
	// Download registered NFT.
	Download(context.Context, *DownloadPayload) (res *FileDownloadResult, err error)
	// Duplication detection output file details
	DdServiceOutputFileDetail(context.Context, *DownloadPayload) (res *DDServiceOutputFileResult, err error)
	// Duplication detection output file
	DdServiceOutputFile(context.Context, *DownloadPayload) (res *DDFPResultFile, err error)
}

// Auther defines the authorization functions to be implemented by the service.
type Auther interface {
	// APIKeyAuth implements the authorization logic for the APIKey security scheme.
	APIKeyAuth(ctx context.Context, key string, schema *security.APIKeyScheme) (context.Context, error)
}

// ServiceName is the name of the service as defined in the design. This is the
// same value that is set in the endpoint request contexts under the ServiceKey
// key.
const ServiceName = "nft"

// MethodNames lists the service method names as defined in the design. These
// are the same values that are set in the endpoint request contexts under the
// MethodKey key.
var MethodNames = [11]string{"register", "registerTaskState", "getTaskHistory", "registerTask", "registerTasks", "uploadImage", "nftSearch", "nftGet", "download", "ddServiceOutputFileDetail", "ddServiceOutputFile"}

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

// NftSearchServerStream is the interface a "nftSearch" endpoint server stream
// must satisfy.
type NftSearchServerStream interface {
	// Send streams instances of "NftSearchResult".
	Send(*NftSearchResult) error
	// Close closes the stream.
	Close() error
}

// NftSearchClientStream is the interface a "nftSearch" endpoint client stream
// must satisfy.
type NftSearchClientStream interface {
	// Recv reads instances of "NftSearchResult" from the stream.
	Recv() (*NftSearchResult, error)
}

type AlternativeNSFWScores struct {
	// drawings nsfw score
	Drawings *float32
	// hentai nsfw score
	Hentai *float32
	// sexy nsfw score
	Sexy *float32
	// porn nsfw score
	Porn *float32
	// neutral nsfw score
	Neutral *float32
}

// DDFPResultFile is the result type of the nft service ddServiceOutputFile
// method.
type DDFPResultFile struct {
	// File downloaded
	File string
}

// DDServiceOutputFileResult is the result type of the nft service
// ddServiceOutputFileDetail method.
type DDServiceOutputFileResult struct {
	// block hash when request submitted
	PastelBlockHashWhenRequestSubmitted *string
	// block Height when request submitted
	PastelBlockHeightWhenRequestSubmitted *string
	// timestamp of request when submitted
	UtcTimestampWhenRequestSubmitted *string
	// pastel id of the submitter
	PastelIDOfSubmitter *string
	// pastel id of registering SN1
	PastelIDOfRegisteringSupernode1 *string
	// pastel id of registering SN2
	PastelIDOfRegisteringSupernode2 *string
	// pastel id of registering SN3
	PastelIDOfRegisteringSupernode3 *string
	// is pastel open API request
	IsPastelOpenapiRequest *bool
	// system version of dupe detection
	DupeDetectionSystemVersion *string
	// is this nft likely a duplicate
	IsLikelyDupe *bool
	// is this nft rare on the internet
	IsRareOnInternet *bool
	// pastel rareness score
	OverallRarenessScore *float32
	// PCT of top 10 most similar with dupe probe above 25 PCT
	PctOfTop10MostSimilarWithDupeProbAbove25pct *float32
	// PCT of top 10 most similar with dupe probe above 33 PCT
	PctOfTop10MostSimilarWithDupeProbAbove33pct *float32
	// PCT of top 10 most similar with dupe probe above 50 PCT
	PctOfTop10MostSimilarWithDupeProbAbove50pct *float32
	// rareness scores table json compressed b64
	RarenessScoresTableJSONCompressedB64 *string
	// open nsfw score
	OpenNsfwScore *float32
	// Image fingerprint of candidate image file
	ImageFingerprintOfCandidateImageFile []float64
	// hash of candidate image file
	HashOfCandidateImageFile *string
	// name of the collection
	CollectionNameString *string
	// open api group id string
	OpenAPIGroupIDString *string
	// rareness score of the group
	GroupRarenessScore *float32
	// candidate image thumbnail as base64 string
	CandidateImageThumbnailWebpAsBase64String *string
	// does not impact collection strings
	DoesNotImpactTheFollowingCollectionStrings *string
	// similarity score to first entry in collection
	SimilarityScoreToFirstEntryInCollection *float32
	// probability of CP
	CpProbability *float32
	// child probability
	ChildProbability *float32
	// file path of the image
	ImageFilePath *string
	// internet rareness
	InternetRareness *InternetRareness
	// alternative NSFW scores
	AlternativeNsfwScores *AlternativeNSFWScores
	// name of the creator
	CreatorName string
	// website of creator
	CreatorWebsite string
	// written statement of creator
	CreatorWrittenStatement string
	// title of NFT
	NftTitle string
	// series name of NFT
	NftSeriesName string
	// nft creation video youtube url
	NftCreationVideoYoutubeURL string
	// keywords for NFT
	NftKeywordSet string
	// total copies of NFT
	TotalCopies int
	// preview hash of NFT
	PreviewHash []byte
	// thumbnail1 hash of NFT
	Thumbnail1Hash []byte
	// thumbnail2 hash of NFT
	Thumbnail2Hash []byte
	// original file size in bytes
	OriginalFileSizeInBytes int
	// type of the file
	FileType string
	// max permitted open NSFW score
	MaxPermittedOpenNsfwScore float64
}

type Details struct {
	// details regarding the status
	Message *string
	// important fields regarding status history
	Fields map[string]any
}

// DownloadPayload is the payload type of the nft service download method.
type DownloadPayload struct {
	// Nft Registration Request transaction ID
	Txid string
	// Owner's PastelID
	Pid string
	// Passphrase of the owner's PastelID
	Key string
}

// FileDownloadResult is the result type of the nft service download method.
type FileDownloadResult struct {
	// File path
	FileID string
}

type FuzzyMatch struct {
	// String that is matched
	Str *string
	// Field that is matched
	FieldType *string
	// The indexes of matched characters. Useful for highlighting matches
	MatchedIndexes []int
	// Score used to rank matches
	Score *int
}

// GetTaskHistoryPayload is the payload type of the nft service getTaskHistory
// method.
type GetTaskHistoryPayload struct {
	// Task ID of the registration process
	TaskID string
}

// ImageRes is the result type of the nft service uploadImage method.
type ImageRes struct {
	// Uploaded image ID
	ImageID string
	// Image expiration
	ExpiresIn string
	// Estimated fee
	EstimatedFee float64
}

type InternetRareness struct {
	// Base64 Compressed JSON Table of Rare On Internet Summary
	RareOnInternetSummaryTableAsJSONCompressedB64 *string
	// Base64 Compressed JSON of Rare On Internet Graph
	RareOnInternetGraphJSONCompressedB64 *string
	// Base64 Compressed Json of Alternative Rare On Internet Dict
	AlternativeRareOnInternetDictAsJSONCompressedB64 *string
	// Minimum Number of Exact Matches on Page
	MinNumberOfExactMatchesInPage *uint32
	// Earliest Available Date of Internet Results
	EarliestAvailableDateOfInternetResults *string
}

// NftDetail is the result type of the nft service nftGet method.
type NftDetail struct {
	// version
	Version *int
	// Green address
	GreenAddress *bool
	// how much artist should get on all future resales
	Royalty *float64
	// Storage fee %
	StorageFee *int
	// NSFW Average score
	NsfwScore float32
	// Average pastel rareness score
	RarenessScore float32
	// Is this image likely a duplicate of another known image
	IsLikelyDupe bool
	// is this nft rare on the internet
	IsRareOnInternet bool
	// nsfw score
	DrawingNsfwScore *float32
	// nsfw score
	NeutralNsfwScore *float32
	// nsfw score
	SexyNsfwScore *float32
	// nsfw score
	PornNsfwScore *float32
	// nsfw score
	HentaiNsfwScore *float32
	// Preview Image
	PreviewThumbnail []byte
	// Base64 Compressed JSON Table of Rare On Internet Summary
	RareOnInternetSummaryTableJSONB64 *string
	// Base64 Compressed JSON of Rare On Internet Graph
	RareOnInternetGraphJSONB64 *string
	// Base64 Compressed Json of Alternative Rare On Internet Dict
	AltRareOnInternetDictJSONB64 *string
	// Minimum Number of Exact Matches on Page
	MinNumExactMatchesOnPage *uint32
	// Earliest Available Date of Internet Results
	EarliestDateOfResults *string
	// Thumbnail_1 image
	Thumbnail1 []byte
	// Thumbnail_2 image
	Thumbnail2 []byte
	// txid
	Txid string
	// Name of the NFT
	Title string
	// Description of the NFT
	Description string
	// Keywords
	Keywords *string
	// Series name
	SeriesName *string
	// Number of copies
	Copies int
	// NFT creation video youtube URL
	YoutubeURL *string
	// Artist's PastelID
	CreatorPastelID string
	// Name of the artist
	CreatorName string
	// Artist website URL
	CreatorWebsiteURL *string
}

// NftGetPayload is the payload type of the nft service nftGet method.
type NftGetPayload struct {
	// Nft Registration Request transaction ID
	Txid string
	// Owner's PastelID
	Pid string
	// Passphrase of the owner's PastelID
	Key string
}

// Request of the registration NFT
type NftRegisterPayload struct {
	// Name of the NFT
	Name string
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
	CreatorPastelID string
	// Name of the NFT creator
	CreatorName string
	// NFT creator website URL
	CreatorWebsiteURL *string
	// Spendable address
	SpendableAddress string
	// Used to find a suitable masternode with a fee equal or less
	MaximumFee float64
	// Percentage the artist received in future sales. If set to 0% he only get
	// paids for the first sale on each copy of the NFT
	Royalty *float64
	// To donate 2% of the sale proceeds on every sale to TeamTrees which plants
	// trees
	Green               *bool
	ThumbnailCoordinate *Thumbnailcoordinate
	// To make it publicly accessible
	MakePubliclyAccessible bool
	// Act Collection TxID to add given ticket in collection
	CollectionActTxid *string
	// OpenAPI GroupID string
	OpenAPIGroupID string
	// Passphrase of the owner's PastelID
	Key string
}

// NftSearchPayload is the payload type of the nft service nftSearch method.
type NftSearchPayload struct {
	// Artist PastelID or special value; mine
	Artist *string
	// Number of results to be return
	Limit int
	// Query is search query entered by user
	Query string
	// Name of the nft creator
	CreatorName bool
	// Title of NFT
	ArtTitle bool
	// NFT series name
	Series bool
	// Artist written statement
	Descr bool
	// Keyword that Artist assigns to NFT
	Keyword bool
	// Minimum blocknum
	MinBlock int
	// Maximum blocknum
	MaxBlock *int
	// Minimum number of created copies
	MinCopies *int
	// Maximum number of created copies
	MaxCopies *int
	// Minimum nsfw score
	MinNsfwScore *float64
	// Maximum nsfw score
	MaxNsfwScore *float64
	// Minimum pastel rareness score
	MinRarenessScore *float64
	// Maximum pastel rareness score
	MaxRarenessScore *float64
	// Is this image likely a duplicate of another known image
	IsLikelyDupe *bool
	// User's PastelID
	UserPastelid *string
	// Passphrase of the User PastelID
	UserPassphrase *string
}

// NftSearchResult is the result type of the nft service nftSearch method.
type NftSearchResult struct {
	// NFT data
	Nft *NftSummary
	// Sort index of the match based on score.This must be used to sort results on
	// UI.
	MatchIndex int
	// Match result details
	Matches []*FuzzyMatch
}

// NFT response
type NftSummary struct {
	// Thumbnail_1 image
	Thumbnail1 []byte
	// Thumbnail_2 image
	Thumbnail2 []byte
	// txid
	Txid string
	// Name of the NFT
	Title string
	// Description of the NFT
	Description string
	// Keywords
	Keywords *string
	// Series name
	SeriesName *string
	// Number of copies
	Copies int
	// NFT creation video youtube URL
	YoutubeURL *string
	// Artist's PastelID
	CreatorPastelID string
	// Name of the artist
	CreatorName string
	// Artist website URL
	CreatorWebsiteURL *string
	// NSFW Average score
	NsfwScore *float32
	// Average pastel rareness score
	RarenessScore *float32
	// Is this image likely a duplicate of another known image
	IsLikelyDupe *bool
}

// RegisterPayload is the payload type of the nft service register method.
type RegisterPayload struct {
	// Uploaded image ID
	ImageID string
	// Name of the NFT
	Name string
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
	CreatorPastelID string
	// Name of the NFT creator
	CreatorName string
	// NFT creator website URL
	CreatorWebsiteURL *string
	// Spendable address
	SpendableAddress string
	// Used to find a suitable masternode with a fee equal or less
	MaximumFee float64
	// Percentage the artist received in future sales. If set to 0% he only get
	// paids for the first sale on each copy of the NFT
	Royalty *float64
	// To donate 2% of the sale proceeds on every sale to TeamTrees which plants
	// trees
	Green               *bool
	ThumbnailCoordinate *Thumbnailcoordinate
	// To make it publicly accessible
	MakePubliclyAccessible bool
	// Act Collection TxID to add given ticket in collection
	CollectionActTxid *string
	// OpenAPI GroupID string
	OpenAPIGroupID string
	// Passphrase of the owner's PastelID
	Key string
}

// RegisterResult is the result type of the nft service register method.
type RegisterResult struct {
	// Task ID of the registration process
	TaskID string
}

// RegisterTaskPayload is the payload type of the nft service registerTask
// method.
type RegisterTaskPayload struct {
	// Task ID of the registration process
	TaskID string
}

// RegisterTaskStatePayload is the payload type of the nft service
// registerTaskState method.
type RegisterTaskStatePayload struct {
	// Task ID of the registration process
	TaskID string
}

// Task is the result type of the nft service registerTask method.
type Task struct {
	// JOb ID of the registration process
	ID string
	// Status of the registration process
	Status string
	// List of states from the very beginning of the process
	States []*TaskState
	// txid
	Txid   *string
	Ticket *NftRegisterPayload
}

// TaskCollection is the result type of the nft service registerTasks method.
type TaskCollection []*Task

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

// TaskState is the result type of the nft service registerTaskState method.
type TaskState struct {
	// Date of the status creation
	Date string
	// Status of the registration process
	Status string
}

// Coordinate of the thumbnail
type Thumbnailcoordinate struct {
	// X coordinate of the thumbnail's top left conner
	TopLeftX int64
	// Y coordinate of the thumbnail's top left conner
	TopLeftY int64
	// X coordinate of the thumbnail's bottom right conner
	BottomRightX int64
	// Y coordinate of the thumbnail's bottom right conner
	BottomRightY int64
}

// UploadImagePayload is the payload type of the nft service uploadImage method.
type UploadImagePayload struct {
	// File to upload
	Bytes []byte
	// For internal use
	Filename *string
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

// MakeUnAuthorized builds a goa.ServiceError from an error.
func MakeUnAuthorized(err error) *goa.ServiceError {
	return goa.NewServiceError(err, "UnAuthorized", false, false, false)
}

// NewRegisterResult initializes result type RegisterResult from viewed result
// type RegisterResult.
func NewRegisterResult(vres *nftviews.RegisterResult) *RegisterResult {
	return newRegisterResult(vres.Projected)
}

// NewViewedRegisterResult initializes viewed result type RegisterResult from
// result type RegisterResult using the given view.
func NewViewedRegisterResult(res *RegisterResult, view string) *nftviews.RegisterResult {
	p := newRegisterResultView(res)
	return &nftviews.RegisterResult{Projected: p, View: "default"}
}

// NewTask initializes result type Task from viewed result type Task.
func NewTask(vres *nftviews.Task) *Task {
	var res *Task
	switch vres.View {
	case "tiny":
		res = newTaskTiny(vres.Projected)
	case "default", "":
		res = newTask(vres.Projected)
	}
	return res
}

// NewViewedTask initializes viewed result type Task from result type Task
// using the given view.
func NewViewedTask(res *Task, view string) *nftviews.Task {
	var vres *nftviews.Task
	switch view {
	case "tiny":
		p := newTaskViewTiny(res)
		vres = &nftviews.Task{Projected: p, View: "tiny"}
	case "default", "":
		p := newTaskView(res)
		vres = &nftviews.Task{Projected: p, View: "default"}
	}
	return vres
}

// NewTaskCollection initializes result type TaskCollection from viewed result
// type TaskCollection.
func NewTaskCollection(vres nftviews.TaskCollection) TaskCollection {
	var res TaskCollection
	switch vres.View {
	case "tiny":
		res = newTaskCollectionTiny(vres.Projected)
	case "default", "":
		res = newTaskCollection(vres.Projected)
	}
	return res
}

// NewViewedTaskCollection initializes viewed result type TaskCollection from
// result type TaskCollection using the given view.
func NewViewedTaskCollection(res TaskCollection, view string) nftviews.TaskCollection {
	var vres nftviews.TaskCollection
	switch view {
	case "tiny":
		p := newTaskCollectionViewTiny(res)
		vres = nftviews.TaskCollection{Projected: p, View: "tiny"}
	case "default", "":
		p := newTaskCollectionView(res)
		vres = nftviews.TaskCollection{Projected: p, View: "default"}
	}
	return vres
}

// NewImageRes initializes result type ImageRes from viewed result type
// ImageRes.
func NewImageRes(vres *nftviews.ImageRes) *ImageRes {
	return newImageRes(vres.Projected)
}

// NewViewedImageRes initializes viewed result type ImageRes from result type
// ImageRes using the given view.
func NewViewedImageRes(res *ImageRes, view string) *nftviews.ImageRes {
	p := newImageResView(res)
	return &nftviews.ImageRes{Projected: p, View: "default"}
}

// newRegisterResult converts projected type RegisterResult to service type
// RegisterResult.
func newRegisterResult(vres *nftviews.RegisterResultView) *RegisterResult {
	res := &RegisterResult{}
	if vres.TaskID != nil {
		res.TaskID = *vres.TaskID
	}
	return res
}

// newRegisterResultView projects result type RegisterResult to projected type
// RegisterResultView using the "default" view.
func newRegisterResultView(res *RegisterResult) *nftviews.RegisterResultView {
	vres := &nftviews.RegisterResultView{
		TaskID: &res.TaskID,
	}
	return vres
}

// newTaskTiny converts projected type Task to service type Task.
func newTaskTiny(vres *nftviews.TaskView) *Task {
	res := &Task{
		Txid: vres.Txid,
	}
	if vres.ID != nil {
		res.ID = *vres.ID
	}
	if vres.Status != nil {
		res.Status = *vres.Status
	}
	if vres.Ticket != nil {
		res.Ticket = transformNftviewsNftRegisterPayloadViewToNftRegisterPayload(vres.Ticket)
	}
	return res
}

// newTask converts projected type Task to service type Task.
func newTask(vres *nftviews.TaskView) *Task {
	res := &Task{
		Txid: vres.Txid,
	}
	if vres.ID != nil {
		res.ID = *vres.ID
	}
	if vres.Status != nil {
		res.Status = *vres.Status
	}
	if vres.States != nil {
		res.States = make([]*TaskState, len(vres.States))
		for i, val := range vres.States {
			res.States[i] = transformNftviewsTaskStateViewToTaskState(val)
		}
	}
	if vres.Ticket != nil {
		res.Ticket = transformNftviewsNftRegisterPayloadViewToNftRegisterPayload(vres.Ticket)
	}
	return res
}

// newTaskViewTiny projects result type Task to projected type TaskView using
// the "tiny" view.
func newTaskViewTiny(res *Task) *nftviews.TaskView {
	vres := &nftviews.TaskView{
		ID:     &res.ID,
		Status: &res.Status,
		Txid:   res.Txid,
	}
	if res.Ticket != nil {
		vres.Ticket = transformNftRegisterPayloadToNftviewsNftRegisterPayloadView(res.Ticket)
	}
	return vres
}

// newTaskView projects result type Task to projected type TaskView using the
// "default" view.
func newTaskView(res *Task) *nftviews.TaskView {
	vres := &nftviews.TaskView{
		ID:     &res.ID,
		Status: &res.Status,
		Txid:   res.Txid,
	}
	if res.States != nil {
		vres.States = make([]*nftviews.TaskStateView, len(res.States))
		for i, val := range res.States {
			vres.States[i] = transformTaskStateToNftviewsTaskStateView(val)
		}
	}
	if res.Ticket != nil {
		vres.Ticket = transformNftRegisterPayloadToNftviewsNftRegisterPayloadView(res.Ticket)
	}
	return vres
}

// newThumbnailcoordinate converts projected type Thumbnailcoordinate to
// service type Thumbnailcoordinate.
func newThumbnailcoordinate(vres *nftviews.ThumbnailcoordinateView) *Thumbnailcoordinate {
	res := &Thumbnailcoordinate{}
	if vres.TopLeftX != nil {
		res.TopLeftX = *vres.TopLeftX
	}
	if vres.TopLeftY != nil {
		res.TopLeftY = *vres.TopLeftY
	}
	if vres.BottomRightX != nil {
		res.BottomRightX = *vres.BottomRightX
	}
	if vres.BottomRightY != nil {
		res.BottomRightY = *vres.BottomRightY
	}
	if vres.TopLeftX == nil {
		res.TopLeftX = 0
	}
	if vres.TopLeftY == nil {
		res.TopLeftY = 0
	}
	if vres.BottomRightX == nil {
		res.BottomRightX = 0
	}
	if vres.BottomRightY == nil {
		res.BottomRightY = 0
	}
	return res
}

// newThumbnailcoordinateView projects result type Thumbnailcoordinate to
// projected type ThumbnailcoordinateView using the "default" view.
func newThumbnailcoordinateView(res *Thumbnailcoordinate) *nftviews.ThumbnailcoordinateView {
	vres := &nftviews.ThumbnailcoordinateView{
		TopLeftX:     &res.TopLeftX,
		TopLeftY:     &res.TopLeftY,
		BottomRightX: &res.BottomRightX,
		BottomRightY: &res.BottomRightY,
	}
	return vres
}

// newTaskCollectionTiny converts projected type TaskCollection to service type
// TaskCollection.
func newTaskCollectionTiny(vres nftviews.TaskCollectionView) TaskCollection {
	res := make(TaskCollection, len(vres))
	for i, n := range vres {
		res[i] = newTaskTiny(n)
	}
	return res
}

// newTaskCollection converts projected type TaskCollection to service type
// TaskCollection.
func newTaskCollection(vres nftviews.TaskCollectionView) TaskCollection {
	res := make(TaskCollection, len(vres))
	for i, n := range vres {
		res[i] = newTask(n)
	}
	return res
}

// newTaskCollectionViewTiny projects result type TaskCollection to projected
// type TaskCollectionView using the "tiny" view.
func newTaskCollectionViewTiny(res TaskCollection) nftviews.TaskCollectionView {
	vres := make(nftviews.TaskCollectionView, len(res))
	for i, n := range res {
		vres[i] = newTaskViewTiny(n)
	}
	return vres
}

// newTaskCollectionView projects result type TaskCollection to projected type
// TaskCollectionView using the "default" view.
func newTaskCollectionView(res TaskCollection) nftviews.TaskCollectionView {
	vres := make(nftviews.TaskCollectionView, len(res))
	for i, n := range res {
		vres[i] = newTaskView(n)
	}
	return vres
}

// newImageRes converts projected type ImageRes to service type ImageRes.
func newImageRes(vres *nftviews.ImageResView) *ImageRes {
	res := &ImageRes{}
	if vres.ImageID != nil {
		res.ImageID = *vres.ImageID
	}
	if vres.ExpiresIn != nil {
		res.ExpiresIn = *vres.ExpiresIn
	}
	if vres.EstimatedFee != nil {
		res.EstimatedFee = *vres.EstimatedFee
	}
	if vres.EstimatedFee == nil {
		res.EstimatedFee = 1
	}
	return res
}

// newImageResView projects result type ImageRes to projected type ImageResView
// using the "default" view.
func newImageResView(res *ImageRes) *nftviews.ImageResView {
	vres := &nftviews.ImageResView{
		ImageID:      &res.ImageID,
		ExpiresIn:    &res.ExpiresIn,
		EstimatedFee: &res.EstimatedFee,
	}
	return vres
}

// transformNftviewsNftRegisterPayloadViewToNftRegisterPayload builds a value
// of type *NftRegisterPayload from a value of type
// *nftviews.NftRegisterPayloadView.
func transformNftviewsNftRegisterPayloadViewToNftRegisterPayload(v *nftviews.NftRegisterPayloadView) *NftRegisterPayload {
	if v == nil {
		return nil
	}
	res := &NftRegisterPayload{
		Name:              *v.Name,
		Description:       v.Description,
		Keywords:          v.Keywords,
		SeriesName:        v.SeriesName,
		IssuedCopies:      v.IssuedCopies,
		YoutubeURL:        v.YoutubeURL,
		CreatorPastelID:   *v.CreatorPastelID,
		CreatorName:       *v.CreatorName,
		CreatorWebsiteURL: v.CreatorWebsiteURL,
		SpendableAddress:  *v.SpendableAddress,
		MaximumFee:        *v.MaximumFee,
		Royalty:           v.Royalty,
		Green:             v.Green,
		CollectionActTxid: v.CollectionActTxid,
		Key:               *v.Key,
	}
	if v.MakePubliclyAccessible != nil {
		res.MakePubliclyAccessible = *v.MakePubliclyAccessible
	}
	if v.OpenAPIGroupID != nil {
		res.OpenAPIGroupID = *v.OpenAPIGroupID
	}
	if v.ThumbnailCoordinate != nil {
		res.ThumbnailCoordinate = transformNftviewsThumbnailcoordinateViewToThumbnailcoordinate(v.ThumbnailCoordinate)
	}
	if v.MakePubliclyAccessible == nil {
		res.MakePubliclyAccessible = false
	}
	if v.OpenAPIGroupID == nil {
		res.OpenAPIGroupID = "PASTEL"
	}

	return res
}

// transformNftviewsThumbnailcoordinateViewToThumbnailcoordinate builds a value
// of type *Thumbnailcoordinate from a value of type
// *nftviews.ThumbnailcoordinateView.
func transformNftviewsThumbnailcoordinateViewToThumbnailcoordinate(v *nftviews.ThumbnailcoordinateView) *Thumbnailcoordinate {
	if v == nil {
		return nil
	}
	res := &Thumbnailcoordinate{
		TopLeftX:     *v.TopLeftX,
		TopLeftY:     *v.TopLeftY,
		BottomRightX: *v.BottomRightX,
		BottomRightY: *v.BottomRightY,
	}

	return res
}

// transformNftviewsTaskStateViewToTaskState builds a value of type *TaskState
// from a value of type *nftviews.TaskStateView.
func transformNftviewsTaskStateViewToTaskState(v *nftviews.TaskStateView) *TaskState {
	if v == nil {
		return nil
	}
	res := &TaskState{
		Date:   *v.Date,
		Status: *v.Status,
	}

	return res
}

// transformNftRegisterPayloadToNftviewsNftRegisterPayloadView builds a value
// of type *nftviews.NftRegisterPayloadView from a value of type
// *NftRegisterPayload.
func transformNftRegisterPayloadToNftviewsNftRegisterPayloadView(v *NftRegisterPayload) *nftviews.NftRegisterPayloadView {
	res := &nftviews.NftRegisterPayloadView{
		Name:                   &v.Name,
		Description:            v.Description,
		Keywords:               v.Keywords,
		SeriesName:             v.SeriesName,
		IssuedCopies:           v.IssuedCopies,
		YoutubeURL:             v.YoutubeURL,
		CreatorPastelID:        &v.CreatorPastelID,
		CreatorName:            &v.CreatorName,
		CreatorWebsiteURL:      v.CreatorWebsiteURL,
		SpendableAddress:       &v.SpendableAddress,
		MaximumFee:             &v.MaximumFee,
		Royalty:                v.Royalty,
		Green:                  v.Green,
		MakePubliclyAccessible: &v.MakePubliclyAccessible,
		CollectionActTxid:      v.CollectionActTxid,
		OpenAPIGroupID:         &v.OpenAPIGroupID,
		Key:                    &v.Key,
	}
	if v.ThumbnailCoordinate != nil {
		res.ThumbnailCoordinate = transformThumbnailcoordinateToNftviewsThumbnailcoordinateView(v.ThumbnailCoordinate)
	}

	return res
}

// transformThumbnailcoordinateToNftviewsThumbnailcoordinateView builds a value
// of type *nftviews.ThumbnailcoordinateView from a value of type
// *Thumbnailcoordinate.
func transformThumbnailcoordinateToNftviewsThumbnailcoordinateView(v *Thumbnailcoordinate) *nftviews.ThumbnailcoordinateView {
	if v == nil {
		return nil
	}
	res := &nftviews.ThumbnailcoordinateView{
		TopLeftX:     &v.TopLeftX,
		TopLeftY:     &v.TopLeftY,
		BottomRightX: &v.BottomRightX,
		BottomRightY: &v.BottomRightY,
	}

	return res
}

// transformTaskStateToNftviewsTaskStateView builds a value of type
// *nftviews.TaskStateView from a value of type *TaskState.
func transformTaskStateToNftviewsTaskStateView(v *TaskState) *nftviews.TaskStateView {
	if v == nil {
		return nil
	}
	res := &nftviews.TaskStateView{
		Date:   &v.Date,
		Status: &v.Status,
	}

	return res
}
