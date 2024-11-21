package transaction

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

	qCreate = `INSERT INTO transactions(
		id, 
		wallet_id, 
		type, 
		status, 
		reference_id, 
		amount, 
		transacted_at, 
		created_at, 
		updated_at, 
		deleted_at,
		is_active
	) VALUES($1,$2,$3,$4,$5,$6, null, $7,$8,null,true)`

	qList = `
	   SELECT 
	   	id, 
		wallet_id, 
		type, 
		status, 
		reference_id, 
		amount, 
		transacted_at, 
		created_at, 
		updated_at 
	   FROM transactions
	   WHERE (id = ANY($1) or $1 IS NULL)
	   AND ( wallet_id = ANY($2) or $2 IS NULL)
	   AND is_active = true
	   ORDER BY $3 DESC
	`

	qUpdate = `
	UPDATE 
		transactions
	SET 
		status = $1,
		transacted_at = $2,
		updated_at = $3
	WHERE 
		id = $4
	`
)

type transactionRepository struct {
	db *sql.DB
}

func NewAccountRepo(db *sql.DB) *transactionRepository {
	return &transactionRepository{db: db}
}

func (a *transactionRepository) List(ctx context.Context, filter internal.TransactionFilter) ([]model.Transaction, error) {
	rows, err := a.db.QueryContext(
		ctx,
		qList,
		pq.Array(filter.IDs),
		pq.Array(filter.WalletIDs),
		defaultOrderColumn,
	)
	if err != nil {
		return nil, err
	}

	transactions := []model.Transaction{}
	for rows.Next() {
		t := model.Transaction{}
		err := rows.Scan(
			&t.ID,
			&t.WalletID,
			&t.Type,
			&t.Status,
			&t.ReferenceID,
			&t.Amount,
			&t.TransactedAt,
			&t.CreatedAt,
			&t.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		transactions = append(transactions, t)
	}

	return transactions, nil
}

func (a *transactionRepository) Create(ctx context.Context, newAcc model.Transaction) (err error) {

	stmt, err := a.db.Prepare(qCreate)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(
		ctx,
		newAcc.ID,
		newAcc.WalletID,
		newAcc.Type,
		newAcc.Status,
		newAcc.ReferenceID,
		newAcc.Amount,
		newAcc.CreatedAt,
		newAcc.UpdatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}

func (a *transactionRepository) UpdateTx(ctx context.Context, tx *sql.Tx, transaction model.Transaction) (err error) {

	stmt, err := tx.Prepare(qUpdate)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(
		ctx,
		transaction.Status,
		transaction.TransactedAt,
		transaction.UpdatedAt,
		transaction.ID,
	)
	if err != nil {
		return err
	}

	return nil
}
