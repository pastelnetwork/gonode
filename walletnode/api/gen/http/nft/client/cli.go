// Code generated by goa v3.14.0, DO NOT EDIT.
//
// nft HTTP client CLI support package
//
// Command:
// $ goa gen github.com/pastelnetwork/gonode/walletnode/api/design

package client

import (
	"encoding/json"
	"fmt"
	"strconv"
	"unicode/utf8"

	nft "github.com/pastelnetwork/gonode/walletnode/api/gen/nft"
	goa "goa.design/goa/v3/pkg"
)

// BuildRegisterPayload builds the payload for the nft register endpoint from
// CLI flags.
func BuildRegisterPayload(nftRegisterBody string, nftRegisterKey string) (*nft.RegisterPayload, error) {
	var err error
	var body RegisterRequestBody
	{
		err = json.Unmarshal([]byte(nftRegisterBody), &body)
		if err != nil {
			return nil, fmt.Errorf("invalid JSON for body, \nerror: %s, \nexample of valid JSON:\n%s", err, "'{\n      \"burn_txid\": \"576e7b824634a488a2f0baacf5a53b237d883029f205df25b300b87c8877ab58\",\n      \"collection_act_txid\": \"576e7b824634a488a2f0baacf5a53b237d883029f205df25b300b87c8877ab58\",\n      \"creator_name\": \"Leonardo da Vinci\",\n      \"creator_pastelid\": \"jXYJud3rmrR1Sk2scvR47N4E4J5Vv48uCC6se2nzHrBRdjaKj3ybPoi1Y2VVoRqi1GnQrYKjSxQAC7NBtvtEdS\",\n      \"creator_website_url\": \"https://www.leonardodavinci.net\",\n      \"description\": \"The Mona Lisa is an oil painting by Italian artist, inventor, and writer Leonardo da Vinci. Likely completed in 1506, the piece features a portrait of a seated woman set against an imaginary landscape.\",\n      \"green\": false,\n      \"image_id\": \"VK7mpAqZ\",\n      \"issued_copies\": 1,\n      \"keywords\": \"Renaissance, sfumato, portrait\",\n      \"make_publicly_accessible\": false,\n      \"maximum_fee\": 100,\n      \"name\": \"Mona Lisa\",\n      \"open_api_group_id\": \"Cum voluptas ut mollitia incidunt aperiam.\",\n      \"royalty\": 12,\n      \"series_name\": \"Famous artist\",\n      \"spendable_address\": \"PtiqRXn2VQwBjp1K8QXR2uW2w2oZ3Ns7N6j\",\n      \"thumbnail_coordinate\": {\n         \"bottom_right_x\": 640,\n         \"bottom_right_y\": 480,\n         \"top_left_x\": 0,\n         \"top_left_y\": 0\n      },\n      \"youtube_url\": \"https://www.youtube.com/watch?v=0xl6Ufo4ZX0\"\n   }'")
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
		if body.IssuedCopies != nil {
			if *body.IssuedCopies > 1000 {
				err = goa.MergeErrors(err, goa.InvalidRangeError("body.issued_copies", *body.IssuedCopies, 1000, false))
			}
		}
		if body.YoutubeURL != nil {
			if utf8.RuneCountInString(*body.YoutubeURL) > 128 {
				err = goa.MergeErrors(err, goa.InvalidLengthError("body.youtube_url", *body.YoutubeURL, utf8.RuneCountInString(*body.YoutubeURL), 128, false))
			}
		}
		err = goa.MergeErrors(err, goa.ValidatePattern("body.creator_pastelid", body.CreatorPastelID, "^[a-zA-Z0-9]+$"))
		if utf8.RuneCountInString(body.CreatorPastelID) < 86 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("body.creator_pastelid", body.CreatorPastelID, utf8.RuneCountInString(body.CreatorPastelID), 86, true))
		}
		if utf8.RuneCountInString(body.CreatorPastelID) > 86 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("body.creator_pastelid", body.CreatorPastelID, utf8.RuneCountInString(body.CreatorPastelID), 86, false))
		}
		if utf8.RuneCountInString(body.CreatorName) > 256 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("body.creator_name", body.CreatorName, utf8.RuneCountInString(body.CreatorName), 256, false))
		}
		if body.CreatorWebsiteURL != nil {
			if utf8.RuneCountInString(*body.CreatorWebsiteURL) > 256 {
				err = goa.MergeErrors(err, goa.InvalidLengthError("body.creator_website_url", *body.CreatorWebsiteURL, utf8.RuneCountInString(*body.CreatorWebsiteURL), 256, false))
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
		if body.Royalty != nil {
			if *body.Royalty > 20 {
				err = goa.MergeErrors(err, goa.InvalidRangeError("body.royalty", *body.Royalty, 20, false))
			}
		}
		if body.BurnTxid != nil {
			if utf8.RuneCountInString(*body.BurnTxid) < 64 {
				err = goa.MergeErrors(err, goa.InvalidLengthError("body.burn_txid", *body.BurnTxid, utf8.RuneCountInString(*body.BurnTxid), 64, true))
			}
		}
		if body.BurnTxid != nil {
			if utf8.RuneCountInString(*body.BurnTxid) > 64 {
				err = goa.MergeErrors(err, goa.InvalidLengthError("body.burn_txid", *body.BurnTxid, utf8.RuneCountInString(*body.BurnTxid), 64, false))
			}
		}
		if err != nil {
			return nil, err
		}
	}
	var key string
	{
		key = nftRegisterKey
	}
	v := &nft.RegisterPayload{
		ImageID:                body.ImageID,
		Name:                   body.Name,
		Description:            body.Description,
		Keywords:               body.Keywords,
		SeriesName:             body.SeriesName,
		IssuedCopies:           body.IssuedCopies,
		YoutubeURL:             body.YoutubeURL,
		CreatorPastelID:        body.CreatorPastelID,
		CreatorName:            body.CreatorName,
		CreatorWebsiteURL:      body.CreatorWebsiteURL,
		SpendableAddress:       body.SpendableAddress,
		MaximumFee:             body.MaximumFee,
		Royalty:                body.Royalty,
		Green:                  body.Green,
		MakePubliclyAccessible: body.MakePubliclyAccessible,
		CollectionActTxid:      body.CollectionActTxid,
		OpenAPIGroupID:         body.OpenAPIGroupID,
		BurnTxid:               body.BurnTxid,
	}
	if body.ThumbnailCoordinate != nil {
		v.ThumbnailCoordinate = marshalThumbnailcoordinateRequestBodyToNftThumbnailcoordinate(body.ThumbnailCoordinate)
	}
	{
		var zero bool
		if v.MakePubliclyAccessible == zero {
			v.MakePubliclyAccessible = false
		}
	}
	{
		var zero string
		if v.OpenAPIGroupID == zero {
			v.OpenAPIGroupID = "PASTEL"
		}
	}
	v.Key = key

	return v, nil
}

// BuildRegisterTaskStatePayload builds the payload for the nft
// registerTaskState endpoint from CLI flags.
func BuildRegisterTaskStatePayload(nftRegisterTaskStateTaskID string) (*nft.RegisterTaskStatePayload, error) {
	var err error
	var taskID string
	{
		taskID = nftRegisterTaskStateTaskID
		if utf8.RuneCountInString(taskID) < 8 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("taskId", taskID, utf8.RuneCountInString(taskID), 8, true))
		}
		if utf8.RuneCountInString(taskID) > 8 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("taskId", taskID, utf8.RuneCountInString(taskID), 8, false))
		}
		if err != nil {
			return nil, err
		}
	}
	v := &nft.RegisterTaskStatePayload{}
	v.TaskID = taskID

	return v, nil
}

// BuildGetTaskHistoryPayload builds the payload for the nft getTaskHistory
// endpoint from CLI flags.
func BuildGetTaskHistoryPayload(nftGetTaskHistoryTaskID string) (*nft.GetTaskHistoryPayload, error) {
	var err error
	var taskID string
	{
		taskID = nftGetTaskHistoryTaskID
		if utf8.RuneCountInString(taskID) < 8 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("taskId", taskID, utf8.RuneCountInString(taskID), 8, true))
		}
		if utf8.RuneCountInString(taskID) > 8 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("taskId", taskID, utf8.RuneCountInString(taskID), 8, false))
		}
		if err != nil {
			return nil, err
		}
	}
	v := &nft.GetTaskHistoryPayload{}
	v.TaskID = taskID

	return v, nil
}

// BuildRegisterTaskPayload builds the payload for the nft registerTask
// endpoint from CLI flags.
func BuildRegisterTaskPayload(nftRegisterTaskTaskID string) (*nft.RegisterTaskPayload, error) {
	var err error
	var taskID string
	{
		taskID = nftRegisterTaskTaskID
		if utf8.RuneCountInString(taskID) < 8 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("taskId", taskID, utf8.RuneCountInString(taskID), 8, true))
		}
		if utf8.RuneCountInString(taskID) > 8 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("taskId", taskID, utf8.RuneCountInString(taskID), 8, false))
		}
		if err != nil {
			return nil, err
		}
	}
	v := &nft.RegisterTaskPayload{}
	v.TaskID = taskID

	return v, nil
}

// BuildUploadImagePayload builds the payload for the nft uploadImage endpoint
// from CLI flags.
func BuildUploadImagePayload(nftUploadImageBody string) (*nft.UploadImagePayload, error) {
	var err error
	var body UploadImageRequestBody
	{
		err = json.Unmarshal([]byte(nftUploadImageBody), &body)
		if err != nil {
			return nil, fmt.Errorf("invalid JSON for body, \nerror: %s, \nexample of valid JSON:\n%s", err, "'{\n      \"file\": \"SWxsbyB2ZXJpdGF0aXMgYXJjaGl0ZWN0byBzaXQgcXVpLg==\"\n   }'")
		}
		if body.Bytes == nil {
			err = goa.MergeErrors(err, goa.MissingFieldError("file", "body"))
		}
		if err != nil {
			return nil, err
		}
	}
	v := &nft.UploadImagePayload{
		Bytes:    body.Bytes,
		Filename: body.Filename,
	}

	return v, nil
}

// BuildNftSearchPayload builds the payload for the nft nftSearch endpoint from
// CLI flags.
func BuildNftSearchPayload(nftNftSearchArtist string, nftNftSearchLimit string, nftNftSearchQuery string, nftNftSearchCreatorName string, nftNftSearchArtTitle string, nftNftSearchSeries string, nftNftSearchDescr string, nftNftSearchKeyword string, nftNftSearchMinCopies string, nftNftSearchMaxCopies string, nftNftSearchMinBlock string, nftNftSearchMaxBlock string, nftNftSearchIsLikelyDupe string, nftNftSearchMinRarenessScore string, nftNftSearchMaxRarenessScore string, nftNftSearchMinNsfwScore string, nftNftSearchMaxNsfwScore string, nftNftSearchUserPastelid string, nftNftSearchUserPassphrase string) (*nft.NftSearchPayload, error) {
	var err error
	var artist *string
	{
		if nftNftSearchArtist != "" {
			artist = &nftNftSearchArtist
			if utf8.RuneCountInString(*artist) > 256 {
				err = goa.MergeErrors(err, goa.InvalidLengthError("artist", *artist, utf8.RuneCountInString(*artist), 256, false))
			}
			if err != nil {
				return nil, err
			}
		}
	}
	var limit int
	{
		if nftNftSearchLimit != "" {
			var v int64
			v, err = strconv.ParseInt(nftNftSearchLimit, 10, strconv.IntSize)
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
		query = nftNftSearchQuery
	}
	var creatorName bool
	{
		if nftNftSearchCreatorName != "" {
			creatorName, err = strconv.ParseBool(nftNftSearchCreatorName)
			if err != nil {
				return nil, fmt.Errorf("invalid value for creatorName, must be BOOL")
			}
		}
	}
	var artTitle bool
	{
		if nftNftSearchArtTitle != "" {
			artTitle, err = strconv.ParseBool(nftNftSearchArtTitle)
			if err != nil {
				return nil, fmt.Errorf("invalid value for artTitle, must be BOOL")
			}
		}
	}
	var series bool
	{
		if nftNftSearchSeries != "" {
			series, err = strconv.ParseBool(nftNftSearchSeries)
			if err != nil {
				return nil, fmt.Errorf("invalid value for series, must be BOOL")
			}
		}
	}
	var descr bool
	{
		if nftNftSearchDescr != "" {
			descr, err = strconv.ParseBool(nftNftSearchDescr)
			if err != nil {
				return nil, fmt.Errorf("invalid value for descr, must be BOOL")
			}
		}
	}
	var keyword bool
	{
		if nftNftSearchKeyword != "" {
			keyword, err = strconv.ParseBool(nftNftSearchKeyword)
			if err != nil {
				return nil, fmt.Errorf("invalid value for keyword, must be BOOL")
			}
		}
	}
	var minCopies *int
	{
		if nftNftSearchMinCopies != "" {
			var v int64
			v, err = strconv.ParseInt(nftNftSearchMinCopies, 10, strconv.IntSize)
			val := int(v)
			minCopies = &val
			if err != nil {
				return nil, fmt.Errorf("invalid value for minCopies, must be INT")
			}
			if *minCopies < 1 {
				err = goa.MergeErrors(err, goa.InvalidRangeError("min_copies", *minCopies, 1, true))
			}
			if *minCopies > 1000 {
				err = goa.MergeErrors(err, goa.InvalidRangeError("min_copies", *minCopies, 1000, false))
			}
			if err != nil {
				return nil, err
			}
		}
	}
	var maxCopies *int
	{
		if nftNftSearchMaxCopies != "" {
			var v int64
			v, err = strconv.ParseInt(nftNftSearchMaxCopies, 10, strconv.IntSize)
			val := int(v)
			maxCopies = &val
			if err != nil {
				return nil, fmt.Errorf("invalid value for maxCopies, must be INT")
			}
			if *maxCopies < 1 {
				err = goa.MergeErrors(err, goa.InvalidRangeError("max_copies", *maxCopies, 1, true))
			}
			if *maxCopies > 1000 {
				err = goa.MergeErrors(err, goa.InvalidRangeError("max_copies", *maxCopies, 1000, false))
			}
			if err != nil {
				return nil, err
			}
		}
	}
	var minBlock int
	{
		if nftNftSearchMinBlock != "" {
			var v int64
			v, err = strconv.ParseInt(nftNftSearchMinBlock, 10, strconv.IntSize)
			minBlock = int(v)
			if err != nil {
				return nil, fmt.Errorf("invalid value for minBlock, must be INT")
			}
			if minBlock < 1 {
				err = goa.MergeErrors(err, goa.InvalidRangeError("min_block", minBlock, 1, true))
			}
			if err != nil {
				return nil, err
			}
		}
	}
	var maxBlock *int
	{
		if nftNftSearchMaxBlock != "" {
			var v int64
			v, err = strconv.ParseInt(nftNftSearchMaxBlock, 10, strconv.IntSize)
			val := int(v)
			maxBlock = &val
			if err != nil {
				return nil, fmt.Errorf("invalid value for maxBlock, must be INT")
			}
			if *maxBlock < 1 {
				err = goa.MergeErrors(err, goa.InvalidRangeError("max_block", *maxBlock, 1, true))
			}
			if err != nil {
				return nil, err
			}
		}
	}
	var isLikelyDupe *bool
	{
		if nftNftSearchIsLikelyDupe != "" {
			var val bool
			val, err = strconv.ParseBool(nftNftSearchIsLikelyDupe)
			isLikelyDupe = &val
			if err != nil {
				return nil, fmt.Errorf("invalid value for isLikelyDupe, must be BOOL")
			}
		}
	}
	var minRarenessScore *float64
	{
		if nftNftSearchMinRarenessScore != "" {
			val, err := strconv.ParseFloat(nftNftSearchMinRarenessScore, 64)
			minRarenessScore = &val
			if err != nil {
				return nil, fmt.Errorf("invalid value for minRarenessScore, must be FLOAT64")
			}
			if *minRarenessScore < 0 {
				err = goa.MergeErrors(err, goa.InvalidRangeError("min_rareness_score", *minRarenessScore, 0, true))
			}
			if *minRarenessScore > 1 {
				err = goa.MergeErrors(err, goa.InvalidRangeError("min_rareness_score", *minRarenessScore, 1, false))
			}
			if err != nil {
				return nil, err
			}
		}
	}
	var maxRarenessScore *float64
	{
		if nftNftSearchMaxRarenessScore != "" {
			val, err := strconv.ParseFloat(nftNftSearchMaxRarenessScore, 64)
			maxRarenessScore = &val
			if err != nil {
				return nil, fmt.Errorf("invalid value for maxRarenessScore, must be FLOAT64")
			}
			if *maxRarenessScore < 0 {
				err = goa.MergeErrors(err, goa.InvalidRangeError("max_rareness_score", *maxRarenessScore, 0, true))
			}
			if *maxRarenessScore > 1 {
				err = goa.MergeErrors(err, goa.InvalidRangeError("max_rareness_score", *maxRarenessScore, 1, false))
			}
			if err != nil {
				return nil, err
			}
		}
	}
	var minNsfwScore *float64
	{
		if nftNftSearchMinNsfwScore != "" {
			val, err := strconv.ParseFloat(nftNftSearchMinNsfwScore, 64)
			minNsfwScore = &val
			if err != nil {
				return nil, fmt.Errorf("invalid value for minNsfwScore, must be FLOAT64")
			}
			if *minNsfwScore < 0 {
				err = goa.MergeErrors(err, goa.InvalidRangeError("min_nsfw_score", *minNsfwScore, 0, true))
			}
			if *minNsfwScore > 1 {
				err = goa.MergeErrors(err, goa.InvalidRangeError("min_nsfw_score", *minNsfwScore, 1, false))
			}
			if err != nil {
				return nil, err
			}
		}
	}
	var maxNsfwScore *float64
	{
		if nftNftSearchMaxNsfwScore != "" {
			val, err := strconv.ParseFloat(nftNftSearchMaxNsfwScore, 64)
			maxNsfwScore = &val
			if err != nil {
				return nil, fmt.Errorf("invalid value for maxNsfwScore, must be FLOAT64")
			}
			if *maxNsfwScore < 0 {
				err = goa.MergeErrors(err, goa.InvalidRangeError("max_nsfw_score", *maxNsfwScore, 0, true))
			}
			if *maxNsfwScore > 1 {
				err = goa.MergeErrors(err, goa.InvalidRangeError("max_nsfw_score", *maxNsfwScore, 1, false))
			}
			if err != nil {
				return nil, err
			}
		}
	}
	var userPastelid *string
	{
		if nftNftSearchUserPastelid != "" {
			userPastelid = &nftNftSearchUserPastelid
			err = goa.MergeErrors(err, goa.ValidatePattern("user_pastelid", *userPastelid, "^[a-zA-Z0-9]+$"))
			if utf8.RuneCountInString(*userPastelid) < 86 {
				err = goa.MergeErrors(err, goa.InvalidLengthError("user_pastelid", *userPastelid, utf8.RuneCountInString(*userPastelid), 86, true))
			}
			if utf8.RuneCountInString(*userPastelid) > 86 {
				err = goa.MergeErrors(err, goa.InvalidLengthError("user_pastelid", *userPastelid, utf8.RuneCountInString(*userPastelid), 86, false))
			}
			if err != nil {
				return nil, err
			}
		}
	}
	var userPassphrase *string
	{
		if nftNftSearchUserPassphrase != "" {
			userPassphrase = &nftNftSearchUserPassphrase
		}
	}
	v := &nft.NftSearchPayload{}
	v.Artist = artist
	v.Limit = limit
	v.Query = query
	v.CreatorName = creatorName
	v.ArtTitle = artTitle
	v.Series = series
	v.Descr = descr
	v.Keyword = keyword
	v.MinCopies = minCopies
	v.MaxCopies = maxCopies
	v.MinBlock = minBlock
	v.MaxBlock = maxBlock
	v.IsLikelyDupe = isLikelyDupe
	v.MinRarenessScore = minRarenessScore
	v.MaxRarenessScore = maxRarenessScore
	v.MinNsfwScore = minNsfwScore
	v.MaxNsfwScore = maxNsfwScore
	v.UserPastelid = userPastelid
	v.UserPassphrase = userPassphrase

	return v, nil
}

// BuildNftGetPayload builds the payload for the nft nftGet endpoint from CLI
// flags.
func BuildNftGetPayload(nftNftGetTxid string, nftNftGetPid string, nftNftGetKey string) (*nft.NftGetPayload, error) {
	var err error
	var txid string
	{
		txid = nftNftGetTxid
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
		pid = nftNftGetPid
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
		key = nftNftGetKey
	}
	v := &nft.NftGetPayload{}
	v.Txid = txid
	v.Pid = pid
	v.Key = key

	return v, nil
}

// BuildDownloadPayload builds the payload for the nft download endpoint from
// CLI flags.
func BuildDownloadPayload(nftDownloadTxid string, nftDownloadPid string, nftDownloadKey string) (*nft.DownloadPayload, error) {
	var err error
	var txid string
	{
		txid = nftDownloadTxid
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
		pid = nftDownloadPid
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
		key = nftDownloadKey
	}
	v := &nft.DownloadPayload{}
	v.Txid = txid
	v.Pid = pid
	v.Key = key

	return v, nil
}

// BuildDdServiceOutputFileDetailPayload builds the payload for the nft
// ddServiceOutputFileDetail endpoint from CLI flags.
func BuildDdServiceOutputFileDetailPayload(nftDdServiceOutputFileDetailTxid string, nftDdServiceOutputFileDetailPid string, nftDdServiceOutputFileDetailKey string) (*nft.DownloadPayload, error) {
	var err error
	var txid string
	{
		txid = nftDdServiceOutputFileDetailTxid
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
		pid = nftDdServiceOutputFileDetailPid
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
		key = nftDdServiceOutputFileDetailKey
	}
	v := &nft.DownloadPayload{}
	v.Txid = txid
	v.Pid = pid
	v.Key = key

	return v, nil
}

// BuildDdServiceOutputFilePayload builds the payload for the nft
// ddServiceOutputFile endpoint from CLI flags.
func BuildDdServiceOutputFilePayload(nftDdServiceOutputFileTxid string, nftDdServiceOutputFilePid string, nftDdServiceOutputFileKey string) (*nft.DownloadPayload, error) {
	var err error
	var txid string
	{
		txid = nftDdServiceOutputFileTxid
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
		pid = nftDdServiceOutputFilePid
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
		key = nftDdServiceOutputFileKey
	}
	v := &nft.DownloadPayload{}
	v.Txid = txid
	v.Pid = pid
	v.Key = key

	return v, nil
}
