package internal

import (
	"context"
	"database/sql"

	"github.com/hokdre/mini-ewallet/internal/model"
)

type AccountFilter struct {
	ExternalIDs []string
	IDs         []string
}

type AccountRepository interface {
	Get(ctx context.Context, filter AccountFilter) (model.Account, error)
	CreateTx(ctx context.Context, tx *sql.Tx, newAcc model.Account) error
}
