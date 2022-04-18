// Code generated by goa v3.6.2, DO NOT EDIT.
//
// cascade HTTP client CLI support package
//
// Command:
// $ goa gen github.com/pastelnetwork/gonode/walletnode/api/design -o api/

package client

import (
	"encoding/json"
	"fmt"
	"unicode/utf8"

	cascade "github.com/pastelnetwork/gonode/walletnode/api/gen/cascade"
	goa "goa.design/goa/v3/pkg"
)

// BuildUploadImagePayload builds the payload for the cascade uploadImage
// endpoint from CLI flags.
func BuildUploadImagePayload(cascadeUploadImageBody string) (*cascade.UploadImagePayload, error) {
	var err error
	var body UploadImageRequestBody
	{
		err = json.Unmarshal([]byte(cascadeUploadImageBody), &body)
		if err != nil {
			return nil, fmt.Errorf("invalid JSON for body, \nerror: %s, \nexample of valid JSON:\n%s", err, "'{\n      \"file\": \"VGVtcG9yYSBub24gY3VscGEgdm9sdXB0YXRpYnVzIHJlcnVtIGVzdC4=\"\n   }'")
		}
		if body.Bytes == nil {
			err = goa.MergeErrors(err, goa.MissingFieldError("file", "body"))
		}
		if err != nil {
			return nil, err
		}
	}
	v := &cascade.UploadImagePayload{
		Bytes:    body.Bytes,
		Filename: body.Filename,
	}

	return v, nil
}

// BuildStartProcessingPayload builds the payload for the cascade
// startProcessing endpoint from CLI flags.
func BuildStartProcessingPayload(cascadeStartProcessingBody string, cascadeStartProcessingImageID string, cascadeStartProcessingAppPastelidPassphrase string) (*cascade.StartProcessingPayload, error) {
	var err error
	var body StartProcessingRequestBody
	{
		err = json.Unmarshal([]byte(cascadeStartProcessingBody), &body)
		if err != nil {
			return nil, fmt.Errorf("invalid JSON for body, \nerror: %s, \nexample of valid JSON:\n%s", err, "'{\n      \"app_pastelid\": \"jXYJud3rmrR1Sk2scvR47N4E4J5Vv48uCC6se2nzHrBRdjaKj3ybPoi1Y2VVoRqi1GnQrYKjSxQAC7NBtvtEdS\",\n      \"burn_txid\": \"576e7b824634a488a2f0baacf5a53b237d883029f205df25b300b87c8877ab58\"\n   }'")
		}
		if utf8.RuneCountInString(body.BurnTxid) < 64 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("body.burn_txid", body.BurnTxid, utf8.RuneCountInString(body.BurnTxid), 64, true))
		}
		if utf8.RuneCountInString(body.BurnTxid) > 64 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("body.burn_txid", body.BurnTxid, utf8.RuneCountInString(body.BurnTxid), 64, false))
		}
		err = goa.MergeErrors(err, goa.ValidatePattern("body.app_pastelid", body.AppPastelID, "^[a-zA-Z0-9]+$"))
		if utf8.RuneCountInString(body.AppPastelID) < 86 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("body.app_pastelid", body.AppPastelID, utf8.RuneCountInString(body.AppPastelID), 86, true))
		}
		if utf8.RuneCountInString(body.AppPastelID) > 86 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("body.app_pastelid", body.AppPastelID, utf8.RuneCountInString(body.AppPastelID), 86, false))
		}
		if err != nil {
			return nil, err
		}
	}
	var imageID string
	{
		imageID = cascadeStartProcessingImageID
		if utf8.RuneCountInString(imageID) < 8 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("imageID", imageID, utf8.RuneCountInString(imageID), 8, true))
		}
		if utf8.RuneCountInString(imageID) > 8 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("imageID", imageID, utf8.RuneCountInString(imageID), 8, false))
		}
		if err != nil {
			return nil, err
		}
	}
	var appPastelidPassphrase string
	{
		appPastelidPassphrase = cascadeStartProcessingAppPastelidPassphrase
	}
	v := &cascade.StartProcessingPayload{
		BurnTxid:    body.BurnTxid,
		AppPastelID: body.AppPastelID,
	}
	v.ImageID = imageID
	v.AppPastelidPassphrase = appPastelidPassphrase

	return v, nil
}

// BuildRegisterTaskStatePayload builds the payload for the cascade
// registerTaskState endpoint from CLI flags.
func BuildRegisterTaskStatePayload(cascadeRegisterTaskStateTaskID string) (*cascade.RegisterTaskStatePayload, error) {
	var err error
	var taskID string
	{
		taskID = cascadeRegisterTaskStateTaskID
		if utf8.RuneCountInString(taskID) < 8 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("taskID", taskID, utf8.RuneCountInString(taskID), 8, true))
		}
		if utf8.RuneCountInString(taskID) > 8 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("taskID", taskID, utf8.RuneCountInString(taskID), 8, false))
		}
		if err != nil {
			return nil, err
		}
	}
	v := &cascade.RegisterTaskStatePayload{}
	v.TaskID = taskID

	return v, nil
}

// BuildGetTaskHistoryPayload builds the payload for the cascade getTaskHistory
// endpoint from CLI flags.
func BuildGetTaskHistoryPayload(cascadeGetTaskHistoryTaskID string) (*cascade.GetTaskHistoryPayload, error) {
	var err error
	var taskID string
	{
		taskID = cascadeGetTaskHistoryTaskID
		if utf8.RuneCountInString(taskID) < 8 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("taskID", taskID, utf8.RuneCountInString(taskID), 8, true))
		}
		if utf8.RuneCountInString(taskID) > 8 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("taskID", taskID, utf8.RuneCountInString(taskID), 8, false))
		}
		if err != nil {
			return nil, err
		}
	}
	v := &cascade.GetTaskHistoryPayload{}
	v.TaskID = taskID

	return v, nil
}

// BuildDownloadPayload builds the payload for the cascade download endpoint
// from CLI flags.
func BuildDownloadPayload(cascadeDownloadTxid string, cascadeDownloadPid string, cascadeDownloadKey string) (*cascade.DownloadPayload, error) {
	var err error
	var txid string
	{
		txid = cascadeDownloadTxid
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
		pid = cascadeDownloadPid
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
		key = cascadeDownloadKey
	}
	v := &cascade.DownloadPayload{}
	v.Txid = txid
	v.Pid = pid
	v.Key = key

	return v, nil
}
