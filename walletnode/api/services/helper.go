package services

import (
	"github.com/pastelnetwork/gonode/common/service/userdata"
	"github.com/pastelnetwork/gonode/walletnode/api/gen/nft"
	"github.com/pastelnetwork/gonode/walletnode/api/gen/userdatas"
	"time"

	"github.com/pastelnetwork/gonode/pastel"

	"github.com/pastelnetwork/gonode/walletnode/services/nftsearch"

	"github.com/pastelnetwork/gonode/common/service/task/state"
)

func safeStr(p *string) string {
	if p != nil {
		return *p
	}
	return ""
}

func toNftStates(statuses []*state.Status) []*nft.TaskState {
	var states []*nft.TaskState

	for _, status := range statuses {
		states = append(states, &nft.TaskState{
			Date:   status.CreatedAt.Format(time.RFC3339),
			Status: status.String(),
		})
	}
	return states
}

// NFT Search
func toArtSearchResult(srch *nftsearch.RegTicketSearch) *nft.NftSearchResult {
	ticketData := srch.RegTicketData.NFTTicketData.AppTicketData
	res := &nft.NftSearchResult{
		Nft: &nft.NftSummary{
			Txid:       srch.TXID,
			Thumbnail1: srch.Thumbnail,
			Thumbnail2: srch.ThumbnailSecondry,
			Title:      ticketData.NFTTitle,

			Copies:           srch.RegTicketData.NFTTicketData.Copies,
			ArtistName:       ticketData.CreatorName,
			YoutubeURL:       &ticketData.NFTCreationVideoYoutubeURL,
			ArtistPastelID:   srch.RegTicketData.NFTTicketData.Author,
			ArtistWebsiteURL: &ticketData.CreatorWebsite,
			Description:      ticketData.CreatorWrittenStatement,
			Keywords:         &ticketData.NFTKeywordSet,
			SeriesName:       &ticketData.NFTSeriesName,
		},

		MatchIndex: srch.MatchIndex,
	}

	res.Matches = []*nft.FuzzyMatch{}
	for _, match := range srch.Matches {
		res.Matches = append(res.Matches, &nft.FuzzyMatch{
			Score:          &match.Score,
			Str:            &match.Str,
			FieldType:      &match.FieldType,
			MatchedIndexes: match.MatchedIndexes,
		})
	}

	return res
}

func toNftDetail(ticket *pastel.RegTicket) *nft.NftDetail {
	return &nft.NftDetail{
		Txid:             ticket.TXID,
		Title:            ticket.RegTicketData.NFTTicketData.AppTicketData.NFTTitle,
		Copies:           ticket.RegTicketData.NFTTicketData.Copies,
		ArtistName:       ticket.RegTicketData.NFTTicketData.AppTicketData.CreatorName,
		YoutubeURL:       &ticket.RegTicketData.NFTTicketData.AppTicketData.NFTCreationVideoYoutubeURL,
		ArtistPastelID:   ticket.RegTicketData.NFTTicketData.Author,
		ArtistWebsiteURL: &ticket.RegTicketData.NFTTicketData.AppTicketData.CreatorWebsite,
		Description:      ticket.RegTicketData.NFTTicketData.AppTicketData.CreatorWrittenStatement,
		Keywords:         &ticket.RegTicketData.NFTTicketData.AppTicketData.NFTKeywordSet,
		SeriesName:       &ticket.RegTicketData.NFTTicketData.AppTicketData.NFTSeriesName,
		Royalty:          &ticket.RegTicketData.Royalty,
		// ------------------- WIP: PSL-142 -------------figure out how to search for these ---------
		//RarenessScore:         ticket.RegTicketData.NFTTicketData.AppTicketData.PastelRarenessScore,
		//NsfwScore:             ticket.RegTicketData.NFTTicketData.AppTicketData.OpenNSFWScore,
		//InternetRarenessScore: &ticket.RegTicketData.NFTTicketData.AppTicketData.InternetRarenessScore,
		Version:    &ticket.RegTicketData.NFTTicketData.Version,
		StorageFee: &ticket.RegTicketData.StorageFee,
		//HentaiNsfwScore:       &ticket.RegTicketData.NFTTicketData.AppTicketData.AlternateNSFWScores.Hentai,
		//DrawingNsfwScore:      &ticket.RegTicketData.NFTTicketData.AppTicketData.AlternateNSFWScores.Drawing,
		//NeutralNsfwScore:      &ticket.RegTicketData.NFTTicketData.AppTicketData.AlternateNSFWScores.Neutral,
		//PornNsfwScore:         &ticket.RegTicketData.NFTTicketData.AppTicketData.AlternateNSFWScores.Porn,
		//SexyNsfwScore:         &ticket.RegTicketData.NFTTicketData.AppTicketData.AlternateNSFWScores.Sexy,
	}
}

// NFT User Data
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

	request.UserPastelID = req.UserPastelID
	request.UserPastelIDPassphrase = req.UserPastelIDPassphrase

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

	request.UserPastelID = req.UserPastelID
	request.UserPastelIDPassphrase = req.UserPastelIDPassphrase

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
	result.UserPastelID = req.UserPastelID
	result.UserPastelIDPassphrase = "" // No need to reture ArtistPastelIDPassphrase

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
