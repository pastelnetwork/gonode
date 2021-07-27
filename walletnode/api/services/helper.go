package services

import (
	"time"

	"github.com/pastelnetwork/gonode/pastel"

	"github.com/pastelnetwork/gonode/walletnode/services/artworksearch"

	"github.com/pastelnetwork/gonode/common/service/task/state"
	"github.com/pastelnetwork/gonode/walletnode/api/gen/artworks"
	"github.com/pastelnetwork/gonode/walletnode/services/artworkregister"

	"github.com/pastelnetwork/gonode/common/service/userdata"
	"github.com/pastelnetwork/gonode/walletnode/api/gen/userdatas"
)

func fromRegisterPayload(payload *artworks.RegisterPayload) *artworkregister.Ticket {
	return &artworkregister.Ticket{
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
	}
}

func toArtworkTicket(ticket *artworkregister.Ticket) *artworks.ArtworkTicket {
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
	ticketData := srch.RegTicketData.ArtTicketData.AppTicketData
	res := &artworks.ArtworkSearchResult{
		Artwork: &artworks.ArtworkSummary{
			Txid:      srch.TXID,
			Thumbnail: srch.Thumbnail,
			Title:     ticketData.ArtworkTitle,

			Copies:           srch.RegTicketData.ArtTicketData.Copies,
			ArtistName:       ticketData.ArtistName,
			YoutubeURL:       &ticketData.ArtworkCreationVideoYoutubeURL,
			ArtistPastelID:   ticketData.AuthorPastelID,
			ArtistWebsiteURL: &ticketData.ArtistWebsite,
			Description:      ticketData.ArtistWrittenStatement,
			Keywords:         &ticketData.ArtworkKeywordSet,
			SeriesName:       &ticketData.ArtworkSeriesName,
		},

		MatchIndex: srch.MatchIndex,
	}

	res.Matches = []*artworks.FuzzyMatch{}
	for _, match := range res.Matches {
		res.Matches = append(res.Matches, &artworks.FuzzyMatch{
			Score:          match.Score,
			Str:            match.Str,
			FieldType:      match.FieldType,
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
	}
}

func toArtworkDetail(ticket *pastel.RegTicket) *artworks.ArtworkDetail {
	return &artworks.ArtworkDetail{
		Txid:             ticket.TXID,
		Title:            ticket.RegTicketData.ArtTicketData.AppTicketData.ArtworkTitle,
		Copies:           ticket.RegTicketData.ArtTicketData.Copies,
		ArtistName:       ticket.RegTicketData.ArtTicketData.AppTicketData.ArtistName,
		YoutubeURL:       &ticket.RegTicketData.ArtTicketData.AppTicketData.ArtworkCreationVideoYoutubeURL,
		ArtistPastelID:   ticket.RegTicketData.ArtTicketData.AppTicketData.AuthorPastelID,
		ArtistWebsiteURL: &ticket.RegTicketData.ArtTicketData.AppTicketData.ArtistWebsite,
		Description:      ticket.RegTicketData.ArtTicketData.AppTicketData.ArtistWrittenStatement,
		Keywords:         &ticket.RegTicketData.ArtTicketData.AppTicketData.ArtworkKeywordSet,
		SeriesName:       &ticket.RegTicketData.ArtTicketData.AppTicketData.ArtworkSeriesName,
		IsGreen:          ticket.RegTicketData.IsGreen,
		Royalty:          float64(ticket.RegTicketData.Royalty),
		RarenessScore:    ticket.RegTicketData.ArtTicketData.AppTicketData.RarenessScore,
		NsfwScore:        ticket.RegTicketData.ArtTicketData.AppTicketData.NSFWScore,
		SeenScore:        ticket.RegTicketData.ArtTicketData.AppTicketData.SeenScore,
		Version:          &ticket.RegTicketData.ArtTicketData.Version,
		StorageFee:       &ticket.RegTicketData.StorageFee,
	}
}

// fromUserdataProcessRequest convert the request receive from swagger api to request object that will send to super nodes
func fromUserdataProcessRequest(req *userdatas.ProcessUserdataPayload) *userdata.UserdataProcessRequest {
	request := &userdata.UserdataProcessRequest{}

	if req.Realname != nil 					{ request.Realname 			= *(req.Realname) }
	if req.FacebookLink != nil 				{ request.FacebookLink 		= *(req.FacebookLink) }
	if req.TwitterLink != nil 				{ request.TwitterLink 		= *(req.TwitterLink) }
	if req.Location != nil 					{ request.Location 			= *(req.Location) }
	if req.PrimaryLanguage != nil 			{ request.PrimaryLanguage 	= *(req.PrimaryLanguage) }
	if req.Categories != nil 				{ request.Categories 		= *(req.Categories) }
	if req.Biography != nil 				{ request.Biography 		= *(req.Biography) }
	if req.AvatarImage != nil { 
		if req.AvatarImage.Content != nil && len(req.AvatarImage.Content) > 0 {
			request.AvatarImage.Content = make ([]byte, len(req.AvatarImage.Content))
			copy(request.AvatarImage.Content,req.AvatarImage.Content)
		}
		request.AvatarImage.Filename = *(req.AvatarImage.Filename)
	}
	if req.CoverPhoto != nil { 
		if req.CoverPhoto.Content != nil && len(req.CoverPhoto.Content) > 0 {
			request.CoverPhoto.Content = make ([]byte, len(req.CoverPhoto.Content))
			copy(request.CoverPhoto.Content,req.CoverPhoto.Content)
		}
		request.CoverPhoto.Filename = *(req.CoverPhoto.Filename)
	}

	request.ArtistPastelID 	= req.ArtistPastelID
	request.ArtistPastelIDPassphrase = req.ArtistPastelIDPassphrase

	request.Timestamp 	 					= time.Now().Unix() 	// The moment request is prepared to send to super nodes
	request.PreviousBlockHash				= ""					// PreviousBlockHash will be generated later when prepare to send
	
	return request
}

// toUserdataProcessResult convert the final response receive from super nodes and reponse to swagger api
func toUserdataProcessResult(result *userdata.UserdataProcessResult) *userdatas.UserdataProcessResult {
	res := &userdatas.UserdataProcessResult{
			ResponseCode: 		int(result.ResponseCode),
			Detail:				result.Detail,
			Realname:			&result.Realname,
			FacebookLink:		&result.FacebookLink,
			TwitterLink:		&result.TwitterLink,
			NativeCurrency:		&result.NativeCurrency,
			Location:			&result.Location,
			PrimaryLanguage:	&result.PrimaryLanguage,
			Categories:			&result.Categories,
			AvatarImage:		&result.AvatarImage,
			CoverPhoto:			&result.CoverPhoto,
		}
	return res
}

func toUserSpecifiedData (req *userdata.UserdataProcessRequest) *userdatas.UserSpecifiedData {
	result := &userdatas.UserSpecifiedData {}
	
	result.Realname = &req.Realname
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
			req.AvatarImage.Content = make ([]byte, len(result.AvatarImage.Content))
			copy(req.AvatarImage.Content,result.AvatarImage.Content)
		}
		req.AvatarImage.Filename = *result.AvatarImage.Filename
	}
	if result.CoverPhoto != nil { 
		if result.CoverPhoto.Content != nil && len(result.CoverPhoto.Content) > 0 {
			req.CoverPhoto.Content = make ([]byte, len(result.CoverPhoto.Content))
			copy(req.CoverPhoto.Content,result.CoverPhoto.Content)
		}
		req.CoverPhoto.Filename = *result.CoverPhoto.Filename
	}

	return result
}