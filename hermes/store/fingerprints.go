package store

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/binary"
	"fmt"
	"math"
	"strings"

	_ "github.com/mattn/go-sqlite3" //go-sqlite3
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/hermes/domain"
)

const (
	createFgTableStatement                       = `CREATE TABLE IF NOT EXISTS image_hash_to_image_fingerprint_table (sha256_hash_of_art_image_file text PRIMARY KEY, path_to_art_image_file text, new_model_image_fingerprint_vector array, datetime_fingerprint_added_to_database text, thumbnail_of_image text, request_type text, open_api_subset_id_string text,open_api_group_id_string text,collection_name_string text, registration_ticket_txid text)`
	getLatestFingerprintStatement                = `SELECT * FROM image_hash_to_image_fingerprint_table ORDER BY datetime_fingerprint_added_to_database DESC LIMIT 1`
	getFingerprintFromHashStatement              = `SELECT * FROM image_hash_to_image_fingerprint_table WHERE sha256_hash_of_art_image_file = ?`
	getFingerprintFromTxidStatement              = `SELECT * FROM image_hash_to_image_fingerprint_table WHERE registration_ticket_txid = ?`
	insertFingerprintStatement                   = `INSERT INTO image_hash_to_image_fingerprint_table(sha256_hash_of_art_image_file, path_to_art_image_file, new_model_image_fingerprint_vector, datetime_fingerprint_added_to_database, thumbnail_of_image, request_type, open_api_subset_id_string,open_api_group_id_string,collection_name_string,registration_ticket_txid) VALUES(?,?,?,?,?,?,?,?,?,?)`
	getNumberOfFingerprintsStatement             = `SELECT COUNT(*) as count FROM image_hash_to_image_fingerprint_table`
	createDoesNotImpactCollectionsTableStatement = `CREATE TABLE does_not_impact_collections_table(id integer not null PRIMARY KEY, collection_name_string text, sha256_hash_of_art_image_file text)`
)

type fingerprints struct {
	Sha256HashOfArtImageFile           string         `db:"sha256_hash_of_art_image_file,omitempty"`
	PathToArtImageFile                 string         `db:"path_to_art_image_file,omitempty"`
	ImageFingerprintVector             []byte         `db:"new_model_image_fingerprint_vector,omitempty"`
	DatetimeFingerprintAddedToDatabase string         `db:"datetime_fingerprint_added_to_database,omitempty"`
	ImageThumbnailAsBase64             string         `db:"thumbnail_of_image,omitempty"`
	RequestType                        string         `db:"request_type,omitempty"`
	IDString                           string         `db:"open_api_subset_id_string,omitempty"`
	OpenAPIGroupIDString               sql.NullString `db:"open_api_group_id_string"`
	CollectionNameString               sql.NullString `db:"collection_name_string"`
	RegistrationTicketTXID             sql.NullString `db:"registration_ticket_txid"`
}

func (r *fingerprints) toDomain() (*domain.DDFingerprints, error) {

	fp, err := readFloat64SliceFromBytes(r.ImageFingerprintVector)
	if err != nil {
		return nil, fmt.Errorf("failed to convert byte to float32: %w", err)
	}

	return &domain.DDFingerprints{
		Sha256HashOfArtImageFile:           r.Sha256HashOfArtImageFile,
		PathToArtImageFile:                 r.PathToArtImageFile,
		ImageFingerprintVector:             fp,
		DatetimeFingerprintAddedToDatabase: r.DatetimeFingerprintAddedToDatabase,
		ImageThumbnailAsBase64:             r.ImageThumbnailAsBase64,
		RequestType:                        r.RequestType,
		IDString:                           r.IDString,
		OpenAPIGroupIDString:               r.OpenAPIGroupIDString.String,
		CollectionNameString:               r.CollectionNameString.String,
		RegTXID:                            r.RegistrationTicketTXID.String,
	}, nil

}

// CheckNonSeedRecord checks if there's non-seed record
func (s *SQLiteStore) CheckNonSeedRecord(_ context.Context) (bool, error) {
	r := []fingerprints{}
	err := s.db.Select(&r, "SELECT * FROM image_hash_to_image_fingerprint_table")
	if err != nil {
		return false, fmt.Errorf("failed to get image_hash_to_image_fingerprint_table: %w", err)
	}

	for _, v := range r {
		if v.IDString != "SEED" {
			return true, nil
		}
	}

	return false, nil
}

// StoreFingerprint stores fingerprint
func (s *SQLiteStore) StoreFingerprint(ctx context.Context, input *domain.DDFingerprints) error {
	fp := writeFloat64SliceToBytes(input.ImageFingerprintVector)

	_, err := s.db.Exec(`INSERT INTO image_hash_to_image_fingerprint_table(sha256_hash_of_art_image_file,
		 path_to_art_image_file, new_model_image_fingerprint_vector, datetime_fingerprint_added_to_database,
		  thumbnail_of_image, request_type, open_api_subset_id_string,open_api_group_id_string,collection_name_string, registration_ticket_txid) VALUES(?,?,?,?,?,?,?,?,?,?)`, input.Sha256HashOfArtImageFile,
		input.PathToArtImageFile, fp, input.DatetimeFingerprintAddedToDatabase, input.ImageThumbnailAsBase64, input.RequestType, input.IDString, input.OpenAPIGroupIDString, input.CollectionNameString, input.RegTXID)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("Failed to insert fingerprint record")
		return err
	}

	if input.DoesNotImpactTheFollowingCollectionsString == "" {
		log.WithContext(ctx).Info("list of non-impacted collection is empty")

		return nil
	}

	commaSeparatedDoesNotImpactCollectionNamesList := strings.Split(input.DoesNotImpactTheFollowingCollectionsString, ",")

	for _, collectionName := range commaSeparatedDoesNotImpactCollectionNamesList {
		_, err = s.db.Exec(`INSERT INTO does_not_impact_collections_table(collection_name_string, sha256_hash_of_art_image_file)
			VALUES(?,?)`, collectionName, input.Sha256HashOfArtImageFile)
		if err != nil {
			log.WithContext(ctx).WithError(err).Error("Failed to insert collection name  in does not impact collection table")
			return err
		}
	}

	return nil
}

// GetLatestFingerprints gets latest fg
func (s *SQLiteStore) GetLatestFingerprints(ctx context.Context) (*domain.DDFingerprints, error) {
	r := fingerprints{}
	err := s.db.Get(&r, getLatestFingerprintStatement)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("dd database is empty")
		}

		return nil, fmt.Errorf("failed to get record by key %w", err)
	}

	dd, err := r.toDomain()
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("converting db data to dd-fingerprint struct failure")
	}

	return dd, nil
}

// IfFingerprintExists checks if fg exists against the hash
func (s *SQLiteStore) IfFingerprintExists(_ context.Context, hash string) (bool, error) {
	r := fingerprints{}
	err := s.db.Get(&r, getFingerprintFromHashStatement, hash)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}

		return false, fmt.Errorf("failed to get record by key %w : key: %s", err, hash)
	}

	if r.Sha256HashOfArtImageFile == "" {
		return false, nil
	}

	return true, nil
}

// IfFingerprintExistsByRegTxid checks if fg exists against the reg txid
func (s *SQLiteStore) IfFingerprintExistsByRegTxid(_ context.Context, txid string) (bool, error) {
	r := fingerprints{}
	err := s.db.Get(&r, getFingerprintFromTxidStatement, txid)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}

		return false, fmt.Errorf("failed to get record by key %w : key: %s", err, txid)
	}

	if r.Sha256HashOfArtImageFile == "" {
		return false, nil
	}

	return true, nil
}

// GetFingerprintsCount gets fingerprint count
func (s *SQLiteStore) GetFingerprintsCount(_ context.Context) (int64, error) {
	var Data struct {
		Count int64 `db:"count"`
	}

	statement := getNumberOfFingerprintsStatement
	err := s.db.Get(&Data, statement)
	if err != nil {
		return 0, errors.Errorf("query: %w", err)
	}

	return Data.Count, nil
}

func float64ToBytes(f float64) []byte {
	bits := math.Float64bits(f)
	bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytes, bits)
	return bytes
}

func bytesToFloat64(b []byte) float64 {
	bits := binary.LittleEndian.Uint64(b)
	return math.Float64frombits(bits)
}

func writeFloat64SliceToBytes(data []float64) []byte {
	header := fmt.Sprintf("{'descr': '<f8', 'fortran_order': False, 'shape': (%d, 1), }", len(data))
	paddedHeader := header + strings.Repeat(" ", (128-len(header)-1)) + "\n"

	var buf bytes.Buffer
	buf.Write([]byte("\x93NUMPY")) // Magic string
	buf.WriteByte(0x01)            // Major version
	buf.WriteByte(0x00)            // Minor version
	buf.Write([]byte{0x76, 0x00})  // Header length (little endian uint16)
	buf.WriteString(paddedHeader)  // Header
	for _, value := range data {
		buf.Write(float64ToBytes(value)) // Data
	}

	return buf.Bytes()
}

func readFloat64SliceFromBytes(content []byte) ([]float64, error) {
	if !bytes.HasPrefix(content, []byte("\x93NUMPY")) {
		return nil, fmt.Errorf("invalid magic string")
	}

	headerLength := binary.LittleEndian.Uint16(content[8:10])
	header := string(content[10 : 10+headerLength])

	var dataLength int
	_, err := fmt.Sscanf(header, "{'descr': '<f8', 'fortran_order': False, 'shape': (%d, 1), }", &dataLength)
	if err != nil {
		return nil, err
	}

	dataBytes := content[10+headerLength:]
	if len(dataBytes)%8 != 0 {
		return nil, fmt.Errorf("invalid byte slice length")
	}

	data := make([]float64, dataLength)
	for i := 0; i < dataLength; i++ {
		data[i] = bytesToFloat64(dataBytes[i*8 : (i+1)*8])
	}

	return data, nil
}
