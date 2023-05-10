package common

// This file contains the keys(tags) of structured logging fields as constants regarding the status history for
// get task history endpoint
// It must be made necessary to define all the field tag names used in structured logs for get task history endpoint

const (
	//FieldBlockHeight represents creator block height
	FieldBlockHeight = "block_height"
	//FieldTaskType represents type of task e.g. cascade registration & nft download etc
	FieldTaskType = "task_type"
	//FieldFileSize represents the size of file in MBs
	FieldFileSize = "file_size"
	//FieldFee represents fee for registration task
	FieldFee = "fee"
	//FieldRegTicketTxnID represents txn id of registered ticket
	FieldRegTicketTxnID = "reg_ticket_txn_id"
	//FieldActivateTicketTxnID represents txn id of the activated ticket
	FieldActivateTicketTxnID = "activate_ticket_txn_id"
	//FieldErrorDetail represents the detail of error occurred
	FieldErrorDetail = "error_detail"
	//FieldBurnTxnID represents burnt txn id
	FieldBurnTxnID = "burn_txn_id"
	//FieldNFTRegisterTaskMaxFee represents the nft register task request fee
	FieldNFTRegisterTaskMaxFee = "nft_register_task_max_fee"
	//FieldSpendableAddress represents the spendable address of the user registering NFT
	FieldSpendableAddress = "spendable_address"
	//FieldMeshNodes represents the ips of nodes establish a mesh in reg endpoint
	FieldMeshNodes = "mesh_nodes_ips"
	// FieldNodes represents the nodes used for the task
	FieldNodes = "nodes"
)
