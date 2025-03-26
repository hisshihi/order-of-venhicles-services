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

// DeleteCategoryTxParams представляет параметры для транзакции удаления категории
type DeleteCategoryTxParams struct {
    CategoryID int64
}

// DeleteCategoryTxResult представляет результат транзакции удаления категории
type DeleteCategoryTxResult struct {
    DeletedOrders     int64
    DeletedServices   int64
    DeletedCategory   bool
}

// DeleteCategoryTx выполняет транзакцию для удаления категории и всех связанных с ней данных
func (store *Store) DeleteCategoryTx(ctx context.Context, arg DeleteCategoryTxParams) (DeleteCategoryTxResult, error) {
    var result DeleteCategoryTxResult

    err := store.execTx(ctx, func(q *sqlc.Queries) error {
        var err error

        // 1. Удаляем заказы, связанные с этой категорией
        deletedOrders, err := q.DeleteOrdersByCategoryID(ctx, arg.CategoryID)
        if err != nil {
            return fmt.Errorf("ошибка при удалении заказов: %w", err)
        }
        result.DeletedOrders = deletedOrders

        // 2. Удаляем услуги, связанные с этой категорией
        deletedServices, err := q.DeleteServicesByCategoryID(ctx, arg.CategoryID)
        if err != nil {
            return fmt.Errorf("ошибка при удалении услуг: %w", err)
        }
        result.DeletedServices = deletedServices

        // 3. Удаляем саму категорию
        deletedCategories, err := q.DeleteServiceCategory(ctx, arg.CategoryID)
        if err != nil {
            return fmt.Errorf("ошибка при удалении категории: %w", err)
        }
        if deletedCategories > 0 {
            result.DeletedCategory = true
        }

        return nil
    })

    return result, err
}
