package transaction

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/hokdre/mini-ewallet/internal"
	"github.com/hokdre/mini-ewallet/internal/model"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestAccountRepository(t *testing.T) {
	t.Run("Create", TestCreate)
	t.Run("UpdateTx", TestUpdateTx)
	t.Run("List", TestList)
}

func TestCreate(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		assert.NoError(t, err)
		defer db.Close()

		timestamp := time.Now()
		newAcc := model.Transaction{
			ID:          uuid.New(),
			WalletID:    uuid.New(),
			Type:        model.TransactionType.Deposit,
			Status:      model.TransactionStatus.Pending,
			ReferenceID: "abc",
			Amount:      10000,
			CreatedAt:   timestamp,
			UpdatedAt:   timestamp,
		}
		mock.
			ExpectPrepare(qCreate).
			ExpectExec().
			WithArgs(
				newAcc.ID,
				newAcc.WalletID,
				newAcc.Type,
				newAcc.Status,
				newAcc.ReferenceID,
				newAcc.Amount,
				newAcc.CreatedAt,
				newAcc.UpdatedAt,
			).
			WillReturnResult(sqlmock.NewResult(0, 1))

		repo := &transactionRepository{db: db}
		errCreate := repo.Create(context.Background(), newAcc)
		assert.NoError(t, errCreate)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Failed Prepare", func(t *testing.T) {
		db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		assert.NoError(t, err)
		defer db.Close()

		newAcc := model.Transaction{}
		errExpected := errors.New("err")
		mock.
			ExpectPrepare(qCreate).
			WillReturnError(errExpected)

		repo := &transactionRepository{db: db}
		errCreate := repo.Create(context.Background(), newAcc)
		assert.Error(t, errCreate, errExpected)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Failed Execute", func(t *testing.T) {
		db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		assert.NoError(t, err)
		defer db.Close()

		timestamp := time.Now()
		newAcc := model.Transaction{
			ID:          uuid.New(),
			WalletID:    uuid.New(),
			Type:        model.TransactionType.Deposit,
			Status:      model.TransactionStatus.Pending,
			ReferenceID: "abc",
			Amount:      10000,
			CreatedAt:   timestamp,
			UpdatedAt:   timestamp,
		}

		errExpect := errors.New("err")
		mock.
			ExpectPrepare(qCreate).
			ExpectExec().
			WithArgs(
				newAcc.ID,
				newAcc.WalletID,
				newAcc.Type,
				newAcc.Status,
				newAcc.ReferenceID,
				newAcc.Amount,
				newAcc.CreatedAt,
				newAcc.UpdatedAt,
			).
			WillReturnError(errExpect)

		repo := &transactionRepository{db: db}
		errCreate := repo.Create(context.Background(), newAcc)
		assert.Error(t, errCreate, errExpect)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestUpdateTx(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		assert.NoError(t, err)
		defer db.Close()
		mock.ExpectBegin()

		timestamp := time.Now()
		newAcc := model.Transaction{
			ID:          uuid.New(),
			WalletID:    uuid.New(),
			Type:        model.TransactionType.Deposit,
			Status:      model.TransactionStatus.Pending,
			ReferenceID: "abc",
			Amount:      10000,
			CreatedAt:   timestamp,
			UpdatedAt:   timestamp,
		}
		mock.
			ExpectPrepare(qUpdate).
			ExpectExec().
			WithArgs(
				newAcc.Status,
				newAcc.TransactedAt,
				newAcc.UpdatedAt,
				newAcc.ID,
			).
			WillReturnResult(sqlmock.NewResult(0, 1))

		tx, err := db.Begin()
		assert.NoError(t, err)

		repo := &transactionRepository{db: db}
		errCreate := repo.UpdateTx(context.Background(), tx, newAcc)
		assert.NoError(t, errCreate)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Failed Prepare", func(t *testing.T) {
		db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		assert.NoError(t, err)
		defer db.Close()
		mock.ExpectBegin()

		newAcc := model.Transaction{}
		errExpected := errors.New("err")
		mock.
			ExpectPrepare(qUpdate).
			WillReturnError(errExpected)

		tx, err := db.Begin()
		assert.NoError(t, err)

		repo := &transactionRepository{db: db}
		errCreate := repo.UpdateTx(context.Background(), tx, newAcc)
		assert.Error(t, errCreate, errExpected)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Failed Execute", func(t *testing.T) {
		db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		assert.NoError(t, err)
		defer db.Close()
		mock.ExpectBegin()

		newAcc := model.Transaction{}
		errExpected := errors.New("err")
		mock.
			ExpectPrepare(qUpdate).
			ExpectExec().
			WithArgs(
				newAcc.Status,
				newAcc.TransactedAt,
				newAcc.UpdatedAt,
				newAcc.ID,
			).
			WillReturnError(err)

		tx, err := db.Begin()
		assert.NoError(t, err)

		repo := &transactionRepository{db: db}
		errCreate := repo.UpdateTx(context.Background(), tx, newAcc)
		assert.Error(t, errCreate, errExpected)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestList(t *testing.T) {
	t.Run("success with filter", func(t *testing.T) {
		db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		assert.NoError(t, err)
		defer db.Close()

		timestamp := time.Now()
		acc := model.Transaction{
			ID:          uuid.New(),
			WalletID:    uuid.New(),
			Type:        model.TransactionType.Deposit,
			Status:      model.TransactionStatus.Pending,
			ReferenceID: "abc",
			Amount:      10000,
			CreatedAt:   timestamp,
			UpdatedAt:   timestamp,
		}

		expectedRow := sqlmock.NewRows([]string{
			"id",
			"wallet_id",
			"type",
			"status",
			"reference_id",
			"amount",
			"transacted_at",
			"created_at",
			"updated_at",
		}).AddRow(
			acc.ID,
			acc.WalletID,
			acc.Type,
			acc.Status,
			acc.ReferenceID,
			acc.Amount,
			acc.TransactedAt,
			acc.CreatedAt,
			acc.UpdatedAt,
		)

		filter := internal.TransactionFilter{
			IDs:       []string{acc.ID.String()},
			WalletIDs: []string{acc.WalletID.String()},
		}
		mock.ExpectQuery(qList).WithArgs(
			pq.Array(filter.IDs),
			pq.Array(filter.WalletIDs),
			defaultOrderColumn,
		).WillReturnRows(expectedRow)

		repo := &transactionRepository{db: db}
		results, err := repo.List(context.Background(), filter)
		assert.NoError(t, err)
		assert.Equal(t, []model.Transaction{acc}, results)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Not Found", func(t *testing.T) {
		db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		assert.NoError(t, err)
		defer db.Close()

		filter := internal.TransactionFilter{
			IDs:       []string{},
			WalletIDs: []string{},
		}
		mock.ExpectQuery(qList).WithArgs(
			pq.Array(filter.IDs),
			pq.Array(filter.WalletIDs),
			defaultOrderColumn,
		).WillReturnError(sql.ErrNoRows)

		repo := &transactionRepository{db: db}
		results, err := repo.List(context.Background(), filter)
		assert.Error(t, err)
		assert.Nil(t, results)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
