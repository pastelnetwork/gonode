package services

import (
	"time"

	"github.com/pastelnetwork/gonode/walletnode/api/gen/nft"

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
			Txid:          srch.TXID,
			Thumbnail1:    srch.Thumbnail,
			Thumbnail2:    srch.ThumbnailSecondry,
			NsfwScore:     &srch.OpenNSFWScore,
			RarenessScore: &srch.RarenessScore,
			IsLikelyDupe:  &srch.IsLikelyDupe,
			Title:         ticketData.NFTTitle,

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

func translateNftSummary(res *nft.DDServiceOutputFileResult, ticket *pastel.RegTicket) *nft.DDServiceOutputFileResult {
	ticketData := ticket.RegTicketData.NFTTicketData.AppTicketData

	res.CreatorName = ticketData.CreatorName
	res.CreatorWebsite = ticketData.CreatorWebsite
	res.CreatorWrittenStatement = ticketData.CreatorWrittenStatement
	res.NftTitle = ticketData.NFTTitle
	res.NftSeriesName = ticketData.NFTSeriesName
	res.NftCreationVideoYoutubeURL = ticketData.NFTCreationVideoYoutubeURL
	res.NftKeywordSet = ticketData.NFTKeywordSet
	res.TotalCopies = ticketData.TotalCopies
	res.PreviewHash = ticketData.PreviewHash
	res.Thumbnail1Hash = ticketData.Thumbnail1Hash
	res.Thumbnail2Hash = ticketData.Thumbnail2Hash
	res.OriginalFileSizeInBytes = ticketData.OriginalFileSizeInBytes
	res.FileType = ticketData.FileType
	res.MaxPermittedOpenNsfwScore = ticketData.MaxPermittedOpenNSFWScore

	return res
}

func translateDDServiceOutputFile(res *nft.DDServiceOutputFileResult, ddAndFpStruct *pastel.DDAndFingerprints) *nft.DDServiceOutputFileResult {
	res.PastelBlockHashWhenRequestSubmitted = &ddAndFpStruct.BlockHash
	res.PastelBlockHeightWhenRequestSubmitted = &ddAndFpStruct.BlockHeight
	res.UtcTimestampWhenRequestSubmitted = &ddAndFpStruct.TimestampOfRequest
	res.PastelIDOfSubmitter = &ddAndFpStruct.SubmitterPastelID
	res.PastelIDOfRegisteringSupernode1 = &ddAndFpStruct.SN1PastelID
	res.PastelIDOfRegisteringSupernode2 = &ddAndFpStruct.SN2PastelID
	res.PastelIDOfRegisteringSupernode3 = &ddAndFpStruct.SN3PastelID
	res.IsPastelOpenapiRequest = &ddAndFpStruct.IsOpenAPIRequest
	res.DupeDetectionSystemVersion = &ddAndFpStruct.DupeDetectionSystemVersion
	res.IsLikelyDupe = &ddAndFpStruct.IsLikelyDupe
	res.IsRareOnInternet = &ddAndFpStruct.IsRareOnInternet
	res.OverallRarenessScore = &ddAndFpStruct.OverallRarenessScore
	res.PctOfTop10MostSimilarWithDupeProbAbove25pct = &ddAndFpStruct.PctOfTop10MostSimilarWithDupeProbAbove25pct
	res.PctOfTop10MostSimilarWithDupeProbAbove33pct = &ddAndFpStruct.PctOfTop10MostSimilarWithDupeProbAbove33pct
	res.PctOfTop10MostSimilarWithDupeProbAbove50pct = &ddAndFpStruct.PctOfTop10MostSimilarWithDupeProbAbove50pct
	res.RarenessScoresTableJSONCompressedB64 = &ddAndFpStruct.RarenessScoresTableJSONCompressedB64
	res.OpenNsfwScore = &ddAndFpStruct.OpenNSFWScore
	res.ImageFingerprintOfCandidateImageFile = ddAndFpStruct.ImageFingerprintOfCandidateImageFile
	res.HashOfCandidateImageFile = &ddAndFpStruct.HashOfCandidateImageFile
	res.RarenessScoresTableJSONCompressedB64 = &ddAndFpStruct.InternetRareness.RareOnInternetSummaryTableAsJSONCompressedB64
	res.CollectionNameString = &ddAndFpStruct.CollectionNameString
	res.OpenAPIGroupIDString = &ddAndFpStruct.OpenAPIGroupIDString
	res.GroupRarenessScore = &ddAndFpStruct.GroupRarenessScore
	res.CandidateImageThumbnailWebpAsBase64String = &ddAndFpStruct.CandidateImageThumbnailWebpAsBase64String
	res.DoesNotImpactTheFollowingCollectionStrings = &ddAndFpStruct.DoesNotImpactTheFollowingCollectionStrings
	res.SimilarityScoreToFirstEntryInCollection = &ddAndFpStruct.SimilarityScoreToFirstEntryInCollection
	res.CpProbability = &ddAndFpStruct.CPProbability
	res.ChildProbability = &ddAndFpStruct.ChildProbability
	res.ImageFilePath = &ddAndFpStruct.ImageFilePath
	res.InternetRareness = &nft.InternetRareness{
		RareOnInternetSummaryTableAsJSONCompressedB64:    &ddAndFpStruct.InternetRareness.RareOnInternetSummaryTableAsJSONCompressedB64,
		RareOnInternetGraphJSONCompressedB64:             &ddAndFpStruct.InternetRareness.RareOnInternetGraphJSONCompressedB64,
		AlternativeRareOnInternetDictAsJSONCompressedB64: &ddAndFpStruct.InternetRareness.AlternativeRareOnInternetDictAsJSONCompressedB64,
		MinNumberOfExactMatchesInPage:                    &ddAndFpStruct.InternetRareness.MinNumberOfExactMatchesInPage,
		EarliestAvailableDateOfInternetResults:           &ddAndFpStruct.InternetRareness.EarliestAvailableDateOfInternetResults,
	}
	res.AlternativeNsfwScores = &nft.AlternativeNSFWScores{
		Drawings: &ddAndFpStruct.AlternativeNSFWScores.Drawings,
		Sexy:     &ddAndFpStruct.AlternativeNSFWScores.Sexy,
		Porn:     &ddAndFpStruct.AlternativeNSFWScores.Porn,
		Hentai:   &ddAndFpStruct.AlternativeNSFWScores.Neutral,
	}

	return res
}
