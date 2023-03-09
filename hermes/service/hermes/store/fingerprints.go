package store

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3" //go-sqlite3
	"github.com/sbinet/npyio"
	"gonum.org/v1/gonum/mat"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/hermes/service/hermes/domain"
)

const (
	createFgTableStatement           = `CREATE TABLE IF NOT EXISTS image_hash_to_image_fingerprint_table (sha256_hash_of_art_image_file text PRIMARY KEY, path_to_art_image_file text, new_model_image_fingerprint_vector array, datetime_fingerprint_added_to_database text, thumbnail_of_image text, request_type text, open_api_subset_id_string text,open_api_group_id_string text,collection_name_string text,does_not_impact_the_following_collection_strings text)`
	getLatestFingerprintStatement    = `SELECT * FROM image_hash_to_image_fingerprint_table ORDER BY datetime_fingerprint_added_to_database DESC LIMIT 1`
	getFingerprintFromHashStatement  = `SELECT * FROM image_hash_to_image_fingerprint_table WHERE sha256_hash_of_art_image_file = ?`
	insertFingerprintStatement       = `INSERT INTO image_hash_to_image_fingerprint_table(sha256_hash_of_art_image_file, path_to_art_image_file, new_model_image_fingerprint_vector, datetime_fingerprint_added_to_database, thumbnail_of_image, request_type, open_api_subset_id_string,open_api_group_id_string,collection_name_string,does_not_impact_the_following_collection_strings) VALUES(?,?,?,?,?,?,?,?,?,?)`
	getNumberOfFingerprintsStatement = `SELECT COUNT(*) as count FROM image_hash_to_image_fingerprint_table`
)

type fingerprints struct {
	Sha256HashOfArtImageFile           string `db:"sha256_hash_of_art_image_file,omitempty"`
	PathToArtImageFile                 string `db:"path_to_art_image_file,omitempty"`
	ImageFingerprintVector             []byte `db:"new_model_image_fingerprint_vector,omitempty"`
	DatetimeFingerprintAddedToDatabase string `db:"datetime_fingerprint_added_to_database,omitempty"`
	ImageThumbnailAsBase64             string `db:"thumbnail_of_image,omitempty"`
	RequestType                        string `db:"request_type,omitempty"`
	IDString                           string `db:"open_api_subset_id_string,omitempty"`
	OpenAPIGroupIDString               string `db:"open_api_group_id_string"`
	CollectionNameString               string `db:"collection_name_string"`
	DoesNotImpactCollectionStrings     string `db:"does_not_impact_the_following_collection_strings"`
}

func (r *fingerprints) toDomain() (*domain.DDFingerprints, error) {
	f := bytes.NewBuffer(r.ImageFingerprintVector)

	var fp []float64
	if err := npyio.Read(f, &fp); err != nil {
		return nil, errors.New("Failed to convert npy to float64")
	}

	return &domain.DDFingerprints{
		Sha256HashOfArtImageFile:                   r.Sha256HashOfArtImageFile,
		PathToArtImageFile:                         r.PathToArtImageFile,
		ImageFingerprintVector:                     fp,
		DatetimeFingerprintAddedToDatabase:         r.DatetimeFingerprintAddedToDatabase,
		ImageThumbnailAsBase64:                     r.ImageThumbnailAsBase64,
		RequestType:                                r.RequestType,
		IDString:                                   r.IDString,
		OpenAPIGroupIDString:                       r.OpenAPIGroupIDString,
		CollectionNameString:                       r.CollectionNameString,
		DoesNotImpactTheFollowingCollectionsString: r.DoesNotImpactCollectionStrings,
	}, nil
}

// StoreFingerprint stores fingerprint
func (s *SQLiteStore) StoreFingerprint(ctx context.Context, input *domain.DDFingerprints) error {
	encodeFloat2Npy := func(v []float64) ([]byte, error) {
		// create numpy matrix Nx1
		m := mat.NewDense(len(v), 1, v)
		f := bytes.NewBuffer(nil)
		if err := npyio.Write(f, m); err != nil {
			return nil, errors.Errorf("encode to npy: %w", err)
		}
		return f.Bytes(), nil
	}

	fp, err := encodeFloat2Npy(input.ImageFingerprintVector)
	if err != nil {
		return err
	}

	_, err = s.db.Exec(`INSERT INTO image_hash_to_image_fingerprint_table(sha256_hash_of_art_image_file,
		 path_to_art_image_file, new_model_image_fingerprint_vector, datetime_fingerprint_added_to_database,
		  thumbnail_of_image, request_type, open_api_subset_id_string,open_api_group_id_string,collection_name_string,does_not_impact_the_following_collection_strings) VALUES(?,?,?,?,?,?,?,?,?,?)`, input.Sha256HashOfArtImageFile,
		input.PathToArtImageFile, fp, input.DatetimeFingerprintAddedToDatabase, input.ImageThumbnailAsBase64, input.RequestType, input.IDString, input.OpenAPIGroupIDString, input.CollectionNameString, input.DoesNotImpactTheFollowingCollectionsString)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("Failed to insert fingerprint record")
		return err
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
		log.WithContext(ctx).WithError(err).Error("converting db data to dd-fingerprints struct failure")
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
