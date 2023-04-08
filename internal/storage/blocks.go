package storage

import (
	"context"
	"database/sql"

	"github.com/goforbroke1006/boatswain/domain"
)

func NewBlockStorage(db *sql.DB) domain.BlockStorage {
	return &blockStorage{db: db}
}

var _ domain.BlockStorage = (*blockStorage)(nil)

type blockStorage struct {
	db *sql.DB
}

func (s blockStorage) GetCount(ctx context.Context) (uint64, error) {
	const query = `SELECT COUNT(*) FROM blocks`

	rows, rowsErr := s.db.QueryContext(ctx, query)
	if rowsErr != nil {
		return 0, rowsErr
	}
	defer func() { _ = rows.Close() }()

	var count uint64

	rows.Next()
	if scanErr := rows.Scan(&count); scanErr != nil {
		return 0, scanErr
	}

	if rows.Err() != nil {
		return 0, rows.Err()
	}

	return count, nil
}

func (s blockStorage) GetLast(ctx context.Context) (*domain.Block, error) {
	const query = `
		SELECT 
			"index", "hash", "previous_hash", "timestamp", "data" 
		FROM blocks 
		ORDER BY "index" DESC 
		LIMIT 1`

	rows, rowsErr := s.db.QueryContext(ctx, query)
	if rowsErr != nil {
		return nil, rowsErr
	}
	defer func() { _ = rows.Close() }()

	var (
		index uint64
		hash  string
		phash string
		ts    int64
		data  string
	)

	var lastBlock *domain.Block

	for rows.Next() {
		if scanErr := rows.Scan(&index, &hash, &phash, &ts, &data); scanErr != nil {
			return nil, scanErr
		}

		lastBlock = &domain.Block{
			Index:        domain.BlockIndex(index),
			Hash:         domain.BlockHash(hash),
			PreviousHash: domain.BlockHash(phash),
			Timestamp:    ts,
			Data:         nil, // TODO: parse data from DB
		}
		break
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return lastBlock, nil
}

func (s blockStorage) Store(ctx context.Context, blocks ...*domain.Block) error {
	tx, txErr := s.db.BeginTx(ctx, nil)
	if txErr != nil {
		return txErr
	}
	defer func() {
		_ = tx.Commit()
	}()

	for _, b := range blocks {
		_, execErr := tx.ExecContext(ctx, `INSERT INTO blocks VALUES (?, ?, ?, ?, ?)`,
			b.Index, b.Hash, b.PreviousHash, b.Timestamp, b.Data)
		if execErr != nil {
			return execErr
		}
	}

	return nil
}
