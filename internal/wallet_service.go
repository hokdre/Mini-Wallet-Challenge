package internal

import (
	"context"

	"github.com/google/uuid"
	"github.com/hokdre/mini-ewallet/internal/model"
)

type WalletService interface {
	Init(ctx context.Context, externalID string) (string, error)
	Enable(ctx context.Context, accountID uuid.UUID) (model.Wallet, error)
	Disable(ctx context.Context, accountID uuid.UUID) (model.Wallet, error)
	Get(ctx context.Context, accountID uuid.UUID) (model.Wallet, error)
	GetTransactions(ctx context.Context, accountID uuid.UUID) ([]model.Transaction, error)
	Deposit(ctx context.Context, accountID uuid.UUID, transaction model.Transaction) (model.Transaction, error)
	Withdrawal(ctx context.Context, accountID uuid.UUID, transaction model.Transaction) (model.Transaction, error)
}
