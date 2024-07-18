package cascaderegister

import (
	"context"
	"database/sql"
	"encoding/base64"
	"time"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/walletnode/api/gen/cascade"
)

const (
	defaultRestoreVolumeInterval = 3 * time.Minute
)

// RestoreService restore the missed volume file
type RestoreService struct{}

// Run runs the restore service.
func (service *RestoreService) Run(ctx context.Context, cascadeRegistrationService CascadeRegistrationService) error {
	if r := recover(); r != nil {
		log.Errorf("Recovered from panic in restore run: %v", r)
	}

	for {
		select {
		case <-ctx.Done():
			log.Info("context canceled, returning restore service")
			return nil
		case <-time.After(defaultRestoreVolumeInterval):
			restoringFilesRes, err := service.restore(ctx, cascadeRegistrationService)
			if err != nil {
				log.WithContext(ctx).WithError(err).Error("Error restoring files")
			}
			log.WithContext(ctx).WithField("restoring_files_res", restoringFilesRes).Debug("successfully called restored files function")
		}
	}
}

func (service *RestoreService) restore(ctx context.Context, cascadeRegistrationService CascadeRegistrationService) ([]*cascade.RestoreFile, error) {
	log.WithContext(ctx).Debug("Restore worker started")

	ucFiles, err := cascadeRegistrationService.GetUnCompletedFiles()
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		log.WithContext(ctx).WithError(err).Error("Error in retrieving un-concluded files")
		return nil, err
	}
	if len(ucFiles) == 0 {
		log.WithContext(ctx).Infof("no un-concluded files found")
		return nil, nil
	}

	runningTasks := cascadeRegistrationService.Worker.Tasks()
	runningTaskIDs := make(map[string]struct{}, len(runningTasks))
	for _, rt := range runningTasks {
		runningTaskIDs[rt.ID()] = struct{}{}
	}

	var restoreFiles []*cascade.RestoreFile
	for _, ucFile := range ucFiles {
		if ucFile.TaskID == "" || ucFile.PastelID == "" || ucFile.Passphrase == "" {
			log.WithContext(ctx).WithField("file_id", ucFile.FileID).
				Debug("un-concluded volume without task-id, pastel-id & passphrase, cannot proceed with recovery")
			continue
		}

		if _, exists := runningTaskIDs[ucFile.TaskID]; exists {
			log.WithContext(ctx).WithField("task_id", ucFile.TaskID).WithField("base_file_id", ucFile.BaseFileID).
				Debug("current task is already in-progress can't execute the recovery-flow")
			continue
		}

		decodedPassphraseBytes, err := base64.StdEncoding.DecodeString(ucFile.Passphrase)
		if err != nil {
			log.WithContext(ctx).WithField("task_id", ucFile.TaskID).WithField("base_file_id", ucFile.BaseFileID).
				Error("error decoding passphrase")
			continue
		}

		restoreFileRes, err := cascadeRegistrationService.RestoreFile(ctx, &cascade.RestorePayload{
			BaseFileID:             ucFile.BaseFileID,
			AppPastelID:            ucFile.PastelID,
			Key:                    string(decodedPassphraseBytes),
			MakePubliclyAccessible: false,
		})
		if err != nil {
			log.WithContext(ctx).WithError(err).WithField("base_file_id", ucFile.BaseFileID).
				WithField("file_id", ucFile.FileID).Error("Error restoring file")
			continue
		}
		restoreFiles = append(restoreFiles, restoreFileRes)
	}

	log.WithContext(ctx).Debug("Restore worker finished")
	return restoreFiles, nil
}

// NewRestoreService returns a new RestoreService instance.
func NewRestoreService() *RestoreService {
	return &RestoreService{}
}
