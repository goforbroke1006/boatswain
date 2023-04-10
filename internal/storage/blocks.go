package storage

import (
	"context"
	"database/sql"
	"encoding/json"

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
		index      uint64
		hash       string
		phash      string
		ts         int64
		dataAsJson string
		data       []*domain.Transaction
	)

	var lastBlock *domain.Block

	rows.Next()

	if scanErr := rows.Scan(&index, &hash, &phash, &ts, &dataAsJson); scanErr != nil {
		return nil, scanErr
	}

	_ = json.Unmarshal([]byte(dataAsJson), &data)

	lastBlock = &domain.Block{
		ID:       domain.BlockIndex(index),
		Hash:     domain.BlockHash(hash),
		PrevHash: domain.BlockHash(phash),
		Ts:       ts,
		Data:     nil, // TODO: parse data from DB
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
		dataAsJson, _ := json.Marshal(b.Data)
		_, execErr := tx.ExecContext(ctx, `INSERT INTO blocks VALUES (?, ?, ?, ?, ?)`,
			b.ID, b.Hash, b.PrevHash, b.Ts, string(dataAsJson))
		if execErr != nil {
			return execErr
		}
	}

	return nil
}

func (s blockStorage) LoadLast(count uint64) ([]*domain.Block, error) {
	//TODO implement me
	_ = count
	panic("implement me")
}
