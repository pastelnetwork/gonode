package ticketstore

import (
	"github.com/pastelnetwork/gonode/common/types"
)

type FilesQueries interface {
	UpsertFile(file types.File) error
	GetFileByID(fileID string) (*types.File, error)
	GetFilesByBaseFileID(baseFileID string) ([]*types.File, error)
	GetFileByTaskID(taskID string) (*types.File, error)
}

// UpsertFile inserts a new file into the files table
func (s *TicketStore) UpsertFile(file types.File) error {
	const upsertQuery = `
        INSERT INTO files (
            file_id, upload_timestamp, path, file_index, base_file_id, task_id, 
            reg_txid, activation_txid, req_burn_txn_amount, burn_txn_id, 
            req_amount, is_concluded, cascade_metadata_ticket_id, uuid_key, 
            hash_of_original_big_file, name_of_original_big_file_with_ext, 
            size_of_original_big_file, data_type_of_original_big_file, 
            start_block, done_block
        ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
        ON CONFLICT(file_id) 
        DO UPDATE SET
            upload_timestamp = COALESCE(excluded.upload_timestamp, files.upload_timestamp),
            path = COALESCE(excluded.path, files.path),
            file_index = COALESCE(excluded.file_index, files.file_index),
            base_file_id = COALESCE(excluded.base_file_id, files.base_file_id),
            task_id = COALESCE(excluded.task_id, files.task_id),
            reg_txid = COALESCE(excluded.reg_txid, files.reg_txid),
            activation_txid = COALESCE(excluded.activation_txid, files.activation_txid),
            req_burn_txn_amount = COALESCE(excluded.req_burn_txn_amount, files.req_burn_txn_amount),
            burn_txn_id = COALESCE(excluded.burn_txn_id, files.burn_txn_id),
            req_amount = COALESCE(excluded.req_amount, files.req_amount),
            is_concluded = COALESCE(excluded.is_concluded, files.is_concluded),
            cascade_metadata_ticket_id = COALESCE(excluded.cascade_metadata_ticket_id, files.cascade_metadata_ticket_id),
            uuid_key = COALESCE(excluded.uuid_key, files.uuid_key),
            hash_of_original_big_file = COALESCE(excluded.hash_of_original_big_file, files.hash_of_original_big_file),
            name_of_original_big_file_with_ext = COALESCE(excluded.name_of_original_big_file_with_ext, files.name_of_original_big_file_with_ext),
            size_of_original_big_file = COALESCE(excluded.size_of_original_big_file, files.size_of_original_big_file),
            data_type_of_original_big_file = COALESCE(excluded.data_type_of_original_big_file, files.data_type_of_original_big_file),
            start_block = COALESCE(excluded.start_block, files.start_block),
            done_block = COALESCE(excluded.done_block, files.done_block);`

	_, err := s.db.Exec(upsertQuery,
		file.FileID, file.UploadTimestamp, file.Path, file.FileIndex, file.BaseFileID, file.TaskID,
		file.RegTxid, file.ActivationTxid, file.ReqBurnTxnAmount, file.BurnTxnID,
		file.ReqAmount, file.IsConcluded, file.CascadeMetadataTicketID, file.UUIDKey,
		file.HashOfOriginalBigFile, file.NameOfOriginalBigFileWithExt,
		file.SizeOfOriginalBigFile, file.DataTypeOfOriginalBigFile,
		file.StartBlock, file.DoneBlock)
	if err != nil {
		return err
	}

	return nil
}

// GetFileByID retrieves a file by its ID from the files table
func (s *TicketStore) GetFileByID(fileID string) (*types.File, error) {
	const selectQuery = `
        SELECT file_id, upload_timestamp, path, file_index, base_file_id, task_id, 
               reg_txid, activation_txid, req_burn_txn_amount, burn_txn_id, 
               req_amount, is_concluded, cascade_metadata_ticket_id, uuid_key, 
               hash_of_original_big_file, name_of_original_big_file_with_ext, 
               size_of_original_big_file, data_type_of_original_big_file, 
               start_block, done_block
        FROM files
        WHERE file_id = ?;`

	row := s.db.QueryRow(selectQuery, fileID)

	var file types.File
	err := row.Scan(
		&file.FileID, &file.UploadTimestamp, &file.Path, &file.FileIndex, &file.BaseFileID, &file.TaskID,
		&file.RegTxid, &file.ActivationTxid, &file.ReqBurnTxnAmount, &file.BurnTxnID,
		&file.ReqAmount, &file.IsConcluded, &file.CascadeMetadataTicketID, &file.UUIDKey,
		&file.HashOfOriginalBigFile, &file.NameOfOriginalBigFileWithExt,
		&file.SizeOfOriginalBigFile, &file.DataTypeOfOriginalBigFile,
		&file.StartBlock, &file.DoneBlock)
	if err != nil {
		return nil, err
	}

	return &file, nil
}

// GetFileByTaskID retrieves a file by its task-id from the files table
func (s *TicketStore) GetFileByTaskID(taskID string) (*types.File, error) {
	const selectQuery = `
        SELECT file_id, upload_timestamp, path, file_index, base_file_id, task_id, 
               reg_txid, activation_txid, req_burn_txn_amount, burn_txn_id, 
               req_amount, is_concluded, cascade_metadata_ticket_id, uuid_key, 
               hash_of_original_big_file, name_of_original_big_file_with_ext, 
               size_of_original_big_file, data_type_of_original_big_file, 
               start_block, done_block
        FROM files
        WHERE task_id = ?;`

	row := s.db.QueryRow(selectQuery, taskID)

	var file types.File
	err := row.Scan(
		&file.FileID, &file.UploadTimestamp, &file.Path, &file.FileIndex, &file.BaseFileID, &file.TaskID,
		&file.RegTxid, &file.ActivationTxid, &file.ReqBurnTxnAmount, &file.BurnTxnID,
		&file.ReqAmount, &file.IsConcluded, &file.CascadeMetadataTicketID, &file.UUIDKey,
		&file.HashOfOriginalBigFile, &file.NameOfOriginalBigFileWithExt,
		&file.SizeOfOriginalBigFile, &file.DataTypeOfOriginalBigFile,
		&file.StartBlock, &file.DoneBlock)
	if err != nil {
		return nil, err
	}

	return &file, nil
}

// GetFilesByBaseFileID retrieves files by base_file_id from the files table
func (s *TicketStore) GetFilesByBaseFileID(baseFileID string) ([]*types.File, error) {
	const selectQuery = `
        SELECT file_id, upload_timestamp, path, file_index, base_file_id, task_id, 
               reg_txid, activation_txid, req_burn_txn_amount, burn_txn_id, 
               req_amount, is_concluded, cascade_metadata_ticket_id, uuid_key, 
               hash_of_original_big_file, name_of_original_big_file_with_ext, 
               size_of_original_big_file, data_type_of_original_big_file, 
               start_block, done_block
        FROM files
        WHERE base_file_id = ?;`

	rows, err := s.db.Query(selectQuery, baseFileID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var files []*types.File
	for rows.Next() {
		var file types.File
		err := rows.Scan(
			&file.FileID, &file.UploadTimestamp, &file.Path, &file.FileIndex, &file.BaseFileID, &file.TaskID,
			&file.RegTxid, &file.ActivationTxid, &file.ReqBurnTxnAmount, &file.BurnTxnID,
			&file.ReqAmount, &file.IsConcluded, &file.CascadeMetadataTicketID, &file.UUIDKey,
			&file.HashOfOriginalBigFile, &file.NameOfOriginalBigFileWithExt,
			&file.SizeOfOriginalBigFile, &file.DataTypeOfOriginalBigFile,
			&file.StartBlock, &file.DoneBlock)
		if err != nil {
			return nil, err
		}
		files = append(files, &file)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return files, nil
}
