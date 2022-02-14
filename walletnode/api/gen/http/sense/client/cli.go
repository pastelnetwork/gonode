// Code generated by goa v3.4.3, DO NOT EDIT.
//
// sense HTTP client CLI support package
//
// Command:
// $ goa gen github.com/pastelnetwork/gonode/walletnode/api/design

package client

import (
	"encoding/json"
	"fmt"
	"unicode/utf8"

	sense "github.com/pastelnetwork/gonode/walletnode/api/gen/sense"
	goa "goa.design/goa/v3/pkg"
)

// BuildUploadImagePayload builds the payload for the sense uploadImage
// endpoint from CLI flags.
func BuildUploadImagePayload(senseUploadImageBody string) (*sense.UploadImagePayload, error) {
	var err error
	var body UploadImageRequestBody
	{
		err = json.Unmarshal([]byte(senseUploadImageBody), &body)
		if err != nil {
			return nil, fmt.Errorf("invalid JSON for body, \nerror: %s, \nexample of valid JSON:\n%s", err, "'{\n      \"file\": \"Vm9sdXB0YXMgdm9sdXB0YXRlcyB0b3RhbSBhdXQgc2VkIHF1byBjb25zZXF1YXR1ci4=\"\n   }'")
		}
		if body.Bytes == nil {
			err = goa.MergeErrors(err, goa.MissingFieldError("file", "body"))
		}
		if err != nil {
			return nil, err
		}
	}
	v := &sense.UploadImagePayload{
		Bytes:    body.Bytes,
		Filename: body.Filename,
	}

	return v, nil
}

// BuildActionDetailsPayload builds the payload for the sense actionDetails
// endpoint from CLI flags.
func BuildActionDetailsPayload(senseActionDetailsBody string, senseActionDetailsImageID string) (*sense.ActionDetailsPayload, error) {
	var err error
	var body ActionDetailsRequestBody
	{
		err = json.Unmarshal([]byte(senseActionDetailsBody), &body)
		if err != nil {
			return nil, fmt.Errorf("invalid JSON for body, \nerror: %s, \nexample of valid JSON:\n%s", err, "'{\n      \"action_data_hash\": \"7ae3874ff2df92df38cce7586c08fe8f3687884edf3b0543f8d9420f4df31265\",\n      \"action_data_signature\": \"bTwvO6UZSvFqHb9qmXbAOg2VmupmP70wfhYsAvwFfeC61cuoL9KIXZtbdQ/Ek8FVNoTpCY5BuxcA6lNjkIOBh4w9/RWtuqF16IaJhAnZ4JbZm1MiDCGcf7x0UU/GNRSk6rNAHlPYkOPdhkha+JjCwD4A\",\n      \"app_pastelid\": \"jXZqaS48TT6LFjxnf9P68hZNQVCFqY631FPz4CtM6VugDi5yLNB51ccy17kgKKCyFuL4qadQsXHH3QwHVuPVyY\"\n   }'")
		}
		err = goa.MergeErrors(err, goa.ValidatePattern("body.app_pastelid", body.PastelID, "^[a-zA-Z0-9]"))
		if utf8.RuneCountInString(body.PastelID) < 86 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("body.app_pastelid", body.PastelID, utf8.RuneCountInString(body.PastelID), 86, true))
		}
		if utf8.RuneCountInString(body.PastelID) > 86 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("body.app_pastelid", body.PastelID, utf8.RuneCountInString(body.PastelID), 86, false))
		}
		err = goa.MergeErrors(err, goa.ValidatePattern("body.action_data_hash", body.ActionDataHash, "^[a-fA-F0-9]"))
		if utf8.RuneCountInString(body.ActionDataHash) < 64 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("body.action_data_hash", body.ActionDataHash, utf8.RuneCountInString(body.ActionDataHash), 64, true))
		}
		if utf8.RuneCountInString(body.ActionDataHash) > 64 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("body.action_data_hash", body.ActionDataHash, utf8.RuneCountInString(body.ActionDataHash), 64, false))
		}
		err = goa.MergeErrors(err, goa.ValidatePattern("body.action_data_signature", body.ActionDataSignature, "^[a-zA-Z0-9\\/+]"))
		if utf8.RuneCountInString(body.ActionDataSignature) < 152 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("body.action_data_signature", body.ActionDataSignature, utf8.RuneCountInString(body.ActionDataSignature), 152, true))
		}
		if utf8.RuneCountInString(body.ActionDataSignature) > 152 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("body.action_data_signature", body.ActionDataSignature, utf8.RuneCountInString(body.ActionDataSignature), 152, false))
		}
		if err != nil {
			return nil, err
		}
	}
	var imageID string
	{
		imageID = senseActionDetailsImageID
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
	v := &sense.ActionDetailsPayload{
		PastelID:            body.PastelID,
		ActionDataHash:      body.ActionDataHash,
		ActionDataSignature: body.ActionDataSignature,
	}
	v.ImageID = imageID

	return v, nil
}

// BuildStartProcessingPayload builds the payload for the sense startProcessing
// endpoint from CLI flags.
func BuildStartProcessingPayload(senseStartProcessingBody string, senseStartProcessingImageID string, senseStartProcessingAppPastelidPassphrase string) (*sense.StartProcessingPayload, error) {
	var err error
	var body StartProcessingRequestBody
	{
		err = json.Unmarshal([]byte(senseStartProcessingBody), &body)
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
		imageID = senseStartProcessingImageID
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
		appPastelidPassphrase = senseStartProcessingAppPastelidPassphrase
	}
	v := &sense.StartProcessingPayload{
		BurnTxid:    body.BurnTxid,
		AppPastelID: body.AppPastelID,
	}
	v.ImageID = imageID
	v.AppPastelidPassphrase = appPastelidPassphrase

	return v, nil
}

// BuildRegisterTaskStatePayload builds the payload for the sense
// registerTaskState endpoint from CLI flags.
func BuildRegisterTaskStatePayload(senseRegisterTaskStateTaskID string) (*sense.RegisterTaskStatePayload, error) {
	var err error
	var taskID string
	{
		taskID = senseRegisterTaskStateTaskID
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
	v := &sense.RegisterTaskStatePayload{}
	v.TaskID = taskID

	return v, nil
}
