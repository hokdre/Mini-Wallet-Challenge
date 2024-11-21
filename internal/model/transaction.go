package model

import (
	"time"

	"github.com/google/uuid"
)

var (
	TransactionType = struct {
		Withdrawal string
		Deposit    string
	}{
		Withdrawal: "withdrawal",
		Deposit:    "deposit",
	}

	TransactionStatus = struct {
		Pending string
		Success string
		Failed  string
	}{
		Pending: "pending",
		Success: "success",
		Failed:  "failed",
	}
)

type Transaction struct {
	ID           uuid.UUID  `json:"id" db:"id" validate:"required"`
	WalletID     uuid.UUID  `json:"wallet_id" db:"wallet_id" validate:"required"`
	Type         string     `json:"type" db:"type" validate:"enumTransactionType"`
	Status       string     `json:"status" db:"status" validate:"enumTransactionStatus"`
	TransactedAt *time.Time `json:"transacted_at" db:"transacted_at"`
	Amount       int64      `json:"amount" db:"amount" validate:"gte=1"`
	ReferenceID  string     `json:"reference_id" db:"reference_id" validate:"required"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at" validate:"required"`
	UpdatedAt    time.Time  `json:"updated_at" db:"updated_at" validate:"required"`
}
