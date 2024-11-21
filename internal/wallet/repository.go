package wallet

import (
	"database/sql"

	"github.com/hokdre/mini-ewallet/internal"
	"github.com/hokdre/mini-ewallet/internal/model"
	"github.com/lib/pq"
	"golang.org/x/net/context"
)

const (
	defaultOrderDirection = "DESC"
	defaultOffset         = 0
	defaultOrderColumn    = "created_at"

	qCreate = `
		INSERT INTO wallets (
			id, owned_by, balance, status, enabled_at, disabled_at, created_at, updated_at, deleted_at, is_active
		) VALUES(
			$1, $2, $3, $4, $5, $6, $7, $8, null, true 
		)
	`
	qGet = `
	   SELECT 
	   	id, owned_by, balance, status, enabled_at, disabled_at, created_at, updated_at
	   FROM wallets
	   WHERE (id = ANY($1) or $1 IS NULL)
	   AND (owned_by = ANY($2) or $2 IS NULL)
	   LIMIT $3
	   OFFSET $4
	`

	qUpdate = `
		UPDATE wallets SET
			status = $1, enabled_at = $2, disabled_at = $3, updated_at = $4
		WHERE id = $5
	`

	qIncrementWallet = `
		UPDATE wallets SET
			balance = balance + $1,
			updated_at = $2
		WHERE id = $3
	`

	qDecrementWallet = `
		UPDATE wallets SET
			balance = balance - $1,
			updated_at = $2
		WHERE id = $3 AND balance >= $1
	`
)

type walletRepository struct {
	db *sql.DB
}

func NewWalletRepository(db *sql.DB) *walletRepository {
	return &walletRepository{db: db}
}

func (w *walletRepository) CreateTx(ctx context.Context, tx *sql.Tx, newWallet model.Wallet) error {
	stmt, err := tx.Prepare(qCreate)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(
		ctx,
		newWallet.ID,
		newWallet.OwnedBy,
		newWallet.Balance,
		newWallet.Status,
		newWallet.EnabledAt,
		newWallet.DisabledAt,
		newWallet.CreatedAt,
		newWallet.UpdatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}

func (a *walletRepository) GetOne(ctx context.Context, filter internal.WalletFilter) (model.Wallet, error) {
	limit := 1
	row := a.db.QueryRowContext(
		ctx,
		qGet,
		pq.Array(filter.IDs),
		pq.Array(filter.OwnedBies),
		limit,
		defaultOffset,
	)

	wallet := model.Wallet{}
	err := row.Scan(
		&wallet.ID,
		&wallet.OwnedBy,
		&wallet.Balance,
		&wallet.Status,
		&wallet.EnabledAt,
		&wallet.DisabledAt,
		&wallet.CreatedAt,
		&wallet.UpdatedAt,
	)
	if err != nil {
		return model.Wallet{}, err
	}

	return wallet, nil
}

func (a *walletRepository) Update(ctx context.Context, wallet model.Wallet) error {
	stmt, err := a.db.Prepare(qUpdate)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(
		ctx,
		wallet.Status,
		wallet.EnabledAt,
		wallet.DisabledAt,
		wallet.UpdatedAt,
		wallet.ID,
	)
	if err != nil {
		return err
	}

	return nil
}

func (a *walletRepository) Increment(ctx context.Context, tx *sql.Tx, wallet model.Wallet, amount int64) (int64, error) {
	stmt, err := tx.Prepare(qIncrementWallet)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(
		ctx,
		amount,
		wallet.UpdatedAt,
		wallet.ID,
	)
	if err != nil {
		return 0, err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	return affected, nil
}

func (a *walletRepository) Decrement(ctx context.Context, tx *sql.Tx, wallet model.Wallet, amount int64) (int64, error) {
	stmt, err := tx.Prepare(qDecrementWallet)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(
		ctx,
		amount,
		wallet.UpdatedAt,
		wallet.ID,
	)
	if err != nil {
		return 0, err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	return affected, nil
}
