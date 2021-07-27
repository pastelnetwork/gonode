

package userdata

// List common code 
const (
	SuccessProcess = iota

	// Error for Secondary Supernode verify process
	ErrorNotEnoughSupernode
	ErrorNotEnoughSupernodeResponse
	ErrorNotEnoughSupernodeConfirm
	ErrorSupernodeVerifyFail

	// Validation userdata
	SuccessValidateContent
	ErrorOnContent

	// Status for Primary Supernode verify process
	SuccessAddDataToPrimarySupernode
	ErrorPrimarySupernodeFailToProcess
	ErrorSignatureMismatchBetweenSupernode
	ErrorUserdataMismatchBetweenSupernode
	SuccessVerifyAllSignature
	ErrorPrimarySupernodeVerifyFail

	// Status for actual write to rqlite db
	SuccessWriteToRQLiteDB
	ErrorWriteToRQLiteDBFail

)

// ResponseCode for Supernode reponse

// Description of ResponseCode
var Description = map[int]string{
	SuccessProcess:               			"User specified data set successfully",
	ErrorNotEnoughSupernode:         		"The process verify signature in primary Supernode success",
	ErrorNotEnoughSupernodeResponse:		"Not enough SuperNodes reply",
	ErrorNotEnoughSupernodeConfirm:			"SuperNodes response mismatch",
	
	SuccessValidateContent:					"Content Validation successfully",
	ErrorOnContent:            				"There is error on field(s) in user specified data",
	
	SuccessAddDataToPrimarySupernode:		"Primary Supernode success to process data signed send by current supernode",
	ErrorPrimarySupernodeFailToProcess:		"Don't have response from Primary SuperNode or response content is empty",
	ErrorSignatureMismatchBetweenSupernode: "There is not enough number of valid signatures from Supernodes",
	ErrorUserdataMismatchBetweenSupernode:  "There is mismatch of Userdata between Supernodes",
	
	SuccessVerifyAllSignature:				"User specified data set successfully",
	ErrorPrimarySupernodeVerifyFail:		"The process verify signature in primary Supernode fail",

	SuccessWriteToRQLiteDB:					"Data is written to rqlite db successfully",
	ErrorWriteToRQLiteDBFail:				"Data fail to write to rqlite db",
}