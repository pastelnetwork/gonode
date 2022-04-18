// Code generated by goa v3.6.2, DO NOT EDIT.
//
// nft HTTP server encoders and decoders
//
// Command:
// $ goa gen github.com/pastelnetwork/gonode/walletnode/api/design -o api/

package server

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"

	nft "github.com/pastelnetwork/gonode/walletnode/api/gen/nft"
	nftviews "github.com/pastelnetwork/gonode/walletnode/api/gen/nft/views"
	goahttp "goa.design/goa/v3/http"
	goa "goa.design/goa/v3/pkg"
)

// EncodeRegisterResponse returns an encoder for responses returned by the nft
// register endpoint.
func EncodeRegisterResponse(encoder func(context.Context, http.ResponseWriter) goahttp.Encoder) func(context.Context, http.ResponseWriter, interface{}) error {
	return func(ctx context.Context, w http.ResponseWriter, v interface{}) error {
		res := v.(*nftviews.RegisterResult)
		enc := encoder(ctx, w)
		body := NewRegisterResponseBody(res.Projected)
		w.WriteHeader(http.StatusCreated)
		return enc.Encode(body)
	}
}

// DecodeRegisterRequest returns a decoder for requests sent to the nft
// register endpoint.
func DecodeRegisterRequest(mux goahttp.Muxer, decoder func(*http.Request) goahttp.Decoder) func(*http.Request) (interface{}, error) {
	return func(r *http.Request) (interface{}, error) {
		var (
			body RegisterRequestBody
			err  error
		)
		err = decoder(r).Decode(&body)
		if err != nil {
			if err == io.EOF {
				return nil, goa.MissingPayloadError()
			}
			return nil, goa.DecodePayloadError(err.Error())
		}
		err = ValidateRegisterRequestBody(&body)
		if err != nil {
			return nil, err
		}
		payload := NewRegisterPayload(&body)

		return payload, nil
	}
}

// EncodeRegisterError returns an encoder for errors returned by the register
// nft endpoint.
func EncodeRegisterError(encoder func(context.Context, http.ResponseWriter) goahttp.Encoder, formatter func(err error) goahttp.Statuser) func(context.Context, http.ResponseWriter, error) error {
	encodeError := goahttp.ErrorEncoder(encoder, formatter)
	return func(ctx context.Context, w http.ResponseWriter, v error) error {
		var en ErrorNamer
		if !errors.As(v, &en) {
			return encodeError(ctx, w, v)
		}
		switch en.ErrorName() {
		case "BadRequest":
			var res *goa.ServiceError
			errors.As(v, &res)
			enc := encoder(ctx, w)
			var body interface{}
			if formatter != nil {
				body = formatter(res)
			} else {
				body = NewRegisterBadRequestResponseBody(res)
			}
			w.Header().Set("goa-error", res.ErrorName())
			w.WriteHeader(http.StatusBadRequest)
			return enc.Encode(body)
		case "InternalServerError":
			var res *goa.ServiceError
			errors.As(v, &res)
			enc := encoder(ctx, w)
			var body interface{}
			if formatter != nil {
				body = formatter(res)
			} else {
				body = NewRegisterInternalServerErrorResponseBody(res)
			}
			w.Header().Set("goa-error", res.ErrorName())
			w.WriteHeader(http.StatusInternalServerError)
			return enc.Encode(body)
		default:
			return encodeError(ctx, w, v)
		}
	}
}

// DecodeRegisterTaskStateRequest returns a decoder for requests sent to the
// nft registerTaskState endpoint.
func DecodeRegisterTaskStateRequest(mux goahttp.Muxer, decoder func(*http.Request) goahttp.Decoder) func(*http.Request) (interface{}, error) {
	return func(r *http.Request) (interface{}, error) {
		var (
			taskID string
			err    error

			params = mux.Vars(r)
		)
		taskID = params["taskId"]
		if utf8.RuneCountInString(taskID) < 8 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("taskID", taskID, utf8.RuneCountInString(taskID), 8, true))
		}
		if utf8.RuneCountInString(taskID) > 8 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("taskID", taskID, utf8.RuneCountInString(taskID), 8, false))
		}
		if err != nil {
			return nil, err
		}
		payload := NewRegisterTaskStatePayload(taskID)

		return payload, nil
	}
}

// EncodeRegisterTaskStateError returns an encoder for errors returned by the
// registerTaskState nft endpoint.
func EncodeRegisterTaskStateError(encoder func(context.Context, http.ResponseWriter) goahttp.Encoder, formatter func(err error) goahttp.Statuser) func(context.Context, http.ResponseWriter, error) error {
	encodeError := goahttp.ErrorEncoder(encoder, formatter)
	return func(ctx context.Context, w http.ResponseWriter, v error) error {
		var en ErrorNamer
		if !errors.As(v, &en) {
			return encodeError(ctx, w, v)
		}
		switch en.ErrorName() {
		case "NotFound":
			var res *goa.ServiceError
			errors.As(v, &res)
			enc := encoder(ctx, w)
			var body interface{}
			if formatter != nil {
				body = formatter(res)
			} else {
				body = NewRegisterTaskStateNotFoundResponseBody(res)
			}
			w.Header().Set("goa-error", res.ErrorName())
			w.WriteHeader(http.StatusNotFound)
			return enc.Encode(body)
		case "InternalServerError":
			var res *goa.ServiceError
			errors.As(v, &res)
			enc := encoder(ctx, w)
			var body interface{}
			if formatter != nil {
				body = formatter(res)
			} else {
				body = NewRegisterTaskStateInternalServerErrorResponseBody(res)
			}
			w.Header().Set("goa-error", res.ErrorName())
			w.WriteHeader(http.StatusInternalServerError)
			return enc.Encode(body)
		default:
			return encodeError(ctx, w, v)
		}
	}
}

// EncodeGetTaskHistoryResponse returns an encoder for responses returned by
// the nft getTaskHistory endpoint.
func EncodeGetTaskHistoryResponse(encoder func(context.Context, http.ResponseWriter) goahttp.Encoder) func(context.Context, http.ResponseWriter, interface{}) error {
	return func(ctx context.Context, w http.ResponseWriter, v interface{}) error {
		res, _ := v.(*nft.TaskHistory)
		enc := encoder(ctx, w)
		body := NewGetTaskHistoryResponseBody(res)
		w.WriteHeader(http.StatusOK)
		return enc.Encode(body)
	}
}

// DecodeGetTaskHistoryRequest returns a decoder for requests sent to the nft
// getTaskHistory endpoint.
func DecodeGetTaskHistoryRequest(mux goahttp.Muxer, decoder func(*http.Request) goahttp.Decoder) func(*http.Request) (interface{}, error) {
	return func(r *http.Request) (interface{}, error) {
		var (
			taskID string
			err    error

			params = mux.Vars(r)
		)
		taskID = params["taskId"]
		if utf8.RuneCountInString(taskID) < 8 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("taskID", taskID, utf8.RuneCountInString(taskID), 8, true))
		}
		if utf8.RuneCountInString(taskID) > 8 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("taskID", taskID, utf8.RuneCountInString(taskID), 8, false))
		}
		if err != nil {
			return nil, err
		}
		payload := NewGetTaskHistoryPayload(taskID)

		return payload, nil
	}
}

// EncodeGetTaskHistoryError returns an encoder for errors returned by the
// getTaskHistory nft endpoint.
func EncodeGetTaskHistoryError(encoder func(context.Context, http.ResponseWriter) goahttp.Encoder, formatter func(err error) goahttp.Statuser) func(context.Context, http.ResponseWriter, error) error {
	encodeError := goahttp.ErrorEncoder(encoder, formatter)
	return func(ctx context.Context, w http.ResponseWriter, v error) error {
		var en ErrorNamer
		if !errors.As(v, &en) {
			return encodeError(ctx, w, v)
		}
		switch en.ErrorName() {
		case "NotFound":
			var res *goa.ServiceError
			errors.As(v, &res)
			enc := encoder(ctx, w)
			var body interface{}
			if formatter != nil {
				body = formatter(res)
			} else {
				body = NewGetTaskHistoryNotFoundResponseBody(res)
			}
			w.Header().Set("goa-error", res.ErrorName())
			w.WriteHeader(http.StatusNotFound)
			return enc.Encode(body)
		case "InternalServerError":
			var res *goa.ServiceError
			errors.As(v, &res)
			enc := encoder(ctx, w)
			var body interface{}
			if formatter != nil {
				body = formatter(res)
			} else {
				body = NewGetTaskHistoryInternalServerErrorResponseBody(res)
			}
			w.Header().Set("goa-error", res.ErrorName())
			w.WriteHeader(http.StatusInternalServerError)
			return enc.Encode(body)
		default:
			return encodeError(ctx, w, v)
		}
	}
}

// EncodeRegisterTaskResponse returns an encoder for responses returned by the
// nft registerTask endpoint.
func EncodeRegisterTaskResponse(encoder func(context.Context, http.ResponseWriter) goahttp.Encoder) func(context.Context, http.ResponseWriter, interface{}) error {
	return func(ctx context.Context, w http.ResponseWriter, v interface{}) error {
		res := v.(*nftviews.Task)
		enc := encoder(ctx, w)
		body := NewRegisterTaskResponseBody(res.Projected)
		w.WriteHeader(http.StatusOK)
		return enc.Encode(body)
	}
}

// DecodeRegisterTaskRequest returns a decoder for requests sent to the nft
// registerTask endpoint.
func DecodeRegisterTaskRequest(mux goahttp.Muxer, decoder func(*http.Request) goahttp.Decoder) func(*http.Request) (interface{}, error) {
	return func(r *http.Request) (interface{}, error) {
		var (
			taskID string
			err    error

			params = mux.Vars(r)
		)
		taskID = params["taskId"]
		if utf8.RuneCountInString(taskID) < 8 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("taskID", taskID, utf8.RuneCountInString(taskID), 8, true))
		}
		if utf8.RuneCountInString(taskID) > 8 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("taskID", taskID, utf8.RuneCountInString(taskID), 8, false))
		}
		if err != nil {
			return nil, err
		}
		payload := NewRegisterTaskPayload(taskID)

		return payload, nil
	}
}

// EncodeRegisterTaskError returns an encoder for errors returned by the
// registerTask nft endpoint.
func EncodeRegisterTaskError(encoder func(context.Context, http.ResponseWriter) goahttp.Encoder, formatter func(err error) goahttp.Statuser) func(context.Context, http.ResponseWriter, error) error {
	encodeError := goahttp.ErrorEncoder(encoder, formatter)
	return func(ctx context.Context, w http.ResponseWriter, v error) error {
		var en ErrorNamer
		if !errors.As(v, &en) {
			return encodeError(ctx, w, v)
		}
		switch en.ErrorName() {
		case "NotFound":
			var res *goa.ServiceError
			errors.As(v, &res)
			enc := encoder(ctx, w)
			var body interface{}
			if formatter != nil {
				body = formatter(res)
			} else {
				body = NewRegisterTaskNotFoundResponseBody(res)
			}
			w.Header().Set("goa-error", res.ErrorName())
			w.WriteHeader(http.StatusNotFound)
			return enc.Encode(body)
		case "InternalServerError":
			var res *goa.ServiceError
			errors.As(v, &res)
			enc := encoder(ctx, w)
			var body interface{}
			if formatter != nil {
				body = formatter(res)
			} else {
				body = NewRegisterTaskInternalServerErrorResponseBody(res)
			}
			w.Header().Set("goa-error", res.ErrorName())
			w.WriteHeader(http.StatusInternalServerError)
			return enc.Encode(body)
		default:
			return encodeError(ctx, w, v)
		}
	}
}

// EncodeRegisterTasksResponse returns an encoder for responses returned by the
// nft registerTasks endpoint.
func EncodeRegisterTasksResponse(encoder func(context.Context, http.ResponseWriter) goahttp.Encoder) func(context.Context, http.ResponseWriter, interface{}) error {
	return func(ctx context.Context, w http.ResponseWriter, v interface{}) error {
		res := v.(nftviews.TaskCollection)
		enc := encoder(ctx, w)
		body := NewTaskResponseTinyCollection(res.Projected)
		w.WriteHeader(http.StatusOK)
		return enc.Encode(body)
	}
}

// EncodeRegisterTasksError returns an encoder for errors returned by the
// registerTasks nft endpoint.
func EncodeRegisterTasksError(encoder func(context.Context, http.ResponseWriter) goahttp.Encoder, formatter func(err error) goahttp.Statuser) func(context.Context, http.ResponseWriter, error) error {
	encodeError := goahttp.ErrorEncoder(encoder, formatter)
	return func(ctx context.Context, w http.ResponseWriter, v error) error {
		var en ErrorNamer
		if !errors.As(v, &en) {
			return encodeError(ctx, w, v)
		}
		switch en.ErrorName() {
		case "InternalServerError":
			var res *goa.ServiceError
			errors.As(v, &res)
			enc := encoder(ctx, w)
			var body interface{}
			if formatter != nil {
				body = formatter(res)
			} else {
				body = NewRegisterTasksInternalServerErrorResponseBody(res)
			}
			w.Header().Set("goa-error", res.ErrorName())
			w.WriteHeader(http.StatusInternalServerError)
			return enc.Encode(body)
		default:
			return encodeError(ctx, w, v)
		}
	}
}

// EncodeUploadImageResponse returns an encoder for responses returned by the
// nft uploadImage endpoint.
func EncodeUploadImageResponse(encoder func(context.Context, http.ResponseWriter) goahttp.Encoder) func(context.Context, http.ResponseWriter, interface{}) error {
	return func(ctx context.Context, w http.ResponseWriter, v interface{}) error {
		res := v.(*nftviews.Image)
		enc := encoder(ctx, w)
		body := NewUploadImageResponseBody(res.Projected)
		w.WriteHeader(http.StatusCreated)
		return enc.Encode(body)
	}
}

// DecodeUploadImageRequest returns a decoder for requests sent to the nft
// uploadImage endpoint.
func DecodeUploadImageRequest(mux goahttp.Muxer, decoder func(*http.Request) goahttp.Decoder) func(*http.Request) (interface{}, error) {
	return func(r *http.Request) (interface{}, error) {
		var payload *nft.UploadImagePayload
		if err := decoder(r).Decode(&payload); err != nil {
			return nil, goa.DecodePayloadError(err.Error())
		}

		return payload, nil
	}
}

// NewNftUploadImageDecoder returns a decoder to decode the multipart request
// for the "nft" service "uploadImage" endpoint.
func NewNftUploadImageDecoder(mux goahttp.Muxer, nftUploadImageDecoderFn NftUploadImageDecoderFunc) func(r *http.Request) goahttp.Decoder {
	return func(r *http.Request) goahttp.Decoder {
		return goahttp.EncodingFunc(func(v interface{}) error {
			mr, merr := r.MultipartReader()
			if merr != nil {
				return merr
			}
			p := v.(**nft.UploadImagePayload)
			if err := nftUploadImageDecoderFn(mr, p); err != nil {
				return err
			}
			return nil
		})
	}
}

// EncodeUploadImageError returns an encoder for errors returned by the
// uploadImage nft endpoint.
func EncodeUploadImageError(encoder func(context.Context, http.ResponseWriter) goahttp.Encoder, formatter func(err error) goahttp.Statuser) func(context.Context, http.ResponseWriter, error) error {
	encodeError := goahttp.ErrorEncoder(encoder, formatter)
	return func(ctx context.Context, w http.ResponseWriter, v error) error {
		var en ErrorNamer
		if !errors.As(v, &en) {
			return encodeError(ctx, w, v)
		}
		switch en.ErrorName() {
		case "BadRequest":
			var res *goa.ServiceError
			errors.As(v, &res)
			enc := encoder(ctx, w)
			var body interface{}
			if formatter != nil {
				body = formatter(res)
			} else {
				body = NewUploadImageBadRequestResponseBody(res)
			}
			w.Header().Set("goa-error", res.ErrorName())
			w.WriteHeader(http.StatusBadRequest)
			return enc.Encode(body)
		case "InternalServerError":
			var res *goa.ServiceError
			errors.As(v, &res)
			enc := encoder(ctx, w)
			var body interface{}
			if formatter != nil {
				body = formatter(res)
			} else {
				body = NewUploadImageInternalServerErrorResponseBody(res)
			}
			w.Header().Set("goa-error", res.ErrorName())
			w.WriteHeader(http.StatusInternalServerError)
			return enc.Encode(body)
		default:
			return encodeError(ctx, w, v)
		}
	}
}

// DecodeNftSearchRequest returns a decoder for requests sent to the nft
// nftSearch endpoint.
func DecodeNftSearchRequest(mux goahttp.Muxer, decoder func(*http.Request) goahttp.Decoder) func(*http.Request) (interface{}, error) {
	return func(r *http.Request) (interface{}, error) {
		var (
			artist           *string
			limit            int
			query            string
			creatorName      bool
			artTitle         bool
			series           bool
			descr            bool
			keyword          bool
			minCopies        *int
			maxCopies        *int
			minBlock         int
			maxBlock         *int
			isLikelyDupe     *bool
			minRarenessScore *float64
			maxRarenessScore *float64
			minNsfwScore     *float64
			maxNsfwScore     *float64
			userPastelid     *string
			userPassphrase   *string
			err              error
		)
		artistRaw := r.URL.Query().Get("artist")
		if artistRaw != "" {
			artist = &artistRaw
		}
		if artist != nil {
			if utf8.RuneCountInString(*artist) > 256 {
				err = goa.MergeErrors(err, goa.InvalidLengthError("artist", *artist, utf8.RuneCountInString(*artist), 256, false))
			}
		}
		{
			limitRaw := r.URL.Query().Get("limit")
			if limitRaw == "" {
				limit = 10
			} else {
				v, err2 := strconv.ParseInt(limitRaw, 10, strconv.IntSize)
				if err2 != nil {
					err = goa.MergeErrors(err, goa.InvalidFieldTypeError("limit", limitRaw, "integer"))
				}
				limit = int(v)
			}
		}
		if limit < 10 {
			err = goa.MergeErrors(err, goa.InvalidRangeError("limit", limit, 10, true))
		}
		if limit > 200 {
			err = goa.MergeErrors(err, goa.InvalidRangeError("limit", limit, 200, false))
		}
		query = r.URL.Query().Get("query")
		if query == "" {
			err = goa.MergeErrors(err, goa.MissingFieldError("query", "query string"))
		}
		{
			creatorNameRaw := r.URL.Query().Get("creator_name")
			if creatorNameRaw == "" {
				creatorName = true
			} else {
				v, err2 := strconv.ParseBool(creatorNameRaw)
				if err2 != nil {
					err = goa.MergeErrors(err, goa.InvalidFieldTypeError("creatorName", creatorNameRaw, "boolean"))
				}
				creatorName = v
			}
		}
		{
			artTitleRaw := r.URL.Query().Get("art_title")
			if artTitleRaw == "" {
				artTitle = true
			} else {
				v, err2 := strconv.ParseBool(artTitleRaw)
				if err2 != nil {
					err = goa.MergeErrors(err, goa.InvalidFieldTypeError("artTitle", artTitleRaw, "boolean"))
				}
				artTitle = v
			}
		}
		{
			seriesRaw := r.URL.Query().Get("series")
			if seriesRaw == "" {
				series = true
			} else {
				v, err2 := strconv.ParseBool(seriesRaw)
				if err2 != nil {
					err = goa.MergeErrors(err, goa.InvalidFieldTypeError("series", seriesRaw, "boolean"))
				}
				series = v
			}
		}
		{
			descrRaw := r.URL.Query().Get("descr")
			if descrRaw == "" {
				descr = true
			} else {
				v, err2 := strconv.ParseBool(descrRaw)
				if err2 != nil {
					err = goa.MergeErrors(err, goa.InvalidFieldTypeError("descr", descrRaw, "boolean"))
				}
				descr = v
			}
		}
		{
			keywordRaw := r.URL.Query().Get("keyword")
			if keywordRaw == "" {
				keyword = true
			} else {
				v, err2 := strconv.ParseBool(keywordRaw)
				if err2 != nil {
					err = goa.MergeErrors(err, goa.InvalidFieldTypeError("keyword", keywordRaw, "boolean"))
				}
				keyword = v
			}
		}
		{
			minCopiesRaw := r.URL.Query().Get("min_copies")
			if minCopiesRaw != "" {
				v, err2 := strconv.ParseInt(minCopiesRaw, 10, strconv.IntSize)
				if err2 != nil {
					err = goa.MergeErrors(err, goa.InvalidFieldTypeError("minCopies", minCopiesRaw, "integer"))
				}
				pv := int(v)
				minCopies = &pv
			}
		}
		if minCopies != nil {
			if *minCopies < 1 {
				err = goa.MergeErrors(err, goa.InvalidRangeError("minCopies", *minCopies, 1, true))
			}
		}
		if minCopies != nil {
			if *minCopies > 1000 {
				err = goa.MergeErrors(err, goa.InvalidRangeError("minCopies", *minCopies, 1000, false))
			}
		}
		{
			maxCopiesRaw := r.URL.Query().Get("max_copies")
			if maxCopiesRaw != "" {
				v, err2 := strconv.ParseInt(maxCopiesRaw, 10, strconv.IntSize)
				if err2 != nil {
					err = goa.MergeErrors(err, goa.InvalidFieldTypeError("maxCopies", maxCopiesRaw, "integer"))
				}
				pv := int(v)
				maxCopies = &pv
			}
		}
		if maxCopies != nil {
			if *maxCopies < 1 {
				err = goa.MergeErrors(err, goa.InvalidRangeError("maxCopies", *maxCopies, 1, true))
			}
		}
		if maxCopies != nil {
			if *maxCopies > 1000 {
				err = goa.MergeErrors(err, goa.InvalidRangeError("maxCopies", *maxCopies, 1000, false))
			}
		}
		{
			minBlockRaw := r.URL.Query().Get("min_block")
			if minBlockRaw == "" {
				minBlock = 1
			} else {
				v, err2 := strconv.ParseInt(minBlockRaw, 10, strconv.IntSize)
				if err2 != nil {
					err = goa.MergeErrors(err, goa.InvalidFieldTypeError("minBlock", minBlockRaw, "integer"))
				}
				minBlock = int(v)
			}
		}
		if minBlock < 1 {
			err = goa.MergeErrors(err, goa.InvalidRangeError("minBlock", minBlock, 1, true))
		}
		{
			maxBlockRaw := r.URL.Query().Get("max_block")
			if maxBlockRaw != "" {
				v, err2 := strconv.ParseInt(maxBlockRaw, 10, strconv.IntSize)
				if err2 != nil {
					err = goa.MergeErrors(err, goa.InvalidFieldTypeError("maxBlock", maxBlockRaw, "integer"))
				}
				pv := int(v)
				maxBlock = &pv
			}
		}
		if maxBlock != nil {
			if *maxBlock < 1 {
				err = goa.MergeErrors(err, goa.InvalidRangeError("maxBlock", *maxBlock, 1, true))
			}
		}
		{
			isLikelyDupeRaw := r.URL.Query().Get("is_likely_dupe")
			if isLikelyDupeRaw != "" {
				v, err2 := strconv.ParseBool(isLikelyDupeRaw)
				if err2 != nil {
					err = goa.MergeErrors(err, goa.InvalidFieldTypeError("isLikelyDupe", isLikelyDupeRaw, "boolean"))
				}
				isLikelyDupe = &v
			}
		}
		{
			minRarenessScoreRaw := r.URL.Query().Get("min_rareness_score")
			if minRarenessScoreRaw != "" {
				v, err2 := strconv.ParseFloat(minRarenessScoreRaw, 64)
				if err2 != nil {
					err = goa.MergeErrors(err, goa.InvalidFieldTypeError("minRarenessScore", minRarenessScoreRaw, "float"))
				}
				minRarenessScore = &v
			}
		}
		if minRarenessScore != nil {
			if *minRarenessScore < 0 {
				err = goa.MergeErrors(err, goa.InvalidRangeError("minRarenessScore", *minRarenessScore, 0, true))
			}
		}
		if minRarenessScore != nil {
			if *minRarenessScore > 1 {
				err = goa.MergeErrors(err, goa.InvalidRangeError("minRarenessScore", *minRarenessScore, 1, false))
			}
		}
		{
			maxRarenessScoreRaw := r.URL.Query().Get("max_rareness_score")
			if maxRarenessScoreRaw != "" {
				v, err2 := strconv.ParseFloat(maxRarenessScoreRaw, 64)
				if err2 != nil {
					err = goa.MergeErrors(err, goa.InvalidFieldTypeError("maxRarenessScore", maxRarenessScoreRaw, "float"))
				}
				maxRarenessScore = &v
			}
		}
		if maxRarenessScore != nil {
			if *maxRarenessScore < 0 {
				err = goa.MergeErrors(err, goa.InvalidRangeError("maxRarenessScore", *maxRarenessScore, 0, true))
			}
		}
		if maxRarenessScore != nil {
			if *maxRarenessScore > 1 {
				err = goa.MergeErrors(err, goa.InvalidRangeError("maxRarenessScore", *maxRarenessScore, 1, false))
			}
		}
		{
			minNsfwScoreRaw := r.URL.Query().Get("min_nsfw_score")
			if minNsfwScoreRaw != "" {
				v, err2 := strconv.ParseFloat(minNsfwScoreRaw, 64)
				if err2 != nil {
					err = goa.MergeErrors(err, goa.InvalidFieldTypeError("minNsfwScore", minNsfwScoreRaw, "float"))
				}
				minNsfwScore = &v
			}
		}
		if minNsfwScore != nil {
			if *minNsfwScore < 0 {
				err = goa.MergeErrors(err, goa.InvalidRangeError("minNsfwScore", *minNsfwScore, 0, true))
			}
		}
		if minNsfwScore != nil {
			if *minNsfwScore > 1 {
				err = goa.MergeErrors(err, goa.InvalidRangeError("minNsfwScore", *minNsfwScore, 1, false))
			}
		}
		{
			maxNsfwScoreRaw := r.URL.Query().Get("max_nsfw_score")
			if maxNsfwScoreRaw != "" {
				v, err2 := strconv.ParseFloat(maxNsfwScoreRaw, 64)
				if err2 != nil {
					err = goa.MergeErrors(err, goa.InvalidFieldTypeError("maxNsfwScore", maxNsfwScoreRaw, "float"))
				}
				maxNsfwScore = &v
			}
		}
		if maxNsfwScore != nil {
			if *maxNsfwScore < 0 {
				err = goa.MergeErrors(err, goa.InvalidRangeError("maxNsfwScore", *maxNsfwScore, 0, true))
			}
		}
		if maxNsfwScore != nil {
			if *maxNsfwScore > 1 {
				err = goa.MergeErrors(err, goa.InvalidRangeError("maxNsfwScore", *maxNsfwScore, 1, false))
			}
		}
		userPastelidRaw := r.Header.Get("user_pastelid")
		if userPastelidRaw != "" {
			userPastelid = &userPastelidRaw
		}
		if userPastelid != nil {
			err = goa.MergeErrors(err, goa.ValidatePattern("userPastelid", *userPastelid, "^[a-zA-Z0-9]+$"))
		}
		if userPastelid != nil {
			if utf8.RuneCountInString(*userPastelid) < 86 {
				err = goa.MergeErrors(err, goa.InvalidLengthError("userPastelid", *userPastelid, utf8.RuneCountInString(*userPastelid), 86, true))
			}
		}
		if userPastelid != nil {
			if utf8.RuneCountInString(*userPastelid) > 86 {
				err = goa.MergeErrors(err, goa.InvalidLengthError("userPastelid", *userPastelid, utf8.RuneCountInString(*userPastelid), 86, false))
			}
		}
		userPassphraseRaw := r.Header.Get("user_passphrase")
		if userPassphraseRaw != "" {
			userPassphrase = &userPassphraseRaw
		}
		if err != nil {
			return nil, err
		}
		payload := NewNftSearchPayload(artist, limit, query, creatorName, artTitle, series, descr, keyword, minCopies, maxCopies, minBlock, maxBlock, isLikelyDupe, minRarenessScore, maxRarenessScore, minNsfwScore, maxNsfwScore, userPastelid, userPassphrase)

		return payload, nil
	}
}

// EncodeNftSearchError returns an encoder for errors returned by the nftSearch
// nft endpoint.
func EncodeNftSearchError(encoder func(context.Context, http.ResponseWriter) goahttp.Encoder, formatter func(err error) goahttp.Statuser) func(context.Context, http.ResponseWriter, error) error {
	encodeError := goahttp.ErrorEncoder(encoder, formatter)
	return func(ctx context.Context, w http.ResponseWriter, v error) error {
		var en ErrorNamer
		if !errors.As(v, &en) {
			return encodeError(ctx, w, v)
		}
		switch en.ErrorName() {
		case "BadRequest":
			var res *goa.ServiceError
			errors.As(v, &res)
			enc := encoder(ctx, w)
			var body interface{}
			if formatter != nil {
				body = formatter(res)
			} else {
				body = NewNftSearchBadRequestResponseBody(res)
			}
			w.Header().Set("goa-error", res.ErrorName())
			w.WriteHeader(http.StatusBadRequest)
			return enc.Encode(body)
		case "InternalServerError":
			var res *goa.ServiceError
			errors.As(v, &res)
			enc := encoder(ctx, w)
			var body interface{}
			if formatter != nil {
				body = formatter(res)
			} else {
				body = NewNftSearchInternalServerErrorResponseBody(res)
			}
			w.Header().Set("goa-error", res.ErrorName())
			w.WriteHeader(http.StatusInternalServerError)
			return enc.Encode(body)
		default:
			return encodeError(ctx, w, v)
		}
	}
}

// EncodeNftGetResponse returns an encoder for responses returned by the nft
// nftGet endpoint.
func EncodeNftGetResponse(encoder func(context.Context, http.ResponseWriter) goahttp.Encoder) func(context.Context, http.ResponseWriter, interface{}) error {
	return func(ctx context.Context, w http.ResponseWriter, v interface{}) error {
		res, _ := v.(*nft.NftDetail)
		enc := encoder(ctx, w)
		body := NewNftGetResponseBody(res)
		w.WriteHeader(http.StatusOK)
		return enc.Encode(body)
	}
}

// DecodeNftGetRequest returns a decoder for requests sent to the nft nftGet
// endpoint.
func DecodeNftGetRequest(mux goahttp.Muxer, decoder func(*http.Request) goahttp.Decoder) func(*http.Request) (interface{}, error) {
	return func(r *http.Request) (interface{}, error) {
		var (
			body NftGetRequestBody
			err  error
		)
		err = decoder(r).Decode(&body)
		if err != nil {
			if err == io.EOF {
				return nil, goa.MissingPayloadError()
			}
			return nil, goa.DecodePayloadError(err.Error())
		}
		err = ValidateNftGetRequestBody(&body)
		if err != nil {
			return nil, err
		}

		var (
			txid string

			params = mux.Vars(r)
		)
		txid = params["txid"]
		if utf8.RuneCountInString(txid) < 64 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("txid", txid, utf8.RuneCountInString(txid), 64, true))
		}
		if utf8.RuneCountInString(txid) > 64 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("txid", txid, utf8.RuneCountInString(txid), 64, false))
		}
		if err != nil {
			return nil, err
		}
		payload := NewNftGetPayload(&body, txid)

		return payload, nil
	}
}

// EncodeNftGetError returns an encoder for errors returned by the nftGet nft
// endpoint.
func EncodeNftGetError(encoder func(context.Context, http.ResponseWriter) goahttp.Encoder, formatter func(err error) goahttp.Statuser) func(context.Context, http.ResponseWriter, error) error {
	encodeError := goahttp.ErrorEncoder(encoder, formatter)
	return func(ctx context.Context, w http.ResponseWriter, v error) error {
		var en ErrorNamer
		if !errors.As(v, &en) {
			return encodeError(ctx, w, v)
		}
		switch en.ErrorName() {
		case "BadRequest":
			var res *goa.ServiceError
			errors.As(v, &res)
			enc := encoder(ctx, w)
			var body interface{}
			if formatter != nil {
				body = formatter(res)
			} else {
				body = NewNftGetBadRequestResponseBody(res)
			}
			w.Header().Set("goa-error", res.ErrorName())
			w.WriteHeader(http.StatusBadRequest)
			return enc.Encode(body)
		case "NotFound":
			var res *goa.ServiceError
			errors.As(v, &res)
			enc := encoder(ctx, w)
			var body interface{}
			if formatter != nil {
				body = formatter(res)
			} else {
				body = NewNftGetNotFoundResponseBody(res)
			}
			w.Header().Set("goa-error", res.ErrorName())
			w.WriteHeader(http.StatusNotFound)
			return enc.Encode(body)
		case "InternalServerError":
			var res *goa.ServiceError
			errors.As(v, &res)
			enc := encoder(ctx, w)
			var body interface{}
			if formatter != nil {
				body = formatter(res)
			} else {
				body = NewNftGetInternalServerErrorResponseBody(res)
			}
			w.Header().Set("goa-error", res.ErrorName())
			w.WriteHeader(http.StatusInternalServerError)
			return enc.Encode(body)
		default:
			return encodeError(ctx, w, v)
		}
	}
}

// EncodeDownloadResponse returns an encoder for responses returned by the nft
// download endpoint.
func EncodeDownloadResponse(encoder func(context.Context, http.ResponseWriter) goahttp.Encoder) func(context.Context, http.ResponseWriter, interface{}) error {
	return func(ctx context.Context, w http.ResponseWriter, v interface{}) error {
		res, _ := v.(*nft.DownloadResult)
		enc := encoder(ctx, w)
		body := NewDownloadResponseBody(res)
		w.WriteHeader(http.StatusOK)
		return enc.Encode(body)
	}
}

// DecodeDownloadRequest returns a decoder for requests sent to the nft
// download endpoint.
func DecodeDownloadRequest(mux goahttp.Muxer, decoder func(*http.Request) goahttp.Decoder) func(*http.Request) (interface{}, error) {
	return func(r *http.Request) (interface{}, error) {
		var (
			txid string
			pid  string
			key  string
			err  error
		)
		txid = r.URL.Query().Get("txid")
		if txid == "" {
			err = goa.MergeErrors(err, goa.MissingFieldError("txid", "query string"))
		}
		if utf8.RuneCountInString(txid) < 64 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("txid", txid, utf8.RuneCountInString(txid), 64, true))
		}
		if utf8.RuneCountInString(txid) > 64 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("txid", txid, utf8.RuneCountInString(txid), 64, false))
		}
		pid = r.URL.Query().Get("pid")
		if pid == "" {
			err = goa.MergeErrors(err, goa.MissingFieldError("pid", "query string"))
		}
		err = goa.MergeErrors(err, goa.ValidatePattern("pid", pid, "^[a-zA-Z0-9]+$"))
		if utf8.RuneCountInString(pid) < 86 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("pid", pid, utf8.RuneCountInString(pid), 86, true))
		}
		if utf8.RuneCountInString(pid) > 86 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("pid", pid, utf8.RuneCountInString(pid), 86, false))
		}
		key = r.Header.Get("Authorization")
		if key == "" {
			err = goa.MergeErrors(err, goa.MissingFieldError("Authorization", "header"))
		}
		if err != nil {
			return nil, err
		}
		payload := NewDownloadPayload(txid, pid, key)
		if strings.Contains(payload.Key, " ") {
			// Remove authorization scheme prefix (e.g. "Bearer")
			cred := strings.SplitN(payload.Key, " ", 2)[1]
			payload.Key = cred
		}

		return payload, nil
	}
}

// EncodeDownloadError returns an encoder for errors returned by the download
// nft endpoint.
func EncodeDownloadError(encoder func(context.Context, http.ResponseWriter) goahttp.Encoder, formatter func(err error) goahttp.Statuser) func(context.Context, http.ResponseWriter, error) error {
	encodeError := goahttp.ErrorEncoder(encoder, formatter)
	return func(ctx context.Context, w http.ResponseWriter, v error) error {
		var en ErrorNamer
		if !errors.As(v, &en) {
			return encodeError(ctx, w, v)
		}
		switch en.ErrorName() {
		case "NotFound":
			var res *goa.ServiceError
			errors.As(v, &res)
			enc := encoder(ctx, w)
			var body interface{}
			if formatter != nil {
				body = formatter(res)
			} else {
				body = NewDownloadNotFoundResponseBody(res)
			}
			w.Header().Set("goa-error", res.ErrorName())
			w.WriteHeader(http.StatusNotFound)
			return enc.Encode(body)
		case "InternalServerError":
			var res *goa.ServiceError
			errors.As(v, &res)
			enc := encoder(ctx, w)
			var body interface{}
			if formatter != nil {
				body = formatter(res)
			} else {
				body = NewDownloadInternalServerErrorResponseBody(res)
			}
			w.Header().Set("goa-error", res.ErrorName())
			w.WriteHeader(http.StatusInternalServerError)
			return enc.Encode(body)
		default:
			return encodeError(ctx, w, v)
		}
	}
}

// unmarshalThumbnailcoordinateRequestBodyToNftThumbnailcoordinate builds a
// value of type *nft.Thumbnailcoordinate from a value of type
// *ThumbnailcoordinateRequestBody.
func unmarshalThumbnailcoordinateRequestBodyToNftThumbnailcoordinate(v *ThumbnailcoordinateRequestBody) *nft.Thumbnailcoordinate {
	if v == nil {
		return nil
	}
	res := &nft.Thumbnailcoordinate{
		TopLeftX:     *v.TopLeftX,
		TopLeftY:     *v.TopLeftY,
		BottomRightX: *v.BottomRightX,
		BottomRightY: *v.BottomRightY,
	}

	return res
}

// marshalNftviewsTaskStateViewToTaskStateResponseBody builds a value of type
// *TaskStateResponseBody from a value of type *nftviews.TaskStateView.
func marshalNftviewsTaskStateViewToTaskStateResponseBody(v *nftviews.TaskStateView) *TaskStateResponseBody {
	if v == nil {
		return nil
	}
	res := &TaskStateResponseBody{
		Date:   *v.Date,
		Status: *v.Status,
	}

	return res
}

// marshalNftviewsNftRegisterPayloadViewToNftRegisterPayloadResponseBody builds
// a value of type *NftRegisterPayloadResponseBody from a value of type
// *nftviews.NftRegisterPayloadView.
func marshalNftviewsNftRegisterPayloadViewToNftRegisterPayloadResponseBody(v *nftviews.NftRegisterPayloadView) *NftRegisterPayloadResponseBody {
	res := &NftRegisterPayloadResponseBody{
		Name:                      *v.Name,
		Description:               v.Description,
		Keywords:                  v.Keywords,
		SeriesName:                v.SeriesName,
		IssuedCopies:              *v.IssuedCopies,
		YoutubeURL:                v.YoutubeURL,
		CreatorPastelID:           *v.CreatorPastelID,
		CreatorPastelIDPassphrase: *v.CreatorPastelIDPassphrase,
		CreatorName:               *v.CreatorName,
		CreatorWebsiteURL:         v.CreatorWebsiteURL,
		SpendableAddress:          *v.SpendableAddress,
		MaximumFee:                *v.MaximumFee,
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
		res.ThumbnailCoordinate = marshalNftviewsThumbnailcoordinateViewToThumbnailcoordinateResponseBody(v.ThumbnailCoordinate)
	}

	return res
}

// marshalNftviewsThumbnailcoordinateViewToThumbnailcoordinateResponseBody
// builds a value of type *ThumbnailcoordinateResponseBody from a value of type
// *nftviews.ThumbnailcoordinateView.
func marshalNftviewsThumbnailcoordinateViewToThumbnailcoordinateResponseBody(v *nftviews.ThumbnailcoordinateView) *ThumbnailcoordinateResponseBody {
	if v == nil {
		return nil
	}
	res := &ThumbnailcoordinateResponseBody{
		TopLeftX:     *v.TopLeftX,
		TopLeftY:     *v.TopLeftY,
		BottomRightX: *v.BottomRightX,
		BottomRightY: *v.BottomRightY,
	}

	return res
}

// marshalNftviewsTaskViewToTaskResponseTiny builds a value of type
// *TaskResponseTiny from a value of type *nftviews.TaskView.
func marshalNftviewsTaskViewToTaskResponseTiny(v *nftviews.TaskView) *TaskResponseTiny {
	res := &TaskResponseTiny{
		ID:     *v.ID,
		Status: *v.Status,
		Txid:   v.Txid,
	}
	if v.Ticket != nil {
		res.Ticket = marshalNftviewsNftRegisterPayloadViewToNftRegisterPayloadResponse(v.Ticket)
	}

	return res
}

// marshalNftviewsNftRegisterPayloadViewToNftRegisterPayloadResponse builds a
// value of type *NftRegisterPayloadResponse from a value of type
// *nftviews.NftRegisterPayloadView.
func marshalNftviewsNftRegisterPayloadViewToNftRegisterPayloadResponse(v *nftviews.NftRegisterPayloadView) *NftRegisterPayloadResponse {
	res := &NftRegisterPayloadResponse{
		Name:                      *v.Name,
		Description:               v.Description,
		Keywords:                  v.Keywords,
		SeriesName:                v.SeriesName,
		IssuedCopies:              *v.IssuedCopies,
		YoutubeURL:                v.YoutubeURL,
		CreatorPastelID:           *v.CreatorPastelID,
		CreatorPastelIDPassphrase: *v.CreatorPastelIDPassphrase,
		CreatorName:               *v.CreatorName,
		CreatorWebsiteURL:         v.CreatorWebsiteURL,
		SpendableAddress:          *v.SpendableAddress,
		MaximumFee:                *v.MaximumFee,
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
		res.ThumbnailCoordinate = marshalNftviewsThumbnailcoordinateViewToThumbnailcoordinateResponse(v.ThumbnailCoordinate)
	}

	return res
}

// marshalNftviewsThumbnailcoordinateViewToThumbnailcoordinateResponse builds a
// value of type *ThumbnailcoordinateResponse from a value of type
// *nftviews.ThumbnailcoordinateView.
func marshalNftviewsThumbnailcoordinateViewToThumbnailcoordinateResponse(v *nftviews.ThumbnailcoordinateView) *ThumbnailcoordinateResponse {
	if v == nil {
		return nil
	}
	res := &ThumbnailcoordinateResponse{
		TopLeftX:     *v.TopLeftX,
		TopLeftY:     *v.TopLeftY,
		BottomRightX: *v.BottomRightX,
		BottomRightY: *v.BottomRightY,
	}

	return res
}

// marshalNftNftSummaryToNftSummaryResponseBody builds a value of type
// *NftSummaryResponseBody from a value of type *nft.NftSummary.
func marshalNftNftSummaryToNftSummaryResponseBody(v *nft.NftSummary) *NftSummaryResponseBody {
	res := &NftSummaryResponseBody{
		Thumbnail1:        v.Thumbnail1,
		Thumbnail2:        v.Thumbnail2,
		Txid:              v.Txid,
		Title:             v.Title,
		Description:       v.Description,
		Keywords:          v.Keywords,
		SeriesName:        v.SeriesName,
		Copies:            v.Copies,
		YoutubeURL:        v.YoutubeURL,
		CreatorPastelID:   v.CreatorPastelID,
		CreatorName:       v.CreatorName,
		CreatorWebsiteURL: v.CreatorWebsiteURL,
		NsfwScore:         v.NsfwScore,
		RarenessScore:     v.RarenessScore,
		IsLikelyDupe:      v.IsLikelyDupe,
	}

	return res
}

// marshalNftFuzzyMatchToFuzzyMatchResponseBody builds a value of type
// *FuzzyMatchResponseBody from a value of type *nft.FuzzyMatch.
func marshalNftFuzzyMatchToFuzzyMatchResponseBody(v *nft.FuzzyMatch) *FuzzyMatchResponseBody {
	res := &FuzzyMatchResponseBody{
		Str:       v.Str,
		FieldType: v.FieldType,
		Score:     v.Score,
	}
	if v.MatchedIndexes != nil {
		res.MatchedIndexes = make([]int, len(v.MatchedIndexes))
		for i, val := range v.MatchedIndexes {
			res.MatchedIndexes[i] = val
		}
	}

	return res
}
