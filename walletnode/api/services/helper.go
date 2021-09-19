package services

import (
	"time"

	"github.com/pastelnetwork/gonode/pastel"

	"github.com/pastelnetwork/gonode/walletnode/services/artworksearch"

	"github.com/pastelnetwork/gonode/common/service/task/state"

	"github.com/pastelnetwork/gonode/common/service/artwork"
	"github.com/pastelnetwork/gonode/walletnode/api/gen/artworks"
	"github.com/pastelnetwork/gonode/walletnode/services/artworkdownload"
	"github.com/pastelnetwork/gonode/walletnode/services/artworkregister"

	"github.com/pastelnetwork/gonode/common/service/userdata"
	"github.com/pastelnetwork/gonode/walletnode/api/gen/userdatas"
)

func safeStr(p *string) string {
	if p != nil {
		return *p
	}
	return ""
}

func fromRegisterPayload(payload *artworks.RegisterPayload) *artworkregister.Request {
	thumbnail := artwork.ThumbnailCoordinate{
		TopLeftX:     payload.ThumbnailCoordinate.TopLeftX,
		TopLeftY:     payload.ThumbnailCoordinate.TopLeftY,
		BottomRightX: payload.ThumbnailCoordinate.BottomRightX,
		BottomRightY: payload.ThumbnailCoordinate.BottomRightY,
	}

	return &artworkregister.Request{
		Name:                     payload.Name,
		Description:              payload.Description,
		Keywords:                 payload.Keywords,
		SeriesName:               payload.SeriesName,
		IssuedCopies:             payload.IssuedCopies,
		YoutubeURL:               payload.YoutubeURL,
		ArtistPastelID:           payload.ArtistPastelID,
		ArtistPastelIDPassphrase: payload.ArtistPastelIDPassphrase,
		ArtistName:               payload.ArtistName,
		ArtistWebsiteURL:         payload.ArtistWebsiteURL,
		SpendableAddress:         payload.SpendableAddress,
		MaximumFee:               payload.MaximumFee,
		Green:                    payload.Green,
		Royalty:                  payload.Royalty,
		Thumbnail:                thumbnail,
	}
}

func toArtworkTicket(ticket *artworkregister.Request) *artworks.ArtworkTicket {
	thumbnail := artworks.Thumbnailcoordinate{
		TopLeftX:     ticket.Thumbnail.TopLeftX,
		TopLeftY:     ticket.Thumbnail.TopLeftY,
		BottomRightY: ticket.Thumbnail.BottomRightX,
		BottomRightX: ticket.Thumbnail.BottomRightY,
	}
	return &artworks.ArtworkTicket{
		Name:                     ticket.Name,
		Description:              ticket.Description,
		Keywords:                 ticket.Keywords,
		SeriesName:               ticket.SeriesName,
		IssuedCopies:             ticket.IssuedCopies,
		YoutubeURL:               ticket.YoutubeURL,
		ArtistPastelID:           ticket.ArtistPastelID,
		ArtistPastelIDPassphrase: ticket.ArtistPastelIDPassphrase,
		ArtistName:               ticket.ArtistName,
		ArtistWebsiteURL:         ticket.ArtistWebsiteURL,
		SpendableAddress:         ticket.SpendableAddress,
		MaximumFee:               ticket.MaximumFee,
		Green:                    ticket.Green,
		Royalty:                  ticket.Royalty,
		ThumbnailCoordinate:      &thumbnail,
	}
}

func toArtworkStates(statuses []*state.Status) []*artworks.TaskState {
	var states []*artworks.TaskState

	for _, status := range statuses {
		states = append(states, &artworks.TaskState{
			Date:   status.CreatedAt.Format(time.RFC3339),
			Status: status.String(),
		})
	}
	return states
}

func toArtSearchResult(srch *artworksearch.RegTicketSearch) *artworks.ArtworkSearchResult {
	ticketData := srch.RegTicketData.NFTTicketData.AppTicketData
	res := &artworks.ArtworkSearchResult{
		Artwork: &artworks.ArtworkSummary{
			Txid:       srch.TXID,
			Thumbnail1: srch.Thumbnail,
			Thumbnail2: srch.ThumbnailSecondry,
			Title:      ticketData.NFTTitle,

			Copies:           srch.RegTicketData.NFTTicketData.Copies,
			ArtistName:       ticketData.CreatorName,
			YoutubeURL:       &ticketData.NFTCreationVideoYoutubeURL,
			ArtistPastelID:   ticketData.AuthorPastelID,
			ArtistWebsiteURL: &ticketData.CreatorWebsite,
			Description:      ticketData.CreatorWrittenStatement,
			Keywords:         &ticketData.NFTKeywordSet,
			SeriesName:       &ticketData.NFTSeriesName,
		},

		MatchIndex: srch.MatchIndex,
	}

	res.Matches = []*artworks.FuzzyMatch{}
	for _, match := range srch.Matches {
		res.Matches = append(res.Matches, &artworks.FuzzyMatch{
			Score:          &match.Score,
			Str:            &match.Str,
			FieldType:      &match.FieldType,
			MatchedIndexes: match.MatchedIndexes,
		})
	}

	return res
}

func fromArtSearchRequest(req *artworks.ArtSearchPayload) *artworksearch.ArtSearchRequest {
	return &artworksearch.ArtSearchRequest{
		Artist:           req.Artist,
		Limit:            req.Limit,
		Query:            req.Query,
		ArtistName:       req.ArtistName,
		ArtTitle:         req.ArtTitle,
		Series:           req.Series,
		Descr:            req.Descr,
		Keyword:          req.Keyword,
		MinBlock:         req.MinBlock,
		MaxBlock:         req.MaxBlock,
		MinCopies:        req.MinCopies,
		MaxCopies:        req.MaxCopies,
		MinNsfwScore:     req.MinNsfwScore,
		MaxNsfwScore:     req.MaxNsfwScore,
		MinRarenessScore: req.MinRarenessScore,
		MaxRarenessScore: req.MaxRarenessScore,
		UserPastelID:     *req.UserPastelid,
		UserPassphrase:   *req.UserPassphrase,
	}
}

func toArtworkDetail(ticket *pastel.RegTicket) *artworks.ArtworkDetail {
	return &artworks.ArtworkDetail{
		Txid:                  ticket.TXID,
		Title:                 ticket.RegTicketData.NFTTicketData.AppTicketData.NFTTitle,
		Copies:                ticket.RegTicketData.NFTTicketData.Copies,
		ArtistName:            ticket.RegTicketData.NFTTicketData.AppTicketData.CreatorName,
		YoutubeURL:            &ticket.RegTicketData.NFTTicketData.AppTicketData.NFTCreationVideoYoutubeURL,
		ArtistPastelID:        ticket.RegTicketData.NFTTicketData.AppTicketData.AuthorPastelID,
		ArtistWebsiteURL:      &ticket.RegTicketData.NFTTicketData.AppTicketData.CreatorWebsite,
		Description:           ticket.RegTicketData.NFTTicketData.AppTicketData.CreatorWrittenStatement,
		Keywords:              &ticket.RegTicketData.NFTTicketData.AppTicketData.NFTKeywordSet,
		SeriesName:            &ticket.RegTicketData.NFTTicketData.AppTicketData.NFTSeriesName,
		Royalty:               &ticket.RegTicketData.Royalty,
		RarenessScore:         ticket.RegTicketData.NFTTicketData.AppTicketData.PastelRarenessScore,
		NsfwScore:             ticket.RegTicketData.NFTTicketData.AppTicketData.OpenNSFWScore,
		InternetRarenessScore: &ticket.RegTicketData.NFTTicketData.AppTicketData.InternetRarenessScore,
		Version:               &ticket.RegTicketData.NFTTicketData.Version,
		StorageFee:            &ticket.RegTicketData.StorageFee,
		HentaiNsfwScore:       &ticket.RegTicketData.NFTTicketData.AppTicketData.AlternateNSFWScores.Hentai,
		DrawingNsfwScore:      &ticket.RegTicketData.NFTTicketData.AppTicketData.AlternateNSFWScores.Drawing,
		NeutralNsfwScore:      &ticket.RegTicketData.NFTTicketData.AppTicketData.AlternateNSFWScores.Neutral,
		PornNsfwScore:         &ticket.RegTicketData.NFTTicketData.AppTicketData.AlternateNSFWScores.Porn,
		SexyNsfwScore:         &ticket.RegTicketData.NFTTicketData.AppTicketData.AlternateNSFWScores.Sexy,
	}
}

func fromDownloadPayload(payload *artworks.ArtworkDownloadPayload) *artworkdownload.Ticket {
	return &artworkdownload.Ticket{
		Txid:               payload.Txid,
		PastelID:           payload.Pid,
		PastelIDPassphrase: payload.Key,
	}
}

// fromUserdataCreateRequest convert the request receive from swagger api to request object that will send to super nodes
func fromUserdataCreateRequest(req *userdatas.CreateUserdataPayload) *userdata.ProcessRequest {
	request := &userdata.ProcessRequest{}

	request.RealName = safeStr(req.RealName)
	request.FacebookLink = safeStr(req.FacebookLink)
	request.TwitterLink = safeStr(req.TwitterLink)
	request.Location = safeStr(req.Location)
	request.PrimaryLanguage = safeStr(req.PrimaryLanguage)
	request.Categories = safeStr(req.Categories)
	request.Biography = safeStr(req.Biography)

	if req.AvatarImage != nil {
		if req.AvatarImage.Content != nil && len(req.AvatarImage.Content) > 0 {
			request.AvatarImage.Content = make([]byte, len(req.AvatarImage.Content))
			copy(request.AvatarImage.Content, req.AvatarImage.Content)
		}
		request.AvatarImage.Filename = safeStr(req.AvatarImage.Filename)
	}
	if req.CoverPhoto != nil {
		if req.CoverPhoto.Content != nil && len(req.CoverPhoto.Content) > 0 {
			request.CoverPhoto.Content = make([]byte, len(req.CoverPhoto.Content))
			copy(request.CoverPhoto.Content, req.CoverPhoto.Content)
		}
		request.CoverPhoto.Filename = safeStr(req.CoverPhoto.Filename)
	}

	request.ArtistPastelID = req.ArtistPastelID
	request.ArtistPastelIDPassphrase = req.ArtistPastelIDPassphrase

	request.Timestamp = time.Now().Unix() // The moment request is prepared to send to super nodes
	request.PreviousBlockHash = ""        // PreviousBlockHash will be generated later when prepare to send

	return request
}

// fromUserdataUpdateRequest convert the request receive from swagger api to request object that will send to super nodes
func fromUserdataUpdateRequest(req *userdatas.UpdateUserdataPayload) *userdata.ProcessRequest {
	request := &userdata.ProcessRequest{}

	request.RealName = safeStr(req.RealName)
	request.FacebookLink = safeStr(req.FacebookLink)
	request.TwitterLink = safeStr(req.TwitterLink)
	request.Location = safeStr(req.Location)
	request.PrimaryLanguage = safeStr(req.PrimaryLanguage)
	request.Categories = safeStr(req.Categories)
	request.Biography = safeStr(req.Biography)

	if req.AvatarImage != nil {
		if req.AvatarImage.Content != nil && len(req.AvatarImage.Content) > 0 {
			request.AvatarImage.Content = make([]byte, len(req.AvatarImage.Content))
			copy(request.AvatarImage.Content, req.AvatarImage.Content)
		}
		request.AvatarImage.Filename = safeStr(req.AvatarImage.Filename)
	}
	if req.CoverPhoto != nil {
		if req.CoverPhoto.Content != nil && len(req.CoverPhoto.Content) > 0 {
			request.CoverPhoto.Content = make([]byte, len(req.CoverPhoto.Content))
			copy(request.CoverPhoto.Content, req.CoverPhoto.Content)
		}
		request.CoverPhoto.Filename = safeStr(req.CoverPhoto.Filename)
	}

	request.ArtistPastelID = req.ArtistPastelID
	request.ArtistPastelIDPassphrase = req.ArtistPastelIDPassphrase

	request.Timestamp = time.Now().Unix() // The moment request is prepared to send to super nodes
	request.PreviousBlockHash = ""        // PreviousBlockHash will be generated later when prepare to send

	return request
}

// toUserdataProcessResult convert the final response receive from super nodes and reponse to swagger api
func toUserdataProcessResult(result *userdata.ProcessResult) *userdatas.UserdataProcessResult {
	res := &userdatas.UserdataProcessResult{
		ResponseCode:    int(result.ResponseCode),
		Detail:          result.Detail,
		RealName:        &result.RealName,
		FacebookLink:    &result.FacebookLink,
		TwitterLink:     &result.TwitterLink,
		NativeCurrency:  &result.NativeCurrency,
		Location:        &result.Location,
		PrimaryLanguage: &result.PrimaryLanguage,
		Categories:      &result.Categories,
		AvatarImage:     &result.AvatarImage,
		CoverPhoto:      &result.CoverPhoto,
	}
	return res
}

func toUserSpecifiedData(req *userdata.ProcessRequest) *userdatas.UserSpecifiedData {
	result := &userdatas.UserSpecifiedData{}

	result.RealName = &req.RealName
	result.FacebookLink = &req.FacebookLink
	result.TwitterLink = &req.TwitterLink
	result.NativeCurrency = &req.NativeCurrency
	result.Location = &req.Location
	result.PrimaryLanguage = &req.PrimaryLanguage
	result.Categories = &req.Categories
	result.Biography = &req.Biography
	result.ArtistPastelID = req.ArtistPastelID
	result.ArtistPastelIDPassphrase = "" // No need to reture ArtistPastelIDPassphrase

	if result.AvatarImage != nil {
		if result.AvatarImage.Content != nil && len(result.AvatarImage.Content) > 0 {
			req.AvatarImage.Content = make([]byte, len(result.AvatarImage.Content))
			copy(req.AvatarImage.Content, result.AvatarImage.Content)
		}
		req.AvatarImage.Filename = safeStr(result.AvatarImage.Filename)
	}
	if result.CoverPhoto != nil {
		if result.CoverPhoto.Content != nil && len(result.CoverPhoto.Content) > 0 {
			req.CoverPhoto.Content = make([]byte, len(result.CoverPhoto.Content))
			copy(req.CoverPhoto.Content, result.CoverPhoto.Content)
		}
		req.CoverPhoto.Filename = safeStr(result.CoverPhoto.Filename)
	}

	return result
}
