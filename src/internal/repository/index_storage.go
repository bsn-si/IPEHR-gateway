package repository

import (
	"context"
	"fmt"

	"github.com/bsn-si/IPEHR-gateway/src/internal/models"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/errors"
	"github.com/jmoiron/sqlx"
)

type IndexStorage struct {
	db *sqlx.DB
}

func NewIndexStorage(db *sqlx.DB) *IndexStorage {
	return &IndexStorage{
		db: db,
	}
}

func (store *IndexStorage) AddNewIndexObject(ctx context.Context, chunk models.IndexChunk) error {
	const query = `INSERT INTO tree_index_chunks (key, group_id, data_id, ehr_id, data, hash)
	VALUES (:key, :group_id, :data_id, :ehr_id, :data, :hash);`

	if _, err := store.db.NamedExecContext(ctx, query, chunk); err != nil {
		return fmt.Errorf("cannot add index into db: %w", err)
	}

	return nil
}

func (store *IndexStorage) GetAllIndexObjects(ctx context.Context) ([]models.IndexChunk, error) {
	const query = `SELECT key, group_id, data_id, ehr_id, data, hash FROM tree_index_chunks ORDER BY created_at;`

	result := []models.IndexChunk{}
	if err := store.db.SelectContext(ctx, &result, query); err != nil {
		return nil, errors.Wrap(err, "cannot get index chunks")
	}

	return result, nil
}
