package account

import (
	"context"
	"database/sql"

	"github.com/hokdre/mini-ewallet/internal"
	"github.com/hokdre/mini-ewallet/internal/model"
	"github.com/lib/pq"
)

const (
	defaultOrderDirection = "DESC"
	defaultOffset         = 0
	defaultOrderColumn    = "created_at"

	qCreate = `INSERT INTO accounts(
		id, external_id, created_at, updated_at, deleted_at, is_active
	) VALUES($1,$2,$3,$4,null,true)`

	qGet = `
	   SELECT 
	   	id, external_id, created_at, updated_at
	   FROM accounts
	   WHERE (id = ANY($1) OR $1 IS NULL)
	   AND ( external_id = ANY($2) OR $2 IS NULL)
	   AND is_active = true
	   LIMIT $3
	   OFFSET $4
	`
)

type accountRepository struct {
	db *sql.DB
}

func NewAccountRepo(db *sql.DB) *accountRepository {
	return &accountRepository{db: db}
}

func (a *accountRepository) Get(ctx context.Context, filter internal.AccountFilter) (model.Account, error) {
	limit := 1
	row := a.db.QueryRowContext(
		ctx,
		qGet,
		pq.Array(filter.IDs),
		pq.Array(filter.ExternalIDs),
		limit,
		defaultOffset,
	)

	acc := model.Account{}
	err := row.Scan(
		&acc.ID,
		&acc.ExternalCustomerID,
		&acc.CreatedAt,
		&acc.UpdatedAt,
	)
	if err != nil {
		return model.Account{}, err
	}

	return acc, nil
}

func (a *accountRepository) CreateTx(ctx context.Context, tx *sql.Tx, newAcc model.Account) (err error) {

	stmt, err := tx.Prepare(qCreate)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(
		ctx,
		newAcc.ID,
		newAcc.ExternalCustomerID,
		newAcc.CreatedAt,
		newAcc.UpdatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}
