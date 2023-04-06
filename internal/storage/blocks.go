package storage

import (
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

func (s blockStorage) Store(b *domain.Block) error {
	_, err := s.db.Exec(`INSERT INTO blocks VALUES (?, ?, ?, ?, ?)`,
		b.Index, b.Hash, b.PreviousHash, b.Timestamp, b.Data)
	return err
}

func (s blockStorage) Load() ([]*domain.Block, error) {
	rows, err := s.db.Query(`
	SELECT 
	    "index", "hash", "previous_hash", "timestamp", "data" 
	FROM blocks 
	ORDER BY "index"`)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var (
		index uint64
		hash  string
		phash string
		ts    int64
		data  string
	)

	var result []*domain.Block

	for rows.Next() {
		if err := rows.Scan(&index, &hash, &phash, &ts, &data); err != nil {
			return nil, err
		}

		b := &domain.Block{
			Index:        domain.BlockIndex(index),
			Hash:         domain.BlockHash(hash),
			PreviousHash: domain.BlockHash(phash),
			Timestamp:    ts,
			Data:         []byte(data),
		}
		result = append(result, b)
	}

	return result, nil
}
