package services

import (
	"context"
	"fmt"
	"github.com/pastelnetwork/gonode/common/types"
	"time"

	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/walletnode/api"
	"github.com/pastelnetwork/gonode/walletnode/api/gen/http/storage_challenge/server"
	"github.com/pastelnetwork/gonode/walletnode/api/gen/metrics"
	storagechallenge "github.com/pastelnetwork/gonode/walletnode/api/gen/storage_challenge"
	storagechallengeregister "github.com/pastelnetwork/gonode/walletnode/services/storagechallenge"

	goahttp "goa.design/goa/v3/http"
	"goa.design/goa/v3/security"
)

// StorageChallengeAPIHandler - StorageChallengeAPIHandler service
type StorageChallengeAPIHandler struct {
	*Common
	storageChallengeService *storagechallengeregister.Service
}

// Mount onfigures the mux to serve the OpenAPI enpoints.
func (service *StorageChallengeAPIHandler) Mount(ctx context.Context, mux goahttp.Muxer) goahttp.Server {
	endpoints := storagechallenge.NewEndpoints(service)

	srv := server.New(
		endpoints,
		mux,
		goahttp.RequestDecoder,
		goahttp.ResponseEncoder,
		api.ErrorHandler,
		nil,
	)

	server.Mount(mux, srv)

	for _, m := range srv.Mounts {
		log.WithContext(ctx).Infof("%q mounted on %s %s", m.Method, m.Verb, m.Pattern)
	}
	return srv
}

// APIKeyAuth implements the authorization logic for the APIKey security scheme.
func (service *StorageChallengeAPIHandler) APIKeyAuth(ctx context.Context, _ string, _ *security.APIKeyScheme) (context.Context, error) {
	return ctx, nil
}

// GetDetailedLogs returns the detailed storage-challenge data logs
func (service *StorageChallengeAPIHandler) GetDetailedLogs(ctx context.Context, p *storagechallenge.GetDetailedLogsPayload) (scDetailedLogMessages []*storagechallenge.StorageMessage, err error) {
	if p.ChallengeID == "" {
		return nil, storagechallenge.MakeBadRequest(fmt.Errorf("challenge_id is required"))
	}

	req := storagechallengeregister.SCDetailedLogRequest{
		ChallengeID: p.ChallengeID,
		PastelID:    p.Pid,
		Passphrase:  p.Key,
	}

	scDetailedMessageDataList, err := service.storageChallengeService.GetDetailedMessageDataList(ctx, req)
	if err != nil {
		return nil, metrics.MakeInternalServerError(fmt.Errorf("failed to get storage-challenge detailed logs: %w", err))
	}

	for _, scl := range scDetailedMessageDataList {

		scMsg := &storagechallenge.StorageMessage{
			ChallengeID: scl.ChallengeID,
			MessageType: scl.MessageType.String(),
			SenderID:    scl.Sender,

			ChallengerID: scl.Data.ChallengerID,
			Observers:    scl.Data.Observers,
			RecipientID:  scl.Data.RecipientID,
		}

		if scl.MessageType == types.ChallengeMessageType {
			scMsg.Challenge = &storagechallenge.ChallengeData{
				Timestamp:  scl.Data.Challenge.Timestamp.Format(time.RFC3339),
				FileHash:   scl.Data.Challenge.FileHash,
				StartIndex: scl.Data.Challenge.StartIndex,
				EndIndex:   scl.Data.Challenge.EndIndex,
			}

			if scl.Data.Challenge.Block != 0 {
				scMsg.Challenge.Block = &scl.Data.Challenge.Block
			}

			if scl.Data.Challenge.Merkelroot != "" {
				scMsg.Challenge.Merkelroot = &scl.Data.Challenge.Merkelroot
			}
		}

		if scl.MessageType == types.ResponseMessageType {
			scMsg.Response = &storagechallenge.ResponseData{
				Timestamp: scl.Data.Response.Timestamp.Format(time.RFC3339),
			}

			if scl.Data.Response.Block != 0 {
				scMsg.Response.Block = &scl.Data.Response.Block
			}

			if scl.Data.Response.Merkelroot != "" {
				scMsg.Response.Merkelroot = &scl.Data.Response.Merkelroot
			}

			if scl.Data.Response.Hash != "" {
				scMsg.Response.Hash = &scl.Data.Response.Hash
			}
		}

		if scl.MessageType == types.EvaluationMessageType {
			scMsg.ChallengerEvaluation = &storagechallenge.EvaluationData{
				Timestamp:  scl.Data.ChallengerEvaluation.Timestamp.Format(time.RFC3339),
				Hash:       scl.Data.ChallengerEvaluation.Hash,
				IsVerified: scl.Data.ChallengerEvaluation.IsVerified,
			}

			if scl.Data.ChallengerEvaluation.Block != 0 {
				scMsg.ChallengerEvaluation.Block = &scl.Data.ChallengerEvaluation.Block
			}

			if scl.Data.ChallengerEvaluation.Merkelroot != "" {
				scMsg.ChallengerEvaluation.Merkelroot = &scl.Data.ChallengerEvaluation.Merkelroot
			}
		}

		if scl.MessageType == types.AffirmationMessageType {
			scMsg.ObserverEvaluation = &storagechallenge.ObserverEvaluationData{
				IsChallengeTimestampOK:  scl.Data.ObserverEvaluation.IsChallengeTimestampOK,
				IsProcessTimestampOK:    scl.Data.ObserverEvaluation.IsProcessTimestampOK,
				IsEvaluationTimestampOK: scl.Data.ObserverEvaluation.IsEvaluationTimestampOK,
				IsRecipientSignatureOK:  scl.Data.ObserverEvaluation.IsRecipientSignatureOK,
				IsChallengerSignatureOK: scl.Data.ObserverEvaluation.IsChallengerSignatureOK,
				IsEvaluationResultOK:    scl.Data.ObserverEvaluation.IsEvaluationResultOK,
				TrueHash:                scl.Data.ObserverEvaluation.TrueHash,
				Timestamp:               scl.Data.ObserverEvaluation.Timestamp.Format(time.RFC3339),
			}

			if scl.Data.ObserverEvaluation.Block != 0 {
				scMsg.ObserverEvaluation.Block = &scl.Data.ObserverEvaluation.Block
			}

			if scl.Data.ObserverEvaluation.Merkelroot != "" {
				scMsg.ObserverEvaluation.Merkelroot = &scl.Data.ObserverEvaluation.Merkelroot
			}

			if scl.Data.ObserverEvaluation.Reason != "" {
				scMsg.ObserverEvaluation.Reason = &scl.Data.ObserverEvaluation.Reason
			}
		}

		scDetailedLogMessages = append(scDetailedLogMessages, scMsg)
	}

	return scDetailedLogMessages, nil
}

// GetSummaryStats returns the stats over a specified time range
func (service *StorageChallengeAPIHandler) GetSummaryStats(ctx context.Context, p *storagechallenge.GetSummaryStatsPayload) (*storagechallenge.SummaryStatsResult, error) {
	var from, to *time.Time

	if p.From != nil {
		fromTime, err := time.Parse(time.RFC3339, *p.From)
		if err != nil {
			return nil, storagechallenge.MakeBadRequest(fmt.Errorf("invalid from time format: %w", err))
		}

		from = &fromTime
	}

	if p.To != nil {
		toTime, err := time.Parse(time.RFC3339, *p.To)
		if err != nil {
			return nil, storagechallenge.MakeBadRequest(fmt.Errorf("invalid to time format: %w", err))
		}
		to = &toTime
	}

	req := storagechallengeregister.GetSummaryStats{
		From:       from,
		To:         to,
		PastelID:   p.Pid,
		Passphrase: p.Key,
	}

	res, err := service.storageChallengeService.GetSummaryStats(ctx, req)
	if err != nil {
		return nil, storagechallenge.MakeInternalServerError(fmt.Errorf("failed to get metrics: %w", err))
	}

	return &storagechallenge.SummaryStatsResult{
		ScSummaryStats: &storagechallenge.SCSummaryStats{
			TotalChallengesIssued:                    res.SCSummaryStats.TotalChallenges,
			TotalChallengesProcessedByRecipient:      res.SCSummaryStats.TotalChallengesProcessed,
			TotalChallengesEvaluatedByChallenger:     res.SCSummaryStats.TotalChallengesEvaluatedByChallenger,
			TotalChallengesVerified:                  res.SCSummaryStats.TotalChallengesVerified,
			NoOfSlowResponsesObservedByObservers:     res.SCSummaryStats.SlowResponsesObservedByObservers,
			NoOfInvalidSignaturesObservedByObservers: res.SCSummaryStats.InvalidSignaturesObservedByObservers,
			NoOfInvalidEvaluationObservedByObservers: res.SCSummaryStats.InvalidEvaluationObservedByObservers,
		},
	}, nil
}

// NewStorageChallengeAPIHandler returns the swagger OpenAPI implementation.
func NewStorageChallengeAPIHandler(srvc *storagechallengeregister.Service) *StorageChallengeAPIHandler {
	return &StorageChallengeAPIHandler{
		storageChallengeService: srvc,
	}
}
