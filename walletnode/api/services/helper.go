package services

import (
	"time"

	"github.com/pastelnetwork/gonode/walletnode/api/gen/nft"

	"github.com/pastelnetwork/gonode/pastel"

	"github.com/pastelnetwork/gonode/walletnode/services/nftsearch"

	"github.com/pastelnetwork/gonode/common/service/task/state"
	"github.com/pastelnetwork/gonode/common/types"
	"github.com/pastelnetwork/gonode/walletnode/api/gen/metrics"
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

func toNFTDDServiceFile(ticketData pastel.AppTicket, ddAndFpStruct *pastel.DDAndFingerprints) *DDServiceOutputFileResult {
	res := &DDServiceOutputFileResult{}
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
	res.CollectionNameString = &ddAndFpStruct.CollectionNameString
	res.OpenAPIGroupIDString = &ddAndFpStruct.OpenAPIGroupIDString
	res.GroupRarenessScore = &ddAndFpStruct.GroupRarenessScore
	res.CandidateImageThumbnailWebpAsBase64String = &ddAndFpStruct.CandidateImageThumbnailWebpAsBase64String
	res.DoesNotImpactTheFollowingCollectionStrings = &ddAndFpStruct.DoesNotImpactTheFollowingCollectionStrings
	res.SimilarityScoreToFirstEntryInCollection = &ddAndFpStruct.SimilarityScoreToFirstEntryInCollection
	res.CpProbability = &ddAndFpStruct.CPProbability
	res.ChildProbability = &ddAndFpStruct.ChildProbability
	res.ImageFilePath = &ddAndFpStruct.ImageFilePath
	res.InternetRareness = &InternetRareness{
		RareOnInternetSummaryTableAsJSONCompressedB64:    &ddAndFpStruct.InternetRareness.RareOnInternetSummaryTableAsJSONCompressedB64,
		RareOnInternetGraphJSONCompressedB64:             &ddAndFpStruct.InternetRareness.RareOnInternetGraphJSONCompressedB64,
		AlternativeRareOnInternetDictAsJSONCompressedB64: &ddAndFpStruct.InternetRareness.AlternativeRareOnInternetDictAsJSONCompressedB64,
		MinNumberOfExactMatchesInPage:                    &ddAndFpStruct.InternetRareness.MinNumberOfExactMatchesInPage,
		EarliestAvailableDateOfInternetResults:           &ddAndFpStruct.InternetRareness.EarliestAvailableDateOfInternetResults,
	}
	res.AlternativeNsfwScores = &AlternativeNSFWScores{
		Drawings: &ddAndFpStruct.AlternativeNSFWScores.Drawings,
		Sexy:     &ddAndFpStruct.AlternativeNSFWScores.Sexy,
		Porn:     &ddAndFpStruct.AlternativeNSFWScores.Porn,
		Hentai:   &ddAndFpStruct.AlternativeNSFWScores.Neutral,
	}

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

// DDServiceOutputFileResult represents DD Service Output File for NFT
type DDServiceOutputFileResult struct {
	// block hash when request submitted
	PastelBlockHashWhenRequestSubmitted *string `json:"pastel_block_hash_when_request_submitted"`
	// block Height when request submitted
	PastelBlockHeightWhenRequestSubmitted *string `json:"pastel_block_height_when_request_submitted"`
	// timestamp of request when submitted
	UtcTimestampWhenRequestSubmitted *string `json:"utc_timestamp_when_request_submitted"`
	// pastel id of the submitter
	PastelIDOfSubmitter *string `json:"pastel_id_of_submitter"`
	// pastel id of registering SN1
	PastelIDOfRegisteringSupernode1 *string `json:"pastel_id_of_registering_supernode_1"`
	// pastel id of registering SN2
	PastelIDOfRegisteringSupernode2 *string `json:"pastel_id_of_registering_supernode_2"`
	// pastel id of registering SN3
	PastelIDOfRegisteringSupernode3 *string `json:"pastel_id_of_registering_supernode_3"`
	// is pastel open API request
	IsPastelOpenapiRequest *bool `json:"is_pastel_openapi_request"`
	// system version of dupe detection
	DupeDetectionSystemVersion *string `json:"dupe_detection_system_version"`
	// is this nft likely a duplicate
	IsLikelyDupe *bool `json:"is_likely_dupe"`
	// is this nft rare on the internet
	IsRareOnInternet *bool `json:"is_rare_on_internet"`
	// pastel rareness score
	OverallRarenessScore *float32 `json:"overall_rareness_score"`
	// PCT of top 10 most similar with dupe probe above 25 PCT
	PctOfTop10MostSimilarWithDupeProbAbove25pct *float32 `json:"pct_of_top_10_most_similar_with_dupe_prob_above_25pct"`
	// PCT of top 10 most similar with dupe probe above 33 PCT
	PctOfTop10MostSimilarWithDupeProbAbove33pct *float32 `json:"pct_of_top_10_most_similar_with_dupe_prob_above_33pct"`
	// PCT of top 10 most similar with dupe probe above 50 PCT
	PctOfTop10MostSimilarWithDupeProbAbove50pct *float32 `json:"pct_of_top_10_most_similar_with_dupe_prob_above_50pct"`
	// rareness scores table json compressed b64
	RarenessScoresTableJSONCompressedB64 *string `json:"rareness_scores_table_json_compressed_b64"`
	// open nsfw score
	OpenNsfwScore *float32 `json:"open_nsfw_score"`
	// Image fingerprint of candidate image file
	ImageFingerprintOfCandidateImageFile []float64 `json:"image_fingerprint_of_candidate_image_file"`
	// hash of candidate image file
	HashOfCandidateImageFile *string `json:"hash_of_candidate_image_file"`
	// name of the collection
	CollectionNameString *string `json:"collection_name_string"`
	// open api group id string
	OpenAPIGroupIDString *string `json:"open_api_group_id_string"`
	// rareness score of the group
	GroupRarenessScore *float32 `json:"group_rareness_score"`
	// candidate image thumbnail as base64 string
	CandidateImageThumbnailWebpAsBase64String *string `json:"candidate_image_thumbnail_webp_as_base64_string"`
	// does not impact collection strings
	DoesNotImpactTheFollowingCollectionStrings *string `json:"does_not_impact_the_following_collection_strings"`
	//is invalid sense request
	IsInvalidSenseRequest *bool `json:"is_invalid_sense_request"`
	//invalid sense request reason
	InvalidSenseRequestReason *string `json:"invalid_sense_request_reason"`
	// similarity score to first entry in collection
	SimilarityScoreToFirstEntryInCollection *float32 `json:"similarity_score_to_first_entry_in_collection"`
	// probability of CP
	CpProbability *float32 `json:"cp_probability"`
	// child probability
	ChildProbability *float32 `json:"child_probability"`
	// file path of the image
	ImageFilePath *string `json:"image_file_path"`
	// internet rareness
	InternetRareness *InternetRareness `json:"internet_rareness"`
	// alternative NSFW scores
	AlternativeNsfwScores *AlternativeNSFWScores `json:"alternative_nsfw_scores"`
	// name of the creator
	CreatorName string `json:"creator_name"`
	// website of creator
	CreatorWebsite string `json:"creator_website"`
	// written statement of creator
	CreatorWrittenStatement string `json:"creator_written_statement"`
	// title of NFT
	NftTitle string `json:"nft_title"`
	// series name of NFT
	NftSeriesName string `json:"nft_series_name"`
	// nft creation video youtube url
	NftCreationVideoYoutubeURL string `json:"nft_creation_video_youtube_url"`
	// keywords for NFT
	NftKeywordSet string `json:"nft_keyword_set"`
	// total copies of NFT
	TotalCopies int `json:"total_copies"`
	// preview hash of NFT
	PreviewHash []byte `json:"preview_hash"`
	// thumbnail1 hash of NFT
	Thumbnail1Hash []byte `json:"thumbnail_1_hash"`
	// thumbnail2 hash of NFT
	Thumbnail2Hash []byte `json:"thumbnail_2_hash"`
	// original file size in bytes
	OriginalFileSizeInBytes int `json:"original_file_size_in_bytes"`
	// type of the file
	FileType string `json:"file_type"`
	// max permitted open NSFW score
	MaxPermittedOpenNsfwScore float64 `json:"max_permitted_open_nsfw_score"`
}

// InternetRareness represents the internet_rareness from DD FP File
type InternetRareness struct {
	// Base64 Compressed JSON Table of Rare On Internet Summary
	RareOnInternetSummaryTableAsJSONCompressedB64 *string `json:"rare_on_internet_summary_table_as_json_compressed_b64"`
	// Base64 Compressed JSON of Rare On Internet Graph
	RareOnInternetGraphJSONCompressedB64 *string `json:"rare_on_internet_graph_json_compressed_b64"`
	// Base64 Compressed Json of Alternative Rare On Internet Dict
	AlternativeRareOnInternetDictAsJSONCompressedB64 *string `json:"alternative_rare_on_internet_dict_as_json_compressed_b64"`
	// Minimum Number of Exact Matches on Page
	MinNumberOfExactMatchesInPage *uint32 `json:"min_number_of_exact_matches_in_page"`
	// Earliest Available Date of Internet Results
	EarliestAvailableDateOfInternetResults *string `json:"earliest_available_date_of_internet_results"`
}

// AlternativeNSFWScores represents the struct from DD & FP file
type AlternativeNSFWScores struct {
	// drawings nsfw score
	Drawings *float32 `json:"drawings"`
	// hentai nsfw score
	Hentai *float32 `json:"hentai"`
	// sexy nsfw score
	Sexy *float32 `json:"sexy"`
	// porn nsfw score
	Porn *float32 `json:"porn"`
	// neutral nsfw score
	Neutral *float32 `json:"neutral"`
}

func toSHChallengeReport(data types.SelfHealingChallengeReports) *metrics.SelfHealingChallengeReports {
	reports := &metrics.SelfHealingChallengeReports{
		Reports: make([]*metrics.SelfHealingChallengeReportKV, 0, len(data)),
	}

	for challengeID, report := range data {
		chlngID := challengeID
		reportKV := &metrics.SelfHealingChallengeReportKV{
			ChallengeID: &chlngID,
			Report:      toSHChallengeReportStruct(report),
		}
		reports.Reports = append(reports.Reports, reportKV)
	}

	return reports
}

func toSHChallengeReportStruct(report types.SelfHealingChallengeReport) *metrics.SelfHealingChallengeReport {
	messages := make([]*metrics.SelfHealingMessageKV, 0, len(report))
	for messageType, msgs := range report {
		msgTyp := messageType
		msgKV := &metrics.SelfHealingMessageKV{
			MessageType: &msgTyp,
			Messages:    toSHMessages(msgs),
		}
		messages = append(messages, msgKV)
	}
	return &metrics.SelfHealingChallengeReport{Messages: messages}
}

func toSHMessages(msgs types.SelfHealingMessages) []*metrics.SelfHealingMessage {
	messages := make([]*metrics.SelfHealingMessage, 0, len(msgs))

	for _, msg := range msgs {
		msgType := msg.MessageType.String()
		message := &metrics.SelfHealingMessage{
			TriggerID:       &msg.TriggerID,
			MessageType:     &msgType,
			Data:            toSHMessageData(msg.SelfHealingMessageData, msg.MessageType),
			SenderID:        &msg.SenderID,
			SenderSignature: msg.SenderSignature,
		}
		messages = append(messages, message)
	}
	return messages
}

func toSHMessageData(data types.SelfHealingMessageData, msgType types.SelfHealingMessageType) *metrics.SelfHealingMessageData {
	ret := &metrics.SelfHealingMessageData{
		ChallengerID: &data.ChallengerID,
		RecipientID:  &data.RecipientID,
	}

	if msgType == types.SelfHealingChallengeMessage {
		ret.Challenge = toSHChallengeData(data.Challenge)
	} else if msgType == types.SelfHealingResponseMessage {
		ret.Response = toSHResponseData(data.Response)
	} else if msgType == types.SelfHealingVerificationMessage {
		ret.Verification = toSHVerificationData(data.Verification)
	}

	return ret
}

func toSHChallengeData(data types.SelfHealingChallengeData) *metrics.SelfHealingChallengeData {
	timestamp := data.Timestamp.Format(time.RFC3339)
	return &metrics.SelfHealingChallengeData{
		Block:            &data.Block,
		Merkelroot:       &data.Merkelroot,
		Timestamp:        &timestamp,
		ChallengeTickets: toChallengeTickets(data.ChallengeTickets),
		NodesOnWatchlist: &data.NodesOnWatchlist,
	}
}

func toChallengeTickets(tickets []types.ChallengeTicket) []*metrics.ChallengeTicket {
	mTickets := make([]*metrics.ChallengeTicket, 0, len(tickets))
	for _, ticket := range tickets {
		tType := ticket.TicketType.String()

		mTicket := &metrics.ChallengeTicket{
			TxID:        &ticket.TxID,
			TicketType:  &tType,
			MissingKeys: ticket.MissingKeys,
			DataHash:    ticket.DataHash,
			Recipient:   &ticket.Recipient,
		}
		mTickets = append(mTickets, mTicket)
	}
	return mTickets
}

func toSHResponseData(data types.SelfHealingResponseData) *metrics.SelfHealingResponseData {
	timestamp := data.Timestamp.Format(time.RFC3339)
	return &metrics.SelfHealingResponseData{
		ChallengeID:     &data.ChallengeID,
		Block:           &data.Block,
		Merkelroot:      &data.Merkelroot,
		Timestamp:       &timestamp,
		RespondedTicket: toRespondedTicket(data.RespondedTicket),
		Verifiers:       data.Verifiers,
	}
}

func toRespondedTicket(ticket types.RespondedTicket) *metrics.RespondedTicket {
	tType := ticket.TicketType.String()

	return &metrics.RespondedTicket{
		TxID:                     &ticket.TxID,
		TicketType:               &tType,
		MissingKeys:              ticket.MissingKeys,
		ReconstructedFileHash:    ticket.ReconstructedFileHash,
		SenseFileIds:             ticket.FileIDs,
		RaptorQSymbols:           ticket.RaptorQSymbols,
		IsReconstructionRequired: &ticket.IsReconstructionRequired,
	}
}

func toSHVerificationData(data types.SelfHealingVerificationData) *metrics.SelfHealingVerificationData {
	timestamp := data.Timestamp.Format(time.RFC3339)
	return &metrics.SelfHealingVerificationData{
		ChallengeID:    &data.ChallengeID,
		Block:          &data.Block,
		Merkelroot:     &data.Merkelroot,
		Timestamp:      &timestamp,
		VerifiedTicket: toVerifiedTicket(data.VerifiedTicket),
		VerifiersData:  data.VerifiersData,
	}
}

func toVerifiedTicket(ticket types.VerifiedTicket) *metrics.VerifiedTicket {
	tType := ticket.TicketType.String()

	return &metrics.VerifiedTicket{
		TxID:                     &ticket.TxID,
		TicketType:               &tType,
		MissingKeys:              ticket.MissingKeys,
		ReconstructedFileHash:    ticket.ReconstructedFileHash,
		IsReconstructionRequired: &ticket.IsReconstructionRequired,
		RaptorQSymbols:           ticket.RaptorQSymbols,
		SenseFileIds:             ticket.FileIDs,
		IsVerified:               &ticket.IsVerified,
		Message:                  &ticket.Message,
	}
}
