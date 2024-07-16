// Code generated by goa v3.15.0, DO NOT EDIT.
//
// cascade HTTP client CLI support package
//
// Command:
// $ goa gen github.com/pastelnetwork/gonode/walletnode/api/design

package client

import (
	"encoding/json"
	"fmt"
	"unicode/utf8"

	cascade "github.com/pastelnetwork/gonode/walletnode/api/gen/cascade"
	goa "goa.design/goa/v3/pkg"
)

// BuildUploadAssetPayload builds the payload for the cascade uploadAsset
// endpoint from CLI flags.
func BuildUploadAssetPayload(cascadeUploadAssetBody string) (*cascade.UploadAssetPayload, error) {
	var err error
	var body UploadAssetRequestBody
	{
		err = json.Unmarshal([]byte(cascadeUploadAssetBody), &body)
		if err != nil {
			return nil, fmt.Errorf("invalid JSON for body, \nerror: %s, \nexample of valid JSON:\n%s", err, "'{\n      \"file\": \"TmFtIGEu\"\n   }'")
		}
		if body.Bytes == nil {
			err = goa.MergeErrors(err, goa.MissingFieldError("file", "body"))
		}
		if err != nil {
			return nil, err
		}
	}
	v := &cascade.UploadAssetPayload{
		Bytes:    body.Bytes,
		Filename: body.Filename,
		Hash:     body.Hash,
		Size:     body.Size,
	}

	return v, nil
}

// BuildUploadAssetV2Payload builds the payload for the cascade uploadAssetV2
// endpoint from CLI flags.
func BuildUploadAssetV2Payload(cascadeUploadAssetV2Body string) (*cascade.UploadAssetV2Payload, error) {
	var err error
	var body UploadAssetV2RequestBody
	{
		err = json.Unmarshal([]byte(cascadeUploadAssetV2Body), &body)
		if err != nil {
			return nil, fmt.Errorf("invalid JSON for body, \nerror: %s, \nexample of valid JSON:\n%s", err, "'{\n      \"file\": \"QSBldmVuaWV0IG51bXF1YW0gdmVsLg==\"\n   }'")
		}
		if body.Bytes == nil {
			err = goa.MergeErrors(err, goa.MissingFieldError("file", "body"))
		}
		if err != nil {
			return nil, err
		}
	}
	v := &cascade.UploadAssetV2Payload{
		Bytes:    body.Bytes,
		Filename: body.Filename,
		Hash:     body.Hash,
		Size:     body.Size,
	}

	return v, nil
}

// BuildStartProcessingPayload builds the payload for the cascade
// startProcessing endpoint from CLI flags.
func BuildStartProcessingPayload(cascadeStartProcessingBody string, cascadeStartProcessingFileID string, cascadeStartProcessingKey string) (*cascade.StartProcessingPayload, error) {
	var err error
	var body StartProcessingRequestBody
	{
		err = json.Unmarshal([]byte(cascadeStartProcessingBody), &body)
		if err != nil {
			return nil, fmt.Errorf("invalid JSON for body, \nerror: %s, \nexample of valid JSON:\n%s", err, "'{\n      \"app_pastelid\": \"jXYJud3rmrR1Sk2scvR47N4E4J5Vv48uCC6se2nzHrBRdjaKj3ybPoi1Y2VVoRqi1GnQrYKjSxQAC7NBtvtEdS\",\n      \"burn_txid\": \"576e7b824634a488a2f0baacf5a53b237d883029f205df25b300b87c8877ab58\",\n      \"burn_txids\": [\n         \"576e7b824634a488a2f0baacf5a53b237d883029f205df25b300b87c8877ab58\"\n      ],\n      \"make_publicly_accessible\": false,\n      \"spendable_address\": \"PtiqRXn2VQwBjp1K8QXR2uW2w2oZ3Ns7N6j\"\n   }'")
		}
		if body.BurnTxid != nil {
			if utf8.RuneCountInString(*body.BurnTxid) < 64 {
				err = goa.MergeErrors(err, goa.InvalidLengthError("body.burn_txid", *body.BurnTxid, utf8.RuneCountInString(*body.BurnTxid), 64, true))
			}
		}
		if body.BurnTxid != nil {
			if utf8.RuneCountInString(*body.BurnTxid) > 64 {
				err = goa.MergeErrors(err, goa.InvalidLengthError("body.burn_txid", *body.BurnTxid, utf8.RuneCountInString(*body.BurnTxid), 64, false))
			}
		}
		err = goa.MergeErrors(err, goa.ValidatePattern("body.app_pastelid", body.AppPastelID, "^[a-zA-Z0-9]+$"))
		if utf8.RuneCountInString(body.AppPastelID) < 86 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("body.app_pastelid", body.AppPastelID, utf8.RuneCountInString(body.AppPastelID), 86, true))
		}
		if utf8.RuneCountInString(body.AppPastelID) > 86 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("body.app_pastelid", body.AppPastelID, utf8.RuneCountInString(body.AppPastelID), 86, false))
		}
		if body.SpendableAddress != nil {
			err = goa.MergeErrors(err, goa.ValidatePattern("body.spendable_address", *body.SpendableAddress, "^[a-zA-Z0-9]+$"))
		}
		if body.SpendableAddress != nil {
			if utf8.RuneCountInString(*body.SpendableAddress) < 35 {
				err = goa.MergeErrors(err, goa.InvalidLengthError("body.spendable_address", *body.SpendableAddress, utf8.RuneCountInString(*body.SpendableAddress), 35, true))
			}
		}
		if body.SpendableAddress != nil {
			if utf8.RuneCountInString(*body.SpendableAddress) > 35 {
				err = goa.MergeErrors(err, goa.InvalidLengthError("body.spendable_address", *body.SpendableAddress, utf8.RuneCountInString(*body.SpendableAddress), 35, false))
			}
		}
		if err != nil {
			return nil, err
		}
	}
	var fileID string
	{
		fileID = cascadeStartProcessingFileID
		if utf8.RuneCountInString(fileID) < 8 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("file_id", fileID, utf8.RuneCountInString(fileID), 8, true))
		}
		if utf8.RuneCountInString(fileID) > 8 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("file_id", fileID, utf8.RuneCountInString(fileID), 8, false))
		}
		if err != nil {
			return nil, err
		}
	}
	var key string
	{
		key = cascadeStartProcessingKey
	}
	v := &cascade.StartProcessingPayload{
		BurnTxid:               body.BurnTxid,
		AppPastelID:            body.AppPastelID,
		MakePubliclyAccessible: body.MakePubliclyAccessible,
		SpendableAddress:       body.SpendableAddress,
	}
	if body.BurnTxids != nil {
		v.BurnTxids = make([]string, len(body.BurnTxids))
		for i, val := range body.BurnTxids {
			v.BurnTxids[i] = val
		}
	}
	{
		var zero bool
		if v.MakePubliclyAccessible == zero {
			v.MakePubliclyAccessible = false
		}
	}
	v.FileID = fileID
	v.Key = key

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
			err = goa.MergeErrors(err, goa.InvalidLengthError("taskId", taskID, utf8.RuneCountInString(taskID), 8, true))
		}
		if utf8.RuneCountInString(taskID) > 8 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("taskId", taskID, utf8.RuneCountInString(taskID), 8, false))
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
			err = goa.MergeErrors(err, goa.InvalidLengthError("taskId", taskID, utf8.RuneCountInString(taskID), 8, true))
		}
		if utf8.RuneCountInString(taskID) > 8 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("taskId", taskID, utf8.RuneCountInString(taskID), 8, false))
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

// BuildDownloadV2Payload builds the payload for the cascade downloadV2
// endpoint from CLI flags.
func BuildDownloadV2Payload(cascadeDownloadV2Txid string, cascadeDownloadV2Pid string, cascadeDownloadV2Key string) (*cascade.DownloadPayload, error) {
	var err error
	var txid string
	{
		txid = cascadeDownloadV2Txid
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
		pid = cascadeDownloadV2Pid
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
		key = cascadeDownloadV2Key
	}
	v := &cascade.DownloadPayload{}
	v.Txid = txid
	v.Pid = pid
	v.Key = key

	return v, nil
}

// BuildGetDownloadTaskStatePayload builds the payload for the cascade
// getDownloadTaskState endpoint from CLI flags.
func BuildGetDownloadTaskStatePayload(cascadeGetDownloadTaskStateFileID string) (*cascade.GetDownloadTaskStatePayload, error) {
	var err error
	var fileID string
	{
		fileID = cascadeGetDownloadTaskStateFileID
		if utf8.RuneCountInString(fileID) < 6 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("file_id", fileID, utf8.RuneCountInString(fileID), 6, true))
		}
		if utf8.RuneCountInString(fileID) > 6 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("file_id", fileID, utf8.RuneCountInString(fileID), 6, false))
		}
		if err != nil {
			return nil, err
		}
	}
	v := &cascade.GetDownloadTaskStatePayload{}
	v.FileID = fileID

	return v, nil
}

// BuildRegistrationDetailsPayload builds the payload for the cascade
// registrationDetails endpoint from CLI flags.
func BuildRegistrationDetailsPayload(cascadeRegistrationDetailsBaseFileID string) (*cascade.RegistrationDetailsPayload, error) {
	var err error
	var baseFileID string
	{
		baseFileID = cascadeRegistrationDetailsBaseFileID
		if utf8.RuneCountInString(baseFileID) > 8 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("base_file_id", baseFileID, utf8.RuneCountInString(baseFileID), 8, false))
		}
		if err != nil {
			return nil, err
		}
	}
	v := &cascade.RegistrationDetailsPayload{}
	v.BaseFileID = baseFileID

	return v, nil
}

// BuildRestorePayload builds the payload for the cascade restore endpoint from
// CLI flags.
func BuildRestorePayload(cascadeRestoreBody string, cascadeRestoreBaseFileID string, cascadeRestoreKey string) (*cascade.RestorePayload, error) {
	var err error
	var body RestoreRequestBody
	{
		err = json.Unmarshal([]byte(cascadeRestoreBody), &body)
		if err != nil {
			return nil, fmt.Errorf("invalid JSON for body, \nerror: %s, \nexample of valid JSON:\n%s", err, "'{\n      \"app_pastelId\": \"jXYJud3rmrR1Sk2scvR47N4E4J5Vv48uCC6se2nzHrBRdjaKj3ybPoi1Y2VVoRqi1GnQrYKjSxQAC7NBtvtEdS\",\n      \"make_publicly_accessible\": false,\n      \"spendable_address\": \"PtiqRXn2VQwBjp1K8QXR2uW2w2oZ3Ns7N6j\"\n   }'")
		}
		err = goa.MergeErrors(err, goa.ValidatePattern("body.app_pastelId", body.AppPastelID, "^[a-zA-Z0-9]+$"))
		if utf8.RuneCountInString(body.AppPastelID) < 86 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("body.app_pastelId", body.AppPastelID, utf8.RuneCountInString(body.AppPastelID), 86, true))
		}
		if utf8.RuneCountInString(body.AppPastelID) > 86 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("body.app_pastelId", body.AppPastelID, utf8.RuneCountInString(body.AppPastelID), 86, false))
		}
		if body.SpendableAddress != nil {
			err = goa.MergeErrors(err, goa.ValidatePattern("body.spendable_address", *body.SpendableAddress, "^[a-zA-Z0-9]+$"))
		}
		if body.SpendableAddress != nil {
			if utf8.RuneCountInString(*body.SpendableAddress) < 35 {
				err = goa.MergeErrors(err, goa.InvalidLengthError("body.spendable_address", *body.SpendableAddress, utf8.RuneCountInString(*body.SpendableAddress), 35, true))
			}
		}
		if body.SpendableAddress != nil {
			if utf8.RuneCountInString(*body.SpendableAddress) > 35 {
				err = goa.MergeErrors(err, goa.InvalidLengthError("body.spendable_address", *body.SpendableAddress, utf8.RuneCountInString(*body.SpendableAddress), 35, false))
			}
		}
		if err != nil {
			return nil, err
		}
	}
	var baseFileID string
	{
		baseFileID = cascadeRestoreBaseFileID
		if utf8.RuneCountInString(baseFileID) > 8 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("base_file_id", baseFileID, utf8.RuneCountInString(baseFileID), 8, false))
		}
		if err != nil {
			return nil, err
		}
	}
	var key string
	{
		key = cascadeRestoreKey
	}
	v := &cascade.RestorePayload{
		AppPastelID:            body.AppPastelID,
		MakePubliclyAccessible: body.MakePubliclyAccessible,
		SpendableAddress:       body.SpendableAddress,
	}
	{
		var zero bool
		if v.MakePubliclyAccessible == zero {
			v.MakePubliclyAccessible = false
		}
	}
	v.BaseFileID = baseFileID
	v.Key = key

	return v, nil
}
