package db

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/hisshihi/order-of-venhicles-services/db/sqlc"
)

// Store предоставляет все фукнции для запросов и объединение запросов в транзакции
type Store struct {
	*sqlc.Queries
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db: db,
		Queries: sqlc.New(db),
	}
}

func (store *Store) execTx(ctx context.Context, fn func(*sqlc.Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := sqlc.New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
	}

	return tx.Commit()
}

// type TransferTxParams struct {
	
// }

// func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {

// }
