package internal

import (
	"context"
	"database/sql"

	"github.com/hokdre/mini-ewallet/internal/model"
)

type WalletFilter struct {
	OwnedBies []string
	IDs       []string
}

type WalletRepository interface {
	GetOne(ctx context.Context, filter WalletFilter) (model.Wallet, error)
	Update(ctx context.Context, wallet model.Wallet) error
	CreateTx(ctx context.Context, tx *sql.Tx, newWallet model.Wallet) error
	Increment(ctx context.Context, tx *sql.Tx, wallet model.Wallet, amount int64) (int64, error)
	Decrement(ctx context.Context, tx *sql.Tx, wallet model.Wallet, amount int64) (int64, error)
}
