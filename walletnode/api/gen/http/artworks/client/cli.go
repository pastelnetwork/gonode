// Code generated by goa v3.4.3, DO NOT EDIT.
//
// artworks HTTP client CLI support package
//
// Command:
// $ goa gen github.com/pastelnetwork/gonode/walletnode/api/design

package client

import (
	"encoding/json"
	"fmt"
	"strconv"
	"unicode/utf8"

	artworks "github.com/pastelnetwork/gonode/walletnode/api/gen/artworks"
	goa "goa.design/goa/v3/pkg"
)

// BuildRegisterPayload builds the payload for the artworks register endpoint
// from CLI flags.
func BuildRegisterPayload(artworksRegisterBody string) (*artworks.RegisterPayload, error) {
	var err error
	var body RegisterRequestBody
	{
		err = json.Unmarshal([]byte(artworksRegisterBody), &body)
		if err != nil {
			return nil, fmt.Errorf("invalid JSON for body, \nerror: %s, \nexample of valid JSON:\n%s", err, "'{\n      \"artist_name\": \"Leonardo da Vinci\",\n      \"artist_pastelid\": \"jXYJud3rmrR1Sk2scvR47N4E4J5Vv48uCC6se2nzHrBRdjaKj3ybPoi1Y2VVoRqi1GnQrYKjSxQAC7NBtvtEdS\",\n      \"artist_pastelid_passphrase\": \"qwerasdf1234\",\n      \"artist_website_url\": \"https://www.leonardodavinci.net\",\n      \"description\": \"The Mona Lisa is an oil painting by Italian artist, inventor, and writer Leonardo da Vinci. Likely completed in 1506, the piece features a portrait of a seated woman set against an imaginary landscape.\",\n      \"green\": false,\n      \"image_id\": \"VK7mpAqZ\",\n      \"issued_copies\": 1,\n      \"keywords\": \"Renaissance, sfumato, portrait\",\n      \"maximum_fee\": 100,\n      \"name\": \"Mona Lisa\",\n      \"royalty\": 12,\n      \"series_name\": \"Famous artist\",\n      \"spendable_address\": \"PtiqRXn2VQwBjp1K8QXR2uW2w2oZ3Ns7N6j\",\n      \"thumbnail_coordinate\": {\n         \"bottom_right_x\": 640,\n         \"bottom_right_y\": 480,\n         \"top_left_x\": 0,\n         \"top_left_y\": 0\n      },\n      \"youtube_url\": \"https://www.youtube.com/watch?v=0xl6Ufo4ZX0\"\n   }'")
		}
		if utf8.RuneCountInString(body.ImageID) < 8 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("body.image_id", body.ImageID, utf8.RuneCountInString(body.ImageID), 8, true))
		}
		if utf8.RuneCountInString(body.ImageID) > 8 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("body.image_id", body.ImageID, utf8.RuneCountInString(body.ImageID), 8, false))
		}
		if utf8.RuneCountInString(body.Name) > 256 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("body.name", body.Name, utf8.RuneCountInString(body.Name), 256, false))
		}
		if body.Description != nil {
			if utf8.RuneCountInString(*body.Description) > 1024 {
				err = goa.MergeErrors(err, goa.InvalidLengthError("body.description", *body.Description, utf8.RuneCountInString(*body.Description), 1024, false))
			}
		}
		if body.Keywords != nil {
			if utf8.RuneCountInString(*body.Keywords) > 256 {
				err = goa.MergeErrors(err, goa.InvalidLengthError("body.keywords", *body.Keywords, utf8.RuneCountInString(*body.Keywords), 256, false))
			}
		}
		if body.SeriesName != nil {
			if utf8.RuneCountInString(*body.SeriesName) > 256 {
				err = goa.MergeErrors(err, goa.InvalidLengthError("body.series_name", *body.SeriesName, utf8.RuneCountInString(*body.SeriesName), 256, false))
			}
		}
		if body.IssuedCopies < 1 {
			err = goa.MergeErrors(err, goa.InvalidRangeError("body.issued_copies", body.IssuedCopies, 1, true))
		}
		if body.IssuedCopies > 1000 {
			err = goa.MergeErrors(err, goa.InvalidRangeError("body.issued_copies", body.IssuedCopies, 1000, false))
		}
		if body.YoutubeURL != nil {
			if utf8.RuneCountInString(*body.YoutubeURL) > 128 {
				err = goa.MergeErrors(err, goa.InvalidLengthError("body.youtube_url", *body.YoutubeURL, utf8.RuneCountInString(*body.YoutubeURL), 128, false))
			}
		}
		err = goa.MergeErrors(err, goa.ValidatePattern("body.artist_pastelid", body.ArtistPastelID, "^[a-zA-Z0-9]+$"))
		if utf8.RuneCountInString(body.ArtistPastelID) < 86 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("body.artist_pastelid", body.ArtistPastelID, utf8.RuneCountInString(body.ArtistPastelID), 86, true))
		}
		if utf8.RuneCountInString(body.ArtistPastelID) > 86 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("body.artist_pastelid", body.ArtistPastelID, utf8.RuneCountInString(body.ArtistPastelID), 86, false))
		}
		if utf8.RuneCountInString(body.ArtistName) > 256 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("body.artist_name", body.ArtistName, utf8.RuneCountInString(body.ArtistName), 256, false))
		}
		if body.ArtistWebsiteURL != nil {
			if utf8.RuneCountInString(*body.ArtistWebsiteURL) > 256 {
				err = goa.MergeErrors(err, goa.InvalidLengthError("body.artist_website_url", *body.ArtistWebsiteURL, utf8.RuneCountInString(*body.ArtistWebsiteURL), 256, false))
			}
		}
		err = goa.MergeErrors(err, goa.ValidatePattern("body.spendable_address", body.SpendableAddress, "^[a-zA-Z0-9]+$"))
		if utf8.RuneCountInString(body.SpendableAddress) < 35 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("body.spendable_address", body.SpendableAddress, utf8.RuneCountInString(body.SpendableAddress), 35, true))
		}
		if utf8.RuneCountInString(body.SpendableAddress) > 35 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("body.spendable_address", body.SpendableAddress, utf8.RuneCountInString(body.SpendableAddress), 35, false))
		}
		if body.MaximumFee < 1e-05 {
			err = goa.MergeErrors(err, goa.InvalidRangeError("body.maximum_fee", body.MaximumFee, 1e-05, true))
		}
		if body.Royalty < 0 {
			err = goa.MergeErrors(err, goa.InvalidRangeError("body.royalty", body.Royalty, 0, true))
		}
		if body.Royalty > 100 {
			err = goa.MergeErrors(err, goa.InvalidRangeError("body.royalty", body.Royalty, 100, false))
		}
		if err != nil {
			return nil, err
		}
	}
	v := &artworks.RegisterPayload{
		ImageID:                  body.ImageID,
		Name:                     body.Name,
		Description:              body.Description,
		Keywords:                 body.Keywords,
		SeriesName:               body.SeriesName,
		IssuedCopies:             body.IssuedCopies,
		YoutubeURL:               body.YoutubeURL,
		ArtistPastelID:           body.ArtistPastelID,
		ArtistPastelIDPassphrase: body.ArtistPastelIDPassphrase,
		ArtistName:               body.ArtistName,
		ArtistWebsiteURL:         body.ArtistWebsiteURL,
		SpendableAddress:         body.SpendableAddress,
		MaximumFee:               body.MaximumFee,
		Royalty:                  body.Royalty,
		Green:                    body.Green,
	}
	{
		var zero float64
		if v.Royalty == zero {
			v.Royalty = 0
		}
	}
	{
		var zero bool
		if v.Green == zero {
			v.Green = false
		}
	}
	if body.ThumbnailCoordinate != nil {
		v.ThumbnailCoordinate = marshalThumbnailcoordinateRequestBodyToArtworksThumbnailcoordinate(body.ThumbnailCoordinate)
	}

	return v, nil
}

// BuildRegisterTaskStatePayload builds the payload for the artworks
// registerTaskState endpoint from CLI flags.
func BuildRegisterTaskStatePayload(artworksRegisterTaskStateTaskID string) (*artworks.RegisterTaskStatePayload, error) {
	var err error
	var taskID string
	{
		taskID = artworksRegisterTaskStateTaskID
		if utf8.RuneCountInString(taskID) < 8 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("taskID", taskID, utf8.RuneCountInString(taskID), 8, true))
		}
		if utf8.RuneCountInString(taskID) > 8 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("taskID", taskID, utf8.RuneCountInString(taskID), 8, false))
		}
		if err != nil {
			return nil, err
		}
	}
	v := &artworks.RegisterTaskStatePayload{}
	v.TaskID = taskID

	return v, nil
}

// BuildRegisterTaskPayload builds the payload for the artworks registerTask
// endpoint from CLI flags.
func BuildRegisterTaskPayload(artworksRegisterTaskTaskID string) (*artworks.RegisterTaskPayload, error) {
	var err error
	var taskID string
	{
		taskID = artworksRegisterTaskTaskID
		if utf8.RuneCountInString(taskID) < 8 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("taskID", taskID, utf8.RuneCountInString(taskID), 8, true))
		}
		if utf8.RuneCountInString(taskID) > 8 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("taskID", taskID, utf8.RuneCountInString(taskID), 8, false))
		}
		if err != nil {
			return nil, err
		}
	}
	v := &artworks.RegisterTaskPayload{}
	v.TaskID = taskID

	return v, nil
}

// BuildUploadImagePayload builds the payload for the artworks uploadImage
// endpoint from CLI flags.
func BuildUploadImagePayload(artworksUploadImageBody string) (*artworks.UploadImagePayload, error) {
	var err error
	var body UploadImageRequestBody
	{
		err = json.Unmarshal([]byte(artworksUploadImageBody), &body)
		if err != nil {
			return nil, fmt.Errorf("invalid JSON for body, \nerror: %s, \nexample of valid JSON:\n%s", err, "'{\n      \"file\": \"QXV0ZW0gZWFxdWUu\"\n   }'")
		}
		if body.Bytes == nil {
			err = goa.MergeErrors(err, goa.MissingFieldError("file", "body"))
		}
		if err != nil {
			return nil, err
		}
	}
	v := &artworks.UploadImagePayload{
		Bytes:    body.Bytes,
		Filename: body.Filename,
	}

	return v, nil
}

// BuildArtSearchPayload builds the payload for the artworks artSearch endpoint
// from CLI flags.
func BuildArtSearchPayload(artworksArtSearchArtist string, artworksArtSearchLimit string, artworksArtSearchQuery string, artworksArtSearchArtistName string, artworksArtSearchArtTitle string, artworksArtSearchSeries string, artworksArtSearchDescr string, artworksArtSearchKeyword string, artworksArtSearchMinCopies string, artworksArtSearchMaxCopies string, artworksArtSearchMinBlock string, artworksArtSearchMaxBlock string, artworksArtSearchMinRarenessScore string, artworksArtSearchMaxRarenessScore string, artworksArtSearchMinNsfwScore string, artworksArtSearchMaxNsfwScore string, artworksArtSearchMinInternetRarenessScore string, artworksArtSearchMaxInternetRarenessScore string) (*artworks.ArtSearchPayload, error) {
	var err error
	var artist *string
	{
		if artworksArtSearchArtist != "" {
			artist = &artworksArtSearchArtist
			if artist != nil {
				if utf8.RuneCountInString(*artist) > 256 {
					err = goa.MergeErrors(err, goa.InvalidLengthError("artist", *artist, utf8.RuneCountInString(*artist), 256, false))
				}
			}
			if err != nil {
				return nil, err
			}
		}
	}
	var limit int
	{
		if artworksArtSearchLimit != "" {
			var v int64
			v, err = strconv.ParseInt(artworksArtSearchLimit, 10, 64)
			limit = int(v)
			if err != nil {
				return nil, fmt.Errorf("invalid value for limit, must be INT")
			}
			if limit < 10 {
				err = goa.MergeErrors(err, goa.InvalidRangeError("limit", limit, 10, true))
			}
			if limit > 200 {
				err = goa.MergeErrors(err, goa.InvalidRangeError("limit", limit, 200, false))
			}
			if err != nil {
				return nil, err
			}
		}
	}
	var query string
	{
		query = artworksArtSearchQuery
	}
	var artistName bool
	{
		if artworksArtSearchArtistName != "" {
			artistName, err = strconv.ParseBool(artworksArtSearchArtistName)
			if err != nil {
				return nil, fmt.Errorf("invalid value for artistName, must be BOOL")
			}
		}
	}
	var artTitle bool
	{
		if artworksArtSearchArtTitle != "" {
			artTitle, err = strconv.ParseBool(artworksArtSearchArtTitle)
			if err != nil {
				return nil, fmt.Errorf("invalid value for artTitle, must be BOOL")
			}
		}
	}
	var series bool
	{
		if artworksArtSearchSeries != "" {
			series, err = strconv.ParseBool(artworksArtSearchSeries)
			if err != nil {
				return nil, fmt.Errorf("invalid value for series, must be BOOL")
			}
		}
	}
	var descr bool
	{
		if artworksArtSearchDescr != "" {
			descr, err = strconv.ParseBool(artworksArtSearchDescr)
			if err != nil {
				return nil, fmt.Errorf("invalid value for descr, must be BOOL")
			}
		}
	}
	var keyword bool
	{
		if artworksArtSearchKeyword != "" {
			keyword, err = strconv.ParseBool(artworksArtSearchKeyword)
			if err != nil {
				return nil, fmt.Errorf("invalid value for keyword, must be BOOL")
			}
		}
	}
	var minCopies *int
	{
		if artworksArtSearchMinCopies != "" {
			var v int64
			v, err = strconv.ParseInt(artworksArtSearchMinCopies, 10, 64)
			val := int(v)
			minCopies = &val
			if err != nil {
				return nil, fmt.Errorf("invalid value for minCopies, must be INT")
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
			if err != nil {
				return nil, err
			}
		}
	}
	var maxCopies *int
	{
		if artworksArtSearchMaxCopies != "" {
			var v int64
			v, err = strconv.ParseInt(artworksArtSearchMaxCopies, 10, 64)
			val := int(v)
			maxCopies = &val
			if err != nil {
				return nil, fmt.Errorf("invalid value for maxCopies, must be INT")
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
			if err != nil {
				return nil, err
			}
		}
	}
	var minBlock int
	{
		if artworksArtSearchMinBlock != "" {
			var v int64
			v, err = strconv.ParseInt(artworksArtSearchMinBlock, 10, 64)
			minBlock = int(v)
			if err != nil {
				return nil, fmt.Errorf("invalid value for minBlock, must be INT")
			}
			if minBlock < 1 {
				err = goa.MergeErrors(err, goa.InvalidRangeError("minBlock", minBlock, 1, true))
			}
			if err != nil {
				return nil, err
			}
		}
	}
	var maxBlock *int
	{
		if artworksArtSearchMaxBlock != "" {
			var v int64
			v, err = strconv.ParseInt(artworksArtSearchMaxBlock, 10, 64)
			val := int(v)
			maxBlock = &val
			if err != nil {
				return nil, fmt.Errorf("invalid value for maxBlock, must be INT")
			}
			if maxBlock != nil {
				if *maxBlock < 1 {
					err = goa.MergeErrors(err, goa.InvalidRangeError("maxBlock", *maxBlock, 1, true))
				}
			}
			if err != nil {
				return nil, err
			}
		}
	}
	var minRarenessScore *float64
	{
		if artworksArtSearchMinRarenessScore != "" {
			val, err := strconv.ParseFloat(artworksArtSearchMinRarenessScore, 64)
			minRarenessScore = &val
			if err != nil {
				return nil, fmt.Errorf("invalid value for minRarenessScore, must be FLOAT64")
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
			if err != nil {
				return nil, err
			}
		}
	}
	var maxRarenessScore *float64
	{
		if artworksArtSearchMaxRarenessScore != "" {
			val, err := strconv.ParseFloat(artworksArtSearchMaxRarenessScore, 64)
			maxRarenessScore = &val
			if err != nil {
				return nil, fmt.Errorf("invalid value for maxRarenessScore, must be FLOAT64")
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
			if err != nil {
				return nil, err
			}
		}
	}
	var minNsfwScore *float64
	{
		if artworksArtSearchMinNsfwScore != "" {
			val, err := strconv.ParseFloat(artworksArtSearchMinNsfwScore, 64)
			minNsfwScore = &val
			if err != nil {
				return nil, fmt.Errorf("invalid value for minNsfwScore, must be FLOAT64")
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
			if err != nil {
				return nil, err
			}
		}
	}
	var maxNsfwScore *float64
	{
		if artworksArtSearchMaxNsfwScore != "" {
			val, err := strconv.ParseFloat(artworksArtSearchMaxNsfwScore, 64)
			maxNsfwScore = &val
			if err != nil {
				return nil, fmt.Errorf("invalid value for maxNsfwScore, must be FLOAT64")
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
			if err != nil {
				return nil, err
			}
		}
	}
	var minInternetRarenessScore *float64
	{
		if artworksArtSearchMinInternetRarenessScore != "" {
			val, err := strconv.ParseFloat(artworksArtSearchMinInternetRarenessScore, 64)
			minInternetRarenessScore = &val
			if err != nil {
				return nil, fmt.Errorf("invalid value for minInternetRarenessScore, must be FLOAT64")
			}
			if minInternetRarenessScore != nil {
				if *minInternetRarenessScore < 0 {
					err = goa.MergeErrors(err, goa.InvalidRangeError("minInternetRarenessScore", *minInternetRarenessScore, 0, true))
				}
			}
			if minInternetRarenessScore != nil {
				if *minInternetRarenessScore > 1 {
					err = goa.MergeErrors(err, goa.InvalidRangeError("minInternetRarenessScore", *minInternetRarenessScore, 1, false))
				}
			}
			if err != nil {
				return nil, err
			}
		}
	}
	var maxInternetRarenessScore *float64
	{
		if artworksArtSearchMaxInternetRarenessScore != "" {
			val, err := strconv.ParseFloat(artworksArtSearchMaxInternetRarenessScore, 64)
			maxInternetRarenessScore = &val
			if err != nil {
				return nil, fmt.Errorf("invalid value for maxInternetRarenessScore, must be FLOAT64")
			}
			if maxInternetRarenessScore != nil {
				if *maxInternetRarenessScore < 0 {
					err = goa.MergeErrors(err, goa.InvalidRangeError("maxInternetRarenessScore", *maxInternetRarenessScore, 0, true))
				}
			}
			if maxInternetRarenessScore != nil {
				if *maxInternetRarenessScore > 1 {
					err = goa.MergeErrors(err, goa.InvalidRangeError("maxInternetRarenessScore", *maxInternetRarenessScore, 1, false))
				}
			}
			if err != nil {
				return nil, err
			}
		}
	}
	v := &artworks.ArtSearchPayload{}
	v.Artist = artist
	v.Limit = limit
	v.Query = query
	v.ArtistName = artistName
	v.ArtTitle = artTitle
	v.Series = series
	v.Descr = descr
	v.Keyword = keyword
	v.MinCopies = minCopies
	v.MaxCopies = maxCopies
	v.MinBlock = minBlock
	v.MaxBlock = maxBlock
	v.MinRarenessScore = minRarenessScore
	v.MaxRarenessScore = maxRarenessScore
	v.MinNsfwScore = minNsfwScore
	v.MaxNsfwScore = maxNsfwScore
	v.MinInternetRarenessScore = minInternetRarenessScore
	v.MaxInternetRarenessScore = maxInternetRarenessScore

	return v, nil
}

// BuildArtworkGetPayload builds the payload for the artworks artworkGet
// endpoint from CLI flags.
func BuildArtworkGetPayload(artworksArtworkGetTxid string) (*artworks.ArtworkGetPayload, error) {
	var err error
	var txid string
	{
		txid = artworksArtworkGetTxid
		if utf8.RuneCountInString(txid) < 64 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("txid", txid, utf8.RuneCountInString(txid), 64, true))
		}
		if utf8.RuneCountInString(txid) > 64 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("txid", txid, utf8.RuneCountInString(txid), 64, false))
		}
		if err != nil {
			return nil, err
		}
	}
	v := &artworks.ArtworkGetPayload{}
	v.Txid = txid

	return v, nil
}

// BuildDownloadPayload builds the payload for the artworks download endpoint
// from CLI flags.
func BuildDownloadPayload(artworksDownloadTxid string, artworksDownloadPid string, artworksDownloadKey string) (*artworks.ArtworkDownloadPayload, error) {
	var err error
	var txid string
	{
		txid = artworksDownloadTxid
		if utf8.RuneCountInString(txid) < 64 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("txid", txid, utf8.RuneCountInString(txid), 64, true))
		}
		if utf8.RuneCountInString(txid) > 64 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("txid", txid, utf8.RuneCountInString(txid), 64, false))
		}
		if err != nil {
			return nil, err
		}
	}
	var pid string
	{
		pid = artworksDownloadPid
		err = goa.MergeErrors(err, goa.ValidatePattern("pid", pid, "^[a-zA-Z0-9]+$"))
		if utf8.RuneCountInString(pid) < 86 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("pid", pid, utf8.RuneCountInString(pid), 86, true))
		}
		if utf8.RuneCountInString(pid) > 86 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("pid", pid, utf8.RuneCountInString(pid), 86, false))
		}
		if err != nil {
			return nil, err
		}
	}
	var key string
	{
		key = artworksDownloadKey
	}
	v := &artworks.ArtworkDownloadPayload{}
	v.Txid = txid
	v.Pid = pid
	v.Key = key

	return v, nil
}
