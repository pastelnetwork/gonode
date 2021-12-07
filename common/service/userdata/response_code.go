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
	ErrorVerifyUserdataFail

	// Status for Primary Supernode verify process
	SuccessAddDataToPrimarySupernode
	ErrorPrimarySupernodeFailToProcess
	ErrorSignatureMismatchBetweenSupernode
	ErrorUserdataMismatchBetweenSupernode
	SuccessVerifyAllSignature
	ErrorPrimarySupernodeVerifyFail

	// Status for actual write to rqlite db
	SuccessWriteToRQLiteDB
	ErrorRQLiteDBNotFound
	ErrorWriteToRQLiteDBFail
)

// Description of ResponseCode
var Description = map[int]string{
	SuccessProcess:                  "User specified data set successfully",
	ErrorNotEnoughSupernode:         "The process verify signature in primary Supernode was successful",
	ErrorNotEnoughSupernodeResponse: "Not enough SuperNodes replied",
	ErrorNotEnoughSupernodeConfirm:  "Not enough SuperNodes confirmed",

	SuccessValidateContent:  "Content Validation successfully",
	ErrorOnContent:          "There is error on field(s) in user specified data",
	ErrorVerifyUserdataFail: "Supernode failed to verify Walletnode's signature",

	SuccessAddDataToPrimarySupernode:       "Primary Supernode successfully processed data signed and sent by current supernode",
	ErrorPrimarySupernodeFailToProcess:     "Don't have response from Primary SuperNode or response content is empty",
	ErrorSignatureMismatchBetweenSupernode: "There are not enough valid signatures from Supernodes",
	ErrorUserdataMismatchBetweenSupernode:  "There is a mismatch of Userdata between Supernodes",

	SuccessVerifyAllSignature:       "The process of verifying the signature in primary Supernode was successful",
	ErrorPrimarySupernodeVerifyFail: "The process of verifying the signature in primary Supernode failed",

	SuccessWriteToRQLiteDB:   "Data was written to rqlite db successfully",
	ErrorRQLiteDBNotFound:    "RQLite Database not found",
	ErrorWriteToRQLiteDBFail: "Data failed to be written to rqlite db",
}
