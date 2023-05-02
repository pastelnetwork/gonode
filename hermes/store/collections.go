package store

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/hermes/domain"
)

const (
	createCollectionsTableStatement = `CREATE TABLE IF NOT EXISTS collections_table (collection_ticket_txid text PRIMARY KEY, collection_name_string text, collection_ticket_activation_block_height int, collection_final_allowed_block_height int, max_permitted_open_nsfw_score float, minimum_similarity_score_to_first_entry_in_collection float, collection_state text, datetime_collection_state_updated text)`
)

type collection struct {
	CollectionTicketTXID                           string         `db:"collection_ticket_txid"`
	CollectionName                                 string         `db:"collection_name_string"`
	CollectionTicketActivationBlockHeight          int            `db:"collection_ticket_activation_block_height"`
	CollectionFinalAllowedBlockHeight              int            `db:"collection_final_allowed_block_height"`
	MaxPermittedOpenNSFWScore                      float64        `db:"max_permitted_open_nsfw_score"`
	MinimumSimilarityScoreToFirstEntryInCollection float64        `db:"minimum_similarity_score_to_first_entry_in_collection"`
	CollectionState                                sql.NullString `db:"collection_state"`
	DatetimeCollectionStateUpdated                 sql.NullString `db:"datetime_collection_state_updated"`
}

func (c *collection) ToDomain() *domain.Collection {
	return &domain.Collection{
		CollectionTicketTXID:                           c.CollectionTicketTXID,
		CollectionName:                                 c.CollectionName,
		CollectionTicketActivationBlockHeight:          c.CollectionTicketActivationBlockHeight,
		CollectionFinalAllowedBlockHeight:              c.CollectionFinalAllowedBlockHeight,
		MaxPermittedOpenNSFWScore:                      c.MaxPermittedOpenNSFWScore,
		MinimumSimilarityScoreToFirstEntryInCollection: c.MinimumSimilarityScoreToFirstEntryInCollection,
		CollectionState:                                domain.CollectionState(c.CollectionState.String),
		DatetimeCollectionStateUpdated:                 c.DatetimeCollectionStateUpdated.String,
	}
}

// IfCollectionExists checks if collection exists against the id
func (s *SQLiteStore) IfCollectionExists(_ context.Context, collectionTxID string) (bool, error) {
	c := collection{}
	getCollectionByIDQuery := `SELECT * FROM collections_table WHERE collection_ticket_txid = ?`
	err := s.db.Get(&c, getCollectionByIDQuery, collectionTxID)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}

		return false, fmt.Errorf("failed to get record: %w : collection_ticket_txid: %s", err, collectionTxID)
	}

	if c.CollectionTicketTXID == "" {
		return false, nil
	}

	return true, nil
}

// StoreCollection store collection object to DB
func (s *SQLiteStore) StoreCollection(_ context.Context, c domain.Collection) error {
	_, err := s.db.Exec(`INSERT INTO collections_table(collection_ticket_txid,
		 collection_name_string, collection_ticket_activation_block_height, collection_final_allowed_block_height,
		  max_permitted_open_nsfw_score, minimum_similarity_score_to_first_entry_in_collection, collection_state, 
          datetime_collection_state_updated) VALUES(?,?,?,?,?,?,?,?)`, c.CollectionTicketTXID,
		c.CollectionName, c.CollectionTicketActivationBlockHeight, c.CollectionFinalAllowedBlockHeight,
		c.MaxPermittedOpenNSFWScore, c.MinimumSimilarityScoreToFirstEntryInCollection, c.CollectionState.String(),
		c.DatetimeCollectionStateUpdated)
	if err != nil {
		return fmt.Errorf("failed to insert collection record: %w", err)
	}

	return nil
}

func (s *SQLiteStore) GetAllInProcessCollections(ctx context.Context) ([]*domain.Collection, error) {
	c := []*collection{}

	getCollectionByStateQuery := `SELECT * FROM collections_table WHERE collection_state = ?`
	err := s.db.GetContext(ctx, &c, getCollectionByStateQuery, domain.InProcessCollectionState)
	if err != nil {
		return nil, fmt.Errorf("failed to get record: %w : collection_state: %s", err, domain.InProcessCollectionState)
	}

	collections := []*domain.Collection{}
	for _, collection := range c {
		collections = append(collections, collection.ToDomain())
	}

	return collections, nil
}

// GetCollection get collection object from DB
func (s *SQLiteStore) GetCollection(ctx context.Context, collectionTxID string) (*domain.Collection, error) {
	c := collection{}

	getCollectionByIDQuery := `SELECT * FROM collections_table WHERE collection_ticket_txid = ?`
	err := s.db.GetContext(ctx, &c, getCollectionByIDQuery, collectionTxID)
	if err != nil {
		return nil, fmt.Errorf("failed to get record: %w : collection_ticket_txid: %s", err, collectionTxID)
	}

	return c.ToDomain(), nil
}

// FinalizeCollectionState finalizes collection state
func (s *SQLiteStore) FinalizeCollectionState(_ context.Context, txid string) error {
	_, err := s.db.Exec(`UPDATE collections_table SET collection_state = ? WHERE collection_ticket_txid `, domain.FinalizedCollectionState, txid)
	if err != nil {
		return fmt.Errorf("failed to update collection record: %w", err)
	}

	return nil
}

type nonImpactedCollections struct {
	ID                       int    `db:"id"`
	CollectionName           string `db:"collection_name_string"`
	Sha256HashOfArtImageFile string `db:"sha256_hash_of_art_image_file"`
}

func (n *nonImpactedCollections) toDomain() *domain.NonImpactedCollection {
	return &domain.NonImpactedCollection{
		ID:                       n.ID,
		CollectionName:           n.CollectionName,
		Sha256HashOfArtImageFile: n.Sha256HashOfArtImageFile,
	}
}

// GetDoesNotImpactCollections retrieves the collection names which should not be impacted linked to image hash
func (s *SQLiteStore) GetDoesNotImpactCollections(ctx context.Context, hash string) (domain.NonImpactedCollections, error) {
	dncs := []nonImpactedCollections{}
	query := `SELECT * from does_not_impact_collections_table where sha256_hash_of_art_image_file = ?`
	err := s.db.SelectContext(ctx, &dncs, query, hash)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.Errorf("no records found")
		}

		return nil, fmt.Errorf("failed to get record by hash %w", err)
	}

	var nonImpactedCollections domain.NonImpactedCollections
	for _, c := range dncs {
		nonImpactedCollections = append(nonImpactedCollections, c.toDomain())
	}

	return nonImpactedCollections, nil
}
