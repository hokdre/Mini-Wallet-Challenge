package internal

import (
	"context"
	"database/sql"

	"github.com/hokdre/mini-ewallet/internal/model"
)

type TransactionFilter struct {
	WalletIDs []string
	IDs       []string
}

type TransactionRepository interface {
	List(ctx context.Context, filter TransactionFilter) ([]model.Transaction, error)
	Create(ctx context.Context, newTransaction model.Transaction) error
	UpdateTx(ctx context.Context, tx *sql.Tx, transaction model.Transaction) (err error)
}
