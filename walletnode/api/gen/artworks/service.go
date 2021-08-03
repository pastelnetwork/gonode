// Code generated by goa v3.4.3, DO NOT EDIT.
//
// artworks service
//
// Command:
// $ goa gen github.com/pastelnetwork/gonode/walletnode/api/design

package artworks

import (
	"context"

	artworksviews "github.com/pastelnetwork/gonode/walletnode/api/gen/artworks/views"
	goa "goa.design/goa/v3/pkg"
	"goa.design/goa/v3/security"
)

// Pastel Artwork
type Service interface {
	// Runs a new registration process for the new artwork.
	Register(context.Context, *RegisterPayload) (res *RegisterResult, err error)
	// Streams the state of the registration process.
	RegisterTaskState(context.Context, *RegisterTaskStatePayload, RegisterTaskStateServerStream) (err error)
	// Returns a single task.
	RegisterTask(context.Context, *RegisterTaskPayload) (res *Task, err error)
	// List of all tasks.
	RegisterTasks(context.Context) (res TaskCollection, err error)
	// Upload the image that is used when registering a new artwork.
	UploadImage(context.Context, *UploadImagePayload) (res *Image, err error)
	// Streams the search result for artwork
	ArtSearch(context.Context, *ArtSearchPayload, ArtSearchServerStream) (err error)
	// Gets the Artwork detail
	ArtworkGet(context.Context, *ArtworkGetPayload) (res *ArtworkDetail, err error)
	// Download registered artwork.
	Download(context.Context, *DownloadPayload, DownloadServerStream) (err error)
}

// Auther defines the authorization functions to be implemented by the service.
type Auther interface {
	// APIKeyAuth implements the authorization logic for the APIKey security scheme.
	APIKeyAuth(ctx context.Context, key string, schema *security.APIKeyScheme) (context.Context, error)
}

// ServiceName is the name of the service as defined in the design. This is the
// same value that is set in the endpoint request contexts under the ServiceKey
// key.
const ServiceName = "artworks"

// MethodNames lists the service method names as defined in the design. These
// are the same values that are set in the endpoint request contexts under the
// MethodKey key.
var MethodNames = [8]string{"register", "registerTaskState", "registerTask", "registerTasks", "uploadImage", "artSearch", "artworkGet", "download"}

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

// ArtSearchServerStream is the interface a "artSearch" endpoint server stream
// must satisfy.
type ArtSearchServerStream interface {
	// Send streams instances of "ArtworkSearchResult".
	Send(*ArtworkSearchResult) error
	// Close closes the stream.
	Close() error
}

// ArtSearchClientStream is the interface a "artSearch" endpoint client stream
// must satisfy.
type ArtSearchClientStream interface {
	// Recv reads instances of "ArtworkSearchResult" from the stream.
	Recv() (*ArtworkSearchResult, error)
}

// DownloadServerStream is the interface a "download" endpoint server stream
// must satisfy.
type DownloadServerStream interface {
	// Send streams instances of "DownloadResult".
	Send(*DownloadResult) error
	// Close closes the stream.
	Close() error
}

// DownloadClientStream is the interface a "download" endpoint client stream
// must satisfy.
type DownloadClientStream interface {
	// Recv reads instances of "DownloadResult" from the stream.
	Recv() (*DownloadResult, error)
}

// RegisterPayload is the payload type of the artworks service register method.
type RegisterPayload struct {
	// Uploaded image ID
	ImageID string
	// Name of the artwork
	Name string
	// Description of the artwork
	Description *string
	// Keywords
	Keywords *string
	// Series name
	SeriesName *string
	// Number of copies issued
	IssuedCopies int
	// Artwork creation video youtube URL
	YoutubeURL *string
	// Artist's PastelID
	ArtistPastelID string
	// Passphrase of the artist's PastelID
	ArtistPastelIDPassphrase string
	// Name of the artist
	ArtistName string
	// Artist website URL
	ArtistWebsiteURL *string
	// Spendable address
	SpendableAddress string
	// Used to find a suitable masternode with a fee equal or less
	MaximumFee float64
	// Percentage the artist received in future sales. If set to 0% he only get
	// paids for the first sale on each copy of the NFT
	Royalty float64
	// To donate 2% of the sale proceeds on every sale to TeamTrees which plants
	// trees
	Green               bool
	ThumbnailCoordinate *Thumbnailcoordinate
}

// RegisterResult is the result type of the artworks service register method.
type RegisterResult struct {
	// Task ID of the registration process
	TaskID string
}

// RegisterTaskStatePayload is the payload type of the artworks service
// registerTaskState method.
type RegisterTaskStatePayload struct {
	// Task ID of the registration process
	TaskID string
}

// TaskState is the result type of the artworks service registerTaskState
// method.
type TaskState struct {
	// Date of the status creation
	Date string
	// Status of the registration process
	Status string
}

// RegisterTaskPayload is the payload type of the artworks service registerTask
// method.
type RegisterTaskPayload struct {
	// Task ID of the registration process
	TaskID string
}

// Task is the result type of the artworks service registerTask method.
type Task struct {
	// JOb ID of the registration process
	ID string
	// Status of the registration process
	Status string
	// List of states from the very beginning of the process
	States []*TaskState
	// txid
	Txid   *string
	Ticket *ArtworkTicket
}

// TaskCollection is the result type of the artworks service registerTasks
// method.
type TaskCollection []*Task

// UploadImagePayload is the payload type of the artworks service uploadImage
// method.
type UploadImagePayload struct {
	// File to upload
	Bytes []byte
	// For internal use
	Filename *string
}

// Image is the result type of the artworks service uploadImage method.
type Image struct {
	// Uploaded image ID
	ImageID string
	// Image expiration
	ExpiresIn string
}

// ArtSearchPayload is the payload type of the artworks service artSearch
// method.
type ArtSearchPayload struct {
	// Artist PastelID or special value; mine
	Artist *string
	// Number of results to be return
	Limit int
	// Query is search query entered by user
	Query string
	// Name of the artist
	ArtistName bool
	// Title of artwork
	ArtTitle bool
	// Artwork series name
	Series bool
	// Artist written statement
	Descr bool
	// Keyword that Artist assigns to Artwork
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
	MinNsfwScore *int
	// Maximum nsfw score
	MaxNsfwScore *int
	// Minimum rareness score
	MinRarenessScore *int
	// Maximum rareness score
	MaxRarenessScore *int
}

// ArtworkSearchResult is the result type of the artworks service artSearch
// method.
type ArtworkSearchResult struct {
	// Artwork data
	Artwork *ArtworkSummary
	// Sort index of the match based on score.This must be used to sort results on
	// UI.
	MatchIndex int
	// Match result details
	Matches []*FuzzyMatch
}

// ArtworkGetPayload is the payload type of the artworks service artworkGet
// method.
type ArtworkGetPayload struct {
	// txid
	Txid string
}

// ArtworkDetail is the result type of the artworks service artworkGet method.
type ArtworkDetail struct {
	// version
	Version *int
	// Green flag
	IsGreen bool
	// how much artist should get on all future resales
	Royalty float64
	// Storage fee
	StorageFee *int
	// nsfw score
	NsfwScore int
	// rareness score
	RarenessScore int
	// seen score
	SeenScore int
	// Thumbnail image
	Thumbnail []byte
	// txid
	Txid string
	// Name of the artwork
	Title string
	// Description of the artwork
	Description string
	// Keywords
	Keywords *string
	// Series name
	SeriesName *string
	// Number of copies
	Copies int
	// Artwork creation video youtube URL
	YoutubeURL *string
	// Artist's PastelID
	ArtistPastelID string
	// Name of the artist
	ArtistName string
	// Artist website URL
	ArtistWebsiteURL *string
}

// DownloadPayload is the payload type of the artworks service download method.
type DownloadPayload struct {
	// Art Registration Ticket transaction ID
	Txid string
	// Owner's PastelID
	Pid string
	// Passphrase of the owner's PastelID
	Key string
}

// DownloadResult is the result type of the artworks service download method.
type DownloadResult struct {
	// File downloaded
	File []byte
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

// Ticket of the registration artwork
type ArtworkTicket struct {
	// Name of the artwork
	Name string
	// Description of the artwork
	Description *string
	// Keywords
	Keywords *string
	// Series name
	SeriesName *string
	// Number of copies issued
	IssuedCopies int
	// Artwork creation video youtube URL
	YoutubeURL *string
	// Artist's PastelID
	ArtistPastelID string
	// Passphrase of the artist's PastelID
	ArtistPastelIDPassphrase string
	// Name of the artist
	ArtistName string
	// Artist website URL
	ArtistWebsiteURL *string
	// Spendable address
	SpendableAddress string
	// Used to find a suitable masternode with a fee equal or less
	MaximumFee float64
	// Percentage the artist received in future sales. If set to 0% he only get
	// paids for the first sale on each copy of the NFT
	Royalty float64
	// To donate 2% of the sale proceeds on every sale to TeamTrees which plants
	// trees
	Green               bool
	ThumbnailCoordinate *Thumbnailcoordinate
}

// Artwork response
type ArtworkSummary struct {
	// Thumbnail image
	Thumbnail []byte
	// txid
	Txid string
	// Name of the artwork
	Title string
	// Description of the artwork
	Description string
	// Keywords
	Keywords *string
	// Series name
	SeriesName *string
	// Number of copies
	Copies int
	// Artwork creation video youtube URL
	YoutubeURL *string
	// Artist's PastelID
	ArtistPastelID string
	// Name of the artist
	ArtistName string
	// Artist website URL
	ArtistWebsiteURL *string
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

// MakeBadRequest builds a goa.ServiceError from an error.
func MakeBadRequest(err error) *goa.ServiceError {
	return &goa.ServiceError{
		Name:    "BadRequest",
		ID:      goa.NewErrorID(),
		Message: err.Error(),
	}
}

// MakeNotFound builds a goa.ServiceError from an error.
func MakeNotFound(err error) *goa.ServiceError {
	return &goa.ServiceError{
		Name:    "NotFound",
		ID:      goa.NewErrorID(),
		Message: err.Error(),
	}
}

// MakeInternalServerError builds a goa.ServiceError from an error.
func MakeInternalServerError(err error) *goa.ServiceError {
	return &goa.ServiceError{
		Name:    "InternalServerError",
		ID:      goa.NewErrorID(),
		Message: err.Error(),
	}
}

// NewRegisterResult initializes result type RegisterResult from viewed result
// type RegisterResult.
func NewRegisterResult(vres *artworksviews.RegisterResult) *RegisterResult {
	return newRegisterResult(vres.Projected)
}

// NewViewedRegisterResult initializes viewed result type RegisterResult from
// result type RegisterResult using the given view.
func NewViewedRegisterResult(res *RegisterResult, view string) *artworksviews.RegisterResult {
	p := newRegisterResultView(res)
	return &artworksviews.RegisterResult{Projected: p, View: "default"}
}

// NewTask initializes result type Task from viewed result type Task.
func NewTask(vres *artworksviews.Task) *Task {
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
func NewViewedTask(res *Task, view string) *artworksviews.Task {
	var vres *artworksviews.Task
	switch view {
	case "tiny":
		p := newTaskViewTiny(res)
		vres = &artworksviews.Task{Projected: p, View: "tiny"}
	case "default", "":
		p := newTaskView(res)
		vres = &artworksviews.Task{Projected: p, View: "default"}
	}
	return vres
}

// NewTaskCollection initializes result type TaskCollection from viewed result
// type TaskCollection.
func NewTaskCollection(vres artworksviews.TaskCollection) TaskCollection {
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
func NewViewedTaskCollection(res TaskCollection, view string) artworksviews.TaskCollection {
	var vres artworksviews.TaskCollection
	switch view {
	case "tiny":
		p := newTaskCollectionViewTiny(res)
		vres = artworksviews.TaskCollection{Projected: p, View: "tiny"}
	case "default", "":
		p := newTaskCollectionView(res)
		vres = artworksviews.TaskCollection{Projected: p, View: "default"}
	}
	return vres
}

// NewImage initializes result type Image from viewed result type Image.
func NewImage(vres *artworksviews.Image) *Image {
	return newImage(vres.Projected)
}

// NewViewedImage initializes viewed result type Image from result type Image
// using the given view.
func NewViewedImage(res *Image, view string) *artworksviews.Image {
	p := newImageView(res)
	return &artworksviews.Image{Projected: p, View: "default"}
}

// newRegisterResult converts projected type RegisterResult to service type
// RegisterResult.
func newRegisterResult(vres *artworksviews.RegisterResultView) *RegisterResult {
	res := &RegisterResult{}
	if vres.TaskID != nil {
		res.TaskID = *vres.TaskID
	}
	return res
}

// newRegisterResultView projects result type RegisterResult to projected type
// RegisterResultView using the "default" view.
func newRegisterResultView(res *RegisterResult) *artworksviews.RegisterResultView {
	vres := &artworksviews.RegisterResultView{
		TaskID: &res.TaskID,
	}
	return vres
}

// newTaskTiny converts projected type Task to service type Task.
func newTaskTiny(vres *artworksviews.TaskView) *Task {
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
		res.Ticket = transformArtworksviewsArtworkTicketViewToArtworkTicket(vres.Ticket)
	}
	return res
}

// newTask converts projected type Task to service type Task.
func newTask(vres *artworksviews.TaskView) *Task {
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
			res.States[i] = transformArtworksviewsTaskStateViewToTaskState(val)
		}
	}
	if vres.Ticket != nil {
		res.Ticket = transformArtworksviewsArtworkTicketViewToArtworkTicket(vres.Ticket)
	}
	return res
}

// newTaskViewTiny projects result type Task to projected type TaskView using
// the "tiny" view.
func newTaskViewTiny(res *Task) *artworksviews.TaskView {
	vres := &artworksviews.TaskView{
		ID:     &res.ID,
		Status: &res.Status,
		Txid:   res.Txid,
	}
	if res.Ticket != nil {
		vres.Ticket = transformArtworkTicketToArtworksviewsArtworkTicketView(res.Ticket)
	}
	return vres
}

// newTaskView projects result type Task to projected type TaskView using the
// "default" view.
func newTaskView(res *Task) *artworksviews.TaskView {
	vres := &artworksviews.TaskView{
		ID:     &res.ID,
		Status: &res.Status,
		Txid:   res.Txid,
	}
	if res.States != nil {
		vres.States = make([]*artworksviews.TaskStateView, len(res.States))
		for i, val := range res.States {
			vres.States[i] = transformTaskStateToArtworksviewsTaskStateView(val)
		}
	}
	if res.Ticket != nil {
		vres.Ticket = transformArtworkTicketToArtworksviewsArtworkTicketView(res.Ticket)
	}
	return vres
}

// newThumbnailcoordinate converts projected type Thumbnailcoordinate to
// service type Thumbnailcoordinate.
func newThumbnailcoordinate(vres *artworksviews.ThumbnailcoordinateView) *Thumbnailcoordinate {
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
func newThumbnailcoordinateView(res *Thumbnailcoordinate) *artworksviews.ThumbnailcoordinateView {
	vres := &artworksviews.ThumbnailcoordinateView{
		TopLeftX:     &res.TopLeftX,
		TopLeftY:     &res.TopLeftY,
		BottomRightX: &res.BottomRightX,
		BottomRightY: &res.BottomRightY,
	}
	return vres
}

// newTaskCollectionTiny converts projected type TaskCollection to service type
// TaskCollection.
func newTaskCollectionTiny(vres artworksviews.TaskCollectionView) TaskCollection {
	res := make(TaskCollection, len(vres))
	for i, n := range vres {
		res[i] = newTaskTiny(n)
	}
	return res
}

// newTaskCollection converts projected type TaskCollection to service type
// TaskCollection.
func newTaskCollection(vres artworksviews.TaskCollectionView) TaskCollection {
	res := make(TaskCollection, len(vres))
	for i, n := range vres {
		res[i] = newTask(n)
	}
	return res
}

// newTaskCollectionViewTiny projects result type TaskCollection to projected
// type TaskCollectionView using the "tiny" view.
func newTaskCollectionViewTiny(res TaskCollection) artworksviews.TaskCollectionView {
	vres := make(artworksviews.TaskCollectionView, len(res))
	for i, n := range res {
		vres[i] = newTaskViewTiny(n)
	}
	return vres
}

// newTaskCollectionView projects result type TaskCollection to projected type
// TaskCollectionView using the "default" view.
func newTaskCollectionView(res TaskCollection) artworksviews.TaskCollectionView {
	vres := make(artworksviews.TaskCollectionView, len(res))
	for i, n := range res {
		vres[i] = newTaskView(n)
	}
	return vres
}

// newImage converts projected type Image to service type Image.
func newImage(vres *artworksviews.ImageView) *Image {
	res := &Image{}
	if vres.ImageID != nil {
		res.ImageID = *vres.ImageID
	}
	if vres.ExpiresIn != nil {
		res.ExpiresIn = *vres.ExpiresIn
	}
	return res
}

// newImageView projects result type Image to projected type ImageView using
// the "default" view.
func newImageView(res *Image) *artworksviews.ImageView {
	vres := &artworksviews.ImageView{
		ImageID:   &res.ImageID,
		ExpiresIn: &res.ExpiresIn,
	}
	return vres
}

// transformArtworksviewsArtworkTicketViewToArtworkTicket builds a value of
// type *ArtworkTicket from a value of type *artworksviews.ArtworkTicketView.
func transformArtworksviewsArtworkTicketViewToArtworkTicket(v *artworksviews.ArtworkTicketView) *ArtworkTicket {
	if v == nil {
		return nil
	}
	res := &ArtworkTicket{
		Name:                     *v.Name,
		Description:              v.Description,
		Keywords:                 v.Keywords,
		SeriesName:               v.SeriesName,
		IssuedCopies:             *v.IssuedCopies,
		YoutubeURL:               v.YoutubeURL,
		ArtistPastelID:           *v.ArtistPastelID,
		ArtistPastelIDPassphrase: *v.ArtistPastelIDPassphrase,
		ArtistName:               *v.ArtistName,
		ArtistWebsiteURL:         v.ArtistWebsiteURL,
		SpendableAddress:         *v.SpendableAddress,
		MaximumFee:               *v.MaximumFee,
	}
	if v.Royalty != nil {
		res.Royalty = *v.Royalty
	}
	if v.Green != nil {
		res.Green = *v.Green
	}
	if v.Royalty == nil {
		res.Royalty = 0
	}
	if v.Green == nil {
		res.Green = false
	}
	if v.ThumbnailCoordinate != nil {
		res.ThumbnailCoordinate = transformArtworksviewsThumbnailcoordinateViewToThumbnailcoordinate(v.ThumbnailCoordinate)
	}

	return res
}

// transformArtworksviewsThumbnailcoordinateViewToThumbnailcoordinate builds a
// value of type *Thumbnailcoordinate from a value of type
// *artworksviews.ThumbnailcoordinateView.
func transformArtworksviewsThumbnailcoordinateViewToThumbnailcoordinate(v *artworksviews.ThumbnailcoordinateView) *Thumbnailcoordinate {
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

// transformArtworksviewsTaskStateViewToTaskState builds a value of type
// *TaskState from a value of type *artworksviews.TaskStateView.
func transformArtworksviewsTaskStateViewToTaskState(v *artworksviews.TaskStateView) *TaskState {
	if v == nil {
		return nil
	}
	res := &TaskState{
		Date:   *v.Date,
		Status: *v.Status,
	}

	return res
}

// transformArtworkTicketToArtworksviewsArtworkTicketView builds a value of
// type *artworksviews.ArtworkTicketView from a value of type *ArtworkTicket.
func transformArtworkTicketToArtworksviewsArtworkTicketView(v *ArtworkTicket) *artworksviews.ArtworkTicketView {
	res := &artworksviews.ArtworkTicketView{
		Name:                     &v.Name,
		Description:              v.Description,
		Keywords:                 v.Keywords,
		SeriesName:               v.SeriesName,
		IssuedCopies:             &v.IssuedCopies,
		YoutubeURL:               v.YoutubeURL,
		ArtistPastelID:           &v.ArtistPastelID,
		ArtistPastelIDPassphrase: &v.ArtistPastelIDPassphrase,
		ArtistName:               &v.ArtistName,
		ArtistWebsiteURL:         v.ArtistWebsiteURL,
		SpendableAddress:         &v.SpendableAddress,
		MaximumFee:               &v.MaximumFee,
		Royalty:                  &v.Royalty,
		Green:                    &v.Green,
	}
	if v.ThumbnailCoordinate != nil {
		res.ThumbnailCoordinate = transformThumbnailcoordinateToArtworksviewsThumbnailcoordinateView(v.ThumbnailCoordinate)
	}

	return res
}

// transformThumbnailcoordinateToArtworksviewsThumbnailcoordinateView builds a
// value of type *artworksviews.ThumbnailcoordinateView from a value of type
// *Thumbnailcoordinate.
func transformThumbnailcoordinateToArtworksviewsThumbnailcoordinateView(v *Thumbnailcoordinate) *artworksviews.ThumbnailcoordinateView {
	if v == nil {
		return nil
	}
	res := &artworksviews.ThumbnailcoordinateView{
		TopLeftX:     &v.TopLeftX,
		TopLeftY:     &v.TopLeftY,
		BottomRightX: &v.BottomRightX,
		BottomRightY: &v.BottomRightY,
	}

	return res
}

// transformTaskStateToArtworksviewsTaskStateView builds a value of type
// *artworksviews.TaskStateView from a value of type *TaskState.
func transformTaskStateToArtworksviewsTaskStateView(v *TaskState) *artworksviews.TaskStateView {
	if v == nil {
		return nil
	}
	res := &artworksviews.TaskStateView{
		Date:   &v.Date,
		Status: &v.Status,
	}

	return res
}
