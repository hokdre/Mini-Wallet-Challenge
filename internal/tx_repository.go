package internal

import (
	"context"
	"database/sql"
)

type TxRepository interface {
	Process(ctx context.Context, f func(context.Context, *sql.Tx) error) error
}

type txRepository struct {
	db *sql.DB
}

func NewTxRepository(db *sql.DB) *txRepository {
	return &txRepository{db: db}
}

func (t *txRepository) Process(ctx context.Context, f func(context.Context, *sql.Tx) error) error {
	tx, err := t.db.Begin()
	if err != nil {
		return err
	}

	err = f(ctx, tx)
	if err != nil {
		_ = tx.Rollback()
	} else {
		_ = tx.Commit()
	}

	return err
}
