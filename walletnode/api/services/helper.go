package services

import (
	"github.com/pastelnetwork/gonode/walletnode/api/gen/nft"
	"time"

	"github.com/pastelnetwork/gonode/pastel"

	"github.com/pastelnetwork/gonode/walletnode/services/nftsearch"

	"github.com/pastelnetwork/gonode/common/service/task/state"
)

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
func toNftSearchResult(srch *nftsearch.RegTicketSearch) *nft.NftSearchResult {
	ticketData := srch.RegTicketData.NFTTicketData.AppTicketData
	res := &nft.NftSearchResult{
		Nft: &nft.NftSummary{
			Txid:       srch.TXID,
			Thumbnail1: srch.Thumbnail,
			Thumbnail2: srch.ThumbnailSecondry,
			Title:      ticketData.NFTTitle,

			Copies:            srch.RegTicketData.NFTTicketData.Copies,
			CreatorName:       ticketData.CreatorName,
			YoutubeURL:        &ticketData.NFTCreationVideoYoutubeURL,
			CreatorPastelID:   srch.RegTicketData.NFTTicketData.Author,
			CreatorWebsiteURL: &ticketData.CreatorWebsite,
			Description:       ticketData.CreatorWrittenStatement,
			Keywords:          &ticketData.NFTKeywordSet,
			SeriesName:        &ticketData.NFTSeriesName,
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
		Txid:              ticket.TXID,
		Title:             ticket.RegTicketData.NFTTicketData.AppTicketData.NFTTitle,
		Copies:            ticket.RegTicketData.NFTTicketData.Copies,
		CreatorName:       ticket.RegTicketData.NFTTicketData.AppTicketData.CreatorName,
		YoutubeURL:        &ticket.RegTicketData.NFTTicketData.AppTicketData.NFTCreationVideoYoutubeURL,
		CreatorPastelID:   ticket.RegTicketData.NFTTicketData.Author,
		CreatorWebsiteURL: &ticket.RegTicketData.NFTTicketData.AppTicketData.CreatorWebsite,
		Description:       ticket.RegTicketData.NFTTicketData.AppTicketData.CreatorWrittenStatement,
		Keywords:          &ticket.RegTicketData.NFTTicketData.AppTicketData.NFTKeywordSet,
		SeriesName:        &ticket.RegTicketData.NFTTicketData.AppTicketData.NFTSeriesName,
		Royalty:           &ticket.RegTicketData.Royalty,
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
