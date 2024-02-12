package selfhealing

import (
	"context"
	json "github.com/json-iterator/go"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/types"
	"github.com/pastelnetwork/gonode/common/utils"
)

// AcknowledgeAndStoreChallengeTickets acknowledges the trigger event and save the challenge tickets in DB
func (task *SHTask) AcknowledgeAndStoreChallengeTickets(ctx context.Context, incomingChallengeMessage types.SelfHealingMessage) error {
	logger := log.WithContext(ctx).WithField("trigger_id", incomingChallengeMessage.TriggerID)

	if task.triggerAcknowledgementMap[incomingChallengeMessage.TriggerID] {
		return nil
	}
	task.triggerAcknowledgementMap[incomingChallengeMessage.TriggerID] = true

	logger.Info("AcknowledgeTriggerEventWorker has been invoked")
	// incoming challenge message validation
	if err := task.validateSelfHealingChallengeIncomingData(ctx, incomingChallengeMessage); err != nil {
		log.WithContext(ctx).WithError(err).Error("Error validating self-healing challenge incoming data: ")
		return err
	}
	logger.Info("self healing challenge message has been validated")

	logger.Info("total challenge tickets received", len(incomingChallengeMessage.SelfHealingMessageData.Challenge.ChallengeTickets))

	batches := splitSelfHealingTickets(incomingChallengeMessage.SelfHealingMessageData.Challenge.ChallengeTickets, 50)

	logger.WithField("total_batches", len(batches)).Info("batches have been created")

	for i := 0; i < len(batches); i++ {
		err := task.StoreSelfHealingChallengeEvents(ctx, incomingChallengeMessage.TriggerID, incomingChallengeMessage.SenderID, incomingChallengeMessage.SenderSignature, batches[i])
		if err != nil {
			logger.WithError(err).Error("Error inserting batch of challenge tickets")
			return err
		}
	}

	logger.Info("all challenge tickets have been stored")

	return nil
}

func splitSelfHealingTickets(tickets []types.ChallengeTicket, batchSize int) [][]types.ChallengeTicket {
	var batches [][]types.ChallengeTicket

	numBatches := (len(tickets) + batchSize - 1) / batchSize // calculate the total number of batches
	batches = make([][]types.ChallengeTicket, 0, numBatches) // pre-allocate memory for the batches

	for batchSize < len(tickets) {
		tickets, batches = tickets[batchSize:], append(batches, tickets[0:batchSize:batchSize])
	}
	batches = append(batches, tickets)

	return batches
}

// StoreSelfHealingChallengeEvents stores the self-healing challenge tickets as challenge events to db for further processing
func (task *SHTask) StoreSelfHealingChallengeEvents(ctx context.Context, triggerID string, senderID string, sig []byte, tickets []types.ChallengeTicket) error {
	var challengeEvents []types.SelfHealingChallengeEvent

	for i := 0; i < len(tickets); i++ {
		ticket := tickets[i]
		challengeID := utils.GetHashFromString(ticket.TxID + triggerID)

		ticketBytes, err := json.Marshal(ticket)
		if err != nil {
			log.WithContext(ctx).WithField("ticket_txid", tickets[i].TxID).Error("error converting ticket to bytes")
			continue
		}

		execMetricMsg := []types.SelfHealingMessage{
			{
				TriggerID:              triggerID,
				MessageType:            types.SelfHealingAcknowledgementMessage,
				SenderID:               task.nodeID,
				SelfHealingMessageData: types.SelfHealingMessageData{},
				SenderSignature:        sig,
			},
		}

		execMetricBytes, err := json.Marshal(execMetricMsg)
		if err != nil {
			log.WithContext(ctx).WithField("ticket_txid", tickets[i].TxID).Error("error converting exec metric to bytes")
			continue
		}

		executionMetric := types.SelfHealingExecutionMetric{
			TriggerID:       triggerID,
			ChallengeID:     challengeID,
			MessageType:     int(types.SelfHealingAcknowledgementMessage),
			Data:            execMetricBytes,
			SenderID:        senderID,
			SenderSignature: sig,
		}

		challengeEvents = append(challengeEvents, types.SelfHealingChallengeEvent{
			TriggerID:   triggerID,
			TicketID:    ticket.TxID,
			ChallengeID: challengeID,
			Data:        ticketBytes,
			SenderID:    senderID,
			ExecMetric:  executionMetric,
		})
	}

	if task.historyDB != nil {
		err := task.historyDB.BatchInsertSelfHealingChallengeEvents(ctx, challengeEvents)
		if err != nil {
			log.WithContext(ctx).WithError(err).Error("error storing self-healing tickets batch to DB")
			return err
		}
	}

	return nil
}
